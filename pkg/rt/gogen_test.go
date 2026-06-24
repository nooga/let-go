/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

// --- helpers ---------------------------------------------------------

// render is a test convenience that wraps cRender and asserts no error.
func render(t *testing.T, node vm.Value) string {
	t.Helper()
	out, err := cRender(node)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}
	return string(out.(vm.String))
}

// renderErr renders and returns the error, asserting the call did error.
func renderErr(t *testing.T, node vm.Value) error {
	t.Helper()
	_, err := cRender(node)
	if err == nil {
		t.Fatalf("expected render error, got nil")
	}
	return err
}

// must asserts no error and returns the value. Use as:
//
//	x := must(t)(cIdent(vm.String("x")))
//
// The currying lets us pass t once and reuse it through the call.
func must(t *testing.T) func(vm.Value, error) vm.Value {
	return func(v vm.Value, err error) vm.Value {
		t.Helper()
		if err != nil {
			t.Fatalf("constructor failed: %v", err)
		}
		return v
	}
}

// errMust asserts the (value, error) pair contains an error with msg.
// Use as: errMust(t, "valid Go identifier")(cIdent(vm.String("1bad")))
// The currying matches the pattern of must(t)(...).
func errMust(t *testing.T, msg string) func(vm.Value, error) {
	return func(_ vm.Value, err error) {
		t.Helper()
		if err == nil {
			t.Fatalf("expected error containing %q, got nil", msg)
		}
		if !strings.Contains(err.Error(), msg) {
			t.Fatalf("expected error containing %q, got: %v", msg, err)
		}
	}
}

// --- literal constructors --------------------------------------------

func TestIntLit(t *testing.T) {
	node := must(t)(cIntLit(vm.Int(42)))
	got := render(t, node)
	if got != "42" {
		t.Errorf("got %q, want %q", got, "42")
	}
}

func TestFloatLit(t *testing.T) {
	node := must(t)(cFloatLit(vm.Float(3.14)))
	got := render(t, node)
	if got != "3.14" {
		t.Errorf("got %q, want %q", got, "3.14")
	}
}

func TestStringLit(t *testing.T) {
	node := must(t)(cStringLit(vm.String(`hello "world"`)))
	got := render(t, node)
	if got != `"hello \"world\""` {
		t.Errorf("got %q, want quoted form", got)
	}
}

func TestCharLit(t *testing.T) {
	node := must(t)(cCharLit(vm.Char('a')))
	got := render(t, node)
	if got != `'a'` {
		t.Errorf("got %q, want %q", got, `'a'`)
	}
}

func TestIdent(t *testing.T) {
	node := must(t)(cIdent(vm.String("myVar")))
	if got := render(t, node); got != "myVar" {
		t.Errorf("got %q, want %q", got, "myVar")
	}
}

func TestIdentRejectsInvalidNames(t *testing.T) {
	cases := []string{"", "1var", "with-hyphen", "has.dot", "has space"}
	for _, name := range cases {
		t.Run(name, func(t *testing.T) {
			errMust(t, "not a valid Go identifier")(cIdent(vm.String(name)))
		})
	}
}

// --- expression constructors ----------------------------------------

func TestBinary(t *testing.T) {
	x := must(t)(cIdent(vm.String("x")))
	y := must(t)(cIdent(vm.String("y")))
	node := must(t)(cBinary(vm.String("+"), x, y))
	if got := render(t, node); got != "x + y" {
		t.Errorf("got %q, want %q", got, "x + y")
	}
}

func TestBinaryUnknownOp(t *testing.T) {
	x := must(t)(cIdent(vm.String("x")))
	y := must(t)(cIdent(vm.String("y")))
	errMust(t, "unknown operator")(cBinary(vm.String("@@"), x, y))
}

func TestUnary(t *testing.T) {
	x := must(t)(cIdent(vm.String("x")))
	node := must(t)(cUnary(vm.String("!"), x))
	if got := render(t, node); got != "!x" {
		t.Errorf("got %q, want %q", got, "!x")
	}
}

