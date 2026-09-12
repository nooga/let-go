package main

import (
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"strings"
)

// Callable extraction for code-quality measurement (see
// docs/superpowers/specs/2026-09-12-quality-bounded-score-design.md).
//
// This file is the AST walk only. It MEASURES and makes no judgments: every
// metric that consumes it -- mass, the CC > 10 split, erosion's ratio,
// weights, thresholds, reporting -- lives in .lg. Cyclomatic complexity is
// computed in Go because it is a pure function of AST node kinds, and this
// is a standalone command rather than a runtime native because a Go parser
// bridged into pkg/rt would be runtime scope creep for a tooling need.

// callable is one measured unit before it is boxed for .lg.
type callable struct {
	name string
	// "func" or "method" -- a callable in the enumeration -- or
	// "dead-closure", which is NOT a callable: it is a named nested binding
	// nothing references, reported only so the dead-code term can charge its
	// lines and the issues list can point at it.
	kind string
	line int
	end  int
	sloc int
	cc   int
}

// codeLines marks, for each 1-based line of src, whether the line carries
// code: not blank, and not occupied solely by a comment.
func codeLines(fset *token.FileSet, file *ast.File, src string) []bool {
	lines := strings.Split(src, "\n")
	code := make([]bool, len(lines)+2)
	for i, l := range lines {
		t := strings.TrimSpace(l)
		code[i+1] = t != ""
	}
	for _, cg := range file.Comments {
		start := fset.Position(cg.Pos())
		end := fset.Position(cg.End())
		first := start.Line
		// A comment trailing real code leaves that line as code; only the
		// lines it occupies alone are cleared.
		if start.Column > 1 && strings.TrimSpace(lines[start.Line-1][:start.Column-1]) != "" {
			first++
		}
		for l := first; l <= end.Line && l < len(code); l++ {
			code[l] = false
		}
	}
	return code
}

// fileCounts is the per-file line census: total lines, code lines (neither
// blank nor comment-only), and blank lines. The remainder -- comment-only
// lines -- is total minus code minus blank.
type fileCounts struct {
	lines int
	code  int
	blank int
}

func countFile(code []bool, src string) fileCounts {
	lines := strings.Split(src, "\n")
	// A trailing newline leaves a final empty element that is not a line.
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	out := fileCounts{lines: len(lines)}
	for i, l := range lines {
		if strings.TrimSpace(l) == "" {
			out.blank++
			continue
		}
		if i+1 < len(code) && code[i+1] {
			out.code++
		}
	}
	return out
}

func countCodeLines(code []bool, from, to int) int {
	n := 0
	for l := from; l <= to && l < len(code); l++ {
		if code[l] {
			n++
		}
	}
	return n
}

// decisionPoints counts cyclomatic decision points in `n`, skipping any
// subtree in `skip` (the named closures, which are separate callables).
//
// One point each for: *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt, an
// *ast.CaseClause with a non-empty List (a bare `default:` adds nothing), an
// *ast.CommClause with a non-nil Comm, and an *ast.BinaryExpr whose Op is
// token.LAND or token.LOR.
func decisionPoints(n ast.Node, skip map[ast.Node]bool) int {
	count := 0
	ast.Inspect(n, func(x ast.Node) bool {
		if x == nil || (skip != nil && skip[x]) {
			return false
		}
		switch t := x.(type) {
		case *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt:
			count++
		case *ast.CaseClause:
			if len(t.List) > 0 {
				count++
			}
		case *ast.CommClause:
			if t.Comm != nil {
				count++
			}
		case *ast.BinaryExpr:
			if t.Op == token.LAND || t.Op == token.LOR {
				count++
			}
		}
		return true
	})
	return count
}

type namedLit struct {
	name  string
	ident ast.Node // the binding occurrence, which is not a "use"
}

