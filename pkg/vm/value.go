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
	Count() Value
	Empty() Collection
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
	case reflect.Map, reflect.Chan:
		if v.IsNil() {
			return NIL, nil
		}
		return NIL, NewTypeError(v, "is not boxable", nil)
	default:
		return NIL, NewTypeError(v, "is not boxable", nil)
	}
}

func IsTruthy(v Value) bool {
	return !(v == NIL || v == FALSE)
}
