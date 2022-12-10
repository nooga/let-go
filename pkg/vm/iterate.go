/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import "fmt"

type theIterateType struct {
	zero *Iterate
}

func (t *theIterateType) String() string     { return t.Name() }
func (t *theIterateType) Type() ValueType    { return IterateType }
func (t *theIterateType) Unbox() interface{} { return nil }

func (t *theIterateType) Name() string                     { return "let-go.lang.Iterate" }
func (t *theIterateType) Box(_ interface{}) (Value, error) { return t.zero, nil }

// Iterate is a Value whose only value is Iterate
type Iterate struct {
	f     Fn
	state Value
}

func (n *Iterate) Cons(value Value) Seq {
	return EmptyList.Cons(value)
}

func (n *Iterate) First() Value {
	return n.state
}

func (n *Iterate) More() Seq {
	return n.Next()
}

func (n *Iterate) Next() Seq {
	nv, err := n.f.Invoke([]Value{n.state})
	if err != nil {
		return NIL
	}
	return &Iterate{
		state: nv,
		f:     n.f,
	}
}

func (n *Iterate) Seq() Seq {
	return n
}

// Type implements Value
func (n *Iterate) Type() ValueType { return IterateType }

// Unbox implements Value
func (n *Iterate) Unbox() interface{} { return nil }

// IterateType is the type of IterateValues
var IterateType *theIterateType = &theIterateType{zero: nil}

func (n *Iterate) String() string {
	return fmt.Sprintf("<let-go.lang.Iterate %s %s>", n.f, n.state)
}

func NewIterate(f Fn, state Value) *Iterate {
	return &Iterate{
		f:     f,
		state: state,
	}
}
