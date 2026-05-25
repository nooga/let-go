/*
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"iter"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIter_OverList(t *testing.T) {
	l, _ := ListType.Box([]Value{Int(1), Int(2), Int(3)})
	var got []Value
	for v := range Iter(l) {
		got = append(got, v)
	}
	assert.Equal(t, []Value{Int(1), Int(2), Int(3)}, got)
}

func TestIter_OverArrayVector(t *testing.T) {
	v := ArrayVector{Int(10), Int(20), Int(30)}
	var got []Value
	for x := range Iter(v) {
		got = append(got, x)
	}
	assert.Equal(t, []Value{Int(10), Int(20), Int(30)}, got)
}

func TestIter_OverNil(t *testing.T) {
	count := 0
	for range Iter(nil) {
		count++
	}
	assert.Equal(t, 0, count)
}

func TestIter_OverNILValue(t *testing.T) {
	count := 0
	for range Iter(NIL) {
		count++
	}
	assert.Equal(t, 0, count)
}

func TestIter_OverEmptyList(t *testing.T) {
	count := 0
	for range Iter(EmptyList) {
		count++
	}
	assert.Equal(t, 0, count)
}

func TestIter_OverNonSeqYieldsNothing(t *testing.T) {
	count := 0
	for range Iter(Int(42)) {
		count++
	}
	assert.Equal(t, 0, count)
}

func TestIter_EarlyBreakPreservesLaziness(t *testing.T) {
	// LazySeq counting realizations: yield element n, advance via a thunk
	// that increments a counter and returns the next chunk.
	realized := 0
	var build func(n int) Seq
	build = func(n int) Seq {
		thunk, _ := NativeFnType.Wrap(func(_ []Value) (Value, error) {
			realized++
			return NewCons(Int(int64(n)), build(n+1)), nil
		})
		return NewLazySeq(thunk.(Fn))
	}
	seq := build(0)

	pulled := 0
	for v := range Iter(seq) {
		pulled++
		if pulled == 3 {
			break
		}
		_ = v
	}
	assert.Equal(t, 3, pulled, "should pull exactly 3 elements")
	// Realization happens on demand; break stops further realization.
	// Allow some over-realization (Iter resolves one extra to detect end-of-step),
	// but it must be far short of running forever.
	assert.Less(t, realized, 10, "early break should stop infinite realization")
}

func TestSeqToSlice_OverList(t *testing.T) {
	l, _ := ListType.Box([]Value{Int(1), Int(2), Int(3)})
	got, err := SeqToSlice(l)
	assert.NoError(t, err)
	assert.Equal(t, []Value{Int(1), Int(2), Int(3)}, got)
}

func TestSeqToSlice_OverArrayVector(t *testing.T) {
	v := ArrayVector{Keyword("a"), Keyword("b")}
	got, err := SeqToSlice(v)
	assert.NoError(t, err)
	assert.Equal(t, []Value{Keyword("a"), Keyword("b")}, got)
}

func TestSeqToSlice_OverNil(t *testing.T) {
	got, err := SeqToSlice(nil)
	assert.NoError(t, err)
	assert.Nil(t, got)
}

func TestSeqToSlice_OverNILValue(t *testing.T) {
	got, err := SeqToSlice(NIL)
	assert.NoError(t, err)
	assert.Nil(t, got)
}

func TestSeqToSlice_OverEmptyList(t *testing.T) {
	got, err := SeqToSlice(EmptyList)
	assert.NoError(t, err)
	assert.Nil(t, got)
}

func TestSeqToSlice_OverNonSeqErrors(t *testing.T) {
	_, err := SeqToSlice(Int(42))
	assert.Error(t, err)
}

func TestSeqIsEmpty(t *testing.T) {
	// Nil and EmptyList both report as empty.
	assert.True(t, SeqIsEmpty(nil))
	assert.True(t, SeqIsEmpty(EmptyList))
	// A non-empty list is not empty.
	l, _ := ListType.Box([]Value{Int(1)})
	assert.False(t, SeqIsEmpty(l.(Seq)))
}

func TestSeqValues_WalksFromSeq(t *testing.T) {
	l, _ := ListType.Box([]Value{Int(1), Int(2), Int(3)})
	var got []Value
	for v := range SeqValues(l.(Seq)) {
		got = append(got, v)
	}
	assert.Equal(t, []Value{Int(1), Int(2), Int(3)}, got)
}

func TestSeqValues_HandlesEmptySentinels(t *testing.T) {
	assert.Equal(t, 0, countValues(SeqValues(nil)))
	assert.Equal(t, 0, countValues(SeqValues(EmptyList)))
}

func countValues(seq iter.Seq[Value]) int {
	n := 0
	for range seq {
		n++
	}
	return n
}
