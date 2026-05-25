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
