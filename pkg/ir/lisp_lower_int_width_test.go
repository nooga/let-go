/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package ir_test

// The lowerer types a let-go :int as Go int64, not Go int. On a host where Go's
// int is 32 bits — TinyGo wasm, linux/386, linux/arm — `int` would narrow a
// value the VM holds at full width, so the lowered core computed different
// answers from the bytecode core on exactly those hosts. It did not even
// compile there: hash-combine's golden-ratio constant 2654435769 exceeds
// math.MaxInt32, so the generated core.go failed with "overflows int".
//
// These tests pin the emitted Go type rather than the arithmetic, because the
// arithmetic is only observable on a 32-bit host and the test suite runs on
// 64-bit. Flip go-type-spec's :int arm back to "int" and both fail here.

import (
	"strings"
	"testing"
)

func TestLoweredIntConstantWiderThanInt32(t *testing.T) {
	ensureLoader()

	// 2654435769 > math.MaxInt32 — the constant that broke the TinyGo build.
	rendered := lowerForms(t, "iwpkg",
		`(defn combine ^long [^long h] (unchecked-add h 2654435769))`)

	if !strings.Contains(rendered, "2654435769") {
		t.Fatalf("expected the wide constant to survive lowering; rendered:\n%s", rendered)
	}
	if !strings.Contains(rendered, "func Combine(ec *vm.ExecContext, arg0 int64) (int64, error)") &&
		!strings.Contains(rendered, "func Combine(ec *vm.ExecContext, arg0 int64) int64") {
		t.Fatalf("expected int64 param and return so the constant fits on a 32-bit host; rendered:\n%s", rendered)
	}
	// A bare `int` declaration is the regression: it truncates wherever Go's
	// int is narrower than the value. Match with the trailing delimiter so
	// "int64" does not count as a hit.
	for _, bad := range []string{" int)", " int,", " int ", " int\n"} {
		if strings.Contains(rendered, bad) {
			t.Errorf("lowered Go declares a host-width int (%q); rendered:\n%s", bad, rendered)
		}
	}
}

func TestLoweredUnsignedShiftIsWidthExplicit(t *testing.T) {
	ensureLoader()

	// unsigned-bit-shift-right went through Go's `uint`, which is host-width:
	// on a 32-bit host it shifted the wrong width and truncated the result.
	rendered := lowerForms(t, "ushpkg",
		`(defn ush ^long [^long x] (unsigned-bit-shift-right x 3))`)

	if strings.Contains(rendered, "uint(") {
		t.Errorf("unsigned shift uses host-width uint; rendered:\n%s", rendered)
	}
	if !strings.Contains(rendered, "uint64(") {
		t.Errorf("expected an explicit uint64 shift; rendered:\n%s", rendered)
	}
}
