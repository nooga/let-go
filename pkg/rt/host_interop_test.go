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

// TestSeqIterator checks the java.util.Iterator shim over a let-go vector.
func TestSeqIterator(t *testing.T) {
	it := &seqIterator{cur: hostNormSeq(vm.ArrayVector{vm.Int(1), vm.Int(2)})}

	h, _ := it.InvokeMethod(vm.Symbol("hasNext"), nil)
	assert.Equal(t, vm.Boolean(true), h)
	v1, _ := it.InvokeMethod(vm.Symbol("next"), nil)
	assert.Equal(t, vm.Int(1), v1)
	v2, _ := it.InvokeMethod(vm.Symbol("next"), nil)
	assert.Equal(t, vm.Int(2), v2)
	hEnd, _ := it.InvokeMethod(vm.Symbol("hasNext"), nil)
	assert.Equal(t, vm.Boolean(false), hEnd)
	_, err := it.InvokeMethod(vm.Symbol("next"), nil)
	assert.Error(t, err)
}
