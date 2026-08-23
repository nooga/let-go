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

		// A trailing `error` is Go's failure channel, so surface it as a throw
		// rather than as a value — but only when a non-error result remains to
		// return. `func() error` (Close, Write, …) keeps handing the error back
		// as an ordinary value, which is what existing let-go code expects; the
		// peel applies from (T, error) upward, exactly as it always has.
		//
		// The type test is against the `error` interface rather than the
		// implemented-by relation: res carries DECLARED result types, so a
		// function returning a concrete *MyError is returning a value, not
		// signalling failure.
		var callErr error
		if n := len(res); n >= 2 && res[n-1].Type() == reflect.TypeFor[error]() {
			if e := res[n-1].Interface(); e != nil {
				callErr = e.(error)
			}
			res = res[:n-1]
		}

		switch len(res) {
		case 0:
			return NIL, callErr
		case 1:
			wv, err := BoxValue(res[0])
			if err != nil {
				return NIL, err
			}
			return wv, callErr
		}

		// Two or more results: box them into a vector. Previously everything
		// past the second was dropped here, on the Go side and before boxing,
		// so a .lg veneer could never recover it. (a, b, ok) and (a, b, err)
		// are ordinary modern Go — strings.Cut returns three.
		vals := make([]Value, len(res))
		for i := range res {
			wv, err := BoxValue(res[i])
			if err != nil {
				return NIL, err
			}
			vals[i] = wv
		}
		return ArrayVector(vals), callErr
	}

	f := &NativeFn{
		arity:       declArgs,
		isVariadric: variadric,
		fn:          fn,
		proxy:       proxy,
	}

	return f, nil
}
