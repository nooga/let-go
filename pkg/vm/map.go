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
var MapType *theMapType = &theMapType{}

// Map is boxed singly linked list that can hold other Values.
type Map map[Value]Value

func (l Map) Conj(value Value) Collection {
	if value.Type() != ArrayVectorType {
		return l
	}
	v := value.(ArrayVector)
	if len(v) != 2 {
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

func (l Map) toEntries() []Value {
	entries := make([]Value, 0, len(l))
	for k, v := range l {
		entries = append(entries, ArrayVector{k, v})
	}
	return entries
}

func (l Map) Seq() Seq {
	if len(l) == 0 {
		return EmptyList
	}
	return &MapSeq{entries: l.toEntries(), i: 0}
}

// More implements Seq
func (l Map) More() Seq {
	if len(l) <= 1 {
		return EmptyList
	}
	return &MapSeq{entries: l.toEntries(), i: 1}
}

// Next implements Seq
func (l Map) Next() Seq {
	if len(l) <= 1 {
		return nil
	}
	return &MapSeq{entries: l.toEntries(), i: 1}
}

// Cons implements Seq
func (l Map) Cons(val Value) Seq {
	return NewCons(val, l.Seq())
}

// MapSeq is a lightweight seq view over a Map's entries
type MapSeq struct {
	entries []Value // cached [k v] pairs
	i       int
}

func (s *MapSeq) String() string {
	b := &strings.Builder{}
	b.WriteRune('(')
	for i := s.i; i < len(s.entries); i++ {
		if i > s.i {
			b.WriteRune(' ')
		}
		b.WriteString(s.entries[i].String())
	}
	b.WriteRune(')')
	return b.String()
}

func (s *MapSeq) Type() ValueType    { return ListType }
func (s *MapSeq) Unbox() interface{} { return s.entries[s.i:] }

func (s *MapSeq) First() Value {
	if s.i >= len(s.entries) {
		return NIL
	}
	return s.entries[s.i]
}

func (s *MapSeq) More() Seq {
	if s.i+1 >= len(s.entries) {
		return EmptyList
	}
	return &MapSeq{entries: s.entries, i: s.i + 1}
}

func (s *MapSeq) Next() Seq {
	if s.i+1 >= len(s.entries) {
		return nil
	}
	return &MapSeq{entries: s.entries, i: s.i + 1}
}

func (s *MapSeq) Cons(val Value) Seq {
	return NewCons(val, s)
}

func (s *MapSeq) Count() Value    { return Int(len(s.entries) - s.i) }
func (s *MapSeq) RawCount() int   { return len(s.entries) - s.i }
func (s *MapSeq) Empty() Collection { return EmptyList }
func (s *MapSeq) Conj(val Value) Collection { return s.Cons(val).(*List) }
func (s *MapSeq) Seq() Seq { return s }

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
	newmap := make(Map)
	for ok, ov := range l {
		newmap[ok] = ov
	}
	newmap[k] = v
	return newmap
}

func (l Map) Dissoc(k Value) Associative {
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

func (l Map) Contains(value Value) Boolean {
	if _, ok := l[value]; ok {
		return TRUE
	}
	return FALSE
}

func NewMap(v []Value) Value {
	if len(v) == 0 {
		return EmptyPersistentMap
	}
	if len(v)%2 != 0 {
		return NIL
	}
	return NewPersistentMap(v)
}

// MapFromGoMap converts a Go map[Value]Value to a *PersistentMap.
func MapFromGoMap(m map[Value]Value) *PersistentMap {
	result := EmptyPersistentMap
	for k, v := range m {
		result = result.Assoc(k, v).(*PersistentMap)
	}
	return result
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

func (l Map) Invoke(pargs []Value) (Value, error) {
	vl := len(pargs)
	if vl < 1 || vl > 2 {
		return NIL, fmt.Errorf("wrong number of arguments %d", vl)
	}
	if vl == 1 {
		return l.ValueAt(pargs[0]), nil
	}
	return l.ValueAtOr(pargs[0], pargs[1]), nil
}
