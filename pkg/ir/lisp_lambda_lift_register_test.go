/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package ir_test

import (
	"regexp"
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

// A lambda-lifted sibling is emitted as a native Go func in the SAME package as
// its outer fn, and referenced from the outer fn's native body via
// rt.LookupVar(<ns>, "<name>__lifted0"). For that lookup to resolve at runtime
// (rather than nil-deref), the sibling's var must be interned — i.e. it must
// appear in the package init()'s rt.RegisterGoOverrides map. Before the fix, the
// go-target lambda-lift path stripped each sibling's result to {:fn decl},
// dropping :fn-name/:override-eligible?, so override-entries never registered it.
func TestLambdaLiftSiblingIsRegisteredNatively(t *testing.T) {
	ensureLoader()

	v := runLispExpr(t, `(do (create-ns (quote liftreg))
      (ir.passes.pipeline/lower-ns-to-go "liftreg" (quote liftreg)
        [(quote (defn cx [p] (identity (fn* [q] q))))]))`)
	s, ok := v.(vm.String)
	if !ok {
		t.Fatalf("expected rendered Go source string, got %T", v)
	}
	src := string(s)

	// Sanity: the sibling decl and the var reference both exist (Part 1).
	if !strings.Contains(src, "func Cx__lifted0(") {
		t.Fatalf("expected the lifted sibling decl `func Cx__lifted0(`:\n%s", src)
	}
	if !strings.Contains(src, `LookupVar("liftreg", "cx__lifted0")`) {
		t.Fatalf("expected the outer body to reference LookupVar(\"liftreg\", \"cx__lifted0\"):\n%s", src)
	}

	// Part 2: the sibling's var must be registered in the init() override map so
	// the LookupVar above resolves at runtime. The map key is the bare var name.
	if !strings.Contains(src, "RegisterGoOverrides") {
		t.Fatalf("expected an init() with rt.RegisterGoOverrides:\n%s", src)
	}
	registered := regexp.MustCompile(`"cx__lifted0":\s*__gogen_wrap`)
	if !registered.MatchString(src) {
		t.Fatalf("lifted sibling cx__lifted0 is NOT registered in the override init map:\n%s", src)
	}
}

// A lifted sibling whose body typeinfer would prove TYPED (here: int? -> bool)
// must STILL end up registered. Siblings are only ever invoked boxed through
// their var (higher-order use), so they are lowered with the uniform vm.Value
// ABI — a typed signature would fail override-uniform-value?, silently skip
// registration, and leave the main's rt.LookupVar dangling (the Part-3 parity
// break: 674 nil-derefs from typeinfer's numeric-op-type__lifted0).
func TestLambdaLiftTypedSiblingStillRegisters(t *testing.T) {
	ensureLoader()

	v := runLispExpr(t, `(do (create-ns (quote liftt))
      (ir.passes.pipeline/lower-ns-to-go "liftt" (quote liftt)
        [(quote (defn cx [p] (some (fn* [q] (or (= q :int) (= q :float))) p)))]))`)
	s, ok := v.(vm.String)
	if !ok {
		t.Fatalf("expected rendered Go source string, got %T", v)
	}
	src := string(s)

	// The invariant: EVERY emitted LookupVar of a lifted sibling has a matching
	// registration. Either the sibling registers (preferred) or the whole form
	// fell back to bytecode (no native decls at all) — never a dangling ref.
	if strings.Contains(src, `"cx__lifted0")`) {
		if !regexp.MustCompile(`"cx__lifted0":\s*__gogen_wrap`).MatchString(src) {
			t.Fatalf("main references cx__lifted0 but the sibling is NOT registered (dangling LookupVar):\n%s", src)
		}
		// And the sibling must present the uniform boxed ABI.
		if regexp.MustCompile(`func Cx__lifted0\([^)]*\) \(?bool`).MatchString(src) {
			t.Fatalf("lifted sibling lowered with a typed (bool) ABI — must be uniform vm.Value:\n%s", src)
		}
	}
}
