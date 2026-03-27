package vm

import (
	"fmt"
	"testing"
)

func TestPersistentMapBasicPutGet(t *testing.T) {
	m := EmptyPersistentMap
	m2 := m.Assoc(Keyword("a"), Int(1)).(*PersistentMap)
	m3 := m2.Assoc(Keyword("b"), Int(2)).(*PersistentMap)

	if m3.ValueAt(Keyword("a")) != Int(1) {
		t.Errorf("expected 1, got %v", m3.ValueAt(Keyword("a")))
	}
	if m3.ValueAt(Keyword("b")) != Int(2) {
		t.Errorf("expected 2, got %v", m3.ValueAt(Keyword("b")))
	}
	if m3.RawCount() != 2 {
		t.Errorf("expected count 2, got %d", m3.RawCount())
	}
}

func TestPersistentMapOverwrite(t *testing.T) {
	m := EmptyPersistentMap
	m2 := m.Assoc(Keyword("a"), Int(1)).(*PersistentMap)
	m3 := m2.Assoc(Keyword("a"), Int(42)).(*PersistentMap)

	if m3.ValueAt(Keyword("a")) != Int(42) {
		t.Errorf("expected 42, got %v", m3.ValueAt(Keyword("a")))
	}
	if m3.RawCount() != 1 {
		t.Errorf("expected count 1, got %d", m3.RawCount())
	}
	// Original unchanged
	if m2.ValueAt(Keyword("a")) != Int(1) {
		t.Errorf("original should still have 1, got %v", m2.ValueAt(Keyword("a")))
	}
}

func TestPersistentMapDissoc(t *testing.T) {
	m := NewPersistentMap([]Value{Keyword("a"), Int(1), Keyword("b"), Int(2), Keyword("c"), Int(3)})
	m2 := m.Dissoc(Keyword("b")).(*PersistentMap)

	if m2.RawCount() != 2 {
		t.Errorf("expected count 2, got %d", m2.RawCount())
	}
	if m2.ValueAt(Keyword("a")) != Int(1) {
		t.Errorf("expected 1, got %v", m2.ValueAt(Keyword("a")))
	}
	if m2.ValueAt(Keyword("b")) != NIL {
		t.Errorf("expected nil, got %v", m2.ValueAt(Keyword("b")))
	}
	if m2.ValueAt(Keyword("c")) != Int(3) {
		t.Errorf("expected 3, got %v", m2.ValueAt(Keyword("c")))
	}
	// Original unchanged
	if m.RawCount() != 3 {
		t.Errorf("original should still have count 3, got %d", m.RawCount())
	}
}

func TestPersistentMapDissocToEmpty(t *testing.T) {
	m := NewPersistentMap([]Value{Keyword("a"), Int(1)})
	m2 := m.Dissoc(Keyword("a")).(*PersistentMap)
	if m2.RawCount() != 0 {
		t.Errorf("expected count 0, got %d", m2.RawCount())
	}
	if m2 != EmptyPersistentMap {
		t.Errorf("expected EmptyPersistentMap")
	}
}

func TestPersistentMapDissocMissing(t *testing.T) {
	m := NewPersistentMap([]Value{Keyword("a"), Int(1)})
	m2 := m.Dissoc(Keyword("nonexistent")).(*PersistentMap)
	if m2 != m {
		t.Errorf("dissoc of missing key should return same map")
	}
}

func TestPersistentMapEmpty(t *testing.T) {
	m := EmptyPersistentMap
	if m.RawCount() != 0 {
		t.Errorf("expected count 0, got %d", m.RawCount())
	}
	if m.ValueAt(Keyword("a")) != NIL {
		t.Errorf("expected nil, got %v", m.ValueAt(Keyword("a")))
	}
}

