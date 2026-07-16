/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"fmt"
	"maps"
	"reflect"
	"sync"
	"sync/atomic"
)

// guardedRootDeviations counts guarded vars (native-primitive roots, see
// GuardRoot) whose root currently differs from its canonical value. Lowered
// Go code consults GuardedRootsIntact() before taking a baked direct call to
// a native primitive; any active override (with-redefs / alter-var-root /
// re-def) flips the counter non-zero and every native direct-call site falls
// back to the var-dispatch trampoline, which observes the override. When the
// override is restored (with-redefs unwinding sets the original root value
// back), the counter returns to zero and the direct fast path resumes.
var guardedRootDeviations atomic.Int64

// GuardedRootsIntact reports that no guarded native-primitive var root is
// currently overridden. One atomic load — cheap enough for emitted per-call
// guards in the lowered tree.
func GuardedRootsIntact() bool { return guardedRootDeviations.Load() == 0 }

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
	// rootBind is this var's ROOT-context binding chain (Phase 2). Head = current
	// binding; nil = unbound. Per-var, so deref reads the head in O(1) with no
	// shared-structure access. Nesting lives in next; readers touch only the head.
	rootBind  atomic.Pointer[rootBindFrame]
	isPrivate bool
	mu        sync.Mutex // guards meta + watches
	watches   map[Value]Fn
	// guardExpected, when non-nil, holds the canonical root of a guarded
	// native-primitive var (see GuardRoot). guardDeviated tracks whether the
	// current root differs from it, mirrored into guardedRootDeviations.
	guardExpected atomic.Pointer[Value]
	guardDeviated atomic.Bool
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
	v.updateGuard(val)
	return v
}

// GuardRoot marks the var's CURRENT root as the canonical value of a native
// primitive. Any later root mutation to a different value (with-redefs,
// alter-var-root, re-def) bumps the global deviation counter, steering
// lowered direct-call sites onto the trampoline until the root is restored.
func (v *Var) GuardRoot() {
	cur := v.Root()
	v.guardExpected.Store(valPtr(cur))
	if v.guardDeviated.Swap(false) {
		guardedRootDeviations.Add(-1)
	}
}

// adoptGuard transfers guard state from an older Var replaced in a namespace
// registry by a re-def (Namespace.Def creates a fresh Var). The new var
// inherits the canonical root and is immediately re-evaluated against it.
func (v *Var) adoptGuard(old *Var) {
	if old == nil {
		return
	}
	p := old.guardExpected.Load()
	if p == nil {
		return
	}
	v.guardExpected.Store(p)
	// The replaced var no longer participates; drop its deviation, then
	// account for the new var's root against the canonical value.
	if old.guardDeviated.Swap(false) {
		guardedRootDeviations.Add(-1)
	}
	v.updateGuard(v.Root())
}

// updateGuard re-evaluates a guarded var's deviation state after a root
// mutation. No-op for unguarded vars (one atomic load).
func (v *Var) updateGuard(newVal Value) {
	p := v.guardExpected.Load()
	if p == nil {
		return
	}
	if sameRootIdentity(*p, newVal) {
		if v.guardDeviated.Swap(false) {
			guardedRootDeviations.Add(-1)
		}
	} else if !v.guardDeviated.Swap(true) {
		guardedRootDeviations.Add(1)
	}
}

// sameRootIdentity reports whether two root values are the same object.
// Interface == would panic on uncomparable dynamic types (funcs, maps), so
// gate on comparability; native adapters are *NativeFn pointers, for which
// this is exact pointer identity.
func sameRootIdentity(a, b Value) bool {
	ta, tb := reflect.TypeOf(a), reflect.TypeOf(b)
	if ta == nil || tb == nil || ta != tb || !ta.Comparable() {
		return ta == nil && tb == nil
	}
	return a == b
}

// derefRoot resolves v against the ROOT (process-global) context. It is the
// single shared root-resolution path — Var.Deref (host/lowered callers) and
// ExecContext.deref's root branch both call it, so the two can never diverge.
func (v *Var) derefRoot() Value {
	if val, ok := rootDerefHead(v); ok {
		return val
	}
	return v.Root()
}

// Deref returns the current value in the root execution context: the dynamic
// top binding if one is active, else the root. Host and lowered callers that
// hold no ExecContext resolve here; the interpreter resolves against its
// frame's context via ExecContext.deref (which shares derefRoot for the root).
func (v *Var) Deref() Value {
	return v.derefRoot()
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
	return rootSnapshot()
}

// RunWithBindings runs fn with snap installed as the root context's dynamic
// bindings, restoring the prior state afterwards. This is the legacy
// process-global bracketing used by spawn sites that have not yet moved to a
// child ExecContext; true per-goroutine isolation comes from running fn under
// ec.Child() instead (see docs/design/exec-context-threading.md).
func RunWithBindings(snap BindingSnapshot, fn func() (Value, error)) (Value, error) {
	saved := rootSnapshot()
	rootInstall(snap)
	out, err := fn()
	rootInstall(saved)
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
	v.updateGuard(result)
	// alter-var-root (and with-redefs, which is built on it) is a root
	// mutation path for *lg-trace* just like binding/set! — arm the trace
	// gate here too, or (alter-var-root #'*lg-trace* (constantly true))
	// silently never traces.
	armTraceIfTruthy(v, result)
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
