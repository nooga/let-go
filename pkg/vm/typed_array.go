/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"fmt"
	"reflect"
	"strings"
)

// ArrayKind discriminates the element type of a TypedArray.
type ArrayKind int

const (
	ArrayByte   ArrayKind = iota // backing: []byte
	ArrayInt                     // backing: []int64
	ArrayFloat                   // backing: []float64
	ArrayObject                  // backing: []Value
)

type theTypedArrayType struct{}

func (t *theTypedArrayType) String() string  { return t.Name() }
func (t *theTypedArrayType) Type() ValueType { return TypeType }
func (t *theTypedArrayType) Unbox() any      { return reflect.TypeFor[*theTypedArrayType]() }
func (t *theTypedArrayType) Name() string    { return "let-go.lang.Array" }
func (t *theTypedArrayType) Box(bare any) (Value, error) {
	return NIL, NewTypeError(bare, "can't be boxed as", t)
}

// TypedArrayType is the singleton type for all typed arrays.
var TypedArrayType *theTypedArrayType = &theTypedArrayType{}

// TypedArray is a mutable, typed array backed by a native Go slice.
// Unlike persistent collections, arrays support in-place mutation via Set.
//
// One typed field per kind rather than a single `data any`: the typed
// accessors for lowered code (AtFloat64 and friends, below) must stay under
// the Go inliner's budget to be worth emitting, and an interface assertion
// per access put them over it (cost 88 against a budget of 80). Exactly one
// field is live, the one `kind` names; the others stay nil.
type TypedArray struct {
	kind   ArrayKind
	bytes  []byte
	ints   []int64
	floats []float64
	objs   []Value
}

// --- Constructors ---

func NewByteArray(size int) *TypedArray {
	return &TypedArray{kind: ArrayByte, bytes: make([]byte, size)}
}

func NewByteArrayFrom(data []byte) *TypedArray {
	return &TypedArray{kind: ArrayByte, bytes: data}
}

func NewIntArray(size int) *TypedArray {
	return &TypedArray{kind: ArrayInt, ints: make([]int64, size)}
}

func NewIntArrayFrom(data []int64) *TypedArray {
	return &TypedArray{kind: ArrayInt, ints: data}
}

func NewFloatArray(size int) *TypedArray {
	return &TypedArray{kind: ArrayFloat, floats: make([]float64, size)}
}

func NewFloatArrayFrom(data []float64) *TypedArray {
	return &TypedArray{kind: ArrayFloat, floats: data}
}

func NewObjectArray(size int) *TypedArray {
	d := make([]Value, size)
	for i := range d {
		d[i] = NIL
	}
	return &TypedArray{kind: ArrayObject, objs: d}
}

func NewObjectArrayFrom(data []Value) *TypedArray {
	return &TypedArray{kind: ArrayObject, objs: data}
}

// --- Value interface ---

func (a *TypedArray) Type() ValueType { return TypedArrayType }

// Unbox returns the underlying Go slice directly for interop: one of []byte,
// []int64, []float64, []Value, per Kind.
func (a *TypedArray) Unbox() any {
	switch a.kind {
	case ArrayByte:
		return a.bytes
	case ArrayInt:
		return a.ints
	case ArrayFloat:
		return a.floats
	case ArrayObject:
		return a.objs
	}
	return nil
}

// Kind returns the element kind.
func (a *TypedArray) Kind() ArrayKind { return a.kind }

func (a *TypedArray) String() string {
	b := &strings.Builder{}
	n := a.Len()
	b.WriteString("#" + a.kindName() + "[")
	for i := range n {
		if i > 0 {
			b.WriteRune(' ')
		}
		b.WriteString(a.Get(i).String())
	}
	b.WriteRune(']')
	return b.String()
}

func (a *TypedArray) kindName() string {
	switch a.kind {
	case ArrayByte:
		return "byte-array"
	case ArrayInt:
		return "int-array"
	case ArrayFloat:
		return "double-array"
	case ArrayObject:
		return "object-array"
	}
	return "array"
}

// Meta implements IMeta — arrays don't carry metadata.
func (a *TypedArray) Meta() Value { return NIL }

// WithMeta implements IMeta — returns self (arrays don't support metadata).
func (a *TypedArray) WithMeta(_ Value) Value { return a }

// --- Element access ---

func (a *TypedArray) Len() int {
	switch a.kind {
	case ArrayByte:
		return len(a.bytes)
	case ArrayInt:
		return len(a.ints)
	case ArrayFloat:
		return len(a.floats)
	case ArrayObject:
		return len(a.objs)
	}
	return 0
}

// Get returns the element at index i as a boxed Value.
func (a *TypedArray) Get(i int) Value {
	switch a.kind {
	case ArrayByte:
		return MakeInt(int(a.bytes[i]))
	case ArrayInt:
		return MakeInt64(a.ints[i])
	case ArrayFloat:
		return Float(a.floats[i])
	case ArrayObject:
		return a.objs[i]
	}
	return NIL
}

