package main

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Every expectation below is hand-computed from the CC rule: 1, plus one per
// if / for / range / non-default case / non-nil comm clause / && / ||.

const callableFixture = `package fixture

import "fmt"

// plain has one if and one && : 1 + 1 + 1 = 3.
func plain(a, b int) int {
	if a > 0 && b > 0 {
		return a + b
	}
	return 0
}

// method on a POINTER receiver: range(1) + if(1) + ||(1) = 1 + 3 = 4.
func (c *Compiler) emit(xs []int) int {
	n := 0
	for _, x := range xs {
		if x < 0 || x > 10 {
			continue
		}
		n += x
	}
	return n
}

// same method name on a different receiver: must not collide with the above.
func (c Renderer) emit() {}

// generic: 1 + range(1) + if(1) = 3.
func mapKeys[K comparable, V any](m map[K]V) []K {
	out := []K{}
	for k := range m {
		if len(out) < 100 {
			out = append(out, k)
		}
	}
	return out
}

// bare default adds nothing: 1 + two non-default cases = 3.
func classify(n int) string {
	switch n {
	case 1:
		return "one"
	case 2:
		return "two"
	default:
		return "many"
	}
}

// select: two comm clauses count, the bare default does not: 1 + 2 = 3.
func pump(a, b chan int) {
	select {
	case <-a:
	case <-b:
	default:
	}
}

// Live nested literals fold into the parent. The call-argument literal
// contributes its if, and the deferred and immediately-invoked ones
// contribute theirs: 1 + if(1) + if(1) + if(1) = 4.
func withInlineLiterals(xs []int) {
	each(xs, func(x int) {
		if x > 0 {
			fmt.Println(x)
		}
	})
	defer func() {
		if len(xs) == 0 {
			fmt.Println("empty")
		}
	}()
	func() {
		if len(xs) > 5 {
			fmt.Println("big")
		}
	}()
}

// RULE CHANGE: a live named local closure used to be split out as its own
// callable (parent 2, closure 3). It now FOLDS into the parent like any other
// live nested function, because the parent really does execute those
// branches: 1 + if(1) + if(1) + &&(1) = 4, and no separate callable.
func withNamedLiteral(xs []int) int {
	score := func(x int) int {
		if x > 0 && x < 10 {
			return x
		}
		return 0
	}
	total := 0
	if len(xs) > 0 {
		total = score(xs[0])
	}
	return total
}

// A named literal nothing references is DEAD. Parent CC stays 1.
func withDeadLiteral() {
	unused := func(x int) int {
		if x > 0 {
			return x
		}
		return -x
	}
	_ = 1
}

// no body: skipped entirely.
func external(x int) int
`

func callablesByName(t *testing.T, src string) map[string]callable {
	t.Helper()
	cs, err := goCallables(src)
	if err != nil {
		t.Fatalf("goCallables: %v", err)
	}
	out := map[string]callable{}
	for _, c := range cs {
		out[c.name] = c
	}
	return out
}

func TestGoCallablesComplexity(t *testing.T) {
	got := callablesByName(t, callableFixture)

	for _, tc := range []struct {
		name string
		cc   int
		kind string
	}{
		{"plain", 3, "func"},
		{"(*Compiler).emit", 4, "method"},
		{"(Renderer).emit", 1, "method"},
		{"mapKeys", 3, "func"},
		{"classify", 3, "func"},
		{"pump", 3, "func"},
		{"withInlineLiterals", 4, "func"},
		{"withNamedLiteral", 4, "func"},
		{"withDeadLiteral", 1, "func"},
		{"withDeadLiteral.unused", 2, "dead-closure"},
	} {
		c, ok := got[tc.name]
		if !ok {
			t.Errorf("%s: not reported", tc.name)
			continue
		}
		if c.cc != tc.cc {
			t.Errorf("%s: cc = %d, want %d", tc.name, c.cc, tc.cc)
		}
		if c.kind != tc.kind {
			t.Errorf("%s: kind = %q, want %q", tc.name, c.kind, tc.kind)
		}
	}
}

func TestGoCallablesMethodsDoNotCollide(t *testing.T) {
	got := callablesByName(t, callableFixture)
	if _, ok := got["(*Compiler).emit"]; !ok {
		t.Error("pointer-receiver method missing its receiver in the name")
	}
	if _, ok := got["(Renderer).emit"]; !ok {
		t.Error("value-receiver method missing its receiver in the name")
	}
	if _, ok := got["emit"]; ok {
		t.Error("a method was reported without its receiver")
	}
}

func TestGoCallablesLiveNestedFunctionsFoldIntoParent(t *testing.T) {
	cs, err := goCallables(callableFixture)
	if err != nil {
		t.Fatal(err)
	}
	// The enumeration is top-level declarations. The ONLY non-declaration
	// entry allowed is a dead named binding, which is not a callable at all.
	for _, c := range cs {
		if c.kind != "func" && c.kind != "method" && c.name != "withDeadLiteral.unused" {
			t.Errorf("nested function reported as its own callable: %s (%s)", c.name, c.kind)
		}
	}
	// A live named closure folds: its branches count toward the declaration.
	got := callablesByName(t, callableFixture)
	if _, ok := got["withNamedLiteral.score"]; ok {
		t.Error("a LIVE named closure must fold into its parent, not be enumerated")
	}
}

