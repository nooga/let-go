/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package ir_test

// Emission half of typed direct calls: a lowered sibling with a TYPED signature
// (here a bool result) is now recorded in the direct-call registry with its real
// specs (override-coercible?, not just override-uniform-value?), so a caller
// emits a direct typed call (boxing the typed result at the site) instead of the
// rt.CachedVarFn / InvokeValue trampoline. Previously such fns were excluded.

import (
	"regexp"
	"testing"
)

func TestLowerNsTypedResultSiblingDirectCall(t *testing.T) {
	ensureLoader()
	// callee returns a bool ((= x 1) → :bool); caller calls it intra-ns.
	v := runLispExpr(t,
		`(do (create-ns (quote typedres))
		     (intern (quote typedres) (quote callee))
		     (intern (quote typedres) (quote caller))
		     (ir.passes.pipeline/lower-ns-to-go "typedres" (quote typedres)
		       [(quote (defn callee [x] (= x 1)))
		        (quote (defn caller [y] (callee y)))]))`)
	src := v.String()

	// callee lowers with a typed bool result (= is total → no error result).
	if !regexp.MustCompile(`func Callee\(ec \*vm\.ExecContext, [a-z0-9_]+ vm\.Value\) bool`).MatchString(src) {
		t.Fatalf("expected Callee to lower with a bool result:\n%s", src)
	}
	// caller emits a DIRECT call to Callee, boxing the typed bool result.
	if !regexp.MustCompile(`vm\.Boolean\(Callee\(ec,`).MatchString(src) {
		t.Fatalf("expected caller to emit a direct vm.Boolean(Callee(ec, ...)) call:\n%s", src)
	}
	if regexp.MustCompile(`InvokeValue\([^\n]*"callee"|CachedVarFn\([^\n]*"callee"`).MatchString(src) {
		t.Fatalf("caller must NOT trampoline to callee:\n%s", src)
	}
}
