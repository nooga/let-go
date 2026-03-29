/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
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
	currentForm    vm.Value // tracks the form being compiled for error source info
	currentList    vm.Value // tracks the enclosing list form for error source info
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

func (c *Context) Consts() *vm.Consts {
	return c.consts
}

func (c *Context) CurrentNS() *vm.Namespace {
	return rt.CurrentNS.Deref().(*vm.Namespace)
}

func (c *Context) SetCurrentNS(ns *vm.Namespace) {
	rt.CurrentNS.SetRoot(ns)
}

func (c *Context) Compile(s string) (*vm.CodeChunk, error) {
	vm.SourceRegistry.Register(c.source, s)
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
	// Buffer source for error display
	srcBytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, vm.NIL, err
	}
	src := string(srcBytes)
	vm.SourceRegistry.Register(c.source, src)
	r := NewLispReader(strings.NewReader(src), c.source)
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
		result, err = f.RunProtected()
		vm.ReleaseFrame(f)
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

// compileError creates a CompileError with source info from the given form.
func compileErrorAt(msg string, form vm.Value) *CompileError {
	info := vm.FormSource.Get(form)
	return NewCompileErrorWithSource(msg, info)
}

// compileError creates a CompileError with source info from the current form context.
func (c *Context) compileError(msg string) *CompileError {
	// Try the current form, then walk up the form stack via parent list
	if info := vm.FormSource.Get(c.currentForm); info != nil {
		return NewCompileErrorWithSource(msg, info)
	}
	if info := vm.FormSource.Get(c.currentList); info != nil {
		return NewCompileErrorWithSource(msg, info)
	}
	return NewCompileError(msg)
}

