/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"fmt"
	"sync/atomic"
)

// TransientMap is a mutable version of PersistentMap for batch operations.
// Mutations modify nodes in place instead of path-copying.
// Call persistent! to freeze back to an immutable PersistentMap.
type TransientMap struct {
	edit  atomic.Bool // true while mutable, false after persistent!
	root  hmapNode
	count int
}

func NewTransientMap(m *PersistentMap) *TransientMap {
	t := &TransientMap{
		root:  m.root,
		count: m.count,
	}
	t.edit.Store(true)
	return t
}

func (t *TransientMap) ensureEditable() error {
	if !t.edit.Load() {
		return fmt.Errorf("transient used after persistent! call")
	}
	return nil
}

func (t *TransientMap) Type() ValueType    { return TransientMapType }
func (t *TransientMap) Unbox() interface{} { return t }
func (t *TransientMap) String() string     { return fmt.Sprintf("<transient-map count=%d>", t.count) }

// Assoc mutates the transient map in place.
func (t *TransientMap) Assoc(key Value, val Value) (*TransientMap, error) {
	if err := t.ensureEditable(); err != nil {
		return nil, err
	}
	addedLeaf := false
	hash := hashValue(key)
	if t.root == nil {
		t.root = (&hmapBitmapNode{edit: &t.edit}).assocTransient(&t.edit, 0, hash, key, val, &addedLeaf)
	} else if bn, ok := t.root.(*hmapBitmapNode); ok {
		t.root = bn.assocTransient(&t.edit, 0, hash, key, val, &addedLeaf)
	} else {
		t.root = t.root.assoc(0, hash, key, val, &addedLeaf)
	}
	if addedLeaf {
		t.count++
	}
	return t, nil
}

// Dissoc mutates the transient map in place.
func (t *TransientMap) Dissoc(key Value) (*TransientMap, error) {
	if err := t.ensureEditable(); err != nil {
		return nil, err
	}
	if t.root == nil {
		return t, nil
	}
	hash := hashValue(key)
	newRoot := t.root.dissoc(0, hash, key)
	if newRoot != t.root {
		t.root = newRoot
		t.count--
	}
	return t, nil
}

// Conj adds a [key val] pair.
func (t *TransientMap) Conj(value Value) (*TransientMap, error) {
	if err := t.ensureEditable(); err != nil {
		return nil, err
	}
	if av, ok := value.(ArrayVector); ok && len(av) == 2 {
		return t.Assoc(av[0], av[1])
	}
	return nil, fmt.Errorf("conj! on transient map expects [key val] pair")
}

// Persistent freezes the transient and returns an immutable PersistentMap.
func (t *TransientMap) Persistent() *PersistentMap {
	t.edit.Store(false)
	return &PersistentMap{
		root:  t.root,
		count: t.count,
	}
}

// ValueAt for lookups during construction.
func (t *TransientMap) ValueAt(key Value) Value {
	return t.ValueAtOr(key, NIL)
}

func (t *TransientMap) ValueAtOr(key Value, notFound Value) Value {
	if t.root == nil {
		return notFound
	}
	v, found := t.root.find(0, hashValue(key), key)
	if !found {
		return notFound
	}
	return v
}

func (t *TransientMap) Count() Value  { return MakeInt(t.count) }
func (t *TransientMap) RawCount() int { return t.count }

// TransientVector is a mutable version of ArrayVector/PersistentVector.
type TransientVector struct {
	edit  atomic.Bool
	array []Value
}

func NewTransientVector(v []Value) *TransientVector {
	t := &TransientVector{
		array: make([]Value, len(v)),
	}
	copy(t.array, v)
	t.edit.Store(true)
	return t
}

func (t *TransientVector) ensureEditable() error {
	if !t.edit.Load() {
		return fmt.Errorf("transient used after persistent! call")
	}
	return nil
}

func (t *TransientVector) Type() ValueType    { return TransientVectorType }
func (t *TransientVector) Unbox() interface{} { return t }
func (t *TransientVector) String() string {
	return fmt.Sprintf("<transient-vector count=%d>", len(t.array))
}

// Conj appends to the transient vector in place.
func (t *TransientVector) Conj(value Value) (*TransientVector, error) {
	if err := t.ensureEditable(); err != nil {
		return nil, err
	}
	t.array = append(t.array, value)
	return t, nil
}

// Assoc sets index to value.
func (t *TransientVector) Assoc(key Value, val Value) (*TransientVector, error) {
	if err := t.ensureEditable(); err != nil {
		return nil, err
	}
	idx, ok := key.(Int)
	if !ok {
		return nil, fmt.Errorf("transient vector assoc! expects Int key")
	}
	i := int(idx)
	if i < 0 || i > len(t.array) {
		return nil, fmt.Errorf("index out of bounds: %d", i)
	}
	if i == len(t.array) {
		t.array = append(t.array, val)
	} else {
		t.array[i] = val
	}
	return t, nil
}

// Persistent freezes the transient and returns an immutable vector.
func (t *TransientVector) Persistent() Value {
	t.edit.Store(false)
	if len(t.array) <= arrayVectorPromotionThreshold {
		result := make(ArrayVector, len(t.array))
		copy(result, t.array)
		return result
	}
	return NewPersistentVector(t.array)
}

func (t *TransientVector) ValueAt(key Value) Value {
	idx, ok := key.(Int)
	if !ok || int(idx) < 0 || int(idx) >= len(t.array) {
		return NIL
	}
	return t.array[int(idx)]
}

func (t *TransientVector) Count() Value  { return MakeInt(len(t.array)) }
func (t *TransientVector) RawCount() int { return len(t.array) }

// --- Type metadata ---

type theTransientMapType struct{}

func (t *theTransientMapType) String() string     { return t.Name() }
func (t *theTransientMapType) Type() ValueType    { return TypeType }
func (t *theTransientMapType) Unbox() interface{} { return nil }
func (t *theTransientMapType) Name() string       { return "let-go.lang.TransientMap" }
func (t *theTransientMapType) Box(bare interface{}) (Value, error) {
	return NIL, NewTypeError(bare, "can't be boxed as", t)
}

var TransientMapType *theTransientMapType = &theTransientMapType{}

type theTransientVectorType struct{}

func (t *theTransientVectorType) String() string     { return t.Name() }
func (t *theTransientVectorType) Type() ValueType    { return TypeType }
func (t *theTransientVectorType) Unbox() interface{} { return nil }
func (t *theTransientVectorType) Name() string       { return "let-go.lang.TransientVector" }
func (t *theTransientVectorType) Box(bare interface{}) (Value, error) {
	return NIL, NewTypeError(bare, "can't be boxed as", t)
}

var TransientVectorType *theTransientVectorType = &theTransientVectorType{}
