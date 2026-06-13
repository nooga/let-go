/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"sync"
	"testing"
)

// TestExecContextSetBinding covers (set! *v* val) semantics at the context
// level: it mutates the top binding frame of THIS context only (thread-local),
// is visible to a subsequent deref, does not leak to a sibling context or the
// root, and reports false (so callers can fall back to the root) when no
// binding is active.
func TestExecContextSetBinding(t *testing.T) {
	v := NewVar(nil, "test", "m")
	v.SetRoot(MakeInt(0))

	c1 := RootExecContext.Child()
	c2 := RootExecContext.Child()
	c1.PushBinding(v, MakeInt(1))
	c2.PushBinding(v, MakeInt(2))

	// set! in c1 mutates c1's frame, visible to c1's deref...
	if ok := c1.setBinding(v, MakeInt(11)); !ok {
		t.Fatal("setBinding should report true when a binding is active")
	}
	if got := c1.deref(v).Unbox(); got != 11 {
		t.Fatalf("c1 deref = %v, want 11 (set! must be visible within the binding)", got)
	}
	// ...but not c2's frame, nor the root.
	if got := c2.deref(v).Unbox(); got != 2 {
		t.Fatalf("c2 deref = %v, want 2 (sibling must not see c1's set!)", got)
	}
	if got := RootExecContext.deref(v).Unbox(); got != 0 {
		t.Fatalf("root deref = %v, want 0 (set! must not leak to root)", got)
	}

	// No active binding for w: setBinding reports false (caller would root-set).
	w := NewVar(nil, "test", "n")
	w.SetRoot(MakeInt(7))
	if ok := RootExecContext.setBinding(w, MakeInt(9)); ok {
		t.Fatal("setBinding should report false when no binding is active")
	}
	if got := RootExecContext.deref(w).Unbox(); got != 7 {
		t.Fatalf("root deref = %v, want 7 (setBinding must not create a binding)", got)
	}
}

// TestExecContextChildIsolation proves two child contexts hold INDEPENDENT
// dynamic bindings of the same var (distinct values), and that neither child's
// push leaks back to the root. Distinct values are what make this a real
// isolation test: a leak would surface as one child reading the other's value.
func TestExecContextChildIsolation(t *testing.T) {
	v := NewVar(nil, "test", "x")
	v.SetRoot(MakeInt(0))

	c1 := RootExecContext.Child()
	c2 := RootExecContext.Child()
	c1.PushBinding(v, MakeInt(1))
	c2.PushBinding(v, MakeInt(2))

	if got := c1.Deref(v).Unbox(); got != 1 {
		t.Fatalf("c1 sees %v, want 1 (c2's push must not be visible to c1)", got)
	}
	if got := c2.Deref(v).Unbox(); got != 2 {
		t.Fatalf("c2 sees %v, want 2 (c1's push must not be visible to c2)", got)
	}
	if got := RootExecContext.Deref(v).Unbox(); got != 0 {
		t.Fatalf("root sees %v, want 0 (child pushes must not leak to root)", got)
	}
}

// TestExecContextChildInheritsParentBinding proves the inheritance half of the
// contract: a child created by ec.Child() carries the parent's bindings that
// were active AT CREATION (snapshot semantics), a later parent push does NOT
// retroactively reach the already-created child, and a child override does not
// leak back up to the parent.
func TestExecContextChildInheritsParentBinding(t *testing.T) {
	v := NewVar(nil, "test", "w")
	v.SetRoot(MakeInt(0))

	parent := RootExecContext.Child()
	parent.PushBinding(v, MakeInt(9))

	child := parent.Child() // created AFTER parent bound v
	if got := child.Deref(v).Unbox(); got != 9 {
		t.Fatalf("child sees %v, want 9 (must inherit the parent's active binding)", got)
	}

	// A later parent push must not retroactively appear in the child: Child()
	// is a snapshot, not a live view.
	parent.PushBinding(v, MakeInt(10))
	if got := child.Deref(v).Unbox(); got != 9 {
		t.Fatalf("child sees %v, want 9 (snapshot must not track later parent pushes)", got)
	}

	// A child override must not leak up to the parent.
	child.PushBinding(v, MakeInt(11))
	if got := parent.Deref(v).Unbox(); got != 10 {
		t.Fatalf("parent sees %v, want 10 (child push must not leak to parent)", got)
	}
}

// TestExecContextConcurrentIsolation runs many goroutines that each spawn a
// child context and bind the same var to a goroutine-UNIQUE value, then read it
// back many times. Because every worker binds a different value, any
// shared-state interleaving shows up as a worker reading a value that isn't its
// own. Run with -race to additionally catch unsynchronised access.
func TestExecContextConcurrentIsolation(t *testing.T) {
	v := NewVar(nil, "test", "y")
	v.SetRoot(MakeInt(-1))

	const workers = 64
	const reads = 200
	var wg sync.WaitGroup
	errs := make(chan string, workers)
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			child := RootExecContext.Child()
			child.PushBinding(v, MakeInt(id))
			for r := 0; r < reads; r++ {
				if got := child.Deref(v).Unbox(); got != id {
					errs <- "cross-talk: a worker observed another worker's binding"
					return
				}
			}
		}(i)
	}
	wg.Wait()
	close(errs)
	if msg, bad := <-errs; bad {
		t.Fatal(msg)
	}
	if got := RootExecContext.Deref(v).Unbox(); got != -1 {
		t.Fatalf("root sees %v, want -1 (worker pushes must not leak to root)", got)
	}
}

// TestExecContextSnapshotSeedsChild proves the bound-fn* / spawn primitive:
// a context built from a captured snapshot re-establishes exactly those
// bindings, independently per construction, with no leak back to root.
func TestExecContextSnapshotSeedsChild(t *testing.T) {
	v := NewVar(nil, "test", "z")
	v.SetRoot(MakeInt(0))

	snap := func() BindingSnapshot {
		RootExecContext.PushBinding(v, MakeInt(42))
		defer RootExecContext.PopBinding(v)
		return RootExecContext.BindingSnapshot()
	}()

	// Each fresh context from the snapshot sees the captured binding...
	if got := NewExecContextFrom(snap).Deref(v).Unbox(); got != 42 {
		t.Fatalf("snapshot child sees %v, want 42", got)
	}
	// ...and a push into one such context does not bleed into another built
	// from the same snapshot.
	a := NewExecContextFrom(snap)
	b := NewExecContextFrom(snap)
	a.PushBinding(v, MakeInt(7))
	if got := a.Deref(v).Unbox(); got != 7 {
		t.Fatalf("a sees %v, want 7", got)
	}
	if got := b.Deref(v).Unbox(); got != 42 {
		t.Fatalf("b sees %v, want 42 (a's push must not leak through the shared snapshot)", got)
	}
	// Root never saw 42 or 7 after the bracket unwound.
	if got := RootExecContext.Deref(v).Unbox(); got != 0 {
		t.Fatalf("root sees %v, want 0", got)
	}
}
