/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package vm

// PersistentVectorSeq implements IChunkedSeq in addition to Seq. The chunked
// view exposes the trie's 32-wide leaves (and the tail) directly, letting
// chunk-aware consumers amortize per-element overhead.
//
// The element-wise Seq methods (First/Next/More) are unaffected — chunked and
// element walks are independent views of the same sequence.

// leafForChunk walks the trie to the leaf containing absolute index i.
// Precondition: 0 <= i < v.tailOff.
func (v *PersistentVector) leafForChunk(i int) []any {
	node := v.root
	level := v.shift
	for level > 0 {
		subidx := (i >> level) & nodeMask
		node = node.array[subidx].(*vnode)
		level -= shift
	}
	return node.array
}

// chunkAt returns the chunk that contains absolute index i, with the chunk
// window starting at i (i.e. any offset within the chunk is already applied).
//
// For trie elements this lifts the leaf into a fresh []Value (one O(32) copy
// per chunk; total walk is still O(n)). For tail elements it reuses the tail
// slice directly.
func (s *PersistentVectorSeq) chunkAt(i int) IChunk {
	v := s.vec
	if i >= v.tailOff {
		// In tail. tail is already []Value, share it.
		return NewArrayChunk(v.tail, i-v.tailOff, len(v.tail))
	}
	leaf := v.leafForChunk(i)
	chunkStart := i &^ nodeMask
	// Lift []any → []Value. We always lift the whole leaf, then the off
	// window inside ArrayChunk handles the in-chunk offset.
	vs := make([]Value, len(leaf))
	for j, x := range leaf {
		vs[j], _ = x.(Value)
	}
	return NewArrayChunk(vs, i-chunkStart, len(vs))
}

// ChunkedFirst returns the chunk starting at the current position.
func (s *PersistentVectorSeq) ChunkedFirst() IChunk {
	return s.chunkAt(s.i)
}

// ChunkedNext advances by one whole chunk, returning the seq positioned at
// the start of the following chunk, or nil at the end.
func (s *PersistentVectorSeq) ChunkedNext() Seq {
	v := s.vec
	var nextStart int
	if s.i >= v.tailOff {
		return nil // tail is the last chunk
	}
	// Advance to next 32-aligned boundary past the current position. The
	// current chunk's window starts at s.i and ends at the next leaf
	// boundary (or tailOff for the last trie chunk).
	nextStart = (s.i &^ nodeMask) + nodeCap
	if nextStart >= v.count {
		return nil
	}
	return &PersistentVectorSeq{vec: v, i: nextStart}
}

// ChunkedMore is like ChunkedNext but returns EmptyList instead of nil at
// the end, matching Clojure's IChunkedSeq.chunkedMore.
func (s *PersistentVectorSeq) ChunkedMore() Seq {
	n := s.ChunkedNext()
	if n == nil {
		return EmptyList
	}
	return n
}
