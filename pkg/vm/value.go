/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"fmt"
	"reflect"
)

// ValueType represents a type of a Value
type ValueType interface {
	Value
	Name() string
	Box(interface{}) (Value, error)
}

// Value is implemented by all LETGO values
type Value interface {
	fmt.Stringer
	Type() ValueType
	Unbox() interface{}
}

// Seq is implemented by all sequence-like values
type Seq interface {
	Value
	Cons(Value) Seq
	First() Value
	More() Seq
	Next() Seq
}

// Collection is implemented by all collections
type Collection interface {
	Value
	RawCount() int
	Count() Value
	Empty() Collection
	Conj(Value) Collection
}

type Fn interface {
	Value
	Invoke([]Value) Value
	Arity() int
}

type Associative interface {
	Value
	Assoc(Value, Value) Associative
	Dissoc(Value) Associative
}

type Lookup interface {
	Value
	ValueAt(Value) Value
	ValueAtOr(Value, Value) Value
}

type Receiver interface {
	Value
	InvokeMethod(Symbol, []Value) Value
}

type Reference interface {
	Deref() Value
}

type theTypeType struct{}

var TypeType *theTypeType

func init() {
	TypeType = &theTypeType{}
}

func (t *theTypeType) String() string     { return t.Name() }
func (t *theTypeType) Type() ValueType    { return t }
func (t *theTypeType) Unbox() interface{} { return reflect.TypeOf(t) }

func (t *theTypeType) Name() string { return "let-go.lang.Type" }
func (t *theTypeType) Box(b interface{}) (Value, error) {
	//FIXME this is probably not accurate
	return NIL, NewTypeError(b, "can't be boxed as", t)
}

func BoxValue(v reflect.Value) (Value, error) {
	if v.CanInterface() {
		rv, ok := v.Interface().(Value)
		if ok {
			return rv, nil
		}
	}
	switch v.Type().Kind() {
	case reflect.Int:
		return IntType.Box(v.Interface())
	case reflect.String:
		return StringType.Box(v.Interface())
	case reflect.Bool:
		return BooleanType.Box(v.Interface())
	case reflect.Func:
		return NativeFnType.Box(v.Interface())
	case reflect.Ptr:
		if v.IsNil() {
			return NIL, nil
		}
		// FIXME check if this is how we should handle pointers
		return BoxValue(v.Elem())
	case reflect.Slice, reflect.Array:
		if v.IsNil() {
			// FIXME not sure if maybe this has to be empty coll in let-go-land
			return NIL, nil
		}
		in := make([]Value, v.Len())
		for i := 0; i < v.Len(); i++ {
			e := v.Index(i)
			mv, err := BoxValue(e)
			if err != nil {
				return NIL, NewTypeError(e, "can't be boxed", nil).Wrap(err)
			}
			in[i] = mv
		}
		return ArrayVector(in), nil
	case reflect.Map:
		if v.IsNil() {
			// FIXME not sure if maybe this has to be empty coll in let-go-land
			return NIL, nil
		}
		in := make(map[Value]Value)
		iter := v.MapRange()
		for iter.Next() {
			k, err := BoxValue(iter.Key())
			if err != nil {
				return NIL, err //FIXME wrap
			}
			v, err := BoxValue(iter.Value())
			if err != nil {
				return NIL, err //FIXME wrap
			}
			in[k] = v
		}
		return MapType.Box(in)
	case reflect.Chan:
		if v.IsNil() {
			return NIL, nil
		}
		return NIL, NewTypeError(v, "is not boxable", nil)
	default:
		if v.CanInterface() {
			return NewBoxed(v.Interface()), nil
		}
		return NIL, NewTypeError(v, "is not boxable", nil)
	}
}

func IsTruthy(v Value) bool {
	return !(v == NIL || v == FALSE)
}
