//go:build runtime_only

/*
 * Copyright (c) 2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 *
 * Runtime-only entry point: executes precompiled bytecode (.lgb) with no reader,
 * compiler, or resolver linked in. Built with `-tags runtime_only`, which drops
 * pkg/compiler and pkg/resolver from the binary entirely (see pkg/rt/boot.go).
 * No eval / load-string / read-string, no dynamic source require — the deployed
 * artifact can run only bytecode compiled ahead of time by a trusted toolchain.
 *
 *   lg <program.lgb>   run a precompiled program
 *   lg                 boot the runtime and exit
 */

package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/nooga/let-go/pkg/bytecode"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

func fatal(prefix string, err error) {
	fmt.Fprintln(os.Stderr, prefix+":", err)
	os.Exit(1)
}

func runChunk(c *vm.CodeChunk) error {
	f := vm.NewFrame(c, nil)
	_, err := f.RunProtected()
	vm.ReleaseFrame(f)
	return err
}

func main() {
	if err := rt.LoadCore(); err != nil {
		fatal("boot", err)
	}
	rt.UseBytecodeNSLoader()

	if len(os.Args) < 2 {
		return // booted; nothing to run
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fatal("read", err)
	}
	resolve := func(nsName, name string) *vm.Var {
		n := rt.DefNSBare(nsName)
		if v := n.LookupLocal(vm.Symbol(name)); v != nil {
			return v
		}
		return n.DefStub(name)
	}
	unit, err := bytecode.DecodeToExecUnit(bytes.NewReader(data), resolve)
	if err != nil {
		fatal("decode", err)
	}
	// Replay the program's namespace chunks in dependency order, then main.
	for _, name := range unit.NSOrder {
		c := unit.NSChunks[name]
		if c == nil || c == unit.MainChunk {
			continue
		}
		if err := runChunk(c); err != nil {
			fatal("load "+name, err)
		}
	}
	if err := runChunk(unit.MainChunk); err != nil {
		fmt.Fprint(os.Stderr, vm.FormatError(err))
		os.Exit(1)
	}
}
