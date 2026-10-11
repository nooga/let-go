/*
 * Copyright (c) 2026 Matt Parrett
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"strings"
	"testing"
)

// The typed accessors ir.lower-go emits for aget/aset on a statically typed
// array: in-bounds access reads and writes the shared backing store, out of
// bounds and a kind mismatch come back as errors (never a panic), and the
// errors match what CoreAgetf/CoreAsetf report on the bytecode path.
func TestTypedArrayLoweredAccessors(t *testing.T) {
	d := NewFloatArray(3)
	if err := d.SetFloat64(1, 2.5); err != nil {
		t.Fatal(err)
	}
	if got, err := d.AtFloat64(1); err != nil || got != 2.5 {
		t.Fatalf("AtFloat64 = %v, %v", got, err)
	}
	if got := d.Get(1); got != Float(2.5) {
		t.Fatalf("boxed Get after SetFloat64 = %v, want the same backing store", got)
	}
	if _, err := d.AtFloat64(3); err == nil || !strings.Contains(err.Error(), "out of bounds for length 3") {
		t.Fatalf("AtFloat64 past the end: %v", err)
	}
	if err := d.SetFloat64(-1, 0); err == nil || !strings.Contains(err.Error(), "out of bounds") {
		t.Fatalf("SetFloat64 at -1: %v", err)
	}
	if _, err := d.AtInt64(0); err == nil || !strings.Contains(err.Error(), "int-array expected") {
		t.Fatalf("AtInt64 on a double-array: %v", err)
	}

	n := NewIntArray(2)
	if err := n.SetInt64(0, -7); err != nil {
		t.Fatal(err)
	}
	if got, err := n.AtInt64(0); err != nil || got != -7 {
		t.Fatalf("AtInt64 = %v, %v", got, err)
	}

	b := NewByteArray(1)
	if err := b.SetByte(0, 300); err != nil {
		t.Fatal(err)
	}
	if got, err := b.AtByte(0); err != nil || got != 44 {
		t.Fatalf("AtByte after SetByte(300) = %v, %v; want 44 (truncated as Set does)", got, err)
	}

	o := NewObjectArray(1)
	if err := o.SetValue(0, String("x")); err != nil {
		t.Fatal(err)
	}
	if got, err := o.AtValue(0); err != nil || got != String("x") {
		t.Fatalf("AtValue = %v, %v", got, err)
	}
	if _, err := o.AtFloat64(0); err == nil {
		t.Fatal("AtFloat64 on an object-array must error, not panic")
	}
}
