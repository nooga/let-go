/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package vm_test

import (
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

// Converting a let-go vector into a Go []any parameter gave up whenever an
// element had no Go counterpart, handed the whole slice over as []vm.Value,
// and the reflect call then failed:
//
//	reflect: Call using []vm.Value as type []interface {}
//
// A let-go nil is the common trigger, which makes every SQL NULL parameter
// fail. []any can hold anything, Go nil included, so the conversion should
// always succeed for it.
//
// This drives the real path — boxArgForReflect on a boxed Go func — rather
// than calling the conversion helper directly, so it fails the way a caller
// would see it.
func TestUnboxVectorIntoAnySlice(t *testing.T) {
	var captured []any
	fn, err := vm.NativeFnType.Box(func(args []any) vm.Value {
		captured = args
		return vm.Int(int64(len(args)))
	})
	if err != nil {
		t.Fatalf("Box: %v", err)
	}

	tests := []struct {
		name string
		in   vm.Value
		want []any
	}{
		// Element values are whatever Value.Unbox yields — vm.Int unboxes to
		// int, not int64. That is existing behavior on the non-nil path and
		// is deliberately left alone here.
		{
			name: "mixed scalars",
			in:   vm.ArrayVector{vm.String("Ada"), vm.Int(36), vm.Float(1.5), vm.TRUE},
			want: []any{"Ada", 36, 1.5, true},
		},
		{
			name: "nil element becomes Go nil (the SQL NULL parameter case)",
			in:   vm.ArrayVector{vm.String("Ada"), vm.NIL, vm.Int(1)},
			want: []any{"Ada", nil, 1},
		},
		{
			name: "all nil",
			in:   vm.ArrayVector{vm.NIL, vm.NIL},
			want: []any{nil, nil},
		},
		{
			name: "empty vector is an empty slice, not nil",
			in:   vm.ArrayVector{},
			want: []any{},
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
				t.Fatalf("callee saw %v args, want %d", got, len(tt.want))
			}
			if captured == nil {
				t.Fatalf("callee received a nil slice, want a non-nil %d-element slice", len(tt.want))
			}
			if len(captured) != len(tt.want) {
				t.Fatalf("captured %#v, want %#v", captured, tt.want)
			}
			for i := range tt.want {
				if captured[i] != tt.want[i] {
					t.Errorf("element %d = %#v (%T), want %#v (%T)",
						i, captured[i], captured[i], tt.want[i], tt.want[i])
				}
			}
		})
	}
}
