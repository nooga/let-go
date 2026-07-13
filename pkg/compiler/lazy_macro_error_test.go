/*
 * Copyright (c) 2026 let-go contributors; see CONTRIBUTORS.
 * SPDX-License-Identifier: MIT
 */

package compiler

import (
	"errors"
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

func TestCompileConvertsLazyMacroExpansionThrowToError(t *testing.T) {
	boom := errors.New("boom during lazy macro expansion")
	thunk, err := vm.NativeFnType.Wrap(func(_ []vm.Value) (vm.Value, error) {
		return vm.NIL, boom
	})
	if err != nil {
		t.Fatal(err)
	}
	macroFn, err := vm.NativeFnType.Wrap(func(_ []vm.Value) (vm.Value, error) {
		return vm.NewLazySeq(thunk.(vm.Fn)), nil
	})
	if err != nil {
		t.Fatal(err)
	}

	ns := vm.NewNamespace("lazy-macro-error-test")
	macroVar := ns.Def("lazy-macro", macroFn)
	macroVar.SetMacro()
	compiler := NewCompiler(vm.NewConsts(), ns)

	_, err = compiler.Compile("(lazy-macro)")
	if !errors.Is(err, boom) {
		t.Fatalf("Compile error = %v, want %v", err, boom)
	}
}
