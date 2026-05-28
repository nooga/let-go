/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package vm

import "strings"

// IChunk is a fixed-size, indexable, sliceable batch of Values.
// Implementations may share backing storage; DropFirst must never mutate.
//
// Note: the per-chunk size is exposed as ChunkCount() rather than Count()
// because concrete chunks (ArrayChunk, ChunkBuffer) also satisfy the Counted
// protocol — which requires Count() Value and RawCount() int — and Go
// interfaces cannot disambiguate two methods with the same name and
// different return types.
type IChunk interface {
	Value
	ChunkCount() int
	Nth(i int) Value
	DropFirst() IChunk
}

// IChunkedSeq is a Seq that can yield a whole chunk at a time.
// Consumers may dispatch on this interface to amortize per-element overhead;
// the Seq methods (First/Next/More) must remain consistent with the chunked
// view: First() == ChunkedFirst().Nth(0), Next() walks within the chunk
// before moving to ChunkedNext().
type IChunkedSeq interface {
	Seq
	ChunkedFirst() IChunk
	ChunkedNext() Seq // nil at end
	ChunkedMore() Seq // EmptyList at end
}

// ArrayChunk is a slice-backed IChunk. It uses an off/end window over a
// shared backing array so DropFirst is allocation-free.
type ArrayChunk struct {
	vs  []Value
	off int
	end int
}

// NewArrayChunk wraps vs[off:end] as an IChunk. The slice is shared; callers
// must treat it as immutable from this point.
func NewArrayChunk(vs []Value, off, end int) *ArrayChunk {
	if off < 0 || end < off || end > len(vs) {
		panic("NewArrayChunk: bad bounds")
	}
	return &ArrayChunk{vs: vs, off: off, end: end}
}

// NewArrayChunkFromValues copies vs into a fresh chunk. Convenience for
// callers that don't already have a stable backing slice.
func NewArrayChunkFromValues(vs []Value) *ArrayChunk {
	cp := make([]Value, len(vs))
	copy(cp, vs)
	return &ArrayChunk{vs: cp, off: 0, end: len(cp)}
}

func (c *ArrayChunk) ChunkCount() int { return c.end - c.off }
func (c *ArrayChunk) Nth(i int) Value { return c.vs[c.off+i] }
func (c *ArrayChunk) Type() ValueType { return ListType }
func (c *ArrayChunk) Unbox() any      { return c.vs[c.off:c.end] }

// Counted lets Lisp `(count chunk)` work on a chunk directly.
func (c *ArrayChunk) Count() Value  { return Int(c.end - c.off) }
func (c *ArrayChunk) RawCount() int { return c.end - c.off }

// Lookup so `(nth chunk i)` and `(get chunk i)` work.
func (c *ArrayChunk) ValueAt(key Value) Value {
	return c.ValueAtOr(key, NIL)
}

func (c *ArrayChunk) ValueAtOr(key Value, dflt Value) Value {
	idx, ok := key.(Int)
	if !ok {
		return dflt
	}
	i := int(idx)
	n := c.end - c.off
	if i < 0 || i >= n {
		return dflt
	}
	return c.vs[c.off+i]
}

func (c *ArrayChunk) DropFirst() IChunk {
	if c.off >= c.end {
		panic("DropFirst on empty chunk")
	}
	return &ArrayChunk{vs: c.vs, off: c.off + 1, end: c.end}
}

func (c *ArrayChunk) String() string {
	var b strings.Builder
	b.WriteByte('[')
	for i := c.off; i < c.end; i++ {
		if i > c.off {
			b.WriteByte(' ')
		}
		b.WriteString(c.vs[i].String())
	}
	b.WriteByte(']')
	return b.String()
}

// ChunkedCons pairs a non-empty IChunk with a tail seq. It is itself both a
// Seq and an IChunkedSeq: the chunk is the head batch, `more` is whatever
// comes after.
type ChunkedCons struct {
	chunk IChunk
	more  Seq
}

// NewChunkedCons returns a ChunkedCons. The chunk must be non-empty; callers
// that may have an empty chunk should use ConsChunk which folds the empty
// case down to `more`.
func NewChunkedCons(chunk IChunk, more Seq) *ChunkedCons {
	if chunk == nil || chunk.ChunkCount() == 0 {
		panic("NewChunkedCons: empty chunk")
	}
	return &ChunkedCons{chunk: chunk, more: more}
}

