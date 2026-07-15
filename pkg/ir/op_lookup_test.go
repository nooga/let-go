/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package ir_test

import (
	"testing"

	"github.com/nooga/let-go/pkg/ir"
)

// Every catalogued op must round-trip through its canonical kebab-case
// keyword (the form ir.build emits and OpKeywords returns) back to its Op.
// This proves the generated opByKeywordExact switch covers the whole
// catalog and agrees with OpKeywords index-for-index — if op-keyword-name
// and the switch ever drift, this fails.
func TestOpByKeywordRoundTrip(t *testing.T) {
	for i, name := range ir.OpKeywords() {
		op, ok := ir.OpByKeyword(name)
		if !ok {
			t.Errorf("OpByKeyword(%q) = _, false; want the op at index %d", name, i)
			continue
		}
		if int(op) != i {
			t.Errorf("OpByKeyword(%q) = %d; want %d", name, int(op), i)
		}
	}
}

// The canonical fast path (the generated string switch) must resolve
// without falling through to the ReplaceAll + scan fallback, so it does
// not allocate. This pins the reason the switch exists.
func TestOpByKeywordCanonicalIsAllocFree(t *testing.T) {
	var sink ir.Op
	allocs := testing.AllocsPerRun(1000, func() {
		sink, _ = ir.OpByKeyword("load-arg")
	})
	if allocs != 0 {
		t.Errorf("OpByKeyword(\"load-arg\") allocated %v times; want 0 (canonical input must hit the switch, not the fallback scan)", allocs)
	}
	_ = sink
}

// The dash/case-insensitive fallback semantics are preserved for
// non-canonical spellings, and an unknown keyword still reports absence.
func TestOpByKeywordFallbackSemantics(t *testing.T) {
	canonical, ok := ir.OpByKeyword("load-arg")
	if !ok {
		t.Fatal("OpByKeyword(\"load-arg\") should resolve")
	}
	for _, spelling := range []string{"loadarg", "LoadArg", "LOAD-ARG"} {
		got, ok := ir.OpByKeyword(spelling)
		if !ok || got != canonical {
			t.Errorf("OpByKeyword(%q) = %d, %v; want %d, true (lenient fallback)", spelling, int(got), ok, int(canonical))
		}
	}
	if op, ok := ir.OpByKeyword("no-such-op"); ok {
		t.Errorf("OpByKeyword(\"no-such-op\") = %d, true; want OpInvalid, false", int(op))
	}
}
