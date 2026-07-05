/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package ir_test

// A small typed-bool sibling is now INLINED into its caller (EPIC-013 combinator
// inlining), superseding the earlier typed-direct-call path: the caller carries
// callee's body (rt.EqValue) directly, with no call or trampoline to callee.
// Direct-call for inline-INELIGIBLE callees is still covered by
// TestLowerNsSeedsNativeRegistry (native_direct_call_test) and the crosspkg tests.

import (
	"regexp"
	"testing"
)

func TestLowerNsInlinesTypedSibling(t *testing.T) {
	ensureLoader()
	// callee returns a bool ((= x 1) → :bool); it is small + non-recursive, so
	// the intra-ns call in caller is inlined.
	v := runLispExpr(t,
		`(do (create-ns (quote typedres))
		     (intern (quote typedres) (quote callee))
		     (intern (quote typedres) (quote caller))
		     (binding [ir.passes.inline/*enable-inline* true]
		       (ir.passes.pipeline/lower-ns-to-go "typedres" (quote typedres)
		         [(quote (defn callee [x] (= x 1)))
		          (quote (defn caller [y] (callee y)))])))`)
	src := v.String()

	// callee still lowers as its own bool-returning fn (public → PascalCase per #371).
	if !regexp.MustCompile(`func Callee\(ec \*vm\.ExecContext, [a-z0-9_]+ vm\.Value\) bool`).MatchString(src) {
		t.Fatalf("expected Callee to lower with a bool result:\n%s", src)
	}
	// caller lowers with a bool result and carries the INLINED body (inline ON).
	if !regexp.MustCompile(`func Caller\(ec \*vm\.ExecContext, [a-z0-9_]+ vm\.Value\) bool`).MatchString(src) {
		t.Fatalf("expected Caller to lower with a bool result:\n%s", src)
	}
	if !regexp.MustCompile(`func Caller\(ec[^}]*rt\.EqValue`).MatchString(src) {
		t.Fatalf("expected Caller to carry callee's inlined body (rt.EqValue):\n%s", src)
	}
	// No call or trampoline to Callee: `Callee(` appears only in its own definition.
	if n := len(regexp.MustCompile(`Callee\(`).FindAllString(src, -1)); n != 1 {
		t.Fatalf("expected Callee( once (its definition only), got %d — caller did not inline:\n%s", n, src)
	}
	if regexp.MustCompile(`InvokeValue\([^\n]*"callee"|CachedVarFn\([^\n]*"callee"`).MatchString(src) {
		t.Fatalf("caller must NOT trampoline to callee after inlining:\n%s", src)
	}
}
