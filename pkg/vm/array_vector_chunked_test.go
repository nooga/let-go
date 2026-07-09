package vm

import (
	"sync"
	"testing"
)

// Compile-time contract: *ArrayVectorSeq must satisfy IChunkedSeq.
var _ IChunkedSeq = (*ArrayVectorSeq)(nil)

// avVec builds an ArrayVector [0,1,...,n-1] as Int values.
func avVec(n int) ArrayVector {
	vec := make(ArrayVector, n)
	for i := 0; i < n; i++ {
		vec[i] = Int(i)
	}
	return vec
}

// The HEAD chunk is size 1 — this is what keeps first/short-circuit
// consumption element-wise (no over-realization).
func TestArrayVectorSeqChunkedFirstIsSizeOne(t *testing.T) {
	s := &ArrayVectorSeq{vec: avVec(100), i: 0}
	c := s.ChunkedFirst()
	if c.ChunkCount() != 1 {
		t.Fatalf("head ChunkCount = %d, want 1", c.ChunkCount())
	}
	if c.Nth(0) != Int(0) {
		t.Fatalf("Nth(0) = %v, want 0", c.Nth(0))
	}
}

// Chunk sizes grow geometrically from 1 (1, 2, 4, 8, …) as the walk advances.
func TestArrayVectorSeqChunkSizesDouble(t *testing.T) {
	var got []int
	var start []int
	s := &ArrayVectorSeq{vec: avVec(100), i: 0}
	for s != nil {
		c := s.ChunkedFirst()
		got = append(got, c.ChunkCount())
		start = append(start, s.i)
		nx := s.ChunkedNext()
		if nx == nil {
			break
		}
		s = nx.(*ArrayVectorSeq)
	}
	wantSizes := []int{1, 2, 4, 8, 16, 32, 37} // 1+2+4+8+16+32 = 63, trailing = 100-63 = 37
	wantStart := []int{0, 1, 3, 7, 15, 31, 63}
	if len(got) != len(wantSizes) {
		t.Fatalf("chunk count = %d %v, want %d %v", len(got), got, len(wantSizes), wantSizes)
	}
	for i := range wantSizes {
		if got[i] != wantSizes[i] || start[i] != wantStart[i] {
			t.Fatalf("chunk %d: size=%d start=%d, want size=%d start=%d", i, got[i], start[i], wantSizes[i], wantStart[i])
		}
	}
}

// The trailing chunk is clamped to the real remainder, never the full quantum.
func TestArrayVectorSeqTrailingChunkClamped(t *testing.T) {
	// size 5: chunks are 1 (@0), 2 (@1), then quantum wants 4 but only 2 remain.
	s := &ArrayVectorSeq{vec: avVec(5), i: 0}
	if c := s.ChunkedFirst(); c.ChunkCount() != 1 {
		t.Fatalf("chunk0 = %d, want 1", c.ChunkCount())
	}
	s2 := s.ChunkedNext().(*ArrayVectorSeq)
	if c := s2.ChunkedFirst(); c.ChunkCount() != 2 {
		t.Fatalf("chunk1 = %d, want 2", c.ChunkCount())
	}
	s3 := s2.ChunkedNext().(*ArrayVectorSeq)
	if c := s3.ChunkedFirst(); c.ChunkCount() != 2 { // quantum 4 clamped to remaining 2
		t.Fatalf("trailing chunk = %d, want 2 (clamped from quantum 4)", c.ChunkCount())
	}
	if s3.ChunkedNext() != nil {
		t.Fatalf("ChunkedNext past end = %v, want nil", s3.ChunkedNext())
	}
	if s3.ChunkedMore() != EmptyList {
		t.Fatalf("ChunkedMore past end = %v, want EmptyList", s3.ChunkedMore())
	}
}

// elementWalk collects First()/Next() over a seq.
func elementWalk(s Seq) []Value {
	var out []Value
	for s != nil {
		out = append(out, s.First())
		s = s.Next()
	}
	return out
}

// chunkedWalk collects ChunkedFirst()/ChunkedNext() over an IChunkedSeq.
func chunkedWalk(cs IChunkedSeq) []Value {
	var out []Value
	for {
		c := cs.ChunkedFirst()
		for j := 0; j < c.ChunkCount(); j++ {
			out = append(out, c.Nth(j))
		}
		nx := cs.ChunkedNext()
		if nx == nil {
			break
		}
		cs = nx.(IChunkedSeq)
	}
	return out
}

func TestArrayVectorSeqChunkedWalkMatchesElementWalk(t *testing.T) {
	for _, n := range []int{1, 31, 32, 33, 100} {
		for _, start := range []int{0, 1, 32} {
			if start >= n {
				continue
			}
			elem := elementWalk(&ArrayVectorSeq{vec: avVec(n), i: start})
			chunked := chunkedWalk(&ArrayVectorSeq{vec: avVec(n), i: start})
			if len(elem) != len(chunked) {
				t.Fatalf("n=%d start=%d: len elem=%d chunked=%d", n, start, len(elem), len(chunked))
			}
			for k := range elem {
				if elem[k] != chunked[k] {
					t.Fatalf("n=%d start=%d idx=%d: elem=%v chunked=%v", n, start, k, elem[k], chunked[k])
				}
			}
		}
	}
}

func TestArrayVectorSeqSeqMethodsUnchanged(t *testing.T) {
	// Regression guard: the per-element Seq contract is untouched.
	s := &ArrayVectorSeq{vec: avVec(3), i: 0}
	if s.First() != Int(0) {
		t.Fatalf("First = %v, want 0", s.First())
	}
	if int(s.Count().(Int)) != 3 {
		t.Fatalf("Count = %v, want 3", s.Count())
	}
	if s.String() != "(0 1 2)" {
		t.Fatalf("String = %q, want (0 1 2)", s.String())
	}
	if s.Nth(2) != Int(2) {
		t.Fatalf("Nth(2) = %v, want 2", s.Nth(2))
	}
}

// TestArrayVectorSeqChunkedFirstRace guards the persistent-collection
// thread-safety invariant: an *ArrayVectorSeq node must stay read/share-safe,
// so ChunkedFirst returns a FRESH chunk rather than mutating shared node state.
// Run under `go test -race`. A prior revision reused an embedded s.chunk field
// and raced here — two goroutines calling ChunkedFirst on the same node wrote
// it concurrently (reported at vector.go ChunkedFirst / chunk.go Nth).
func TestArrayVectorSeqChunkedFirstRace(t *testing.T) {
	s := &ArrayVectorSeq{vec: avVec(64), i: 0}
	const goroutines = 8
	const iters = 2000
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for g := 0; g < goroutines; g++ {
		go func() {
			defer wg.Done()
			for i := 0; i < iters; i++ {
				if v := s.ChunkedFirst().Nth(0); v != Int(0) {
					t.Errorf("ChunkedFirst().Nth(0) = %v, want 0", v)
					return
				}
			}
		}()
	}
	wg.Wait()
}
