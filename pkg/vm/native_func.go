/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated
 * documentation files (the "Software"), to deal in the Software without restriction, including without limitation the
 * rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit
 * persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the
 * Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE
 * WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
 * COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
 * OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package vm

import (
	"fmt"
	"reflect"
)

type theNativeFnType struct{}

func (t *theNativeFnType) Name() string { return "NativeFn" }
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
			rawArgs[i] = reflect.ValueOf(args[i].Unbox())
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

var NativeFnType *theNativeFnType

func init() {
	NativeFnType = &theNativeFnType{}
}

type NativeFn struct {
	arity       int
	isVariadric bool
	fn          interface{}
	proxy       func([]Value) Value
}

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
	return fmt.Sprintf("<native-fn %p>", l)
}
