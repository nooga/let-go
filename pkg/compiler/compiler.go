/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package compiler

import (
	"fmt"
	"io"
	"strings"

	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

type Context struct {
	parent         *Context
	consts         *vm.Consts
	chunk          *vm.CodeChunk
	formalArgs     map[vm.Symbol]int
	source         string
	variadric      bool
	locals         []map[vm.Symbol]int
	sp             int
	spMax          int
	isFunction     bool
	isClosure      bool
	closedOversC   int
	closedOvers    map[vm.Symbol]*closureCell
	closedOversSeq []vm.Symbol
	recurPoints    []*recurPoint
	tailPosition   bool
	debug          bool
	defName        string
}

func NewCompiler(consts *vm.Consts, ns *vm.Namespace) *Context {
	rt.CurrentNS.SetRoot(ns)
	return &Context{
		consts:      consts,
		source:      "<default>",
		locals:      []map[vm.Symbol]int{},
		closedOvers: map[vm.Symbol]*closureCell{},
		debug:       false,
	}
}

func NewDebugCompiler(consts *vm.Consts, ns *vm.Namespace) *Context {
	c := NewCompiler(consts, ns)
	c.debug = true
	return c
}

func (c *Context) SetSource(source string) *Context {
	c.source = source
	return c
}

func (c *Context) CurrentNS() *vm.Namespace {
	return rt.CurrentNS.Deref().(*vm.Namespace)
}

func (c *Context) Compile(s string) (*vm.CodeChunk, error) {
	r := NewLispReader(strings.NewReader(s), c.source)
	o, err := r.Read()
	if err != nil {
		return nil, err
	}
	c.resetSP()
	c.chunk = vm.NewCodeChunk(c.consts)
	err = c.compileForm(o)
	c.chunk.SetMaxStack(c.spMax)
	if err != nil {
		return nil, err
	}
	c.emit(vm.OP_RETURN)
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
			chunk.Append(vm.OP_POP)
		}
		formchunk := vm.NewCodeChunk(c.consts)
		c.chunk = formchunk
		c.resetSP()
		err = c.compileForm(o)
		c.chunk.SetMaxStack(c.spMax)
		if err != nil {
			return nil, result, err
		}
		chunk.AppendChunk(formchunk)

		formchunk.Append(vm.OP_RETURN)
		var f *vm.Frame
		if c.debug {
			f = vm.NewDebugFrame(formchunk, nil)
		} else {
			f = vm.NewFrame(formchunk, nil)
		}
		result, err = f.Run()
		if err != nil {
			return nil, result, err
		}
		compiledForms++
	}

	c.chunk = chunk

	c.emit(vm.OP_RETURN)
	c.decSP(1)
	return c.chunk, result, nil
}

func (c *Context) emit(op int32) {
	c.chunk.Append(op | int32(c.sp<<16))
}

func (c *Context) emitWithArg(op int32, arg int) {
	c.chunk.Append(op | int32(c.sp<<16))
	c.chunk.Append32(arg)
}

