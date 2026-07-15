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

// evalCoreCompat compiles+evaluates an expression against the core NS (mirrors
// evalMedley; kept self-contained in package test to avoid a pkg/compiler ->
// pkg/rt import cycle).
func evalCoreCompat(expr string) (vm.Value, error) {
	ctx := compiler.NewCompiler(vm.NewConsts(), rt.NS(rt.NameCoreNS))
	_, out, err := ctx.CompileMultiple(strings.NewReader(expr))
	if err != nil {
		return vm.NIL, err
	}
	return out, nil
}

func TestCoreCompatFns(t *testing.T) {
	t.Run("indexed?", func(t *testing.T) {
		v, err := evalCoreCompat(`[(indexed? [1]) (indexed? {}) (indexed? "s")]`)
		assert.NoError(t, err)
		assert.Equal(t, "[true false false]", v.String())
	})
	t.Run("class aliases type", func(t *testing.T) {
		v, err := evalCoreCompat(`(= (class 1) (type 1))`)
		assert.NoError(t, err)
		assert.Equal(t, vm.TRUE, v)
	})
	t.Run("uri? is always false", func(t *testing.T) {
		v, err := evalCoreCompat(`(uri? "x")`)
		assert.NoError(t, err)
		assert.Equal(t, vm.FALSE, v)
	})
	t.Run("monitor-enter/exit no-op", func(t *testing.T) {
		v, err := evalCoreCompat(`[(monitor-enter :x) (monitor-exit :x)]`)
		assert.NoError(t, err)
		assert.Equal(t, "[nil nil]", v.String())
	})
}
