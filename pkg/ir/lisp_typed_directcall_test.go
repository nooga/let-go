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
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/vm"
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
	vs, ok := v.(vm.String)
	if !ok {
		t.Fatalf("expected rendered Go source string, got %T", v)
	}
	src := string(vs)

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
	// No call to Callee from Caller's body — the actual inlining invariant.
	// Scoped to Caller rather than counting `Callee(` across the file: since
	// #554 a typed (bool) return also registers a boxed override, so the file
	// legitimately contains `vm.Boolean(Callee(ec, args[0]))` in the init() map.
	callerBody := regexp.MustCompile(`(?s)func Caller\(ec.*?\n\}`).FindString(src)
	if callerBody == "" {
		t.Fatalf("could not isolate Caller's body:\n%s", src)
	}
	if strings.Contains(callerBody, "Callee(") {
		t.Fatalf("Caller still calls Callee — it did not inline:\n%s", callerBody)
	}
	if regexp.MustCompile(`InvokeValue\([^\n]*"callee"|CachedVarFn\([^\n]*"callee"`).MatchString(src) {
		t.Fatalf("caller must NOT trampoline to callee after inlining:\n%s", src)
	}
	// Both fns return bool, so both register a boxed override (#554). Before
	// that they lowered and then registered nothing, leaving the native code
	// unreachable from a dynamic call.
	for _, name := range []string{"callee", "caller"} {
		if !regexp.MustCompile(`"` + name + `":\s*__gogen_wrap`).MatchString(src) {
			t.Fatalf("expected %q to register a boxed override (#554):\n%s", name, src)
		}
	}
}
