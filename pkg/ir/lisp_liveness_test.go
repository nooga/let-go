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

// tryBuildLispIR builds IR like buildLispIR but returns an error instead of
// failing the test — some corpus entries (e.g. destructuring) need pipeline
// stub context this bare builder does not provide, and are skipped, not failed.
func tryBuildLispIR(src string) (vm.Value, error) {
	consts := vm.NewConsts()
	c := compiler.NewCompiler(consts, rt.NS(rt.NameCoreNS))
	c.SetSource("lisp-liveness-build")
	expr := fmt.Sprintf(`(ir.build/build-fn (quote %s))`, src)
	_, result, err := c.CompileMultiple(strings.NewReader(expr))
	return result, err
}

// The correctness oracle for ir.passes.liveness. The load-bearing check is
// livenessInvariantExpr: an edge-consistency law re-derived from the IR
// structure INDEPENDENTLY of the pass's own dataflow. For every edge B -> tgt:
//   (1) (live-in[tgt] - params[tgt]) must be a subset of live-out[B]
//   (2) for each edge-arg i, if the target param i is live-in to tgt, then the
//       edge-arg value must be in live-out[B]
// A pass that drops edge-args, mishandles block-params, or fails to propagate
// a cross-block use produces a non-zero violation count. This is deliberately
// not the pass's own logic — it is the ground truth the pass must satisfy.

const livenessInvariantExpr = `
(let [f %s
      live (ir.passes.liveness/compute-liveness f)]
  (reduce
    (fn [violations bid]
      (let [lo (:live-out (get live bid))
            term (ir/block-term bid f)
            targets (if (nil? term) []
                      (let [op (ir/op term f) aux (ir/aux term f)]
                        (cond
                          (= op :branch) [aux]
                          (= op :branch-if) [(ir/cond-target-true aux) (ir/cond-target-false aux)]
                          :else [])))]
        (reduce
          (fn [vs bt]
            (let [tgt (ir/branch-target-target bt)
                  args (vec (ir/branch-target-args bt))
                  params (vec (ir/block-params tgt f))
                  pset (set params)
                  tgt-in (:live-in (get live tgt))
                  leak (reduce (fn [n v]
                                 (if (and (not (contains? pset v)) (not (contains? lo v)))
                                   (inc n) n))
                               0 tgt-in)
                  earg (reduce (fn [n i]
                                 (if (and (< i (count args))
                                          (contains? tgt-in (nth params i))
                                          (not (contains? lo (nth args i))))
                                   (inc n) n))
                               0 (range (count params)))]
              (+ vs leak earg)))
          violations targets)))
    0 (ir/blocks f)))`

// livenessEntryLiveInExpr returns the count of values live-in to the entry block
// (block 0) that are NOT :load-arg insts. A well-formed function uses nothing
// before it is defined, so entry live-in must contain only (or no) arguments;
// any other value live at entry is a use-before-def bug.
const livenessEntryLiveInExpr = `
(let [f %s
      live (ir.passes.liveness/compute-liveness f)
      ein (:live-in (get live 0))]
  (reduce (fn [n v] (if (= :load-arg (ir/op v f)) n (inc n))) 0 ein))`

// livenessDeterminismExpr computes liveness twice with the cache cleared before
// each run (so the equality reflects the fixpoint, not a cache hit) and returns
// whether the two results are equal.
const livenessDeterminismExpr = `
(do
  (ir.passes.liveness/reset-liveness-cache!)
  (let [f %s
        a (ir.passes.liveness/compute-liveness f)
        _ (ir.passes.liveness/reset-liveness-cache!)
        b (ir.passes.liveness/compute-liveness f)]
    (= a b)))`

// livenessFixtures spans the shapes that exercise block-arg liveness: straight
// line, if-merge, loop (loop-carried value), nested loop, terminator-condition
// uses, and multi-use.
var livenessFixtures = []struct {
	name string
	src  string
}{
	{"straight-line", `(defn f [x] (+ x 1))`},
	{"if-merge", `(defn f [x] (if (< x 0) (- x) x))`},
	{"loop-accum", `(defn f [n] (loop [i 0 acc 0] (if (< i n) (recur (inc i) (+ acc i)) acc)))`},
	{"nested-loop", `(defn f [n] (loop [i 0 s 0] (if (< i n) (recur (inc i) (loop [j 0 t s] (if (< j n) (recur (inc j) (+ t j)) t))) s)))`},
	{"cond-uses-args", `(defn f [a b] (if (< a b) a b))`},
	{"multi-use", `(defn f [x] (let [y (* x x)] (+ y (+ y y))))`},
}

func lispInt(t *testing.T, v vm.Value) int {
	t.Helper()
	i, ok := v.(vm.Int)
	if !ok {
		t.Fatalf("expected Int, got %T (%v)", v, v)
	}
	return int(i)
}

