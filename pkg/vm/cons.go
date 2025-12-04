/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import "strings"

// Cons is a cons cell - a pair of (first, rest) where rest is a seq.
// This is the fundamental building block for consing onto any seq.
type Cons struct {
	first Value
	more  Seq
}

func NewCons(first Value, more Seq) *Cons {
	return &Cons{first: first, more: more}
}

func (c *Cons) String() string {
	b := &strings.Builder{}
	b.WriteRune('(')
	b.WriteString(c.first.String())
	for s := c.more; s != nil && s != EmptyList; s = s.Next() {
		b.WriteRune(' ')
		b.WriteString(s.First().String())
	}
	b.WriteRune(')')
	return b.String()
}

func (c *Cons) Type() ValueType    { return ListType }
func (c *Cons) Unbox() interface{} { return c }

func (c *Cons) First() Value {
	return c.first
}

func (c *Cons) More() Seq {
	if c.more == nil {
		return EmptyList
	}
	return c.more
}

func (c *Cons) Next() Seq {
	return c.more
}

func (c *Cons) Cons(val Value) Seq {
	return NewCons(val, c)
}

func (c *Cons) Seq() Seq {
	return c
}

func (c *Cons) Count() Value {
	n := 1
	for s := c.more; s != nil && s != EmptyList; s = s.Next() {
		n++
	}
	return Int(n)
}

func (c *Cons) RawCount() int {
	n := 1
	for s := c.more; s != nil && s != EmptyList; s = s.Next() {
		n++
	}
	return n
}

func (c *Cons) Empty() Collection {
	return EmptyList
}

func (c *Cons) Conj(val Value) Collection {
	return NewCons(val, c)
}

