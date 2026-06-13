/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package vm

import "sync"

// bindingStack is the dynamic-var binding state owned by a single
// ExecContext. Unlike the process-global store it replaces, it is keyed by
// nothing — it belongs to one execution, so per-goroutine isolation comes from
// each goroutine carrying its own ExecContext, with no goroutine-id lookup and
// no cross-goroutine reuse hazard. The mutex guards against a value escaping to
// a helper goroutine that shares the same context (rare, but cheap to be safe).
type bindingStack struct {
	mu       sync.Mutex
	bindings BindingSnapshot // *Var -> stack of bound values (top = last)
}

func newBindingStack() *bindingStack {
	return &bindingStack{bindings: make(BindingSnapshot)}
}

func (b *bindingStack) push(v *Var, val Value) {
	b.mu.Lock()
	b.bindings[v] = append(b.bindings[v], val)
	b.mu.Unlock()
}

func (b *bindingStack) pop(v *Var) {
	b.mu.Lock()
	stack := b.bindings[v]
	if n := len(stack); n > 0 {
		if n == 1 {
			delete(b.bindings, v)
		} else {
			b.bindings[v] = stack[:n-1]
		}
	}
	b.mu.Unlock()
}

func (b *bindingStack) current(v *Var) (Value, bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if stack := b.bindings[v]; len(stack) > 0 {
		return stack[len(stack)-1], true
	}
	return nil, false
}

// setCurrent replaces the value of v's top dynamic binding in place, returning
// true if a binding existed. This is the (set! *v* val) primitive: it mutates
// only THIS context's top frame (thread-local in Clojure terms) and never the
// root. A child context's frame is its own copy (Child snapshots it), so the
// mutation stays isolated to this execution and does not leak to siblings or
// the parent.
func (b *bindingStack) setCurrent(v *Var, val Value) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	if stack := b.bindings[v]; len(stack) > 0 {
		stack[len(stack)-1] = val
		return true
	}
	return false
}

func (b *bindingStack) hasBinding(v *Var) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return len(b.bindings[v]) > 0
}

func (b *bindingStack) snapshot() BindingSnapshot {
	b.mu.Lock()
	defer b.mu.Unlock()
	out := make(BindingSnapshot, len(b.bindings))
	for v, stack := range b.bindings {
		out[v] = append([]Value(nil), stack...)
	}
	return out
}

func (b *bindingStack) installSnapshot(snap BindingSnapshot) {
	b.mu.Lock()
	out := make(BindingSnapshot, len(snap))
	for v, stack := range snap {
		out[v] = append([]Value(nil), stack...)
	}
	b.bindings = out
	b.mu.Unlock()
}
