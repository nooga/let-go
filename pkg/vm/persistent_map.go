/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"fmt"
	"math/bits"
	"reflect"
	"strings"
	"sync/atomic"
)

const (
	hmapShift = 5
	hmapMask  = 0x1f // 31
)

// --- Type metadata ---

type thePersistentMapType struct{}

func (t *thePersistentMapType) String() string  { return t.Name() }
func (t *thePersistentMapType) Type() ValueType { return TypeType }
func (t *thePersistentMapType) Unbox() any      { return reflect.TypeFor[*thePersistentMapType]() }
func (t *thePersistentMapType) Name() string    { return "let-go.lang.PersistentHashMap" }

func (t *thePersistentMapType) Box(bare any) (Value, error) {
	if m, ok := bare.(*PersistentMap); ok {
		return m, nil
	}
	return NIL, NewTypeError(bare, "can't be boxed as", t)
}

// PersistentMapType is the type of PersistentMap values.
var PersistentMapType *thePersistentMapType = &thePersistentMapType{}

// --- Node interface ---

// hmapNode is the internal node interface for the HAMT.
type hmapNode interface {
	find(shift uint, hash uint32, key Value) (Value, bool)
	findKeyword(shift uint, hash uint32, k Keyword) (Value, bool)
	assoc(shift uint, hash uint32, key Value, val Value, addedLeaf *bool) hmapNode
	dissoc(shift uint, hash uint32, key Value) hmapNode
	nodeSeq() []MapEntry
	// each visits every (key, value) pair in place without allocating an
	// entry slice. It returns false as soon as fn returns false (allowing
	// short-circuit), true once the whole subtree has been visited.
	each(fn func(key, val Value) bool) bool
}

// MapEntry is a key-value pair.
type MapEntry struct {
	Key   Value
	Value Value
}

func (e MapEntry) Type() ValueType { return ArrayVectorType }
func (e MapEntry) Unbox() any      { return []Value{e.Key, e.Value} }
func (e MapEntry) String() string {
	return "[" + e.Key.String() + " " + e.Value.String() + "]"
}
func (e MapEntry) First() Value { return e.Key }
func (e MapEntry) More() Seq    { return EmptyList.Cons(e.Value) }
func (e MapEntry) Next() Seq    { return EmptyList.Cons(e.Value) }
func (e MapEntry) Cons(v Value) Seq {
	return NewCons(v, e.Seq())
}
func (e MapEntry) Seq() Seq {
	return ArrayVector{e.Key, e.Value}.Seq()
}
func (e MapEntry) Count() Value  { return Int(2) }
func (e MapEntry) RawCount() int { return 2 }
func (e MapEntry) Empty() Collection {
	return ArrayVector{}
}
func (e MapEntry) Conj(v Value) Collection {
	return ArrayVector{e.Key, e.Value}.Conj(v)
}
func (e MapEntry) ValueAt(k Value) Value {
	return e.ValueAtOr(k, NIL)
}
func (e MapEntry) ValueAtOr(k Value, dflt Value) Value {
	idx, ok := k.(Int)
	if !ok {
		return dflt
	}
	switch idx {
	case 0:
		return e.Key
	case 1:
		return e.Value
	default:
		return dflt
	}
}
func (e MapEntry) Contains(k Value) Boolean {
	idx, ok := k.(Int)
	return Boolean(ok && idx >= 0 && idx < 2)
}
func (e MapEntry) Equals(other Value) bool {
	return ArrayVector{e.Key, e.Value}.Equals(other)
}
func (e MapEntry) Hash() uint32 {
	return ArrayVector{e.Key, e.Value}.Hash()
}

func MapEntryKV(entry Value) (Value, Value, bool) {
	switch e := entry.(type) {
	case MapEntry:
		return e.Key, e.Value, true
	case ArrayVector:
		if len(e) == 2 {
			return e[0], e[1], true
		}
	}
	return NIL, NIL, false
}

// --- Bit helpers ---

func hmapMaskFn(hash uint32, shift uint) uint32 {
	return (hash >> shift) & hmapMask
}

func hmapBitpos(hash uint32, shift uint) uint32 {
	return 1 << hmapMaskFn(hash, shift)
}

