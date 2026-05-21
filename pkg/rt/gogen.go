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

	"github.com/nooga/let-go/pkg/vm"
)

// --- goASTValue: boxed *ast.Node as a vm.Value -----------------------

type theGoASTType struct{}

func (t *theGoASTType) String() string     { return t.Name() }
func (t *theGoASTType) Type() vm.ValueType { return vm.TypeType }
func (t *theGoASTType) Unbox() interface{} { return reflect.TypeOf(t) }
func (t *theGoASTType) Name() string       { return "let-go.lang.GoAST" }
func (t *theGoASTType) Box(_ interface{}) (vm.Value, error) {
	return vm.NIL, fmt.Errorf("gogen: GoAST values are constructed via gogen/* fns, not boxed")
}

var GoASTType *theGoASTType = &theGoASTType{}

type goASTValue struct{ node ast.Node }

func (g *goASTValue) String() string     { return fmt.Sprintf("#<go-ast %T>", g.node) }
func (g *goASTValue) Type() vm.ValueType { return GoASTType }
func (g *goASTValue) Unbox() interface{} { return g.node }

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

// type-expr: (gogen/type "spec") -> parsed type expression
// Uses go/parser so the full Go type grammar is supported.
func cType(v vm.Value) (vm.Value, error) {
	s, err := asString(v)
	if err != nil {
		return vm.NIL, err
	}
	expr, err := parser.ParseExpr(s)
	if err != nil {
		return vm.NIL, fmt.Errorf("gogen: parsing type %q: %w", s, err)
	}
	return box(expr), nil
}

// int-lit / float-lit / string-lit
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

// nolint
func installGogenNS() {
	ns := DefNSBare("gogen")

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
		mk(wrap1Named("type", cType)),
		mk(wrap1Named("int-lit", cIntLit)),
		mk(wrap1Named("float-lit", cFloatLit)),
		mk(wrap1Named("string-lit", cStringLit)),
		mk(wrap1Named("expr-stmt", cExprStmt)),
		mk(wrap1Named("return-stmt", cReturn)),

		mk(wrap2Named("unary", cUnary)),
		mk(wrap2Named("index", cIndex)),
		mk(wrap2Named("field-sel", cFieldSel)),
		mk(wrap2Named("call", cCall)),
		mk(wrap2Named("cast", cCast)),
		mk(wrap2Named("param", cParam)),
		mk(wrap1Named("result", cResult)),
		mk(wrap2Named("type-assert", cTypeAssert)),
		mk(wrap2Named("kv-expr", cKVExpr)),
		mk(wrap2Named("composite-lit", cCompositeLit)),
		mk(wrap2Named("case-clause", cCaseClause)),

		mk(wrap3Named("binary", cBinary)),
		mk(wrap3Named("assign", cAssign)),
		mk(wrap3Named("multi-assign", cMultiAssign)),
		mk(wrap3Named("var-decl", cVarDecl)),
		mk(wrap3Named("file", cFile)),
		mk(wrap3Named("func-lit", cFuncLit)),
		mk(wrap3Named("switch-stmt", cSwitchStmt)),

		mk(wrap4Named("if-stmt", cIfStmt)),
		mk(wrap4Named("for-stmt", cForStmt)),
		mk(wrap4Named("func-decl", cFuncDecl)),

		mk("import-spec", makeVariadic("import-spec", cImportSpec), nil),
	}

	for _, e := range entries {
		ns.Def(e.name, e.fn)
	}

	MarkNSNeedsLoad("gogen")
}

// helpers to keep the entries table readable.
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

func makeVariadic(name string, fn func(...vm.Value) (vm.Value, error)) vm.Value {
	v, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		return fn(vs...)
	})
	return v
}
