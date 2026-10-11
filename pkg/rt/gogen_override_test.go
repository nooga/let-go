/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

// fnVal builds a throwaway NativeFn value for assertions.
func fnVal(t *testing.T) vm.Value {
	t.Helper()
	v, err := vm.NativeFnType.Wrap(func(_ []vm.Value) (vm.Value, error) { return vm.NIL, nil })
	if err != nil {
		t.Fatalf("wrap: %v", err)
	}
	return v
}

// When the target namespace does NOT yet exist, RegisterGoOverrides must
// QUEUE the defs; ApplyGoOverrides drains them onto the ns once it is
// created. This is the bundle-replay-then-blank-import ordering.
func TestRegisterGoOverridesQueuesUntilApply(t *testing.T) {
	const ns = "test-gogen-queue"
	delete(pendingGoOverrides, ns) // isolate from other tests
	fn := fnVal(t)

	RegisterGoOverrides(ns, map[string]vm.Value{"frob": fn})

	if LookupNS(ns) != nil {
		t.Fatalf("precondition: ns %q should not exist yet", ns)
	}
	if pendingGoOverrides[ns]["frob"] != fn {
		t.Fatalf("override should be queued before the ns exists")
	}

	target := vm.NewNamespace(ns)
	ApplyGoOverrides(target)

	if got := target.LookupLocal(vm.Symbol("frob")); got == nil || got.Deref() != fn {
		t.Fatalf("ApplyGoOverrides did not Def the queued override onto the ns")
	}
	if _, still := pendingGoOverrides[ns]; still {
		t.Fatalf("pending overrides for %q should be drained after apply", ns)
	}
}

// The init-order hazard guard: if the namespace ALREADY exists when
// RegisterGoOverrides is called (the lowered package's init() runs AFTER
// postCoreInit already drained an empty queue), the defs must be applied
// IMMEDIATELY rather than queued forever.
func TestRegisterGoOverridesAppliesImmediatelyWhenNSExists(t *testing.T) {
	const ns = "test-gogen-immediate"
	delete(pendingGoOverrides, ns)
	existing := NS(ns) // create/register the namespace up front
	fn := fnVal(t)

	RegisterGoOverrides(ns, map[string]vm.Value{"baz": fn})

	if got := existing.LookupLocal(vm.Symbol("baz")); got == nil || got.Deref() != fn {
		t.Fatalf("override should be applied immediately when the ns already exists")
	}
	if _, queued := pendingGoOverrides[ns]; queued {
		t.Fatalf("nothing should be queued when applied immediately")
	}
}

// A lazily loaded bundled namespace is registered (DefNSBare) long before
// its chunk runs. Overrides registered in that window must wait for the
// load: applied to the bare namespace, they would be redefined by the
// chunk's own defs over its clojure.core refers, with a warning per name.
func TestRegisterGoOverridesQueuesWhileNSNeedsLoad(t *testing.T) {
	const ns = "test-gogen-needs-load"
	delete(pendingGoOverrides, ns)
	bare := DefNSBare(ns)
	MarkNSNeedsLoad(ns)
	t.Cleanup(func() { ClearNSNeedsLoad(ns) })
	fn := fnVal(t)

	RegisterGoOverrides(ns, map[string]vm.Value{"qux": fn})

	if got := bare.LookupLocal(vm.Symbol("qux")); got != nil {
		t.Fatalf("override was applied to a namespace whose chunk has not run")
	}
	if pendingGoOverrides[ns]["qux"] != fn {
		t.Fatalf("override should stay queued until the namespace loads")
	}

	ClearNSNeedsLoad(ns)
	ApplyGoOverrides(bare)
	if got := bare.LookupLocal(vm.Symbol("qux")); got == nil || got.Deref() != fn {
		t.Fatalf("ApplyGoOverrides after the load did not Def the queued override")
	}
}
