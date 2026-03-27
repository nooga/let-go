/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import "fmt"

// MultiFn implements Clojure-style multimethods.
// It holds a dispatch function and a map of dispatch-value → method.
type MultiFn struct {
	name       string
	dispatchFn Fn
	methods    *PersistentMap
	defaultVal Value // dispatch value for the default method
}

func NewMultiFn(name string, dispatchFn Fn, defaultVal Value) *MultiFn {
	return &MultiFn{
		name:       name,
		dispatchFn: dispatchFn,
		methods:    EmptyPersistentMap,
		defaultVal: defaultVal,
	}
}

func (m *MultiFn) Type() ValueType    { return MultiFnType }
func (m *MultiFn) Unbox() interface{} { return m }

func (m *MultiFn) String() string {
	return fmt.Sprintf("<multifn %s>", m.name)
}

// AddMethod registers an implementation for a dispatch value.
func (m *MultiFn) AddMethod(dispatchVal Value, method Fn) *MultiFn {
	return &MultiFn{
		name:       m.name,
		dispatchFn: m.dispatchFn,
		methods:    m.methods.Assoc(dispatchVal, method).(*PersistentMap),
		defaultVal: m.defaultVal,
	}
}

// RemoveMethod unregisters an implementation.
func (m *MultiFn) RemoveMethod(dispatchVal Value) *MultiFn {
	return &MultiFn{
		name:       m.name,
		dispatchFn: m.dispatchFn,
		methods:    m.methods.Dissoc(dispatchVal).(*PersistentMap),
		defaultVal: m.defaultVal,
	}
}

// Arity returns -1 (variadic — arity depends on the method).
func (m *MultiFn) Arity() int { return -1 }

// Invoke calls the dispatch function, looks up the method, and calls it.
func (m *MultiFn) Invoke(args []Value) (Value, error) {
	// Call dispatch function
	dv, err := m.dispatchFn.Invoke(args)
	if err != nil {
		return NIL, fmt.Errorf("multimethod %s dispatch failed: %w", m.name, err)
	}

	// Look up method for dispatch value
	method := m.methods.ValueAt(dv)
	if method == NIL {
		// Try default
		method = m.methods.ValueAt(m.defaultVal)
		if method == NIL {
			return NIL, fmt.Errorf("no method in multimethod '%s' for dispatch value: %s", m.name, dv)
		}
	}

	fn, ok := method.(Fn)
	if !ok {
		return NIL, fmt.Errorf("multimethod '%s' method is not a function", m.name)
	}

	return fn.Invoke(args)
}

// Methods returns the method map.
func (m *MultiFn) Methods() *PersistentMap {
	return m.methods
}

type theMultiFnType struct{}

func (t *theMultiFnType) String() string     { return t.Name() }
func (t *theMultiFnType) Type() ValueType    { return TypeType }
func (t *theMultiFnType) Unbox() interface{} { return nil }
func (t *theMultiFnType) Name() string       { return "let-go.lang.MultiFn" }
func (t *theMultiFnType) Box(bare interface{}) (Value, error) {
	return NIL, NewTypeError(bare, "can't be boxed as", t)
}

var MultiFnType *theMultiFnType = &theMultiFnType{}
