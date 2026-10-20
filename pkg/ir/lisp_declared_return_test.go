/*
 * Copyright (c) 2026 Marcin Gasperowicz <xnooga@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package ir_test

// #599: a self-recursive integer fn used to lower with unboxed params but a
// boxed return — `(vm.Value, error)` plus `rt.AddValue` for the addition —
// because the self-call typed as :unknown and poisoned the numeric chain.
// Typing the self-calls as T is sound when the body then returns T on every
// path, so T is a guess the body must confirm: the declared return type
// (`(defn f ^long [n] …)`) when there is one, else :int, then :float. These
// tests pin the resulting shape, that a guess the body contradicts is dropped,
// and that the seed does not leak across function boundaries.

import (
	"strings"
	"testing"
)

func TestDeclaredReturnTypeUnboxesSelfRecursion(t *testing.T) {
	ensureLoader()

	rendered := lowerForms(t, "drpkg",
		`(defn drfib ^long [^long n] (if (< n 2) n (+ (drfib (- n 1)) (drfib (- n 2)))))`)

	if !strings.Contains(rendered, "func Drfib(ec *vm.ExecContext, arg0 int64) (int64, error)") {
		t.Fatalf("expected a natively-typed int64 return; rendered:\n%s", rendered)
	}
	// The payoff: generic boxed addition is gone in favour of int + int.
	if strings.Contains(rendered, "rt.AddValue") {
		t.Errorf("addition still dispatches through rt.AddValue; rendered:\n%s", rendered)
	}
}

// Without a hint the self-calls are guessed :int, and the body confirms it.
func TestUnhintedRecursionInfersReturn(t *testing.T) {
	ensureLoader()

	rendered := lowerForms(t, "drpkg",
		`(defn drplain [^long n] (if (< n 2) n (+ (drplain (- n 1)) (drplain (- n 2)))))`)

	if !strings.Contains(rendered, "func Drplain(ec *vm.ExecContext, arg0 int64) (int64, error)") {
		t.Fatalf("expected an inferred int64 return; rendered:\n%s", rendered)
	}
	if strings.Contains(rendered, "rt.AddValue") {
		t.Errorf("addition still dispatches through rt.AddValue; rendered:\n%s", rendered)
	}
}

// cond's :else arm used to leave a dead nil arm in the IR, so the return
// joined int with nil and stayed boxed. Branch folding drops it.
func TestCondElseRecursionInfersReturn(t *testing.T) {
	ensureLoader()

	rendered := lowerForms(t, "drpkg",
		`(defn drcond [n] (cond (= n 0) 1 (= n 1) 1 :else (+ (drcond (- n 2)) (drcond (- n 1)))))`)

	if !strings.Contains(rendered, "func Drcond(ec *vm.ExecContext, arg0 int64) (int64, error)") {
		t.Fatalf("expected an inferred int64 return; rendered:\n%s", rendered)
	}
	if strings.Contains(rendered, `vm.Keyword("else")`) {
		t.Errorf(":else test survived lowering; rendered:\n%s", rendered)
	}
}

// A guess the body contradicts is dropped. With the nil arm drnil really
// returns vm.Value, so its self-call must stay boxed too: typing it from the
// hint assigned Drnil's vm.Value result into an int64 temp, which does not
// compile, and added with a raw +.
func TestContradictedReturnHintStaysBoxed(t *testing.T) {
	ensureLoader()

	rendered := lowerForms(t, "drpkg",
		`(defn drnil ^long [n] (if (pos? n) (+ 1 (drnil (dec n))) nil))`)

	if !strings.Contains(rendered, "(vm.Value, error)") {
		t.Fatalf("expected a boxed return; rendered:\n%s", rendered)
	}
	if !strings.Contains(rendered, "rt.AddValue") {
		t.Errorf("addition of the boxed self-call result went native; rendered:\n%s", rendered)
	}
}

// The seed is restricted to self-calls. A declaration the body contradicts must
// not reach a *caller's* inference: before that restriction, drcaller returned
// (int, error) and assigned vm.String(...) into an int temp, which does not
// compile. drliar itself is free to lower honestly as string.
func TestContradictoryHintDoesNotCorruptCallers(t *testing.T) {
	ensureLoader()

	rendered := lowerForms(t, "drpkg",
		`(defn drliar ^long [n] "not-an-int") (defn drcaller [n] (drliar n))`)

	if strings.Contains(rendered, "func Drcaller(ec *vm.ExecContext, arg0 vm.Value) (int, error)") {
		t.Fatalf("callee's contradictory hint leaked into the caller's return type; rendered:\n%s", rendered)
	}
	// The give-away statement from the broken version: a string boxed into an
	// int-typed temp.
	for _, bad := range []string{"var v2 int\n\tv2 = vm.String("} {
		if strings.Contains(rendered, bad) {
			t.Fatalf("string result assigned to an int temp; rendered:\n%s", rendered)
		}
	}
}

// A float hint takes the same path, so the mechanism is not int-specific.
func TestDeclaredFloatReturnType(t *testing.T) {
	ensureLoader()

	rendered := lowerForms(t, "drpkg",
		`(defn drdec ^double [^double x] (if (< x 1.0) x (drdec (- x 1.0))))`)

	if !strings.Contains(rendered, "func Drdec(ec *vm.ExecContext, arg0 float64) (float64, error)") {
		t.Fatalf("expected a natively-typed float64 return; rendered:\n%s", rendered)
	}
}
