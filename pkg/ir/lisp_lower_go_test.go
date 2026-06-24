/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package ir_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

func optimizeLispIR(t *testing.T, f vm.Value) vm.Value {
	t.Helper()
	passVarCounter++
	varName := fmt.Sprintf("*lower-go-opt-%d*", passVarCounter)
	rt.NS(rt.NameCoreNS).Def(varName, f)
	runLispExpr(t, fmt.Sprintf(`(ir.passes.pipeline/optimize-fn %s)`, varName))
	return f
}

func lowerGo(t *testing.T, f vm.Value, mode string) *vm.PersistentMap {
	t.Helper()
	passVarCounter++
	varName := fmt.Sprintf("*lower-go-fn-%d*", passVarCounter)
	rt.NS(rt.NameCoreNS).Def(varName, f)
	v := runLispExpr(t, fmt.Sprintf(`(ir.lower-go/lower %s %s)`, varName, mode))
	m, ok := v.(*vm.PersistentMap)
	if !ok {
		t.Fatalf("expected lower-go result map, got %T", v)
	}
	return m
}

func renderGoDecl(t *testing.T, result *vm.PersistentMap) string {
	t.Helper()
	decl := result.ValueAt(vm.Keyword("decl"))
	if decl == vm.NIL {
		t.Fatalf("lower-go result missing :decl")
	}
	rendered := runLispExpr(t, fmt.Sprintf(`(gogen/render *lower-go-decl-%d*)`, passVarCounter))
	if s, ok := rendered.(vm.String); ok {
		return string(s)
	}
	t.Fatalf("expected gogen/render to return string, got %T", rendered)
	return ""
}

func runLispExprErr(expr string) error {
	consts := vm.NewConsts()
	c := compiler.NewCompiler(consts, rt.NS(rt.NameCoreNS))
	c.SetSource("lisp-lower-go-error")
	_, _, err := c.CompileMultiple(strings.NewReader(expr))
	return err
}

func bindAndRenderGoDecl(t *testing.T, result *vm.PersistentMap) string {
	t.Helper()
	passVarCounter++
	varName := fmt.Sprintf("*lower-go-decl-%d*", passVarCounter)
	rt.NS(rt.NameCoreNS).Def(varName, result.ValueAt(vm.Keyword("decl")))
	rendered := runLispExpr(t, fmt.Sprintf(`(gogen/render %s)`, varName))
	s, ok := rendered.(vm.String)
	if !ok {
		t.Fatalf("expected gogen/render string, got %T", rendered)
	}
	return string(s)
}

func bindAndRenderGoFile(t *testing.T, file vm.Value) string {
	t.Helper()
	passVarCounter++
	varName := fmt.Sprintf("*lower-go-file-%d*", passVarCounter)
	rt.NS(rt.NameCoreNS).Def(varName, file)
	rendered := runLispExpr(t, fmt.Sprintf(`(gogen/render %s)`, varName))
	s, ok := rendered.(vm.String)
	if !ok {
		t.Fatalf("expected gogen/render string, got %T", rendered)
	}
	return string(s)
}

func TestLowerGoStrictArithmeticLowersToFuncDeclAST(t *testing.T) {
	ensureLoader()

	fn := buildLispIR(t, `(defn add [x y] (+ x y))`)
	seedArgTypes(t, fn, "[:int :int]")
	optimizeLispIR(t, fn)
	result := lowerGo(t, fn, ":strict")

	if got := result.ValueAt(vm.Keyword("status")); got != vm.Keyword("lowered") {
		t.Fatalf("expected :lowered status, got %v", got)
	}

	rendered := bindAndRenderGoDecl(t, result)
	if !strings.Contains(rendered, "func add(") {
		t.Fatalf("expected rendered Go func decl\n--- go ---\n%s", rendered)
	}
	if !strings.Contains(rendered, "int") || !strings.Contains(rendered, "return") || !strings.Contains(rendered, "+") {
		t.Fatalf("expected typed arithmetic lowering\n--- go ---\n%s", rendered)
	}
}

// Regression for PR #235 review (precision): mixed int/float arithmetic obeys
// numeric contagion — (+ 1 2.0) is a Float at runtime, so its result type is
// :float, lowering to a native float64. Previously typeinfer joined :int and
// :float to the ambiguous :number, which has no native Go type and made strict
// lower-go throw "unsupported result type".
func TestLowerGoStrictMixedIntFloatInfersNativeFloat(t *testing.T) {
	ensureLoader()

	fn := buildLispIR(t, `(defn mixed [] (+ 1 2.0))`)
	optimizeLispIR(t, fn)
	result := lowerGo(t, fn, ":strict")

	if got := result.ValueAt(vm.Keyword("status")); got != vm.Keyword("lowered") {
		t.Fatalf("expected :lowered status for a mixed int/float result, got %v (reason=%v)",
			got, result.ValueAt(vm.Keyword("reason")))
	}
	rendered := bindAndRenderGoDecl(t, result)
	if !strings.Contains(rendered, "float64") {
		t.Fatalf("expected (+ 1 2.0) to lower with a native float64 result (numeric contagion), got:\n%s", rendered)
	}
}

// Regression for PR #235 review (robustness): a genuinely ambiguous :number
// result — e.g. (+ x x) where x is only known to be {int,float} — has no native
// Go type, so go-type-spec must give it a lowering target (vm.Value) instead of
// crashing strict lower-go. This is the backstop for honest :number results
// that contagion can't narrow.
func TestLowerGoStrictAmbiguousNumberResultBoxesToValue(t *testing.T) {
	ensureLoader()

	fn := buildLispIR(t, `(defn addn [x] (+ x x))`)
	seedArgTypes(t, fn, "[:number]")
	optimizeLispIR(t, fn)
	result := lowerGo(t, fn, ":strict")

	if got := result.ValueAt(vm.Keyword("status")); got != vm.Keyword("lowered") {
		t.Fatalf("expected :lowered status for an ambiguous :number result, got %v (reason=%v)",
			got, result.ValueAt(vm.Keyword("reason")))
	}
}

// Regression for PR #235: quot is now a first-class IR op. int/int quot lowers
// to native Go `/` (Go integer division truncates toward zero exactly like
// clojure.core/quot) and is typed :int.
func TestLowerGoStrictQuotIntIntLowersNative(t *testing.T) {
	ensureLoader()

	fn := buildLispIR(t, `(defn qii [x y] (quot x y))`)
	seedArgTypes(t, fn, "[:int :int]")
	optimizeLispIR(t, fn)
	result := lowerGo(t, fn, ":strict")

	if got := result.ValueAt(vm.Keyword("status")); got != vm.Keyword("lowered") {
		t.Fatalf("expected :lowered status for int/int quot, got %v (reason=%v)",
			got, result.ValueAt(vm.Keyword("reason")))
	}
	rendered := bindAndRenderGoDecl(t, result)
	if !strings.Contains(rendered, "/") || strings.Contains(rendered, "QuotValue") {
		t.Fatalf("expected int/int quot to lower to native Go division, got:\n%s", rendered)
	}
}

// Regression for PR #235: a mixed-operand quot can NOT use native Go `/` (float
// `/` does not truncate; int/float `/` isn't valid Go), so it must route through
// rt.QuotValue, which truncates per the numeric tower at runtime.
func TestLowerGoStrictQuotMixedRoutesThroughQuotValue(t *testing.T) {
	ensureLoader()

	fn := buildLispIR(t, `(defn qif [x] (quot x 2.0))`)
	seedArgTypes(t, fn, "[:int]")
	optimizeLispIR(t, fn)
	result := lowerGo(t, fn, ":strict")

	if got := result.ValueAt(vm.Keyword("status")); got != vm.Keyword("lowered") {
		t.Fatalf("expected :lowered status for mixed quot, got %v (reason=%v)",
			got, result.ValueAt(vm.Keyword("reason")))
	}
	rendered := bindAndRenderGoDecl(t, result)
	if !strings.Contains(rendered, "QuotValue") {
		t.Fatalf("expected mixed int/float quot to route through rt.QuotValue, got:\n%s", rendered)
	}
}

