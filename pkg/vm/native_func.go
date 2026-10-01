/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"fmt"
	"os"
	"reflect"
)

type theNativeFnType struct{}

func (t *theNativeFnType) String() string  { return t.Name() }
func (t *theNativeFnType) Type() ValueType { return TypeType }
func (t *theNativeFnType) Unbox() any      { return reflect.TypeFor[*theNativeFnType]() }

func (t *theNativeFnType) Name() string { return "let-go.lang.NativeFn" }

// fastArity reports the arity of a Go func whose shape call dispatches
// directly, without reflection: Values in, a Value out, with or without a
// trailing error. Variadic shapes report -1. ok is false for any other shape.
func fastArity(fn any) (arity int, variadic bool, ok bool) {
	switch fn.(type) {
	case func() (Value, error), func() Value:
		return 0, false, true
	case func(Value) (Value, error), func(Value) Value:
		return 1, false, true
	case func(Value, Value) (Value, error), func(Value, Value) Value:
		return 2, false, true
	case func(Value, Value, Value) (Value, error), func(Value, Value, Value) Value:
		return 3, false, true
	case func(Value, Value, Value, Value) (Value, error), func(Value, Value, Value, Value) Value:
		return 4, false, true
	case func(...Value) (Value, error):
		return -1, true, true
	}
	return 0, false, false
}

// WrongArgCount is the error a fixed-arity native fn returns when called with
// the wrong number of arguments: the ExecutionError the bytecode VM raises for
// a bytecode fn of that arity, naming fn as the VM prints it. Lowered Go
// closures (NativeBody) and typed override wrappers report their own arity
// check through it.
func WrongArgCount(fn Value, got, want int) error {
	return NewExecutionError(fmt.Sprintf("function %s expected %d args, got %d", fn, want, got))
}

// TooFewArgs is the error a variadic native fn returns when called with fewer
// arguments than its fixed parameters, as the bytecode VM reports it.
func TooFewArgs(fn Value, got, min int) error {
	return NewExecutionError(fmt.Sprintf("function %s expected at least %d args, got %d", fn, min, got))
}

// orNIL maps the nil interface a Go func may return in place of NIL to NIL,
// as the reflect path's BoxValue does, so no fast shape hands nil to callers.
func orNIL(v Value) Value {
	if v == nil {
		return NIL
	}
	return v
}

// orNILErr is orNIL for a (Value, error) result.
func orNILErr(v Value, err error) (Value, error) {
	return orNIL(v), err
}

// NativeBody is a callable that lives inside a larger Go struct, so that one
// struct allocation holds the NativeFn, the body and everything the body
// closes over. Lowered Go emits one such struct per fn literal: the captured
// values are its fields and CallNative is the lowered body.
type NativeBody interface {
	CallNative(args []Value) (Value, error)
}

// InitNativeFn sets up n in place, without allocating, to dispatch to body
// with the declared arity. n is normally the NativeFn embedded in the struct
// that implements body, so the caller returns &n as the let-go value and the
// whole closure is one allocation. arity counts the fixed parameters; for a
// variadic body it is the minimum call width, which is what multi-arity
// dispatch (MakeMultiArity) reads off a variadic native.
func InitNativeFn(n *NativeFn, arity int, variadic bool, body NativeBody) {
	n.arity = arity
	n.isVariadric = variadic
	n.fn = body
}

