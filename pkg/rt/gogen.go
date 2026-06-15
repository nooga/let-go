/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

// Package rt — gogen.go installs the `gogen` namespace's Go AST
// constructors and renderer.
//
// Architecture (refactored): Clojure code calls constructors that
// build real *go/ast nodes, boxed as goASTValue. There is no
// intermediate map representation; the macro layer produces calls
// to these constructors, and the renderer is just go/format.Node.
//
// This means:
//   - Whatever go/ast can express, gogen can express.
//   - Errors surface at construction (e.g., invalid identifiers).
//   - New Go language features are exposed by adding constructors,
//     no protocol changes needed.

package rt

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"reflect"
	"strconv"
	"strings"

	"github.com/nooga/let-go/pkg/vm"
)

// --- goASTValue: boxed *ast.Node as a vm.Value -----------------------

type theGoASTType struct{}

func (t *theGoASTType) String() string     { return t.Name() }
func (t *theGoASTType) Type() vm.ValueType { return vm.TypeType }
func (t *theGoASTType) Unbox() any         { return reflect.TypeFor[*theGoASTType]() }
func (t *theGoASTType) Name() string       { return "let-go.lang.GoAST" }
func (t *theGoASTType) Box(_ any) (vm.Value, error) {
	return vm.NIL, fmt.Errorf("gogen: GoAST values are constructed via gogen/* fns, not boxed")
}

var GoASTType *theGoASTType = &theGoASTType{}

type goASTValue struct{ node ast.Node }

func (g *goASTValue) String() string     { return fmt.Sprintf("#<go-ast %T>", g.node) }
func (g *goASTValue) Type() vm.ValueType { return GoASTType }
func (g *goASTValue) Unbox() any         { return g.node }

func box(n ast.Node) vm.Value {
	if n == nil {
		return vm.NIL
	}
	return &goASTValue{node: n}
}

func unboxNode(v vm.Value) (ast.Node, error) {
	if v == vm.NIL {
		return nil, nil
	}
	g, ok := v.(*goASTValue)
	if !ok {
		return nil, fmt.Errorf("gogen: expected go-ast value, got %s", v.Type().Name())
	}
	return g.node, nil
}

func unboxExpr(v vm.Value) (ast.Expr, error) {
	n, err := unboxNode(v)
	if err != nil {
		return nil, err
	}
	if n == nil {
		return nil, nil
	}
	e, ok := n.(ast.Expr)
	if !ok {
		return nil, fmt.Errorf("gogen: expected go expression, got %T", n)
	}
	return e, nil
}

func unboxStmt(v vm.Value) (ast.Stmt, error) {
	n, err := unboxNode(v)
	if err != nil {
		return nil, err
	}
	if n == nil {
		return nil, nil
	}
	s, ok := n.(ast.Stmt)
	if !ok {
		return nil, fmt.Errorf("gogen: expected go statement, got %T", n)
	}
	return s, nil
}

func unboxDecl(v vm.Value) (ast.Decl, error) {
	n, err := unboxNode(v)
	if err != nil {
		return nil, err
	}
	if n == nil {
		return nil, nil
	}
	d, ok := n.(ast.Decl)
	if !ok {
		return nil, fmt.Errorf("gogen: expected go decl, got %T", n)
	}
	return d, nil
}

// seqToValues flattens a Clojure sequable into a []vm.Value.
func seqToValues(v vm.Value) ([]vm.Value, error) {
	if v == vm.NIL {
		return nil, nil
	}
	sq, ok := v.(vm.Sequable)
	if !ok {
		return nil, fmt.Errorf("gogen: expected seq, got %s", v.Type().Name())
	}
	var out []vm.Value
	for s := sq.Seq(); s != nil && s != vm.EmptyList; s = s.Next() {
		out = append(out, s.First())
	}
	return out, nil
}

// asString extracts a string from a vm.String / vm.Symbol / vm.Keyword.
func asString(v vm.Value) (string, error) {
	switch x := v.(type) {
	case vm.String:
		return string(x), nil
	case vm.Symbol:
		return string(x), nil
	case vm.Keyword:
		return string(x), nil
	}
	return "", fmt.Errorf("gogen: expected string-like, got %s", v.Type().Name())
}

// --- identifier validation -------------------------------------------
//
// go/ast happily accepts garbage in *ast.Ident.Name and *ast.SelectorExpr.Sel.
// We validate up front so users get errors at construction, not at go build.

func validIdent(s string) bool {
	if s == "" {
		return false
	}
	for i, r := range s {
		if r == '_' || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			continue
		}
		if i > 0 && r >= '0' && r <= '9' {
			continue
		}
		return false
	}
	return true
}

// --- operator token table --------------------------------------------

var opTokens = map[string]token.Token{
	"+":  token.ADD,
	"-":  token.SUB,
	"*":  token.MUL,
	"/":  token.QUO,
	"%":  token.REM,
	"<":  token.LSS,
	"<=": token.LEQ,
	">":  token.GTR,
	">=": token.GEQ,
	"==": token.EQL,
	"!=": token.NEQ,
	"&&": token.LAND,
	"||": token.LOR,
	"&":  token.AND,
	"|":  token.OR,
	"^":  token.XOR,
	"<<": token.SHL,
	">>": token.SHR,
	"&^": token.AND_NOT,
	"!":  token.NOT,
	"+=": token.ADD_ASSIGN,
	"-=": token.SUB_ASSIGN,
	"*=": token.MUL_ASSIGN,
	"/=": token.QUO_ASSIGN,
	"%=": token.REM_ASSIGN,
	"=":  token.ASSIGN,
	":=": token.DEFINE,
}

func opTokenOrErr(op string) (token.Token, error) {
	if t, ok := opTokens[op]; ok {
		return t, nil
	}
	return token.ILLEGAL, fmt.Errorf("gogen: unknown operator %q", op)
}

// --- constructor helpers ---------------------------------------------

func wrap0(name string, fn func() (vm.Value, error)) (vm.Value, error) {
	return vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 0 {
			return vm.NIL, fmt.Errorf("gogen/%s: expected 0 args, got %d", name, len(vs))
		}
		return fn()
	})
}

func wrap1(name string, fn func(vm.Value) (vm.Value, error)) (vm.Value, error) {
	return vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("gogen/%s: expected 1 arg, got %d", name, len(vs))
		}
		return fn(vs[0])
	})
}

func wrap2(name string, fn func(vm.Value, vm.Value) (vm.Value, error)) (vm.Value, error) {
	return vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("gogen/%s: expected 2 args, got %d", name, len(vs))
		}
		return fn(vs[0], vs[1])
	})
}

func wrap3(name string, fn func(vm.Value, vm.Value, vm.Value) (vm.Value, error)) (vm.Value, error) {
	return vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 3 {
			return vm.NIL, fmt.Errorf("gogen/%s: expected 3 args, got %d", name, len(vs))
		}
		return fn(vs[0], vs[1], vs[2])
	})
}

func wrap4(name string, fn func(vm.Value, vm.Value, vm.Value, vm.Value) (vm.Value, error)) (vm.Value, error) {
	return vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 4 {
			return vm.NIL, fmt.Errorf("gogen/%s: expected 4 args, got %d", name, len(vs))
		}
		return fn(vs[0], vs[1], vs[2], vs[3])
	})
}

func wrap5(name string, fn func(vm.Value, vm.Value, vm.Value, vm.Value, vm.Value) (vm.Value, error)) (vm.Value, error) {
	return vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 5 {
			return vm.NIL, fmt.Errorf("gogen/%s: expected 5 args, got %d", name, len(vs))
		}
		return fn(vs[0], vs[1], vs[2], vs[3], vs[4])
	})
}

// --- ast constructors (each returns a boxed go/ast node) -------------

// ident: (gogen/ident "name") -> *ast.Ident
func cIdent(v vm.Value) (vm.Value, error) {
	s, err := asString(v)
	if err != nil {
		return vm.NIL, err
	}
	if !validIdent(s) {
		return vm.NIL, fmt.Errorf("gogen: %q is not a valid Go identifier", s)
	}
	return box(ast.NewIdent(s)), nil
}

