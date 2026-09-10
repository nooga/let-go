package vm

import "testing"

func TestVarMetadata(t *testing.T) {
	v := NewVar(NewNamespace("test.var-meta"), "test.var-meta", "x")
	if got := v.Meta(); got != NIL {
		t.Fatalf("new var meta = %v, want nil", got)
	}

	meta := NewPersistentMap([]Value{Keyword("doc"), String("the doc")})
	v.SetMeta(meta)
	if got := v.Meta(); got != meta {
		t.Fatalf("var meta = %v, want %v", got, meta)
	}
}

func TestVarCompactMetadataIsMaterializedLazily(t *testing.T) {
	v := NewVar(NewNamespace("test.var-compact-meta"), "test.var-compact-meta", "x")
	pairs := DefMetaPairs{
		Keyword("doc"), String("the doc"),
		Keyword("dynamic"), TRUE,
	}

	v.SetMetaPairs(pairs)
	stored, ok := v.meta.(DefMetaPairs)
	if !ok {
		t.Fatalf("SetMetaPairs stored %T, want compact DefMetaPairs", v.meta)
	}
	if got := len(stored); got != len(pairs) {
		t.Fatalf("stored pair count = %d, want %d", got, len(pairs))
	}

	meta := v.Meta()
	m, ok := meta.(*PersistentMap)
	if !ok {
		t.Fatalf("Meta() = %T, want *PersistentMap", meta)
	}
	if got := m.ValueAt(Keyword("doc")); got != String("the doc") {
		t.Fatalf(":doc = %v, want the doc", got)
	}
	if got := m.ValueAt(Keyword("dynamic")); got != TRUE {
		t.Fatalf(":dynamic = %v, want true", got)
	}
	if v.meta != meta {
		t.Fatalf("Meta() did not replace compact pairs with the cached map")
	}
	if got := v.Meta(); got != meta {
		t.Fatalf("second Meta() returned a different map: %p != %p", got, meta)
	}

	replacement := NewPersistentMap([]Value{Keyword("private"), TRUE})
	v.SetMeta(replacement)
	if _, stale := v.meta.(DefMetaPairs); stale {
		t.Fatalf("SetMeta retained stale compact pairs")
	}
	if got := v.Meta(); got != replacement {
		t.Fatalf("replacement meta = %v, want %v", got, replacement)
	}
}

func TestVarAlterMetaMaterializesCompactPairs(t *testing.T) {
	v := NewVar(NewNamespace("test.var-alter-meta"), "test.var-alter-meta", "x")
	v.SetMetaPairs(DefMetaPairs{Keyword("doc"), String("before")})

	seenMap := false
	boxed, err := NativeFnType.Box(func(meta Value) (Value, error) {
		_, seenMap = meta.(*PersistentMap)
		return NewPersistentMap([]Value{Keyword("doc"), String("after")}), nil
	})
	if err != nil {
		t.Fatalf("boxing alter-meta fn: %v", err)
	}
	updated, err := v.AlterMeta(boxed.(Fn), nil)
	if err != nil {
		t.Fatalf("AlterMeta: %v", err)
	}
	if !seenMap {
		t.Fatal("AlterMeta fn observed compact pairs instead of a metadata map")
	}
	if got := updated.(*PersistentMap).ValueAt(Keyword("doc")); got != String("after") {
		t.Fatalf("updated :doc = %v, want after", got)
	}
	if got := v.Meta(); got != updated {
		t.Fatalf("stored metadata = %v, want %v", got, updated)
	}
}

func TestVarEmptyCompactMetadataMaterializesEmptyMap(t *testing.T) {
	v := NewVar(NewNamespace("test.var-empty-meta"), "test.var-empty-meta", "x")
	v.SetMetaPairs(DefMetaPairs{})
	meta, ok := v.Meta().(*PersistentMap)
	if !ok || meta.RawCount() != 0 {
		t.Fatalf("Meta() = %v (%T), want empty PersistentMap", meta, meta)
	}
}

func TestDefMetaPairsHasDistinctInternalType(t *testing.T) {
	pairs := DefMetaPairs{Keyword("doc"), String("doc")}
	if pairs.Type() == ArrayVectorType {
		t.Fatal("DefMetaPairs claims ArrayVectorType without implementing vector contracts")
	}
	if got, want := pairs.Type().Name(), "let-go.internal.DefMetaPairs"; got != want {
		t.Fatalf("DefMetaPairs type name = %q, want %q", got, want)
	}
}