// namedLiterals finds the function literals BOUND TO AN IDENTIFIER inside
// `decl` -- `f := func(){}` and `var f = func(){}`.
//
// Binding form does not decide complexity attribution -- live nested
// functions fold into the enclosing declaration however they are written.
// This exists for one purpose: only a binding WITH A NAME can be checked for
// references, and an unreferenced one is dead code.
func namedLiterals(decl *ast.FuncDecl) map[*ast.FuncLit]namedLit {
	out := map[*ast.FuncLit]namedLit{}
	ast.Inspect(decl, func(x ast.Node) bool {
		switch t := x.(type) {
		case *ast.AssignStmt:
			if len(t.Lhs) == len(t.Rhs) {
				for i, r := range t.Rhs {
					lit, ok := r.(*ast.FuncLit)
					if !ok {
						continue
					}
					if id, ok := t.Lhs[i].(*ast.Ident); ok {
						out[lit] = namedLit{name: id.Name, ident: id}
					}
				}
			}
		case *ast.ValueSpec:
			for i, v := range t.Values {
				lit, ok := v.(*ast.FuncLit)
				if !ok || i >= len(t.Names) {
					continue
				}
				out[lit] = namedLit{name: t.Names[i].Name, ident: t.Names[i]}
			}
		}
		return true
	})
	return out
}

// identUses counts how many times `name` appears as an identifier anywhere in
// `decl`, excluding the binding occurrence itself -- both the bound IDENT on
// the left of the assignment and the literal on its right, neither of which
// is a use of the binding.
func identUses(decl *ast.FuncDecl, name string, lit ast.Node, ident ast.Node) int {
	n := 0
	ast.Inspect(decl, func(x ast.Node) bool {
		if x == lit || x == ident {
			return false
		}
		if id, ok := x.(*ast.Ident); ok && id.Name == name {
			n++
		}
		return true
	})
	return n
}

func receiverName(fset *token.FileSet, recv *ast.Field) string {
	var b strings.Builder
	// printer keeps `*Compiler` and `Compiler[T]` faithful without pulling in
	// the type checker.
	if err := printer.Fprint(&b, fset, recv.Type); err != nil {
		return "?"
	}
	return b.String()
}

// goCallables parses Go source and returns one entry per callable: every
// top-level function and method declaration, plus every NAMED function
// literal inside them.
func goCallables(src string) ([]callable, error) {
	cs, _, err := goAnalyze(src)
	return cs, err
}

// goAnalyze parses Go source once and returns everything measured from it:
// the callables and the file's line census. Keep new per-file measurements
// coming out of here rather than re-parsing.
func goAnalyze(src string) ([]callable, fileCounts, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "src.go", src, parser.ParseComments)
	if err != nil {
		return nil, fileCounts{}, err
	}
	code := codeLines(fset, file, src)
	counts := countFile(code, src)
	var out []callable

	for _, d := range file.Decls {
		decl, ok := d.(*ast.FuncDecl)
		if !ok || decl.Body == nil {
			// A declaration with no body (assembly or external stub) has no
			// complexity to measure.
			continue
		}
		name := decl.Name.Name
		kind := "func"
		if decl.Recv != nil && len(decl.Recv.List) > 0 {
			// Methods carry their receiver so two types' `emit` never collide.
			name = "(" + receiverName(fset, decl.Recv.List[0]) + ")." + name
			kind = "method"
		}
		start := fset.Position(decl.Pos()).Line
		end := fset.Position(decl.End()).Line

		// Every LIVE nested function -- argument position, immediately
		// invoked, deferred, `go`, or bound to a local name -- folds into
		// this declaration. Its branches are branches this declaration
		// executes; you cannot call it separately or reason about it
		// separately, so its complexity belongs here.
		//
		// The one exception is a named binding nothing references: it is
		// never executed, so it is dead code rather than complexity of the
		// enclosing function, and it is excluded from both this
		// declaration's CC and its SLOC.
		dead := map[ast.Node]bool{}
		deadLines := 0
		for lit, nl := range namedLiterals(decl) {
			if identUses(decl, nl.name, lit, nl.ident) > 0 {
				continue
			}
			dead[lit] = true
			ls := fset.Position(lit.Pos()).Line
			le := fset.Position(lit.End()).Line
			n := countCodeLines(code, ls, le)
			deadLines += n
			out = append(out, callable{
				name: name + "." + nl.name,
				kind: "dead-closure",
				line: ls, end: le,
				sloc: n,
				cc:   1 + decisionPoints(lit, nil),
			})
		}
		sloc := countCodeLines(code, start, end) - deadLines
		if sloc < 1 {
			sloc = 1
		}
		out = append(out, callable{
			name: name, kind: kind,
			line: start, end: end,
			sloc: sloc,
			cc:   1 + decisionPoints(decl, dead),
		})
	}
	return out, counts, nil
}
