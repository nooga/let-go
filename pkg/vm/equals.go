/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

// ValueEquals performs deep structural equality comparison.
// Used by the = operator and OP_EQ opcode.
var ValueEquals func(a, b Value) bool

// SetValueEquals sets the equality function (called by rt package during init).
func SetValueEquals(fn func(a, b Value) bool) {
	ValueEquals = fn
}
