/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"fmt"
	"reflect"
	"strings"
)

// --- Type metadata ---

type theSortedMapType struct{}

func (t *theSortedMapType) String() string     { return t.Name() }
func (t *theSortedMapType) Type() ValueType    { return TypeType }
func (t *theSortedMapType) Unbox() interface{} { return reflect.TypeOf(t) }
func (t *theSortedMapType) Name() string       { return "let-go.lang.PersistentTreeMap" }

func (t *theSortedMapType) Box(bare interface{}) (Value, error) {
	if m, ok := bare.(*SortedMap); ok {
		return m, nil
	}
	return NIL, NewTypeError(bare, "can't be boxed as", t)
}

var SortedMapType *theSortedMapType = &theSortedMapType{}

// --- LLRB tree node ---

const (
	rbRed   = true
	rbBlack = false
)

type rbNode struct {
	key   Value
	val   Value
	left  *rbNode
	right *rbNode
	color bool // true = red, false = black
}

func newRBNode(key, val Value) *rbNode {
	return &rbNode{key: key, val: val, color: rbRed}
}

func isRed(n *rbNode) bool {
	return n != nil && n.color
}

func (n *rbNode) rotateLeft() *rbNode {
	x := n.right
	n.right = x.left
	x.left = n
	x.color = n.color
	n.color = rbRed
	return x
}

func (n *rbNode) rotateRight() *rbNode {
	x := n.left
	n.left = x.right
	x.right = n
	x.color = n.color
	n.color = rbRed
	return x
}

func (n *rbNode) flipColors() {
	n.color = !n.color
	n.left.color = !n.left.color
	n.right.color = !n.right.color
}

func (n *rbNode) fixUp() *rbNode {
	if isRed(n.right) && !isRed(n.left) {
		n = n.rotateLeft()
	}
	if isRed(n.left) && isRed(n.left.left) {
		n = n.rotateRight()
	}
	if isRed(n.left) && isRed(n.right) {
		n.flipColors()
	}
	return n
}

// insert returns a new tree with key/val inserted. replaced is true if key already existed.
func (n *rbNode) insert(key, val Value, cmp Comparator) (*rbNode, bool) {
	if n == nil {
		return newRBNode(key, val), false
	}
	c, err := cmp(key, n.key)
	if err != nil {
		// Fall back to hash comparison for incomparable types
		h1, h2 := hashValue(key), hashValue(n.key)
		if h1 < h2 {
			c = -1
		} else if h1 > h2 {
			c = 1
		} else {
			c = 0
		}
	}

	var replaced bool
	switch {
	case c < 0:
		n = n.clone()
		n.left, replaced = n.left.insert(key, val, cmp)
	case c > 0:
		n = n.clone()
		n.right, replaced = n.right.insert(key, val, cmp)
	default:
		n = n.clone()
		n.val = val
		replaced = true
	}
	return n.fixUp(), replaced
}

func (n *rbNode) clone() *rbNode {
	cp := *n
	return &cp
}

func (n *rbNode) find(key Value, cmp Comparator) (Value, bool) {
	cur := n
	for cur != nil {
		c, err := cmp(key, cur.key)
		if err != nil {
			h1, h2 := hashValue(key), hashValue(cur.key)
			if h1 < h2 {
				c = -1
			} else if h1 > h2 {
				c = 1
			} else {
				c = 0
			}
		}
		switch {
		case c < 0:
			cur = cur.left
		case c > 0:
			cur = cur.right
		default:
			return cur.val, true
		}
	}
	return NIL, false
}

// moveRedLeft moves red from right to left during deletion.
func (n *rbNode) moveRedLeft() *rbNode {
	n.flipColors()
	if n.right != nil && isRed(n.right.left) {
		n.right = n.right.rotateRight()
		n = n.rotateLeft()
		n.flipColors()
	}
	return n
}

// moveRedRight moves red from left to right during deletion.
func (n *rbNode) moveRedRight() *rbNode {
	n.flipColors()
	if n.left != nil && isRed(n.left.left) {
		n = n.rotateRight()
		n.flipColors()
	}
	return n
}

func (n *rbNode) min() *rbNode {
	cur := n
	for cur.left != nil {
		cur = cur.left
	}
	return cur
}

func (n *rbNode) deleteMin() *rbNode {
	if n.left == nil {
		return nil
	}
	n = n.clone()
	if !isRed(n.left) && !isRed(n.left.left) {
		n = n.moveRedLeft()
	}
	n.left = n.left.deleteMin()
	return n.fixUp()
}

