/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package ir_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

// ITER-0034 / STORY-0035 (EPIC-013, Component 2 of the multiblock-inliner +
// loop-unroll design). Task 2a: recognize the fold-over-rest idiom in IR.
//
// The idiom (fixtures below) is a variadic combinator whose body loops over its
// rest-arg seq, consuming (first v) / (rest v) and terminating on (empty? v):
//
// The element is the callee applied to a fixed input arg (the yamlstar
// combinator shape), so after specialization each (first v) becomes a known
// rule operand -> a direct call:
//   any: (loop [v fs] (if (empty? v) false (if ((first v) input) true (recur (rest v)))))
//   all: (loop [v fs] (if (empty? v) true  (if ((first v) input) (recur (rest v)) false)))

func foldDesc(t *testing.T, src string, field string) string {
	t.Helper()
	f := buildLispIR(t, src)
	return strings.TrimSpace(lispEval(t, "("+field+" (ir.passes.inline/fold-over-rest? %s))", f))
}

func foldRecognized(t *testing.T, src string) string {
	t.Helper()
	f := buildLispIR(t, src)
	return strings.TrimSpace(lispEval(t, `(some? (ir.passes.inline/fold-over-rest? %s))`, f))
}

const anycSrc = `(defn anyc [input & fs] (loop [v fs] (if (empty? v) false (if ((first v) input) true (recur (rest v))))))`
const allcSrc = `(defn allc [input & fs] (loop [v fs] (if (empty? v) true (if ((first v) input) (recur (rest v)) false))))`

func TestFoldOverRestRecognizesAny(t *testing.T) {
	ensureLoader()
	if got := foldRecognized(t, anycSrc); got != "true" {
		t.Fatalf("anyc should be recognized as fold-over-rest, got %s", got)
	}
	if got := foldDesc(t, anycSrc, ":kind"); got != ":any" {
		t.Fatalf("anyc kind should be :any, got %s", got)
	}
	// rest-index is the last load-arg (arity-1); for [p & xs] arity=2 -> 1.
	if got := foldDesc(t, anycSrc, ":rest-index"); got != "1" {
		t.Fatalf("anyc rest-index should be 1, got %s", got)
	}
}

func TestFoldOverRestRecognizesAll(t *testing.T) {
	ensureLoader()
	if got := foldRecognized(t, allcSrc); got != "true" {
		t.Fatalf("allc should be recognized, got %s", got)
	}
	if got := foldDesc(t, allcSrc, ":kind"); got != ":all" {
		t.Fatalf("allc kind should be :all, got %s", got)
	}
}

func TestFoldOverRestRejectsCounterLoop(t *testing.T) {
	ensureLoader()
	// A counter loop over a fixed arg — NOT a fold over the rest seq.
	src := `(defn cu [n] (loop [i 0] (if (< i n) (recur (inc i)) i)))`
	if got := foldRecognized(t, src); got != "false" {
		t.Fatalf("counter loop must NOT be recognized as fold-over-rest, got %s", got)
	}
}

func TestFoldOverRestRejectsNonVariadic(t *testing.T) {
	ensureLoader()
	// Non-variadic, no rest arg at all.
	if got := foldRecognized(t, `(defn tiny [x] (+ x 1))`); got != "false" {
		t.Fatalf("non-variadic fn must NOT be recognized, got %s", got)
	}
}

func TestFoldOverRestRejectsVariadicNonFold(t *testing.T) {
	ensureLoader()
	// Variadic but does not loop over the rest seq (just returns it).
	if got := foldRecognized(t, `(defn vf [p & xs] xs)`); got != "false" {
		t.Fatalf("variadic non-fold fn must NOT be recognized, got %s", got)
	}
}

// ── T2: specialize a recognized fold-over-rest to a fixed-arity unrolled fn ──

