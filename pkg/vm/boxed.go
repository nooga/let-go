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

type aBoxedType struct {
	typ     reflect.Type
	methods map[Symbol]*NativeFn
}

func (t *aBoxedType) String() string     { return t.Name() }
func (t *aBoxedType) Type() ValueType    { return TypeType }
func (t *aBoxedType) Unbox() interface{} { return t.typ }

func (t *aBoxedType) Name() string { return "go." + t.typ.String() }
func (t *aBoxedType) Box(value interface{}) (Value, error) {
	if !reflect.TypeOf(value).ConvertibleTo(t.typ) {
		return NIL, NewTypeError(value, "can't be boxed as", t)
	}
	return &Boxed{value, t}, nil
}

type Boxed struct {
	value interface{}
	typ   *aBoxedType
}

// Type implements Value
func (n *Boxed) Type() ValueType { return n.typ }

// Unbox implements Value
func (n *Boxed) Unbox() interface{} { return n.value }

func (n *Boxed) String() string {
	return fmt.Sprintf("<%s %v>", n.typ.Name(), n.value)
}

func (n *Boxed) InvokeMethod(methodName Symbol, args []Value) Value {
	if n.typ.methods == nil {
		// FIXME error :P
		fmt.Println("methods nil")
		return NIL
	}
	method, ok := n.typ.methods[methodName]
	if !ok {
		// FIXME error :P
		fmt.Println("method", methodName, "not found")
		return NIL
	}
	return method.Invoke(append([]Value{n}, args...))
}

// BoxedType is the type of NilValues
var BoxedTypes map[string]*aBoxedType

func init() {
	BoxedTypes = map[string]*aBoxedType{}
}

func valueType(value interface{}) *aBoxedType {
	reflected := reflect.TypeOf(value)
	t, ok := BoxedTypes[reflected.Name()]
	if ok {
		return t
	}
	t = &aBoxedType{
		typ:     reflected,
		methods: nil,
	}
	methodc := reflected.NumMethod()
	if methodc > 0 {
		t.methods = map[Symbol]*NativeFn{}
		for i := 0; i < methodc; i++ {
			m := reflected.Method(i)
			me, err := NativeFnType.Box(m.Func.Interface())
			if err != nil {
				// FIXME notice this somehow
				continue
			}
			mef, ok := me.(*NativeFn)
			if !ok {
				// FIXME notice this somehow
				continue
			}
			t.methods[Symbol(m.Name)] = mef
		}
	}
	BoxedTypes[reflected.Name()] = t
	return t
}

func NewBoxed(value interface{}) *Boxed {
	return &Boxed{value: value, typ: valueType(value)}
}
