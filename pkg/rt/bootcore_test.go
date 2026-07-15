//go:build !bootstrap

/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

// BootCore must bind the .lg-DEFINED core surface, not just the native
// installers. `frequencies` is defined in clojure.core's .lg source, so it is
// bound only after the embedded core chunk runs — on a fresh ExecContext it is
// a nil stub, which is exactly the AOT-binary nil-panic BootCore exists to fix.
func TestBootCoreBindsLgDefinedCore(t *testing.T) {
	ec, err := BootCore()
	if err != nil {
		t.Fatalf("BootCore: %v", err)
	}
	if ec == nil {
		t.Fatal("BootCore returned a nil ExecContext")
	}

	v := LookupCoreVar("frequencies")
	if v == nil || !v.IsBound() {
		t.Fatal("clojure.core/frequencies unbound after BootCore")
	}

	// It actually runs: (frequencies ()) => {} with no error.
	if _, err := v.Invoke([]vm.Value{vm.EmptyList}); err != nil {
		t.Fatalf("invoke frequencies after BootCore: %v", err)
	}
}

// The eager chunk loops mutate the global *ns* root (in-ns with no ExecContext
// falls through to CurrentNS.SetRoot). Without a restore, the returned context
// would deref *ns* pointing at whichever hybrid chunk iterated last — a value
// that used to vary with map order between runs. BootCore must leave *ns* at
// its pre-boot root.
func TestBootCoreRestoresCurrentNS(t *testing.T) {
	before := CurrentNS.Deref()
	if _, err := BootCore(); err != nil {
		t.Fatalf("BootCore: %v", err)
	}
	if after := CurrentNS.Deref(); after != before {
		t.Fatalf("*ns* not restored: before=%v after=%v", before, after)
	}
}
