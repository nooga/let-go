/*
 * Copyright (c) 2026 let-go contributors; see CONTRIBUTORS.
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"slices"
	"sync"
)

// A defrecord or deftype type is made when its defining form runs, but a
// compiled unit can carry instances of it as constants, and a unit's
// constants are decoded before its forms run. The registry names each defined
// type by namespace and name, so a decoded constant gets the live type, as
// Clojure reconstructs a record constant through its class.
//
// Each decoded unit has a TypeScope holding the types its constants resolved
// to: the registered type when the fields match, else one made for the unit.
// While the unit runs (WithTypeScope), its own defrecord or deftype adopts the
// scope's type with the same fields, once, so the constants and the defined
// type agree on the first load, on a repeated load, and after a field change.
// Any other evaluation of the form makes a new type, as redefining a record in
// Clojure makes a new class, and constants decoded after that resolve to it.

type typeKey struct {
	ns, name string
	deftype  bool
}

// TypeScope holds the types one decoded unit's constants resolved to, for that
// unit's defining forms to adopt while it runs.
type TypeScope struct {
	types map[typeKey]ValueType
}

// NewTypeScope returns an empty scope for a unit about to be decoded.
func NewTypeScope() *TypeScope { return &TypeScope{} }

var (
	typeRegistryMu sync.Mutex
	typeRegistry   = map[typeKey]ValueType{}
	activeScopes   []*TypeScope
)

// WithTypeScope runs run with s active, so a defrecord or deftype it evaluates
// can adopt a type s holds. A nil or empty scope costs nothing.
//
// Active scopes are process-wide: a defining form evaluated on another
// goroutine while s is active can adopt one of its types too.
func WithTypeScope(s *TypeScope, run func() error) error {
	typeRegistryMu.Lock()
	if s == nil || len(s.types) == 0 {
		typeRegistryMu.Unlock()
		return run()
	}
	activeScopes = append(activeScopes, s)
	typeRegistryMu.Unlock()
	defer func() {
		typeRegistryMu.Lock()
		defer typeRegistryMu.Unlock()
		for i := len(activeScopes) - 1; i >= 0; i-- {
			if activeScopes[i] == s {
				activeScopes = slices.Delete(activeScopes, i, i+1)
				break
			}
		}
	}()
	return run()
}

// registeredType is a type the registry keeps: a *RecordType or a *DType.
type registeredType[F comparable] interface {
	ValueType
	Fields() []F
}

// defineType returns the type a defining form evaluated in key.ns makes,
// adopting an active scope's type with the same fields, and registers it.
func defineType[T registeredType[F], F comparable](key typeKey, fields []F, newType func() T) T {
	if key.ns == "" {
		return newType()
	}
	typeRegistryMu.Lock()
	defer typeRegistryMu.Unlock()
	for i := len(activeScopes) - 1; i >= 0; i-- {
		s := activeScopes[i]
		if t, ok := s.types[key].(T); ok && slices.Equal(t.Fields(), fields) {
			delete(s.types, key)
			typeRegistry[key] = t
			return t
		}
	}
	t := newType()
	typeRegistry[key] = t
	return t
}

// resolveType returns the type a constant decoded into scope names: the
// registered type if its fields match, else the scope's type for it, made if
// need be. A nil scope (a decode with no unit to run) gets an unscoped type.
func resolveType[T registeredType[F], F comparable](scope *TypeScope, key typeKey, fields []F, newType func() T) T {
	if key.ns == "" {
		return newType()
	}
	typeRegistryMu.Lock()
	defer typeRegistryMu.Unlock()
	t, ok := typeRegistry[key].(T)
	if !ok || !slices.Equal(t.Fields(), fields) {
		if scope == nil {
			return newType()
		}
		if t, ok = scope.types[key].(T); ok && slices.Equal(t.Fields(), fields) {
			return t
		}
		t = newType()
	}
	if scope != nil {
		if scope.types == nil {
			scope.types = map[typeKey]ValueType{}
		}
		scope.types[key] = t
	}
	return t
}

// DefineRecordType returns the type for (defrecord name fields) evaluated in
// ns.
func DefineRecordType(ns, name string, fields []Keyword) *RecordType {
	return defineType(typeKey{ns: ns, name: name}, fields, func() *RecordType { return newRecordTypeIn(ns, name, fields) })
}

// ResolveRecordType returns the type a record constant decoded into scope
// names.
func ResolveRecordType(scope *TypeScope, ns, name string, fields []Keyword) *RecordType {
	return resolveType(scope, typeKey{ns: ns, name: name}, fields, func() *RecordType { return newRecordTypeIn(ns, name, fields) })
}

// DefineDType is DefineRecordType for deftype.
func DefineDType(ns, name string, fields []Symbol) *DType {
	return defineType(typeKey{ns: ns, name: name, deftype: true}, fields, func() *DType { return newDTypeIn(ns, name, fields) })
}

// ResolveDType is ResolveRecordType for deftype.
func ResolveDType(scope *TypeScope, ns, name string, fields []Symbol) *DType {
	return resolveType(scope, typeKey{ns: ns, name: name, deftype: true}, fields, func() *DType { return newDTypeIn(ns, name, fields) })
}

func newRecordTypeIn(ns, name string, fields []Keyword) *RecordType {
	t := NewRecordType(name, fields)
	t.ns = ns
	return t
}

func newDTypeIn(ns, name string, fields []Symbol) *DType {
	t := NewDType(name, fields)
	t.ns = ns
	return t
}
