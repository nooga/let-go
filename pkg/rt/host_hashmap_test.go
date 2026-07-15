/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"testing"

	"github.com/nooga/let-go/pkg/vm"
	"github.com/stretchr/testify/assert"
)

// TestHostHashMap exercises the mutable java.util.HashMap compat shim used by
// metosin/malli's fast-registry: (doto (HashMap. n f) (.putAll m)) then (.get fm k).
func TestHostHashMap(t *testing.T) {
	h := newHostHashMap()
	assert.Equal(t, "java.util.HashMap", h.Type().Name())

	m := vm.EmptyPersistentMap.
		Assoc(vm.Keyword("a"), vm.Int(1)).
		Assoc(vm.Keyword("b"), vm.Int(2))
	_, err := h.InvokeMethod(vm.Symbol("putAll"), []vm.Value{m})
	assert.NoError(t, err)

	got, err := h.InvokeMethod(vm.Symbol("get"), []vm.Value{vm.Keyword("a")})
	assert.NoError(t, err)
	assert.Equal(t, vm.Int(1), got)

	absent, err := h.InvokeMethod(vm.Symbol("get"), []vm.Value{vm.Keyword("z")})
	assert.NoError(t, err)
	assert.Equal(t, vm.NIL, absent)

	// A second putAll accumulates (Java semantics) — the first key survives.
	m2 := vm.EmptyPersistentMap.Assoc(vm.Keyword("c"), vm.Int(3))
	_, err = h.InvokeMethod(vm.Symbol("putAll"), []vm.Value{m2})
	assert.NoError(t, err)
	a, _ := h.InvokeMethod(vm.Symbol("get"), []vm.Value{vm.Keyword("a")})
	c, _ := h.InvokeMethod(vm.Symbol("get"), []vm.Value{vm.Keyword("c")})
	assert.Equal(t, vm.Int(1), a)
	assert.Equal(t, vm.Int(3), c)

	// putAll of a non-enumerable source errors loudly (no silent data loss),
	// and leaves the map unchanged.
	_, err = h.InvokeMethod(vm.Symbol("putAll"), []vm.Value{vm.Int(7)})
	assert.Error(t, err)
	stillA, _ := h.InvokeMethod(vm.Symbol("get"), []vm.Value{vm.Keyword("a")})
	assert.Equal(t, vm.Int(1), stillA)
}
