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
func (t *theNativeFnType) Box(fn any) (Value, error) {
	ty := reflect.TypeOf(fn)
	if ty.Kind() != reflect.Func {
		return NIL, NewTypeError(fn, "can't be boxed into", t)
	}

	variadric := ty.IsVariadic()
	declArgs := ty.NumIn()
	v := reflect.ValueOf(fn)

	proxy := func(args []Value) (Value, error) {
		rawArgs := make([]reflect.Value, len(args))

		for i := range args {
			// For variadic fns called via reflect.Call (not CallSlice),
			// each variadic arg slot expects the slice's ELEMENT type, not
			// the slice type itself. Previously the loop used the slice
			// type ([]vm.Value) for variadic args, which sent vm.String
			// through the slice-target branch of boxArgForReflect and out
			// the Unbox fallback — surfacing a Go primitive that reflect
			// rejected when dispatching through the let-go Value interface.
			var in reflect.Type
			if variadric && i >= declArgs-1 {
				in = ty.In(declArgs - 1).Elem()
			} else {
				in = ty.In(i)
			}
			if args[i] != NIL {
				rawArgs[i] = boxArgForReflect(args[i], in)
				// Skip the .Convert() step when the prepared value is
				// already assignable to the param's interface type — Convert
				// to an interface erases the dynamic type info reflect.Call
				// needs to dispatch through the let-go Value interface.
				if rawArgs[i].IsValid() && rawArgs[i].Type().AssignableTo(in) {
					// already valid as-is
				} else if rawArgs[i].CanConvert(in) {
					rawArgs[i] = rawArgs[i].Convert(in)
				}
			} else {
				rawArgs[i] = reflect.Zero(in)
			}
		}
		res := v.Call(rawArgs)
		lr := len(res)
		if lr == 0 {
			return NIL, nil
		}
		if lr == 1 {
			wv, err := BoxValue(res[0])
			if err != nil {
				return NIL, err
			}
			return wv, nil
		}
		// assume lr == 2 && res[1] is error
		wv, err := BoxValue(res[0])
		if err != nil {
			return NIL, err
		}
		errorInterface := reflect.TypeFor[error]()
		if res[1].Type() == errorInterface && res[1].Interface() != nil {
			return wv, res[1].Interface().(error)
		}
		return wv, nil
	}

	f := &NativeFn{
		arity:       declArgs,
		isVariadric: variadric,
		fn:          fn,
		proxy:       proxy,
	}

	return f, nil
}

// boxArgForReflect prepares a let-go Value for reflect.Call into a Go fn.
//
// When the Go parameter is a slice/array kind, we want per-element
// conversion (so e.g. []vm.Int can flow into []int). The struct_mapping
// machinery already does this via unboxSliceInto, so we delegate to it.
// For non-slice targets and for boxed Go values, plain Unbox is correct.
func boxArgForReflect(v Value, target reflect.Type) reflect.Value {
	if debugBoxArgs {
		fmt.Fprintf(os.Stderr, "[boxArgForReflect] v=%T target=%s kind=%s\n", v, target.String(), target.Kind())
	}
	if target.Kind() == reflect.Slice || target.Kind() == reflect.Array {
		if sq, ok := v.(Sequable); ok {
			out := reflect.New(target).Elem()
			if err := unboxSliceInto(out, sq.Seq()); err == nil {
				return out
			}
		}
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
			return rv
		}
	}
	if debugBoxArgs {
		fmt.Fprintf(os.Stderr, "[boxArgForReflect]   FALLBACK Unbox: v.Unbox()=%T\n", v.Unbox())
	}
	return reflect.ValueOf(v.Unbox())
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
}

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
	defer recoverThrownPanic(&err)
	return l.proxy(args)
}

func (l *NativeFn) String() string {
	if len(l.name) > 0 {
		return fmt.Sprintf("<native-fn %s %p>", l.name, l)
	}
	return fmt.Sprintf("<native-fn %p>", l)
}
