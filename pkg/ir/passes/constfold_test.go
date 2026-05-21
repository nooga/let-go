/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package passes

import (
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/ir"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

// TestConstFoldSimple: (defn f [] (+ 2 3)) should fold to Const 5
// after the pass. Before: Const 2; Const 3; Add → Return. After:
// Const 5 → Return (with the Add op overwritten as a Const).
func TestConstFoldSimple(t *testing.T) {
	consts := vm.NewConsts()
	ns := rt.NS(rt.NameCoreNS)
	ctx := compiler.NewCompiler(consts, ns)
	_, _, err := ctx.CompileMultiple(strings.NewReader(`(defn cfsimple [] (+ 2 3))`))
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	v := ns.Lookup(vm.Symbol("cfsimple")).(*vm.Var).Deref()
	fn := v.(*vm.Func)
	irFn, err := ir.Build(fn.Chunk(), "cfsimple", 0, false)
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	t.Logf("before:\n%s", ir.Dump(irFn))

	changed := ConstFold(irFn)
	if !changed {
		t.Error("expected ConstFold to change the IR")
	}
	t.Logf("after:\n%s", ir.Dump(irFn))

	// Verify the Add node is now a Const with value 5.
	for _, n := range irFn.Nodes {
		if n.Op == ir.OpAdd {
			t.Errorf("expected no Add nodes after fold, found one")
		}
	}
	// And check there's a Const(5) somewhere.
	foundFive := false
	for _, n := range irFn.Nodes {
		if n.Op == ir.OpConst {
			if v, ok := n.Aux.(vm.Int); ok && int(v) == 5 {
				foundFive = true
				break
			}
		}
	}
	if !foundFive {
		t.Errorf("expected a Const(5) node after fold")
	}
}

// TestRoundtripOptimized: full pipeline — build → fold → DCE → lower → run.
// The optimized chunk should produce the same result as the original.
func TestRoundtripOptimized(t *testing.T) {
	consts := vm.NewConsts()
	ns := rt.NS(rt.NameCoreNS)
	ctx := compiler.NewCompiler(consts, ns)
	_, _, err := ctx.CompileMultiple(strings.NewReader(`(defn rtopt [] (+ (* 2 3) (- 10 5)))`))
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	v := ns.Lookup(vm.Symbol("rtopt")).(*vm.Var).Deref()
	fn := v.(*vm.Func)

	// Original.
	frame := vm.NewFrame(fn.Chunk(), nil)
	orig, err := frame.Run()
	vm.ReleaseFrame(frame)
	if err != nil {
		t.Fatalf("orig run: %v", err)
	}

	// Build → optimize → lower.
	irFn, err := ir.Build(fn.Chunk(), "rtopt", 0, false)
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	ConstFold(irFn)
	DCE(irFn)
	t.Logf("optimized IR:\n%s", ir.Dump(irFn))

	lowered, err := ir.Lower(irFn)
	if err != nil {
		t.Fatalf("Lower: %v", err)
	}

	// Run lowered.
	frame2 := vm.NewFrame(lowered, nil)
	optResult, err := frame2.Run()
	vm.ReleaseFrame(frame2)
	if err != nil {
		t.Fatalf("optimized run: %v", err)
	}

	t.Logf("orig = %s, optimized = %s", orig, optResult)
	if orig.String() != optResult.String() {
		t.Errorf("mismatch after optimization: orig %s vs lowered %s", orig, optResult)
	}

	// Bonus: optimized should have fewer instructions in its chunk.
	origLen := len(fn.Chunk().Code())
	optLen := len(lowered.Code())
	t.Logf("chunk size: orig %d ops, optimized %d ops", origLen, optLen)
	if optLen >= origLen {
		t.Logf("note: optimized chunk isn't smaller (%d vs %d) — lowering doesn't yet skip OpInvalid nodes optimally", optLen, origLen)
	}
}

// TestConstFoldDCE: const-fold + DCE should leave only Const 11 and Return.
func TestConstFoldDCE(t *testing.T) {
	consts := vm.NewConsts()
	ns := rt.NS(rt.NameCoreNS)
	ctx := compiler.NewCompiler(consts, ns)
	_, _, err := ctx.CompileMultiple(strings.NewReader(`(defn cfdce [] (+ (* 2 3) (- 10 5)))`))
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	v := ns.Lookup(vm.Symbol("cfdce")).(*vm.Var).Deref()
	fn := v.(*vm.Func)
	irFn, err := ir.Build(fn.Chunk(), "cfdce", 0, false)
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	ConstFold(irFn)
	DCE(irFn)
	t.Logf("after fold + DCE:\n%s", ir.Dump(irFn))

	// Count live (non-dead, non-terminator) nodes in entry block.
	live := 0
	for _, nid := range irFn.Blocks[irFn.Entry].Nodes {
		n := irFn.Node(nid)
		if n.Op != ir.OpInvalid {
			live++
		}
	}
	if live != 1 {
		t.Errorf("expected 1 live body node after fold+DCE, got %d", live)
	}
}

// TestConstFoldChained: (+ (* 2 3) (- 10 5)) should fold to Const 11.
// Two-pass: first Mul and Sub fold, then Add over the resulting consts.
func TestConstFoldChained(t *testing.T) {
	consts := vm.NewConsts()
	ns := rt.NS(rt.NameCoreNS)
	ctx := compiler.NewCompiler(consts, ns)
	_, _, err := ctx.CompileMultiple(strings.NewReader(`(defn cfchain [] (+ (* 2 3) (- 10 5)))`))
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	v := ns.Lookup(vm.Symbol("cfchain")).(*vm.Var).Deref()
	fn := v.(*vm.Func)
	irFn, err := ir.Build(fn.Chunk(), "cfchain", 0, false)
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	t.Logf("before:\n%s", ir.Dump(irFn))

	ConstFold(irFn)
	t.Logf("after:\n%s", ir.Dump(irFn))

	// Want a Const(11) somewhere.
	for _, n := range irFn.Nodes {
		if n.Op == ir.OpConst {
			if v, ok := n.Aux.(vm.Int); ok && int(v) == 11 {
				return // success
			}
		}
	}
	t.Error("expected a Const(11) after chained fold")
}