func hmapIndex(bitmap uint32, bit uint32) int {
	return bits.OnesCount32(bitmap & (bit - 1))
}

// --- Bitmap Indexed Node (Clojure's BitmapIndexedNode) ---

type hmapBitmapNode struct {
	bitmap uint32
	array  []any        // pairs: [key-or-nil, val-or-child, ...]
	edit   *atomic.Bool // non-nil for transient-owned nodes
	// If array[2*i] == nil, then array[2*i+1] is an hmapNode child.
	// If array[2*i] != nil, then array[2*i] is key and array[2*i+1] is value.
}

func (n *hmapBitmapNode) find(shift uint, hash uint32, key Value) (Value, bool) {
	bit := hmapBitpos(hash, shift)
	if n.bitmap&bit == 0 {
		return nil, false
	}
	idx := hmapIndex(n.bitmap, bit)
	keyOrNil := n.array[2*idx]
	valOrNode := n.array[2*idx+1]
	if keyOrNil == nil {
		// Sub-node
		return valOrNode.(hmapNode).find(shift+hmapShift, hash, key)
	}
	if valueEquiv(key, keyOrNil.(Value)) {
		return valOrNode.(Value), true
	}
	return nil, false
}

func (n *hmapBitmapNode) findKeyword(shift uint, hash uint32, k Keyword) (Value, bool) {
	bit := hmapBitpos(hash, shift)
	if n.bitmap&bit == 0 {
		return nil, false
	}
	idx := hmapIndex(n.bitmap, bit)
	keyOrNil := n.array[2*idx]
	valOrNode := n.array[2*idx+1]
	if keyOrNil == nil {
		return valOrNode.(hmapNode).findKeyword(shift+hmapShift, hash, k)
	}
	// Box-free compare: a keyword only ever equals a keyword with the same
	// string; non-keyword stored keys (incl. hash collisions) never match.
	if sk, ok := keyOrNil.(Keyword); ok && sk == k {
		return valOrNode.(Value), true
	}
	return nil, false
}

func (n *hmapBitmapNode) assoc(shift uint, hash uint32, key Value, val Value, addedLeaf *bool) hmapNode {
	bit := hmapBitpos(hash, shift)
	idx := hmapIndex(n.bitmap, bit)

	if n.bitmap&bit != 0 {
		// Slot exists
		keyOrNil := n.array[2*idx]
		valOrNode := n.array[2*idx+1]

		if keyOrNil == nil {
			// Sub-node — recurse
			child := valOrNode.(hmapNode)
			newChild := child.assoc(shift+hmapShift, hash, key, val, addedLeaf)
			if newChild == child {
				return n
			}
			return n.cloneAndSet(2*idx+1, newChild)
		}

		existingKey := keyOrNil.(Value)
		if valueEquiv(key, existingKey) {
			// Same key — replace value
			existingVal := valOrNode.(Value)
			if valueEquiv(existingVal, val) {
				return n
			}
			return n.cloneAndSet(2*idx+1, val)
		}

		// Hash collision at this level — push both entries down
		*addedLeaf = true
		existingHash := hashValue(existingKey)
		return n.cloneAndSet2(
			2*idx, nil,
			2*idx+1, createNode(shift+hmapShift, existingHash, existingKey, valOrNode.(Value), hash, key, val),
		)
	}

	// New entry — expand the array
	*addedLeaf = true
	nEntries := bits.OnesCount32(n.bitmap)
	newArray := make([]any, 2*(nEntries+1))
	// Copy entries before idx
	copy(newArray, n.array[:2*idx])
	// Insert new entry
	newArray[2*idx] = key
	newArray[2*idx+1] = val
	// Copy entries after idx
	copy(newArray[2*(idx+1):], n.array[2*idx:])
	return &hmapBitmapNode{
		bitmap: n.bitmap | bit,
		array:  newArray,
	}
}

