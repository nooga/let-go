/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"fmt"
	"reflect"
)

type theFuncType struct{}

func (t *theFuncType) String() string  { return t.Name() }
func (t *theFuncType) Type() ValueType { return TypeType }
func (t *theFuncType) Unbox() any      { return reflect.TypeFor[*theFuncType]() }

func (t *theFuncType) Name() string { return "let-go.lang.Fn" }
func (t *theFuncType) Box(fn any) (Value, error) {
	return NIL, NewTypeError(fn, "can't be boxed as", t)
}

var FuncType *theFuncType = &theFuncType{}

type Func struct {
	name        string
	arity       int
	isVariadric bool
	chunk       *CodeChunk
	meta        Value // IMeta support; nil until (with-meta fn m) copies one in
}

func MakeFunc(arity int, variadric bool, c *CodeChunk) *Func {
	return &Func{
		arity:       arity,
		isVariadric: variadric,
		chunk:       c,
	}
}

func (l *Func) SetName(n string) {
	l.name = n
}

func (l *Func) Type() ValueType { return FuncType }

// Meta implements IMeta.
func (l *Func) Meta() Value {
	if l.meta == nil {
		return NIL
	}
	return l.meta
}

// WithMeta implements IMeta. Returns a copy carrying m. A capture-free fn is a
// shared compiler constant pushed by OP_LOAD_CONST, so this MUST copy rather
// than mutate — otherwise every evaluation of the fn literal would observe the
// metadata. The chunk is immutable, so sharing it across the copy is safe.
func (l *Func) WithMeta(m Value) Value {
	cp := *l
	cp.meta = m
	return &cp
}

type FuncInterface func(any)

// Unbox implements Unbox
func (l *Func) Unbox() any {
	proxy := func(in []reflect.Value) []reflect.Value {
		args := make([]Value, len(in))
		for i := range in {
			a, _ := BoxValue(in[i]) // error not propagatable through reflect proxy
			args[i] = a
		}
		f := NewFrame(l.chunk, args)
		out, _ := f.Run() // error not propagatable through reflect proxy
		return []reflect.Value{reflect.ValueOf(out.Unbox())}
	}
	return func(fptr any) {
		fn := reflect.ValueOf(fptr).Elem()
		v := reflect.MakeFunc(fn.Type(), proxy)
		fn.Set(v)
	}
}

func (l *Func) Arity() int {
	return l.arity
}

// boxRest packages variadic rest args for the rest parameter. Empty
// rest binds nil, like Clojure: (defn f [& xs] xs) (f) => nil.
func boxRest(rest []Value) (Value, error) {
	if len(rest) == 0 {
		return NIL, nil
	}
	return ListType.Box(rest)
}

// bytecodeCallTarget is the resolved, frame-independent description of a
// directly executable bytecode call. Frame allocation and same-frame tail
// replacement are deliberately separate transitions so this resolver can be
// shared with the constant-space tail-call work in #620.
type bytecodeCallTarget struct {
	fn          *Func
	args        []Value
	closedOvers []Value
}

// resolveBytecodeCall unwraps metadata, selects multi-arity variants,
// preserves closure captures, validates arity, and packs variadic arguments.
// It returns false for callables that must continue through ExecContext.Invoke.
//
// Fixed-arity args may still borrow the caller's operand stack. A transition
// that reuses that same frame must copy them into frame-owned storage before
// resetting the operand stack. Variadic packing always returns a fresh slice.
func resolveBytecodeCall(fn Fn, args []Value) (bytecodeCallTarget, bool, error) {
	display := fn
	var closedOvers []Value
	for {
		switch t := fn.(type) {
		case *MetaFn:
			fn = t.Wrapped()
		case *MultiArityFn:
			variant, err := t.variantFor(display, len(args))
			if err != nil {
				return bytecodeCallTarget{}, false, err
			}
			fn = variant
		case *Closure:
			closedOvers = t.closedOvers
			fn = t.fn
		case *Func:
			prepared := args
			if t.isVariadric {
				if len(prepared) < t.arity-1 {
					return bytecodeCallTarget{}, false, NewExecutionError(fmt.Sprintf("function %s expected at least %d args, got %d", display, t.arity-1, len(prepared)))
				}
				restlist, boxErr := boxRest(prepared[t.arity-1:])
				if boxErr != nil {
					return bytecodeCallTarget{}, false, boxErr
				}
				// Never append into prepared: it may be a window into the
				// caller's operand stack or a host-owned argument slice.
				packed := make([]Value, t.arity)
				copy(packed, prepared[:t.arity-1])
				packed[t.arity-1] = restlist
				prepared = packed
			} else if len(prepared) != t.arity {
				return bytecodeCallTarget{}, false, NewExecutionError(fmt.Sprintf("function %s expected %d args, got %d", display, t.arity, len(prepared)))
			}
			return bytecodeCallTarget{
				fn:          t,
				args:        prepared,
				closedOvers: closedOvers,
			}, true, nil
		default:
			return bytecodeCallTarget{}, false, nil
		}
	}
}

