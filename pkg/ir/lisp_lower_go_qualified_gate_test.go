/*
 * Copyright (c) 2026 Matt Parrett
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package ir_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

// #1045: syntax-quote qualifies every symbol a macro emits, so dotimes and
// friends hand the IR core/< and core/inc (clojure.core is an alias of the
// core ns). Those heads must get the same builtin ops as the unqualified
// ones, or the loop's locals fall back to vm.Value and every op becomes a var
// call.

// buildLispIRInNS builds src in the namespace nsName after evaluating setup
// there, so a test can define same-named vars or aliases the gate must see.
func buildLispIRInNS(t *testing.T, nsName, setup, src string) vm.Value {
	t.Helper()
	c := compiler.NewCompiler(vm.NewConsts(), rt.NS(nsName))
	c.SetSource("qualified-gate-build")
	_, result, err := c.CompileMultiple(strings.NewReader(
		fmt.Sprintf("%s\n(ir.build/build-fn (quote %s))", setup, src)))
	if err != nil {
		t.Fatalf("build in %s: %v", nsName, err)
	}
	return result
}

func lowerBuiltFn(t *testing.T, fn vm.Value) string {
	t.Helper()
	optimizeLispIR(t, fn)
	result := lowerGo(t, fn, ":strict")
	if got := result.ValueAt(vm.Keyword("status")); got != vm.Keyword("lowered") {
		t.Fatalf("expected :lowered status, got %v (reason=%v)", got, result.ValueAt(vm.Keyword("reason")))
	}
	return bindAndRenderGoDecl(t, result)
}

// The issue's `qual` repro lowers exactly like `plain`: native int64 locals
// and Go operators, no cached-var trampolines for the core ops.
func TestLowerGoQualifiedCoreHeadsGetBuiltinOps(t *testing.T) {
	ensureLoader()
	rendered := lowerArrayFn(t, `(defn qual [n]
	  (loop [j 0 acc 0]
	    (if (clojure.core/< j n)
	      (recur (clojure.core/inc j) (clojure.core/+ acc j))
	      acc)))`)
	mustContain(t, rendered, "j int64", "acc int64", "j < arg0")
	mustNotContain(t, rendered, "CachedVarFn", "j vm.Value")
}

// core/ is the spelling syntax-quote actually emits.
func TestLowerGoCoreNsHeadsGetBuiltinOps(t *testing.T) {
	ensureLoader()
	rendered := lowerArrayFn(t, `(defn qual [n]
	  (loop [j 0 acc 0]
	    (if (core/< j n) (recur (core/inc j) (core/+ acc j)) acc)))`)
	mustContain(t, rendered, "j int64", "acc int64", "j < arg0")
	mustNotContain(t, rendered, "CachedVarFn", "j vm.Value")
}

// An alias to clojure.core resolves to the same var, so it passes too.
func TestLowerGoAliasedCoreHeadGetsBuiltinOp(t *testing.T) {
	ensureLoader()
	fn := buildLispIRInNS(t, "qgatealias1045", "(alias (quote c) (quote clojure.core))",
		`(defn count-up [n] (loop [j 0] (if (c/< j n) (recur (c/inc j)) j)))`)
	rendered := lowerBuiltFn(t, fn)
	mustContain(t, rendered, "j int64", "j < arg0")
	mustNotContain(t, rendered, "CachedVarFn")
}

// dotimes is the case the issue measured: its counter is native once the
// expansion's qualified < and inc reach the builtin ops. lower-ns-to-go
// macroexpands, as lg compile does; build-fn on the raw form would not.
func TestLowerGoDotimesCounterIsNative(t *testing.T) {
	ensureLoader()
	rendered := runLispString(t,
		`(do (create-ns (quote qgatedotimes1045))
		     (intern (quote qgatedotimes1045) (quote spin))
		     (ir.passes.pipeline/lower-ns-to-go "qgatedotimes1045" (quote qgatedotimes1045)
		       [(quote (defn spin [n] (dotimes [j n] (identity j)) n))]))`)
	mustContain(t, rendered, "func Spin", "j int64", "j < ")
	mustNotContain(t, rendered, `"clojure.core", "<"`, `"clojure.core", "inc"`, "j vm.Value")
}

// A qualified head that resolves to some other namespace's inc is that var,
// not the core op, and stays a call through it.
func TestLowerGoQualifiedNonCoreHeadStaysCall(t *testing.T) {
	ensureLoader()
	fn := buildLispIRInNS(t, "qgateother1045", "(def inc (fn* [x] x))",
		`(defn f [n] (loop [j 0] (if (< j n) (recur (qgateother1045/inc j)) j)))`)
	rendered := lowerBuiltFn(t, fn)
	mustContain(t, rendered, `"qgateother1045", "inc"`)
	mustNotContain(t, rendered, "j = j + 1")
}

// An alias named core that points at a user ns wins over the core ns, so
// core/inc there is the user's inc and stays a call.
func TestLowerGoCoreAliasToUserNsStaysCall(t *testing.T) {
	ensureLoader()
	buildLispIRInNS(t, "qgateuser1045", "(def inc (fn* [x] x))", `(fn* [] nil)`)
	fn := buildLispIRInNS(t, "qgatealiascore1045", "(alias (quote core) (quote qgateuser1045))",
		`(defn f [n] (loop [j 0] (if (< j n) (recur (core/inc j)) j)))`)
	rendered := lowerBuiltFn(t, fn)
	mustContain(t, rendered, `"qgateuser1045", "inc"`)
	mustNotContain(t, rendered, "j = j + 1")
}
