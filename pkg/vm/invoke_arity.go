/*
 * Copyright (c) 2026 let-go contributors; see CONTRIBUTORS.
 * SPDX-License-Identifier: MIT
 */

package vm

// Arity-typed invocation.
//
// Invoke1..Invoke4 pass a callee one to four arguments as parameters, never
// as a slice, and the callee borrows them for the call and never retains
// them. They reach every callee that can take arguments that way — a lowered
// closure or override wrapper (InvokerN), a fast-shape Go func, a keyword
// lookup, or a bytecode fn through a PreparedCall whose argument slots live
// in the pooled frame. A callee of any other kind, or of another arity, takes
// the slice path, ec.Invoke, whose behavior it then has in every respect.

// Invoker1 through Invoker4 are the arity-typed entry points a lowered
// closure struct or a typed override wrapper implements beside CallNative.
// ec is the caller's context; a lowered closure ignores it in favor of the
// context it captured lexically, as its CallNative does.
type Invoker1 interface {
	Invoke1(ec *ExecContext, a Value) (Value, error)
}

type Invoker2 interface {
	Invoke2(ec *ExecContext, a, b Value) (Value, error)
}

type Invoker3 interface {
	Invoke3(ec *ExecContext, a, b, c Value) (Value, error)
}

type Invoker4 interface {
	Invoke4(ec *ExecContext, a, b, c, d Value) (Value, error)
}

// Invoke1 runs fn with one argument in this context; see Invoke.
func (ec *ExecContext) Invoke1(fn Fn, a Value) (Value, error) {
	ec = ec.orRoot()
	switch f := fn.(type) {
	case *MetaFn:
		return ec.Invoke1(f.Wrapped(), a)
	case *NativeFn:
		return f.invoke1(ec, a)
	case *Func, *Closure:
		var p PreparedCall
		if ec.PrepareCallInto(&p, fn, 1) {
			ret, err := p.Call1(a)
			p.Release()
			return ret, err
		}
	case *MultiArityFn:
		variant, err := f.variantFor(f, 1)
		if err != nil {
			return NIL, err
		}
		return ec.Invoke1(variant, a)
	case Keyword:
		return f.lookup(a), nil
	}
	return ec.Invoke(fn, []Value{a})
}

// Invoke2 runs fn with two arguments in this context; see Invoke.
func (ec *ExecContext) Invoke2(fn Fn, a, b Value) (Value, error) {
	ec = ec.orRoot()
	switch f := fn.(type) {
	case *MetaFn:
		return ec.Invoke2(f.Wrapped(), a, b)
	case *NativeFn:
		return f.invoke2(ec, a, b)
	case *Func, *Closure:
		var p PreparedCall
		if ec.PrepareCallInto(&p, fn, 2) {
			ret, err := p.Call2(a, b)
			p.Release()
			return ret, err
		}
	case *MultiArityFn:
		variant, err := f.variantFor(f, 2)
		if err != nil {
			return NIL, err
		}
		return ec.Invoke2(variant, a, b)
	case Keyword:
		return f.lookupOr(a, b), nil
	}
	return ec.Invoke(fn, []Value{a, b})
}

// Invoke3 runs fn with three arguments in this context; see Invoke.
func (ec *ExecContext) Invoke3(fn Fn, a, b, c Value) (Value, error) {
	ec = ec.orRoot()
	switch f := fn.(type) {
	case *MetaFn:
		return ec.Invoke3(f.Wrapped(), a, b, c)
	case *NativeFn:
		return f.invoke3(ec, a, b, c)
	case *Func, *Closure:
		var p PreparedCall
		if ec.PrepareCallInto(&p, fn, 3) {
			ret, err := p.Call3(a, b, c)
			p.Release()
			return ret, err
		}
	case *MultiArityFn:
		variant, err := f.variantFor(f, 3)
		if err != nil {
			return NIL, err
		}
		return ec.Invoke3(variant, a, b, c)
	}
	return ec.Invoke(fn, []Value{a, b, c})
}

// Invoke4 runs fn with four arguments in this context; see Invoke.
func (ec *ExecContext) Invoke4(fn Fn, a, b, c, d Value) (Value, error) {
	ec = ec.orRoot()
	switch f := fn.(type) {
	case *MetaFn:
		return ec.Invoke4(f.Wrapped(), a, b, c, d)
	case *NativeFn:
		return f.invoke4(ec, a, b, c, d)
	case *Func, *Closure:
		var p PreparedCall
		if ec.PrepareCallInto(&p, fn, 4) {
			ret, err := p.Call4(a, b, c, d)
			p.Release()
			return ret, err
		}
	case *MultiArityFn:
		variant, err := f.variantFor(f, 4)
		if err != nil {
			return NIL, err
		}
		return ec.Invoke4(variant, a, b, c, d)
	}
	return ec.Invoke(fn, []Value{a, b, c, d})
}

