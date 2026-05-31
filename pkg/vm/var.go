/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"fmt"
	"maps"
	"sync"
	"sync/atomic"
)

type Var struct {
	// root and curr are atomic so Deref — by far the hottest var
	// operation — is lock-free. Previously every Deref (even of a var
	// with no dynamic binding) took the global bindingsMu just to check
	// the binding stack, serializing all var reads across goroutines.
	// curr holds the current dynamic top binding (nil = use root); it is
	// kept in sync with the bindings stack by push/pop/RunWithBindings
	// under bindingsMu (the cold path).
	root      atomic.Pointer[Value] // root binding (lock-free read)
	curr      atomic.Pointer[Value] // current dynamic top binding; nil = use root
	bindings  []Value               // full dynamic binding stack (guarded by bindingsMu)
	nsref     *Namespace
	ns        string
	name      string
	meta      Value
	isMacro   bool
	isDynamic bool
	isPrivate bool
	mu        sync.Mutex // guards meta + watches
	watches   map[Value]Fn
}

var (
	bindingsMu sync.Mutex
	activeVars = map[*Var]struct{}{}
)

// valPtr boxes a Value for storage in an atomic.Pointer[Value].
func valPtr(v Value) *Value { return &v }

// syncCurrLocked refreshes the atomic current-binding pointer from the
// bindings stack. Must be called while holding bindingsMu.
func (v *Var) syncCurrLocked() {
	if len(v.bindings) == 0 {
		v.curr.Store(nil)
	} else {
		v.curr.Store(valPtr(v.bindings[len(v.bindings)-1]))
	}
}

type BindingSnapshot map[*Var][]Value

func (v *Var) Invoke(values []Value) (Value, error) {
	root := v.Root()
	f, ok := root.(Fn)
	if !ok {
		return NIL, fmt.Errorf("%v root does not implement Fn", root)
	}
	return f.Invoke(values)
}

func (v *Var) Arity() int {
	f, ok := v.Root().(Fn)
	if !ok {
		return 0
	}
	return f.Arity()
}

func NewVar(nsref *Namespace, ns string, name string) *Var {
	v := &Var{
		nsref: nsref,
		ns:    ns,
		name:  name,
	}
	// Leave the root pointer unset (unbound). Deref/Root already fall back to
	// NIL when no root is stored, so reads are unchanged, but HasRoot() now
	// correctly reports false until a real assignment — which is what `defonce`
	// needs to tell a "compiler-interned forward ref" from "actually defined".
	return v
}

func (v *Var) SetRoot(val Value) *Var {
	v.root.Store(valPtr(val))
	return v
}

// Deref returns the current value: the dynamic top binding if one is
// active, else the root. Lock-free — two atomic loads, no mutex — so it
// scales across goroutines.
func (v *Var) Deref() Value {
	if p := v.curr.Load(); p != nil {
		return *p
	}
	if p := v.root.Load(); p != nil {
		return *p
	}
	return NIL
}

// Root returns the var's root binding directly, bypassing any current
// dynamic binding on the stack. Use this where Clojure semantics require
// the root (e.g. alter-var-root) rather than the currently visible deref
// value.
func (v *Var) Root() Value {
	if p := v.root.Load(); p != nil {
		return *p
	}
	return NIL
}

// IsBound reports whether the var has any bound value — a root binding OR an
// active dynamic binding — matching Clojure's bound?. A var interned by the
// compiler for a forward `(def x v)` (before it runs) has neither yet, which
// distinguishes "declared but unset" from "set", as `defonce` needs.
func (v *Var) IsBound() bool {
	return v.curr.Load() != nil || v.root.Load() != nil
}

// PushBinding pushes a dynamic binding value.
func (v *Var) PushBinding(val Value) {
	bindingsMu.Lock()
	defer bindingsMu.Unlock()
	v.bindings = append(v.bindings, val)
	v.syncCurrLocked()
	activeVars[v] = struct{}{}
}

// PopBinding removes the most recent dynamic binding.
func (v *Var) PopBinding() {
	bindingsMu.Lock()
	defer bindingsMu.Unlock()
	if len(v.bindings) > 0 {
		v.bindings = v.bindings[:len(v.bindings)-1]
	}
	v.syncCurrLocked()
	if len(v.bindings) == 0 {
		delete(activeVars, v)
	}
}

