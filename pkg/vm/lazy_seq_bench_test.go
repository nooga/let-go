/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package vm

import "testing"

// lazyChain builds an n-element chain of nested LazySeqs — each cell a
// LazySeq realizing to a Cons whose tail is the next LazySeq — the shape
// map/filter produce.
func lazyChain(n int) *LazySeq {
	var mk func(i int) *LazySeq
	mk = func(i int) *LazySeq {
		return NewLazySeq(thunkOf(func() Value {
			if i >= n {
				return NIL
			}
			return NewCons(Int(i), mk(i+1))
		}))
	}
	return mk(0)
}

func walkSum(s Seq) int {
	total := 0
	for ; s != nil; s = s.Next() {
		if v, ok := s.First().(Int); ok {
			total += int(v)
		}
	}
	return total
}

// BenchmarkLazySeqRealizedWalk measures re-walking an already-realized lazy
// chain: the per-element accessor cost after realization (Cons.Next resolves
// its lazy tail on every call, so each cell visit goes through LazySeq.seq).
func BenchmarkLazySeqRealizedWalk(b *testing.B) {
	const n = 1024
	ls := lazyChain(n)
	want := n * (n - 1) / 2
	if got := walkSum(ls.Resolve()); got != want {
		b.Fatalf("realization walk sum = %d, want %d", got, want)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if got := walkSum(ls.Resolve()); got != want {
			b.Fatalf("re-walk sum = %d, want %d", got, want)
		}
	}
}

// BenchmarkLazySeqFirstRealized measures the single-accessor cost on a
// realized lazy seq (the get-one-element case).
func BenchmarkLazySeqFirstRealized(b *testing.B) {
	ls := lazyChain(4)
	ls.Resolve()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if ls.First() != Int(0) {
			b.Fatal("bad First")
		}
	}
}
