/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package test

import (
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
	"github.com/stretchr/testify/assert"
)

// evalMutableColls compiles+evaluates an expression against the core NS (mirrors
// evalMedley; self-contained to avoid a pkg/compiler -> pkg/rt import cycle).
func evalMutableColls(expr string) (vm.Value, error) {
	ctx := compiler.NewCompiler(vm.NewConsts(), rt.NS(rt.NameCoreNS))
	_, out, err := ctx.CompileMultiple(strings.NewReader(expr))
	if err != nil {
		return vm.NIL, err
	}
	return out, nil
}

// TestMutableCollections exercises the ctor forms + methods of the mutable
// java.util.HashMap / ArrayDeque shims end-to-end through the compiler. (Their
// method semantics are unit-tested in pkg/rt/host_*_test.go.)
func TestMutableCollections(t *testing.T) {
	t.Run("HashMap putAll/get", func(t *testing.T) {
		v, err := evalMutableColls(`(let [m (java.util.HashMap. 4 0.5)] (.putAll m {:a 1}) (.get m :a))`)
		assert.NoError(t, err)
		assert.Equal(t, vm.Int(1), v)
	})
	t.Run("ArrayDeque LIFO", func(t *testing.T) {
		v, err := evalMutableColls(`(let [d (java.util.ArrayDeque.)]
                                     (.push d 1) (.push d 2)
                                     [(.peek d) (.pop d) (.pop d) (.isEmpty d)])`)
		assert.NoError(t, err)
		assert.Equal(t, "[2 2 1 true]", v.String())
	})
}
