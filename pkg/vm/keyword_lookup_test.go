package vm

import "testing"

// Keyword("x").Hash() must equal hashValue(boxed) — the typed path relies on it.
func TestKeywordHashMatchesHashValue(t *testing.T) {
	k := Keyword("some-field")
	if k.Hash() != hashValue(k) {
		t.Fatalf("k.Hash()=%d != hashValue(k)=%d", k.Hash(), hashValue(k))
	}
}

// ValueAtKeyword must be identical to ValueAt for a keyword key, across:
// present, missing, mixed keyword/non-keyword maps.
func TestPersistentMapValueAtKeywordMatchesValueAt(t *testing.T) {
	m := EmptyPersistentMap.
		Assoc(Keyword("a"), Int(1)).
		Assoc(Keyword("b"), Int(2)).
		Assoc(Int(42), Int(99)).                    // non-keyword key
		Assoc(String("s"), Int(7)).(*PersistentMap) // non-keyword key
	cases := []Keyword{"a", "b", "missing"}
	for _, k := range cases {
		got := m.ValueAtKeyword(k)
		want := m.ValueAt(k)
		if got != want {
			t.Fatalf("ValueAtKeyword(%q)=%v, ValueAt=%v", k, got, want)
		}
	}
	// default variant
	if got := m.ValueAtKeywordOr("missing", Int(-1)); got != m.ValueAtOr(Keyword("missing"), Int(-1)) {
		t.Fatalf("ValueAtKeywordOr default mismatch: %v", got)
	}
}

func TestRecordValueAtKeywordMatchesValueAt(t *testing.T) {
	rt := NewRecordType("Point", []Keyword{"x", "y"})
	data := EmptyPersistentMap.Assoc(Keyword("x"), Int(3)).Assoc(Keyword("y"), Int(4)).(*PersistentMap)
	r := NewRecord(rt, data) // NewRecord(rtype *RecordType, data *PersistentMap) *Record
	for _, k := range []Keyword{"x", "y", "z"} {
		if got, want := r.ValueAtKeyword(k), r.ValueAt(k); got != want {
			t.Fatalf("Record.ValueAtKeyword(%q)=%v, ValueAt=%v", k, got, want)
		}
	}
	if got, want := r.ValueAtKeywordOr("z", Int(-1)), r.ValueAtOr(Keyword("z"), Int(-1)); got != want {
		t.Fatalf("Record.ValueAtKeywordOr default: %v != %v", got, want)
	}
}

// A keyword and a non-keyword key forced into the same bucket must not
// cross-match. We can't easily force a raw hash collision, so assert the
// weaker invariant that a non-keyword key with the same STRING content as a
// keyword is not returned by ValueAtKeyword.
func TestValueAtKeywordIgnoresNonKeywordKeys(t *testing.T) {
	m := EmptyPersistentMap.Assoc(String("x"), Int(1)).(*PersistentMap) // string key "x"
	if got := m.ValueAtKeyword(Keyword("x")); got != NIL {
		t.Fatalf("keyword :x matched a string key \"x\": %v", got)
	}
	if got := m.ValueAt(Keyword("x")); got != NIL {
		t.Fatalf("sanity: ValueAt(:x) should also be NIL: %v", got)
	}
}

func TestKeywordInvokeUsesFastPathAndFallback(t *testing.T) {
	m := EmptyPersistentMap.Assoc(Keyword("a"), Int(1)).(*PersistentMap)
	// fast path (map implements KeywordLookup)
	if v, _ := Keyword("a").Invoke([]Value{m}); v != Int(1) {
		t.Fatalf("fast path (:a m) = %v, want 1", v)
	}
	if v, _ := Keyword("z").Invoke([]Value{m, Int(-1)}); v != Int(-1) {
		t.Fatalf("fast path (:z m -1) = %v, want -1", v)
	}
	// fallback path: a Lookup that does NOT implement KeywordLookup.
	// A *List implements Lookup (index lookup) but not KeywordLookup; use it to
	// exercise the fallback branch without regressing behavior.
	lst := NewList([]Value{Int(10), Int(20)})
	if v, _ := Keyword("a").Invoke([]Value{lst, Int(-1)}); v != Int(-1) {
		t.Fatalf("fallback (:a list -1) = %v, want -1 (keyword not a list key)", v)
	}
}

func BenchmarkKeywordInvoke(b *testing.B) {
	m := EmptyPersistentMap.Assoc(Keyword("a"), Int(1)).(*PersistentMap)
	args := []Value{m}
	kw := Keyword("a")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = kw.Invoke(args)
	}
}
