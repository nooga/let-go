/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package ir_test

import (
	"testing"
)

// TestTypeinferDrainCapSemantics pins the work-unit guard's contract (review
// round 3 on #446): a cap of N performs EXACTLY N drains before bailing (the
// prior `>` check admitted N+1, and a cap of 0 still drained once), and a
// queue that empties on its own never trips the cap. :drains is recorded via
// lat/*ti-counters* on both the bail and the normal exit, so the counter is
// the observable for both sides.
func TestTypeinferDrainCapSemantics(t *testing.T) {
	ensureLoader()
	got := runLispString(t, `(pr-str
		(let [run (fn [cap]
		            ;; fresh function per run: typeinfer flushes inferred types
		            ;; back into the function, so a reused one converges early.
		            (let [f (ir.build/build-fn '(defn t446-cap [x y]
		                                          (let [a (+ x 1)
		                                                b (* y 2)]
		                                            (- a b))))]
		              (binding [ir.lattice/*ti-counters* (atom {})
		                        ir.passes.typeinfer/*typeinfer-max-drains* cap]
		                (ir.passes.typeinfer/typeinfer f)
		                (:drains (deref ir.lattice/*ti-counters*)))))
		      uncapped (run nil)]
		  [;; cap 0: zero drains, not one
		   (run 0)
		   ;; cap 2: exactly two drains, not three
		   (run 2)
		   ;; the fixture genuinely needs more than 2 drains, so cap 2 truncated
		   (> uncapped 2)
		   ;; an ample cap must not fire: same drain count as unbounded
		   (= uncapped (run 1000000))]))`)
	want := `[0 2 true true]`
	if got != want {
		t.Fatalf("drain-cap semantics mismatch:\n got: %s\nwant: %s", got, want)
	}
}
