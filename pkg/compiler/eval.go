/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package compiler

import (
	"bytes"
	"strings"

	"github.com/nooga/let-go/pkg/bytecode"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

var consts *vm.Consts

func Eval(src string) (vm.Value, error) {
	ns := rt.NS(rt.NameCoreNS)
	compiler := NewCompiler(consts, ns)

	_, out, err := compiler.CompileMultiple(strings.NewReader(src))
	if err != nil {
		return vm.NIL, err
	}

	return out, nil
}

// ReadString parses a string into a let-go Value.
func ReadString(s string) (vm.Value, error) {
	reader := NewLispReader(strings.NewReader(s), "<read-string>")
	return reader.Read()
}

func evalInit() {
	consts = vm.NewConsts()

	// Try loading pre-compiled core bytecode
	if len(rt.CoreCompiledLGB) > 0 {
		if err := loadPrecompiledCore(); err == nil {
			postCoreInit()
			return
		}
		// Fall through to source compilation on error
	}

	// Original path: compile from source
	_, err := Eval(rt.CoreSrc)
	if err != nil {
		panic("core.lg compilation failed: " + err.Error())
	}
	postCoreInit()
}

func loadPrecompiledCore() error {
	resolve := func(ns, name string) *vm.Var {
		n := rt.NS(ns)
		v := n.Lookup(vm.Symbol(name))
		if v == vm.NIL {
			return n.Def(name, vm.NIL)
		}
		return v.(*vm.Var)
	}
	unit, err := bytecode.DecodeToExecUnit(bytes.NewReader(rt.CoreCompiledLGB), resolve)
	if err != nil {
		return err
	}

	// Use the decoded const pool as the global pool
	consts = unit.Consts

	// Execute the main chunk to replay all def/defn/defmacro definitions
	f := vm.NewFrame(unit.MainChunk, nil)
	_, err = f.RunProtected()
	vm.ReleaseFrame(f)
	return err
}

func postCoreInit() {
	// Register read-string (needs the reader which lives in the compiler package)
	readStringFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, nil
		}
		s, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, nil
		}
		return ReadString(string(s))
	})
	rt.NS(rt.NameCoreNS).Def("read-string", readStringFn)

	// Wire up EDN reader for pod support
	rt.SetReadEDN(func(s string) (vm.Value, error) {
		return ReadString(s)
	})

	// Wire up namespace-aware eval for pod client-side code
	rt.SetEvalInNS(func(code string, ns *vm.Namespace) (vm.Value, error) {
		c := NewCompiler(consts, ns)
		_, out, err := c.CompileMultiple(strings.NewReader(code))
		return out, err
	})

	// test, walk, etc. are demand-loaded via resolver when required
}
