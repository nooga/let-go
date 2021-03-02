/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated
 * documentation files (the "Software"), to deal in the Software without restriction, including without limitation the
 * rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit
 * persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the
 * Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE
 * WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
 * COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
 * OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package rt

import (
	_ "embed"
	"fmt"
	"github.com/nooga/let-go/pkg/vm"
)

var nsRegistry map[string]*vm.Namespace

func init() {
	nsRegistry = make(map[string]*vm.Namespace)

	installLangNS()
}

func NS(name string) *vm.Namespace {
	return nsRegistry[name]
}

func RegisterNS(namespace *vm.Namespace) *vm.Namespace {
	nsRegistry[namespace.Name()] = namespace
	return namespace
}

//go:embed core/core.lg
var CoreSrc string

func installLangNS() {
	plus, err := vm.NativeFnType.Box(func(a int, b int) int { return a + b })
	mul, err := vm.NativeFnType.Box(func(a int, b int) int { return a * b })
	sub, err := vm.NativeFnType.Box(func(a int, b int) int { return a - b })

	equals, err := vm.NativeFnType.Wrap(func(vs []vm.Value) vm.Value {
		for i := 0; i < len(vs)-1; i++ {
			if vs[i] != vs[i+1] {
				return vm.FALSE
			}
		}
		return vm.TRUE
	})

	setMacro, err := vm.NativeFnType.Wrap(func(vs []vm.Value) vm.Value {
		m := vs[0].(*vm.Var)
		m.SetMacro()
		return m
	})

	vector, err := vm.NativeFnType.Wrap(vm.NewArrayVector)
	list, err := vm.NativeFnType.Wrap(vm.NewList)

	cons, err := vm.NativeFnType.Wrap(func(vs []vm.Value) vm.Value {
		elem := vs[0]
		seq, ok := vs[1].(vm.Seq)
		if !ok {
			// FIXME make this an error (we need to handle exceptions first)
			return vm.NIL
		}
		return seq.Cons(elem)
	})

	first, err := vm.NativeFnType.Wrap(func(vs []vm.Value) vm.Value {
		seq, ok := vs[0].(vm.Seq)
		if !ok {
			// FIXME make this an error (we need to handle exceptions first)
			return vm.NIL
		}
		return seq.First()
	})

	second, err := vm.NativeFnType.Wrap(func(vs []vm.Value) vm.Value {
		seq, ok := vs[0].(vm.Seq)
		if !ok {
			// FIXME make this an error (we need to handle exceptions first)
			return vm.NIL
		}
		return seq.Next().First()
	})

	next, err := vm.NativeFnType.Wrap(func(vs []vm.Value) vm.Value {
		seq, ok := vs[0].(vm.Seq)
		if !ok {
			// FIXME make this an error (we need to handle exceptions first)
			return vm.NIL
		}

		n := seq.Next()

		// FIXME move that to Seq.Next()
		if n.(vm.Collection).Count().(vm.Int) == 0 {
			return vm.NIL
		}
		return n
	})

	printlnf, err := vm.NativeFnType.Box(fmt.Println)
	if err != nil {
		panic("lang NS init failed")
	}

	ns := vm.NewNamespace("lang")
	ns.Def("+", plus)
	ns.Def("*", mul)
	ns.Def("-", sub)
	ns.Def("=", equals)

	ns.Def("set-macro!", setMacro)

	ns.Def("vector", vector)
	ns.Def("list", list)
	ns.Def("cons", cons)
	ns.Def("first", first)
	ns.Def("second", second)
	ns.Def("next", next)

	ns.Def("println", printlnf)

	RegisterNS(ns)
}
