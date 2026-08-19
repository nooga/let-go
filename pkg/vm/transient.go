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
//
// Like its persistent counterpart it has two modes: an insertion-ordered
// array-map mode (hamt == false, entries in okvs) used while the map stays
// at or below arrayMapMaxEntries, and the HAMT mode past that. This keeps
// the into/transduce path order-preserving for small maps, matching
// Clojure's transient array maps.
type TransientMap struct {
	edit  atomic.Bool // true while mutable, false after persistent!
	root  hmapNode
	okvs  []Value // ordered mode storage (owned, mutated in place)
	hamt  bool    // true once promoted (or seeded from a HAMT-mode map)
	count int
}

func NewTransientMap(m *PersistentMap) *TransientMap {
	t := &TransientMap{count: m.count}
	if m.ordered() {
		t.okvs = make([]Value, len(m.okvs), 2*(arrayMapMaxEntries+1))
		copy(t.okvs, m.okvs)
	} else {
		t.root = m.root
		t.hamt = true
	}
	t.edit.Store(true)
	return t
}

// forceHAMT switches an ordered transient to the HAMT mode, pushing the
// ordered entries into the trie.
func (t *TransientMap) forceHAMT() {
	if t.hamt {
		return
	}
	for i := 0; i < len(t.okvs); i += 2 {
		addedLeaf := false
		key, val := t.okvs[i], t.okvs[i+1]
		hash := hashValue(key)
		if t.root == nil {
			t.root = (&hmapBitmapNode{edit: &t.edit}).assocTransient(&t.edit, 0, hash, key, val, &addedLeaf)
		} else if bn, ok := t.root.(*hmapBitmapNode); ok {
			t.root = bn.assocTransient(&t.edit, 0, hash, key, val, &addedLeaf)
		} else {
			t.root = t.root.assoc(0, hash, key, val, &addedLeaf)
		}
	}
	t.okvs = nil
	t.hamt = true
}

// okvsIndex linear-scans the ordered entries for key.
func (t *TransientMap) okvsIndex(key Value) int {
	for i := 0; i < len(t.okvs); i += 2 {
		if valueEquiv(key, t.okvs[i]) {
			return i
		}
	}
	return -1
}

func (t *TransientMap) ensureEditable() error {
	if !t.edit.Load() {
		return fmt.Errorf("transient used after persistent! call")
	}
	return nil
}

func (t *TransientMap) Type() ValueType { return TransientMapType }
func (t *TransientMap) Unbox() any      { return t }
func (t *TransientMap) String() string  { return fmt.Sprintf("<transient-map count=%d>", t.count) }

