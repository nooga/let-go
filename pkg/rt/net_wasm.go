//go:build js && wasm

/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

/*
 * net namespace — browser stub.
 *
 * Go's GOOS=js net stack is an in-process fake: net.Dial there resolves against
 * nothing real, so a browser build would "succeed" and then never talk to a
 * host. Register the namespace with stubs that fail loudly instead, matching
 * the native surface (net.go) name-for-name so a require doesn't break and the
 * error only fires if a program actually reaches for the network.
 */

package rt

import (
	"fmt"

	"github.com/nooga/let-go/pkg/vm"
)

func init() { RegisterInstaller(installNetNS) }

func installNetNS() {
	unsupported := func(name string) vm.Value {
		fn, _ := vm.NativeFnType.Wrap(func(_ []vm.Value) (vm.Value, error) {
			return vm.NIL, fmt.Errorf("net/%s: not supported in the browser (no host TCP)", name)
		})
		return fn
	}

	ns := vm.NewNamespace("net")
	for _, name := range []string{"dial", "write!", "read!", "close!"} {
		ns.Def(name, unsupported(name))
	}
	RegisterNS(ns)
}
