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

func runTypeInfer(t *testing.T, f vm.Value) vm.Value {
	t.Helper()
	return runLispPass(t, "ir.passes.typeinfer", "typeinfer", f)
}

func seedArgTypes(t *testing.T, f vm.Value, seedExpr string) {
	t.Helper()
	passVarCounter++
	varName := fmt.Sprintf("*typeinfer-seed-%d*", passVarCounter)
	rt.NS(rt.NameCoreNS).Def(varName, f)
	runLispExpr(t, fmt.Sprintf(`(swap! %s assoc :arg-types %s)`, varName, seedExpr))
}

func TestTypeInferDumpsPrettyUnionTypes(t *testing.T) {
	ensureLoader()

	fn := buildLispIR(t, `(defn branchy [flag x]
	                       (if flag
	                         x
	                         nil))`)
	seedArgTypes(t, fn, "[:bool :int]")
	runTypeInfer(t, fn)
	dump := lispDump(t, fn)

	if !strings.Contains(dump, "union{int,nil}") {
		t.Fatalf("expected pretty union type in dump\n--- dump ---\n%s", dump)
	}
}

func TestTypeInferUsesArgSeedsForLoadArg(t *testing.T) {
	ensureLoader()

	fn := buildLispIR(t, `(defn seeded-add [x y] (+ x y))`)
	seedArgTypes(t, fn, "[:int :int]")
	runTypeInfer(t, fn)
	dump := lispDump(t, fn)

	if !strings.Contains(dump, "v0 = LoadArg ; 0 : int") || !strings.Contains(dump, "v1 = LoadArg ; 1 : int") {
		t.Fatalf("expected seeded load args to infer as int\n--- dump ---\n%s", dump)
	}
	if !strings.Contains(dump, "Add v0 v1 : int") {
		t.Fatalf("expected add of seeded int args to infer as int\n--- dump ---\n%s", dump)
	}
}

func TestTypeInferJoinsBranchValuesIntoBlockParam(t *testing.T) {
	ensureLoader()

	fn := buildLispIR(t, `(defn maybe-inc [flag x]
	                       (+ (if flag x nil) 1))`)
	seedArgTypes(t, fn, "[:bool :int]")
	runTypeInfer(t, fn)
	dump := lispDump(t, fn)

	if !strings.Contains(dump, "union{int,nil}") {
		t.Fatalf("expected join block param to infer union{int,nil}\n--- dump ---\n%s", dump)
	}
}

func TestTypeInferConvergesLoopCarriedBlockParams(t *testing.T) {
	ensureLoader()

	fn := buildLispIR(t, `(defn sum-to-n [n]
	                       (loop* [i 0 acc 0]
	                         (if (< i n)
	                           (recur (+ i 1) (+ acc i))
	                           acc)))`)
	seedArgTypes(t, fn, "[:int]")
	runTypeInfer(t, fn)
	dump := lispDump(t, fn)

	if !strings.Contains(dump, "fn sum-to-n(arity=1, variadic=false):") {
		t.Fatalf("expected valid dump output\n--- dump ---\n%s", dump)
	}
	if strings.Contains(dump, "union{bottom") || strings.Contains(dump, ": bottom") {
		t.Fatalf("expected loop inference to converge beyond bottom\n--- dump ---\n%s", dump)
	}
	if strings.Count(dump, ": int") < 5 {
		t.Fatalf("expected loop-carried values to stabilize as ints\n--- dump ---\n%s", dump)
	}
}

func TestTypeInferNormalizesNumericUnionToNumber(t *testing.T) {
	ensureLoader()

	fn := buildLispIR(t, `(defn numeric-join [flag]
	                       (if flag
	                         1
	                         1.5))`)
	runTypeInfer(t, fn)
	dump := lispDump(t, fn)

	if !strings.Contains(dump, "number") {
		t.Fatalf("expected int/float join to normalize to number\n--- dump ---\n%s", dump)
	}
	if strings.Contains(dump, "union{float,int}") || strings.Contains(dump, "union{int,float}") {
		t.Fatalf("expected numeric union to collapse to number\n--- dump ---\n%s", dump)
	}
}

