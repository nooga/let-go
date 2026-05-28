/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

// incFn returns a vm.Fn that adds 1 to its single Int argument.
func incFn() vm.Fn {
	f, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		return vs[0].(vm.Int) + 1, nil
	})
	return f.(vm.Fn)
}

// addFn returns a 2-arg vm.Fn that adds two Ints.
func addFn() vm.Fn {
	f, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		return vs[0].(vm.Int) + vs[1].(vm.Int), nil
	})
	return f.(vm.Fn)
}

func TestMapLazy1PreservesChunkednessFromRange(t *testing.T) {
	r := vm.NewRange(vm.Int(0), vm.Int(100), vm.Int(1)).(vm.Seq)
	out := mapLazy1(incFn(), r)
	// out is a LazySeq; its resolved value should be a ChunkedCons.
	cs, ok := vm.AsChunkedSeq(out)
	if !ok {
		t.Fatalf("map over chunked range should produce IChunkedSeq, got %T", out)
	}
	first := cs.ChunkedFirst()
	if first.ChunkCount() != 32 {
		t.Fatalf("first chunk count: got %d, want 32", first.ChunkCount())
	}
	if first.Nth(0) != vm.Int(1) || first.Nth(31) != vm.Int(32) {
		t.Fatalf("first chunk values: 0->%v 31->%v", first.Nth(0), first.Nth(31))
	}
}

func TestMapLazy1WalkAcrossChunks(t *testing.T) {
	r := vm.NewRange(vm.Int(0), vm.Int(100), vm.Int(1)).(vm.Seq)
	out := mapLazy1(incFn(), r)
	got := []int{}
	for v := range vm.Iter(out) {
		got = append(got, int(v.(vm.Int)))
	}
	if len(got) != 100 {
		t.Fatalf("len: got %d, want 100", len(got))
	}
	for i, v := range got {
		if v != i+1 {
			t.Fatalf("[%d]: got %d, want %d", i, v, i+1)
		}
	}
}

func TestMapLazy1NonChunkedInputFallback(t *testing.T) {
	// Build a Cons chain (not chunked) and confirm map still works.
	var s vm.Seq = nil
	for i := 4; i >= 0; i-- {
		s = vm.NewCons(vm.Int(i), s)
	}
	out := mapLazy1(incFn(), s)
	got := []int{}
	for v := range vm.Iter(out) {
		got = append(got, int(v.(vm.Int)))
	}
	want := []int{1, 2, 3, 4, 5}
	if len(got) != len(want) {
		t.Fatalf("len: got %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("[%d]: got %d, want %d", i, got[i], want[i])
		}
	}
}

func TestMapLazyNChunkedFastPath(t *testing.T) {
	r1 := vm.NewRange(vm.Int(0), vm.Int(50), vm.Int(1)).(vm.Seq)
	r2 := vm.NewRange(vm.Int(100), vm.Int(150), vm.Int(1)).(vm.Seq)
	out := mapLazyN(addFn(), []vm.Seq{r1, r2})
	cs, ok := vm.AsChunkedSeq(out)
	if !ok {
		t.Fatalf("mapN over two chunked ranges should be chunked, got %T", out)
	}
	first := cs.ChunkedFirst()
	if first.ChunkCount() != 32 {
		t.Fatalf("first chunk count: got %d, want 32", first.ChunkCount())
	}
	if first.Nth(0) != vm.Int(100) {
		t.Fatalf("first chunk[0]: got %v, want 100", first.Nth(0))
	}
}

func TestMapLazyNWalkAcrossChunks(t *testing.T) {
	r1 := vm.NewRange(vm.Int(0), vm.Int(100), vm.Int(1)).(vm.Seq)
	r2 := vm.NewRange(vm.Int(0), vm.Int(100), vm.Int(2)).(vm.Seq) // 50 elements
	out := mapLazyN(addFn(), []vm.Seq{r1, r2})
	got := []int{}
	for v := range vm.Iter(out) {
		got = append(got, int(v.(vm.Int)))
	}
	// Shortest wins: 50 results, got[i] = i + 2i = 3i
	if len(got) != 50 {
		t.Fatalf("len: got %d, want 50", len(got))
	}
	for i, v := range got {
		if v != 3*i {
			t.Fatalf("[%d]: got %d, want %d", i, v, 3*i)
		}
	}
}

func TestMapLazyNMixedChunkedAndNonChunkedFallsBack(t *testing.T) {
	// First input is chunked; second is a Cons chain. The chunkedHeads
	// guard requires all-or-nothing, so this falls back to element-wise
	// map. Result still must be correct.
	r := vm.NewRange(vm.Int(0), vm.Int(5), vm.Int(1)).(vm.Seq)
	var c vm.Seq = nil
	for i := 4; i >= 0; i-- {
		c = vm.NewCons(vm.Int(i*10), c)
	}
	out := mapLazyN(addFn(), []vm.Seq{r, c})
	got := []int{}
	for v := range vm.Iter(out) {
		got = append(got, int(v.(vm.Int)))
	}
	want := []int{0, 11, 22, 33, 44}
	if len(got) != len(want) {
		t.Fatalf("len: got %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("[%d]: got %d, want %d", i, got[i], want[i])
		}
	}
}
