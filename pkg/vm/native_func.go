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

	proxy := func(args []Value) Value {
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
		if len(res) == 0 {
			return NIL
		}
		wv, err := BoxValue(res[0])
		if err != nil {
			return NIL
		}
		return wv
	}

	f := &NativeFn{
		arity:       ty.NumIn(),
		isVariadric: variadric,
		fn:          fn,
		proxy:       proxy,
	}

	return f, nil
}

func (t *theNativeFnType) Wrap(fn func(args []Value) Value) (Value, error) {
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

var NativeFnType *theNativeFnType

func init() {
	NativeFnType = &theNativeFnType{}
}

type NativeFn struct {
	name        string
	arity       int
	isVariadric bool
	fn          interface{}
	proxy       func([]Value) Value
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

func (l *NativeFn) Invoke(args []Value) Value {
	return l.proxy(args)
}

func (l *NativeFn) String() string {
	if len(l.name) > 0 {
		return fmt.Sprintf("<native-fn %s %p>", l.name, l)
	}
	return fmt.Sprintf("<native-fn %p>", l)
}
