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

func (t *theCharType) String() string  { return t.Name() }
func (t *theCharType) Type() ValueType { return TypeType }
func (t *theCharType) Unbox() any      { return reflect.TypeFor[*theCharType]() }

func (t *theCharType) Name() string { return "let-go.lang.Character" }

func (t *theCharType) Box(bare any) (Value, error) {
	raw, ok := bare.(rune)
	if !ok {
		return CharType.zero, NewTypeError(bare, "can't be boxed as", t)
	}
	return Char(raw), nil
}

// CharType is the type of CharValues
var CharType *theCharType = &theCharType{zero: utf8.RuneError}

// Char is boxed rune
type Char rune

// Hash implements Hashable.
func (l Char) Hash() uint32 { return hashUint64(uint64(l)) }

// Type implements Value
func (l Char) Type() ValueType { return CharType }

// Unbox implements Unbox
func (l Char) Unbox() any {
	return rune(l)
}

func (l Char) String() string {
	switch rune(l) {
	case ' ':
		return "\\space"
	case '\n':
		return "\\newline"
	case '\t':
		return "\\tab"
	case '\r':
		return "\\return"
	}
	return "\\" + string(l)
}
