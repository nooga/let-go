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

// TestBuildFnDefnAttrMap covers the defn header grammar
// (defn name docstring? attr-map? [params] body) in the IR builder (#289):
// the attr-map position ({:added "1.0"} etc.) is var metadata the lowering
// must skip, not an arg pattern. Before the fix, build-fn treated the map as
// the arglist and threw "unsupported destructuring pattern: :added".
func TestBuildFnDefnAttrMap(t *testing.T) {
	ensureLoader()

	cases := []struct {
		label string
		src   string
	}{
		{"doc+attr single-arity", `(defn meta-doc-fn "defn doc" {:added "1.0"} [x] x)`},
		{"attr only single-arity", `(defn attr-fn {:added "1.0"} [x] x)`},
		{"doc+attr multi-arity", `(defn meta-doc-multi "doc" {:added "2.0"} ([x] x) ([x y] y))`},
		{"attr only multi-arity", `(defn attr-multi {:added "2.0"} ([x] x) ([x y] y))`},
		{"doc only (regression guard)", `(defn doc-fn "doc" [x] x)`},
		{"bare (regression guard)", `(defn bare-fn [x] x)`},
	}
	for _, tc := range cases {
		f := buildLispIR(t, tc.src)
		if f == vm.NIL {
			t.Errorf("%s: build-fn returned nil for %s", tc.label, tc.src)
		}
	}
}
