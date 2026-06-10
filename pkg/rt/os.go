/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"runtime"

	"github.com/nooga/let-go/pkg/vm"
)

func init() { RegisterInstaller(installOsNS) }

// nolint
func installOsNS() {
	getenv, err := vm.NativeFnType.Box(os.Getenv)
	execf, err := vm.NativeFnType.Box(exec.Command)
	tempDir, err := vm.NativeFnType.Box(os.TempDir)
	args, err := vm.ToLetGo(os.Args)
	withStdin, err := vm.NativeFnType.Wrap(func(v []vm.Value) (vm.Value, error) {
		var cmd = v[0].Unbox().(*exec.Cmd)
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

	// os/exec* — (os/exec* cmd & args) → exit code (Int). Unlike os/sh, the
	// child inherits the parent's stdin/stdout/stderr, so output streams and the
	// child stays interactive (e.g. launching a REPL). Returns the exit code.
	execStar, err := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 1 {
			return vm.NIL, fmt.Errorf("os/exec* expects at least 1 arg")
		}
		cmdName, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("os/exec* expected String command")
		}
		cmdArgs := make([]string, len(vs)-1)
		for i := 1; i < len(vs); i++ {
			if s, ok := vs[i].(vm.String); ok {
				cmdArgs[i-1] = string(s)
			} else {
				cmdArgs[i-1] = vs[i].String()
			}
		}
		cmd := exec.Command(string(cmdName), cmdArgs...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		exitCode := 0
		if runErr := cmd.Run(); runErr != nil {
			if exitErr, ok := runErr.(*exec.ExitError); ok {
				exitCode = exitErr.ExitCode()
			} else {
				return vm.NIL, runErr
			}
		}
		return vm.Int(exitCode), nil
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
	ns.Def("exec*", execStar)

	// os/free-port — (os/free-port) → an OS-assigned free TCP port (Int).
	// Binds 127.0.0.1:0, reads the assigned port, and releases the listener.
	// Check-then-use: the port could in principle be taken between this call
	// and the caller's own bind, so treat it as a strong hint rather than a
	// reservation (kernels avoid immediately reissuing a just-released port).
	ns.Def("free-port", mustWrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 0 {
			return vm.NIL, fmt.Errorf("os/free-port expects no args")
		}
		l, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			return vm.NIL, err
		}
		port := l.Addr().(*net.TCPAddr).Port
		if err := l.Close(); err != nil {
			return vm.NIL, err
		}
		return vm.Int(port), nil
	}))

	// os/os-name — (os/os-name) → "linux", "darwin", "windows", ...
	ns.Def("os-name", mustWrap(func(vs []vm.Value) (vm.Value, error) {
		return vm.String(runtime.GOOS), nil
	}))

	// os/arch — (os/arch) → "amd64", "arm64", ...
	ns.Def("arch", mustWrap(func(vs []vm.Value) (vm.Value, error) {
		return vm.String(runtime.GOARCH), nil
	}))

	// os/user-name — (os/user-name)
	ns.Def("user-name", mustWrap(func(vs []vm.Value) (vm.Value, error) {
		if _, name := currentUser(); name != "" {
			return vm.String(name), nil
		}
		return vm.String(os.Getenv("USER")), nil
	}))

	// os/hostname — (os/hostname)
	ns.Def("hostname", mustWrap(func(vs []vm.Value) (vm.Value, error) {
		h, err := os.Hostname()
		if err != nil {
			return vm.NIL, err
		}
		return vm.String(h), nil
	}))

	// os/file-separator, os/path-separator, os/line-separator
	ns.Def("file-separator", vm.String(string(os.PathSeparator)))
	ns.Def("path-separator", vm.String(string(os.PathListSeparator)))
	ns.Def("line-separator", vm.String(lineSeparator()))

	RegisterNS(ns)
}

func lineSeparator() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"
}

func mustWrap(fn func([]vm.Value) (vm.Value, error)) vm.Value {
	v, err := vm.NativeFnType.Wrap(fn)
	if err != nil {
		panic(err)
	}
	return v
}
