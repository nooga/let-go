/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import "testing"

// Decoding the core bundle calls AddSourceInfo / AddLocalVar once per
// instruction. Without a pre-sized backing array these grow by repeated
// append, and that reallocation churn was ~47% of startup heap allocation
// (SourceMap.Add alone was 39%). Reserve lets the decoder size the slice
// once from the known entry count. These tests pin the no-realloc property.

func TestSourceMapReserveAvoidsRealloc(t *testing.T) {
	const n = 256
	sm := NewSourceMap()
	sm.Reserve(n)
	sm.Add(0, SourceInfo{File: "a.lg"})
	first := &sm.entries[0]
	for i := 1; i < n; i++ {
		sm.Add(i, SourceInfo{File: "a.lg", Line: i})
	}
	if &sm.entries[0] != first {
		t.Fatalf("SourceMap reallocated despite Reserve(%d)", n)
	}
	if len(sm.entries) != n {
		t.Fatalf("expected %d entries, got %d", n, len(sm.entries))
	}
}

func TestCodeChunkReserveSourceMapAvoidsRealloc(t *testing.T) {
	const n = 256
	c := NewCodeChunk(NewConsts())
	c.ReserveSourceMap(n)
	c.AddSourceInfo(SourceInfo{File: "a.lg"})
	first := &c.sourceMap.entries[0]
	for i := 1; i < n; i++ {
		c.AddSourceInfo(SourceInfo{File: "a.lg", Line: i})
	}
	if &c.sourceMap.entries[0] != first {
		t.Fatalf("CodeChunk sourceMap reallocated despite ReserveSourceMap(%d)", n)
	}
}

func TestCodeChunkReserveLocalVarsAvoidsRealloc(t *testing.T) {
	const n = 128
	c := NewCodeChunk(NewConsts())
	c.ReserveLocalVars(n)
	c.AddLocalVar(0, "x0")
	first := &c.localVars[0]
	for i := 1; i < n; i++ {
		c.AddLocalVar(i, "x")
	}
	if &c.localVars[0] != first {
		t.Fatalf("localVars reallocated despite ReserveLocalVars(%d)", n)
	}
	if len(c.localVars) != n {
		t.Fatalf("expected %d localVars, got %d", n, len(c.localVars))
	}
}
