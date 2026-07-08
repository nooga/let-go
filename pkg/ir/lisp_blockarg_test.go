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

// firstBlockArg binds `ba` to the first :block-arg InstId in f (nil if none).
const firstBlockArg = `
      ba (loop [bs (ir/blocks f)]
           (if (empty? bs) nil
             (let [found (loop [ns (ir/block-insts (first bs) f)]
                           (cond
                             (empty? ns) nil
                             (= :block-arg (ir/op (first ns) f)) (first ns)
                             :else (recur (rest ns))))]
               (if found found (recur (rest bs))))))`

// resolvesThrough: true iff the first block-arg is single-source and
// resolve-single-source sees through it to a non-block-arg definition.
const resolvesThroughExpr = `(let [f %s` + firstBlockArg + `]
  (and (not (nil? ba))
       (not= ba (ir/resolve-single-source f ba))
       (not= :block-arg (ir/op (ir/resolve-single-source f ba) f))))`

// mergeUnchanged: true iff the first block-arg is a genuine merge and
// resolve-single-source leaves it untouched.
const mergeUnchangedExpr = `(let [f %s` + firstBlockArg + `]
  (and (not (nil? ba))
       (= ba (ir/resolve-single-source f ba))))`

// hasBlockArgExpr: true iff f has at least one :block-arg.
const hasBlockArgExpr = `(let [f %s` + firstBlockArg + `] (not (nil? ba)))`

// TestResolveSingleSourceBlockArg: single-source block-args resolve through to
// their source; genuine merges do not. Pure analysis — the IR is not mutated
// (guarded separately by the byte-identical regenerated core tree).
func TestResolveSingleSourceBlockArg(t *testing.T) {
	ensureLoader()

	// Single-source: both arms of the if yield the same value (x), so the join
	// block-arg is a trivial copy of x — resolve-single-source must see through it.
	t.Run("single-source-resolves", func(t *testing.T) {
		f := buildLispIR(t, `(defn f [x] (if (< x 0) x x))`)
		if lispEvalReturn(t, f, hasBlockArgExpr) != vm.TRUE {
			t.Skip("fixture produced no block-arg (folded)")
		}
		if lispEvalReturn(t, f, resolvesThroughExpr) != vm.TRUE {
			t.Fatalf("single-source block-arg did NOT resolve through\nIR:\n%s", lispDump(t, f))
		}
	})

	// Genuine merge: the arms yield DIFFERENT values, so the join block-arg is a
	// real merge — resolve-single-source must leave it untouched.
	t.Run("merge-unchanged", func(t *testing.T) {
		f := buildLispIR(t, `(defn f [x] (if (< x 0) (- x) x))`)
		if lispEvalReturn(t, f, hasBlockArgExpr) != vm.TRUE {
			t.Fatalf("expected a merge block-arg, got none\nIR:\n%s", lispDump(t, f))
		}
		if lispEvalReturn(t, f, mergeUnchangedExpr) != vm.TRUE {
			t.Fatalf("genuine merge block-arg was wrongly resolved\nIR:\n%s", lispDump(t, f))
		}
	})
}

// pureCallThroughBlockArgExpr finds the first :call whose callee (refs[0]) is a
// :block-arg and reports whether purity/pure-call? recognizes it as pure — i.e.
// whether the callee resolution sees through the block-arg to the pure builtin.
const pureCallThroughBlockArgExpr = `(let [f %s
      vf (ir.passes.mutability/analyze-var-stability f)
      call (loop [bs (ir/blocks f)]
             (if (empty? bs) nil
               (let [found (loop [ns (ir/block-insts (first bs) f)]
                             (cond
                               (empty? ns) nil
                               (and (= :call (ir/op (first ns) f))
                                    (pos? (count (ir/refs (first ns) f)))
                                    (= :block-arg (ir/op (first (ir/refs (first ns) f)) f)))
                               (first ns)
                               :else (recur (rest ns))))]
                 (if found found (recur (rest bs))))))]
  (if (nil? call) :no-blockarg-call (ir.passes.purity/pure-call? f call vf)))`

// TestPureCallThroughBlockArg (AC-BA.2): a call whose callee is a block-arg whose
// edges all denote the same pure builtin — here (if b count count), which yields
// two distinct :load-var insts of #'count in different blocks — must be recognized
// as pure. This is the common source-level shape; the resolver must see through
// the (equivalent, not identical) sources. A regression that fails to resolve
// makes this RED — do not weaken it.
func TestPureCallThroughBlockArg(t *testing.T) {
	ensureLoader()
	f := buildLispIR(t, `(defn f [xs b] (let [g (if b count count)] (g xs)))`)
	res := lispEvalReturn(t, f, pureCallThroughBlockArgExpr)
	switch res {
	case vm.TRUE:
		// callee resolved through the block-arg to #'count → recognized pure
	case vm.FALSE:
		t.Fatalf("pure-call? did NOT recognize a block-arg'd pure callee (resolution failed)\nIR:\n%s", lispDump(t, f))
	default:
		t.Fatalf("expected a block-arg'd callee in the fixture, got %v\nIR:\n%s", res, lispDump(t, f))
	}
}
