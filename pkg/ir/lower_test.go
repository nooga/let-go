/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package ir

import (
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

// TestLower_CrossBlockMultiUse_Errors: construct an IR by hand where a
// value defined in block A is used as a Refs operand by a node in block B
// (not just as a BlockArg). The spike's Lower cannot express this; assert
// it errors cleanly with a message mentioning "cross-block use".
func TestLower_CrossBlockMultiUse_Errors(t *testing.T) {
	consts := vm.NewConsts()
	f := NewFunction("xb", 1, false, consts)
	b0 := f.Entry
	b1 := f.AddBlock()

	// v0 = LoadArg 0 — explicit fn-arg reference in the entry block.
	v0 := f.AddNode(Node{Op: OpLoadArg, Aux: 0, Block: b0})
	f.AppendToBlock(b0, v0)

	// Use 1 (in-block, Refs): Inc inside b0
	incID := f.AddNode(Node{Op: OpInc, Refs: []NodeID{v0}, Block: b0})
	f.AppendToBlock(b0, incID)

	// Pop the inc result so it doesn't dangle.
	popID := f.AddNode(Node{Op: OpPop, Refs: []NodeID{incID}, Block: b0})
	f.AppendToBlock(b0, popID)

	// b0 terminator: unconditional branch to b1 (no args — v0 will be
	// referenced directly by b1's Return, which is the cross-block sin).
	bt := &BranchTarget{Target: b1}
	termID := f.AddNode(Node{Op: OpBranch, Aux: bt, Block: b0})
	f.SetTerminator(b0, termID)
	f.AddPred(b1, b0)

	// b1: returns v0 directly — a cross-block Refs use. This is what we
	// expect Lower to reject.
	retID := f.AddNode(Node{Op: OpReturn, Refs: []NodeID{v0}, Block: b1})
	f.SetTerminator(b1, retID)

	_, err := Lower(f)
	if err == nil {
		t.Fatal("expected Lower to error on cross-block multi-use, got nil")
	}
	if !strings.Contains(err.Error(), "cross-block use") {
		t.Errorf("expected error mentioning 'cross-block use', got: %v", err)
	}
}

// TestLower_CrossBlockSingleUse_Errors: a value with exactly one use,
// and that use is a cross-block direct Refs (not a branch-arg). The
// spike's lowering can't express this; assert Lower errors.
//
// Earlier versions of the check gated on len(uses) >= 2 and would have
// silently miscompiled this case. This test prevents regression.
func TestLower_CrossBlockSingleUse_Errors(t *testing.T) {
	consts := vm.NewConsts()
	f := NewFunction("xb1", 1, false, consts)
	b0 := f.Entry
	b1 := f.AddBlock()

	v0 := f.AddNode(Node{Op: OpLoadArg, Aux: 0, Block: b0})
	f.AppendToBlock(b0, v0)

	// b0 terminator: Branch to b1 with no args.
	bt := &BranchTarget{Target: b1}
	termID := f.AddNode(Node{Op: OpBranch, Aux: bt, Block: b0})
	f.SetTerminator(b0, termID)
	f.AddPred(b1, b0)

	// b1 returns v0 — v0's only use, and it's cross-block Refs.
	retID := f.AddNode(Node{Op: OpReturn, Refs: []NodeID{v0}, Block: b1})
	f.SetTerminator(b1, retID)

	_, err := Lower(f)
	if err == nil {
		t.Fatal("expected Lower to error on cross-block single-use, got nil")
	}
	if !strings.Contains(err.Error(), "cross-block use") {
		t.Errorf("expected error mentioning 'cross-block use', got: %v", err)
	}
}

// TestLower_LegitimateCrossBlockViaBranchArg: a value defined in b0,
// passed to b1 as a branch-arg, used in b1 as a BlockArg. This is the
// legitimate SSA cross-block pattern; Lower must NOT reject it.
func TestLower_LegitimateCrossBlockViaBranchArg(t *testing.T) {
	consts := vm.NewConsts()
	f := NewFunction("xbok", 1, false, consts)
	b0 := f.Entry
	b1 := f.AddBlock()

	v0 := f.AddNode(Node{Op: OpLoadArg, Aux: 0, Block: b0})
	f.AppendToBlock(b0, v0)

	// b0 terminator: Branch to b1, passing v0 as an arg.
	bt := &BranchTarget{Target: b1, Args: []NodeID{v0}}
	termID := f.AddNode(Node{Op: OpBranch, Aux: bt, Block: b0})
	f.SetTerminator(b0, termID)
	f.AddPred(b1, b0)

	// b1 declares one BlockArg and returns it.
	argID := f.AddNode(Node{Op: OpBlockArg, Aux: 0, Block: b1})
	f.Blocks[b1].Params = []NodeID{argID}
	retID := f.AddNode(Node{Op: OpReturn, Refs: []NodeID{argID}, Block: b1})
	f.SetTerminator(b1, retID)

	// We don't necessarily expect Lower to succeed (the spike has other
	// quirks), but it must NOT reject with the cross-block error.
	_, err := Lower(f)
	if err != nil && strings.Contains(err.Error(), "cross-block use") {
		t.Errorf("Lower incorrectly rejected legitimate branch-arg cross-block use: %v", err)
	}
}

// TestLower_InBlockMultiUse_NotRejectedByCrossBlockCheck: a value with
// multiple uses all in the same block. The cross-block check must not
// fire here; Task 3 will handle the actual multi-use lowering.
func TestLower_InBlockMultiUse_NotRejectedByCrossBlockCheck(t *testing.T) {
	consts := vm.NewConsts()
	f := NewFunction("inblock", 1, false, consts)
	b0 := f.Entry
	v0 := f.AddNode(Node{Op: OpLoadArg, Aux: 0, Block: b0})
	f.AppendToBlock(b0, v0)

	// Two in-block uses of v0.
	inc1 := f.AddNode(Node{Op: OpInc, Refs: []NodeID{v0}, Block: b0})
	f.AppendToBlock(b0, inc1)
	inc2 := f.AddNode(Node{Op: OpInc, Refs: []NodeID{v0}, Block: b0})
	f.AppendToBlock(b0, inc2)

	// Return one of them so the function is well-formed.
	retID := f.AddNode(Node{Op: OpReturn, Refs: []NodeID{inc2}, Block: b0})
	f.SetTerminator(b0, retID)

	// Lower may fail later (spike doesn't yet handle multi-use), but
	// it must NOT fail with the cross-block error.
	_, err := Lower(f)
	if err != nil && strings.Contains(err.Error(), "cross-block use") {
		t.Errorf("Lower incorrectly rejected in-block multi-use as cross-block: %v", err)
	}
}

// TestLower_MultiUseAdd_DUPNTH: hand-built IR where an Add result is
// used twice. Lower must emit DUP_NTH to share it (Add isn't cheap to
// re-materialize). Verifies the DUP_NTH branch of materialize.
func TestLower_MultiUseAdd_DUPNTH(t *testing.T) {
	consts := vm.NewConsts()
	f := NewFunction("dupadd", 1, false, consts)
	b0 := f.Entry
	// Use OpLoadArg (not the entry BlockArg param) so the lowerer emits
	// OP_LOAD_ARG bytecode. BlockArg params are treated as "virtually on
	// stack" by the lowerer but the VM starts with sp=0, so the entry
	// block's args must be loaded explicitly via OpLoadArg.
	arg := f.AddNode(Node{Op: OpLoadArg, Aux: 0, Block: b0})
	f.AppendToBlock(b0, arg)
	c1 := f.AddNode(Node{Op: OpConst, Aux: vm.Int(1), Block: b0})
	f.AppendToBlock(b0, c1)
	// shared = arg + 1
	shared := f.AddNode(Node{Op: OpAdd, Refs: []NodeID{arg, c1}, Block: b0})
	f.AppendToBlock(b0, shared)
	// result = shared + shared
	result := f.AddNode(Node{Op: OpAdd, Refs: []NodeID{shared, shared}, Block: b0})
	f.AppendToBlock(b0, result)
	ret := f.AddNode(Node{Op: OpReturn, Refs: []NodeID{result}, Block: b0})
	f.SetTerminator(b0, ret)

	chunk, err := Lower(f)
	if err != nil {
		t.Fatalf("Lower: %v", err)
	}
	// Run with arg=3: (3+1) + (3+1) = 8
	frame := vm.NewFrame(chunk, []vm.Value{vm.Int(3)})
	r, err := frame.Run()
	vm.ReleaseFrame(frame)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if r.String() != "8" {
		t.Errorf("expected 8, got %s", r)
	}

	// Sanity: assert DUP_NTH appears in the emitted bytecode.
	code := chunk.Code()
	foundDup := false
	for i := 0; i < len(code); {
		op := code[i] & 0xff
		if op == vm.OP_DUP_NTH {
			foundDup = true
			break
		}
		stride := instStride(op)
		if stride <= 0 {
			stride = 1
		}
		i += stride
	}
	if !foundDup {
		t.Error("expected emitted bytecode to contain DUP_NTH for multi-use Add; not found")
	}
}

// TestLower_MultiUseConst: a Const referenced twice within a block.
// After Lower + Run, the result must be correct.
func TestLower_MultiUseConst(t *testing.T) {
	consts := vm.NewConsts()
	f := NewFunction("twoones", 0, false, consts)
	// %0 = Const 1; result = Add %0 %0; Return
	b0 := f.Entry
	c1 := f.AddNode(Node{Op: OpConst, Aux: vm.Int(1), Block: b0})
	f.AppendToBlock(b0, c1)
	add := f.AddNode(Node{Op: OpAdd, Refs: []NodeID{c1, c1}, Block: b0})
	f.AppendToBlock(b0, add)
	ret := f.AddNode(Node{Op: OpReturn, Refs: []NodeID{add}, Block: b0})
	f.SetTerminator(b0, ret)

	chunk, err := Lower(f)
	if err != nil {
		t.Fatalf("Lower: %v", err)
	}
	frame := vm.NewFrame(chunk, nil)
	result, err := frame.Run()
	vm.ReleaseFrame(frame)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if result.String() != "2" {
		t.Errorf("expected 2, got %s", result)
	}
}
