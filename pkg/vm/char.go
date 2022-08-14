/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"reflect"
	"unicode/utf8"
)

type theCharType struct {
	zero Char
}

func (t *theCharType) String() string     { return t.Name() }
func (t *theCharType) Type() ValueType    { return TypeType }
func (t *theCharType) Unbox() interface{} { return reflect.TypeOf(t) }

func (lt *theCharType) Name() string { return "let-go.lang.Character" }

func (lt *theCharType) Box(bare interface{}) (Value, error) {
	raw, ok := bare.(rune)
	if !ok {
		return CharType.zero, NewTypeError(bare, "can't be boxed as", lt)
	}
	return Char(raw), nil
}

// CharType is the type of CharValues
var CharType *theCharType = &theCharType{zero: utf8.RuneError}

// Char is boxed rune
type Char rune

// Type implements Value
func (l Char) Type() ValueType { return CharType }

// Unbox implements Unbox
func (l Char) Unbox() interface{} {
	return rune(l)
}

func (l Char) String() string {
	return "\\" + string(l)
}
