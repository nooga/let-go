/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package vm

import "sync"

// rootBindFrame is one immutable link in a var's ROOT-context binding chain.
// Frames are never mutated after publication, so readers load the head
// atomically and read val without locking. Nesting (shadowed bindings) lives in
// next; the head is always the current binding.
type rootBindFrame struct {
	val  Value
	next *rootBindFrame
}

// rootRegistry is the set of vars that currently have a live root binding. It
// exists so snapshot() can enumerate live root bindings without a shared chain
// to walk. Written ONLY on bind/unbind transitions (first bind adds, last
// unbind removes) and read ONLY by snapshot — the deref hot path never touches
// it. A single mutex is enough: writes are rare (binding establishment, not
// reads) and concurrent ROOT binding is itself unusual (thread-local binding
// uses child contexts, which never touch this registry).
var (
	rootRegistryMu sync.Mutex
	rootRegistry   = map[*Var]struct{}{}
)

func registryAdd(v *Var) {
	rootRegistryMu.Lock()
	rootRegistry[v] = struct{}{}
	rootRegistryMu.Unlock()
}

func registryRemove(v *Var) {
	rootRegistryMu.Lock()
	delete(rootRegistry, v)
	rootRegistryMu.Unlock()
}

// rootPush makes val v's current root binding. O(1). Serialized on v.mu with
// pop/setCurrent. Registers v on its first live binding. Note: v.mu is released
// BEFORE touching the registry so we never nest v.mu inside rootRegistryMu
// (rootInstall nests the other way).
func rootPush(v *Var, val Value) {
	v.mu.Lock()
	old := v.rootBind.Load()
	v.rootBind.Store(&rootBindFrame{val: val, next: old})
	v.mu.Unlock()
	if old == nil {
		registryAdd(v)
	}
}

// rootPop removes v's current root binding, restoring the shadowed one. O(1).
// Deregisters v when its last binding is popped.
func rootPop(v *Var) {
	v.mu.Lock()
	cur := v.rootBind.Load()
	var next *rootBindFrame
	if cur != nil {
		next = cur.next
	}
	v.rootBind.Store(next)
	v.mu.Unlock()
	if cur != nil && next == nil {
		registryRemove(v)
	}
}

// rootDerefHead returns (value, true) if v has a live root binding, else
// (_, false). O(1), lock-free — the hot deref path.
func rootDerefHead(v *Var) (Value, bool) {
	if f := v.rootBind.Load(); f != nil {
		return f.val, true
	}
	return nil, false
}

// rootHasBinding reports whether v has a live root binding.
func rootHasBinding(v *Var) bool { return v.rootBind.Load() != nil }

// rootSetCurrent replaces v's current root binding value in place (Clojure
// set!), returning false if v has none. Rebuilds the head frame so frames stay
// immutable for lock-free readers.
func rootSetCurrent(v *Var, val Value) bool {
	v.mu.Lock()
	defer v.mu.Unlock()
	cur := v.rootBind.Load()
	if cur == nil {
		return false
	}
	v.rootBind.Store(&rootBindFrame{val: val, next: cur.next})
	return true
}

// rootSnapshot captures all live root bindings as a BindingSnapshot
// (map[*Var][]Value, top = last). Reads the registry under lock, then each
// member's chain. A binding racing establishment/teardown may be included or
// omitted — snapshotting concurrently with binding is inherently racy, and each
// captured chain is internally consistent (atomic head load).
func rootSnapshot() BindingSnapshot {
	rootRegistryMu.Lock()
	vars := make([]*Var, 0, len(rootRegistry))
	for v := range rootRegistry {
		vars = append(vars, v)
	}
	rootRegistryMu.Unlock()

	out := BindingSnapshot{}
	for _, v := range vars {
		var vals []Value // head→tail = innermost→outermost
		for f := v.rootBind.Load(); f != nil; f = f.next {
			vals = append(vals, f.val)
		}
		if len(vals) == 0 {
			continue // raced an unbind
		}
		for i, j := 0, len(vals)-1; i < j; i, j = i+1, j-1 { // reverse → top-last
			vals[i], vals[j] = vals[j], vals[i]
		}
		out[v] = vals
	}
	return out
}

// rootInstall replaces the entire root binding state with target: sets each
// var's chain from target[v] (top-last) and fixes registry membership. Used by
// RunWithBindings' whole-state swap. Holds rootRegistryMu across the swap so
// membership stays consistent; takes each v.mu inside it (the one allowed
// nesting direction — see rootPush).
func rootInstall(target BindingSnapshot) {
	rootRegistryMu.Lock()
	defer rootRegistryMu.Unlock()

	// Clear vars currently bound but absent from target.
	for v := range rootRegistry {
		if _, ok := target[v]; !ok {
			v.mu.Lock()
			v.rootBind.Store(nil)
			v.mu.Unlock()
			delete(rootRegistry, v)
		}
	}
	// Set each target var's chain (bottom-to-top so the last value is head).
	for v, vals := range target {
		var head *rootBindFrame
		for _, val := range vals {
			head = &rootBindFrame{val: val, next: head}
		}
		v.mu.Lock()
		v.rootBind.Store(head)
		v.mu.Unlock()
		if head != nil {
			rootRegistry[v] = struct{}{}
		} else {
			delete(rootRegistry, v)
		}
	}
}
