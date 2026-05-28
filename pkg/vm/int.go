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

func (t *theIntType) String() string  { return t.Name() }
func (t *theIntType) Type() ValueType { return TypeType }
func (t *theIntType) Unbox() any      { return reflect.TypeFor[*theIntType]() }

func (t *theIntType) Name() string { return "let-go.lang.Int" }

func (t *theIntType) Box(bare any) (Value, error) {
	switch v := bare.(type) {
	case int:
		return Int(v), nil
	case int8:
		return Int(v), nil
	case int16:
		return Int(v), nil
	case int32:
		return Int(v), nil
	case int64:
		return Int(v), nil
	case uint:
		return Int(v), nil
	case uint8:
		return Int(v), nil
	case uint16:
		return Int(v), nil
	case uint32:
		return Int(v), nil
	case uint64:
		return Int(v), nil
	}
	return IntType.zero, NewTypeError(bare, "can't be boxed as", t)
}

// IntType is the type of IntValues
var IntType *theIntType = &theIntType{zero: 0}

// Int is boxed int
type Int int

// Hash implements Hashable.
func (l Int) Hash() uint32 { return hashUint64(uint64(l)) }

// Type implements Value
func (l Int) Type() ValueType { return IntType }

// Unbox implements Unbox
func (l Int) Unbox() any {
	return int(l)
}

func (l Int) String() string {
	return fmt.Sprintf("%d", int(l))
}
