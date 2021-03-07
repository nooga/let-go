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
	"io"
	"strings"
)

type Context struct {
	ns           *vm.Namespace
	parent       *Context
	consts       *[]vm.Value
	chunk        *vm.CodeChunk
	formalArgs   map[vm.Symbol]int
	source       string
	variadric    bool
	locals       []map[vm.Symbol]int
	sp           int
	spMax        int
	isFunction   bool
	isClosure    bool
	closedOversC int
	closedOvers  map[vm.Symbol]*closureCell
	recurPoints  []*recurPoint
	tailPosition bool
}

// FIXME this is unacceptable hax
var globalConsts *[]vm.Value

func init() {
	globalConsts = &[]vm.Value{}
}

func NewCompiler(ns *vm.Namespace) *Context {
	return &Context{
		ns:          ns,
		consts:      globalConsts,
		source:      "<default>",
		locals:      []map[vm.Symbol]int{},
		closedOvers: map[vm.Symbol]*closureCell{},
	}
}

func (c *Context) SetSource(source string) *Context {
	c.source = source
	return c
}

func (c *Context) CurrentNS() *vm.Namespace {
	return c.ns
}

func (c *Context) Compile(s string) (*vm.CodeChunk, error) {
	r := NewLispReader(strings.NewReader(s), c.source)
	o, err := r.Read()
	if err != nil {
		return nil, err
	}

	c.chunk = vm.NewCodeChunk(c.consts)
	err = c.compileForm(o)
	c.chunk.SetMaxStack(c.spMax)
	if err != nil {
		return nil, err
	}
	c.emit(vm.OPRET)
	c.decSP(1)
	return c.chunk, nil
}

func (c *Context) CompileMultiple(reader io.Reader) (*vm.CodeChunk, vm.Value, error) {
	r := NewLispReader(reader, c.source)
	chunk := vm.NewCodeChunk(c.consts)
	var result vm.Value = vm.NIL
	compiledForms := 0
	for {
		o, err := r.Read()
		if err != nil {
			if isErrorEOF(err) {
				break
			}
			return nil, result, err
		}
		if compiledForms > 0 {
			chunk.Append(vm.OPPOP)
		}
		formchunk := vm.NewCodeChunk(c.consts)
		c.chunk = formchunk
		err = c.compileForm(o)
		c.chunk.SetMaxStack(c.spMax)
		if err != nil {
			return nil, result, err
		}
		chunk.AppendChunk(formchunk)

		formchunk.Append(vm.OPRET)
		f := vm.NewFrame(formchunk, nil)
		result, err = f.Run()
		if err != nil {
			return nil, result, err
		}
		compiledForms++
	}

	c.chunk = chunk

	c.emit(vm.OPRET)
	c.decSP(1)
	return c.chunk, result, nil
}

func (c *Context) emit(op uint8) {
	c.chunk.Append(op)
}

func (c *Context) emitWithArg(op uint8, arg int) {
	c.chunk.Append(op)
	c.chunk.Append32(arg)
}

func (c *Context) constant(v vm.Value) int {
	for i := range *c.consts {
		if (*c.consts)[i] == v {
			return i
		}
	}
	*c.consts = append(*c.consts, v)
	return len(*c.consts) - 1
}

func (c *Context) arg(v vm.Symbol) int {
	n, ok := c.formalArgs[v]
	if !ok {
		return -1
	}
	return n
}

func (c *Context) enterFn(args []vm.Value) (*Context, error) {
	fchunk := vm.NewCodeChunk(c.consts)

	fc := &Context{
		ns:           c.ns,
		parent:       c,
		consts:       c.consts,
		chunk:        fchunk,
		formalArgs:   make(map[vm.Symbol]int),
		locals:       []map[vm.Symbol]int{},
		closedOvers:  make(map[vm.Symbol]*closureCell),
		isFunction:   true,
		tailPosition: true,
	}

	for i := range args {
		a := args[i]
		s, ok := a.(vm.Symbol)
		if !ok {
			return nil, NewCompileError("all fn formal arguments must be symbols")
		}
		if s == "&" {
			if fc.variadric {
				return nil, NewCompileError("only one rest argument allowed")
			}
			fc.variadric = true
			continue
		}
		if fc.variadric {
			if i < len(args)-1 {
				return nil, NewCompileError("only one argument allowed after &")
			}
			i = i - 1
		}
		fc.formalArgs[s] = i
	}
	return fc, nil
}

