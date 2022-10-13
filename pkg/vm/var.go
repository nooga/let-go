/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import "fmt"

type Var struct {
	root    Value
	nsref   *Namespace
	ns      string
	name    string
	isMacro bool
}

func (v *Var) Invoke(values []Value) (Value, error) {
	f, ok := v.root.(Fn)
	if !ok {
		return NIL, fmt.Errorf("%v root does not implement Fn") // FIXME this should be an error
	}
	return f.Invoke(values)
}

func (v *Var) Arity() int {
	f, ok := v.root.(Fn)
	if !ok {
		return 0 // FIXME this should be an error
	}
	return f.Arity()
}

func NewVar(nsref *Namespace, ns string, name string) *Var {
	return &Var{
		nsref:   nsref,
		ns:      ns,
		name:    name,
		root:    NIL,
		isMacro: false,
	}
}

func (v *Var) SetRoot(val Value) *Var {
	v.root = val
	return v
}

func (v *Var) Deref() Value {
	return v.root
}

func (v *Var) Type() ValueType {
	return v.Deref().Type()
}

func (v *Var) Unbox() interface{} {
	return v.Deref().Unbox()
}

func (v *Var) String() string {
	return fmt.Sprintf("#'%s/%s", v.ns, v.name)
}

func (v *Var) IsMacro() bool {
	return v.isMacro
}

func (v *Var) SetMacro() {
	v.isMacro = true
}
