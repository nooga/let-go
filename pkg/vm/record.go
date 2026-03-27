/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"fmt"
	"strings"
)

// RecordType is a dynamically-created ValueType for a defrecord.
// Each defrecord call creates a unique RecordType instance.
type RecordType struct {
	typeName string
	fields   []Keyword      // ordered field names
	fieldIdx map[Keyword]int // field name → index for O(1) access
}

func NewRecordType(name string, fields []Keyword) *RecordType {
	idx := make(map[Keyword]int, len(fields))
	for i, f := range fields {
		idx[f] = i
	}
	return &RecordType{typeName: name, fields: fields, fieldIdx: idx}
}

func (t *RecordType) String() string     { return t.Name() }
func (t *RecordType) Type() ValueType    { return TypeType }
func (t *RecordType) Unbox() interface{} { return t }
func (t *RecordType) Name() string       { return t.typeName }
func (t *RecordType) Fields() []Keyword  { return t.fields }

func (t *RecordType) Box(bare interface{}) (Value, error) {
	return NIL, NewTypeError(bare, "can't be boxed as", t)
}

// Record is a value created by defrecord.
// Fixed fields are stored in a flat array for O(1) access.
// Extra keys (from assoc) go into an overflow map.
// All map operations work — seq, keys, count include both fixed and extra fields.
type Record struct {
	rtype  *RecordType
	fields []Value        // fixed fields, indexed by RecordType.fieldIdx
	extra  *PersistentMap // overflow for assoc'd keys not in the record definition
	meta   Value
	origin interface{}    // original Go struct for fast roundtrip (nil if mutated)
}

func NewRecord(rtype *RecordType, data *PersistentMap) *Record {
	r := &Record{
		rtype:  rtype,
		fields: make([]Value, len(rtype.fields)),
		extra:  EmptyPersistentMap,
	}
	// Extract fixed fields from the map
	for i, kw := range rtype.fields {
		v := data.ValueAt(kw)
		if v != NIL {
			r.fields[i] = v
		}
	}
	// Anything not a fixed field goes to extra
	s := data.Seq()
	for s != nil && s != EmptyList {
		entry := s.First().(ArrayVector)
		if kw, ok := entry[0].(Keyword); ok {
			if _, isField := rtype.fieldIdx[kw]; isField {
				s = s.Next()
				continue
			}
		}
		r.extra = r.extra.Assoc(entry[0], entry[1]).(*PersistentMap)
		s = s.Next()
	}
	return r
}

// --- Value interface ---

func (r *Record) Type() ValueType    { return r.rtype }
func (r *Record) Unbox() interface{} {
	if r.origin != nil {
		return r.origin
	}
	return r
}

// Origin returns the original Go struct if this Record was created from one, or nil.
func (r *Record) Origin() interface{} { return r.origin }

func (r *Record) String() string {
	b := &strings.Builder{}
	b.WriteString("#")
	b.WriteString(r.rtype.typeName)
	b.WriteString("{")
	first := true
	// Print fixed fields in definition order
	for i, kw := range r.rtype.fields {
		if !first {
			b.WriteString(", ")
		}
		b.WriteRune(':')
		b.WriteString(string(kw))
		b.WriteRune(' ')
		v := r.fields[i]
		if v == nil {
			b.WriteString("nil")
		} else {
			b.WriteString(v.String())
		}
		first = false
	}
	// Print extra fields
	s := r.extra.Seq()
	for s != nil && s != EmptyList {
		entry := s.First().(ArrayVector)
		if !first {
			b.WriteString(", ")
		}
		b.WriteString(entry[0].String())
		b.WriteRune(' ')
		b.WriteString(entry[1].String())
		first = false
		s = s.Next()
	}
	b.WriteString("}")
	return b.String()
}

// --- Hashable ---

func (r *Record) Hash() uint32 {
	// Combine field hashes + extra map hash
	var h uint32
	for i, v := range r.fields {
		if v != nil {
			h += hashValue(r.rtype.fields[i]) ^ hashValue(v)
		}
	}
	h += r.extra.Hash()
	return mixFinish(h)
}

// --- Equals ---

func (r *Record) Equals(other Value) bool {
	o, ok := other.(*Record)
	if !ok {
		return false
	}
	if r.rtype != o.rtype {
		return false
	}
	for i := range r.fields {
		rv, ov := r.fields[i], o.fields[i]
		if rv == nil && ov == nil {
			continue
		}
		if rv == nil || ov == nil {
			return false
		}
		if !valueEquiv(rv, ov) {
			return false
		}
	}
	return r.extra.Equals(o.extra)
}

// --- IMeta ---

