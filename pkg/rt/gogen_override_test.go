/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"fmt"
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

// Top-level def var-initializers are queued until ApplyGoVarInits runs them
// with an ec and Defs the produced roots. They run in registration (source)
// order; a forward reference (a later init reading an earlier init's var)
// resolves because phase 1 pre-interns every target var.
func TestApplyGoVarInitsRunsInOrderWithForwardRefs(t *testing.T) {
	const ns = "test-gogen-varinit"
	delete(pendingGoVarInits, ns)
	target := vm.NewNamespace(ns)

	// Each initializer Defs its own var (mirroring the lowered :def side
	// effect). A := 41 ; B := (deref A) + 1 — B reads A's var, set by A's init.
	RegisterGoVarInits(ns, []GoVarInit{
		{Name: "A", Init: func(_ *vm.ExecContext) (vm.Value, error) {
			v := vm.Int(41)
			return target.Def("A", v).Deref(), nil
		}},
		{Name: "B", Init: func(ec *vm.ExecContext) (vm.Value, error) {
			a := target.LookupLocal(vm.Symbol("A"))
			if a == nil {
				return nil, fmt.Errorf("A not interned before B runs")
			}
			v := vm.Int(int(a.Deref().(vm.Int)) + 1)
			return target.Def("B", v).Deref(), nil
		}},
	})

	if target.LookupLocal(vm.Symbol("A")) != nil {
		t.Fatalf("precondition: vars should not exist before ApplyGoVarInits")
	}

	ApplyGoVarInits(target, nil) // nil ec → RootExecContext

	if got := target.LookupLocal(vm.Symbol("A")); got == nil || got.Deref() != vm.Int(41) {
		t.Fatalf("A not set to 41: %v", got)
	}
	if got := target.LookupLocal(vm.Symbol("B")); got == nil || got.Deref() != vm.Int(42) {
		t.Fatalf("B not set to 42 (forward ref to A failed): %v", got)
	}
	if _, still := pendingGoVarInits[ns]; still {
		t.Fatalf("var-inits for %q should be drained after apply", ns)
	}
}

// An initializer that errors is reported and skipped; the others still apply.
func TestApplyGoVarInitsSkipsFailingInit(t *testing.T) {
	const ns = "test-gogen-varinit-err"
	delete(pendingGoVarInits, ns)
	target := vm.NewNamespace(ns)
	RegisterGoVarInits(ns, []GoVarInit{
		{Name: "good", Init: func(_ *vm.ExecContext) (vm.Value, error) {
			return target.Def("good", vm.Int(7)).Deref(), nil
		}},
		{Name: "bad", Init: func(_ *vm.ExecContext) (vm.Value, error) { return nil, fmt.Errorf("boom") }},
	})
	ApplyGoVarInits(target, nil)
	if got := target.LookupLocal(vm.Symbol("good")); got == nil || got.Deref() != vm.Int(7) {
		t.Fatalf("good var should be set despite sibling failure: %v", got)
	}
	// "bad" was pre-interned to NIL in phase 1 and left at NIL after its init failed.
	if got := target.LookupLocal(vm.Symbol("bad")); got == nil || got.Deref() != vm.NIL {
		t.Fatalf("bad var should remain NIL after failed init: %v", got)
	}
}
