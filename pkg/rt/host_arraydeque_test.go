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

// TestHostArrayDeque exercises the mutable java.util.ArrayDeque compat shim used
// as the CPS regex engine's LIFO backtracking stack in metosin/malli.
func TestHostArrayDeque(t *testing.T) {
	d := newHostArrayDeque()
	assert.Equal(t, "java.util.ArrayDeque", d.Type().Name())

	empty, _ := d.InvokeMethod(vm.Symbol("isEmpty"), nil)
	assert.Equal(t, vm.Boolean(true), empty)

	for i := 1; i <= 3; i++ {
		_, err := d.InvokeMethod(vm.Symbol("push"), []vm.Value{vm.Int(i)})
		assert.NoError(t, err)
	}

	peek, _ := d.InvokeMethod(vm.Symbol("peek"), nil)
	assert.Equal(t, vm.Int(3), peek)

	// LIFO order.
	got := []vm.Value{}
	for i := 0; i < 3; i++ {
		v, err := d.InvokeMethod(vm.Symbol("pop"), nil)
		assert.NoError(t, err)
		got = append(got, v)
	}
	assert.Equal(t, []vm.Value{vm.Int(3), vm.Int(2), vm.Int(1)}, got)

	notEmpty, _ := d.InvokeMethod(vm.Symbol("isEmpty"), nil)
	assert.Equal(t, vm.Boolean(true), notEmpty)

	_, err := d.InvokeMethod(vm.Symbol("pop"), nil)
	assert.Error(t, err) // pop on empty throws
}
