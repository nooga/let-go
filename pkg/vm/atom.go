/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
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

// Atom is a thread-safe mutable reference.
// Swap uses optimistic concurrency with a generation counter — no value comparison needed.
// The function may be called multiple times under contention.
type Atom struct {
	val      Value
	gen      uint64 // generation counter — incremented on every mutation
	mu       sync.Mutex
	meta     Value
	watches  map[Value]Fn // key → watch fn
}

func NewAtom(root Value) *Atom {
	return &Atom{val: root}
}

func (a *Atom) notifyWatches(oldVal, newVal Value) {
	if len(a.watches) == 0 {
		return
	}
	for key, fn := range a.watches {
		fn.Invoke([]Value{key, a, oldVal, newVal})
	}
}

func (a *Atom) AddWatch(key Value, fn Fn) {
	a.mu.Lock()
	if a.watches == nil {
		a.watches = make(map[Value]Fn)
	}
	a.watches[key] = fn
	a.mu.Unlock()
}

func (a *Atom) RemoveWatch(key Value) {
	a.mu.Lock()
	delete(a.watches, key)
	a.mu.Unlock()
}

func (a *Atom) Meta() Value {
	if a.meta == nil {
		return NIL
	}
	return a.meta
}

func (a *Atom) WithMeta(m Value) Value {
	a.mu.Lock()
	defer a.mu.Unlock()
	return &Atom{val: a.val, gen: a.gen, meta: m, watches: a.watches}
}

func (a *Atom) AlterMeta(fn Fn, args []Value) (Value, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	m := a.meta
	if m == nil {
		m = NIL
	}
	allArgs := append([]Value{m}, args...)
	newMeta, err := fn.Invoke(allArgs)
	if err != nil {
		return NIL, err
	}
	a.meta = newMeta
	return newMeta, nil
}

func (a *Atom) Reset(newVal Value) Value {
	a.mu.Lock()
	oldVal := a.val
	a.val = newVal
	a.gen++
	watches := a.watches
	a.mu.Unlock()
	if len(watches) > 0 {
		a.notifyWatches(oldVal, newVal)
	}
	return newVal
}

// Swap applies fn to the current value and atomically sets the result.
// The fn is called outside the lock; if the value changed during computation,
// fn is retried with the new value (like Clojure's swap!).
func (a *Atom) Swap(fn Fn, args []Value) (Value, error) {
	for {
		// Snapshot current value and generation
		a.mu.Lock()
		oldVal := a.val
		oldGen := a.gen
		a.mu.Unlock()

		// Compute new value without holding the lock
		newVal, err := fn.Invoke(append([]Value{oldVal}, args...))
		if err != nil {
			return NIL, err
		}

		// Try to set — only if generation hasn't changed
		a.mu.Lock()
		if a.gen == oldGen {
			a.val = newVal
			a.gen++
			watches := a.watches
			a.mu.Unlock()
			if len(watches) > 0 {
				a.notifyWatches(oldVal, newVal)
			}
			return newVal, nil
		}
		a.mu.Unlock()
		// Generation changed — another goroutine mutated, retry
	}
}

func (a *Atom) Deref() Value {
	a.mu.Lock()
	v := a.val
	a.mu.Unlock()
	return v
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
