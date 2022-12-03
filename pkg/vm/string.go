/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"fmt"
	"reflect"
)

type theStringType struct {
	zero String
}

func (t *theStringType) String() string     { return t.Name() }
func (t *theStringType) Type() ValueType    { return TypeType }
func (t *theStringType) Unbox() interface{} { return reflect.TypeOf(t) }

func (t *theStringType) Name() string { return "let-go.lang.String" }

func (t *theStringType) Box(bare interface{}) (Value, error) {
	raw, ok := bare.(string)
	if !ok {
		return StringType.zero, NewTypeError(bare, "can't be boxed as", t)
	}
	return String(raw), nil
}

// StringType is the type of StringValues
var StringType *theStringType = &theStringType{zero: ""}

// String is boxed int
type String string

func (l String) Conj(value Value) Collection {
	return String(string(l) + value.String())
}

func (l String) RawCount() int {
	return len(l)
}

func (l String) Count() Value {
	return Int(len(l))
}

func (l String) Empty() Collection {
	return String("")
}

// Type implements Value
func (l String) Type() ValueType { return StringType }

// Unbox implements Unbox
func (l String) Unbox() interface{} {
	return string(l)
}

// First implements Seq
func (l String) First() Value {
	for _, r := range l {
		return Char(r)
	}
	return NIL
}

// More implements Seq
func (l String) More() Seq {
	return l.Next()
}

// Next implements Seq
func (l String) Next() Seq {
	if len(l) <= 1 {
		return NIL
	}
	ret := EmptyList
	s := []rune(l)
	for i := len(s) - 1; i >= 1; i-- {
		ret = ret.Conj(Char(s[i])).(*List)
	}
	return ret
}

// Cons implements Seq
func (l String) Cons(val Value) Seq {
	return NIL
}

func (l String) Seq() Seq {
	if len(l) <= 1 {
		return NIL
	}
	ret := EmptyList
	s := []rune(l)
	for i := len(s) - 1; i >= 0; i-- {
		ret = ret.Conj(Char(s[i])).(*List)
	}
	return ret
}

func (l String) ValueAt(key Value) Value {
	return l.ValueAtOr(key, NIL)
}

func (l String) ValueAtOr(key Value, dflt Value) Value {
	if key == NIL {
		return dflt
	}
	r := []rune(l)
	numkey, ok := key.(Int)
	if !ok || numkey < 0 || int(numkey) >= len(r) {
		return dflt
	}
	return Char(r[numkey])
}

func (l String) String() string {
	return fmt.Sprintf("%q", string(l))
}
