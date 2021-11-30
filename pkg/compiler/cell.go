/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
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
	c.scope.emitWithArg(vm.OP_DUP_NTH, c.scope.sp-1-c.local)
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
	c.scope.emitWithArg(vm.OP_LOAD_ARG, c.arg)
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
//	c.scope.emitWithArg(vm.OP_LOAD_ARG, c.arg)
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
	c.scope.emitWithArg(vm.OP_LOAD_CLOSEDOVER, c.closure)
	c.scope.incSP(1)
	return nil
}