func (n *hmapBitmapNode) dissoc(shift uint, hash uint32, key Value) hmapNode {
	bit := hmapBitpos(hash, shift)
	if n.bitmap&bit == 0 {
		return n // Key not present
	}
	idx := hmapIndex(n.bitmap, bit)
	keyOrNil := n.array[2*idx]
	valOrNode := n.array[2*idx+1]

	if keyOrNil == nil {
		// Sub-node — recurse
		child := valOrNode.(hmapNode)
		newChild := child.dissoc(shift+hmapShift, hash, key)
		if newChild == child {
			return n
		}
		if newChild != nil {
			return n.cloneAndSet(2*idx+1, newChild)
		}
		// Child became empty — remove this slot
		if n.bitmap == bit {
			// This was the only entry
			return nil
		}
		return n.removePair(idx)
	}

	if valueEquiv(key, keyOrNil.(Value)) {
		// Found key — remove it
		if n.bitmap == bit {
			return nil
		}
		return n.removePair(idx)
	}

	return n // Key not present
}

func (n *hmapBitmapNode) nodeSeq() []MapEntry {
	var entries []MapEntry
	nSlots := len(n.array) / 2
	for i := range nSlots {
		keyOrNil := n.array[2*i]
		valOrNode := n.array[2*i+1]
		if keyOrNil != nil {
			entries = append(entries, MapEntry{Key: keyOrNil.(Value), Value: valOrNode.(Value)})
		} else if valOrNode != nil {
			entries = append(entries, valOrNode.(hmapNode).nodeSeq()...)
		}
	}
	if allocAttrEnabled {
		recordAllocAttr(akMapNodeSeq, len(entries)*32+24)
	}
	return entries
}

func (n *hmapBitmapNode) each(fn func(key, val Value) bool) bool {
	nSlots := len(n.array) / 2
	for i := 0; i < nSlots; i++ {
		keyOrNil := n.array[2*i]
		valOrNode := n.array[2*i+1]
		if keyOrNil != nil {
			if !fn(keyOrNil.(Value), valOrNode.(Value)) {
				return false
			}
		} else if valOrNode != nil {
			if !valOrNode.(hmapNode).each(fn) {
				return false
			}
		}
	}
	return true
}

func (n *hmapBitmapNode) cloneAndSet(i int, val any) *hmapBitmapNode {
	if allocAttrEnabled {
		recordAllocAttr(akMapCloneAndSet, len(n.array)*16+48)
	}
	newArray := make([]any, len(n.array))
	copy(newArray, n.array)
	newArray[i] = val
	return &hmapBitmapNode{bitmap: n.bitmap, array: newArray}
}

func (n *hmapBitmapNode) cloneAndSet2(i int, a any, j int, b any) *hmapBitmapNode {
	if allocAttrEnabled {
		recordAllocAttr(akMapCloneAndSet, len(n.array)*16+48)
	}
	newArray := make([]any, len(n.array))
	copy(newArray, n.array)
	newArray[i] = a
	newArray[j] = b
	return &hmapBitmapNode{bitmap: n.bitmap, array: newArray}
}

// editAndSet mutates in place if this node is owned by the given edit token,
// otherwise clones and sets.
func (n *hmapBitmapNode) editAndSet(edit *atomic.Bool, i int, val any) *hmapBitmapNode {
	if n.edit == edit {
		n.array[i] = val
		return n
	}
	return n.cloneAndSet(i, val)
}

func (n *hmapBitmapNode) editAndSet2(edit *atomic.Bool, i int, a any, j int, b any) *hmapBitmapNode {
	if n.edit == edit {
		n.array[i] = a
		n.array[j] = b
		return n
	}
	return n.cloneAndSet2(i, a, j, b)
}

// ensureEditable returns n if already owned by edit, otherwise a mutable clone.
func (n *hmapBitmapNode) ensureEditable(edit *atomic.Bool) *hmapBitmapNode {
	if n.edit == edit {
		return n
	}
	nEntries := bits.OnesCount32(n.bitmap)
	// Allocate with room to grow (avoid frequent realloc during transient batch ops)
	newLen := max(2*(nEntries+1), len(n.array))
	newArray := make([]any, newLen)
	copy(newArray, n.array)
	return &hmapBitmapNode{bitmap: n.bitmap, array: newArray[:len(n.array)], edit: edit}
}

