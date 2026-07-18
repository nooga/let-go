/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package vm

import "testing"

// BenchmarkNamespaceLookupStatsDisabled measures Namespace.Lookup with
// lookup stats off (the production case): the noteLookup gate is the only
// stats cost on this path.
func BenchmarkNamespaceLookupStatsDisabled(b *testing.B) {
	ns := NewNamespace("bench-ns")
	ns.Def("x", Int(42))
	sym := Symbol("x")
	SetLookupStatsEnabled(false)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if v := ns.Lookup(sym); v == nil {
			b.Fatal("lookup miss")
		}
	}
}

// BenchmarkNamespaceLookupStatsDisabledParallel measures the contended
// case: with stats off, concurrent lookups previously serialized on
// lookupStatsMu; after the atomic gate they contend only on the
// namespace's own read locks.
func BenchmarkNamespaceLookupStatsDisabledParallel(b *testing.B) {
	ns := NewNamespace("bench-ns")
	ns.Def("x", Int(42))
	sym := Symbol("x")
	SetLookupStatsEnabled(false)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if v := ns.Lookup(sym); v == nil {
				b.Fatal("lookup miss")
			}
		}
	})
}

// BenchmarkNamespaceLookupStatsEnabled keeps the enabled path honest: the
// gate must not make the measuring configuration slower than the mutex it
// replaces on the fast path.
func BenchmarkNamespaceLookupStatsEnabled(b *testing.B) {
	ns := NewNamespace("bench-ns")
	ns.Def("x", Int(42))
	sym := Symbol("x")
	SetLookupStatsEnabled(true)
	defer SetLookupStatsEnabled(false)
	defer ResetLookupStats()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if v := ns.Lookup(sym); v == nil {
			b.Fatal("lookup miss")
		}
	}
}
