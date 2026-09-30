/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package ir_test

import (
	"strings"
	"testing"
)

// validate-fn! judges whether IR is well formed, for every consumer. A
// :branch-if whose true and false edges pass different args is well formed:
// each edge's args match its own target's params, and ir.passes.cleanup
// produces exactly that when it compacts dead params per target. Only the
// bytecode lowerer needs the two edges to agree, because it places one arg
// shape on the stack under the cond; ir.lower checks that itself.

// cleanedFixtures all leave a :branch-if with asymmetric edge args after
// cleanup (a loop exit passes fewer values than the back edge).
var cleanedFixtures = []struct {
	name string
	src  string
}{
	{"dead-param", deadParamFixture},
	{"sum-loop", `(defn ksum [n] (loop [i 0 acc 0] (if (< i n) (recur (+ i 1) (+ acc i)) acc)))`},
	{"vec-loop", `(defn kvec [v] (loop [i 0 acc 0] (if (< i (count v)) (recur (+ i 1) (+ acc (nth v i))) acc)))`},
	{"nested-loop", `(defn f [n] (loop [i 0 s 0] (if (< i n) (recur (inc i) (loop [j 0 t s] (if (< j n) (recur (inc j) (+ t j)) t))) s)))`},
}

const cleanedFormExpr = `(ir.passes.cleanup/cleanup (ir.passes.pipeline/optimize-fn %s))`

func TestValidateAcceptsAsymmetricBranchIfArgs(t *testing.T) {
	ensureLoader()
	for _, fx := range cleanedFixtures {
		t.Run(fx.name, func(t *testing.T) {
			f, err := tryBuildLispIR(fx.src)
			if err != nil {
				t.Fatalf("cannot build IR: %v", err)
			}
			got := lispEvalOn(t, f, `(try (do (ir.validate/validate-fn! `+cleanedFormExpr+` "cleanup") "ok")
   (catch Exception e (.getMessage e)))`)
			if got != `"ok"` {
				t.Fatalf("validate-fn! rejects well-formed cleaned IR: %s", got)
			}
		})
	}
}

// Without its own check the bytecode lowerer lowers the true edge's args onto
// both edges.
func TestLowerRejectsAsymmetricBranchIfArgs(t *testing.T) {
	ensureLoader()
	f, err := tryBuildLispIR(cleanedFixtures[1].src)
	if err != nil {
		t.Fatalf("cannot build IR: %v", err)
	}
	got := lispEvalOn(t, f, `(try (do (ir.lower/lower `+cleanedFormExpr+`) "lowered")
   (catch Exception e (.getMessage e)))`)
	if !strings.Contains(got, "asymmetric branch-target args") {
		t.Fatalf("ir.lower/lower on an asymmetric :branch-if: want a refusal, got %s", got)
	}
}
