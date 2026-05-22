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

type theStringType struct {
	zero String
}

func (t *theStringType) String() string     { return t.Name() }
func (t *theStringType) Type() ValueType    { return TypeType }
func (t *theStringType) Unbox() interface{} { return reflect.TypeOf(t) }

func (t *theStringType) Name() string { return "let-go.lang.String" }

func (t *theStringType) Box(bare interface{}) (Value, error) {
	raw, ok := bare.(string)
	if !ok {
		return StringType.zero, NewTypeError(bare, "can't be boxed as", t)
	}
	return String(raw), nil
}

// StringType is the type of StringValues
var StringType *theStringType = &theStringType{zero: ""}

// String is boxed int
type String string

func (l String) Conj(value Value) Collection {
	return String(string(l) + value.String())
}

func (l String) RawCount() int {
	return len(l)
}

func (l String) Count() Value {
	return Int(len(l))
}

func (l String) Empty() Collection {
	return String("")
}

// Hash implements Hashable for fast map lookups.
func (l String) Hash() uint32 { return hashString(string(l)) }

// Type implements Value
func (l String) Type() ValueType { return StringType }

// Unbox implements Unbox
func (l String) Unbox() interface{} {
	return string(l)
}

func (l String) InvokeMethod(name Symbol, args []Value) (Value, error) {
	switch name {
	case "replace":
		if len(args) != 2 {
			return NIL, fmt.Errorf("String.replace expected 2 arguments")
		}
		old, ok := args[0].(String)
		if !ok {
			return NIL, fmt.Errorf("String.replace expected String target")
		}
		repl, ok := args[1].(String)
		if !ok {
			return NIL, fmt.Errorf("String.replace expected String replacement")
		}
		return String(strings.ReplaceAll(string(l), string(old), string(repl))), nil
	case "getBytes":
		if len(args) == 0 {
			// No-arg form: let-go String is Go's UTF-8 string, so the
			// default encoding is UTF-8 by construction.
			return NewByteArrayFrom([]byte(string(l))), nil
		}
		if len(args) == 1 {
			enc, ok := args[0].(String)
			if !ok {
				return NIL, fmt.Errorf("String.getBytes encoding must be String, got %s", args[0].Type().Name())
			}
			e := strings.ToUpper(string(enc))
			if e != "UTF-8" && e != "UTF8" {
				return NIL, fmt.Errorf("String.getBytes: only UTF-8 supported, got %q", string(enc))
			}
			return NewByteArrayFrom([]byte(string(l))), nil
		}
		return NIL, fmt.Errorf("String.getBytes expected 0 or 1 argument, got %d", len(args))
	default:
		return NIL, fmt.Errorf("method %s not found on String", name)
	}
}

// First implements Seq
func (l String) First() Value {
	for _, r := range l {
		return Char(r)
	}
	return NIL
}

// More implements Seq
func (l String) More() Seq {
	r := l.Next()
	if r == nil {
		return EmptyList
	}
	return r
}

// Next implements Seq
func (l String) Next() Seq {
	if len(l) <= 1 {
		return nil
	}
	ret := EmptyList
	s := []rune(l)
	for i := len(s) - 1; i >= 1; i-- {
		ret = ret.Conj(Char(s[i])).(*List)
	}
	return ret
}

// Cons implements Seq
func (l String) Cons(val Value) Seq {
	return NIL
}

func (l String) Seq() Seq {
	if len(l) == 0 {
		return nil
	}
	ret := EmptyList
	s := []rune(l)
	for i := len(s) - 1; i >= 0; i-- {
		ret = ret.Conj(Char(s[i])).(*List)
	}
	return ret
}

func (l String) ValueAt(key Value) Value {
	return l.ValueAtOr(key, NIL)
}

func (l String) ValueAtOr(key Value, dflt Value) Value {
	if key == NIL {
		return dflt
	}
	r := []rune(l)
	numkey, ok := key.(Int)
	if !ok || numkey < 0 || int(numkey) >= len(r) {
		return dflt
	}
	return Char(r[numkey])
}

func (l String) String() string {
	return fmt.Sprintf("%q", string(l))
}
