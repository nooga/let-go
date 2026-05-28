//go:build !linux

/*
 * Copyright (c) 2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"fmt"

	"github.com/nooga/let-go/pkg/vm"
)

func init() { RegisterInstaller(installUnixNS) }

func installUnixNS() {
	unsupported := func(name string) vm.Value {
		fn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
			return vm.NIL, fmt.Errorf("unix/%s is only supported on Linux", name)
		})
		return fn
	}

	ns := vm.NewNamespace("unix")
	for _, name := range []string{"listen", "accept", "connect", "send", "recv", "close", "fd"} {
		ns.Def(name, unsupported(name))
	}
	RegisterNS(ns)
}
