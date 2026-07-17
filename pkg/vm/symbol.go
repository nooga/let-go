/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"fmt"
	"reflect"
	"strings"
)

type theSymbolType struct {
	zero Symbol
}

func (t *theSymbolType) String() string  { return t.Name() }
func (t *theSymbolType) Type() ValueType { return TypeType }
func (t *theSymbolType) Unbox() any      { return reflect.TypeFor[*theSymbolType]() }

func (t *theSymbolType) Name() string { return "let-go.lang.Symbol" }

func (t *theSymbolType) Box(bare any) (Value, error) {
	raw, ok := bare.(fmt.Stringer)
	if !ok {
		return BooleanType.zero, NewTypeError(bare, "can't be boxed as", t)
	}
	return Symbol(raw.String()), nil
}

// SymbolType is the type of Symbol values
var SymbolType *theSymbolType = &theSymbolType{zero: "????BADSYMBOL????"}

// Symbol is a string
type Symbol string

// Hash implements Hashable for fast map lookups.
func (l Symbol) Hash() uint32 { return hashUnencodedChars(string(l)) }

// Type implements Value
func (l Symbol) Type() ValueType { return SymbolType }

// Unbox implements Unbox
func (l Symbol) Unbox() any {
	return string(l)
}

func (l Symbol) String() string {
	return string(l)
}

// NamespacedRaw splits "ns/name" WITHOUT allocating: IndexByte finds the
// separator and the two parts are substrings (which share the receiver's
// string backing — no new allocation), returned as raw Symbols rather than
// boxed into Value interfaces. hasNS is false for an unqualified symbol and
// for the bare "/" symbol. This is the hot path — it does NOT touch the
// lookup-stats mutex; the boxing Namespaced() below carries that.
func (l Symbol) NamespacedRaw() (ns Symbol, name Symbol, hasNS bool) {
	s := string(l)
	if s == "/" {
		return "", l, false
	}
	i := strings.IndexByte(s, '/')
	if i < 0 {
		return "", l, false
	}
	return Symbol(s[:i]), Symbol(s[i+1:]), true
}

func (l Symbol) Namespaced() (Value, Value) {
	noteNamespaced(string(l))
	ns, name, hasNS := l.NamespacedRaw()
	if !hasNS {
		return NIL, name
	}
	return ns, name
}

func (l Symbol) Name() Value {
	_, name, _ := l.NamespacedRaw()
	return String(name)
}

func (l Symbol) Namespace() Value {
	ns, _, hasNS := l.NamespacedRaw()
	if !hasNS {
		return NIL
	}
	return String(ns)
}
