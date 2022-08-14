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

func (l String) String() string {
	return fmt.Sprintf("%q", string(l))
}
