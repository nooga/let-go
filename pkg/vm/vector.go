/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"reflect"
	"strings"
)

type theArrayVectorType struct{}

func (t *theArrayVectorType) String() string     { return t.Name() }
func (t *theArrayVectorType) Type() ValueType    { return TypeType }
func (t *theArrayVectorType) Unbox() interface{} { return reflect.TypeOf(t) }

func (lt *theArrayVectorType) Name() string { return "let-go.lang.ArrayVector" }

func (lt *theArrayVectorType) Box(bare interface{}) (Value, error) {
	arr, ok := bare.([]Value)
	if !ok {
		return NIL, NewTypeError(bare, "can't be boxed as", lt)
	}

	return ArrayVector(arr), nil
}

// ArrayVectorType is the type of ArrayVectors
var ArrayVectorType *theArrayVectorType = &theArrayVectorType{}

// ArrayVector is boxed singly linked list that can hold other Values.
type ArrayVector []Value

func (l ArrayVector) Conj(val Value) Collection {
	ret := make([]Value, len(l)+1)
	copy(ret, l)
	ret[len(ret)-1] = val
	return ArrayVector(ret)
}

// Type implements Value
func (l ArrayVector) Type() ValueType { return ArrayVectorType }

// Unbox implements Value
func (l ArrayVector) Unbox() interface{} {
	return []Value(l)
}

// First implements Seq
func (l ArrayVector) First() Value {
	if len(l) == 0 {
		return NIL
	}
	return l[0]
}

// More implements Seq
func (l ArrayVector) More() Seq {
	if len(l) <= 1 {
		return EmptyList
	}
	newl, _ := ListType.Box([]Value(l[1:]))
	return newl.(*List)
}

// Next implements Seq
func (l ArrayVector) Next() Seq {
	return l.More()
}

// Cons implements Seq
func (l ArrayVector) Cons(val Value) Seq {
	ret := EmptyList
	n := len(l) - 1
	for i := range l {
		ret = ret.Cons(l[n-i]).(*List)
	}
	return ret.Cons(val)
}

func (l ArrayVector) Seq() Seq {
	return l
}

// Count implements Collection
func (l ArrayVector) Count() Value {
	return Int(len(l))
}

func (l ArrayVector) RawCount() int {
	return len(l)
}

// Empty implements Collection
func (l ArrayVector) Empty() Collection {
	return make(ArrayVector, 0)
}

func NewArrayVector(v []Value) Value {
	vk := make([]Value, len(v))
	copy(vk, v)
	return ArrayVector(vk)
}

func (l ArrayVector) ValueAt(key Value) Value {
	return l.ValueAtOr(key, NIL)
}

func (l ArrayVector) ValueAtOr(key Value, dflt Value) Value {
	if key == NIL {
		return dflt
	}
	numkey, ok := key.(Int)
	if !ok || numkey < 0 || int(numkey) >= len(l) {
		return dflt
	}
	return l[int(numkey)]
}

func (l ArrayVector) Assoc(k Value, v Value) Associative {
	var new ArrayVector = NewArrayVector(l).(ArrayVector)
	ik, ok := k.(Int)
	if !ok {
		return NIL
	}
	if ik < 0 || int(ik) >= len(new) {
		return NIL
	}
	new[ik] = v
	return new
}

func (l ArrayVector) Dissoc(k Value) Associative {
	return NIL
}

func (l ArrayVector) String() string {
	b := &strings.Builder{}
	b.WriteRune('[')
	n := len(l)
	for i := range l {
		b.WriteString(l[i].String())
		if i < n-1 {
			b.WriteRune(' ')
		}
	}
	b.WriteRune(']')
	return b.String()
}
