/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package compiler

import (
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
	"github.com/stretchr/testify/assert"
)

func TestContext_Compile(t *testing.T) {
	tests := map[string]interface{}{
		"(+ (* 2 20) 2)":          42,
		"(- 10 2)":                8,
		`(if true "big" "meh")`:   "big",
		`(if false "big" "meh")`:  "meh",
		`(if nil 1 2)`:            2,
		`(if true 101)`:           101,
		`(if false 101)`:          nil,
		`(do 1 2 3)`:              3,
		`(do (+ 1 2))`:            3,
		`(do)`:                    nil,
		`(do (def x 40) (+ x 2))`: 42,
		`(do (def x true)
			 (def y (if x :big :meh))
		   y)`: "big",
		`(if \P \N \P)`:                         'N',
		`(println "hello" "world")`:             nil,
		`(do (def sq (fn [x] (* x x))) (sq 9))`: 81,
		`[1 2 (+ 1 2)]`:                         []vm.Value{vm.Int(1), vm.Int(2), vm.Int(3)},
		`'foo`:                                  "foo",
		`(quote foo)`:                           "foo",
		`{}`:                                    map[vm.Value]vm.Value{},
		`{:a 1}`:                                map[vm.Value]vm.Value{vm.Keyword("a"): vm.Int(1)},
	}
	for k, v := range tests {
		out, err := Eval(k)
		assert.NoError(t, err)
		assert.Equal(t, v, out.Unbox())
	}
}

func TestContext_CompileFn(t *testing.T) {
	out, err := Eval("(fn [x] (+ x 1))")
	assert.NoError(t, err)

	var inc func(int) int
	out.Unbox().(func(interface{}))(&inc)

	assert.NotNil(t, inc)

	x := 999
	assert.Equal(t, x, inc(x)-1)
}

func TestContext_CompileFnPoly(t *testing.T) {
	out, err := Eval("(fn [x] x)")
	assert.NoError(t, err)

	identity := out.Unbox().(func(interface{}))

	var intIdentity func(int) int
	identity(&intIdentity)

	var strIdentity func(string) string
	identity(&strIdentity)

	x := 999
	y := "foobar"
	assert.Equal(t, x, intIdentity(x))
	assert.Equal(t, y, strIdentity(y))
}

func TestContext_CompileMultiple(t *testing.T) {
	src := `(def parens 20)
			(def fun 1) 
			(def double (fn [a] (* a 2)))
			(def hey! (fn [a _ b] (+ a b)))
			(println (hey! (double parens) 'equals (double fun)))`

	cp := vm.NewConsts()
	ns := rt.NS(rt.NameCoreNS)
	assert.NotNil(t, ns)
	ctx := NewCompiler(cp, ns)

	chunk, _, err := ctx.CompileMultiple(strings.NewReader(src))
	assert.NoError(t, err)

	_, err = vm.NewFrame(chunk, nil).Run()
	assert.NoError(t, err)
}

func TestContext_CompileVar(t *testing.T) {
	v := vm.NewVar(rt.NS(rt.NameCoreNS), rt.NameCoreNS, "foo")

	out, err := Eval("(var foo)")
	assert.NoError(t, err)
	assert.Equal(t, v, out)

	out, err = Eval("#'foo")
	assert.NoError(t, err)
	assert.Equal(t, v, out)
}

func TestContext_CompileMultiArityFn(t *testing.T) {
	src := `(def f (fn* ([a] (+ a 1)) 
						([a b] (+ a b)) 
						([a b & r] (+ a b (second r)))))
			(and (= 2 (f 1))
			     (= 3 (f 1 2))
				 (= 6 (f 1 2 4 3)))`

	cp := vm.NewConsts()
	ns := rt.NS(rt.NameCoreNS)
	assert.NotNil(t, ns)
	ctx := NewCompiler(cp, ns)

	chunk, _, err := ctx.CompileMultiple(strings.NewReader(src))
	assert.NoError(t, err)

	val, err := vm.NewFrame(chunk, nil).Run()
	assert.NoError(t, err)
	assert.Equal(t, vm.TRUE, val)
}
