/*
 * Copyright (c) 2026 Marcin Gasperowicz <xnooga@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package ir_test

// #598: `(defn f ^long [n] …)` — the idiomatic return-hint position — reads as
// `(defn f (with-meta [n] {:tag …}) …)`. The lowering pipeline split single- from
// multi-arity on a bare `vector?` test, so the hinted form fell through to the
// multi-arity arm, threw, and was demoted to a warning: the fn silently vanished
// from the lowered package and the program ran on the VM instead. These tests
// assert the hinted form lowers, that the spellings which already worked still
// do, and that real multi-arity is still routed as multi-arity.

import (
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

// Intern the defn before lowering it: build resolves same-ns call targets
// through the namespace, so a self-recursive fn is an unresolved symbol
// otherwise. (lisp_defmulti_native_test.go does the same for its multimethod.)
func lowerForms(t *testing.T, pkg string, forms string) string {
	t.Helper()
	runLispExpr(t, `(do `+forms+`)`)
	v := runLispExpr(t, `(ir.passes.pipeline/lower-ns-to-go "`+pkg+`" (quote core) (quote [`+forms+`]))`)
	s, ok := v.(vm.String)
	if !ok {
		t.Fatalf("expected rendered Go string, got %T: %v", v, v)
	}
	return string(s)
}

func TestReturnHintedArityLowers(t *testing.T) {
	ensureLoader()

	rendered := lowerForms(t, "rhpkg",
		`(defn rhfib ^long [^long n] (if (< n 2) n (+ (rhfib (- n 1)) (rhfib (- n 2)))))`)

	if !strings.Contains(rendered, "func Rhfib(") {
		t.Fatalf("return-hinted defn did not lower to a native fn; rendered:\n%s", rendered)
	}
	// The param hint still has to survive the unwrap — a lowered fn taking
	// vm.Value would mean the arity vector reached build-fn stripped of its
	// element metadata.
	if !strings.Contains(rendered, "arg0 int") {
		t.Errorf("expected an unboxed int param on the lowered fn; rendered:\n%s", rendered)
	}
}

// The spellings the issue reported as already working. They route through the
// same branch, so a fix that only special-cased the broken one would show up here.
func TestUnhintedAndAlternateHintSpellingsStillLower(t *testing.T) {
	ensureLoader()

	cases := []struct {
		name  string
		form  string
		wantF string
	}{
		{"no hint", `(defn rhplain [^long n] (+ n 1))`, "func Rhplain("},
		{"param hint only", `(defn rhparam [^long n] (* n 2))`, "func Rhparam("},
		{"hint on the name", `(defn ^long rhname [^long n] (- n 1))`, "func Rhname("},
		// Wrapped in an arity list, so it takes the multi-arity path and gets
		// that path's per-arity name (`_1`) rather than a bare `Rhwrapped`.
		{"hint on a wrapped arity vector", `(defn rhwrapped (^long [^long n] (+ n 3)))`, "func Rhwrapped_1("},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rendered := lowerForms(t, "rhpkg", tc.form)
			if !strings.Contains(rendered, tc.wantF) {
				t.Fatalf("%s did not lower; rendered:\n%s", tc.name, rendered)
			}
		})
	}
}

// Guard the other direction: widening the single-arity predicate must not
// swallow a genuine multi-arity defn, whose first form is a list too.
func TestRealMultiArityStillRoutesAsMultiArity(t *testing.T) {
	ensureLoader()

	rendered := lowerForms(t, "rhpkg", `(defn rhmulti ([a] a) ([a b] (+ a b)))`)

	// Multi-arity lowers to per-arity fns, never to a single `func Rhmulti(`
	// carrying one signature — that would mean it took the single-arity arm and
	// dropped an arity on the floor.
	if strings.Contains(rendered, "func Rhmulti(ec *vm.ExecContext, arg0 vm.Value) (vm.Value, error)") {
		t.Fatalf("multi-arity defn collapsed to a single arity; rendered:\n%s", rendered)
	}
	if !strings.Contains(rendered, "Rhmulti") {
		t.Fatalf("multi-arity defn vanished from the lowered package; rendered:\n%s", rendered)
	}
}
