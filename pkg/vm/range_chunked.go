/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package vm

// rangeChunkSize matches the leaf width used elsewhere (PersistentVector,
// ChunkedCons consumers). Keeping it 32 means a chunked-seq fast path over a
// range mirrors the cost profile of a chunked seq over a vector.
const rangeChunkSize = nodeCap // 32

// ChunkedFirst returns up to rangeChunkSize values starting from l.start,
// stopping at the range's end. The chunk holds materialized Int values; this
// is one O(chunkSize) allocation per chunk in exchange for amortizing the
// element-wise allocation overhead of l.Next().
func (l *Range) ChunkedFirst() IChunk {
	n := l.RawCount()
	if n == 0 {
		return NewArrayChunk([]Value{}, 0, 0)
	}
	if n > rangeChunkSize {
		n = rangeChunkSize
	}
	vs := make([]Value, n)
	v := l.start
	for i := 0; i < n; i++ {
		vs[i] = Int(v)
		v += l.step
	}
	return NewArrayChunk(vs, 0, n)
}

// ChunkedNext returns the range positioned at the element after the current
// chunk, or nil if the current chunk consumes the whole range.
func (l *Range) ChunkedNext() Seq {
	n := l.RawCount()
	if n <= rangeChunkSize {
		return nil
	}
	return &Range{
		start: l.start + rangeChunkSize*l.step,
		end:   l.end,
		step:  l.step,
	}
}

func (l *Range) ChunkedMore() Seq {
	n := l.ChunkedNext()
	if n == nil {
		return EmptyList
	}
	return n
}

// InfiniteRange's chunked view emits a full rangeChunkSize chunk every time
// and never returns nil from ChunkedNext/More — consumers must short-circuit.
func (r *InfiniteRange) ChunkedFirst() IChunk {
	vs := make([]Value, rangeChunkSize)
	v := r.start
	for i := 0; i < rangeChunkSize; i++ {
		vs[i] = Int(v)
		v += r.step
	}
	return NewArrayChunk(vs, 0, rangeChunkSize)
}

func (r *InfiniteRange) ChunkedNext() Seq {
	return &InfiniteRange{
		start: r.start + rangeChunkSize*r.step,
		step:  r.step,
	}
}

func (r *InfiniteRange) ChunkedMore() Seq { return r.ChunkedNext() }