// TestLivenessEdgeConsistency is the correctness gate: over the fixture battery,
// the edge-consistency invariant must have ZERO violations. A wrong pass fails
// loudly here (the check is independent of the pass's implementation).
func TestLivenessEdgeConsistency(t *testing.T) {
	ensureLoader()
	for _, fx := range livenessFixtures {
		t.Run(fx.name, func(t *testing.T) {
			f := buildLispIR(t, fx.src)
			v := lispEvalReturn(t, f, livenessInvariantExpr)
			if n := lispInt(t, v); n != 0 {
				t.Fatalf("%s: %d edge-consistency violation(s) — pass is WRONG.\nIR dump:\n%s",
					fx.name, n, lispDump(t, f))
			}
		})
	}
}

// TestLivenessEntrySound: nothing but arguments may be live-in to the entry.
func TestLivenessEntrySound(t *testing.T) {
	ensureLoader()
	for _, fx := range livenessFixtures {
		t.Run(fx.name, func(t *testing.T) {
			f := buildLispIR(t, fx.src)
			v := lispEvalReturn(t, f, livenessEntryLiveInExpr)
			if n := lispInt(t, v); n != 0 {
				t.Fatalf("%s: %d non-argument value(s) live-in to the entry block (use-before-def).\nIR dump:\n%s",
					fx.name, n, lispDump(t, f))
			}
		})
	}
}

// TestLivenessDeterminism: the fixpoint is deterministic (cache cleared between
// runs so equality reflects the computation, not the memo).
func TestLivenessDeterminism(t *testing.T) {
	ensureLoader()
	for _, fx := range livenessFixtures {
		t.Run(fx.name, func(t *testing.T) {
			f := buildLispIR(t, fx.src)
			v := lispEvalReturn(t, f, livenessDeterminismExpr)
			if v != vm.TRUE {
				t.Fatalf("%s: liveness not deterministic across cache-cleared recomputation", fx.name)
			}
		})
	}
}

// livenessUseSoundnessExpr: an INDEPENDENT soundness law. Every value referenced
// by an inst or terminator in a block must be either a non-param result of that
// block (defined locally) or live-in to it. A pass that fails to make a used
// parameter (or any cross-block value) live-in produces a violation. This does
// not consult the pass's edge logic — it re-derives the requirement from refs.
const livenessUseSoundnessExpr = `
(let [f %s
      live (ir.passes.liveness/compute-liveness f)]
  (reduce
    (fn [violations bid]
      (let [li (:live-in (get live bid))
            results (set (filter (fn [nid] (not (= :block-arg (ir/op nid f))))
                                 (ir/block-insts bid f)))
            inst-refs (mapcat (fn [nid] (ir/refs nid f)) (ir/block-insts bid f))
            term (ir/block-term bid f)
            term-refs (if (nil? term) [] (ir/refs term f))]
        (+ violations
           (reduce (fn [n v]
                     (if (or (contains? results v) (contains? li v)) n (inc n)))
                   0 (concat inst-refs term-refs)))))
    0 (ir/blocks f)))`

// TestLivenessUseSoundness: every used value is defined-here or live-in.
func TestLivenessUseSoundness(t *testing.T) {
	ensureLoader()
	for _, fx := range livenessFixtures {
		t.Run(fx.name, func(t *testing.T) {
			f := buildLispIR(t, fx.src)
			v := lispEvalReturn(t, f, livenessUseSoundnessExpr)
			if n := lispInt(t, v); n != 0 {
				t.Fatalf("%s: %d used value(s) neither defined-here nor live-in (unsound liveness).\nIR dump:\n%s",
					fx.name, n, lispDump(t, f))
			}
		})
	}
}

// TestLivenessHandTrace freezes the exact per-block (live-in, live-out) sets for
// a hand-traced fixture (AC-LV.2c). Fixture: (defn f [x] (if (< x 0) (- x) x)).
//
//	b0 entry: v0=LoadArg, v1=Const, v2=Lt v0 v1, BranchIf v2 -> b1 : b2
//	b1:       v4=Const, v5=Sub v4 v0, Branch -> b3(v5)
//	b2:       Branch -> b3(v0)
//	b3(v8):   v8=BlockArg, Return v8
//
// Hand trace: b3 returns v8, so v8 is live-in to b3; its edge-args (v5 from b1,
// v0 from b2) are therefore live-out of the arms; x (v0) is live-out of the
// entry (both arms need it) and live-in to each arm; nothing is live-in to the
// entry (v0 is defined there).
const livenessHandTraceExpr = `
(let [f %s
      live (ir.passes.liveness/compute-liveness f)
      expect {0 {:live-in #{}   :live-out #{0}}
              1 {:live-in #{0}  :live-out #{5}}
              2 {:live-in #{0}  :live-out #{0}}
              3 {:live-in #{8}  :live-out #{}}}]
  (reduce (fn [n bid]
            (let [e (get expect bid) a (get live bid)]
              (+ n (if (= (:live-in e) (:live-in a)) 0 1)
                   (if (= (:live-out e) (:live-out a)) 0 1))))
          0 (keys expect)))`