// ConsChunk produces the natural seq of `chunk` followed by `more`. If
// chunk is empty/nil, returns `more` unchanged.
func ConsChunk(chunk IChunk, more Seq) Seq {
	if chunk == nil || chunk.ChunkCount() == 0 {
		if more == nil {
			return nil
		}
		return more
	}
	return &ChunkedCons{chunk: chunk, more: more}
}

func (c *ChunkedCons) First() Value { return c.chunk.Nth(0) }

func (c *ChunkedCons) More() Seq {
	if c.chunk.ChunkCount() > 1 {
		return &ChunkedCons{chunk: c.chunk.DropFirst(), more: c.more}
	}
	if c.more == nil {
		return EmptyList
	}
	return c.more
}

func (c *ChunkedCons) Next() Seq {
	if c.chunk.ChunkCount() > 1 {
		return &ChunkedCons{chunk: c.chunk.DropFirst(), more: c.more}
	}
	// chunk has exactly one element left; advance into `more`
	if c.more == nil || c.more == EmptyList {
		return nil
	}
	// If tail is a LazySeq, resolve so we don't expose a wrapper that
	// realizes to empty as a "non-empty" seq. Mirrors Cons.Next.
	if ls, ok := c.more.(*LazySeq); ok {
		s := ls.Resolve()
		if s == nil {
			return nil
		}
		return s
	}
	return c.more
}

func (c *ChunkedCons) Cons(val Value) Seq {
	// Consing onto the head defeats chunking for that one element; that's
	// fine — return a regular Cons.
	return NewCons(val, c)
}

func (c *ChunkedCons) ChunkedFirst() IChunk { return c.chunk }

func (c *ChunkedCons) ChunkedNext() Seq {
	if c.more == nil || c.more == EmptyList {
		return nil
	}
	if ls, ok := c.more.(*LazySeq); ok {
		s := ls.Resolve()
		if s == nil {
			return nil
		}
		return s
	}
	return c.more
}

func (c *ChunkedCons) ChunkedMore() Seq {
	if c.more == nil {
		return EmptyList
	}
	return c.more
}

func (c *ChunkedCons) Seq() Seq        { return c }
func (c *ChunkedCons) Type() ValueType { return ListType }
func (c *ChunkedCons) Unbox() any      { return c }

func (c *ChunkedCons) String() string {
	var b strings.Builder
	b.WriteByte('(')
	first := true
	for s := Seq(c); s != nil && s != EmptyList; s = s.Next() {
		if !first {
			b.WriteByte(' ')
		}
		first = false
		b.WriteString(s.First().String())
	}
	b.WriteByte(')')
	return b.String()
}

// Hash implements Hashable, mirroring Cons/List ordered hashing so the three
// seq surfaces equate in sets and maps.
func (c *ChunkedCons) Hash() uint32 { return hashOrdered(c) }

// ChunkBuffer is a mutable builder for an ArrayChunk. Use Append to fill, then
// Chunk() to finalize — the returned IChunk takes ownership of the buffer.
type ChunkBuffer struct {
	vs []Value
}

// NewChunkBuffer returns a buffer pre-sized for cap elements. Append grows
// beyond cap if needed.
func NewChunkBuffer(cap int) *ChunkBuffer {
	if cap < 0 {
		cap = 0
	}
	return &ChunkBuffer{vs: make([]Value, 0, cap)}
}

func (b *ChunkBuffer) Append(v Value) { b.vs = append(b.vs, v) }

// Counted: lets `(count buf)` return the number of elements buffered so far.
func (b *ChunkBuffer) Count() Value  { return Int(len(b.vs)) }
func (b *ChunkBuffer) RawCount() int { return len(b.vs) }

// Chunk finalizes the buffer into an IChunk and clears the buffer so further
// Append calls operate on a fresh backing slice.
func (b *ChunkBuffer) Chunk() IChunk {
	out := &ArrayChunk{vs: b.vs, off: 0, end: len(b.vs)}
	b.vs = nil
	return out
}

// ChunkBuffer satisfies Value so it can be threaded through Lisp.
func (b *ChunkBuffer) Type() ValueType { return ListType }
func (b *ChunkBuffer) Unbox() any      { return b }
func (b *ChunkBuffer) String() string  { return "#<chunk-buffer>" }
