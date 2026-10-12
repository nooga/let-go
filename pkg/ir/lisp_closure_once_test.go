/*
 * Copyright (c) 2026 let-go contributors; see CONTRIBUTORS.
 * SPDX-License-Identifier: MIT
 */

package ir_test

import (
	"go/ast"
	"testing"
)

// countBoxNativeFn counts the rt.BoxNativeFn calls in node, nested closures
// included: every call is one closure allocation at run time.
func countBoxNativeFn(node ast.Node) int {
	n := 0
	ast.Inspect(node, func(x ast.Node) bool {
		if ce, ok := x.(*ast.CallExpr); ok && isSelector(ce.Fun, "rt", "BoxNativeFn") {
			n++
		}
		return true
	})
	return n
}

// A capturing closure is allocated once, where its fn form sits, and every
// later use reads that one value. The bytecode VM allocates exactly one Fn per
// evaluation of the form; the lowered Go must match it, both so `identical?`
// agrees and so a use never costs an allocation. Each case names the number
// of closure forms in the source, which is the number of rt.BoxNativeFn calls
// the lowered function may contain.
func TestCapturingClosureBoxedOncePerForm(t *testing.T) {
	ensureLoader()
	cases := []struct {
		name  string
		goFn  string
		form  string
		boxes int
	}{
		{"let-bound local used twice", "SameLocalCap",
			`(defn same-local-cap [y] (let [f (fn [x] (+ x y))] (identical? f f)))`, 1},
		{"loop-carried local", "LoopCarriedSame",
			`(defn loop-carried-same [y]
			   (loop [f (fn [x] (+ x y)) i 0]
			     (if (< i 2) (recur f (inc i)) (identical? f f))))`, 1},
		{"captured by another closure", "CapturedByClosure",
			`(defn captured-by-closure [y]
			   (let [f (fn [x] (+ x y)) g (fn [] f)] (identical? f (g))))`, 2},
		{"stored twice in a collection", "StoredInCollection",
			`(defn stored-in-collection [y]
			   (let [f (fn [x] (+ x y)) v [f f]] (identical? (nth v 0) (nth v 1))))`, 1},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			source := runLispString(t, `(do
			  (create-ns 'closureonce)
			  (ir.passes.pipeline/lower-ns-to-go "closureonce" 'closureonce '[`+tc.form+`]))`)
			fn, found := findFunc(parseLoweredGo(t, source), func(name string) bool { return name == tc.goFn })
			if !found {
				t.Fatalf("%s was not lowered\n%s", tc.goFn, source)
			}
			if got := countBoxNativeFn(fn.Body); got != tc.boxes {
				t.Fatalf("got %d rt.BoxNativeFn calls, want %d (one per closure form)\n%s", got, tc.boxes, source)
			}
		})
	}
}
