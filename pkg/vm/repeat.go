/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import "strings"

type theRepeatType struct {
	zero *Repeat
}

func (t *theRepeatType) String() string     { return t.Name() }
func (t *theRepeatType) Type() ValueType    { return RepeatType }
func (t *theRepeatType) Unbox() interface{} { return nil }

func (t *theRepeatType) Name() string                     { return "let-go.lang.Repeat" }
func (t *theRepeatType) Box(_ interface{}) (Value, error) { return t.zero, nil }

// Repeat is a Value whose only value is Repeat
type Repeat struct {
	i   int
	val Value
}

func (n *Repeat) Cons(value Value) Seq {
	l := EmptyList
	for i := n.i; i > 0; i-- {
		l = l.Cons(n.val).(*List)
	}
	return l.Cons(value)
}

func (n *Repeat) First() Value {
	return n.val
}

func (n *Repeat) More() Seq {
	return n.Next()
}

func (n *Repeat) Next() Seq {
	i := n.i
	if i == 0 {
		return NIL
	}
	if i > 0 {
		i--
	}
	return &Repeat{
		val: n.val,
		i:   i,
	}
}

func (n *Repeat) Seq() Seq {
	return n
}

func (n *Repeat) Count() Value {
	return Int(n.i)
}

func (n *Repeat) RawCount() int {
	return n.i
}

// Type implements Value
func (n *Repeat) Type() ValueType { return RepeatType }

// Unbox implements Value
func (n *Repeat) Unbox() interface{} { return nil }

// RepeatType is the type of RepeatValues
var RepeatType *theRepeatType = &theRepeatType{zero: nil}

func (n *Repeat) String() string {
	b := &strings.Builder{}
	b.WriteRune('(')
	for i := n.i; i > 0; i-- {
		b.WriteString(n.val.String())
		if i > 1 {
			b.WriteRune(' ')
		}
	}
	b.WriteRune(')')
	return b.String()
}

func NewRepeat(val Value, i int) *Repeat {
	return &Repeat{
		i:   i,
		val: val,
	}
}
