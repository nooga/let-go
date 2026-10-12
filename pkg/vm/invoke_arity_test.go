/*
 * Copyright (c) 2026 let-go contributors; see CONTRIBUTORS.
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"testing"
)

// pairClosure is the shape lowered Go emits for a two-argument fn literal:
// CallNative for the slice path, Invoke2 for the typed one.
type pairClosure struct {
	NativeFn
	y Value
}

func (c *pairClosure) CallNative(args []Value) (Value, error) {
	if len(args) != 2 {
		return NIL, WrongArgCount(&c.NativeFn, len(args), 2)
	}
	return c.Invoke2(nil, args[0], args[1])
}

func (c *pairClosure) Invoke2(_ *ExecContext, a, b Value) (Value, error) {
	return Int(int(a.(Int)) + int(b.(Int)) + int(c.y.(Int))), nil
}

func newPairClosure(y Value) Value {
	c := &pairClosure{y: y}
	InitNativeFn(&c.NativeFn, 2, false, c)
	return &c.NativeFn
}

// An arity-typed call reaches every callee kind lowered code meets without
// allocating for its arguments, and a call of the wrong arity through the
// typed entry fails with exactly the slice path's error.
func TestInvokeArityAllocatesNothingForArgs(t *testing.T) {
	ec := RootExecContext
	one, two := Int(1), Int(2)
	fast2, err := NativeFnType.Box(func(a, b Value) (Value, error) { return a, nil })
	if err != nil {
		t.Fatal(err)
	}
	fast1, err := NativeFnType.Box(func(a Value) Value { return a })
	if err != nil {
		t.Fatal(err)
	}
	typed := NewTypedCtxNativeFn("typed", func(_ *ExecContext, a, b Value) (Value, error) { return b, nil })
	lowered := newPairClosure(Int(10)).(*NativeFn)
	m := NewPersistentMap([]Value{Keyword("k"), one})

	callees := []struct {
		name string
		call func() (Value, error)
		want Value
	}{
		{"fast shape, 2 args", func() (Value, error) { return ec.Invoke2(fast2.(Fn), one, two) }, one},
		{"fast shape, 1 arg no error", func() (Value, error) { return ec.Invoke1(fast1.(Fn), two) }, two},
		{"typed ctx wrapper", func() (Value, error) { return ec.Invoke2(typed, one, two) }, two},
		{"lowered closure struct", func() (Value, error) { return ec.Invoke2(lowered, one, two) }, Int(13)},
		{"keyword lookup", func() (Value, error) { return ec.Invoke1(Keyword("k"), m) }, one},
		{"keyword lookup with default", func() (Value, error) { return ec.Invoke2(Keyword("z"), m, two) }, two},
	}
	for _, c := range callees {
		got, err := c.call()
		if err != nil || got != c.want {
			t.Fatalf("%s: got %v, %v; want %v", c.name, got, err, c.want)
		}
		if n := testing.AllocsPerRun(200, func() { _, _ = c.call() }); n != 0 {
			t.Errorf("%s: %v allocations per call, want 0", c.name, n)
		}
	}

	// Wrong arity through the typed entry: the same error as the slice path.
	pairs := []struct {
		name  string
		typed func() (Value, error)
		slice func() (Value, error)
	}{
		{"fast shape", func() (Value, error) { return ec.Invoke1(fast2.(Fn), one) },
			func() (Value, error) { return ec.Invoke(fast2.(Fn), []Value{one}) }},
		{"typed ctx wrapper", func() (Value, error) { return ec.Invoke1(typed, one) },
			func() (Value, error) { return ec.Invoke(typed, []Value{one}) }},
		{"lowered closure struct", func() (Value, error) { return ec.Invoke3(lowered, one, two, one) },
			func() (Value, error) { return ec.Invoke(lowered, []Value{one, two, one}) }},
		{"keyword", func() (Value, error) { return ec.Invoke3(Keyword("k"), m, one, two) },
			func() (Value, error) { return ec.Invoke(Keyword("k"), []Value{m, one, two}) }},
	}
	for _, p := range pairs {
		_, typedErr := p.typed()
		_, sliceErr := p.slice()
		if typedErr == nil || sliceErr == nil || typedErr.Error() != sliceErr.Error() {
			t.Errorf("%s: arity errors differ: typed %v, slice %v", p.name, typedErr, sliceErr)
		}
	}

	// A native panic is recovered into an error on the typed path as on the
	// slice path.
	panicky, _ := NativeFnType.Box(func(a Value) Value { panic("boom") })
	_, typedErr := ec.Invoke1(panicky.(Fn), one)
	_, sliceErr := ec.Invoke(panicky.(Fn), []Value{one})
	if typedErr == nil || sliceErr == nil || typedErr.Error() != sliceErr.Error() {
		t.Errorf("panic recovery differs: typed %v, slice %v", typedErr, sliceErr)
	}

	// The zero-argument slice literal lowered code keeps using does not
	// allocate, so no Invoke0 is needed.
	fast0, _ := NativeFnType.Box(func() Value { return one })
	if n := testing.AllocsPerRun(200, func() { _, _ = ec.Invoke(fast0.(Fn), []Value{}) }); n != 0 {
		t.Errorf("zero-arg slice literal: %v allocations, want 0", n)
	}
}