// delete removes key from tree. Returns new root and whether key was found.
func (n *rbNode) delete(key Value, cmp Comparator) (*rbNode, bool) {
	if n == nil {
		return nil, false
	}

	c, err := cmp(key, n.key)
	if err != nil {
		h1, h2 := hashValue(key), hashValue(n.key)
		if h1 < h2 {
			c = -1
		} else if h1 > h2 {
			c = 1
		} else {
			c = 0
		}
	}

	found := false
	n = n.clone()

	if c < 0 {
		if n.left == nil {
			return n, false
		}
		if !isRed(n.left) && !isRed(n.left.left) {
			n = n.moveRedLeft()
		}
		n.left, found = n.left.delete(key, cmp)
	} else {
		if isRed(n.left) {
			n = n.rotateRight()
			// Recompute comparison after rotation
			c, err = cmp(key, n.key)
			if err != nil {
				h1, h2 := hashValue(key), hashValue(n.key)
				if h1 < h2 {
					c = -1
				} else if h1 > h2 {
					c = 1
				} else {
					c = 0
				}
			}
		}
		if c == 0 && n.right == nil {
			return nil, true
		}
		if n.right != nil && !isRed(n.right) && !isRed(n.right.left) {
			n = n.moveRedRight()
			// Recompute comparison after restructuring
			c, err = cmp(key, n.key)
			if err != nil {
				h1, h2 := hashValue(key), hashValue(n.key)
				if h1 < h2 {
					c = -1
				} else if h1 > h2 {
					c = 1
				} else {
					c = 0
				}
			}
		}
		if c == 0 {
			found = true
			if n.right != nil {
				m := n.right.min()
				n.key = m.key
				n.val = m.val
				n.right = n.right.deleteMin()
			} else {
				return nil, true
			}
		} else {
			if n.right != nil {
				n.right, found = n.right.delete(key, cmp)
			}
		}
	}
	return n.fixUp(), found
}

// inorder collects entries in sorted order.
func (n *rbNode) inorder(out *[]MapEntry) {
	if n == nil {
		return
	}
	n.left.inorder(out)
	*out = append(*out, MapEntry{Key: n.key, Value: n.val})
	n.right.inorder(out)
}

// reverseOrder collects entries in reverse sorted order.
func (n *rbNode) reverseOrder(out *[]MapEntry) {
	if n == nil {
		return
	}
	n.right.reverseOrder(out)
	*out = append(*out, MapEntry{Key: n.key, Value: n.val})
	n.left.reverseOrder(out)
}

// --- SortedMap ---

type SortedMap struct {
	root     *rbNode
	count    int
	cmp      Comparator
	meta     Value
	_hash    uint32
	_hasHash bool
}

var EmptySortedMap = &SortedMap{cmp: DefaultCompare}

func NewSortedMap(cmp Comparator, kvs []Value) *SortedMap {
	if cmp == nil {
		cmp = DefaultCompare
	}
	m := &SortedMap{cmp: cmp}
	for i := 0; i+1 < len(kvs); i += 2 {
		m = m.assocImpl(kvs[i], kvs[i+1])
	}
	return m
}

func (m *SortedMap) assocImpl(key, val Value) *SortedMap {
	newRoot, replaced := m.root.insert(key, val, m.cmp)
	if newRoot != nil {
		newRoot.color = rbBlack
	}
	newCount := m.count
	if !replaced {
		newCount++
	}
	return &SortedMap{root: newRoot, count: newCount, cmp: m.cmp}
}

// --- Value ---

func (m *SortedMap) Type() ValueType    { return SortedMapType }
func (m *SortedMap) Unbox() interface{} { return m.entries() }

func (m *SortedMap) String() string {
	b := &strings.Builder{}
	b.WriteRune('{')
	entries := m.entries()
	for i, e := range entries {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(e.Key.String())
		b.WriteRune(' ')
		b.WriteString(e.Value.String())
	}
	b.WriteRune('}')
	return b.String()
}

// --- Hashable ---

func (m *SortedMap) Hash() uint32 {
	if m._hasHash {
		return m._hash
	}
	var h uint32
	entries := m.entries()
	for _, e := range entries {
		h += hashValue(e.Key) ^ hashValue(e.Value)
	}
	m._hash = mixFinish(h)
	m._hasHash = true
	return m._hash
}

// --- IMeta ---

func (m *SortedMap) Meta() Value {
	if m.meta == nil {
		return NIL
	}
	return m.meta
}

func (m *SortedMap) WithMeta(meta Value) Value {
	cp := *m
	cp.meta = meta
	cp._hash = 0
	cp._hasHash = false
	return &cp
}

// --- Counted ---

func (m *SortedMap) RawCount() int { return m.count }
func (m *SortedMap) Count() Value  { return MakeInt(m.count) }

// --- Collection ---

func (m *SortedMap) Empty() Collection { return EmptySortedMap }

