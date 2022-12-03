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

func (l Set) toList() *List {
	lst := l.keys()
	ret, _ := ListType.Box(lst)
	return ret.(*List)
}

func (l Set) Seq() Seq {
	return l.toList()
}

// More implements Seq
func (l Set) More() Seq {
	if len(l) == 1 {
		return EmptyList
	}
	ret := l.toList()
	return ret.More()
}

// Next implements Seq
func (l Set) Next() Seq {
	return l.More()
}

// Cons implements Seq
func (l Set) Cons(val Value) Seq {
	return l.toList().Cons(val)
}

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
