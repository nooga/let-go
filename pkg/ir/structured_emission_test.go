/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package ir_test

// Tests for the structured Go emission path (ir.structurize + lower_go
// tree-walk). The premise: build.lg only emits reducible CFGs, so the
// lowerer can walk the structured control tree and emit if/for/break/
// continue instead of goto+labels — and the emitted Go must be VALID
// (the structurizer must not drop the continuation after an if-join,
// which would produce an empty if + missing return + unused locals).
//
// RED-first: on the goto-based emitter these assertions fail because the
// output still contains `goto`. After the structured emitter is ported,
// the complex-CFG cases must additionally type-check (no dropped
// continuation).

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

// structuredCorpus mirrors the structurize-pass corpus but is chosen to
// exercise the lower_go tree-walk, including continuation-bearing shapes
// (a value-if whose join binds a local that later statements consume).
var structuredCorpus = []struct {
	name string
	src  string
}{
	{"straight", `(defn straight [x] (+ x 1))`},
	{"one-if", `(defn one-if [x] (if (< x 0) 1 2))`},
	// classify: outer value-if binds y, then a SECOND if consumes y.
	// If the structurizer drops the continuation after the first if,
	// the binding of y is emitted but never returned -> missing return
	// + y declared-and-not-used.
	{"classify", `(defn classify [x]
                     (let [y (if (< x 0) (- x) x)]
                       (if (> y 100) :big (if (> y 10) :med :small))))`},
	{"sum", `(defn sum [n] (loop [i 0 s 0] (if (< i n) (recur (inc i) (+ s i)) s)))`},
	// loop-if: an if INSIDE a loop body whose arms both continue, with a
	// post-loop continuation (return t). A prime continuation-drop shape.
	{"loop-if", `(defn loop-if [n]
                    (loop [i 0 t 0]
                      (if (< i n) (recur (inc i) (if (even? i) (+ t i) t)) t)))`},
}

// renderStructuredFile builds, optimizes and lowers `src` in :strict mode,
// then renders the full Go file to a string.
func renderStructuredFile(t *testing.T, name, src string) string {
	t.Helper()
	fn := buildLispIR(t, src)
	optimizeLispIR(t, fn)
	result := lowerGo(t, fn, ":strict")
	if got := result.ValueAt(vm.Keyword("status")); got != vm.Keyword("lowered") {
		t.Fatalf("%s: expected :lowered status, got %v", name, got)
	}
	passVarCounter++
	varName := fmt.Sprintf("*structured-file-%d*", passVarCounter)
	rt.NS(rt.NameCoreNS).Def(varName, result)
	file := runLispExpr(t, fmt.Sprintf(`(ir.lower-go/file "p" [%s])`, varName))
	return bindAndRenderGoFile(t, file)
}

func TestStructuredEmissionIsGotoFree(t *testing.T) {
	ensureLoader()
	for _, c := range structuredCorpus {
		t.Run(c.name, func(t *testing.T) {
			rendered := renderStructuredFile(t, c.name, c.src)
			if strings.Contains(rendered, "goto ") || strings.Contains(rendered, "_blk:") {
				t.Fatalf("%s: expected structured (goto-free) Go, found goto/label\n--- go ---\n%s", c.name, rendered)
			}
		})
	}
}

func TestStructuredEmissionTypeChecks(t *testing.T) {
	ensureLoader()
	for _, c := range structuredCorpus {
		t.Run(c.name, func(t *testing.T) {
			rendered := renderStructuredFile(t, c.name, c.src)
			if err := typeCheckGoSource(rendered); err != nil {
				t.Fatalf("%s: rendered Go does not type-check (dropped continuation?): %v\n--- go ---\n%s", c.name, err, rendered)
			}
		})
	}
}

// typeCheckGoSource parses and type-checks a complete Go source file using
// the source importer, so "missing return" and "declared and not used"
// (the continuation-drop signature) surface as errors.
func typeCheckGoSource(src string) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "rendered.go", src, parser.AllErrors)
	if err != nil {
		return fmt.Errorf("parse: %w", err)
	}
	conf := types.Config{Importer: importer.ForCompiler(fset, "source", nil)}
	var firstErr error
	conf.Error = func(e error) {
		if firstErr == nil {
			firstErr = e
		}
	}
	_, _ = conf.Check("p", fset, []*ast.File{f}, nil)
	return firstErr
}
