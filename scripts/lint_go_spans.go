// Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
// SPDX-License-Identifier: MIT

package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

type span struct {
	path       string
	start, end int
}

func spans(paths []string) ([]span, error) {
	fset := token.NewFileSet()
	result := []span{}
	for _, path := range paths {
		if strings.ContainsAny(path, "\t\n\r") {
			return nil, fmt.Errorf("cannot parse Go file %q: path contains a tab or newline", path)
		}
		file, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			return nil, fmt.Errorf("cannot parse Go file %q: %w", path, err)
		}
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Body == nil {
				continue
			}
			result = append(result, span{
				path:  path,
				start: fset.Position(fn.Pos()).Line,
				end:   fset.Position(fn.Body.End()).Line,
			})
		}
	}
	return result, nil
}

func main() {
	paths := os.Args[1:]
	if len(paths) > 0 && paths[0] == "--" {
		paths = paths[1:]
	}
	rows, err := spans(paths)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	for _, row := range rows {
		fmt.Printf("%s\t%d\t%d\n", row.path, row.start, row.end)
	}
}
