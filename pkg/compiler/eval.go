/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package compiler

import (
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
	"strings"
)

func Eval(src string) (vm.Value, error) {
	ns := rt.NS(rt.NameCoreNS)
	compiler := NewCompiler(ns)

	_, out, err := compiler.CompileMultiple(strings.NewReader(src))
	if err != nil {
		return vm.NIL, err
	}

	//frame := vm.NewFrame(chunk, nil)
	//out, err := frame.Run()

	//chunk.Debug()
	//fmt.Println("eval: ", src, "=>", out)

	return out, nil
}

func evalInit() {
	_, err := Eval(rt.CoreSrc)
	if err != nil {
		panic(err)
	}
}
