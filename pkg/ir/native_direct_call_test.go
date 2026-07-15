/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package ir_test

// Tests for cross-namespace direct calls into native modules (lg-dq58).
//
// main already has the typed direct-call emission (emit-typed-direct-call)
// and cross-ns resolution (resolve-call-entry handles :native? entries +
// records the import). The missing piece is seeding *lowered-registry* with
// the native modules exposed via (rt/native-modules), so a call to a native
// clojure.core fn (e.g. seq -> corefns.Seq) resolves to a direct Go call
// instead of an rt.CachedVarFn / rt.InvokeValue trampoline.

import (
	"fmt"
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/rt"
	// Link corefns so its init() registers clojure.core native modules
	// (e.g. seq -> corefns.Seq) into (rt/native-modules). The real lowering
	// path (lgbgen) links this transitively; the unit test must do so too.
	_ "github.com/nooga/let-go/pkg/rt/corefns"
	"github.com/nooga/let-go/pkg/vm"
)

// lowerWithNativeRegistry lowers the IR fn bound to varName with
// *lowered-registry* seeded from (rt/native-modules), mirroring what
// pipeline.lg does for a real package lowering.
func lowerWithNativeRegistry(t *testing.T, varName string) *vm.PersistentMap {
	t.Helper()
	v := runLispExpr(t, fmt.Sprintf(
		`(binding [ir.lower-go/*lowered-registry* (ir.lower-go/native-registry (rt/native-modules))]
		   (ir.lower-go/lower %s :strict))`, varName))
	m, ok := v.(*vm.PersistentMap)
	if !ok {
		t.Fatalf("expected lower to return a map, got %T", v)
	}
	return m
}

func TestNativeModuleSeedsCrossNsDirectCall(t *testing.T) {
	ensureLoader()
	fn := buildLispIR(t, `(defn use-seq [x] (seq x))`)
	optimizeLispIR(t, fn)
	passVarCounter++
	varName := fmt.Sprintf("*native-dc-fn-%d*", passVarCounter)
	rt.NS(rt.NameCoreNS).Def(varName, fn)

	result := lowerWithNativeRegistry(t, varName)
	rendered := bindAndRenderGoDecl(t, result)

	if !strings.Contains(rendered, "corefns.Seq(") {
		t.Fatalf("expected cross-ns native direct call corefns.Seq(...):\n--- go ---\n%s", rendered)
	}
	// A native call site is root-guarded: the static call runs while the
	// primitive var roots are intact, with the var-dispatch trampoline as the
	// else-branch so with-redefs/alter-var-root in a caller's dynamic extent
	// stays observable. Both branches must be present.
	if !strings.Contains(rendered, "rt.NativePrimsIntact()") {
		t.Fatalf("expected the native direct call to be root-guarded:\n--- go ---\n%s", rendered)
	}
	if !strings.Contains(rendered, "InvokeValue") {
		t.Fatalf("expected a trampoline fallback branch beside the direct call:\n--- go ---\n%s", rendered)
	}
}

// TestLowerNsSeedsNativeRegistry checks the production path: lower-ns-to-go
// itself seeds *lowered-registry* from (rt/native-modules), so a whole-ns
// lowering emits native direct calls (corefns.Seq) without the caller having
// to bind the registry manually (that's what slice 1 did). This is what makes
// the generated core_go_lowered tree carry native direct calls.
func TestLowerNsSeedsNativeRegistry(t *testing.T) {
	ensureLoader()
	rendered := runLispString(t,
		`(do (create-ns (quote nativeseedns))
		     (intern (quote nativeseedns) (quote use-seq))
		     (ir.passes.pipeline/lower-ns-to-go "nativeseedns" (quote nativeseedns)
		       [(quote (defn use-seq [x] (seq x)))]))`)

	if !strings.Contains(rendered, "corefns.Seq(") {
		t.Fatalf("lower-ns-to-go did not seed native registry; no corefns.Seq(...):\n--- go ---\n%s", rendered)
	}
}

func runLispString(t *testing.T, expr string) string {
	t.Helper()
	v := runLispExpr(t, expr)
	s, ok := v.(vm.String)
	if !ok {
		t.Fatalf("expected lower-ns-to-go to return a string, got %T", v)
	}
	return string(s)
}

