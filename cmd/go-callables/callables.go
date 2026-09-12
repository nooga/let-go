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
	// loopDepth is the maximum nesting of `for`/`range` inside the
	// declaration, and selfCall reports whether it calls itself. Both feed
	// the .lg cost model, which estimates algorithmic complexity rather than
	// branch count. The recursion SHRINK shape is not read for Go, so a
	// self-recursive Go function is reported as unknown rather than guessed.
	loopDepth     int
	untracedLoops int
	selfCall      bool
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

// paramDerived collects the identifiers of `decl` that carry input size: its
// parameters, and locals assigned from something parameter-derived.
//
// THE POINT OF THE RULE. A degree is a function of INPUT SIZE, and a
// function's inputs are its parameters. A loop over a package-level table, a
// literal, or `len(vs)` where vs is a variadic call's own arguments is
// CONSTANT in the input however deeply it nests. installLangNS -- four nested
// loops over fixed registration tables -- read as n^4 for exactly this
// reason.
func paramDerived(decl *ast.FuncDecl) map[string]string {
	// name -> the ORIGIN parameter it derives from. Attributing to the origin
	// is what makes nesting compose: `for row := range xs { for x := range row }`
	// is xs^2, because row is part of xs -- not xs^1 * row^1.
	derived := map[string]string{}
	if decl.Type.Params != nil {
		for _, f := range decl.Type.Params.List {
			for _, n := range f.Names {
				derived[n.Name] = n.Name
			}
		}
	}
	if decl.Recv != nil {
		for _, f := range decl.Recv.List {
			for _, n := range f.Names {
				derived[n.Name] = n.Name
			}
		}
	}
	// One forward pass: a local assigned from a derived expression becomes
	// derived itself, which follows chains like `ys := xs[1:]`.
	ast.Inspect(decl, func(n ast.Node) bool {
		switch t := n.(type) {
		case *ast.AssignStmt:
			if len(t.Lhs) == len(t.Rhs) {
				for i, r := range t.Rhs {
					if id, ok := t.Lhs[i].(*ast.Ident); ok {
						if root := derivedRoot(r, derived); root != "" {
							derived[id.Name] = root
						}
					}
				}
			}
		case *ast.RangeStmt:
			// An element of a derived collection carries its size, and its
			// ORIGIN: iterating the elements of an element of xs is xs again.
			if root := derivedRoot(t.X, derived); root != "" {
				if k, ok := t.Key.(*ast.Ident); ok {
					derived[k.Name] = root
				}
				if v, ok := t.Value.(*ast.Ident); ok {
					derived[v.Name] = root
				}
			}
		}
		return true
	})
	return derived
}

// derivedRoot is the parameter-derived identifier an expression traces to,
// or "" when none does. It is what a loop's degree is a degree IN.
func derivedRoot(e ast.Expr, derived map[string]string) string {
	root := ""
	ast.Inspect(e, func(n ast.Node) bool {
		if id, ok := n.(*ast.Ident); ok && root == "" {
			if r, ok := derived[id.Name]; ok {
				root = r
			}
		}
		return root == ""
	})
	return root
}

// exprDerived reports whether any identifier in e is parameter-derived.
func exprDerived(e ast.Expr, derived map[string]string) bool {
	return derivedRoot(e, derived) != ""
}

// loopSource is what a loop iterates over, or nil when it cannot be read:
// `range X` gives X, and `for i := 0; i < len(X); i++` gives X.
func loopSource(n ast.Node) ast.Expr {
	switch t := n.(type) {
	case *ast.RangeStmt:
		return t.X
	case *ast.ForStmt:
		if bin, ok := t.Cond.(*ast.BinaryExpr); ok {
			for _, side := range []ast.Expr{bin.X, bin.Y} {
				if call, ok := side.(*ast.CallExpr); ok {
					if id, ok := call.Fun.(*ast.Ident); ok && id.Name == "len" && len(call.Args) == 1 {
						return call.Args[0]
					}
				}
			}
			return bin.Y
		}
	}
	return nil
}

// inductionVars are the variables a loop advances: the ones its init binds
// and its post updates, or a range statement's key and value.
func inductionVars(n ast.Node) map[string]bool {
	out := map[string]bool{}
	switch t := n.(type) {
	case *ast.RangeStmt:
		if id, ok := t.Key.(*ast.Ident); ok {
			out[id.Name] = true
		}
		if id, ok := t.Value.(*ast.Ident); ok {
			out[id.Name] = true
		}
	case *ast.ForStmt:
		if as, ok := t.Init.(*ast.AssignStmt); ok {
			for _, l := range as.Lhs {
				if id, ok := l.(*ast.Ident); ok {
					out[id.Name] = true
				}
			}
		}
		if t.Post != nil {
			ast.Inspect(t.Post, func(c ast.Node) bool {
				switch p := c.(type) {
				case *ast.IncDecStmt:
					if id, ok := p.X.(*ast.Ident); ok {
						out[id.Name] = true
					}
				case *ast.AssignStmt:
					for _, l := range p.Lhs {
						if id, ok := l.(*ast.Ident); ok {
							out[id.Name] = true
						}
					}
				}
				return true
			})
		}
	}
	return out
}

// advancesAny reports whether `n` ASSIGNS to any of `vars` -- the test for an
// inner loop that moves the enclosing loop's cursor.
func advancesAny(n ast.Node, vars map[string]bool) bool {
	found := false
	ast.Inspect(n, func(c ast.Node) bool {
		switch p := c.(type) {
		case *ast.IncDecStmt:
			if id, ok := p.X.(*ast.Ident); ok && vars[id.Name] {
				found = true
			}
		case *ast.AssignStmt:
			for _, l := range p.Lhs {
				if id, ok := l.(*ast.Ident); ok && vars[id.Name] {
					found = true
				}
			}
		}
		return !found
	})
	return found
}

