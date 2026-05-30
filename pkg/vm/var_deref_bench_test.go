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

// Root-only (the common case: fn vars, config) — no dynamic binding.
func BenchmarkVarDerefRoot(b *testing.B) {
	v := newRootVar()
	for i := 0; i < b.N; i++ {
		_ = v.Deref()
	}
}

func BenchmarkVarDerefRootParallel(b *testing.B) {
	v := newRootVar()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = v.Deref()
		}
	})
}

// With an active dynamic binding.
func BenchmarkVarDerefBound(b *testing.B) {
	v := newBoundVar()
	for i := 0; i < b.N; i++ {
		_ = v.Deref()
	}
}

func BenchmarkVarDerefBoundParallel(b *testing.B) {
	v := newBoundVar()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = v.Deref()
		}
	})
}

// Distinct vars per worker, all dereffed concurrently — proves the old
// contention was the GLOBAL bindingsMu, not per-var.
func BenchmarkVarDerefDistinctParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		v := newRootVar()
		for pb.Next() {
			_ = v.Deref()
		}
	})
}