func specializeFold(t *testing.T, comboSrc string, n int) vm.Value {
	t.Helper()
	f := buildLispIR(t, comboSrc)
	return lispEvalReturn(t, f,
		fmt.Sprintf(`(ir.passes.inline/specialize-fold (ir.passes.inline/fold-over-rest? %%s) %d)`, n))
}

// assertValidatesVal runs ir.validate/validate-fn! on a Function value; a
// throw (malformed IR) surfaces as an eval error and fails the test.
func assertValidatesVal(t *testing.T, f vm.Value, label string) {
	t.Helper()
	passVarCounter++
	varName := fmt.Sprintf("*vf-%d*", passVarCounter)
	coreNS := rt.NS(rt.NameCoreNS)
	coreNS.Def(varName, f)
	consts := vm.NewConsts()
	c := compiler.NewCompiler(consts, coreNS)
	c.SetSource("validate-val")
	expr := fmt.Sprintf(`(ir.validate/validate-fn! %s "%s")`, varName, label)
	if _, _, err := c.CompileMultiple(strings.NewReader(expr)); err != nil {
		t.Fatalf("specialized IR failed validation (%s): %v", label, err)
	}
}

func TestSpecializeFoldAnyUnrolls(t *testing.T) {
	ensureLoader()
	spec := specializeFold(t, anycSrc, 3)
	// arity = num-fixed(1: input) + N(3 elements) = 4; non-variadic.
	if got := strings.TrimSpace(lispEvalOn(t, spec, "(ir/fn-arity %s)")); got != "4" {
		t.Fatalf("specialized any arity should be 4, got %s", got)
	}
	if got := strings.TrimSpace(lispEvalOn(t, spec, "(ir/fn-variadic? %s)")); got != "false" {
		t.Fatalf("specialized fn should be non-variadic, got %s", got)
	}
	dump := lispDump(t, spec)
	// Unrolled any(N=3) = 3 short-circuit branch-ifs, loop-free.
	if n := strings.Count(dump, "BranchIf"); n != 3 {
		t.Fatalf("expected 3 BranchIf (unrolled chain), got %d:\n%s", n, dump)
	}
	assertValidatesVal(t, spec, "specialize-any")
}

func TestSpecializeFoldAllUnrolls(t *testing.T) {
	ensureLoader()
	spec := specializeFold(t, allcSrc, 2)
	if got := strings.TrimSpace(lispEvalOn(t, spec, "(ir/fn-arity %s)")); got != "3" {
		t.Fatalf("specialized all arity should be 3, got %s", got)
	}
	dump := lispDump(t, spec)
	if n := strings.Count(dump, "BranchIf"); n != 2 {
		t.Fatalf("expected 2 BranchIf (unrolled chain), got %d:\n%s", n, dump)
	}
	assertValidatesVal(t, spec, "specialize-all")
}

// ── T5: wire specialize+unroll into inline-pass (+ *max-unroll* cap) ──

// inlineFoldDump builds a caller + stubs in a fresh ns, runs inline-with under
// an optional *max-unroll* binding, and returns the dump. cap<=0 means default.
func inlineFoldDump(t *testing.T, cap int, callerSrc string, stubs map[string]string) string {
	t.Helper()
	passVarCounter++
	nsName := fmt.Sprintf("foldwire%d", passVarCounter)
	ns := rt.NS(nsName)
	regEntries := make([]string, 0, len(stubs))
	build := func(src string) vm.Value {
		consts := vm.NewConsts()
		c := compiler.NewCompiler(consts, ns)
		c.SetSource("foldwire-build")
		_, res, err := c.CompileMultiple(strings.NewReader(fmt.Sprintf("(ir.build/build-fn (quote %s))", src)))
		if err != nil {
			t.Fatalf("build %q: %v", src, err)
		}
		return res
	}
	for name, src := range stubs {
		res := build(src)
		vn := fmt.Sprintf("*stub-%s*", name)
		ns.Def(vn, res)
		ns.Def(name, res)
		regEntries = append(regEntries, fmt.Sprintf("'%s/%s %s", nsName, name, vn))
	}
	caller := build(callerSrc)
	ns.Def("*caller-func*", caller)
	regMap := fmt.Sprintf("(hash-map %s)", strings.Join(regEntries, " "))
	var body string
	if cap > 0 {
		body = fmt.Sprintf(`(binding [ir.passes.inline/*max-unroll* %d] (ir.passes.inline/inline-with *caller-func* %s))`, cap, regMap)
	} else {
		body = fmt.Sprintf(`(ir.passes.inline/inline-with *caller-func* %s)`, regMap)
	}
	consts := vm.NewConsts()
	c := compiler.NewCompiler(consts, ns)
	c.SetSource("foldwire-inline")
	if _, _, err := c.CompileMultiple(strings.NewReader(fmt.Sprintf("(do %s *caller-func*)", body))); err != nil {
		t.Fatalf("inline-with: %v", err)
	}
	return lispDump(t, caller)
}

