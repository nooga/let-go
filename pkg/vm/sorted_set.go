/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"fmt"
	"strings"
)

// --- Type metadata ---

type theSortedSetType struct{}

func (t *theSortedSetType) String() string     { return t.Name() }
func (t *theSortedSetType) Type() ValueType    { return TypeType }
func (t *theSortedSetType) Unbox() interface{} { return t }
func (t *theSortedSetType) Name() string       { return "let-go.lang.PersistentTreeSet" }

func (t *theSortedSetType) Box(bare interface{}) (Value, error) {
	if s, ok := bare.(*SortedSet); ok {
		return s, nil
	}
	return NIL, NewTypeError(bare, "can't be boxed as", t)
}

var SortedSetType *theSortedSetType = &theSortedSetType{}

// --- SortedSet ---

type SortedSet struct {
	impl     *SortedMap
	meta     Value
	_hash    uint32
	_hasHash bool
}

var EmptySortedSet = &SortedSet{impl: EmptySortedMap}

func NewSortedSet(cmp Comparator, vals []Value) *SortedSet {
	if cmp == nil {
		cmp = DefaultCompare
	}
	m := &SortedMap{cmp: cmp}
	for _, v := range vals {
		m = m.assocImpl(v, v)
	}
	return &SortedSet{impl: m}
}

// --- Value ---

func (s *SortedSet) Type() ValueType    { return SortedSetType }
func (s *SortedSet) Unbox() interface{} { return s.elements() }

func (s *SortedSet) String() string {
	b := &strings.Builder{}
	b.WriteString("#{")
	seq := s.Seq()
	first := true
	for seq != nil && seq != EmptyList {
		if !first {
			b.WriteRune(' ')
		}
		b.WriteString(seq.First().String())
		first = false
		seq = seq.Next()
	}
	b.WriteRune('}')
	return b.String()
}

// --- Hashable ---

func (s *SortedSet) Hash() uint32 {
	if s._hasHash {
		return s._hash
	}
	var h uint32
	seq := s.Seq()
	for seq != nil && seq != EmptyList {
		h += hashValue(seq.First())
		seq = seq.Next()
	}
	s._hash = mixFinish(h)
	s._hasHash = true
	return s._hash
}

// --- IMeta ---

func (s *SortedSet) Meta() Value {
	if s.meta == nil {
		return NIL
	}
	return s.meta
}

func (s *SortedSet) WithMeta(m Value) Value {
	cp := *s
	cp.meta = m
	cp._hash = 0
	cp._hasHash = false
	return &cp
}

// --- Counted ---

func (s *SortedSet) RawCount() int { return s.impl.count }
func (s *SortedSet) Count() Value  { return MakeInt(s.impl.count) }

// --- Collection ---

func (s *SortedSet) Empty() Collection { return EmptySortedSet }

func (s *SortedSet) Conj(value Value) Collection {
	newImpl := s.impl.assocImpl(value, value)
	if newImpl.count == s.impl.count {
		// Check if it was actually the same — might have been a replacement
		return s
	}
	return &SortedSet{impl: newImpl}
}

// --- Set-specific ---

func (s *SortedSet) Disj(value Value) *SortedSet {
	newImpl := s.impl.Dissoc(value).(*SortedMap)
	if newImpl.count == s.impl.count {
		return s
	}
	return &SortedSet{impl: newImpl}
}

func (s *SortedSet) Contains(value Value) Boolean {
	return s.impl.Contains(value)
}

// --- Lookup ---

func (s *SortedSet) ValueAt(key Value) Value {
	if s.impl.root == nil {
		return NIL
	}
	_, found := s.impl.root.find(key, s.impl.cmp)
	if found {
		return key
	}
	return NIL
}

func (s *SortedSet) ValueAtOr(key Value, dflt Value) Value {
	if s.impl.root == nil {
		return dflt
	}
	_, found := s.impl.root.find(key, s.impl.cmp)
	if found {
		return key
	}
	return dflt
}

// --- Sequable ---

func (s *SortedSet) Seq() Seq {
	if s.impl.count == 0 {
		return EmptyList
	}
	entries := s.impl.entries()
	if len(entries) == 0 {
		return EmptyList
	}
	keys := make([]Value, len(entries))
	for i, e := range entries {
		keys[i] = e.Key
	}
	return &SortedSetSeq{keys: keys, i: 0}
}

// RSeq returns elements in reverse sorted order.
func (s *SortedSet) RSeq() Seq {
	if s.impl.count == 0 {
		return EmptyList
	}
	var entries []MapEntry
	s.impl.root.reverseOrder(&entries)
	if len(entries) == 0 {
		return EmptyList
	}
	keys := make([]Value, len(entries))
	for i, e := range entries {
		keys[i] = e.Key
	}
	return &SortedSetSeq{keys: keys, i: 0}
}

func (s *SortedSet) elements() []Value {
	entries := s.impl.entries()
	keys := make([]Value, len(entries))
	for i, e := range entries {
		keys[i] = e.Key
	}
	return keys
}

// --- Seq (for direct iteration) ---

func (s *SortedSet) First() Value {
	seq := s.Seq()
	if seq == EmptyList {
		return NIL
	}
	return seq.First()
}

func (s *SortedSet) More() Seq {
	seq := s.Seq()
	if seq == EmptyList {
		return EmptyList
	}
	return seq.More()
}

func (s *SortedSet) Next() Seq {
	seq := s.Seq()
	if seq == EmptyList {
		return nil
	}
	return seq.Next()
}

func (s *SortedSet) Cons(val Value) Seq {
	return NewCons(val, s.Seq())
}

// --- Fn (set as function) ---

func (s *SortedSet) Arity() int { return 1 }

func (s *SortedSet) Invoke(pargs []Value) (Value, error) {
	if len(pargs) != 1 {
		return NIL, fmt.Errorf("wrong number of arguments %d", len(pargs))
	}
	if s.impl.root == nil {
		return NIL, nil
	}
	_, found := s.impl.root.find(pargs[0], s.impl.cmp)
	if found {
		return pargs[0], nil
	}
	return NIL, nil
}

// --- SortedSetSeq ---

type SortedSetSeq struct {
	keys []Value
	i    int
}

func (s *SortedSetSeq) Type() ValueType    { return ListType }
func (s *SortedSetSeq) Unbox() interface{} { return s.keys[s.i:] }

func (s *SortedSetSeq) String() string {
	b := &strings.Builder{}
	b.WriteRune('(')
	for j := s.i; j < len(s.keys); j++ {
		if j > s.i {
			b.WriteRune(' ')
		}
		b.WriteString(s.keys[j].String())
	}
	b.WriteRune(')')
	return b.String()
}

func (s *SortedSetSeq) First() Value { return s.keys[s.i] }

func (s *SortedSetSeq) More() Seq {
	if s.i+1 >= len(s.keys) {
		return EmptyList
	}
	return &SortedSetSeq{keys: s.keys, i: s.i + 1}
}

func (s *SortedSetSeq) Next() Seq {
	if s.i+1 >= len(s.keys) {
		return nil
	}
	return &SortedSetSeq{keys: s.keys, i: s.i + 1}
}

func (s *SortedSetSeq) Cons(val Value) Seq {
	return NewCons(val, s)
}

func (s *SortedSetSeq) RawCount() int { return len(s.keys) - s.i }
func (s *SortedSetSeq) Count() Value  { return MakeInt(len(s.keys) - s.i) }
