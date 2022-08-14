/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"fmt"
	"reflect"
)

type theIntType struct {
	zero Int
}

func (t *theIntType) String() string     { return t.Name() }
func (t *theIntType) Type() ValueType    { return TypeType }
func (t *theIntType) Unbox() interface{} { return reflect.TypeOf(t) }

func (t *theIntType) Name() string { return "let-go.lang.Int" }

func (t *theIntType) Box(bare interface{}) (Value, error) {
	raw, ok := bare.(int)
	if !ok {
		return IntType.zero, NewTypeError(bare, "can't be boxed as", t)
	}
	return Int(raw), nil
}

// IntType is the type of IntValues
var IntType *theIntType = &theIntType{zero: 0}

// Int is boxed int
type Int int

// Type implements Value
func (l Int) Type() ValueType { return IntType }

// Unbox implements Unbox
func (l Int) Unbox() interface{} {
	return int(l)
}

func (l Int) String() string {
	return fmt.Sprintf("%d", int(l))
}