// Regression for PR #235: div (/) is a first-class IR op. float/float division
// has a native Go type (float64) and lowers to native `/` (matches clojure.core//
// on floats).
func TestLowerGoStrictDivFloatFloatLowersNative(t *testing.T) {
	ensureLoader()

	fn := buildLispIR(t, `(defn dff [x y] (/ x y))`)
	seedArgTypes(t, fn, "[:float :float]")
	optimizeLispIR(t, fn)
	result := lowerGo(t, fn, ":strict")

	if got := result.ValueAt(vm.Keyword("status")); got != vm.Keyword("lowered") {
		t.Fatalf("expected :lowered for float/float div, got %v (reason=%v)",
			got, result.ValueAt(vm.Keyword("reason")))
	}
	rendered := bindAndRenderGoDecl(t, result)
	if !strings.Contains(rendered, "/") || strings.Contains(rendered, "DivValue") {
		t.Fatalf("expected float/float div to lower to native Go division, got:\n%s", rendered)
	}
}

// Regression for PR #235: int/int div yields a Ratio (or Int when exact), which
// has no native Go scalar type, so it must route through rt.DivValue (matching
// clojure.core//) rather than native Go `/` (which would be integer division).
func TestLowerGoStrictDivIntIntRoutesThroughDivValue(t *testing.T) {
	ensureLoader()

	fn := buildLispIR(t, `(defn dii [x y] (/ x y))`)
	seedArgTypes(t, fn, "[:int :int]")
	optimizeLispIR(t, fn)
	result := lowerGo(t, fn, ":strict")

	if got := result.ValueAt(vm.Keyword("status")); got != vm.Keyword("lowered") {
		t.Fatalf("expected :lowered for int/int div, got %v (reason=%v)",
			got, result.ValueAt(vm.Keyword("reason")))
	}
	rendered := bindAndRenderGoDecl(t, result)
	if !strings.Contains(rendered, "DivValue") {
		t.Fatalf("expected int/int div (Ratio-producing) to route through rt.DivValue, got:\n%s", rendered)
	}
}

func TestLowerGoFileRendersFullGoFile(t *testing.T) {
	ensureLoader()

	fn := buildLispIR(t, `(defn add [x y] (+ x y))`)
	seedArgTypes(t, fn, "[:int :int]")
	optimizeLispIR(t, fn)
	result := lowerGo(t, fn, ":strict")

	passVarCounter++
	varName := fmt.Sprintf("*lower-go-file-src-%d*", passVarCounter)
	rt.NS(rt.NameCoreNS).Def(varName, result)
	file := runLispExpr(t, fmt.Sprintf(`(ir.lower-go/file "main" [%s])`, varName))
	rendered := bindAndRenderGoFile(t, file)

	if !strings.HasPrefix(rendered, "package main") {
		t.Fatalf("expected package header\n--- go ---\n%s", rendered)
	}
	if !strings.Contains(rendered, "func add(") {
		t.Fatalf("expected func decl in file\n--- go ---\n%s", rendered)
	}
}

func TestLowerGoStrictBranchLowersToIfStmt(t *testing.T) {
	ensureLoader()

	fn := buildLispIR(t, `(defn choose [flag x] (if flag x 0))`)
	seedArgTypes(t, fn, "[:bool :int]")
	optimizeLispIR(t, fn)
	result := lowerGo(t, fn, ":strict")

	if got := result.ValueAt(vm.Keyword("status")); got != vm.Keyword("lowered") {
		t.Fatalf("expected :lowered status, got %v", got)
	}

	rendered := bindAndRenderGoDecl(t, result)
	if !strings.Contains(rendered, "if ") {
		t.Fatalf("expected branch lowering to contain if\n--- go ---\n%s", rendered)
	}
	if !strings.Contains(rendered, "return") {
		t.Fatalf("expected branch lowering to return from each arm\n--- go ---\n%s", rendered)
	}
}

func TestLowerGoStrictCallLowersViaInvokeValue(t *testing.T) {
	ensureLoader()

	fn := buildLispIR(t, `(defn stringify [x] (str x))`)
	optimizeLispIR(t, fn)
	result := lowerGo(t, fn, ":strict")

	if got := result.ValueAt(vm.Keyword("status")); got != vm.Keyword("lowered") {
		t.Fatalf("expected :lowered status, got %v", got)
	}

	rendered := bindAndRenderGoDecl(t, result)
	if !strings.Contains(rendered, "rt.InvokeValue") {
		t.Fatalf("expected generic call lowering via rt.InvokeValue\n--- go ---\n%s", rendered)
	}
	if !strings.Contains(rendered, "error") || !strings.Contains(rendered, "!= nil") {
		t.Fatalf("expected call lowering to thread error handling\n--- go ---\n%s", rendered)
	}

	passVarCounter++
	varName := fmt.Sprintf("*lower-go-call-file-%d*", passVarCounter)
	rt.NS(rt.NameCoreNS).Def(varName, result)
	file := runLispExpr(t, fmt.Sprintf(`(ir.lower-go/file "main" [%s])`, varName))
	fileSrc := bindAndRenderGoFile(t, file)
	if !strings.Contains(fileSrc, `"github.com/nooga/let-go/pkg/rt"`) ||
		!strings.Contains(fileSrc, `"github.com/nooga/let-go/pkg/vm"`) {
		t.Fatalf("expected rt/vm imports in rendered file\n--- go ---\n%s", fileSrc)
	}
}

func TestLowerGoStrictDecUnknownValueUsesRuntimeHelper(t *testing.T) {
	ensureLoader()

	fn := buildLispIR(t, `(defn dec-count [xs] (dec (count xs)))`)
	optimizeLispIR(t, fn)
	result := lowerGo(t, fn, ":strict")

	if got := result.ValueAt(vm.Keyword("status")); got != vm.Keyword("lowered") {
		t.Fatalf("expected :lowered status, got %v", got)
	}

	rendered := bindAndRenderGoDecl(t, result)
	if !strings.Contains(rendered, "rt.SubValue") {
		t.Fatalf("expected vm.Value dec to lower via rt.SubValue\n--- go ---\n%s", rendered)
	}
	if strings.Contains(rendered, " - 1") {
		t.Fatalf("vm.Value dec lowered to native subtraction\n--- go ---\n%s", rendered)
	}
}

func TestLowerGoStrictDynamicFnArgCallLowersViaInvokeValue(t *testing.T) {
	ensureLoader()

	fn := buildLispIR(t, `(defn apply1 [f x] (f x))`)
	optimizeLispIR(t, fn)
	result := lowerGo(t, fn, ":strict")

	if got := result.ValueAt(vm.Keyword("status")); got != vm.Keyword("lowered") {
		t.Fatalf("expected :lowered status, got %v", got)
	}

	rendered := bindAndRenderGoDecl(t, result)
	if !strings.Contains(rendered, "rt.InvokeValue") {
		t.Fatalf("expected dynamic fn arg call lowering via rt.InvokeValue\n--- go ---\n%s", rendered)
	}
}

