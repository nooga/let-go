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

type theArrayVectorType struct{}

func (t *theArrayVectorType) String() string  { return t.Name() }
func (t *theArrayVectorType) Type() ValueType { return TypeType }
func (t *theArrayVectorType) Unbox() any      { return reflect.TypeFor[*theArrayVectorType]() }

func (t *theArrayVectorType) Name() string { return "let-go.lang.ArrayVector" }

func (t *theArrayVectorType) Box(bare any) (Value, error) {
	arr, ok := bare.([]Value)
	if !ok {
		return NIL, NewTypeError(bare, "can't be boxed as", t)
	}

	return ArrayVector(arr), nil
}

// ArrayVectorType is the type of ArrayVectors
var ArrayVectorType *theArrayVectorType = &theArrayVectorType{}

// ArrayVector is boxed singly linked list that can hold other Values.
type ArrayVector []Value

const arrayVectorPromotionThreshold = 32

func (l ArrayVector) Conj(val Value) Collection {
	newLen := len(l) + 1
	// Promote to PersistentVector when exceeding threshold
	if newLen > arrayVectorPromotionThreshold {
		values := make([]Value, newLen)
		copy(values, l)
		values[newLen-1] = val
		return NewPersistentVector(values).(Collection)
	}
	ret := make([]Value, newLen)
	copy(ret, l)
	ret[newLen-1] = val
	return ArrayVector(ret)
}

// Type implements Value
// Hash implements Hashable. Computed from elements.
func (l ArrayVector) Hash() uint32 {
	h := uint32(1)
	for _, v := range l {
		h = 31*h + hashValue(v)
	}
	return mixFinish(h)
}

// Meta implements IMeta. ArrayVector doesn't store meta, always returns NIL.
func (l ArrayVector) Meta() Value { return NIL }

// WithMeta implements IMeta. Promotes to PersistentVector to carry meta.
func (l ArrayVector) WithMeta(m Value) Value {
	pv := NewPersistentVector([]Value(l)).(PersistentVector)
	pv.meta = m
	return pv
}

func (l ArrayVector) Type() ValueType { return ArrayVectorType }

// Unbox implements Value
func (l ArrayVector) Unbox() any {
	return []Value(l)
}

// Equals implements value equality for ArrayVector
func (l ArrayVector) Equals(other Value) bool {
	switch o := other.(type) {
	case MapEntry:
		return l.Equals(ArrayVector{o.Key, o.Value})
	case ArrayVector:
		if len(l) != len(o) {
			return false
		}
		for i, v := range l {
			if nilListEquivalent(v, o[i]) {
				continue
			}
			if ValueEquals != nil && ValueEquals(v, o[i]) {
				continue
			}
			if eq, ok := v.(interface{ Equals(Value) bool }); ok {
				if !eq.Equals(o[i]) {
					return false
				}
			} else if v != o[i] {
				return false
			}
		}
		return true
	case PersistentVector:
		if len(l) != o.count {
			return false
		}
		for i, v := range l {
			ov := o.ValueAt(Int(i))
			if nilListEquivalent(v, ov) {
				continue
			}
			if ValueEquals != nil && ValueEquals(v, ov) {
				continue
			}
			if eq, ok := v.(interface{ Equals(Value) bool }); ok {
				if !eq.Equals(ov) {
					return false
				}
			} else if v != ov {
				return false
			}
		}
		return true
	default:
		return false
	}
}

func nilListEquivalent(a, b Value) bool {
	return (a == NIL && b == EmptyList) || (a == EmptyList && b == NIL)
}

// First implements Seq
func (l ArrayVector) First() Value {
	if len(l) == 0 {
		return NIL
	}
	return l[0]
}

// More implements Seq
func (l ArrayVector) More() Seq {
	if len(l) <= 1 {
		return EmptyList
	}
	return &ArrayVectorSeq{vec: l, i: 1}
}

// Next implements Seq
func (l ArrayVector) Next() Seq {
	if len(l) <= 1 {
		return nil
	}
	return &ArrayVectorSeq{vec: l, i: 1}
}

// Cons implements Seq
func (l ArrayVector) Cons(val Value) Seq {
	return NewCons(val, l.Seq())
}

func (l ArrayVector) Seq() Seq {
	if len(l) == 0 {
		return EmptyList
	}
	if allocAttrEnabled {
		recordAllocAttr(akArrayVectorSeq, 64)
	}
	return &ArrayVectorSeq{vec: l, i: 0}
}

