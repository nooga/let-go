/*
 * Copyright (c) 2026 let-go contributors
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

func evalStringInterop(expr string) (vm.Value, error) {
	ctx := compiler.NewCompiler(vm.NewConsts(), rt.NS(rt.NameCoreNS))
	_, out, err := ctx.CompileMultiple(strings.NewReader(expr))
	if err != nil {
		return vm.NIL, err
	}
	return out, nil
}

// TestStringHostMethods exercises the java.lang.String surface Clojure
// libraries reach on their :clj branches (honeysql walks .length/.charAt/
// .indexOf over entities and upper-cases keywords with a Locale argument).
// Indices are rune indices, matching Java's UTF-16 units for BMP text.
func TestStringHostMethods(t *testing.T) {
	cases := []struct {
		expr string
		want vm.Value
	}{
		{`(.toString "x")`, vm.String("x")},
		{`(.length "héllo")`, vm.Int(5)},
		{`(.isEmpty "")`, vm.Boolean(true)},
		{`(.isEmpty "a")`, vm.Boolean(false)},
		{`(.charAt "héllo" 1)`, vm.Char('é')},
		{`(.indexOf "hello" "ll")`, vm.Int(2)},
		{`(.indexOf "hello" "zz")`, vm.Int(-1)},
		{`(.indexOf "abc" 98)`, vm.Int(1)}, // Java's indexOf(int ch) codepoint form
		{`(.indexOf "aba" \a 1)`, vm.Int(2)},
		{`(.concat "a" "b")`, vm.String("ab")},
		{`(.substring "hello" 1)`, vm.String("ello")},
		{`(.substring "hello" 1 3)`, vm.String("el")},
		{`(.startsWith "hello" "he")`, vm.Boolean(true)},
		{`(.endsWith "hello" "lo")`, vm.Boolean(true)},
		{`(.contains "hello" "ell")`, vm.Boolean(true)},
		{`(.toUpperCase "ab")`, vm.String("AB")},
		{`(.toUpperCase "ab" java.util.Locale/US)`, vm.String("AB")},
		{`(.toLowerCase "AB" Locale/US)`, vm.String("ab")},
		{`(.trim "  a ")`, vm.String("a")},
	}
	for _, tc := range cases {
		t.Run(tc.expr, func(t *testing.T) {
			v, err := evalStringInterop(tc.expr)
			assert.NoError(t, err)
			assert.Equal(t, tc.want, v)
		})
	}

	t.Run("charAt out of bounds fails loudly", func(t *testing.T) {
		_, err := evalStringInterop(`(.charAt "ab" 5)`)
		assert.Error(t, err)
	})

	// Locale/getDefault-setDefault-forLanguageTag are honest no-ops: a
	// default-locale change cannot affect Go's case mapping, which is what
	// locale-independence tests (honeysql's Turkish-locale regression)
	// assert. The round-trip must resolve and casing must stay identical.
	t.Run("default-locale round-trip cannot change casing", func(t *testing.T) {
		v, err := evalStringInterop(`
			(let [orig (java.util.Locale/getDefault)]
			  (try
			    (java.util.Locale/setDefault (java.util.Locale/forLanguageTag "tr"))
			    (.toUpperCase "inner join")
			    (finally (java.util.Locale/setDefault orig))))`)
		assert.NoError(t, err)
		assert.Equal(t, vm.String("INNER JOIN"), v)
	})
}

// TestKeywordHostMethods pins clojure.lang.Keyword's accessors. honeysql's
// :clj branch converts keywords to symbols via (.sym k).
func TestKeywordHostMethods(t *testing.T) {
	cases := []struct {
		expr string
		want vm.Value
	}{
		{`(.sym :a/b)`, vm.Symbol("a/b")},
		{`(.sym :plain)`, vm.Symbol("plain")},
		{`(.getName :a/b)`, vm.String("b")},
		{`(.getNamespace :a/b)`, vm.String("a")},
		{`(.getNamespace :plain)`, vm.NIL},
		{`(.toString :a/b)`, vm.String(":a/b")},
	}
	for _, tc := range cases {
		t.Run(tc.expr, func(t *testing.T) {
			v, err := evalStringInterop(tc.expr)
			assert.NoError(t, err)
			assert.Equal(t, tc.want, v)
		})
	}

	// A ^Tag hint in expression position compiles to a runtime with-meta
	// wrap, and keywords (being invokable) wrap into MetaFn — the decorator
	// must keep delegating method dispatch to the wrapped value.
	t.Run("hinted receiver keeps its methods", func(t *testing.T) {
		v, err := evalStringInterop(`(let [k :a/b] (.sym ^clojure.lang.Keyword k))`)
		assert.NoError(t, err)
		assert.Equal(t, vm.Symbol("a/b"), v)
	})
}

// TestStringBuilderShim exercises the (StringBuilder.) ctor forms + methods
// end-to-end through the compiler.
func TestStringBuilderShim(t *testing.T) {
	t.Run("append chain and toString", func(t *testing.T) {
		v, err := evalStringInterop(
			`(let [sb (StringBuilder.)] (.append sb "x") (.append sb 42) (.append sb \!) (.toString sb))`)
		assert.NoError(t, err)
		assert.Equal(t, vm.String("x42!"), v)
	})
	t.Run("append returns the builder for chaining", func(t *testing.T) {
		v, err := evalStringInterop(`(.toString (.append (.append (StringBuilder.) "a") "b"))`)
		assert.NoError(t, err)
		assert.Equal(t, vm.String("ab"), v)
	})
	t.Run("initial content and length", func(t *testing.T) {
		v, err := evalStringInterop(`(.length (java.lang.StringBuilder. "héllo"))`)
		assert.NoError(t, err)
		assert.Equal(t, vm.Int(5), v)
	})
	t.Run("capacity hint is accepted", func(t *testing.T) {
		v, err := evalStringInterop(`(.toString (StringBuilder. 16))`)
		assert.NoError(t, err)
		assert.Equal(t, vm.String(""), v)
	})
	t.Run("nil appends Java's null", func(t *testing.T) {
		v, err := evalStringInterop(`(let [sb (StringBuilder.)] (.append sb nil) (.toString sb))`)
		assert.NoError(t, err)
		assert.Equal(t, vm.String("null"), v)
	})
}

// TestUniversalToString pins the .toString fallback: every JVM object answers
// toString, so any let-go value maps it to str semantics — except nil, which
// fails loudly like the NPE it would be on the JVM.
func TestUniversalToString(t *testing.T) {
	cases := []struct {
		expr string
		want vm.Value
	}{
		{`(.toString 42)`, vm.String("42")},
		{`(.toString 1.5)`, vm.String("1.5")},
		{`(.toString true)`, vm.String("true")},
		{`(.toString (quote s))`, vm.String("s")},
	}
	for _, tc := range cases {
		t.Run(tc.expr, func(t *testing.T) {
			v, err := evalStringInterop(tc.expr)
			assert.NoError(t, err)
			assert.Equal(t, tc.want, v)
		})
	}
	t.Run("toString on nil fails loudly", func(t *testing.T) {
		_, err := evalStringInterop(`(.toString nil)`)
		assert.Error(t, err)
	})
}

// TestUncheckedCoercionsAcceptChars pins Java's char->int widening, which
// Clojure's unchecked-* inherit: honeysql's alphanumeric? state machine feeds
// .charAt results straight into (unchecked-long (unchecked-int c)).
func TestUncheckedCoercionsAcceptChars(t *testing.T) {
	cases := []struct {
		expr string
		want vm.Value
	}{
		{`(unchecked-int \a)`, vm.Int(97)},
		{`(unchecked-long \a)`, vm.Int(97)},
		{`(unchecked-short \a)`, vm.Int(97)},
		{`(unchecked-byte \a)`, vm.Int(97)},
	}
	for _, tc := range cases {
		t.Run(tc.expr, func(t *testing.T) {
			v, err := evalStringInterop(tc.expr)
			assert.NoError(t, err)
			assert.Equal(t, tc.want, v)
		})
	}
}