func TestCall(t *testing.T) {
	fn := must(t)(cIdent(vm.String("foo")))
	a := must(t)(cIntLit(vm.Int(1)))
	b := must(t)(cIntLit(vm.Int(2)))
	args := vm.NewArrayVector([]vm.Value{a, b})
	node := must(t)(cCall(fn, args))
	if got := render(t, node); got != "foo(1, 2)" {
		t.Errorf("got %q, want %q", got, "foo(1, 2)")
	}
}

func TestIndex(t *testing.T) {
	arr := must(t)(cIdent(vm.String("arr")))
	idx := must(t)(cIntLit(vm.Int(0)))
	node := must(t)(cIndex(arr, idx))
	if got := render(t, node); got != "arr[0]" {
		t.Errorf("got %q, want %q", got, "arr[0]")
	}
}

func TestFieldSelector(t *testing.T) {
	recv := must(t)(cIdent(vm.String("obj")))
	node := must(t)(cFieldSel(recv, vm.String("Name")))
	if got := render(t, node); got != "obj.Name" {
		t.Errorf("got %q, want %q", got, "obj.Name")
	}
}

func TestFieldSelectorRejectsDottedName(t *testing.T) {
	// Reviewer finding: (field a 'X.Y') would silently produce invalid Go.
	recv := must(t)(cIdent(vm.String("obj")))
	errMust(t, "not a valid identifier")(cFieldSel(recv, vm.String("X.Y")))
}

func TestTypeAssert(t *testing.T) {
	x := must(t)(cIdent(vm.String("v")))
	intT := must(t)(cType(vm.String("int")))
	node := must(t)(cTypeAssert(x, intT))
	if got := render(t, node); got != "v.(int)" {
		t.Errorf("got %q, want %q", got, "v.(int)")
	}
}

func TestCast(t *testing.T) {
	x := must(t)(cIdent(vm.String("v")))
	intT := must(t)(cType(vm.String("int")))
	node := must(t)(cCast(intT, x))
	if got := render(t, node); got != "int(v)" {
		t.Errorf("got %q, want %q", got, "int(v)")
	}
}

func TestTypeSwitchStmt(t *testing.T) {
	// switch x := v.(type) { case *vm.ArrayVector: return x; default: return nil }
	v := must(t)(cIdent(vm.String("v")))
	avT := must(t)(cType(vm.String("*vm.ArrayVector")))
	xIdent := must(t)(cIdent(vm.String("x")))
	retX := must(t)(cReturn(vm.NewArrayVector([]vm.Value{xIdent})))
	caseArm := must(t)(cCaseClause(
		vm.NewArrayVector([]vm.Value{avT}),
		vm.NewArrayVector([]vm.Value{retX})))
	nilIdent := must(t)(cIdent(vm.String("nil")))
	retNil := must(t)(cReturn(vm.NewArrayVector([]vm.Value{nilIdent})))
	defArm := must(t)(cCaseClause(vm.NIL, vm.NewArrayVector([]vm.Value{retNil})))
	node := must(t)(cTypeSwitchStmt(vm.String("x"), v,
		vm.NewArrayVector([]vm.Value{caseArm, defArm})))
	got := render(t, node)
	for _, want := range []string{"switch x := v.(type)", "case *vm.ArrayVector:", "default:", "return x", "return nil"} {
		if !strings.Contains(got, want) {
			t.Errorf("render %q missing %q", got, want)
		}
	}
}

func TestTypeSwitchStmtNoBinding(t *testing.T) {
	// switch v.(type) { default: return nil }  (no x := binding)
	v := must(t)(cIdent(vm.String("v")))
	nilIdent := must(t)(cIdent(vm.String("nil")))
	retNil := must(t)(cReturn(vm.NewArrayVector([]vm.Value{nilIdent})))
	defArm := must(t)(cCaseClause(vm.NIL, vm.NewArrayVector([]vm.Value{retNil})))
	node := must(t)(cTypeSwitchStmt(vm.NIL, v, vm.NewArrayVector([]vm.Value{defArm})))
	got := render(t, node)
	if !strings.Contains(got, "switch v.(type)") {
		t.Errorf("render %q missing %q", got, "switch v.(type)")
	}
}

// --- statement constructors -----------------------------------------

func TestVarDecl(t *testing.T) {
	intT := must(t)(cType(vm.String("int")))
	zero := must(t)(cIntLit(vm.Int(0)))
	node := must(t)(cVarDecl(vm.String("x"), intT, zero))
	got := render(t, node)
	if got != "var x int = 0" {
		t.Errorf("got %q, want %q", got, "var x int = 0")
	}
}

