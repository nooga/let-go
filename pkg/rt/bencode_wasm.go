//go:build js && wasm

/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

/*
 * bencode namespace — browser stub. bencode framing runs over a net conn,
 * which the browser doesn't have (see net_wasm.go), so the operations can't
 * work here. Register the namespace with stubs that fail loudly.
 */

package rt

import (
	"fmt"

	"github.com/nooga/let-go/pkg/vm"
)

func init() { RegisterInstaller(installBencodeNS) }

func installBencodeNS() {
	unsupported := func(name string) vm.Value {
		fn, _ := vm.NativeFnType.Wrap(func(_ []vm.Value) (vm.Value, error) {
			return vm.NIL, fmt.Errorf("bencode/%s: not supported in the browser (requires a net conn)", name)
		})
		return fn
	}

	ns := vm.NewNamespace("bencode")
	for _, name := range []string{"write!", "read!"} {
		ns.Def(name, unsupported(name))
	}
	RegisterNS(ns)
}
