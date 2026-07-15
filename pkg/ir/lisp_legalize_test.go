/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package ir_test

import (
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

// A hand-built 2-block Function carrying a cross-block ref:
//
//	block0: v0 = load-arg 0        (cheap-load, exempt)
//	        v1 = (+ v0 v0)         (non-cheap, defined in block0)
//	        branch -> block1
//	block1: v2 = (+ v1 v1)         (refs v1 ACROSS blocks — illegal SSA)
//	        return v2
//
// This is the minimal shape ir.validate/check-no-cross-block-refs! rejects,
// and the shape a real build (e.g. try-assign-stmts) trips on. Constructed
// directly so the fixture does not depend on build.lg's threading behaviour.
const xblockFixture = `
  (let [f  (ir/new-fn "xblock" 1 false)
        v0 (ir/add-inst f 0 :load-arg [] 0)
        v1 (ir/add-inst f 0 :add [v0 v0] nil)
        b1 (ir/add-block f)]
    (ir/add-pred! f b1 0)
    (ir/add-terminator! f 0 :branch [] (ir/new-branch-target b1 []))
    (let [v2 (ir/add-inst f b1 :add [v1 v1] nil)]
      (ir/add-terminator! f b1 :return [v2] nil))
    f)`

// Guard: the fixture is genuinely invalid SSA (a cross-block ref), so the
// GREEN test below is not vacuous.
func TestLegalizeFixtureIsGenuinelyCrossBlock(t *testing.T) {
	ensureLoader()
	got := string(runLispExpr(t, `
      (do (require 'ir.data) (require 'ir.validate)
        (let [f `+xblockFixture+`]
          (try (ir.validate/validate-fn! f "probe") "UNEXPECTEDLY-VALID"
               (catch e (str e)))))`).(vm.String))
	if !strings.Contains(got, "cross-block ref") {
		t.Fatalf("fixture should carry a cross-block ref; validate said: %s", got)
	}
}

// RED→GREEN driver: the legalize pass must thread the cross-block value into a
// fresh block-arg so the function validates, and block1 must gain exactly one
// param (the threaded arg) — proving the value was threaded, not deleted.
func TestLegalizeThreadsCrossBlockRefIntoBlockArg(t *testing.T) {
	ensureLoader()
	got := string(runLispExpr(t, `
      (do (require 'ir.data) (require 'ir.validate)
        (let [f `+xblockFixture+`]
          (ir.passes.legalize/legalize-fn! f)
          (try (ir.validate/validate-fn! f "legalize-test")
               (pr-str [:validated (count (ir/block-params 1 f))])
               (catch e (str "STILL-INVALID: " e)))))`).(vm.String))
	if got != `[:validated 1]` {
		t.Fatalf("legalize should thread the cross-block ref into a block1 block-arg; got %s", got)
	}
}