// tracedLoopNesting is the maximum nesting of loops WHOSE ITERATION SOURCE
// TRACES TO A PARAMETER, plus the number of loops whose source could not be
// traced -- reported rather than silently counted or silently dropped.
func tracedLoopNesting(decl *ast.FuncDecl) (int, int) {
	derived := paramDerived(decl)
	untraced := 0
	// Exponents PER PARAMETER. Nesting a loop over `xs` inside a loop over
	// `ys` is n*m, which is neither n nor n^2; collapsing it to n^2 is the
	// same error as counting nesting depth. The reported degree is the
	// maximum single-parameter exponent, and a definition with more than one
	// is multivariate -- the single number is a projection.
	perRoot := map[string]int{}
	max := 0
	var walk func(ast.Node, map[string]int, map[string]bool)
	walk = func(x ast.Node, depth map[string]int, enclosing map[string]bool) {
		if x == nil {
			return
		}
		ast.Inspect(x, func(c ast.Node) bool {
			if c == nil || c == x {
				return c == x
			}
			switch c.(type) {
			case *ast.ForStmt, *ast.RangeStmt:
				vars := inductionVars(c)
				// SHARED INDUCTION: an inner loop that advances the enclosing
				// loop's cursor is one pass over the collection, not a
				// product. Two-pointer and cursor-scanning shapes are this:
				// `for i := 0; i < len(s); { ... for i < len(s) && p(s[i]) { i++ } }`
				// visits each element once, so it is n, not n^2.
				shared := len(enclosing) > 0 && advancesAny(c, enclosing)
				src := loopSource(c)
				next := map[string]int{}
				for k, v := range depth {
					next[k] = v
				}
				switch {
				case shared:
					// contributes nothing: the same pass as the enclosing loop
				case src == nil:
					untraced++
				default:
					if root := derivedRoot(src, derived); root != "" {
						next[root] = next[root] + 1
						if next[root] > perRoot[root] {
							perRoot[root] = next[root]
						}
						if next[root] > max {
							max = next[root]
						}
						// This loop's own counter now carries the collection's
						// size, so an inner `for j := 0; j < i; j++` is a
						// second pass over the SAME parameter rather than an
						// untraceable loop over an integer.
						for v := range vars {
							if _, seen := derived[v]; !seen {
								derived[v] = root
							}
						}
					}
					// no derived root: constant in the input
				}
				merged := map[string]bool{}
				for k := range enclosing {
					merged[k] = true
				}
				for k := range vars {
					merged[k] = true
				}
				walk(c, next, merged)
				return false
			}
			return true
		})
	}
	walk(decl, map[string]int{}, map[string]bool{})
	return max, untraced
}

// callsItself reports// exprMentions reports whether `name` appears anywhere in `e`.
func exprMentions(e ast.Expr, name string) bool {
	found := false
	ast.Inspect(e, func(n ast.Node) bool {
		if id, ok := n.(*ast.Ident); ok && id.Name == name {
			found = true
		}
		return !found
	})
	return found
}

// callsItself reports whether the declaration calls ITSELF: a plain call to
// its own name, or -- for a method -- a call on its own receiver.
//
// Matching any `x.emit()` inside `func (t *T) emit()` would be wrong: a call
// to a DIFFERENT value's method of the same name is not recursion. That
// over-match reported most Go methods as self-recursive, and since Go
// recursion shrink is not read, every one of them became `unknown`.
func callsItself(decl *ast.FuncDecl, name string) bool {
	recv := ""
	if decl.Recv != nil && len(decl.Recv.List) > 0 && len(decl.Recv.List[0].Names) > 0 {
		recv = decl.Recv.List[0].Names[0].Name
	}
	found := false
	ast.Inspect(decl, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		switch f := call.Fun.(type) {
		case *ast.Ident:
			if f.Name == name {
				found = true
			}
		case *ast.SelectorExpr:
			// `t.next.walk()` IS self-recursion -- the same method on another
			// node of the same structure -- so the receiver EXPRESSION need
			// only involve the receiver. `o.emit()` on a parameter does not.
			if f.Sel.Name == name && recv != "" && exprMentions(f.X, recv) {
				found = true
			}
		}
		return true
	})
	return found
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
		loopDepth, untraced := tracedLoopNesting(decl)
		out = append(out, callable{
			name: name, kind: kind,
			line: start, end: end,
			sloc: sloc,
			cc:   1 + decisionPoints(decl, dead),

			loopDepth:     loopDepth,
			untracedLoops: untraced,
			selfCall:      callsItself(decl, decl.Name.Name),
		})
	}
	return out, counts, nil
}

// extent is one structural node's line range: a STATEMENT or a DECLARATION,
// which is the granularity someone can actually extract. Expression nodes are
// deliberately excluded -- a duplicate region snapped to the smallest
// enclosing expression is noise, not a finding.
type extent struct {
	from int
	to   int
}

// goExtents returns the line range of every statement and declaration in the
// file, so a token-derived duplicate region can be snapped to real structure.
// Go has no s-expressions, but it has an AST, and the parse has already
// happened.
func goExtents(src string) ([]extent, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "src.go", src, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	seen := map[extent]bool{}
	var out []extent
	ast.Inspect(file, func(n ast.Node) bool {
		if n == nil {
			return false
		}
		switch n.(type) {
		case ast.Stmt, ast.Decl:
		default:
			return true
		}
		e := extent{fset.Position(n.Pos()).Line, fset.Position(n.End()).Line}
		if !seen[e] {
			seen[e] = true
			out = append(out, e)
		}
		return true
	})
	return out, nil
}
