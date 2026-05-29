/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package ir_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

func optimizeWithPipeline(t *testing.T, f vm.Value) vm.Value {
	t.Helper()
	passVarCounter++
	varName := fmt.Sprintf("*pipeline-fn-%d*", passVarCounter)
	rt.NS(rt.NameCoreNS).Def(varName, f)
	runLispExpr(t, fmt.Sprintf(`(ir.passes.pipeline/optimize-fn %s)`, varName))
	return f
}

func validateIR(t *testing.T, f vm.Value, label string) {
	t.Helper()
	passVarCounter++
	varName := fmt.Sprintf("*validate-fn-%d*", passVarCounter)
	rt.NS(rt.NameCoreNS).Def(varName, f)
	runLispExpr(t, fmt.Sprintf(`(ir.validate/validate-fn! %s %q)`, varName, label))
}

func TestLispCSERewritesDuplicateSameBlockPureExpr(t *testing.T) {
	ensureLoader()

	fn := buildLispIR(t, `(defn cse-square [x] (+ (* x x) (* x x)))`)
	runLispPass(t, "ir.passes.cse", "cse", fn)
	validateIR(t, fn, "test-cse")
	dump := lispDump(t, fn)

	if !strings.Contains(dump, "Add v1 v1") {
		t.Fatalf("expected CSE to rewrite Add refs to the same Mul result\n--- dump ---\n%s", dump)
	}
}

func TestPipelineOptimizeRunsCSEForSameBlockPureExpr(t *testing.T) {
	ensureLoader()

	fn := buildLispIR(t, `(defn cse-poly [x] (+ (* x x) (* x x) 5))`)
	optimizeWithPipeline(t, fn)
	dump := lispDump(t, fn)

	if strings.Count(dump, "Mul v0 v0") != 1 {
		t.Fatalf("expected pipeline CSE+DCE to leave one multiply\n--- dump ---\n%s", dump)
	}
	if !strings.Contains(dump, "Add v1 v1") {
		t.Fatalf("expected pipeline CSE to rewrite duplicate multiply uses\n--- dump ---\n%s", dump)
	}
}

func TestPipelineConstFoldCanonicalizationFeedsCSE(t *testing.T) {
	ensureLoader()

	fn := buildLispIR(t, `(defn cse-canon [x] (+ (+ x 5) (+ 5 x)))`)
	optimizeWithPipeline(t, fn)
	dump := lispDump(t, fn)

	if strings.Count(dump, "Add v0 v1") != 1 {
		t.Fatalf("expected canonicalized inner adds to be CSE'd into one Add v0 v1\n--- dump ---\n%s", dump)
	}
	if !strings.Contains(dump, "Add v2 v2") {
		t.Fatalf("expected outer add to reuse the single canonical inner add\n--- dump ---\n%s", dump)
	}
}

func TestLispCSEMergesStableLoadVar(t *testing.T) {
	ensureLoader()

	fn := buildLispIR(t, `(defn cse-load-var [x] (do (str x) (str x)))`)
	runLispPass(t, "ir.passes.cse", "cse", fn)
	validateIR(t, fn, "test-cse-load-var")
	dump := lispDump(t, fn)

	if !strings.Contains(dump, "Call v1 v0 ; 1") || !strings.Contains(dump, "Return v4") {
		t.Fatalf("expected calls to use the surviving LoadVar without changing call structure\n--- dump ---\n%s", dump)
	}
	if strings.Contains(dump, "Call v3 v0 ; 1") {
		t.Fatalf("expected duplicate LoadVar use to be rewritten to the first LoadVar\n--- dump ---\n%s", dump)
	}
}

func TestLispCSEDoesNotMergeLoadVarWhenMutatingBuiltinCalled(t *testing.T) {
	ensureLoader()

	fn := buildLispIR(t, `(defn cse-load-var-mutating [x]
	                       (do (alter-var-root x identity)
	                           (str x)
	                           (str x)))`)
	runLispPass(t, "ir.passes.cse", "cse", fn)
	validateIR(t, fn, "test-cse-load-var-mutating")
	dump := lispDump(t, fn)

	if strings.Count(dump, "LoadVar ; #'core/str") != 2 {
		t.Fatalf("expected LoadVar #'core/str not to be CSE'd after mutating builtin\n--- dump ---\n%s", dump)
	}
}

func TestMutabilityAnalysisReportsStableLoadVar(t *testing.T) {
	ensureLoader()

	fn := buildLispIR(t, `(defn mut-stable [x] (do (str x) (str x)))`)
	passVarCounter++
	varName := fmt.Sprintf("*mut-stable-fn-%d*", passVarCounter)
	rt.NS(rt.NameCoreNS).Def(varName, fn)

	result := runLispExpr(t, fmt.Sprintf(`
		(do
		  (require 'ir.passes.mutability)
		  (let [facts (ir.passes.mutability/analyze-var-stability %s)]
		    (ir.passes.mutability/stable-load-var? facts #'core/str)))`, varName))

	if result != vm.TRUE {
		t.Fatalf("expected #'core/str to be stable, got %v", result)
	}
}

func TestMutabilityAnalysisReportsUnknownAfterMutatingBuiltin(t *testing.T) {
	ensureLoader()

	fn := buildLispIR(t, `(defn mut-unknown [x]
	                       (do (alter-var-root x identity)
	                           (str x)
	                           (str x)))`)
	passVarCounter++
	varName := fmt.Sprintf("*mut-unknown-fn-%d*", passVarCounter)
	rt.NS(rt.NameCoreNS).Def(varName, fn)

	result := runLispExpr(t, fmt.Sprintf(`
		(do
		  (require 'ir.passes.mutability)
		  (let [facts (ir.passes.mutability/analyze-var-stability %s)]
		    (ir.passes.mutability/stable-load-var? facts #'core/str)))`, varName))

	if result != vm.FALSE {
		t.Fatalf("expected #'core/str to be unstable after mutating builtin, got %v", result)
	}
}
