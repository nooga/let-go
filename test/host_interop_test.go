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

// evalInterop compiles+evaluates an expression against the core NS (mirrors
// evalMedley; self-contained to avoid a pkg/compiler -> pkg/rt import cycle).
func evalInterop(expr string) (vm.Value, error) {
	ctx := compiler.NewCompiler(vm.NewConsts(), rt.NS(rt.NameCoreNS))
	_, out, err := ctx.CompileMultiple(strings.NewReader(expr))
	if err != nil {
		return vm.NIL, err
	}
	return out, nil
}

// TestCollectionInterop exercises the generic JVM collection-interface `.`-method
// dispatch on let-go collections.
func TestCollectionInterop(t *testing.T) {
	t.Run(".valAt", func(t *testing.T) {
		v, err := evalInterop(`(.valAt {:a 1} :a)`)
		assert.NoError(t, err)
		assert.Equal(t, vm.Int(1), v)
	})
	t.Run(".valAt not-found", func(t *testing.T) {
		v, err := evalInterop(`(.valAt {:a 1} :b :nf)`)
		assert.NoError(t, err)
		assert.Equal(t, vm.Keyword("nf"), v)
	})
	t.Run(".nth / .count", func(t *testing.T) {
		v, err := evalInterop(`[(.nth [10 20] 1) (.count [1 2 3])]`)
		assert.NoError(t, err)
		assert.Equal(t, "[20 3]", v.String())
	})
	t.Run(".assoc", func(t *testing.T) {
		v, err := evalInterop(`(.valAt (.assoc {} :a 1) :a)`)
		assert.NoError(t, err)
		assert.Equal(t, vm.Int(1), v)
	})
	t.Run(".iterator over a vector", func(t *testing.T) {
		v, err := evalInterop(`(let [it (.iterator [1 2])]
                                [(.hasNext it) (.next it) (.next it) (.hasNext it)])`)
		assert.NoError(t, err)
		assert.Equal(t, "[true 1 2 false]", v.String())
	})
	t.Run(".iterator over an empty lazy seq", func(t *testing.T) {
		v, err := evalInterop(`(.hasNext (.iterator (lazy-seq nil)))`)
		assert.NoError(t, err)
		assert.Equal(t, vm.FALSE, v)
	})
	// Strings are not real collections here: the shim must not shadow core's
	// rune-based count/nth with String's byte-based Counted/Indexed.
	t.Run("strings keep core rune semantics", func(t *testing.T) {
		v, err := evalInterop(`[(.count "café") (.count "é") (= (.nth "café" 3) \é)]`)
		assert.NoError(t, err)
		assert.Equal(t, "[4 1 true]", v.String())
	})
}
