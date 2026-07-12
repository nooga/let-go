package vm

import "testing"

func TestRootBindDepthGate(t *testing.T) {
	v := NewVar(nil, "t", "x")
	v.SetRoot(Int(1))

	// unbound (also covers non-dynamic): reads root, gate skips the stack
	if got := v.Deref(); got != Int(1) {
		t.Fatalf("unbound deref = %v, want 1", got)
	}
	if d := v.rootBindDepth.Load(); d != 0 {
		t.Fatalf("rootBindDepth unbound = %d, want 0", d)
	}

	// bound: deref sees the binding
	v.PushBinding(Int(2))
	if d := v.rootBindDepth.Load(); d != 1 {
		t.Fatalf("rootBindDepth bound = %d, want 1", d)
	}
	if got := v.Deref(); got != Int(2) {
		t.Fatalf("bound deref = %v, want 2", got)
	}

	// nested shadow
	v.PushBinding(Int(3))
	if got := v.Deref(); got != Int(3) {
		t.Fatalf("nested deref = %v, want 3", got)
	}
	v.PopBinding()
	if got := v.Deref(); got != Int(2) {
		t.Fatalf("after inner pop deref = %v, want 2 (restored)", got)
	}

	// fully unbound: back to root, counter zero (the PreviouslyBound case)
	v.PopBinding()
	if d := v.rootBindDepth.Load(); d != 0 {
		t.Fatalf("rootBindDepth after unbind = %d, want 0", d)
	}
	if got := v.Deref(); got != Int(1) {
		t.Fatalf("previously-bound deref = %v, want 1 (root)", got)
	}
}

func TestRunWithBindingsCounterConsistent(t *testing.T) {
	v := NewVar(nil, "t", "y")
	v.SetRoot(Int(10))
	v.SetDynamic()

	snap := BindingSnapshot{v: {Int(99)}}
	_, _ = RunWithBindings(snap, func() (Value, error) {
		if got := v.Deref(); got != Int(99) {
			t.Errorf("inside RunWithBindings deref = %v, want 99", got)
		}
		if d := v.rootBindDepth.Load(); d != 1 {
			t.Errorf("inside RunWithBindings depth = %d, want 1", d)
		}
		return NIL, nil
	})
	if got := v.Deref(); got != Int(10) {
		t.Errorf("after RunWithBindings deref = %v, want 10 (root)", got)
	}
	if d := v.rootBindDepth.Load(); d != 0 {
		t.Errorf("after RunWithBindings depth = %d, want 0", d)
	}
}
