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
)

const (
	hmapShift = 5
	hmapMask  = 0x1f // 31
)

// --- Type metadata ---

type thePersistentMapType struct{}

func (t *thePersistentMapType) String() string     { return t.Name() }
func (t *thePersistentMapType) Type() ValueType    { return TypeType }
func (t *thePersistentMapType) Unbox() interface{} { return reflect.TypeOf(t) }
func (t *thePersistentMapType) Name() string        { return "let-go.lang.PersistentHashMap" }

func (t *thePersistentMapType) Box(bare interface{}) (Value, error) {
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
	assoc(shift uint, hash uint32, key Value, val Value, addedLeaf *bool) hmapNode
	dissoc(shift uint, hash uint32, key Value) hmapNode
	nodeSeq() []MapEntry
}

// MapEntry is a key-value pair.
type MapEntry struct {
	Key   Value
	Value Value
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
	array  []interface{} // pairs: [key-or-nil, val-or-child, ...]
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
	newArray := make([]interface{}, 2*(nEntries+1))
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
	for i := 0; i < nSlots; i++ {
		keyOrNil := n.array[2*i]
		valOrNode := n.array[2*i+1]
		if keyOrNil != nil {
			entries = append(entries, MapEntry{Key: keyOrNil.(Value), Value: valOrNode.(Value)})
		} else if valOrNode != nil {
			entries = append(entries, valOrNode.(hmapNode).nodeSeq()...)
		}
	}
	return entries
}

func (n *hmapBitmapNode) cloneAndSet(i int, val interface{}) *hmapBitmapNode {
	newArray := make([]interface{}, len(n.array))
	copy(newArray, n.array)
	newArray[i] = val
	return &hmapBitmapNode{bitmap: n.bitmap, array: newArray}
}

func (n *hmapBitmapNode) cloneAndSet2(i int, a interface{}, j int, b interface{}) *hmapBitmapNode {
	newArray := make([]interface{}, len(n.array))
	copy(newArray, n.array)
	newArray[i] = a
	newArray[j] = b
	return &hmapBitmapNode{bitmap: n.bitmap, array: newArray}
}

func (n *hmapBitmapNode) removePair(idx int) *hmapBitmapNode {
	newArray := make([]interface{}, len(n.array)-2)
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
	for i := 0; i < idx; i++ {
		b &= b - 1 // Clear lowest set bit
	}
	return b & (^b + 1) // Isolate lowest set bit
}

func (n *hmapBitmapNode) indexToBit(idx int) int {
	b := n.bitmap
	for i := 0; i < idx; i++ {
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
			array: []interface{}{key1, val1, key2, val2},
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
	array []interface{} // pairs: [key0, val0, key1, val1, ...]
}

func (n *hmapCollisionNode) find(shift uint, hash uint32, key Value) (Value, bool) {
	idx := n.findIndex(key)
	if idx < 0 {
		return nil, false
	}
	return n.array[idx+1].(Value), true
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
			newArray := make([]interface{}, len(n.array))
			copy(newArray, n.array)
			newArray[idx+1] = val
			return &hmapCollisionNode{hash: n.hash, count: n.count, array: newArray}
		}
		// New key in collision bucket
		*addedLeaf = true
		newArray := make([]interface{}, len(n.array)+2)
		copy(newArray, n.array)
		newArray[len(n.array)] = key
		newArray[len(n.array)+1] = val
		return &hmapCollisionNode{hash: n.hash, count: n.count + 1, array: newArray}
	}
	// Different hash — wrap in a bitmap node and insert both
	*addedLeaf = true
	newNode := &hmapBitmapNode{
		bitmap: hmapBitpos(n.hash, shift),
		array:  []interface{}{nil, n},
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
	newArray := make([]interface{}, len(n.array)-2)
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
	count int
	root  hmapNode
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

// --- Value interface ---

func (m *PersistentMap) Type() ValueType    { return MapType }
func (m *PersistentMap) Unbox() interface{} { return m }

func (m *PersistentMap) String() string {
	b := &strings.Builder{}
	b.WriteRune('{')
	entries := m.entries()
	for i, e := range entries {
		entry := e.(ArrayVector)
		b.WriteString(entry[0].String())
		b.WriteRune(' ')
		b.WriteString(entry[1].String())
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
	return EmptyPersistentMap
}

func (m *PersistentMap) Conj(value Value) Collection {
	if value.Type() != ArrayVectorType {
		return m
	}
	v := value.(ArrayVector)
	if len(v) != 2 {
		return m
	}
	return m.Assoc(v[0], v[1]).(*PersistentMap)
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
	return &PersistentMap{count: newCount, root: newRoot}
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
		return EmptyPersistentMap
	}
	return &PersistentMap{count: m.count - 1, root: newRoot}
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
	result := make([]Value, len(mes))
	for i, e := range mes {
		result[i] = ArrayVector{e.Key, e.Value}
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
	if m.root == nil && o.root == nil {
		return true
	}
	// Check that every entry in m exists with the same value in o
	entries := m.root.nodeSeq()
	for _, e := range entries {
		v, found := o.root.find(0, hashValue(e.Key), e.Key)
		if !found || !valueEquiv(e.Value, v) {
			return false
		}
	}
	return true
}
