/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"fmt"
	"reflect"
)

type theChanType struct {
}

func (t *theChanType) String() string     { return t.Name() }
func (t *theChanType) Type() ValueType    { return TypeType }
func (t *theChanType) Unbox() interface{} { return reflect.TypeOf(t) }

func (t *theChanType) Name() string { return "let-go.lang.Chan" }
func (t *theChanType) Box(b interface{}) (Value, error) {
	rb, ok := b.(chan Value)
	if !ok {
		return nil, NewTypeError(b, "can't be boxed as", t)
	}
	return Chan(rb), nil
}

// Chan is either TRUE or FALSE
type Chan chan Value

// Type implements Value
func (n Chan) Type() ValueType { return ChanType }

// Unbox implements Value
func (n Chan) Unbox() interface{} { return n }

// ChanType is the type of Chan
var ChanType *theChanType = &theChanType{}

func (n Chan) String() string {
	return fmt.Sprintf("<chan %p>", n)
}
