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

func (t *theFuncType) String() string     { return t.Name() }
func (t *theFuncType) Type() ValueType    { return TypeType }
func (t *theFuncType) Unbox() interface{} { return reflect.TypeOf(t) }

func (t *theFuncType) Name() string { return "let-go.lang.Fn" }
func (t *theFuncType) Box(fn interface{}) (Value, error) {
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

type FuncInterface func(interface{})

// Unbox implements Unbox
func (l *Func) Unbox() interface{} {
	proxy := func(in []reflect.Value) []reflect.Value {
		args := make([]Value, len(in))
		for i := range in {
			a, _ := BoxValue(in[i]) // FIXME handle error
			args[i] = a
		}
		f := NewFrame(l.chunk, args)
		out, _ := f.Run() // FIXME handle error
		return []reflect.Value{reflect.ValueOf(out.Unbox())}
	}
	return func(fptr interface{}) {
		fn := reflect.ValueOf(fptr).Elem()
		v := reflect.MakeFunc(fn.Type(), proxy)
		fn.Set(v)
	}
}

func (l *Func) Arity() int {
	return l.arity
}

func (l *Func) Invoke(pargs []Value) (Value, error) {
	args := pargs
	if l.isVariadric {
		// pretty sure variadric should guarantee arity >= 1
		sargs := args[0 : l.arity-1]
		rest := args[l.arity-1:]
		restlist, err := ListType.Box(rest)
		if err != nil {
			return NIL, err
		}
		args = append(sargs, restlist)
	}
	f := NewFrame(l.chunk, args)
	return f.Run()
}

func (l *Func) String() string {
	if len(l.name) > 0 {
		return fmt.Sprintf("<fn %s %p>", l.name, l)
	}
	return fmt.Sprintf("<fn %p>", l)
}

func (l *Func) MakeClosure() Fn {
	return &Closure{
		closedOvers: nil,
		fn:          l,
	}
}

type Closure struct {
	closedOvers []Value
	fn          *Func
}

func (l *Closure) Type() ValueType { return FuncType }

// Unbox implements Unbox
func (l *Closure) Unbox() interface{} {
	proxy := func(in []reflect.Value) []reflect.Value {
		args := make([]Value, len(in))
		for i := range in {
			a, _ := BoxValue(in[i]) // FIXME handle error
			args[i] = a
		}
		f := NewFrame(l.fn.chunk, args)
		f.closedOvers = l.closedOvers
		out, _ := f.Run() // FIXME handle error
		return []reflect.Value{reflect.ValueOf(out.Unbox())}
	}
	return func(fptr interface{}) {
		fn := reflect.ValueOf(fptr).Elem()
		v := reflect.MakeFunc(fn.Type(), proxy)
		fn.Set(v)
	}
}

func (l *Closure) Arity() int {
	return l.fn.arity
}

func (l *Closure) Invoke(pargs []Value) (Value, error) {
	args := pargs
	if l.fn.isVariadric {
		// pretty sure variadric should guarantee arity >= 1
		sargs := args[0 : l.fn.arity-1]
		rest := args[l.fn.arity-1:]
		// FIXME don't swallow the error, make invoke return an error
		restlist, err := ListType.Box(rest)
		if err != nil {
			return NIL, err
		}
		args = append(sargs, restlist)
	}
	f := NewFrame(l.fn.chunk, args)
	f.closedOvers = l.closedOvers
	// FIXME don't swallow the error, make invoke return an error
	return f.Run()
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
func (l *MultiArityFn) Unbox() interface{} {
	proxy := func(in []reflect.Value) []reflect.Value {
		args := make([]Value, len(in))
		for i := range in {
			a, _ := BoxValue(in[i]) // FIXME handle error
			args[i] = a
		}
		out, _ := l.Invoke(args)
		return []reflect.Value{reflect.ValueOf(out.Unbox())}
	}
	return func(fptr interface{}) {
		fn := reflect.ValueOf(fptr).Elem()
		v := reflect.MakeFunc(fn.Type(), proxy)
		fn.Set(v)
	}
}

func (l *MultiArityFn) Arity() int {
	return l.arity
}

func (l *MultiArityFn) Invoke(pargs []Value) (Value, error) {
	le := len(pargs)
	if f, ok := l.fns[le]; ok {
		return f.Invoke(pargs)
	}
	if l.rest != nil && le >= l.rest.Arity() {
		return l.rest.Invoke(pargs)
	}
	return NIL, NewExecutionError(fmt.Sprintf("function %s doesn't have a %d-arity variant", l, le))
}

func (l *MultiArityFn) String() string {
	return fmt.Sprintf("<mfn %s %p>", l.name, l)
}

func makeMultiArity(fns []Value) (*MultiArityFn, error) {
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
		} else {
			ma.fns[a] = f
		}
	}
	return ma, nil
}