func (c *Context) leaveFn(ctx *Context) {
	fnchunk := ctx.chunk
	fnchunk.SetMaxStack(ctx.spMax)
	f := vm.MakeFunc(len(ctx.formalArgs), ctx.variadric, fnchunk)

	n := c.constant(f)
	c.emitWithArg(vm.OPLDC, n)
	c.incSP(1)

	// if we have a closure on our hands then add closed overs
	if ctx.isClosure {
		for _, clo := range ctx.closedOvers {
			_ = clo.source().emit()
			c.emit(vm.OPPAK)
		}
	}
}

func (c *Context) symbolLookup(s vm.Symbol) cell {
	if c.isClosure {
		clo := c.closedOvers[s]
		if clo != nil {
			return clo
		}
	}
	local := c.lookupLocal(s)
	if local >= 0 {
		// we have a local symbol in scope
		return &localCell{
			scope: c,
			local: local,
		}
	}
	arg := c.arg(s)
	if arg >= 0 {
		return &argCell{
			scope: c,
			arg:   arg,
		}
	}
	if c.parent == nil {
		return nil
	}
	outer := c.parent.symbolLookup(s)
	if outer != nil {
		c.isClosure = true
		newClosedOver := c.closedOversC
		c.closedOversC++
		c.closedOvers[s] = &closureCell{
			src:     outer,
			scope:   c,
			closure: newClosedOver,
		}
		return c.closedOvers[s]
	}
	return nil
}

func (c *Context) compileForm(o vm.Value) error {
	switch o.Type() {
	case vm.IntType, vm.StringType, vm.NilType, vm.BooleanType, vm.KeywordType, vm.CharType, vm.VoidType:
		n := c.constant(o)
		c.emitWithArg(vm.OPLDC, n)
		c.incSP(1)
	case vm.SymbolType:
		cel := c.symbolLookup(o.(vm.Symbol))
		if cel != nil {
			return cel.emit()
		}
		// if symbol not found so far then we have a free variable on our hands
		varn := c.constant(c.ns.LookupOrAdd(o.(vm.Symbol)))
		c.emitWithArg(vm.OPLDC, varn)
		c.emit(vm.OPLDV)
		c.incSP(1)
	case vm.ArrayVectorType:
		tp := c.tailPosition
		c.tailPosition = false
		v := o.(vm.ArrayVector)
		// FIXME detect const vectors and push them like this
		if len(v) == 0 {
			n := c.constant(v)
			c.emitWithArg(vm.OPLDC, n)
			c.incSP(1)
			return nil
		}
		vector := c.constant(c.ns.LookupOrAdd("vector"))
		c.emitWithArg(vm.OPLDC, vector)
		c.incSP(1)
		for i := range v {
			err := c.compileForm(v[i])
			if err != nil {
				return NewCompileError("compiling vector elements").Wrap(err)
			}
		}
		c.emitWithArg(vm.OPINV, len(v))
		c.decSP(len(v) + 1)
		c.tailPosition = tp
	case vm.ListType:
		fn := o.(*vm.List).First()
		// check if we're looking at a special form
		if fn.Type() == vm.SymbolType {
			formCompiler, ok := specialForms[fn.(vm.Symbol)]
			if ok {
				return formCompiler(c, o)
			}

			fvar, ok := c.ns.Lookup(fn.(vm.Symbol)).(*vm.Var)
			if ok && fvar.IsMacro() {
				argvec := o.(*vm.List).Next().(*vm.List).Unbox().([]vm.Value)
				newform := fvar.Invoke(argvec)
				return c.compileForm(newform)
			}
		}

		tp := c.tailPosition
		c.tailPosition = false

		// treat as function invocation if this is not a special form
		err := c.compileForm(fn)
		if err != nil {
			return NewCompileError("compiling function position").Wrap(err)
		}

		args := o.(*vm.List).Next()
		argc := args.(vm.Collection).Count().Unbox().(int)
		for args != vm.EmptyList {
			err := c.compileForm(args.First())
			if err != nil {
				return NewCompileError("compiling arguments").Wrap(err)
			}
			args = args.Next()
		}

		c.emitWithArg(vm.OPINV, argc)
		c.decSP(argc)

		c.tailPosition = tp
	}
	return nil
}

