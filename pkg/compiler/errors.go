/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package compiler

import (
	"fmt"
	"github.com/nooga/let-go/pkg/errors"
	"io"
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
	if r.cause != nil {
		c, ok := r.cause.(*ReaderError)
		if ok {
			return c.IsEOF()
		}
	}
	return r.cause == io.EOF
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
	cause   error
}

func NewCompileError(message string) *CompileError {
	return &CompileError{
		message: message,
	}
}

func (r *CompileError) Error() string {
	return errors.AddCause(r,
		fmt.Sprintf("CompileError: %s", r.message))
}

func (r *CompileError) Wrap(err error) errors.Error {
	r.cause = err
	return r
}

func (r *CompileError) GetCause() error {
	return r.cause
}

func isErrorEOF(err error) bool {
	if err == io.EOF {
		return true
	}
	rerr, ok := err.(*ReaderError)
	if ok {
		return rerr.IsEOF()
	}
	return false
}
