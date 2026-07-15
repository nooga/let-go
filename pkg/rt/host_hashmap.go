/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"fmt"
	"sync/atomic"

	"github.com/nooga/let-go/pkg/vm"
)

// theHostHashMapType is the ValueType for the mutable java.util.HashMap compat
// shim. It is a JVM-interop shim (mutable, not a persistent let-go value), so it
// lives in the compat layer rather than in pkg/vm beside the real collections.
type theHostHashMapType struct{}

func (t *theHostHashMapType) String() string     { return t.Name() }
func (t *theHostHashMapType) Type() vm.ValueType { return vm.TypeType }
func (t *theHostHashMapType) Unbox() any         { return nil }
func (t *theHostHashMapType) Name() string       { return "java.util.HashMap" }
func (t *theHostHashMapType) Box(any) (vm.Value, error) {
	return vm.NIL, fmt.Errorf("java.util.HashMap cannot be boxed")
}

// hostHashMapType is the singleton type for hostHashMap values.
var hostHashMapType = &theHostHashMapType{}

// hostHashMap is a mutable map backing (java.util.HashMap.). malli.registry's
// fast-registry builds one with (doto (HashMap. n f) (.putAll m)) and reads it
// with (.get fm type) — so only putAll and get are exercised. The backing map is
// an immutable let-go map held in an atomic.Value: putAll accumulates into a new
// map and swaps it in, get reads the current one. atomic.Value keeps the
// write/read pair safe (let-go futures run on goroutines) rather than relying on
// callers happening to build-once-then-read. Every stored value is a
// *vm.PersistentMap, so atomic.Value's same-concrete-type rule holds.
type hostHashMap struct{ m atomic.Value }

func newHostHashMap() *hostHashMap {
	h := &hostHashMap{}
	h.m.Store(vm.EmptyPersistentMap)
	return h
}

func (h *hostHashMap) Type() vm.ValueType { return hostHashMapType }
func (h *hostHashMap) Unbox() any         { return h }
func (h *hostHashMap) String() string     { return "#<java.util.HashMap>" }

func (h *hostHashMap) InvokeMethod(name vm.Symbol, args []vm.Value) (vm.Value, error) {
	switch string(name) {
	case "putAll":
		if len(args) == 1 {
			// Java putAll(Map) copies the source's entries or throws — it must not
			// silently succeed with an unchanged map. A source we can't enumerate
			// (not Sequable — e.g. another java.util.HashMap) is a hard error, not
			// a no-op, so the data loss is never silent.
			sq, ok := args[0].(vm.Sequable)
			if !ok {
				return vm.NIL, fmt.Errorf("java.util.HashMap.putAll: source %s is not enumerable", args[0].Type().Name())
			}
			acc, ok := h.m.Load().(vm.Associative)
			if !ok {
				acc = vm.EmptyPersistentMap
			}
			for s := sq.Seq(); !vm.SeqIsEmpty(s); s = s.Next() {
				if k, v, ok := vm.MapEntryKV(s.First()); ok {
					acc = acc.Assoc(k, v)
				}
			}
			h.m.Store(acc)
			return vm.NIL, nil
		}
	case "get":
		if len(args) == 1 {
			if l, ok := h.m.Load().(vm.Lookup); ok {
				return l.ValueAt(args[0]), nil
			}
			return vm.NIL, nil
		}
	}
	return vm.NIL, fmt.Errorf("java.util.HashMap has no method .%s/%d", name, len(args))
}

// installHostHashMap registers the (java.util.HashMap. …) constructor forms.
// let-go desugars (HashMap. …) to (->HashMap …); malli uses the bare name
// (imported), other code may use the fully-qualified one.
func installHostHashMap(ns *vm.Namespace) {
	ctor := mustWrap(func(vs []vm.Value) (vm.Value, error) {
		// JVM HashMap(initialCapacity, loadFactor) — sizing args are ignored.
		return newHostHashMap(), nil
	})
	for _, n := range []string{"HashMap.", "->HashMap", "java.util.HashMap.", "->java.util.HashMap"} {
		ns.Def(n, ctor)
	}
}
