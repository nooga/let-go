/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"reflect"
	"strconv"
)

type theFloatType struct {
	zero Float
}

func (t *theFloatType) String() string     { return t.Name() }
func (t *theFloatType) Type() ValueType    { return TypeType }
func (t *theFloatType) Unbox() interface{} { return reflect.TypeOf(t) }

func (t *theFloatType) Name() string { return "let-go.lang.Float" }

func (t *theFloatType) Box(bare interface{}) (Value, error) {
	raw, ok := bare.(float64)
	if !ok {
		return FloatType.zero, NewTypeError(bare, "can't be boxed as", t)
	}
	return Float(raw), nil
}

// FloatType is the type of Float values
var FloatType *theFloatType = &theFloatType{zero: 0}

// Float is boxed float64
type Float float64

// Type implements Value
func (l Float) Type() ValueType { return FloatType }

// Unbox implements Unbox
func (l Float) Unbox() interface{} {
	return float64(l)
}

func (l Float) String() string {
	return strconv.FormatFloat(float64(l), 'g', -1, 64)
}
