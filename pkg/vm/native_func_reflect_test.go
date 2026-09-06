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
