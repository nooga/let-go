package ir_test

import (
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

// isValidGoIdent reports whether s is a legal (exported) Go identifier:
// starts with a letter, then letters/digits/underscore only. Mirrors the
// property `go build` enforces on generated func names.
func isValidGoIdent(s string) bool {
	if len(s) == 0 {
		return false
	}
	for i, r := range s {
		letter := (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '_'
		digit := r >= '0' && r <= '9'
		if i == 0 && !letter {
			return false
		}
		if i > 0 && !letter && !digit {
			return false
		}
	}
	return true
}

func exportName(t *testing.T, src string) string {
	t.Helper()
	ensureLoader()
	v := runLispExpr(t, `(ir.lower-go/go-export-name "`+src+`")`)
	s, ok := v.(vm.String)
	if !ok {
		t.Fatalf("go-export-name %q: expected String, got %T (%v)", src, v, v)
	}
	return string(s)
}

// P1 finding #1: a source name containing a Lisp-legal char with no explicit
// munge token (@, #, ., etc.) must still produce a VALID Go identifier, not a
// verbatim illegal one like "Foo@bar". Names already handled today (predicates,
// bang-mutators, operators) must keep their exact PascalCase.
func TestGoExportNameAlwaysValidIdent(t *testing.T) {
	// Illegal-char cases the current munge map misses (the RED cases).
	for _, c := range []string{"foo@bar", "ns.member", "a#b"} {
		got := exportName(t, c)
		if !isValidGoIdent(got) {
			t.Errorf("go-export-name %q -> %q is not a valid Go identifier", c, got)
		}
	}
	// Byte-identity guard: currently-handled names must not drift.
	stable := map[string]string{
		"vector?": "VectorQmark", "conj!": "ConjBang", "map*": "MapStar",
		"do-thing": "DoThing", "+": "Plus", "empty?": "EmptyQmark",
	}
	for src, want := range stable {
		if got := exportName(t, src); got != want {
			t.Errorf("go-export-name %q drifted: got %q, want %q", src, got, want)
		}
	}
}

// P1 finding #2: PascalCase is not injective — foo? and foo-qmark both fall to
// FooQmark. When BOTH live in one namespace the lowering must emit two DISTINCT
// exported func names (uncompilable duplicate decls otherwise). A namespace with
// only one of them stays byte-identical (FooQmark).
func TestExportNameCollisionResolvedWithinNamespace(t *testing.T) {
	ensureLoader()
	src := lowerOneNsGo(t, "coll",
		`(defn foo? [x] x)`,
		`(defn foo-qmark [x] x)`)
	names := funcDeclNames(src)
	if len(names) < 2 {
		t.Fatalf("expected >=2 top-level func decls, got %d in:\n%s", len(names), src)
	}
	seen := map[string]bool{}
	for _, n := range names {
		if seen[n] {
			t.Fatalf("duplicate exported func %q — collision not resolved:\n%s", n, src)
		}
		seen[n] = true
	}
}

// lowerOneNsGo renders a single namespace's forms to Go via the pipeline and
// returns the source string. Models the existing crosspkg wrapper test.
func lowerOneNsGo(t *testing.T, ns string, forms ...string) string {
	t.Helper()
	quoted := make([]string, len(forms))
	for i, f := range forms {
		quoted[i] = "(quote " + f + ")"
	}
	// The namespace must exist (interned defns) before lowering so `the-ns`
	// resolves it — declare it, then switch back to user, then lower.
	setup := `(ns ` + ns + `) ` + strings.Join(forms, " ") + ` (ns user) `
	expr := setup + `(nth (ir.passes.pipeline/lower-all-ns-to-go
	  [["` + ns + `" (quote ` + ns + `) [` + strings.Join(quoted, " ") + `] "ex/` + ns + `"]]) 0)`
	v := runLispExpr(t, expr)
	s, ok := v.(vm.String)
	if !ok {
		t.Fatalf("lower-all-ns-to-go: expected String, got %T (%v)", v, v)
	}
	return string(s)
}

// funcDeclNames returns the names of top-level `func <Name>(` declarations
// (skips methods `func (recv T) Name(`).
func funcDeclNames(src string) []string {
	var out []string
	for _, line := range strings.Split(src, "\n") {
		rest, ok := strings.CutPrefix(line, "func ")
		if !ok || strings.HasPrefix(rest, "(") {
			continue
		}
		if i := strings.IndexByte(rest, '('); i > 0 {
			out = append(out, strings.TrimSpace(rest[:i]))
		}
	}
	return out
}
