//go:build !tinygo

/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"errors"
	"testing"
)

// The reflect proxy assumed a Go function returns at most two values, and
// silently dropped anything past the second. That drop happens on the Go side,
// before boxing, so no amount of .lg veneer can recover the lost values — and a
// wrapper library has no Go of its own to shim with. (a, b, ok) and
// (a, b, err) are ordinary modern Go: strings.Cut returns three.
func TestBoxReflectFuncMultiReturn(t *testing.T) {
	boom := errors.New("boom")

	tests := []struct {
		name    string
		fn      any
		want    string
		wantErr error
	}{
		// Existing shapes must not move.
		{name: "no results", fn: func() {}, want: "nil"},
		{name: "one result", fn: func() int { return 7 }, want: "7"},
		{name: "value and nil error", fn: func() (int, error) { return 7, nil }, want: "7"},
		{
			name:    "value and non-nil error throws",
			fn:      func() (int, error) { return 0, boom },
			want:    "0",
			wantErr: boom,
		},

		// Previously dropped.
		{name: "two non-error results", fn: func() (int, string) { return 1, "a" }, want: `[1 "a"]`},
		{name: "three non-error results", fn: func() (int, string, bool) { return 1, "a", true }, want: `[1 "a" true]`},
		{
			name: "strings.Cut shape: two values and a bool",
			fn:   func() (string, string, bool) { return "k", "v", true },
			want: `["k" "v" true]`,
		},
		{
			name: "two values and a nil error",
			fn:   func() (int, string, error) { return 1, "a", nil },
			want: `[1 "a"]`,
		},
		{
			name:    "two values and a non-nil error throws",
			fn:      func() (int, string, error) { return 1, "a", boom },
			want:    `[1 "a"]`,
			wantErr: boom,
		},
		{
			name: "four values",
			fn:   func() (int, int, int, int) { return 1, 2, 3, 4 },
			want: "[1 2 3 4]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			boxed, err := NativeFnType.Box(tt.fn)
			if err != nil {
				t.Fatalf("Box: %v", err)
			}
			fn, ok := boxed.(*NativeFn)
			if !ok {
				t.Fatalf("Box returned %T, want *NativeFn", boxed)
			}
			got, err := fn.proxy(nil)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("proxy err = %v, want %v", err, tt.wantErr)
			}
			rendered, rerr := SafeString(got)
			if rerr != nil {
				t.Fatalf("SafeString: %v", rerr)
			}
			if rendered != tt.want {
				t.Errorf("proxy() = %s, want %s", rendered, tt.want)
			}
		})
	}
}

func TestBoxReflectFuncTypedArrayPreservesBackingSlice(t *testing.T) {
	invoke := func(t *testing.T, fn any, arg Value) {
		t.Helper()
		boxed, err := NativeFnType.Box(fn)
		if err != nil {
			t.Fatalf("Box: %v", err)
		}
		if _, err := boxed.(*NativeFn).proxy([]Value{arg}); err != nil {
			t.Fatalf("proxy: %v", err)
		}
	}

	t.Run("byte array", func(t *testing.T) {
		backing := []byte{1, 2}
		arr := NewByteArrayFrom(backing)
		sameBacking := false
		invoke(t, func(xs []byte) {
			sameBacking = &xs[0] == &backing[0]
			xs[0] = 9
		}, arr)
		if !sameBacking || backing[0] != 9 || arr.Get(0) != Int(9) {
			t.Fatalf("byte-array backing was copied: same=%v backing=%v value=%v", sameBacking, backing, arr.Get(0))
		}
	})

	t.Run("int array", func(t *testing.T) {
		backing := []int64{1, 2}
		arr := NewIntArrayFrom(backing)
		sameBacking := false
		invoke(t, func(xs []int64) {
			sameBacking = &xs[0] == &backing[0]
			xs[0] = 9
		}, arr)
		if !sameBacking || backing[0] != 9 || arr.Get(0) != Int(9) {
			t.Fatalf("int-array backing was copied: same=%v backing=%v value=%v", sameBacking, backing, arr.Get(0))
		}
	})

	t.Run("double array", func(t *testing.T) {
		backing := []float64{1, 2}
		arr := NewFloatArrayFrom(backing)
		sameBacking := false
		invoke(t, func(xs []float64) {
			sameBacking = &xs[0] == &backing[0]
			xs[0] = 9.5
		}, arr)
		if !sameBacking || backing[0] != 9.5 || arr.Get(0) != Float(9.5) {
			t.Fatalf("double-array backing was copied: same=%v backing=%v value=%v", sameBacking, backing, arr.Get(0))
		}
	})

	t.Run("object array", func(t *testing.T) {
		backing := []Value{Int(1), Int(2)}
		arr := NewObjectArrayFrom(backing)
		sameBacking := false
		invoke(t, func(xs []Value) {
			sameBacking = &xs[0] == &backing[0]
			xs[0] = String("changed")
		}, arr)
		if !sameBacking || backing[0] != String("changed") || arr.Get(0) != String("changed") {
			t.Fatalf("object-array backing was copied: same=%v backing=%v value=%v", sameBacking, backing, arr.Get(0))
		}
	})

	t.Run("compatible named slice", func(t *testing.T) {
		type namedBytes []byte
		backing := []byte{1, 2}
		arr := NewByteArrayFrom(backing)
		sameBacking := false
		invoke(t, func(xs namedBytes) {
			sameBacking = &xs[0] == &backing[0]
			xs[0] = 9
		}, arr)
		if !sameBacking || arr.Get(0) != Int(9) {
			t.Fatalf("named slice lost typed-array backing: same=%v value=%v", sameBacking, arr.Get(0))
		}
	})
}