func TestLowerGoStrictIdentityClosureLowersToWrappedGoClosure(t *testing.T) {
	ensureLoader()

	fn := buildLispIR(t, `(defn make-id [] (fn* [x] x))`)
	optimizeLispIR(t, fn)
	result := lowerGo(t, fn, ":strict")

	if got := result.ValueAt(vm.Keyword("status")); got != vm.Keyword("lowered") {
		t.Fatalf("expected :lowered status, got %v", got)
	}

	rendered := bindAndRenderGoDecl(t, result)
	if !strings.Contains(rendered, "rt.BoxNativeFn") {
		t.Fatalf("expected closure lowering to wrap a Go func literal\n--- go ---\n%s", rendered)
	}
	if !strings.Contains(rendered, "func(arg0 vm.Value) vm.Value") {
		t.Fatalf("expected identity closure to lower to a vm.Value-typed Go closure\n--- go ---\n%s", rendered)
	}
	if !strings.Contains(rendered, "return arg0") {
		t.Fatalf("expected identity closure body to return the inner argument\n--- go ---\n%s", rendered)
	}
}

// TestLowerGoClosureCapturingLoopCarriedBlockParam: a closure defined inside a
// loop, capturing the loop-carried binding, must lower to Go. The optimizer
// threads the closure's fn-template AND its captured loop var through the loop
// header as block-parameters (:block-arg); closure-info must follow that
// block-param lineage back through branch-target-args to the originating
// :const template. Before the fix it hit :else nil and the whole body fell back
// as "unsupported function body shape".
func TestLowerGoClosureCapturingLoopCarriedBlockParam(t *testing.T) {
	ensureLoader()

	fn := buildLispIR(t, `(defn lc [xs]
		(loop [seen #{}]
			(let [more (filter (fn* [x] (not (seen x))) xs)]
				(if (empty? more) seen (recur (into seen more))))))`)
	optimizeLispIR(t, fn)
	result := lowerGo(t, fn, ":bridge")

	if got := result.ValueAt(vm.Keyword("status")); got != vm.Keyword("lowered") {
		reason := result.ValueAt(vm.Keyword("reason"))
		t.Fatalf("expected closure-over-loop-block-param to lower; got status=%v reason=%v", got, reason)
	}
}

func TestLowerGoStrictCapturedClosureUsesOuterGoLocal(t *testing.T) {
	ensureLoader()

	fn := buildLispIR(t, `(defn make-const [x] (fn* [] x))`)
	optimizeLispIR(t, fn)
	result := lowerGo(t, fn, ":strict")

	if got := result.ValueAt(vm.Keyword("status")); got != vm.Keyword("lowered") {
		t.Fatalf("expected :lowered status, got %v", got)
	}

	rendered := bindAndRenderGoDecl(t, result)
	if !strings.Contains(rendered, "rt.BoxNativeFn") {
		t.Fatalf("expected captured closure to lower via rt.BoxNativeFn\n--- go ---\n%s", rendered)
	}
	if !strings.Contains(rendered, "func() vm.Value") {
		t.Fatalf("expected captured closure to lower to a zero-arg Go closure\n--- go ---\n%s", rendered)
	}
	if !strings.Contains(rendered, "return arg0") {
		t.Fatalf("expected captured closure to close over the outer Go local\n--- go ---\n%s", rendered)
	}
}

func TestLowerGoStrictKeywordClosureLowersToVmKeyword(t *testing.T) {
	ensureLoader()

	fn := buildLispIR(t, `(defn make-keyword [] (fn* [] :ok))`)
	optimizeLispIR(t, fn)
	result := lowerGo(t, fn, ":strict")

	if got := result.ValueAt(vm.Keyword("status")); got != vm.Keyword("lowered") {
		t.Fatalf("expected :lowered status, got %v", got)
	}

	rendered := bindAndRenderGoDecl(t, result)
	if !strings.Contains(rendered, `vm.Keyword("ok")`) {
		t.Fatalf("expected keyword literal to lower through vm.Keyword\n--- go ---\n%s", rendered)
	}
	if !strings.Contains(rendered, "rt.BoxNativeFn") {
		t.Fatalf("expected keyword-returning closure to lower via Go closure wrapping\n--- go ---\n%s", rendered)
	}
}

func TestLowerGoStrictMultiArityClosureLowersToNativeMultiArity(t *testing.T) {
	ensureLoader()

	fn := buildLispIR(t, `(defn make-multi [] (fn* ([] :zero) ([x] x)))`)
	optimizeLispIR(t, fn)
	result := lowerGo(t, fn, ":strict")

	if got := result.ValueAt(vm.Keyword("status")); got != vm.Keyword("lowered") {
		t.Fatalf("expected :lowered status, got %v", got)
	}

	rendered := bindAndRenderGoDecl(t, result)
	if !strings.Contains(rendered, "rt.MakeNativeMultiArity") {
		t.Fatalf("expected multi-arity closure to lower via rt.MakeNativeMultiArity\n--- go ---\n%s", rendered)
	}
	if !strings.Contains(rendered, "func() vm.Value") || !strings.Contains(rendered, "func(arg0 vm.Value) vm.Value") {
		t.Fatalf("expected both arity branches to lower as Go closures\n--- go ---\n%s", rendered)
	}
}

func TestLowerGoStrictCapturedMultiArityClosureUsesOuterGoLocals(t *testing.T) {
	ensureLoader()

	fn := buildLispIR(t, `(defn make-multi [x] (fn* ([] x) ([y] y)))`)
	optimizeLispIR(t, fn)
	result := lowerGo(t, fn, ":strict")

	if got := result.ValueAt(vm.Keyword("status")); got != vm.Keyword("lowered") {
		t.Fatalf("expected :lowered status, got %v", got)
	}

	rendered := bindAndRenderGoDecl(t, result)
	if !strings.Contains(rendered, "rt.MakeNativeMultiArity") {
		t.Fatalf("expected captured multi-arity closure to lower via rt.MakeNativeMultiArity\n--- go ---\n%s", rendered)
	}
	if !strings.Contains(rendered, "return arg0") {
		t.Fatalf("expected zero-arity branch to capture the outer Go local\n--- go ---\n%s", rendered)
	}
}

func TestLowerGoStrictLiteralMapLowersToVmPersistentMap(t *testing.T) {
	ensureLoader()

	fn := buildLispIR(t, `(defn literal-map [] {:a 1 :b 2})`)
	optimizeLispIR(t, fn)
	result := lowerGo(t, fn, ":strict")

	if got := result.ValueAt(vm.Keyword("status")); got != vm.Keyword("lowered") {
		t.Fatalf("expected :lowered status, got %v", got)
	}

	rendered := bindAndRenderGoDecl(t, result)
	if !strings.Contains(rendered, "vm.NewPersistentMap") {
		t.Fatalf("expected literal map to lower through vm.NewPersistentMap\n--- go ---\n%s", rendered)
	}
	if !strings.Contains(rendered, `vm.Keyword("a")`) || !strings.Contains(rendered, `vm.Int(1)`) {
		t.Fatalf("expected literal map entries to lower as boxed vm values\n--- go ---\n%s", rendered)
	}
}

func TestLowerGoStrictQuotedListLowersToVmList(t *testing.T) {
	ensureLoader()

	fn := buildLispIR(t, `(defn quoted-list [] '(1 2))`)
	optimizeLispIR(t, fn)
	result := lowerGo(t, fn, ":strict")

	if got := result.ValueAt(vm.Keyword("status")); got != vm.Keyword("lowered") {
		t.Fatalf("expected :lowered status, got %v", got)
	}

	rendered := bindAndRenderGoDecl(t, result)
	if !strings.Contains(rendered, "vm.NewList") {
		t.Fatalf("expected quoted list to lower through vm.NewList\n--- go ---\n%s", rendered)
	}
	if !strings.Contains(rendered, `vm.Int(1)`) || !strings.Contains(rendered, `vm.Int(2)`) {
		t.Fatalf("expected quoted list elements to lower as boxed vm values\n--- go ---\n%s", rendered)
	}
}