func TestPersistentMapLarge(t *testing.T) {
	m := EmptyPersistentMap
	n := 1000
	for i := 0; i < n; i++ {
		m = m.Assoc(Int(i), Int(i*10)).(*PersistentMap)
	}
	if m.RawCount() != n {
		t.Errorf("expected count %d, got %d", n, m.RawCount())
	}
	for i := 0; i < n; i++ {
		v := m.ValueAt(Int(i))
		if v != Int(i*10) {
			t.Errorf("at key %d: expected %d, got %v", i, i*10, v)
		}
	}
}

func TestPersistentMapImmutability(t *testing.T) {
	m1 := NewPersistentMap([]Value{Keyword("a"), Int(1)})
	m2 := m1.Assoc(Keyword("b"), Int(2)).(*PersistentMap)

	if m1.RawCount() != 1 {
		t.Errorf("m1 count should be 1, got %d", m1.RawCount())
	}
	if m2.RawCount() != 2 {
		t.Errorf("m2 count should be 2, got %d", m2.RawCount())
	}
	if m1.ValueAt(Keyword("b")) != NIL {
		t.Errorf("m1 should not have key :b")
	}
}

func TestPersistentMapVariousKeyTypes(t *testing.T) {
	m := EmptyPersistentMap
	m = m.Assoc(Int(42), String("int-key")).(*PersistentMap)
	m = m.Assoc(String("hello"), String("string-key")).(*PersistentMap)
	m = m.Assoc(Keyword("kw"), String("keyword-key")).(*PersistentMap)
	m = m.Assoc(Symbol("sym"), String("symbol-key")).(*PersistentMap)

	if m.ValueAt(Int(42)) != String("int-key") {
		t.Errorf("int key lookup failed")
	}
	if m.ValueAt(String("hello")) != String("string-key") {
		t.Errorf("string key lookup failed")
	}
	if m.ValueAt(Keyword("kw")) != String("keyword-key") {
		t.Errorf("keyword key lookup failed")
	}
	if m.ValueAt(Symbol("sym")) != String("symbol-key") {
		t.Errorf("symbol key lookup failed")
	}
}

func TestPersistentMapContains(t *testing.T) {
	m := NewPersistentMap([]Value{Keyword("a"), Int(1), Keyword("b"), Int(2)})
	if m.Contains(Keyword("a")) != TRUE {
		t.Errorf("should contain :a")
	}
	if m.Contains(Keyword("b")) != TRUE {
		t.Errorf("should contain :b")
	}
	if m.Contains(Keyword("c")) != FALSE {
		t.Errorf("should not contain :c")
	}
}

func TestPersistentMapConj(t *testing.T) {
	m := NewPersistentMap([]Value{Keyword("a"), Int(1)})
	m2 := m.Conj(ArrayVector{Keyword("b"), Int(2)}).(*PersistentMap)
	if m2.RawCount() != 2 {
		t.Errorf("expected count 2, got %d", m2.RawCount())
	}
	if m2.ValueAt(Keyword("b")) != Int(2) {
		t.Errorf("expected 2, got %v", m2.ValueAt(Keyword("b")))
	}
}

func TestPersistentMapAsFn(t *testing.T) {
	m := NewPersistentMap([]Value{Keyword("a"), Int(1), Keyword("b"), Int(2)})

	// 1-arity: lookup
	v, err := m.Invoke([]Value{Keyword("a")})
	if err != nil {
		t.Fatal(err)
	}
	if v != Int(1) {
		t.Errorf("expected 1, got %v", v)
	}

	// 2-arity: lookup with default
	v, err = m.Invoke([]Value{Keyword("c"), Int(99)})
	if err != nil {
		t.Fatal(err)
	}
	if v != Int(99) {
		t.Errorf("expected 99, got %v", v)
	}

	// Wrong arity
	_, err = m.Invoke([]Value{})
	if err == nil {
		t.Errorf("expected error for 0 args")
	}
}

