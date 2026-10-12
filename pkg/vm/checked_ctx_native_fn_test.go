/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"strings"
	"testing"
)

// A checked context-aware native rejects a call the shape doesn't accept with
// the bytecode VM's error for that shape, and hands every accepted call to fn.
func TestCheckedCtxNativeFnReportsArityAsTheVM(t *testing.T) {
	count := func(_ *ExecContext, args []Value) (Value, error) { return Int(len(args)), nil }
	cases := []struct {
		name    string
		fixed   []int
		restMin int
		ok      []int
		bad     int
		want    string
	}{
		{"one fixed arity", []int{2}, -1, []int{2}, 1, "expected 2 args, got 1"},
		{"rest only", nil, 1, []int{1, 3}, 0, "expected at least 1 args, got 0"},
		{"several arities", []int{0, 2}, -1, []int{0, 2}, 1, "doesn't have a 1-arity variant"},
		{"fixed and rest", []int{1}, 3, []int{1, 3, 5}, 2, "doesn't have a 2-arity variant"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := NewCheckedCtxNativeFn("", tc.fixed, tc.restMin, count)
			for _, n := range tc.ok {
				got, err := f.Invoke(make([]Value, n))
				if err != nil || got != Int(n) {
					t.Fatalf("%d args: got %v, %v", n, got, err)
				}
			}
			_, err := f.Invoke(make([]Value, tc.bad))
			if err == nil || !strings.Contains(err.Error(), "function "+f.String()+" "+tc.want) {
				t.Fatalf("%d args: got %v, want %q", tc.bad, err, tc.want)
			}
		})
	}
}