func TestLowerGoStrictCharConstLowersToVmChar(t *testing.T) {
	ensureLoader()

	fn := buildLispIR(t, `(defn char-lit [] (quote \a))`)
	optimizeLispIR(t, fn)
	result := lowerGo(t, fn, ":strict")

	if got := result.ValueAt(vm.Keyword("status")); got != vm.Keyword("lowered") {
		t.Fatalf("expected :lowered status, got %v", got)
	}

	rendered := bindAndRenderGoDecl(t, result)
	if !strings.Contains(rendered, "vm.Char('a')") {
		t.Fatalf("expected char literal to lower through vm.Char with a Go rune literal\n--- go ---\n%s", rendered)
	}
}

func TestLowerGoStrictJoinValueFeedsLaterArithmetic(t *testing.T) {
	ensureLoader()

	fn := buildLispIR(t, `(defn maybe-inc [flag x]
	                       (+ (if flag x 0) 1))`)
	seedArgTypes(t, fn, "[:bool :int]")
	optimizeLispIR(t, fn)
	result := lowerGo(t, fn, ":strict")

	if got := result.ValueAt(vm.Keyword("status")); got != vm.Keyword("lowered") {
		t.Fatalf("expected :lowered status, got %v", got)
	}

	rendered := bindAndRenderGoDecl(t, result)
	// Structured emission lowers the value-if to an if/else that assigns the
	// join local, then the continuation consumes it — no goto/labels.
	if !strings.Contains(rendered, "} else {") {
		t.Fatalf("expected join lowering to use a structured if/else\n--- go ---\n%s", rendered)
	}
	if strings.Contains(rendered, "goto ") {
		t.Fatalf("expected structured (goto-free) join lowering\n--- go ---\n%s", rendered)
	}
	// The join value must still feed the later arithmetic. Depending on
	// whether typeinfer narrows the join param to int, this is either a
	// Go-native `+ 1` or rt.AddValue.
	if !strings.Contains(rendered, "+ 1") && !strings.Contains(rendered, "rt.AddValue") {
		t.Fatalf("expected joined value to feed later arithmetic\n--- go ---\n%s", rendered)
	}
}

func TestLowerGoStrictLoopLowersToCFG(t *testing.T) {
	ensureLoader()

	fn := buildLispIR(t, `(defn sum-to-n [n]
	                       (loop* [i 0 acc 0]
	                         (if (< i n)
	                           (recur (+ i 1) (+ acc i))
	                           acc)))`)
	seedArgTypes(t, fn, "[:int]")
	optimizeLispIR(t, fn)
	result := lowerGo(t, fn, ":strict")

	if got := result.ValueAt(vm.Keyword("status")); got != vm.Keyword("lowered") {
		t.Fatalf("expected :lowered status, got %v", got)
	}

	rendered := bindAndRenderGoDecl(t, result)
	// Structured emission lowers loop*/recur to a `for { ... }` with the
	// recur arm as `continue` and the exit arm as `break` — no goto/labels.
	if !strings.Contains(rendered, "for {") {
		t.Fatalf("expected loop lowering to use a structured for-loop\n--- go ---\n%s", rendered)
	}
	if !strings.Contains(rendered, "continue") || !strings.Contains(rendered, "break") {
		t.Fatalf("expected loop lowering to use continue/break for the back-edge and exit\n--- go ---\n%s", rendered)
	}
	if strings.Contains(rendered, "goto ") {
		t.Fatalf("expected structured (goto-free) loop lowering\n--- go ---\n%s", rendered)
	}
	if !strings.Contains(rendered, "if ") || !strings.Contains(rendered, "return") {
		t.Fatalf("expected loop lowering to preserve conditional exit\n--- go ---\n%s", rendered)
	}
}

func TestLowerGoStrictMultiArityDefnLowersToNativeMultiArity(t *testing.T) {
	ensureLoader()

	fn := buildLispIR(t, `(defn foo ([x] x) ([x y] (+ x y)))`)
	optimizeLispIR(t, fn)
	result := lowerGo(t, fn, ":strict")

	if got := result.ValueAt(vm.Keyword("status")); got != vm.Keyword("lowered") {
		t.Fatalf("expected :lowered status, got %v (reason: %v)", got, result.ValueAt(vm.Keyword("reason")))
	}

	rendered := bindAndRenderGoDecl(t, result)
	if !strings.Contains(rendered, "rt.MakeNativeMultiArity") {
		t.Fatalf("expected multi-arity defn to lower via rt.MakeNativeMultiArity\n--- go ---\n%s", rendered)
	}
	if !strings.Contains(rendered, "func(arg0 vm.Value) vm.Value") {
		t.Fatalf("expected single-arg branch to lower as Go closure\n--- go ---\n%s", rendered)
	}
	if !strings.Contains(rendered, "func(arg0 vm.Value, arg1 vm.Value) vm.Value") {
		t.Fatalf("expected two-arg branch to lower as Go closure\n--- go ---\n%s", rendered)
	}
}

func TestLowerGoBridgeLowersUUIDConst(t *testing.T) {
	ensureLoader()

	fn := buildLispIR(t, `(defn a-uuid [] #uuid "123e4567-e89b-12d3-a456-426614174000")`)
	optimizeLispIR(t, fn)
	result := lowerGo(t, fn, ":bridge")

	if got := result.ValueAt(vm.Keyword("status")); got != vm.Keyword("lowered") {
		t.Fatalf("expected :lowered status, got %v", got)
	}
	rendered := bindAndRenderGoDecl(t, result)
	if !strings.Contains(rendered, `vm.ParseUUID("123e4567-e89b-12d3-a456-426614174000")`) {
		t.Fatalf("expected UUID const to lower through vm.ParseUUID with the bare canonical string\n--- go ---\n%s", rendered)
	}
}

func TestLowerGoBridgeLowersBigDecimalConst(t *testing.T) {
	ensureLoader()

	fn := buildLispIR(t, `(defn a-bigdec [] 123.456M)`)
	optimizeLispIR(t, fn)
	result := lowerGo(t, fn, ":bridge")

	if got := result.ValueAt(vm.Keyword("status")); got != vm.Keyword("lowered") {
		t.Fatalf("expected :lowered status, got %v", got)
	}
	rendered := bindAndRenderGoDecl(t, result)
	// str on a BigDecimal drops the M suffix, so the emitted parse string is the
	// plain numeric form that vm.MustBigDecimalFromString round-trips.
	if !strings.Contains(rendered, `vm.MustBigDecimalFromString("123.456")`) {
		t.Fatalf("expected BigDecimal const to lower through vm.MustBigDecimalFromString\n--- go ---\n%s", rendered)
	}
}

func TestLowerGoBridgeLowersRatioConst(t *testing.T) {
	ensureLoader()

	fn := buildLispIR(t, `(defn a-ratio [] 1/2)`)
	optimizeLispIR(t, fn)
	result := lowerGo(t, fn, ":bridge")

	if got := result.ValueAt(vm.Keyword("status")); got != vm.Keyword("lowered") {
		t.Fatalf("expected :lowered status, got %v", got)
	}
	rendered := bindAndRenderGoDecl(t, result)
	if !strings.Contains(rendered, `vm.MustRatioFromString("1/2")`) {
		t.Fatalf("expected Ratio const to lower through vm.MustRatioFromString\n--- go ---\n%s", rendered)
	}
}

