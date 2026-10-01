/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"strings"
	"testing"
)

// adderClosure is the shape lowered Go emits for a fn literal: the NativeFn
// embedded beside the captured values, with the body as a method.
type adderClosure struct {
	NativeFn
	y Value
}

func (c *adderClosure) CallNative(args []Value) (Value, error) {
	if len(args) != 1 {
		return NIL, WrongArgCount(&c.NativeFn, len(args), 1)
	}
	return Int(int(args[0].(Int)) + int(c.y.(Int))), nil
}

func newAdderClosure(y Value) Value {
	c := &adderClosure{y: y}
	InitNativeFn(&c.NativeFn, 1, false, c)
	return &c.NativeFn
}

// A struct-embedded closure is one allocation, calls its body with the
// captures, and stays a *NativeFn for everything that inspects one.
func TestInitNativeFnEmbeddedClosure(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { _ = newAdderClosure(Int(2)) }); n != 1 {
		t.Errorf("creating the closure: %v allocations, want 1", n)
	}
	v := newAdderClosure(Int(2))
	fn, ok := v.(*NativeFn)
	if !ok {
		t.Fatalf("closure is %T, want *NativeFn", v)
	}
	if fn.Arity() != 1 || fn.IsVariadic() {
		t.Errorf("arity %d variadic %v, want 1 false", fn.Arity(), fn.IsVariadic())
	}
	args := []Value{Int(40)}
	if n := testing.AllocsPerRun(100, func() { _, _ = fn.Invoke(args) }); n != 0 {
		t.Errorf("calling the closure: %v allocations, want 0", n)
	}
	got, err := fn.Invoke(args)
	if err != nil || got != Int(42) {
		t.Fatalf("Invoke = %v, %v; want 42", got, err)
	}

	// The arity error is the bytecode VM's, naming the closure itself.
	_, bodyErr := fn.Invoke(nil)
	want := "function " + fn.String() + " expected 1 args, got 0"
	if bodyErr == nil || !strings.Contains(bodyErr.Error(), want) {
		t.Errorf("arity error = %v, want %q", bodyErr, want)
	}

	// SetName and printing work on the embedded NativeFn.
	fn.SetName("adder")
	if s := fn.String(); !strings.HasPrefix(s, "<native-fn adder ") {
		t.Errorf("String() = %q, want a named native-fn", s)
	}

	// with-meta copies the NativeFn; the copy still calls the same body and
	// reads the same captures.
	m := Keyword("k")
	withMeta := fn.WithMeta(m).(*NativeFn)
	if withMeta == fn {
		t.Fatal("WithMeta returned the same object")
	}
	if withMeta.Meta() != m {
		t.Errorf("Meta() = %v, want %v", withMeta.Meta(), m)
	}
	if fn.Meta() != NIL {
		t.Errorf("the original gained meta: %v", fn.Meta())
	}
	got, err = withMeta.Invoke(args)
	if err != nil || got != Int(42) {
		t.Fatalf("copy Invoke = %v, %v; want 42", got, err)
	}
	if s := withMeta.String(); !strings.HasPrefix(s, "<native-fn adder ") {
		t.Errorf("copy String() = %q, want the same name", s)
	}
}
