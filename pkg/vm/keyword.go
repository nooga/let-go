/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"fmt"
	"reflect"
)

type theKeywordType struct {
	zero Keyword
}

func (t *theKeywordType) String() string  { return t.Name() }
func (t *theKeywordType) Type() ValueType { return TypeType }
func (t *theKeywordType) Unbox() any      { return reflect.TypeFor[*theKeywordType]() }

func (t *theKeywordType) Name() string { return "let-go.lang.Keyword" }

func (t *theKeywordType) Box(bare any) (Value, error) {
	raw, ok := bare.(fmt.Stringer)
	if !ok {
		return BooleanType.zero, NewTypeError(bare, "can't be boxed as", t)
	}
	return Keyword(raw.String()), nil
}

// KeywordType is the type of KeywordValues
var KeywordType *theKeywordType = &theKeywordType{zero: "????BADKeyword????"}

// Keyword is boxed int
type Keyword string

// Hash implements Hashable for fast map lookups.
func (l Keyword) Hash() uint32 { return hashUnencodedChars(string(l)) + 0x9e3779b9 }

// Type implements Value
func (l Keyword) Type() ValueType { return KeywordType }

// Unbox implements Unbox
func (l Keyword) Unbox() any {
	return string(l)
}

func (l Keyword) String() string {
	return ":" + string(l)
}

func (l Keyword) Arity() int {
	return -1
}

func (l Keyword) Invoke(pargs []Value) (Value, error) {
	vl := len(pargs)
	if vl < 1 || vl > 2 {
		return NIL, fmt.Errorf("wrong number of arguments %d", vl)
	}
	as, ok := pargs[0].(Lookup)
	if !ok {
		if vl == 2 {
			return pargs[1], nil
		}
		return NIL, nil
	}
	if kl, ok := as.(KeywordLookup); ok {
		if vl == 1 {
			return kl.ValueAtKeyword(l), nil
		}
		return kl.ValueAtKeywordOr(l, pargs[1]), nil
	}
	if vl == 1 {
		return as.ValueAt(l), nil
	}
	return as.ValueAtOr(l, pargs[1]), nil
}

func (l Keyword) NamespacedRaw() (ns Keyword, name Keyword, hasNS bool) {
	nsString, nameString, hasNS := splitNamespaced(string(l))
	return Keyword(nsString), Keyword(nameString), hasNS
}

func (l Keyword) Namespaced() (Value, Value) {
	ns, name, hasNS := l.NamespacedRaw()
	if !hasNS {
		return NIL, Symbol(name)
	}
	return Symbol(ns), Symbol(name)
}

func (l Keyword) Name() Value {
	_, name, _ := l.NamespacedRaw()
	return String(name)
}

func (l Keyword) Namespace() Value {
	ns, _, hasNS := l.NamespacedRaw()
	if !hasNS {
		return NIL
	}
	return String(ns)
}
