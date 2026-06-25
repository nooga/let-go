/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

// A lowered package with native multimethod arms registers its multifn names so
// the runtime can freeze them as the native baseline once the namespace has
// finished loading (all build-time defmethods applied). Mirrors the
// RegisterGoOverrides queue-until-apply ordering.
func TestRegisterNativeMultiFnsFreezesAtApply(t *testing.T) {
	const ns = "test-native-mm-queue"
	delete(pendingGoOverrides, ns)
	delete(pendingNativeMultiFns, ns)

	RegisterNativeMultiFns(ns, []string{"mm"})

	if LookupNS(ns) != nil {
		t.Fatalf("precondition: ns %q should not exist yet", ns)
	}

	target := vm.NewNamespace(ns)
	mm := vm.NewMultiFn("mm", nil, vm.Keyword("default"))
	target.Def("mm", mm)

	ApplyGoOverrides(target)

	if !mm.IsNativeFrozen() {
		t.Fatal("ApplyGoOverrides should freeze the registered native multifn")
	}
	if _, still := pendingNativeMultiFns[ns]; still {
		t.Fatalf("pending native multifns for %q should be drained after apply", ns)
	}
}

// When the namespace already exists at registration time (the lowered init runs
// after the bundle replay already finished), the multifn must be frozen
// immediately, matching the RegisterGoOverrides immediate-apply path.
func TestRegisterNativeMultiFnsFreezesImmediatelyWhenNSExists(t *testing.T) {
	const ns = "test-native-mm-immediate"
	delete(pendingNativeMultiFns, ns)
	target := NS(ns)
	mm := vm.NewMultiFn("mm", nil, vm.Keyword("default"))
	target.Def("mm", mm)

	RegisterNativeMultiFns(ns, []string{"mm"})

	if !mm.IsNativeFrozen() {
		t.Fatal("should freeze immediately when ns + multifn already exist")
	}
	if _, queued := pendingNativeMultiFns[ns]; queued {
		t.Fatalf("nothing should be queued when frozen immediately")
	}
}

// End-to-end runtime behavior of the generated guard: after freezing,
// MultiFnNativeFrozen reports true (native arms current); after a late
// defmethod replaces the var's root with an AddMethod result, it reports false
// (the generated switch then falls back to runtime dispatch). This is exactly
// the var → *MultiFn flow the lowered call site executes.
func TestMultiFnNativeFrozenGuardFlow(t *testing.T) {
	const ns = "test-mm-guard-flow"
	delete(pendingNativeMultiFns, ns)
	target := NS(ns)
	mm := vm.NewMultiFn("mm", nil, vm.Keyword("default"))
	target.Def("mm", mm)

	RegisterNativeMultiFns(ns, []string{"mm"}) // freezes immediately (ns exists)

	var vp *vm.Var
	if !MultiFnNativeFrozen(&vp, ns, "mm") {
		t.Fatal("guard must report frozen for the captured native baseline")
	}

	// Late defmethod: set! the multifn var to the AddMethod result. defmethod
	// expands to set!, which mutates the existing var's root (SetRoot) rather
	// than re-interning — so the memoized vp still resolves and now sees an
	// unfrozen MultiFn. (This is the same var-stability the default-arm
	// CachedVarFn trampoline already relies on.)
	v := target.LookupLocal(vm.Symbol("mm"))
	if v == nil {
		t.Fatal("multifn var should exist")
	}
	v.SetRoot(mm.AddMethod(vm.Keyword("x"), nil))

	if MultiFnNativeFrozen(&vp, ns, "mm") {
		t.Fatal("guard must report NOT frozen after a late defmethod")
	}
}
