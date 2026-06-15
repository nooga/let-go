/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import "testing"

// squareTV is a test record type embedding RecordBase
type squareTV struct {
	RecordBase
	side Value
}

func TestRecordBaseSatisfiesValue(t *testing.T) {
	// Create a test record value
	v := &squareTV{
		RecordBase: RecordBase{TypeName: "Square"},
		side:       MakeInt(3),
	}

	// Verify it satisfies Value interface
	var val Value = v
	_ = val

	// Check String() is not empty
	s := v.String()
	if s == "" {
		t.Fatalf("String() returned empty string")
	}

	// Check Type() returns something
	vt := v.Type()
	if vt == nil {
		t.Fatalf("Type() returned nil")
	}

	// Check Unbox() returns something
	unboxed := v.Unbox()
	if unboxed == nil {
		t.Fatalf("Unbox() returned nil")
	}
}

// TestRecordBaseTypeAndString verifies the TypeName is used correctly
func TestRecordBaseTypeAndString(t *testing.T) {
	v := &squareTV{
		RecordBase: RecordBase{TypeName: "Square"},
		side:       MakeInt(5),
	}

	s := v.String()
	if s != "Square" {
		t.Fatalf("expected String() = \"Square\", got %q", s)
	}

	vt := v.Type()
	if vt == nil {
		t.Fatalf("Type() returned nil")
	}
}
