/*
 * Copyright (c) 2026 let-go contributors; see CONTRIBUTORS.
 * SPDX-License-Identifier: MIT
 */

package vm

import "testing"

// defineIn evaluates a defrecord while scope is active, as a unit's own
// defrecord runs.
func defineIn(scope *TypeScope, ns, name string, fields []Keyword) (t *RecordType) {
	_ = WithTypeScope(scope, func() error {
		t = DefineRecordType(ns, name, fields)
		return nil
	})
	return t
}

func TestResolveRecordTypeFindsDefinedType(t *testing.T) {
	fields := []Keyword{"x", "y"}
	defined := DefineRecordType("registry.defined", "Pt", fields)
	if got := ResolveRecordType(NewTypeScope(), "registry.defined", "Pt", fields); got != defined {
		t.Fatalf("decoded constant got a new type, want the defined one")
	}
	if got := ResolveRecordType(nil, "registry.defined", "Pt", fields); got != defined {
		t.Fatalf("a decode with no unit got a new type, want the defined one")
	}
	if got := ResolveRecordType(NewTypeScope(), "registry.other", "Pt", fields); got == defined {
		t.Fatalf("a same-named type in another namespace resolved to this one")
	}
	if got := ResolveRecordType(NewTypeScope(), "registry.defined", "Pt", []Keyword{"x"}); got == defined {
		t.Fatalf("a constant with other fields resolved to the defined type")
	}
	if got := ResolveRecordType(NewTypeScope(), "registry.defined", "Pt", fields); got != defined {
		t.Fatalf("a field mismatch replaced the registered type")
	}
}

// A unit's constants decode before its defrecord runs; that defrecord adopts
// their type once, and evaluating it again makes a new type.
func TestUnitDefinitionAdoptsDecodedTypeOnce(t *testing.T) {
	fields := []Keyword{"x"}
	scope := NewTypeScope()
	decoded := ResolveRecordType(scope, "registry.adopt", "Pt", fields)
	if again := ResolveRecordType(scope, "registry.adopt", "Pt", fields); again != decoded {
		t.Fatalf("two constants of one unit got different types")
	}
	var first, second *RecordType
	_ = WithTypeScope(scope, func() error {
		first = DefineRecordType("registry.adopt", "Pt", fields)
		second = DefineRecordType("registry.adopt", "Pt", fields)
		return nil
	})
	if first != decoded {
		t.Fatalf("the unit's defrecord did not adopt the type its constant decoded to")
	}
	if second == decoded {
		t.Fatalf("redefining the record kept the old type")
	}
	if got := ResolveRecordType(NewTypeScope(), "registry.adopt", "Pt", fields); got != second {
		t.Fatalf("a constant decoded after redefinition did not get the new type")
	}
}

// A unit loaded again decodes its constant while its type is defined, then
// replays the defining form, which must keep that type.
func TestReplayedDefinitionKeepsDecodedType(t *testing.T) {
	fields := []Keyword{"x"}
	defined := DefineRecordType("registry.replay", "Pt", fields)
	for load := 1; load <= 2; load++ {
		scope := NewTypeScope()
		decoded := ResolveRecordType(scope, "registry.replay", "Pt", fields)
		if got := defineIn(scope, "registry.replay", "Pt", fields); got != decoded || got != defined {
			t.Fatalf("load %d: the replayed defrecord replaced the type its constant decoded to", load)
		}
	}
	if DefineRecordType("registry.replay", "Pt", fields) == defined {
		t.Fatalf("redefining the record outside a load kept the old type")
	}
}

// A unit whose record has other fields than the live one decodes its constant
// to a type of its own, and its defrecord adopts that type.
func TestFieldChangeBetweenLoadsKeepsOneType(t *testing.T) {
	live := DefineRecordType("registry.fields", "Pt", []Keyword{"x", "y"})
	fields := []Keyword{"x", "y", "z"}
	scope := NewTypeScope()
	decoded := ResolveRecordType(scope, "registry.fields", "Pt", fields)
	if decoded == live {
		t.Fatalf("a constant with other fields got the live type")
	}
	if got := defineIn(scope, "registry.fields", "Pt", fields); got != decoded {
		t.Fatalf("the unit's defrecord made a third type instead of adopting the decoded one")
	}
	if got := ResolveRecordType(NewTypeScope(), "registry.fields", "Pt", fields); got != decoded {
		t.Fatalf("the adopted type was not registered")
	}
}

// Decoding a unit that is never run must not reach a later definition.
func TestDecodeAloneDoesNotReachALaterDefinition(t *testing.T) {
	fields := []Keyword{"v"}
	defined := DefineRecordType("registry.decode-only", "Box", fields)
	scope := NewTypeScope()
	if ResolveRecordType(scope, "registry.decode-only", "Box", fields) != defined {
		t.Fatalf("decoded constant got a new type, want the defined one")
	}
	if DefineRecordType("registry.decode-only", "Box", fields) == defined {
		t.Fatalf("a decode that never ran made the next defrecord reuse the old type")
	}
	if got := defineIn(NewTypeScope(), "registry.decode-only", "Box", fields); got == defined {
		t.Fatalf("another unit's scope made its defrecord reuse the old type")
	}
}

func TestTypeScopeEndsWithItsRun(t *testing.T) {
	fields := []Keyword{"x"}
	scope := NewTypeScope()
	decoded := ResolveRecordType(scope, "registry.ended", "Pt", fields)
	_ = WithTypeScope(scope, func() error { return nil })
	if DefineRecordType("registry.ended", "Pt", fields) == decoded {
		t.Fatalf("a scope stayed active after its run returned")
	}
}

func TestTypesWithoutNamespaceAreNotRegistered(t *testing.T) {
	a := DefineRecordType("", "Pt", []Keyword{"x"})
	if got := ResolveRecordType(NewTypeScope(), "", "Pt", []Keyword{"x"}); got == a {
		t.Fatalf("a type with no namespace was registered")
	}
	d := DefineDType("", "Box", []Symbol{"v"})
	if got := ResolveDType(NewTypeScope(), "", "Box", []Symbol{"v"}); got == d {
		t.Fatalf("a deftype with no namespace was registered")
	}
}

func TestDTypeRegistry(t *testing.T) {
	fields := []Symbol{"v"}
	scope := NewTypeScope()
	decoded := ResolveDType(scope, "registry.dtype", "Box", fields)
	var defined *DType
	_ = WithTypeScope(scope, func() error {
		defined = DefineDType("registry.dtype", "Box", fields)
		return nil
	})
	if defined != decoded {
		t.Fatalf("deftype did not adopt the type its decoded constant made")
	}
	for load := 1; load <= 2; load++ {
		scope := NewTypeScope()
		if got := ResolveDType(scope, "registry.dtype", "Box", fields); got != defined {
			t.Fatalf("load %d: decoded constant got a new type, want the defined one", load)
		}
		var replayed *DType
		_ = WithTypeScope(scope, func() error {
			replayed = DefineDType("registry.dtype", "Box", fields)
			return nil
		})
		if replayed != defined {
			t.Fatalf("load %d: the replayed deftype replaced the type its constant decoded to", load)
		}
	}
	DefineRecordType("registry.dtype", "Box", []Keyword{"v"})
	if got := ResolveDType(NewTypeScope(), "registry.dtype", "Box", fields); got != defined {
		t.Fatalf("a record of the same name displaced the deftype")
	}
	if DefineDType("registry.dtype", "Box", fields) == defined {
		t.Fatalf("redefining the deftype outside a load kept the old type")
	}
}