func TestTypeInferTracksBooleanLiteralsPrecisely(t *testing.T) {
	ensureLoader()

	fn := buildLispIR(t, `(defn bool-consts []
	                       (if true false true))`)
	runTypeInfer(t, fn)
	dump := lispDump(t, fn)

	if !strings.Contains(dump, ": true") || !strings.Contains(dump, ": false") {
		t.Fatalf("expected true/false literals to retain precise types\n--- dump ---\n%s", dump)
	}
}

func TestTypeInferTracksTypedZeroLiterals(t *testing.T) {
	ensureLoader()

	fn := buildLispIR(t, `(defn zeroes []
	                       (if true 0 0.0))`)
	runTypeInfer(t, fn)
	dump := lispDump(t, fn)

	if !strings.Contains(dump, "int(0)") {
		t.Fatalf("expected integer zero to dump as typed zero\n--- dump ---\n%s", dump)
	}
	if !strings.Contains(dump, "float(0.0)") {
		t.Fatalf("expected float zero to dump as typed zero\n--- dump ---\n%s", dump)
	}
}

func TestTypeInferRefinesTruthyEdgeForMaybeNilArithmetic(t *testing.T) {
	ensureLoader()

	fn := buildLispIR(t, `(defn maybe-plus-one [x]
	                       (if x
	                         (+ x 1)
	                         0))`)
	seedArgTypes(t, fn, "[[:union :nil :int]]")
	runTypeInfer(t, fn)
	dump := lispDump(t, fn)

	if !strings.Contains(dump, "union{int,nil}") {
		t.Fatalf("expected joined result to remain union{int,nil}\n--- dump ---\n%s", dump)
	}
	if !strings.Contains(dump, "Add v") || !strings.Contains(dump, ": int") {
		t.Fatalf("expected truthy-edge arithmetic to infer int\n--- dump ---\n%s", dump)
	}
}

func TestLatticePreservesExactDTypeFacts(t *testing.T) {
	ensureLoader()

	got := runLispExpr(t, `(pr-str
		[(ir.lattice/normalize-type [:dtype 'Square])
		 (ir.lattice/type-join [:dtype 'Square] [:dtype 'Square])
		 (ir.lattice/type-join [:dtype 'Square] :int)])`)

	if string(got.(vm.String)) != "[[:dtype Square] [:dtype Square] :any]" {
		t.Fatalf("unexpected dtype lattice behavior: %s", got)
	}
}

func TestTypeInfraStateSeedsWithoutTouchingInstTypesUntilFlush(t *testing.T) {
	ensureLoader()

	fn := buildLispIR(t, `(defn seeded-add [x y] (+ x y))`)

	passVarCounter++
	fnVar := fmt.Sprintf("*typeinfra-fn-%d*", passVarCounter)
	rt.NS(rt.NameCoreNS).Def(fnVar, fn)

	state := runLispExpr(t, fmt.Sprintf(`(let [s (ir.lattice/new-typeinfra-state %s)]
		(ir.lattice/seed-state-from-inst-types! s %s)
		(ir.passes.infer-arg-types/infer-arg-types %s s)
		s)`, fnVar, fnVar, fnVar))

	passVarCounter++
	stateVar := fmt.Sprintf("*typeinfra-state-%d*", passVarCounter)
	rt.NS(rt.NameCoreNS).Def(stateVar, state)

	before := lispDump(t, fn)
	if strings.Contains(before, "v0 = LoadArg ; 0 : int") || strings.Contains(before, "v1 = LoadArg ; 1 : int") {
		t.Fatalf("expected shared state seeding to avoid direct inst type writes before flush\n--- dump ---\n%s", before)
	}

	seeded := runLispExpr(t, fmt.Sprintf(`(pr-str
		[(ir.lattice/state-type %s (ir/fn-load-arg %s 0))
		 (ir.lattice/state-type %s (ir/fn-load-arg %s 1))])`,
		stateVar, fnVar, stateVar, fnVar))
	if string(seeded.(vm.String)) != "[:int :int]" {
		beforeState := runLispExpr(t, fmt.Sprintf(`(let [s (ir.lattice/new-typeinfra-state %s)]
			(ir.lattice/seed-state-from-inst-types! s %s)
			(pr-str
			  [(ir.lattice/state-type s (ir/fn-load-arg %s 0))
			   (ir.lattice/state-type s (ir/fn-load-arg %s 1))
			   (do (ir.lattice/join-inst-type! s (ir/fn-load-arg %s 0) :int)
			       (ir.lattice/state-type s (ir/fn-load-arg %s 0)))]))`,
			fnVar, fnVar, fnVar, fnVar, fnVar, fnVar))
		scan := runLispExpr(t, fmt.Sprintf(`(pr-str
			(map (fn [i] [(ir/op i %s) (ir/refs i %s)])
			     (range (ir/inst-count %s))))`, fnVar, fnVar, fnVar))
		runLispExpr(t, fmt.Sprintf(`(ir.passes.infer-arg-types/infer-arg-types %s)`, fnVar))
		t.Fatalf("unexpected state seeds: %s\n--- state-debug ---\n%s\n--- scan ---\n%s\n--- dump-after-wrapper ---\n%s", seeded, beforeState, scan, lispDump(t, fn))
	}

	runLispExpr(t, fmt.Sprintf(`(ir.lattice/flush-state-types! %s %s)`, stateVar, fnVar))
	after := lispDump(t, fn)
	if !strings.Contains(after, "v0 = LoadArg ; 0 : int") || !strings.Contains(after, "v1 = LoadArg ; 1 : int") {
		t.Fatalf("expected flush to materialize inferred arg seeds\n--- dump ---\n%s", after)
	}
}

