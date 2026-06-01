/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"fmt"
	"strings"
	"unsafe"
)

// DType is a dynamically-created ValueType for a deftype.
// Unlike RecordType, instances of a DType do NOT implement map behavior
// (no Associative, Lookup, Seqable, or Collection interfaces).
type DType struct {
	typeName string
	fields   []Symbol
	fieldIdx map[Symbol]int
}

func NewDType(name string, fields []Symbol) *DType {
	idx := make(map[Symbol]int, len(fields))
	for i, f := range fields {
		idx[f] = i
	}
	return &DType{typeName: name, fields: fields, fieldIdx: idx}
}

func (t *DType) String() string   { return t.typeName }
func (t *DType) Type() ValueType  { return TypeType }
func (t *DType) Unbox() any       { return t }
func (t *DType) Name() string     { return t.typeName }
func (t *DType) Fields() []Symbol { return t.fields }

func (t *DType) Box(bare any) (Value, error) {
	return NIL, NewTypeError(bare, "can't be boxed as", t)
}

// DTypeInstance is a value created by a deftype constructor.
// Fields are stored in a flat array, accessed by symbol name via Receiver.
type DTypeInstance struct {
	dtype  *DType
	fields []Value
}

func NewDTypeInstance(dt *DType, fields []Value) *DTypeInstance {
	if len(fields) != len(dt.fields) {
		// Pad or truncate to match field count
		padded := make([]Value, len(dt.fields))
		copy(padded, fields)
		for i := range padded {
			if padded[i] == nil {
				padded[i] = NIL
			}
		}
		return &DTypeInstance{dtype: dt, fields: padded}
	}
	return &DTypeInstance{dtype: dt, fields: fields}
}

func (d *DTypeInstance) Type() ValueType { return d.dtype }
func (d *DTypeInstance) Unbox() any      { return d }

func (d *DTypeInstance) String() string {
	b := &strings.Builder{}
	b.WriteString("#<")
	b.WriteString(d.dtype.typeName)
	for i, f := range d.dtype.fields {
		b.WriteRune(' ')
		b.WriteString(string(f))
		b.WriteRune('=')
		v := d.fields[i]
		if v == nil {
			b.WriteString("nil")
		} else {
			b.WriteString(v.String())
		}
	}
	b.WriteString(">")
	return b.String()
}

// Equals: reference identity, matching JVM deftype semantics.
func (d *DTypeInstance) Equals(other Value) bool {
	o, ok := other.(*DTypeInstance)
	if !ok {
		return false
	}
	return d == o
}

// Hash: identity-based (matches JVM default hashCode).
func (d *DTypeInstance) Hash() uint32 {
	return mixFinish(uint32(uintptr(unsafe.Pointer(d))))
}

// InvokeMethod implements Receiver so (.fieldname instance) returns the field.
// Both (.field x) and the explicit field-access form (.-field x) are accepted:
// the reader hands the latter through as the member symbol "-field", so a
// leading "-" is stripped before the field lookup (Clojure's .-field syntax).
func (d *DTypeInstance) InvokeMethod(name Symbol, args []Value) (Value, error) {
	if idx, ok := d.dtype.fieldIdx[name]; ok {
		return d.fields[idx], nil
	}
	if len(name) > 1 && name[0] == '-' {
		if idx, ok := d.dtype.fieldIdx[name[1:]]; ok {
			return d.fields[idx], nil
		}
	}
	return NIL, fmt.Errorf("no field %s on deftype %s", name, d.dtype.typeName)
}

// Fields returns the field values.
func (d *DTypeInstance) Fields() []Value { return d.fields }

// SetField mutates a field in place by name. Backs cljs-style ^:mutable
// deftype fields. UNSYNCHRONIZED: not safe for concurrent access; callers must
// ensure exclusive access (single-threaded use, atoms, or explicit locks).
// A Go nil is normalized to NIL to preserve the field-value invariant.
func (d *DTypeInstance) SetField(name Symbol, val Value) error {
	idx, ok := d.dtype.fieldIdx[name]
	if !ok {
		return fmt.Errorf("no field %s on deftype %s", name, d.dtype.typeName)
	}
	if val == nil {
		val = NIL
	}
	d.fields[idx] = val
	return nil
}

// DType returns the type.
func (d *DTypeInstance) DType() *DType { return d.dtype }
