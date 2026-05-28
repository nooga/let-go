/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package vm

import "testing"

func TestRangeImplementsIChunkedSeq(t *testing.T) {
	r := NewRange(Int(0), Int(10), Int(1))
	if _, ok := r.(IChunkedSeq); !ok {
		t.Fatalf("Range should implement IChunkedSeq, got %T", r)
	}
}

func TestRangeChunkedFirstSmall(t *testing.T) {
	r := NewRange(Int(0), Int(5), Int(1)).(*Range)
	c := r.ChunkedFirst()
	if c.ChunkCount() != 5 {
		t.Fatalf("count: got %d want 5", c.ChunkCount())
	}
	for i := 0; i < 5; i++ {
		if c.Nth(i) != Int(i) {
			t.Fatalf("Nth(%d): got %v want %d", i, c.Nth(i), i)
		}
	}
	if r.ChunkedNext() != nil {
		t.Fatal("small range: ChunkedNext should be nil")
	}
}

func TestRangeChunkedFirstFull(t *testing.T) {
	r := NewRange(Int(0), Int(32), Int(1)).(*Range)
	c := r.ChunkedFirst()
	if c.ChunkCount() != 32 {
		t.Fatalf("count: got %d want 32", c.ChunkCount())
	}
	if r.ChunkedNext() != nil {
		t.Fatal("exactly-full range: ChunkedNext should be nil")
	}
}

func TestRangeChunkedMulti(t *testing.T) {
	// 100 elements → chunks of 32, 32, 32, 4
	r := NewRange(Int(0), Int(100), Int(1)).(*Range)
	got := []int{}
	var s Seq = r
	for s != nil {
		cs := s.(IChunkedSeq)
		c := cs.ChunkedFirst()
		for i := 0; i < c.ChunkCount(); i++ {
			got = append(got, int(c.Nth(i).(Int)))
		}
		s = cs.ChunkedNext()
	}
	if len(got) != 100 {
		t.Fatalf("len: got %d want 100", len(got))
	}
	for i := 0; i < 100; i++ {
		if got[i] != i {
			t.Fatalf("[%d]: got %d want %d", i, got[i], i)
		}
	}
}

func TestRangeWithStep(t *testing.T) {
	// (range 0 20 2) → 0,2,4,...,18 = 10 elements
	r := NewRange(Int(0), Int(20), Int(2)).(*Range)
	c := r.ChunkedFirst()
	if c.ChunkCount() != 10 {
		t.Fatalf("count: got %d want 10", c.ChunkCount())
	}
	if c.Nth(0) != Int(0) || c.Nth(9) != Int(18) {
		t.Fatalf("step range values wrong: 0->%v 9->%v", c.Nth(0), c.Nth(9))
	}
	if r.ChunkedNext() != nil {
		t.Fatal("step range should fit in one chunk")
	}
}

func TestRangeWithNegativeStep(t *testing.T) {
	r := NewRange(Int(10), Int(0), Int(-1)).(*Range)
	c := r.ChunkedFirst()
	if c.ChunkCount() != 10 {
		t.Fatalf("count: got %d want 10", c.ChunkCount())
	}
	if c.Nth(0) != Int(10) || c.Nth(9) != Int(1) {
		t.Fatalf("negative-step values wrong: 0->%v 9->%v", c.Nth(0), c.Nth(9))
	}
}

func TestRangeChunkedNextStartsCorrectly(t *testing.T) {
	r := NewRange(Int(0), Int(64), Int(1)).(*Range)
	n := r.ChunkedNext().(*Range)
	if n.start != 32 {
		t.Fatalf("next chunk start: got %d want 32", n.start)
	}
	c := n.ChunkedFirst()
	if c.ChunkCount() != 32 || c.Nth(0) != Int(32) || c.Nth(31) != Int(63) {
		t.Fatal("next chunk values wrong")
	}
}

func TestRangeChunkedMoreEmptyListAtEnd(t *testing.T) {
	r := NewRange(Int(0), Int(5), Int(1)).(*Range)
	if more := r.ChunkedMore(); more != EmptyList {
		t.Fatalf("ChunkedMore at end: got %v want EmptyList", more)
	}
}

func TestRangeElementWalkConsistentWithChunked(t *testing.T) {
	r := NewRange(Int(0), Int(100), Int(1)).(*Range)
	// Element walk
	elem := []int{}
	for s := Seq(r); s != nil; s = s.Next() {
		elem = append(elem, int(s.First().(Int)))
	}
	// Chunked walk
	chunked := []int{}
	var s Seq = r
	for s != nil {
		cs := s.(IChunkedSeq)
		c := cs.ChunkedFirst()
		for i := 0; i < c.ChunkCount(); i++ {
			chunked = append(chunked, int(c.Nth(i).(Int)))
		}
		s = cs.ChunkedNext()
	}
	if len(elem) != len(chunked) {
		t.Fatalf("len mismatch: elem=%d chunked=%d", len(elem), len(chunked))
	}
	for i := range elem {
		if elem[i] != chunked[i] {
			t.Fatalf("[%d]: elem=%d chunked=%d", i, elem[i], chunked[i])
		}
	}
}

func TestInfiniteRangeChunked(t *testing.T) {
	r := NewInfiniteRange(0, 1)
	c := r.ChunkedFirst()
	if c.ChunkCount() != 32 {
		t.Fatalf("infinite chunk count: got %d want 32", c.ChunkCount())
	}
	if c.Nth(0) != Int(0) || c.Nth(31) != Int(31) {
		t.Fatal("infinite chunk values wrong")
	}
	next := r.ChunkedNext()
	if next == nil {
		t.Fatal("InfiniteRange.ChunkedNext must never be nil")
	}
	c2 := next.(*InfiniteRange).ChunkedFirst()
	if c2.Nth(0) != Int(32) {
		t.Fatalf("second infinite chunk start: got %v want 32", c2.Nth(0))
	}
}

func TestInfiniteRangeAsChunkedSeq(t *testing.T) {
	r := NewInfiniteRange(0, 2)
	cs, ok := AsChunkedSeq(r)
	if !ok {
		t.Fatal("InfiniteRange should be IChunkedSeq via AsChunkedSeq")
	}
	c := cs.ChunkedFirst()
	if c.Nth(1) != Int(2) {
		t.Fatalf("step-2 infinite range chunk wrong: %v", c.Nth(1))
	}
}
