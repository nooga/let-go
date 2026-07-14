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

// evalJVMStubs compiles+evaluates an expression against the core NS (mirrors
// evalMedley; self-contained to avoid a pkg/compiler -> pkg/rt import cycle).
func evalJVMStubs(expr string) (vm.Value, error) {
	ctx := compiler.NewCompiler(vm.NewConsts(), rt.NS(rt.NameCoreNS))
	_, out, err := ctx.CompileMultiple(strings.NewReader(expr))
	if err != nil {
		return vm.NIL, err
	}
	return out, nil
}

func TestJVMStubs(t *testing.T) {
	// Loud ctor stubs must RESOLVE (compile) but throw if run — probe resolution
	// from an unevaluated fn body so the loud error doesn't fire.
	t.Run("timeout stubs resolve, throw if run", func(t *testing.T) {
		v, err := evalJVMStubs(`[(fn? (fn [] (FutureTask. nil))) (= TimeUnit/MILLISECONDS TimeUnit/MILLISECONDS)]`)
		assert.NoError(t, err)
		assert.Equal(t, "[true true]", v.String())
		_, runErr := evalJVMStubs(`(FutureTask. nil)`)
		assert.Error(t, runErr)
	})
	t.Run("java.time chain threads at load, .format throws", func(t *testing.T) {
		v, err := evalJVMStubs(`(some? (-> (DateTimeFormatterBuilder.) (.appendPattern "x") (.toFormatter)))`)
		assert.NoError(t, err)
		assert.Equal(t, vm.TRUE, v)
		_, runErr := evalJVMStubs(`(.format (DateTimeFormatterBuilder.) "x")`)
		assert.Error(t, runErr)
	})
	t.Run("BigDecimal./URI. resolve, throw if run", func(t *testing.T) {
		v, err := evalJVMStubs(`[(fn? (fn [] (BigDecimal. "1"))) (fn? (fn [] (URI. "x")))]`)
		assert.NoError(t, err)
		assert.Equal(t, "[true true]", v.String())
	})
	t.Run("clojure.lang.IDeref/IFn resolve as protocols", func(t *testing.T) {
		v, err := evalJVMStubs(`[(some? clojure.lang.IDeref) (some? clojure.lang.IFn)]`)
		assert.NoError(t, err)
		assert.Equal(t, "[true true]", v.String())
	})
}
