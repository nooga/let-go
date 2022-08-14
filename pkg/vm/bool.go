/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import "reflect"

type theBooleanType struct {
	zero Boolean
}

func (t *theBooleanType) String() string     { return t.Name() }
func (t *theBooleanType) Type() ValueType    { return TypeType }
func (t *theBooleanType) Unbox() interface{} { return reflect.TypeOf(t) }

func (t *theBooleanType) Name() string { return "let-go.lang.Boolean" }
func (t *theBooleanType) Box(b interface{}) (Value, error) {
	rb, ok := b.(bool)
	if !ok {
		return BooleanType.zero, NewTypeError(b, "can't be boxed as", t)
	}
	return Boolean(rb), nil
}

// Boolean is either TRUE or FALSE
type Boolean bool

// Type implements Value
func (n Boolean) Type() ValueType { return BooleanType }

// Unbox implements Value
func (n Boolean) Unbox() interface{} { return bool(n) }

// BooleanType is the type of Boolean
var BooleanType *theBooleanType = &theBooleanType{zero: FALSE}

// TRUE is Boolean
const TRUE Boolean = true

// FALSE is Boolean
const FALSE Boolean = false

func (n Boolean) String() string {
	if n == TRUE {
		return "true"
	}
	return "false"
}
