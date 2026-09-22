/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package vm

// NewCancelled builds the scope-cancellation condition that a blocking
// native (sleep, channel take/put, ...) raises through its error slot when
// the ExecContext's scope is done, instead of quietly returning nil (issue
// #920). cause is the underlying ctx.Err() — context.Canceled or
// context.DeadlineExceeded — kept as the ExInfo's cause so it survives
// ex-cause/.getCause.
//
// Reuses *ExInfo (the existing exception-class machinery — see
// NewExInfoWithClass) rather than a parallel value type: it already gives
// the condition a message, a cause slot, Error()/String() and the
// getMessage/getCause interop methods every other class-tagged exception
// value has. What makes it a DISTINCT condition, not just another
// exception, is its class: ClassCancelled is never registered in
// vm.ExceptionClasses, so pkg/rt's installExceptionClasses never Defs a
// Lisp-callable constructor for it. The runtime — this function alone — is
// the only source of a value whose Type() is ClassCancelled.
func NewCancelled(cause error) *ExInfo {
	return NewExInfoWithClass(ClassCancelled, "scope cancelled", cause)
}

// IsCancelled reports whether v is exactly the scope-cancellation
// condition — identity via its class, not ancestry. Used at the two places
// that must treat cancellation specially: catch-matches? (pkg/rt/exceptions.go),
// which excludes it from Throwable/Exception's catch-everything role, and
// the go*/future*-style absorption boundaries that must swallow it silently
// on ordinary scope teardown instead of routing it to *err*.
func IsCancelled(v Value) bool {
	if v == nil {
		return false
	}
	return v.Type() == ValueType(ClassCancelled)
}