func TestBoxReflectFuncPersistentCollectionsStillConvertToSlices(t *testing.T) {
	for _, tt := range []struct {
		name string
		arg  Value
	}{
		{name: "vector", arg: NewArrayVector([]Value{Int(1), Int(2)})},
		{name: "list", arg: NewList([]Value{Int(1), Int(2)})},
	} {
		t.Run(tt.name, func(t *testing.T) {
			boxed, err := NativeFnType.Box(func(xs []int) int {
				xs[0] = 9
				return xs[0] + xs[1]
			})
			if err != nil {
				t.Fatalf("Box: %v", err)
			}
			got, err := boxed.(*NativeFn).proxy([]Value{tt.arg})
			if err != nil {
				t.Fatalf("proxy: %v", err)
			}
			if got != Int(11) {
				t.Fatalf("converted slice result = %v, want 11", got)
			}
			if first := tt.arg.(Sequable).Seq().First(); first != Int(1) {
				t.Fatalf("persistent source was mutated through copied []int: %v", first)
			}
		})
	}
}

// func() error must keep returning the error as a VALUE rather than throwing.
// Every reflect-boxed Close/Write/Flush has this shape, so peeling here would
// silently change how a large amount of existing interop behaves. The peel
// starts at (T, error), exactly where it always did.
func TestBoxReflectFuncSoleErrorStaysAValue(t *testing.T) {
	boom := errors.New("boom")
	for _, tt := range []struct {
		name string
		fn   any
		want string
	}{
		{name: "non-nil error", fn: func() error { return boom }, want: "<go.*errors.errorString boom>"},
		{name: "nil error", fn: func() error { return nil }, want: "nil"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			boxed, err := NativeFnType.Box(tt.fn)
			if err != nil {
				t.Fatalf("Box: %v", err)
			}
			got, err := boxed.(*NativeFn).proxy(nil)
			if err != nil {
				t.Fatalf("func() error must not throw, got %v", err)
			}
			rendered, rerr := SafeString(got)
			if rerr != nil {
				t.Fatalf("SafeString: %v", rerr)
			}
			if rendered != tt.want {
				t.Errorf("proxy() = %s, want %s", rendered, tt.want)
			}
		})
	}
}

// Known limitation, pinned so a future fix is deliberate rather than
// accidental: BoxValue does not support Go ARRAY values. Its
// `case reflect.Slice, reflect.Array` calls IsNil — invalid for arrays — and
// the branches under it assert []int64 / call Bytes(), neither of which holds
// for [N]T. This predates multi-return: func() [1]int panics the same way and
// always has. NativeFn.Invoke recovers it into an error, so callers see a
// failure rather than a crash.
//
// Multi-return makes the existing hole reachable through one more shape. That
// is a visible error instead of the silent wrong answer it used to give (the
// array result was simply dropped), so it is not a regression — but it is not
// a fix either, and fixing BoxValue is out of scope here.
func TestBoxReflectFuncArrayResultIsAKnownLimitation(t *testing.T) {
	for _, tt := range []struct {
		name string
		fn   any
	}{
		{name: "sole array result (pre-existing)", fn: func() [1]int { return [1]int{9} }},
		{name: "array among several results", fn: func() (int, [1]int) { return 1, [1]int{9} }},
	} {
		t.Run(tt.name, func(t *testing.T) {
			boxed, err := NativeFnType.Box(tt.fn)
			if err != nil {
				t.Fatalf("Box: %v", err)
			}
			// Invoke, not proxy: Invoke is what recovers the panic into an error.
			if _, err := boxed.(*NativeFn).Invoke(nil); err == nil {
				t.Errorf("expected an error for an array result; if BoxValue " +
					"learned to box arrays, update this test to assert the value")
			}
		})
	}
}

// A trailing error is peeled off as a throw only when it is actually the last
// result. An error in any other position is an ordinary value.
func TestBoxReflectFuncErrorOnlyPeeledWhenTrailing(t *testing.T) {
	boom := errors.New("boom")
	boxed, err := NativeFnType.Box(func() (error, int) { return boom, 3 })
	if err != nil {
		t.Fatalf("Box: %v", err)
	}
	got, err := boxed.(*NativeFn).proxy(nil)
	if err != nil {
		t.Fatalf("a non-trailing error must not throw, got %v", err)
	}
	v, ok := got.(ArrayVector)
	if !ok || len(v) != 2 {
		t.Fatalf("got %#v, want a 2-element vector", got)
	}
}