// callFast invokes l's func, which has a fastArity shape or is the NativeBody
// of a lowered closure. A type switch on the stored func replaces a per-box
// adapter closure, so boxing costs one allocation. The trailing error case is
// a guard: Box admits only these shapes.
func (l *NativeFn) callFast(a []Value) (Value, error) {
	switch f := l.fn.(type) {
	case NativeBody:
		return f.CallNative(a)
	case func() (Value, error):
		if len(a) != 0 {
			return NIL, WrongArgCount(l, len(a), 0)
		}
		return orNILErr(f())
	case func(Value) (Value, error):
		if len(a) != 1 {
			return NIL, WrongArgCount(l, len(a), 1)
		}
		return orNILErr(f(a[0]))
	case func(Value, Value) (Value, error):
		if len(a) != 2 {
			return NIL, WrongArgCount(l, len(a), 2)
		}
		return orNILErr(f(a[0], a[1]))
	case func(Value, Value, Value) (Value, error):
		if len(a) != 3 {
			return NIL, WrongArgCount(l, len(a), 3)
		}
		return orNILErr(f(a[0], a[1], a[2]))
	case func(Value, Value, Value, Value) (Value, error):
		if len(a) != 4 {
			return NIL, WrongArgCount(l, len(a), 4)
		}
		return orNILErr(f(a[0], a[1], a[2], a[3]))
	case func() Value:
		if len(a) != 0 {
			return NIL, WrongArgCount(l, len(a), 0)
		}
		return orNIL(f()), nil
	case func(Value) Value:
		if len(a) != 1 {
			return NIL, WrongArgCount(l, len(a), 1)
		}
		return orNIL(f(a[0])), nil
	case func(Value, Value) Value:
		if len(a) != 2 {
			return NIL, WrongArgCount(l, len(a), 2)
		}
		return orNIL(f(a[0], a[1])), nil
	case func(Value, Value, Value) Value:
		if len(a) != 3 {
			return NIL, WrongArgCount(l, len(a), 3)
		}
		return orNIL(f(a[0], a[1], a[2])), nil
	case func(Value, Value, Value, Value) Value:
		if len(a) != 4 {
			return NIL, WrongArgCount(l, len(a), 4)
		}
		return orNIL(f(a[0], a[1], a[2], a[3])), nil
	case func(...Value) (Value, error):
		return orNILErr(f(a...))
	}
	return NIL, fmt.Errorf("callFast: unsupported shape %T", l.fn)
}

func (t *theNativeFnType) Box(fn any) (Value, error) {
	ty := reflect.TypeOf(fn)
	if ty == nil || ty.Kind() != reflect.Func {
		return NIL, NewTypeError(fn, "can't be boxed into", t)
	}

	// Fast path: the fastArity shapes — func(Value…N) returning Value or
	// (Value, error), and func(...Value) (Value, error) — dispatch directly
	// through callFast, no reflection. Everything else uses the reflect proxy
	// below.
	if arity, variadic, ok := fastArity(fn); ok {
		return &NativeFn{arity: arity, isVariadric: variadic, fn: fn}, nil
	}

	// Signature inspection + reflect.Call dispatch live in a build-tagged
	// helper: the stock path reflects the Go signature and calls through it,
	// while TinyGo's reflect can't do IsVariadic/NumIn/In/Value.Call (all
	// unimplemented), so it returns a stub that errors only if actually
	// invoked — letting interop namespaces install at boot without trapping.
	return boxReflectFunc(t, fn, ty)
}

