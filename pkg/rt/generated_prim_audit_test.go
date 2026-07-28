/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

// TestAuditGeneratedPrimitivesCleanAfterInit is the permanent regression guard
// for the alias-resolution bug (LookupOrRegisterNSNoLoad registering into a
// phantom "clojure.core" instead of canonical "core"): every generated
// primitive MUST be resolvable in its canonical namespace after init. This is
// also the gate to run after hoisting native closures to //lg:native decls —
// if any batch misregisters, this names each offender instead of leaving the
// bootstrap compile to die on the first `Can't resolve X`.
func TestAuditGeneratedPrimitivesCleanAfterInit(t *testing.T) {
	if probs := AuditGeneratedPrimitives(); len(probs) > 0 {
		t.Fatalf("generated primitives not all resolvable in their canonical namespace:\n  %s",
			strings.Join(probs, "\n  "))
	}
}

// TestAuditGeneratedPrimitivesFlagsMisregistration proves the audit is a
// non-terminating diagnostic: given a primitive recorded under a canonical ns
// but actually bound in a different ns (the exact shape of the alias bug), it
// reports the miss AND points at where the name actually landed.
func TestAuditGeneratedPrimitivesFlagsMisregistration(t *testing.T) {
	const probe = "audit-probe-should-not-collide"
	adapter, err := vm.NativeFnType.Wrap(func(_ []vm.Value) (vm.Value, error) { return vm.NIL, nil })
	if err != nil {
		t.Fatalf("wrap adapter: %v", err)
	}

	// Bind the probe in a phantom namespace (simulating a prim that landed in
	// the wrong ns), and register it so whereNameBound can locate it.
	phantom := vm.NewNamespace("clojure.audit.phantom")
	phantom.Def(probe, adapter)
	RegisterNS(phantom)

	// Record it under canonical "core", which does NOT bind the probe.
	genPrimMu.Lock()
	if genPrimBindings[NameCoreNS] == nil {
		genPrimBindings[NameCoreNS] = map[string]vm.Value{}
	}
	genPrimBindings[NameCoreNS][probe] = adapter
	genPrimMu.Unlock()
	defer func() {
		genPrimMu.Lock()
		delete(genPrimBindings[NameCoreNS], probe)
		genPrimMu.Unlock()
	}()

	var line string
	for _, p := range AuditGeneratedPrimitives() {
		if strings.Contains(p, probe) {
			line = p
		}
	}
	if line == "" {
		t.Fatalf("audit did not flag misregistered %q", probe)
	}
	if !strings.Contains(line, "not bound in canonical ns") {
		t.Errorf("audit line missing reason: %s", line)
	}
	if !strings.Contains(line, "clojure.audit.phantom") {
		t.Errorf("audit line did not point at the phantom ns: %s", line)
	}
}