// ident?: (gogen/ident? value) -> true if value is an *ast.Ident, false otherwise
func cIdentP(v vm.Value) (vm.Value, error) {
	n, err := unboxNode(v)
	if err != nil {
		return vm.FALSE, nil
	}
	if _, ok := n.(*ast.Ident); ok {
		return vm.TRUE, nil
	}
	return vm.FALSE, nil
}

// ident-name: (gogen/ident-name ident) -> name of the identifier as a string
func cIdentName(v vm.Value) (vm.Value, error) {
	n, err := unboxNode(v)
	if err != nil {
		return vm.NIL, fmt.Errorf("gogen/ident-name: %v", err)
	}
	if n == nil {
		return vm.NIL, fmt.Errorf("gogen/ident-name: got nil node")
	}
	id, ok := n.(*ast.Ident)
	if !ok {
		return vm.NIL, fmt.Errorf("gogen/ident-name: not an *ast.Ident: %T", n)
	}
	return vm.String(id.Name), nil
}

// type-expr: (gogen/type "spec") -> parsed type expression
// Uses go/parser so the full Go type grammar is supported.
//
// We strip all positions from the parsed expression so its positions
// (which start at 1, relative to the parser's own internal scratch
// FileSet) don't interfere with hand-allocated positions when this
// type is spliced into other nodes.
func cType(v vm.Value) (vm.Value, error) {
	s, err := asString(v)
	if err != nil {
		return vm.NIL, err
	}
	expr, err := parser.ParseExpr(s)
	if err != nil {
		return vm.NIL, fmt.Errorf("gogen: parsing type %q: %w", s, err)
	}
	stripPositions(expr)
	return box(expr), nil
}

// int-lit / float-lit / string-lit / char-lit
func cIntLit(v vm.Value) (vm.Value, error) {
	switch x := v.(type) {
	case vm.Int:
		return box(&ast.BasicLit{Kind: token.INT, Value: strconv.FormatInt(int64(x), 10)}), nil
	}
	return vm.NIL, fmt.Errorf("gogen/int-lit: expected Int, got %s", v.Type().Name())
}

func cFloatLit(v vm.Value) (vm.Value, error) {
	switch x := v.(type) {
	case vm.Float:
		return box(&ast.BasicLit{Kind: token.FLOAT, Value: strconv.FormatFloat(float64(x), 'g', -1, 64)}), nil
	case vm.Int:
		// Allow Int → float literal coercion for ergonomics: (float-lit 0) emits "0.0".
		return box(&ast.BasicLit{Kind: token.FLOAT, Value: strconv.FormatFloat(float64(int64(x)), 'g', -1, 64) + ".0"}), nil
	}
	return vm.NIL, fmt.Errorf("gogen/float-lit: expected number, got %s", v.Type().Name())
}

func cStringLit(v vm.Value) (vm.Value, error) {
	s, err := asString(v)
	if err != nil {
		return vm.NIL, err
	}
	return box(&ast.BasicLit{Kind: token.STRING, Value: strconv.Quote(s)}), nil
}

func cCharLit(v vm.Value) (vm.Value, error) {
	switch x := v.(type) {
	case vm.Char:
		return box(&ast.BasicLit{Kind: token.CHAR, Value: strconv.QuoteRune(rune(x))}), nil
	}
	return vm.NIL, fmt.Errorf("gogen/char-lit: expected Char, got %s", v.Type().Name())
}

// binary: (gogen/binary "+" left right)
func cBinary(opV, leftV, rightV vm.Value) (vm.Value, error) {
	op, err := asString(opV)
	if err != nil {
		return vm.NIL, err
	}
	tok, err := opTokenOrErr(op)
	if err != nil {
		return vm.NIL, err
	}
	l, err := unboxExpr(leftV)
	if err != nil {
		return vm.NIL, fmt.Errorf("gogen/binary: left: %w", err)
	}
	r, err := unboxExpr(rightV)
	if err != nil {
		return vm.NIL, fmt.Errorf("gogen/binary: right: %w", err)
	}
	return box(&ast.BinaryExpr{X: l, Op: tok, Y: r}), nil
}

// unary: (gogen/unary "!" x)  (also "-" "&" "*" "^")
func cUnary(opV, xV vm.Value) (vm.Value, error) {
	op, err := asString(opV)
	if err != nil {
		return vm.NIL, err
	}
	tok, err := opTokenOrErr(op)
	if err != nil {
		return vm.NIL, err
	}
	x, err := unboxExpr(xV)
	if err != nil {
		return vm.NIL, err
	}
	return box(&ast.UnaryExpr{Op: tok, X: x}), nil
}

// index: (gogen/index recv idx)
func cIndex(recvV, idxV vm.Value) (vm.Value, error) {
	r, err := unboxExpr(recvV)
	if err != nil {
		return vm.NIL, err
	}
	i, err := unboxExpr(idxV)
	if err != nil {
		return vm.NIL, err
	}
	return box(&ast.IndexExpr{X: r, Index: i}), nil
}

// field-sel (selector): (gogen/field-sel recv name)
// recv is an expression node, name is a string (must be a valid identifier).
func cFieldSel(recvV, nameV vm.Value) (vm.Value, error) {
	r, err := unboxExpr(recvV)
	if err != nil {
		return vm.NIL, err
	}
	name, err := asString(nameV)
	if err != nil {
		return vm.NIL, err
	}
	if !validIdent(name) {
		return vm.NIL, fmt.Errorf("gogen/field-sel: %q is not a valid identifier", name)
	}
	return box(&ast.SelectorExpr{X: r, Sel: ast.NewIdent(name)}), nil
}

// call: (gogen/call fn-expr [arg-exprs...])
func cCall(fnV, argsV vm.Value) (vm.Value, error) {
	fn, err := unboxExpr(fnV)
	if err != nil {
		return vm.NIL, err
	}
	argVals, err := seqToValues(argsV)
	if err != nil {
		return vm.NIL, err
	}
	args := make([]ast.Expr, 0, len(argVals))
	for i, av := range argVals {
		e, err := unboxExpr(av)
		if err != nil {
			return vm.NIL, fmt.Errorf("gogen/call: arg %d: %w", i, err)
		}
		args = append(args, e)
	}
	return box(&ast.CallExpr{Fun: fn, Args: args}), nil
}

// cast: T(x) — same as call with a type expr as fn.
func cCast(typeV, xV vm.Value) (vm.Value, error) {
	t, err := unboxExpr(typeV)
	if err != nil {
		return vm.NIL, err
	}
	x, err := unboxExpr(xV)
	if err != nil {
		return vm.NIL, err
	}
	return box(&ast.CallExpr{Fun: t, Args: []ast.Expr{x}}), nil
}

// type-assert: (gogen/type-assert x type) -> x.(type) expression
// Used in if-init clauses: `if ai, ok := a.(Int); ok { ... }`.
func cTypeAssert(xV, typeV vm.Value) (vm.Value, error) {
	x, err := unboxExpr(xV)
	if err != nil {
		return vm.NIL, err
	}
	t, err := unboxExpr(typeV)
	if err != nil {
		return vm.NIL, err
	}
	return box(&ast.TypeAssertExpr{X: x, Type: t}), nil
}

// multi-assign: (gogen/multi-assign "=" [lhs-exprs] [rhs-exprs]) — supports
// any number of LHS/RHS expressions. Use this for `:=` (short var decl) and
// tuple-style assignment.
//
// Examples:
//
//	(multi-assign ":=" [ai ok] [(type-assert a (type "Int"))])
//	(multi-assign "=" [a b] [b a])  // swap
func cMultiAssign(opV, lhsV, rhsV vm.Value) (vm.Value, error) {
	op, err := asString(opV)
	if err != nil {
		return vm.NIL, err
	}
	tok, err := opTokenOrErr(op)
	if err != nil {
		return vm.NIL, err
	}
	lhsVals, err := seqToValues(lhsV)
	if err != nil {
		return vm.NIL, fmt.Errorf("gogen/multi-assign: lhs: %w", err)
	}
	rhsVals, err := seqToValues(rhsV)
	if err != nil {
		return vm.NIL, fmt.Errorf("gogen/multi-assign: rhs: %w", err)
	}
	lhs := make([]ast.Expr, 0, len(lhsVals))
	for i, v := range lhsVals {
		e, err := unboxExpr(v)
		if err != nil {
			return vm.NIL, fmt.Errorf("gogen/multi-assign: lhs[%d]: %w", i, err)
		}
		lhs = append(lhs, e)
	}
	rhs := make([]ast.Expr, 0, len(rhsVals))
	for i, v := range rhsVals {
		e, err := unboxExpr(v)
		if err != nil {
			return vm.NIL, fmt.Errorf("gogen/multi-assign: rhs[%d]: %w", i, err)
		}
		rhs = append(rhs, e)
	}
	return box(&ast.AssignStmt{Lhs: lhs, Tok: tok, Rhs: rhs}), nil
}

