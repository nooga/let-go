/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"math"
	"reflect"
	"strconv"
	"unsafe"
)

type theFloatType struct {
	zero Float
}

func (t *theFloatType) String() string  { return t.Name() }
func (t *theFloatType) Type() ValueType { return TypeType }
func (t *theFloatType) Unbox() any      { return reflect.TypeFor[*theFloatType]() }

func (t *theFloatType) Name() string { return "let-go.lang.Float" }

func (t *theFloatType) Box(bare any) (Value, error) {
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
type Float32 float64

// Hash implements Hashable.
func (l Float) Hash() uint32 {
	f := float64(l)
	return hashUint64(*(*uint64)(unsafe.Pointer(&f)))
}

func (l Float32) Hash() uint32 {
	f := float64(l)
	return hashUint64(*(*uint64)(unsafe.Pointer(&f)))
}

// Type implements Value
func (l Float) Type() ValueType   { return FloatType }
func (l Float32) Type() ValueType { return FloatType }

// Unbox implements Unbox
func (l Float) Unbox() any {
	return float64(l)
}

func (l Float32) Unbox() any {
	return float64(l)
}

func (l Float) String() string {
	f := float64(l)
	if !math.IsInf(f, 0) && !math.IsNaN(f) && math.Trunc(f) == f {
		return strconv.FormatFloat(f, 'f', 1, 64)
	}
	return strconv.FormatFloat(f, 'g', -1, 64)
}

func (l Float32) String() string {
	return Float(l).String()
}