func (c *Context) emitWithArgPlaceholder(inst uint8) int {
	placeholder := c.currentAddress()
	c.emitWithArg(inst, 0)
	return placeholder
}

func (c *Context) currentAddress() int {
	return c.chunk.Length()
}

func (c *Context) updatePlaceholderArg(placeholder int, arg int) {
	c.chunk.Update32(placeholder+1, arg)
}

func (c *Context) pushLocals() {
	c.locals = append(c.locals, map[vm.Symbol]int{})
}

func (c *Context) popLocals() {
	c.locals = c.locals[0 : len(c.locals)-1]
}

func (c *Context) addLocal(name vm.Symbol) {
	c.locals[len(c.locals)-1][name] = c.sp - 1
}

func (c *Context) incSP(i int) {
	c.sp += i
	if c.sp > c.spMax {
		c.spMax = c.sp
	}
}

func (c *Context) decSP(i int) {
	c.sp -= i
}

func (c *Context) lookupLocal(symbol vm.Symbol) int {
	if len(c.locals) < 1 {
		return -1
	}
	for i := len(c.locals) - 1; i >= 0; i-- {
		local, ok := c.locals[i][symbol]
		if ok {
			return local
		}
	}
	return -1
}

type recurPoint struct {
	address int
	argsc   int
}

func (c *Context) pushRecurPoint(argsc int) {
	c.recurPoints = append(c.recurPoints, &recurPoint{
		address: c.currentAddress(),
		argsc:   argsc,
	})
}

func (c *Context) popRecurPoint() {
	if len(c.recurPoints) > 0 {
		c.recurPoints = c.recurPoints[:len(c.recurPoints)-1]
	}
}

func (c *Context) currentRecurPoint() *recurPoint {
	if len(c.recurPoints) > 0 {
		return c.recurPoints[len(c.recurPoints)-1]
	}
	return nil
}

type formCompilerFunc func(*Context, vm.Value) error

var specialForms map[vm.Symbol]formCompilerFunc

func compilerInit() {
	specialForms = map[vm.Symbol]formCompilerFunc{
		"if":    ifCompiler,
		"do":    doCompiler,
		"def":   defCompiler,
		"fn":    fnCompiler,
		"quote": quoteCompiler,
		"var":   varCompiler,
		"let":   letCompiler,
		"loop":  loopCompiler,
		"recur": recurCompiler,
	}
}

func recurCompiler(c *Context, form vm.Value) error {
	if !c.tailPosition {
		return NewCompileError("recur is only allowed in tail position")
	}
	rp := c.currentRecurPoint()

	tp := c.tailPosition
	c.tailPosition = false

	args := form.(*vm.List).Next()
	argc := args.(vm.Collection).Count().Unbox().(int)

	if rp != nil {
		if argc != rp.argsc {
			return NewCompileError("recur argument count must match loop bindings count")
		}
	} else {
		if !c.isFunction {
			return NewCompileError("recur is only allowed inside loops and functions")
		}
		if argc != len(c.formalArgs) {
			return NewCompileError("recur argument count must match function argument count")
		}
	}

	for args != vm.EmptyList {
		err := c.compileForm(args.First())
		if err != nil {
			return NewCompileError("compiling recur arguments").Wrap(err)
		}
		args = args.Next()
	}

	if !c.isFunction {
		c.emitWithArg(vm.OPREC, c.currentAddress()-rp.address)
		c.chunk.Append32(argc)
	} else {
		c.emitWithArg(vm.OPREF, argc)
	}
	c.tailPosition = tp
	c.decSP(argc - 1) // this is needed to keep the balance of if branches
	return nil
}

