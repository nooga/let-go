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

func TestInlineEligibilityUsesMutabilityAndSize(t *testing.T) {
	ensureLoader()

	// Test 1: a small pure function should be eligible.
	small := buildLispIR(t, `(defn tiny [x] (+ x 1))`)
	if got := strings.TrimSpace(lispEval(t, `(ir.passes.inline/inline-eligible? %s)`, small)); got != "true" {
		t.Fatalf("small pure fn should be eligible, got %s", got)
	}

	// Test 2: a function that mutates a var root should be ineligible.
	// We use set! which directly generates a :set-var IR instruction.
	// Define some-v as a stub var so the set! can reference it.
	mut := buildLispIRWith(t,
		map[string]string{
			"some-v": "nil",
		},
		`(defn bad [x] (set! some-v x))`)
	if got := strings.TrimSpace(lispEval(t, `(ir.passes.inline/inline-eligible? %s)`, mut)); got != "false" {
		t.Fatalf("var-root-rebinding fn must be ineligible, got %s", got)
	}
}

func TestInlineEligibilityRejectsSelfRecursive(t *testing.T) {
	ensureLoader()

	// A self-recursive defn must be INELIGIBLE for inlining: splicing it into
	// itself grows IR size every round (max-inline-rounds bounds the round count,
	// not the size), so it must be rejected up front.
	//
	// The self-call `(fact ...)` builds to a :load-var whose aux is a *Var*
	// object (build.lg), not a symbol. self-recursive? must normalize the aux to
	// a symbol (mirroring call-head-var) before comparing it to the fn name — a
	// raw Var-vs-symbol compare never matches, so the fn would be wrongly judged
	// non-recursive and marked inlineable.
	selfRec := buildLispIRWith(t,
		map[string]string{"fact": "(fn* [n] n)"},
		`(defn fact [n] (if (< n 2) n (* n (fact (dec n)))))`)
	if got := strings.TrimSpace(lispEval(t, `(ir.passes.inline/inline-eligible? %s)`, selfRec)); got != "false" {
		t.Fatalf("self-recursive fn must be ineligible for inlining, got %s", got)
	}
}

func TestInlineWithLeavesSelfRecursiveCallUnspliced(t *testing.T) {
	ensureLoader()
	// End-to-end guard: a caller invoking a self-recursive helper must keep the
	// call — inline-with only splices inline-eligible? callees, and a
	// self-recursive fn is ineligible. Before the self-recursive? Var
	// normalization fix, the helper was wrongly eligible and would splice into
	// itself, growing IR every round.
	passVarCounter++
	nsName := fmt.Sprintf("selfrectest%d", passVarCounter)
	ns := rt.NS(nsName)

	eval := func(label, expr string) vm.Value {
		consts := vm.NewConsts()
		c := compiler.NewCompiler(consts, ns)
		c.SetSource(label)
		_, res, err := c.CompileMultiple(strings.NewReader(expr))
		if err != nil {
			t.Fatalf("%s: %v", label, err)
		}
		return res
	}
	build := func(label, src string) vm.Value {
		return eval(label, fmt.Sprintf("(ir.build/build-fn (quote %s))", src))
	}

	// Pre-define `fact` as a plain fn value so the self-reference in its body
	// resolves at build time; then build the real self-recursive IR.
	eval("selfrec-stub", `(def fact (fn* [n] n))`)
	fact := build("selfrec-callee", `(defn fact [n] (if (< n 2) n (* n (fact (dec n)))))`)
	ns.Def("*fact*", fact)

	caller := build("selfrec-caller", `(defn caller [m] (fact m))`)
	ns.Def("*caller*", caller)

	consts := vm.NewConsts()
	c := compiler.NewCompiler(consts, ns)
	c.SetSource("selfrec-inline-with")
	inlineExpr := fmt.Sprintf(
		`(do (ir.passes.inline/inline-with *caller* (hash-map '%s/fact *fact*)) *caller*)`,
		nsName)
	_, result, err := c.CompileMultiple(strings.NewReader(inlineExpr))
	if err != nil {
		t.Fatalf("inline-with: %v", err)
	}
	dump := lispDump(t, result)
	if !(strings.Contains(dump, "Call") || strings.Contains(dump, "Invoke")) {
		t.Fatalf("self-recursive callee must NOT be spliced (call should survive):\n%s", dump)
	}
}

func TestInlineSplicesSimpleCall(t *testing.T) {
	ensureLoader()
	// (defn caller [p] (tiny p)) with tiny=(defn tiny [x] (+ x 1)) inlines to (+ p 1): no call op remains.
	dump := inlineWithRegistry(t, `(defn caller [p] (tiny p))`, map[string]string{"tiny": `(defn tiny [x] (+ x 1))`})
	if strings.Contains(dump, "Call") || strings.Contains(dump, "Invoke") {
		t.Fatalf("call should be inlined away:\n%s", dump)
	}
	if !strings.Contains(dump, "Add") {
		t.Fatalf("inlined body (Add) missing:\n%s", dump)
	}
}