// assocTransient is like assoc but mutates in place when the node is owned.
func (n *hmapBitmapNode) assocTransient(edit *atomic.Bool, shift uint, hash uint32, key Value, val Value, addedLeaf *bool) *hmapBitmapNode {
	bit := hmapBitpos(hash, shift)
	idx := hmapIndex(n.bitmap, bit)

	if n.bitmap&bit != 0 {
		keyOrNil := n.array[2*idx]
		valOrNode := n.array[2*idx+1]

		if keyOrNil == nil {
			child := valOrNode.(hmapNode)
			var newChild hmapNode
			if bn, ok := child.(*hmapBitmapNode); ok {
				newChild = bn.assocTransient(edit, shift+hmapShift, hash, key, val, addedLeaf)
			} else {
				newChild = child.assoc(shift+hmapShift, hash, key, val, addedLeaf)
			}
			if newChild == child {
				return n
			}
			return n.editAndSet(edit, 2*idx+1, newChild)
		}

		existingKey := keyOrNil.(Value)
		if valueEquiv(key, existingKey) {
			existingVal := valOrNode.(Value)
			if valueEquiv(existingVal, val) {
				return n
			}
			return n.editAndSet(edit, 2*idx+1, val)
		}

		*addedLeaf = true
		existingHash := hashValue(existingKey)
		editable := n.ensureEditable(edit)
		editable.array[2*idx] = nil
		editable.array[2*idx+1] = createNode(shift+hmapShift, existingHash, existingKey, valOrNode.(Value), hash, key, val)
		return editable
	}

	// New entry — expand
	*addedLeaf = true
	nEntries := bits.OnesCount32(n.bitmap)
	editable := n.ensureEditable(edit)
	// Need to grow the array
	newArray := make([]any, 2*(nEntries+1))
	copy(newArray, editable.array[:2*idx])
	newArray[2*idx] = key
	newArray[2*idx+1] = val
	copy(newArray[2*(idx+1):], editable.array[2*idx:2*nEntries])
	editable.array = newArray
	editable.bitmap = n.bitmap | bit
	return editable
}

func (n *hmapBitmapNode) removePair(idx int) *hmapBitmapNode {
	newArray := make([]any, len(n.array)-2)
	copy(newArray, n.array[:2*idx])
	copy(newArray[2*idx:], n.array[2*(idx+1):])
	bit := uint32(1) << uint(n.indexToBit(idx))
	_ = bit
	// Recompute the bit to clear from the bitmap
	return &hmapBitmapNode{
		bitmap: n.bitmap & ^n.bitAtIndex(idx),
		array:  newArray,
	}
}

// bitAtIndex returns the bit in the bitmap corresponding to array index idx.
func (n *hmapBitmapNode) bitAtIndex(idx int) uint32 {
	// Walk the bitmap to find the idx-th set bit
	b := n.bitmap
	for range idx {
		b &= b - 1 // Clear lowest set bit
	}
	return b & (^b + 1) // Isolate lowest set bit
}

func (n *hmapBitmapNode) indexToBit(idx int) int {
	b := n.bitmap
	for range idx {
		b &= b - 1
	}
	return bits.TrailingZeros32(b)
}

// createNode creates a sub-node containing two key-value pairs.
func createNode(shift uint, hash1 uint32, key1 Value, val1 Value, hash2 uint32, key2 Value, val2 Value) hmapNode {
	if hash1 == hash2 {
		// Full hash collision — use collision node
		return &hmapCollisionNode{
			hash:  hash1,
			count: 2,
			array: []any{key1, val1, key2, val2},
		}
	}
	addedLeaf := false
	n1 := (&hmapBitmapNode{}).assoc(shift, hash1, key1, val1, &addedLeaf)
	return n1.assoc(shift, hash2, key2, val2, &addedLeaf)
}

// --- Collision Node ---

type hmapCollisionNode struct {
	hash  uint32
	count int
	array []any // pairs: [key0, val0, key1, val1, ...]
}

func (n *hmapCollisionNode) find(shift uint, hash uint32, key Value) (Value, bool) {
	idx := n.findIndex(key)
	if idx < 0 {
		return nil, false
	}
	return n.array[idx+1].(Value), true
}

func (n *hmapCollisionNode) findKeyword(shift uint, hash uint32, k Keyword) (Value, bool) {
	// same linear scan as find(), comparing keys typed:
	for i := 0; i < len(n.array); i += 2 {
		if sk, ok := n.array[i].(Keyword); ok && sk == k {
			return n.array[i+1].(Value), true
		}
	}
	return nil, false
}

