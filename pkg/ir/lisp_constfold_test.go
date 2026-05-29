/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package ir_test

import (
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/resolver"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

// Lisp-pass regression tests. Each entry builds an IR via the Lisp
// `ir.build/build-fn`, runs a Lisp pass against that Function (via
// the runtime), and asserts behavioral properties of the resulting
// dump. The Go-side bytecode→IR path was retired; these tests are
// the per-pass coverage for the Lisp passes that remain.

// loaderOnce wires the on-demand namespace loader once. Without it,
// `ir.build/build-fn` resolves to a stub Var whose root is nil.
var loaderOnce sync.Once

func ensureLoader() {
	loaderOnce.Do(func() {
		consts := vm.NewConsts()
		ctx := compiler.NewCompiler(consts, rt.NS(rt.NameCoreNS))
		res := resolver.NewNSResolver(ctx, []string{"."})
		rt.SetNSLoader(res)
		// Phase F: ir.data must load BEFORE any other IR ns whose
		// bytecode resolves ir/* symbols (op, refs, aux, etc.) — those
		// are interned into the `ir` namespace by data.lg's bottom
		// block.
		if res.Load("ir.data") == nil {
			panic("ir.data namespace failed to load — Phase F not active")
		}
		if res.Load("ir.build") == nil {
			panic("ir.build namespace failed to load — bundle missing or corrupt")
		}
		for _, ns := range []string{
			"ir.zipper", "ir.passes",
			"ir.passes.dce", "ir.passes.constfold",
			"ir.passes.mutability", "ir.passes.cse",
			"ir.passes.typeinfer", "ir.passes.infer-arg-types",
			"ir.passes.licm",
			"ir.passes.pipeline", "ir.dump", "ir.dominance", "ir.lower-go"} {
			if res.Load(ns) == nil {
				panic("namespace failed to load: " + ns)
			}
		}
	})
}

// buildLispIR builds an IR Function via the Lisp ir.build/build-fn.
//
// Phase F: the result is a Lisp atom (vm.Value) wrapping the
// data-shape map — opaque to Go. Earlier this returned *ir.Function
// (Go struct); see history.
func buildLispIR(t *testing.T, src string) vm.Value {
	t.Helper()
	consts := vm.NewConsts()
	c := compiler.NewCompiler(consts, rt.NS(rt.NameCoreNS))
	c.SetSource("lisp-pass-build")
	expr := fmt.Sprintf(`(ir.build/build-fn (quote %s))`, src)
	_, result, err := c.CompileMultiple(strings.NewReader(expr))
	if err != nil {
		t.Fatalf("lisp build: %v", err)
	}
	return result
}

// lispDump evaluates `(ir.dump/dump f)` via the runtime and returns
// the resulting String. The Go ir.Dump was retired in Phase D.
func lispDump(t *testing.T, f vm.Value) string {
	t.Helper()
	passVarCounter++
	varName := fmt.Sprintf("*dump-fn-%d*", passVarCounter)
	coreNS := rt.NS(rt.NameCoreNS)
	coreNS.Def(varName, f)

	consts := vm.NewConsts()
	c := compiler.NewCompiler(consts, coreNS)
	c.SetSource("lisp-dump")
	expr := fmt.Sprintf(`(ir.dump/dump %s)`, varName)
	_, result, err := c.CompileMultiple(strings.NewReader(expr))
	if err != nil {
		t.Fatalf("eval %s: %v", expr, err)
	}
	return string(result.(vm.String))
}

// runLispPass installs `f` as a Var in the core ns under a fresh
// name, evals `(passNS/passFn the-var)` against the runtime, and
// returns the (same, mutated-in-place) Function value.
var passVarCounter int

func runLispPass(t *testing.T, passNS, passFn string, f vm.Value) vm.Value {
	t.Helper()
	passVarCounter++
	varName := fmt.Sprintf("*pass-fn-%d*", passVarCounter)

	coreNS := rt.NS(rt.NameCoreNS)
	coreNS.Def(varName, f)

	expr := fmt.Sprintf(`(%s/%s %s)`, passNS, passFn, varName)
	consts := vm.NewConsts()
	c := compiler.NewCompiler(consts, coreNS)
	c.SetSource("run-lisp-pass")
	if _, _, err := c.CompileMultiple(strings.NewReader(expr)); err != nil {
		t.Fatalf("eval %s: %v", expr, err)
	}
	return f
}

// runLispExpr evaluates a Lisp expression via the runtime compiler
// and returns the result. Used by tests that need direct expression evaluation
// without constructing IR via ir.build/build-fn.
func runLispExpr(t *testing.T, expr string) vm.Value {
	t.Helper()
	consts := vm.NewConsts()
	c := compiler.NewCompiler(consts, rt.NS(rt.NameCoreNS))
	c.SetSource("lisp-ir-test")
	_, result, err := c.CompileMultiple(strings.NewReader(expr))
	if err != nil {
		t.Fatalf("eval %s: %v", expr, err)
	}
	return result
}

// constFoldCase — one regression case. `mustContain` lists substrings
// the dump MUST contain after constfold; `mustNotContain` lists
// substrings it must NOT contain. Behavioral assertion, robust to
// inst-id renumbering and incidental dump-format tweaks.
type constFoldCase struct {
	name           string
	src            string
	mustContain    []string
	mustNotContain []string
}

func TestLispConstFold(t *testing.T) {
	ensureLoader()

	cases := []constFoldCase{
		// Strategy 1 — primitive fold
		{name: "const-arith", src: `(defn const-arith [] (+ 1 (* 2 3)))`,
			mustContain:    []string{"Const ; 7"},
			mustNotContain: []string{"Add", "Mul"}},
		{name: "const-sub", src: `(defn const-sub [] (- 10 4))`,
			mustContain:    []string{"Const ; 6"},
			mustNotContain: []string{"Sub"}},
		// Strategy 2 — algebraic identity. The dead Add/Mul stays in
		// the inst list (DCE's job to remove it later); what matters
		// is that uses get redirected past it. Assert via the Return.
		{name: "add-zero", src: `(defn add-zero [x] (+ x 0))`,
			mustContain: []string{"v0 = LoadArg ; 0", "Return v0"}},
		{name: "mul-one", src: `(defn mul-one [x] (* x 1))`,
			mustContain: []string{"v0 = LoadArg ; 0", "Return v0"}},
		// mul-zero IS rewritten in place (fold-this! writes Const ; 0
		// over the Mul), so the Mul opcode itself disappears.
		{name: "mul-zero", src: `(defn mul-zero [x] (* x 0))`,
			mustContain:    []string{"Const ; 0"},
			mustNotContain: []string{"Mul"}},
		// Strategy 3 — commutative canonicalization (Const goes right).
		// Const operand should be the SECOND ref of the Add/Mul.
		{name: "canon-add", src: `(defn canon-add [x] (+ 5 x))`,
			mustContain: []string{"Add v0 v1", "v0 = LoadArg", "v1 = Const ; 5"}},
		{name: "canon-mul", src: `(defn canon-mul [x] (* 7 x))`,
			mustContain: []string{"Mul v0 v1", "v0 = LoadArg", "v1 = Const ; 7"}},
		// Non-foldable shapes — pass must round-trip without error
		// and preserve the un-foldable op.
		{name: "id", src: `(defn id [x] x)`,
			mustContain: []string{"LoadArg"}},
		{name: "add", src: `(defn add [a b] (+ a b))`,
			mustContain: []string{"Add"}},
		{name: "use-let", src: `(defn use-let [x] (let [y 1] (+ x y)))`,
			mustContain: []string{"Add"}},
		// Numeric tower — BigInt i64 overflow refuses to fold (matches
		// Go's foldNumeric ok=false branch). safe-apply in constfold.lg
		// resolves spec open question #1.
		{name: "bigint-fold", src: `(defn bigint-fold [] (+ 9223372036854775807 1))`,
			mustContain: []string{"Add"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fn := buildLispIR(t, tc.src)
			runLispPass(t, "ir.passes.constfold", "constfold", fn)
			dump := lispDump(t, fn)
			for _, want := range tc.mustContain {
				if !strings.Contains(dump, want) {
					t.Errorf("dump missing %q\n--- dump ---\n%s", want, dump)
				}
			}
			for _, unwanted := range tc.mustNotContain {
				if strings.Contains(dump, unwanted) {
					t.Errorf("dump unexpectedly contains %q\n--- dump ---\n%s", unwanted, dump)
				}
			}
		})
	}
}