func TestLowerGoStrictVecDestructureDefn(t *testing.T) {
	ensureLoader()

	fn := buildLispIR(t, `(defn vec-dest [x [a b]] (+ x a b))`)
	optimizeLispIR(t, fn)
	result := lowerGo(t, fn, ":strict")

	if got := result.ValueAt(vm.Keyword("status")); got != vm.Keyword("lowered") {
		t.Fatalf("expected :lowered status, got %v (reason: %v)", got, result.ValueAt(vm.Keyword("reason")))
	}

	rendered := bindAndRenderGoDecl(t, result)
	if !strings.Contains(rendered, "func vec_") {
		t.Fatalf("expected func decl\n--- go ---\n%s", rendered)
	}
	if !strings.Contains(rendered, "rt.AddValue") && !strings.Contains(rendered, "+") {
		t.Fatalf("expected arithmetic in lowered body\n--- go ---\n%s", rendered)
	}
}

func TestLowerGoStrictMapKeysDestructureDefn(t *testing.T) {
	ensureLoader()

	fn := buildLispIR(t, `(defn map-dest [{:keys [a b]}] (+ a b))`)
	optimizeLispIR(t, fn)
	result := lowerGo(t, fn, ":strict")

	if got := result.ValueAt(vm.Keyword("status")); got != vm.Keyword("lowered") {
		t.Fatalf("expected :lowered status, got %v (reason: %v)", got, result.ValueAt(vm.Keyword("reason")))
	}

	rendered := bindAndRenderGoDecl(t, result)
	// Destructured locals don't carry an inferred numeric type, so + on
	// them routes through rt.AddValue instead of a Go-native `+`.
	if !strings.Contains(rendered, "rt.AddValue") && !strings.Contains(rendered, "+") {
		t.Fatalf("expected arithmetic in lowered body\n--- go ---\n%s", rendered)
	}
}

func TestLowerGoStrictNestedDestructureDefn(t *testing.T) {
	ensureLoader()

	fn := buildLispIR(t, `(defn nested-dest [[a {:keys [b]}]] (+ a b))`)
	optimizeLispIR(t, fn)
	result := lowerGo(t, fn, ":strict")

	if got := result.ValueAt(vm.Keyword("status")); got != vm.Keyword("lowered") {
		t.Fatalf("expected :lowered status, got %v (reason: %v)", got, result.ValueAt(vm.Keyword("reason")))
	}

	rendered := bindAndRenderGoDecl(t, result)
	if !strings.Contains(rendered, "rt.AddValue") && !strings.Contains(rendered, "+") {
		t.Fatalf("expected arithmetic in lowered body\n--- go ---\n%s", rendered)
	}
}

func TestLowerGoNoSelfAssignment(t *testing.T) {
	ensureLoader()
	fn := buildLispIR(t, `(defn test-loop [x y]
	                       (loop [a x b y]
	                         (if (not= a 0)
	                           (recur a b)
	                           b)))`)
	seedArgTypes(t, fn, "[:int :int]")
	optimizeLispIR(t, fn)
	result := lowerGo(t, fn, ":strict")
	rendered := bindAndRenderGoDecl(t, result)
	for _, line := range strings.Split(rendered, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.Contains(trimmed, "=") && !strings.HasPrefix(trimmed, "//") {
			parts := strings.Split(trimmed, "=")
			if len(parts) == 2 {
				lhs := strings.TrimSpace(parts[0])
				rhs := strings.TrimSpace(parts[1])
				if lhs != "" && lhs == rhs {
					t.Fatalf("self-assignment emitted: %q\n--- go ---\n%s", trimmed, rendered)
				}
			}
		}
	}
}

func TestBlockParamsCarrySourceName(t *testing.T) {
	ensureLoader()
	fn := buildLispIR(t, `(defn sum [n] (loop [i 0 acc 0] (if (= i n) acc (recur (+ i 1) (+ acc i)))))`)
	seedArgTypes(t, fn, "[:int]")
	passVarCounter++
	v := fmt.Sprintf("*bp-%d*", passVarCounter)
	rt.NS(rt.NameCoreNS).Def(v, fn)
	out := runLispExpr(t, fmt.Sprintf(`
	  (let [f %s
	        params (mapcat (fn [b] (ir/block-params b f)) (ir/blocks f))
	        named  (filter (fn [p]
	                         (some (fn [si] (ir/source-info-symbol si))
	                               (ir/source-infos p f)))
	                       params)]
	    (pr-str [(count params) (count named)]))`, v))
	s := string(out.(vm.String))
	if strings.HasSuffix(s, " 0]") {
		t.Fatalf("no block params carry a source name: %s", s)
	}
}

func TestLowerGoCoalescesLineageByName(t *testing.T) {
	ensureLoader()
	fn := buildLispIR(t, `(defn sum [n] (loop [i 0 acc 0] (if (= i n) acc (recur (+ i 1) (+ acc i)))))`)
	seedArgTypes(t, fn, "[:int]")
	optimizeLispIR(t, fn)
	rendered := bindAndRenderGoDecl(t, lowerGo(t, fn, ":strict"))
	if regexp.MustCompile(`acc_\d+`).MatchString(rendered) {
		t.Fatalf("expected coalesced `acc`, found versioned acc_NN:\n%s", rendered)
	}
	if !strings.Contains(rendered, "acc") {
		t.Fatalf("expected `acc`:\n%s", rendered)
	}
}

func TestLowerGoReassignsInPlaceNoDiscardNet(t *testing.T) {
	ensureLoader()
	fn := buildLispIR(t, `(defn sum [n] (loop [i 0 acc 0] (if (= i n) acc (recur (+ i 1) (+ acc i)))))`)
	seedArgTypes(t, fn, "[:int]")
	optimizeLispIR(t, fn)
	rendered := bindAndRenderGoDecl(t, lowerGo(t, fn, ":strict"))
	// The `_, _, ... = v3, v6` discard net masks declared-but-unread locals
	// instead of not declaring them. Read-set gating must make it unnecessary.
	if regexp.MustCompile(`(?m)^\s*_(, _)+ = `).MatchString(rendered) {
		t.Fatalf("discard net still present\n--- go ---\n%s", rendered)
	}
}

// writeOnlyLocals returns declared locals (`var NAME T`) that are never read in
// the rendered Go: every occurrence is either the declaration or an assignment
// LHS. Such a local is a dead store — Go would reject it as declared-and-not-used
// (or ineffassign would flag it). callErr is excused (error plumbing).
func writeOnlyLocals(rendered string) []string {
	declRe := regexp.MustCompile(`^\s*var (\w+) `)
	declared := []string{}
	reads := map[string]int{}
	lines := strings.Split(rendered, "\n")
	for _, ln := range lines {
		if m := declRe.FindStringSubmatch(ln); m != nil {
			if m[1] != "callErr" {
				declared = append(declared, m[1])
			}
			continue
		}
		// Assignment? Split on the first " = " (not ":=", not "==").
		rhs := ln
		if idx := strings.Index(ln, " = "); idx >= 0 && !strings.Contains(ln[:idx], ":") {
			rhs = ln[idx+3:]
		}
		for _, name := range declared {
			if regexp.MustCompile(`\b` + regexp.QuoteMeta(name) + `\b`).MatchString(rhs) {
				reads[name]++
			}
		}
	}
	var dead []string
	for _, name := range declared {
		if reads[name] == 0 {
			dead = append(dead, name)
		}
	}
	return dead
}

func TestLowerGoDeadCallResultRoutedToBlank(t *testing.T) {
	ensureLoader()
	// `(vec (reverse path))` is built with the `(reverse path)` subexpression
	// duplicated; one copy is dead. The dead side-effecting call must be
	// discarded via `_, callErr =`, never stored in a declared-but-unread local.
	fn := buildLispIR(t, `(defn f [path] (vec (reverse path)))`)
	seedArgTypes(t, fn, "[:string]")
	optimizeLispIR(t, fn)
	rendered := bindAndRenderGoDecl(t, lowerGo(t, fn, ":strict"))
	if dead := writeOnlyLocals(rendered); len(dead) > 0 {
		t.Fatalf("declared-but-unread locals (dead stores): %v\n--- go ---\n%s", dead, rendered)
	}
}

