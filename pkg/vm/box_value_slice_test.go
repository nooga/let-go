/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package vm_test

import (
	"reflect"
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

// BoxValue's slice branch switches on the slice's STATIC element type. For a
// []any that kind is Interface, so every element misses the string/int/float
// paths and lands as a vm.Boxed — which prints as <go.string Ada> and compares
// equal to nothing in let-go:
//
//	(= "Ada" (first (sql/ScanRow rows)))   ; => false
//
// []any is the natural Go return type for "a row of unknown column types", so
// this makes the documented database/sql wrapper shape unusable. Elements must
// be boxed by their DYNAMIC type instead.
func TestBoxValueInterfaceSliceElements(t *testing.T) {
	row := []any{"Ada", int64(36), 1.5, true, nil, []byte("bin")}

	got, err := vm.BoxValue(reflectOf(row))
	if err != nil {
		t.Fatalf("BoxValue: %v", err)
	}
	vec, ok := got.(vm.ArrayVector)
	if !ok {
		t.Fatalf("expected ArrayVector, got %T", got)
	}
	if len(vec) != len(row) {
		t.Fatalf("expected %d elements, got %d", len(row), len(vec))
	}

	tests := []struct {
		idx  int
		want vm.Value
	}{
		{0, vm.String("Ada")},
		{1, vm.Int(36)},
		{2, vm.Float(1.5)},
		{3, vm.TRUE},
		// A nil interface element must reach let-go as NIL, not as a boxing
		// error. This is the SQL NULL case.
		{4, vm.NIL},
		// []byte keeps its own fast path once unwrapped.
		{5, vm.String("bin")},
	}
	for _, tt := range tests {
		el := vec[tt.idx]
		if _, boxed := el.(*vm.Boxed); boxed {
			t.Errorf("element %d is a vm.Boxed (%v), want native %T", tt.idx, el, tt.want)
			continue
		}
		if el != tt.want {
			t.Errorf("element %d = %v (%T), want %v (%T)", tt.idx, el, el, tt.want, tt.want)
		}
	}
}

// A []any nested inside a []any must box recursively, not become one opaque box.
func TestBoxValueNestedInterfaceSlice(t *testing.T) {
	got, err := vm.BoxValue(reflectOf([]any{"outer", []any{"inner", int64(1)}}))
	if err != nil {
		t.Fatalf("BoxValue: %v", err)
	}
	vec, ok := got.(vm.ArrayVector)
	if !ok || len(vec) != 2 {
		t.Fatalf("expected a 2-element ArrayVector, got %T %v", got, got)
	}
	inner, ok := vec[1].(vm.ArrayVector)
	if !ok {
		t.Fatalf("nested element = %T (%v), want ArrayVector", vec[1], vec[1])
	}
	if inner[0] != vm.String("inner") || inner[1] != vm.Int(1) {
		t.Errorf("nested elements = %v, %v; want \"inner\", 1", inner[0], inner[1])
	}
}

// The concrete-element fast paths specialize on the slice's static type and
// must not change: []byte is a String, []int64 and []float64 keep their packed
// array representations rather than becoming generic vectors.
func TestBoxValueConcreteSliceFastPathsUnchanged(t *testing.T) {
	s, err := vm.BoxValue(reflectOf([]byte("hello")))
	if err != nil {
		t.Fatalf("BoxValue([]byte): %v", err)
	}
	if s != vm.String("hello") {
		t.Errorf("[]byte boxed to %v (%T), want String", s, s)
	}

	ints, err := vm.BoxValue(reflectOf([]int64{1, 2, 3}))
	if err != nil {
		t.Fatalf("BoxValue([]int64): %v", err)
	}
	if _, ok := ints.(vm.ArrayVector); ok {
		t.Errorf("[]int64 fell through to the generic vector path (%T)", ints)
	}

	floats, err := vm.BoxValue(reflectOf([]float64{1.5, 2.5}))
	if err != nil {
		t.Fatalf("BoxValue([]float64): %v", err)
	}
	if _, ok := floats.(vm.ArrayVector); ok {
		t.Errorf("[]float64 fell through to the generic vector path (%T)", floats)
	}
}

// Boxing an element by its dynamic type must never make a previously-boxable
// []any worse. Two of BoxValue's own paths panic rather than erroring — the
// slice/array case calls IsNil (invalid for arrays), and the
// []byte/[]int64/[]float64 fast paths assert the exact slice type — and a
// defined scalar like json.Number is rejected outright by StringType.Box.
// Before dynamic boxing every one of these stayed an opaque Boxed; none may
// now panic, and none may fail the whole slice.
func TestBoxValueDynamicElementsNeverRegress(t *testing.T) {
	tests := []struct {
		name string
		in   any
	}{
		{"Go array element (BoxValue's IsNil is invalid for arrays)", []any{[1]int{1}}},
		{"defined slice over []int64 (fast path asserts the exact type)", []any{definedInts{1}}},
		{"slice of a defined element type", []any{[]definedInt{1}}},
		{"defined string type (StringType.Box wants exactly string)", []any{jsonNumber("1")}},
		{"defined bool type", []any{definedBool(true)}},
		// Composites reach the unsafe paths recursively, so a top-level-only
		// check is not enough.
		{"slice of arrays (nested IsNil)", []any{[][1]int{{1}}}},
		{"map of a defined slice (nested assertion)", []any{map[string]definedInts{"x": {1}}}},
		// The nastiest: ChanType.Box spawns a goroutine that calls Recv, which
		// for a send-only channel panics in ANOTHER goroutine and kills the
		// process. No recover at the call site can catch that.
		{"send-only channel", []any{(chan<- int)(make(chan int, 1))}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("panicked: %v", r)
				}
			}()
			got, err := vm.BoxValue(reflectOf(tt.in))
			if err != nil {
				t.Fatalf("boxing failed for the whole slice: %v", err)
			}
			vec, ok := got.(vm.ArrayVector)
			if !ok || len(vec) != 1 {
				t.Fatalf("got %T %v, want a 1-element ArrayVector", got, got)
			}
			if vec[0] == vm.NIL {
				t.Errorf("element was dropped to NIL, want it preserved")
			}
		})
	}
}

type definedInts []int64
type definedInt int64
type jsonNumber string
type definedBool bool

func reflectOf(v any) reflect.Value { return reflect.ValueOf(v) }
