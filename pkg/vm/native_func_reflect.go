//go:build !tinygo

/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import "reflect"

// boxReflectFunc builds a NativeFn that dispatches to a typed Go function via
// reflection: it reads the signature (variadic-ness, declared arity, per-param
// types) and, on invoke, marshals let-go Values into reflect.Values and calls
// through reflect.Value.Call. This is the stock-Go path; TinyGo can't reflect
// function signatures or call through them, so it uses the stub in
// native_func_tinygo.go instead.
func boxReflectFunc(t *theNativeFnType, fn any, ty reflect.Type) (Value, error) {
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
				// NIL to an interface param: pass vm.NIL (falsy) instead of a
				// nil interface (which IsTruthy treats as truthy, breaking
				// (or nil []) patterns) — but ONLY when *vm.Nil actually
				// satisfies the param interface (vm.Value or a super-interface).
				// For unrelated interfaces (error, io.Reader, …) *vm.Nil is not
				// assignable, so a genuine nil interface (reflect.Zero) is
				// required or reflect.Call panics ("using *vm.Nil as type error").
				nilVal := reflect.ValueOf(NIL)
				if in.Kind() == reflect.Interface && nilVal.Type().AssignableTo(in) {
					rawArgs[i] = nilVal
				} else {
					rawArgs[i] = reflect.Zero(in)
				}
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
