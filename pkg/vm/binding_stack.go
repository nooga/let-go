/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"sync"
	"sync/atomic"
)

// bindingStack is the dynamic-var binding state owned by a single ExecContext.
// It is a persistent stack of immutable (var, value) frames behind an atomic
// head pointer. push prepends a frame; the topmost frame for a var is its
// current dynamic binding, and the frames beneath it are that var's shadowed
// outer bindings.
//
// Because frames are never mutated after publication, a reader walks the chain
// from an atomically-loaded head without locking — so the hot deref path of
// every ^:dynamic var (which must consult the stack on each read, since the var
// is *declared* dynamic whether or not it currently has a binding) stays off
// the mutex. Writers (push/pop/setCurrent/installSnapshot) serialize on mu,
// build a new chain that shares the unchanged tail, and atomic-swap the head,
// so a reader sees a consistent old-or-new chain, never a torn state. The mutex
// only serializes the rare case of two writers sharing one context (a value
// escaping to a helper goroutine).
//
// Versus a copy-on-write map, push is one allocation and O(1), and a child
// context can share the parent's chain by pointer; the cost is that reads walk
// the active bindings (typically one to a few) rather than hashing.
type bindingStack struct {
	mu   sync.Mutex // serializes writers
	head atomic.Pointer[bindingFrame]
}

// bindingFrame is one immutable link in a binding stack: a (var, value) binding
// over the frame below it.
type bindingFrame struct {
	v    *Var
	val  Value
	next *bindingFrame
}

func newBindingStack() *bindingStack { return &bindingStack{} }

func (b *bindingStack) push(v *Var, val Value) {
	b.mu.Lock()
	b.head.Store(&bindingFrame{v: v, val: val, next: b.head.Load()})
	b.mu.Unlock()
}

func (b *bindingStack) pop(v *Var) {
	b.mu.Lock()
	b.head.Store(removeTopBinding(b.head.Load(), v))
	b.mu.Unlock()
}

func (b *bindingStack) current(v *Var) (Value, bool) {
	for f := b.head.Load(); f != nil; f = f.next {
		if f.v == v {
			return f.val, true
		}
	}
	return nil, false
}

// setCurrent replaces the value of v's top dynamic binding, returning true if a
// binding existed. This is the (set! *v* val) primitive: it mutates only THIS
// context's top frame for v (thread-local in Clojure terms) and never the root.
// Implemented immutably — it rebuilds the frames above v's topmost binding over
// a replacement frame — so a concurrent reader on the old chain is unaffected
// and child contexts that share frames via snapshot stay isolated.
func (b *bindingStack) setCurrent(v *Var, val Value) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	nh, ok := replaceTopBinding(b.head.Load(), v, val)
	if ok {
		b.head.Store(nh)
	}
	return ok
}

func (b *bindingStack) hasBinding(v *Var) bool {
	for f := b.head.Load(); f != nil; f = f.next {
		if f.v == v {
			return true
		}
	}
	return false
}

// removeTopBinding returns head with the topmost frame for v spliced out,
// rebuilding only the frames above it (each a fresh frame) and sharing the tail
// below. If v is not bound, head is returned unchanged (no allocation).
func removeTopBinding(head *bindingFrame, v *Var) *bindingFrame {
	if head == nil {
		return nil
	}
	if head.v == v {
		return head.next
	}
	rest := removeTopBinding(head.next, v)
	if rest == head.next {
		return head
	}
	return &bindingFrame{v: head.v, val: head.val, next: rest}
}

// replaceTopBinding returns head with v's topmost frame replaced by one holding
// val, rebuilding the frames above it. ok is false (and head unchanged) if v is
// not bound.
func replaceTopBinding(head *bindingFrame, v *Var, val Value) (*bindingFrame, bool) {
	if head == nil {
		return nil, false
	}
	if head.v == v {
		return &bindingFrame{v: v, val: val, next: head.next}, true
	}
	rest, ok := replaceTopBinding(head.next, v, val)
	if !ok {
		return head, false
	}
	return &bindingFrame{v: head.v, val: head.val, next: rest}, true
}

// snapshot converts the frame chain to the public BindingSnapshot map
// (*Var -> value stack, top = last). Off the hot path; used at goroutine
// boundaries to seed a child context. Lock-free: the head load yields a
// consistent immutable chain.
func (b *bindingStack) snapshot() BindingSnapshot {
	out := BindingSnapshot{}
	for f := b.head.Load(); f != nil; f = f.next {
		// head is the innermost binding; prepend so the current (innermost)
		// value lands last, matching the map's top-is-last convention.
		out[f.v] = append([]Value{f.val}, out[f.v]...)
	}
	return out
}

// installSnapshot replaces the chain with one rebuilt from snap. Each var's
// values are pushed bottom-to-top so its current (last) value ends up topmost;
// order across distinct vars is irrelevant to resolution.
func (b *bindingStack) installSnapshot(snap BindingSnapshot) {
	var head *bindingFrame
	for v, stack := range snap {
		for _, val := range stack {
			head = &bindingFrame{v: v, val: val, next: head}
		}
	}
	b.mu.Lock()
	b.head.Store(head)
	b.mu.Unlock()
}
