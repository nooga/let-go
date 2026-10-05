/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package compiler

import (
	"testing"

	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

// TestInitFromLGBLeavesLazyNamespacesUnloaded pins what BenchmarkInitFromLGB
// measures: process startup runs clojure.core and the lg baseline namespaces
// and leaves every other bundled namespace marked needs-load until it is
// required. The benchmark's per-iteration path must do the same, or code
// added to a lazily loaded namespace would move a startup gate that a real
// `lg` process never pays.
func TestInitFromLGBLeavesLazyNamespacesUnloaded(t *testing.T) {
	if len(rt.CoreCompiledLGB) == 0 {
		t.Skip("no precompiled core_compiled.lgb")
	}
	const lazy = "edn" // bundled, not core, not a baseline namespace
	for _, b := range rt.LgBaselineNSNames() {
		if b == lazy {
			t.Fatalf("%s is a baseline namespace; pick a lazily loaded one", lazy)
		}
	}

	// Whether the lazy namespace's defs are already bound depends on what
	// earlier tests in this process required, so only the transition is
	// pinned: an unbound def must stay unbound across the init path.
	boundBefore := lazyDefBound(lazy, "pretty")

	// Package init already left the marker set; drop it so that only an init
	// path that marks the namespace itself can satisfy the check below, and
	// put it back if the path under test did not, so later tests still load
	// the namespace on require.
	rt.ClearNSNeedsLoad(lazy)
	defer func() {
		if !boundBefore && !rt.NSNeedsLoad(lazy) {
			rt.MarkNSNeedsLoad(lazy)
		}
	}()

	initFromLGB(t)
	t.Logf("%s/pretty bound before=%v after=%v needs-load=%v", lazy, boundBefore, lazyDefBound(lazy, "pretty"), rt.NSNeedsLoad(lazy))

	if rt.LookupNS(lazy) == nil {
		t.Fatalf("%s is not registered after init; is it still in the bundle?", lazy)
	}
	if !rt.NSNeedsLoad(lazy) {
		t.Fatalf("%s is not marked needs-load after init: the init path ran (or dropped) a lazily loaded namespace", lazy)
	}
	if !boundBefore && lazyDefBound(lazy, "pretty") {
		t.Fatalf("%s/pretty became bound during init: the init path ran a lazily loaded namespace", lazy)
	}
	for _, b := range rt.LgBaselineNSNames() {
		if rt.NSNeedsLoad(b) {
			t.Fatalf("baseline namespace %s is marked needs-load after init; startup runs it eagerly", b)
		}
	}
}

func lazyDefBound(ns, name string) bool {
	n := rt.LookupNS(ns)
	if n == nil {
		return false
	}
	v := n.LookupLocal(vm.Symbol(name))
	return v != nil && v.IsBound()
}
