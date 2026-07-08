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

	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

// ITER-0034 / STORY-0035 (EPIC-013, Component 2 of the multiblock-inliner +
// loop-unroll design). Task 2a: recognize the fold-over-rest idiom in IR.
//
// The idiom (fixtures below) is a variadic combinator whose body loops over its
// rest-arg seq, consuming (first v) / (rest v) and terminating on (empty? v):
//
// The element is the callee applied to a fixed input arg (the yamlstar
// combinator shape), so after specialization each (first v) becomes a known
// rule operand -> a direct call:
//   any: (loop [v fs] (if (empty? v) false (if ((first v) input) true (recur (rest v)))))
//   all: (loop [v fs] (if (empty? v) true  (if ((first v) input) (recur (rest v)) false)))

func foldDesc(t *testing.T, src string, field string) string {
	t.Helper()
	f := buildLispIR(t, src)
	return strings.TrimSpace(lispEval(t, "("+field+" (ir.passes.inline/fold-over-rest? %s))", f))
}

func foldRecognized(t *testing.T, src string) string {
	t.Helper()
	f := buildLispIR(t, src)
	return strings.TrimSpace(lispEval(t, `(some? (ir.passes.inline/fold-over-rest? %s))`, f))
}

const anycSrc = `(defn anyc [input & fs] (loop [v fs] (if (empty? v) false (if ((first v) input) true (recur (rest v))))))`
const allcSrc = `(defn allc [input & fs] (loop [v fs] (if (empty? v) true (if ((first v) input) (recur (rest v)) false))))`

func TestFoldOverRestRecognizesAny(t *testing.T) {
	ensureLoader()
	if got := foldRecognized(t, anycSrc); got != "true" {
		t.Fatalf("anyc should be recognized as fold-over-rest, got %s", got)
	}
	if got := foldDesc(t, anycSrc, ":kind"); got != ":any" {
		t.Fatalf("anyc kind should be :any, got %s", got)
	}
	// rest-index is the last load-arg (arity-1); for [p & xs] arity=2 -> 1.
	if got := foldDesc(t, anycSrc, ":rest-index"); got != "1" {
		t.Fatalf("anyc rest-index should be 1, got %s", got)
	}
}

func TestFoldOverRestRecognizesAll(t *testing.T) {
	ensureLoader()
	if got := foldRecognized(t, allcSrc); got != "true" {
		t.Fatalf("allc should be recognized, got %s", got)
	}
	if got := foldDesc(t, allcSrc, ":kind"); got != ":all" {
		t.Fatalf("allc kind should be :all, got %s", got)
	}
}

func TestFoldOverRestRejectsCounterLoop(t *testing.T) {
	ensureLoader()
	// A counter loop over a fixed arg — NOT a fold over the rest seq.
	src := `(defn cu [n] (loop [i 0] (if (< i n) (recur (inc i)) i)))`
	if got := foldRecognized(t, src); got != "false" {
		t.Fatalf("counter loop must NOT be recognized as fold-over-rest, got %s", got)
	}
}

func TestFoldOverRestRejectsNonVariadic(t *testing.T) {
	ensureLoader()
	// Non-variadic, no rest arg at all.
	if got := foldRecognized(t, `(defn tiny [x] (+ x 1))`); got != "false" {
		t.Fatalf("non-variadic fn must NOT be recognized, got %s", got)
	}
}

func TestFoldOverRestRejectsVariadicNonFold(t *testing.T) {
	ensureLoader()
	// Variadic but does not loop over the rest seq (just returns it).
	if got := foldRecognized(t, `(defn vf [p & xs] xs)`); got != "false" {
		t.Fatalf("variadic non-fold fn must NOT be recognized, got %s", got)
	}
}

// ── T2: specialize a recognized fold-over-rest to a fixed-arity unrolled fn ──

func specializeFold(t *testing.T, comboSrc string, n int) vm.Value {
	t.Helper()
	f := buildLispIR(t, comboSrc)
	return lispEvalReturn(t, f,
		fmt.Sprintf(`(ir.passes.inline/specialize-fold (ir.passes.inline/fold-over-rest? %%s) %d)`, n))
}

// assertValidatesVal runs ir.validate/validate-fn! on a Function value; a
// throw (malformed IR) surfaces as an eval error and fails the test.
func assertValidatesVal(t *testing.T, f vm.Value, label string) {
	t.Helper()
	passVarCounter++
	varName := fmt.Sprintf("*vf-%d*", passVarCounter)
	coreNS := rt.NS(rt.NameCoreNS)
	coreNS.Def(varName, f)
	consts := vm.NewConsts()
	c := compiler.NewCompiler(consts, coreNS)
	c.SetSource("validate-val")
	expr := fmt.Sprintf(`(ir.validate/validate-fn! %s "%s")`, varName, label)
	if _, _, err := c.CompileMultiple(strings.NewReader(expr)); err != nil {
		t.Fatalf("specialized IR failed validation (%s): %v", label, err)
	}
}

func TestSpecializeFoldAnyUnrolls(t *testing.T) {
	ensureLoader()
	spec := specializeFold(t, anycSrc, 3)
	// arity = num-fixed(1: input) + N(3 elements) = 4; non-variadic.
	if got := strings.TrimSpace(lispEvalOn(t, spec, "(ir/fn-arity %s)")); got != "4" {
		t.Fatalf("specialized any arity should be 4, got %s", got)
	}
	if got := strings.TrimSpace(lispEvalOn(t, spec, "(ir/fn-variadic? %s)")); got != "false" {
		t.Fatalf("specialized fn should be non-variadic, got %s", got)
	}
	dump := lispDump(t, spec)
	// Unrolled any(N=3) = 3 short-circuit branch-ifs, loop-free.
	if n := strings.Count(dump, "BranchIf"); n != 3 {
		t.Fatalf("expected 3 BranchIf (unrolled chain), got %d:\n%s", n, dump)
	}
	assertValidatesVal(t, spec, "specialize-any")
}

func TestSpecializeFoldAllUnrolls(t *testing.T) {
	ensureLoader()
	spec := specializeFold(t, allcSrc, 2)
	if got := strings.TrimSpace(lispEvalOn(t, spec, "(ir/fn-arity %s)")); got != "3" {
		t.Fatalf("specialized all arity should be 3, got %s", got)
	}
	dump := lispDump(t, spec)
	if n := strings.Count(dump, "BranchIf"); n != 2 {
		t.Fatalf("expected 2 BranchIf (unrolled chain), got %d:\n%s", n, dump)
	}
	assertValidatesVal(t, spec, "specialize-all")
}