// SeqFrom returns a seq over v starting at index i, or nil when i is past
// the end. This lets callers step into a vector (e.g. next/rest fast
// paths) with a single allocation instead of Seq()+Next().
func (l ArrayVector) SeqFrom(i int) Seq {
	if i >= len(l) {
		return nil
	}
	return &ArrayVectorSeq{vec: l, i: i}
}

// ArrayVectorSeq is a lightweight seq view over an ArrayVector.
// The embedded chunk lets ChunkedFirst return a pointer into this
// already-heap-allocated seq node instead of allocating a fresh *ArrayChunk
// per chunk — so chunk-walking a (typically tiny, single-chunk) vector costs
// no extra allocation beyond the seq node itself.
//
// Chunk sizes GROW geometrically from 1 (1, 2, 4, 8, … clamped to length)
// rather than a fixed 32 quantum. The size-1 head chunk keeps `first`/short-
// circuit consumption element-wise (no over-realization), while the doubling
// amortizes longer walks — the allocation-minimizing shape given that
// ArrayVectors are overwhelmingly tiny (avg ~1.4 elements).
type ArrayVectorSeq struct {
	vec      ArrayVector
	i        int
	chunkLen int // intended size of THIS chunk; 0 == 1, doubles per ChunkedNext
	chunk    ArrayChunk
}

// curChunkLen is the intended chunk quantum at this position (min 1). A seq
// produced by element-wise Next/More leaves chunkLen zero; a chunk-walk that
// starts from it restarts the geometric growth at 1.
func (s *ArrayVectorSeq) curChunkLen() int {
	if s.chunkLen < 1 {
		return 1
	}
	return s.chunkLen
}

// *ArrayVectorSeq yields whole chunks so map/reduce fast-paths avoid
// per-element seq-node allocation.
var _ IChunkedSeq = (*ArrayVectorSeq)(nil)

func (s *ArrayVectorSeq) String() string {
	str := s.vec[s.i:].String()
	return "(" + str[1:len(str)-1] + ")" // reuse vector's string but change brackets
}

func (s *ArrayVectorSeq) Type() ValueType {
	return ListType // seqs print as lists
}

func (s *ArrayVectorSeq) Unbox() any {
	return []Value(s.vec[s.i:])
}

func (s *ArrayVectorSeq) First() Value {
	if s.i >= len(s.vec) {
		return NIL
	}
	return s.vec[s.i]
}

func (s *ArrayVectorSeq) More() Seq {
	if s.i+1 >= len(s.vec) {
		return EmptyList
	}
	return &ArrayVectorSeq{vec: s.vec, i: s.i + 1}
}

func (s *ArrayVectorSeq) Next() Seq {
	if s.i+1 >= len(s.vec) {
		return nil
	}
	if allocAttrEnabled {
		recordAllocAttr(akArrayVectorSeq, 64)
	}
	return &ArrayVectorSeq{vec: s.vec, i: s.i + 1}
}

// ChunkedFirst returns up to a nodeCap-wide (32) window over the backing
// array starting at the current position. No copy — ArrayChunk shares the
// backing slice; ArrayVector is immutable so aliasing is safe.
func (s *ArrayVectorSeq) ChunkedFirst() IChunk {
	// Clamp to length: the trailing chunk spans only the real remainder, never
	// the full quantum — so ChunkCount() reports the actual size and consumers
	// don't over-allocate/over-realize past the vector's end.
	end := s.i + s.curChunkLen()
	if end > len(s.vec) {
		end = len(s.vec)
	}
	// Populate the embedded chunk and hand back a pointer into this seq node
	// (no per-chunk heap allocation). The chunk is a read-only window over the
	// immutable backing array, so sharing it with the consumer is safe.
	s.chunk = ArrayChunk{vs: []Value(s.vec), off: s.i, end: end}
	return &s.chunk
}

// ChunkedNext advances past this chunk and DOUBLES the quantum, returning the
// seq positioned at the start of the next (larger) chunk, or nil at the end.
func (s *ArrayVectorSeq) ChunkedNext() Seq {
	cl := s.curChunkLen()
	next := s.i + cl
	if next >= len(s.vec) {
		return nil
	}
	return &ArrayVectorSeq{vec: s.vec, i: next, chunkLen: cl * 2}
}

