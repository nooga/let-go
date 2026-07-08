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
