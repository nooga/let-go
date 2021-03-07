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
)

type cell interface {
	source() cell
	emit() error
}

type localCell struct {
	scope *Context
	local int
}

func (c *localCell) source() cell {
	return nil
}

func (c *localCell) emit() error {
	c.scope.emitWithArg(vm.OPDPN, c.scope.sp-1-c.local)
	c.scope.incSP(1)
	return nil
}

type argCell struct {
	scope *Context
	arg   int
}

func (c *argCell) source() cell {
	return nil
}

func (c *argCell) emit() error {
	c.scope.emitWithArg(vm.OPLDA, c.arg)
	c.scope.incSP(1)
	return nil
}

// might come in handy later

//type varCell struct {
//	scope *Context
//	arg   int
//}
//
//func (c *varCell) source() cell {
//	return nil
//}
//
//func (c *varCell) emit() error {
//	c.scope.emitWithArg(vm.OPLDA, c.arg)
//	c.scope.incSP(1)
//	return nil
//}

type closureCell struct {
	src     cell
	scope   *Context
	closure int
}

func (c *closureCell) source() cell {
	return c.src
}

func (c *closureCell) emit() error {
	c.scope.emitWithArg(vm.OPLDK, c.closure)
	c.scope.incSP(1)
	return nil
}