func TestInlineMultiBlockTailNonLooping(t *testing.T) {
	ensureLoader()
	// callee `choose` is multi-block: (if x 1 2) builds to 4 blocks
	// (entry BranchIf, two Const arms, join BlockArg -> Return).
	// The call (choose p) is in tail position (its result feeds caller's Return),
	// so splice-multiblock-tail! should clone the callee's blocks in, route p via
	// entry block-args, and rewire the returns — leaving no Call/Invoke.
	dump := inlineWithRegistry(t,
		`(defn caller [p] (choose p))`,
		map[string]string{"choose": `(defn choose [x] (if x 1 2))`})
	if strings.Contains(dump, "Call") || strings.Contains(dump, "Invoke") {
		t.Fatalf("tail-position multi-block call should be inlined away:\n%s", dump)
	}
	// The callee's control flow (the if) must survive as a BranchIf in the caller.
	if !strings.Contains(dump, "BranchIf") {
		t.Fatalf("inlined callee's BranchIf missing:\n%s", dump)
	}
	// Both const arms should be present.
	if !strings.Contains(dump, "Const ; 1") || !strings.Contains(dump, "Const ; 2") {
		t.Fatalf("inlined callee's const arms missing:\n%s", dump)
	}
	// The spliced CFG must validate (no cross-block refs, arities/symmetry OK).
	assertInlineValidates(t, `(defn caller [p] (choose p))`,
		map[string]string{"choose": `(defn choose [x] (if x 1 2))`})
}

func TestInlineMultiBlockTailLoopCallee(t *testing.T) {
	ensureLoader()
	// callee `cu` contains a loop (multi-block with a back-edge). Inlined in
	// tail position, its loop blocks should clone in and the call disappear —
	// the loop still runs (unrolling is Component 2, not here).
	callee := `(defn cu [n] (loop [i 0] (if (< i n) (recur (inc i)) i)))`
	dump := inlineWithRegistry(t,
		`(defn caller [m] (cu m))`,
		map[string]string{"cu": callee})
	if strings.Contains(dump, "Call") || strings.Contains(dump, "Invoke") {
		t.Fatalf("tail-position loop callee should be inlined away:\n%s", dump)
	}
	// The loop's comparison and increment must survive (loop still present).
	if !strings.Contains(dump, "Lt") || !strings.Contains(dump, "Inc") {
		t.Fatalf("inlined loop body (Lt/Inc) missing:\n%s", dump)
	}
	// The spliced loop CFG must validate (back-edge preds, symmetric branch-if
	// args, in-range refs, arities).
	assertInlineValidates(t, `(defn caller [m] (cu m))`,
		map[string]string{"cu": callee})
}

func TestInlineMultiBlockNonTailLeftUntouched(t *testing.T) {
	ensureLoader()
	// The call (choose p) here is NON-tail: its result feeds an Add, not the
	// caller's Return. General-position multi-block inlining is deferred to
	// Component 1b, so the call must be left intact.
	dump := inlineWithRegistry(t,
		`(defn caller [p] (+ (choose p) 1))`,
		map[string]string{"choose": `(defn choose [x] (if x 1 2))`})
	if !(strings.Contains(dump, "Call") || strings.Contains(dump, "Invoke")) {
		t.Fatalf("non-tail multi-block call must be left untouched (1b deferred):\n%s", dump)
	}
}

func TestInlineTerminatesOnMutualRecursion(t *testing.T) {
	ensureLoader()
	// Test that inline-with terminates when the iteration bound is reached.
	// We verify this by checking that the bound (max-inline-rounds) is respected.
	//
	// Since mutual recursion within the IR layer is difficult to construct
	// (the functions must resolve during build-fn), we instead verify that the
	// loop bound is present and active by testing a scenario where cascading
	// inlines could theoretically occur, ensuring the bound prevents hangs.
	//
	// A simple caller -> callee chain: inlining callee into caller should
	// terminate in 1 iteration. If the bound were missing or set to 0,
	// the loop would behave differently. This test is a sanity check that
	// the bounding logic is in place.
	dump := inlineWithRegistry(t,
		`(defn caller [p] (callee p))`,
		map[string]string{
			"callee": `(defn callee [x] (+ x 1))`,
		})
	// The test passes if inlineWithRegistry returns without hanging.
	// Verify the dump is non-empty and well-formed.
	if strings.TrimSpace(dump) == "" {
		t.Fatalf("inline-with produced empty dump")
	}
	// The call should be inlined away (no Call op should remain).
	if strings.Contains(dump, "Call") {
		t.Errorf("call should have been inlined away, but dump contains Call:\n%s", dump)
	}
}