func (n *hmapCollisionNode) assoc(shift uint, hash uint32, key Value, val Value, addedLeaf *bool) hmapNode {
	if hash == n.hash {
		// Same hash bucket
		idx := n.findIndex(key)
		if idx >= 0 {
			// Key exists — replace value
			if valueEquiv(n.array[idx+1].(Value), val) {
				return n
			}
			newArray := make([]any, len(n.array))
			copy(newArray, n.array)
			newArray[idx+1] = val
			return &hmapCollisionNode{hash: n.hash, count: n.count, array: newArray}
		}
		// New key in collision bucket
		*addedLeaf = true
		newArray := make([]any, len(n.array)+2)
		copy(newArray, n.array)
		newArray[len(n.array)] = key
		newArray[len(n.array)+1] = val
		return &hmapCollisionNode{hash: n.hash, count: n.count + 1, array: newArray}
	}
	// Different hash — wrap in a bitmap node and insert both
	*addedLeaf = true
	newNode := &hmapBitmapNode{
		bitmap: hmapBitpos(n.hash, shift),
		array:  []any{nil, n},
	}
	return newNode.assoc(shift, hash, key, val, addedLeaf)
}

func (n *hmapCollisionNode) dissoc(shift uint, hash uint32, key Value) hmapNode {
	idx := n.findIndex(key)
	if idx < 0 {
		return n // Key not present
	}
	if n.count == 1 {
		return nil // Last entry removed
	}
	if n.count == 2 {
		// Promote remaining entry to a bitmap node
		remainIdx := 0
		if idx == 0 {
			remainIdx = 2
		}
		remainKey := n.array[remainIdx].(Value)
		remainVal := n.array[remainIdx+1].(Value)
		addedLeaf := false
		return (&hmapBitmapNode{}).assoc(0, n.hash, remainKey, remainVal, &addedLeaf)
	}
	// Remove the pair at idx
	newArray := make([]any, len(n.array)-2)
	copy(newArray, n.array[:idx])
	copy(newArray[idx:], n.array[idx+2:])
	return &hmapCollisionNode{hash: n.hash, count: n.count - 1, array: newArray}
}

func (n *hmapCollisionNode) nodeSeq() []MapEntry {
	entries := make([]MapEntry, n.count)
	for i := 0; i < n.count; i++ {
		entries[i] = MapEntry{Key: n.array[2*i].(Value), Value: n.array[2*i+1].(Value)}
	}
	return entries
}

func (n *hmapCollisionNode) each(fn func(key, val Value) bool) bool {
	for i := 0; i < n.count; i++ {
		if !fn(n.array[2*i].(Value), n.array[2*i+1].(Value)) {
			return false
		}
	}
	return true
}

func (n *hmapCollisionNode) findIndex(key Value) int {
	for i := 0; i < n.count; i++ {
		if valueEquiv(key, n.array[2*i].(Value)) {
			return 2 * i
		}
	}
	return -1
}

// --- PersistentMap ---

// PersistentMap is an immutable hash-array mapped trie (HAMT).
type PersistentMap struct {
	count    int
	root     hmapNode
	meta     Value
	_hash    uint32
	_hasHash bool
}

// Meta implements IMeta.
func (m *PersistentMap) Meta() Value {
	if m.meta == nil {
		return NIL
	}
	return m.meta
}

// WithMeta implements IMeta.
func (m *PersistentMap) WithMeta(meta Value) Value {
	cp := *m
	cp.meta = meta
	cp._hash = 0
	cp._hasHash = false
	return &cp
}

// EmptyPersistentMap is the canonical empty persistent map.
var EmptyPersistentMap = &PersistentMap{count: 0, root: nil}

// NewPersistentMap creates a PersistentMap from an alternating key-value slice.
func NewPersistentMap(kvs []Value) *PersistentMap {
	if len(kvs) == 0 {
		return EmptyPersistentMap
	}
	if len(kvs)%2 != 0 {
		return EmptyPersistentMap
	}
	m := EmptyPersistentMap
	for i := 0; i < len(kvs); i += 2 {
		m = m.Assoc(kvs[i], kvs[i+1]).(*PersistentMap)
	}
	return m
}

