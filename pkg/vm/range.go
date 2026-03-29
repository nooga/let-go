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
var RangeType *theRangeType = &theRangeType{}

// Range is a lazy integer sequence with start, end, and step.
type Range struct {
	start int
	end   int
	step  int
}

// inBounds reports whether val is within [start, end) respecting step direction.
func (l *Range) inBounds(val int) bool {
	if l.step > 0 {
		return val < l.end
	}
	return val > l.end
}

// Type implements Value
func (l *Range) Type() ValueType { return RangeType }

// Unbox implements Value
func (l *Range) Unbox() interface{} {
	return nil
}

// First implements Seq
func (l *Range) First() Value {
	return Int(l.start)
}

// More implements Seq
func (l *Range) More() Seq {
	nexts := l.start + l.step
	if l.inBounds(nexts) {
		return &Range{nexts, l.end, l.step}
	}
	return EmptyList
}

// Next implements Seq
func (l *Range) Next() Seq {
	nexts := l.start + l.step
	if l.inBounds(nexts) {
		return &Range{nexts, l.end, l.step}
	}
	return nil
}

func (l *Range) Seq() Seq {
	if l.RawCount() == 0 {
		return nil
	}
	return l
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
	diff := l.end - l.start
	if l.step == 1 || l.step == -1 {
		if diff < 0 {
			return -diff
		}
		return diff
	}
	// Integer ceiling division: (diff + step - sign(step)) / step
	if l.step > 0 {
		return (diff + l.step - 1) / l.step
	}
	// step < 0, diff < 0
	return (diff + l.step + 1) / l.step
}

// Empty implements Collection
func (l *Range) Empty() Collection {
	return EmptyList
}

func (l *Range) Conj(val Value) Collection {
	return l.Cons(val).(Collection)
}

func (l *Range) String() string {
	out := "("
	for s := Seq(l); s != nil; s = s.Next() {
		if out != "(" {
			out += " "
		}
		out += s.First().String()
	}
	return out + ")"
}

func (l *Range) ValueAt(key Value) Value {
	return l.ValueAtOr(key, NIL)
}

func (l *Range) ValueAtOr(key Value, dflt Value) Value {
	if key == NIL {
		return dflt
	}
	numkey, ok := key.(Int)
	if !ok {
		return dflt
	}
	idx := int(numkey)
	if idx < 0 || idx >= l.RawCount() {
		return dflt
	}
	return Int(l.start + idx*l.step)
}

// InfiniteRange is a lazy infinite integer sequence.
type InfiniteRange struct {
	start int
	step  int
}

func NewInfiniteRange(start, step int) *InfiniteRange {
	return &InfiniteRange{start: start, step: step}
}

func (r *InfiniteRange) Type() ValueType    { return RangeType }
func (r *InfiniteRange) Unbox() interface{} { return nil }
func (r *InfiniteRange) First() Value       { return Int(r.start) }

func (r *InfiniteRange) Next() Seq {
	return &InfiniteRange{start: r.start + r.step, step: r.step}
}

func (r *InfiniteRange) More() Seq {
	return r.Next()
}

func (r *InfiniteRange) Seq() Seq { return r }

func (r *InfiniteRange) Cons(val Value) Seq {
	return NewCons(val, r)
}

func (r *InfiniteRange) String() string {
	// Don't realize - just show a hint
	return "(range ...)"
}

func NewRange(start, end, step Int) Value {
	s, e, st := int(start), int(end), int(step)
	if st == 0 {
		return EmptyList
	}
	if st > 0 && e > s {
		return &Range{s, e, st}
	}
	if st < 0 && e < s {
		return &Range{s, e, st}
	}
	return EmptyList
}