// ChunkedMore is ChunkedNext but returns EmptyList (not nil) at the end.
func (s *ArrayVectorSeq) ChunkedMore() Seq {
	n := s.ChunkedNext()
	if n == nil {
		return EmptyList
	}
	return n
}

func (s *ArrayVectorSeq) Cons(val Value) Seq {
	return NewCons(val, s)
}

func (s *ArrayVectorSeq) Count() Value {
	return Int(len(s.vec) - s.i)
}

func (s *ArrayVectorSeq) RawCount() int {
	return len(s.vec) - s.i
}

func (s *ArrayVectorSeq) Empty() Collection {
	return EmptyList
}

func (s *ArrayVectorSeq) Conj(val Value) Collection {
	values := make([]Value, len(s.vec)-s.i+1)
	values[0] = val
	copy(values[1:], s.vec[s.i:])
	return NewList(values).(*List)
}

func (s *ArrayVectorSeq) Seq() Seq {
	return s
}

// Nth implements Indexed: positional access by integer index.
func (s *ArrayVectorSeq) Nth(i int) Value { return s.ValueAt(Int(i)) }

// ValueAt implements Lookup for ArrayVectorSeq so that `get` works on seq views.
func (s *ArrayVectorSeq) ValueAt(key Value) Value {
	return s.ValueAtOr(key, NIL)
}

// ValueAtOr implements Lookup for ArrayVectorSeq.
func (s *ArrayVectorSeq) ValueAtOr(key Value, dflt Value) Value {
	if key == NIL {
		return dflt
	}
	idx, ok := key.(Int)
	if !ok || idx < 0 {
		return dflt
	}
	absIdx := s.i + int(idx)
	if absIdx >= len(s.vec) {
		return dflt
	}
	return s.vec[absIdx]
}

// Count implements Collection
func (l ArrayVector) Count() Value {
	return Int(len(l))
}

func (l ArrayVector) RawCount() int {
	return len(l)
}

// Empty implements Collection
func (l ArrayVector) Empty() Collection {
	return make(ArrayVector, 0)
}

func NewArrayVector(v []Value) Value {
	if allocAttrEnabled {
		recordAllocAttr(akNewArrayVector, len(v)*16+24)
	}
	vk := make([]Value, len(v))
	copy(vk, v)
	return ArrayVector(vk)
}

func (l ArrayVector) Nth(i int) Value { return l.ValueAt(Int(i)) }

func (l ArrayVector) ValueAt(key Value) Value {
	return l.ValueAtOr(key, NIL)
}

func (l ArrayVector) ValueAtOr(key Value, dflt Value) Value {
	if key == NIL {
		return dflt
	}
	numkey, ok := key.(Int)
	if !ok || numkey < 0 || int(numkey) >= len(l) {
		return dflt
	}
	return l[int(numkey)]
}

func (l ArrayVector) Contains(value Value) Boolean {
	numkey, ok := value.(Int)
	if !ok || numkey < 0 || int(numkey) >= len(l) {
		return FALSE
	}
	return TRUE
}

func (l ArrayVector) Assoc(k Value, v Value) Associative {
	new := NewArrayVector(l).(ArrayVector)
	ik, ok := k.(Int)
	if !ok {
		return NIL
	}
	if ik < 0 || int(ik) > len(new) {
		return NIL
	}
	if int(ik) == len(new) {
		return new.Conj(v).(Associative)
	}
	new[ik] = v
	return new
}

func (l ArrayVector) Dissoc(k Value) Associative {
	return NIL
}

func (l ArrayVector) Arity() int {
	return 1
}

func (l ArrayVector) Invoke(pargs []Value) (Value, error) {
	vl := len(pargs)
	if vl != 1 {
		return NIL, fmt.Errorf("wrong number of arguments %d", vl)
	}
	idx, ok := pargs[0].(Int)
	if !ok {
		return NIL, fmt.Errorf("vector key must be Int")
	}
	i := int(idx)
	if i < 0 || i >= len(l) {
		return NIL, fmt.Errorf("index out of bounds: %d", i)
	}
	return l[i], nil
}

func (l ArrayVector) String() string {
	b := &strings.Builder{}
	b.WriteRune('[')
	n := len(l)
	for i := range l {
		b.WriteString(l[i].String())
		if i < n-1 {
			b.WriteRune(' ')
		}
	}
	b.WriteRune(']')
	return b.String()
}
