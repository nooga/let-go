//go:build js && wasm

/*
 * Copyright (c) 2026 Matt Parrett
 * SPDX-License-Identifier: MIT
 *
 * Browser bridge namespace. Lets .lg code dispatch structured events to JS
 * without caring whether the let-go VM is running on the main thread (where
 * the DOM is directly reachable) or inside a Web Worker (where it isn't).
 *
 * Single hidden contract with the bundle bootstrap: a global function
 *   _lgEmit(name string, dataJson string)
 * which dispatches a CustomEvent named `name` whose .detail is the parsed
 * dataJson. The bootstrap defines _lgEmit per mode — main-thread version
 * calls window.dispatchEvent directly; worker version postMessages to main,
 * which dispatches there. From .lg code, both look the same.
 */

package rt

import (
	"syscall/js"

	"github.com/nooga/let-go/pkg/vm"
)

func init() { RegisterInstaller(installJSNS) }

func installJSNS() {
	ns := vm.NewNamespace("js")

	// (js/emit event-name data) -> nil
	//
	// Fire-and-forget. Silently drops if _lgEmit isn't wired up (running
	// outside the official bundle, e.g. via a host that doesn't define it).
	emitFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		name, dataJSON, err := prepareEmit(vs)
		if err != nil {
			return vm.NIL, err
		}
		emit := js.Global().Get("_lgEmit")
		if emit.IsUndefined() {
			return vm.NIL, nil
		}
		emit.Invoke(name, dataJSON)
		return vm.NIL, nil
	})
	ns.Def("emit", emitFn)

	RegisterNS(ns)
}
