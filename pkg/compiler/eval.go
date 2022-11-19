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

func evalInit() {
	consts = vm.NewConsts()
	_, err := Eval(rt.CoreSrc)
	if err != nil {
		panic(err)
	}
}
