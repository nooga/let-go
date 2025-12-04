/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"fmt"
	"io"
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
	withStdin, err := vm.NativeFnType.Wrap(func(v []vm.Value) (vm.Value, error) {
		var cmd *exec.Cmd = v[0].Unbox().(*exec.Cmd)
		s := string(v[1].(vm.String))
		stdin, err := cmd.StdinPipe()
		if err != nil {
			return vm.NIL, err
		}
		go func() {
			defer stdin.Close()
			io.WriteString(stdin, s)
		}()
		return v[0], nil
	})
	if err != nil {
		panic(fmt.Sprintf("os NS init failed: %e", err))
	}

	ns := vm.NewNamespace("os")

	ns.Def("getenv", getenv)
	ns.Def("exec", execf)
	ns.Def("with-stdin", withStdin)
	ns.Def("temp-dir", tempDir)
	ns.Def("args", args)
	RegisterNS(ns)
}
