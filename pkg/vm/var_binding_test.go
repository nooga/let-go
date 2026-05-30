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

// TestVarBindingStackSemantics guards the atomic curr-pointer that backs
// the lock-free Deref: nested push/pop must shadow and restore correctly,
// and Deref must reflect the current top (or root when empty).
func TestVarBindingStackSemantics(t *testing.T) {
	v := NewVar(nil, "test", "*x*")
	v.SetRoot(Int(0))

	if got := v.Deref(); got != Int(0) {
		t.Fatalf("root deref: want 0, got %v", got)
	}
	if got := v.Root(); got != Int(0) {
		t.Fatalf("Root(): want 0, got %v", got)
	}

	v.PushBinding(Int(1))
	if got := v.Deref(); got != Int(1) {
		t.Fatalf("after push 1: want 1, got %v", got)
	}
	// Root() must still see the root, not the binding.
	if got := v.Root(); got != Int(0) {
		t.Fatalf("Root() under binding: want 0, got %v", got)
	}

	v.PushBinding(Int(2))
	if got := v.Deref(); got != Int(2) {
		t.Fatalf("after push 2: want 2, got %v", got)
	}

	v.PopBinding()
	if got := v.Deref(); got != Int(1) {
		t.Fatalf("after pop: want 1 (restored), got %v", got)
	}

	v.PopBinding()
	if got := v.Deref(); got != Int(0) {
		t.Fatalf("after pop to empty: want root 0, got %v", got)
	}
}

// TestVarConcurrentDerefDuringBinding stresses the lock-free read path
// against concurrent push/pop. Under -race this proves Deref never races
// the binding mutation, and every observed value is a legitimate one
// (root or a pushed value) — never torn.
func TestVarConcurrentDerefDuringBinding(t *testing.T) {
	v := NewVar(nil, "test", "*y*")
	v.SetRoot(Int(-1))

	var wg sync.WaitGroup
	stop := make(chan struct{})

	// Readers: spin on Deref, asserting only legitimate values.
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-stop:
					return
				default:
				}
				got := v.Deref()
				n, ok := got.(Int)
				if !ok || (n != Int(-1) && n != Int(7)) {
					t.Errorf("torn/illegal deref value: %v", got)
					return
				}
			}
		}()
	}

	// Writer: push/pop a single value many times.
	for i := 0; i < 50000; i++ {
		v.PushBinding(Int(7))
		_ = v.Deref()
		v.PopBinding()
	}
	close(stop)
	wg.Wait()

	if got := v.Deref(); got != Int(-1) {
		t.Fatalf("after all pops: want root -1, got %v", got)
	}
}

// TestVarRunWithBindingsKeepsCurrInSync verifies the snapshot/restore
// path keeps the atomic curr pointer consistent (future/go inherit the
// dynamic environment through this).
func TestVarRunWithBindingsKeepsCurrInSync(t *testing.T) {
	v := NewVar(nil, "test", "*z*")
	v.SetRoot(Int(0))
	v.PushBinding(Int(5))
	defer v.PopBinding()

	snap := SnapshotBindings()

	// Run a fn under a different (empty-for-this-var) environment, then
	// confirm the binding is restored and visible via the lock-free read.
	_, _ = RunWithBindings(BindingSnapshot{}, func() (Value, error) {
		if got := v.Deref(); got != Int(0) {
			t.Errorf("inside empty snapshot: want root 0, got %v", got)
		}
		return NIL, nil
	})
	if got := v.Deref(); got != Int(5) {
		t.Fatalf("after RunWithBindings restore: want 5, got %v", got)
	}

	// And re-running with the captured snapshot re-establishes the binding.
	_, _ = RunWithBindings(snap, func() (Value, error) {
		if got := v.Deref(); got != Int(5) {
			t.Errorf("inside captured snapshot: want 5, got %v", got)
		}
		return NIL, nil
	})
}
