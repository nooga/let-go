/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"fmt"
	"reflect"
)

type theNativeFnType struct{}

func (t *theNativeFnType) String() string     { return t.Name() }
func (t *theNativeFnType) Type() ValueType    { return TypeType }
func (t *theNativeFnType) Unbox() interface{} { return reflect.TypeOf(t) }

func (t *theNativeFnType) Name() string { return "let-go.lang.NativeFn" }
func (t *theNativeFnType) Box(fn interface{}) (Value, error) {
	ty := reflect.TypeOf(fn)
	if ty.Kind() != reflect.Func {
		return NIL, NewTypeError(fn, "can't be boxed into", t)
	}

	variadric := ty.IsVariadic()

	v := reflect.ValueOf(fn)

	proxy := func(args []Value) (Value, error) {
		rawArgs := make([]reflect.Value, len(args))
		for i := range args {
			if args[i] != NIL {
				rawArgs[i] = reflect.ValueOf(args[i].Unbox())
				// FIXME handle variadric
				if rawArgs[i].CanConvert(ty.In(i)) {
					rawArgs[i] = rawArgs[i].Convert(ty.In(i))
				}
			} else {
				//FIXME handle variadric
				rawArgs[i] = reflect.Zero(ty.In(i))
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
		errorInterface := reflect.TypeOf((*error)(nil)).Elem()
		if res[1].Type() == errorInterface && res[1].Interface() != nil {
			return wv, res[1].Interface().(error)
		}
		return wv, nil
	}

	f := &NativeFn{
		arity:       ty.NumIn(),
		isVariadric: variadric,
		fn:          fn,
		proxy:       proxy,
	}

	return f, nil
}

func (t *theNativeFnType) WrapNoErr(fn func([]Value) Value) (Value, error) {
	// FIXME this is ugly and unnecessary wrap in closure
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
	fn          interface{}
	proxy       func([]Value) (Value, error)
}

func (l *NativeFn) SetName(n string) { l.name = n }

func (l *NativeFn) Type() ValueType { return NativeFnType }

// Unbox implements Unbox
func (l *NativeFn) Unbox() interface{} {
	return l.fn
}

func (l *NativeFn) Arity() int {
	return l.arity
}

func (l *NativeFn) Invoke(args []Value) (Value, error) {
	return l.proxy(args)
}

func (l *NativeFn) String() string {
	if len(l.name) > 0 {
		return fmt.Sprintf("<native-fn %s %p>", l.name, l)
	}
	return fmt.Sprintf("<native-fn %p>", l)
}