// Assoc mutates the transient map in place.
func (t *TransientMap) Assoc(key Value, val Value) (*TransientMap, error) {
	if err := t.ensureEditable(); err != nil {
		return nil, err
	}
	if !t.hamt {
		if i := t.okvsIndex(key); i >= 0 {
			t.okvs[i+1] = val
			return t, nil
		}
		if t.count < arrayMapMaxEntries {
			t.okvs = append(t.okvs, key, val)
			t.count++
			return t, nil
		}
		t.forceHAMT()
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
	if !t.hamt {
		if i := t.okvsIndex(key); i >= 0 {
			t.okvs = append(t.okvs[:i], t.okvs[i+2:]...)
			t.count--
		}
		return t, nil
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
	if value == NIL {
		return t, nil
	}
	switch m := value.(type) {
	case *PersistentMap:
		for _, e := range m.entries() {
			k, v, ok := MapEntryKV(e)
			if !ok {
				continue
			}
			var err error
			t, err = t.Assoc(k, v)
			if err != nil {
				return nil, err
			}
		}
		return t, nil
	case *SortedMap:
		for _, e := range m.entries() {
			var err error
			t, err = t.Assoc(e.Key, e.Value)
			if err != nil {
				return nil, err
			}
		}
		return t, nil
	case Map:
		for k, v := range m {
			var err error
			t, err = t.Assoc(k, v)
			if err != nil {
				return nil, err
			}
		}
		return t, nil
	}
	if k, v, ok := MapEntryKV(value); ok {
		return t.Assoc(k, v)
	}
	if v, ok := value.(PersistentVector); ok && v.RawCount() == 2 {
		return t.Assoc(v.ValueAt(Int(0)), v.ValueAt(Int(1)))
	}
	return nil, fmt.Errorf("conj! on transient map expects [key val] pair")
}

// Persistent freezes the transient and returns an immutable PersistentMap.
func (t *TransientMap) Persistent() (*PersistentMap, error) {
	if err := t.ensureEditable(); err != nil {
		return nil, err
	}
	t.edit.Store(false)
	if !t.hamt {
		// Hand the slice off: the transient is frozen, so nothing can
		// mutate it after this point.
		return &PersistentMap{count: t.count, okvs: t.okvs}, nil
	}
	return &PersistentMap{
		root:  t.root,
		count: t.count,
	}, nil
}

// ValueAt for lookups during construction.
func (t *TransientMap) ValueAt(key Value) Value {
	return t.ValueAtOr(key, NIL)
}

func (t *TransientMap) ValueAtOr(key Value, notFound Value) Value {
	if !t.hamt {
		if i := t.okvsIndex(key); i >= 0 {
			return t.okvs[i+1]
		}
		return notFound
	}
	if t.root == nil {
		return notFound
	}
	v, found := t.root.find(0, hashValue(key), key)
	if !found {
		return notFound
	}
	return v
}

func (t *TransientMap) Seq() Seq {
	if t.count == 0 {
		return EmptyList
	}
	if !t.hamt {
		result := make([]Value, 0, t.count)
		for i := 0; i < len(t.okvs); i += 2 {
			result = append(result, MapEntry{Key: t.okvs[i], Value: t.okvs[i+1]})
		}
		return &MapSeq{entries: result, i: 0}
	}
	if t.root == nil {
		return EmptyList
	}
	mes := t.root.nodeSeq()
	result := make([]Value, len(mes))
	for i, e := range mes {
		result[i] = e
	}
	return &MapSeq{entries: result, i: 0}
}

func (t *TransientMap) Count() Value  { return MakeInt(t.count) }
func (t *TransientMap) RawCount() int { return t.count }

func (t *TransientMap) Contains(key Value) Boolean {
	if !t.hamt {
		return Boolean(t.okvsIndex(key) >= 0)
	}
	if t.root == nil {
		return FALSE
	}
	_, found := t.root.find(0, hashValue(key), key)
	return Boolean(found)
}

func (t *TransientMap) Invoke(args []Value) (Value, error) {
	switch len(args) {
	case 1:
		return t.ValueAt(args[0]), nil
	case 2:
		return t.ValueAtOr(args[0], args[1]), nil
	}
	return NIL, fmt.Errorf("transient map invoke expects 1 or 2 args")
}
func (t *TransientMap) Arity() int { return 1 }

// TransientVector is a mutable version of ArrayVector/PersistentVector.
type TransientVector struct {
	edit  atomic.Bool
	array []Value
}

// NewTransientVectorOfNils returns an editable transient vector of n NIL
// slots in a single allocation. Callers that need an indexed scratch table
// (ir.direct's per-call value array) would otherwise build it as
// (vec (repeat n nil)) — a lazy seq walked element by element — which costs
// one cons per slot before the transient copy even starts.
func NewTransientVectorOfNils(n int) *TransientVector {
	array := make([]Value, n)
	for i := range array {
		array[i] = NIL
	}
	t := &TransientVector{array: array}
	t.edit.Store(true)
	return t
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

func (t *TransientVector) Type() ValueType { return TransientVectorType }
func (t *TransientVector) Unbox() any      { return t }
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

// Pop removes the last element from the transient vector.
func (t *TransientVector) Pop() (*TransientVector, error) {
	if err := t.ensureEditable(); err != nil {
		return nil, err
	}
	if len(t.array) == 0 {
		return nil, fmt.Errorf("can't pop empty transient vector")
	}
	t.array = t.array[:len(t.array)-1]
	return t, nil
}

// Persistent freezes the transient and returns an immutable vector.
func (t *TransientVector) Persistent() (Value, error) {
	if err := t.ensureEditable(); err != nil {
		return NIL, err
	}
	t.edit.Store(false)
	if len(t.array) <= arrayVectorPromotionThreshold {
		result := make(ArrayVector, len(t.array))
		copy(result, t.array)
		return result, nil
	}
	return NewPersistentVector(t.array), nil
}

// Nth implements Indexed: positional access by integer index. TransientVector
// is positional but not seqable, so `nth` must reach it through this rather
// than seq traversal.
func (t *TransientVector) Nth(i int) Value { return t.ValueAt(Int(i)) }

func (t *TransientVector) ValueAt(key Value) Value {
	idx, ok := key.(Int)
	if !ok || int(idx) < 0 || int(idx) >= len(t.array) {
		return NIL
	}
	return t.array[int(idx)]
}

func (t *TransientVector) ValueAtOr(key Value, notFound Value) Value {
	idx, ok := key.(Int)
	if !ok || int(idx) < 0 || int(idx) >= len(t.array) {
		return notFound
	}
	return t.array[int(idx)]
}

func (t *TransientVector) Contains(key Value) Boolean {
	idx, ok := key.(Int)
	if !ok {
		return FALSE
	}
	return Boolean(int(idx) >= 0 && int(idx) < len(t.array))
}

func (t *TransientVector) Invoke(args []Value) (Value, error) {
	if len(args) != 1 {
		return NIL, fmt.Errorf("transient vector invoke expects 1 arg")
	}
	idx, ok := args[0].(Int)
	if !ok {
		return NIL, fmt.Errorf("transient vector key must be Int")
	}
	i := int(idx)
	if i < 0 || i >= len(t.array) {
		return NIL, fmt.Errorf("index out of bounds: %d", i)
	}
	return t.array[i], nil
}
func (t *TransientVector) Arity() int { return 1 }

func (t *TransientVector) Count() Value  { return MakeInt(len(t.array)) }
func (t *TransientVector) RawCount() int { return len(t.array) }

// --- Type metadata ---

type theTransientMapType struct{}

func (t *theTransientMapType) String() string  { return t.Name() }
func (t *theTransientMapType) Type() ValueType { return TypeType }
func (t *theTransientMapType) Unbox() any      { return nil }
func (t *theTransientMapType) Name() string    { return "let-go.lang.TransientMap" }
func (t *theTransientMapType) Box(bare any) (Value, error) {
	return NIL, NewTypeError(bare, "can't be boxed as", t)
}

var TransientMapType *theTransientMapType = &theTransientMapType{}

type theTransientVectorType struct{}

func (t *theTransientVectorType) String() string  { return t.Name() }
func (t *theTransientVectorType) Type() ValueType { return TypeType }
func (t *theTransientVectorType) Unbox() any      { return nil }
func (t *theTransientVectorType) Name() string    { return "let-go.lang.TransientVector" }
func (t *theTransientVectorType) Box(bare any) (Value, error) {
	return NIL, NewTypeError(bare, "can't be boxed as", t)
}

var TransientVectorType *theTransientVectorType = &theTransientVectorType{}

// TransientSet is a mutable version of PersistentSet for batch ops.
type TransientSet struct {
	edit atomic.Bool
	tm   *TransientMap // value is sentinel; reuse map machinery
}

var transientSetSentinel = Boolean(true)

func NewTransientSet(s *PersistentSet) *TransientSet {
	var tm *TransientMap
	if s.impl != nil {
		tm = NewTransientMap(s.impl)
	} else {
		tm = NewTransientMap(EmptyPersistentMap)
	}
	t := &TransientSet{tm: tm}
	t.edit.Store(true)
	return t
}

func (t *TransientSet) ensureEditable() error {
	if !t.edit.Load() {
		return fmt.Errorf("transient used after persistent! call")
	}
	return nil
}

func (t *TransientSet) Type() ValueType { return TransientSetType }
func (t *TransientSet) Unbox() any      { return t }
func (t *TransientSet) String() string {
	return fmt.Sprintf("<transient-set count=%d>", t.tm.count)
}

func (t *TransientSet) Conj(value Value) (*TransientSet, error) {
	if err := t.ensureEditable(); err != nil {
		return nil, err
	}
	_, _ = t.tm.Assoc(value, transientSetSentinel)
	return t, nil
}

func (t *TransientSet) Disj(value Value) (*TransientSet, error) {
	if err := t.ensureEditable(); err != nil {
		return nil, err
	}
	_, _ = t.tm.Dissoc(value)
	return t, nil
}

func (t *TransientSet) Persistent() (*PersistentSet, error) {
	if err := t.ensureEditable(); err != nil {
		return nil, err
	}
	t.edit.Store(false)
	impl, err := t.tm.Persistent()
	if err != nil {
		return nil, err
	}
	return &PersistentSet{impl: impl}, nil
}

func (t *TransientSet) ValueAt(key Value) Value {
	if t.tm.Contains(key) {
		return key
	}
	return NIL
}

func (t *TransientSet) ValueAtOr(key Value, notFound Value) Value {
	if t.tm.Contains(key) {
		return key
	}
	return notFound
}

func (t *TransientSet) Contains(key Value) Boolean { return t.tm.Contains(key) }

func (t *TransientSet) Invoke(args []Value) (Value, error) {
	if len(args) != 1 {
		return NIL, fmt.Errorf("transient set invoke expects 1 arg")
	}
	return t.ValueAt(args[0]), nil
}
func (t *TransientSet) Arity() int { return 1 }

func (t *TransientSet) Count() Value  { return MakeInt(t.tm.count) }
func (t *TransientSet) RawCount() int { return t.tm.count }

type theTransientSetType struct{}

func (t *theTransientSetType) String() string  { return t.Name() }
func (t *theTransientSetType) Type() ValueType { return TypeType }
func (t *theTransientSetType) Unbox() any      { return nil }
func (t *theTransientSetType) Name() string    { return "let-go.lang.TransientSet" }
func (t *theTransientSetType) Box(bare any) (Value, error) {
	return NIL, NewTypeError(bare, "can't be boxed as", t)
}

var TransientSetType *theTransientSetType = &theTransientSetType{}