// var-decl: (gogen/var-decl "name" type-expr init-expr-or-nil) -> ast.DeclStmt
func cVarDecl(nameV, typeV, initV vm.Value) (vm.Value, error) {
	name, err := asString(nameV)
	if err != nil {
		return vm.NIL, err
	}
	if !validIdent(name) {
		return vm.NIL, fmt.Errorf("gogen/var-decl: %q is not a valid identifier", name)
	}
	t, err := unboxExpr(typeV)
	if err != nil {
		return vm.NIL, err
	}
	spec := &ast.ValueSpec{
		Names: []*ast.Ident{ast.NewIdent(name)},
		Type:  t,
	}
	if initV != vm.NIL {
		init, err := unboxExpr(initV)
		if err != nil {
			return vm.NIL, err
		}
		spec.Values = []ast.Expr{init}
	}
	return box(&ast.DeclStmt{Decl: &ast.GenDecl{Tok: token.VAR, Specs: []ast.Spec{spec}}}), nil
}

// assign: (gogen/assign "=" lhs-expr rhs-expr)  — also "+=" etc.
func cAssign(opV, lhsV, rhsV vm.Value) (vm.Value, error) {
	op, err := asString(opV)
	if err != nil {
		return vm.NIL, err
	}
	tok, err := opTokenOrErr(op)
	if err != nil {
		return vm.NIL, err
	}
	lhs, err := unboxExpr(lhsV)
	if err != nil {
		return vm.NIL, err
	}
	rhs, err := unboxExpr(rhsV)
	if err != nil {
		return vm.NIL, err
	}
	return box(&ast.AssignStmt{Lhs: []ast.Expr{lhs}, Tok: tok, Rhs: []ast.Expr{rhs}}), nil
}

// return: (gogen/return-stmt [exprs])
func cReturn(valsV vm.Value) (vm.Value, error) {
	vs, err := seqToValues(valsV)
	if err != nil {
		return vm.NIL, err
	}
	exprs := make([]ast.Expr, 0, len(vs))
	for i, v := range vs {
		e, err := unboxExpr(v)
		if err != nil {
			return vm.NIL, fmt.Errorf("gogen/return-stmt: value %d: %w", i, err)
		}
		exprs = append(exprs, e)
	}
	return box(&ast.ReturnStmt{Results: exprs}), nil
}

// if: (gogen/if-stmt init-or-nil cond [then-stmts] else-or-nil)
// init covers the common `if x, ok := foo(); ok { ... }` pattern.
// Pass nil for init when not needed. else may be nil, a single stmt, or
// a sequence of stmts.
func cIfStmt(initV, condV, thenV, elseV vm.Value) (vm.Value, error) {
	var init ast.Stmt
	if initV != vm.NIL {
		s, err := unboxStmt(initV)
		if err != nil {
			return vm.NIL, fmt.Errorf("gogen/if-stmt: init: %w", err)
		}
		init = s
	}
	cond, err := unboxExpr(condV)
	if err != nil {
		return vm.NIL, fmt.Errorf("gogen/if-stmt: cond: %w", err)
	}
	thenStmts, err := stmtSlice(thenV)
	if err != nil {
		return vm.NIL, err
	}
	stmt := &ast.IfStmt{Init: init, Cond: cond, Body: &ast.BlockStmt{List: thenStmts}}
	if elseV != vm.NIL {
		elseStmts, err := stmtSlice(elseV)
		if err != nil {
			return vm.NIL, err
		}
		stmt.Else = &ast.BlockStmt{List: elseStmts}
	}
	return box(stmt), nil
}

// for: (gogen/for-stmt init-or-nil cond-or-nil post-or-nil [body-stmts])
// All three loop clauses are optional (Go allows `for cond {}` and `for {}`).
func cForStmt(initV, condV, postV, bodyV vm.Value) (vm.Value, error) {
	var init, post ast.Stmt
	var cond ast.Expr
	if initV != vm.NIL {
		s, err := unboxStmt(initV)
		if err != nil {
			return vm.NIL, err
		}
		init = s
	}
	if condV != vm.NIL {
		c, err := unboxExpr(condV)
		if err != nil {
			return vm.NIL, err
		}
		cond = c
	}
	if postV != vm.NIL {
		s, err := unboxStmt(postV)
		if err != nil {
			return vm.NIL, err
		}
		post = s
	}
	body, err := stmtSlice(bodyV)
	if err != nil {
		return vm.NIL, err
	}
	return box(&ast.ForStmt{
		Init: init,
		Cond: cond,
		Post: post,
		Body: &ast.BlockStmt{List: body},
	}), nil
}

func cExprStmt(v vm.Value) (vm.Value, error) {
	e, err := unboxExpr(v)
	if err != nil {
		return vm.NIL, err
	}
	return box(&ast.ExprStmt{X: e}), nil
}

// block-stmt: (gogen/block-stmt [stmts]) -> a bare `{ … }` block, introducing a
// new lexical scope. Use to scope short-lived locals so a `goto` elsewhere in
// the function does not jump over their declarations (Go forbids that).
func cBlockStmt(v vm.Value) (vm.Value, error) {
	stmts, err := stmtSlice(v)
	if err != nil {
		return vm.NIL, err
	}
	return box(&ast.BlockStmt{List: stmts}), nil
}

func cGotoStmt(nameV vm.Value) (vm.Value, error) {
	name, err := asString(nameV)
	if err != nil {
		return vm.NIL, err
	}
	if !validIdent(name) {
		return vm.NIL, fmt.Errorf("gogen/goto-stmt: %q is not a valid identifier", name)
	}
	return box(&ast.BranchStmt{Tok: token.GOTO, Label: ast.NewIdent(name)}), nil
}

func cContinueStmt() (vm.Value, error) {
	return box(&ast.BranchStmt{Tok: token.CONTINUE}), nil
}

func cBreakStmt() (vm.Value, error) {
	return box(&ast.BranchStmt{Tok: token.BREAK}), nil
}

func cLabelStmt(nameV, stmtV vm.Value) (vm.Value, error) {
	name, err := asString(nameV)
	if err != nil {
		return vm.NIL, err
	}
	if !validIdent(name) {
		return vm.NIL, fmt.Errorf("gogen/label-stmt: %q is not a valid identifier", name)
	}
	var stmt ast.Stmt
	if stmtV == vm.NIL {
		stmt = &ast.EmptyStmt{Implicit: true}
	} else {
		s, err := unboxStmt(stmtV)
		if err != nil {
			return vm.NIL, fmt.Errorf("gogen/label-stmt: stmt: %w", err)
		}
		stmt = s
	}
	return box(&ast.LabeledStmt{Label: ast.NewIdent(name), Stmt: stmt}), nil
}

// stmtSlice accepts either a single boxed stmt or a sequable of boxed stmts.
func stmtSlice(v vm.Value) ([]ast.Stmt, error) {
	if v == vm.NIL {
		return nil, nil
	}
	if _, ok := v.(*goASTValue); ok {
		// single statement
		s, err := unboxStmt(v)
		if err != nil {
			return nil, err
		}
		return []ast.Stmt{s}, nil
	}
	vs, err := seqToValues(v)
	if err != nil {
		return nil, err
	}
	out := make([]ast.Stmt, 0, len(vs))
	for i, vv := range vs {
		s, err := unboxStmt(vv)
		if err != nil {
			return nil, fmt.Errorf("stmt %d: %w", i, err)
		}
		out = append(out, s)
	}
	return out, nil
}

