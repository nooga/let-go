/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import "testing"

// A lazy SourceMap defers decoding its entries until the first Lookup/Entries —
// bundle load (decode of every chunk's source map) is the dominant startup heap
// churn, yet source maps are only read on an error/stack-trace. The decode
// closure must run at most once, on first access, and produce identical results
// to an eagerly-built map.
func TestLazySourceMapDefersUntilFirstAccess(t *testing.T) {
	decodes := 0
	sm := NewLazySourceMap(func() []SourceMapEntry {
		decodes++
		return []SourceMapEntry{
			{StartIP: 0, Info: SourceInfo{File: "a.lg", Line: 2}},
			{StartIP: 5, Info: SourceInfo{File: "a.lg", Line: 7}},
		}
	})

	if decodes != 0 {
		t.Fatalf("lazy decode ran at construction (decodes=%d)", decodes)
	}

	// First Lookup triggers exactly one decode and returns the right entry.
	got := sm.Lookup(6)
	if decodes != 1 {
		t.Fatalf("expected exactly 1 decode after first Lookup, got %d", decodes)
	}
	if got == nil || got.Line != 7 || got.File != "a.lg" {
		t.Fatalf("Lookup(6) = %v, want {a.lg line 7}", got)
	}

	// Further accesses must reuse the materialized entries (no re-decode).
	if got := sm.Lookup(0); got == nil || got.Line != 2 {
		t.Fatalf("Lookup(0) = %v, want line 2", got)
	}
	if es := sm.Entries(); len(es) != 2 {
		t.Fatalf("Entries() len = %d, want 2", len(es))
	}
	if decodes != 1 {
		t.Fatalf("lazy decode re-ran (decodes=%d, want 1)", decodes)
	}
}

// A lazy SourceMap must be Lookup-equivalent to the eager map built from the
// same entries — same result for every IP, including the pre-first-entry gap.
func TestLazySourceMapEquivalentToEager(t *testing.T) {
	entries := []SourceMapEntry{
		{StartIP: 3, Info: SourceInfo{File: "f.lg", Line: 1, Column: 4}},
		{StartIP: 10, Info: SourceInfo{File: "f.lg", Line: 9, Column: 2}},
	}
	eager := NewSourceMap()
	for _, e := range entries {
		eager.Add(e.StartIP, e.Info)
	}
	lazy := NewLazySourceMap(func() []SourceMapEntry { return entries })

	for ip := 0; ip < 15; ip++ {
		a, b := eager.Lookup(ip), lazy.Lookup(ip)
		if (a == nil) != (b == nil) {
			t.Fatalf("ip=%d: eager=%v lazy=%v (nil mismatch)", ip, a, b)
		}
		if a != nil && (a.Line != b.Line || a.Column != b.Column || a.File != b.File) {
			t.Fatalf("ip=%d: eager=%v lazy=%v", ip, a, b)
		}
	}
}