// TestBuildCallEvaluatesCalleeFirst guards call evaluation order: the bytecode
// compiler evaluates the CALLEE before the arguments, so the gogen lowering must
// too. build-call builds the head first, then build-call-with-head (which builds
// the args once and threads the head across argument control-flow joins). If it
// instead built arguments first, an argument carrying a side effect would run
// before the callee — observable here as the ARGMARK println emitted ahead of the
// HEADMARK println in the lowered Go.
func TestBuildCallEvaluatesCalleeFirst(t *testing.T) {
	ensureLoader()
	rendered := runLispString(t,
		`(do (create-ns (quote evordns))
		     (intern (quote evordns) (quote probe))
		     (ir.passes.pipeline/lower-ns-to-go "evordns" (quote evordns)
		       [(quote (defn probe [x]
		                 ((do (println :HEADMARK) identity)
		                  (if x (do (println :ARGMARK) 1) 2))))]))`)

	head := strings.Index(rendered, `vm.Keyword("HEADMARK")`)
	arg := strings.Index(rendered, `vm.Keyword("ARGMARK")`)
	if head < 0 || arg < 0 {
		t.Fatalf("expected both HEADMARK and ARGMARK in lowered Go (head=%d arg=%d):\n%s", head, arg, rendered)
	}
	if head > arg {
		t.Fatalf("callee must be evaluated before arguments: HEADMARK (%d) emitted after ARGMARK (%d):\n%s", head, arg, rendered)
	}
}

// TestWithRedefsDisablesNativeDirect guards the Var override seam: a body that
// rebinds a core var via with-redefs must NOT lower calls to that var into a
// baked native-direct call (corefns.Count), because with-redefs mutates the
// var root at runtime and the native call would ignore it. The call must fall
// back to the cached-var / InvokeValue trampoline, which re-reads the root.
func TestWithRedefsDisablesNativeDirect(t *testing.T) {
	ensureLoader()
	rendered := runLispString(t,
		`(do (create-ns (quote withredefseedns))
		     (intern (quote withredefseedns) (quote probe))
		     (ir.passes.pipeline/lower-ns-to-go "withredefseedns" (quote withredefseedns)
		       [(quote (defn probe [x] (with-redefs [count (fn [_] 42)] (count x))))]))`)

	if strings.Contains(rendered, "corefns.Count(") {
		t.Fatalf("with-redefs over count must disable native-direct; found baked corefns.Count(...):\n--- go ---\n%s", rendered)
	}
	// And the call must still happen — through the var-mediated cached-var
	// trampoline, which re-reads count's root each call so the redef is seen.
	if !strings.Contains(rendered, `rt.CachedVarFn(&__v_clojure_core_count, "clojure.core", "count")`) {
		t.Fatalf("expected count to fall back to the cached-var trampoline:\n--- go ---\n%s", rendered)
	}
}

// TestWithRedefsDisablesListIntrinsic is the list-intrinsic analogue of
// TestWithRedefsDisablesNativeDirect: the (clojure.core/list …) lowering
// intrinsic bakes vm.EmptyList.Cons(…) directly, bypassing the Var. Inside a
// (with-redefs [list …] …) body — which rebinds the var root at runtime — the
// intrinsic must stand down and dispatch through the var, or the rebinding (and
// alter-var-root / intern) is silently ignored.
func TestWithRedefsDisablesListIntrinsic(t *testing.T) {
	ensureLoader()
	rendered := runLispString(t,
		`(do (create-ns (quote withredeflistns))
		     (intern (quote withredeflistns) (quote probe))
		     (ir.passes.pipeline/lower-ns-to-go "withredeflistns" (quote withredeflistns)
		       [(quote (defn probe [x] (with-redefs [list (fn [& _] :overridden)] (list x))))]))`)

	if strings.Contains(rendered, "vm.EmptyList.Cons(") {
		t.Fatalf("with-redefs over list must disable the list intrinsic; found baked vm.EmptyList.Cons(...):\n--- go ---\n%s", rendered)
	}
	if !strings.Contains(rendered, `rt.CachedVarFn(&__v_clojure_core_list, "clojure.core", "list")`) {
		t.Fatalf("expected list to fall back to the cached-var trampoline:\n--- go ---\n%s", rendered)
	}
}
