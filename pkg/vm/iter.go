/*
 * SPDX-License-Identifier: MIT
 */

package vm

import "iter"

// SeqIsEmpty reports whether a Seq is exhausted. Encapsulates the two
// sentinels — Go nil and the EmptyList singleton — that let-go uses
// interchangeably to signal end-of-sequence.
//
// Use this to avoid the common `for s != nil && s != EmptyList` pattern.
func SeqIsEmpty(s Seq) bool {
	return s == nil || s == EmptyList
}

// SeqValues returns an iter.Seq[Value] that walks a Seq from a known
// Seq starting point. Use when you already have a Seq; use Iter when
// you have a Value (which may be Sequable or non-seq).
func SeqValues(s Seq) iter.Seq[Value] {
	return func(yield func(Value) bool) {
		for !SeqIsEmpty(s) {
			if !yield(s.First()) {
				return
			}
			s = s.Next()
		}
	}
}

// Iter returns a Go-1.23+ iter.Seq[Value] over any let-go sequence value.
// Yields nothing if v is nil/NIL or isn't a sequence.
//
// Preserves laziness: the underlying LazySeq thunks are realized one step at
// a time as the consumer pulls values, and an early break stops realization.
//
// Idiomatic use:
//
//	for v := range vm.Iter(result) {
//	    // process v
//	}
//
// Composes with slices.Collect, slices.Filter, etc.
func Iter(v Value) iter.Seq[Value] {
	return func(yield func(Value) bool) {
		s := seqFromValue(v)
		if s == nil {
			return
		}
		for !SeqIsEmpty(s) {
			if !yield(s.First()) {
				return
			}
			s = s.Next()
		}
	}
}

// SeqToSlice realizes a let-go sequence value into a []Value. Useful when a
// Go caller has a vm.Value and explicitly wants a slice. Returns nil for
// nil/NIL/empty input. Returns an error if v isn't a sequence.
//
// For infinite sequences this will not terminate — use Iter and break instead.
func SeqToSlice(v Value) ([]Value, error) {
	if v == nil || v == NIL {
		return nil, nil
	}
	s := seqFromValue(v)
	if s == nil {
		return nil, NewTypeError(v, "is not a sequence", ListType)
	}
	var out []Value
	for !SeqIsEmpty(s) {
		out = append(out, s.First())
		s = s.Next()
	}
	return out, nil
}

// seqFromValue extracts a Seq from a Value, accepting either Sequable
// (the common case) or a value that itself implements Seq. Returns nil
// if v is nil, NIL, or neither.
func seqFromValue(v Value) Seq {
	if v == nil || v == NIL {
		return nil
	}
	if sq, ok := v.(Sequable); ok {
		return sq.Seq()
	}
	if s, ok := v.(Seq); ok {
		return s
	}
	return nil
}