// Set sets the element at index i, coercing v to the element type.
func (a *TypedArray) Set(i int, v Value) error {
	switch a.kind {
	case ArrayByte:
		n, ok := v.(Int)
		if !ok {
			return fmt.Errorf("byte-array expects Int, got %s", v.Type().Name())
		}
		a.bytes[i] = byte(n)
	case ArrayInt:
		switch n := v.(type) {
		case Int:
			a.ints[i] = int64(n)
		case *BigInt:
			v64, ok := n.ToInt64()
			if !ok {
				return fmt.Errorf("bigint too large for int-array")
			}
			a.ints[i] = v64
		default:
			return fmt.Errorf("int-array expects Int, got %s", v.Type().Name())
		}
	case ArrayFloat:
		f, ok := ToFloat(v)
		if !ok {
			return fmt.Errorf("double-array expects numeric, got %s", v.Type().Name())
		}
		a.floats[i] = f
	case ArrayObject:
		a.objs[i] = v
	}
	return nil
}

// Clone returns a shallow copy.
func (a *TypedArray) Clone() *TypedArray {
	switch a.kind {
	case ArrayByte:
		dst := make([]byte, len(a.bytes))
		copy(dst, a.bytes)
		return &TypedArray{kind: ArrayByte, bytes: dst}
	case ArrayInt:
		dst := make([]int64, len(a.ints))
		copy(dst, a.ints)
		return &TypedArray{kind: ArrayInt, ints: dst}
	case ArrayFloat:
		dst := make([]float64, len(a.floats))
		copy(dst, a.floats)
		return &TypedArray{kind: ArrayFloat, floats: dst}
	case ArrayObject:
		dst := make([]Value, len(a.objs))
		copy(dst, a.objs)
		return &TypedArray{kind: ArrayObject, objs: dst}
	}
	return nil
}

// --- Counted interface ---

func (a *TypedArray) Count() Value      { return Int(a.Len()) }
func (a *TypedArray) RawCount() int     { return a.Len() }
func (a *TypedArray) Empty() Collection { return NewObjectArray(0) }
func (a *TypedArray) Conj(v Value) Collection {
	// Conj on an array creates a new object-array with element appended
	n := a.Len()
	vals := make([]Value, n+1)
	for i := range n {
		vals[i] = a.Get(i)
	}
	vals[n] = v
	return NewObjectArrayFrom(vals)
}

// --- Sequable interface ---

func (a *TypedArray) Seq() Seq {
	if a.Len() == 0 {
		return nil
	}
	return &TypedArraySeq{arr: a, i: 0}
}

// --- Lookup interface (for nth/get fast path) ---

// Nth implements Indexed: positional access by integer index.
func (a *TypedArray) Nth(i int) Value { return a.ValueAt(Int(i)) }

func (a *TypedArray) ValueAt(key Value) Value {
	return a.ValueAtOr(key, NIL)
}

func (a *TypedArray) ValueAtOr(key Value, dflt Value) Value {
	idx, ok := key.(Int)
	if !ok || int(idx) < 0 || int(idx) >= a.Len() {
		return dflt
	}
	return a.Get(int(idx))
}

// --- Fn interface (arrays as functions of their index) ---

func (a *TypedArray) Arity() int { return 1 }

func (a *TypedArray) Invoke(args []Value) (Value, error) {
	if len(args) != 1 {
		return NIL, fmt.Errorf("wrong number of arguments %d", len(args))
	}
	idx, ok := args[0].(Int)
	if !ok {
		return NIL, fmt.Errorf("array index must be Int")
	}
	i := int(idx)
	if i < 0 || i >= a.Len() {
		return NIL, fmt.Errorf("array index %d out of bounds for length %d", i, a.Len())
	}
	return a.Get(i), nil
}

// ============================================================
// TypedArraySeq — lightweight seq view over a TypedArray
// ============================================================

type TypedArraySeq struct {
	arr *TypedArray
	i   int
}

func (s *TypedArraySeq) Type() ValueType        { return ListType }
func (s *TypedArraySeq) Unbox() any             { return s }
func (s *TypedArraySeq) Meta() Value            { return NIL }
func (s *TypedArraySeq) WithMeta(_ Value) Value { return s }

func (s *TypedArraySeq) First() Value {
	if s.i >= s.arr.Len() {
		return NIL
	}
	return s.arr.Get(s.i)
}

func (s *TypedArraySeq) More() Seq {
	if s.i+1 >= s.arr.Len() {
		return EmptyList
	}
	return &TypedArraySeq{arr: s.arr, i: s.i + 1}
}

func (s *TypedArraySeq) Next() Seq {
	if s.i+1 >= s.arr.Len() {
		return nil
	}
	return &TypedArraySeq{arr: s.arr, i: s.i + 1}
}

