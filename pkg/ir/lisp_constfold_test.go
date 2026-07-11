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
			"ir.passes.licm", "ir.passes.lambda-lift", "ir.passes.inline", "ir.passes.fusion",
			"ir.passes.liveness", "ir.passes.blockarg", "ir.passes.cleanup", "ir.passes.pipeline", "ir.dump", "ir.dominance", "ir.lower-go"} {
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

// buildLispIRWith interns each stub as a var in a FRESH uniquely-named ns,
// then builds src in that ns so build-fn can resolve the stub symbols.
// stubs maps a symbol name -> a lambda source string, e.g. {"chr": "(fn* [p c] c)"}.
func buildLispIRWith(t *testing.T, stubs map[string]string, src string) vm.Value {
	t.Helper()
	passVarCounter++
	nsName := fmt.Sprintf("lltest%d", passVarCounter)
	ns := rt.NS(nsName)
	consts := vm.NewConsts()
	c := compiler.NewCompiler(consts, ns)
	c.SetSource("ll-build")
	var b strings.Builder
	for name, lam := range stubs {
		fmt.Fprintf(&b, "(def %s %s)\n", name, lam)
	}
	fmt.Fprintf(&b, "(ir.build/build-fn (quote %s))", src)
	_, result, err := c.CompileMultiple(strings.NewReader(b.String()))
	if err != nil {
		t.Fatalf("buildLispIRWith: %v", err)
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

// runLispPassValue is like runLispPass but returns the pass's return value
// instead of the mutated input. Used for passes that return a complex result
// (e.g., lambda-lift returning {:main f :lifted [...]}).
func runLispPassValue(t *testing.T, passNS, passFn string, f vm.Value) vm.Value {
	t.Helper()
	passVarCounter++
	varName := fmt.Sprintf("*pass-fn-%d*", passVarCounter)

	coreNS := rt.NS(rt.NameCoreNS)
	coreNS.Def(varName, f)

	expr := fmt.Sprintf(`(%s/%s %s)`, passNS, passFn, varName)
	consts := vm.NewConsts()
	c := compiler.NewCompiler(consts, coreNS)
	c.SetSource("run-lisp-pass")
	_, result, err := c.CompileMultiple(strings.NewReader(expr))
	if err != nil {
		t.Fatalf("eval %s: %v", expr, err)
	}
	return result
}

// lispEvalOn installs a vm.Value as a fresh core var, evaluates the
// format string as a Lisp expression, and returns its pr-str output.
func lispEvalOn(t *testing.T, val vm.Value, format string) string {
	t.Helper()
	passVarCounter++
	varName := fmt.Sprintf("*eval-val-%d*", passVarCounter)

	coreNS := rt.NS(rt.NameCoreNS)
	coreNS.Def(varName, val)

	// Wrap the expression in (pr-str ...) to convert result to string
	expr := fmt.Sprintf("(pr-str %s)", fmt.Sprintf(format, varName))
	consts := vm.NewConsts()
	c := compiler.NewCompiler(consts, coreNS)
	c.SetSource("lisp-eval-on")
	_, result, err := c.CompileMultiple(strings.NewReader(expr))
	if err != nil {
		t.Fatalf("eval %s: %v", expr, err)
	}
	return string(result.(vm.String))
}

// lispDumpValue evaluates (ir.dump/dump val) and returns the result as a string.
// Wrapper around lispDump for clearer naming in tests.
func lispDumpValue(t *testing.T, val vm.Value) string {
	t.Helper()
	return lispDump(t, val)
}

// lispEvalReturn is like lispEvalOn but returns the raw vm.Value result
// instead of its pr-str.
func lispEvalReturn(t *testing.T, val vm.Value, format string) vm.Value {
	t.Helper()
	passVarCounter++
	varName := fmt.Sprintf("*eval-val-%d*", passVarCounter)

	coreNS := rt.NS(rt.NameCoreNS)
	coreNS.Def(varName, val)

	expr := fmt.Sprintf(format, varName)
	consts := vm.NewConsts()
	c := compiler.NewCompiler(consts, coreNS)
	c.SetSource("lisp-eval-return")
	_, result, err := c.CompileMultiple(strings.NewReader(expr))
	if err != nil {
		t.Fatalf("eval %s: %v", expr, err)
	}
	return result
}

// lispEval installs a vm.Value as a fresh core var, evaluates the format
// string as a Lisp expression with the var substituted, and returns the
// pr-str output of the result. The format string uses %s as a placeholder
// for the var name (e.g., "(count %s)").
func lispEval(t *testing.T, format string, val vm.Value) string {
	t.Helper()
	passVarCounter++
	varName := fmt.Sprintf("*eval-val-%d*", passVarCounter)

	coreNS := rt.NS(rt.NameCoreNS)
	coreNS.Def(varName, val)

	// Wrap the expression in (pr-str ...) to convert result to string
	expr := fmt.Sprintf("(pr-str %s)", fmt.Sprintf(format, varName))
	consts := vm.NewConsts()
	c := compiler.NewCompiler(consts, coreNS)
	c.SetSource("lisp-eval")
	_, result, err := c.CompileMultiple(strings.NewReader(expr))
	if err != nil {
		t.Fatalf("eval %s: %v", expr, err)
	}
	return string(result.(vm.String))
}

// inlineWithRegistry builds a caller function and its callees (as stubs),
// constructs a symbol->Function registry, applies ir.passes.inline/inline-with,
// and returns the caller's dump. stubs maps callee symbol names to defn source,
// e.g., {"tiny": "(defn tiny [x] (+ x 1))"}.
func inlineWithRegistry(t *testing.T, callerSrc string, stubs map[string]string) string {
	t.Helper()
	passVarCounter++
	nsName := fmt.Sprintf("inlinetest%d", passVarCounter)
	ns := rt.NS(nsName)

	// First, build all stubs as Functions. We need to define the stub symbols
	// so that the caller can reference them during build.
	regEntries := make([]string, 0, len(stubs))
	for name, src := range stubs {
		// Build the stub function
		consts := vm.NewConsts()
		c := compiler.NewCompiler(consts, ns)
		c.SetSource("inline-stub-build")
		expr := fmt.Sprintf("(ir.build/build-fn (quote %s))", src)
		_, result, err := c.CompileMultiple(strings.NewReader(expr))
		if err != nil {
			t.Fatalf("build stub %s: %v", name, err)
		}
		// Store the Function in the namespace so we can reference it
		varName := fmt.Sprintf("*stub-%s*", name)
		ns.Def(varName, result)
		// Also define the stub symbol itself so build-fn can resolve it when building the caller
		ns.Def(name, result)
		// Add to registry entries using the fully qualified symbol (ns/name)
		// This matches what call-head-var will return when converting a Var to a symbol
		qualifiedName := fmt.Sprintf("%s/%s", nsName, name)
		regEntries = append(regEntries, fmt.Sprintf("'%s %s", qualifiedName, varName))
	}

	// Build the caller function in the same namespace
	consts := vm.NewConsts()
	c := compiler.NewCompiler(consts, ns)
	c.SetSource("inline-caller-build")
	expr := fmt.Sprintf("(ir.build/build-fn (quote %s))", callerSrc)
	_, callerFunc, err := c.CompileMultiple(strings.NewReader(expr))
	if err != nil {
		t.Fatalf("build caller: %v", err)
	}

	// Store caller in namespace
	varName := fmt.Sprintf("*caller-func*")
	ns.Def(varName, callerFunc)

	// Now run inline-with with the registry
	regMapExpr := fmt.Sprintf("(hash-map %s)", strings.Join(regEntries, " "))
	inlineExpr := fmt.Sprintf(
		`(do
		  (ir.passes.inline/inline-with %s %s)
		  %s)`,
		varName, regMapExpr, varName)

	consts = vm.NewConsts()
	c = compiler.NewCompiler(consts, ns)
	c.SetSource("inline-with-registry")
	_, resultFunc, err := c.CompileMultiple(strings.NewReader(inlineExpr))
	if err != nil {
		t.Fatalf("inline-with: %v", err)
	}

	// Dump and return
	return lispDump(t, resultFunc)
}

// optimizeFnWithRegistry builds a caller + stubs like inlineWithRegistry, but
// instead of calling inline-with directly it binds ir.passes.inline/*inline-registry*
// to the registry and runs the FULL ir.passes.pipeline/optimize-fn. This verifies
// the pass is wired into the pipeline (not just callable in isolation).
func optimizeFnWithRegistry(t *testing.T, callerSrc string, stubs map[string]string) string {
	t.Helper()
	passVarCounter++
	nsName := fmt.Sprintf("optinlinetest%d", passVarCounter)
	ns := rt.NS(nsName)

	regEntries := make([]string, 0, len(stubs))
	for name, src := range stubs {
		consts := vm.NewConsts()
		c := compiler.NewCompiler(consts, ns)
		c.SetSource("optinline-stub-build")
		expr := fmt.Sprintf("(ir.build/build-fn (quote %s))", src)
		_, result, err := c.CompileMultiple(strings.NewReader(expr))
		if err != nil {
			t.Fatalf("build stub %s: %v", name, err)
		}
		varName := fmt.Sprintf("*stub-%s*", name)
		ns.Def(varName, result)
		ns.Def(name, result)
		qualifiedName := fmt.Sprintf("%s/%s", nsName, name)
		regEntries = append(regEntries, fmt.Sprintf("'%s %s", qualifiedName, varName))
	}

	consts := vm.NewConsts()
	c := compiler.NewCompiler(consts, ns)
	c.SetSource("optinline-caller-build")
	expr := fmt.Sprintf("(ir.build/build-fn (quote %s))", callerSrc)
	_, callerFunc, err := c.CompileMultiple(strings.NewReader(expr))
	if err != nil {
		t.Fatalf("build caller: %v", err)
	}
	varName := "*caller-func*"
	ns.Def(varName, callerFunc)

	regMapExpr := fmt.Sprintf("(hash-map %s)", strings.Join(regEntries, " "))
	optExpr := fmt.Sprintf(
		`(do
		  (binding [ir.passes.inline/*enable-inline* true
		            ir.passes.inline/*inline-registry* %s]
		    (ir.passes.pipeline/optimize-fn %s))
		  %s)`,
		regMapExpr, varName, varName)

	consts = vm.NewConsts()
	c = compiler.NewCompiler(consts, ns)
	c.SetSource("optimize-fn-with-registry")
	_, resultFunc, err := c.CompileMultiple(strings.NewReader(optExpr))
	if err != nil {
		t.Fatalf("optimize-fn: %v", err)
	}
	return lispDump(t, resultFunc)
}

// assertInlineValidates builds caller+callees, runs inline-with, then runs
// ir.validate/validate-fn! on the result. validate-fn! throws on malformed IR
// (cross-block refs, branch-arg arity, branch-if asymmetry, …), which surfaces
// here as a compile/eval error and fails the test.
func assertInlineValidates(t *testing.T, callerSrc string, stubs map[string]string) {
	t.Helper()
	passVarCounter++
	nsName := fmt.Sprintf("inlinevalidtest%d", passVarCounter)
	ns := rt.NS(nsName)

	regEntries := make([]string, 0, len(stubs))
	for name, src := range stubs {
		consts := vm.NewConsts()
		c := compiler.NewCompiler(consts, ns)
		c.SetSource("inline-valid-stub")
		expr := fmt.Sprintf("(ir.build/build-fn (quote %s))", src)
		_, result, err := c.CompileMultiple(strings.NewReader(expr))
		if err != nil {
			t.Fatalf("build stub %s: %v", name, err)
		}
		varName := fmt.Sprintf("*stub-%s*", name)
		ns.Def(varName, result)
		ns.Def(name, result)
		qualifiedName := fmt.Sprintf("%s/%s", nsName, name)
		regEntries = append(regEntries, fmt.Sprintf("'%s %s", qualifiedName, varName))
	}

	consts := vm.NewConsts()
	c := compiler.NewCompiler(consts, ns)
	c.SetSource("inline-valid-caller")
	expr := fmt.Sprintf("(ir.build/build-fn (quote %s))", callerSrc)
	_, callerFunc, err := c.CompileMultiple(strings.NewReader(expr))
	if err != nil {
		t.Fatalf("build caller: %v", err)
	}
	ns.Def("*caller-func*", callerFunc)

	regMapExpr := fmt.Sprintf("(hash-map %s)", strings.Join(regEntries, " "))
	validateExpr := fmt.Sprintf(
		`(do
		  (ir.passes.inline/inline-with *caller-func* %s)
		  (ir.validate/validate-fn! *caller-func* "inline-multiblock-test"))`,
		regMapExpr)

	consts = vm.NewConsts()
	c = compiler.NewCompiler(consts, ns)
	c.SetSource("inline-validate")
	if _, _, err := c.CompileMultiple(strings.NewReader(validateExpr)); err != nil {
		t.Fatalf("inlined IR failed validation: %v", err)
	}
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
		{name: "const-quot-int", src: `(defn const-quot-int [] (quot 7 3))`,
			mustContain:    []string{"Const ; 2"},
			mustNotContain: []string{"Quot"}},
		{name: "const-div-float", src: `(defn const-div-float [] (/ 7.5 2.5))`,
			mustContain:    []string{"Const ; 3"},
			mustNotContain: []string{"Div"}},
		{name: "const-bit-and", src: `(defn const-bit-and [] (bit-and 14 12))`,
			mustContain:    []string{"Const ; 12"},
			mustNotContain: []string{"Call", "LoadVar"}},
		{name: "const-bit-or", src: `(defn const-bit-or [] (bit-or 5 3))`,
			mustContain:    []string{"Const ; 7"},
			mustNotContain: []string{"Call", "LoadVar"}},
		{name: "const-bit-xor", src: `(defn const-bit-xor [] (bit-xor 5 3))`,
			mustContain:    []string{"Const ; 6"},
			mustNotContain: []string{"Call", "LoadVar"}},
		{name: "const-bit-not", src: `(defn const-bit-not [] (bit-not 1))`,
			mustContain:    []string{"Const ; -2"},
			mustNotContain: []string{"Call", "LoadVar"}},
		{name: "const-bit-shift-left", src: `(defn const-bit-shift-left [] (bit-shift-left 1 4))`,
			mustContain:    []string{"Const ; 16"},
			mustNotContain: []string{"Call", "LoadVar"}},
		{name: "const-bit-shift-right", src: `(defn const-bit-shift-right [] (bit-shift-right 8 2))`,
			mustContain:    []string{"Const ; 2"},
			mustNotContain: []string{"Call", "LoadVar"}},
		{name: "const-unsigned-bit-shift-right", src: `(defn const-unsigned-bit-shift-right [] (unsigned-bit-shift-right -1 1))`,
			mustContain:    []string{"Const ; 9223372036854775807"},
			mustNotContain: []string{"Call", "LoadVar"}},
		{name: "const-bit-and-not", src: `(defn const-bit-and-not [] (bit-and-not 7 3))`,
			mustContain:    []string{"Const ; 4"},
			mustNotContain: []string{"Call", "LoadVar"}},
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
		{name: "shadowed-builtin-stays-call", src: `(defn shadowed-builtin-stays-call [bit-and] (bit-and 14 12))`,
			mustContain: []string{"LoadArg", "Call"}},
		// Numeric tower — BigInt i64 overflow refuses to fold (matches
		// Go's foldNumeric ok=false branch). safe-apply in constfold.lg
		// resolves spec open question #1.
		{name: "bigint-fold", src: `(defn bigint-fold [] (+ 9223372036854775807 1))`,
			mustContain: []string{"Add"}},
		{name: "quot-divide-by-zero", src: `(defn quot-divide-by-zero [] (quot 1 0))`,
			mustContain: []string{"Quot"}},
		{name: "div-int-int-stays-div", src: `(defn div-int-int-stays-div [] (/ 3 2))`,
			mustContain: []string{"Div"}},
		{name: "quot-mixed-stays-quot", src: `(defn quot-mixed-stays-quot [] (quot 7 2.0))`,
			mustContain: []string{"Quot"}},
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

func TestLambdaLiftFindsOneCandidate(t *testing.T) {
	ensureLoader()
	// `identity` is a core fn (resolves); the inner fn* is the candidate.
	f := buildLispIR(t, `(defn outer [p] (identity (fn* [q] q)))`)
	n := lispEval(t, `(count (ir.passes.lambda-lift/candidate-info %s))`, f)
	if strings.TrimSpace(n) != "1" {
		t.Fatalf("want 1 candidate, got %s", n)
	}
}

func TestInlineResolvesLoadVarCallee(t *testing.T) {
	ensureLoader()
	// Build a function that calls a stub callee.
	// Use buildLispIRWith to define the callee as a stub.
	f := buildLispIRWith(t, map[string]string{"callee": "(fn* [p] p)"}, `(defn caller [p] (callee p))`)

	// Test 1: find the first :call instruction
	callNIDStr := lispEval(t, `(ir.passes.inline/first-call-nid %s)`, f)
	callNIDStr = strings.TrimSpace(callNIDStr)
	if callNIDStr == "nil" {
		t.Fatalf("no call instruction found in caller function")
	}

	// Test 2: Install f as a var for use in multiple expressions
	passVarCounter++
	varName := fmt.Sprintf("*inline-test-f-%d*", passVarCounter)
	coreNS := rt.NS(rt.NameCoreNS)
	coreNS.Def(varName, f)

	// Test call-head-var: extract the symbol from the :load-var
	// The aux from buildLispIRWith will be a qualified Var, which (symbol aux) converts to 'lltest*/callee.
	// We want just the unqualified name 'callee'.
	expr := fmt.Sprintf(`(pr-str (ir.passes.inline/call-head-var %s (ir.passes.inline/first-call-nid %s)))`, varName, varName)
	consts := vm.NewConsts()
	c := compiler.NewCompiler(consts, coreNS)
	c.SetSource("inline-test")
	_, headVarResult, err := c.CompileMultiple(strings.NewReader(expr))
	if err != nil {
		t.Fatalf("eval call-head-var: %v", err)
	}
	headVar := strings.TrimSpace(string(headVarResult.(vm.String)))
	// call-head-var returns a qualified symbol 'lltest*/callee because the callee is in a test namespace.
	// We just need to verify it contains 'callee'.
	if !strings.Contains(headVar, "callee") {
		t.Errorf("call-head-var should return a symbol containing 'callee', got %s", headVar)
	}

	// Test resolve-callee with empty registry
	emptyRegistryExpr := fmt.Sprintf(`(pr-str (ir.passes.inline/resolve-callee %s (ir.passes.inline/first-call-nid %s) {}))`, varName, varName)
	consts = vm.NewConsts()
	c = compiler.NewCompiler(consts, coreNS)
	c.SetSource("inline-test")
	_, emptyResult, err := c.CompileMultiple(strings.NewReader(emptyRegistryExpr))
	if err != nil {
		t.Fatalf("eval resolve-callee empty: %v", err)
	}
	emptyResultStr := strings.TrimSpace(string(emptyResult.(vm.String)))
	if emptyResultStr != "nil" {
		t.Errorf("resolve-callee with empty registry should return nil, got %s", emptyResultStr)
	}

	// Test resolve-callee with a populated registry
	// The key must match the qualified symbol returned by call-head-var
	calleeF := buildLispIR(t, `(defn callee [p] p)`)
	passVarCounter++
	calleeVarName := fmt.Sprintf("*inline-callee-%d*", passVarCounter)
	coreNS.Def(calleeVarName, calleeF)

	// Build the registry with the qualified symbol as the key
	// We use the call-head-var result as the key
	registryExpr := fmt.Sprintf(`(pr-str (let [key (ir.passes.inline/call-head-var %s (ir.passes.inline/first-call-nid %s))]
	                              (boolean (ir.passes.inline/resolve-callee %s (ir.passes.inline/first-call-nid %s) {key %s}))))`, varName, varName, varName, varName, calleeVarName)
	consts = vm.NewConsts()
	c = compiler.NewCompiler(consts, coreNS)
	c.SetSource("inline-test")
	_, registryResult, err := c.CompileMultiple(strings.NewReader(registryExpr))
	if err != nil {
		t.Fatalf("eval resolve-callee with registry: %v", err)
	}
	registryResultStr := strings.TrimSpace(string(registryResult.(vm.String)))
	if registryResultStr != "true" {
		t.Errorf("resolve-callee with matching registry entry should return non-nil, got %s", registryResultStr)
	}
}

// compileFormToGo compiles a defn form through the Go target pipeline
// (via compile-form* with *target* bound to :go) and returns the concatenated
// Go declarations as a string. Used to test integration of passes like lambda-lift.
func compileFormToGo(t *testing.T, src string) string {
	t.Helper()
	passVarCounter++
	formName := fmt.Sprintf("*compile-form-%d*", passVarCounter)

	coreNS := rt.NS(rt.NameCoreNS)

	// Parse the form
	consts := vm.NewConsts()
	c := compiler.NewCompiler(consts, coreNS)
	c.SetSource("compile-form-to-go-parse")
	expr := fmt.Sprintf(`(quote %s)`, src)
	_, parsed, err := c.CompileMultiple(strings.NewReader(expr))
	if err != nil {
		t.Fatalf("parse form: %v", err)
	}

	// Install as a Var
	coreNS.Def(formName, parsed)

	// Compile with *target* bound to :go
	consts = vm.NewConsts()
	c = compiler.NewCompiler(consts, coreNS)
	c.SetSource("compile-form-to-go")
	compileExpr := fmt.Sprintf(`(binding [ir.passes.pipeline/*target* :go]
                                   (ir.passes.pipeline/compile-form %s))`, formName)
	_, result, err := c.CompileMultiple(strings.NewReader(compileExpr))
	if err != nil {
		t.Fatalf("compile-form with *target* :go: %v", err)
	}

	// Handle the result:
	// - If :status :lowered, extract :decl and render
	// - If :kind :multi-fn-template, flatten and render each decl
	declVarName := fmt.Sprintf("*compile-result-%d*", passVarCounter)
	coreNS.Def(declVarName, result)

	// Render each decl via gogen/render, handling both single and multi results
	renderExpr := fmt.Sprintf(`
(let [decls (if (= :multi-fn-template (:kind %s))
              (mapv (fn* [entry] (:fn entry)) (:fns %s))
              (if (= :lowered (:status %s))
                [(:decl %s)]
                []))]
  (apply str (mapv (fn* [d] (gogen/render d)) decls)))`, declVarName, declVarName, declVarName, declVarName)

	consts = vm.NewConsts()
	c = compiler.NewCompiler(consts, coreNS)
	c.SetSource("render-decls")
	_, rendered, err := c.CompileMultiple(strings.NewReader(renderExpr))
	if err != nil {
		t.Fatalf("render decls: %v", err)
	}

	if s, ok := rendered.(vm.String); ok {
		return string(s)
	}
	t.Fatalf("expected gogen/render to return string, got %T", rendered)
	return ""
}

func TestInlineSeesThroughNameStar(t *testing.T) {
	ensureLoader()
	dump := inlineWithRegistry(t,
		`(defn caller [p] (callthru p (name* "r" tiny "r")))`,
		map[string]string{
			"tiny":     `(defn tiny [x] (+ x 1))`,
			"callthru": `(defn callthru [p f] (f p))`, // a 'call'-like applier
			"name*":    `(defn name* [a b c] b)`,      // stub: returns middle arg
		})
	if strings.Contains(dump, "Invoke") || strings.Contains(dump, "Call") {
		t.Fatalf("name*-wrapped callee not seen through:\n%s", dump)
	}
}
