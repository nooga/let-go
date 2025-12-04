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
	array []interface{} // can hold either Values or other nodes
}

func newNode() *vnode {
	return &vnode{array: make([]interface{}, 0, nodeCap)}
}

type thePersistentVectorType struct{}

func (t *thePersistentVectorType) String() string     { return t.Name() }
func (t *thePersistentVectorType) Type() ValueType    { return TypeType }
func (t *thePersistentVectorType) Unbox() interface{} { return reflect.TypeOf(t) }
func (t *thePersistentVectorType) Name() string       { return "let-go.lang.PersistentVector" }

func (t *thePersistentVectorType) Box(bare interface{}) (Value, error) {
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
	count   int
	shift   uint
	root    *vnode
	tail    []Value // Last node is stored separately for efficiency
	tailOff int
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
func (v PersistentVector) Unbox() interface{} {
	ret := make([]Value, v.count)
	for i := 0; i < v.count; i++ {
		ret[i] = v.ValueAt(Int(i))
	}
	return ret
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
	if s.inTail {
		return s.vec.tail[s.nodeIdx]
	}
	return s.node.array[s.nodeIdx].(Value)
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
		return NIL
	}
	return s.nextSeq()
}

func (s *PersistentVectorSeq) nextSeq() *PersistentVectorSeq {
	nextI := s.i + 1

	// Check if we need to move to tail
	if nextI >= s.vec.tailOff && !s.inTail {
		return &PersistentVectorSeq{
			vec:     s.vec,
			i:       nextI,
			nodeIdx: nextI - s.vec.tailOff,
			inTail:  true,
		}
	}

	// If we're already in tail, just increment index
	if s.inTail {
		return &PersistentVectorSeq{
			vec:     s.vec,
			i:       nextI,
			nodeIdx: s.nodeIdx + 1,
			inTail:  true,
		}
	}

	// If we're at the end of current node, find next node
	if s.nodeIdx+1 >= len(s.node.array) {
		newNode := s.findNextNode(nextI)
		return &PersistentVectorSeq{
			vec:     s.vec,
			i:       nextI,
			node:    newNode,
			nodeIdx: 0,
			inTail:  false,
		}
	}

	// Just move to next element in current node
	return &PersistentVectorSeq{
		vec:     s.vec,
		i:       nextI,
		node:    s.node,
		nodeIdx: s.nodeIdx + 1,
		inTail:  false,
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
		count:   0,
		shift:   shift,
		root:    newNode(),
		tail:    make([]Value, 0, nodeCap),
		tailOff: 0,
	}
}

// Conj implements Collection
func (v PersistentVector) Conj(val Value) Collection {
	// Special case for empty vector
	if v.count == 0 {
		return PersistentVector{
			count:   1,
			shift:   shift,
			root:    newNode(),
			tail:    []Value{val},
			tailOff: 0,
		}
	}

	// If tail is not full, just append to tail
	if len(v.tail) < nodeCap {
		newTail := make([]Value, len(v.tail)+1)
		copy(newTail, v.tail)
		newTail[len(v.tail)] = val
		return PersistentVector{
			count:   v.count + 1,
			shift:   v.shift,
			root:    v.root,
			tail:    newTail,
			tailOff: v.tailOff,
		}
	}

	// Need to push tail into tree and create new tail
	newTail := []Value{val}
	newRoot := v.root
	newShift := v.shift

	// Check if we need to grow the tree height
	if (v.count-v.tailOff)>>v.shift > (1 << v.shift) {
		newRoot = newNode()
		newRoot.array = append(newRoot.array, v.root)
		newShift += shift
	}

	// Push current tail into tree
	newRoot = v.pushTail(newRoot, newShift, v.count-1, v.tail)

	return PersistentVector{
		count:   v.count + 1,
		shift:   newShift,
		root:    newRoot,
		tail:    newTail,
		tailOff: v.count,
	}
}

