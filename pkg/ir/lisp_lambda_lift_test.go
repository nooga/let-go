/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package ir_test

import (
	"strings"
	"testing"
)

func TestLambdaLiftNoClosuresIsIdentity(t *testing.T) {
	ensureLoader()
	f := buildLispIR(t, `(defn add1 [x] (+ x 1))`)
	// lambda-lift returns a map; with no inner closures, :lifted is empty
	// and :main dumps identically to the input.
	before := lispDump(t, f)
	out := runLispPassValue(t, "ir.passes.lambda-lift", "lambda-lift", f)
	lifted := lispEvalOn(t, out, `(count (:lifted %s))`)
	if strings.TrimSpace(lifted) != "0" {
		t.Fatalf("expected 0 lifted siblings, got %s", lifted)
	}
	after := lispDumpValue(t, lispEvalReturn(t, out, `(:main %s)`))
	if before != after {
		t.Fatalf("identity broken:\nbefore:\n%s\nafter:\n%s", before, after)
	}
}

func TestLambdaLiftCaptureFreeAndCapturing(t *testing.T) {
	ensureLoader()
	// Capture-free: inner fn* uses only its own params (q)
	free := buildLispIR(t, `(defn o1 [p] (identity (fn* [q] q)))`)
	if got := strings.TrimSpace(lispEval(t,
		`(:captures (first (ir.passes.lambda-lift/candidate-info %s)))`, free)); got != "[]" {
		t.Fatalf("capture-free should have [] captures, got %s", got)
	}
	// Capturing: inner fn* closes over let binding y (which comes from param p)
	cap := buildLispIR(t, `(defn o2 [p] (let [y p] (identity (fn* [q] (vector q y)))))`)
	if got := strings.TrimSpace(lispEval(t,
		`(count (:captures (first (ir.passes.lambda-lift/candidate-info %s))))`, cap)); got == "0" {
		t.Fatalf("capturing closure should report >0 captures, got %s", got)
	}
}

func TestLambdaLiftBuildsSiblingDeterministicName(t *testing.T) {
	ensureLoader()
	out := runLispPassValue(t, "ir.passes.lambda-lift", "lambda-lift",
		buildLispIR(t, `(defn c_x [p] (identity (fn* [q] q)))`))
	names := lispEvalOn(t, out, `(mapv (fn* [g] (ir/fn-name g)) (:lifted %s))`)
	if !strings.Contains(names, "c_x__lifted0") {
		t.Fatalf("want a sibling named c_x__lifted0, got %s", names)
	}
}

func TestLambdaLiftRewritesSiteToLoadVar(t *testing.T) {
	ensureLoader()
	out := runLispPassValue(t, "ir.passes.lambda-lift", "lambda-lift",
		buildLispIR(t, `(defn c_x [p] (identity (fn* [q] q)))`))
	mainDump := lispDumpValue(t, lispEvalReturn(t, out, `(:main %s)`))
	if strings.Contains(mainDump, ":fn-template") {
		t.Fatalf(":main still contains an inline fn-template const after lifting:\n%s", mainDump)
	}
	if !strings.Contains(mainDump, "c_x__lifted0") {
		t.Fatalf(":main should reference the lifted sibling by name:\n%s", mainDump)
	}
}

func TestCompileFormEmitsLiftedSiblingDecl(t *testing.T) {
	ensureLoader()
	// Compile a go-target form with one inner capture-free fn used as a value,
	// assert the lowered output contains BOTH the outer decl and the lifted sibling.
	out := compileFormToGo(t, `(defn c_x [p] (identity (fn* [q] q)))`)
	// Public defn lowers to exported PascalCase, and its lifted sibling with it (#371).
	if !strings.Contains(out, "C_x__lifted0") || !strings.Contains(out, "func C_x") {
		t.Fatalf("expected both outer and lifted decls in lowered Go:\n%s", out)
	}
}
