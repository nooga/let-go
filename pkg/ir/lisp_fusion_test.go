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

// TestFusionAddInstBefore: add-inst-before! inserts a new instruction
// immediately before a target nid in its block's :insts list, with
// append-only flat :insts array (new nid = pre-insert count).
func TestFusionAddInstBefore(t *testing.T) {
	ensureLoader()

	// Build a simple function with a :call instruction.
	// (defn f [xs b] (let [g (if b count count)] (g xs)))
	// This creates a call to a variable 'g', which will be a :call inst.
	f := buildLispIR(t, `(defn f [xs b] (let [g (if b count count)] (g xs)))`)

	// Store f in a var for multiple operations
	passVarCounter++
	varNameF := fmt.Sprintf("*test-fn-%d*", passVarCounter)
	coreNS := rt.NS(rt.NameCoreNS)
	coreNS.Def(varNameF, f)

	// Before insertion: capture the nid count.
	countExpr := `(count (:insts @` + varNameF + `))`
	consts := vm.NewConsts()
	c := compiler.NewCompiler(consts, coreNS)
	c.SetSource("count-insts")
	_, countResult, err := c.CompileMultiple(strings.NewReader(countExpr))
	if err != nil {
		t.Fatalf("count insts: %v", err)
	}
	preInsertNid := int(countResult.(vm.Int))

	// Get the first :call inst in any block to use as target.
	findCallExpr := `(let [bs (ir/blocks ` + varNameF + `)]
		  (loop [blocks bs]
		    (if (empty? blocks) nil
		      (let [insts (ir/block-insts (first blocks) ` + varNameF + `)]
		        (let [found (loop [ns insts]
		                      (cond
		                        (empty? ns) nil
		                        (= :call (ir/op (first ns) ` + varNameF + `)) (first ns)
		                        :else (recur (rest ns))))]
		          (if found found (recur (rest blocks))))))))`
	consts = vm.NewConsts()
	c = compiler.NewCompiler(consts, coreNS)
	c.SetSource("find-target-nid")
	_, targetNidVal, err := c.CompileMultiple(strings.NewReader(findCallExpr))
	if err != nil {
		t.Fatalf("find target nid: %v", err)
	}

	if targetNidVal == vm.NIL {
		t.Skip("fixture produced no :call inst")
	}

	// Now we need to call add-inst-before! with the target nid.
	passVarCounter++
	varNameTarget := fmt.Sprintf("*test-target-%d*", passVarCounter)
	coreNS.Def(varNameTarget, targetNidVal)

	addInstExpr := fmt.Sprintf(`(ir.data/add-inst-before! %s %s :const [] 42)`, varNameF, varNameTarget)
	consts = vm.NewConsts()
	c = compiler.NewCompiler(consts, coreNS)
	c.SetSource("add-inst-before-test")
	_, result, err := c.CompileMultiple(strings.NewReader(addInstExpr))
	if err != nil {
		t.Fatalf("eval add-inst-before!: %v", err)
	}
	newNid := int(result.(vm.Int))

	// Assert 1: new nid should equal pre-insert count
	if newNid != preInsertNid {
		t.Fatalf("new nid (%d) should equal pre-insert count (%d)", newNid, preInsertNid)
	}

	// Assert 2: the new inst appears immediately before targetNid in block 0's :insts list.
	targetNid := int(targetNidVal.(vm.Int))

	// Use Lisp to verify the positions.
	// We need to find which block contains the target-nid and check positions there.
	passVarCounter++
	varNameNewNid := fmt.Sprintf("*new-nid-%d*", passVarCounter)
	coreNS.Def(varNameNewNid, vm.Int(newNid))

	verifyExpr := `(let [target ` + varNameTarget + `
		       bid (ir.data/block-of target ` + varNameF + `)
		       insts (vec (ir/block-insts bid ` + varNameF + `))
		       new-pos (first (keep-indexed (fn [i v] (when (= v ` + varNameNewNid + `) i)) insts))
		       tgt-pos (first (keep-indexed (fn [i v] (when (= v ` + varNameTarget + `) i)) insts))]
		   (if (and new-pos tgt-pos (= (inc new-pos) tgt-pos))
		     :ok
		     :mismatch))`

	consts = vm.NewConsts()
	c = compiler.NewCompiler(consts, coreNS)
	c.SetSource("add-inst-before-verify")
	_, result, err = c.CompileMultiple(strings.NewReader(verifyExpr))
	if err != nil {
		t.Fatalf("verify expr failed: %v", err)
	}

	if result != vm.Keyword("ok") {
		t.Fatalf("new inst at nid %d was not inserted immediately before target nid %d", newNid, targetNid)
	}

	t.Logf("✓ add-inst-before! inserted nid %d immediately before nid %d", newNid, targetNid)
}