// param: (gogen/param "name" type-expr)
// Wraps as an *ast.Field with a single name.
func cParam(nameV, typeV vm.Value) (vm.Value, error) {
	name, err := asString(nameV)
	if err != nil {
		return vm.NIL, err
	}
	if !validIdent(name) {
		return vm.NIL, fmt.Errorf("gogen/param: %q is not a valid identifier", name)
	}
	t, err := unboxExpr(typeV)
	if err != nil {
		return vm.NIL, err
	}
	return box(&ast.Field{
		Names: []*ast.Ident{ast.NewIdent(name)},
		Type:  t,
	}), nil
}

// variadic-param: (gogen/variadic-param "name" type-expr)
// Like param but wraps the type in *ast.Ellipsis for Go variadic params (e.g. args ...int).
func cVariadicParam(nameV, typeV vm.Value) (vm.Value, error) {
	name, err := asString(nameV)
	if err != nil {
		return vm.NIL, err
	}
	if !validIdent(name) {
		return vm.NIL, fmt.Errorf("gogen/variadic-param: %q is not a valid identifier", name)
	}
	t, err := unboxExpr(typeV)
	if err != nil {
		return vm.NIL, err
	}
	return box(&ast.Field{
		Names: []*ast.Ident{ast.NewIdent(name)},
		Type:  &ast.Ellipsis{Elt: t},
	}), nil
}

// result: (gogen/result type-expr)
// Anonymous result for multi-return signatures.
func cResult(typeV vm.Value) (vm.Value, error) {
	t, err := unboxExpr(typeV)
	if err != nil {
		return vm.NIL, err
	}
	return box(&ast.Field{Type: t}), nil
}

// func-decl: (gogen/func-decl "name" [params] [results] [body])
// results may be empty (void function).
func cFuncDecl(nameV, paramsV, resultsV, bodyV vm.Value) (vm.Value, error) {
	name, err := asString(nameV)
	if err != nil {
		return vm.NIL, err
	}
	if !validIdent(name) {
		return vm.NIL, fmt.Errorf("gogen/func-decl: %q is not a valid identifier", name)
	}
	params, err := fieldSlice(paramsV)
	if err != nil {
		return vm.NIL, fmt.Errorf("gogen/func-decl: params: %w", err)
	}
	results, err := fieldSlice(resultsV)
	if err != nil {
		return vm.NIL, fmt.Errorf("gogen/func-decl: results: %w", err)
	}
	body, err := stmtSlice(bodyV)
	if err != nil {
		return vm.NIL, fmt.Errorf("gogen/func-decl: body: %w", err)
	}
	funcType := &ast.FuncType{
		Params: &ast.FieldList{List: params},
	}
	if len(results) > 0 {
		funcType.Results = &ast.FieldList{List: results}
	}
	return box(&ast.FuncDecl{
		Name: ast.NewIdent(name),
		Type: funcType,
		Body: &ast.BlockStmt{List: body},
	}), nil
}

func fieldSlice(v vm.Value) ([]*ast.Field, error) {
	if v == vm.NIL {
		return nil, nil
	}
	vs, err := seqToValues(v)
	if err != nil {
		return nil, err
	}
	out := make([]*ast.Field, 0, len(vs))
	for i, vv := range vs {
		n, err := unboxNode(vv)
		if err != nil {
			return nil, fmt.Errorf("field %d: %w", i, err)
		}
		f, ok := n.(*ast.Field)
		if !ok {
			return nil, fmt.Errorf("field %d: expected *ast.Field, got %T", i, n)
		}
		out = append(out, f)
	}
	return out, nil
}

// func-lit: (gogen/func-lit [params] [results-or-nil] [body-stmts]) -> *ast.FuncLit
// A function literal — an anonymous function used as an expression
// (the value, not the declaration). Returns an expression, so it can
// be passed as a call argument, assigned, etc.
//
// Example uses:
//
//	cb := func(x int) int { return x * 2 }
//	vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) { ... })
func cFuncLit(paramsV, resultsV, bodyV vm.Value) (vm.Value, error) {
	params, err := fieldSlice(paramsV)
	if err != nil {
		return vm.NIL, fmt.Errorf("gogen/func-lit: params: %w", err)
	}
	results, err := fieldSlice(resultsV)
	if err != nil {
		return vm.NIL, fmt.Errorf("gogen/func-lit: results: %w", err)
	}
	body, err := stmtSlice(bodyV)
	if err != nil {
		return vm.NIL, fmt.Errorf("gogen/func-lit: body: %w", err)
	}
	funcType := &ast.FuncType{
		Params: &ast.FieldList{List: params},
	}
	if len(results) > 0 {
		funcType.Results = &ast.FieldList{List: results}
	}
	return box(&ast.FuncLit{
		Type: funcType,
		Body: &ast.BlockStmt{List: body},
	}), nil
}

// case-clause: (gogen/case-clause [values] [body-stmts]) -> *ast.CaseClause
// A single arm of a switch statement. Empty values means the default arm.
// Multiple values produce the `case a, b, c:` form.
func cCaseClause(valuesV, bodyV vm.Value) (vm.Value, error) {
	var vals []ast.Expr
	if valuesV != vm.NIL {
		vs, err := seqToValues(valuesV)
		if err != nil {
			return vm.NIL, fmt.Errorf("gogen/case-clause: values: %w", err)
		}
		vals = make([]ast.Expr, 0, len(vs))
		for i, v := range vs {
			e, err := unboxExpr(v)
			if err != nil {
				return vm.NIL, fmt.Errorf("gogen/case-clause: value %d: %w", i, err)
			}
			vals = append(vals, e)
		}
	}
	body, err := stmtSlice(bodyV)
	if err != nil {
		return vm.NIL, fmt.Errorf("gogen/case-clause: body: %w", err)
	}
	// Note: go/ast represents the default clause as a CaseClause with
	// List==nil (not an empty slice). Normalize empty-slice → nil so the
	// renderer emits `default:` instead of `case :`.
	if len(vals) == 0 {
		vals = nil
	}
	return box(&ast.CaseClause{List: vals, Body: body}), nil
}

// switch-stmt: (gogen/switch-stmt init-or-nil tag-or-nil [case-clauses]) -> *ast.SwitchStmt
// A Go `switch` statement. Pass nil for `tag` to get the tagless form
// (`switch { case cond1: ... case cond2: ... }`). Each clause must be
// a *ast.CaseClause (produced by cCaseClause).
func cSwitchStmt(initV, tagV, clausesV vm.Value) (vm.Value, error) {
	var init ast.Stmt
	if initV != vm.NIL {
		s, err := unboxStmt(initV)
		if err != nil {
			return vm.NIL, fmt.Errorf("gogen/switch-stmt: init: %w", err)
		}
		init = s
	}
	var tag ast.Expr
	if tagV != vm.NIL {
		e, err := unboxExpr(tagV)
		if err != nil {
			return vm.NIL, fmt.Errorf("gogen/switch-stmt: tag: %w", err)
		}
		tag = e
	}
	clauseVals, err := seqToValues(clausesV)
	if err != nil {
		return vm.NIL, fmt.Errorf("gogen/switch-stmt: clauses: %w", err)
	}
	body := make([]ast.Stmt, 0, len(clauseVals))
	for i, cv := range clauseVals {
		n, err := unboxNode(cv)
		if err != nil {
			return vm.NIL, fmt.Errorf("gogen/switch-stmt: clause %d: %w", i, err)
		}
		cc, ok := n.(*ast.CaseClause)
		if !ok {
			return vm.NIL, fmt.Errorf("gogen/switch-stmt: clause %d: expected *ast.CaseClause, got %T", i, n)
		}
		body = append(body, cc)
	}
	return box(&ast.SwitchStmt{
		Init: init,
		Tag:  tag,
		Body: &ast.BlockStmt{List: body},
	}), nil
}

// kv-expr: (gogen/kv-expr key value) -> *ast.KeyValueExpr
// A key/value pair, used inside composite literals (map and struct).
func cKVExpr(keyV, valueV vm.Value) (vm.Value, error) {
	k, err := unboxExpr(keyV)
	if err != nil {
		return vm.NIL, fmt.Errorf("gogen/kv-expr: key: %w", err)
	}
	v, err := unboxExpr(valueV)
	if err != nil {
		return vm.NIL, fmt.Errorf("gogen/kv-expr: value: %w", err)
	}
	return box(&ast.KeyValueExpr{Key: k, Value: v}), nil
}

