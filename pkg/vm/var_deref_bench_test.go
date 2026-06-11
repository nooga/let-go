/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package vm

import "testing"

// These benchmarks quantify Var.Deref — the hottest var operation. Before
// this change every Deref took the global bindingsMu (even a var with NO
// dynamic binding, just to check the stack was empty), so all var reads
// across all goroutines serialized on one mutex. The Parallel variants
// expose that contention; after making root/curr atomic, Deref is a
// couple of atomic loads and scales.

// derefSink defeats dead-code elimination: without a package-level sink the
// compiler may elide the Deref call entirely, producing impossible sub-ns
// numbers that say nothing about the real cost.
var derefSink Value

func newRootVar() *Var {
	v := NewVar(nil, "bench", "x")
	v.SetRoot(Int(42))
	return v
}

func newBoundVar() *Var {
	v := newRootVar()
	v.PushBinding(Int(7))
	return v
}

// newPreviouslyBoundVar models the most common read of a declared-dynamic var:
// it was bound at least once (so isDynamic is permanently set, as ^:dynamic and
// any binding both set it), but no binding is active now. Every such deref still
// consults the binding stack — the cost this optimization targets.
func newPreviouslyBoundVar() *Var {
	v := newBoundVar()
	v.PopBinding()
	return v
}

// Root-only (the common case: fn vars, config) — no dynamic binding.
func BenchmarkVarDerefRoot(b *testing.B) {
	v := newRootVar()
	for i := 0; i < b.N; i++ {
		derefSink = v.Deref()
	}
}

func BenchmarkVarDerefRootParallel(b *testing.B) {
	v := newRootVar()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			derefSink = v.Deref()
		}
	})
}

// Declared dynamic, bound once, now unbound — the common steady-state read of
// *out*/*ns* outside any binding. isDynamic stays set, so the deref still has to
// consult the (empty-for-this-var) stack.
func BenchmarkVarDerefPreviouslyBound(b *testing.B) {
	v := newPreviouslyBoundVar()
	for i := 0; i < b.N; i++ {
		derefSink = v.Deref()
	}
}

func BenchmarkVarDerefPreviouslyBoundParallel(b *testing.B) {
	v := newPreviouslyBoundVar()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			derefSink = v.Deref()
		}
	})
}

// With an active dynamic binding.
func BenchmarkVarDerefBound(b *testing.B) {
	v := newBoundVar()
	for i := 0; i < b.N; i++ {
		derefSink = v.Deref()
	}
}

func BenchmarkVarDerefBoundParallel(b *testing.B) {
	v := newBoundVar()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			derefSink = v.Deref()
		}
	})
}

// Distinct vars per worker, all dereffed concurrently — proves the old
// contention was the GLOBAL bindingsMu, not per-var.
func BenchmarkVarDerefDistinctParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		v := newRootVar()
		for pb.Next() {
			derefSink = v.Deref()
		}
	})
}

// BenchmarkBindingPushPop measures one full binding extent (the (binding […])
// write path). Copy-on-write makes reads lock-free at the cost of allocating a
// fresh map per push/pop, so this is the side of the trade that gets more
// expensive — binding establishment is far rarer than the reads inside it.
func BenchmarkBindingPushPop(b *testing.B) {
	v := newRootVar()
	for i := 0; i < b.N; i++ {
		v.PushBinding(Int(7))
		v.PopBinding()
	}
}