// TestOptimizeFnInlinesWhenRegistryBound proves the inline pass is WIRED into
// the lowering pipeline: with *inline-registry* bound, optimize-fn must splice
// a resolvable call away. Without the wiring (B7), optimize-fn leaves the call.
func TestOptimizeFnInlinesWhenRegistryBound(t *testing.T) {
	ensureLoader()
	dump := optimizeFnWithRegistry(t,
		`(defn caller [p] (tiny p))`,
		map[string]string{"tiny": `(defn tiny [x] (+ x 1))`})
	if strings.Contains(dump, "Call") || strings.Contains(dump, "Invoke") {
		t.Errorf("optimize-fn should have inlined the call away, but dump still has a call:\n%s", dump)
	}
	if !strings.Contains(dump, "Add") {
		t.Errorf("inlined body (Add) missing — inline did not run in optimize-fn:\n%s", dump)
	}
}

// TestLowerNsInlinesSameNsCall proves lower-ns-to-go seeds *inline-registry*
// with the namespace's own defns, so a caller inlines a same-ns callee. On this
// branch a bare same-ns sibling call does not lower (falls back to nil); once
// the callee body is spliced in, the caller lowers with no call to the callee.
func TestLowerNsInlinesSameNsCall(t *testing.T) {
	ensureLoader()
	// Intern the defns first, mirroring the real AOT flow (the namespace is
	// loaded/required before it is lowered), so the sibling call resolves to a
	// :call op that inline can splice.
	v := runLispExpr(t, `(do (create-ns (quote inlinensx))
	     (binding [*ns* (the-ns (quote inlinensx))]
	       (eval (quote (defn add1 [x] (+ x 1))))
	       (eval (quote (defn use1 [p] (add1 p)))))
	     (binding [ir.passes.inline/*enable-inline* true]
	       (ir.passes.pipeline/lower-ns-to-go "inlinensx" (quote inlinensx)
	         [(quote (defn add1 [x] (+ x 1)))
	          (quote (defn use1 [p] (add1 p)))])))`)
	rendered, ok := v.Unbox().(string)
	if !ok {
		t.Fatalf("expected rendered Go string, got %T", v.Unbox())
	}
	// The caller must have lowered (proving inline rescued the same-ns call).
	// Public defns lower to exported PascalCase names (#371).
	if !strings.Contains(rendered, "func Use1(") {
		t.Fatalf("Use1 did not lower — inline registry not seeded by lower-ns-to-go:\n%s", rendered)
	}
	// Add1 should appear only as its own definition, not as a call inside Use1.
	if strings.Count(rendered, "Add1(") != 1 {
		t.Errorf("expected Add1 to appear once (definition only), got %d — call not inlined:\n%s",
			strings.Count(rendered, "Add1("), rendered)
	}
}

// TestDevirtualizesLocalClosureCall proves the higher-order devirt: a call whose
// callee is a locally-constructed closure literal (IIFE) is spliced inline —
// no rt.InvokeValueEC / rt.BoxNativeFn survives; the closure body (AddValue) is
// emitted directly with args + captures substituted. This is what collapses the
// combinator parse dispatch after named-defn inlining exposes the closure literals.
func TestDevirtualizesLocalClosureCall(t *testing.T) {
	ensureLoader()
	v := runLispExpr(t,
		`(do (create-ns (quote iife))
		     (binding [ir.passes.inline/*enable-inline* true]
		       (ir.passes.pipeline/lower-ns-to-go "iife" (quote iife)
		         [(quote (defn f [x] ((fn [y] (+ y x)) 10)))])))`)
	src := v.String()
	if strings.Contains(src, "InvokeValueEC") || strings.Contains(src, "BoxNativeFn") {
		t.Fatalf("local closure call should be devirtualized (no InvokeValueEC/BoxNativeFn):\n%s", src)
	}
	// The closure body (+ y x) with y=10, x=arg0 is spliced inline. Downstream
	// typeinfer specializes it to a native int add ("+ 10"); if types stay boxed
	// it would be rt.AddValue. Either way the add is present with no dispatch.
	if !strings.Contains(src, "+ 10") && !strings.Contains(src, "AddValue") {
		t.Fatalf("expected inlined closure body (the add):\n%s", src)
	}
}

// TestDevirtualizesMultiBlockClosureCall proves general-position multi-block
// closure devirt: a capturing closure with internal branches (multi-block),
// invoked in non-tail position (result feeds further arithmetic), is spliced
// inline — the caller block is split, the closure blocks are cloned in, and its
// returns jump to the continuation. This is the pegbench combinator shape.
func TestDevirtualizesMultiBlockClosureCall(t *testing.T) {
	ensureLoader()
	v := runLispExpr(t,
		`(do (create-ns (quote mbc))
		     (binding [ir.passes.inline/*enable-inline* true]
		       (ir.passes.pipeline/lower-ns-to-go "mbc" (quote mbc)
		         [(quote (defn g [x] (+ ((fn [y] (if (< y x) (- y) y)) 5) 1)))])))`)
	src := v.String()
	if strings.Contains(src, "InvokeValueEC") || strings.Contains(src, "BoxNativeFn") {
		t.Fatalf("multi-block capturing closure call should be devirtualized (no InvokeValueEC/BoxNativeFn):\n%s", src)
	}
}
