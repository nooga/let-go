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
	"fmt"
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
	err = c.compileForm(o)
	if err != nil {
		return nil, err
	}
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

func (c *Context) compileForm(o vm.Value) error {
	switch o.Type() {
	case vm.IntType, vm.StringType, vm.NilType, vm.BooleanType, vm.KeywordType, vm.CharType:
		n := c.Constant(o)
		c.EmitWithArg(vm.OPLDC, n)
	case vm.SymbolType:
		n := c.Constant(c.ns.LookupOrAdd(o.(vm.Symbol)))
		c.EmitWithArg(vm.OPLDC, n)
	//case vm.ArrayVectorType:
	//	v := o.(vm.ArrayVector)
	//	// FIXME detect const vectors and push them like this
	//	if len(v) == 0 {
	//		n := c.Constant(v)
	//		c.EmitWithArg(vm.OPLDC, n)
	//		return nil
	//	}
	//	for i := range v {
	//		err := c.compileForm(v[i])
	//		if err != nil {
	//			return NewCompileError("compiling vector members").Wrap(err)
	//		}
	//	}
	//  c.EmitWithArg(vm.OPVEC, len(v)) // FIXME this should not be a special instruction
	case vm.ListType:
		fn := o.(*vm.List).First()
		// check if we're looking at a special form
		if fn.Type() == vm.SymbolType {
			formCompiler, ok := specialForms[fn.(vm.Symbol)]
			if ok {
				return formCompiler(c, o)
			}
		}
		// treat as function invocation if this is not a special form
		args := o.(*vm.List).Next().Unbox().([]vm.Value)
		for i := len(args) - 1; i >= 0; i-- {
			err := c.compileForm(args[i])
			if err != nil {
				return NewCompileError("compiling arguments").Wrap(err)
			}
		}
		err := c.compileForm(fn)
		if err != nil {
			return NewCompileError("compiling function position").Wrap(err)
		}
		c.Emit(vm.OPINV)
	}
	return nil
}

func (c *Context) EmitWithArgPlaceholder(inst uint8) int {
	placeholder := c.CurrentAddress()
	c.EmitWithArg(inst, 0)
	return placeholder
}

func (c *Context) CurrentAddress() int {
	return c.chunk.Length()
}

func (c *Context) UpdatePlaceholderArg(placeholder int, arg int) {
	c.chunk.Update32(placeholder+1, arg)
}

type formCompilerFunc func(*Context, vm.Value) error

var specialForms map[vm.Symbol]formCompilerFunc

func init() {
	specialForms = map[vm.Symbol]formCompilerFunc{
		"if":  ifCompiler,
		"do":  doCompiler,
		"def": defCompiler,
	}
}

func ifCompiler(c *Context, form vm.Value) error {
	args := form.(*vm.List).Next().Unbox().([]vm.Value)
	l := len(args)
	if l < 2 || l > 3 {
		return NewCompileError(fmt.Sprintf("if: wrong number of forms (%d), need 2 or 3", l))
	}
	// compile condition
	err := c.compileForm(args[0])
	if err != nil {
		return NewCompileError("compiling if condition").Wrap(err)
	}
	elseJumpStart := c.EmitWithArgPlaceholder(vm.OPBRF)
	// compile then branch
	err = c.compileForm(args[1])
	if err != nil {
		return NewCompileError("compiling if then branch").Wrap(err)
	}
	finJumpStart := c.EmitWithArgPlaceholder(vm.OPJMP)
	elseJumpEnd := c.CurrentAddress()
	c.UpdatePlaceholderArg(elseJumpStart, elseJumpEnd-elseJumpStart)
	if l == 3 {
		err = c.compileForm(args[2])
		if err != nil {
			return NewCompileError("compiling if else branch").Wrap(err)
		}
	} else {
		c.EmitWithArg(vm.OPLDC, c.Constant(vm.NIL))
	}
	finJumpEnd := c.CurrentAddress()
	c.UpdatePlaceholderArg(finJumpStart, finJumpEnd-finJumpStart)
	return nil
}

func doCompiler(c *Context, form vm.Value) error {
	args := form.(*vm.List).Next().Unbox().([]vm.Value)
	l := len(args)
	if l == 0 {
		c.EmitWithArg(vm.OPLDC, c.Constant(vm.NIL))
		return nil
	}
	for i := range args {
		err := c.compileForm(args[i])
		if err != nil {
			return NewCompileError("compiling do member").Wrap(err)
		}
		if i < l-1 {
			c.Emit(vm.OPPOP)
		}
	}
	return nil
}

func defCompiler(c *Context, form vm.Value) error {
	args := form.(*vm.List).Next().Unbox().([]vm.Value)
	l := len(args)
	if l != 2 {
		return NewCompileError(fmt.Sprintf("def: wrong number of forms (%d), need 2", l))
	}
	sym := args[0]
	val := args[1]
	if sym.Type() != vm.SymbolType {
		return NewCompileError(fmt.Sprintf("def: first argument must be a symbol, got (%v)", sym))
	}
	varr := c.Constant(c.ns.LookupOrAdd(sym.(vm.Symbol)))
	c.EmitWithArg(vm.OPLDC, varr)

	err := c.compileForm(val)
	if err != nil {
		return NewCompileError("compiling def value").Wrap(err)
	}
	c.Emit(vm.OPSTV)

	return nil
}