func loopCompiler(c *Context, form vm.Value) error {
	bindings := form.(*vm.List).Next()
	binds, ok := bindings.First().(vm.ArrayVector)
	if !ok {
		return NewCompileError("let bindings should be a vector")
	}
	body := bindings.Next()
	c.pushLocals()
	tp := c.tailPosition
	c.tailPosition = false
	bindn := 0
	for i := 0; i < len(binds); i += 2 {
		name := binds[i]
		if name.Type() != vm.SymbolType {
			return NewCompileError("loop binding name must be a symbol")
		}
		if i+1 >= len(binds) {
			return NewCompileError("loop bindings must have even number of forms")
		}
		value := binds[i+1]
		err := c.compileForm(value)
		if err != nil {
			return NewCompileError("compiling let binding").Wrap(err)
		}
		c.addLocal(name.(vm.Symbol))
		bindn++
	}
	c.pushRecurPoint(bindn)
	if body == vm.EmptyList {
		c.emitWithArg(vm.OPLDC, c.constant(vm.NIL))
		c.incSP(1)
	} else {
		for b := body; b != vm.EmptyList; b = b.Next() {
			if b.Next() == vm.EmptyList {
				c.tailPosition = true
			}
			err := c.compileForm(b.First())
			if err != nil {
				return NewCompileError("compiling let body").Wrap(err)
			}
			if b.Next() != vm.EmptyList {
				c.emit(vm.OPPOP)
				c.decSP(1)
			}
		}
	}
	c.popLocals()
	c.popRecurPoint()
	c.emitWithArg(vm.OPPON, bindn)
	c.decSP(bindn)
	c.tailPosition = tp
	return nil
}

func letCompiler(c *Context, form vm.Value) error {
	bindings := form.(*vm.List).Next()
	binds, ok := bindings.First().(vm.ArrayVector)
	if !ok {
		return NewCompileError("let bindings should be a vector")
	}
	body := bindings.Next()
	c.pushLocals()
	tc := c.tailPosition
	c.tailPosition = false
	bindn := 0
	for i := 0; i < len(binds); i += 2 {
		name := binds[i]
		if name.Type() != vm.SymbolType {
			return NewCompileError("let binding name must be a symbol")
		}
		if i+1 >= len(binds) {
			return NewCompileError("let bindings must have even number of forms")
		}
		value := binds[i+1]
		err := c.compileForm(value)
		if err != nil {
			return NewCompileError("compiling let binding").Wrap(err)
		}
		c.addLocal(name.(vm.Symbol))
		bindn++
	}
	if body == vm.EmptyList {
		c.emitWithArg(vm.OPLDC, c.constant(vm.NIL))
		c.incSP(1)
	} else {
		for b := body; b != vm.EmptyList; b = b.Next() {
			if tc && b.Next() == vm.EmptyList {
				c.tailPosition = true
			}
			err := c.compileForm(b.First())
			if err != nil {
				return NewCompileError("compiling let body").Wrap(err)
			}
			if b.Next() != vm.EmptyList {
				c.emit(vm.OPPOP)
				c.decSP(1)
			}
		}
	}
	c.popLocals()
	c.emitWithArg(vm.OPPON, bindn)
	c.decSP(bindn)
	c.tailPosition = tc
	return nil
}

func quoteCompiler(c *Context, form vm.Value) error {
	n := c.constant(form.(vm.Seq).Next().First())
	c.emitWithArg(vm.OPLDC, n)
	c.incSP(1)
	return nil
}

