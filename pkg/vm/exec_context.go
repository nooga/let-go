/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package vm

import "context"

// ExecContext is the single object that resolves an execution's dynamic state —
// both the dynamic-var binding stack and the structured-concurrency Scope (see
// docs/design/exec-context-threading.md). It is the *implementer* of invocation
// and of dynamic resolution: the eval loop, builtins, and goroutine spawn all
// go through an ExecContext rather than looking state up by goroutine id.
//
// There is always an ExecContext in play. RootExecContext is the process
// default — its binding stack is the shared global one and its scope is the
// root Goroutines scope, so code that never asks for isolation behaves exactly
// as it always has. Per-goroutine isolation is explicit: a spawn hands the
// child a fresh ExecContext seeded from a snapshot of the parent's
// (ec.Child()) — carrying the parent's bindings and scope — so goroutines
// cannot interleave each other's dynamic state, with no goroutine-id lookup and
// no reuse hazard.
type ExecContext struct {
	bindings *bindingStack
	// scope is the structured-concurrency scope this execution runs under
	// (cancellation tree + tracked goroutines). nil means the root Goroutines
	// scope; Scope() normalises that.
	scope *Scope
}

// globalBindingStack backs RootExecContext (and the package-level Var binding
// API used by host/lowered code). Sharing it is what makes the root context
// behave like the historical process-global binding store.
var globalBindingStack = newBindingStack()

// RootExecContext is the always-available default context. Frames with no more
// specific context resolve against it.
var RootExecContext = &ExecContext{bindings: globalBindingStack}

// NewExecContext returns a context with a fresh, empty binding stack.
func NewExecContext() *ExecContext {
	return &ExecContext{bindings: newBindingStack()}
}

// Child returns a fresh ExecContext seeded from a snapshot of this one's
// bindings — the explicit-propagation primitive used at goroutine boundaries.
// A nil receiver seeds from the root context.
func (ec *ExecContext) Child() *ExecContext {
	src := ec
	if src == nil {
		src = RootExecContext
	}
	c := NewExecContext()
	c.bindings.installSnapshot(src.bindings.snapshot())
	c.scope = src.scope
	return c
}

// Scope returns the structured-concurrency scope this context runs under,
// normalising the nil (root) context to the process-wide Goroutines scope.
// A nil receiver also resolves to the root.
func (ec *ExecContext) Scope() *Scope {
	if ec == nil || ec.scope == nil {
		return Goroutines
	}
	return ec.scope
}

// Context returns the cancellation context blocking native ops (channel
// take/put, alts!) select on — the scope's context.
func (ec *ExecContext) Context() context.Context {
	return ec.Scope().Context()
}

// SetScope installs s as this context's scope. Used by with-scope to open a
// child scope for a dynamic extent and to restore the previous one on close.
func (ec *ExecContext) SetScope(s *Scope) {
	ec.scope = s
}

// OpenChildEC opens a child of ec's current scope and installs it on ec for a
// synchronous dynamic extent (the with-scope body, which runs in the same
// frame, hence the same ec). It returns the child; pair with CloseScoped, which
// runs the restore that reinstates the previous scope. This is the
// explicit-context replacement for the old goroutine-id-keyed OpenChild.
func OpenChildEC(ec *ExecContext) *Scope {
	ec = ec.orRoot()
	prev := ec.scope
	c := ec.Scope().Child()
	c.closeRestore = func() { ec.scope = prev }
	ec.SetScope(c)
	return c
}

// BindingSnapshot returns a snapshot of this context's current bindings,
// suitable for passing to NewExecContextFrom to isolate a child context.
// A nil receiver uses the root context.
func (ec *ExecContext) BindingSnapshot() BindingSnapshot {
	src := ec
	if src == nil {
		src = RootExecContext
	}
	return src.bindings.snapshot()
}

