/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package ir_test

import (
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

// lispString evaluates expr and returns its value as a Go string.
func lispString(t *testing.T, expr string) string {
	t.Helper()
	v := runLispExpr(t, expr)
	s, ok := v.(vm.String)
	if !ok {
		t.Fatalf("%s returned %v (%T), want a string", expr, v, v)
	}
	return string(s)
}

// A sorted map is not interchangeable with a plain one: same entries, different
// contract. pipeline/expand-all rebuilt every map it walked with (into {} …),
// so a sorted map spliced into a macroexpansion reached the IR builder as an
// array-map. The fn compiled and ran; it just answered with a collection that
// had quietly lost its ordering, which no error reports and `=` does not catch.
func TestExpandAllKeepsSortedMapOrdering(t *testing.T) {
	ensureLoader()

	got := lispString(t, `(str (type (ir.passes.pipeline/expand-all (sorted-map :c 3 :a 1))))`)
	if got != "let-go.lang.PersistentTreeMap" {
		t.Fatalf("expand-all turned a sorted map into a %s", got)
	}

	// A sorted-map-by comparator is the part with no other home: the entries
	// survive any rebuild, the ordering only survives if `empty` carries the
	// comparator through.
	got = lispString(t, `(pr-str (ir.passes.pipeline/expand-all (sorted-map-by > 1 :a 3 :c 2 :b)))`)
	if got != "{3 :c, 2 :b, 1 :a}" {
		t.Fatalf("expand-all reordered a sorted-map-by map: %s", got)
	}

	// Entries are still macroexpanded — the point is to preserve the map type
	// while walking it, not to stop walking.
	runLispExpr(t, `(defmacro sorted-entry-macro [] 42)`)
	got = lispString(t, `(pr-str (ir.passes.pipeline/expand-all
	                              (sorted-map :a (list 'sorted-entry-macro))))`)
	if got != "{:a 42}" {
		t.Fatalf("expand-all stopped expanding sorted-map entries: %s", got)
	}

	// The plain-map path is unchanged.
	got = lispString(t, `(str (type (ir.passes.pipeline/expand-all {:a 1})))`)
	if got != "let-go.lang.Map" {
		t.Fatalf("expand-all changed the type of a plain map literal: %s", got)
	}
}

// build-map emits a map whose entries are all literals as a :const carrying the
// real collection, so the sorted type and its comparator ride along untouched.
func TestBuildMapEmitsSortedMapConstIntact(t *testing.T) {
	ensureLoader()

	runLispExpr(t, `(def *sorted-const-form* (list 'defn 'smfn [] (sorted-map-by > 1 :a 3 :c 2 :b)))`)
	consts := lispString(t, `(let [f (ir.build/build-fn *sorted-const-form*)]
	                           (pr-str (filterv sorted?
	                                            (mapv (fn [nid] (ir/aux nid f))
	                                                  (range (ir/inst-count f))))))`)
	if !strings.Contains(consts, "{3 :c, 2 :b, 1 :a}") {
		t.Fatalf("no sorted-map const with its own ordering in the built IR: %s", consts)
	}
}

// A sorted map with a non-literal entry has to be rebuilt by a call, and there
// is no call that reproduces a comparator. build-map declines instead of
// emitting an array-map rebuild that would silently reorder it; declining hands
// the fn to the bytecode compiler, which handles both cases correctly.
func TestBuildMapDeclinesSortedMapItWouldHaveToRebuild(t *testing.T) {
	ensureLoader()

	runLispExpr(t, `(def *sorted-subform-form*
	                  (list 'defn 'smfn [] (assoc (sorted-map :c 3 :a 1) :b 'conj)))`)
	err := runLispExprErr(`(ir.build/build-fn *sorted-subform-form*)`)
	if err == nil {
		t.Fatal("build-fn rebuilt a sorted map with a non-literal entry instead of declining")
	}
	if !strings.Contains(err.Error(), "sorted map") {
		t.Fatalf("decline does not say what it declined: %v", err)
	}

	// The plain-map rebuild it is modelled on still happens.
	runLispExpr(t, `(def *plain-subform-form* (list 'defn 'pmfn [] {:a 'conj}))`)
	if err := runLispExprErr(`(ir.build/build-fn *plain-subform-form*)`); err != nil {
		t.Fatalf("a plain map with a non-literal entry should still build: %v", err)
	}
}

// The Go backend renders a map const with vm.NewPersistentMap, which is the
// same flattening one layer down: a sorted map would come out of an AOT binary
// as a plain map. Nothing in the emitted Go can carry a comparator, so the
// lowerer reports "not a literal" (nil) and the fn falls back.
func TestLowerGoDoesNotRenderSortedMapAsPlainMap(t *testing.T) {
	ensureLoader()

	runLispExpr(t, `(def *sorted-lower-form* (list 'defn 'smfn [] (sorted-map :c 3 :a 1)))`)
	fn := runLispExpr(t, `(ir.build/build-fn *sorted-lower-form*)`)
	optimizeLispIR(t, fn)

	// :bridge mode, so an unlowerable body comes back as a nil :decl instead
	// of throwing.
	result := lowerGo(t, fn, ":bridge")
	if decl := result.ValueAt(vm.Keyword("decl")); decl != vm.NIL {
		t.Fatalf("a sorted-map const lowered to Go instead of declining:\n--- go ---\n%s",
			bindAndRenderGoDecl(t, result))
	}

	// The plain-map const it is modelled on still lowers.
	runLispExpr(t, `(def *plain-lower-form* (list 'defn 'pmfn [] {:c 3 :a 1}))`)
	plain := runLispExpr(t, `(ir.build/build-fn *plain-lower-form*)`)
	optimizeLispIR(t, plain)
	plainResult := lowerGo(t, plain, ":bridge")
	if plainResult.ValueAt(vm.Keyword("decl")) == vm.NIL {
		t.Fatal("a plain map const should still lower to a Go literal")
	}
	if rendered := bindAndRenderGoDecl(t, plainResult); !strings.Contains(rendered, "NewPersistentMap") {
		t.Fatalf("plain map const no longer renders as a Go map literal:\n--- go ---\n%s", rendered)
	}
}
