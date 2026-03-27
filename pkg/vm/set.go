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

type theSetType struct{}

func (t *theSetType) String() string     { return t.Name() }
func (t *theSetType) Type() ValueType    { return TypeType }
func (t *theSetType) Unbox() interface{} { return reflect.TypeOf(t) }

func (t *theSetType) Name() string { return "let-go.lang.Set" }

func (t *theSetType) Box(bare interface{}) (Value, error) {
	// FIXME make this work
	return NIL, NewTypeError(bare, "can't be boxed as", t)
}

// ArrayVectorType is the type of ArrayVectors
var SetType *theSetType = &theSetType{}

// Set is boxed singly linked list that can hold other Values.
type Set map[Value]struct{}

func (l Set) Conj(value Value) Collection {
	ret := make(Set, len(l)+1)
	for k := range l {
		ret[k] = l[k]
	}
	ret[value] = struct{}{}
	return ret
}

func (l Set) Disj(value Value) Set {
	ret := make(Set, len(l)+1)
	for k := range l {
		if k == value {
			continue
		}
		ret[k] = l[k]
	}
	return ret
}

func (l Set) Contains(value Value) Boolean {
	if _, ok := l[value]; ok {
		return TRUE
	}
	return FALSE
}

// Hash implements Hashable. Unordered hash over elements.
func (l Set) Hash() uint32 {
	var h uint32
	for k := range l {
		h += hashValue(k)
	}
	return mixFinish(h)
}

// Type implements Value
func (l Set) Type() ValueType { return SetType }

func (l Set) keys() []Value {
	ks := make([]Value, len(l))
	i := 0
	for k := range l {
		ks[i] = k
		i++
	}
	return ks
}

// Unbox implements Value
func (l Set) Unbox() interface{} {
	return l.keys()
}

// First implements Seq
func (l Set) First() Value {
	for k := range l {
		return k
	}
	return NIL
}

func (l Set) Seq() Seq {
	if len(l) == 0 {
		return EmptyList
	}
	return &SetSeq{keys: l.keys(), i: 0}
}

// More implements Seq
func (l Set) More() Seq {
	if len(l) <= 1 {
		return EmptyList
	}
	return &SetSeq{keys: l.keys(), i: 1}
}

// Next implements Seq
func (l Set) Next() Seq {
	if len(l) <= 1 {
		return nil
	}
	return &SetSeq{keys: l.keys(), i: 1}
}

// Cons implements Seq
func (l Set) Cons(val Value) Seq {
	return NewCons(val, l.Seq())
}

// SetSeq is a lightweight seq view over a Set's keys
type SetSeq struct {
	keys []Value
	i    int
}

func (s *SetSeq) String() string {
	b := &strings.Builder{}
	b.WriteRune('(')
	for i := s.i; i < len(s.keys); i++ {
		if i > s.i {
			b.WriteRune(' ')
		}
		b.WriteString(s.keys[i].String())
	}
	b.WriteRune(')')
	return b.String()
}

func (s *SetSeq) Type() ValueType    { return ListType }
func (s *SetSeq) Unbox() interface{} { return s.keys[s.i:] }

func (s *SetSeq) First() Value {
	if s.i >= len(s.keys) {
		return NIL
	}
	return s.keys[s.i]
}

func (s *SetSeq) More() Seq {
	if s.i+1 >= len(s.keys) {
		return EmptyList
	}
	return &SetSeq{keys: s.keys, i: s.i + 1}
}

func (s *SetSeq) Next() Seq {
	if s.i+1 >= len(s.keys) {
		return nil
	}
	return &SetSeq{keys: s.keys, i: s.i + 1}
}

func (s *SetSeq) Cons(val Value) Seq {
	return NewCons(val, s)
}

func (s *SetSeq) Count() Value      { return Int(len(s.keys) - s.i) }
func (s *SetSeq) RawCount() int     { return len(s.keys) - s.i }
func (s *SetSeq) Empty() Collection { return EmptyList }
func (s *SetSeq) Conj(val Value) Collection { return s.Cons(val).(*List) }
func (s *SetSeq) Seq() Seq { return s }

// Count implements Collection
func (l Set) Count() Value {
	return Int(len(l))
}

func (l Set) RawCount() int {
	return len(l)
}

// Empty implements Collection
func (l Set) Empty() Collection {
	return make(Set)
}

func NewSet(v []Value) Value {
	if len(v) == 0 {
		return make(Set)
	}
	newmap := make(Set)
	for _, k := range v {
		newmap[k] = struct{}{}
	}
	return newmap
}

func (l Set) String() string {
	b := &strings.Builder{}
	b.WriteString("#{")
	i := 0
	n := len(l)
	for k := range l {
		b.WriteString(k.String())
		if i < n-1 {
			b.WriteRune(' ')
		}
		i++
	}
	b.WriteRune('}')
	return b.String()
}

func (l Set) Arity() int {
	return 1
}

func (l Set) Invoke(pargs []Value) (Value, error) {
	if len(pargs) != 1 {
		return NIL, fmt.Errorf("wrong number of arguments %d", len(pargs))
	}
	if _, ok := l[pargs[0]]; ok {
		return pargs[0], nil
	}
	return NIL, nil
}