func TestLowerGoCapturedNameNotShadowedByBlockParam(t *testing.T) {
	ensureLoader()
	// `bad` is captured by the closure AND used inside an if-branch, so it is
	// threaded through a block-param. If that param coalesces to "bad", it
	// shadows the lexical capture and (the init copy being a self-assign) is
	// never assigned — the closure then reads a zero value, and the outer `bad`
	// becomes unused. Captured names must stay versioned so the copy is real.
	fn := buildLispIR(t, `(defn h [xs] (let [bad (count xs)] (map (fn* [s] (if s (contains? bad s) nil)) xs)))`)
	optimizeLispIR(t, fn)
	rendered := bindAndRenderGoDecl(t, lowerGo(t, fn, ":strict"))
	if n := strings.Count(rendered, "var bad vm.Value"); n > 1 {
		t.Fatalf("captured `bad` shadowed by an inner block-param local (%d decls)\n--- go ---\n%s", n, rendered)
	}
}

func TestLowerGoInlinesConstantBlockParams(t *testing.T) {
	ensureLoader()
	// `:k` is carried through a staircase of merge block-params, each holding
	// only the constant. A constant block-param is freely re-materializable
	// (consts are block-free), so it must lower to the literal inline — no local,
	// no `vN = vm.Keyword("k")` forwarding copy. The function just returns :k.
	fn := buildLispIR(t, `(defn f [a b] (if a :k (if b :k :k)))`)
	seedArgTypes(t, fn, "[:any :any]")
	optimizeLispIR(t, fn)
	rendered := bindAndRenderGoDecl(t, lowerGo(t, fn, ":strict"))
	if m := regexp.MustCompile(`(?m)^\s*v\d+ = vm\.Keyword\("k"\)$`).FindString(rendered); m != "" {
		t.Fatalf("constant block-param not inlined: %q\n--- go ---\n%s", strings.TrimSpace(m), rendered)
	}
	if !strings.Contains(rendered, `return vm.Keyword("k")`) {
		t.Fatalf("expected `return vm.Keyword(\"k\")` (inlined const param)\n--- go ---\n%s", rendered)
	}
}

func TestLowerGoLowersVarOperandInSideEffectingCall(t *testing.T) {
	ensureLoader()
	// A `(var *ns*)` operand now lowers to rt.LookupVar(ns,name) (RC2), so a
	// push-binding!/pop-binding! call carrying it lowers cleanly instead of
	// forcing whole-function fallback. The side-effecting calls must still be
	// PRESENT in the rendered Go (never silently dropped) — that invariant is
	// the original point of this test and is preserved by the assertions below.
	fn := buildLispIR(t, `(defn cf [form caller-ns]
	  (do
	    (push-binding! (var *ns*) (first caller-ns))
	    (let [r (do (count form))]
	      (do (pop-binding! (var *ns*)) r))))`)
	optimizeLispIR(t, fn)
	result := lowerGo(t, fn, ":bridge")
	if got := result.ValueAt(vm.Keyword("status")); got != vm.Keyword("lowered") {
		t.Fatalf("expected :lowered now that (var *ns*) lowers, got %v reason=%v",
			got, result.ValueAt(vm.Keyword("reason")))
	}
	rendered := bindAndRenderGoDecl(t, result)
	// The var operand renders via rt.LookupVar; both side-effecting binding
	// calls survive (not dropped).
	for _, want := range []string{
		`rt.LookupVar("core", "*ns*")`,
		`"push-binding!"`,
		`"pop-binding!"`,
	} {
		if !strings.Contains(rendered, want) {
			t.Fatalf("rendered Go missing %q\n%s", want, rendered)
		}
	}
}

// TestLowerGoCharLiteral (ITER-0001): character literals lower to vm.Char, and
// in a primitive-typed (char-returning) context they UNBOX to a vm.Char-typed
// return rather than vm.Value — consistent with int. Before the fix, build-form
// had no char? case and threw "unrecognized form".
func TestLowerGoCharLiteral(t *testing.T) {
	ensureLoader()

	fn := buildLispIR(t, `(defn whitespace? [c] (or (= c \space) (= c \tab)))`)
	optimizeLispIR(t, fn)
	result := lowerGo(t, fn, ":bridge")
	if got := result.ValueAt(vm.Keyword("status")); got != vm.Keyword("lowered") {
		t.Fatalf("char-literal defn should lower; got %v reason=%v", got, result.ValueAt(vm.Keyword("reason")))
	}
	if r := bindAndRenderGoDecl(t, result); !strings.Contains(r, `vm.Char(' ')`) {
		t.Fatalf("expected vm.Char(' ')\n%s", r)
	}

	fn2 := buildLispIR(t, `(defn ch [] \a)`)
	optimizeLispIR(t, fn2)
	result2 := lowerGo(t, fn2, ":bridge")
	r2 := bindAndRenderGoDecl(t, result2)
	if !strings.Contains(r2, "func ch(ec *vm.ExecContext) vm.Char") {
		t.Fatalf("expected char return to UNBOX to vm.Char\n%s", r2)
	}
}

// With a same-ns single-arity sibling registered in *lowered-registry*, an
// intra-ns call lowers to a direct Go call (no InvokeValue / LookupVar).
func TestLowerGoIntraNsCallLiftsToDirect(t *testing.T) {
	ensureLoader()
	// Define a dummy callee function in the core namespace so the symbol resolves.
	runLispExpr(t, `(defn callee [x] x)`)
	fn := buildLispIR(t, `(defn caller [y] (callee y))`)
	optimizeLispIR(t, fn)
	passVarCounter++
	varName := fmt.Sprintf("*lower-go-direct-%d*", passVarCounter)
	rt.NS(rt.NameCoreNS).Def(varName, fn)

	// Get the var reference and set it directly
	lowerGoNS := rt.NS("ir.lower-go")
	if lowerGoNS == nil {
		t.Fatalf("ir.lower-go namespace not found")
	}
	loweredRegistryVar := lowerGoNS.LookupLocal(vm.Symbol("*lowered-registry*"))
	if loweredRegistryVar == nil {
		t.Fatalf("*lowered-registry* var not found")
	}

	// Register `callee` as a same-ns single-arity sibling, matching the shape
	// lower-go/registry-entry-from-result produces. Key is [ns-sym name arity];
	// an unqualified intra-ns call matches by name+arity (ns optional), so the
	// key's ns is not load-bearing here.
	registryKey := vm.NewPersistentVector([]vm.Value{
		vm.Symbol("core"), vm.String("callee"), vm.Int(1),
	})
	registryEntry := vm.NewPersistentMap([]vm.Value{
		vm.Keyword("go-name"), vm.String("callee"),
		vm.Keyword("arity"), vm.Int(1),
		vm.Keyword("needs-error?"), vm.FALSE,
		vm.Keyword("param-specs"), vm.NewPersistentVector([]vm.Value{vm.String("vm.Value")}),
		vm.Keyword("result-spec"), vm.String("vm.Value"),
		vm.Keyword("native?"), vm.FALSE,
		vm.Keyword("go-pkg"), vm.NIL,
	})
	registryMap := vm.NewPersistentMap([]vm.Value{registryKey, registryEntry})
	oldVal := loweredRegistryVar.Deref()
	loweredRegistryVar.SetRoot(registryMap)
	defer loweredRegistryVar.SetRoot(oldVal)

	v := runLispExpr(t, fmt.Sprintf(
		`(gogen/render (:decl (ir.lower-go/lower %s :bridge)))`, varName))
	s, ok := v.(vm.String)
	if !ok {
		t.Fatalf("expected rendered Go string, got %T", v)
	}
	got := string(s)
	if !strings.Contains(got, "callee(") {
		t.Fatalf("expected a direct callee(...) call:\n%s", got)
	}
	if strings.Contains(got, "InvokeValue") {
		t.Fatalf("expected NO InvokeValue for the lifted intra-ns call:\n%s", got)
	}
}

