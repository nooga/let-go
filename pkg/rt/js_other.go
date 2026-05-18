//go:build !(js && wasm)

/*
 * Copyright (c) 2026 Matt Parrett
 * SPDX-License-Identifier: MIT
 *
 * Native (non-WASM) stub for the `js` namespace. The namespace IS installed
 * so .lg code that uses (js/emit ...) parses and runs identically on native
 * and in the browser — emits just no-op on native, mirroring the way
 * `term/raw-mode!` is a no-op in WASM. Arg validation still runs so type
 * bugs are caught during native dev, not just at runtime in the browser.
 */

package rt

import (
	"github.com/nooga/let-go/pkg/vm"
)

func init() { RegisterInstaller(installJSNS) }

func installJSNS() {
	ns := vm.NewNamespace("js")

	// (js/emit event-name data) -> nil — validates args, then no-ops.
	emitFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if _, _, err := prepareEmit(vs); err != nil {
			return vm.NIL, err
		}
		return vm.NIL, nil
	})
	ns.Def("emit", emitFn)

	RegisterNS(ns)
}