func (r *Record) Meta() Value {
	if r.meta == nil {
		return NIL
	}
	return r.meta
}

func (r *Record) WithMeta(m Value) Value {
	return &Record{rtype: r.rtype, fields: r.fields, extra: r.extra, meta: m, origin: r.origin}
}

// --- Lookup (get, keyword access) ---

func (r *Record) ValueAt(key Value) Value {
	return r.ValueAtOr(key, NIL)
}

func (r *Record) ValueAtOr(key Value, notFound Value) Value {
	// Fast path: keyword field access
	if kw, ok := key.(Keyword); ok {
		if idx, ok := r.rtype.fieldIdx[kw]; ok {
			v := r.fields[idx]
			if v == nil {
				return notFound
			}
			return v
		}
	}
	// Fallback: check extra map
	return r.extra.ValueAtOr(key, notFound)
}

// --- Associative ---

func (r *Record) Assoc(key Value, val Value) Associative {
	// If it's a fixed field, update the fields array
	if kw, ok := key.(Keyword); ok {
		if idx, ok := r.rtype.fieldIdx[kw]; ok {
			newFields := make([]Value, len(r.fields))
			copy(newFields, r.fields)
			newFields[idx] = val
			return &Record{rtype: r.rtype, fields: newFields, extra: r.extra}
		}
	}
	// Otherwise, add to extra
	return &Record{rtype: r.rtype, fields: r.fields, extra: r.extra.Assoc(key, val).(*PersistentMap)}
}

func (r *Record) Dissoc(key Value) Associative {
	// Can't remove a fixed field — set to nil
	if kw, ok := key.(Keyword); ok {
		if idx, ok := r.rtype.fieldIdx[kw]; ok {
			newFields := make([]Value, len(r.fields))
			copy(newFields, r.fields)
			newFields[idx] = nil
			return &Record{rtype: r.rtype, fields: newFields, extra: r.extra}
		}
	}
	return &Record{rtype: r.rtype, fields: r.fields, extra: r.extra.Dissoc(key).(*PersistentMap)}
}

// --- Collection ---

func (r *Record) Count() Value  { return MakeInt(r.RawCount()) }
func (r *Record) RawCount() int {
	n := r.extra.RawCount()
	for _, v := range r.fields {
		if v != nil {
			n++
		}
	}
	return n
}

func (r *Record) Empty() Collection {
	return &Record{rtype: r.rtype, fields: make([]Value, len(r.rtype.fields)), extra: EmptyPersistentMap}
}

func (r *Record) Conj(value Value) Collection {
	// Conj a [k v] pair
	if av, ok := value.(ArrayVector); ok && len(av) == 2 {
		return r.Assoc(av[0], av[1]).(Collection)
	}
	return r
}

// --- Keyed ---

func (r *Record) Contains(key Value) Boolean {
	if kw, ok := key.(Keyword); ok {
		if idx, ok := r.rtype.fieldIdx[kw]; ok {
			if r.fields[idx] != nil {
				return TRUE
			}
			return FALSE
		}
	}
	return r.extra.Contains(key)
}

// --- Sequable ---

// Seq returns all entries — fixed fields first (in definition order), then extras.
func (r *Record) Seq() Seq {
	if r.RawCount() == 0 {
		return EmptyList
	}
	var entries []Value
	for i, kw := range r.rtype.fields {
		if r.fields[i] != nil {
			entries = append(entries, ArrayVector{kw, r.fields[i]})
		}
	}
	// Append extra entries
	s := r.extra.Seq()
	for s != nil && s != EmptyList {
		entries = append(entries, s.First())
		s = s.Next()
	}
	if len(entries) == 0 {
		return EmptyList
	}
	result, _ := ListType.Box(entries)
	return result.(*List)
}

// --- Fn (record as function, like maps) ---

func (r *Record) Arity() int { return -1 }

func (r *Record) Invoke(args []Value) (Value, error) {
	if len(args) < 1 || len(args) > 2 {
		return NIL, fmt.Errorf("wrong number of arguments %d", len(args))
	}
	if len(args) == 1 {
		return r.ValueAt(args[0]), nil
	}
	return r.ValueAtOr(args[0], args[1]), nil
}

// --- Receiver (interop .field access) ---

func (r *Record) InvokeMethod(name Symbol, args []Value) (Value, error) {
	kw := Keyword(name)
	v := r.ValueAt(kw)
	if v != NIL {
		return v, nil
	}
	return NIL, fmt.Errorf("no method %s on record %s", name, r.rtype.typeName)
}

// --- RecordType accessor ---

func (r *Record) RecordType() *RecordType {
	return r.rtype
}
