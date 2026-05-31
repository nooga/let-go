/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
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

// evalMedley compiles and evaluates a single expression against the core NS,
// returning the resulting value and any compile/eval error. It mirrors the
// helper used by language_test.go but lives here so the medley-compat suite is
// self-contained. Placed in package test (not package rt) to avoid the
// pkg/compiler -> pkg/rt import cycle.
func evalMedley(expr string) (vm.Value, error) {
	ctx := compiler.NewCompiler(vm.NewConsts(), rt.NS(rt.NameCoreNS))
	_, out, err := ctx.CompileMultiple(strings.NewReader(expr))
	if err != nil {
		return vm.NIL, err
	}
	return out, nil
}

// TestMedleyCompat exercises the Clojure-compat aliases added so that
// weavejester/medley loads under let-go. Each case is a JVM reference medley
// reaches on its :clj / :default reader-conditional branches; without the
// alias the namespace fails to compile with "Can't resolve ...".
func TestMedleyCompat(t *testing.T) {
	// Blocker #1: clojure.lang.PersistentQueue marker + EMPTY stub.
	// The marker resolves so (instance? clojure.lang.PersistentQueue x)
	// compiles; nothing carries the marker as an ancestor, so it is false
	// (load-only semantics — queue?/queue are degraded).
	t.Run("PersistentQueue marker instance? is false", func(t *testing.T) {
		v, err := evalMedley(`(instance? clojure.lang.PersistentQueue [1 2])`)
		assert.NoError(t, err)
		assert.Equal(t, vm.FALSE, v)
	})

	t.Run("PersistentQueue/EMPTY resolves without compile error", func(t *testing.T) {
		_, err := evalMedley(`clojure.lang.PersistentQueue/EMPTY`)
		assert.NoError(t, err)
	})

	// EMPTY is a load-only stub bound to a non-collection marker symbol, so
	// medley's (queue coll) = (into (queue) coll) must FAIL LOUDLY rather than
	// silently return a reversed list. Conjing onto it errors at runtime.
	t.Run("PersistentQueue/EMPTY fails loudly when conj'd", func(t *testing.T) {
		_, err := evalMedley(`(into clojure.lang.PersistentQueue/EMPTY [1 2 3])`)
		assert.Error(t, err)
	})

	// Blocker #2: (java.util.ArrayList.) / (java.util.ArrayList. n). medley's
	// partition-between / sliding construct one on their :clj branch. let-go
	// has no mutable ArrayList, so this is a load-only ctor stub: the defn must
	// COMPILE (the constructor symbol must resolve); .add/.toArray throw at
	// runtime if ever called (out of scope).
	t.Run("ArrayList zero-arg ctor compiles", func(t *testing.T) {
		_, err := evalMedley(`(defn f [] (java.util.ArrayList.))`)
		assert.NoError(t, err)
	})

	t.Run("ArrayList one-arg ctor compiles", func(t *testing.T) {
		_, err := evalMedley(`(defn g [n] (java.util.ArrayList. n))`)
		assert.NoError(t, err)
	})
}
