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

package compiler

import (
	"github.com/nooga/let-go/pkg/vm"
	"strings"
)

type Context struct {
	ns     *vm.Namespace
	parent *Context
	consts []vm.Value
	chunk  *vm.CodeChunk
}

func (c *Context) Compile(s string) (*vm.CodeChunk, error) {
	r := NewLispReader(strings.NewReader(s), "<reader>")
	o, err := r.Read()
	if err != nil {
		return nil, err
	}

	c.chunk = vm.NewCodeChunk(&c.consts)
	c.compileForm(o)
	c.Emit(vm.OPRET)
	return c.chunk, nil
}

func (c *Context) Emit(op uint8) {
	c.chunk.Append(op)
}

func (c *Context) EmitWithArg(op uint8, arg int) {
	c.chunk.Append(op)
	c.chunk.Append32(arg)
}

func (c *Context) Constant(v vm.Value) int {
	for i := range c.consts {
		if c.consts[i] == v {
			return i
		}
	}
	c.consts = append(c.consts, v)
	return len(c.consts) - 1
}

func (c *Context) compileForm(o vm.Value) {
	switch o.Type() {
	case vm.IntType, vm.StringType, vm.NilType, vm.BooleanType:
		n := c.Constant(o)
		c.EmitWithArg(vm.OPLDC, n)
	case vm.SymbolType:
		n := c.Constant(c.ns.LookupOrAdd(o.(vm.Symbol)))
		c.EmitWithArg(vm.OPLDC, n)
	case vm.ListType:
		fn := o.(*vm.List).First()
		args := o.(*vm.List).Next().Unbox().([]vm.Value)

		for i := len(args) - 1; i >= 0; i-- {
			c.compileForm(args[i])
		}

		c.compileForm(fn)
		c.Emit(vm.OPINV)
	}
}
