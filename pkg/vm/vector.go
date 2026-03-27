/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"fmt"
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

const arrayVectorPromotionThreshold = 32

func (l ArrayVector) Conj(val Value) Collection {
	newLen := len(l) + 1
	// Promote to PersistentVector when exceeding threshold
	if newLen > arrayVectorPromotionThreshold {
		values := make([]Value, newLen)
		copy(values, l)
		values[newLen-1] = val
		return NewPersistentVector(values).(Collection)
	}
	ret := make([]Value, newLen)
	copy(ret, l)
	ret[newLen-1] = val
	return ArrayVector(ret)
}

// Type implements Value
func (l ArrayVector) Type() ValueType { return ArrayVectorType }

// Unbox implements Value
func (l ArrayVector) Unbox() interface{} {
	return []Value(l)
}

// Equals implements value equality for ArrayVector
func (l ArrayVector) Equals(other Value) bool {
	switch o := other.(type) {
	case ArrayVector:
		if len(l) != len(o) {
			return false
		}
		for i, v := range l {
			if eq, ok := v.(interface{ Equals(Value) bool }); ok {
				if !eq.Equals(o[i]) {
					return false
				}
			} else if v != o[i] {
				return false
			}
		}
		return true
	case PersistentVector:
		if len(l) != o.count {
			return false
		}
		for i, v := range l {
			ov := o.ValueAt(Int(i))
			if eq, ok := v.(interface{ Equals(Value) bool }); ok {
				if !eq.Equals(ov) {
					return false
				}
			} else if v != ov {
				return false
			}
		}
		return true
	default:
		return false
	}
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
	return &ArrayVectorSeq{vec: l, i: 1}
}

// Next implements Seq
func (l ArrayVector) Next() Seq {
	if len(l) <= 1 {
		return nil
	}
	return &ArrayVectorSeq{vec: l, i: 1}
}

// Cons implements Seq
func (l ArrayVector) Cons(val Value) Seq {
	return NewCons(val, l.Seq())
}

func (l ArrayVector) Seq() Seq {
	if len(l) == 0 {
		return EmptyList
	}
	return &ArrayVectorSeq{vec: l, i: 0}
}

// ArrayVectorSeq is a lightweight seq view over an ArrayVector
type ArrayVectorSeq struct {
	vec ArrayVector
	i   int
}

func (s *ArrayVectorSeq) String() string {
	return "(" + s.vec[s.i:].String()[1:] // reuse vector's string but change brackets
}

func (s *ArrayVectorSeq) Type() ValueType {
	return ListType // seqs print as lists
}

func (s *ArrayVectorSeq) Unbox() interface{} {
	return []Value(s.vec[s.i:])
}

func (s *ArrayVectorSeq) First() Value {
	if s.i >= len(s.vec) {
		return NIL
	}
	return s.vec[s.i]
}

func (s *ArrayVectorSeq) More() Seq {
	if s.i+1 >= len(s.vec) {
		return EmptyList
	}
	return &ArrayVectorSeq{vec: s.vec, i: s.i + 1}
}

func (s *ArrayVectorSeq) Next() Seq {
	if s.i+1 >= len(s.vec) {
		return nil
	}
	return &ArrayVectorSeq{vec: s.vec, i: s.i + 1}
}

func (s *ArrayVectorSeq) Cons(val Value) Seq {
	return NewCons(val, s)
}

func (s *ArrayVectorSeq) Count() Value {
	return Int(len(s.vec) - s.i)
}

func (s *ArrayVectorSeq) RawCount() int {
	return len(s.vec) - s.i
}

func (s *ArrayVectorSeq) Empty() Collection {
	return EmptyList
}

func (s *ArrayVectorSeq) Conj(val Value) Collection {
	// Conj on a seq creates a new seq with val at front
	return s.Cons(val).(*List)
}

func (s *ArrayVectorSeq) Seq() Seq {
	return s
}

// ValueAt implements Lookup for ArrayVectorSeq so that `get` works on seq views.
func (s *ArrayVectorSeq) ValueAt(key Value) Value {
	return s.ValueAtOr(key, NIL)
}

// ValueAtOr implements Lookup for ArrayVectorSeq.
func (s *ArrayVectorSeq) ValueAtOr(key Value, dflt Value) Value {
	if key == NIL {
		return dflt
	}
	idx, ok := key.(Int)
	if !ok || idx < 0 {
		return dflt
	}
	absIdx := s.i + int(idx)
	if absIdx >= len(s.vec) {
		return dflt
	}
	return s.vec[absIdx]
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

func (l ArrayVector) Contains(value Value) Boolean {
	numkey, ok := value.(Int)
	if !ok || numkey < 0 || int(numkey) >= len(l) {
		return FALSE
	}
	return TRUE
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

func (l ArrayVector) Arity() int {
	return 1
}

func (l ArrayVector) Invoke(pargs []Value) (Value, error) {
	vl := len(pargs)
	if vl != 1 {
		return NIL, fmt.Errorf("wrong number of arguments %d", vl)
	}
	return l.ValueAt(pargs[0]), nil
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
