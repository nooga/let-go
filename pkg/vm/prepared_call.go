/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import "strconv"

// PreparedCall caches the resolution of a fixed-arity bytecode callable and
// reuses one frame across many calls from native code. It exists for
// per-element callback loops (reduce/some/filter and friends): resolution,
// frame-pool traffic, and frame initialization are per-call-site constants,
// so paying them once per seq operation instead of once per element removes
// most of the host->VM entry cost.
//
// A PreparedCall is single-owner: it must be used by one native call
// activation and never shared or called reentrantly. Nested uses of the same
// lg function each prepare their own.
type PreparedCall struct {
	fn          *Func
	closedOvers []Value
	ec          *ExecContext
	frame       *Frame
	args        []Value
	constsc     int
	stackLen    int
}

// maxPreparedArity is the widest arity a CallN entry point can fully
// populate — Call1 and Call2 today. Widen this as CallN methods land.
const maxPreparedArity = 2

// PrepareCallInto resolves fn for repeated arity-n invocation into p, which
// the caller owns — typically a stack variable in a native loop. It reports
// false, leaving p unusable, for targets that are not plain fixed-arity
// bytecode callables (variadic, native, protocol, ...) and for arities
// without a matching CallN method; callers fall back to ec.Invoke.
//
// Nothing here allocates once the frame pool is warm: p is the caller's,
// the frame comes from the pool, and the argument slots live in the frame's
// prepArgs, which survives pool reuse. They deliberately do not alias
// argbuf, which installBytecodeCall overwrites on a same-frame tail call.
func (ec *ExecContext) PrepareCallInto(p *PreparedCall, fn Fn, arity int) bool {
	if arity < 1 || arity > maxPreparedArity {
		return false
	}
	// Inspect the target without resolveBytecodeCall: that path retains the
	// argument slice it is given (so a stack probe would escape) and packs a
	// rest list for variadic fns before we could reject them. unwrapBytecodeFn
	// answers "which *Func, is it fixed-arity n" without allocating.
	target, closedOvers, _, ok, err := unwrapBytecodeFn(fn, arity)
	if err != nil || !ok || target.isVariadric || target.arity != arity {
		return false
	}
	ec = ec.orRoot()
	f := NewFrame(target.chunk, nil)
	f.closedOvers = closedOvers
	f.ec = ec
	// The argument slots live in the pooled frame and survive pool reuse, so
	// they cost nothing after the first use of each frame. They deliberately
	// do not alias argbuf, which installBytecodeCall overwrites on a
	// same-frame tail call.
	if cap(f.prepArgs) < arity {
		f.prepArgs = make([]Value, maxPreparedArity)
	}
	args := f.prepArgs[:arity]
	*p = PreparedCall{
		fn:          target,
		closedOvers: closedOvers,
		ec:          ec,
		frame:       f,
		args:        args,
		constsc:     target.chunk.consts.count(),
		stackLen:    max(target.chunk.maxStack, 4),
	}
	return true
}

// PrepareCall is the heap-allocating form of PrepareCallInto: it returns a
// new PreparedCall, or nil when fn cannot be prepared. Hot loops should
// prefer PrepareCallInto with a stack variable.
func (ec *ExecContext) PrepareCall(fn Fn, arity int) *PreparedCall {
	p := new(PreparedCall)
	if !ec.PrepareCallInto(p, fn, arity) {
		return nil
	}
	return p
}

// Call1 invokes the prepared unary callable. Calling it on a preparation of
// any other arity is a contract violation and returns an error.
func (p *PreparedCall) Call1(a Value) (Value, error) {
	if len(p.args) != 1 {
		return NIL, NewExecutionError("PreparedCall: Call1 on a preparation of arity " + strconv.Itoa(len(p.args)))
	}
	p.args[0] = a
	return p.call()
}

// Call2 invokes the prepared binary callable. Calling it on a preparation of
// any other arity is a contract violation and returns an error.
func (p *PreparedCall) Call2(a, b Value) (Value, error) {
	if len(p.args) != 2 {
		return NIL, NewExecutionError("PreparedCall: Call2 on a preparation of arity " + strconv.Itoa(len(p.args)))
	}
	p.args[0] = a
	p.args[1] = b
	return p.call()
}

// call resets the owned frame and runs it. The reset must cover code/consts/
// closedOvers, not just args/ip/sp: a tail call in the callee body rebinds the
// frame to the tail target via installBytecodeCall, and an error unwind can
// leave stale handlers.
func (p *PreparedCall) call() (Value, error) {
	f := p.frame
	f.args = p.args
	f.argc = len(p.args)
	f.closedOvers = p.closedOvers
	f.code = p.fn.chunk
	f.consts = f.code.consts
	f.constsc = p.constsc
	f.ip = 0
	f.sp = 0
	f.debug = false
	f.parent = nil
	f.stack = f.stack[:p.stackLen]
	if len(f.handlers) > 0 {
		f.handlers = f.handlers[:0]
	}
	state := frameRunState{root: f, current: f}
	if allocAttrEnabled {
		return runChainProtected(&state)
	}
	return runChain(&state)
}

// Release returns the owned frame to the pool. The PreparedCall must not be
// used afterwards.
func (p *PreparedCall) Release() {
	if p.frame != nil {
		ReleaseFrame(p.frame)
		p.frame = nil
	}
	p.args = nil // the slots belong to the frame, which the pool may hand elsewhere
}
