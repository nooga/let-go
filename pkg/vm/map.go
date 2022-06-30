/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"reflect"
	"strings"
)

type theMapType struct{}

func (t *theMapType) String() string     { return t.Name() }
func (t *theMapType) Type() ValueType    { return TypeType }
func (t *theMapType) Unbox() interface{} { return reflect.TypeOf(t) }

func (t *theMapType) Name() string { return "let-go.lang.Map" }

func (t *theMapType) Box(bare interface{}) (Value, error) {
	casted, ok := bare.(map[Value]Value)
	if !ok {
		return NIL, NewTypeError(bare, "can't be boxed as", t)
	}

	return Map(casted), nil
}

// ArrayVectorType is the type of ArrayVectors
var MapType *theMapType

func init() {
	MapType = &theMapType{}
}

// Map is boxed singly linked list that can hold other Values.
type Map map[Value]Value

func (l Map) Conj(value Value) Collection {
	// FIXME this needs MapEntry
	if value.Type() != ArrayVectorType {
		// FIXME this is error
		return l
	}
	v := value.(ArrayVector)
	if len(v) != 2 {
		// FIXME this is error
		return l
	}
	ret := make(Map, len(l)+1)
	for k := range l {
		ret[k] = l[k]
	}
	ret[v[0]] = v[1]
	return ret
}

// Type implements Value
func (l Map) Type() ValueType { return MapType }

// Unbox implements Value
func (l Map) Unbox() interface{} {
	return map[Value]Value(l)
}

// First implements Seq
func (l Map) First() Value {
	if len(l) == 0 {
		return NIL
	}
	for k, v := range l {
		return ArrayVector{k, v}
	}
	return NIL // unreachable
}

func toList(l Map) *List {
	var lst []Value
	for k, v := range l {
		lst = append(lst, ArrayVector{k, v})
	}
	ret, _ := ListType.Box(lst)
	return ret.(*List)
}

func (l Map) Seq() Seq {
	return toList(l)
}

// More implements Seq
func (l Map) More() Seq {
	if len(l) == 1 {
		return EmptyList
	}
	ret := toList(l)
	return ret.More()
}

// Next implements Seq
func (l Map) Next() Seq {
	return l.More()
}

// Cons implements Seq
func (l Map) Cons(val Value) Seq {
	return toList(l).Cons(val)
}

// Count implements Collection
func (l Map) Count() Value {
	return Int(len(l))
}

func (l Map) RawCount() int {
	return len(l)
}

// Empty implements Collection
func (l Map) Empty() Collection {
	return make(Map)
}

func (l Map) Assoc(k Value, v Value) Associative {
	// FIXME implement persistent maps :P
	newmap := make(Map)
	for ok, ov := range l {
		newmap[ok] = ov
	}
	newmap[k] = v
	return newmap
}

func (l Map) Dissoc(k Value) Associative {
	// FIXME implement persistent maps :P
	newmap := make(Map)
	for ok, ov := range l {
		if ok == k {
			continue
		}
		newmap[ok] = ov
	}
	return newmap
}

func (l Map) ValueAt(key Value) Value {
	return l.ValueAtOr(key, NIL)
}

func (l Map) ValueAtOr(key Value, dflt Value) Value {
	if key == NIL {
		return dflt
	}
	ret, ok := l[key]
	if !ok {
		return dflt
	}
	return ret
}

func NewMap(v []Value) Value {
	if len(v) == 0 {
		return make(Map)
	}
	if len(v)%2 != 0 {
		// FIXME this is an error
		return NIL
	}
	newmap := make(Map)
	for i := 0; i < len(v); i += 2 {
		newmap[v[i]] = v[i+1]
	}
	return newmap
}

func (l Map) String() string {
	b := &strings.Builder{}
	b.WriteRune('{')
	i := 0
	n := len(l)
	for k, v := range l {
		b.WriteString(k.String())
		b.WriteRune(' ')
		b.WriteString(v.String())
		if i < n-1 {
			b.WriteRune(' ')
		}
		i++
	}
	b.WriteRune('}')
	return b.String()
}

func (l Map) Arity() int {
	return -1
}

func (l Map) Invoke(pargs []Value) Value {
	vl := len(pargs)
	if vl < 1 || vl > 2 {
		// FIXME return error
		return NIL
	}
	if vl == 1 {
		return l.ValueAt(pargs[0])
	}
	return l.ValueAtOr(pargs[0], pargs[1])
}
