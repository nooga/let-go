package vm

import "testing"

func TestLookup_AliasTriggersNamespaceMaterialization(t *testing.T) {
	caller := NewNamespace("caller")
	placeholder := NewNamespace("xsofy.det")
	caller.Alias(Symbol("det"), placeholder)

	loaded := NewNamespace("xsofy.det")
	loaded.Def("int-in", Int(42))

	prev := nsLookup
	defer SetNSLookup(prev)
	SetNSLookup(func(name string) *Namespace {
		if name == "xsofy.det" {
			return loaded
		}
		return nil
	})

	v := caller.Lookup(Symbol("det/int-in"))
	if v == NIL {
		t.Fatalf("expected aliased var to resolve after namespace materialization")
	}

	resolved, ok := v.(*Var)
	if !ok {
		t.Fatalf("expected *Var, got %T", v)
	}
	if got := resolved.Deref(); got != Int(42) {
		t.Fatalf("resolved var value = %v, want 42", got)
	}
}

func TestLookup_QualifiedAliasFollowsTargetRefers(t *testing.T) {
	// Reproduce the portability/jank regression:
	// - lib namespace defines big-int?
	// - target namespace refers big-int? from lib
	// - caller aliases target as "p"
	// - caller looks up p/big-int? — should find the referred var in target

	lib := NewNamespace("lib")
	lib.Def("big-int?", TRUE)

	target := NewNamespace("target")
	target.Refer(lib, "", true) // :refer :all

	caller := NewNamespace("caller")
	caller.Alias(Symbol("p"), target)

	// Lookup p/big-int? should resolve via target's refers
	v := caller.Lookup(Symbol("p/big-int?"))
	if v == NIL {
		t.Fatalf("qualified alias p/big-int? should resolve via target's refers, got NIL")
	}

	resolved, ok := v.(*Var)
	if !ok {
		t.Fatalf("expected *Var, got %T", v)
	}
	if got := resolved.Deref(); got != TRUE {
		t.Fatalf("resolved var value = %v, want TRUE", got)
	}
}

func TestFuzzySymbolLookup_CyclicRefersTerminate(t *testing.T) {
	// The real refer graph is cyclic: clojure.core requires let-go.types, and
	// RegisterNS auto-refers clojure.core back into it. Completion only reads
	// each refer's locals (not the target's refers), so the cycle cannot make
	// the walk recurse — a still sees b's defcmd, and each name once.
	a := NewNamespace("cyc-a")
	b := NewNamespace("cyc-b")
	a.Refer(b, "", true)
	b.Refer(a, "", true)
	a.Def("defer-me", TRUE)
	b.Def("defcmd", TRUE)

	got := FuzzySymbolLookup(a, Symbol("def"), true)
	names := map[Symbol]int{}
	for _, s := range got {
		names[s]++
	}
	if names["defer-me"] != 1 || names["defcmd"] != 1 {
		t.Fatalf("want each symbol once across the cycle, got %v", got)
	}
}

func TestFuzzySymbolLookup_ReferNotTransitive(t *testing.T) {
	// Refer visibility is one hop: user → mid does not expose mid's refers.
	// A transitive walk would surface "cymbal" here.
	leaf := NewNamespace("leaf")
	leaf.Def("cymbal", TRUE)

	mid := NewNamespace("mid")
	mid.Refer(leaf, "", true)

	user := NewNamespace("user-nt")
	user.Refer(mid, "", true)
	user.Def("cycle", TRUE)

	got := FuzzySymbolLookup(user, Symbol("c"), true)
	names := map[Symbol]bool{}
	for _, s := range got {
		names[s] = true
	}
	if !names["cycle"] {
		t.Fatalf("want user's own cycle, got %v", got)
	}
	if names["cymbal"] {
		t.Fatalf("want cymbal excluded (mid's refer, not user's), got %v", got)
	}
}

func TestFuzzySymbolLookup_HonorsReferOnly(t *testing.T) {
	// (:require [lib :refer [foo]]) should only bring foo into scope, not
	// every other public symbol lib happens to define.
	lib := NewNamespace("only-lib")
	lib.Def("foo", TRUE)
	lib.Def("fred", TRUE)
	lib.Def("flip", TRUE)

	user := NewNamespace("only-user")
	user.ReferList(lib, []Symbol{"foo"})

	got := FuzzySymbolLookup(user, Symbol("f"), true)
	names := map[Symbol]bool{}
	for _, s := range got {
		names[s] = true
	}
	if !names["foo"] {
		t.Fatalf("want foo (explicitly referred) in results, got %v", got)
	}
	if names["fred"] || names["flip"] {
		t.Fatalf("want fred/flip excluded (not in :refer list), got %v", got)
	}
}

func TestFuzzySymbolLookup_HonorsUnmap(t *testing.T) {
	// After (ns-unmap *ns* 'foo), foo should no longer resolve in this
	// namespace even though it's still referred from lib.
	lib := NewNamespace("unmap-lib")
	lib.Def("foo", TRUE)

	user := NewNamespace("unmap-user")
	user.Refer(lib, "", true) // :refer :all
	user.Unmap("foo")

	got := FuzzySymbolLookup(user, Symbol("f"), true)
	for _, s := range got {
		if s == "foo" {
			t.Fatalf("want foo excluded after ns-unmap, got %v", got)
		}
	}
}
