/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"fmt"
	"testing"
)

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

// With an active dynamic binding. b.Cleanup pops it so the binding does not
// leak onto the shared root stack and contaminate later benchmarks' write costs.
func BenchmarkVarDerefBound(b *testing.B) {
	v := newBoundVar()
	b.Cleanup(v.PopBinding)
	for i := 0; i < b.N; i++ {
		derefSink = v.Deref()
	}
}

func BenchmarkVarDerefBoundParallel(b *testing.B) {
	v := newBoundVar()
	b.Cleanup(v.PopBinding)
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

// bindDepth binds `depth` distinct dynamic vars on the root context (vars[0]
// pushed first sits at the bottom of the stack, vars[depth-1] at the top) and
// pops them all on cleanup. It returns the vars plus a declared-but-unbound var
// for the miss case.
func bindDepth(b *testing.B, depth int) (vars []*Var, miss *Var) {
	vars = make([]*Var, depth)
	for i := range vars {
		vars[i] = newRootVar()
		vars[i].PushBinding(Int(i))
	}
	b.Cleanup(func() {
		for i := depth - 1; i >= 0; i-- {
			vars[i].PopBinding()
		}
	})
	miss = newRootVar()
	miss.SetDynamic() // declared dynamic, never bound → deref walks the whole stack and misses
	return vars, miss
}

// BenchmarkVarDerefDepth measures deref cost as a function of how many unrelated
// bindings are active above the target. This is where the frame-list (O(active
// depth) walk) and the copy-on-write map (O(1) hash) profiles diverge: head hits
// stay cheap for the list, but tail hits and misses pay the walk.
func BenchmarkVarDerefDepth(b *testing.B) {
	for _, depth := range []int{1, 4, 16, 64} {
		vars, miss := bindDepth(b, depth)
		targets := map[string]*Var{
			"head": vars[depth-1], // most-recently bound (top of stack)
			"tail": vars[0],       // first bound (bottom of stack)
			"miss": miss,          // not on the stack — walks all frames
		}
		for _, pos := range []string{"head", "tail", "miss"} {
			t := targets[pos]
			b.Run(fmt.Sprintf("d%d/%s", depth, pos), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					derefSink = t.Deref()
				}
			})
		}
	}
}

// BenchmarkNsWorkload models nnunley's requested scenario: a nest of dynamic
// bindings (think thread-local *ns* across nested namespaces) accessed in three
// read/write ratios. Each iteration does readsPerWrite reads spread across the
// active stack (a mix of head and tail hits) and one write cycle (push a new
// binding, read under it, pop) — so it exercises both the deref walk and the
// binding-establishment path that the three designs trade off differently.
func benchNsWorkload(b *testing.B, depth, readsPerWrite int) {
	vars, _ := bindDepth(b, depth)
	w := newRootVar()
	w.SetDynamic()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for r := 0; r < readsPerWrite; r++ {
			derefSink = vars[r%depth].Deref()
		}
		w.PushBinding(Int(i))
		derefSink = w.Deref()
		w.PopBinding()
	}
}

func BenchmarkNsWorkload(b *testing.B) {
	profiles := []struct {
		name          string
		readsPerWrite int
	}{
		{"read-only", 64},
		{"balanced", 4},
		{"write-heavy", 1},
	}
	// Several nested-extent depths: deeper stacks stress the frame-list read
	// walk while the copy-on-write map pays its whole-map copy on every write.
	for _, depth := range []int{8, 16, 32} {
		for _, p := range profiles {
			b.Run(fmt.Sprintf("d%d/%s", depth, p.name), func(b *testing.B) {
				benchNsWorkload(b, depth, p.readsPerWrite)
			})
		}
	}
}
