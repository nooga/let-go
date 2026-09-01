/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package vm_test

import (
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

// BoxValue's map branch boxes each key and value with plain BoxValue, which
// reads the map's STATIC element type. For a map[string]any that kind is
// Interface, so every value misses the string/int/bool fast paths and lands as
// an opaque vm.Boxed — printing as <go.string Ada> and comparing equal to
// nothing in let-go:
//
//	(= "Ada" (get m "name"))   ; => false
//
// map[string]any is the natural Go type for "a bag of options" or "a JSON
// object", which is exactly the shape a Wails v3 binding hands back. Values
// must be boxed by their DYNAMIC type instead, the way slice elements already
// are.
func TestBoxValueMapInterfaceValues(t *testing.T) {
	got, err := vm.BoxValue(reflectOf(map[string]any{
		"name": "Ada",
		"age":  int64(36),
		"ok":   true,
		"none": nil,
	}))
	if err != nil {
		t.Fatalf("BoxValue: %v", err)
	}
	m, ok := got.(*vm.PersistentMap)
	if !ok {
		t.Fatalf("expected *vm.PersistentMap, got %T", got)
	}
	if m.RawCount() != 4 {
		t.Fatalf("expected 4 entries, got %d", m.RawCount())
	}

	// Compared by value, not by type name, and with an explicit Boxed
	// rejection: an opaque Boxed can carry the right Go value and still be
	// equal to nothing, which is the actual bug.
	tests := []struct {
		key  string
		want vm.Value
	}{
		{"name", vm.String("Ada")},
		{"age", vm.Int(36)},
		{"ok", vm.TRUE},
		{"none", vm.NIL},
	}
	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			v := m.ValueAt(vm.String(tt.key))
			if _, boxed := v.(*vm.Boxed); boxed {
				t.Fatalf("value at %q is a vm.Boxed (%v), want native %T", tt.key, v, tt.want)
			}
			if v != tt.want {
				t.Fatalf("value at %q: want %s (%T), got %s (%T)", tt.key, tt.want, tt.want, v, v)
			}
		})
	}
}

// A Go string key boxes to a let-go STRING, not a keyword: Go map keys may hold
// spaces and dots, which do not make valid keywords. So a map round-trips as
// {"name" "Ada"} and is read with (get m "name"), matching clojure.data.json
// with no :key-fn.
func TestBoxValueMapKeysAreStrings(t *testing.T) {
	got, err := vm.BoxValue(reflectOf(map[string]any{"first name": "Ada"}))
	if err != nil {
		t.Fatalf("BoxValue: %v", err)
	}
	m, ok := got.(*vm.PersistentMap)
	if !ok {
		t.Fatalf("expected *vm.PersistentMap, got %T", got)
	}
	if v := m.ValueAt(vm.String("first name")); v != vm.String("Ada") {
		t.Fatalf("string key lookup: got %s (%T)", v, v)
	}
	if v := m.ValueAt(vm.Keyword("first name")); v != vm.NIL {
		t.Fatalf("keyword lookup should miss, got %s (%T)", v, v)
	}
}

// A nil Go map is let-go nil; an EMPTY Go map is an empty let-go map, which is
// a different value and must not collapse to nil.
func TestBoxValueMapNilAndEmpty(t *testing.T) {
	var nilMap map[string]any
	got, err := vm.BoxValue(reflectOf(nilMap))
	if err != nil {
		t.Fatalf("BoxValue(nil map): %v", err)
	}
	if got != vm.NIL {
		t.Fatalf("nil map: want NIL, got %s (%T)", got, got)
	}

	got, err = vm.BoxValue(reflectOf(map[string]any{}))
	if err != nil {
		t.Fatalf("BoxValue(empty map): %v", err)
	}
	m, ok := got.(*vm.PersistentMap)
	if !ok {
		t.Fatalf("empty map: expected *vm.PersistentMap, got %T", got)
	}
	if m.RawCount() != 0 {
		t.Fatalf("empty map: expected 0 entries, got %d", m.RawCount())
	}
}

// Nesting works because the allowlist admits map[string]any and []any, whose
// elements are themselves interfaces and recurse through it.
func TestBoxValueMapNested(t *testing.T) {
	got, err := vm.BoxValue(reflectOf(map[string]any{
		"inner": map[string]any{"b": int64(1)},
		"xs":    []any{"a", int64(2)},
	}))
	if err != nil {
		t.Fatalf("BoxValue: %v", err)
	}
	m, ok := got.(*vm.PersistentMap)
	if !ok {
		t.Fatalf("expected *vm.PersistentMap, got %T", got)
	}

	inner, ok := m.ValueAt(vm.String("inner")).(*vm.PersistentMap)
	if !ok {
		t.Fatalf("inner: expected *vm.PersistentMap, got %T", m.ValueAt(vm.String("inner")))
	}
	if v := inner.ValueAt(vm.String("b")); v != vm.Int(1) {
		t.Fatalf("inner b: want 1, got %s (%T)", v, v)
	}

	xs, ok := m.ValueAt(vm.String("xs")).(vm.ArrayVector)
	if !ok {
		t.Fatalf("xs: expected vm.ArrayVector, got %T", m.ValueAt(vm.String("xs")))
	}
	if len(xs) != 2 || xs[0] != vm.String("a") || xs[1] != vm.Int(2) {
		t.Fatalf("xs: want [\"a\" 2], got %s", xs)
	}
}

// A []any holding a map[string]any exercises the mirror of the case above: the
// slice branch's unwrapping has to admit maps for a row of JSON objects to
// arrive as let-go maps.
func TestBoxValueMapInsideSlice(t *testing.T) {
	got, err := vm.BoxValue(reflectOf([]any{map[string]any{"b": int64(1)}}))
	if err != nil {
		t.Fatalf("BoxValue: %v", err)
	}
	vec, ok := got.(vm.ArrayVector)
	if !ok {
		t.Fatalf("expected vm.ArrayVector, got %T", got)
	}
	if len(vec) != 1 {
		t.Fatalf("expected 1 element, got %d", len(vec))
	}
	m, ok := vec[0].(*vm.PersistentMap)
	if !ok {
		t.Fatalf("element: expected *vm.PersistentMap, got %T", vec[0])
	}
	if v := m.ValueAt(vm.String("b")); v != vm.Int(1) {
		t.Fatalf("b: want 1, got %s (%T)", v, v)
	}
}
