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
	// clojure.lang.PersistentQueue marker + EMPTY stub.
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

	// (java.util.ArrayList.) / (java.util.ArrayList. n). medley's
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

	// Throwable bare in (instance? Throwable ex). The marker is false
	// for non-exceptions (matching medley's "nil for all other types" fallback)
	// and true for let-go ex-info values, which ARE-A Throwable (see hierarchy).
	t.Run("Throwable marker instance? is false for non-exceptions", func(t *testing.T) {
		v, err := evalMedley(`(instance? Throwable nil)`)
		assert.NoError(t, err)
		assert.Equal(t, vm.FALSE, v)
	})

	t.Run("Throwable marker instance? is true for ex-info", func(t *testing.T) {
		v, err := evalMedley(`(instance? Throwable (ex-info "boom" {}))`)
		assert.NoError(t, err)
		assert.Equal(t, vm.TRUE, v)
	})

	// .getMessage / .getCause dispatch on ex-info via ExInfo.InvokeMethod. Use
	// the TYPE-HINTED form medley actually emits — (.getMessage ^Throwable ex) —
	// which compiles to with-meta on the value, so this also exercises ExInfo's
	// IMeta support. An unhinted call would miss the real medley path.
	t.Run("getMessage dispatches on type-hinted ex-info", func(t *testing.T) {
		v, err := evalMedley(`(.getMessage ^Throwable (ex-info "boom" {}))`)
		assert.NoError(t, err)
		assert.Equal(t, vm.String("boom"), v)
	})

	t.Run("getCause returns nil when no cause (type-hinted)", func(t *testing.T) {
		v, err := evalMedley(`(.getCause ^Throwable (ex-info "boom" {}))`)
		assert.NoError(t, err)
		assert.Equal(t, vm.NIL, v)
	})

	// instance? Throwable must remain true after a type hint adds metadata.
	t.Run("Throwable marker instance? true for hinted ex-info", func(t *testing.T) {
		v, err := evalMedley(`(instance? Throwable ^Throwable (ex-info "boom" {}))`)
		assert.NoError(t, err)
		assert.Equal(t, vm.TRUE, v)
	})

	// java.util.UUID/fromString reuses vm.ParseUUID → m/uuid works.
	t.Run("UUID/fromString returns a UUID", func(t *testing.T) {
		v, err := evalMedley(`(java.util.UUID/fromString "00000000-0000-0000-0000-000000000000")`)
		assert.NoError(t, err)
		assert.Equal(t, vm.UUIDType, v.Type())
	})

	// java.util.UUID/randomUUID reuses the random-uuid builtin → m/random-uuid works.
	t.Run("UUID/randomUUID returns a UUID", func(t *testing.T) {
		v, err := evalMedley(`(java.util.UUID/randomUUID)`)
		assert.NoError(t, err)
		assert.Equal(t, vm.UUIDType, v.Type())
	})

	// java.util.regex.Pattern aliased to vm.RegexType → m/regexp?
	// works against real let-go regexes.
	t.Run("Pattern instance? is true for a regex", func(t *testing.T) {
		v, err := evalMedley(`(instance? java.util.regex.Pattern #"x")`)
		assert.NoError(t, err)
		assert.Equal(t, vm.TRUE, v)
	})

	// medley's deref-swap! uses compare-and-set!, which let-go lacked.
	// It's a real atom primitive, not a stub. Success path swaps and returns true;
	// mismatch returns false.
	t.Run("compare-and-set! succeeds on matching value", func(t *testing.T) {
		v, err := evalMedley(`(let [a (atom 1)] [(compare-and-set! a 1 2) @a])`)
		assert.NoError(t, err)
		assert.Equal(t, "[true 2]", v.String())
	})

	t.Run("compare-and-set! fails on mismatched value", func(t *testing.T) {
		v, err := evalMedley(`(let [a (atom 1)] [(compare-and-set! a 99 2) @a])`)
		assert.NoError(t, err)
		assert.Equal(t, "[false 1]", v.String())
	})
}