func TestUnrollWiredIntoInlineAny(t *testing.T) {
	ensureLoader()
	dump := inlineFoldDump(t, 0,
		`(defn rule [input] (anyc input r0 r1 r2))`,
		map[string]string{"anyc": anycSrc, "r0": `(defn r0 [x] x)`, "r1": `(defn r1 [x] x)`, "r2": `(defn r2 [x] x)`})
	// Unrolled any(N=3) => 3 short-circuit branch-ifs in the caller.
	if n := strings.Count(dump, "BranchIf"); n != 3 {
		t.Fatalf("expected 3 BranchIf (unrolled any), got %d:\n%s", n, dump)
	}
	// The loop mechanics must be gone (specialized away, not inlined as a loop).
	for _, loopism := range []string{"empty?", "first", "rest"} {
		if strings.Contains(dump, loopism) {
			t.Fatalf("loop mechanic %q should be gone after unroll:\n%s", loopism, dump)
		}
	}
}

func TestUnrollWiredIntoInlineAll(t *testing.T) {
	ensureLoader()
	dump := inlineFoldDump(t, 0,
		`(defn rule [input] (allc input r0 r1))`,
		map[string]string{"allc": allcSrc, "r0": `(defn r0 [x] x)`, "r1": `(defn r1 [x] x)`})
	if n := strings.Count(dump, "BranchIf"); n != 2 {
		t.Fatalf("expected 2 BranchIf (unrolled all), got %d:\n%s", n, dump)
	}
	for _, loopism := range []string{"empty?", "first", "rest"} {
		if strings.Contains(dump, loopism) {
			t.Fatalf("loop mechanic %q should be gone after unroll:\n%s", loopism, dump)
		}
	}
}

func TestUnrollRespectsMaxCap(t *testing.T) {
	ensureLoader()
	// *max-unroll* = 2, but N = 3 -> over cap -> left as a runtime call, NOT unrolled.
	dump := inlineFoldDump(t, 2,
		`(defn rule [input] (anyc input r0 r1 r2))`,
		map[string]string{"anyc": anycSrc, "r0": `(defn r0 [x] x)`, "r1": `(defn r1 [x] x)`, "r2": `(defn r2 [x] x)`})
	if n := strings.Count(dump, "BranchIf"); n != 0 {
		t.Fatalf("over-cap fold must NOT be unrolled (expected 0 BranchIf), got %d:\n%s", n, dump)
	}
	if !strings.Contains(dump, "Call") {
		t.Fatalf("over-cap fold should remain a runtime call:\n%s", dump)
	}
}

// ── E1: lowered-Go cross-surface proof (AC-FR.1 integration / AC-PO-FR.1) ──

// sliceFunc returns the body of `func <name>(` .. up to the next top-level
// `func ` in rendered Go (or end), for per-function assertions.
func sliceFunc(src, name string) string {
	start := strings.Index(src, "func "+name+"(")
	if start < 0 {
		return ""
	}
	rest := src[start+len("func "+name+"("):]
	if i := strings.Index(rest, "\nfunc "); i >= 0 {
		return rest[:i]
	}
	return rest
}

