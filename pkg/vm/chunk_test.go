/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package vm

import "testing"

func ints(n int) []Value {
	out := make([]Value, n)
	for i := range out {
		out[i] = Int(i)
	}
	return out
}

func TestArrayChunkBasics(t *testing.T) {
	c := NewArrayChunk(ints(5), 0, 5)
	if c.ChunkCount() != 5 {
		t.Fatalf("count: got %d, want 5", c.ChunkCount())
	}
	for i := 0; i < 5; i++ {
		if c.Nth(i) != Int(i) {
			t.Fatalf("Nth(%d) = %v, want %d", i, c.Nth(i), i)
		}
	}
}

func TestArrayChunkDropFirst(t *testing.T) {
	c := NewArrayChunk(ints(5), 0, 5)
	d := c.DropFirst()
	if d.ChunkCount() != 4 {
		t.Fatalf("DropFirst count: got %d, want 4", d.ChunkCount())
	}
	if d.Nth(0) != Int(1) {
		t.Fatalf("DropFirst Nth(0): got %v, want 1", d.Nth(0))
	}
	// original is unchanged (no mutation)
	if c.ChunkCount() != 5 || c.Nth(0) != Int(0) {
		t.Fatalf("DropFirst mutated original")
	}
}

func TestArrayChunkBadBounds(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic")
		}
	}()
	NewArrayChunk(ints(3), 0, 10)
}

func TestDropFirstOnEmptyPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic")
		}
	}()
	c := NewArrayChunk(ints(1), 0, 1)
	c.DropFirst().DropFirst()
}

func TestChunkedConsSeq(t *testing.T) {
	c := NewArrayChunk(ints(3), 0, 3)
	cc := NewChunkedCons(c, nil)
	got := []Value{}
	for s := Seq(cc); s != nil && s != EmptyList; s = s.Next() {
		got = append(got, s.First())
	}
	if len(got) != 3 {
		t.Fatalf("len: got %d, want 3", len(got))
	}
	for i, v := range got {
		if v != Int(i) {
			t.Fatalf("got[%d] = %v, want %d", i, v, i)
		}
	}
}

func TestChunkedConsWithTail(t *testing.T) {
	tail := NewChunkedCons(NewArrayChunk(ints(2), 0, 2), nil) // 0,1
	head := NewChunkedCons(NewArrayChunk([]Value{Int(10), Int(11), Int(12)}, 0, 3), tail)
	want := []int{10, 11, 12, 0, 1}
	got := []int{}
	for s := Seq(head); s != nil && s != EmptyList; s = s.Next() {
		got = append(got, int(s.First().(Int)))
	}
	if len(got) != len(want) {
		t.Fatalf("len: got %d, want %d (got=%v)", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got[%d] = %d, want %d", i, got[i], want[i])
		}
	}
}

func TestChunkedConsChunkedView(t *testing.T) {
	tail := NewChunkedCons(NewArrayChunk(ints(2), 0, 2), nil)
	head := NewChunkedCons(NewArrayChunk([]Value{Int(10), Int(11)}, 0, 2), tail)
	// First chunk
	c1 := head.ChunkedFirst()
	if c1.ChunkCount() != 2 || c1.Nth(0) != Int(10) || c1.Nth(1) != Int(11) {
		t.Fatalf("ChunkedFirst wrong: %v", c1)
	}
	// Advance
	n := head.ChunkedNext()
	if n == nil {
		t.Fatal("ChunkedNext nil, want next chunk")
	}
	cn, ok := n.(IChunkedSeq)
	if !ok {
		t.Fatalf("ChunkedNext not IChunkedSeq: %T", n)
	}
	c2 := cn.ChunkedFirst()
	if c2.ChunkCount() != 2 || c2.Nth(0) != Int(0) || c2.Nth(1) != Int(1) {
		t.Fatalf("second chunk wrong: %v", c2)
	}
	if cn.ChunkedNext() != nil {
		t.Fatal("expected ChunkedNext nil at end")
	}
}

func TestConsChunkEmptyFolds(t *testing.T) {
	empty := NewArrayChunk([]Value{}, 0, 0)
	if got := ConsChunk(empty, nil); got != nil {
		t.Fatalf("ConsChunk(empty, nil): got %v, want nil", got)
	}
	tail := NewChunkedCons(NewArrayChunk(ints(1), 0, 1), nil)
	if got := ConsChunk(empty, tail); got != tail {
		t.Fatal("ConsChunk(empty, tail) should pass tail through")
	}
}

func TestChunkedConsImplementsIChunkedSeq(t *testing.T) {
	var _ IChunkedSeq = (*ChunkedCons)(nil)
	var _ Seq = (*ChunkedCons)(nil)
}

func TestChunkBufferAppendChunk(t *testing.T) {
	b := NewChunkBuffer(4)
	for i := 0; i < 4; i++ {
		b.Append(Int(i * 10))
	}
	c := b.Chunk()
	if c.ChunkCount() != 4 {
		t.Fatalf("count: got %d, want 4", c.ChunkCount())
	}
	for i := 0; i < 4; i++ {
		if c.Nth(i) != Int(i*10) {
			t.Fatalf("Nth(%d) = %v, want %d", i, c.Nth(i), i*10)
		}
	}
	// buffer is now empty for reuse
	if b.RawCount() != 0 {
		t.Fatalf("buffer not reset: %d", b.RawCount())
	}
	b.Append(Int(99))
	c2 := b.Chunk()
	if c2.ChunkCount() != 1 || c2.Nth(0) != Int(99) {
		t.Fatal("reused buffer wrong")
	}
	// first chunk must be unchanged
	if c.ChunkCount() != 4 || c.Nth(0) != Int(0) {
		t.Fatal("ChunkBuffer.Chunk did not detach from buffer")
	}
}

func TestChunkBufferGrowsBeyondCap(t *testing.T) {
	b := NewChunkBuffer(2)
	for i := 0; i < 10; i++ {
		b.Append(Int(i))
	}
	c := b.Chunk()
	if c.ChunkCount() != 10 {
		t.Fatalf("count after grow: got %d, want 10", c.ChunkCount())
	}
	if c.Nth(9) != Int(9) {
		t.Fatal("last element wrong after grow")
	}
}

func TestChunkedConsMoreReturnsEmptyListAtEnd(t *testing.T) {
	c := NewChunkedCons(NewArrayChunk(ints(1), 0, 1), nil)
	more := c.More()
	if more != EmptyList {
		t.Fatalf("More at end: got %v, want EmptyList", more)
	}
}
