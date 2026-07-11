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

// TestFusionMatchReduceMap: match-chain recognizes (reduce + 0 (map inc coll))
func TestFusionMatchReduceMap(t *testing.T) {
	ensureLoader()

	f := buildLispIR(t, `(defn f [coll] (reduce + 0 (map inc coll)))`)

	// Store f in a var for Lisp evaluation
	passVarCounter++
	varNameF := fmt.Sprintf("*test-fn-%d*", passVarCounter)
	coreNS := rt.NS(rt.NameCoreNS)
	coreNS.Def(varNameF, f)

	// Find the reduce call (aux = 3 for reduce)
	findConsumerExpr := `(let [bs (ir/blocks ` + varNameF + `)]
		  (first (filter (fn [nid] (and (= :call (ir/op nid ` + varNameF + `))
		                                  (= 3 (ir/aux nid ` + varNameF + `))))
		                  (ir/block-insts (first bs) ` + varNameF + `))))`
	consts := vm.NewConsts()
	c := compiler.NewCompiler(consts, coreNS)
	c.SetSource("find-consumer")
	_, consumerNidVal, err := c.CompileMultiple(strings.NewReader(findConsumerExpr))
	if err != nil {
		t.Fatalf("find consumer nid: %v", err)
	}
	if consumerNidVal == vm.NIL {
		t.Skip("fixture produced no :call inst")
	}

	passVarCounter++
	varNameConsumer := fmt.Sprintf("*consumer-%d*", passVarCounter)
	coreNS.Def(varNameConsumer, consumerNidVal)

	// Call match-chain on the consumer
	matchExpr := `(ir.passes.fusion/match-chain ` + varNameF + ` ` + varNameConsumer + `)`
	consts = vm.NewConsts()
	c = compiler.NewCompiler(consts, coreNS)
	c.SetSource("test-match-chain")
	_, matchResult, err := c.CompileMultiple(strings.NewReader(matchExpr))
	if err != nil {
		t.Fatalf("eval match-chain: %v", err)
	}

	if matchResult == vm.NIL {
		t.Fatalf("match-chain returned nil for (reduce + 0 (map inc coll))")
	}

	// Assert key presence and correct producer-kind
	checkKeyExpr := `(let [m (ir.passes.fusion/match-chain ` + varNameF + ` ` + varNameConsumer + `)]
	                   (and (contains? m :consumer-kind)
	                        (contains? m :producer-kind)
	                        (contains? m :consumer)
	                        (contains? m :producer)
	                        (contains? m :g)
	                        (contains? m :init)
	                        (contains? m :xform-fn)
	                        (contains? m :coll)
	                        (= (m :producer-kind) "map")))`
	consts = vm.NewConsts()
	c = compiler.NewCompiler(consts, coreNS)
	c.SetSource("check-match-keys")
	_, checkResult, err := c.CompileMultiple(strings.NewReader(checkKeyExpr))
	if err != nil {
		t.Fatalf("check match keys: %v", err)
	}

	if checkResult != vm.TRUE {
		t.Fatalf("match map did not have expected keys or producer-kind='map': %v", matchResult)
	}

	t.Logf("✓ match-chain recognized (reduce + 0 (map inc coll))")
}

// TestFusionMatchNoProducer: match-chain returns nil for (reduce + 0 coll) with no producer
func TestFusionMatchNoProducer(t *testing.T) {
	ensureLoader()

	f := buildLispIR(t, `(defn f [coll] (reduce + 0 coll))`)

	passVarCounter++
	varNameF := fmt.Sprintf("*test-fn-%d*", passVarCounter)
	coreNS := rt.NS(rt.NameCoreNS)
	coreNS.Def(varNameF, f)

	// Find the reduce call (aux = 3 for reduce)
	findConsumerExpr := `(let [bs (ir/blocks ` + varNameF + `)]
		  (first (filter (fn [nid] (and (= :call (ir/op nid ` + varNameF + `))
		                                  (= 3 (ir/aux nid ` + varNameF + `))))
		                  (ir/block-insts (first bs) ` + varNameF + `))))`
	consts := vm.NewConsts()
	c := compiler.NewCompiler(consts, coreNS)
	c.SetSource("find-consumer-no-prod")
	_, consumerNidVal, err := c.CompileMultiple(strings.NewReader(findConsumerExpr))
	if err != nil {
		t.Fatalf("find consumer nid: %v", err)
	}
	if consumerNidVal == vm.NIL {
		t.Skip("fixture produced no :call inst")
	}

	passVarCounter++
	varNameConsumer := fmt.Sprintf("*consumer-%d*", passVarCounter)
	coreNS.Def(varNameConsumer, consumerNidVal)

	// Call match-chain on the consumer
	matchExpr := `(ir.passes.fusion/match-chain ` + varNameF + ` ` + varNameConsumer + `)`
	consts = vm.NewConsts()
	c = compiler.NewCompiler(consts, coreNS)
	c.SetSource("test-match-chain-no-prod")
	_, matchResult, err := c.CompileMultiple(strings.NewReader(matchExpr))
	if err != nil {
		t.Fatalf("eval match-chain: %v", err)
	}

	if matchResult != vm.NIL {
		t.Fatalf("match-chain should return nil for (reduce + 0 coll) with no producer, got: %v", matchResult)
	}

	t.Logf("✓ match-chain correctly returned nil for (reduce + 0 coll)")
}

