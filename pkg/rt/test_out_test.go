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
