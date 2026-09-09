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

// evalJVMStatics compiles+evaluates an expression against the core NS (mirrors
// evalMedley; self-contained to avoid a pkg/compiler -> pkg/rt import cycle).
func evalJVMStatics(expr string) (vm.Value, error) {
	ctx := compiler.NewCompiler(vm.NewConsts(), rt.NS(rt.NameCoreNS))
	_, out, err := ctx.CompileMultiple(strings.NewReader(expr))
	if err != nil {
		return vm.NIL, err
	}
	return out, nil
}

func TestJVMStatics(t *testing.T) {
	t.Run("host-class markers", func(t *testing.T) {
		v, err := evalJVMStatics(`[(instance? java.util.Map {}) (instance? CharSequence "s") (instance? Pattern #"x")]`)
		assert.NoError(t, err)
		assert.Equal(t, "[true true true]", v.String())
	})
	t.Run("createOwning -> vector", func(t *testing.T) {
		v, err := evalJVMStatics(`(LazilyPersistentVector/createOwning (object-array 0))`)
		assert.NoError(t, err)
		assert.Equal(t, "[]", v.String())
	})
	t.Run("createWithCheck -> map, dup + odd throw", func(t *testing.T) {
		v, err := evalJVMStatics(`(let [a (object-array 4)]
                                   (aset a 0 :a) (aset a 1 1) (aset a 2 :b) (aset a 3 2)
                                   (clojure.lang.PersistentArrayMap/createWithCheck a))`)
		assert.NoError(t, err)
		assert.Equal(t, vm.Int(1), v.(vm.Lookup).ValueAt(vm.Keyword("a")))
		_, dupErr := evalJVMStatics(`(let [a (object-array 4)]
                                      (aset a 0 :a) (aset a 1 1) (aset a 2 :a) (aset a 3 2)
                                      (clojure.lang.PersistentArrayMap/createWithCheck a))`)
		assert.Error(t, dupErr)
		_, oddErr := evalJVMStatics(`(clojure.lang.PersistentArrayMap/createWithCheck (object-array 3))`)
		assert.Error(t, oddErr)
	})
	t.Run("Array/newInstance + negative throws", func(t *testing.T) {
		v, err := evalJVMStatics(`(count (Array/newInstance java.lang.Object 3))`)
		assert.NoError(t, err)
		assert.Equal(t, vm.Int(3), v)
		_, negErr := evalJVMStatics(`(Array/newInstance java.lang.Object -1)`)
		assert.Error(t, negErr)
	})
	t.Run("Util/hashCombine + Murmur3/hashLong -> int", func(t *testing.T) {
		v, err := evalJVMStatics(`[(int? (Util/hashCombine 1 2)) (int? (Murmur3/hashLong 7))]`)
		assert.NoError(t, err)
		assert.Equal(t, "[true true]", v.String())
	})
	t.Run("number parses + UUID/fromString", func(t *testing.T) {
		v, err := evalJVMStatics(`[(Long/parseLong "42") (Float/parseFloat "3.5") (Double/parseDouble "3.5")
                                   (uuid? (UUID/fromString "00000000-0000-0000-0000-000000000000"))]`)
		assert.NoError(t, err)
		assert.Equal(t, "[42 3.5 3.5 true]", v.String())
	})
	t.Run("System/arraycopy over object arrays", func(t *testing.T) {
		v, err := evalJVMStatics(`(let [src (object-array 3) dst (object-array 3)]
                                   (aset src 0 :x) (aset src 1 :y) (aset src 2 :z)
                                   (System/arraycopy src 0 dst 0 3)
                                   (vec dst))`)
		assert.NoError(t, err)
		assert.Equal(t, "[:x :y :z]", v.String())
		_, oobErr := evalJVMStatics(`(System/arraycopy (object-array 2) 0 (object-array 2) 0 5)`)
		assert.Error(t, oobErr)
	})
}

// TestURLEncoderEncode pins java.net.URLEncoder/encode: Go's url.QueryEscape
// implements the same application/x-www-form-urlencoded rules (space -> '+').
func TestURLEncoderEncode(t *testing.T) {
	v, err := evalJVMStatics(`[(URLEncoder/encode "a b") (java.net.URLEncoder/encode "привіт" "UTF-8")]`)
	assert.NoError(t, err)
	assert.Equal(t, `["a+b" "%D0%BF%D1%80%D0%B8%D0%B2%D1%96%D1%82"]`, v.String())

	_, err = evalJVMStatics(`(URLEncoder/encode "x" "latin-1")`)
	assert.Error(t, err, "non-UTF-8 charsets must fail loudly")
}