func TestPersistentMapSeq(t *testing.T) {
	m := NewPersistentMap([]Value{Keyword("a"), Int(1), Keyword("b"), Int(2), Keyword("c"), Int(3)})
	seq := m.Seq()
	seen := make(map[string]bool)
	count := 0
	for seq != nil && seq != EmptyList {
		entry := seq.First().(ArrayVector)
		key := entry[0].(Keyword)
		seen[string(key)] = true
		count++
		seq = seq.Next()
	}
	if count != 3 {
		t.Errorf("expected 3 entries in seq, got %d", count)
	}
	for _, k := range []string{"a", "b", "c"} {
		if !seen[k] {
			t.Errorf("missing key :%s in seq", k)
		}
	}
}

func TestPersistentMapSeqEmpty(t *testing.T) {
	m := EmptyPersistentMap
	seq := m.Seq()
	if seq != EmptyList {
		t.Errorf("expected EmptyList for empty map seq")
	}
}

func TestPersistentMapEquals(t *testing.T) {
	m1 := NewPersistentMap([]Value{Keyword("a"), Int(1), Keyword("b"), Int(2)})
	m2 := NewPersistentMap([]Value{Keyword("b"), Int(2), Keyword("a"), Int(1)})
	m3 := NewPersistentMap([]Value{Keyword("a"), Int(1), Keyword("b"), Int(3)})

	if !m1.Equals(m2) {
		t.Errorf("m1 and m2 should be equal")
	}
	if m1.Equals(m3) {
		t.Errorf("m1 and m3 should not be equal")
	}
	if !EmptyPersistentMap.Equals(EmptyPersistentMap) {
		t.Errorf("empty maps should be equal")
	}
}

func TestPersistentMapString(t *testing.T) {
	m := EmptyPersistentMap
	if m.String() != "{}" {
		t.Errorf("expected {}, got %s", m.String())
	}

	m2 := NewPersistentMap([]Value{Keyword("a"), Int(1)})
	s := m2.String()
	if s != "{:a 1}" {
		t.Errorf("expected {:a 1}, got %s", s)
	}
}

func TestPersistentMapNewPersistentMap(t *testing.T) {
	m := NewPersistentMap([]Value{})
	if m != EmptyPersistentMap {
		t.Errorf("expected EmptyPersistentMap for empty kvs")
	}
	m = NewPersistentMap([]Value{Keyword("a")})
	if m != EmptyPersistentMap {
		t.Errorf("expected EmptyPersistentMap for odd-length kvs")
	}
}

func TestPersistentMapValueAtOr(t *testing.T) {
	m := NewPersistentMap([]Value{Keyword("a"), Int(1)})
	if m.ValueAtOr(Keyword("a"), Int(99)) != Int(1) {
		t.Errorf("expected 1")
	}
	if m.ValueAtOr(Keyword("missing"), Int(99)) != Int(99) {
		t.Errorf("expected 99")
	}
}

func TestPersistentMapType(t *testing.T) {
	m := EmptyPersistentMap
	if m.Type() != MapType {
		t.Errorf("wrong type")
	}
	if PersistentMapType.Name() != "let-go.lang.PersistentHashMap" {
		t.Errorf("wrong type name")
	}
}

func TestPersistentMapEmptyMethod(t *testing.T) {
	m := NewPersistentMap([]Value{Keyword("a"), Int(1)})
	empty := m.Empty().(*PersistentMap)
	if empty.RawCount() != 0 {
		t.Errorf("Empty() should return empty map")
	}
}

