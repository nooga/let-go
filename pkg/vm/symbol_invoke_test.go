package vm

import "testing"

// Symbols are invokable for collection lookup exactly like keywords:
// ('a {'a 1}) → 1. HoneySQL relies on this (symbols as lookup fns in
// format-do-update-set).
func TestSymbolInvoke(t *testing.T) {
	fn, ok := Value(Symbol("a")).(Fn)
	if !ok {
		t.Fatal("Symbol must implement Fn")
	}

	m := Map{Symbol("a"): Int(1)}

	got, err := fn.Invoke([]Value{m})
	if err != nil {
		t.Fatalf("('a {'a 1}): %v", err)
	}
	if got != Int(1) {
		t.Fatalf("('a {'a 1}) = %v, want 1", got)
	}

	missing := Value(Symbol("missing")).(Fn)
	got, err = missing.Invoke([]Value{m, Keyword("default")})
	if err != nil {
		t.Fatalf("('missing m :default): %v", err)
	}
	if got != Keyword("default") {
		t.Fatalf("('missing m :default) = %v, want :default", got)
	}

	got, err = missing.Invoke([]Value{m})
	if err != nil {
		t.Fatalf("('missing m): %v", err)
	}
	if got != NIL {
		t.Fatalf("('missing m) = %v, want nil", got)
	}

	s := NewSet([]Value{Symbol("a")})
	got, err = fn.Invoke([]Value{s})
	if err != nil {
		t.Fatalf("('a #{'a}): %v", err)
	}
	if got != Symbol("a") {
		t.Fatalf("('a #{'a}) = %v, want 'a", got)
	}

	// Non-Lookup argument mirrors Keyword: nil for 1-arity, default for 2.
	got, err = fn.Invoke([]Value{Int(42)})
	if err != nil {
		t.Fatalf("('a 42): %v", err)
	}
	if got != NIL {
		t.Fatalf("('a 42) = %v, want nil", got)
	}
	got, err = fn.Invoke([]Value{Int(42), Keyword("d")})
	if err != nil {
		t.Fatalf("('a 42 :d): %v", err)
	}
	if got != Keyword("d") {
		t.Fatalf("('a 42 :d) = %v, want :d", got)
	}

	// Arity errors: 0 and 3+ args.
	if _, err = fn.Invoke([]Value{}); err == nil {
		t.Fatal("('a) with no args must error")
	}
	if _, err = fn.Invoke([]Value{m, NIL, NIL}); err == nil {
		t.Fatal("('a m nil nil) must error")
	}
}
