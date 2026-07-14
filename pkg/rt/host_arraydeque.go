/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"fmt"

	"github.com/nooga/let-go/pkg/vm"
)

// theHostArrayDequeType is the ValueType for the mutable java.util.ArrayDeque
// compat shim — a JVM-interop shim (mutable, not a persistent let-go value), so
// it lives in the compat layer rather than in pkg/vm.
type theHostArrayDequeType struct{}

func (t *theHostArrayDequeType) String() string     { return t.Name() }
func (t *theHostArrayDequeType) Type() vm.ValueType { return vm.TypeType }
func (t *theHostArrayDequeType) Unbox() any         { return nil }
func (t *theHostArrayDequeType) Name() string       { return "java.util.ArrayDeque" }
func (t *theHostArrayDequeType) Box(any) (vm.Value, error) {
	return vm.NIL, fmt.Errorf("java.util.ArrayDeque cannot be boxed")
}

// hostArrayDequeType is the singleton type for hostArrayDeque values.
var hostArrayDequeType = &theHostArrayDequeType{}

// hostArrayDeque is a mutable LIFO stack backing (java.util.ArrayDeque.).
// malli.impl.regex/make-stack uses it as the CPS regex engine's backtracking
// stack via .push / .pop / .peek / .isEmpty (Java Deque used stack-style).
type hostArrayDeque struct{ items []vm.Value }

func newHostArrayDeque() *hostArrayDeque { return &hostArrayDeque{} }

func (d *hostArrayDeque) Type() vm.ValueType { return hostArrayDequeType }
func (d *hostArrayDeque) Unbox() any         { return d }
func (d *hostArrayDeque) String() string     { return "#<java.util.ArrayDeque>" }

func (d *hostArrayDeque) InvokeMethod(name vm.Symbol, args []vm.Value) (vm.Value, error) {
	switch string(name) {
	case "push":
		if len(args) == 1 {
			d.items = append(d.items, args[0])
			return vm.NIL, nil
		}
	case "pop":
		if len(args) == 0 {
			n := len(d.items)
			if n == 0 {
				return vm.NIL, fmt.Errorf("java.util.ArrayDeque: pop on empty deque")
			}
			v := d.items[n-1]
			d.items = d.items[:n-1]
			return v, nil
		}
	case "peek":
		if len(args) == 0 {
			if len(d.items) == 0 {
				return vm.NIL, nil
			}
			return d.items[len(d.items)-1], nil
		}
	case "isEmpty":
		if len(args) == 0 {
			return vm.Boolean(len(d.items) == 0), nil
		}
	}
	return vm.NIL, fmt.Errorf("java.util.ArrayDeque has no method .%s/%d", name, len(args))
}

// installHostArrayDeque registers the (java.util.ArrayDeque.) constructor forms.
func installHostArrayDeque(ns *vm.Namespace) {
	ctor := mustWrap(func(vs []vm.Value) (vm.Value, error) {
		return newHostArrayDeque(), nil
	})
	for _, n := range []string{"ArrayDeque.", "->ArrayDeque", "java.util.ArrayDeque.", "->java.util.ArrayDeque"} {
		ns.Def(n, ctor)
	}
}
