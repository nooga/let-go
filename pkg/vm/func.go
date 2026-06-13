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

func (l *Func) Invoke(pargs []Value) (result Value, err error) {
	return l.invokeIn(RootExecContext, pargs)
}

// invokeIn runs the function with the given ExecContext active in its frame,
// so dynamic bindings propagate into the call. Invoke is invokeIn against the
// root context.
func (l *Func) invokeIn(ec *ExecContext, pargs []Value) (result Value, err error) {
	args := pargs
	if l.isVariadric {
		if len(args) < l.arity-1 {
			return NIL, NewExecutionError(fmt.Sprintf("function %s expected at least %d args, got %d", l, l.arity-1, len(args)))
		}
		sargs := args[0 : l.arity-1]
		rest := args[l.arity-1:]
		restlist, boxErr := ListType.Box(rest)
		if boxErr != nil {
			return NIL, boxErr
		}
		args = append(sargs, restlist)
	} else if len(args) != l.arity {
		return NIL, NewExecutionError(fmt.Sprintf("function %s expected %d args, got %d", l, l.arity, len(args)))
	}
	f := NewFrame(l.chunk, args)
	f.ec = ec
	result, err = f.Run()
	ReleaseFrame(f)
	return result, err
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
}

func (l *Closure) Type() ValueType { return FuncType }

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
	if f, ok := l.fn.(*Func); ok {
		args := pargs
		if f.isVariadric {
			if len(args) < f.arity-1 {
				return NIL, NewExecutionError(fmt.Sprintf("function %s expected at least %d args, got %d", l, f.arity-1, len(args)))
			}
			sargs := args[0 : f.arity-1]
			rest := args[f.arity-1:]
			restlist, boxErr := ListType.Box(rest)
			if boxErr != nil {
				return NIL, boxErr
			}
			args = append(sargs, restlist)
		} else if len(args) != f.arity {
			return NIL, NewExecutionError(fmt.Sprintf("function %s expected %d args, got %d", l, f.arity, len(args)))
		}
		frame := NewFrame(f.chunk, args)
		frame.closedOvers = l.closedOvers
		frame.ec = ec
		result, err = frame.Run()
		ReleaseFrame(frame)
		return result, err
	}

	if mfn, ok := l.fn.(*MultiArityFn); ok {
		le := len(pargs)
		var variant Fn
		if f, ok := mfn.fns[le]; ok {
			variant = f
		} else if mfn.rest != nil && le >= mfn.rest.Arity() {
			variant = mfn.rest
		} else {
			return NIL, NewExecutionError(fmt.Sprintf("function %s doesn't have a %d-arity variant", l, le))
		}

		if f, ok := variant.(*Func); ok {
			subClosure := &Closure{
				closedOvers: l.closedOvers,
				fn:          f,
			}
			return subClosure.invokeIn(ec, pargs)
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
}

func (l *MultiArityFn) Type() ValueType { return FuncType }

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

// invokeIn runs the multi-arity function with the given ExecContext active,
// so dynamic bindings propagate into the selected variant's call. Invoke delegates
// to invokeIn against the root context.
func (l *MultiArityFn) invokeIn(ec *ExecContext, pargs []Value) (Value, error) {
	le := len(pargs)
	if f, ok := l.fns[le]; ok {
		return ec.Invoke(f, pargs)
	}
	if l.rest != nil && le >= l.rest.Arity() {
		return ec.Invoke(l.rest, pargs)
	}
	return NIL, NewExecutionError(fmt.Sprintf("function %s doesn't have a %d-arity variant", l, le))
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