// With *lowered-registry* empty (the default), behavior is unchanged: the call
// stays on the rt.InvokeValue / LookupVar path. Guards the bytecode path.
func TestLowerGoCallDefaultStaysInvokeValue(t *testing.T) {
	ensureLoader()
	// Define a dummy callee function so the symbol resolves.
	runLispExpr(t, `(defn callee [x] x)`)
	fn := buildLispIR(t, `(defn caller [y] (callee y))`)
	optimizeLispIR(t, fn)
	passVarCounter++
	varName := fmt.Sprintf("*lower-go-noopt-%d*", passVarCounter)
	rt.NS(rt.NameCoreNS).Def(varName, fn)
	v := runLispExpr(t, fmt.Sprintf(
		`(gogen/render (:decl (ir.lower-go/lower %s :bridge)))`, varName))
	got := string(v.(vm.String))
	if !strings.Contains(got, "InvokeValue") {
		t.Fatalf("expected InvokeValue when *lowered-registry* is empty:\n%s", got)
	}
}

// End-to-end: lower-ns-to-go discovers the lowered sibling and lifts the
// intra-ns call (caller -> callee) to a direct Go call.
func TestLowerNsToGoLiftsIntraNsCall(t *testing.T) {
	ensureLoader()
	// Intern the siblings into the target ns BEFORE lowering, mirroring the real
	// loader (which defs every fn into its ns before lowering runs). This makes
	// the unqualified `callee` call resolve to directtest/callee — not the
	// clojure.core fallback — so the intra-ns direct-call match fires. Done in a
	// single eval with explicit (intern ns sym) so it doesn't depend on in-ns
	// persisting across separate runLispExpr calls.
	v := runLispExpr(t,
		`(do (create-ns (quote directtest))
		     (intern (quote directtest) (quote callee))
		     (intern (quote directtest) (quote caller))
		     (ir.passes.pipeline/lower-ns-to-go "directtest" (quote directtest)
		       [(quote (defn callee [x] x)) (quote (defn caller [y] (callee y)))]))`)
	got := string(v.(vm.String))
	// With two-pass lowering in place:
	// - pass 1 discovers callee as override-eligible (single-arity, vm.Value uniform)
	// - pass 2 re-lowers with *lowered-registry* bound
	// - caller should emit a direct callee(...) call, NOT InvokeValue+LookupVar
	hasDirectCall := strings.Contains(got, "callee(") && !regexp.MustCompile(`InvokeValue\([^\n]*LookupVar\([^\n]*"callee"`).MatchString(got)
	if !hasDirectCall {
		t.Fatalf("expected direct callee(...) call (no InvokeValue+LookupVar for callee) in:\n%s", got)
	}
}

// A cross-namespace callee (clojure.core/count) is lifted to cached-var IFn
// dispatch — it no longer pays a per-call rt.LookupVar at the call site.
// (Superseded the old "stays on InvokeValue" scope guard: the IFn-dispatch
// slice is exactly what lifts cross-ns calls.)
func TestLowerNsToGoLiftsCrossNsCall(t *testing.T) {
	ensureLoader()
	runLispExpr(t, `(create-ns (quote directtest2))`)
	v := runLispExpr(t,
		`(ir.passes.pipeline/lower-ns-to-go "directtest2" (quote directtest2)
		   [(quote (defn caller2 [y] (count y)))])`)
	src := string(v.(vm.String))
	f := parseLoweredGo(t, src)
	caller, ok := findFunc(f, func(n string) bool { return n == "caller2" })
	if !ok {
		t.Fatalf("no caller2 in:\n%s", src)
	}
	if _, ok := findIFnDispatch(caller.Body); !ok {
		t.Fatalf("expected cross-ns count to lift to cached-var IFn dispatch:\n%s", src)
	}
	if callsLookupVarNamed(caller, "count") {
		t.Fatalf("caller2 must not do a per-call rt.LookupVar for count:\n%s", src)
	}
}

// IFn-dispatch lowering: a cross-ns callee (clojure.core/count) makes the
// lowered package declare a package-level cached `var __v_* *vm.Var`. Resolution
// + memoization happen at the call site via the shared rt.CachedVarFn helper (no
// per-callee accessor func is emitted), and NOT eagerly in init() — blank-
// imported lowered packages run init() before the bundle replay loads
// namespaces, where the lookup would return nil.
//
// Assertions are structural (go/ast), not regex over rendered text.
func TestLowerGoCrossNsEmitsCachedVar(t *testing.T) {
	ensureLoader()
	runLispExpr(t, `(create-ns (quote crossnsvar))`)
	v := runLispExpr(t,
		`(ir.passes.pipeline/lower-ns-to-go "crossnsvar" (quote crossnsvar)
		   [(quote (defn caller3 [y] (count y)))])`)
	src := string(v.(vm.String))
	f := parseLoweredGo(t, src)

	isCountVar := func(n string) bool {
		return strings.HasPrefix(n, "__v_") && strings.HasSuffix(n, "count")
	}

	// 1. package-level cached *vm.Var decl for count.
	name, typ, ok := findPkgVar(f, isCountVar)
	if !ok {
		t.Fatalf("expected package-level cached var for count:\n%s", src)
	}
	if typ != "*vm.Var" {
		t.Fatalf("cached var %s has type %q, want *vm.Var", name, typ)
	}

	// 2. the call site resolves it via the shared rt.CachedVarFn helper, and NO
	//    per-callee accessor func is emitted (that bloat is what rt.CachedVarFn
	//    replaces).
	if _, ok := findIFnDispatch(f); !ok {
		t.Fatalf("expected rt.CachedVarFn IFn dispatch for count:\n%s", src)
	}
	if _, ok := findFunc(f, func(n string) bool { return strings.HasPrefix(n, "__getv_") }); ok {
		t.Fatalf("no per-callee accessor func should be emitted (use rt.CachedVarFn):\n%s", src)
	}

	// 3. resolution is lazy: no init() may eagerly call rt.LookupVar (that would
	//    cache nil before bundle replay loads namespaces).
	for _, in := range initFuncs(f) {
		if callsLookupVarNamed(in, "count") {
			t.Fatalf("cached var must NOT be resolved eagerly in init():\n%s", src)
		}
	}
}

// Task 2 (IFn-dispatch lowering): a cross-ns call lowers its CALL SITE to the
// cached-var IFn dispatch
//
//	rt.CachedVarFn(&__v_<ns>_<name>, "ns", "name").Invoke([]vm.Value{...})
//
// keeping the per-call .Deref() as the override seam, and NO per-call
// rt.LookupVar at the call site.
func TestLowerGoCrossNsCallLiftsToIFn(t *testing.T) {
	ensureLoader()
	runLispExpr(t, `(create-ns (quote crossnsifn))`)
	v := runLispExpr(t,
		`(ir.passes.pipeline/lower-ns-to-go "crossnsifn" (quote crossnsifn)
		   [(quote (defn caller4 [y] (count y)))])`)
	src := string(v.(vm.String))
	f := parseLoweredGo(t, src)

	caller, ok := findFunc(f, func(n string) bool { return n == "caller4" })
	if !ok {
		t.Fatalf("no caller4 in:\n%s", src)
	}
	d, ok := findIFnDispatch(caller.Body)
	if !ok {
		t.Fatalf("expected cached-var IFn dispatch in caller4:\n%s", src)
	}
	if d.nargs != 1 {
		t.Fatalf("Invoke got %d args, want 1", d.nargs)
	}
	if !strings.HasSuffix(d.varName, "count") {
		t.Fatalf("dispatch var %q does not target count", d.varName)
	}
	if d.nsArg != "clojure.core" || d.nameArg != "count" {
		t.Fatalf("CachedVarFn resolves (%q,%q), want (clojure.core,count)", d.nsArg, d.nameArg)
	}
	// The call site must NOT do a per-call rt.LookupVar for count anymore.
	if callsLookupVarNamed(caller, "count") {
		t.Fatalf("caller4 still does a per-call rt.LookupVar for count:\n%s", src)
	}
}