// composite-lit: (gogen/composite-lit type-or-nil [elements]) -> *ast.CompositeLit
// Composite literal, e.g. `[]int{1,2,3}`, `Point{X:1,Y:2}`, `map[string]int{"a":1}`.
//
// type-or-nil is the explicit type (a type expression). It can be nil
// in contexts where Go infers the element type (e.g. nested literals).
//
// elements is a sequence of expression nodes. For map/struct-style
// literals, each element should be a *ast.KeyValueExpr (built via
// cKVExpr). For slice/array-style literals, elements are bare exprs.
// The two styles can be mixed only as Go permits (struct field omission etc.).
func cCompositeLit(typeV, elementsV vm.Value) (vm.Value, error) {
	var typ ast.Expr
	if typeV != vm.NIL {
		t, err := unboxExpr(typeV)
		if err != nil {
			return vm.NIL, fmt.Errorf("gogen/composite-lit: type: %w", err)
		}
		typ = t
	}
	var elts []ast.Expr
	if elementsV != vm.NIL {
		evals, err := seqToValues(elementsV)
		if err != nil {
			return vm.NIL, fmt.Errorf("gogen/composite-lit: elements: %w", err)
		}
		elts = make([]ast.Expr, 0, len(evals))
		for i, ev := range evals {
			e, err := unboxExpr(ev)
			if err != nil {
				return vm.NIL, fmt.Errorf("gogen/composite-lit: element %d: %w", i, err)
			}
			elts = append(elts, e)
		}
	}
	return box(&ast.CompositeLit{Type: typ, Elts: elts}), nil
}

// composite-lit-multi: (gogen/composite-lit-multi type-or-nil [elements])
// -> *ast.CompositeLit, with Lbrace/Rbrace minted from the synthetic
// FileSet and each element placed on a distinct line.
//
// This is the variant to use for file-level array/slice/map literals
// that should render multi-line. gogen/composite-lit (the basic
// variant) leaves Lbrace/Rbrace at NoPos, which can collapse to a
// single line — fine for nested literals, wrong for top-level tables.
//
// Position scheme: Lbrace one line before the first element; each
// element on a fresh line; Rbrace one line past the last element.
// Keeping Lbrace close to the first element prevents the printer
// from inserting blank lines after the `{`.
func cCompositeLitMulti(typeV, elementsV vm.Value) (vm.Value, error) {
	var typ ast.Expr
	if typeV != vm.NIL {
		t, err := unboxExpr(typeV)
		if err != nil {
			return vm.NIL, fmt.Errorf("gogen/composite-lit-multi: type: %w", err)
		}
		typ = t
	}
	var elts []ast.Expr
	var minPos, maxPos token.Pos
	if elementsV != vm.NIL {
		evals, err := seqToValues(elementsV)
		if err != nil {
			return vm.NIL, fmt.Errorf("gogen/composite-lit-multi: elements: %w", err)
		}
		elts = make([]ast.Expr, 0, len(evals))
		for i, ev := range evals {
			e, err := unboxExpr(ev)
			if err != nil {
				return vm.NIL, fmt.Errorf("gogen/composite-lit-multi: element %d: %w", i, err)
			}
			// Place each element on a distinct line so the printer
			// breaks them apart. Rewrite the lead-token position of
			// the element when it has no position yet.
			linePos := allocPos()
			switch x := e.(type) {
			case *ast.KeyValueExpr:
				if id, ok := x.Key.(*ast.Ident); ok && id.NamePos == token.NoPos {
					id.NamePos = linePos
				}
			case *ast.Ident:
				if x.NamePos == token.NoPos {
					x.NamePos = linePos
				}
			}
			if minPos == 0 || linePos < minPos {
				minPos = linePos
			}
			if linePos > maxPos {
				maxPos = linePos
			}
			elts = append(elts, e)
		}
	}
	lbrace := token.Pos(1)
	if minPos > 1 {
		lbrace = minPos - 1
	}
	rbrace := maxPos + 1
	if rbrace == 0 {
		rbrace = allocPos()
	}
	return box(&ast.CompositeLit{
		Type:   typ,
		Lbrace: lbrace,
		Elts:   elts,
		Rbrace: rbrace,
	}), nil
}

// --- top-level decl constructors (type/const/var/method) -------------

// type-decl: (gogen/type-decl "Name" type-expr) -> *ast.GenDecl
// Top-level type declaration, e.g. `type Op uint16` or
// `type opInfo struct { ... }`. The Type slot can be any type
// expression (identifier, struct, interface, channel, etc.).
func cTypeDecl(nameV, typeV vm.Value) (vm.Value, error) {
	name, err := asString(nameV)
	if err != nil {
		return vm.NIL, err
	}
	if !validIdent(name) {
		return vm.NIL, fmt.Errorf("gogen/type-decl: %q is not a valid identifier", name)
	}
	t, err := unboxExpr(typeV)
	if err != nil {
		return vm.NIL, fmt.Errorf("gogen/type-decl: type: %w", err)
	}
	return box(&ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{&ast.TypeSpec{
			Name: ast.NewIdent(name),
			Type: t,
		}},
	}), nil
}

// field-decl: (gogen/field-decl "name" type-expr "trailing comment") -> *ast.Field
// A struct field. Comment, if non-empty, is attached as a trailing
// `// comment` line comment via positioned *ast.CommentGroup.
func cFieldDecl(nameV, typeV, commentV vm.Value) (vm.Value, error) {
	name, err := asString(nameV)
	if err != nil {
		return vm.NIL, err
	}
	if !validIdent(name) {
		return vm.NIL, fmt.Errorf("gogen/field-decl: %q is not a valid identifier", name)
	}
	t, err := unboxExpr(typeV)
	if err != nil {
		return vm.NIL, fmt.Errorf("gogen/field-decl: type: %w", err)
	}
	comment, err := asString(commentV)
	if err != nil {
		return vm.NIL, fmt.Errorf("gogen/field-decl: comment: %w", err)
	}
	// Allocate a fresh line for this field so adjacent fields land on
	// separate lines in the rendered output.
	pos := allocPos()
	f := &ast.Field{
		Names: []*ast.Ident{{Name: name, NamePos: pos}},
		Type:  t,
	}
	if comment != "" {
		f.Comment = &ast.CommentGroup{List: []*ast.Comment{
			{Slash: pos, Text: "// " + comment},
		}}
	}
	return box(f), nil
}

// embed-field: (gogen/embed-field type-expr) -> *ast.Field
// An embedded (anonymous) struct field, e.g. `vm.RecordBase` inside a struct.
// The field has no Names, so the type itself is promoted — the mechanism by
// which a generated record struct inherits vm.RecordBase's vm.Value contract.
func cEmbedField(typeV vm.Value) (vm.Value, error) {
	t, err := unboxExpr(typeV)
	if err != nil {
		return vm.NIL, fmt.Errorf("gogen/embed-field: type: %w", err)
	}
	// Allocate a fresh line so the embedded field lands on its own line.
	pos := allocPos()
	if id, ok := t.(*ast.Ident); ok {
		id.NamePos = pos
	}
	return box(&ast.Field{Type: t}), nil
}

// struct-type: (gogen/struct-type [fields]) -> *ast.StructType
// A struct type expression, suitable as the type slot in a type-decl
// or as a nested type.
//
// Position scheme: Struct/Opening go one line BEFORE the first field;
// Closing goes ONE LINE AFTER the last field. This makes the printer
// emit each field on its own line without leading or trailing blank
// lines inside the braces.
func cStructType(fieldsV vm.Value) (vm.Value, error) {
	fields, err := fieldSlice(fieldsV)
	if err != nil {
		return vm.NIL, fmt.Errorf("gogen/struct-type: fields: %w", err)
	}
	// Find min/max position among fields to bracket Opening/Closing.
	var minPos, maxPos token.Pos
	for _, f := range fields {
		for _, n := range f.Names {
			if minPos == 0 || n.NamePos < minPos {
				minPos = n.NamePos
			}
			if n.NamePos > maxPos {
				maxPos = n.NamePos
			}
		}
	}
	openPos := token.Pos(1)
	if minPos > 1 {
		openPos = minPos - 1
	}
	closePos := maxPos + 1
	if closePos == 0 {
		closePos = allocPos()
	}
	return box(&ast.StructType{
		Struct: openPos,
		Fields: &ast.FieldList{
			Opening: openPos,
			List:    fields,
			Closing: closePos,
		},
	}), nil
}

