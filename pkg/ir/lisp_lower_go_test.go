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
	if !strings.Contains(rendered, "goto ") {
		t.Fatalf("expected join lowering to use CFG goto\n--- go ---\n%s", rendered)
	}
	// Join block params are typed vm.Value (typeinfer doesn't currently
	// narrow them from int-typed feeds), so the trailing arithmetic
	// lowers via rt.AddValue rather than a Go-native `+ 1`.
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
	if !strings.Contains(rendered, "goto ") || !strings.Contains(rendered, ":") {
		t.Fatalf("expected loop lowering to use labels/gotos\n--- go ---\n%s", rendered)
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

func TestLowerGoBridgeFallsBackOnUnsupportedQuotedUUIDConst(t *testing.T) {
	ensureLoader()

	fn := buildLispIR(t, `(defn quoted-uuid [] (quote #uuid "123e4567-e89b-12d3-a456-426614174000"))`)
	optimizeLispIR(t, fn)
	result := lowerGo(t, fn, ":bridge")

	if got := result.ValueAt(vm.Keyword("status")); got != vm.Keyword("fallback") {
		t.Fatalf("expected :fallback status, got %v", got)
	}
	if result.ValueAt(vm.Keyword("decl")) != vm.NIL {
		t.Fatalf("expected bridge fallback to omit :decl")
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

func TestLowerGoFallsBackOnUnloweredSideEffectingCall(t *testing.T) {
	ensureLoader()
	// push-binding! takes a `(var *ns*)` arg that gogen cannot lower, so
	// call-assign-stmts returns nil. A side-effecting call must NOT be silently
	// dropped (which would change behavior and leave its arg a dead store) — the
	// whole function must fall back to bytecode instead.
	fn := buildLispIR(t, `(defn cf [form caller-ns]
	  (do
	    (push-binding! (var *ns*) (first caller-ns))
	    (let [r (do (count form))]
	      (do (pop-binding! (var *ns*)) r))))`)
	optimizeLispIR(t, fn)
	result := lowerGo(t, fn, ":bridge")
	if got := result.ValueAt(vm.Keyword("status")); got != vm.Keyword("fallback") {
		t.Fatalf("expected :fallback for un-lowerable side-effecting push-binding!, got %v\n%s",
			got, bindAndRenderGoDecl(t, result))
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
	if !strings.Contains(r2, "func ch() vm.Char") {
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

// Scope guard: a cross-namespace callee (clojure.core/count) is NOT lifted —
// it stays on the InvokeValue path (slice 1 is intra-ns only).
func TestLowerNsToGoLeavesCrossNsCall(t *testing.T) {
	ensureLoader()
	runLispExpr(t, `(create-ns (quote directtest2))`)
	v := runLispExpr(t,
		`(ir.passes.pipeline/lower-ns-to-go "directtest2" (quote directtest2)
		   [(quote (defn caller2 [y] (count y)))])`)
	got := string(v.(vm.String))
	if !regexp.MustCompile(`InvokeValue\([^\n]*LookupVar\([^\n]*"count"`).MatchString(got) {
		t.Fatalf("expected cross-ns count call to stay on InvokeValue:\n%s", got)
	}
}
