/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/nooga/let-go/pkg/vm"
)

// nolint
func installOsNS() {
	getenv, err := vm.NativeFnType.Box(os.Getenv)
	execf, err := vm.NativeFnType.Box(exec.Command)
	tempDir, err := vm.NativeFnType.Box(os.TempDir)
	args, err := vm.ToLetGo(os.Args)
	if err != nil {
		panic(fmt.Sprintf("os NS init failed: %e", err))
	}

	ns := vm.NewNamespace("os")

	// vars
	CurrentNS = ns.Def("*ns*", ns)

	ns.Def("getenv", getenv)
	ns.Def("exec", execf)
	ns.Def("temp-dir", tempDir)
	ns.Def("args", args)
	RegisterNS(ns)
}