func (c *Context) constant(v vm.Value) int {
	return c.consts.Intern(v)
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
		parent:         c,
		consts:         c.consts,
		chunk:          fchunk,
		formalArgs:     make(map[vm.Symbol]int),
		locals:         []map[vm.Symbol]int{},
		closedOvers:    make(map[vm.Symbol]*closureCell),
		closedOversSeq: []vm.Symbol{},
		isFunction:     true,
		tailPosition:   true,
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
	f.SetName(c.defName)
	n := c.constant(f)
	c.emitWithArg(vm.OP_LOAD_CONST, n)
	c.incSP(1)

	// if we have a closure on our hands then add closed overs
	if ctx.isClosure {
		c.emit(vm.OP_MAKE_CLOSURE)
		for _, s := range ctx.closedOversSeq {
			clo := ctx.closedOvers[s]
			_ = clo.source().emit()
			c.emit(vm.OP_PUSH_CLOSEDOVER)
			c.decSP(1)
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
		c.closedOversSeq = append(c.closedOversSeq, s)
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
	case vm.IntType, vm.StringType, vm.NilType, vm.BooleanType, vm.KeywordType, vm.CharType, vm.VoidType, vm.FuncType:
		n := c.constant(o)
		c.emitWithArg(vm.OP_LOAD_CONST, n)
		c.incSP(1)
	case vm.SymbolType:
		cel := c.symbolLookup(o.(vm.Symbol))
		if cel != nil {
			return cel.emit()
		}
		// when symbol not found so far we have a free variable on our hands
		v := c.CurrentNS().Lookup(o.(vm.Symbol))
		if v == vm.NIL {
			return NewCompileError("Can't resolve " + string(o.(vm.Symbol)) + " in this context")
		}
		varn := c.constant(v)
		c.emitWithArg(vm.OP_LOAD_VAR, varn)
		c.incSP(1)
	case vm.ArrayVectorType:
		tp := c.tailPosition
		c.tailPosition = false
		v := o.(vm.ArrayVector)
		// FIXME detect const vectors and push them like this
		//if len(v) == 0 {
		//	n := c.constant(v)
		//	c.emitWithArg(vm.OP_LOAD_CONST, n)
		//	c.incSP(1)
		//	return nil
		//}
		vector := c.constant(rt.CoreNS.Lookup("vector"))
		c.emitWithArg(vm.OP_LOAD_CONST, vector)
		c.incSP(1)
		for i := range v {
			err := c.compileForm(v[i])
			if err != nil {
				return NewCompileError("compiling vector elements").Wrap(err)
			}
		}
		c.emitWithArg(vm.OP_INVOKE, len(v))
		c.decSP(len(v))
		c.tailPosition = tp
	case vm.MapType:
		tp := c.tailPosition
		c.tailPosition = false
		v := o.(vm.Map)
		// FIXME detect const maps and push them like this
		//if len(v) == 0 {
		//	n := c.constant(v)
		//	c.emitWithArg(vm.OP_LOAD_CONST, n)
		//	c.incSP(1)
		//	return nil
		//}
		hashMap := c.constant(rt.CoreNS.Lookup("hash-map"))
		c.emitWithArg(vm.OP_LOAD_CONST, hashMap)
		c.incSP(1)
		for k, val := range v {
			err := c.compileForm(k)
			if err != nil {
				return NewCompileError("compiling map key").Wrap(err)
			}
			err = c.compileForm(val)
			if err != nil {
				return NewCompileError("compiling map value").Wrap(err)
			}
		}
		c.emitWithArg(vm.OP_INVOKE, len(v)*2)
		c.decSP(len(v) * 2)
		c.tailPosition = tp
	case vm.ListType:
		if o == vm.EmptyList {
			c.emitWithArg(vm.OP_LOAD_CONST, c.constant(vm.EmptyList))
			c.incSP(1)
			return nil
		}
		fn := o.(*vm.List).First()
		// check if we're looking at a special form
		if fn.Type() == vm.SymbolType {
			fnsym := fn.(vm.Symbol)
			formCompiler, ok := specialForms[fnsym]
			if ok {
				return formCompiler(c, o)
			}

			if fnsym[0] == '.' && len(fnsym) > 1 {
				newform := o.(*vm.List).Next()
				if newform.(vm.Collection).RawCount() < 1 {
					return NewCompileError("Malformed member expression, expecting (.member target ...)")
				}
				instance := newform.First()
				member := vm.EmptyList.Cons(fnsym[1:]).Cons(vm.Symbol("quote"))
				newform = newform.Next().Cons(member).Cons(instance).Cons(vm.Symbol("."))
				return c.compileForm(newform)
			}

			fvar := c.CurrentNS().Lookup(fnsym)
			if fvar != vm.NIL && fvar.(*vm.Var).IsMacro() {
				argvec := o.(*vm.List).Next().(*vm.List).Unbox().([]vm.Value)
				newform, err := fvar.(*vm.Var).Invoke(argvec)
				if err != nil {
					return NewCompileError("Executing macro failed").Wrap(err)
				}
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
				return NewCompileError("compiling arguments " + args.First().String()).Wrap(err)
			}
			args = args.Next()
		}

		c.emitWithArg(vm.OP_INVOKE, argc)
		c.decSP(argc)

		c.tailPosition = tp
	}
	return nil
}

func (c *Context) emitWithArgPlaceholder(inst int32) int {
	placeholder := c.currentAddress()
	c.emitWithArg(inst, 0)
	return placeholder
}

func (c *Context) currentAddress() int {
	return c.chunk.Length()
}

func (c *Context) updatePlaceholderArg(placeholder int, arg int) {
	c.chunk.Update32(placeholder+1, int32(arg))
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

func (c *Context) resetSP() {
	c.sp = 0
	c.spMax = 0
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
	locals  int
	argsc   int
}

func (c *Context) pushRecurPoint(argsc int) {
	c.recurPoints = append(c.recurPoints, &recurPoint{
		address: c.currentAddress(),
		locals:  len(c.locals),
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
		"set!":  setBangCompiler,
		"fn*":   fnCompiler,
		"quote": quoteCompiler,
		"var":   varCompiler,
		"let*":  letCompiler,
		"loop*": loopCompiler,
		"recur": recurCompiler,
		"trace": traceCompiler,
	}
}

func traceCompiler(c *Context, form vm.Value) error {
	args := form.(*vm.List).Next()
	c.emit(vm.OP_TRACE_ENABLE)
	for args != vm.EmptyList {
		err := c.compileForm(args.First())
		if err != nil {
			return NewCompileError("compiling trace arguments").Wrap(err)
		}
		args = args.Next()
	}
	c.emit(vm.OP_TRACE_DISABLE)
	return nil
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

	if rp != nil {
		passedScopes := len(c.locals) - rp.locals
		ignore := 0
		if passedScopes > 0 {
			passedLocals := 0
			for i := 0; i < passedScopes; i++ {
				passedLocals += len(c.locals[len(c.locals)-i-1])
			}
			ignore += passedLocals
		}
		c.emitWithArg(vm.OP_RECUR, c.currentAddress()-rp.address)
		c.chunk.Append32(argc)
		c.chunk.Append32(ignore)
	} else if c.isFunction {
		c.emitWithArg(vm.OP_RECUR_FN, argc)
	}
	c.tailPosition = tp
	c.decSP(argc - 1) // this is needed to keep the balance of if branches
	return nil
}

func loopCompiler(c *Context, form vm.Value) error {
	bindings := form.(*vm.List).Next()
	binds, ok := bindings.First().(vm.ArrayVector)
	if !ok {
		return NewCompileError("loop bindings should be a vector")
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
		c.emitWithArg(vm.OP_LOAD_CONST, c.constant(vm.NIL))
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
				c.emit(vm.OP_POP)
				c.decSP(1)
			}
		}
	}
	c.popLocals()
	c.popRecurPoint()
	if bindn > 0 {
		c.emitWithArg(vm.OP_POP_N, bindn)
		c.decSP(bindn)
	}
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
			return NewCompileError("let binding name must be a symbol: " + name.String())
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
		c.emitWithArg(vm.OP_LOAD_CONST, c.constant(vm.NIL))
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
				c.emit(vm.OP_POP)
				c.decSP(1)
			}
		}
	}
	c.popLocals()
	if bindn > 0 {
		c.emitWithArg(vm.OP_POP_N, bindn)
		c.decSP(bindn)
	}
	c.tailPosition = tc
	return nil
}

