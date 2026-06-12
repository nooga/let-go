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
	"go/printer"
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

// goShapes is a histogram of the control-flow node kinds present in a
// rendered lowered-Go function, recovered by traversing the parsed go/ast
// rather than substring-matching the rendered text. Substring checks for
// "goto " / "_blk:" are brittle (false-positives inside string literals or
// comments, and couple to the label-naming convention); an AST traversal
// asserts the actual emitted structure.
type goShapes struct {
	forLoops, ifs, continues, breaks, gotos, labels int
}

// shapesOf parses rendered lowered Go and counts the control-flow node kinds
// reachable from the file. Reuses parseLoweredGo (lisp_lower_go_ast_test.go).
func shapesOf(t *testing.T, rendered string) goShapes {
	t.Helper()
	f := parseLoweredGo(t, rendered)
	var s goShapes
	ast.Inspect(f, func(n ast.Node) bool {
		switch b := n.(type) {
		case *ast.ForStmt, *ast.RangeStmt:
			s.forLoops++
		case *ast.IfStmt:
			s.ifs++
		case *ast.LabeledStmt:
			s.labels++
		case *ast.BranchStmt:
			switch b.Tok {
			case token.CONTINUE:
				s.continues++
			case token.BREAK:
				s.breaks++
			case token.GOTO:
				s.gotos++
			}
		}
		return true
	})
	return s
}

// TestStructuredEmissionIsGotoFree asserts, at the AST level, that the
// structured emitter produced zero `goto` statements and zero labeled
// statements for every corpus case.
func TestStructuredEmissionIsGotoFree(t *testing.T) {
	ensureLoader()
	for _, c := range structuredCorpus {
		t.Run(c.name, func(t *testing.T) {
			rendered := renderStructuredFile(t, c.name, c.src)
			s := shapesOf(t, rendered)
			if s.gotos != 0 || s.labels != 0 {
				t.Fatalf("%s: expected goto-free structured Go, found %d goto / %d labels\n--- go ---\n%s",
					c.name, s.gotos, s.labels, rendered)
			}
		})
	}
}

// TestStructuredEmissionHasExpectedStructure verifies the emitter produced the
// CONTROL STRUCTURES the corpus is designed to exercise — not merely that no
// goto remains. A goto-free-but-degenerate emission (e.g. structurize silently
// falling back to an empty body, or emitting an `if` where a `for` was
// required) is goto-free yet wrong; this gate catches that by asserting the
// presence of for/if/continue per case. Expectations are the observed shapes
// of the structured emitter (see the probe in the PR review).
func TestStructuredEmissionHasExpectedStructure(t *testing.T) {
	ensureLoader()
	// want* are MINIMUMS / required presence, not exact counts, so the test
	// is robust to incidental emission changes while still pinning the
	// structural intent of each case.
	cases := map[string]struct {
		wantLoop     bool // recur/loop must lower to a Go for-loop
		minIfs       int  // conditionals must lower to if statements
		wantContinue bool // a recur back-edge must emit `continue`
		wantNoLoop   bool // straight-line code must NOT introduce a loop
	}{
		"straight": {wantNoLoop: true},
		"one-if":   {minIfs: 1, wantNoLoop: true},
		"classify": {minIfs: 2, wantNoLoop: true}, // nested if; second if consumes the join binding
		"sum":      {wantLoop: true, minIfs: 1, wantContinue: true},
		"loop-if":  {wantLoop: true, minIfs: 2, wantContinue: true},
	}
	for _, c := range structuredCorpus {
		want, ok := cases[c.name]
		if !ok {
			t.Fatalf("no structural expectation for corpus case %q", c.name)
		}
		t.Run(c.name, func(t *testing.T) {
			s := shapesOf(t, renderStructuredFile(t, c.name, c.src))
			if want.wantLoop && s.forLoops < 1 {
				t.Errorf("%s: expected a for-loop (recur→for), got %d", c.name, s.forLoops)
			}
			if want.wantNoLoop && s.forLoops != 0 {
				t.Errorf("%s: expected no loop in straight-line/conditional code, got %d for-loops", c.name, s.forLoops)
			}
			if s.ifs < want.minIfs {
				t.Errorf("%s: expected >=%d if statements, got %d", c.name, want.minIfs, s.ifs)
			}
			if want.wantContinue && s.continues < 1 {
				t.Errorf("%s: expected a continue (recur back-edge), got %d", c.name, s.continues)
			}
		})
	}
}

// stmtStr renders a statement back to source for duplicate detection.
func stmtStr(s ast.Stmt) string {
	var sb strings.Builder
	_ = printer.Fprint(&sb, token.NewFileSet(), s)
	return sb.String()
}

// TestStructuredEmissionNoDuplicateRecurCopies guards against the back-edge
// parallel-copy being emitted twice. A recur lowers to `i = <tmp>; s = <tmp>`
// before `continue`; if both the enclosing `:if`/`:seq` node AND the child
// `:continue` node emit the edge copies, the assignment block is duplicated
// (i = v15; s = v16; i = v15; s = v16). It is idempotent here so it type-checks
// and is goto-free — invisible to the other gates — but it is still wrong
// emission. Assert no statement block contains the same assignment twice.
func TestStructuredEmissionNoDuplicateRecurCopies(t *testing.T) {
	ensureLoader()
	for _, c := range structuredCorpus {
		t.Run(c.name, func(t *testing.T) {
			f := parseLoweredGo(t, renderStructuredFile(t, c.name, c.src))
			ast.Inspect(f, func(n ast.Node) bool {
				blk, ok := n.(*ast.BlockStmt)
				if !ok {
					return true
				}
				seen := map[string]bool{}
				for _, s := range blk.List {
					a, ok := s.(*ast.AssignStmt)
					if !ok {
						continue
					}
					key := stmtStr(a)
					if seen[key] {
						t.Errorf("%s: duplicate assignment %q in a single block (recur copy emitted twice)", c.name, key)
					}
					seen[key] = true
				}
				return true
			})
		})
	}
}

