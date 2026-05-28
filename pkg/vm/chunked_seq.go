/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package vm

// AsChunkedSeq returns the chunked view of s if available, resolving through
// LazySeq wrappers. Returns (cs, true) if s ultimately resolves to a non-empty
// IChunkedSeq, (nil, false) otherwise (including for empty seqs).
//
// Callers use this to choose a chunked fast path:
//
//	if cs, ok := AsChunkedSeq(s); ok {
//	    // process cs.ChunkedFirst() in bulk, recurse on cs.ChunkedMore()
//	}
//	// fallback: one-at-a-time
func AsChunkedSeq(s Seq) (IChunkedSeq, bool) {
	for {
		if s == nil || s == EmptyList {
			return nil, false
		}
		if ls, ok := s.(*LazySeq); ok {
			s = ls.Resolve()
			continue
		}
		cs, ok := s.(IChunkedSeq)
		return cs, ok
	}
}
