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

func (t *theAtomType) String() string  { return t.Name() }
func (t *theAtomType) Type() ValueType { return TypeType }
func (t *theAtomType) Unbox() any      { return reflect.TypeFor[*theAtomType]() }

func (t *theAtomType) Name() string { return "let-go.lang.Atom" }
func (t *theAtomType) Box(b any) (Value, error) {
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
	val       Value
	gen       uint64 // generation counter — incremented on every mutation
	mu        sync.Mutex
	meta      Value
	validator Fn
	watches   map[Value]Fn // key → watch fn
}

func NewAtom(root Value) *Atom {
	return &Atom{val: root}
}

func NewAtomWithMetaValidator(root Value, meta Value, validator Fn) (*Atom, error) {
	a := &Atom{val: root, meta: meta, validator: validator}
	if err := a.validate(root); err != nil {
		return nil, err
	}
	return a, nil
}

func (a *Atom) validate(newVal Value) error {
	if a.validator == nil {
		return nil
	}
	result, err := a.validator.Invoke([]Value{newVal})
	if err != nil {
		return err
	}
	if !IsTruthy(result) {
		return fmt.Errorf("validator rejected reference state")
	}
	return nil
}

func (a *Atom) notifyWatches(oldVal, newVal Value) error {
	if len(a.watches) == 0 {
		return nil
	}
	for key, fn := range a.watches {
		if _, err := fn.Invoke([]Value{key, a, oldVal, newVal}); err != nil {
			return err
		}
	}
	return nil
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

func (a *Atom) Validator() Value {
	if a.validator == nil {
		return NIL
	}
	return a.validator
}

func (a *Atom) WithMeta(m Value) Value {
	a.mu.Lock()
	defer a.mu.Unlock()
	return &Atom{val: a.val, gen: a.gen, meta: m, validator: a.validator, watches: a.watches}
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

func (a *Atom) Reset(newVal Value) (Value, error) {
	if err := a.validate(newVal); err != nil {
		return NIL, err
	}
	a.mu.Lock()
	oldVal := a.val
	a.val = newVal
	a.gen++
	watches := a.watches
	a.mu.Unlock()
	if len(watches) > 0 {
		if err := a.notifyWatches(oldVal, newVal); err != nil {
			return NIL, err
		}
	}
	return newVal, nil
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
		if err := a.validate(newVal); err != nil {
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
				if err := a.notifyWatches(oldVal, newVal); err != nil {
					return NIL, err
				}
			}
			return newVal, nil
		}
		a.mu.Unlock()
		// Generation changed — another goroutine mutated, retry
	}
}

// CompareAndSet atomically sets the value to newVal only if the current value
// is identical to oldVal, returning true on success. Identity comparison (==),
// matching Clojure's compare-and-set!.
func (a *Atom) CompareAndSet(oldVal, newVal Value) (bool, error) {
	if err := a.validate(newVal); err != nil {
		return false, err
	}
	a.mu.Lock()
	if a.val != oldVal {
		a.mu.Unlock()
		return false, nil
	}
	prev := a.val
	a.val = newVal
	a.gen++
	watches := a.watches
	a.mu.Unlock()
	if len(watches) > 0 {
		if err := a.notifyWatches(prev, newVal); err != nil {
			return false, err
		}
	}
	return true, nil
}

func (a *Atom) Deref() Value {
	a.mu.Lock()
	v := a.val
	a.mu.Unlock()
	return v
}

func (a *Atom) Type() ValueType {
	return AtomType
}

func (a *Atom) Unbox() any {
	return a
}

func (a *Atom) Hash() uint32 {
	return hashString(fmt.Sprintf("%p", a))
}

func (a *Atom) String() string {
	return fmt.Sprintf("<%s %s>", AtomType, a.Deref())
}
