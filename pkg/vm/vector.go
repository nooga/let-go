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
	"strings"
)

type theArrayVectorType struct{}

func (lt *theArrayVectorType) Name() string { return "let-go.lang.ArrayVector" }

func (lt *theArrayVectorType) Box(bare interface{}) (Value, error) {
	arr, ok := bare.([]Value)
	if !ok {
		return NIL, NewTypeError(bare, "can't be boxed as", lt)
	}

	return ArrayVector(arr), nil
}

// ArrayVectorType is the type of ArrayVectors
var ArrayVectorType *theArrayVectorType

func init() {
	ArrayVectorType = &theArrayVectorType{}
}

// ArrayVector is boxed singly linked list that can hold other Values.
type ArrayVector []Value

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
	if len(l) == 1 {
		return ArrayVector{}
	}
	return l[1:]
}

// Next implements Seq
func (l ArrayVector) Next() Seq {
	return l.More()
}

// Cons implements Seq
func (l ArrayVector) Cons(val Value) Seq {
	return append(l, val)
}

// Count implements Collection
func (l ArrayVector) Count() Value {
	return Int(len(l))
}

// Empty implements Collection
func (l ArrayVector) Empty() Collection {
	return make(ArrayVector, 0)
}

func NewArrayVector(v []Value) Value {
	return ArrayVector(v)
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
