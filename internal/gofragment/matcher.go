/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * SPDX-License-Identifier: MIT
 */

// Package gofragment is the Go AST fragment oracle used by the native-entry
// gate and the jank direct-ABI test.
//
// It compares a committed inline Go fragment against generated Go by parsing
// both sides with go/parser and comparing canonical AST renderings at a named
// point (a function declaration, or an expression inside one), requiring that
// point to occur exactly once. Matching on AST shape rather than text means
// formatting, comments and whitespace cannot make an assertion vacuous, and a
// lowering change that renames generated temporaries fails loudly with a diff.
package gofragment

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"reflect"
	"strings"
)

type GoFragmentKind string

const (
	GoExpression  GoFragmentKind = "expression"
	GoStatements  GoFragmentKind = "statements"
	GoFunction    GoFragmentKind = "function"
	GoDeclaration GoFragmentKind = "declaration"
)

// GoMatchRequest identifies one expected inline Go fragment and the named
// point in an independently parsed generated file where it must occur.
type GoMatchRequest struct {
	Kind     GoFragmentKind
	Expected string
	Target   string
}

// MatchGeneratedGoFragment parses both sides through go/parser, requires the
// named point to occur exactly once, then compares canonical AST renderings.
func MatchGeneratedGoFragment(req GoMatchRequest, generated string) error {
	want, err := parseExpectedGoFragment(req.Kind, req.Expected)
	if err != nil {
		return fmt.Errorf("parse expected %s: %w", req.Kind, err)
	}

	generatedSet := token.NewFileSet()
	file, err := parser.ParseFile(
		generatedSet,
		"generated.go",
		generated,
		parser.SkipObjectResolution,
	)
	if err != nil {
		return fmt.Errorf("parse generated Go: %w", err)
	}

	scopes := selectGeneratedGoScopes(req.Kind, req.Target, file)
	if len(scopes) != 1 {
		return fragmentCardinalityError(req, len(scopes))
	}

	wantText, err := normalizeGoAST(want)
	if err != nil {
		return fmt.Errorf("normalize expected %s: %w", req.Kind, err)
	}
	if req.Kind == GoExpression || req.Kind == GoStatements {
		matches, candidates, err := matchGoSubtrees(req.Kind, want, wantText, scopes[0])
		if err != nil {
			return err
		}
		if len(matches) != 1 {
			cardinality := fragmentCardinalityError(req, len(matches))
			if len(matches) == 0 && len(candidates) > 0 {
				gotText, err := normalizeGoAST(candidates[0])
				if err != nil {
					return err
				}
				return fmt.Errorf(
					"%w\n%s target %q AST mismatch:\n%s",
					cardinality,
					req.Kind,
					req.Target,
					goFragmentDiff(wantText, gotText),
				)
			}
			return cardinality
		}
		return nil
	}

	gotText, err := normalizeGoAST(scopes[0])
	if err != nil {
		return fmt.Errorf("normalize generated %s target %q: %w", req.Kind, req.Target, err)
	}
	if wantText == gotText {
		return nil
	}

	diff := goFragmentDiff(wantText, gotText)
	return fmt.Errorf("%s target %q AST mismatch:\n%s", req.Kind, req.Target, diff)
}

func fragmentCardinalityError(req GoMatchRequest, found int) error {
	return fmt.Errorf("%s target %q: found %d, want exactly 1", req.Kind, req.Target, found)
}

