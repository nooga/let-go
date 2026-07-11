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
