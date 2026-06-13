/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package ir_test

// Structural assertions over rendered lowered-Go source. Parsing the output
// into a go/ast and matching node shapes is robust to formatting/whitespace
// and expresses intent far more precisely than regex over rendered text.

import (
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"strings"
	"testing"
)

// parseLoweredGo parses a rendered lowered package into an *ast.File, failing
// the test (with the offending source) if it is not syntactically valid Go.
func parseLoweredGo(t *testing.T, src string) *ast.File {
	t.Helper()
	f, err := parser.ParseFile(token.NewFileSet(), "lowered.go", src, parser.AllErrors)
	if err != nil {
		t.Fatalf("rendered lowered Go did not parse: %v\n----\n%s", err, src)
	}
	return f
}

// exprStr renders an ast.Expr back to source for readable assertion messages.
func exprStr(e ast.Expr) string {
	if e == nil {
		return "<nil>"
	}
	var sb strings.Builder
	_ = printer.Fprint(&sb, token.NewFileSet(), e)
	return sb.String()
}

// findPkgVar returns the first package-level `var` spec whose name satisfies
// pred, along with the rendered source of its declared type.
func findPkgVar(f *ast.File, pred func(name string) bool) (name, typ string, ok bool) {
	for _, d := range f.Decls {
		gd, isGen := d.(*ast.GenDecl)
		if !isGen || gd.Tok != token.VAR {
			continue
		}
		for _, spec := range gd.Specs {
			vs, isVal := spec.(*ast.ValueSpec)
			if !isVal {
				continue
			}
			for _, n := range vs.Names {
				if pred(n.Name) {
					return n.Name, exprStr(vs.Type), true
				}
			}
		}
	}
	return "", "", false
}

// findFunc returns the first top-level (non-method) func decl whose name
// satisfies pred.
func findFunc(f *ast.File, pred func(name string) bool) (*ast.FuncDecl, bool) {
	for _, d := range f.Decls {
		if fd, ok := d.(*ast.FuncDecl); ok && fd.Recv == nil && pred(fd.Name.Name) {
			return fd, true
		}
	}
	return nil, false
}

// funcResultTypes returns the rendered result types of a func decl.
func funcResultTypes(fd *ast.FuncDecl) []string {
	var out []string
	if fd.Type.Results == nil {
		return out
	}
	for _, r := range fd.Type.Results.List {
		out = append(out, exprStr(r.Type))
	}
	return out
}

// isSelector reports whether e is `<x>.<sel>` where x renders to xStr.
func isSelector(e ast.Expr, xStr, sel string) bool {
	se, ok := e.(*ast.SelectorExpr)
	return ok && se.Sel.Name == sel && exprStr(se.X) == xStr
}

// callsLookupVarNamed reports whether node contains a call
// `rt.LookupVar(<anything>, "<name>")`.
func callsLookupVarNamed(node ast.Node, name string) bool {
	found := false
	ast.Inspect(node, func(n ast.Node) bool {
		ce, ok := n.(*ast.CallExpr)
		if !ok || len(ce.Args) != 2 {
			return true
		}
		if !isSelector(ce.Fun, "rt", "LookupVar") {
			return true
		}
		if lit, ok := ce.Args[1].(*ast.BasicLit); ok &&
			lit.Kind == token.STRING && strings.Trim(lit.Value, `"`) == name {
			found = true
		}
		return true
	})
	return found
}

// assertsType reports whether node contains a type assertion `x.(<typ>)`
// whose target type renders to typ (e.g. "vm.IDeref").
func assertsType(node ast.Node, typ string) bool {
	found := false
	ast.Inspect(node, func(n ast.Node) bool {
		if ta, ok := n.(*ast.TypeAssertExpr); ok && ta.Type != nil && exprStr(ta.Type) == typ {
			found = true
		}
		return true
	})
	return found
}

// callsInvokeValue reports whether node contains a trampoline call
// `rt.InvokeValue(...)` or its ExecContext-threaded form `rt.InvokeValueEC(...)`
// (lowered code emits the EC variant so dynamic vars resolve against the
// running context).
func callsInvokeValue(node ast.Node) bool {
	found := false
	ast.Inspect(node, func(n ast.Node) bool {
		if ce, ok := n.(*ast.CallExpr); ok &&
			(isSelector(ce.Fun, "rt", "InvokeValue") || isSelector(ce.Fun, "rt", "InvokeValueEC")) {
			found = true
		}
		return true
	})
	return found
}

// initFuncs returns all `func init()` decls in the file.
func initFuncs(f *ast.File) []*ast.FuncDecl {
	var out []*ast.FuncDecl
	for _, d := range f.Decls {
		if fd, ok := d.(*ast.FuncDecl); ok && fd.Recv == nil && fd.Name.Name == "init" {
			out = append(out, fd)
		}
	}
	return out
}

// ifnDispatch describes a matched cached-var IFn dispatch call chain
//
//	rt.CachedVarFn(&<varName>, "ns", "name").Invoke(<args>)
type ifnDispatch struct {
	varName string // cached var base, e.g. "__v_clojure_core_count"
	nsArg   string // ns string literal passed to CachedVarFn
	nameArg string // name string literal passed to CachedVarFn
	nargs   int    // number of invoke args
}

func basicStr(e ast.Expr) string {
	if lit, ok := e.(*ast.BasicLit); ok && lit.Kind == token.STRING {
		return strings.Trim(lit.Value, `"`)
	}
	return ""
}

// findIFnDispatch searches node for the cached-var IFn dispatch shape. Lowered
// code threads the caller's ExecContext, so the call is
//
//	ec.Invoke(rt.CachedVarFn(&__v_*, "ns", "name"), ARGS...)
//
// (the older method form rt.CachedVarFn(...).Invoke(ARGS...) is still
// recognised). Returns the first match.
func findIFnDispatch(node ast.Node) (ifnDispatch, bool) {
	var res ifnDispatch
	found := false
	ast.Inspect(node, func(n ast.Node) bool {
		if found {
			return false
		}
		// Outermost: <recv>.Invoke(...)
		invoke, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		invSel, ok := invoke.Fun.(*ast.SelectorExpr)
		if !ok || invSel.Sel.Name != "Invoke" {
			return true
		}
		// Locate the rt.CachedVarFn(&__v_*, ns, name) call. Two shapes:
		//   ec.Invoke(rt.CachedVarFn(...), ARGS...)  -> first call arg
		//   rt.CachedVarFn(...).Invoke(ARGS...)      -> the receiver
		var cached *ast.CallExpr
		invokeArgs := invoke.Args
		if c, ok := invSel.X.(*ast.CallExpr); ok && isSelector(c.Fun, "rt", "CachedVarFn") {
			cached = c
		} else if len(invoke.Args) >= 1 {
			if c, ok := invoke.Args[0].(*ast.CallExpr); ok && isSelector(c.Fun, "rt", "CachedVarFn") {
				cached = c
				invokeArgs = invoke.Args[1:]
			}
		}
		if cached == nil || len(cached.Args) != 3 {
			return true
		}
		amp, ok := cached.Args[0].(*ast.UnaryExpr)
		if !ok || amp.Op != token.AND {
			return true
		}
		varIdent, ok := amp.X.(*ast.Ident)
		if !ok || !strings.HasPrefix(varIdent.Name, "__v_") {
			return true
		}
		res = ifnDispatch{
			varName: varIdent.Name,
			nsArg:   basicStr(cached.Args[1]),
			nameArg: basicStr(cached.Args[2]),
			nargs:   len(invokeArgs),
		}
		found = true
		return false
	})
	return res, found
}