func matchGoSubtrees(
	kind GoFragmentKind,
	want ast.Node,
	wantText string,
	scope ast.Node,
) ([]ast.Node, []ast.Node, error) {
	var candidates []ast.Node
	switch kind {
	case GoExpression:
		wantExpr := want.(ast.Expr)
		ast.Inspect(scope, func(node ast.Node) bool {
			expr, ok := node.(ast.Expr)
			if ok && reflect.TypeOf(expr) == reflect.TypeOf(wantExpr) {
				candidates = append(candidates, expr)
			}
			return true
		})
	case GoStatements:
		wantBlock := want.(*ast.BlockStmt)
		width := len(wantBlock.List)
		if width == 0 {
			return nil, nil, fmt.Errorf("expected statements fragment is empty")
		}
		ast.Inspect(scope, func(node ast.Node) bool {
			block, ok := node.(*ast.BlockStmt)
			if !ok {
				return true
			}
			for start := 0; start+width <= len(block.List); start++ {
				candidates = append(candidates, &ast.BlockStmt{
					List: block.List[start : start+width],
				})
			}
			return true
		})
	default:
		return nil, nil, fmt.Errorf("%s does not support subtree matching", kind)
	}

	var matches []ast.Node
	for _, candidate := range candidates {
		candidateText, err := normalizeGoAST(candidate)
		if err != nil {
			return nil, nil, fmt.Errorf("normalize generated %s candidate: %w", kind, err)
		}
		if candidateText == wantText {
			matches = append(matches, candidate)
		}
	}
	return matches, candidates, nil
}

func parseExpectedGoFragment(kind GoFragmentKind, fragment string) (ast.Node, error) {
	fileSet := token.NewFileSet()
	switch kind {
	case GoExpression:
		return parser.ParseExprFrom(
			fileSet,
			"expected.go",
			fragment,
			parser.SkipObjectResolution,
		)
	case GoStatements:
		file, err := parser.ParseFile(
			fileSet,
			"expected.go",
			"package probe\nfunc _() {\n"+fragment+"\n}",
			parser.SkipObjectResolution,
		)
		if err != nil {
			return nil, err
		}
		return file.Decls[0].(*ast.FuncDecl).Body, nil
	case GoFunction:
		return parseExpectedGoDeclaration(fileSet, fragment, true)
	case GoDeclaration:
		return parseExpectedGoDeclaration(fileSet, fragment, false)
	default:
		return nil, fmt.Errorf("unknown Go fragment kind %q", kind)
	}
}

func parseExpectedGoDeclaration(fileSet *token.FileSet, fragment string, wantFunction bool) (ast.Node, error) {
	file, err := parser.ParseFile(
		fileSet,
		"expected.go",
		"package probe\n"+fragment,
		parser.SkipObjectResolution,
	)
	if err != nil {
		return nil, err
	}
	if len(file.Decls) != 1 {
		return nil, fmt.Errorf("found %d declarations, want exactly 1", len(file.Decls))
	}
	_, isFunction := file.Decls[0].(*ast.FuncDecl)
	if wantFunction != isFunction {
		want := "non-function declaration"
		if wantFunction {
			want = "function"
		}
		return nil, fmt.Errorf("expected one %s", want)
	}
	return file.Decls[0], nil
}

func selectGeneratedGoScopes(kind GoFragmentKind, target string, file *ast.File) []ast.Node {
	var matches []ast.Node
	for _, decl := range file.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok && fn.Name.Name == target {
			switch kind {
			case GoExpression, GoStatements, GoFunction:
				matches = append(matches, fn)
			}
		}

		gen, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range gen.Specs {
			switch spec := spec.(type) {
			case *ast.ValueSpec:
				for _, name := range spec.Names {
					if kind == GoExpression && name.Name == target {
						matches = append(matches, gen)
					}
				}
			case *ast.TypeSpec:
				if kind == GoDeclaration && spec.Name.Name == target {
					matches = append(matches, gen)
				}
			}
		}
	}
	return matches
}

func normalizeGoAST(node ast.Node) (string, error) {
	var out bytes.Buffer
	if err := format.Node(&out, token.NewFileSet(), node); err != nil {
		return "", err
	}
	return out.String(), nil
}

func goFragmentDiff(want, got string) string {
	var out strings.Builder
	out.WriteString("--- expected\n+++ generated\n")
	for _, line := range strings.Split(want, "\n") {
		fmt.Fprintf(&out, "-%s\n", line)
	}
	for _, line := range strings.Split(got, "\n") {
		fmt.Fprintf(&out, "+%s\n", line)
	}
	return out.String()
}
