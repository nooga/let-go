/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

func TestRegisterTestOutUsesHostStdoutOrOutRoot(t *testing.T) {
	outVar := NS(NameCoreNS).LookupLocal(vm.Symbol("*out*"))
	if outVar == nil {
		t.Fatal("*out* not installed")
	}
	v := vm.NewVar(NS(NameCoreNS), "test", "*test-out*")
	v.SetDynamic()
	RegisterTestOut(v)
	if v.Deref() != outVar.Deref() {
		t.Fatalf("register-test-out before host install: got %v want *out* root %v", v.Deref(), outVar.Deref())
	}
	h := &IOHandle{name: "host"}
	InstallHostOutputRoots(h, h)
	if got := v.Deref().(*vm.Boxed).Unbox(); got != h {
		t.Fatalf("install-host-output-roots did not update *test-out*: %v", got)
	}
	RegisterTestOut(v)
	if got := v.Deref().(*vm.Boxed).Unbox(); got != h {
		t.Fatalf("register-test-out after host install should use host stdout: %v", got)
	}
}

func TestRegisterTestOutUsesOutRootNotDynamicBinding(t *testing.T) {
	outVar := NS(NameCoreNS).LookupLocal(vm.Symbol("*out*"))
	if outVar == nil {
		t.Fatal("*out* not installed")
	}
	// No host installed for this test: registeredHostStdout must be nil so
	// RegisterTestOut falls back to *out*'s root.
	saved := registeredHostStdout
	registeredHostStdout = nil
	defer func() { registeredHostStdout = saved }()

	wantRoot := outVar.Root()

	h2 := &IOHandle{name: "per-run-dynamic-binding"}
	outVar.PushBinding(vm.NewBoxed(h2))
	defer outVar.PopBinding()

	v := vm.NewVar(NS(NameCoreNS), "test", "*test-out*")
	v.SetDynamic()
	RegisterTestOut(v)

	got := v.Root()
	if got != wantRoot {
		t.Fatalf("register-test-out with a pushed *out* binding: got root %v, want *out*'s ROOT %v", got, wantRoot)
	}
	if b, ok := got.(*vm.Boxed); ok && b.Unbox() == h2 {
		t.Fatalf("register-test-out must not bake the dynamic top binding %v into *test-out*'s root", h2)
	}
}
