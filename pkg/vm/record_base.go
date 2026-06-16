/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"sync"
)

// RecordBase is embedded by gogen-generated record/op structs so they satisfy
// vm.Value without per-type boilerplate. Field-specific behavior (field access,
// equality over fields) is added by the generated struct; this covers the
// type-identity and stringification contract.
//
// Usage in generated code:
//
//	type MyRecord struct {
//	    RecordBase
//	    fieldOne Value
//	    fieldTwo Value
//	}
//
// The embedding struct automatically satisfies the Value interface.
type RecordBase struct {
	TypeName string
	typeOnce sync.Once
	typeVal  ValueType
}

// String returns the TypeName. Required by fmt.Stringer (part of Value interface).
func (r *RecordBase) String() string {
	return r.TypeName
}

// Type returns a ValueType representing this record's type.
// The ValueType is lazily created and cached.
func (r *RecordBase) Type() ValueType {
	r.typeOnce.Do(func() {
		r.typeVal = &recordType{name: r.TypeName}
	})
	return r.typeVal
}

// Unbox returns the embedded *RecordBase, NOT the outer embedding struct.
// Because Unbox is a method on RecordBase promoted onto the embedding type, its
// receiver is always *RecordBase — Go method promotion cannot recover the outer
// value. Callers that need the concrete record (e.g. *Square) must hold a typed
// reference; they cannot recover it through this Value-level Unbox. This is
// sufficient for the type-identity/stringification contract RecordBase covers;
// field access on native records goes through the generated struct directly.
func (r *RecordBase) Unbox() any {
	return r
}

// recordType is a simple ValueType that wraps a record type name.
// It's created lazily by RecordBase.Type() and cached.
type recordType struct {
	name string
}

func (rt *recordType) String() string {
	return rt.name
}

func (rt *recordType) Type() ValueType {
	// A ValueType's type is itself (self-describing).
	return rt
}

func (rt *recordType) Unbox() any {
	return rt
}

func (rt *recordType) Name() string {
	return rt.name
}

func (rt *recordType) Box(v any) (Value, error) {
	// Records cannot be boxed; they are created by their constructor.
	return NIL, NewTypeError(v, "can't be boxed as", rt)
}
