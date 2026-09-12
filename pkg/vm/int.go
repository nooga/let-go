/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"reflect"
	"strconv"
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

// Int is the boxed integer. It is 64 bits wide on every host, like Clojure's
// long: a value never depends on the platform's int width, so a 32-bit-int
// target (linux/386, TinyGo wasm) computes exactly what a 64-bit one does.
type Int int64

// Hash implements Hashable.
func (l Int) Hash() uint32 { return hashUint64(uint64(l)) }

// Type implements Value
func (l Int) Type() ValueType { return IntType }

// Unbox implements Unbox. It returns the platform int because it feeds Go
// interop (reflect proxies, struct mapping) whose declared types are int; on a
// 32-bit host a wide value narrows here. ToInt64 returns the full value.
func (l Int) Unbox() any {
	return int(l)
}

func (l Int) String() string {
	return strconv.FormatInt(int64(l), 10)
}
