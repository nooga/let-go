/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package compiler

import (
	"fmt"
	"io"
	"reflect"

	"github.com/nooga/let-go/pkg/errors"
	"github.com/nooga/let-go/pkg/vm"
)

type ReaderError struct {
	inputName string
	message   string
	pos       int
	line      int
	column    int
	cause     error
}

func NewReaderError(r *LispReader, message string) *ReaderError {
	return &ReaderError{
		inputName: r.inputName,
		pos:       r.pos,
		line:      r.line,
		column:    r.column,
		message:   message,
	}
}

func (r *ReaderError) IsEOF() bool {
	return errorIsEOF(r.cause)
}

func (r *ReaderError) Error() string {
	return errors.AddCause(r,
		fmt.Sprintf(
			"Syntax error reading source at (%s:%d:%d).\n%s",
			r.inputName,
			r.line+1,
			r.column+1,
			r.message,
		))
}

func (r *ReaderError) Wrap(err error) errors.Error {
	r.cause = err
	return r
}

func (r *ReaderError) GetCause() error {
	return r.cause
}

type CompileError struct {
	message string
	source  *vm.SourceInfo
	cause   error
}

func NewCompileError(message string) *CompileError {
	return &CompileError{
		message: message,
	}
}

func NewCompileErrorWithSource(message string, info *vm.SourceInfo) *CompileError {
	return &CompileError{
		message: message,
		source:  info,
	}
}

func (r *CompileError) Error() string {
	return errors.AddCause(r,
		fmt.Sprintf("CompileError: %s", r.message))
}

func (r *CompileError) Source() *vm.SourceInfo {
	return r.source
}

func (r *CompileError) Message() string { return r.message }

// InnermostSource walks the error chain and returns the deepest source info found.
func (r *CompileError) InnermostSource() *vm.SourceInfo {
	if r.source != nil {
		return r.source
	}
	if c, ok := r.cause.(*CompileError); ok {
		return c.InnermostSource()
	}
	return nil
}

// InnermostMessage walks the error chain and returns the deepest error message.
func (r *CompileError) InnermostMessage() string {
	if r.cause == nil {
		return r.message
	}
	if c, ok := r.cause.(*CompileError); ok {
		return c.InnermostMessage()
	}
	return r.cause.Error()
}

func (r *CompileError) Wrap(err error) errors.Error {
	r.cause = err
	return r
}

func (r *CompileError) GetCause() error {
	return r.cause
}

func isErrorEOF(err error) bool {
	return errorIsEOF(err)
}

func errorIsEOF(err error) bool {
	pending := []error{err}
	seen := make(map[error]struct{})
	for len(pending) > 0 {
		err = pending[len(pending)-1]
		pending = pending[:len(pending)-1]
		if err == nil {
			continue
		}
		if err == io.EOF {
			return true
		}
		// Most errors are pointers or comparable values. Avoid revisiting them so
		// a malformed cyclic Unwrap graph cannot loop forever; guard the map use
		// because an error implementation is allowed to have a non-comparable
		// concrete value.
		if reflect.TypeOf(err).Comparable() {
			if _, ok := seen[err]; ok {
				continue
			}
			seen[err] = struct{}{}
		}
		if caused, ok := err.(errors.Error); ok {
			pending = append(pending, caused.GetCause())
		}
		if wrapped, ok := err.(interface{ Unwrap() []error }); ok {
			pending = append(pending, wrapped.Unwrap()...)
		} else if wrapped, ok := err.(interface{ Unwrap() error }); ok {
			pending = append(pending, wrapped.Unwrap())
		}
	}
	return false
}

// IsErrorEOF reports whether err represents end-of-input: a plain or
// standard-wrapped io.EOF, or a compiler error whose nested cause is io.EOF.
// Use this instead of errors.Is(err, io.EOF) for reader output — ReaderError
// keeps its cause in GetCause rather than implementing Unwrap.
func IsErrorEOF(err error) bool { return isErrorEOF(err) }