func quoteCompiler(c *Context, form vm.Value) error {
	n := c.constant(form.(vm.Seq).Next().First())
	c.emitWithArg(vm.OP_LOAD_CONST, n)
	c.incSP(1)
	return nil
}

func fnFormCompiler(c *Context, args vm.ArrayVector, bodyf *vm.List) error {
	fc, err := c.enterFn(args)
	if err != nil {
		return NewCompileError("compiling fn args").Wrap(err)
	}
	defer c.leaveFn(fc)

	body := bodyf.Unbox().([]vm.Value)
	l := len(body)
	if l == 0 {
		fc.emitWithArg(vm.OP_LOAD_CONST, fc.constant(vm.NIL))
		fc.incSP(1)
		fc.emit(vm.OP_RETURN)
		return nil
	}
	for i := range body {
		err := fc.compileForm(body[i])
		if err != nil {
			return NewCompileError("compiling fn body").Wrap(err)
		}
		if i < l-1 {
			fc.emit(vm.OP_POP)
			fc.decSP(1)
		}
	}
	fc.emit(vm.OP_RETURN)
	return nil
}

func fnCompiler(c *Context, form vm.Value) error {
	f := form.(*vm.List).Next()

	if args, ok := f.First().(vm.ArrayVector); ok {
		// we have (fn* [args] body)
		body := f.Next().(*vm.List)
		return fnFormCompiler(c, args, body)
	} else if _, ok := f.First().(vm.Seq); ok {
		// we have (fn* ([] ...))
		i := 0
		for b := f; b != vm.EmptyList; b = b.Next() {
			e := b.First().(*vm.List)
			args := e.First().(vm.ArrayVector)
			body := e.Next().(*vm.List)
			err := fnFormCompiler(c, args, body)
			if err != nil {
				return err
			}
			i++
		}
		c.emitWithArg(vm.OP_MAKE_MULTI_ARITY, i)
		c.decSP(i - 1)
	} else {
		return NewCompileError("unexpected fn form")
	}

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
	elseJumpStart := c.emitWithArgPlaceholder(vm.OP_BRANCH_FALSE)
	c.decSP(1)
	c.tailPosition = tc

	// compile then branch
	err = c.compileForm(args[1])
	c.decSP(1)

	if err != nil {
		return NewCompileError("compiling if then branch").Wrap(err)
	}
	finJumpStart := c.emitWithArgPlaceholder(vm.OP_JUMP)
	elseJumpEnd := c.currentAddress()
	c.updatePlaceholderArg(elseJumpStart, elseJumpEnd-elseJumpStart)
	if l == 3 {
		err = c.compileForm(args[2])

		if err != nil {
			return NewCompileError("compiling if else branch").Wrap(err)
		}
	} else {
		c.emitWithArg(vm.OP_LOAD_CONST, c.constant(vm.NIL))
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
		c.emitWithArg(vm.OP_LOAD_CONST, c.constant(vm.NIL))
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
			c.emit(vm.OP_POP)
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
		fmt.Println(args)
		return NewCompileError(fmt.Sprintf("def: wrong number of forms (%d), need 2", l))
	}
	var meta vm.Value = vm.NIL
	sym := args[0]
	val := args[1]
	if sym.Type() == vm.ListType {
		ss := sym.(vm.Seq)
		if ss.First() != vm.Symbol("with-meta") {
			return NewCompileError(fmt.Sprintf("def: first argument must be a symbol, got (%v)", sym))
		}
		ss = ss.Next()
		sym = ss.First()
		meta = ss.Next().First()
	}
	if sym.Type() != vm.SymbolType {
		return NewCompileError(fmt.Sprintf("def: first argument must be a symbol, got (%v)", sym))
	}
	c.defName = sym.String()
	varr := c.CurrentNS().LookupOrAdd(sym.(vm.Symbol))
	if meta != vm.NIL {
		m := meta.(vm.Map)
		if m[vm.Keyword("dynamic")] == vm.TRUE {
			varr.(*vm.Var).SetDynamic()
		}
		if m[vm.Keyword("private")] == vm.TRUE {
			varr.(*vm.Var).SetPrivate()
		}
	}
	c.emitWithArg(vm.OP_LOAD_CONST, c.constant(varr))
	c.incSP(1)
	err := c.compileForm(val)
	if err != nil {
		return NewCompileError("compiling def value").Wrap(err)
	}
	c.emit(vm.OP_SET_VAR)
	c.decSP(1)
	c.tailPosition = tc
	c.defName = ""
	return nil
}

