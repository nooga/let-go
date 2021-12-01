/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"reflect"
)

type theRangeType struct{}

func (t *theRangeType) String() string     { return t.Name() }
func (t *theRangeType) Type() ValueType    { return TypeType }
func (t *theRangeType) Unbox() interface{} { return reflect.TypeOf(t) }

func (t *theRangeType) Name() string { return "let-go.lang.Range" }

func (t *theRangeType) Box(bare interface{}) (Value, error) {
	return NIL, NewTypeError(bare, "can't be boxed as", t)
}

// RangeType is the type of Lists
var RangeType *theRangeType

func init() {
	RangeType = &theRangeType{}
}

// Range is boxed singly linked list that can hold other Values.
type Range struct {
	start int
	end   int
	step  int
}

// Type implements Value
func (l *Range) Type() ValueType { return RangeType }

// Unbox implements Value
func (l *Range) Unbox() interface{} {
	return l.Seq().Unbox()
}

// First implements Seq
func (l *Range) First() Value {
	return Int(l.start)
}

// More implements Seq
func (l *Range) More() Seq {
	nexts := l.start + l.step
	if nexts < l.end {
		return &Range{nexts, l.end, l.step}
	}
	return EmptyList
}

// Next implements Seq
func (l *Range) Next() Seq {
	nexts := l.start + l.step
	if nexts < l.end {
		return &Range{nexts, l.end, l.step}
	}
	return EmptyList
}

func (l *Range) Seq() Seq {
	var r Seq = EmptyList
	n := l.end - l.start - 1
	top := l.start + (n/l.step)*l.step
	for i := top; i >= l.start; i -= l.step {
		r = r.Cons(Int(i))
	}
	return r
}

// Cons implements Seq
func (l *Range) Cons(val Value) Seq {
	return l.Seq().Cons(val)
}

// Count implements Collection
func (l *Range) Count() Value {
	return Int(l.RawCount())
}

func (l *Range) RawCount() int {
	if l.step == 1 {
		return l.end - l.start
	}
	return (l.end - l.start + 1) / l.step
}

// Empty implements Collection
func (l *Range) Empty() Collection {
	return EmptyList
}

func (l *Range) Conj(val Value) Collection {
	return l.Cons(val).(Collection)
}

func (l *Range) String() string {
	return l.Seq().String()
}

func (l *Range) ValueAt(key Value) Value {
	return l.ValueAtOr(key, NIL)
}

func (l *Range) ValueAtOr(key Value, dflt Value) Value {
	// FIXME: assumes positive step
	if key == NIL {
		return dflt
	}
	numkey, ok := key.(Int)
	if !ok {
		return dflt
	}
	nth := l.start + int(numkey)*l.step
	if nth <= l.end && nth >= l.start {
		return Int(nth)
	}
	return dflt
}

func NewRange(start, end, step Int) Value {
	// FIXME: Add support for negative step
	if end > start && step > 0 {
		return &Range{
			int(start), int(end), int(step),
		}
	}
	return EmptyList
}
