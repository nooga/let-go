/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package ir_test

import (
	"testing"
)

// TestTypeSideTableContract pins the :types side-table invariants introduced
// when inferred types moved off the inst tuple (PR #558). type-of, the
// lattice seed, and clone-inst each implement the same precedence rule —
// side table first, tuple slot 5 (construction-time :unknown) as fallback —
// and the clone carry was a review catch, so lock the contract down:
//
//  1. fallback: an inst with no side-table entry reads :unknown
//  2. set-type! writes are read back through type-of (side table wins)
//  3. overwrite: the last write wins
//  4. sparse growth: setting only a high nid grows the table; untouched
//     nids in the gap still read the fallback
//  5. clone-inst carries the side-table entry (typed clone keeps its type,
//     untyped clone stays :unknown)
//  6. seed-state-from-inst-types! reads the same merged view
func TestTypeSideTableContract(t *testing.T) {
	ensureLoader()
	got := runLispString(t, `(pr-str
		(let [f        (ir.build/build-fn '(defn t558 [x] (+ x 1)))
		      nid      (first (ir/block-insts (ir/entry-block f) f))
		      last-nid (- (ir/inst-count f) 1)
		      mid-nid  (- last-nid 1)]
		  [;; 1. fallback: nothing set yet -> construction default
		   (ir/type-of nid f)
		   ;; 2. side table wins once set
		   (do (ir/set-type! f nid :int) (ir/type-of nid f))
		   ;; 3. overwrite: last write wins
		   (do (ir/set-type! f nid :float) (ir/type-of nid f))
		   ;; 4. sparse: only the high nid is set; the gap nid still falls back
		   (do (ir/set-type! f last-nid :string)
		       [(ir/type-of last-nid f) (ir/type-of mid-nid f)])
		   ;; 5. clone carry: typed clone keeps the type, untyped stays :unknown
		   (ir/type-of (ir/clone-inst f nid) f)
		   (ir/type-of (ir/clone-inst f mid-nid) f)
		   ;; 6. the typeinfra seed reads the merged view (side table first)
		   (let [s (ir.lattice/seed-state-from-inst-types! (ir.lattice/new-typeinfra-state f) f)]
		     (ir.lattice/state-type s nid))]))`)
	want := `[:unknown :int :float [:string :unknown] :float :unknown :float]`
	if got != want {
		t.Fatalf("side-table contract mismatch:\n got: %s\nwant: %s", got, want)
	}
}
