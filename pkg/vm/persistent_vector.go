/*
 * Copyright (c) 2024 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	shift    = 5
	nodeCap  = 32 // 1 << shift
	nodeMask = 31 // nodeCap - 1
)

type vnode struct {
	array []any // can hold either Values or other nodes
}

func newNode() *vnode {
	return &vnode{array: make([]any, 0, nodeCap)}
}

type thePersistentVectorType struct{}

func (t *thePersistentVectorType) String() string  { return t.Name() }
func (t *thePersistentVectorType) Type() ValueType { return TypeType }
func (t *thePersistentVectorType) Unbox() any      { return reflect.TypeFor[*thePersistentVectorType]() }
func (t *thePersistentVectorType) Name() string    { return "let-go.lang.PersistentVector" }

func (t *thePersistentVectorType) Box(bare any) (Value, error) {
	arr, ok := bare.([]Value)
	if !ok {
		return NIL, NewTypeError(bare, "can't be boxed as", t)
	}
	return NewPersistentVector(arr), nil
}

// PersistentVectorType is the type of PersistentVectors
var PersistentVectorType *thePersistentVectorType = &thePersistentVectorType{}

// PersistentVector is a persistent vector implementation using a bit-partitioned trie
type PersistentVector struct {
	count    int
	shift    uint
	root     *vnode
	tail     []Value // Last node is stored separately for efficiency
	tailOff  int
	meta     Value
	identity *persistentVectorIdentity
}

// persistentVectorIdentity is deliberately non-zero-sized so separately
// allocated tokens cannot share an address. PersistentVector copies retain the
// pointer; every operation that creates a logically new vector replaces it.
type persistentVectorIdentity struct {
	_ byte
}

func newPersistentVectorIdentity() *persistentVectorIdentity {
	return &persistentVectorIdentity{}
}

// IdentityKey returns an opaque key for the vector's private identity token.
// A zero-value PersistentVector has no token and returns zero.
func (v PersistentVector) IdentityKey() uintptr {
	if v.identity == nil {
		return 0
	}
	return reflect.ValueOf(v.identity).Pointer()
}

// SameIdentity reports whether two PersistentVector values are copies of the
// same logical vector. The token itself remains private to the vm package.
func (v PersistentVector) SameIdentity(other PersistentVector) bool {
	return v.identity != nil && v.identity == other.identity
}

// Hash implements Hashable. Computed from elements.
func (v PersistentVector) Hash() uint32 {
	return hashOrdered(v.Seq())
}

// Meta implements IMeta.
func (v PersistentVector) Meta() Value {
	if v.meta == nil {
		return NIL
	}
	return v.meta
}

// WithMeta implements IMeta.
func (v PersistentVector) WithMeta(m Value) Value {
	v.meta = m
	v.identity = newPersistentVectorIdentity()
	return v
}

// Type implements Value
func (v PersistentVector) Type() ValueType { return PersistentVectorType }

// String implements Value
func (v PersistentVector) String() string {
	var b strings.Builder
	b.WriteRune('[')
	for i := 0; i < v.count; i++ {
		if i > 0 {
			b.WriteRune(' ')
		}
		b.WriteString(v.ValueAt(Int(i)).String())
	}
	b.WriteRune(']')
	return b.String()
}

// Unbox implements Value
func (v PersistentVector) Unbox() any {
	ret := make([]Value, v.count)
	for i := 0; i < v.count; i++ {
		ret[i] = v.ValueAt(Int(i))
	}
	return ret
}

// Equals implements value equality for PersistentVector
func (v PersistentVector) Equals(other Value) bool {
	switch o := other.(type) {
	case MapEntry:
		return v.Equals(ArrayVector{o.Key, o.Value})
	case PersistentVector:
		if v.count != o.count {
			return false
		}
		for i := 0; i < v.count; i++ {
			vv := v.ValueAt(Int(i))
			ov := o.ValueAt(Int(i))
			if nilListEquivalent(vv, ov) {
				continue
			}
			if ValueEquals != nil && ValueEquals(vv, ov) {
				continue
			}
			if eq, ok := vv.(interface{ Equals(Value) bool }); ok {
				if !eq.Equals(ov) {
					return false
				}
			} else if vv != ov {
				return false
			}
		}
		return true
	case ArrayVector:
		if v.count != len(o) {
			return false
		}
		for i := 0; i < v.count; i++ {
			vv := v.ValueAt(Int(i))
			if nilListEquivalent(vv, o[i]) {
				continue
			}
			if ValueEquals != nil && ValueEquals(vv, o[i]) {
				continue
			}
			if eq, ok := vv.(interface{ Equals(Value) bool }); ok {
				if !eq.Equals(o[i]) {
					return false
				}
			} else if vv != o[i] {
				return false
			}
		}
		return true
	default:
		return false
	}
}

// PersistentVectorSeq represents a sequence view of a PersistentVector
type PersistentVectorSeq struct {
	vec     *PersistentVector
	i       int    // Overall index
	node    *vnode // Current node in the tree
	nodeIdx int    // Index within current node
	inTail  bool   // Whether we're in the tail section
}

// First implements Seq
func (s *PersistentVectorSeq) First() Value {
	if s.i >= s.vec.count {
		return NIL
	}
	return s.vec.ValueAt(Int(s.i))
}

// More implements Seq
func (s *PersistentVectorSeq) More() Seq {
	if s.i+1 >= s.vec.count {
		return EmptyList
	}
	return s.nextSeq()
}

// Next implements Seq
func (s *PersistentVectorSeq) Next() Seq {
	if s.i+1 >= s.vec.count {
		return nil
	}
	return s.nextSeq()
}

func (s *PersistentVectorSeq) nextSeq() *PersistentVectorSeq {
	return &PersistentVectorSeq{
		vec: s.vec,
		i:   s.i + 1,
	}
}

func (s *PersistentVectorSeq) findNextNode(index int) *vnode {
	node := s.vec.root
	level := s.vec.shift
	for level > 0 {
		node = node.array[(index>>level)&nodeMask].(*vnode)
		level -= shift
	}
	return node
}

// Cons implements Seq
func (s *PersistentVectorSeq) Cons(val Value) Seq {
	return NewCons(val, s)
}

// Nth implements Indexed: positional access by integer index.
func (s *PersistentVectorSeq) Nth(i int) Value { return s.ValueAt(Int(i)) }

// ValueAt implements Lookup for PersistentVectorSeq so that `get` works on seq views.
func (s *PersistentVectorSeq) ValueAt(key Value) Value {
	return s.ValueAtOr(key, NIL)
}

// ValueAtOr implements Lookup for PersistentVectorSeq.
func (s *PersistentVectorSeq) ValueAtOr(key Value, dflt Value) Value {
	if key == NIL {
		return dflt
	}
	idx, ok := key.(Int)
	if !ok || idx < 0 {
		return dflt
	}
	absIdx := s.i + int(idx)
	if absIdx >= s.vec.count {
		return dflt
	}
	return s.vec.ValueAtOr(Int(absIdx), dflt)
}

// Modify PersistentVector's Seq method:
func (v PersistentVector) Seq() Seq {
	if v.count == 0 {
		return EmptyList
	}

	// If vector only has tail elements
	if v.count <= nodeCap {
		return &PersistentVectorSeq{
			vec:     &v,
			i:       0,
			nodeIdx: 0,
			inTail:  true,
		}
	}

	// Start with root node
	return &PersistentVectorSeq{
		vec:     &v,
		i:       0,
		node:    v.getFirstLeafNode(),
		nodeIdx: 0,
		inTail:  false,
	}
}

// Helper method to get the first leaf node
func (v PersistentVector) getFirstLeafNode() *vnode {
	node := v.root
	level := v.shift
	for level > 0 {
		node = node.array[0].(*vnode)
		level -= shift
	}
	return node
}

// Count implements Collection
func (v PersistentVector) Count() Value {
	return Int(v.count)
}

func (v PersistentVector) RawCount() int {
	return v.count
}

// Empty implements Collection
func (v PersistentVector) Empty() Collection {
	return PersistentVector{
		count:    0,
		shift:    shift,
		root:     newNode(),
		tail:     make([]Value, 0, nodeCap),
		tailOff:  0,
		meta:     v.meta,
		identity: newPersistentVectorIdentity(),
	}
}

// Conj implements Collection
func (v PersistentVector) Conj(val Value) Collection {
	// Special case for empty vector
	if v.count == 0 {
		return PersistentVector{
			count:    1,
			shift:    shift,
			root:     newNode(),
			tail:     []Value{val},
			tailOff:  0,
			meta:     v.meta,
			identity: newPersistentVectorIdentity(),
		}
	}

	// If tail is not full, just append to tail
	if len(v.tail) < nodeCap {
		newTail := make([]Value, len(v.tail)+1)
		copy(newTail, v.tail)
		newTail[len(v.tail)] = val
		return PersistentVector{
			count:    v.count + 1,
			shift:    v.shift,
			root:     v.root,
			tail:     newTail,
			tailOff:  v.tailOff,
			meta:     v.meta,
			identity: newPersistentVectorIdentity(),
		}
	}

	// Need to push tail into tree and create new tail
	newTail := []Value{val}
	newShift := v.shift

	// Check if we need to grow the tree height.
	// Tree overflow: tailOff indexes need more bits than the current shift allows.
	if (v.tailOff >> shift) >= (1 << v.shift) {
		newRoot := newNode()
		newRoot.array = append(newRoot.array, v.root)
		newRoot.array = append(newRoot.array, newPath(v.shift, v.tail))
		newShift += shift
		return PersistentVector{
			count:    v.count + 1,
			shift:    newShift,
			root:     newRoot,
			tail:     newTail,
			tailOff:  v.count,
			meta:     v.meta,
			identity: newPersistentVectorIdentity(),
		}
	}

	// Push current tail into tree as a new leaf node
	newRoot := pushTail(v.shift, v.root, v.tailOff, v.tail)

	return PersistentVector{
		count:    v.count + 1,
		shift:    newShift,
		root:     newRoot,
		tail:     newTail,
		tailOff:  v.count,
		meta:     v.meta,
		identity: newPersistentVectorIdentity(),
	}
}

// pushTail inserts a tail chunk into the trie at the position indicated by tailOff.
func pushTail(level uint, parent *vnode, tailOff int, tail []Value) *vnode {
	subidx := (tailOff >> level) & nodeMask

	// Copy parent
	ret := &vnode{array: make([]any, len(parent.array))}
	copy(ret.array, parent.array)

	if level == shift {
		// At the lowest internal level — append a new leaf node
		leafNode := newNode()
		for _, val := range tail {
			leafNode.array = append(leafNode.array, val)
		}
		ret.array = append(ret.array, leafNode)
		return ret
	}

	// Recurse deeper
	if subidx < len(parent.array) {
		// Subtree exists — recurse into it
		child := parent.array[subidx].(*vnode)
		ret.array[subidx] = pushTail(level-shift, child, tailOff, tail)
	} else {
		// No subtree at this position — create a new path
		ret.array = append(ret.array, newPath(level-shift, tail))
	}

	return ret
}

// newPath creates a chain of single-child nodes from level down to a leaf containing tail.
func newPath(level uint, tail []Value) *vnode {
	if level == 0 {
		leafNode := newNode()
		for _, val := range tail {
			leafNode.array = append(leafNode.array, val)
		}
		return leafNode
	}
	ret := newNode()
	ret.array = append(ret.array, newPath(level-shift, tail))
	return ret
}

// Pop returns a new vector with the last element removed, in O(1) amortized
// time (mirrors Clojure's PersistentVector.pop). The previous behavior — the
// `pop` builtin doing Unbox()+NewPersistentVector(all-but-last) — was O(n) per
// call, which made vector-as-stack uses (e.g. the typeinfer worklist drain via
// peek/pop) O(n^2) and a major allocation source. Returns the empty vector for
// count <= 1; callers that must reject popping an empty vector check count
// first.
func (v PersistentVector) Pop() PersistentVector {
	if v.count <= 1 {
		return v.Empty().(PersistentVector)
	}
	// More than one element in the tail: just shrink the tail (O(tail) <= 32).
	if len(v.tail) > 1 {
		newTail := make([]Value, len(v.tail)-1)
		copy(newTail, v.tail[:len(v.tail)-1])
		return PersistentVector{
			count:    v.count - 1,
			shift:    v.shift,
			root:     v.root,
			tail:     newTail,
			tailOff:  v.tailOff,
			meta:     v.meta,
			identity: newPersistentVectorIdentity(),
		}
	}
	// The tail holds 0 or 1 elements, so the new tail is pulled from the
	// rightmost trie leaf. (NewPersistentVector leaves an empty tail when the
	// count is a multiple of nodeCap, so the empty case is real.) removeIdx is
	// the index of an element in that leaf: count-2 when the tail had the last
	// element, count-1 when the last element lives in the leaf itself.
	removeIdx := v.count - 1 - len(v.tail)
	newTail := v.leafArrayFor(removeIdx)
	if len(v.tail) == 0 {
		// The popped element was the last of this leaf — drop it from the tail.
		newTail = newTail[:len(newTail)-1]
	}
	newShift := v.shift
	newRoot := v.popTail(v.shift, v.root, removeIdx)
	if newRoot == nil {
		newRoot = newNode()
	}
	// Collapse a now-redundant root level (only one child left).
	if v.shift > shift && len(newRoot.array) == 1 {
		newRoot = newRoot.array[0].(*vnode)
		newShift -= shift
	}
	return PersistentVector{
		count:    v.count - 1,
		shift:    newShift,
		root:     newRoot,
		tail:     newTail,
		tailOff:  (v.count - 1) - len(newTail),
		meta:     v.meta,
		identity: newPersistentVectorIdentity(),
	}
}

// leafArrayFor returns a fresh []Value copy of the leaf array holding index i
// (i must be in the trie portion, i.e. i < tailOff).
func (v PersistentVector) leafArrayFor(i int) []Value {
	node := v.root
	for level := v.shift; level > 0; level -= shift {
		node = node.array[(i>>level)&nodeMask].(*vnode)
	}
	out := make([]Value, len(node.array))
	for j, x := range node.array {
		out[j] = x.(Value)
	}
	return out
}

// popTail removes the leaf containing idx from the trie rooted at node,
// returning the new subtree (nil when it becomes empty). Mirrors Clojure's
// popTail, adapted to the dynamically-sized node arrays used here: the
// rightmost child is dropped by truncation rather than nulled in a fixed
// 32-slot array.
func (v PersistentVector) popTail(level uint, node *vnode, idx int) *vnode {
	subidx := (idx >> level) & nodeMask
	if level > shift {
		newChild := v.popTail(level-shift, node.array[subidx].(*vnode), idx)
		if newChild == nil && subidx == 0 {
			return nil
		}
		ret := &vnode{array: make([]any, len(node.array))}
		copy(ret.array, node.array)
		if newChild == nil {
			ret.array = ret.array[:subidx]
		} else {
			ret.array[subidx] = newChild
		}
		return ret
	}
	if subidx == 0 {
		return nil
	}
	// Leaf-parent level: drop the rightmost leaf.
	ret := &vnode{array: make([]any, subidx)}
	copy(ret.array, node.array[:subidx])
	return ret
}

// Nth implements Indexed: positional access by integer index.
func (v PersistentVector) Nth(i int) Value { return v.ValueAt(Int(i)) }

// ValueAt implements Associative
func (v PersistentVector) ValueAt(key Value) Value {
	return v.ValueAtOr(key, NIL)
}

// ValueAtOr implements Associative
func (v PersistentVector) ValueAtOr(key Value, dflt Value) Value {
	if key == NIL {
		return dflt
	}
	idx, ok := key.(Int)
	if !ok || idx < 0 || int(idx) >= v.count {
		return dflt
	}

	if int(idx) >= v.tailOff {
		return v.tail[int(idx)-v.tailOff]
	}

	node := v.root
	level := v.shift
	for level > 0 {
		if len(node.array) == 0 {
			return dflt
		}
		subidx := (int(idx) >> level) & nodeMask
		if subidx >= len(node.array) {
			return dflt
		}
		nextNode, ok := node.array[subidx].(*vnode)
		if !ok {
			return dflt
		}
		node = nextNode
		level -= shift
	}

	finalIdx := int(idx) & nodeMask
	if finalIdx >= len(node.array) {
		return dflt
	}
	val, ok := node.array[finalIdx].(Value)
	if !ok {
		return dflt
	}
	return val
}

// Contains implements Associative
func (v PersistentVector) Contains(key Value) Boolean {
	idx, ok := key.(Int)
	if !ok || idx < 0 || int(idx) >= v.count {
		return FALSE
	}
	return TRUE
}

// Assoc implements Associative
func (v PersistentVector) Assoc(key Value, val Value) Associative {
	idx, ok := key.(Int)
	if !ok || idx < 0 || int(idx) > v.count {
		return NIL
	}
	if int(idx) == v.count {
		return v.Conj(val).(Associative)
	}

	if int(idx) >= v.tailOff {
		newTail := make([]Value, len(v.tail))
		copy(newTail, v.tail)
		newTail[int(idx)-v.tailOff] = val
		return PersistentVector{
			count:    v.count,
			shift:    v.shift,
			root:     v.root,
			tail:     newTail,
			tailOff:  v.tailOff,
			meta:     v.meta,
			identity: newPersistentVectorIdentity(),
		}
	}

	return PersistentVector{
		count:    v.count,
		shift:    v.shift,
		root:     v.doAssoc(v.root, v.shift, int(idx), val),
		tail:     v.tail,
		tailOff:  v.tailOff,
		meta:     v.meta,
		identity: newPersistentVectorIdentity(),
	}
}

func (v PersistentVector) doAssoc(n *vnode, level uint, idx int, val Value) *vnode {
	newNode := &vnode{array: make([]any, len(n.array))}
	copy(newNode.array, n.array)

	if level == 0 {
		newNode.array[idx&nodeMask] = val
	} else {
		subidx := (idx >> level) & nodeMask
		newNode.array[subidx] = v.doAssoc(n.array[subidx].(*vnode), level-shift, idx, val)
	}
	return newNode
}

// Dissoc implements Associative
func (v PersistentVector) Dissoc(key Value) Associative {
	return NIL // Vectors don't support removal
}

// Arity implements IFn
func (v PersistentVector) Arity() int {
	return 1
}

// Invoke implements IFn
func (v PersistentVector) Invoke(args []Value) (Value, error) {
	if len(args) != 1 {
		return NIL, fmt.Errorf("wrong number of arguments: %d, expected: 1", len(args))
	}
	idx, ok := args[0].(Int)
	if !ok {
		return NIL, fmt.Errorf("vector key must be Int")
	}
	i := int(idx)
	if i < 0 || i >= v.count {
		return NIL, fmt.Errorf("index out of bounds: %d", i)
	}
	return v.ValueAt(idx), nil
}

func NewPersistentVector(values []Value) Value {
	if len(values) == 0 {
		return PersistentVector{
			count:    0,
			shift:    shift,
			root:     newNode(),
			tail:     make([]Value, 0, nodeCap),
			tailOff:  0,
			identity: newPersistentVectorIdentity(),
		}
	}

	// Calculate how many elements will be in the tree vs tail
	tailSize := len(values) % nodeCap
	treeSize := len(values) - tailSize

	// Create tail
	tail := make([]Value, tailSize)
	copy(tail, values[treeSize:])

	// If we only have tail elements
	if treeSize == 0 {
		return PersistentVector{
			count:    len(values),
			shift:    shift,
			root:     newNode(),
			tail:     tail,
			tailOff:  0,
			identity: newPersistentVectorIdentity(),
		}
	}

	// Build leaf nodes
	leafCount := treeSize / nodeCap
	leaves := make([]*vnode, leafCount)
	for i := range leafCount {
		node := &vnode{array: make([]any, nodeCap)}
		base := i * nodeCap
		for j := range nodeCap {
			node.array[j] = values[base+j]
		}
		leaves[i] = node
	}

	// Build tree bottom-up: group nodes into parents until we have a single root
	currentLevel := leaves
	treeShift := uint(shift)
	for len(currentLevel) > nodeCap {
		parentCount := (len(currentLevel) + nodeCap - 1) / nodeCap
		parents := make([]*vnode, parentCount)
		for i := range parentCount {
			start := i * nodeCap
			end := min(start+nodeCap, len(currentLevel))
			node := &vnode{array: make([]any, end-start)}
			for j := start; j < end; j++ {
				node.array[j-start] = currentLevel[j]
			}
			parents[i] = node
		}
		currentLevel = parents
		treeShift += shift
	}

	// Final root node wrapping the top-level nodes
	root := &vnode{array: make([]any, len(currentLevel))}
	for i, n := range currentLevel {
		root.array[i] = n
	}

	return PersistentVector{
		count:    len(values),
		shift:    treeShift,
		root:     root,
		tail:     tail,
		tailOff:  treeSize,
		identity: newPersistentVectorIdentity(),
	}
}

// String implements Value
func (s *PersistentVectorSeq) String() string {
	return "(seq " + s.vec.String() + ")"
}

// Type implements Value
func (s *PersistentVectorSeq) Type() ValueType {
	return SequenceType
}

// Unbox implements Value
func (s *PersistentVectorSeq) Unbox() any {
	vals := make([]Value, s.vec.count-s.i)
	for i := range vals {
		vals[i] = s.vec.ValueAt(Int(s.i + i))
	}
	return vals
}

// Count implements Collection
func (s *PersistentVectorSeq) Count() Value {
	return Int(s.vec.count - s.i)
}

func (s *PersistentVectorSeq) RawCount() int {
	return s.vec.count - s.i
}

// Empty implements Collection
func (s *PersistentVectorSeq) Empty() Collection {
	return EmptyList
}

// Conj implements Collection
func (s *PersistentVectorSeq) Conj(val Value) Collection {
	values := make([]Value, s.vec.count-s.i+1)
	values[0] = val
	for i := s.i; i < s.vec.count; i++ {
		values[i-s.i+1] = s.vec.ValueAt(Int(i))
	}
	return NewList(values).(*List)
}

type theSequenceType struct{}

func (t *theSequenceType) String() string  { return t.Name() }
func (t *theSequenceType) Type() ValueType { return TypeType }
func (t *theSequenceType) Unbox() any      { return reflect.TypeFor[*theSequenceType]() }
func (t *theSequenceType) Name() string    { return "let-go.lang.Sequence" }

func (t *theSequenceType) Box(bare any) (Value, error) {
	return nil, NewTypeError(bare, "can't be boxed as", t)
}

// SequenceType is the type of Sequences
var SequenceType *theSequenceType = &theSequenceType{}
