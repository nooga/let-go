/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"fmt"
	"sync"

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
// with (.get fm type). The backing map is an immutable let-go map behind a
// mutex: writes build a new map and swap it in, reads take the current one, so
// a write racing a read from another goroutine (let-go futures) is safe.
type hostHashMap struct {
	mu sync.Mutex
	m  vm.Associative
}

func newHostHashMap() *hostHashMap {
	return &hostHashMap{m: vm.EmptyPersistentMap}
}

func (h *hostHashMap) Type() vm.ValueType { return hostHashMapType }
func (h *hostHashMap) Unbox() any         { return h }
func (h *hostHashMap) String() string     { return "#<java.util.HashMap>" }
func (h *hostHashMap) HostShim()          {}

func (h *hostHashMap) load() vm.Associative {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.m
}

// update swaps in f's result, returning what f returned besides the map.
func (h *hostHashMap) update(f func(vm.Associative) (vm.Associative, vm.Value)) vm.Value {
	h.mu.Lock()
	defer h.mu.Unlock()
	m, ret := f(h.m)
	h.m = m
	return ret
}

// notFound marks an absent key in a lookup, so a key mapped to nil is told
// apart from no key.
var hashMapNotFound = vm.Keyword("let-go.host/not-found")

func (h *hostHashMap) InvokeMethod(name vm.Symbol, args []vm.Value) (vm.Value, error) {
	switch string(name) {
	case "putAll":
		if len(args) == 1 {
			// Java putAll(Map) copies the source's entries or throws — it must not
			// silently succeed with an unchanged map.
			sq, ok := args[0].(vm.Sequable)
			if !ok || args[0] == vm.NIL {
				return vm.NIL, throwHost(vm.ClassClassCastException,
					"class %s cannot be cast to class java.util.Map", args[0].Type().Name())
			}
			h.update(func(acc vm.Associative) (vm.Associative, vm.Value) {
				for s := sq.Seq(); !vm.SeqIsEmpty(s); s = s.Next() {
					if k, v, ok := vm.MapEntryKV(s.First()); ok {
						acc = acc.Assoc(k, v)
					}
				}
				return acc, vm.NIL
			})
			return vm.NIL, nil
		}
	case "get":
		if len(args) == 1 {
			return h.lookup(args[0], vm.NIL), nil
		}
	case "getOrDefault":
		if len(args) == 2 {
			return h.lookup(args[0], args[1]), nil
		}
	case "put":
		if len(args) == 2 {
			// Returns the previous value, or nil when the key was absent.
			return h.update(func(m vm.Associative) (vm.Associative, vm.Value) {
				prev := valueAtOr(m, args[0], vm.NIL)
				return m.Assoc(args[0], args[1]), prev
			}), nil
		}
	case "remove":
		if len(args) == 1 {
			return h.update(func(m vm.Associative) (vm.Associative, vm.Value) {
				prev := valueAtOr(m, args[0], vm.NIL)
				return m.Dissoc(args[0]), prev
			}), nil
		}
	case "containsKey":
		if len(args) == 1 {
			return vm.Boolean(h.lookup(args[0], hashMapNotFound) != hashMapNotFound), nil
		}
	case "size":
		if len(args) == 0 {
			return countOf(h.load()), nil
		}
	case "isEmpty":
		if len(args) == 0 {
			c, ok := h.load().(vm.Counted)
			return vm.Boolean(!ok || c.RawCount() == 0), nil
		}
	case "clear":
		if len(args) == 0 {
			h.update(func(vm.Associative) (vm.Associative, vm.Value) {
				return vm.EmptyPersistentMap, vm.NIL
			})
			return vm.NIL, nil
		}
	}
	return vm.NIL, fmt.Errorf("java.util.HashMap has no method .%s/%d", name, len(args))
}

func (h *hostHashMap) lookup(k, notFound vm.Value) vm.Value {
	return valueAtOr(h.load(), k, notFound)
}

func valueAtOr(m vm.Associative, k, notFound vm.Value) vm.Value {
	if l, ok := m.(vm.Lookup); ok {
		return l.ValueAtOr(k, notFound)
	}
	return notFound
}

func countOf(m vm.Associative) vm.Value {
	if c, ok := m.(vm.Counted); ok {
		return c.Count()
	}
	return vm.Int(0)
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
