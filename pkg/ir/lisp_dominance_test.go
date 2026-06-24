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

// defDomFn installs an IR Function as a fresh core-ns var and returns its name
// so dominance/licm queries can reference it from runLispExpr.
func defDomFn(t *testing.T, src string) string {
	t.Helper()
	passVarCounter++
	name := fmt.Sprintf("*dom-fn-%d*", passVarCounter)
	rt.NS(rt.NameCoreNS).Def(name, buildLispIR(t, src))
	return name
}

// dominates-idom? takes a precomputed idom array so a hot caller can compute the
// dominator tree ONCE instead of rebuilding it per dominance query (the cause of
// LICM's quadratic blowup). It must agree with dominates? for every block pair.
func TestDominatesIdomMatchesDominates(t *testing.T) {
	ensureLoader()
	fn := defDomFn(t, `(defn f [n] (loop [i 0] (if (< i n) (recur (inc i)) i)))`)

	got := runLispExpr(t, fmt.Sprintf(`
      (let [f %s
            idom (ir.dominance/dominators f)
            bs   (ir/blocks f)]
        (every? (fn [a]
                  (every? (fn [b]
                            (= (ir.dominance/dominates? f a b)
                               (ir.dominance/dominates-idom? idom a b)))
                          bs))
                bs))`, fn))
	if got != vm.TRUE {
		t.Fatalf("dominates-idom? must match dominates? for every block pair, got %v", got)
	}
}

// back-edges must still find the loop's back-edge after the dominators-once
// refactor — and report none for straight-line code.
func TestBackEdgesUnchangedByMemoization(t *testing.T) {
	ensureLoader()

	loopFn := defDomFn(t, `(defn f [n] (loop [i 0] (if (< i n) (recur (inc i)) i)))`)
	if got := runLispExpr(t, fmt.Sprintf(`(pos? (count (ir.passes.licm/back-edges %s)))`, loopFn)); got != vm.TRUE {
		t.Fatalf("expected the loop fn to have at least one back-edge, got %v", got)
	}

	straightFn := defDomFn(t, `(defn g [a b] (+ a b))`)
	if got := runLispExpr(t, fmt.Sprintf(`(zero? (count (ir.passes.licm/back-edges %s)))`, straightFn)); got != vm.TRUE {
		t.Fatalf("expected straight-line fn to have no back-edges, got %v", got)
	}
}