func (m *SortedMap) Conj(value Value) Collection {
	// Accept [k v] vectors or MapEntry-like things
	if av, ok := value.(ArrayVector); ok && len(av) == 2 {
		return m.assocImpl(av[0], av[1])
	}
	// Accept another map
	if om, ok := value.(*SortedMap); ok {
		result := m
		entries := om.entries()
		for _, e := range entries {
			result = result.assocImpl(e.Key, e.Value)
		}
		return result
	}
	if om, ok := value.(*PersistentMap); ok {
		result := m
		entries := om.entries()
		for _, e := range entries {
			if av, ok := e.(ArrayVector); ok && len(av) == 2 {
				result = result.assocImpl(av[0], av[1])
			}
		}
		return result
	}
	// Try as seq of [k v]
	if s, ok := value.(Sequable); ok {
		seq := s.Seq()
		if seq != nil && seq != EmptyList {
			f := seq.First()
			if av, ok := f.(ArrayVector); ok && len(av) == 2 {
				return m.assocImpl(av[0], av[1])
			}
		}
	}
	return m
}

// --- Associative ---

func (m *SortedMap) Assoc(key, val Value) Associative {
	return m.assocImpl(key, val)
}

func (m *SortedMap) Dissoc(key Value) Associative {
	if m.root == nil {
		return m
	}
	newRoot, found := m.root.delete(key, m.cmp)
	if !found {
		return m
	}
	if newRoot != nil {
		newRoot.color = rbBlack
	}
	return &SortedMap{root: newRoot, count: m.count - 1, cmp: m.cmp}
}

// --- Lookup ---

func (m *SortedMap) ValueAt(key Value) Value {
	if m.root == nil {
		return NIL
	}
	v, found := m.root.find(key, m.cmp)
	if !found {
		return NIL
	}
	return v
}

func (m *SortedMap) ValueAtOr(key, dflt Value) Value {
	if m.root == nil {
		return dflt
	}
	v, found := m.root.find(key, m.cmp)
	if !found {
		return dflt
	}
	return v
}

// --- Keyed ---

func (m *SortedMap) Contains(key Value) Boolean {
	if m.root == nil {
		return FALSE
	}
	_, found := m.root.find(key, m.cmp)
	if found {
		return TRUE
	}
	return FALSE
}

// --- Sequable ---

func (m *SortedMap) Seq() Seq {
	if m.count == 0 {
		return EmptyList
	}
	entries := m.entries()
	if len(entries) == 0 {
		return EmptyList
	}
	return &SortedMapSeq{entries: entries, i: 0}
}

// RSeq returns entries in reverse sorted order.
func (m *SortedMap) RSeq() Seq {
	if m.count == 0 {
		return EmptyList
	}
	var entries []MapEntry
	m.root.reverseOrder(&entries)
	if len(entries) == 0 {
		return EmptyList
	}
	return &SortedMapSeq{entries: entries, i: 0}
}

func (m *SortedMap) entries() []MapEntry {
	if m.root == nil {
		return nil
	}
	var entries []MapEntry
	m.root.inorder(&entries)
	return entries
}

// --- Fn (map as function) ---

func (m *SortedMap) Arity() int { return -1 }

func (m *SortedMap) Invoke(args []Value) (Value, error) {
	switch len(args) {
	case 1:
		return m.ValueAt(args[0]), nil
	case 2:
		return m.ValueAtOr(args[0], args[1]), nil
	}
	return NIL, fmt.Errorf("wrong number of arguments %d to sorted-map lookup", len(args))
}

// --- Seq implementation for sorted entries ---

// SortedMapSeq implements Seq over sorted map entries.
type SortedMapSeq struct {
	entries []MapEntry
	i       int
}

func (s *SortedMapSeq) Type() ValueType    { return ListType }
func (s *SortedMapSeq) Unbox() interface{} { return s.entries[s.i:] }

func (s *SortedMapSeq) String() string {
	b := &strings.Builder{}
	b.WriteRune('(')
	for j := s.i; j < len(s.entries); j++ {
		if j > s.i {
			b.WriteRune(' ')
		}
		b.WriteRune('[')
		b.WriteString(s.entries[j].Key.String())
		b.WriteRune(' ')
		b.WriteString(s.entries[j].Value.String())
		b.WriteRune(']')
	}
	b.WriteRune(')')
	return b.String()
}

func (s *SortedMapSeq) First() Value {
	e := s.entries[s.i]
	return ArrayVector{e.Key, e.Value}
}

func (s *SortedMapSeq) More() Seq {
	if s.i+1 >= len(s.entries) {
		return EmptyList
	}
	return &SortedMapSeq{entries: s.entries, i: s.i + 1}
}

func (s *SortedMapSeq) Next() Seq {
	if s.i+1 >= len(s.entries) {
		return nil
	}
	return &SortedMapSeq{entries: s.entries, i: s.i + 1}
}

func (s *SortedMapSeq) Cons(val Value) Seq {
	return NewCons(val, s)
}

func (s *SortedMapSeq) RawCount() int { return len(s.entries) - s.i }
func (s *SortedMapSeq) Count() Value  { return MakeInt(len(s.entries) - s.i) }
