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

type theKeywordType struct {
	zero Keyword
}

func (t *theKeywordType) String() string     { return t.Name() }
func (t *theKeywordType) Type() ValueType    { return TypeType }
func (t *theKeywordType) Unbox() interface{} { return reflect.TypeOf(t) }

func (t *theKeywordType) Name() string { return "let-go.lang.Keyword" }

func (t *theKeywordType) Box(bare interface{}) (Value, error) {
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

// Type implements Value
func (l Keyword) Type() ValueType { return KeywordType }

// Unbox implements Unbox
func (l Keyword) Unbox() interface{} {
	return string(l)
}

func (l Keyword) String() string {
	return fmt.Sprintf(":%s", string(l))
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
		// FIXME return error
		return NIL, fmt.Errorf("Keyword expected Lookup")
	}
	if vl == 1 {
		return as.ValueAt(l), nil
	}
	return as.ValueAtOr(l, pargs[1]), nil
}

func (l Keyword) Namespaced() (Value, Value) {
	x := strings.Split(string(l), "/")
	if len(x) == 2 {
		return Symbol(x[0]), Symbol(x[1])
	}
	return NIL, Symbol(x[0])
}

// FIXME make it work the other way round
func (l Keyword) Name() Value {
	_, n := l.Namespaced()
	if n == NIL {
		return NIL
	}
	return String(n.(Symbol))
}

func (l Keyword) Namespace() Value {
	n, _ := l.Namespaced()
	if n == NIL {
		return NIL
	}
	return String(n.(Symbol))
}