func fnCompiler(c *Context, form vm.Value) error {
	f := form.(*vm.List).Next()

	args := f.First().(vm.ArrayVector).Unbox().([]vm.Value)

	fc, err := c.enterFn(args)
	if err != nil {
		return NewCompileError("compiling fn args").Wrap(err)
	}
	defer c.leaveFn(fc)

	body := f.(*vm.List).Next().Unbox().([]vm.Value)
	l := len(body)
	if l == 0 {
		fc.emitWithArg(vm.OPLDC, fc.constant(vm.NIL))
		fc.incSP(1)
		fc.emit(vm.OPRET)
		return nil
	}
	for i := range body {
		err := fc.compileForm(body[i])
		if err != nil {
			return NewCompileError("compiling do member").Wrap(err)
		}
		if i < l-1 {
			fc.emit(vm.OPPOP)
			fc.decSP(1)
		}
	}
	fc.emit(vm.OPRET)

	return nil
}

func ifCompiler(c *Context, form vm.Value) error {
	tc := c.tailPosition
	c.tailPosition = tc

	args := form.(*vm.List).Next().Unbox().([]vm.Value)
	l := len(args)
	if l < 2 || l > 3 {
		return NewCompileError(fmt.Sprintf("if: wrong number of forms (%d), need 2 or 3", l))
	}
	c.tailPosition = true
	// compile condition
	err := c.compileForm(args[0])
	if err != nil {
		return NewCompileError("compiling if condition").Wrap(err)
	}
	elseJumpStart := c.emitWithArgPlaceholder(vm.OPBRF)
	c.decSP(1)
	c.tailPosition = tc

	// compile then branch
	err = c.compileForm(args[1])
	c.decSP(1)

	if err != nil {
		return NewCompileError("compiling if then branch").Wrap(err)
	}
	finJumpStart := c.emitWithArgPlaceholder(vm.OPJMP)
	elseJumpEnd := c.currentAddress()
	c.updatePlaceholderArg(elseJumpStart, elseJumpEnd-elseJumpStart)
	if l == 3 {
		err = c.compileForm(args[2])

		if err != nil {
			return NewCompileError("compiling if else branch").Wrap(err)
		}
	} else {
		c.emitWithArg(vm.OPLDC, c.constant(vm.NIL))
		c.incSP(1)
	}
	finJumpEnd := c.currentAddress()
	c.updatePlaceholderArg(finJumpStart, finJumpEnd-finJumpStart)
	return nil
}

func doCompiler(c *Context, form vm.Value) error {
	args := form.(*vm.List).Next().Unbox().([]vm.Value)
	l := len(args)
	tc := c.tailPosition
	c.tailPosition = false
	if l == 0 {
		c.emitWithArg(vm.OPLDC, c.constant(vm.NIL))
		c.incSP(1)
		c.tailPosition = tc
		return nil
	}
	for i := range args {
		if tc && i == l-1 {
			c.tailPosition = true
		}
		err := c.compileForm(args[i])
		if err != nil {
			return NewCompileError("compiling do member").Wrap(err)
		}
		if i < l-1 {
			c.emit(vm.OPPOP)
			c.decSP(1)
		}
	}
	c.tailPosition = tc
	return nil
}

func defCompiler(c *Context, form vm.Value) error {
	tc := c.tailPosition
	c.tailPosition = false
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
	varr := c.constant(c.ns.LookupOrAdd(sym.(vm.Symbol)))
	c.emitWithArg(vm.OPLDC, varr)
	c.incSP(1)
	err := c.compileForm(val)
	if err != nil {
		return NewCompileError("compiling def value").Wrap(err)
	}
	c.emit(vm.OPSTV)
	c.decSP(1)
	c.tailPosition = tc
	return nil
}

func varCompiler(c *Context, form vm.Value) error {
	sym := form.(*vm.List).Next().First().(vm.Symbol)
	varr := c.constant(c.ns.LookupOrAdd(sym))
	c.emitWithArg(vm.OPLDC, varr)
	c.incSP(1)
	return nil
}
