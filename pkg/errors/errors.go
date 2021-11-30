/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package errors

import "fmt"

type Error interface {
	error
	Wrap(error) Error
	GetCause() error
}

func IsCausedBy(me error, e error) bool {
	if me == nil {
		return false
	}
	if me == e {
		return true
	}
	mec, ok := me.(Error)
	if !ok {
		return false
	}
	return IsCausedBy(mec.GetCause(), e)
}

func AddCause(e Error, s string) string {
	cause := e.GetCause()
	if cause == nil {
		return s
	}
	return fmt.Sprintf("%s\n\tcaused by %s", s, cause.Error())
}
