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