func runBytecodeCallTarget(ec *ExecContext, target bytecodeCallTarget) (Value, error) {
	f := newFrameForBytecodeCall(target, ec)
	result, err := f.Run()
	ReleaseFrame(f)
	return result, err
}

func (l *Func) Invoke(pargs []Value) (result Value, err error) {
	return l.invokeIn(RootExecContext, pargs)
}

// invokeIn runs the function with the given ExecContext active in its frame,
// so dynamic bindings propagate into the call. Invoke is invokeIn against the
// root context.
func (l *Func) invokeIn(ec *ExecContext, pargs []Value) (result Value, err error) {
	target, ok, err := resolveBytecodeCall(l, pargs)
	if err != nil {
		return NIL, err
	}
	if !ok {
		return NIL, NewExecutionError("unsupported function type")
	}
	return runBytecodeCallTarget(ec, target)
}

func (l *Func) String() string {
	if len(l.name) > 0 {
		return fmt.Sprintf("<fn %s %p>", l.name, l)
	}
	return fmt.Sprintf("<fn %p>", l)
}

// Chunk returns the code chunk.
func (l *Func) Chunk() *CodeChunk { return l.chunk }

// FuncName returns the function name.
func (l *Func) FuncName() string { return l.name }

// IsVariadic returns whether the function is variadic.
func (l *Func) IsVariadic() bool { return l.isVariadric }

func (l *Func) MakeClosure() Fn {
	return &Closure{
		closedOvers: nil,
		fn:          l,
	}
}

type Closure struct {
	closedOvers []Value
	fn          Fn
	meta        Value // IMeta support: (with-meta (fn ...) m); nil until set
}

func (l *Closure) Type() ValueType { return FuncType }

// Meta implements IMeta.
func (l *Closure) Meta() Value {
	if l.meta == nil {
		return NIL
	}
	return l.meta
}

// WithMeta implements IMeta. Returns a copy carrying m; the compiled fn and the
// closed-over slots are immutable at call time, so sharing them is safe and only
// the metadata differs. Lets (with-meta (fn ...) m) round-trip through meta,
// which spying/instrumentation libraries (e.g. bond) rely on.
func (l *Closure) WithMeta(m Value) Value {
	cp := *l
	cp.meta = m
	return &cp
}

// Unbox implements Unbox
func (l *Closure) Unbox() any {
	proxy := func(in []reflect.Value) []reflect.Value {
		args := make([]Value, len(in))
		for i := range in {
			a, _ := BoxValue(in[i]) // error not propagatable through reflect proxy
			args[i] = a
		}
		out, _ := l.Invoke(args)
		return []reflect.Value{reflect.ValueOf(out.Unbox())}
	}
	return func(fptr any) {
		fn := reflect.ValueOf(fptr).Elem()
		v := reflect.MakeFunc(fn.Type(), proxy)
		fn.Set(v)
	}
}

func (l *Closure) Arity() int {
	return l.fn.Arity()
}

func (l *Closure) Invoke(pargs []Value) (result Value, err error) {
	return l.invokeIn(RootExecContext, pargs)
}

// invokeIn runs the closure with the given ExecContext active in its frame,
// so dynamic bindings propagate into the call. Invoke delegates to invokeIn
// against the root context.
func (l *Closure) invokeIn(ec *ExecContext, pargs []Value) (result Value, err error) {
	target, ok, err := resolveBytecodeCall(l, pargs)
	if err != nil {
		return NIL, err
	}
	if ok {
		return runBytecodeCallTarget(ec, target)
	}

	if mfn, ok := l.fn.(*MultiArityFn); ok {
		variant, variantErr := mfn.variantFor(l, len(pargs))
		if variantErr != nil {
			return NIL, variantErr
		}
		return ec.Invoke(variant, pargs)
	}

	return NIL, NewExecutionError("unsupported closure function type")
}

func (l *Closure) String() string {
	return l.fn.String()
}

