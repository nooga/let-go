/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package ir_test

import (
	"fmt"
	"testing"

	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

// The fold-over-rest unroll fires on the runtime bytecode path when the
// per-defn compile pipeline seeds the inline registry from previously
// compiled same-ns defns (SCENARIO-0027). The flags stay opt-in: tests set
// the var roots around the compiles and restore them.
//
// Note: `binding` around a defn cannot enable the IR path — the defn macro
// consults *ir-compile* at macroexpansion time, before any runtime binding
// takes effect — so the tests flip the var ROOTS and restore them.

func setVarRoot(t *testing.T, nsName, varName string, val vm.Value) vm.Value {
	t.Helper()
	ns := rt.NS(nsName)
	if ns == nil {
		t.Fatalf("namespace %s not found", nsName)
	}
	v, ok := ns.Lookup(vm.Symbol(varName)).(*vm.Var)
	if !ok || v == nil {
		t.Fatalf("var %s/%s not found", nsName, varName)
	}
	old := v.Deref()
	v.SetRoot(val)
	return old
}

// withIRInline runs body with *ir-compile* + *enable-inline* var roots set,
// restoring the previous roots afterward.
func withIRInline(t *testing.T, body func()) {
	t.Helper()
	oldIR := setVarRoot(t, "clojure.core", "*ir-compile*", vm.TRUE)
	oldInline := setVarRoot(t, "ir.passes.inline", "*enable-inline*", vm.TRUE)
	defer func() {
		setVarRoot(t, "clojure.core", "*ir-compile*", oldIR)
		setVarRoot(t, "ir.passes.inline", "*enable-inline*", oldInline)
	}()
	body()
}

func defRules(t *testing.T, nsName, prefix, body string, n int) string {
	t.Helper()
	args := ""
	for i := 0; i < n; i++ {
		runInNs(t, nsName, fmt.Sprintf("(defn %s%d [x] %s)", prefix, i, body))
		args += fmt.Sprintf(" %s%d", prefix, i)
	}
	return args
}

const loopAnySrc = `(defn loopc [input & fs] (loop [v fs] (if (empty? v) false (if ((first v) input) true (recur (rest v))))))`
const loopAllSrc = `(defn allc [input & fs] (loop [v fs] (if (empty? v) true (if ((first v) input) (recur (rest v)) false))))`

// TestUnrollBytecodePipelineFires: a caller of a fold-over-rest combinator,
// compiled through the runtime IR pipeline with inlining enabled, executes
// the unrolled op profile and returns results identical to the loop build.
func TestUnrollBytecodePipelineFires(t *testing.T) {
	ensureLoader()
	const N = 8
	nsName := "bcunroll"

	// Rules: first one matches x=1, the rest never match.
	runInNs(t, nsName, "(defn rt0 [x] (= x 1))")
	args := " rt0" + defRules(t, nsName, "rf", "false", N-1)

	// Loop-combinator baseline: everything compiled with the flags OFF.
	runInNs(t, nsName, loopAnySrc)
	runInNs(t, nsName, fmt.Sprintf("(defn call-loop [x] (loopc x%s))", args))

	// The real pipeline: recompile the combinator and define the caller with
	// *ir-compile* + *enable-inline* on. Compiling loopc under the flags
	// stashes its IR in the runtime cache; compiling call-unroll then seeds
	// the registry from it and the fold recognizer unrolls the call.
	withIRInline(t, func() {
		runInNs(t, nsName, loopAnySrc)
		runInNs(t, nsName, fmt.Sprintf("(defn call-unroll [x] (loopc x%s))", args))
	})

	// Result parity: early match (x=1 short-circuits at rule 0) and full scan.
	for _, x := range []int{1, 2} {
		loopRes := runInNs(t, nsName, fmt.Sprintf("(call-loop %d)", x))
		unrollRes := runInNs(t, nsName, fmt.Sprintf("(call-unroll %d)", x))
		if loopRes != unrollRes {
			t.Fatalf("parity FAIL x=%d: loop=%v unroll=%v", x, loopRes, unrollRes)
		}
		t.Logf("parity OK x=%d: %v", x, loopRes)
	}

	// Hand-unrolled nested-if baseline, compiled with the flags off: the
	// build-stable yardstick for the specialization's op profile.
	runInNs(t, nsName,
		"(defn call-hand [x] (if (rt0 x) true (if (rf0 x) true (if (rf1 x) true (if (rf2 x) true (if (rf3 x) true (if (rf4 x) true (if (rf5 x) true (if (rf6 x) true false)))))))))")
	assertUnrolled(t, nsName, "loopc$any$8", 8, "(call-loop 2)", "(call-unroll 2)", "(call-hand 2)")
}

// TestUnrollBytecodePipelineAll: same proof for the :all-shaped combinator.
func TestUnrollBytecodePipelineAll(t *testing.T) {
	ensureLoader()
	const N = 8
	nsName := "bcunrollall"

	// Rules all true for x!=0; rule 0 false for x=99 exercises early exit.
	// The filler prefix must differ from the special rule's name — a shared
	// "at" prefix would make defRules emit its own at0 and clobber the
	// early-exit rule.
	runInNs(t, nsName, "(defn at0 [x] (not= x 99))")
	args := " at0" + defRules(t, nsName, "af", "(not= x 0)", N-1)

	runInNs(t, nsName, loopAllSrc)
	runInNs(t, nsName, fmt.Sprintf("(defn call-loop [x] (allc x%s))", args))

	withIRInline(t, func() {
		runInNs(t, nsName, loopAllSrc)
		runInNs(t, nsName, fmt.Sprintf("(defn call-unroll [x] (allc x%s))", args))
	})

	for _, x := range []int{1, 99} {
		loopRes := runInNs(t, nsName, fmt.Sprintf("(call-loop %d)", x))
		unrollRes := runInNs(t, nsName, fmt.Sprintf("(call-unroll %d)", x))
		if loopRes != unrollRes {
			t.Fatalf("parity FAIL x=%d: loop=%v unroll=%v", x, loopRes, unrollRes)
		}
	}

	runInNs(t, nsName,
		"(defn call-hand [x] (if (at0 x) (if (af0 x) (if (af1 x) (if (af2 x) (if (af3 x) (if (af4 x) (if (af5 x) (if (af6 x) true false) false) false) false) false) false) false) false))")
	assertUnrolled(t, nsName, "allc$all$8", 8, "(call-loop 1)", "(call-unroll 1)", "(call-hand 1)")
}

// assertUnrolled verifies the outlined specialization fired and that its op
// profile is in hand-written territory. The former assertion compared the
// unrolled caller against the loop combinator at a fixed 2x ratio, but the
// loop baseline's VM-op cost is build-dependent: under -tags gogen_ir the
// per-element seq walk (empty?/first/rest) dispatches to native Go and
// executes zero VM ops, shrinking the baseline — and it keeps shrinking as
// more core fns go native. The unrolled path's cost IS build-stable (it
// executes only user bytecode and native builtins on both builds), so assert
// (a) the callee var was rewritten to the interned specialization, (b) the
// unrolled caller strictly beats the loop build, and (c) it stays within an
// additive budget of a hand-unrolled nested-if baseline compiled with the
// flags off. The budget models what outlining adds over hand-written code —
// one extra fixed-arity call hop carrying the n+1 flat operands (measured
// ~37 ops at n=8) — so it is 6*(n+1); a specialization that kept any
// per-element seq walk would blow past it (the walk alone costs ~10 ops per
// element).
func assertUnrolled(t *testing.T, nsName, specName string, n int, loopProg, unrollProg, handProg string) {
	t.Helper()
	if v := runInNs(t, nsName, fmt.Sprintf("(some? (resolve (quote %s)))", specName)); v != vm.TRUE {
		t.Fatalf("specialization %s was not interned — unroll did not fire", specName)
	}
	opsLoop := profiledOps(t, nsName, loopProg)
	opsUnroll := profiledOps(t, nsName, unrollProg)
	opsHand := profiledOps(t, nsName, handProg)
	budget := opsHand + uint64(6*(n+1))
	t.Logf("full-scan ops: loop=%d unroll=%d hand=%d budget=%d", opsLoop, opsUnroll, opsHand, budget)
	if opsUnroll >= opsLoop {
		t.Fatalf("unrolled caller is no cheaper than the loop: ops(unroll)=%d vs ops(loop)=%d", opsUnroll, opsLoop)
	}
	if opsUnroll > budget {
		t.Fatalf("unrolled caller not in hand-written territory: ops(unroll)=%d > ops(hand)=%d + outlining budget %d", opsUnroll, opsHand, 6*(n+1))
	}
}

// TestUnrollBytecodePipelineOverCap: a call with more operands than
// *max-unroll* stays the loop (no unrolled profile).
func TestUnrollBytecodePipelineOverCap(t *testing.T) {
	ensureLoader()
	const N = 40 // over the *max-unroll* cap of 32
	nsName := "bcunrollcap"

	args := defRules(t, nsName, "rc", "false", N)
	runInNs(t, nsName, loopAnySrc)
	runInNs(t, nsName, fmt.Sprintf("(defn call-loop [x] (loopc x%s))", args))

	withIRInline(t, func() {
		runInNs(t, nsName, loopAnySrc)
		runInNs(t, nsName, fmt.Sprintf("(defn call-capped [x] (loopc x%s))", args))
	})

	if res := runInNs(t, nsName, "(call-capped 2)"); res != vm.FALSE {
		t.Fatalf("over-cap call wrong result: %v", res)
	}
	opsLoop := profiledOps(t, nsName, "(call-loop 2)")
	opsCapped := profiledOps(t, nsName, "(call-capped 2)")
	t.Logf("over-cap ops: loop=%d capped=%d", opsLoop, opsCapped)
	if opsCapped*5 < opsLoop*4 { // capped must NOT be meaningfully cheaper (>20% drop = unrolled)
		t.Fatalf("over-cap call appears unrolled: ops(capped)=%d vs ops(loop)=%d", opsCapped, opsLoop)
	}
}

// TestUnrollBytecodePipelineNameCollision: a pre-existing user var that
// happens to carry the deterministic specialization name is neither reused
// nor clobbered — the call stays a loop and returns correct results.
func TestUnrollBytecodePipelineNameCollision(t *testing.T) {
	ensureLoader()
	const N = 3
	nsName := "bcunrollclash"

	runInNs(t, nsName, "(defn ct0 [x] (= x 1))")
	args := " ct0" + defRules(t, nsName, "cf", "false", N-1)
	// Hijack the deterministic name with a wrong-arity user fn.
	runInNs(t, nsName, "(defn loopc$any$3 [a b c d e f g] :hijacked)")
	runInNs(t, nsName, loopAnySrc)

	withIRInline(t, func() {
		runInNs(t, nsName, loopAnySrc)
		runInNs(t, nsName, fmt.Sprintf("(defn call-clash [x] (loopc x%s))", args))
	})

	for x, want := range map[int]vm.Value{1: vm.TRUE, 2: vm.FALSE} {
		if got := runInNs(t, nsName, fmt.Sprintf("(call-clash %d)", x)); got != want {
			t.Fatalf("collision case x=%d: got %v want %v", x, got, want)
		}
	}
	// The user's var must be untouched.
	if got := runInNs(t, nsName, "(loopc$any$3 1 2 3 4 5 6 7)"); got != vm.Keyword("hijacked") {
		t.Fatalf("user var clobbered: %v", got)
	}
}