func TestVarDeclNoInit(t *testing.T) {
	intT := must(t)(cType(vm.String("int")))
	node := must(t)(cVarDecl(vm.String("x"), intT, vm.NIL))
	got := render(t, node)
	if got != "var x int" {
		t.Errorf("got %q, want %q", got, "var x int")
	}
}

func TestAssign(t *testing.T) {
	tgt := must(t)(cIdent(vm.String("x")))
	v := must(t)(cIntLit(vm.Int(1)))
	node := must(t)(cAssign(vm.String("="), tgt, v))
	if got := render(t, node); got != "x = 1" {
		t.Errorf("got %q, want %q", got, "x = 1")
	}
}

func TestAssignCompound(t *testing.T) {
	tgt := must(t)(cIdent(vm.String("sum")))
	v := must(t)(cIntLit(vm.Int(5)))
	node := must(t)(cAssign(vm.String("+="), tgt, v))
	if got := render(t, node); got != "sum += 5" {
		t.Errorf("got %q, want %q", got, "sum += 5")
	}
}

func TestMultiAssign(t *testing.T) {
	r := must(t)(cIdent(vm.String("r")))
	ok := must(t)(cIdent(vm.String("ok")))
	lhs := vm.NewArrayVector([]vm.Value{r, ok})
	fn := must(t)(cIdent(vm.String("foo")))
	call := must(t)(cCall(fn, vm.NewArrayVector(nil)))
	rhs := vm.NewArrayVector([]vm.Value{call})
	node := must(t)(cMultiAssign(vm.String(":="), lhs, rhs))
	if got := render(t, node); got != "r, ok := foo()" {
		t.Errorf("got %q, want %q", got, "r, ok := foo()")
	}
}

func TestReturn(t *testing.T) {
	r := must(t)(cIdent(vm.String("r")))
	vals := vm.NewArrayVector([]vm.Value{r})
	node := must(t)(cReturn(vals))
	if got := render(t, node); got != "return r" {
		t.Errorf("got %q, want %q", got, "return r")
	}
}

func TestReturnMultiValue(t *testing.T) {
	r := must(t)(cIdent(vm.String("r")))
	nilExpr := must(t)(cIdent(vm.String("nil")))
	vals := vm.NewArrayVector([]vm.Value{r, nilExpr})
	node := must(t)(cReturn(vals))
	if got := render(t, node); got != "return r, nil" {
		t.Errorf("got %q, want %q", got, "return r, nil")
	}
}