type MultiArityFn struct {
	fns   map[int]Fn
	rest  Fn
	arity int
	name  string
	meta  Value // IMeta support; nil until (with-meta fn m) copies one in
}

func (l *MultiArityFn) Type() ValueType { return FuncType }

// Meta implements IMeta.
func (l *MultiArityFn) Meta() Value {
	if l.meta == nil {
		return NIL
	}
	return l.meta
}

// WithMeta implements IMeta. Returns a copy carrying m; the arity variants are
// immutable at call time, so sharing them across the copy is safe.
func (l *MultiArityFn) WithMeta(m Value) Value {
	cp := *l
	cp.meta = m
	return &cp
}

// Unbox implements Unbox
func (l *MultiArityFn) Unbox() any {
	proxy := func(in []reflect.Value) []reflect.Value {
		args := make([]Value, len(in))
		for i := range in {
			a, _ := BoxValue(in[i]) // error not propagatable through reflect proxy
			args[i] = a
		}
		out, _ := l.Invoke(args)
		return []reflect.Value{reflect.ValueOf(out.Unbox())}
	}
	return func(fptr any) {
		fn := reflect.ValueOf(fptr).Elem()
		v := reflect.MakeFunc(fn.Type(), proxy)
		fn.Set(v)
	}
}

func (l *MultiArityFn) Arity() int {
	return l.arity
}

func (l *MultiArityFn) Invoke(pargs []Value) (Value, error) {
	return l.invokeIn(RootExecContext, pargs)
}

func (l *MultiArityFn) variantFor(display Fn, arity int) (Fn, error) {
	if f, ok := l.fns[arity]; ok {
		return f, nil
	}
	if l.rest != nil && arity >= restMinArgs(l.rest) {
		return l.rest, nil
	}
	return nil, NewExecutionError(fmt.Sprintf("function %s doesn't have a %d-arity variant", display, arity))
}

// restMinArgs is the minimum call width the variadic rest arm accepts. A
// variadic *Func's arity counts its rest SLOT — ([a b & c]) has arity 3 with
// two fixed args — so exactly-min is arity-1, matching resolveBytecodeCall's
// `len < arity-1` check. A NativeFn rest arm (NewArityNativeFn) declares its
// variadic minimum directly, so its Arity() already IS the floor (and
// NewCtxNativeFn's -1 sentinel accepts any width).
func restMinArgs(rest Fn) int {
	f := rest
	if c, ok := f.(*Closure); ok {
		f = c.fn
	}
	if ff, ok := f.(*Func); ok {
		return ff.arity - 1
	}
	return f.Arity()
}

// invokeIn runs the multi-arity function with the given ExecContext active,
// so dynamic bindings propagate into the selected variant's call. Invoke delegates
// to invokeIn against the root context.
func (l *MultiArityFn) invokeIn(ec *ExecContext, pargs []Value) (Value, error) {
	variant, err := l.variantFor(l, len(pargs))
	if err != nil {
		return NIL, err
	}
	return ec.Invoke(variant, pargs)
}

func (l *MultiArityFn) String() string {
	return fmt.Sprintf("<mfn %s %p>", l.name, l)
}

func MakeMultiArity(fns []Value) (*MultiArityFn, error) {
	ma := &MultiArityFn{
		arity: 0,
		fns:   map[int]Fn{},
		name:  "",
	}
	for i := range fns {
		e := fns[i]
		f, ok := e.(Fn)
		if !ok {
			return nil, NewExecutionError("making multi-arity function failed")
		}
		if ff, ok := f.(*Func); ok {
			ma.name = ff.name
		}
		a := f.Arity()
		if a > ma.arity {
			ma.arity = a
		}
		if rest, ok := f.(*Func); ok && rest.isVariadric {
			ma.rest = rest
		} else if rest, ok := f.(*NativeFn); ok && rest.isVariadric {
			// A wrapped callable declaring a variadic tail (ir.direct's
			// invokers via NewArityNativeFn). Without this arm it would be
			// filed as a FIXED arity in ma.fns, so (f a b c) past the
			// declared minimum would find no match instead of the rest-arm.
			ma.rest = rest
		} else if rest, ok := f.(*Closure); ok {
			if ff, ok := rest.fn.(*Func); ok && ff.isVariadric {
				ma.rest = rest
			} else {
				ma.fns[a] = f
			}
		} else {
			ma.fns[a] = f
		}
	}
	return ma, nil
}

func (l *MultiArityFn) MakeClosure() Fn {
	return &Closure{
		closedOvers: nil,
		fn:          l,
	}
}