func (c *Context) compileForm(o vm.Value) error {
	// Track current form for error reporting
	prevForm := c.currentForm
	c.currentForm = o
	defer func() { c.currentForm = prevForm }()

	// Emit source location for this form
	if info := vm.FormSource.Get(o); info != nil {
		c.chunk.AddSourceInfo(*info)
	}
	switch o.Type() {
	case vm.IntType, vm.FloatType, vm.StringType, vm.NilType, vm.BooleanType, vm.KeywordType, vm.CharType, vm.VoidType, vm.FuncType, vm.BigIntType:
		n := c.constant(o)
		c.emitWithArg(vm.OP_LOAD_CONST, n)
		c.incSP(1)
	case vm.SymbolType:
		symVal := o.(vm.Symbol)
		// If qualified like ns/sym
		if sns, inner := symVal.Namespaced(); sns != vm.NIL {
			// Resolve core/* via global core ns so (ns ...) expansion works before refers
			if string(sns.(vm.Symbol)) == rt.NameCoreNS {
				target := rt.NS(rt.NameCoreNS)
				v := target.Lookup(inner.(vm.Symbol))
				if v == vm.NIL {
					return c.compileError(fmt.Sprintf("Can't resolve %s in this context", symVal))
				}
				varn := c.constant(v)
				c.emitWithArg(vm.OP_LOAD_VAR, varn)
				c.incSP(1)
				return nil
			}
			// Non-core qualified: honor aliases and refers in current ns
			v := c.CurrentNS().Lookup(symVal)
			if v == vm.NIL {
				return c.compileError(fmt.Sprintf("Can't resolve %s in this context", symVal))
			}
			varn := c.constant(v)
			c.emitWithArg(vm.OP_LOAD_VAR, varn)
			c.incSP(1)
			return nil
		}

		cel := c.symbolLookup(symVal)
		if cel != nil {
			return cel.emit()
		}
		// when symbol not found so far we have a free variable on our hands
		v := c.CurrentNS().Lookup(symVal)
		if v == vm.NIL {
			return c.compileError(fmt.Sprintf("Can't resolve %s in this context", symVal))
		}
		varn := c.constant(v)
		c.emitWithArg(vm.OP_LOAD_VAR, varn)
		c.incSP(1)
	case vm.ArrayVectorType:
		tp := c.tailPosition
		c.tailPosition = false
		v := o.(vm.ArrayVector)
		// Optimization: const vectors could be pushed as constants
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

		hashMap := c.constant(rt.CoreNS.Lookup("hash-map"))
		c.emitWithArg(vm.OP_LOAD_CONST, hashMap)
		c.incSP(1)

		// Get entries via Seq for both Map and PersistentMap
		var count int
		if sq, ok := o.(vm.Sequable); ok {
			s := sq.Seq()
			var entries []vm.Value
			for s != nil && s != vm.EmptyList {
				entry := s.First().(vm.ArrayVector)
				entries = append(entries, entry[0], entry[1])
				s = s.Next()
			}
			count = len(entries) / 2
			for _, e := range entries {
				err := c.compileForm(e)
				if err != nil {
					return NewCompileError("compiling map entry").Wrap(err)
				}
			}
		}

		c.emitWithArg(vm.OP_INVOKE, count*2)
		c.decSP(count * 2)
		c.tailPosition = tp
	case vm.ListType:
		prevList := c.currentList
		c.currentList = o
		defer func() { c.currentList = prevList }()
		if o == vm.EmptyList {
			c.emitWithArg(vm.OP_LOAD_CONST, c.constant(vm.EmptyList))
			c.incSP(1)
			return nil
		}
		lst, isList := o.(*vm.List)
		if !isList {
			if seq, ok := o.(vm.Seq); ok {
				var vals []vm.Value
				for s := seq; s != nil; s = s.Next() {
					vals = append(vals, s.First())
				}
				realized, _ := vm.ListType.Box(vals)
				return c.compileForm(realized)
			}
			n := c.constant(o)
			c.emitWithArg(vm.OP_LOAD_CONST, n)
			c.incSP(1)
			return nil
		}
		fn := lst.First()
		// check if we're looking at a special form
		if fn.Type() == vm.SymbolType {
			fnsym := fn.(vm.Symbol)
			formCompiler, ok := specialForms[fnsym]
			if ok {
				return formCompiler(c, o)
			}

			if fnsym[0] == '.' && len(fnsym) > 1 {
				newform := lst.Next()
				if newform == nil {
					return NewCompileError("Malformed member expression, expecting (.member target ...)")
				}
				if coll, ok := newform.(vm.Collection); ok && coll.RawCount() < 1 {
					return NewCompileError("Malformed member expression, expecting (.member target ...)")
				}
				instance := newform.First()
				member := vm.EmptyList.Cons(fnsym[1:]).Cons(vm.Symbol("quote"))
				nxt := newform.Next()
				if nxt == nil {
					newform = vm.EmptyList.Cons(member).Cons(instance).Cons(vm.Symbol("."))
				} else {
					newform = nxt.Cons(member).Cons(instance).Cons(vm.Symbol("."))
				}
				return c.compileForm(newform)
			}

			fvar := c.CurrentNS().Lookup(fnsym)
			if fvar != vm.NIL && fvar.(*vm.Var).IsMacro() {
				nxt := lst.Next()
				var argvec []vm.Value
				if nxt != nil {
					if nl, ok := nxt.(*vm.List); ok {
						argvec = nl.Unbox().([]vm.Value)
					} else {
						for s := nxt; s != nil; s = s.Next() {
							argvec = append(argvec, s.First())
						}
					}
				}
				newform, err := fvar.(*vm.Var).Deref().(vm.Fn).Invoke(argvec)
				if err != nil {
					return NewCompileError(fmt.Sprintf("Executing macro %s (%s) failed", fvar, fvar.(*vm.Var).Deref())).Wrap(err)
				}
				return c.compileForm(newform)
			}
		}

		tp := c.tailPosition
		c.tailPosition = false

		args := lst.Next()
		argc := 0
		if args != nil {
			if coll, ok := args.(vm.Collection); ok {
				argc = coll.Count().Unbox().(int)
			} else {
				for s := args; s != nil; s = s.Next() {
					argc++
				}
			}
		}

		// Try to emit a specialized opcode for known core builtins
		if fn.Type() == vm.SymbolType {
			if fastOp := c.tryFastOpcode(fn.(vm.Symbol), argc); fastOp != 0 {
				// Compile arguments only (no function position on stack)
				for a := lst.Next(); a != nil; a = a.Next() {
					err := c.compileForm(a.First())
					if err != nil {
						return NewCompileError("compiling arguments " + a.First().String()).Wrap(err)
					}
				}
				c.emit(fastOp)
				if argc == 2 {
					c.decSP(1) // binary: 2 args -> 1 result
				}
				// unary (inc/dec): 1 arg -> 1 result, no SP change
				c.tailPosition = tp
				return nil
			}
		}

		// treat as function invocation if this is not a special form
		err := c.compileForm(fn)
		if err != nil {
			return NewCompileError("compiling function position").Wrap(err)
		}

		for a := lst.Next(); a != nil; a = a.Next() {
			err := c.compileForm(a.First())
			if err != nil {
				return NewCompileError("compiling arguments " + a.First().String()).Wrap(err)
			}
		}

		if tp && c.currentRecurPoint() == nil {
			c.emitWithArg(vm.OP_TAIL_CALL, argc)
		} else {
			c.emitWithArg(vm.OP_INVOKE, argc)
		}
		c.decSP(argc)

		c.tailPosition = tp
	}
	return nil
}