// The drain is LIFO, so the call to transient is visited before the [] it
// reads. Its type has to wait for that operand; settling on :unknown first
// leaves the call untyped until a second epoch, so the facts would depend on
// visit order.
func TestTypeInferTypesTransientCallOnFirstEpoch(t *testing.T) {
	ensureLoader()

	fn := buildLispIR(t, `(defn mapv-shaped [f coll]
	                       (persistent! (reduce (fn* [acc x] (conj! acc (f x))) (transient []) coll)))`)

	passVarCounter++
	fnVar := fmt.Sprintf("*typeinfra-fn-%d*", passVarCounter)
	rt.NS(rt.NameCoreNS).Def(fnVar, fn)

	for epoch := 1; epoch <= 2; epoch++ {
		runLispExpr(t, fmt.Sprintf(`(let [s (ir.lattice/new-typeinfra-state %[1]s)
		                                  s (ir.lattice/seed-state-from-inst-types! s %[1]s)
		                                  s (ir.passes.infer-arg-types/infer-arg-types %[1]s s)
		                                  s (ir.passes.typeinfer/typeinfer %[1]s s)]
		                              (ir.lattice/flush-state-types! s %[1]s))`, fnVar))
		got := runLispExpr(t, fmt.Sprintf(`(pr-str
			(some (fn* [i]
			        (let [callee (first (ir/refs i %[1]s))]
			          (when (and (= :call (ir/op i %[1]s))
			                     (= :load-var (ir/op callee %[1]s))
			                     (= "transient" (name (symbol (ir/aux callee %[1]s)))))
			            (ir/type-of i %[1]s))))
			      (range (ir/inst-count %[1]s))))`, fnVar))
		if string(got.(vm.String)) != ":transient-vector" {
			t.Fatalf("after epoch %d, (transient []) is typed %s, want :transient-vector", epoch, got)
		}
	}
}

// optimize-fn skips its trailing typeinfra epoch when no middle pass edited the
// function, and runs it when one did. Counting typeinfer calls observes the
// skip without touching run-typeinfra-epoch, whose identity the skip checks.
func TestOptimizeFnSkipsTrailingEpochOnlyWhenUnedited(t *testing.T) {
	ensureLoader()
	got := runLispString(t, `(pr-str
		(let [epochs (fn [form]
		               (let [n  (atom 0)
		                     ti ir.passes.typeinfer/typeinfer
		                     f  (ir.build/build-fn form)]
		                 (with-redefs [ir.passes.typeinfer/typeinfer
		                               (fn [& args] (swap! n inc) (apply ti args))]
		                   (ir.passes.pipeline/optimize-fn f))
		                 @n))]
		  [;; no middle pass changes this one
		   (epochs '(defn unedited [x y] (+ x y)))
		   ;; constfold rewrites (+ 1 2)
		   (epochs '(defn folded [x] (+ x (+ 1 2))))]))`)
	if got != "[1 2]" {
		t.Fatalf("typeinfer epochs per optimize-fn: got %s, want [1 2] (unedited, folded)", got)
	}
}
