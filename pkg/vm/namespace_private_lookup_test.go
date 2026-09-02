package vm

import "testing"

// LookupIncludingPrivate backs the var special form (#'ns/sym): it must
// resolve private vars in foreign namespaces, while Lookup — used for bare
// qualified symbol references — keeps hiding them.
func TestLookupIncludingPrivate_AliasPath(t *testing.T) {
	target := NewNamespace("private-lookup-target")
	priv := target.Def("hidden", Int(42))
	priv.SetPrivate()
	target.Def("visible", Int(7))

	user := NewNamespace("private-lookup-user")
	user.Alias(Symbol("p"), target)

	if v := user.Lookup(Symbol("p/hidden")); v != NIL {
		t.Fatalf("Lookup must hide a foreign private var, got %v", v)
	}
	v := user.LookupIncludingPrivate(Symbol("p/hidden"))
	if v == NIL {
		t.Fatal("LookupIncludingPrivate must resolve a foreign private var via an alias")
	}
	if v.(*Var) != priv {
		t.Fatalf("resolved the wrong var: %v", v)
	}
	if got := user.LookupIncludingPrivate(Symbol("p/visible")); got == NIL {
		t.Fatal("LookupIncludingPrivate must still resolve public vars")
	}
}

func TestLookupIncludingPrivate_GlobalRegistryPath(t *testing.T) {
	target := NewNamespace("private-lookup-global")
	priv := target.Def("hidden", Int(42))
	priv.SetPrivate()

	orig := nsLookup
	defer func() { nsLookup = orig }()
	nsLookup = func(name string) *Namespace {
		if name == "private-lookup-global" {
			return target
		}
		if orig != nil {
			return orig(name)
		}
		return nil
	}

	user := NewNamespace("private-lookup-global-user")
	if v := user.Lookup(Symbol("private-lookup-global/hidden")); v != NIL {
		t.Fatalf("Lookup must hide a foreign private var, got %v", v)
	}
	v := user.LookupIncludingPrivate(Symbol("private-lookup-global/hidden"))
	if v == NIL {
		t.Fatal("LookupIncludingPrivate must resolve a foreign private var via the global registry")
	}
	if v.(*Var) != priv {
		t.Fatalf("resolved the wrong var: %v", v)
	}
}
