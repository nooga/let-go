/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"fmt"
	"strings"
)

// PersistentSet is an immutable set backed by a PersistentMap.
// Each element is stored as a key mapping to itself (like Clojure's PersistentHashSet).
type PersistentSet struct {
	impl     *PersistentMap
	meta     Value
	_hash    uint32
	_hasHash bool
}

// EmptyPersistentSet is the canonical empty set.
var EmptyPersistentSet = &PersistentSet{impl: EmptyPersistentMap}

func NewPersistentSet(vals []Value) *PersistentSet {
	if len(vals) == 0 {
		return EmptyPersistentSet
	}
	m := EmptyPersistentMap
	for _, v := range vals {
		m = m.Assoc(v, v).(*PersistentMap)
	}
	return &PersistentSet{impl: m}
}

// --- Value interface ---

func (s *PersistentSet) Type() ValueType    { return SetType }
func (s *PersistentSet) Unbox() interface{} { return s.keys() }

func (s *PersistentSet) String() string {
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

func (s *PersistentSet) Hash() uint32 {
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

func (s *PersistentSet) Meta() Value {
	if s.meta == nil {
		return NIL
	}
	return s.meta
}

func (s *PersistentSet) WithMeta(m Value) Value {
	cp := *s
	cp.meta = m
	cp._hash = 0
	cp._hasHash = false
	return &cp
}

// --- Collection ---

func (s *PersistentSet) Count() Value  { return s.impl.Count() }
func (s *PersistentSet) RawCount() int { return s.impl.RawCount() }

func (s *PersistentSet) Empty() Collection {
	return EmptyPersistentSet
}

func (s *PersistentSet) Conj(value Value) Collection {
	newImpl := s.impl.Assoc(value, value).(*PersistentMap)
	if newImpl == s.impl {
		return s // already contains
	}
	return &PersistentSet{impl: newImpl}
}

// --- Set-specific ---

func (s *PersistentSet) Disj(value Value) *PersistentSet {
	newImpl := s.impl.Dissoc(value).(*PersistentMap)
	if newImpl == s.impl {
		return s // wasn't present
	}
	return &PersistentSet{impl: newImpl}
}

func (s *PersistentSet) Contains(value Value) Boolean {
	if s.impl.root == nil {
		return FALSE
	}
	_, found := s.impl.root.find(0, hashValue(value), value)
	if found {
		return TRUE
	}
	return FALSE
}

// --- Sequable ---

func (s *PersistentSet) Seq() Seq {
	if s.impl.count == 0 {
		return EmptyList
	}
	// Materialize keys from the map's entries
	entries := s.entries()
	if len(entries) == 0 {
		return EmptyList
	}
	return &SetSeq{keys: entries, i: 0}
}

func (s *PersistentSet) keys() []Value {
	return s.entries()
}

func (s *PersistentSet) entries() []Value {
	if s.impl.root == nil {
		return nil
	}
	mes := s.impl.root.nodeSeq()
	result := make([]Value, len(mes))
	for i, e := range mes {
		result[i] = e.Key
	}
	return result
}

// --- Seq interface (for direct iteration) ---

func (s *PersistentSet) First() Value {
	seq := s.Seq()
	if seq == EmptyList {
		return NIL
	}
	return seq.First()
}

func (s *PersistentSet) More() Seq {
	seq := s.Seq()
	if seq == EmptyList {
		return EmptyList
	}
	return seq.More()
}

func (s *PersistentSet) Next() Seq {
	seq := s.Seq()
	if seq == EmptyList {
		return nil
	}
	return seq.Next()
}

func (s *PersistentSet) Cons(val Value) Seq {
	return NewCons(val, s.Seq())
}

// --- Fn (set as function) ---

func (s *PersistentSet) Arity() int { return 1 }

func (s *PersistentSet) Invoke(pargs []Value) (Value, error) {
	if len(pargs) != 1 {
		return NIL, fmt.Errorf("wrong number of arguments %d", len(pargs))
	}
	if s.impl.root == nil {
		return NIL, nil
	}
	_, found := s.impl.root.find(0, hashValue(pargs[0]), pargs[0])
	if found {
		return pargs[0], nil
	}
	return NIL, nil
}