// iface-method: (gogen/iface-method "name" [params] [results]) -> *ast.Field
// An interface method signature. Returns an *ast.Field with a FuncType.
// The field's Name list contains the method name, and the Type is a
// *ast.FuncType with Params and Results.
func cIfaceMethod(nameV, paramsV, resultsV vm.Value) (vm.Value, error) {
	name, err := asString(nameV)
	if err != nil {
		return vm.NIL, err
	}
	if !validIdent(name) {
		return vm.NIL, fmt.Errorf("gogen/iface-method: %q is not a valid identifier", name)
	}
	params, err := fieldSlice(paramsV)
	if err != nil {
		return vm.NIL, fmt.Errorf("gogen/iface-method: params: %w", err)
	}
	results, err := fieldSlice(resultsV)
	if err != nil {
		return vm.NIL, fmt.Errorf("gogen/iface-method: results: %w", err)
	}
	// Allocate a fresh position for the method name.
	pos := allocPos()
	funcType := &ast.FuncType{
		Params: &ast.FieldList{List: params},
	}
	if len(results) > 0 {
		funcType.Results = &ast.FieldList{List: results}
	}
	return box(&ast.Field{
		Names: []*ast.Ident{{Name: name, NamePos: pos}},
		Type:  funcType,
	}), nil
}

// interface-type: (gogen/interface-type [methods]) -> *ast.InterfaceType
// An interface type expression, suitable as the type slot in a type-decl
// or as a nested type.
//
// Position scheme: Interface/Opening go one line BEFORE the first method;
// Closing goes ONE LINE AFTER the last method. This makes the printer
// emit each method on its own line without leading or trailing blank
// lines inside the braces.
func cInterfaceType(methodsV vm.Value) (vm.Value, error) {
	methods, err := fieldSlice(methodsV)
	if err != nil {
		return vm.NIL, fmt.Errorf("gogen/interface-type: methods: %w", err)
	}
	// Find min/max position among methods to bracket Opening/Closing.
	var minPos, maxPos token.Pos
	for _, m := range methods {
		for _, n := range m.Names {
			if minPos == 0 || n.NamePos < minPos {
				minPos = n.NamePos
			}
			if n.NamePos > maxPos {
				maxPos = n.NamePos
			}
		}
	}
	openPos := token.Pos(1)
	if minPos > 1 {
		openPos = minPos - 1
	}
	closePos := maxPos + 1
	if closePos == 0 {
		closePos = allocPos()
	}
	return box(&ast.InterfaceType{
		Interface: openPos,
		Methods: &ast.FieldList{
			Opening: openPos,
			List:    methods,
			Closing: closePos,
		},
	}), nil
}

// const-spec: (gogen/const-spec "name" type-or-nil value-or-nil "trailing comment")
// -> *ast.ValueSpec
// One row inside a const block. Type and value may both be nil for
// rows that inherit from the previous row (Go semantics).
func cConstSpec(nameV, typeV, valueV, commentV vm.Value) (vm.Value, error) {
	name, err := asString(nameV)
	if err != nil {
		return vm.NIL, err
	}
	if !validIdent(name) {
		return vm.NIL, fmt.Errorf("gogen/const-spec: %q is not a valid identifier", name)
	}
	var t ast.Expr
	if typeV != vm.NIL {
		te, err := unboxExpr(typeV)
		if err != nil {
			return vm.NIL, fmt.Errorf("gogen/const-spec: type: %w", err)
		}
		t = te
	}
	var values []ast.Expr
	if valueV != vm.NIL {
		ve, err := unboxExpr(valueV)
		if err != nil {
			return vm.NIL, fmt.Errorf("gogen/const-spec: value: %w", err)
		}
		values = []ast.Expr{ve}
	}
	comment, err := asString(commentV)
	if err != nil {
		return vm.NIL, fmt.Errorf("gogen/const-spec: comment: %w", err)
	}
	pos := allocPos()
	spec := &ast.ValueSpec{
		Names:  []*ast.Ident{{Name: name, NamePos: pos}},
		Type:   t,
		Values: values,
	}
	if comment != "" {
		spec.Comment = &ast.CommentGroup{List: []*ast.Comment{
			{Slash: pos, Text: "// " + comment},
		}}
	}
	return box(spec), nil
}

// const-block: (gogen/const-block [specs]) -> *ast.GenDecl
// A parenthesized const declaration. Lparen/Rparen are set to non-zero
// positions so the printer renders as a parenthesized block rather than
// collapsing onto a single line.
//
// Position scheme: TokPos/Lparen one line BEFORE the first spec;
// Rparen one line past the last spec. This avoids leading and trailing
// blank lines inside the block.
func cConstBlock(specsV vm.Value) (vm.Value, error) {
	vs, err := seqToValues(specsV)
	if err != nil {
		return vm.NIL, fmt.Errorf("gogen/const-block: specs: %w", err)
	}
	specs := make([]ast.Spec, 0, len(vs))
	var minPos, maxPos token.Pos
	for i, v := range vs {
		n, err := unboxNode(v)
		if err != nil {
			return vm.NIL, fmt.Errorf("gogen/const-block: spec %d: %w", i, err)
		}
		s, ok := n.(*ast.ValueSpec)
		if !ok {
			return vm.NIL, fmt.Errorf("gogen/const-block: spec %d: expected *ast.ValueSpec, got %T", i, n)
		}
		for _, name := range s.Names {
			if minPos == 0 || name.NamePos < minPos {
				minPos = name.NamePos
			}
			if name.NamePos > maxPos {
				maxPos = name.NamePos
			}
		}
		specs = append(specs, s)
	}
	lparen := token.Pos(1)
	if minPos > 1 {
		lparen = minPos - 1
	}
	rparen := maxPos + 1
	if rparen == 0 {
		rparen = allocPos()
	}
	return box(&ast.GenDecl{
		Tok:    token.CONST,
		TokPos: lparen,
		Lparen: lparen, // non-zero -> parenthesized form
		Specs:  specs,
		Rparen: rparen,
	}), nil
}

// top-var-decl: (gogen/top-var-decl "name" type-expr init-expr-or-nil)
// -> *ast.GenDecl (decl-level var declaration)
//
// Distinct from gogen/var-decl (which builds a *statement-level* DeclStmt
// for use inside function bodies). top-var-decl produces a decl that can
// live at file scope.
func cTopVarDecl(nameV, typeV, initV vm.Value) (vm.Value, error) {
	name, err := asString(nameV)
	if err != nil {
		return vm.NIL, err
	}
	if !validIdent(name) {
		return vm.NIL, fmt.Errorf("gogen/top-var-decl: %q is not a valid identifier", name)
	}
	var t ast.Expr
	if typeV != vm.NIL {
		te, err := unboxExpr(typeV)
		if err != nil {
			return vm.NIL, fmt.Errorf("gogen/top-var-decl: type: %w", err)
		}
		t = te
	}
	var values []ast.Expr
	if initV != vm.NIL {
		ve, err := unboxExpr(initV)
		if err != nil {
			return vm.NIL, fmt.Errorf("gogen/top-var-decl: init: %w", err)
		}
		values = []ast.Expr{ve}
	}
	return box(&ast.GenDecl{
		Tok: token.VAR,
		Specs: []ast.Spec{&ast.ValueSpec{
			Names:  []*ast.Ident{ast.NewIdent(name)},
			Type:   t,
			Values: values,
		}},
	}), nil
}