// FIXME this is just def with a different name basically
func setBangCompiler(c *Context, form vm.Value) error {
	tc := c.tailPosition
	c.tailPosition = false
	args := form.(*vm.List).Next().Unbox().([]vm.Value)
	l := len(args)
	if l != 2 {
		return NewCompileError(fmt.Sprintf("set!: wrong number of forms (%d), need 2", l))
	}
	sym := args[0]
	val := args[1]
	if sym.Type() != vm.SymbolType {
		return NewCompileError(fmt.Sprintf("set!: first argument must be a symbol, got (%v)", sym))
	}
	varr := c.constant(c.CurrentNS().Lookup(sym.(vm.Symbol)))
	c.emitWithArg(vm.OP_LOAD_CONST, varr)
	c.incSP(1)
	err := c.compileForm(val)
	if err != nil {
		return NewCompileError("compiling set! value").Wrap(err)
	}
	c.emit(vm.OP_SET_VAR)
	c.decSP(1)
	c.tailPosition = tc
	return nil
}

func varCompiler(c *Context, form vm.Value) error {
	sym := form.(*vm.List).Next().First().(vm.Symbol)
	varr := c.constant(c.CurrentNS().LookupOrAdd(sym))
	c.emitWithArg(vm.OP_LOAD_CONST, varr)
	c.incSP(1)
	return nil
}