// TestJVMNumericStatics pins the java.lang.Math / Long / Integer / Double /
// Character / Byte statics that Clojure libraries reach on their :clj branches.
// Go's math package is not a drop-in for java.lang.Math: round, abs and
// getExponent all differ in return type or rounding direction, so those are
// asserted against JVM behaviour rather than Go's defaults.
func TestJVMNumericStatics(t *testing.T) {
	t.Run("Math/round is half-up and returns a long", func(t *testing.T) {
		// JVM rounds half toward positive infinity: -2.5 -> -2, not -3.
		v, err := evalJVMStatics(`[(Math/round 2.5) (Math/round -2.5) (Math/round 2.4) (Math/round -2.6)]`)
		assert.NoError(t, err)
		assert.Equal(t, "[3 -2 2 -3]", v.String())
	})
	t.Run("Math/abs preserves long vs double", func(t *testing.T) {
		v, err := evalJVMStatics(`[(Math/abs -5) (Math/abs 5) (Math/abs -2.5)]`)
		assert.NoError(t, err)
		assert.Equal(t, "[5 5 2.5]", v.String())
	})
	t.Run("Math/floor and Math/ceil on negatives", func(t *testing.T) {
		v, err := evalJVMStatics(`[(Math/floor -2.5) (Math/ceil -2.5) (Math/floor 2.5) (Math/ceil 2.5)]`)
		assert.NoError(t, err)
		assert.Equal(t, "[-3.0 -2.0 2.0 3.0]", v.String())
	})
	t.Run("Math/getExponent incl. zero and NaN", func(t *testing.T) {
		v, err := evalJVMStatics(`[(Math/getExponent 8.0) (Math/getExponent 1.0) (Math/getExponent 0.0) (Math/getExponent Double/NaN)]`)
		assert.NoError(t, err)
		assert.Equal(t, "[3 0 -1023 1024]", v.String())
	})
	t.Run("Math/pow log exp sqrt scalb", func(t *testing.T) {
		v, err := evalJVMStatics(`[(Math/pow 2 10) (Math/exp 0) (Math/sqrt 9) (Math/scalb 1.0 3) (Math/log 1)]`)
		assert.NoError(t, err)
		assert.Equal(t, "[1024.0 1.0 3.0 8.0 0.0]", v.String())
	})
	t.Run("Math/log of a negative is NaN, not an error", func(t *testing.T) {
		v, err := evalJVMStatics(`(Double/isNaN (Math/log -1))`)
		assert.NoError(t, err)
		assert.Equal(t, vm.TRUE, v)
	})
	t.Run("Long/bitCount and Long/reverse", func(t *testing.T) {
		v, err := evalJVMStatics(`[(Long/bitCount -1) (Long/bitCount 0) (Long/bitCount 7) (Long/reverse (Long/reverse 12345))]`)
		assert.NoError(t, err)
		assert.Equal(t, "[64 0 3 12345]", v.String())
	})
	t.Run("Integer and Byte bounds", func(t *testing.T) {
		v, err := evalJVMStatics(`[Integer/MAX_VALUE Integer/MIN_VALUE Byte/MAX_VALUE Byte/MIN_VALUE]`)
		assert.NoError(t, err)
		assert.Equal(t, "[2147483647 -2147483648 127 -128]", v.String())
	})
	t.Run("Double infinities and NaN", func(t *testing.T) {
		v, err := evalJVMStatics(`[(> Double/POSITIVE_INFINITY 1e308) (< Double/NEGATIVE_INFINITY -1e308) (= Double/NaN Double/NaN)]`)
		assert.NoError(t, err)
		assert.Equal(t, "[true true false]", v.String())
	})
	t.Run("Character/isDigit", func(t *testing.T) {
		v, err := evalJVMStatics(`[(Character/isDigit \5) (Character/isDigit \x)]`)
		assert.NoError(t, err)
		assert.Equal(t, "[true false]", v.String())
	})
}
