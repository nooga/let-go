/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package compiler

import (
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/vm"
	"github.com/stretchr/testify/assert"
)

func TestReaderBasic(t *testing.T) {
	cases := map[string]vm.Value{
		"1":                          vm.Int(1),
		"+1":                         vm.Int(1),
		"-1":                         vm.Int(-1),
		"987654321":                  vm.Int(987654321),
		"+987654321":                 vm.Int(987654321),
		"-987654321":                 vm.Int(-987654321),
		"true":                       vm.TRUE,
		"false":                      vm.FALSE,
		"nil":                        vm.NIL,
		"foo":                        vm.Symbol("foo"),
		"()":                         vm.EmptyList,
		"(    )":                     vm.EmptyList,
		"(1 2)":                      vm.EmptyList.Cons(vm.Int(2)).Cons(vm.Int(1)),
		"\"hello\"":                  vm.String("hello"),
		"\"h\\\"el\\tl\\\\o\"":       vm.String("h\"el\tl\\o"),
		":foo":                       vm.Keyword("foo"),
		"\\F":                        vm.Char('F'),
		"\\newline":                  vm.Char('\n'),
		"\\u1234":                    vm.Char('\u1234'),
		"\\o300":                     vm.Char(rune(0300)),
		"\\u03A9":                    vm.Char('Ω'),
		"[]":                         vm.ArrayVector{},
		"[1 :foo true]":              vm.ArrayVector{vm.Int(1), vm.Keyword("foo"), vm.TRUE},
		"'foo":                       vm.EmptyList.Cons(vm.Symbol("foo")).Cons(vm.Symbol("quote")),
		"^:foo zoo":                  vm.EmptyList.Cons(vm.Map{vm.Keyword("foo"): vm.TRUE}).Cons(vm.Symbol("zoo")).Cons(vm.Symbol("with-meta")),
		"^:foo ^:bar zoo":            vm.EmptyList.Cons(vm.Map{vm.Keyword("foo"): vm.TRUE, vm.Keyword("bar"): vm.TRUE}).Cons(vm.Symbol("zoo")).Cons(vm.Symbol("with-meta")),
		"^{:foo 1 :baz 2} ^:bar zoo": vm.EmptyList.Cons(vm.Map{vm.Keyword("foo"): vm.Int(1), vm.Keyword("baz"): vm.Int(2), vm.Keyword("bar"): vm.TRUE}).Cons(vm.Symbol("zoo")).Cons(vm.Symbol("with-meta")),
		"^:bar ^{:foo 1 :baz 2} zoo": vm.EmptyList.Cons(vm.Map{vm.Keyword("foo"): vm.Int(1), vm.Keyword("baz"): vm.Int(2), vm.Keyword("bar"): vm.TRUE}).Cons(vm.Symbol("zoo")).Cons(vm.Symbol("with-meta")),
	}

	for p, e := range cases {
		r := NewLispReader(strings.NewReader(p), "<reader>")
		o, err := r.Read()
		assert.NoError(t, err)
		assert.Equal(t, e, o)
	}
}

func TestSimpleCall(t *testing.T) {
	p := "(+ 40 2)"
	r := NewLispReader(strings.NewReader(p), "<reader>")
	o, err := r.Read()
	assert.NoError(t, err)

	out, err := vm.ListType.Box([]vm.Value{
		vm.Symbol("+"),
		vm.Int(40),
		vm.Int(2),
	})

	assert.NoError(t, err)
	assert.Equal(t, out, o)
}

func TestReadToken(t *testing.T) {
	r := NewLispReader(strings.NewReader("(foo bar)"), "<reader>")
	token, err := r.readToken()
	assert.NoError(t, err)
	assert.Equal(t, vm.Symbol("foo"), token)

	token, err = r.readToken()
	assert.NoError(t, err)
	assert.Equal(t, vm.Symbol("bar"), token)
}

func TestReadNumber(t *testing.T) {
	r := NewLispReader(strings.NewReader("(123 456)"), "<reader>")
	number, err := r.readNumber()
	assert.NoError(t, err)
	assert.Equal(t, vm.Int(123), number)

	number, err = r.readNumber()
	assert.NoError(t, err)
	assert.Equal(t, vm.Int(456), number)
}

func TestReadString(t *testing.T) {
	r := NewLispReader(strings.NewReader("(\"hello\" \"world\")"), "<reader>")
	str, err := r.readString()
	assert.NoError(t, err)
	assert.Equal(t, vm.String("hello"), str)

	str, err = r.readString()
	assert.NoError(t, err)
	assert.Equal(t, vm.String("world"), str)
}

func TestReadChar(t *testing.T) {
	r := NewLispReader(strings.NewReader("(\\a \\b)"), "<reader>")
	ch, err := r.readChar()
	assert.NoError(t, err)
	assert.Equal(t, vm.Char('a'), ch)

	ch, err = r.readChar()
	assert.NoError(t, err)
	assert.Equal(t, vm.Char('b'), ch)
}

func TestReadError(t *testing.T) {
	r := NewLispReader(strings.NewReader("(foo \\x)"), "<reader>")
	_, err := r.readToken()
	assert.NoError(t, err)

	_, err = r.readChar()
	assert.Error(t, err)
}

// Fix the deleted closing bracket
func TestReaderBasicError(t *testing.T) {
	r := NewLispReader(strings.NewReader("(foo"), "<reader>")
	_, err := r.readToken()
	assert.Error(t, err)
}
</user_code>

