/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
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
	if !strings.Contains(rendered, "+ 1") {
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
	if !strings.Contains(rendered, "+") {
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
	if !strings.Contains(rendered, "+") {
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
	if !strings.Contains(rendered, "+") {
		t.Fatalf("expected arithmetic in lowered body\n--- go ---\n%s", rendered)
	}
}