func (s *TypedArraySeq) Cons(val Value) Seq {
	return NewCons(val, s)
}

func (s *TypedArraySeq) Seq() Seq { return s }

func (s *TypedArraySeq) Count() Value      { return Int(s.arr.Len() - s.i) }
func (s *TypedArraySeq) RawCount() int     { return s.arr.Len() - s.i }
func (s *TypedArraySeq) Empty() Collection { return EmptyList }
func (s *TypedArraySeq) Conj(val Value) Collection {
	return s.Cons(val).(*List)
}

// Nth implements Indexed: positional access by integer index.
func (s *TypedArraySeq) Nth(i int) Value { return s.ValueAt(Int(i)) }

func (s *TypedArraySeq) ValueAt(key Value) Value {
	return s.ValueAtOr(key, NIL)
}

func (s *TypedArraySeq) ValueAtOr(key Value, dflt Value) Value {
	idx, ok := key.(Int)
	if !ok || idx < 0 {
		return dflt
	}
	absIdx := s.i + int(idx)
	if absIdx >= s.arr.Len() {
		return dflt
	}
	return s.arr.Get(absIdx)
}

func (s *TypedArraySeq) String() string {
	b := &strings.Builder{}
	b.WriteRune('(')
	for i := s.i; i < s.arr.Len(); i++ {
		if i > s.i {
			b.WriteRune(' ')
		}
		b.WriteString(s.arr.Get(i).String())
	}
	b.WriteRune(')')
	return b.String()
}

// --- Typed element access for lowered code ---
//
// ir.lower-go emits these where the array kind is known statically (a
// ^doubles/^longs param hint, or a double-array/int-array constructor in the
// same function), so an aget/aset is one bounds check and the access, with no
// boxing of the index or the element. The kind is still checked rather than
// trusted: the static kind is an inference, and a mismatch must surface as an
// error the lowered function returns, never a panic. A wrong kind finds its
// field nil (length 0), takes the cold path, and accessError names the kind.
//
// Each accessor keeps its hot path inlinable: one unsigned bounds compare and
// the access. Check with `go build -gcflags=-m=2 ./pkg/vm | grep AtFloat64`
// after editing.

// accessError is the cold path shared by the typed accessors: a kind mismatch
// (the array is not of kind `want`) or, for the right kind, an index out of
// bounds, worded as CoreAgetf/CoreAsetf word it.
func (a *TypedArray) accessError(i int64, want ArrayKind) error {
	if a.kind != want {
		return fmt.Errorf("%s expected, got %s", (&TypedArray{kind: want}).kindName(), a.kindName())
	}
	return fmt.Errorf("array index %d out of bounds for length %d", i, a.Len())
}

// AtFloat64 reads element i of a double-array.
func (a *TypedArray) AtFloat64(i int64) (float64, error) {
	if uint64(i) >= uint64(len(a.floats)) {
		return 0, a.accessError(i, ArrayFloat)
	}
	return a.floats[i], nil
}

// SetFloat64 writes element i of a double-array.
func (a *TypedArray) SetFloat64(i int64, v float64) error {
	if uint64(i) >= uint64(len(a.floats)) {
		return a.accessError(i, ArrayFloat)
	}
	a.floats[i] = v
	return nil
}

// AtInt64 reads element i of an int-array.
func (a *TypedArray) AtInt64(i int64) (int64, error) {
	if uint64(i) >= uint64(len(a.ints)) {
		return 0, a.accessError(i, ArrayInt)
	}
	return a.ints[i], nil
}

// SetInt64 writes element i of an int-array.
func (a *TypedArray) SetInt64(i int64, v int64) error {
	if uint64(i) >= uint64(len(a.ints)) {
		return a.accessError(i, ArrayInt)
	}
	a.ints[i] = v
	return nil
}

// AtByte reads element i of a byte-array, widened to int64 (aget on a
// byte-array returns an Int).
func (a *TypedArray) AtByte(i int64) (int64, error) {
	if uint64(i) >= uint64(len(a.bytes)) {
		return 0, a.accessError(i, ArrayByte)
	}
	return int64(a.bytes[i]), nil
}

// SetByte writes element i of a byte-array, truncating as Set does.
func (a *TypedArray) SetByte(i int64, v int64) error {
	if uint64(i) >= uint64(len(a.bytes)) {
		return a.accessError(i, ArrayByte)
	}
	a.bytes[i] = byte(v)
	return nil
}

// AtValue reads element i of an object-array.
func (a *TypedArray) AtValue(i int64) (Value, error) {
	if uint64(i) >= uint64(len(a.objs)) {
		return NIL, a.accessError(i, ArrayObject)
	}
	return a.objs[i], nil
}

// SetValue writes element i of an object-array.
func (a *TypedArray) SetValue(i int64, v Value) error {
	if uint64(i) >= uint64(len(a.objs)) {
		return a.accessError(i, ArrayObject)
	}
	a.objs[i] = v
	return nil
}
