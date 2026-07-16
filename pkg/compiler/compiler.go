/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package compiler

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"sync/atomic"
	"unicode"

	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

type Context struct {
	parent     *Context
	consts     *vm.Consts
	chunk      *vm.CodeChunk
	formalArgs map[vm.Symbol]int
	argCount   int // total fixed-arity parameter slots, including `_`s
	source     string
	variadric  bool
	locals     []map[vm.Symbol]int
	// localSlotCounts mirrors locals: each entry is the raw count of stack
	// slots in that scope. We track this separately because shadowed bindings
	// (e.g. `(let [[a w] ... [b w] ...])`) overwrite the symbol→slot map
	// entry while the older stack slot is still live; len(locals[i]) under-
	// counts in that case, but the stack footprint is bindn (matching
	// let/loop's OP_POP_N). recurCompiler reads from here to compute its
	// `ignore` operand correctly.
	localSlotCounts []int
	sp              int
	spMax           int
	isFunction      bool
	isClosure       bool
	closedOversC    int
	closedOvers     map[vm.Symbol]*closureCell
	closedOversSeq  []vm.Symbol
	recurPoints     []*recurPoint
	tailPosition    bool
	debug           bool
	defName         string
	currentForm     vm.Value // tracks the form being compiled for error source info
	currentList     vm.Value // tracks the enclosing list form for error source info
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

// NewTransientCompiler returns a compiler for ONE transient evaluation (a
// REPL/nREPL input, load-string, eval, a pod or host request): constants the
// evaluation introduces intern into a CHILD pool layered on parent, so they
// stay reachable only through the chunks and functions that use them and are
// collected with them — instead of rooting the parent pool forever, one
// entry per eval, unbounded in a long-lived session.
//
// Every transient-eval site must build its compiler here; a bare NewCompiler
// call is reserved for pools that SHOULD live as long as the process (boot,
// require/ns loading, build tools). Grep for NewCompiler( to audit.
func NewTransientCompiler(parent *vm.Consts, ns *vm.Namespace) *Context {
	return NewCompiler(vm.NewChildConsts(parent), ns)
}

// ChildForEval returns a fresh compilation context for one transient
// top-level evaluation layered on this context's pool: same namespace and
// debug mode, child constant pool (see NewTransientCompiler).
//
// Field semantics: debug and source are COPIED (user setting + error-report
// label); the namespace is inherited through the process-global CurrentNS
// (continuity across evals); everything else (chunk, locals, stack state) is
// freshly initialized — each eval is independent. Repeated calls produce
// SIBLING children of this context's pool, not a nested chain: each transient
// eval's constants live exactly as long as its own chunks, and siblings
// cannot alias each other's index space (Intern only reuses parent indices
// below the child's base).
func (c *Context) ChildForEval() *Context {
	child := NewTransientCompiler(c.consts, c.CurrentNS())
	child.debug = c.debug
	child.source = c.source
	return child
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
	var evalTopForm func(o vm.Value) error
	evalTopForm = func(o vm.Value) error {
		// Clojure's Compiler.eval special-cases a top-level (do ...): each
		// subform is compiled AND evaluated in turn, so an earlier subform's
		// compile-time effect (a require loading a namespace, an in-ns, an
		// intern) is visible when a later subform compiles. This is what
		// makes (do (require 'x) (x/f)) work at the REPL (#195). An empty
		// (do) falls through to compileForm and evaluates to nil.
		if lst, ok := o.(*vm.List); ok && lst.Next() != nil {
			if first, isSym := lst.First().(vm.Symbol); isSym && first == "do" {
				for s := lst.Next(); s != nil; s = s.Next() {
					if err := evalTopForm(s.First()); err != nil {
						return err
					}
				}
				return nil
			}
		}
		if compiledForms > 0 {
			chunk.Append(vm.OP_POP)
		}
		formchunk := vm.NewCodeChunk(c.consts)
		c.chunk = formchunk
		c.resetSP()
		err := c.compileForm(o)
		c.chunk.SetMaxStack(c.spMax)
		if err != nil {
			return err
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
			return err
		}
		compiledForms++
		return nil
	}
	for {
		o, err := r.Read()
		if err != nil {
			if isErrorEOF(err) {
				break
			}
			return nil, result, err
		}
		if err := evalTopForm(o); err != nil {
			return nil, result, err
		}
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
		// Strip metadata wrappers from arg symbols: `[^String s]` is read as
		// `[(with-meta s {:tag String})]`. We don't yet attach meta to locals,
		// just drop it so the symbol check below succeeds.
		if lst, ok := a.(*vm.List); ok && lst.First() == vm.Symbol("with-meta") {
			if rest := lst.Next(); rest != nil {
				a = rest.First()
			}
		}
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
		// `_` is the conventional "ignored" placeholder but is still an ordinary
		// symbol in Clojure: it IS bound and may be referenced. When `_` repeats,
		// the last occurrence wins, which the Symbol-keyed map gives us for free
		// by overwriting on assignment. argCount tracks the total slot count,
		// including `_`s, for arity checks.
		fc.argCount++
		fc.formalArgs[s] = i
	}
	return fc, nil
}

func (c *Context) leaveFn(ctx *Context) {
	fnchunk := ctx.chunk
	// Record parameter names as debug info (slot -> name), sorted by slot for
	// deterministic output (formalArgs is a map). `_` and `&` are already
	// excluded from formalArgs.
	if len(ctx.formalArgs) > 0 {
		params := make([]vm.LocalVar, 0, len(ctx.formalArgs))
		for name, slot := range ctx.formalArgs {
			params = append(params, vm.LocalVar{Slot: slot, Name: string(name)})
		}
		sort.Slice(params, func(i, j int) bool { return params[i].Slot < params[j].Slot })
		for _, p := range params {
			fnchunk.AddLocalVar(p.Slot, p.Name)
		}
	}
	fnchunk.SetMaxStack(ctx.spMax)
	f := vm.MakeFunc(ctx.argCount, ctx.variadric, fnchunk)
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
	// Locals and args in the current scope shadow a closed-over variable of the
	// same name: `(let [v (f v)] v)` where v is also captured from an enclosing
	// scope must see the NEW binding in the body, not the captured value.
	// Checking closedOvers first (as before) made the let binding a no-op.
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
	if c.isClosure {
		clo := c.closedOvers[s]
		if clo != nil {
			return clo
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
	case vm.IntType, vm.FloatType, vm.StringType, vm.NilType, vm.BooleanType, vm.KeywordType, vm.CharType, vm.VoidType, vm.FuncType, vm.NativeFnType, vm.BigIntType, vm.RatioType, vm.BigDecimalType, vm.UUIDType, vm.InstantType, vm.RegexType:
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
				if hc, ok := rt.LookupHostClass(string(symVal)); ok {
					n := c.constant(hc)
					c.emitWithArg(vm.OP_LOAD_CONST, n)
					c.incSP(1)
					return nil
				}
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
			// Host-class fallback: a bare class symbol used as a value (e.g.
			// java.util.Map, CharSequence in `(instance? java.util.Map x)`)
			// resolves to a registered let-go value (typically a type).
			if hc, ok := rt.LookupHostClass(string(symVal)); ok {
				n := c.constant(hc)
				c.emitWithArg(vm.OP_LOAD_CONST, n)
				c.incSP(1)
				return nil
			}
			return c.compileError(fmt.Sprintf("Can't resolve %s in this context", symVal))
		}
		varn := c.constant(v)
		c.emitWithArg(vm.OP_LOAD_VAR, varn)
		c.incSP(1)
	case vm.ArrayVectorType:
		tp := c.tailPosition
		c.tailPosition = false
		v, ok := o.(vm.ArrayVector)
		if !ok {
			if me, ok := o.(vm.MapEntry); ok {
				v = vm.ArrayVector{me.Key, me.Value}
			} else {
				return c.compileError("expected vector form")
			}
		}
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

		arrayMap := c.constant(rt.CoreNS.Lookup("array-map"))
		c.emitWithArg(vm.OP_LOAD_CONST, arrayMap)
		c.incSP(1)

		// Get entries via Seq for both Map and PersistentMap
		var count int
		if sq, ok := o.(vm.Sequable); ok {
			s := sq.Seq()
			var entries []vm.Value
			for s != nil && s != vm.EmptyList {
				k, v, ok := vm.MapEntryKV(s.First())
				if !ok {
					s = s.Next()
					continue
				}
				entries = append(entries, k, v)
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
	case vm.SetType:
		tp := c.tailPosition
		c.tailPosition = false

		hashSet := c.constant(rt.CoreNS.Lookup("hash-set"))
		c.emitWithArg(vm.OP_LOAD_CONST, hashSet)
		c.incSP(1)

		count := 0
		if sq, ok := o.(vm.Sequable); ok {
			for s := sq.Seq(); s != nil && s != vm.EmptyList; s = s.Next() {
				if err := c.compileForm(s.First()); err != nil {
					return NewCompileError("compiling set element").Wrap(err)
				}
				count++
			}
		}

		c.emitWithArg(vm.OP_INVOKE, count)
		c.decSP(count)
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

			// (Name. args...) — record/type constructor shorthand.
			// Rewrites to (->Name args...) which defrecord defines.
			if len(fnsym) > 1 &&
				fnsym[len(fnsym)-1] == '.' &&
				fnsym[0] != '.' &&
				!strings.HasPrefix(string(fnsym), "->") {
				stem := string(fnsym[:len(fnsym)-1])
				ctor := vm.Symbol("->" + stem)
				args := lst.Next()
				var rewritten vm.Value
				if args == nil {
					rewritten = vm.EmptyList.Cons(ctor)
				} else {
					rewritten = args.Cons(ctor)
				}
				if info := vm.FormSource.Get(o); info != nil {
					vm.FormSource.Set(rewritten, *info)
				}
				return c.compileForm(rewritten)
			}

			if len(fnsym) > 1 && fnsym[0] == '.' && fnsym[1] != '.' {
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
				if info := vm.FormSource.Get(o); info != nil {
					vm.FormSource.Set(newform, *info)
				}
				return c.compileForm(newform)
			}

			// Locals shadow macros: skip macro expansion if name is bound in the
			// enclosing lexical scope (local, arg, or captured by a closure).
			fvar := vm.Value(vm.NIL)
			if !c.resolvesAsLexical(fnsym) {
				fvar = c.CurrentNS().Lookup(fnsym)
			}
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
				if info := vm.FormSource.Get(o); info != nil {
					vm.FormSource.Set(newform, *info)
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
		case "bit-and":
			return vm.OP_BIT_AND
		case "bit-or":
			return vm.OP_BIT_OR
		case "bit-xor":
			return vm.OP_BIT_XOR
		case "bit-and-not":
			return vm.OP_BIT_AND_NOT
		case "bit-shift-left":
			return vm.OP_BIT_SHIFT_LEFT
		case "bit-shift-right":
			return vm.OP_BIT_SHIFT_RIGHT
		case "unsigned-bit-shift-right":
			return vm.OP_UNSIGNED_BIT_SHIFT_RIGHT
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
		case "bit-not":
			return vm.OP_BIT_NOT
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
	c.localSlotCounts = append(c.localSlotCounts, 0)
}

func (c *Context) popLocals() {
	c.locals = c.locals[0 : len(c.locals)-1]
	c.localSlotCounts = c.localSlotCounts[0 : len(c.localSlotCounts)-1]
}

func (c *Context) addLocal(name vm.Symbol) {
	c.locals[len(c.locals)-1][name] = c.sp - 1
	// Record the source name for this slot as debug info (slot -> name), so it
	// survives into the bundle and can name locals in crash traces.
	if name != "_" {
		c.chunk.AddLocalVar(c.sp-1, string(name))
	}
	// Count every binding, even shadowed ones: the older slot is still on the
	// stack and counts toward this scope's footprint. recurCompiler relies on
	// this to compute `ignore` correctly when crossing a shadowing scope.
	c.localSlotCounts[len(c.localSlotCounts)-1]++
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

// resolvesAsLexical reports whether symbol is bound as a local, formal arg, or
// closed-over variable anywhere in the enclosing lexical scope, WITHOUT the
// capture side effects of symbolLookup. Used to decide head-position macro
// shadowing: a lexical binding named like a macro/special form (e.g. a local
// `fn`) must shadow it even when the use site is inside a nested fn/closure
// that captures it. lookupLocal alone only sees the current frame, so a
// captured binding would otherwise be mistaken for the macro.
func (c *Context) resolvesAsLexical(symbol vm.Symbol) bool {
	for ctx := c; ctx != nil; ctx = ctx.parent {
		if ctx.isClosure && ctx.closedOvers[symbol] != nil {
			return true
		}
		if ctx.lookupLocal(symbol) >= 0 {
			return true
		}
		if ctx.arg(symbol) >= 0 {
			return true
		}
	}
	return false
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
		"try":   tryCompiler,
	}
}

// caughtSymCounter feeds freshCaughtSym. Atomic because the stdlib
// determinism test compiles concurrently.
var caughtSymCounter atomic.Int64

// freshCaughtSym returns a binding symbol for a desugared catch handler that
// cannot collide with user bindings or with an enclosing try's own handler.
// The trailing tab makes the name unforgeable: the reader cannot produce a
// symbol containing whitespace, so no user source can shadow or reference it.
func freshCaughtSym() vm.Symbol {
	return vm.Symbol(fmt.Sprintf("__lg-caught-%d\t", caughtSymCounter.Add(1)))
}

func makeList(items ...vm.Value) (vm.Value, error) {
	return vm.ListType.Box(items)
}

func mustQuote(v vm.Value) vm.Value {
	q, err := makeList(vm.Symbol("quote"), v)
	if err != nil {
		panic(err) // two-element list construction cannot fail
	}
	return q
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

	// Separate body, catch clauses, and finally forms
	type catchClause struct {
		class   vm.Symbol // "" for let-go's bare (catch bind-sym body...)
		binding vm.Symbol
		body    []vm.Value
	}
	var bodyForms []vm.Value
	var catchSym vm.Symbol
	var catchForms []vm.Value
	var finallyForms []vm.Value
	var catchClauses []catchClause
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
					// Clojure-compatible form: (catch ClassSym bind-sym body...)
					// vs let-go's bare (catch bind-sym body...). A binding is
					// always a simple unqualified symbol, so a class-shaped
					// first token — qualified/dotted, or an uppercase simple
					// symbol — followed by a symbol marks a typed clause.
					// There is deliberately NO token-count fallback: under
					// class dispatch it would misread the bare
					// (catch e foo bar) as a typed clause on binding e.
					// JVM class names are never lowercase-simple.
					classSym := vm.Symbol("")
					if rest.Next() != nil {
						firstSym, firstIsSym := rest.First().(vm.Symbol)
						_, secondIsSym := rest.Next().First().(vm.Symbol)
						qualified := firstIsSym && strings.ContainsAny(string(firstSym), "./")
						upperSimple := false
						if firstIsSym && !qualified {
							for _, r := range string(firstSym) {
								upperSimple = unicode.IsUpper(r)
								break
							}
						}
						if firstIsSym && secondIsSym && (qualified || upperSimple) {
							classSym = firstSym
							rest = rest.Next() // drop the leading class symbol
						}
					}
					bindSym, ok := rest.First().(vm.Symbol)
					if !ok {
						return NewCompileError("catch requires a binding symbol")
					}
					var clauseBody []vm.Value
					for cb := rest.Next(); cb != nil; cb = cb.Next() {
						clauseBody = append(clauseBody, cb.First())
					}
					catchClauses = append(catchClauses, catchClause{
						class:   classSym,
						binding: bindSym,
						body:    clauseBody,
					})
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

	if len(catchClauses) == 1 && catchClauses[0].class == "" {
		// let-go's native bare catch compiles directly as the handler.
		catchSym = catchClauses[0].binding
		catchForms = catchClauses[0].body
	} else if len(catchClauses) > 0 {
		// Desugar typed and/or multiple clauses into a single handler that
		// dispatches on exception class in source order and rethrows when
		// nothing matches (the IR builder's parse-try generates the same
		// form — keep them in lockstep):
		//
		//   (catch <caught>
		//     (if (core/instance? Class1 <caught>) (let* [b1 <caught>] body1...)
		//       ... (core/throw <caught>)))
		//
		// core/-qualified so a user shadowing instance? or throw cannot
		// capture the dispatch. A bare clause tests as always-true, making
		// any later clauses dead, like Clojure's ordered catch clauses.
		caught := freshCaughtSym()
		acc, err := makeList(vm.Symbol("core/throw"), caught)
		if err != nil {
			return NewCompileError("building catch dispatch").Wrap(err)
		}
		for i := len(catchClauses) - 1; i >= 0; i-- {
			cl := catchClauses[i]
			armParts := append([]vm.Value{vm.Symbol("let*"), vm.ArrayVector{cl.binding, caught}}, cl.body...)
			arm, err := makeList(armParts...)
			if err != nil {
				return NewCompileError("building catch dispatch").Wrap(err)
			}
			if cl.class == "" {
				acc = arm
				continue
			}
			// catch-matches? resolves the class at dispatch time (a JVM-only
			// class let-go does not model never matches, rather than failing
			// compilation) and gives Throwable its catch-everything role.
			test, err := makeList(vm.Symbol("core/catch-matches?"), mustQuote(cl.class), caught)
			if err != nil {
				return NewCompileError("building catch dispatch").Wrap(err)
			}
			acc, err = makeList(vm.Symbol("if"), test, arm, acc)
			if err != nil {
				return NewCompileError("building catch dispatch").Wrap(err)
			}
		}
		catchSym = caught
		catchForms = []vm.Value{acc}
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

	// Layout (see exHandler/handleError in pkg/vm for the runtime side):
	//
	//   TRY_PUSH catchOff finOff       ; either offset 0 means "absent"
	//     <body>
	//   TRY_POP                        ; normal completion retires the handler
	//   JUMP finally (or after)
	//   catch:                         ; only when hasCatch
	//     <bind, catch body>
	//   TRY_POP                        ; only when a finally follows: the
	//                                  ; handler outlives catch entry so a
	//                                  ; throw from the catch body routes
	//                                  ; through the finally
	//   finally:                       ; single copy, shared by both paths
	//     <finally forms>
	//   FINALLY_END finallyOff         ; abnormal entry: rethrow the pending
	//   after:                         ; error; normal entry: no-op
	//
	// The normal path enters the finally with the try/catch result on the
	// stack; an abnormal entry (handleError) cuts the stack to savedSP and
	// pushes a placeholder so both paths run the finally at the same depth.
	// FINALLY_END carries the finally address so it can tell its own
	// abnormally-entered handler from an outer one that is also mid-finally.
	hasFinally := len(finallyForms) > 0

	// Normal completion: pop handler, jump to the finally (or the end)
	c.emit(vm.OP_TRY_POP)
	jumpOverCatchAddr := c.currentAddress()
	c.emitWithArg(vm.OP_JUMP, 0) // placeholder

	// Catch block starts here. The VM restored SP to savedSP and pushed the
	// thrown value, so the runtime SP matches the static sp (savedSP+1).
	catchAddr := c.currentAddress()

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

		if hasFinally {
			// The handler survived catch entry (handleError keeps it when a
			// finally exists); retire it now that the catch completed, then
			// fall through into the finally.
			c.emit(vm.OP_TRY_POP)
		}
	}

	// catchOffset stays 0 (sentinel: no catch) for a finally-only try, so a
	// throw in the body routes to the finally instead of swallowing the error.
	if hasCatch {
		c.chunk.Update32(tryPushAddr+1, int32(catchAddr-tryPushAddr))
	}

	if hasFinally {
		finallyAddr := c.currentAddress()
		c.chunk.Update32(tryPushAddr+2, int32(finallyAddr-tryPushAddr))
		c.chunk.Update32(jumpOverCatchAddr+1, int32(finallyAddr-jumpOverCatchAddr))

		for i, ff := range finallyForms {
			if err := c.compileForm(ff); err != nil {
				return err
			}
			if i < len(finallyForms)-1 {
				c.emit(vm.OP_POP)
				c.decSP(1)
			}
		}
		// Discard the finally result, keep the try/catch result (or the
		// abnormal-entry placeholder).
		c.emit(vm.OP_POP)
		c.decSP(1)

		finallyEndAddr := c.currentAddress()
		c.emitWithArg(vm.OP_FINALLY_END, finallyAddr-finallyEndAddr)
	} else {
		afterAddr := c.currentAddress()
		c.chunk.Update32(jumpOverCatchAddr+1, int32(afterAddr-jumpOverCatchAddr))
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
		if argc != c.argCount {
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
			for i := range passedScopes {
				// Use the per-scope slot count rather than len(map): a scope
				// with shadowed names has more live stack slots than distinct
				// symbols, and OP_RECUR must drop all of them.
				passedLocals += c.localSlotCounts[len(c.localSlotCounts)-i-1]
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

// parseBindingsVector extracts binding forms from a vector.
func parseBindingsVector(val vm.Value) ([]vm.Value, error) {
	switch bv := val.(type) {
	case vm.ArrayVector:
		return []vm.Value(bv), nil
	case vm.PersistentVector:
		return bv.Unbox().([]vm.Value), nil
	default:
		return nil, fmt.Errorf("bindings should be a vector, got %T: %v", val, val)
	}
}

// compileBindings validates and compiles let/loop bindings sequentially.
func compileBindings(c *Context, binds []vm.Value, opName string) (int, error) {
	bindn := 0
	for i := 0; i < len(binds); i += 2 {
		name := binds[i]
		// Strip a metadata wrapper from the binding name: `^long x` is read as
		// `(with-meta x {:tag long})`. As with fn params, we don't yet attach
		// the tag to the local (a future hook for the IR typeinfer pass), just
		// unwrap to the bare symbol so the check below succeeds.
		if lst, ok := name.(*vm.List); ok && lst.First() == vm.Symbol("with-meta") {
			if rest := lst.Next(); rest != nil {
				name = rest.First()
			}
		}
		if name.Type() != vm.SymbolType {
			return 0, NewCompileError(fmt.Sprintf("%s binding name must be a symbol: %v", opName, name))
		}
		if i+1 >= len(binds) {
			return 0, NewCompileError(fmt.Sprintf("%s bindings must have even number of forms", opName))
		}
		value := binds[i+1]
		err := c.compileForm(value)
		if err != nil {
			return 0, NewCompileError(fmt.Sprintf("compiling %s binding", opName)).Wrap(err)
		}
		c.addLocal(name.(vm.Symbol))
		bindn++
	}
	return bindn, nil
}

func loopCompiler(c *Context, form vm.Value) error {
	bindings := form.(*vm.List).Next()
	if bindings == nil {
		return NewCompileError("loop requires bindings")
	}
	binds, err := parseBindingsVector(bindings.First())
	if err != nil {
		return NewCompileError("loop bindings should be a vector").Wrap(err)
	}
	body := bindings.Next()
	c.pushLocals()
	tp := c.tailPosition
	c.tailPosition = false
	bindn, err := compileBindings(c, binds, "loop")
	if err != nil {
		return err
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
	binds, err := parseBindingsVector(bindings.First())
	if err != nil {
		return NewCompileError("let bindings should be a vector").Wrap(err)
	}
	body := bindings.Next()
	c.pushLocals()
	tc := c.tailPosition
	c.tailPosition = false
	bindn, err := compileBindings(c, binds, "let")
	if err != nil {
		return err
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

// listArgs extracts argument values from a list form's body.
func listArgs(form vm.Value) []vm.Value {
	nxt := form.(*vm.List).Next()
	if nxt == nil {
		return nil
	}
	if nl, ok := nxt.(*vm.List); ok {
		return nl.Unbox().([]vm.Value)
	}
	var args []vm.Value
	for s := nxt; s != nil; s = s.Next() {
		args = append(args, s.First())
	}
	return args
}

func ifCompiler(c *Context, form vm.Value) error {
	tc := c.tailPosition
	//c.tailPosition = tc

	args := listArgs(form)
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

// tryDetectInNS detects top-level (in-ns 'foo) and updates the compiler namespace early.
// Use LookupOrRegisterNSNoLoad rather than rt.NS so we don't trigger the resolver to load 'foo from disk.
func tryDetectInNS(c *Context, arg vm.Value) {
	if arg.Type() != vm.ListType {
		return
	}
	lst := arg.(vm.Seq)
	first := lst.First()
	if first == nil || first.Type() != vm.SymbolType || vm.Symbol(first.(vm.Symbol)) != vm.Symbol("in-ns") {
		return
	}
	alist := lst.Next()
	if alist == nil {
		return
	}
	q := alist.First()
	if q == nil || q.Type() != vm.ListType {
		return
	}
	qq := q.(vm.Seq)
	if qq.First() != vm.Symbol("quote") {
		return
	}
	qqN := qq.Next()
	if qqN == nil {
		return
	}
	namev := qqN.First()
	if namev == nil || namev.Type() != vm.SymbolType {
		return
	}
	ns := rt.LookupOrRegisterNSNoLoad(string(namev.(vm.Symbol)))
	if ns == nil {
		return
	}
	c.SetCurrentNS(ns)
}

// trySimulateAlias handles core/alias compile-time simulation.
func trySimulateAlias(c *Context, lst vm.Seq) error {
	asArgs := lst.Next()
	if asArgs == nil {
		return nil
	}
	qa := asArgs.First()
	asArgs = asArgs.Next()
	if asArgs == nil {
		return nil
	}
	qb := asArgs.First()
	if qa == nil || qb == nil || qa.Type() != vm.ListType || qb.Type() != vm.ListType {
		return nil
	}
	qqa := qa.(vm.Seq)
	qqb := qb.(vm.Seq)
	if qqa.First() != vm.Symbol("quote") || qqb.First() != vm.Symbol("quote") {
		return nil
	}
	qqaN := qqa.Next()
	qqbN := qqb.Next()
	if qqaN == nil || qqbN == nil {
		return nil
	}
	alias, okAlias := qqaN.First().(vm.Symbol)
	nsname, okNS := qqbN.First().(vm.Symbol)
	if !okAlias || !okNS {
		return nil
	}
	target, err := rt.RequireNS(string(nsname))
	if err != nil {
		return c.compileError(err.Error())
	}
	c.CurrentNS().Alias(alias, target)
	return nil
}

// trySimulateRefer handles core/refer compile-time simulation.
func trySimulateRefer(c *Context, lst vm.Seq) error {
	rArgs := lst.Next()
	if rArgs == nil {
		return nil
	}
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
	if nsQ == nil || nsQ.Type() != vm.ListType {
		return nil
	}
	qq := nsQ.(vm.Seq)
	if qq.First() != vm.Symbol("quote") {
		return nil
	}
	qqN := qq.Next()
	if qqN == nil {
		return nil
	}
	nsname, ok := qqN.First().(vm.Symbol)
	if !ok {
		return nil
	}
	target, err := rt.RequireNS(string(nsname))
	if err != nil {
		return c.compileError(err.Error())
	}
	c.CurrentNS().Refer(target, aliasStr, all)
	return nil
}

// trySimulateImportVar handles core/import-var compile-time simulation.
func trySimulateImportVar(c *Context, lst vm.Seq) error {
	ivArgs := lst.Next()
	if ivArgs == nil {
		return nil
	}
	qn := ivArgs.First()
	ivArgs = ivArgs.Next()
	if ivArgs == nil {
		return nil
	}
	qfrom := ivArgs.First()
	ivArgs = ivArgs.Next()
	if ivArgs == nil {
		return nil
	}
	qto := ivArgs.First()

	if qn == nil || qfrom == nil || qto == nil {
		return nil
	}
	if qn.Type() != vm.ListType || qfrom.Type() != vm.ListType || qto.Type() != vm.ListType {
		return nil
	}
	qnn := qn.(vm.Seq)
	qff := qfrom.(vm.Seq)
	qtt := qto.(vm.Seq)
	if qnn.First() != vm.Symbol("quote") || qff.First() != vm.Symbol("quote") || qtt.First() != vm.Symbol("quote") {
		return nil
	}
	qnnN := qnn.Next()
	qffN := qff.Next()
	qttN := qtt.Next()
	if qnnN == nil || qffN == nil || qttN == nil {
		return nil
	}
	nsname, okNS := qnnN.First().(vm.Symbol)
	from, okFrom := qffN.First().(vm.Symbol)
	to, okTo := qttN.First().(vm.Symbol)
	if !okNS || !okFrom || !okTo {
		return nil
	}
	fromNs, err := rt.RequireNS(string(nsname))
	if err != nil {
		return c.compileError(err.Error())
	}
	c.CurrentNS().ImportVar(fromNs, from, to)
	return nil
}

// trySimulateUse handles (use 'ns) compile-time simulation.
func trySimulateUse(c *Context, lst vm.Seq) error {
	uArgs := lst.Next()
	for uArgs != nil {
		qa := uArgs.First()
		if qa == nil || qa.Type() != vm.ListType {
			uArgs = uArgs.Next()
			continue
		}
		qq := qa.(vm.Seq)
		if qq.First() != vm.Symbol("quote") {
			uArgs = uArgs.Next()
			continue
		}
		qqN := qq.Next()
		if qqN == nil {
			uArgs = uArgs.Next()
			continue
		}
		nsname, ok := qqN.First().(vm.Symbol)
		if !ok {
			uArgs = uArgs.Next()
			continue
		}
		target, err := rt.RequireNS(string(nsname))
		if err != nil {
			return c.compileError(err.Error())
		}
		c.CurrentNS().Refer(target, "", true)
		uArgs = uArgs.Next()
	}
	return nil
}

// trySimulateNsHelper delegates simulation to helper functions for namespace operations.
func trySimulateNsHelper(c *Context, arg vm.Value) error {
	if arg.Type() != vm.ListType {
		return nil
	}
	lst := arg.(vm.Seq)
	first := lst.First()
	if first == nil || first.Type() != vm.SymbolType {
		return nil
	}
	fname := vm.Symbol(first.(vm.Symbol))

	switch fname {
	case vm.Symbol("core/alias"):
		return trySimulateAlias(c, lst)
	case vm.Symbol("core/refer"):
		return trySimulateRefer(c, lst)
	case vm.Symbol("core/import-var"):
		return trySimulateImportVar(c, lst)
	case vm.Symbol("use"):
		return trySimulateUse(c, lst)
	}
	return nil
}

func doCompiler(c *Context, form vm.Value) error {
	args := listArgs(form)
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
		if i == 0 {
			tryDetectInNS(c, args[i])
		}
		// Simulate ns helpers at compile-time so later forms in the same do can resolve
		err := trySimulateNsHelper(c, args[i])
		if err != nil {
			return err
		}
		if tc && i == l-1 {
			c.tailPosition = true
		}
		err = c.compileForm(args[i])
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

func assocMeta(meta vm.Value, key vm.Value, val vm.Value) vm.Value {
	if meta == nil || meta == vm.NIL {
		return vm.NewPersistentMap([]vm.Value{key, val})
	}
	if m, ok := meta.(*vm.PersistentMap); ok {
		return m.Assoc(key, val).(vm.Value)
	}
	if m, ok := meta.(vm.Map); ok {
		return m.Assoc(key, val).(vm.Value)
	}
	return vm.NewPersistentMap([]vm.Value{key, val})
}

func metaValueAt(meta vm.Value, key vm.Value) vm.Value {
	if meta == nil || meta == vm.NIL {
		return vm.NIL
	}
	if m, ok := meta.(*vm.PersistentMap); ok {
		return m.ValueAt(key)
	}
	if m, ok := meta.(vm.Map); ok {
		return m.ValueAt(key)
	}
	return vm.NIL
}

func defCompiler(c *Context, form vm.Value) error {
	tc := c.tailPosition
	c.tailPosition = false
	args := listArgs(form)
	l := len(args)
	if l < 1 || l > 3 {
		return NewCompileError(fmt.Sprintf("def: wrong number of forms (%d), need 1, 2 or 3", l))
	}
	var doc vm.Value = vm.NIL
	if l == 3 {
		if docString, ok := args[1].(vm.String); ok {
			doc = docString
			args = []vm.Value{args[0], args[2]}
			l = 2
		} else {
			return NewCompileError("def: 3-arg form requires a docstring (String) as second argument")
		}
	}
	var meta vm.Value = vm.NIL
	sym := args[0]
	var val vm.Value = vm.NIL
	if l == 2 {
		val = args[1]
	}
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
	if doc != vm.NIL {
		meta = assocMeta(meta, vm.Keyword("doc"), doc)
	}
	c.defName = sym.String()
	varr := c.CurrentNS().LookupOrAdd(sym.(vm.Symbol))
	if meta != vm.NIL {
		v := varr.(*vm.Var)
		v.SetMeta(meta)
		if vm.IsTruthy(metaValueAt(meta, vm.Keyword("dynamic"))) {
			v.SetDynamic()
		}
		if vm.IsTruthy(metaValueAt(meta, vm.Keyword("private"))) {
			v.SetPrivate()
		}
	}
	c.emitWithArg(vm.OP_LOAD_CONST, c.constant(varr))
	c.incSP(1)
	if l == 1 {
		// No-init form (def x): intern the var but leave its root binding
		// UNAFFECTED, matching Clojure ("If init is not supplied, the root
		// binding of the var is unaffected"). This is a forward declaration /
		// promise — never a write — so it must not clobber an existing root
		// or a value a later (def x v) will provide. The var itself is the
		// expression result, already on the stack from OP_LOAD_CONST above.
		c.tailPosition = tc
		c.defName = ""
		return nil
	}
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
	args := listArgs(form)
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