func (v PersistentVector) pushTail(node *vnode, level uint, index int, tail []Value) *vnode {
	anewNode := &vnode{array: make([]interface{}, len(node.array))}
	copy(anewNode.array, node.array)

	if level == shift {
		// At leaf level, store Values individually
		for _, val := range tail {
			anewNode.array = append(anewNode.array, val)
		}
	} else {
		subidx := (index >> level) & nodeMask
		// Create new path if needed
		if len(node.array) <= int(subidx) {
			child := newNode()
			anewNode.array = append(anewNode.array, child)
			anewNode.array[len(anewNode.array)-1] = v.pushTail(child, level-shift, index, tail)
		} else {
			// Update existing path
			child := node.array[subidx].(*vnode)
			anewNode.array[subidx] = v.pushTail(child, level-shift, index, tail)
		}
	}

	return anewNode
}

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
	if !ok || idx < 0 || int(idx) >= v.count {
		return NIL
	}

	if int(idx) >= v.tailOff {
		newTail := make([]Value, len(v.tail))
		copy(newTail, v.tail)
		newTail[int(idx)-v.tailOff] = val
		return PersistentVector{
			count:   v.count,
			shift:   v.shift,
			root:    v.root,
			tail:    newTail,
			tailOff: v.tailOff,
		}
	}

	return PersistentVector{
		count:   v.count,
		shift:   v.shift,
		root:    v.doAssoc(v.root, v.shift, int(idx), val),
		tail:    v.tail,
		tailOff: v.tailOff,
	}
}

func (v PersistentVector) doAssoc(n *vnode, level uint, idx int, val Value) *vnode {
	newNode := &vnode{array: make([]interface{}, len(n.array))}
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
	return v.ValueAt(args[0]), nil
}

func NewPersistentVector(values []Value) Value {
	if len(values) == 0 {
		return PersistentVector{
			count:   0,
			shift:   shift,
			root:    newNode(),
			tail:    make([]Value, 0, nodeCap),
			tailOff: 0,
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
			count:   len(values),
			shift:   shift,
			root:    newNode(),
			tail:    tail,
			tailOff: 0,
		}
	}

	// Build tree levels bottom-up
	root := newNode()
	currentLevel := []*vnode{root}
	nodesInLevel := (treeSize + nodeCap - 1) / nodeCap

	// Fill leaf level
	for i := 0; i < treeSize; i += nodeCap {
		node := newNode()
		end := i + nodeCap
		if end > treeSize {
			end = treeSize
		}
		for j := i; j < end; j++ {
			node.array = append(node.array, values[j])
		}
		root.array = append(root.array, node)
	}

	// Calculate required height
	shift := uint(5)
	for nodesInLevel > nodeCap {
		newLevel := newNode()
		for i := 0; i < len(currentLevel); i += nodeCap {
			node := newNode()
			end := i + nodeCap
			if end > len(currentLevel) {
				end = len(currentLevel)
			}
			node.array = make([]interface{}, end-i)
			for j := i; j < end; j++ {
				node.array[j-i] = currentLevel[j]
			}
			newLevel.array = append(newLevel.array, node)
		}
		currentLevel = []*vnode{newLevel}
		nodesInLevel = (nodesInLevel + nodeCap - 1) / nodeCap
		shift += 5
	}

	return PersistentVector{
		count:   len(values),
		shift:   shift,
		root:    currentLevel[0],
		tail:    tail,
		tailOff: treeSize,
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
func (s *PersistentVectorSeq) Unbox() interface{} {
	vals := make([]Value, s.vec.count-s.i)
	for i := 0; i < len(vals); i++ {
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
	// Use EmptyList and build the collection
	result := EmptyList.Cons(val).(Collection)

	// Add all remaining elements from the sequence
	for i := s.i; i < s.vec.count; i++ {
		result = result.Conj(s.vec.ValueAt(Int(i)))
	}

	return result
}

type theSequenceType struct{}

func (t *theSequenceType) String() string     { return t.Name() }
func (t *theSequenceType) Type() ValueType    { return TypeType }
func (t *theSequenceType) Unbox() interface{} { return reflect.TypeOf(t) }
func (t *theSequenceType) Name() string       { return "let-go.lang.Sequence" }

func (t *theSequenceType) Box(bare interface{}) (Value, error) {
	return nil, NewTypeError(bare, "can't be boxed as", t)
}

// SequenceType is the type of Sequences
var SequenceType *theSequenceType = &theSequenceType{}
