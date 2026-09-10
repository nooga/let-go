/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package vm

import "testing"

// Sets must never inherit the map's insertion-ordered array mode: Clojure has
// no array-set contract, and an ordered impl made small #{...} literals
// sequence in source order (reversed, even) while (set [...]) sequenced
// forward. Every construction path has to land on a HAMT-backed impl so set
// iteration order depends on element hashes alone. White-box on impl.ordered();
// the observable end-to-end shape lives in test/set_order_test.lg.
func TestSetsStayHAMTBacked(t *testing.T) {
	vals := []Value{Keyword("a"), Keyword("b"), Keyword("c")}

	assertHAMT := func(name string, s *PersistentSet) {
		t.Helper()
		if s.impl.ordered() {
			t.Errorf("%s: set impl is in ordered array-map mode; sets must stay HAMT-backed", name)
		}
	}

	assertHAMT("EmptyPersistentSet", EmptyPersistentSet)
	assertHAMT("NewPersistentSet", NewPersistentSet(vals))

	conj := EmptyPersistentSet
	for _, v := range vals {
		conj = conj.Conj(v).(*PersistentSet)
	}
	assertHAMT("Conj chain from empty", conj)

	ts := NewTransientSet(EmptyPersistentSet)
	for _, v := range vals {
		if _, err := ts.Conj(v); err != nil {
			t.Fatal(err)
		}
	}
	viaTransient, err := ts.Persistent()
	if err != nil {
		t.Fatal(err)
	}
	assertHAMT("transient round-trip", viaTransient)

	empty, err := NewTransientSet(EmptyPersistentSet).Persistent()
	if err != nil {
		t.Fatal(err)
	}
	assertHAMT("empty transient round-trip", empty)
	assertHAMT("re-grown from empty transient", empty.Conj(Keyword("x")).(*PersistentSet))

	drained := NewPersistentSet(vals)
	for _, v := range vals {
		drained = drained.Disj(v)
	}
	assertHAMT("drained to empty via Disj", drained)
	assertHAMT("re-grown after drain", drained.Conj(Keyword("x")).(*PersistentSet))
}

// Iteration order of equal sets is identical no matter how they were built —
// the observable consequence of the invariant above.
func TestSetIterationOrderIsInsertionIndependent(t *testing.T) {
	fwd := NewPersistentSet([]Value{Keyword("a"), Keyword("b"), Keyword("c")})
	rev := NewPersistentSet([]Value{Keyword("c"), Keyword("b"), Keyword("a")})
	a, b := fwd.Seq(), rev.Seq()
	for a != nil && a != EmptyList && b != nil && b != EmptyList {
		if !valueEquiv(a.First(), b.First()) {
			t.Fatalf("iteration order differs between insertion orders: %v vs %v", a.First(), b.First())
		}
		a, b = a.Next(), b.Next()
	}
	if (a != nil && a != EmptyList) != (b != nil && b != EmptyList) {
		t.Fatal("sequences have different lengths")
	}
}