// with-doc: (gogen/with-doc decl-or-spec "comment text") -> same node with
// a leading Doc comment attached. Accepts either a *ast.GenDecl,
// *ast.FuncDecl, *ast.TypeSpec, *ast.ValueSpec, or *ast.Field.
//
// Multi-line comments are split on newlines; each line becomes a
// separate "// line" comment in the group.
//
// Positioning is delicate: the Doc comment's Slash positions must be
// BEFORE the decl's first token (the printer uses this to decide
// "comment above decl" vs "comment trailing previous decl"). We mint a
// fresh block of positions and re-anchor the decl's leading position
// just past them.
func cWithDoc(declV, commentV vm.Value) (vm.Value, error) {
	if declV == vm.NIL {
		return vm.NIL, nil
	}
	comment, err := asString(commentV)
	if err != nil {
		return vm.NIL, fmt.Errorf("gogen/with-doc: comment: %w", err)
	}
	if comment == "" {
		return declV, nil
	}
	n, err := unboxNode(declV)
	if err != nil {
		return vm.NIL, err
	}
	lines := strings.Split(comment, "\n")
	startPos := allocPos()
	for i := 1; i < len(lines); i++ {
		_ = allocPos() // reserve one fresh line per additional comment line
	}
	declPos := allocPos() // line immediately after the comment block
	cmts := make([]*ast.Comment, 0, len(lines))
	for i, ln := range lines {
		cmts = append(cmts, &ast.Comment{
			Slash: startPos + token.Pos(i),
			Text:  "// " + ln,
		})
	}
	group := &ast.CommentGroup{List: cmts}
	switch x := n.(type) {
	case *ast.GenDecl:
		x.Doc = group
		x.TokPos = declPos
		// Keep Lparen consistent with TokPos when it was already set.
		if x.Lparen != token.NoPos {
			x.Lparen = declPos
		}
	case *ast.FuncDecl:
		x.Doc = group
		// FuncDecl.Pos() looks at Type.Func first (the `func` keyword
		// position). Set that to land after the comment block so the
		// printer places the doc directly above. Also push Recv and
		// Name positions forward if they predate declPos.
		if x.Type != nil {
			x.Type.Func = declPos
		}
		if x.Recv != nil {
			x.Recv.Opening = declPos
			for _, f := range x.Recv.List {
				for _, nm := range f.Names {
					if nm.NamePos < declPos {
						nm.NamePos = declPos
					}
				}
			}
		}
		if x.Name != nil && x.Name.NamePos < declPos {
			x.Name.NamePos = declPos
		}
	case *ast.TypeSpec:
		x.Doc = group
		if x.Name != nil && x.Name.NamePos < declPos {
			x.Name.NamePos = declPos
		}
	case *ast.ValueSpec:
		x.Doc = group
		for _, nm := range x.Names {
			if nm.NamePos < declPos {
				nm.NamePos = declPos
			}
		}
	case *ast.Field:
		x.Doc = group
		for _, nm := range x.Names {
			if nm.NamePos < declPos {
				nm.NamePos = declPos
			}
		}
	default:
		return vm.NIL, fmt.Errorf("gogen/with-doc: don't know how to attach doc to %T", n)
	}
	return declV, nil
}

// method-decl: (gogen/method-decl [recv-fields] "Name" [params] [results] [body])
// -> *ast.FuncDecl with Recv set.
//
// Same shape as func-decl but adds a receiver. The receiver is a list
// of one field (typically one named binding to a named type), e.g.
// `(op Op)`.
func cMethodDecl(recvV, nameV, paramsV, resultsV, bodyV vm.Value) (vm.Value, error) {
	recv, err := fieldSlice(recvV)
	if err != nil {
		return vm.NIL, fmt.Errorf("gogen/method-decl: recv: %w", err)
	}
	if len(recv) == 0 {
		return vm.NIL, fmt.Errorf("gogen/method-decl: receiver list must contain at least one field")
	}
	name, err := asString(nameV)
	if err != nil {
		return vm.NIL, err
	}
	if !validIdent(name) {
		return vm.NIL, fmt.Errorf("gogen/method-decl: %q is not a valid identifier", name)
	}
	params, err := fieldSlice(paramsV)
	if err != nil {
		return vm.NIL, fmt.Errorf("gogen/method-decl: params: %w", err)
	}
	results, err := fieldSlice(resultsV)
	if err != nil {
		return vm.NIL, fmt.Errorf("gogen/method-decl: results: %w", err)
	}
	body, err := stmtSlice(bodyV)
	if err != nil {
		return vm.NIL, fmt.Errorf("gogen/method-decl: body: %w", err)
	}
	funcType := &ast.FuncType{
		Params: &ast.FieldList{List: params},
	}
	if len(results) > 0 {
		funcType.Results = &ast.FieldList{List: results}
	}
	return box(&ast.FuncDecl{
		Recv: &ast.FieldList{List: recv},
		Name: ast.NewIdent(name),
		Type: funcType,
		Body: &ast.BlockStmt{List: body},
	}), nil
}

// import-spec: (gogen/import-spec "path") or (gogen/import-spec "path" "alias")
func cImportSpec(args ...vm.Value) (vm.Value, error) {
	if len(args) != 1 && len(args) != 2 {
		return vm.NIL, fmt.Errorf("gogen/import-spec: expected 1 or 2 args, got %d", len(args))
	}
	path, err := asString(args[0])
	if err != nil {
		return vm.NIL, err
	}
	spec := &ast.ImportSpec{
		Path: &ast.BasicLit{Kind: token.STRING, Value: strconv.Quote(path)},
	}
	if len(args) == 2 {
		alias, err := asString(args[1])
		if err != nil {
			return vm.NIL, err
		}
		if !validIdent(alias) {
			return vm.NIL, fmt.Errorf("gogen/import-spec: %q is not a valid identifier", alias)
		}
		spec.Name = ast.NewIdent(alias)
	}
	return box(spec), nil
}

// file: (gogen/file "package-name" [imports] [decls]) -> *ast.File
func cFile(pkgV, importsV, declsV vm.Value) (vm.Value, error) {
	pkg, err := asString(pkgV)
	if err != nil {
		return vm.NIL, err
	}
	if !validIdent(pkg) {
		return vm.NIL, fmt.Errorf("gogen/file: %q is not a valid package name", pkg)
	}
	var decls []ast.Decl
	// Imports first as a GenDecl block (idiomatic Go).
	if importsV != vm.NIL {
		impVals, err := seqToValues(importsV)
		if err != nil {
			return vm.NIL, fmt.Errorf("gogen/file: imports: %w", err)
		}
		if len(impVals) > 0 {
			var specs []ast.Spec
			for i, iv := range impVals {
				n, err := unboxNode(iv)
				if err != nil {
					return vm.NIL, fmt.Errorf("gogen/file: import %d: %w", i, err)
				}
				is, ok := n.(*ast.ImportSpec)
				if !ok {
					return vm.NIL, fmt.Errorf("gogen/file: import %d: expected *ast.ImportSpec, got %T", i, n)
				}
				specs = append(specs, is)
			}
			decls = append(decls, &ast.GenDecl{Tok: token.IMPORT, Specs: specs})
		}
	}
	declVals, err := seqToValues(declsV)
	if err != nil {
		return vm.NIL, fmt.Errorf("gogen/file: decls: %w", err)
	}
	for i, dv := range declVals {
		d, err := unboxDecl(dv)
		if err != nil {
			return vm.NIL, fmt.Errorf("gogen/file: decl %d: %w", i, err)
		}
		decls = append(decls, d)
	}
	return box(&ast.File{
		Name:  ast.NewIdent(pkg),
		Decls: decls,
	}), nil
}

// --- render: the one and only output fn ------------------------------

var goFset = token.NewFileSet()

// Synthetic file used to mint token positions on demand. Many printer
// behaviors (multi-line const blocks, trailing line comments, multi-line
// composite literals) depend on adjacent nodes having distinct line
// positions, even if the actual offsets don't refer to real source.
//
// Line offsets are registered lazily in allocPos so processes that never
// touch gogen pay no startup cost. Each call to allocPos consumes one
// line.
var (
	goSynthFile = goFset.AddFile("gogen.synthetic.go", -1, 1<<28)
	goPosNext   = 1 // next byte offset to hand out (1-based)
)

