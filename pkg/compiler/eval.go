/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package compiler

import (
	"strings"

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
	// core is loaded eagerly – its macros are needed everywhere
	_, err := Eval(rt.CoreSrc)
	if err != nil {
		// Provide helpful debugging info for core.lg compilation failures
		panic("core.lg compilation failed: " + err.Error())
	}
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

	// test, walk, etc. are demand-loaded via resolver when required
}
