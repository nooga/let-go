/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package ir_test

import (
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

// EPIC-012 / ITER-0025, T3: a direct-callable lowered fn gets an exported thin
// wrapper (LG_<go-name>) so another lowered Go package can call it, WITHOUT
// renaming the unexported fn (intra-package call sites stay byte-stable). The
// wrapper forwards ec + args positionally to the unexported fn. Emission is
// gated by *emit-exported-wrappers* (default off, so the committed tree stays
// byte-identical until the whole-program collector wires real cross-package
// calls).
func TestExportedWrapperEmittedWhenEnabled(t *testing.T) {
	ensureLoader()

	withFlag := runLispExpr(t, `(do (create-ns (quote watest))
      (binding [ir.passes.pipeline/*emit-exported-wrappers* true]
        (ir.passes.pipeline/lower-ns-to-go "watest" (quote watest) [(quote (defn add1 [x] (+ x 1)))])))`)
	s, ok := withFlag.(vm.String)
	if !ok {
		t.Fatalf("expected rendered Go source string, got %T", withFlag)
	}
	src := string(s)
	if !strings.Contains(src, "func LG_add1(") {
		t.Fatalf("expected exported wrapper `func LG_add1(`\n--- go ---\n%s", src)
	}
	if !strings.Contains(src, "return add1(ec") {
		t.Fatalf("expected wrapper to forward to the unexported `add1(ec, …)`\n--- go ---\n%s", src)
	}
}

// EPIC-012 / ITER-0025, SCENARIO-0017 (hermetic core): the whole-program
// collector resolves a call into another lowered package to a direct
// `pkg.LG_<go>(ec, …)` call (no cached-var trampoline) and adds the import.
// aaa/pick is vm.Value-uniform so coercion is trivial; bbb/use-it calls it.
func TestCrossPackageDirectCallEmitted(t *testing.T) {
	ensureLoader()

	// Define the fns in their namespaces first (interns aaa/pick so bbb's
	// (aaa/pick x) resolves), then lower the whole program. runLispExpr compiles
	// multiple top-level forms; the last form's value is returned.
	v := runLispExpr(t, `
      (ns aaa)
      (defn pick [x] x)
      (ns bbb)
      (defn use-it [x] (aaa/pick x))
      (ns user)
      (nth (ir.passes.pipeline/lower-all-ns-to-go
             [["aaa" (quote aaa) [(quote (defn pick [x] x))] "ex/aaa"]
              ["bbb" (quote bbb) [(quote (defn use-it [x] (aaa/pick x)))] "ex/bbb"]]) 1)`)
	s, ok := v.(vm.String)
	if !ok {
		t.Fatalf("expected bbb source string, got %T", v)
	}
	src := string(s)
	if !strings.Contains(src, "aaa.LG_pick(ec,") {
		t.Fatalf("expected cross-package direct call `aaa.LG_pick(ec, …)`\n--- go ---\n%s", src)
	}
	if strings.Contains(src, "CachedVarFn(&__v_aaa_pick") {
		t.Fatalf("cross-package call must NOT trampoline\n--- go ---\n%s", src)
	}
}

// The default (flag off) must NOT emit the wrapper — this is the byte-identity
// guard that keeps the bootstrap codegen green until the collector lands.
func TestExportedWrapperSuppressedByDefault(t *testing.T) {
	ensureLoader()

	dflt := runLispExpr(t, `(do (create-ns (quote watest))
      (ir.passes.pipeline/lower-ns-to-go "watest" (quote watest) [(quote (defn add1 [x] (+ x 1)))]))`)
	s, ok := dflt.(vm.String)
	if !ok {
		t.Fatalf("expected rendered Go source string, got %T", dflt)
	}
	if strings.Contains(string(s), "LG_add1") {
		t.Fatalf("flag OFF must not emit an exported wrapper\n--- go ---\n%s", string(s))
	}
}