// tryFastOpcode returns a specialized opcode for known core builtins,
// or 0 if no fast path is available. Only emits for binary (arity 2)
// and unary (arity 1) cases with known symbols.
func (c *Context) tryFastOpcode(sym vm.Symbol, argc int) int32 {
	// Only optimize unqualified symbols that resolve to core vars
	if sym.Namespace() != vm.NIL {
		return 0
	}
	// Check that the symbol resolves to a core var (not a local binding)
	if c.symbolLookup(sym) != nil {
		return 0 // local binding shadows the core var
	}
	v := c.CurrentNS().Lookup(sym)
	if v == vm.NIL {
		return 0
	}

	switch argc {
	case 2:
		switch sym {
		case "+":
			return vm.OP_ADD
		case "-":
			return vm.OP_SUB
		case "*":
			return vm.OP_MUL
		case "<":
			return vm.OP_LT
		case "<=":
			return vm.OP_LTE
		case ">":
			return vm.OP_GT
		case ">=":
			return vm.OP_GTE
		case "=":
			return vm.OP_EQ
		}
	case 1:
		switch sym {
		case "inc":
			return vm.OP_INC
		case "dec":
			return vm.OP_DEC
		}
	}
	return 0
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
		"try":   tryCompiler,
	}
}

func traceCompiler(c *Context, form vm.Value) error {
	args := form.(*vm.List).Next()
	c.emit(vm.OP_TRACE_ENABLE)
	for args != nil {
		err := c.compileForm(args.First())
		if err != nil {
			return NewCompileError("compiling trace arguments").Wrap(err)
		}
		args = args.Next()
	}
	c.emit(vm.OP_TRACE_DISABLE)
	return nil
}

