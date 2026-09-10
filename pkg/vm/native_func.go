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

// fastValueProxy returns a non-reflection proxy for the signatures the
// IR-lowered code boxes most often: func(Value…N) (Value, error), all args and
// the result already let-go Values. A direct typed call skips the per-invocation
// []reflect.Value allocation, per-arg reflect boxing, and reflect.Call (whose
// reflect.unsafe_New dominated the native optimize-pass alloc profile). Returns
// (nil, 0) for any other signature, so Box falls through to the reflect proxy.
func fastValueProxy(fn any) (func([]Value) (Value, error), int) {
	switch f := fn.(type) {
	case func() (Value, error):
		return func(a []Value) (Value, error) {
			if len(a) != 0 {
				return NIL, fmt.Errorf("wrong number of args (%d), expected 0", len(a))
			}
			return f()
		}, 0
	case func(Value) (Value, error):
		return func(a []Value) (Value, error) {
			if len(a) != 1 {
				return NIL, fmt.Errorf("wrong number of args (%d), expected 1", len(a))
			}
			return f(a[0])
		}, 1
	case func(Value, Value) (Value, error):
		return func(a []Value) (Value, error) {
			if len(a) != 2 {
				return NIL, fmt.Errorf("wrong number of args (%d), expected 2", len(a))
			}
			return f(a[0], a[1])
		}, 2
	case func(Value, Value, Value) (Value, error):
		return func(a []Value) (Value, error) {
			if len(a) != 3 {
				return NIL, fmt.Errorf("wrong number of args (%d), expected 3", len(a))
			}
			return f(a[0], a[1], a[2])
		}, 3
	case func(Value, Value, Value, Value) (Value, error):
		return func(a []Value) (Value, error) {
			if len(a) != 4 {
				return NIL, fmt.Errorf("wrong number of args (%d), expected 4", len(a))
			}
			return f(a[0], a[1], a[2], a[3])
		}, 4
	}
	return nil, 0
}

func (t *theNativeFnType) Box(fn any) (Value, error) {
	ty := reflect.TypeOf(fn)
	if ty == nil || ty.Kind() != reflect.Func {
		return NIL, NewTypeError(fn, "can't be boxed into", t)
	}

	// Fast path: exact func(Value…N) (Value, error) shapes dispatch directly,
	// no reflection. Everything else uses the reflect proxy below.
	if fp, arity := fastValueProxy(fn); fp != nil {
		return &NativeFn{arity: arity, isVariadric: false, fn: fn, proxy: fp}, nil
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
	return l.proxy(args)
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
