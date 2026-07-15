/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

// TestReapplyGeneratedPrimitivesRestoresClobberedRoot guards the #438
// coordination the unified eager-hybrid pass relies on (#506). When a hybrid
// namespace's bundle chunk re-Defs a generated native-primitive var with its
// bootstrap closure, ReapplyGeneratedPrimitives must restore the recorded
// adapter as the var's root — otherwise the root stays deviated and lowered
// direct-call sites are stranded on the var-dispatch trampoline (the native
// fast path never re-engages). LoadCoreBundle calls this after every eager
// hybrid chunk; the on-demand loader calls it after a lazy load.
//
// The mechanism is exercised directly against a synthetic namespace on purpose:
// no real namespace is both an eager hybrid AND a carrier of generated
// primitives today (the only eager hybrid, async, has none; clojure.string has
// one but loads lazily), so a bundle-driven assertion would pass even with the
// reapply removed. This keeps the guard meaningful the moment a generated-
// primitive-bearing hybrid is added.
func TestReapplyGeneratedPrimitivesRestoresClobberedRoot(t *testing.T) {
	const nsName = "test.reapplyhybrid"
	ns := DefNSBare(nsName)

	// Stand-ins for the native adapter and the bootstrap closure a hybrid chunk
	// would Def over it. Distinct comparable Values suffice: the guard compares
	// roots by identity (vm.sameRootIdentity) and Deref equality is exact here.
	adapter := vm.String("native-adapter")
	bootstrap := vm.String("bootstrap-closure")

	// Bind + guard + record the adapter exactly as RegisterGeneratedPrimitives.
	defGeneratedPrimitive(ns, nsName, "prim", adapter)
	v := ns.LookupLocal(vm.Symbol("prim"))
	if v == nil {
		t.Fatal("prim var not defined")
	}
	if got := v.Deref(); got != adapter {
		t.Fatalf("adapter not bound as root: got %v, want %v", got, adapter)
	}

	// A hybrid chunk's (def prim …) overwrites the adapter root in place.
	setPrimitiveRoot(ns, "prim", bootstrap)
	if got := v.Deref(); got != bootstrap {
		t.Fatalf("clobber did not take: got %v, want %v", got, bootstrap)
	}

	// The eager-hybrid pass (and the on-demand loader) restore it.
	ReapplyGeneratedPrimitives(nsName)
	if got := v.Deref(); got != adapter {
		t.Fatalf("reapply did not restore adapter: got %v, want %v", got, adapter)
	}
}
