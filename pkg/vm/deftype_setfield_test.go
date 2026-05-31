/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import "testing"

// DTypeInstance.SetField mutates a field in place by name, backing cljs-style
// ^:mutable deftype fields.
func TestDTypeInstance_SetField(t *testing.T) {
	dt := NewDType("Point", []Symbol{"x", "y"})
	inst := NewDTypeInstance(dt, []Value{Int(1), Int(2)})

	if err := inst.SetField("x", Int(99)); err != nil {
		t.Fatalf("SetField x: %v", err)
	}
	if got := inst.Fields()[0]; got != Int(99) {
		t.Fatalf("x = %v, want 99", got)
	}
	// Other field untouched.
	if got := inst.Fields()[1]; got != Int(2) {
		t.Fatalf("y = %v, want 2 (should be untouched)", got)
	}
	// Read-back through the same field-access path the (.field this) read uses.
	if got, _ := inst.InvokeMethod("x", nil); got != Int(99) {
		t.Fatalf("InvokeMethod x = %v, want 99", got)
	}
}

// A Go nil written through SetField is normalized to the NIL sentinel, so the
// field-value invariant (no raw Go nil) holds.
func TestDTypeInstance_SetField_NilNormalized(t *testing.T) {
	dt := NewDType("Point", []Symbol{"x"})
	inst := NewDTypeInstance(dt, []Value{Int(1)})

	if err := inst.SetField("x", nil); err != nil {
		t.Fatalf("SetField nil: %v", err)
	}
	if got := inst.Fields()[0]; got != NIL {
		t.Fatalf("x = %#v, want NIL (Go nil should be normalized)", got)
	}
}

func TestDTypeInstance_SetField_UnknownField(t *testing.T) {
	dt := NewDType("Point", []Symbol{"x"})
	inst := NewDTypeInstance(dt, []Value{Int(1)})

	if err := inst.SetField("z", Int(5)); err == nil {
		t.Fatal("SetField on unknown field should error, got nil")
	}
}
