/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package passes

import (
	"testing"

	"github.com/nooga/let-go/pkg/ir"
	"github.com/nooga/let-go/pkg/vm"
)

// TestConstFold_PreservesOperandSpans: folding (+ 1 2) into Const 3 should
// keep spans from both operands plus the original Add.
func TestConstFold_PreservesOperandSpans(t *testing.T) {
	consts := vm.NewConsts()
	f := ir.NewFunction("fold", 0, false, consts)
	b0 := f.Entry

	siA := vm.SourceInfo{File: "f.lg", Line: 1, Column: 1}
	siB := vm.SourceInfo{File: "f.lg", Line: 1, Column: 5}
	siAdd := vm.SourceInfo{File: "f.lg", Line: 1, Column: 3}

	a := f.AddNode(ir.Node{Op: ir.OpConst, Aux: vm.Int(1), Block: b0, SourceInfos: []vm.SourceInfo{siA}})
	f.AppendToBlock(b0, a)
	b := f.AddNode(ir.Node{Op: ir.OpConst, Aux: vm.Int(2), Block: b0, SourceInfos: []vm.SourceInfo{siB}})
	f.AppendToBlock(b0, b)
	add := f.AddNode(ir.Node{Op: ir.OpAdd, Refs: []ir.NodeID{a, b}, Block: b0, SourceInfos: []vm.SourceInfo{siAdd}})
	f.AppendToBlock(b0, add)
	ret := f.AddNode(ir.Node{Op: ir.OpReturn, Refs: []ir.NodeID{add}, Block: b0})
	f.SetTerminator(b0, ret)

	if !ConstFold(f) {
		t.Fatal("ConstFold did not change anything")
	}
	folded := f.Node(add)
	if folded.Op != ir.OpConst {
		t.Fatalf("expected folded node to be OpConst, got %s", folded.Op)
	}
	// Must contain all three spans.
	seen := map[vm.SourceInfo]bool{}
	for _, si := range folded.SourceInfos {
		seen[si] = true
	}
	for _, expected := range []vm.SourceInfo{siA, siB, siAdd} {
		if !seen[expected] {
			t.Errorf("folded node missing source span %+v; got %+v", expected, folded.SourceInfos)
		}
	}
}

// TestCSE_MergesSpans: two equivalent computations with different spans
// should have BOTH spans on the surviving node after CSE.
func TestCSE_MergesSpans(t *testing.T) {
	consts := vm.NewConsts()
	f := ir.NewFunction("cse", 1, false, consts)
	b0 := f.Entry

	si1 := vm.SourceInfo{File: "f.lg", Line: 2, Column: 1}
	si2 := vm.SourceInfo{File: "f.lg", Line: 3, Column: 1}

	a := f.AddNode(ir.Node{Op: ir.OpLoadArg, Aux: 0, Block: b0, SourceInfos: []vm.SourceInfo{si1}})
	f.AppendToBlock(b0, a)
	b := f.AddNode(ir.Node{Op: ir.OpLoadArg, Aux: 0, Block: b0, SourceInfos: []vm.SourceInfo{si2}})
	f.AppendToBlock(b0, b)
	// Use both somewhere so they aren't dead.
	add := f.AddNode(ir.Node{Op: ir.OpAdd, Refs: []ir.NodeID{a, b}, Block: b0})
	f.AppendToBlock(b0, add)
	ret := f.AddNode(ir.Node{Op: ir.OpReturn, Refs: []ir.NodeID{add}, Block: b0})
	f.SetTerminator(b0, ret)

	if !CSE(f) {
		t.Fatal("CSE did not change anything (expected two LoadArg 0s to merge)")
	}
	// The survivor is `a`. It must carry both si1 and si2.
	survivor := f.Node(a)
	seen := map[vm.SourceInfo]bool{}
	for _, si := range survivor.SourceInfos {
		seen[si] = true
	}
	if !seen[si1] {
		t.Errorf("CSE survivor missing original span si1; got %+v", survivor.SourceInfos)
	}
	if !seen[si2] {
		t.Errorf("CSE survivor missing merged span si2; got %+v", survivor.SourceInfos)
	}
}
