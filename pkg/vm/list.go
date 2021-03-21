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
var ListType *theListType

// EmptyList is an empty List
var EmptyList *List

func init() {
	ListType = &theListType{}
	EmptyList = &List{count: 0}
}

// List is boxed singly linked list that can hold other Values.
type List struct {
	first Value
	next  *List
	count int
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