// boxArgForReflect prepares a let-go Value for reflect.Call into a Go fn.
//
// When the Go parameter is a slice/array or map kind, we want per-element
// conversion (so e.g. []vm.Int can flow into []int, and a let-go map into a
// map[string]any). The struct_mapping machinery already does this via
// unboxSliceInto and unboxMapInto, so we delegate to those. For other targets
// and for boxed Go values, plain Unbox is correct.
//
// A failed conversion is reported only when the Unbox fallback cannot satisfy
// the target either. The fallback is legitimate on its own — a *Boxed already
// holding a Go map takes it and succeeds — so failing eagerly would break
// working calls. But when neither path produces a usable value, returning the
// conversion error is what lets the caller say
//
//	cannot convert map key 65 (let-go.lang.Integer) to string
//
// instead of letting reflect.Call panic with the diagnostic that started this
// whole issue:
//
//	reflect: Call using *vm.PersistentMap as type map[string]interface {}
func boxArgForReflect(v Value, target reflect.Type) (reflect.Value, error) {
	if debugBoxArgs {
		fmt.Fprintf(os.Stderr, "[boxArgForReflect] v=%T target=%s kind=%s\n", v, target.String(), target.Kind())
	}
	// Remembered from a slice or map attempt, surfaced only at the end.
	var convErr error

	if target.Kind() == reflect.Slice || target.Kind() == reflect.Array {
		// TypedArray is the mutable array boundary: when its native backing
		// slice already satisfies the Go parameter, pass that slice through so
		// writes made by APIs such as io.Reader.Read remain visible to let-go.
		// Incompatible element types and fixed arrays still use the ordinary
		// per-element conversion below, as do persistent collections.
		if target.Kind() == reflect.Slice {
			if arr, ok := v.(*TypedArray); ok && arr != nil {
				backing := reflect.ValueOf(arr.Unbox())
				if usableAsReflectArg(backing, target) {
					return backing, nil
				}
			}
		}
		if sq, ok := v.(Sequable); ok {
			out := reflect.New(target).Elem()
			err := unboxSliceInto(out, sq.Seq())
			if err == nil {
				return out, nil
			}
			convErr = err
		}
	}
	if target.Kind() == reflect.Map {
		// Without this, a map target fell through to the Unbox fallback below
		// and reflect.Call died, because (*PersistentMap).Unbox returns the
		// map itself.
		out := reflect.New(target).Elem()
		err := unboxMapInto(out, v)
		if err == nil {
			return out, nil
		}
		convErr = err
	}
	// When the Go param is an interface (typically vm.Value itself), pass
	// the boxed Value directly. Unboxing first would surface a Go-native
	// type (int64, string, []any, …) that reflect.Call can't assign to a
	// vm.Value-typed slot. The Generated Go IR-stack code (lowered defns
	// wrapping inner closures via BoxNativeFn) relies on this path.
	if target.Kind() == reflect.Interface {
		rv := reflect.ValueOf(v)
		if debugBoxArgs {
			fmt.Fprintf(os.Stderr, "[boxArgForReflect]   interface path: rv.Type=%s assignable=%v\n", rv.Type().String(), rv.Type().AssignableTo(target))
		}
		if rv.IsValid() && rv.Type().AssignableTo(target) {
			return rv, nil
		}
	}
	if debugBoxArgs {
		fmt.Fprintf(os.Stderr, "[boxArgForReflect]   FALLBACK Unbox: v.Unbox()=%T\n", v.Unbox())
	}
	out := reflect.ValueOf(v.Unbox())
	if convErr != nil && !usableAsReflectArg(out, target) {
		return reflect.Value{}, convErr
	}
	return out, nil
}

// usableAsReflectArg mirrors what boxReflectFunc does with a prepared
// argument: it is passed through when assignable, and converted when
// convertible. Anything else makes reflect.Call panic.
func usableAsReflectArg(v reflect.Value, target reflect.Type) bool {
	return v.IsValid() && (v.Type().AssignableTo(target) || v.CanConvert(target))
}

var debugBoxArgs = os.Getenv("LG_BOXARGS_DEBUG") != ""

func (t *theNativeFnType) WrapNoErr(fn func([]Value) Value) (Value, error) {
	return t.Wrap(func(args []Value) (Value, error) {
		return fn(args), nil
	})
}

func (t *theNativeFnType) Wrap(fn func([]Value) (Value, error)) (Value, error) {
	f := &NativeFn{
		arity:       -1,
		isVariadric: false,
		fn:          fn,
		proxy:       fn,
	}

	return f, nil
}

func (l *NativeFn) WithArity(arity int, variadric bool) *NativeFn {
	l.arity = arity
	l.isVariadric = variadric
	return l
}

var NativeFnType *theNativeFnType = &theNativeFnType{}

