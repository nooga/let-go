/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

type theNilType struct {
	zero *Nil
}

func (t *theNilType) String() string     { return t.Name() }
func (t *theNilType) Type() ValueType    { return NilType }
func (t *theNilType) Unbox() interface{} { return nil }

func (t *theNilType) Name() string                     { return "nil" }
func (t *theNilType) Box(_ interface{}) (Value, error) { return t.zero, nil }

// Nil is a Value whose only value is Nil
type Nil struct{}

func (n *Nil) Conj(value Value) Collection {
	return EmptyList.Conj(value)
}

func (n *Nil) Count() Value {
	return Int(0)
}

func (n *Nil) RawCount() int {
	return 0
}

func (n *Nil) Empty() Collection {
	return n
}

func (n *Nil) Cons(value Value) Seq {
	return EmptyList.Cons(value)
}

func (n *Nil) First() Value {
	return n
}

func (n *Nil) More() Seq {
	return n
}

func (n *Nil) Next() Seq {
	return n
}

func (n *Nil) Seq() Seq {
	return n
}

// Type implements Value
func (n *Nil) Type() ValueType { return NilType }

// Unbox implements Value
func (n *Nil) Unbox() interface{} { return nil }

// NilType is the type of NilValues
var NilType *theNilType

// NIL is the only value of NilType (and only instance of NilValue)
var NIL *Nil

func init() {
	NilType = &theNilType{zero: nil}
	NIL = NilType.zero
}

func (n *Nil) String() string {
	return "nil"
}
