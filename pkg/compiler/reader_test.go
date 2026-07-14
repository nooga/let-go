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
		"%":                          vm.Symbol("%"),
		"%1":                         vm.Symbol("%1"),
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
		"^:foo zoo":                  vm.EmptyList.Cons(vm.NewPersistentMap([]vm.Value{vm.Keyword("foo"), vm.TRUE})).Cons(vm.Symbol("zoo")).Cons(vm.Symbol("with-meta")),
		"^:foo ^:bar zoo":            vm.EmptyList.Cons(vm.NewPersistentMap([]vm.Value{vm.Keyword("foo"), vm.TRUE, vm.Keyword("bar"), vm.TRUE})).Cons(vm.Symbol("zoo")).Cons(vm.Symbol("with-meta")),
		"^{:foo 1 :baz 2} ^:bar zoo": vm.EmptyList.Cons(vm.NewPersistentMap([]vm.Value{vm.Keyword("foo"), vm.Int(1), vm.Keyword("baz"), vm.Int(2), vm.Keyword("bar"), vm.TRUE})).Cons(vm.Symbol("zoo")).Cons(vm.Symbol("with-meta")),
		"^:bar ^{:foo 1 :baz 2} zoo": vm.EmptyList.Cons(vm.NewPersistentMap([]vm.Value{vm.Keyword("bar"), vm.TRUE, vm.Keyword("foo"), vm.Int(1), vm.Keyword("baz"), vm.Int(2)})).Cons(vm.Symbol("zoo")).Cons(vm.Symbol("with-meta")),
	}

	for p, e := range cases {
		r := NewLispReader(strings.NewReader(p), "<reader>")
		o, err := r.Read()
		assert.NoError(t, err)
		assert.Equal(t, e, o)
	}
}

func TestReaderKeywordInternalColons(t *testing.T) {
	valid := map[string]vm.Keyword{
		`:/`:        vm.Keyword(`/`),
		`:.:.`:      vm.Keyword(`.:.`),
		`:a:b`:      vm.Keyword(`a:b`),
		`:a:b/c`:    vm.Keyword(`a:b/c`),
		`:foo/:bar`: vm.Keyword(`foo/:bar`),
	}
	for input, want := range valid {
		r := NewLispReader(strings.NewReader(input), "keyword-colons.lg")
		got, err := r.Read()
		assert.NoError(t, err, input)
		assert.Equal(t, want, got, input)
	}

	for _, input := range []string{`:a:`, `:a::b`, `:a:/b`} {
		r := NewLispReader(strings.NewReader(input), "keyword-colons.lg")
		_, err := r.Read()
		assert.Error(t, err, input)
	}

	current := rt.CurrentNS.Deref().(*vm.Namespace)
	current.Alias(vm.Symbol("keyword-colon-alias"), vm.NewNamespace("keyword-target"))
	r := NewLispReader(strings.NewReader(`::keyword-colon-alias/:bar`), "keyword-colons.lg")
	got, err := r.Read()
	assert.NoError(t, err)
	assert.Equal(t, vm.Keyword(`keyword-target/:bar`), got)
}

func TestReaderSkipsLeadingNoValueForms(t *testing.T) {
	// ReadSkipNoValue (the read-string entry) skips a leading no-value reader
	// macro (line comment, #_ discard) and returns the next real form, not the
	// VOID sentinel. Regression: read-string of a string that begins with a
	// ';;' comment used to return VOID. (Read itself still returns VOID — see
	// TestReaderMap* below, which collection readers depend on.)
	want := vm.EmptyList.Cons(vm.Int(2)).Cons(vm.Int(1)).Cons(vm.Symbol("+"))
	cases := map[string]vm.Value{
		";; leading comment\n(+ 1 2)":       want,
		";; c with ) paren inside\n(+ 1 2)": want,
		"   ;; indented comment\n(+ 1 2)":   want,
		";; one\n;; two\n(+ 1 2)":           want,
		"#_ discarded (+ 1 2)":              want,
	}
	for p, e := range cases {
		r := NewLispReader(strings.NewReader(p), "<reader>")
		o, err := r.ReadSkipNoValue()
		assert.NoError(t, err, "input %q", p)
		assert.Equal(t, e, o, "input %q", p)
	}
}

func TestReaderEOFAfterCommentOnly(t *testing.T) {
	// Comment-only / whitespace-only input has no form: ReadSkipNoValue must
	// surface EOF rather than silently returning VOID.
	for _, p := range []string{";; only a comment", "   \n  ", "#_ 1"} {
		r := NewLispReader(strings.NewReader(p), "<reader>")
		_, err := r.ReadSkipNoValue()
		assert.Error(t, err, "input %q should be EOF", p)
	}
}

func TestReaderMapDiscardDoesNotConsumeNextKey(t *testing.T) {
	// Read must keep returning VOID in value position so readMap can drop the
	// orphaned key: {:a #_ 1 :b 2} => {:b 2}. Skipping VOID here would read :b
	// as :a's value and leave 2 as an odd key.
	r := NewLispReader(strings.NewReader("{:a #_ 1 :b 2}"), "<reader>")
	o, err := r.Read()
	assert.NoError(t, err)
	assert.Equal(t, vm.NewArrayMap([]vm.Value{vm.Keyword("b"), vm.Int(2)}), o)
}

func TestReaderMapUnmatchedConditionalDoesNotConsumeNextKey(t *testing.T) {
	// Same invariant for a reader conditional with no matching branch.
	r := NewLispReader(strings.NewReader("{:a #?(:cljs 1) :b 2}"), "<reader>")
	o, err := r.Read()
	assert.NoError(t, err)
	assert.Equal(t, vm.NewArrayMap([]vm.Value{vm.Keyword("b"), vm.Int(2)}), o)
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

func TestReaderConditionalSplicing(t *testing.T) {
	cases := map[string]vm.Value{
		"(a #?@(:cljs [nil] :default []) b)": vm.EmptyList.Cons(vm.Symbol("b")).Cons(vm.Symbol("a")),
		"(a #?@(:cljs [] :default [x y]) b)": vm.EmptyList.Cons(vm.Symbol("b")).Cons(vm.Symbol("y")).Cons(vm.Symbol("x")).Cons(vm.Symbol("a")),
	}

	for p, e := range cases {
		r := NewLispReader(strings.NewReader(p), "<reader>")
		o, err := r.Read()
		assert.NoError(t, err)
		assert.Equal(t, e, o)
	}
}
