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
