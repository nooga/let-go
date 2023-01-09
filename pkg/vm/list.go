/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"reflect"
	"strings"
)

type theListType struct{}

func (t *theListType) String() string     { return t.Name() }
func (t *theListType) Type() ValueType    { return TypeType }
func (t *theListType) Unbox() interface{} { return reflect.TypeOf(t) }

func (lt *theListType) Name() string { return "let-go.lang.PersistentList" }

func (lt *theListType) Box(bare interface{}) (Value, error) {
	arr, ok := bare.([]Value)
	if !ok {
		return EmptyList, NewTypeError(bare, "can't be boxed as", lt)
	}
	var ret Seq = EmptyList
	n := len(arr)
	if n == 0 {
		return ret.(*List), nil
	}
	for i := range arr {
		ret = ret.Cons(arr[n-i-1])
	}
	return ret.(*List), nil
}

// ListType is the type of Lists
var ListType *theListType = &theListType{}

// EmptyList is an empty List
var EmptyList *List = &List{count: 0}

// List is boxed singly linked list that can hold other Values.
type List struct {
	first Value
	next  *List
	count int
}

func (l *List) Conj(value Value) Collection {
	return l.Cons(value).(*List)
}

// Type implements Value
func (l *List) Type() ValueType { return ListType }

// Unbox implements Value
func (l *List) Unbox() interface{} {
	if l.count == 0 {
		return []Value{}
	}
	bare := make([]Value, l.count)
	l.unboxInternal(&bare, 0)
	return bare
}

func (l *List) unboxInternal(target *[]Value, idx int) {
	if l.count == 0 {
		return
	}
	(*target)[idx] = l.first
	l.next.unboxInternal(target, idx+1)
}

// First implements Seq
func (l *List) First() Value {
	if l.count == 0 {
		return NIL
	}
	return l.first
}

// More implements Seq
func (l *List) More() Seq {
	if l.count == 0 {
		return l
	}
	return l.next
}

// Next implements Seq
func (l *List) Next() Seq {
	if l.count == 0 {
		return l
	}
	return l.next
}

// Cons implements Seq
func (l *List) Cons(val Value) Seq {
	return &List{
		first: val,
		next:  l,
		count: l.count + 1,
	}
}

func (l *List) Seq() Seq {
	return l
}

// Count implements Collection
func (l *List) Count() Value {
	return Int(l.count)
}

func (l *List) RawCount() int {
	return l.count
}

// Empty implements Collection
func (l *List) Empty() Collection {
	return EmptyList
}

func (l *List) ValueAt(key Value) Value {
	return l.ValueAtOr(key, NIL)
}

func (l *List) ValueAtOr(key Value, dflt Value) Value {
	if key == NIL {
		return dflt
	}
	numkey, ok := key.(Int)
	if !ok || numkey < 0 {
		return dflt
	}
	li := l
	for i := 0; i < int(numkey); i++ {
		li = li.next
		if li == nil {
			return dflt
		}
	}
	if li.first == nil {
		return NIL
	}
	return li.first
}

func (l *List) String() string {
	b := &strings.Builder{}
	b.WriteRune('(')
	li := l.Unbox().([]Value)
	n := len(li)
	for i := range li {
		b.WriteString(li[i].String())
		if i < n-1 {
			b.WriteRune(' ')
		}
	}
	b.WriteRune(')')
	return b.String()
}

func NewList(vs []Value) Value {
	li, err := ListType.Box(vs)
	if err != nil {
		return EmptyList
	}
	return li
}
