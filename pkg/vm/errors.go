/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"fmt"
	"github.com/nooga/let-go/pkg/errors"
	"reflect"
)

// TypeError is a LETGO type error which mostly happens when there is a type mismatch between
// either LETGO values or LETGO values and Go values.
// These errors print as:
//		TypeError: (encountered type name) ... message ... (expected type name)
type TypeError struct {
	message  string
	value    interface{}
	expected ValueType
	cause    error
}

// NewTypeError creates a new type error. This error will print the
// problematic value's (either interface{} or Value) type name, a message, and expected type name.
func NewTypeError(value interface{}, message string, expected ValueType) *TypeError {
	return &TypeError{
		message:  message,
		expected: expected,
		value:    value,
	}
}

// Error implements error
func (te *TypeError) Error() string {
	var s string

	ex := ""
	if te.expected != nil {
		ex = " " + te.expected.Name()
	}

	switch te.value.(type) {
	case Value:
		s = fmt.Sprintf("TypeError: %s %s %s", te.value.(Value).Type().Name(), te.message, ex)
	default:
		s = fmt.Sprintf("TypeError: %s %s %s", reflect.TypeOf(te.value).Name(), te.message, ex)
	}
	return errors.AddCause(te, s)
}

func (te *TypeError) Wrap(e error) errors.Error {
	te.cause = e
	return te
}

func (te *TypeError) GetCause() error {
	return te.cause
}

type ExecutionError struct {
	message string
	cause   error
}

func NewExecutionError(m string) *ExecutionError {
	return &ExecutionError{message: m}
}

func (ve *ExecutionError) Error() string {
	return errors.AddCause(ve, fmt.Sprintf("ExecutionError: %s", ve.message))
}

func (ve *ExecutionError) Wrap(e error) errors.Error {
	ve.cause = e
	return ve
}

func (ve *ExecutionError) GetCause() error {
	return ve.cause
}
