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

// PrepareCall resolves fn for repeated arity-n invocation. It returns nil for
// targets that are not plain fixed-arity bytecode callables (variadic, native,
// protocol, ...), and for arities without a matching CallN method; callers
// fall back to ec.Invoke.
func (ec *ExecContext) PrepareCall(fn Fn, arity int) *PreparedCall {
	// Only arities a CallN entry point can fully populate are prepared —
	// Call1 and Call2 today. Widen this as CallN methods land.
	if arity < 1 || arity > 2 {
		return nil
	}
	args := make([]Value, arity)
	target, direct, err := resolveBytecodeCall(fn, args)
	if err != nil || !direct || target.fn.isVariadric {
		return nil
	}
	ec = ec.orRoot()
	f := NewFrame(target.fn.chunk, nil)
	f.closedOvers = target.closedOvers
	f.ec = ec
	return &PreparedCall{
		fn:          target.fn,
		closedOvers: target.closedOvers,
		ec:          ec,
		frame:       f,
		args:        args,
		constsc:     target.fn.chunk.consts.count(),
		stackLen:    max(target.fn.chunk.maxStack, 4),
	}
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
}
