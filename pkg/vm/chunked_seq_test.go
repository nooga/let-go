/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package vm

import "testing"

// thunkOf wraps a Value-producing function as a vm.Fn for LazySeq tests.
func thunkOf(produce func() Value) Fn {
	f, _ := NativeFnType.Wrap(func(_ []Value) (Value, error) {
		return produce(), nil
	})
	return f.(Fn)
}

func TestAsChunkedSeqDirect(t *testing.T) {
	cc := NewChunkedCons(NewArrayChunk(ints(3), 0, 3), nil)
	cs, ok := AsChunkedSeq(cc)
	if !ok {
		t.Fatal("direct ChunkedCons should be IChunkedSeq")
	}
	if cs.ChunkedFirst().ChunkCount() != 3 {
		t.Fatal("ChunkedFirst count wrong")
	}
}

func TestAsChunkedSeqThroughLazySeq(t *testing.T) {
	cc := NewChunkedCons(NewArrayChunk(ints(3), 0, 3), nil)
	ls := NewLazySeq(thunkOf(func() Value { return cc }))
	cs, ok := AsChunkedSeq(ls)
	if !ok {
		t.Fatal("LazySeq wrapping ChunkedCons should resolve to IChunkedSeq")
	}
	if cs.ChunkedFirst().Nth(0) != Int(0) {
		t.Fatal("resolved chunk wrong")
	}
}

func TestAsChunkedSeqThroughNestedLazySeq(t *testing.T) {
	cc := NewChunkedCons(NewArrayChunk(ints(2), 0, 2), nil)
	inner := NewLazySeq(thunkOf(func() Value { return cc }))
	outer := NewLazySeq(thunkOf(func() Value { return inner }))
	cs, ok := AsChunkedSeq(outer)
	if !ok {
		t.Fatal("nested LazySeq should resolve to IChunkedSeq")
	}
	if cs.ChunkedFirst().ChunkCount() != 2 {
		t.Fatal("resolved chunk count wrong")
	}
}

func TestAsChunkedSeqOnNonChunked(t *testing.T) {
	if _, ok := AsChunkedSeq(nil); ok {
		t.Fatal("nil should not be chunked")
	}
	if _, ok := AsChunkedSeq(EmptyList); ok {
		t.Fatal("EmptyList should not be chunked")
	}
	c := NewCons(Int(1), nil)
	if _, ok := AsChunkedSeq(c); ok {
		t.Fatal("plain Cons should not be chunked")
	}
}

func TestConsOverChunkedTailWalksCorrectly(t *testing.T) {
	// (cons 99 <chunked 0,1,2>) → 99, 0, 1, 2
	tail := NewChunkedCons(NewArrayChunk(ints(3), 0, 3), nil)
	head := NewCons(Int(99), tail)
	want := []int{99, 0, 1, 2}
	got := []int{}
	for s := Seq(head); s != nil && s != EmptyList; s = s.Next() {
		got = append(got, int(s.First().(Int)))
	}
	if len(got) != len(want) {
		t.Fatalf("len: got %d want %d (got=%v)", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got[%d]=%d want %d", i, got[i], want[i])
		}
	}
}

func TestLazySeqOfChunkedConsWalksCorrectly(t *testing.T) {
	// Realizing a LazySeq that produces a ChunkedCons should iterate
	// through all chunk elements without losing any.
	tail := NewChunkedCons(NewArrayChunk([]Value{Int(20), Int(21)}, 0, 2), nil)
	cc := NewChunkedCons(NewArrayChunk([]Value{Int(10), Int(11), Int(12)}, 0, 3), tail)
	ls := NewLazySeq(thunkOf(func() Value { return cc }))
	want := []int{10, 11, 12, 20, 21}
	got := []int{}
	for v := range Iter(ls) {
		got = append(got, int(v.(Int)))
	}
	if len(got) != len(want) {
		t.Fatalf("len: got %d want %d (got=%v)", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got[%d]=%d want %d", i, got[i], want[i])
		}
	}
}

func TestChunkedNextThroughLazySeqTail(t *testing.T) {
	// ChunkedCons with a LazySeq tail that resolves to another ChunkedCons
	// — ChunkedNext should walk through.
	tail := NewChunkedCons(NewArrayChunk(ints(2), 0, 2), nil) // 0,1
	lazyTail := NewLazySeq(thunkOf(func() Value { return tail }))
	head := NewChunkedCons(NewArrayChunk([]Value{Int(100)}, 0, 1), lazyTail)
	cn := head.ChunkedNext()
	if cn == nil {
		t.Fatal("ChunkedNext through LazySeq tail returned nil")
	}
	cs, ok := AsChunkedSeq(cn)
	if !ok {
		t.Fatalf("ChunkedNext didn't yield IChunkedSeq, got %T", cn)
	}
	if cs.ChunkedFirst().ChunkCount() != 2 {
		t.Fatal("second chunk count wrong")
	}
}
