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
	// root is atomic so Deref — by far the hottest var operation — is
	// lock-free. Dynamic (thread-local) bindings no longer live on the Var;
	// they are held by the ExecContext binding stack and resolved via
	// RootExecContext.deref / ec.deref. Only the process-wide root binding
	// remains here.
	root    atomic.Pointer[Value] // root binding (lock-free read)
	nsref   *Namespace
	ns      string
	name    string
	meta    Value
	isMacro bool
	// isDynamic is atomic: push-binding/deref of a var can happen on different
	// goroutines concurrently (e.g. two futures both `(binding [*v* ...] ...)`),
	// so the dynamic flag is read on the hot deref path while being set by a
	// concurrent bind. A plain bool here is a data race.
	isDynamic atomic.Bool
	isPrivate bool
	mu        sync.Mutex // guards meta + watches
	watches   map[Value]Fn
}

// valPtr boxes a Value for storage in an atomic.Pointer[Value].
func valPtr(v Value) *Value { return &v }

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

// Deref returns the current value in the root execution context: the dynamic
// top binding if one is active, else the root. Host and lowered callers that
// hold no ExecContext resolve here; the interpreter resolves against its
// frame's context (which is the root context unless a child was installed).
func (v *Var) Deref() Value {
	return RootExecContext.deref(v)
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
	return RootExecContext.hasBinding(v) || v.root.Load() != nil
}

// PushBinding pushes a dynamic binding value in the root execution context.
func (v *Var) PushBinding(val Value) {
	RootExecContext.pushBinding(v, val)
}

// PopBinding removes the most recent dynamic binding in the root context.
func (v *Var) PopBinding() {
	RootExecContext.popBinding(v)
}

// SnapshotBindings captures the root context's dynamic bindings — the
// explicit-propagation primitive a goroutine spawn hands to its child.
func SnapshotBindings() BindingSnapshot {
	return globalBindingStack.snapshot()
}

// RunWithBindings runs fn with snap installed as the root context's dynamic
// bindings, restoring the prior state afterwards. This is the legacy
// process-global bracketing used by spawn sites that have not yet moved to a
// child ExecContext; true per-goroutine isolation comes from running fn under
// ec.Child() instead (see docs/design/exec-context-threading.md).
func RunWithBindings(snap BindingSnapshot, fn func() (Value, error)) (Value, error) {
	saved := globalBindingStack.snapshot()
	globalBindingStack.installSnapshot(snap)
	out, err := fn()
	globalBindingStack.installSnapshot(saved)
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
	return v.isDynamic.Load()
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
	v.isDynamic.Store(true)
}

func (v *Var) SetPrivate() {
	v.isPrivate = true
}
