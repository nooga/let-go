/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"bytes"
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

	// os/exit — (os/exit code)
	exitf, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("os/exit expects 1 arg")
		}
		code, ok := vs[0].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("os/exit expected Int")
		}
		os.Exit(int(code))
		return vm.NIL, nil
	})

	// os/cwd — (os/cwd)
	cwd, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		d, err := os.Getwd()
		if err != nil {
			return vm.NIL, err
		}
		return vm.String(d), nil
	})

	// os/setenv — (os/setenv key val)
	setenv, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 2 {
			return vm.NIL, fmt.Errorf("os/setenv expects 2 args")
		}
		k, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("os/setenv expected String key")
		}
		v, ok := vs[1].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("os/setenv expected String value")
		}
		return vm.NIL, os.Setenv(string(k), string(v))
	})

	// os/ls — (os/ls path)
	ls, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("os/ls expects 1 arg")
		}
		path, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("os/ls expected String path")
		}
		entries, err := os.ReadDir(string(path))
		if err != nil {
			return vm.NIL, err
		}
		result := make([]vm.Value, len(entries))
		for i, e := range entries {
			result[i] = vm.String(e.Name())
		}
		return vm.NewArrayVector(result), nil
	})

	// os/stat — (os/stat path)
	stat, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("os/stat expects 1 arg")
		}
		path, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("os/stat expected String path")
		}
		info, err := os.Stat(string(path))
		if err != nil {
			if os.IsNotExist(err) {
				return vm.NIL, nil
			}
			return vm.NIL, err
		}
		return fileStatMapping.StructToRecord(FileStat{
			Name:    info.Name(),
			Size:    info.Size(),
			IsDir:   info.IsDir(),
			ModTime: info.ModTime().String(),
		}), nil
	})

	// os/sh — (os/sh cmd & args) → {:exit 0 :out "..." :err "..."}
	sh, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 1 {
			return vm.NIL, fmt.Errorf("os/sh expects at least 1 arg")
		}
		cmdName, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("os/sh expected String command")
		}
		args := make([]string, len(vs)-1)
		for i := 1; i < len(vs); i++ {
			if s, ok := vs[i].(vm.String); ok {
				args[i-1] = string(s)
			} else {
				args[i-1] = vs[i].String()
			}
		}
		cmd := exec.Command(string(cmdName), args...)
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err := cmd.Run()
		exitCode := 0
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				exitCode = exitErr.ExitCode()
			} else {
				return vm.NIL, err
			}
		}
		return shellResultMapping.StructToRecord(ShellResult{
			Exit: exitCode,
			Out:  stdout.String(),
			Err:  stderr.String(),
		}), nil
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
	ns.Def("exit", exitf)
	ns.Def("cwd", cwd)
	ns.Def("setenv", setenv)
	ns.Def("ls", ls)
	ns.Def("stat", stat)
	ns.Def("sh", sh)
	RegisterNS(ns)
}