func TestGoCallablesNamedLiteralDeadness(t *testing.T) {
	got := callablesByName(t, callableFixture)
	if _, ok := got["withNamedLiteral.score"]; ok {
		t.Error("a referenced named closure is live and must not be reported as dead")
	}
	if c, ok := got["withDeadLiteral.unused"]; !ok || c.kind != "dead-closure" {
		t.Error("a named closure nothing references must be reported as dead")
	}
	// A dead closure's complexity must NOT inflate its parent.
	if c := got["withDeadLiteral"]; c.cc != 1 {
		t.Errorf("parent cc = %d, want 1: the dead closure's branch leaked in", c.cc)
	}
}

func TestGoCallablesSLOCExcludesBlankAndComments(t *testing.T) {
	src := `package p

// doc comment, not counted
func f() int {

	// interior comment, not counted
	x := 1

	return x
}
`
	got := callablesByName(t, src)
	// code lines: "func f() int {", "x := 1", "return x", "}" = 4
	if c := got["f"]; c.sloc != 4 {
		t.Errorf("sloc = %d, want 4", c.sloc)
	}
}

func TestGoCallablesParentSLOCExcludesOnlyDeadClosures(t *testing.T) {
	got := callablesByName(t, callableFixture)
	// A live nested closure's lines stay with the parent.
	live := got["withNamedLiteral"]
	if live.sloc != live.end-live.line+1 {
		t.Errorf("live nesting: sloc %d should be every code line of the span (%d)",
			live.sloc, live.end-live.line+1)
	}
	// A dead one's lines are charged to it, not to the parent.
	parent := got["withDeadLiteral"]
	deadC := got["withDeadLiteral.unused"]
	if deadC.sloc < 1 {
		t.Fatalf("dead closure sloc = %d", deadC.sloc)
	}
	if parent.sloc >= parent.end-parent.line+1 {
		t.Errorf("parent sloc %d must exclude the dead closure's %d lines",
			parent.sloc, deadC.sloc)
	}
}

func TestGoCallablesSkipsBodylessDeclarations(t *testing.T) {
	got := callablesByName(t, callableFixture)
	if _, ok := got["external"]; ok {
		t.Error("a declaration with no body has no complexity and must be skipped")
	}
}

func TestGoCallablesRejectsUnparseableSource(t *testing.T) {
	if _, err := goCallables("package p\nfunc ("); err == nil {
		t.Error("expected a parse error")
	}
}

func TestEdnStringEscapes(t *testing.T) {
	for _, tc := range []struct{ in, want string }{
		{`plain`, `"plain"`},
		{`has "quotes"`, `"has \"quotes\""`},
		{`back\slash`, `"back\\slash"`},
		{"tab\there", `"tab\there"`},
	} {
		if got := ednString(tc.in); got != tc.want {
			t.Errorf("ednString(%q) = %s, want %s", tc.in, got, tc.want)
		}
	}
}

func TestGoFilesSkipsGeneratedAndLoweredTrees(t *testing.T) {
	dir := t.TempDir()
	mustWrite := func(rel string) {
		p := dir + "/" + rel
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte("package p\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	mustWrite("keep.go")
	mustWrite("sub/also_keep.go")
	mustWrite("skip_generated.go")
	mustWrite("core_go_lowered/lowered.go")
	mustWrite("notes.txt")

	got, err := goFiles([]string{dir})
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, g := range got {
		names = append(names, strings.TrimPrefix(g, dir+"/"))
	}
	want := []string{"keep.go", "sub/also_keep.go"}
	if len(names) != len(want) {
		t.Fatalf("got %v, want %v", names, want)
	}
	for i := range want {
		if names[i] != want[i] {
			t.Errorf("got %v, want %v", names, want)
			break
		}
	}
}

func TestCountFileCensus(t *testing.T) {
	src := `package p

// a comment line
func f() int {
	return 1
}
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "src.go", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}
	got := countFile(codeLines(fset, file, src), src)
	// 6 lines: package, blank, comment, func, return, }
	if got.lines != 6 {
		t.Errorf("lines = %d, want 6", got.lines)
	}
	if got.blank != 1 {
		t.Errorf("blank = %d, want 1", got.blank)
	}
	// code = package, func, return, } = 4 (the comment line is not code)
	if got.code != 4 {
		t.Errorf("sloc = %d, want 4", got.code)
	}
	if comment := got.lines - got.code - got.blank; comment != 1 {
		t.Errorf("comment lines = %d, want 1", comment)
	}
}

func TestGoExtentsAreStatementsAndDeclarations(t *testing.T) {
	src := `package p

import "fmt"

func f(xs []int) int {
	total := 0
	for _, x := range xs {
		total += x + 1
	}
	return total
}
`
	es, err := goExtents(src)
	if err != nil {
		t.Fatal(err)
	}
	has := func(from, to int) bool {
		for _, e := range es {
			if e.from == from && e.to == to {
				return true
			}
		}
		return false
	}
	// the declaration itself
	if !has(5, 11) {
		t.Errorf("missing the func declaration extent, got %v", es)
	}
	// a multi-line statement someone could lift out
	if !has(7, 9) {
		t.Errorf("missing the for-statement extent, got %v", es)
	}
	// the import declaration
	if !has(3, 3) {
		t.Errorf("missing the import declaration extent, got %v", es)
	}
	// EXPRESSIONS are deliberately absent: snapping to the smallest enclosing
	// expression would be noise. `x + 1` sits inside line 8, which is a
	// statement extent, but must not appear as its own zero-width entry more
	// than the statement already provides.
	count8 := 0
	for _, e := range es {
		if e.from == 8 && e.to == 8 {
			count8++
		}
	}
	if count8 != 1 {
		t.Errorf("line 8 should contribute exactly one statement extent, got %d", count8)
	}
}
