/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"fmt"
	"strings"
	"testing"
)

// Lowered Go boxes a closure per fn literal, so a box must cost one allocation
// (the NativeFn itself) and a call none, for every shape the lowering emits.
func TestNativeFnBoxAndCallAllocations(t *testing.T) {
	x := Int(1)
	shapes := []struct {
		name string
		fn   any
		args []Value
	}{
		{"0 args with error", func() (Value, error) { return x, nil }, nil},
		{"1 arg with error", func(a Value) (Value, error) { return a, nil }, []Value{x}},
		{"2 args with error", func(a, b Value) (Value, error) { return b, nil }, []Value{x, x}},
		{"3 args with error", func(a, b, c Value) (Value, error) { return c, nil }, []Value{x, x, x}},
		{"4 args with error", func(a, b, c, d Value) (Value, error) { return d, nil }, []Value{x, x, x, x}},
		{"0 args, no error", func() Value { return x }, nil},
		{"1 arg, no error", func(a Value) Value { return a }, []Value{x}},
		{"2 args, no error", func(a, b Value) Value { return b }, []Value{x, x}},
		{"3 args, no error", func(a, b, c Value) Value { return c }, []Value{x, x, x}},
		{"4 args, no error", func(a, b, c, d Value) Value { return d }, []Value{x, x, x, x}},
		{"variadic with error", func(as ...Value) (Value, error) { return x, nil }, []Value{x, x, x}},
	}
	for _, s := range shapes {
		t.Run(s.name, func(t *testing.T) {
			if n := testing.AllocsPerRun(100, func() { _, _ = NativeFnType.Box(s.fn) }); n != 1 {
				t.Errorf("Box: %v allocations, want 1", n)
			}
			boxed, err := NativeFnType.Box(s.fn)
			if err != nil {
				t.Fatalf("Box: %v", err)
			}
			f := boxed.(*NativeFn)
			if got, err := f.Invoke(s.args); err != nil || got != x {
				t.Fatalf("Invoke: got %v, %v", got, err)
			}
			if n := testing.AllocsPerRun(100, func() { _, _ = f.Invoke(s.args) }); n != 0 {
				t.Errorf("Invoke: %v allocations, want 0", n)
			}
		})
	}
}

// A fixed-arity fast shape called with the wrong number of arguments fails as
// the bytecode VM does for a fn of that arity, and a nil interface result is
// NIL, as the reflect path's BoxValue makes it.
func TestNativeFnFastShapesCheckArityAndNormalizeNil(t *testing.T) {
	x := Int(1)
	shapes := []struct {
		name  string
		fn    any
		arity int
	}{
		{"0 with error", func() (Value, error) { return nil, nil }, 0},
		{"1 with error", func(a Value) (Value, error) { return nil, nil }, 1},
		{"2 with error", func(a, b Value) (Value, error) { return nil, nil }, 2},
		{"3 with error", func(a, b, c Value) (Value, error) { return nil, nil }, 3},
		{"4 with error", func(a, b, c, d Value) (Value, error) { return nil, nil }, 4},
		{"0 no error", func() Value { return nil }, 0},
		{"1 no error", func(a Value) Value { return nil }, 1},
		{"2 no error", func(a, b Value) Value { return nil }, 2},
		{"3 no error", func(a, b, c Value) Value { return nil }, 3},
		{"4 no error", func(a, b, c, d Value) Value { return nil }, 4},
	}
	for _, s := range shapes {
		t.Run(s.name, func(t *testing.T) {
			boxed, err := NativeFnType.Box(s.fn)
			if err != nil {
				t.Fatalf("Box: %v", err)
			}
			f := boxed.(*NativeFn)
			args := make([]Value, s.arity)
			for i := range args {
				args[i] = x
			}
			if got, err := f.Invoke(args); err != nil || got != NIL {
				t.Fatalf("Invoke: got %#v, %v; want NIL", got, err)
			}
			_, err = f.Invoke(append(args, x))
			want := fmt.Sprintf("function %s expected %d args, got %d", f, s.arity, s.arity+1)
			if err == nil || !strings.Contains(err.Error(), want) {
				t.Fatalf("wrong arity: got %v, want %q", err, want)
			}
		})
	}
}

// A boxed func(...Value) takes any number of arguments, zero included, so it
// declares arity -1 and variadic. Multi-arity dispatch reads that as a rest
// arm with no fixed parameters; the reflect path's NumIn() of 1 counted the
// variadic parameter as a fixed one.
func TestVariadicFastShapeDeclaresAnyArity(t *testing.T) {
	boxed, err := NativeFnType.Box(func(as ...Value) (Value, error) { return Int(len(as)), nil })
	if err != nil {
		t.Fatal(err)
	}
	f := boxed.(*NativeFn)
	if f.Arity() != -1 || !f.IsVariadic() {
		t.Fatalf("Arity() = %d, IsVariadic() = %v; want -1, true", f.Arity(), f.IsVariadic())
	}
	if got, err := f.Invoke(nil); err != nil || got != Int(0) {
		t.Fatalf("Invoke() = %v, %v; want 0", got, err)
	}
}
