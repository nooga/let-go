/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package ir_test

import (
	"fmt"
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

func TestBytecodeFallthroughSize(t *testing.T) {
	ensureLoader()
	for _, tc := range []struct {
		name, form, calls, want string
	}{
		{"choice", `(defn choice [x] (if x 1 2))`, `[(f true) (f false) (f nil)]`, `[44 [1 2 2]]`},
		{"loop", `(defn loop-sum [n] (loop [i 0 acc 0] (if (< i n) (recur (inc i) (+ acc i)) acc)))`, `[(f 0) (f 1) (f 10) (f 1000)]`, `[136 [0 0 45 499500]]`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// Lower directly so compiler fallback cannot hide a lowering failure.
			// Each code word is int32; the byte count excludes constants and debug data.
			got := runLispExpr(t, fmt.Sprintf(`
			  (let [ir-fn (-> (quote %s) ir.build/build-fn ir.passes.pipeline/optimize-fn)
			        chunk (ir.lower/lower ir-fn)
			        f (chunk->fn 1 false chunk)]
			    (pr-str [(* 4 (ir/chunk-length chunk)) %s]))`, tc.form, tc.calls))
			if s, ok := got.(vm.String); !ok || string(s) != tc.want {
				t.Fatalf("[bytecode bytes, results] = %v, want %s", got, tc.want)
			}
		})
	}
}
