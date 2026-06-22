/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import "testing"

func mkVec(n int) PersistentVector {
	vals := make([]Value, n)
	for i := range vals {
		vals[i] = Int(i)
	}
	return NewPersistentVector(vals).(PersistentVector)
}

// TestPersistentVectorPopBoundaries pops one element at sizes that straddle the
// tail/trie boundary (32), the second trie level (1024), and height growth.
func TestPersistentVectorPopBoundaries(t *testing.T) {
	for _, n := range []int{2, 31, 32, 33, 34, 63, 64, 65, 1023, 1024, 1025, 2000} {
		v := mkVec(n)
		p := v.Pop()
		if p.RawCount() != n-1 {
			t.Fatalf("n=%d: popped count = %d, want %d", n, p.RawCount(), n-1)
		}
		for i := 0; i < n-1; i++ {
			if p.ValueAt(Int(i)) != Int(i) {
				t.Fatalf("n=%d: after pop, index %d = %v, want %d", n, i, p.ValueAt(Int(i)), i)
			}
		}
		// Original is unchanged (persistence).
		if v.RawCount() != n || v.ValueAt(Int(n-1)) != Int(n-1) {
			t.Fatalf("n=%d: pop mutated the original", n)
		}
	}
}

// TestPersistentVectorPopFullDrain drains a large vector one element at a time,
// verifying count and every remaining index at each step. This exercises
// popTail, leafArrayFor, every tail/trie boundary on the way down, and the
// root height-collapse path.
func TestPersistentVectorPopFullDrain(t *testing.T) {
	const n = 2000
	v := mkVec(n)
	for remaining := n; remaining > 0; remaining-- {
		if v.RawCount() != remaining {
			t.Fatalf("expected count %d, got %d", remaining, v.RawCount())
		}
		// Verify every remaining index — bulletproof against any popTail /
		// height-collapse / tail-pull mistake at any boundary on the way down.
		for i := 0; i < remaining; i++ {
			if v.ValueAt(Int(i)) != Int(i) {
				t.Fatalf("remaining=%d: index %d = %v, want %d", remaining, i, v.ValueAt(Int(i)), i)
			}
		}
		v = v.Pop()
	}
	if v.RawCount() != 0 {
		t.Fatalf("fully drained vector should be empty, got count %d", v.RawCount())
	}
}

func TestPersistentVectorPopToEmpty(t *testing.T) {
	v := mkVec(1)
	p := v.Pop()
	if p.RawCount() != 0 {
		t.Fatalf("pop of single-element vector should be empty, got count %d", p.RawCount())
	}
	// pop of empty returns empty (callers guard count==0 themselves).
	if p.Pop().RawCount() != 0 {
		t.Fatalf("pop of empty should stay empty")
	}
}
