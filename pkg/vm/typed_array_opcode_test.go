/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"testing"
)

// runArrayOp assembles LOAD_CONST for each operand, the opcode, RETURN, and
// runs the chunk.
func runArrayOp(t *testing.T, op int32, operands ...Value) (Value, error) {
	t.Helper()
	consts := NewConsts()
	chunk := NewCodeChunk(consts)
	chunk.SetMaxStack(len(operands) + 1)
	for _, v := range operands {
		chunk.Append(OP_LOAD_CONST)
		chunk.Append32(consts.Intern(v))
	}
	chunk.Append(op)
	chunk.Append(OP_RETURN)
	return NewFrame(chunk, nil).Run()
}

func TestOpAgetReadsEveryKind(t *testing.T) {
	floats := NewFloatArray(3)
	floats.floats[1] = 2.5
	ints := NewIntArray(3)
	ints.ints[2] = 1 << 40
	bytes := NewByteArrayFrom([]byte{7, 8, 9})
	objs := NewObjectArray(2)
	objs.objs[0] = String("x")

	cases := []struct {
		name string
		arr  *TypedArray
		idx  int
		want Value
	}{
		{"float", floats, 1, Float(2.5)},
		{"int", ints, 2, Int(1 << 40)},
		{"byte", bytes, 0, Int(7)},
		{"object", objs, 0, String("x")},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := runArrayOp(t, OP_AGET, tc.arr, Int(tc.idx))
			if err != nil {
				t.Fatalf("OP_AGET: %v", err)
			}
			if got != tc.want {
				t.Fatalf("OP_AGET = %v (%T), want %v", got, got, tc.want)
			}
		})
	}
}

// aset returns the value as given, not the coerced element, exactly like
// clojure.core/aset.
func TestOpAsetStoresAndReturnsTheGivenValue(t *testing.T) {
	arr := NewFloatArray(2)
	got, err := runArrayOp(t, OP_ASET, arr, Int(1), Int(3))
	if err != nil {
		t.Fatalf("OP_ASET: %v", err)
	}
	if got != Int(3) {
		t.Fatalf("OP_ASET result = %v, want the given Int 3", got)
	}
	if arr.floats[1] != 3.0 {
		t.Fatalf("element = %v, want 3.0", arr.floats[1])
	}
}

// Every rejection reports the same text CoreAgetf / CoreAsetf produce.
func TestOpAgetAsetErrorsMatchTheCoreFns(t *testing.T) {
	arr := NewFloatArray(2)
	cases := []struct {
		name     string
		op       int32
		operands []Value
		want     string
	}{
		{"aget not array", OP_AGET, []Value{String("s"), Int(0)}, "aget expects array, got let-go.lang.String"},
		{"aget index not int", OP_AGET, []Value{arr, Float(0)}, "aget index must be Int"},
		{"aget negative", OP_AGET, []Value{arr, Int(-1)}, "array index -1 out of bounds for length 2"},
		{"aget past end", OP_AGET, []Value{arr, Int(2)}, "array index 2 out of bounds for length 2"},
		{"aset not array", OP_ASET, []Value{NIL, Int(0), Int(1)}, "aset expects array, got nil"},
		{"aset index not int", OP_ASET, []Value{arr, String("0"), Int(1)}, "aset index must be Int"},
		{"aset past end", OP_ASET, []Value{arr, Int(5), Int(1)}, "array index 5 out of bounds for length 2"},
		{"aset wrong element type", OP_ASET, []Value{arr, Int(0), String("no")}, "double-array expects numeric, got let-go.lang.String"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := runArrayOp(t, tc.op, tc.operands...)
			if err == nil {
				t.Fatal("expected an error")
			}
			if err.Error() != tc.want {
				t.Fatalf("error = %q, want %q", err.Error(), tc.want)
			}
		})
	}
}
