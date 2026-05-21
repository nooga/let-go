/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package ir

import (
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

// TestSourceMap_RoundTrip: compiles a small fn, asserts the original
// bytecode has source info, then builds IR + lowers back to bytecode
// and asserts the lowered chunk has source info for its emitted ops.
//
// Doesn't require IP-for-IP equivalence (optimization may shift offsets)
// — just that *some* SourceInfo survives for the non-trivial emitted
// instructions, with the same File as the original.
func TestSourceMap_RoundTrip(t *testing.T) {
	consts := vm.NewConsts()
	ns := rt.NS(rt.NameCoreNS)
	ctx := compiler.NewCompiler(consts, ns)
	src := `(defn add1 [x] (+ x 1))`
	_, _, err := ctx.CompileMultiple(strings.NewReader(src))
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	v := ns.Lookup(vm.Symbol("add1")).(*vm.Var).Deref()
	fn := v.(*vm.Func)
	origChunk := fn.Chunk()
	origMap := origChunk.GetSourceMap()
	if origMap == nil || len(origMap.Entries()) == 0 {
		t.Skip("compiler did not populate source map for this fixture; nothing to round-trip")
		return
	}
	// Pick a File name from the first original entry — that's what we
	// expect the lowered chunk to carry too.
	origFile := origMap.Entries()[0].Info.File

	irFn, err := Build(origChunk, "add1", 1, false)
	if err != nil {
		t.Fatalf("Build: %v", err)
	}

	// Spot-check: at least one IR Node should have SourceInfo set.
	hasNodeSI := false
	for _, n := range irFn.Nodes {
		if len(n.SourceInfos) > 0 {
			hasNodeSI = true
			break
		}
	}
	if !hasNodeSI {
		t.Error("expected at least one IR Node with SourceInfo after Build, got none")
	}

	loweredChunk, err := Lower(irFn)
	if err != nil {
		t.Fatalf("Lower: %v", err)
	}
	loweredMap := loweredChunk.GetSourceMap()
	if loweredMap == nil || len(loweredMap.Entries()) == 0 {
		t.Error("expected lowered chunk to have a non-empty SourceMap, got none")
		return
	}
	// Confirm File names round-trip.
	if loweredMap.Entries()[0].Info.File != origFile {
		t.Errorf("source file mismatch after round-trip: original %q, lowered %q",
			origFile, loweredMap.Entries()[0].Info.File)
	}
}

// TestSourceMap_HandSetRoundTrip confirms that a hand-constructed IR Node
// with SourceInfo set will have its SourceInfo present in the lowered chunk.
// This test is unconditional and does not depend on the compiler emitting
// source info.
func TestSourceMap_HandSetRoundTrip(t *testing.T) {
	consts := vm.NewConsts()
	f := NewFunction("manual", 0, false, consts)
	b0 := f.Entry
	si := vm.SourceInfo{File: "manual.lg", Line: 4, Column: 12}
	c := f.AddNode(Node{Op: OpConst, Aux: vm.Int(42), Block: b0, SourceInfos: []vm.SourceInfo{si}})
	f.AppendToBlock(b0, c)
	ret := f.AddNode(Node{Op: OpReturn, Refs: []NodeID{c}, Block: b0})
	f.SetTerminator(b0, ret)

	chunk, err := Lower(f)
	if err != nil {
		t.Fatalf("Lower: %v", err)
	}
	sm := chunk.GetSourceMap()
	if sm == nil {
		t.Fatal("expected lowered chunk to have a SourceMap")
	}
	entries := sm.Entries()
	if len(entries) == 0 {
		t.Fatal("expected at least one SourceMap entry")
	}
	if entries[0].Info.File != "manual.lg" || entries[0].Info.Line != 4 || entries[0].Info.Column != 12 {
		t.Errorf("source info did not round-trip: got %+v", entries[0].Info)
	}
}
