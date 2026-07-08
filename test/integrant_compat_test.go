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

// evalCompat compiles and evaluates one expression against the core NS,
// mirroring evalMedley. Used to check the compile-only Clojure-compat stubs
// resolve, so a namespace that references them on its :clj branch loads.
func evalCompat(expr string) (vm.Value, error) {
	ctx := compiler.NewCompiler(vm.NewConsts(), rt.NS(rt.NameCoreNS))
	_, out, err := ctx.CompileMultiple(strings.NewReader(expr))
	if err != nil {
		return vm.NIL, err
	}
	return out, nil
}

// TestClasspathScanStubs exercises the compile-only stubs for JVM-only
// classpath scanning that weavejester/integrant reaches on its :clj branch
// (resources/load-hierarchy/load-annotations). They must resolve so the
// namespace loads; they are never usefully callable under let-go.
func TestClasspathScanStubs(t *testing.T) {
	t.Run("enumeration-seq resolves", func(t *testing.T) {
		_, err := evalCompat(`(defn f [e] (enumeration-seq e))`)
		assert.NoError(t, err)
	})

	t.Run("clojure.lang.RT/baseLoader resolves", func(t *testing.T) {
		_, err := evalCompat(`(defn f [] (clojure.lang.RT/baseLoader))`)
		assert.NoError(t, err)
	})

	// The exact shape integrant's `resources` uses on its :clj branch.
	t.Run("integrant resources shape compiles", func(t *testing.T) {
		_, err := evalCompat(`(defn f [path] (enumeration-seq (.getResources (clojure.lang.RT/baseLoader) path)))`)
		assert.NoError(t, err)
	})
}
