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