// Task 2: a callee that is NOT a static load-var (here, an invoked parameter)
// must stay on the rt.InvokeValue trampoline — cached-var dispatch only applies
// to resolvable vars.
func TestLowerGoNonVarCallStaysTrampoline(t *testing.T) {
	ensureLoader()
	runLispExpr(t, `(create-ns (quote nonvarcall))`)
	v := runLispExpr(t,
		`(ir.passes.pipeline/lower-ns-to-go "nonvarcall" (quote nonvarcall)
		   [(quote (defn apply1 [f y] (f y)))])`)
	src := string(v.(vm.String))
	f := parseLoweredGo(t, src)

	caller, ok := findFunc(f, func(n string) bool { return n == "apply1" })
	if !ok {
		t.Fatalf("no apply1 in:\n%s", src)
	}
	if _, ok := findIFnDispatch(caller.Body); ok {
		t.Fatalf("calling a parameter must NOT use cached-var dispatch:\n%s", src)
	}
	if !callsInvokeValue(caller) {
		t.Fatalf("expected rt.InvokeValue trampoline for a parameter call:\n%s", src)
	}
}

// Task 3 (decision: fall back): (deref x) is itself a cross-ns call to
// clojure.core/deref, so it is already lifted to cached-var IFn dispatch by
// the Task-2 path — eliminating the per-call rt.LookupVar. We deliberately do
// NOT lower it to a direct x.(vm.IDeref).Deref(): there is no vm.IDeref Go
// interface, and protocol-based derefables (IDerefProtocol/iderefAdapter) have
// no Go Deref() method, so a type-assert would panic. This pins that the
// optimized cached-var dispatch is used and no IDeref assertion is emitted.
func TestLowerGoDerefUsesCachedVarDispatch(t *testing.T) {
	ensureLoader()
	runLispExpr(t, `(create-ns (quote derefns))`)
	v := runLispExpr(t,
		`(ir.passes.pipeline/lower-ns-to-go "derefns" (quote derefns)
		   [(quote (defn dr [a] (deref a)))])`)
	src := string(v.(vm.String))
	f := parseLoweredGo(t, src)

	caller, ok := findFunc(f, func(n string) bool { return n == "dr" })
	if !ok {
		t.Fatalf("no dr in:\n%s", src)
	}
	d, ok := findIFnDispatch(caller.Body)
	if !ok {
		t.Fatalf("expected deref to lower to cached-var IFn dispatch:\n%s", src)
	}
	if !strings.HasSuffix(d.varName, "deref") {
		t.Fatalf("dispatch var %q does not target deref", d.varName)
	}
	if callsLookupVarNamed(caller, "deref") {
		t.Fatalf("dr must not do a per-call rt.LookupVar for deref:\n%s", src)
	}
	if assertsType(caller.Body, "vm.IDeref") {
		t.Fatalf("deref must NOT lower to a direct .(vm.IDeref) assertion:\n%s", src)
	}
}

func TestLowerGoStrictExactDTypeArgStaysConcrete(t *testing.T) {
	ensureLoader()

	fn := buildLispIR(t, `(defn id-square [x] x)`)

	passVarCounter++
	fnVar := fmt.Sprintf("*lower-go-dtype-fn-%d*", passVarCounter)
	rt.NS(rt.NameCoreNS).Def(fnVar, fn)

	runLispExpr(t, fmt.Sprintf(`(let [s (ir.lattice/new-typeinfra-state %s)
	                                arg0 (ir/fn-load-arg %s 0)]
	                            (ir.lattice/seed-state-from-inst-types! s %s)
	                            (ir.lattice/join-inst-type! s arg0 [:dtype 'Square])
	                            (ir.passes.typeinfer/typeinfer %s s)
	                            (ir.lattice/flush-state-types! s %s))`,
		fnVar, fnVar, fnVar, fnVar, fnVar))

	result := lowerGo(t, fn, ":strict")
	if got := result.ValueAt(vm.Keyword("status")); got != vm.Keyword("lowered") {
		t.Fatalf("expected :lowered status, got %v (reason: %v)", got, result.ValueAt(vm.Keyword("reason")))
	}

	rendered := bindAndRenderGoDecl(t, result)
	if !strings.Contains(rendered, "arg0 *Square") {
		t.Fatalf("expected exact dtype arg to lower as *Square, not vm.Value\n--- dump ---\n%s\n--- go ---\n%s", lispDump(t, fn), rendered)
	}
	if !strings.Contains(rendered, ") *Square") {
		t.Fatalf("expected exact dtype return to stay concrete\n--- go ---\n%s", rendered)
	}
	if strings.Contains(rendered, "arg0 vm.Value") {
		t.Fatalf("dtype arg widened back to vm.Value\n--- go ---\n%s", rendered)
	}
}

// A capturing closure nested inside another capturing closure must get a
// LEXICALLY cumulative arg prefix (c<outer>_c<inner>_...), not a flat
// c<nid>_. nid is unique only within one function IR, so two nested closure
// templates can share a nid; a flat prefix then makes the inner closure's
// param collide with the captured outer param it shadows — the same
// "captured name shadowed by block param" miscompile this PR fixes, just one
// level deeper. (Reported on PR #247.)
func TestLowerGoNestedCapturedClosurePrefixesAreLexical(t *testing.T) {
	ensureLoader()

	fn := buildLispIR(t, `(defn nested-captures [x]
		(fn* [y]
			(fn* [z]
				[x y z])))`)
	optimizeLispIR(t, fn)
	result := lowerGo(t, fn, ":strict")

	if got := result.ValueAt(vm.Keyword("status")); got != vm.Keyword("lowered") {
		t.Fatalf("expected :lowered status, got %v (reason: %v)", got, result.ValueAt(vm.Keyword("reason")))
	}

	rendered := bindAndRenderGoDecl(t, result)

	// Each capturing closure names its first param <prefix>arg0 where prefix is
	// one or more c<nid>_ segments. The two nested capturing closures appear in
	// source order: outer first, inner second. The inner prefix must LEXICALLY
	// EXTEND the outer (c<o>_ -> c<o>_c<i>_), which guarantees distinctness
	// regardless of whether the two templates happen to share a nid. A flat
	// per-nid scheme would instead emit sibling prefixes (c<o>_, c<i>_) that are
	// only accidentally distinct.
	re := regexp.MustCompile(`func\(((?:c[0-9]+_)+)arg0 vm\.Value\)`)
	ms := re.FindAllStringSubmatch(rendered, -1)
	if len(ms) < 2 {
		t.Fatalf("expected two nested prefixed closure params, found %d\n--- go ---\n%s", len(ms), rendered)
	}
	outer, inner := ms[0][1], ms[1][1]
	if inner == outer || !strings.HasPrefix(inner, outer) {
		t.Fatalf("inner closure prefix %q must lexically extend outer %q (else inner param can shadow captured outer param)\n--- go ---\n%s", inner, outer, rendered)
	}
}
