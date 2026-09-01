/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package vm_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

// Handing a let-go map to a Go function declaring map[string]any died with:
//
//	reflect: Call using *vm.PersistentMap as type map[string]interface {}
//
// boxArgForReflect special-cases slice/array and interface targets; a MAP
// target fell through to reflect.ValueOf(v.Unbox()), and (*PersistentMap).Unbox
// returns the map itself. map[string]any is the shape a Wails v3 binding
// declares for an options bag, so a wrapper author had to take a vm.Value and
// write the recursive converter by hand.
//
// These drive the real path — boxArgForReflect on a boxed Go func — rather than
// calling the conversion helper directly, so they fail the way a caller would
// see it.
func TestUnboxMapIntoStringAnyMap(t *testing.T) {
	var captured map[string]any
	fn, err := vm.NativeFnType.Box(func(m map[string]any) vm.Value {
		captured = m
		return vm.Int(int64(len(m)))
	})
	if err != nil {
		t.Fatalf("Box: %v", err)
	}

	tests := []struct {
		name string
		in   vm.Value
		want map[string]any
	}{
		// A Keyword stores its name without the leading colon, and unboxInto's
		// existing reflect.String case already accepts one. So keyword keys
		// becoming string keys is an existing convention, not a new rule.
		{
			name: "keyword keys",
			in: vm.NewPersistentMap([]vm.Value{
				vm.Keyword("name"), vm.String("Ada"),
				vm.Keyword("age"), vm.Int(36),
			}),
			// Values are whatever Value.Unbox yields — vm.Int unboxes to int,
			// not int64. Existing behaviour, deliberately left alone.
			want: map[string]any{"name": "Ada", "age": 36},
		},
		{
			name: "string keys",
			in: vm.NewPersistentMap([]vm.Value{
				vm.String("name"), vm.String("Ada"),
			}),
			want: map[string]any{"name": "Ada"},
		},
		{
			name: "nil value becomes Go nil",
			in: vm.NewPersistentMap([]vm.Value{
				vm.Keyword("a"), vm.NIL,
				vm.Keyword("b"), vm.Int(1),
			}),
			want: map[string]any{"a": nil, "b": 1},
		},
		{
			name: "empty map is an empty map, not nil",
			in:   vm.EmptyPersistentMap,
			want: map[string]any{},
		},
		// Every let-go map type has to convert, and they are identified by
		// TYPE, not by inspecting the first element: an empty map of any type
		// yields EmptyList from Seq() and would be indistinguishable from an
		// empty vector.
		{
			name: "plain Map",
			in:   vm.Map{vm.Keyword("a"): vm.Int(1)},
			want: map[string]any{"a": 1},
		},
		{
			name: "sorted map",
			in:   vm.NewSortedMap(nil, []vm.Value{vm.String("a"), vm.Int(1), vm.String("b"), vm.Int(2)}),
			want: map[string]any{"a": 1, "b": 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			captured = nil
			got, err := fn.(*vm.NativeFn).Invoke([]vm.Value{tt.in})
			if err != nil {
				t.Fatalf("Invoke: %v", err)
			}
			if int(got.(vm.Int)) != len(tt.want) {
				t.Fatalf("callee saw %v entries, want %d", got, len(tt.want))
			}
			if captured == nil {
				t.Fatalf("callee received a nil map, want a non-nil %d-entry map", len(tt.want))
			}
			if !reflect.DeepEqual(captured, tt.want) {
				t.Fatalf("captured %#v, want %#v", captured, tt.want)
			}
		})
	}
}

// A typed map target converts element by element the same way a typed slice
// does, so map[string]string gets Go strings rather than vm.String.
func TestUnboxMapIntoTypedMap(t *testing.T) {
	var captured map[string]string
	fn, err := vm.NativeFnType.Box(func(m map[string]string) vm.Value {
		captured = m
		return vm.Int(int64(len(m)))
	})
	if err != nil {
		t.Fatalf("Box: %v", err)
	}

	in := vm.NewPersistentMap([]vm.Value{
		vm.Keyword("host"), vm.String("localhost"),
		vm.String("port"), vm.String("8080"),
	})
	if _, err := fn.(*vm.NativeFn).Invoke([]vm.Value{in}); err != nil {
		t.Fatalf("Invoke: %v", err)
	}
	want := map[string]string{"host": "localhost", "port": "8080"}
	if !reflect.DeepEqual(captured, want) {
		t.Fatalf("captured %#v, want %#v", captured, want)
	}
}

// Go reports int64 as ConvertibleTo string and converts it to a RUNE, so the
// generic fallback would silently turn the let-go key 65 into the Go key "A".
// Reject instead: a wrong key that looks plausible is worse than a failure.
func TestUnboxMapRejectsNumericKeyIntoStringKey(t *testing.T) {
	var captured map[string]any
	fn, err := vm.NativeFnType.Box(func(m map[string]any) vm.Value {
		captured = m
		return vm.Int(int64(len(m)))
	})
	if err != nil {
		t.Fatalf("Box: %v", err)
	}

	in := vm.NewPersistentMap([]vm.Value{vm.Int(65), vm.String("Ada")})

	// The conversion failing leaves boxArgForReflect on its existing fallback,
	// which hands reflect.Call a *vm.PersistentMap and panics. Either outcome
	// is a rejection; silently producing map[string]any{"A": "Ada"} is not.
	func() {
		defer func() { _ = recover() }()
		if _, err := fn.(*vm.NativeFn).Invoke([]vm.Value{in}); err != nil {
			return
		}
	}()

	if captured != nil {
		t.Fatalf("numeric key was accepted: callee saw %#v, want the conversion rejected", captured)
	}
}

// A map[any]any key target accepts any dynamic type, and a let-go vector key
// unboxes to a []vm.Value — which SetMapIndex cannot hash. That panicked with
// "hash of unhashable type []vm.Value" instead of failing the conversion.
// unboxMapInto is also reached from RecordToStruct, which has no recover of its
// own, so the panic is not guaranteed to be contained.
func TestUnboxMapRejectsUnhashableKey(t *testing.T) {
	var captured map[any]any
	fn, err := vm.NativeFnType.Box(func(m map[any]any) vm.Value {
		captured = m
		return vm.Int(int64(len(m)))
	})
	if err != nil {
		t.Fatalf("Box: %v", err)
	}

	in := vm.NewPersistentMap([]vm.Value{vm.ArrayVector{vm.Int(1)}, vm.String("v")})
	_, invokeErr := fn.(*vm.NativeFn).Invoke([]vm.Value{in})
	if invokeErr == nil {
		t.Fatalf("expected the call to fail, callee saw %#v", captured)
	}
	if captured != nil {
		t.Fatalf("unhashable key was accepted: callee saw %#v", captured)
	}
	if strings.Contains(invokeErr.Error(), "hash of unhashable type") {
		t.Fatalf("conversion panicked instead of failing cleanly: %v", invokeErr)
	}

	// A comparable dynamic type still works, so the guard is not blanket.
	captured = nil
	if _, err := fn.(*vm.NativeFn).Invoke([]vm.Value{
		vm.NewPersistentMap([]vm.Value{vm.Keyword("a"), vm.Int(1)}),
	}); err != nil {
		t.Fatalf("Invoke: %v", err)
	}
	if want := (map[any]any{"a": 1}); !reflect.DeepEqual(captured, want) {
		t.Fatalf("captured %#v, want %#v", captured, want)
	}
}
