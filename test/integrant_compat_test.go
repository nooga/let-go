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

// evalIntegrant compiles and evaluates one expression against the core NS,
// mirroring evalMedley — used to check that find-var / get-method resolve.
// Runtime behavior is covered by test/integrant_compat_test.lg.
func evalIntegrant(expr string) (vm.Value, error) {
	ctx := compiler.NewCompiler(vm.NewConsts(), rt.NS(rt.NameCoreNS))
	_, out, err := ctx.CompileMultiple(strings.NewReader(expr))
	if err != nil {
		return vm.NIL, err
	}
	return out, nil
}

// TestFindVarGetMethodCompat checks the var/multimethod introspection fns
// resolve. weavejester/integrant references find-var (default init-key) and
// get-method (can-expand-key?); an unresolved symbol there fails the whole
// namespace compile.
func TestFindVarGetMethodCompat(t *testing.T) {
	t.Run("find-var resolves", func(t *testing.T) {
		_, err := evalIntegrant(`(defn f [s] (find-var s))`)
		assert.NoError(t, err)
	})

	t.Run("get-method resolves", func(t *testing.T) {
		_, err := evalIntegrant(`(defn f [m v] (get-method m v))`)
		assert.NoError(t, err)
	})
}