func tryCompiler(c *Context, form vm.Value) error {
	// Parse: (try body... (catch sym catch-body...) (finally finally-body...))
	nxt := form.(*vm.List).Next()
	if nxt == nil {
		return NewCompileError("try requires a body")
	}

	// Collect all args into a slice, handling both List and Cons/LazySeq
	var allForms []vm.Value
	for s := nxt; s != nil; s = s.Next() {
		allForms = append(allForms, s.First())
	}

	// Separate body, catch, and finally forms
	var bodyForms []vm.Value
	var catchSym vm.Symbol
	var catchForms []vm.Value
	var finallyForms []vm.Value
	hasCatch := false

	for _, f := range allForms {
		if seq, ok := f.(vm.Seq); ok {
			first := seq.First()
			if first != nil && first.Type() == vm.SymbolType {
				sym := first.(vm.Symbol)
				if sym == "catch" {
					hasCatch = true
					rest := seq.Next()
					if rest == nil {
						return NewCompileError("catch requires a binding symbol")
					}
					bindSym, ok := rest.First().(vm.Symbol)
					if !ok {
						return NewCompileError("catch requires a binding symbol")
					}
					catchSym = bindSym
					for cb := rest.Next(); cb != nil; cb = cb.Next() {
						catchForms = append(catchForms, cb.First())
					}
					continue
				}
				if sym == "finally" {
					for fb := seq.Next(); fb != nil; fb = fb.Next() {
						finallyForms = append(finallyForms, fb.First())
					}
					continue
				}
			}
		}
		bodyForms = append(bodyForms, f)
	}

	if !hasCatch && len(finallyForms) == 0 {
		// No catch or finally — just compile body as do
		for i, bf := range bodyForms {
			err := c.compileForm(bf)
			if err != nil {
				return err
			}
			if i < len(bodyForms)-1 {
				c.emit(vm.OP_POP)
				c.decSP(1)
			}
		}
		return nil
	}

	// Emit: OP_TRY_PUSH catchOffset finallyOffset
	tryPushAddr := c.currentAddress()
	c.emit(vm.OP_TRY_PUSH)
	c.chunk.Append32(0) // placeholder for catchOffset
	c.chunk.Append32(0) // placeholder for finallyOffset

	// Compile body
	tc := c.tailPosition
	c.tailPosition = false
	if len(bodyForms) == 0 {
		c.emitWithArg(vm.OP_LOAD_CONST, c.constant(vm.NIL))
		c.incSP(1)
	} else {
		for i, bf := range bodyForms {
			err := c.compileForm(bf)
			if err != nil {
				return err
			}
			if i < len(bodyForms)-1 {
				c.emit(vm.OP_POP)
				c.decSP(1)
			}
		}
	}

	// Normal completion: pop handler, jump over catch
	c.emit(vm.OP_TRY_POP)
	jumpOverCatchAddr := c.currentAddress()
	c.emitWithArg(vm.OP_JUMP, 0) // placeholder

	// Catch block starts here
	catchAddr := c.currentAddress()

	// At this point, the VM has pushed the thrown value on the stack
	// and restored SP to savedSP before pushing, so sp == savedSP+1.
	// We need to account for that in our SP tracking.
	// The body left SP at savedSP+1 (body result), but the VM restored to savedSP
	// and pushed the thrown value, so the net SP is the same.

	if hasCatch {
		// Bind the thrown value as a local
		c.pushLocals()
		c.addLocal(catchSym)

		if len(catchForms) == 0 {
			// No catch body — push nil as catch result
			c.emitWithArg(vm.OP_LOAD_CONST, c.constant(vm.NIL))
			c.incSP(1)
		} else {
			for i, cf := range catchForms {
				err := c.compileForm(cf)
				if err != nil {
					return err
				}
				if i < len(catchForms)-1 {
					c.emit(vm.OP_POP)
					c.decSP(1)
				}
			}
		}

		// Pop catch binding, keep catch result
		c.emitWithArg(vm.OP_POP_N, 1)
		c.decSP(1)
		c.popLocals()
	}
	// else: no catch, the thrown value is already on stack as the result

	// Patch jump-over-catch
	afterCatch := c.currentAddress()
	c.chunk.Update32(jumpOverCatchAddr+1, int32(afterCatch-jumpOverCatchAddr))

	// Patch TRY_PUSH catchOffset (relative to TRY_PUSH instruction)
	c.chunk.Update32(tryPushAddr+1, int32(catchAddr-tryPushAddr))

	// Finally block (if present)
	if len(finallyForms) > 0 {
		finallyAddr := c.currentAddress()
		// Patch TRY_PUSH finallyOffset
		c.chunk.Update32(tryPushAddr+2, int32(finallyAddr-tryPushAddr))

		for i, ff := range finallyForms {
			err := c.compileForm(ff)
			if err != nil {
				return err
			}
			if i < len(finallyForms)-1 {
				c.emit(vm.OP_POP)
				c.decSP(1)
			}
		}
		// Discard finally result, keep try/catch result
		c.emit(vm.OP_POP)
		c.decSP(1)
	}

	c.tailPosition = tc
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
	argc := 0
	if args != nil {
		if coll, ok := args.(vm.Collection); ok {
			argc = coll.Count().Unbox().(int)
		} else {
			for s := args; s != nil; s = s.Next() {
				argc++
			}
		}
	}

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

	for args != nil {
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
	if bindings == nil {
		return NewCompileError("loop requires bindings")
	}
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
			return NewCompileError("compiling loop binding").Wrap(err)
		}
		c.addLocal(name.(vm.Symbol))
		bindn++
	}
	c.pushRecurPoint(bindn)
	if body == nil || body == vm.EmptyList {
		c.emitWithArg(vm.OP_LOAD_CONST, c.constant(vm.NIL))
		c.incSP(1)
	} else {
		for b := body; b != nil; b = b.Next() {
			if b.Next() == nil {
				c.tailPosition = true
			}
			err := c.compileForm(b.First())
			if err != nil {
				return NewCompileError("compiling loop body").Wrap(err)
			}
			if b.Next() != nil {
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
	if bindings == nil {
		return NewCompileError("let requires bindings")
	}
	binds, ok := bindings.First().(vm.ArrayVector)
	if !ok {
		return NewCompileError(fmt.Sprintf("let bindings should be a vector, got %T: %v", bindings.First(), bindings.First()))
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
	if body == nil || body == vm.EmptyList {
		c.emitWithArg(vm.OP_LOAD_CONST, c.constant(vm.NIL))
		c.incSP(1)
	} else {
		for b := body; b != nil; b = b.Next() {
			if tc && b.Next() == nil {
				c.tailPosition = true
			}
			err := c.compileForm(b.First())
			if err != nil {
				return NewCompileError("compiling let body").Wrap(err)
			}
			if b.Next() != nil {
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
	nxt := form.(vm.Seq).Next()
	if nxt == nil {
		n := c.constant(vm.NIL)
		c.emitWithArg(vm.OP_LOAD_CONST, n)
		c.incSP(1)
		return nil
	}
	n := c.constant(nxt.First())
	c.emitWithArg(vm.OP_LOAD_CONST, n)
	c.incSP(1)
	return nil
}

func fnFormCompiler(c *Context, args vm.ArrayVector, bodyf vm.Seq) error {
	fc, err := c.enterFn(args)
	if err != nil {
		return NewCompileError("compiling fn args").Wrap(err)
	}
	defer c.leaveFn(fc)

	// Realize body to slice
	var body []vm.Value
	for s := bodyf; s != nil; s = s.Next() {
		body = append(body, s.First())
	}
	l := len(body)
	if l == 0 {
		fc.emitWithArg(vm.OP_LOAD_CONST, fc.constant(vm.NIL))
		fc.incSP(1)
		fc.emit(vm.OP_RETURN)
		return nil
	}
	// Only the last form is in tail position
	fc.tailPosition = false
	for i := range body {
		if i == l-1 {
			fc.tailPosition = true
		}
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
	if f == nil {
		return NewCompileError("unexpected fn form")
	}

	if args, ok := f.First().(vm.ArrayVector); ok {
		// we have (fn* [args] body)
		body := f.Next()
		if body == nil {
			body = vm.EmptyList
		}
		return fnFormCompiler(c, args, body)
	} else if _, ok := f.First().(vm.Seq); ok {
		// we have (fn* ([] ...))
		i := 0
		for b := f; b != nil; b = b.Next() {
			e := b.First().(vm.Seq)
			args := e.First().(vm.ArrayVector)
			ebody := e.Next()
			if ebody == nil {
				ebody = vm.EmptyList
			}
			err := fnFormCompiler(c, args, ebody)
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
	//c.tailPosition = tc

	nxt := form.(*vm.List).Next()
	var args []vm.Value
	if nxt != nil {
		if nl, ok := nxt.(*vm.List); ok {
			args = nl.Unbox().([]vm.Value)
		} else {
			for s := nxt; s != nil; s = s.Next() {
				args = append(args, s.First())
			}
		}
	}
	l := len(args)
	if l < 2 || l > 3 {
		return NewCompileError(fmt.Sprintf("if: wrong number of forms (%d), need 2 or 3", l))
	}
	c.tailPosition = false
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
	nxt := form.(*vm.List).Next()
	var args []vm.Value
	if nxt != nil {
		if nl, ok := nxt.(*vm.List); ok {
			args = nl.Unbox().([]vm.Value)
		} else {
			for s := nxt; s != nil; s = s.Next() {
				args = append(args, s.First())
			}
		}
	}
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
		// Detect top-level (in-ns 'foo) and update compiler namespace early
		if i == 0 && args[i].Type() == vm.ListType {
			lst := args[i].(*vm.List)
			if lst.First().Type() == vm.SymbolType && vm.Symbol(lst.First().(vm.Symbol)) == vm.Symbol("in-ns") {
				alist := lst.Next()
				if alist != nil {
					q := alist.First()
					if q.Type() == vm.ListType {
						qq := q.(vm.Seq)
						if qq.First() == vm.Symbol("quote") {
							qqN := qq.Next()
							if qqN != nil {
								namev := qqN.First()
								if namev.Type() == vm.SymbolType {
									if ns := rt.NS(string(namev.(vm.Symbol))); ns != nil {
										c.SetCurrentNS(ns)
									}
								}
							}
						}
					}
				}
			}
		}
		// Simulate ns helpers at compile-time so later forms in the same do can resolve
		if args[i].Type() == vm.ListType {
			lst := args[i].(*vm.List)
			if lst.First().Type() == vm.SymbolType {
				fname := vm.Symbol(lst.First().(vm.Symbol))
				// core/alias
				if fname == vm.Symbol("core/alias") {
					asArgs := lst.Next()
					if asArgs != nil {
						qa := asArgs.First()
						asArgs = asArgs.Next()
						if asArgs != nil {
							qb := asArgs.First()
							if qa.Type() == vm.ListType && qb.Type() == vm.ListType {
								qqa := qa.(vm.Seq)
								qqb := qb.(vm.Seq)
								if qqa.First() == vm.Symbol("quote") && qqb.First() == vm.Symbol("quote") {
									qqaN := qqa.Next()
									qqbN := qqb.Next()
									if qqaN != nil && qqbN != nil {
										alias := qqaN.First().(vm.Symbol)
										nsname := qqbN.First().(vm.Symbol)
										if target := rt.NS(string(nsname)); target != nil {
											c.CurrentNS().Alias(alias, target)
										}
									}
								}
							}
						}
					}
				}
				// core/refer (ns, alias, all)
				if fname == vm.Symbol("core/refer") {
					rArgs := lst.Next()
					if rArgs != nil {
						nsQ := rArgs.First()
						rArgs = rArgs.Next()
						aliasStr := ""
						all := true
						if rArgs != nil {
							if s, ok := rArgs.First().(vm.String); ok {
								aliasStr = string(s)
							}
							rArgs = rArgs.Next()
						}
						if rArgs != nil {
							if b, ok := rArgs.First().(vm.Boolean); ok {
								all = bool(b)
							}
						}
						if nsQ.Type() == vm.ListType {
							qq := nsQ.(vm.Seq)
							if qq.First() == vm.Symbol("quote") {
								qqN := qq.Next()
								if qqN != nil {
									nsname := qqN.First().(vm.Symbol)
									if target := rt.NS(string(nsname)); target != nil {
										c.CurrentNS().Refer(target, aliasStr, all)
									}
								}
							}
						}
					}
				}
				// core/import-var (from-ns, from, to)
				if fname == vm.Symbol("core/import-var") {
					ivArgs := lst.Next()
					if ivArgs != nil {
						qn := ivArgs.First()
						ivArgs = ivArgs.Next()
						if ivArgs != nil {
							qfrom := ivArgs.First()
							ivArgs = ivArgs.Next()
							if ivArgs != nil {
								qto := ivArgs.First()
								if qn.Type() == vm.ListType && qfrom.Type() == vm.ListType && qto.Type() == vm.ListType {
									qnn := qn.(vm.Seq)
									qff := qfrom.(vm.Seq)
									qtt := qto.(vm.Seq)
									qnnN := qnn.Next()
									qffN := qff.Next()
									qttN := qtt.Next()
									if qnnN != nil && qffN != nil && qttN != nil && qnn.First() == vm.Symbol("quote") && qff.First() == vm.Symbol("quote") && qtt.First() == vm.Symbol("quote") {
										fromNs := rt.NS(string(qnnN.First().(vm.Symbol)))
										from := qffN.First().(vm.Symbol)
										to := qttN.First().(vm.Symbol)
										if fromNs != nil {
											c.CurrentNS().ImportVar(fromNs, from, to)
										}
									}
								}
							}
						}
					}
				}
				// (use 'ns) — load namespace and refer all at compile time
				if fname == vm.Symbol("use") {
					uArgs := lst.Next()
					for uArgs != nil {
						qa := uArgs.First()
						if qa.Type() == vm.ListType {
							qq := qa.(vm.Seq)
							if qq.First() == vm.Symbol("quote") {
								qqN := qq.Next()
								if qqN != nil {
									nsname := qqN.First().(vm.Symbol)
									if target := rt.NS(string(nsname)); target != nil {
										c.CurrentNS().Refer(target, "", true)
									}
								}
							}
						}
						uArgs = uArgs.Next()
					}
				}
			}
		}
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
	nxt := form.(*vm.List).Next()
	var args []vm.Value
	if nxt != nil {
		if nl, ok := nxt.(*vm.List); ok {
			args = nl.Unbox().([]vm.Value)
		} else {
			for s := nxt; s != nil; s = s.Next() {
				args = append(args, s.First())
			}
		}
	}
	l := len(args)
	if l != 2 {
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
		if m, ok := meta.(*vm.PersistentMap); ok {
			if vm.IsTruthy(m.ValueAt(vm.Keyword("dynamic"))) {
				varr.(*vm.Var).SetDynamic()
			}
			if vm.IsTruthy(m.ValueAt(vm.Keyword("private"))) {
				varr.(*vm.Var).SetPrivate()
			}
		} else if m, ok := meta.(vm.Map); ok {
			if m[vm.Keyword("dynamic")] == vm.TRUE {
				varr.(*vm.Var).SetDynamic()
			}
			if m[vm.Keyword("private")] == vm.TRUE {
				varr.(*vm.Var).SetPrivate()
			}
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

func setBangCompiler(c *Context, form vm.Value) error {
	tc := c.tailPosition
	c.tailPosition = false
	nxt := form.(*vm.List).Next()
	var args []vm.Value
	if nxt != nil {
		if nl, ok := nxt.(*vm.List); ok {
			args = nl.Unbox().([]vm.Value)
		} else {
			for s := nxt; s != nil; s = s.Next() {
				args = append(args, s.First())
			}
		}
	}
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
	// Try compile-time resolution only
	v := c.CurrentNS().Lookup(sym)
	if v == vm.NIL {
		return c.compileError(fmt.Sprintf("Can't resolve %s in this context", sym))
	}
	varr := c.constant(v)
	c.emitWithArg(vm.OP_LOAD_CONST, varr)
	c.incSP(1)
	return nil
}
