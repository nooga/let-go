/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package vm

import "testing"

func buildVec(n int) PersistentVector {
	var v Collection = PersistentVector{}.Empty()
	for i := 0; i < n; i++ {
		v = v.Conj(Int(i))
	}
	return v.(PersistentVector)
}

func walkChunked(t *testing.T, s Seq) []int {
	t.Helper()
	out := []int{}
	for s != nil && s != EmptyList {
		cs, ok := s.(IChunkedSeq)
		if !ok {
			t.Fatalf("expected IChunkedSeq, got %T", s)
		}
		c := cs.ChunkedFirst()
		for i := 0; i < c.ChunkCount(); i++ {
			out = append(out, int(c.Nth(i).(Int)))
		}
		s = cs.ChunkedNext()
	}
	return out
}

func TestVectorSeqImplementsIChunkedSeq(t *testing.T) {
	v := buildVec(5)
	s := v.Seq()
	if _, ok := s.(IChunkedSeq); !ok {
		t.Fatalf("PersistentVectorSeq should implement IChunkedSeq, got %T", s)
	}
}

func TestChunkedWalkSmallTailOnly(t *testing.T) {
	// count < 32: everything is in the tail
	v := buildVec(7)
	got := walkChunked(t, v.Seq())
	want := []int{0, 1, 2, 3, 4, 5, 6}
	if len(got) != len(want) {
		t.Fatalf("len: got %d want %d (%v)", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("[%d] got %d want %d", i, got[i], want[i])
		}
	}
}

func TestChunkedWalkExactlyOneLeafPlusTail(t *testing.T) {
	v := buildVec(40) // 32 in trie leaf + 8 in tail
	got := walkChunked(t, v.Seq())
	if len(got) != 40 {
		t.Fatalf("len: got %d want 40", len(got))
	}
	for i := 0; i < 40; i++ {
		if got[i] != i {
			t.Fatalf("[%d] got %d want %d", i, got[i], i)
		}
	}
}

func TestChunkedWalkLargeMultiLevel(t *testing.T) {
	// 1024 elements forces a second level in the trie (32^2 = 1024)
	const n = 1024
	v := buildVec(n)
	got := walkChunked(t, v.Seq())
	if len(got) != n {
		t.Fatalf("len: got %d want %d", len(got), n)
	}
	for i := 0; i < n; i++ {
		if got[i] != i {
			t.Fatalf("[%d] got %d want %d", i, got[i], i)
		}
	}
}

func TestChunkedFirstChunkSizeIs32(t *testing.T) {
	v := buildVec(100)
	cs := v.Seq().(IChunkedSeq)
	first := cs.ChunkedFirst()
	if first.ChunkCount() != 32 {
		t.Fatalf("first chunk count: got %d want 32", first.ChunkCount())
	}
	// after one ChunkedNext, still 32
	cs = cs.ChunkedNext().(IChunkedSeq)
	second := cs.ChunkedFirst()
	if second.ChunkCount() != 32 {
		t.Fatalf("second chunk count: got %d want 32", second.ChunkCount())
	}
}

func TestChunkedMoreEmptyListAtEnd(t *testing.T) {
	v := buildVec(5) // tail only, single chunk
	cs := v.Seq().(IChunkedSeq)
	more := cs.ChunkedMore()
	if more != EmptyList {
		t.Fatalf("ChunkedMore at end: got %v, want EmptyList", more)
	}
	if cs.ChunkedNext() != nil {
		t.Fatal("ChunkedNext at end should be nil")
	}
}

func TestVectorSeqElementWalkUnaffected(t *testing.T) {
	// Element-by-element walk via Seq.Next() must remain consistent with
	// the chunked view.
	v := buildVec(100)
	got := []int{}
	for s := v.Seq(); s != nil && s != EmptyList; s = s.Next() {
		got = append(got, int(s.First().(Int)))
	}
	if len(got) != 100 {
		t.Fatalf("len: got %d want 100", len(got))
	}
	for i := 0; i < 100; i++ {
		if got[i] != i {
			t.Fatalf("[%d] got %d want %d", i, got[i], i)
		}
	}
}

func TestVectorChunkedViaAsChunkedSeq(t *testing.T) {
	v := buildVec(50)
	cs, ok := AsChunkedSeq(v.Seq())
	if !ok {
		t.Fatal("AsChunkedSeq should accept PersistentVectorSeq")
	}
	if cs.ChunkedFirst().ChunkCount() != 32 {
		t.Fatal("chunk size wrong via AsChunkedSeq")
	}
}

func TestChunkedFirstAfterPartialAdvance(t *testing.T) {
	// If a seq is mid-chunk (advanced via Next), ChunkedFirst should return
	// the remainder of the current chunk, not start over.
	v := buildVec(50)
	s := v.Seq()
	s = s.Next().Next().Next() // now at index 3
	cs := s.(IChunkedSeq)
	c := cs.ChunkedFirst()
	if c.ChunkCount() != 32-3 {
		t.Fatalf("partial-advance chunk count: got %d want %d", c.ChunkCount(), 32-3)
	}
	if c.Nth(0) != Int(3) {
		t.Fatalf("partial-advance Nth(0): got %v want 3", c.Nth(0))
	}
}