// TestPersistentMapHashCollision tests behavior when keys produce the same hash.
// We craft this by using the collision node path.
func TestPersistentMapHashCollision(t *testing.T) {
	// We can't easily force collisions with real hash functions,
	// so we test indirectly by inserting many keys and verifying all are retrievable.
	// With 1000+ keys, some will share prefix bits and exercise deeper trie levels.
	m := EmptyPersistentMap
	n := 500
	for i := 0; i < n; i++ {
		k := String(fmt.Sprintf("key-%d", i))
		m = m.Assoc(k, Int(i)).(*PersistentMap)
	}
	if m.RawCount() != n {
		t.Fatalf("expected count %d, got %d", n, m.RawCount())
	}
	for i := 0; i < n; i++ {
		k := String(fmt.Sprintf("key-%d", i))
		v := m.ValueAt(k)
		if v != Int(i) {
			t.Errorf("at key %s: expected %d, got %v", k, i, v)
		}
	}
	// Dissoc half and verify
	for i := 0; i < n; i += 2 {
		k := String(fmt.Sprintf("key-%d", i))
		m = m.Dissoc(k).(*PersistentMap)
	}
	expectedCount := n / 2
	if m.RawCount() != expectedCount {
		t.Fatalf("after dissoc expected count %d, got %d", expectedCount, m.RawCount())
	}
	for i := 0; i < n; i++ {
		k := String(fmt.Sprintf("key-%d", i))
		v := m.ValueAt(k)
		if i%2 == 0 {
			if v != NIL {
				t.Errorf("key %s should be dissoc'd, got %v", k, v)
			}
		} else {
			if v != Int(i) {
				t.Errorf("key %s should still be %d, got %v", k, i, v)
			}
		}
	}
}

// TestPersistentMapCollisionNode directly tests collision node behavior
// by manually constructing one.
func TestPersistentMapCollisionNode(t *testing.T) {
	// Create a collision node with two entries sharing the same hash
	cn := &hmapCollisionNode{
		hash:  12345,
		count: 2,
		array: []interface{}{Keyword("a"), Int(1), Keyword("b"), Int(2)},
	}

	// find
	v, ok := cn.find(0, 12345, Keyword("a"))
	if !ok || v != Int(1) {
		t.Errorf("find :a failed")
	}
	v, ok = cn.find(0, 12345, Keyword("b"))
	if !ok || v != Int(2) {
		t.Errorf("find :b failed")
	}
	_, ok = cn.find(0, 12345, Keyword("c"))
	if ok {
		t.Errorf("find :c should fail")
	}

	// assoc existing key
	addedLeaf := false
	cn2 := cn.assoc(0, 12345, Keyword("a"), Int(99), &addedLeaf).(*hmapCollisionNode)
	if addedLeaf {
		t.Errorf("should not set addedLeaf for existing key")
	}
	v, _ = cn2.find(0, 12345, Keyword("a"))
	if v != Int(99) {
		t.Errorf("expected 99, got %v", v)
	}

	// assoc new key with same hash
	addedLeaf = false
	cn3 := cn.assoc(0, 12345, Keyword("c"), Int(3), &addedLeaf).(*hmapCollisionNode)
	if !addedLeaf {
		t.Errorf("should set addedLeaf for new key")
	}
	if cn3.count != 3 {
		t.Errorf("expected count 3, got %d", cn3.count)
	}

	// dissoc
	cn4 := cn.dissoc(0, 12345, Keyword("a"))
	// With count dropping to 1, it should become a bitmap node
	if _, ok := cn4.(*hmapBitmapNode); !ok {
		t.Errorf("expected bitmap node after dissoc to 1 entry, got %T", cn4)
	}
	v, ok = cn4.find(0, 12345, Keyword("b"))
	if !ok || v != Int(2) {
		t.Errorf("remaining entry should be :b -> 2")
	}
}

func TestPersistentMapLargeDissoc(t *testing.T) {
	m := EmptyPersistentMap
	n := 200
	for i := 0; i < n; i++ {
		m = m.Assoc(Int(i), Int(i)).(*PersistentMap)
	}
	// Remove all
	for i := 0; i < n; i++ {
		m = m.Dissoc(Int(i)).(*PersistentMap)
	}
	if m.RawCount() != 0 {
		t.Errorf("expected empty map, got count %d", m.RawCount())
	}
}