func TestLivenessHandTrace(t *testing.T) {
	ensureLoader()
	f := buildLispIR(t, `(defn f [x] (if (< x 0) (- x) x))`)
	v := lispEvalReturn(t, f, livenessHandTraceExpr)
	if n := lispInt(t, v); n != 0 {
		t.Fatalf("hand-trace: %d block/field(s) differ from the frozen expected sets.\ncomputed: %s\nIR:\n%s",
			n, lispEvalOn(t, f, "(ir.passes.liveness/compute-liveness %s)"), lispDump(t, f))
	}
}

// TestLivenessCorpus runs the invariant oracles over the shared pipelineCorpus
// (~20 real functions) — the AC-LV.2b (edge-arg / block-arg-sources contract)
// and AC-PO-LV.1 census. Every function must satisfy edge-consistency,
// use-soundness, and entry-soundness with zero violations. This is the real
// "broader usage" evidence: the pass is exercised and validated on real code,
// and it is the contract the block-arg stories (STORY-0070/0071) depend on.
func TestLivenessCorpus(t *testing.T) {
	ensureLoader()
	validated := 0
	for _, entry := range pipelineCorpus {
		e := entry
		t.Run(e.name, func(t *testing.T) {
			f, err := tryBuildLispIR(e.src)
			if err != nil {
				t.Skipf("build not supported in bare harness (needs pipeline context): %v", err)
			}
			edge := lispInt(t, lispEvalReturn(t, f, livenessInvariantExpr))
			use := lispInt(t, lispEvalReturn(t, f, livenessUseSoundnessExpr))
			ent := lispInt(t, lispEvalReturn(t, f, livenessEntryLiveInExpr))
			if v := edge + use + ent; v != 0 {
				t.Fatalf("%s: %d invariant violation(s) (edge=%d use=%d entry=%d)\nIR:\n%s",
					e.name, v, edge, use, ent, lispDump(t, f))
			}
			validated++
		})
	}
	t.Logf("liveness corpus census: %d/%d functions validated (edge-consistency + use-soundness + entry-soundness), 0 violations",
		validated, len(pipelineCorpus))
}

// reachableConsistentExpr ties the two liveness-family analyses together: every
// value compute-liveness marks live across a block boundary (∈ any block's
// live-in ∪ live-out) MUST be in reachable-nids' DCE read-set — otherwise the
// lowering DCE would drop a value that is genuinely used across blocks. A
// self-contained invariant (no retained reference code) that catches silent
// divergence between the two.
const reachableConsistentExpr = `
(let [f %s
      R (ir.passes.liveness/reachable-nids f)
      live (ir.passes.liveness/compute-liveness f)]
  (reduce (fn [v bid]
            (let [e (get live bid)
                  crossing (into (:live-in e) (:live-out e))]
              (reduce (fn [n x] (if (contains? R x) n (inc n))) v crossing)))
          0 (ir/blocks f)))`

// TestReachableNidsConsistent: reachable-nids ⊇ every cross-block-live value, over
// the corpus. (Exactness vs the old live-nids is locked separately by the
// byte-identical regenerated core_go_lowered/ tree; this is the fast, precise,
// no-lingering regression signal that keeps the two analyses from diverging.)
func TestReachableNidsConsistent(t *testing.T) {
	ensureLoader()
	for _, entry := range pipelineCorpus {
		e := entry
		t.Run(e.name, func(t *testing.T) {
			f, err := tryBuildLispIR(e.src)
			if err != nil {
				t.Skipf("build not supported in bare harness: %v", err)
			}
			if n := lispInt(t, lispEvalReturn(t, f, reachableConsistentExpr)); n != 0 {
				t.Fatalf("%s: %d cross-block-live value(s) missing from reachable-nids (DCE would drop a live value)\nIR:\n%s",
					e.name, n, lispDump(t, f))
			}
		})
	}
}

// TestLivenessMemoNonCollision: two structurally-different functions must not
// return each other's cached liveness.
func TestLivenessMemoNonCollision(t *testing.T) {
	ensureLoader()
	fa := buildLispIR(t, `(defn f [x] (if (< x 0) (- x) x))`)
	fb := buildLispIR(t, `(defn f [x] (if (< x 5) (+ x 1) (* x 2)))`)
	la := lispEvalOn(t, fa, "(ir.passes.liveness/compute-liveness %s)")
	lb := lispEvalOn(t, fb, "(ir.passes.liveness/compute-liveness %s)")
	if la == lb {
		t.Fatalf("memo collision: structurally-different functions produced identical liveness:\n%s", la)
	}
}