// allocPos returns a fresh token.Pos on a new synthetic line. Used by
// constructors that need to position adjacent nodes (e.g. const specs in
// a block, fields in a struct, elements of a composite literal) on
// distinct lines so go/format.Node emits them on separate output lines.
//
// The synthetic file's line table is extended on demand: the first call
// registers offset 1 as the start of line 2, the next registers offset 2
// as the start of line 3, and so on. token.File.AddLine ignores
// duplicate or out-of-order offsets, so repeated calls past the high
// water mark remain safe.
func allocPos() token.Pos {
	p := token.Pos(goPosNext)
	goSynthFile.AddLine(goPosNext)
	goPosNext++
	return p
}

// stripPositions zeroes out every *Pos field reachable from an AST
// node. We do this on parser output before splicing it into hand-built
// nodes, so the parser's positions (which start at 1 and have no
// relationship to our synthetic FileSet) don't interfere with comment
// placement and line-break decisions.
func stripPositions(n ast.Node) {
	if n == nil {
		return
	}
	ast.Inspect(n, func(x ast.Node) bool {
		switch v := x.(type) {
		case *ast.Ident:
			v.NamePos = token.NoPos
		case *ast.BasicLit:
			v.ValuePos = token.NoPos
		case *ast.CompositeLit:
			v.Lbrace = token.NoPos
			v.Rbrace = token.NoPos
		case *ast.ArrayType:
			v.Lbrack = token.NoPos
		case *ast.StructType:
			v.Struct = token.NoPos
			if v.Fields != nil {
				v.Fields.Opening = token.NoPos
				v.Fields.Closing = token.NoPos
			}
		case *ast.InterfaceType:
			v.Interface = token.NoPos
			if v.Methods != nil {
				v.Methods.Opening = token.NoPos
				v.Methods.Closing = token.NoPos
			}
		case *ast.MapType:
			v.Map = token.NoPos
		case *ast.ChanType:
			v.Begin = token.NoPos
			v.Arrow = token.NoPos
		case *ast.FuncType:
			v.Func = token.NoPos
			if v.Params != nil {
				v.Params.Opening = token.NoPos
				v.Params.Closing = token.NoPos
			}
			if v.Results != nil {
				v.Results.Opening = token.NoPos
				v.Results.Closing = token.NoPos
			}
		case *ast.StarExpr:
			v.Star = token.NoPos
		case *ast.Ellipsis:
			v.Ellipsis = token.NoPos
		case *ast.ParenExpr:
			v.Lparen = token.NoPos
			v.Rparen = token.NoPos
		case *ast.SelectorExpr:
			// nothing — Sel and X positions handled by descent
		}
		return true
	})
}

func cRender(v vm.Value) (result vm.Value, retErr error) {
	defer func() {
		if r := recover(); r != nil {
			result = vm.NIL
			retErr = fmt.Errorf("gogen/render: panic during format: %v", r)
		}
	}()
	n, err := unboxNode(v)
	if err != nil {
		return vm.NIL, err
	}
	if n == nil {
		return vm.String(""), nil
	}
	var buf bytes.Buffer
	if err := format.Node(&buf, goFset, n); err != nil {
		return vm.NIL, fmt.Errorf("gogen/render: format: %w", err)
	}
	return vm.String(buf.String()), nil
}

// --- namespace install -----------------------------------------------

func init() { RegisterInstaller(installGogenNS) }

// nolint
func installGogenNS() {
	ns := DefNSBare("gogen")

	// Intentional shadows of clojure.core names — suppress warn-on-shadow.
	ns.Exclude("type")

	type entry struct {
		name string
		fn   vm.Value
	}
	mk := func(name string, fn vm.Value, err error) entry {
		if err != nil {
			panic(fmt.Sprintf("gogen: install %s: %v", name, err))
		}
		return entry{name, fn}
	}

	render, err := wrap1("render", cRender)
	entries := []entry{
		mk("render", render, err),

		mk(wrap1Named("ident", cIdent)),
		mk(wrap1Named("ident?", cIdentP)),
		mk(wrap1Named("ident-name", cIdentName)),
		mk(wrap1Named("type", cType)),
		mk(wrap1Named("int-lit", cIntLit)),
		mk(wrap1Named("float-lit", cFloatLit)),
		mk(wrap1Named("string-lit", cStringLit)),
		mk(wrap1Named("char-lit", cCharLit)),
		mk(wrap1Named("expr-stmt", cExprStmt)),
		mk(wrap1Named("block-stmt", cBlockStmt)),
		mk(wrap1Named("return-stmt", cReturn)),
		mk(wrap1Named("goto-stmt", cGotoStmt)),
		mk(wrap0Named("continue-stmt", cContinueStmt)),
		mk(wrap0Named("break-stmt", cBreakStmt)),

		mk(wrap2Named("unary", cUnary)),
		mk(wrap2Named("index", cIndex)),
		mk(wrap2Named("field-sel", cFieldSel)),
		mk(wrap2Named("call", cCall)),
		mk(wrap2Named("cast", cCast)),
		mk(wrap2Named("param", cParam)),
		mk(wrap2Named("variadic-param", cVariadicParam)),
		mk(wrap1Named("result", cResult)),
		mk(wrap2Named("type-assert", cTypeAssert)),
		mk(wrap2Named("kv-expr", cKVExpr)),
		mk(wrap2Named("composite-lit", cCompositeLit)),
		mk(wrap2Named("composite-lit-multi", cCompositeLitMulti)),
		mk(wrap2Named("case-clause", cCaseClause)),
		mk(wrap2Named("type-decl", cTypeDecl)),
		mk(wrap1Named("struct-type", cStructType)),
		mk(wrap1Named("interface-type", cInterfaceType)),
		mk(wrap1Named("const-block", cConstBlock)),
		mk(wrap2Named("with-doc", cWithDoc)),
		mk(wrap2Named("label-stmt", cLabelStmt)),

		mk(wrap3Named("binary", cBinary)),
		mk(wrap3Named("assign", cAssign)),
		mk(wrap3Named("multi-assign", cMultiAssign)),
		mk(wrap3Named("var-decl", cVarDecl)),
		mk(wrap3Named("file", cFile)),
		mk(wrap3Named("func-lit", cFuncLit)),
		mk(wrap3Named("switch-stmt", cSwitchStmt)),
		mk(wrap3Named("field-decl", cFieldDecl)),
		mk(wrap1Named("embed-field", cEmbedField)),
		mk(wrap3Named("top-var-decl", cTopVarDecl)),
		mk(wrap3Named("iface-method", cIfaceMethod)),

		mk(wrap4Named("if-stmt", cIfStmt)),
		mk(wrap4Named("for-stmt", cForStmt)),
		mk(wrap4Named("func-decl", cFuncDecl)),
		mk(wrap4Named("const-spec", cConstSpec)),

		mk(wrap5Named("method-decl", cMethodDecl)),

		mk("import-spec", makeVariadic("import-spec", cImportSpec), nil),
	}

	for _, e := range entries {
		ns.Def(e.name, e.fn)
	}

	MarkNSNeedsLoad("gogen")
}

// helpers to keep the entries table readable.
func wrap0Named(name string, fn func() (vm.Value, error)) (string, vm.Value, error) {
	v, err := wrap0(name, fn)
	return name, v, err
}
func wrap1Named(name string, fn func(vm.Value) (vm.Value, error)) (string, vm.Value, error) {
	v, err := wrap1(name, fn)
	return name, v, err
}
func wrap2Named(name string, fn func(vm.Value, vm.Value) (vm.Value, error)) (string, vm.Value, error) {
	v, err := wrap2(name, fn)
	return name, v, err
}
func wrap3Named(name string, fn func(vm.Value, vm.Value, vm.Value) (vm.Value, error)) (string, vm.Value, error) {
	v, err := wrap3(name, fn)
	return name, v, err
}
func wrap4Named(name string, fn func(vm.Value, vm.Value, vm.Value, vm.Value) (vm.Value, error)) (string, vm.Value, error) {
	v, err := wrap4(name, fn)
	return name, v, err
}
func wrap5Named(name string, fn func(vm.Value, vm.Value, vm.Value, vm.Value, vm.Value) (vm.Value, error)) (string, vm.Value, error) {
	v, err := wrap5(name, fn)
	return name, v, err
}

func makeVariadic(name string, fn func(...vm.Value) (vm.Value, error)) vm.Value {
	v, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		return fn(vs...)
	})
	return v
}