// isTempIdent reports whether name is a generated SSA temp local (vNN / _pcNN).
func isTempIdent(name string) bool {
	for _, pre := range []string{"v", "_pc"} {
		if rest := strings.TrimPrefix(name, pre); rest != name && rest != "" {
			allDigits := true
			for _, r := range rest {
				if r < '0' || r > '9' {
					allDigits = false
					break
				}
			}
			if allDigits {
				return true
			}
		}
	}
	return false
}

// forSelfUpdate reports whether a for-loop body contains a direct
// self-referential update `x = x <op> ...` for the named variable.
func forSelfUpdate(fr *ast.ForStmt, name string) bool {
	found := false
	ast.Inspect(fr.Body, func(m ast.Node) bool {
		a, ok := m.(*ast.AssignStmt)
		if !ok || len(a.Lhs) != 1 || len(a.Rhs) != 1 {
			return true
		}
		lhs, lok := a.Lhs[0].(*ast.Ident)
		if !lok || lhs.Name != name {
			return true
		}
		bin, bok := a.Rhs[0].(*ast.BinaryExpr)
		if !bok {
			return true
		}
		// one operand is the variable itself: x = x + 1 (self-update)
		for _, operand := range []ast.Expr{bin.X, bin.Y} {
			if id, ok := operand.(*ast.Ident); ok && id.Name == name {
				found = true
			}
		}
		return true
	})
	return found
}

// forStmtOf returns the (single) for-loop in a parsed lowered function.
func forStmtOf(t *testing.T, f *ast.File) *ast.ForStmt {
	t.Helper()
	var fr *ast.ForStmt
	ast.Inspect(f, func(n ast.Node) bool {
		if x, ok := n.(*ast.ForStmt); ok && fr == nil {
			fr = x
		}
		return true
	})
	if fr == nil {
		t.Fatal("no for-loop in lowered output")
	}
	return fr
}

// loopVarTempAssign reports the first loop-carried var routed through a temp
// (`x = vNN`), or "" if none.
func loopVarTempAssign(fr *ast.ForStmt) string {
	bad := ""
	ast.Inspect(fr.Body, func(m ast.Node) bool {
		a, ok := m.(*ast.AssignStmt)
		if !ok || len(a.Lhs) != 1 || len(a.Rhs) != 1 || bad != "" {
			return true
		}
		lhs, lok := a.Lhs[0].(*ast.Ident)
		rhs, rok := a.Rhs[0].(*ast.Ident)
		if lok && rok && !isTempIdent(lhs.Name) && isTempIdent(rhs.Name) {
			bad = stmtStr(a)
		}
		return true
	})
	return bad
}

// TestStructuredEmissionInlinesRecurUpdates asserts pure-arithmetic recur
// updates are inlined directly (i = i + 1, s = s + i) rather than routed through
// single-use temps. `sum` is all pure arithmetic so it fully collapses; the
// accumulator `s = s + i` must be emitted BEFORE `i = i + 1` so it reads the OLD
// i — the parallel-rebind invariant (cycle-safe sequencing). `loop-if`'s `t` is
// a conditional (phi), not pure-inline, so it legitimately keeps its temp.
func TestStructuredEmissionInlinesRecurUpdates(t *testing.T) {
	ensureLoader()

	t.Run("sum-full-collapse", func(t *testing.T) {
		var src string
		for _, c := range structuredCorpus {
			if c.name == "sum" {
				src = c.src
			}
		}
		fr := forStmtOf(t, parseLoweredGo(t, renderStructuredFile(t, "sum", src)))
		if bad := loopVarTempAssign(fr); bad != "" {
			t.Errorf("sum: loop var still routed through a temp: %q (expected i = i + 1 / s = s + i inlined)", bad)
		}
		if !forSelfUpdate(fr, "i") {
			t.Errorf("sum: expected inlined self-update i = i + 1")
		}
		// parallel rebind: the accumulator s must read the OLD i, so `s = s + i`
		// must be sequenced before `i = i + 1`.
		body := stmtStr(fr.Body)
		si := strings.Index(body, "s = s + i")
		ii := strings.Index(body, "i = i + 1")
		if si < 0 || ii < 0 {
			t.Errorf("sum: expected both `s = s + i` and `i = i + 1` inlined; got:\n%s", body)
		} else if si > ii {
			t.Errorf("sum: `s = s + i` must precede `i = i + 1` (parallel rebind: s reads OLD i); got:\n%s", body)
		}
	})

	t.Run("loop-if-induction", func(t *testing.T) {
		var src string
		for _, c := range structuredCorpus {
			if c.name == "loop-if" {
				src = c.src
			}
		}
		fr := forStmtOf(t, parseLoweredGo(t, renderStructuredFile(t, "loop-if", src)))
		if !forSelfUpdate(fr, "i") {
			t.Errorf("loop-if: expected inlined self-update i = i + 1")
		}
	})
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
