/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package ir_test

import (
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

// STORY-0059 — pre-execution cleanup pass (ir.passes.cleanup).
//
// The proven-needed deliverable is dead block-param + edge-arg compaction: DCE
// removes a dead block-arg from a block's inst list and marks it :invalid, but
// leaves it in the block's :params with every in-edge still threading a value
// into it (a per-iteration write nothing reads). `cleanup` drops those params
// AND the correspondingly-positioned edge-args, keeping arity aligned.
//
// Fixture: a loop binding (`unused`) that is re-threaded as a constant on recur
// and never read. optimize-fn's DCE tombstones its block-arg but leaves the
// param + edge-args behind — the exact shape compactDeadParams targeted in the
// ITER-0057 spike.
const deadParamFixture = `(defn f [n] (loop [i 0 unused 7] (if (< i n) (recur (+ i 1) 99) i)))`

// invalidParamCountExpr counts, over an optimized-then-cleaned f, block params
// whose defining inst is a tombstone (:invalid). Must be 0 after cleanup.
const invalidParamCountExpr = `(let [f (ir.passes.cleanup/cleanup (ir.passes.pipeline/optimize-fn %s))]
  (reduce + 0
    (map (fn [bid]
           (count (filter (fn [p] (= :invalid (ir/op p f)))
                          (ir/block-params bid f))))
         (ir/blocks f))))`

// arityAlignedExpr re-derives validate.lg's arity invariant (check-branch-arg-arity-
// for-target!) INDEPENDENTLY: for every terminator edge, the edge's arg count must
// equal the target block's param count. This is the un-gameable core — a pass that
// drops a param without realigning its in-edges makes this FALSE (and would also
// throw in validate.lg). Do not weaken it.
const arityAlignedExpr = `(let [f (ir.passes.cleanup/cleanup (ir.passes.pipeline/optimize-fn %s))]
  (every? identity
    (map (fn [bid]
           (let [term (ir/block-term bid f)]
             (if (nil? term) true
               (let [o (ir/op term f) a (ir/aux term f)]
                 (cond
                   (= o :branch)
                   (= (count (:args a))
                      (count (ir/block-params (:target a) f)))
                   (= o :branch-if)
                   (and (= (count (:args (:true-target a)))
                           (count (ir/block-params (:target (:true-target a)) f)))
                        (= (count (:args (:false-target a)))
                           (count (ir/block-params (:target (:false-target a)) f))))
                   :else true)))))
         (ir/blocks f))))`

// tombstoneRefCountExpr counts, over an optimized-then-cleaned f, refs of LIVE
// body insts (op not :invalid/:block-arg) that point at a tombstone (:invalid).
// Must be 0 — this is where the DCE-sweep-after-compaction earns its keep: an
// inst whose only use was a dropped edge-arg is orphaned; if the sweep missed it,
// a live consumer could still reference a tombstone.
const tombstoneRefCountExpr = `(let [f (ir.passes.cleanup/cleanup (ir.passes.pipeline/optimize-fn %s))
      tomb? (fn [nid] (= :invalid (ir/op nid f)))]
  (reduce + 0
    (map (fn [bid]
           (reduce + 0
             (map (fn [nid]
                    (if (contains? #{:invalid :block-arg} (ir/op nid f))
                      0
                      (count (filter tomb? (ir/refs nid f)))))
                  (ir/block-insts bid f))))
         (ir/blocks f))))`

// TestCleanupCompactsDeadBlockParams (AC-CL.1): cleanup removes every :invalid
// (DCE-tombstoned) block-param and realigns the in-edge args so arity stays valid.
func TestCleanupCompactsDeadBlockParams(t *testing.T) {
	ensureLoader()
	f := buildLispIR(t, deadParamFixture)

	// Precondition sanity: the fixture actually produces dead params pre-cleanup
	// (guards against the fixture silently folding away — otherwise the test would
	// pass vacuously).
	preInvalid := `(let [f (ir.passes.pipeline/optimize-fn %s)]
      (reduce + 0 (map (fn [bid]
                         (count (filter (fn [p] (= :invalid (ir/op p f)))
                                        (ir/block-params bid f))))
                       (ir/blocks f))))`
	if got := lispInt(t, lispEvalReturn(t, buildLispIR(t, deadParamFixture), preInvalid)); got == 0 {
		t.Fatalf("fixture produced no dead params pre-cleanup; test would be vacuous\nIR:\n%s", lispDump(t, f))
	}

	if got := lispInt(t, lispEvalReturn(t, f, invalidParamCountExpr)); got != 0 {
		t.Fatalf("cleanup left %d :invalid block-params\nIR:\n%s", got, lispDump(t, f))
	}
	if lispEvalReturn(t, buildLispIR(t, deadParamFixture), arityAlignedExpr) != vm.TRUE {
		t.Fatalf("cleanup broke edge-arg/param arity\nIR:\n%s", lispDump(t, f))
	}
}

// TestCleanupCorpusInvariants (SCENARIO-0033 / AC-PO-CL.1): over the whole
// pipelineCorpus, cleanup produces a tombstone-free, arity-aligned execution form
// — no :invalid params, every edge's arg count equals its target's param count,
// and no live inst references a tombstone. Also censuses the params stripped. This
// is un-gameable: it re-derives the invariants structurally over real corpus IR,
// so an implementation that hardcodes or fakes one fixture cannot pass it.
func TestCleanupCorpusInvariants(t *testing.T) {
	ensureLoader()
	// optimize-fn and cleanup mutate their input IR in place and are NOT idempotent,
	// so every eval below gets a FRESH build of the source (never a reused value).
	fresh := func(src string) (vm.Value, bool) {
		v, err := tryBuildLispIR(src)
		return v, err == nil
	}
	var paramsBefore, paramsAfter, funcs int
	for _, e := range pipelineCorpus {
		if _, ok := fresh(e.src); !ok {
			continue
		}
		funcs++
		v0, _ := fresh(e.src)
		if got := lispInt(t, lispEvalReturn(t, v0, invalidParamCountExpr)); got != 0 {
			t.Fatalf("%s: cleanup left %d :invalid params", e.name, got)
		}
		v1, _ := fresh(e.src)
		if lispEvalReturn(t, v1, arityAlignedExpr) != vm.TRUE {
			t.Fatalf("%s: cleanup broke edge/param arity\nIR:\n%s", e.name, lispDump(t, v1))
		}
		v2, _ := fresh(e.src)
		if got := lispInt(t, lispEvalReturn(t, v2, tombstoneRefCountExpr)); got != 0 {
			t.Fatalf("%s: %d live-inst refs point at a tombstone after cleanup", e.name, got)
		}
		v3, _ := fresh(e.src)
		before := lispInt(t, lispEvalReturn(t, v3,
			`(let [f (ir.passes.pipeline/optimize-fn %s)]
               (reduce + 0 (map (fn [bid] (count (ir/block-params bid f))) (ir/blocks f))))`))
		v4, _ := fresh(e.src)
		after := lispInt(t, lispEvalReturn(t, v4,
			`(let [f (ir.passes.cleanup/cleanup (ir.passes.pipeline/optimize-fn %s))]
               (reduce + 0 (map (fn [bid] (count (ir/block-params bid f))) (ir/blocks f))))`))
		if after > before {
			t.Fatalf("%s: cleanup INCREASED params %d -> %d", e.name, before, after)
		}
		paramsBefore += before
		paramsAfter += after
	}
	if funcs == 0 {
		t.Skip("corpus produced no functions")
	}
	t.Logf("cleanup census over pipelineCorpus: %d fns; block-params %d -> %d (%d dead params+edge-arg positions stripped)",
		funcs, paramsBefore, paramsAfter, paramsBefore-paramsAfter)
}
