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

func (t *theSymbolType) String() string     { return t.Name() }
func (t *theSymbolType) Type() ValueType    { return TypeType }
func (t *theSymbolType) Unbox() interface{} { return reflect.TypeOf(t) }

func (lt *theSymbolType) Name() string { return "let-go.lang.Symbol" }

func (lt *theSymbolType) Box(bare interface{}) (Value, error) {
	raw, ok := bare.(fmt.Stringer)
	if !ok {
		return BooleanType.zero, NewTypeError(bare, "can't be boxed as", lt)
	}
	return Symbol(raw.String()), nil
}

// SymbolType is the type of SymbolValues
var SymbolType *theSymbolType = &theSymbolType{zero: "????BADSYMBOL????"}

// Symbol is boxed int
type Symbol string

// Type implements Value
func (l Symbol) Type() ValueType { return SymbolType }

// Unbox implements Unbox
func (l Symbol) Unbox() interface{} {
	return string(l)
}

func (l Symbol) String() string {
	return string(l)
}

func (l Symbol) Namespaced() (Value, Value) {
	if string(l) == "/" {
		return NIL, l
	}
	x := strings.SplitN(string(l), "/", 2)
	if len(x) == 2 {
		return Symbol(x[0]), Symbol(x[1])
	}
	return NIL, Symbol(x[0])
}

// FIXME make it work the other way round
func (l Symbol) Name() Value {
	_, n := l.Namespaced()
	if n == NIL {
		return NIL
	}
	return String(n.(Symbol))
}

func (l Symbol) Namespace() Value {
	n, _ := l.Namespaced()
	if n == NIL {
		return NIL
	}
	return String(n.(Symbol))
}
