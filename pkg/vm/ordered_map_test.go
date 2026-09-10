package vm

import "testing"

func orderedKeys(m *PersistentMap) []Value {
	var ks []Value
	for s := m.Seq(); s != nil && s != Seq(EmptyList); s = s.Next() {
		k, _, ok := MapEntryKV(s.First())
		if !ok {
			break
		}
		ks = append(ks, k)
	}
	return ks
}

func wantKeys(t *testing.T, m *PersistentMap, want ...Value) {
	t.Helper()
	got := orderedKeys(m)
	if len(got) != len(want) {
		t.Fatalf("key count = %d, want %d (%v)", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("keys = %v, want %v", got, want)
		}
	}
}

// Small maps (≤8 entries) preserve insertion order, like Clojure's
// PersistentArrayMap behind map literals and array-map.
func TestOrderedMapLiteralOrder(t *testing.T) {
	m := NewArrayMap([]Value{Keyword("_id"), Int(1), Keyword("name"), String("cat")})
	wantKeys(t, m, Keyword("_id"), Keyword("name"))

	m3 := NewArrayMap([]Value{Keyword("a"), Int(1), Keyword("b"), Int(2), Keyword("c"), Int(3)})
	wantKeys(t, m3, Keyword("a"), Keyword("b"), Keyword("c"))
}

func TestOrderedMapAssocChain(t *testing.T) {
	keys := []Keyword{"h", "g", "f", "e", "d", "c", "b", "a"}
	m := EmptyPersistentMap
	for i, k := range keys {
		m = m.Assoc(k, Int(i)).(*PersistentMap)
	}
	want := make([]Value, len(keys))
	for i, k := range keys {
		want[i] = k
	}
	wantKeys(t, m, want...)

	// Replacing an existing key keeps its position.
	m2 := m.Assoc(Keyword("f"), Int(99)).(*PersistentMap)
	wantKeys(t, m2, want...)
	if m2.ValueAt(Keyword("f")) != Int(99) {
		t.Fatalf("replace failed: %v", m2.ValueAt(Keyword("f")))
	}
}

func TestOrderedMapDissocKeepsOrder(t *testing.T) {
	m := NewArrayMap([]Value{Keyword("a"), Int(1), Keyword("b"), Int(2), Keyword("c"), Int(3)})
	m2 := m.Dissoc(Keyword("b")).(*PersistentMap)
	wantKeys(t, m2, Keyword("a"), Keyword("c"))
	if m2.RawCount() != 2 {
		t.Fatalf("count = %d", m2.RawCount())
	}
}

// The 9th entry promotes to the hash representation: order is no longer
// guaranteed, but lookup/count semantics are unchanged and no demotion
// happens on dissoc back below the threshold (matching Clojure).
func TestOrderedMapPromotionBoundary(t *testing.T) {
	m := EmptyPersistentMap
	ks := []Keyword{"a", "b", "c", "d", "e", "f", "g", "h"}
	for i, k := range ks {
		m = m.Assoc(k, Int(i)).(*PersistentMap)
	}
	if m.root != nil {
		t.Fatal("8-entry map must still be in ordered mode")
	}
	m9 := m.Assoc(Keyword("i"), Int(8)).(*PersistentMap)
	if m9.root == nil {
		t.Fatal("9-entry map must be promoted to the hash representation")
	}
	if m9.RawCount() != 9 {
		t.Fatalf("count = %d", m9.RawCount())
	}
	for i, k := range ks {
		if m9.ValueAt(k) != Int(i) {
			t.Fatalf("lost %v after promotion", k)
		}
	}
	back := m9.Dissoc(Keyword("i")).(*PersistentMap)
	if back.root == nil {
		t.Fatal("dissoc below the threshold must not demote")
	}
	if back.RawCount() != 8 {
		t.Fatalf("count = %d", back.RawCount())
	}
	for i, k := range ks {
		if back.ValueAt(k) != Int(i) {
			t.Fatalf("lost %v after dissoc", k)
		}
	}
}

// Draining a promoted map to empty deliberately lands back in the ordered
// mode: an empty map has no order, and rebuilding from it behaves exactly
// like building from a fresh {}.
func TestOrderedMapDrainedHAMTRebuildsOrdered(t *testing.T) {
	m := EmptyPersistentMap
	ks := []Keyword{"a", "b", "c", "d", "e", "f", "g", "h", "i"}
	for i, k := range ks {
		m = m.Assoc(k, Int(i)).(*PersistentMap)
	}
	if m.root == nil {
		t.Fatal("expected HAMT mode at 9 entries")
	}
	for _, k := range ks {
		m = m.Dissoc(k).(*PersistentMap)
	}
	if m.RawCount() != 0 || m.root != nil {
		t.Fatalf("drained map must be the empty ordered mode (count=%d)", m.RawCount())
	}
	re := m.Assoc(Keyword("z"), Int(1)).(*PersistentMap).Assoc(Keyword("y"), Int(2)).(*PersistentMap)
	wantKeys(t, re, Keyword("z"), Keyword("y"))
}

// Equality and hash stay order-insensitive and representation-insensitive.
func TestOrderedMapEqualityAcrossOrderAndMode(t *testing.T) {
	a := NewArrayMap([]Value{Keyword("a"), Int(1), Keyword("b"), Int(2)})
	b := NewArrayMap([]Value{Keyword("b"), Int(2), Keyword("a"), Int(1)})
	if !a.Equals(b) || !b.Equals(a) {
		t.Fatal("insertion order must not affect equality")
	}
	if a.Hash() != b.Hash() {
		t.Fatal("insertion order must not affect hash")
	}

	// Same 8 entries, one ordered, one HAMT (via promote + dissoc).
	ks := []Keyword{"a", "b", "c", "d", "e", "f", "g", "h"}
	ordered := EmptyPersistentMap
	for i, k := range ks {
		ordered = ordered.Assoc(k, Int(i)).(*PersistentMap)
	}
	hashed := ordered.Assoc(Keyword("i"), Int(8)).(*PersistentMap).Dissoc(Keyword("i")).(*PersistentMap)
	if hashed.root == nil {
		t.Fatal("expected HAMT-mode comparand")
	}
	if !ordered.Equals(hashed) || !hashed.Equals(ordered) {
		t.Fatal("representation must not affect equality")
	}
	if ordered.Hash() != hashed.Hash() {
		t.Fatal("representation must not affect hash")
	}
}

// Transients preserve insertion order for small maps (the into path) and
// promote past 8 like the persistent Assoc.
func TestOrderedMapTransient(t *testing.T) {
	tm := NewTransientMap(EmptyPersistentMap)
	var err error
	for _, k := range []Keyword{"z", "y", "x"} {
		if tm, err = tm.Assoc(k, k); err != nil {
			t.Fatal(err)
		}
	}
	m, err := tm.Persistent()
	if err != nil {
		t.Fatal(err)
	}
	wantKeys(t, m, Keyword("z"), Keyword("y"), Keyword("x"))

	tm2 := NewTransientMap(EmptyPersistentMap)
	for i := 0; i < 10; i++ {
		if tm2, err = tm2.Assoc(Int(i), Int(i)); err != nil {
			t.Fatal(err)
		}
	}
	big, err := tm2.Persistent()
	if err != nil {
		t.Fatal(err)
	}
	if big.RawCount() != 10 {
		t.Fatalf("count = %d", big.RawCount())
	}
	for i := 0; i < 10; i++ {
		if big.ValueAt(Int(i)) != Int(i) {
			t.Fatalf("lost %d after transient promotion", i)
		}
	}
}
