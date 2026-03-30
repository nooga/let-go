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
	source  *SourceInfo
	cause   error
}

func NewExecutionError(m string) *ExecutionError {
	return &ExecutionError{message: m}
}

func (ve *ExecutionError) Error() string {
	return errors.AddCause(ve, fmt.Sprintf("ExecutionError: %s", ve.message))
}

func (ve *ExecutionError) WithSource(info *SourceInfo) *ExecutionError {
	ve.source = info
	return ve
}

func (ve *ExecutionError) Wrap(e error) errors.Error {
	ve.cause = e
	return ve
}

func (ve *ExecutionError) GetCause() error {
	return ve.cause
}

// ThrownError wraps a Value that was explicitly thrown.
// It propagates through the normal error return path.
type ThrownError struct {
	Value Value // the thrown value (ExInfo or any Value)
}

func NewThrownError(v Value) *ThrownError {
	return &ThrownError{Value: v}
}

func (e *ThrownError) Error() string {
	if ei, ok := e.Value.(*ExInfo); ok {
		return ei.message
	}
	return e.Value.String()
}

// errorToValue extracts the catchable Value from an error.
// For ThrownError, returns the thrown Value.
// For any other error, wraps the error message as an ExInfo.
func errorToValue(err error) Value {
	// Walk the chain to find a ThrownError
	current := err
	for current != nil {
		if te, ok := current.(*ThrownError); ok {
			return te.Value
		}
		if ee, ok := current.(*ExecutionError); ok {
			current = ee.cause
		} else {
			break
		}
	}
	// Not a ThrownError — extract the root cause message
	msg := innermostMessage(err)
	return NewExInfo(msg, EmptyPersistentMap, nil)
}

// thrownPanic is used to propagate errors through native Go code (map, filter, sort).
// It's caught by recoverThrownPanic in Func/Closure.Invoke.
type thrownPanic struct {
	err error
}

// recoverThrownPanic catches a thrownPanic and converts it back to an error return.
// It also catches arbitrary Go panics and converts them to ExecutionErrors so that
// they produce let-go errors instead of crashing with Go stack traces.
// Call as: defer recoverThrownPanic(&err) at the top of Invoke methods.
func recoverThrownPanic(errp *error) {
	if r := recover(); r != nil {
		if tp, ok := r.(*thrownPanic); ok {
			*errp = tp.err
		} else {
			// Convert arbitrary Go panics to let-go errors
			*errp = fmt.Errorf("%v", r)
		}
	}
}

// innermostMessage extracts the deepest error message from a chain.
func innermostMessage(err error) string {
	for {
		if ee, ok := err.(*ExecutionError); ok && ee.cause != nil {
			err = ee.cause
		} else {
			return err.Error()
		}
	}
}

// unwrapThrown finds a ThrownError anywhere in the error chain.
func unwrapThrown(err error) (*ThrownError, bool) {
	for err != nil {
		if te, ok := err.(*ThrownError); ok {
			return te, true
		}
		// Unwrap through our error types
		switch e := err.(type) {
		case *ExecutionError:
			err = e.cause
		default:
			return nil, false
		}
	}
	return nil, false
}