func SnapshotBindings() BindingSnapshot {
	bindingsMu.Lock()
	defer bindingsMu.Unlock()
	snap := BindingSnapshot{}
	for v := range activeVars {
		if len(v.bindings) == 0 {
			continue
		}
		bs := make([]Value, len(v.bindings))
		copy(bs, v.bindings)
		snap[v] = bs
	}
	return snap
}

func RunWithBindings(snap BindingSnapshot, fn func() (Value, error)) (Value, error) {
	bindingsMu.Lock()
	saved := BindingSnapshot{}
	for v := range activeVars {
		bs := make([]Value, len(v.bindings))
		copy(bs, v.bindings)
		saved[v] = bs
	}
	for v := range snap {
		if _, ok := saved[v]; !ok {
			saved[v] = nil
		}
	}
	for v := range saved {
		if bs, ok := snap[v]; ok {
			v.bindings = append([]Value(nil), bs...)
			if len(v.bindings) > 0 {
				activeVars[v] = struct{}{}
			} else {
				delete(activeVars, v)
			}
		} else {
			v.bindings = nil
			delete(activeVars, v)
		}
		v.syncCurrLocked()
	}
	bindingsMu.Unlock()

	out, err := fn()

	bindingsMu.Lock()
	for v, bs := range saved {
		v.bindings = append([]Value(nil), bs...)
		if len(v.bindings) > 0 {
			activeVars[v] = struct{}{}
		} else {
			delete(activeVars, v)
		}
		v.syncCurrLocked()
	}
	bindingsMu.Unlock()
	return out, err
}

func (v *Var) notifyWatches(oldVal, newVal Value) error {
	v.mu.Lock()
	if len(v.watches) == 0 {
		v.mu.Unlock()
		return nil
	}
	watches := make(map[Value]Fn, len(v.watches))
	maps.Copy(watches, v.watches)
	v.mu.Unlock()
	for key, fn := range watches {
		if _, err := fn.Invoke([]Value{key, v, oldVal, newVal}); err != nil {
			return err
		}
	}
	return nil
}

func (v *Var) AlterRoot(fn Fn) (Value, error) {
	return v.AlterRootArgs(fn, nil)
}

func (v *Var) AlterRootArgs(fn Fn, args []Value) (Value, error) {
	old := v.Root()
	result, err := fn.Invoke(append([]Value{old}, args...))
	if err != nil {
		return NIL, err
	}
	v.root.Store(valPtr(result))
	if err := v.notifyWatches(old, result); err != nil {
		return NIL, err
	}
	return result, nil
}

func (v *Var) AddWatch(key Value, fn Fn) {
	v.mu.Lock()
	if v.watches == nil {
		v.watches = make(map[Value]Fn)
	}
	v.watches[key] = fn
	v.mu.Unlock()
}

func (v *Var) RemoveWatch(key Value) {
	v.mu.Lock()
	delete(v.watches, key)
	v.mu.Unlock()
}

func (v *Var) Type() ValueType {
	return v.Deref().Type()
}

func (v *Var) Unbox() any {
	return v.Deref().Unbox()
}

func (v *Var) String() string {
	return fmt.Sprintf("#'%s/%s", v.ns, v.name)
}

func (v *Var) IsMacro() bool {
	return v.isMacro
}

func (v *Var) IsDynamic() bool {
	return v.isDynamic
}

func (v *Var) IsPrivate() bool {
	return v.isPrivate
}

func (v *Var) Meta() Value {
	v.mu.Lock()
	defer v.mu.Unlock()
	if v.meta == nil {
		return NIL
	}
	return v.meta
}

func (v *Var) SetMeta(meta Value) {
	v.mu.Lock()
	v.meta = meta
	v.mu.Unlock()
}

func (v *Var) AlterMeta(fn Fn, args []Value) (Value, error) {
	v.mu.Lock()
	meta := v.meta
	if meta == nil {
		meta = NIL
	}
	v.mu.Unlock()

	newMeta, err := fn.Invoke(append([]Value{meta}, args...))
	if err != nil {
		return NIL, err
	}
	v.SetMeta(newMeta)
	return newMeta, nil
}

// NS returns the namespace name.
func (v *Var) NS() string { return v.ns }

// VarName returns the var name.
func (v *Var) VarName() string { return v.name }

func (v *Var) SetMacro() {
	v.isMacro = true
}

func (v *Var) SetDynamic() {
	v.isDynamic = true
}

func (v *Var) SetPrivate() {
	v.isPrivate = true
}
