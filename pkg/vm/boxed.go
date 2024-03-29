/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
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

func (n *Boxed) InvokeMethod(methodName Symbol, args []Value) (Value, error) {
	if n.typ.methods == nil {
		return NIL, fmt.Errorf("%v doesn't have any methods", n.typ)
	}
	method, ok := n.typ.methods[methodName]
	if !ok {
		return NIL, fmt.Errorf("method %s not found in %v", methodName, n.typ)
	}
	return method.Invoke(append([]Value{n}, args...))
}

func (n *Boxed) ValueAt(key Value) Value {
	return n.ValueAtOr(key, NIL)
}

func (n *Boxed) ValueAtOr(key Value, dflt Value) Value {
	name, ok := key.Unbox().(string)
	if !ok {
		return dflt
	}
	v, err := BoxValue(reflect.ValueOf(n.value).FieldByName(name))
	if err != nil {
		return dflt
	}
	return v
}

// BoxedType is the type of NilValues
var BoxedTypes map[string]*aBoxedType = map[string]*aBoxedType{}

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