// NewArrayMap creates a PersistentMap from an alternating key-value slice.
// Orderless by design: Clojure map order is an emergent array-map
// implementation detail (dropped past 8 entries), not a guarantee — let-go
// treats maps as unordered (cf. Go/Abseil map-iteration randomization).
// Order-dependent suite tests get :lg reader-cond branches.
func NewArrayMap(kvs []Value) *PersistentMap {
	if len(kvs) == 0 {
		return EmptyPersistentMap
	}
	if len(kvs)%2 != 0 {
		return EmptyPersistentMap
	}
	// Build on a transient to avoid a HAMT-node clone per pair (every map
	// literal flows through here).
	//
	// The Assoc/Persistent errors are provably unreachable: both fail only via
	// TransientMap.ensureEditable(), which errors on use-after-persistent or
	// cross-goroutine use. This transient is freshly created, single-owner,
	// mutated only in this loop, and persisted exactly once at the end — so
	// neither condition can hold. A non-nil error therefore means that
	// invariant was broken; fail loud rather than return a half-built map.
	t := NewTransientMap(EmptyPersistentMap)
	for i := 0; i < len(kvs); i += 2 {
		var err error
		if t, err = t.Assoc(kvs[i], kvs[i+1]); err != nil {
			panic("NewArrayMap: transient Assoc failed: " + err.Error())
		}
	}
	m, err := t.Persistent()
	if err != nil {
		panic("NewArrayMap: transient Persistent failed: " + err.Error())
	}
	return m
}

// --- Value interface ---

// Hash implements Hashable. Cached after first computation.
// Uses unordered hashing (order-independent) since maps are unordered.
func (m *PersistentMap) Hash() uint32 {
	if m._hasHash {
		return m._hash
	}
	// Hash each key-value pair and combine order-independently
	var h uint32
	s := m.Seq()
	for s != nil && s != EmptyList {
		entry := s.First()
		var key, val Value
		key, val, ok := MapEntryKV(entry)
		if !ok {
			s = s.Next()
			continue
		}
		// Combine key and value hashes into a pair hash, then XOR+add into total
		h += hashValue(key) ^ hashValue(val)
		s = s.Next()
	}
	m._hash = mixFinish(h)
	m._hasHash = true
	return m._hash
}

func (m *PersistentMap) Type() ValueType { return MapType }
func (m *PersistentMap) Unbox() any      { return m }

func (m *PersistentMap) String() string {
	b := &strings.Builder{}
	b.WriteRune('{')
	entries := m.entries()
	for i, e := range entries {
		entry := e.(MapEntry)
		b.WriteString(entry.Key.String())
		b.WriteRune(' ')
		b.WriteString(entry.Value.String())
		if i < len(entries)-1 {
			b.WriteRune(' ')
		}
	}
	b.WriteRune('}')
	return b.String()
}

// --- Collection interface ---

func (m *PersistentMap) Count() Value  { return Int(m.count) }
func (m *PersistentMap) RawCount() int { return m.count }

func (m *PersistentMap) Empty() Collection {
	if m.meta != nil {
		return EmptyPersistentMap.WithMeta(m.meta).(*PersistentMap)
	}
	return EmptyPersistentMap
}

func (m *PersistentMap) Conj(value Value) Collection {
	// Map + Map: merge entries
	if om, ok := value.(*PersistentMap); ok {
		result := m
		entries := om.entries()
		for _, e := range entries {
			switch av := e.(type) {
			case ArrayVector:
				if len(av) != 2 {
					continue
				}
				result = result.Assoc(av[0], av[1]).(*PersistentMap)
			case MapEntry:
				result = result.Assoc(av.Key, av.Value).(*PersistentMap)
			}
		}
		return result
	}
	if av, ok := value.(ArrayVector); ok && len(av) == 2 {
		return m.Assoc(av[0], av[1]).(*PersistentMap)
	}
	// Handle PersistentVector or any 2-element collection with Lookup
	if l, ok := value.(Lookup); ok {
		if c, ok := value.(Counted); ok && c.RawCount() == 2 {
			return m.Assoc(l.ValueAt(Int(0)), l.ValueAt(Int(1))).(*PersistentMap)
		}
	}
	return m
}

