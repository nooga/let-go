/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"fmt"
	"reflect"
	"sync"
)

type theAtomType struct {
}

func (t *theAtomType) String() string     { return t.Name() }
func (t *theAtomType) Type() ValueType    { return TypeType }
func (t *theAtomType) Unbox() interface{} { return reflect.TypeOf(t) }

func (t *theAtomType) Name() string { return "let-go.lang.Atom" }
func (t *theAtomType) Box(b interface{}) (Value, error) {
	val, err := BoxValue(reflect.ValueOf(b))
	if err != nil {
		return NIL, err
	}
	return NewAtom(val), nil
}

var AtomType *theAtomType = &theAtomType{}

type Atom struct {
	root Value
	mu   sync.Mutex
}

func NewAtom(root Value) *Atom {
	return &Atom{
		root: root,
	}
}

func (v *Atom) Swap(fn Fn, args []Value) Value {
	v.mu.Lock()
	ret := fn.Invoke(append([]Value{v.root}, args...))
	v.root = ret
	v.mu.Unlock()
	return ret
}

func (v *Atom) Deref() Value {
	v.mu.Lock()
	val := v.root
	v.mu.Unlock()
	return val
}

func (v *Atom) Type() ValueType {
	return AtomType
}

func (v *Atom) Unbox() interface{} {
	return v.Deref().Unbox()
}

func (v *Atom) String() string {
	return fmt.Sprintf("<%s %s>", AtomType, v.Deref())
}