// invoke1..invoke4 dispatch a native by arity without an argument slice:
// through its InvokerN (a lowered closure or typed override wrapper), a typed
// context-aware func, or a fast Go func shape, each under the same panic
// recovery as Invoke. A native of any other shape, or of another arity, takes
// the slice path, which produces its usual error.
func (l *NativeFn) invoke1(ec *ExecContext, a Value) (ret Value, err error) {
	defer RecoverPanic(&err)
	switch f := l.fn.(type) {
	case Invoker1:
		return f.Invoke1(ec, a)
	case func(*ExecContext, Value) (Value, error):
		return f(ec, a)
	case func(Value) (Value, error):
		return orNILErr(f(a))
	case func(Value) Value:
		return orNIL(f(a)), nil
	}
	return ec.Invoke(l, []Value{a})
}

func (l *NativeFn) invoke2(ec *ExecContext, a, b Value) (ret Value, err error) {
	defer RecoverPanic(&err)
	switch f := l.fn.(type) {
	case Invoker2:
		return f.Invoke2(ec, a, b)
	case func(*ExecContext, Value, Value) (Value, error):
		return f(ec, a, b)
	case func(Value, Value) (Value, error):
		return orNILErr(f(a, b))
	case func(Value, Value) Value:
		return orNIL(f(a, b)), nil
	}
	return ec.Invoke(l, []Value{a, b})
}

func (l *NativeFn) invoke3(ec *ExecContext, a, b, c Value) (ret Value, err error) {
	defer RecoverPanic(&err)
	switch f := l.fn.(type) {
	case Invoker3:
		return f.Invoke3(ec, a, b, c)
	case func(*ExecContext, Value, Value, Value) (Value, error):
		return f(ec, a, b, c)
	case func(Value, Value, Value) (Value, error):
		return orNILErr(f(a, b, c))
	case func(Value, Value, Value) Value:
		return orNIL(f(a, b, c)), nil
	}
	return ec.Invoke(l, []Value{a, b, c})
}

func (l *NativeFn) invoke4(ec *ExecContext, a, b, c, d Value) (ret Value, err error) {
	defer RecoverPanic(&err)
	switch f := l.fn.(type) {
	case Invoker4:
		return f.Invoke4(ec, a, b, c, d)
	case func(*ExecContext, Value, Value, Value, Value) (Value, error):
		return f(ec, a, b, c, d)
	case func(Value, Value, Value, Value) (Value, error):
		return orNILErr(f(a, b, c, d))
	case func(Value, Value, Value, Value) Value:
		return orNIL(f(a, b, c, d)), nil
	}
	return ec.Invoke(l, []Value{a, b, c, d})
}

// NewTypedCtxNativeFn builds a context-aware native from a typed Go func
// func(ec *ExecContext, a0 … Value) (Value, error) of one to four fixed
// parameters, so that InvokeN reaches it without an argument slice. The slice
// path (ec.Invoke, plain Invoke) goes through an adapter that reports a wrong
// argument count as the bytecode VM does (WrongArgCount), naming the fn. The
// declared arity is the variadic -1 that NewCtxNativeFn declares, so multi-
// arity dispatch over it is unchanged. Lowered packages wrap their override
// adapters with this; any other shape panics, since the shape is fixed at
// code generation.
func NewTypedCtxNativeFn(name string, fn any) *NativeFn {
	n := &NativeFn{name: name, arity: -1, isVariadric: true, fn: fn}
	var adapter func(*ExecContext, []Value) (Value, error)
	switch f := fn.(type) {
	case func(*ExecContext, Value) (Value, error):
		adapter = func(ec *ExecContext, args []Value) (Value, error) {
			if len(args) != 1 {
				return NIL, WrongArgCount(n, len(args), 1)
			}
			return f(ec, args[0])
		}
	case func(*ExecContext, Value, Value) (Value, error):
		adapter = func(ec *ExecContext, args []Value) (Value, error) {
			if len(args) != 2 {
				return NIL, WrongArgCount(n, len(args), 2)
			}
			return f(ec, args[0], args[1])
		}
	case func(*ExecContext, Value, Value, Value) (Value, error):
		adapter = func(ec *ExecContext, args []Value) (Value, error) {
			if len(args) != 3 {
				return NIL, WrongArgCount(n, len(args), 3)
			}
			return f(ec, args[0], args[1], args[2])
		}
	case func(*ExecContext, Value, Value, Value, Value) (Value, error):
		adapter = func(ec *ExecContext, args []Value) (Value, error) {
			if len(args) != 4 {
				return NIL, WrongArgCount(n, len(args), 4)
			}
			return f(ec, args[0], args[1], args[2], args[3])
		}
	default:
		panic("NewTypedCtxNativeFn: unsupported shape " + name)
	}
	n.ctxProxy = adapter
	n.proxy = func(args []Value) (Value, error) { return adapter(RootExecContext, args) }
	return n
}