// TestFusionMatchNoConsumer: match-chain returns nil for (map inc coll) with no consumer
func TestFusionMatchNoConsumer(t *testing.T) {
	ensureLoader()

	f := buildLispIR(t, `(defn f [coll] (map inc coll))`)

	passVarCounter++
	varNameF := fmt.Sprintf("*test-fn-%d*", passVarCounter)
	coreNS := rt.NS(rt.NameCoreNS)
	coreNS.Def(varNameF, f)

	// Find the :call (the map call)
	findMapExpr := `(let [bs (ir/blocks ` + varNameF + `)]
		  (first (filter (fn [nid] (= :call (ir/op nid ` + varNameF + `)))
		                  (ir/block-insts (first bs) ` + varNameF + `))))`
	consts := vm.NewConsts()
	c := compiler.NewCompiler(consts, coreNS)
	c.SetSource("find-map")
	_, mapNidVal, err := c.CompileMultiple(strings.NewReader(findMapExpr))
	if err != nil {
		t.Fatalf("find map nid: %v", err)
	}
	if mapNidVal == vm.NIL {
		t.Skip("fixture produced no :call inst")
	}

	passVarCounter++
	varNameMap := fmt.Sprintf("*map-%d*", passVarCounter)
	coreNS.Def(varNameMap, mapNidVal)

	// Call match-chain on the map (should be nil since map is not a valid consumer)
	matchExpr := `(ir.passes.fusion/match-chain ` + varNameF + ` ` + varNameMap + `)`
	consts = vm.NewConsts()
	c = compiler.NewCompiler(consts, coreNS)
	c.SetSource("test-match-chain-no-consumer")
	_, matchResult, err := c.CompileMultiple(strings.NewReader(matchExpr))
	if err != nil {
		t.Fatalf("eval match-chain: %v", err)
	}

	if matchResult != vm.NIL {
		t.Fatalf("match-chain should return nil for (map inc coll) without a consumer, got: %v", matchResult)
	}

	t.Logf("✓ match-chain correctly returned nil for (map inc coll)")
}

// TestFusionMatchReduceFilter: match-chain recognizes (reduce + 0 (filter odd? coll))
func TestFusionMatchReduceFilter(t *testing.T) {
	ensureLoader()

	f := buildLispIR(t, `(defn f [coll] (reduce + 0 (filter odd? coll)))`)

	passVarCounter++
	varNameF := fmt.Sprintf("*test-fn-%d*", passVarCounter)
	coreNS := rt.NS(rt.NameCoreNS)
	coreNS.Def(varNameF, f)

	// Find the reduce call (aux = 3 for reduce)
	findConsumerExpr := `(let [bs (ir/blocks ` + varNameF + `)]
		  (first (filter (fn [nid] (and (= :call (ir/op nid ` + varNameF + `))
		                                  (= 3 (ir/aux nid ` + varNameF + `))))
		                  (ir/block-insts (first bs) ` + varNameF + `))))`
	consts := vm.NewConsts()
	c := compiler.NewCompiler(consts, coreNS)
	c.SetSource("find-consumer-filter")
	_, consumerNidVal, err := c.CompileMultiple(strings.NewReader(findConsumerExpr))
	if err != nil {
		t.Fatalf("find consumer nid: %v", err)
	}
	if consumerNidVal == vm.NIL {
		t.Skip("fixture produced no :call inst")
	}

	passVarCounter++
	varNameConsumer := fmt.Sprintf("*consumer-%d*", passVarCounter)
	coreNS.Def(varNameConsumer, consumerNidVal)

	// Check that match-chain returns a match with producer-kind "filter"
	checkFilterExpr := `(let [m (ir.passes.fusion/match-chain ` + varNameF + ` ` + varNameConsumer + `)]
	                      (and (some? m)
	                           (= (m :producer-kind) "filter")))`
	consts = vm.NewConsts()
	c = compiler.NewCompiler(consts, coreNS)
	c.SetSource("test-match-filter")
	_, filterResult, err := c.CompileMultiple(strings.NewReader(checkFilterExpr))
	if err != nil {
		t.Fatalf("check filter match: %v", err)
	}

	if filterResult != vm.TRUE {
		t.Fatalf("match-chain should return a match with producer-kind='filter' for (reduce + 0 (filter odd? coll))")
	}

	t.Logf("✓ match-chain recognized (reduce + 0 (filter odd? coll))")
}