// lowerFoldNs lowers a namespace defining recursive rules r0..r(n-1), the anyc
// combinator, and a `parse` rule that applies anyc to those rules, with inline
// enabled (unroll on). Returns the rendered Go.
func lowerFoldNs(t *testing.T, n int) string {
	t.Helper()
	return lowerFoldNsMode(t, n, true)
}

// lowerFoldNsMode is lowerFoldNs with an explicit unroll toggle (*enable-inline*).
func lowerFoldNsMode(t *testing.T, n int, unroll bool) string {
	t.Helper()
	var rules, forms strings.Builder
	for i := 0; i < n; i++ {
		fmt.Fprintf(&rules, `(eval (quote (defn r%d [x] (if (< x 0) x (r%d (- x 1))))))`+"\n", i, i)
		fmt.Fprintf(&forms, `(quote (defn r%d [x] (if (< x 0) x (r%d (- x 1)))))`+"\n", i, i)
	}
	rulesArgs := ""
	for i := 0; i < n; i++ {
		rulesArgs += fmt.Sprintf(" r%d", i)
	}
	v := runLispExpr(t, fmt.Sprintf(`(do (create-ns (quote ysx))
	     (binding [*ns* (the-ns (quote ysx))]
	       %s
	       (eval (quote (defn anyc [input & fs] (loop [v fs] (if (empty? v) false (if ((first v) input) true (recur (rest v))))))))
	       (eval (quote (defn parse [input] (anyc input%s)))))
	     (binding [ir.passes.inline/*enable-inline* %t]
	       (ir.passes.pipeline/lower-ns-to-go "ysx" (quote ysx)
	         [%s
	          (quote (defn anyc [input & fs] (loop [v fs] (if (empty? v) false (if ((first v) input) true (recur (rest v)))))))
	          (quote (defn parse [input] (anyc input%s)))])))`,
		rules.String(), rulesArgs, unroll, forms.String(), rulesArgs))
	s, ok := v.Unbox().(string)
	if !ok {
		t.Fatalf("expected rendered Go string, got %T", v.Unbox())
	}
	return s
}

func TestUnrollLoweredGoEliminatesLoop(t *testing.T) {
	ensureLoader()
	src := lowerFoldNs(t, 3)
	parse := sliceFunc(src, "Parse")
	if parse == "" {
		t.Fatalf("Parse did not lower:\n%s", src)
	}
	// AC-FR.1 / AC-PO-FR.1 (achieved): in the caller the fold loop is unrolled
	// away — no seq iteration (first/rest/empty?) and no per-node closure alloc.
	for _, loopism := range []string{"corefns.First", "corefns.Rest", "empty_QMARK", "BoxNativeFn"} {
		if strings.Contains(parse, loopism) {
			t.Fatalf("unrolled Parse must not contain %q (loop/closure-alloc eliminated):\n%s", loopism, parse)
		}
	}
	// The unrolled short-circuit chain: N truthiness tests over the N elements.
	if got := strings.Count(parse, "IsTruthy"); got != 3 {
		t.Fatalf("expected 3 short-circuit tests in unrolled Parse, got %d:\n%s", got, parse)
	}
	// KNOWN GAP (follow-up): element applications still lower to rt.InvokeValueEC
	// rather than direct RN(...) calls, because the multi-block tail splice routes
	// operands through block-args and validate.lg forbids threading a :load-var
	// cross-block — so the direct-call devirt (which keys on load-var heads) can't
	// fire. Full "direct rule-to-rule calls" (design §5.3) is deferred to a
	// block-arg-aware devirt follow-up; the profiler (ITER-0035) quantifies whether
	// the residual dispatch matters vs the loop version. This assertion documents
	// current reality so a future devirt fix flips it deliberately.
	if got := strings.Count(parse, "InvokeValueEC"); got != 3 {
		t.Fatalf("expected 3 element dispatches (devirt gap documented), got %d:\n%s", got, parse)
	}
}