// --- Associative interface ---

func (m *PersistentMap) Assoc(key Value, val Value) Associative {
	if key == nil {
		return m
	}
	hash := hashValue(key)
	addedLeaf := false
	var newRoot hmapNode
	if m.root == nil {
		newRoot = (&hmapBitmapNode{}).assoc(0, hash, key, val, &addedLeaf)
	} else {
		newRoot = m.root.assoc(0, hash, key, val, &addedLeaf)
	}
	if newRoot == m.root {
		return m
	}
	newCount := m.count
	if addedLeaf {
		newCount++
	}
	return &PersistentMap{count: newCount, root: newRoot, meta: m.meta}
}

func (m *PersistentMap) Dissoc(key Value) Associative {
	if key == nil || m.root == nil {
		return m
	}
	hash := hashValue(key)
	newRoot := m.root.dissoc(0, hash, key)
	if newRoot == m.root {
		return m
	}
	if newRoot == nil {
		return &PersistentMap{count: 0, root: nil, meta: m.meta}
	}
	return &PersistentMap{count: m.count - 1, root: newRoot, meta: m.meta}
}

// --- Lookup interface ---

func (m *PersistentMap) ValueAt(key Value) Value {
	return m.ValueAtOr(key, NIL)
}

func (m *PersistentMap) ValueAtOr(key Value, dflt Value) Value {
	if key == nil || m.root == nil {
		return dflt
	}
	hash := hashValue(key)
	val, ok := m.root.find(0, hash, key)
	if !ok {
		return dflt
	}
	return val
}

func (m *PersistentMap) ValueAtKeyword(k Keyword) Value {
	return m.ValueAtKeywordOr(k, NIL)
}

func (m *PersistentMap) ValueAtKeywordOr(k Keyword, dflt Value) Value {
	if m.root == nil {
		return dflt
	}
	val, ok := m.root.findKeyword(0, k.Hash(), k)
	if !ok {
		return dflt
	}
	return val
}

// --- Keyed interface ---

func (m *PersistentMap) Contains(key Value) Boolean {
	if key == nil || m.root == nil {
		return FALSE
	}
	hash := hashValue(key)
	_, ok := m.root.find(0, hash, key)
	if ok {
		return TRUE
	}
	return FALSE
}

// --- Fn interface ---

func (m *PersistentMap) Arity() int { return -1 }

func (m *PersistentMap) Invoke(pargs []Value) (Value, error) {
	vl := len(pargs)
	if vl < 1 || vl > 2 {
		return NIL, fmt.Errorf("wrong number of arguments %d", vl)
	}
	if vl == 1 {
		return m.ValueAt(pargs[0]), nil
	}
	return m.ValueAtOr(pargs[0], pargs[1]), nil
}

// --- Sequable interface ---

func (m *PersistentMap) Seq() Seq {
	if m.count == 0 {
		return EmptyList
	}
	entries := m.entries()
	return &MapSeq{entries: entries, i: 0}
}

func (m *PersistentMap) entries() []Value {
	if m.root == nil {
		return nil
	}
	mes := m.root.nodeSeq()
	if allocAttrEnabled {
		recordAllocAttr(akMapEntries, len(mes)*16+24)
	}
	result := make([]Value, len(mes))
	for i, e := range mes {
		result[i] = e
	}
	return result
}

// --- Equals ---

func (m *PersistentMap) Equals(other Value) bool {
	o, ok := other.(*PersistentMap)
	if !ok {
		return false
	}
	if m.count != o.count {
		return false
	}
	// count equality + the invariant (count == 0 ⟺ root == nil) means both
	// roots are non-nil past this point.
	if m.count == 0 {
		return true
	}
	// Walk m in place (no entry-slice materialization — the old nodeSeq()
	// here dominated lowering allocation): every key in m must exist in o
	// with an equiv value. Sizes already match, so this is sufficient.
	return m.root.each(func(k, v Value) bool {
		ov, found := o.root.find(0, hashValue(k), k)
		return found && valueEquiv(v, ov)
	})
}