type NativeFn struct {
	name        string
	arity       int
	isVariadric bool
	fn          any
	proxy       func([]Value) (Value, error)
	// ctxProxy, when non-nil, is the ExecContext-aware entry point. ec.Invoke
	// routes the live context through it; plain Invoke calls it with the root
	// context. Builtins that read dynamic vars (print → *out*, push-binding!,
	// …) set this.
	ctxProxy func(*ExecContext, []Value) (Value, error)
	meta     Value // IMeta support; nil until (with-meta fn m) copies one in
}

// HasCtx reports whether this native takes an ExecContext.
func (l *NativeFn) HasCtx() bool { return l.ctxProxy != nil }

// invokeCtx runs the context-aware entry point with panic recovery.
func (l *NativeFn) invokeCtx(ec *ExecContext, args []Value) (ret Value, err error) {
	defer RecoverPanic(&err)
	return l.ctxProxy(ec, args)
}

// NewCtxNativeFn builds a context-aware native builtin. Its plain Invoke
// resolves against the root context (host/reflection callers); ec.Invoke
// routes the real context in.
func NewCtxNativeFn(name string, fn func(ec *ExecContext, args []Value) (Value, error)) *NativeFn {
	n := &NativeFn{name: name, arity: -1, isVariadric: true, ctxProxy: fn}
	n.proxy = func(args []Value) (Value, error) { return fn(RootExecContext, args) }
	return n
}

// NewArityNativeFn builds a context-aware native that DECLARES a fixed arity
// (or a variadic minimum) rather than the -1/variadic default NewCtxNativeFn
// uses. The declared arity is what MakeMultiArity dispatches on: an arm whose
// Arity() is -1 and isVariadric is true is indistinguishable from a rest-arm,
// so a wrapper built the ordinary way collapses every arm of a multi-arity fn
// into ma.rest and leaves ma.fns empty. Callers that wrap an opaque callable
// (ir.direct's invokers) use this to keep the arity visible to dispatch.
func NewArityNativeFn(name string, arity int, variadic bool, fn func(ec *ExecContext, args []Value) (Value, error)) *NativeFn {
	n := &NativeFn{name: name, arity: arity, isVariadric: variadic, ctxProxy: fn}
	n.proxy = func(args []Value) (Value, error) { return fn(RootExecContext, args) }
	return n
}

// IsVariadic reports whether this native declares a variadic tail.
func (l *NativeFn) IsVariadic() bool { return l.isVariadric }

func (l *NativeFn) SetName(n string) { l.name = n }

func (l *NativeFn) Type() ValueType { return NativeFnType }

// Unbox implements Unbox
func (l *NativeFn) Unbox() any {
	return l.fn
}

func (l *NativeFn) Arity() int {
	return l.arity
}

func (l *NativeFn) Invoke(args []Value) (ret Value, err error) {
	defer RecoverPanic(&err)
	return l.call(args)
}

// call runs the native without panic recovery: through its proxy when it has
// one (reflected and context-aware natives), else straight to its Go func — a
// fast shape, or the NativeBody a lowered closure embeds this NativeFn in.
func (l *NativeFn) call(args []Value) (Value, error) {
	if l.proxy != nil {
		return l.proxy(args)
	}
	return l.callFast(args)
}

func (l *NativeFn) String() string {
	if len(l.name) > 0 {
		return fmt.Sprintf("<native-fn %s %p>", l.name, l)
	}
	return fmt.Sprintf("<native-fn %p>", l)
}

// Meta implements IMeta.
func (l *NativeFn) Meta() Value {
	if l.meta == nil {
		return NIL
	}
	return l.meta
}

// WithMeta implements IMeta. Returns a copy carrying m. Native builtins are
// shared singletons (e.g. one core/+), so this MUST copy rather than mutate;
// the wrapped fn/proxies are immutable, so sharing them across the copy is safe.
func (l *NativeFn) WithMeta(m Value) Value {
	cp := *l
	cp.meta = m
	return &cp
}
