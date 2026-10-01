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

// Compiled code holds the Var it resolved at load time (a bundle's var
// references are decoded before any override applies), so an override must
// land on that Var, not on a replacement: through the queue and immediately.
func TestGoOverridesKeepTheExistingVar(t *testing.T) {
	for _, tc := range []struct {
		name    string
		install func(ns string, defs map[string]vm.Value)
	}{
		{"queued", func(ns string, defs map[string]vm.Value) {
			delete(pendingGoOverrides, ns)
			pendingGoOverrides[ns] = defs
			ApplyGoOverrides(LookupNS(ns))
		}},
		{"immediate", RegisterGoOverrides},
	} {
		t.Run(tc.name, func(t *testing.T) {
			nsName := "test-gogen-keep-var-" + tc.name
			ns := NS(nsName)
			held := ns.Def("g", vm.Int(1))
			fn := fnVal(t)

			tc.install(nsName, map[string]vm.Value{"g": fn})

			if held.Deref() != fn {
				t.Fatalf("a reference to the existing #'%s/g still sees %v after the override", nsName, held.Deref())
			}
			if got := ns.LookupLocal(vm.Symbol("g")); got != held {
				t.Fatalf("the override replaced #'%s/g with a new Var", nsName)
			}
		})
	}
}

// The core bundle decodes every namespace it holds as stub Vars before any
// lowered package registers, so an override lands on a stub and the
// namespace's chunk replays later, setting each Var's root to the bytecode fn.
// The resolver runs ApplyGoOverrides after that replay, which must put the
// override back on the same Var.
func TestGoOverridesSurviveTheNamespaceChunkReplay(t *testing.T) {
	const nsName = "test-gogen-replay"
	ns := NS(nsName)
	stub := ns.Def("g", vm.NIL)
	fn := fnVal(t)

	RegisterGoOverrides(nsName, map[string]vm.Value{"g": fn})
	stub.SetRoot(vm.Int(1)) // the chunk replay's def of g

	ApplyGoOverrides(ns)

	if stub.Deref() != fn {
		t.Fatalf("#'%s/g is %v after the chunk replay, want the override", nsName, stub.Deref())
	}
	if got := ns.LookupLocal(vm.Symbol("g")); got != stub {
		t.Fatalf("restoring the override replaced #'%s/g with a new Var", nsName)
	}
}
