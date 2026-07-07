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

// BRICK 3 / ITER-0037: a public lowered fn is itself the exported callable.
// Previously wrapped via `LG_<go-name>` forwarding, the fn IS now directly
// exposed with its own PascalCase exported name (e.g., `Add1`). No wrapper
// indirection. The *emit-exported-wrappers* binding is now a legacy no-op.
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
	if !strings.Contains(src, "func Add1(") {
		t.Fatalf("expected exported PascalCase func `func Add1(`\n--- go ---\n%s", src)
	}
	if strings.Contains(src, "LG_add1") {
		t.Fatalf("exported wrapper LG_ prefix should NOT exist anymore\n--- go ---\n%s", src)
	}
}

// BRICK 3 / ITER-0037: the whole-program collector resolves a call into
// another lowered package to a direct `pkg.Pick(ec, …)` call using the public
// fn's OWN PascalCase name (no `LG_` wrapper, no cached-var trampoline).
// aaa/pick is public and vm.Value-uniform so coercion is trivial.
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
	if !strings.Contains(src, "aaa.Pick(ec,") {
		t.Fatalf("expected cross-package direct call `aaa.Pick(ec, …)` using public fn's own name\n--- go ---\n%s", src)
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

// P1 regression: a mutual cross-package call cycle must NOT be lowered into
// reciprocal direct calls — `cyca` importing `cycb` while `cycb` imports `cyca`
// is a Go import cycle (uncompilable). The cycle-closing edges stay on the
// cached-var trampoline.
func TestCrossPackageCyclicEdgesStayTrampolined(t *testing.T) {
	ensureLoader()

	// Setup forms (intern from-a/from-b so the mutual calls resolve) are separate
	// top-level forms; only the final lower-all-ns-to-go call is wrapped in nth.
	setup := `
      (ns cyca)
      (defn from-a [x] x)
      (ns cycb)
      (defn from-b [x] x)
      (ns user)
      `
	call := `(ir.passes.pipeline/lower-all-ns-to-go
        [["cyca" (quote cyca) [(quote (defn from-a [x] (cycb/from-b x)))] "ex/cyca"]
         ["cycb" (quote cycb) [(quote (defn from-b [x] (cyca/from-a x)))] "ex/cycb"]])`

	caV := runLispExpr(t, setup+"(nth "+call+" 0)")
	cbV := runLispExpr(t, setup+"(nth "+call+" 1)")
	ca, ok1 := caV.(vm.String)
	cb, ok2 := cbV.(vm.String)
	if !ok1 || !ok2 {
		t.Fatalf("expected two source strings, got %T and %T", caV, cbV)
	}
	// A cross-package DIRECT call now emits `pkg.<PascalCase>(` (public fns are
	// their own exported name — the `LG_` forwarding prefix was removed). Grep the
	// real direct-call form so a regressed acyclic guard (both directions emitting
	// direct calls → Go import cycle) still trips this test.
	caImportsCb := strings.Contains(string(ca), "cycb.FromB(")
	cbImportsCa := strings.Contains(string(cb), "cyca.FromA(")
	if caImportsCb && cbImportsCa {
		t.Fatalf("mutual cross-package direct calls form a Go import cycle:\n--- cyca ---\n%s\n--- cycb ---\n%s", ca, cb)
	}
}

// P2 regression: a whole-program lowering with NO cross-package references must
// emit NO exported wrappers — an empty collected target set means "export
// nothing", not "export everything" (which would add dead exported API).
func TestWholeProgramNoCrossRefEmitsNoWrappers(t *testing.T) {
	ensureLoader()

	v := runLispExpr(t, `
      (ns solo)
      (defn only-here [x] (+ x 1))
      (ns user)
      (nth (ir.passes.pipeline/lower-all-ns-to-go
             [["solo" (quote solo) [(quote (defn only-here [x] (+ x 1)))] "ex/solo"]]) 0)`)
	s, ok := v.(vm.String)
	if !ok {
		t.Fatalf("expected solo source string, got %T", v)
	}
	if strings.Contains(string(s), "func LG_") {
		t.Fatalf("no fn is called cross-package, so no exported wrapper should be emitted:\n--- go ---\n%s", string(s))
	}
}
