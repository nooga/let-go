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

// TestCSE_Simple: (defn f [x] (+ (+ x 1) (+ x 1)))
//
// Before CSE: two separate (+ x 1) computations (Add with same args).
// After CSE: one Add, referenced twice by the outer Add.
//
// The CSE itself succeeds (collapses 3 Adds to 2). Lowering then fails
// because the spike's naive Lower assumes each value is consumed once;
// a shared value would need DUP or local-slot emission. See task #30.
// We assert the IR transformation only, not the round-trip.
func TestCSE_Simple_IRonly(t *testing.T) {
	consts := vm.NewConsts()
	ns := rt.NS(rt.NameCoreNS)
	ctx := compiler.NewCompiler(consts, ns)
	src := `(defn csesimple [x] (+ (+ x 1) (+ x 1)))`
	_, _, err := ctx.CompileMultiple(strings.NewReader(src))
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	v := ns.Lookup(vm.Symbol("csesimple")).(*vm.Var).Deref()
	fn := v.(*vm.Func)
	irFn, err := ir.Build(fn.Chunk(), "csesimple", 1, false)
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	t.Logf("before CSE:\n%s", ir.Dump(irFn))

	// Count Add nodes before.
	addsBefore := countOp(irFn, ir.OpAdd)
	if addsBefore != 3 {
		t.Logf("expected 3 Add nodes initially, got %d (may be ok if compiler already shared)", addsBefore)
	}

	changed := CSE(irFn)
	if !changed {
		t.Error("expected CSE to change something")
	}
	t.Logf("after CSE:\n%s", ir.Dump(irFn))

	addsAfter := countOp(irFn, ir.OpAdd)
	t.Logf("Add nodes: before=%d after=%d", addsBefore, addsAfter)
	if addsAfter >= addsBefore {
		t.Errorf("expected fewer Add nodes after CSE; before=%d after=%d", addsBefore, addsAfter)
	}

	// IR-only validation. Round-trip blocked by task #30 (spike lower
	// doesn't handle multi-use values).
}

// TestCSE_Simple: full pipeline test for CSE + multi-use lowering.
// (defn f [x] (+ (+ x 1) (+ x 1))) — the inner (+ x 1) appears twice
// and CSE merges them; the outer + then has both operands referencing
// the same shared value (>1 use). Lower must emit DUP_NTH (or equivalent)
// to handle the multi-use.
func TestCSE_Simple(t *testing.T) {
	consts := vm.NewConsts()
	ns := rt.NS(rt.NameCoreNS)
	ctx := compiler.NewCompiler(consts, ns)
	src := `(defn csefn [x] (+ (+ x 1) (+ x 1)))`
	_, _, err := ctx.CompileMultiple(strings.NewReader(src))
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	v := ns.Lookup(vm.Symbol("csefn")).(*vm.Var).Deref()
	fn := v.(*vm.Func)
	irFn, err := ir.Build(fn.Chunk(), "csefn", 1, false)
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	t.Logf("before CSE:\n%s", ir.Dump(irFn))
	CSE(irFn)
	t.Logf("after CSE:\n%s", ir.Dump(irFn))
	DCE(irFn)
	chunk, err := ir.Lower(irFn)
	if err != nil {
		t.Fatalf("Lower: %v", err)
	}
	// Run with x=3: (+ (+ 3 1) (+ 3 1)) = 8
	frame := vm.NewFrame(chunk, []vm.Value{vm.Int(3)})
	result, err := frame.Run()
	vm.ReleaseFrame(frame)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if result.String() != "8" {
		t.Errorf("expected 8, got %s", result)
	}
}

// TestCSE_ConstantConsolidation: two Const nodes with the same value
// (e.g., two literal 1s emitted by the compiler) should collapse to one.
func TestCSE_ConstantConsolidation(t *testing.T) {
	consts := vm.NewConsts()
	ns := rt.NS(rt.NameCoreNS)
	ctx := compiler.NewCompiler(consts, ns)
	src := `(defn cseconst [x] (- (+ x 1) (+ x 1)))`
	_, _, err := ctx.CompileMultiple(strings.NewReader(src))
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	v := ns.Lookup(vm.Symbol("cseconst")).(*vm.Var).Deref()
	fn := v.(*vm.Func)
	irFn, err := ir.Build(fn.Chunk(), "cseconst", 1, false)
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	constsBefore := countOp(irFn, ir.OpConst)
	CSE(irFn)
	constsAfter := countOp(irFn, ir.OpConst)
	t.Logf("Const nodes: before=%d after=%d", constsBefore, constsAfter)
	if constsAfter >= constsBefore {
		t.Errorf("expected fewer Const nodes after CSE; before=%d after=%d", constsBefore, constsAfter)
	}
}

// TestCSE_FoldThenCSE: (* 2 3) folded to 6 in both arms of (+ (* 2 3) (* 2 3))
// — after fold both arms become Const(6) — then CSE merges them to one.
func TestCSE_FoldThenCSE(t *testing.T) {
	consts := vm.NewConsts()
	ns := rt.NS(rt.NameCoreNS)
	ctx := compiler.NewCompiler(consts, ns)
	src := `(defn fcse [] (+ (* 2 3) (* 2 3)))`
	_, _, err := ctx.CompileMultiple(strings.NewReader(src))
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	v := ns.Lookup(vm.Symbol("fcse")).(*vm.Var).Deref()
	fn := v.(*vm.Func)
	irFn, err := ir.Build(fn.Chunk(), "fcse", 0, false)
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	ConstFold(irFn)
	CSE(irFn)
	DCE(irFn)
	t.Logf("after fold+CSE+DCE:\n%s", ir.Dump(irFn))

	// Roundtrip and check correctness.
	lowered, err := ir.Lower(irFn)
	if err != nil {
		t.Fatalf("Lower: %v", err)
	}
	frame := vm.NewFrame(lowered, nil)
	result, err := frame.Run()
	vm.ReleaseFrame(frame)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if result.String() != "12" {
		t.Errorf("expected 12, got %s", result)
	}
}

func countOp(f *ir.Function, op ir.Op) int {
	n := 0
	for _, node := range f.Nodes {
		if node.Op == op {
			n++
		}
	}
	return n
}
