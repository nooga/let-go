/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/nooga/let-go/pkg/vm"
)

// Version and Commit are set by lg.go from goreleaser ldflags at startup.
var (
	Version = "dev"
	Commit  = "none"
)

func systemProperties() *vm.PersistentMap {
	home, uname := currentUser()
	if home == "" {
		home = os.Getenv("HOME")
	}
	if uname == "" {
		uname = os.Getenv("USER")
	}
	cwd, _ := os.Getwd()

	pairs := []vm.Value{
		vm.String("user.home"), vm.String(home),
		vm.String("user.name"), vm.String(uname),
		vm.String("user.dir"), vm.String(cwd),
		vm.String("java.io.tmpdir"), vm.String(os.TempDir()),
		vm.String("os.name"), vm.String(runtime.GOOS),
		vm.String("os.arch"), vm.String(runtime.GOARCH),
		vm.String("os.version"), vm.String(""),
		vm.String("file.separator"), vm.String(string(os.PathSeparator)),
		vm.String("path.separator"), vm.String(string(os.PathListSeparator)),
		vm.String("line.separator"), vm.String(lineSeparator()),
		vm.String("file.encoding"), vm.String("UTF-8"),
		vm.String("let-go.version"), vm.String(Version),
		vm.String("let-go.commit"), vm.String(Commit),
	}
	return vm.NewPersistentMap(pairs)
}

func init() { RegisterInstaller(installSystemNS) }

// nolint
func installSystemNS() {
	ns := vm.NewNamespace("System")

	// System/getProperty — (System/getProperty key) → string or nil
	ns.Def("getProperty", mustWrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) < 1 || len(vs) > 2 {
			return vm.NIL, fmt.Errorf("System/getProperty expects 1 or 2 args")
		}
		key, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("System/getProperty expected String key")
		}
		props := systemProperties()
		v := props.ValueAt(vs[0])
		if v == vm.NIL && len(vs) == 2 {
			return vs[1], nil
		}
		_ = key
		return v, nil
	}))

	// System/getProperties — (System/getProperties) → map
	ns.Def("getProperties", mustWrap(func(vs []vm.Value) (vm.Value, error) {
		return systemProperties(), nil
	}))

	// System/arraycopy(src sPos dst dPos len) — element copy between arrays.
	// malli.impl.util/-eager-entry-parser builds entry arrays with it.
	ns.Def("arraycopy", mustWrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 5 {
			return vm.NIL, fmt.Errorf("System/arraycopy expects 5 args")
		}
		src, ok1 := vs[0].(*vm.TypedArray)
		dst, ok2 := vs[2].(*vm.TypedArray)
		sPos, ok3 := vs[1].(vm.Int)
		dPos, ok4 := vs[3].(vm.Int)
		n, ok5 := vs[4].(vm.Int)
		if !(ok1 && ok2 && ok3 && ok4 && ok5) {
			return vm.NIL, fmt.Errorf("System/arraycopy: unsupported arg types")
		}
		if sPos < 0 || dPos < 0 || n < 0 ||
			int(sPos)+int(n) > src.Len() || int(dPos)+int(n) > dst.Len() {
			return vm.NIL, fmt.Errorf("System/arraycopy: index out of bounds")
		}
		// Read all source elements first so an overlapping same-array copy is
		// correct (mirrors java.lang.System.arraycopy).
		buf := make([]vm.Value, int(n))
		for i := 0; i < int(n); i++ {
			buf[i] = src.Get(int(sPos) + i)
		}
		for i := 0; i < int(n); i++ {
			if err := dst.Set(int(dPos)+i, buf[i]); err != nil {
				return vm.NIL, err
			}
		}
		return vm.NIL, nil
	}))

	// System/getenv — (System/getenv name) → string or nil
	ns.Def("getenv", mustWrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("System/getenv expects 1 arg")
		}
		k, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("System/getenv expected String name")
		}
		val, present := os.LookupEnv(string(k))
		if !present {
			return vm.NIL, nil
		}
		return vm.String(val), nil
	}))

	// System/exit — (System/exit code)
	ns.Def("exit", mustWrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("System/exit expects 1 arg")
		}
		code, ok := vs[0].(vm.Int)
		if !ok {
			return vm.NIL, fmt.Errorf("System/exit expected Int")
		}
		RunExitHooks() // flush profiling etc.; os.Exit skips deferred funcs
		os.Exit(int(code))
		return vm.NIL, nil
	}))

	// System/lineSeparator
	ns.Def("lineSeparator", mustWrap(func(vs []vm.Value) (vm.Value, error) {
		return vm.String(lineSeparator()), nil
	}))

	// System/currentTimeMillis
	ns.Def("currentTimeMillis", mustWrap(func(vs []vm.Value) (vm.Value, error) {
		return vm.MakeInt(int(time.Now().UnixMilli())), nil
	}))

	// System/nanoTime
	ns.Def("nanoTime", mustWrap(func(vs []vm.Value) (vm.Value, error) {
		return vm.MakeInt(int(time.Now().UnixNano())), nil
	}))

	RegisterNS(ns)
}