func TestIfStmt(t *testing.T) {
	cond := must(t)(cIdent(vm.String("ok")))
	body := must(t)(cReturn(vm.NewArrayVector([]vm.Value{
		must(t)(cIntLit(vm.Int(1))),
	})))
	then := vm.NewArrayVector([]vm.Value{body})
	node := must(t)(cIfStmt(vm.NIL, cond, then, vm.NIL))
	got := render(t, node)
	want := "if ok {\n\treturn 1\n}"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestIfStmtWithInit(t *testing.T) {
	// Build: if x := foo(); x { return 1 }
	x := must(t)(cIdent(vm.String("x")))
	fn := must(t)(cIdent(vm.String("foo")))
	call := must(t)(cCall(fn, vm.NewArrayVector(nil)))
	init := must(t)(cMultiAssign(vm.String(":="),
		vm.NewArrayVector([]vm.Value{x}),
		vm.NewArrayVector([]vm.Value{call})))
	body := must(t)(cReturn(vm.NewArrayVector([]vm.Value{
		must(t)(cIntLit(vm.Int(1))),
	})))
	then := vm.NewArrayVector([]vm.Value{body})
	node := must(t)(cIfStmt(init, x, then, vm.NIL))
	got := render(t, node)
	want := "if x := foo(); x {\n\treturn 1\n}"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestGotoStmt(t *testing.T) {
	node := must(t)(cGotoStmt(vm.String("done")))
	if got := render(t, node); got != "goto done" {
		t.Errorf("got %q, want %q", got, "goto done")
	}
}

func TestLabelStmt(t *testing.T) {
	stmt := must(t)(cReturn(vm.NewArrayVector([]vm.Value{
		must(t)(cIntLit(vm.Int(1))),
	})))
	node := must(t)(cLabelStmt(vm.String("done"), stmt))
	got := render(t, node)
	want := "done:\n\treturn 1"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

// --- function declarations ------------------------------------------

func TestFuncDecl(t *testing.T) {
	// func Add(x int, y int) int { return x + y }
	intT := must(t)(cType(vm.String("int")))
	x := must(t)(cIdent(vm.String("x")))
	y := must(t)(cIdent(vm.String("y")))
	xParam := must(t)(cParam(vm.String("x"), must(t)(cType(vm.String("int")))))
	yParam := must(t)(cParam(vm.String("y"), must(t)(cType(vm.String("int")))))
	result := must(t)(cResult(intT))
	body := must(t)(cReturn(vm.NewArrayVector([]vm.Value{
		must(t)(cBinary(vm.String("+"), x, y)),
	})))
	node := must(t)(cFuncDecl(
		vm.String("Add"),
		vm.NewArrayVector([]vm.Value{xParam, yParam}),
		vm.NewArrayVector([]vm.Value{result}),
		vm.NewArrayVector([]vm.Value{body}),
	))
	got := render(t, node)
	want := "func Add(x int, y int) int {\n\treturn x + y\n}"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFuncDeclMultiReturn(t *testing.T) {
	// func F() (int, error) { return 0, nil }
	intT := must(t)(cType(vm.String("int")))
	errT := must(t)(cType(vm.String("error")))
	zero := must(t)(cIntLit(vm.Int(0)))
	nilE := must(t)(cIdent(vm.String("nil")))
	body := must(t)(cReturn(vm.NewArrayVector([]vm.Value{zero, nilE})))
	node := must(t)(cFuncDecl(
		vm.String("F"),
		vm.NewArrayVector(nil),
		vm.NewArrayVector([]vm.Value{
			must(t)(cResult(intT)),
			must(t)(cResult(errT)),
		}),
		vm.NewArrayVector([]vm.Value{body}),
	))
	got := render(t, node)
	want := "func F() (int, error) {\n\treturn 0, nil\n}"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

// --- file emission --------------------------------------------------

func TestFile(t *testing.T) {
	imp := must(t)(cImportSpec(vm.String("fmt")))
	body := must(t)(cReturn(vm.NewArrayVector([]vm.Value{
		must(t)(cIntLit(vm.Int(42))),
	})))
	fn := must(t)(cFuncDecl(
		vm.String("Answer"),
		vm.NewArrayVector(nil),
		vm.NewArrayVector([]vm.Value{
			must(t)(cResult(must(t)(cType(vm.String("int"))))),
		}),
		vm.NewArrayVector([]vm.Value{body}),
	))
	node := must(t)(cFile(
		vm.String("main"),
		vm.NewArrayVector([]vm.Value{imp}),
		vm.NewArrayVector([]vm.Value{fn}),
	))
	got := render(t, node)
	if !strings.HasPrefix(got, "package main") {
		t.Errorf("expected file to start with 'package main', got: %s", got)
	}
	if !strings.Contains(got, `import "fmt"`) {
		t.Errorf("expected import fmt, got: %s", got)
	}
	if !strings.Contains(got, "func Answer() int") {
		t.Errorf("expected func Answer, got: %s", got)
	}
}

// --- determinism ----------------------------------------------------

// TestRenderIsDeterministic builds a complex AST and renders it many
// times; if any output differs from the first, the renderer has a
// nondeterminism bug (probably Go-map iteration).
func TestRenderIsDeterministic(t *testing.T) {
	// Build something with enough surface to exercise dispatch:
	// func F(a int, b int) (int, error) {
	//     if x, ok := foo(a); ok {
	//         return x + b, nil
	//     }
	//     return 0, nil
	// }
	intT := func() vm.Value { return must(t)(cType(vm.String("int"))) }
	errT := must(t)(cType(vm.String("error")))
	a := func() vm.Value { return must(t)(cIdent(vm.String("a"))) }
	b := func() vm.Value { return must(t)(cIdent(vm.String("b"))) }
	x := func() vm.Value { return must(t)(cIdent(vm.String("x"))) }
	ok := func() vm.Value { return must(t)(cIdent(vm.String("ok"))) }
	nilE := func() vm.Value { return must(t)(cIdent(vm.String("nil"))) }
	zero := func() vm.Value { return must(t)(cIntLit(vm.Int(0))) }

	buildAST := func() vm.Value {
		fn := must(t)(cIdent(vm.String("foo")))
		call := must(t)(cCall(fn, vm.NewArrayVector([]vm.Value{a()})))
		init := must(t)(cMultiAssign(vm.String(":="),
			vm.NewArrayVector([]vm.Value{x(), ok()}),
			vm.NewArrayVector([]vm.Value{call})))
		thenRet := must(t)(cReturn(vm.NewArrayVector([]vm.Value{
			must(t)(cBinary(vm.String("+"), x(), b())),
			nilE(),
		})))
		ifs := must(t)(cIfStmt(init, ok(),
			vm.NewArrayVector([]vm.Value{thenRet}),
			vm.NIL))
		tailRet := must(t)(cReturn(vm.NewArrayVector([]vm.Value{zero(), nilE()})))
		return must(t)(cFuncDecl(
			vm.String("F"),
			vm.NewArrayVector([]vm.Value{
				must(t)(cParam(vm.String("a"), intT())),
				must(t)(cParam(vm.String("b"), intT())),
			}),
			vm.NewArrayVector([]vm.Value{
				must(t)(cResult(intT())),
				must(t)(cResult(errT)),
			}),
			vm.NewArrayVector([]vm.Value{ifs, tailRet}),
		))
	}

	first := render(t, buildAST())
	for i := range 100 {
		got := render(t, buildAST())
		if got != first {
			t.Fatalf("render %d differs from first:\nfirst:\n%s\nlater:\n%s", i, first, got)
		}
	}
}

// --- render error paths ---------------------------------------------

func TestRenderRejectsNonAST(t *testing.T) {
	_, err := cRender(vm.String("not an ast"))
	if err == nil {
		t.Fatal("expected error rendering a non-AST value")
	}
}

func TestRenderHandlesNil(t *testing.T) {
	out, err := cRender(vm.NIL)
	if err != nil {
		t.Fatalf("rendering NIL should be empty string, got error: %v", err)
	}
	if string(out.(vm.String)) != "" {
		t.Errorf("expected empty string, got %q", out)
	}
}

func TestGogenIdentAccessors(t *testing.T) {
	// Test (gogen/ident? (gogen/ident "foo")) returns true
	identNode := must(t)(cIdent(vm.String("foo")))
	isIdentResult, err := cIdentP(identNode)
	if err != nil {
		t.Fatalf("cIdentP failed: %v", err)
	}
	if isIdentResult != vm.TRUE {
		t.Errorf("cIdentP on ident node: got %v, want vm.TRUE", isIdentResult)
	}

	// Test (gogen/ident? "foo") returns false (non-ident value)
	isIdentResult2, err := cIdentP(vm.String("foo"))
	if err != nil {
		t.Fatalf("cIdentP on non-ident: got error %v, expected nil", err)
	}
	if isIdentResult2 != vm.FALSE {
		t.Errorf("cIdentP on non-ident value: got %v, want vm.FALSE", isIdentResult2)
	}

	// Test (gogen/ident-name (gogen/ident "foo")) returns "foo"
	nameResult, err := cIdentName(identNode)
	if err != nil {
		t.Fatalf("cIdentName failed: %v", err)
	}
	nameStr, ok := nameResult.(vm.String)
	if !ok {
		t.Fatalf("cIdentName returned non-string: %T", nameResult)
	}
	if string(nameStr) != "foo" {
		t.Errorf("cIdentName: got %q, want %q", string(nameStr), "foo")
	}

	// Test (gogen/ident-name "foo") returns error (non-ident)
	_, err = cIdentName(vm.String("foo"))
	if err == nil {
		t.Fatalf("cIdentName on non-ident: expected error, got nil")
	}
	// The error should mention either the type mismatch or *ast.Ident
	errMsg := err.Error()
	if !strings.Contains(errMsg, "String") && !strings.Contains(errMsg, "*ast.Ident") {
		t.Errorf("cIdentName error should mention type issue, got: %v", err)
	}
}