// NewExecContextFrom returns a fresh context whose binding stack is seeded from
// snap. It is the per-call isolation primitive: bound-fn* captures a snapshot
// once and builds a fresh context from it on every invocation, so the wrapped
// function always re-establishes exactly the captured bindings and never leaks
// pushes between calls or across goroutines.
func NewExecContextFrom(snap BindingSnapshot) *ExecContext {
	c := NewExecContext()
	c.bindings.installSnapshot(snap)
	return c
}

// orRoot normalises a possibly-nil context to the root.
func (ec *ExecContext) orRoot() *ExecContext {
	if ec == nil {
		return RootExecContext
	}
	return ec
}

// --- dynamic-var resolution (ec owns the binding state) ---------------------

// deref returns v's current value in this context: the top dynamic binding if
// any, else v's root.
func (ec *ExecContext) deref(v *Var) Value {
	ec = ec.orRoot()
	if v.isDynamic.Load() {
		if val, ok := ec.bindings.current(v); ok {
			return val
		}
	}
	return v.Root()
}

func (ec *ExecContext) pushBinding(v *Var, val Value) {
	v.isDynamic.Store(true)
	ec.orRoot().bindings.push(v, val)
}

func (ec *ExecContext) popBinding(v *Var) {
	ec.orRoot().bindings.pop(v)
}

func (ec *ExecContext) hasBinding(v *Var) bool {
	return ec.orRoot().bindings.hasBinding(v)
}

// setBinding mutates v's top dynamic binding in this context (Clojure's
// thread-local set!), returning false if v has no active binding here. Callers
// fall back to mutating the root for the no-binding case.
func (ec *ExecContext) setBinding(v *Var, val Value) bool {
	return ec.orRoot().bindings.setCurrent(v, val)
}

// Exported entry points for runtime builtins (pkg/rt) that resolve dynamic
// vars against an ExecContext handed to them by ec.Invoke.
func (ec *ExecContext) PushBinding(v *Var, val Value) { ec.pushBinding(v, val) }
func (ec *ExecContext) PopBinding(v *Var)             { ec.popBinding(v) }
func (ec *ExecContext) Deref(v *Var) Value            { return ec.deref(v) }

// SetBinding mutates v's top dynamic binding in this context (thread-local
// set!), returning false if v has no active binding here so the caller can fall
// back to the root. Exported for runtime builtins like in-ns.
func (ec *ExecContext) SetBinding(v *Var, val Value) bool { return ec.setBinding(v, val) }

// --- invocation (ec is the implementer) -------------------------------------

// Invoke runs fn with args in this context. Closures propagate the context to
// their child frame; context-aware natives receive it; pure functions are
// invoked unchanged. This is the single dispatch site the eval loop uses.
func (ec *ExecContext) Invoke(fn Fn, args []Value) (Value, error) {
	ec = ec.orRoot()
	switch f := fn.(type) {
	case *Func:
		return f.invokeIn(ec, args)
	case *Closure:
		return f.invokeIn(ec, args)
	case *MultiArityFn:
		return f.invokeIn(ec, args)
	case *NativeFn:
		if f.ctxProxy != nil {
			return f.invokeCtx(ec, args)
		}
		return f.Invoke(args)
	case *ProtocolFn:
		return f.invokeIn(ec, args)
	case *MultiFn:
		return f.invokeIn(ec, args)
	default:
		return fn.Invoke(args)
	}
}

// Bind pins fn to this context: the returned Fn routes every later call
// through ec.Invoke, so a callback handed to context-free machinery (Atom
// retry loops, sort comparators, stored transducer steps) still resolves
// dynamic vars against the caller's context. At the root this is the identity
// — plain Invoke already resolves there — so root-context callers pay nothing.
func (ec *ExecContext) Bind(fn Fn) Fn {
	ec = ec.orRoot()
	if ec == RootExecContext {
		return fn
	}
	bound := &NativeFn{name: "ec-bound-fn", arity: -1, isVariadric: true}
	bound.proxy = func(args []Value) (Value, error) { return ec.Invoke(fn, args) }
	bound.ctxProxy = func(_ *ExecContext, args []Value) (Value, error) { return ec.Invoke(fn, args) }
	return bound
}
