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

func (t *theTypedArrayType) String() string     { return t.Name() }
func (t *theTypedArrayType) Type() ValueType    { return TypeType }
func (t *theTypedArrayType) Unbox() interface{} { return reflect.TypeOf(t) }
func (t *theTypedArrayType) Name() string       { return "let-go.lang.Array" }
func (t *theTypedArrayType) Box(bare interface{}) (Value, error) {
	return NIL, NewTypeError(bare, "can't be boxed as", t)
}

// TypedArrayType is the singleton type for all typed arrays.
var TypedArrayType *theTypedArrayType = &theTypedArrayType{}

// TypedArray is a mutable, typed array backed by a native Go slice.
// Unlike persistent collections, arrays support in-place mutation via Set.
type TypedArray struct {
	kind ArrayKind
	data interface{} // one of: []byte, []int64, []float64, []Value
}

// --- Constructors ---

func NewByteArray(size int) *TypedArray {
	return &TypedArray{kind: ArrayByte, data: make([]byte, size)}
}

func NewByteArrayFrom(data []byte) *TypedArray {
	return &TypedArray{kind: ArrayByte, data: data}
}

func NewIntArray(size int) *TypedArray {
	return &TypedArray{kind: ArrayInt, data: make([]int64, size)}
}

func NewIntArrayFrom(data []int64) *TypedArray {
	return &TypedArray{kind: ArrayInt, data: data}
}

func NewFloatArray(size int) *TypedArray {
	return &TypedArray{kind: ArrayFloat, data: make([]float64, size)}
}

func NewFloatArrayFrom(data []float64) *TypedArray {
	return &TypedArray{kind: ArrayFloat, data: data}
}

func NewObjectArray(size int) *TypedArray {
	d := make([]Value, size)
	for i := range d {
		d[i] = NIL
	}
	return &TypedArray{kind: ArrayObject, data: d}
}

func NewObjectArrayFrom(data []Value) *TypedArray {
	return &TypedArray{kind: ArrayObject, data: data}
}

// --- Value interface ---

func (a *TypedArray) Type() ValueType { return TypedArrayType }

// Unbox returns the underlying Go slice directly for interop.
func (a *TypedArray) Unbox() interface{} { return a.data }

// Kind returns the element kind.
func (a *TypedArray) Kind() ArrayKind { return a.kind }

func (a *TypedArray) String() string {
	b := &strings.Builder{}
	n := a.Len()
	switch a.kind {
	case ArrayByte:
		b.WriteString("#byte-array[")
	case ArrayInt:
		b.WriteString("#int-array[")
	case ArrayFloat:
		b.WriteString("#double-array[")
	case ArrayObject:
		b.WriteString("#object-array[")
	}
	for i := 0; i < n; i++ {
		if i > 0 {
			b.WriteRune(' ')
		}
		b.WriteString(a.Get(i).String())
	}
	b.WriteRune(']')
	return b.String()
}

// Meta implements IMeta — arrays don't carry metadata.
func (a *TypedArray) Meta() Value { return NIL }

// WithMeta implements IMeta — returns self (arrays don't support metadata).
func (a *TypedArray) WithMeta(_ Value) Value { return a }

// --- Element access ---

func (a *TypedArray) Len() int {
	switch a.kind {
	case ArrayByte:
		return len(a.data.([]byte))
	case ArrayInt:
		return len(a.data.([]int64))
	case ArrayFloat:
		return len(a.data.([]float64))
	case ArrayObject:
		return len(a.data.([]Value))
	}
	return 0
}

// Get returns the element at index i as a boxed Value.
func (a *TypedArray) Get(i int) Value {
	switch a.kind {
	case ArrayByte:
		return MakeInt(int(a.data.([]byte)[i]))
	case ArrayInt:
		return MakeInt(int(a.data.([]int64)[i]))
	case ArrayFloat:
		return Float(a.data.([]float64)[i])
	case ArrayObject:
		return a.data.([]Value)[i]
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
		a.data.([]byte)[i] = byte(n)
	case ArrayInt:
		switch n := v.(type) {
		case Int:
			a.data.([]int64)[i] = int64(n)
		case *BigInt:
			v64, ok := n.ToInt64()
			if !ok {
				return fmt.Errorf("bigint too large for int-array")
			}
			a.data.([]int64)[i] = v64
		default:
			return fmt.Errorf("int-array expects Int, got %s", v.Type().Name())
		}
	case ArrayFloat:
		f, ok := ToFloat(v)
		if !ok {
			return fmt.Errorf("double-array expects numeric, got %s", v.Type().Name())
		}
		a.data.([]float64)[i] = f
	case ArrayObject:
		a.data.([]Value)[i] = v
	}
	return nil
}

// Clone returns a shallow copy.
func (a *TypedArray) Clone() *TypedArray {
	switch a.kind {
	case ArrayByte:
		src := a.data.([]byte)
		dst := make([]byte, len(src))
		copy(dst, src)
		return &TypedArray{kind: ArrayByte, data: dst}
	case ArrayInt:
		src := a.data.([]int64)
		dst := make([]int64, len(src))
		copy(dst, src)
		return &TypedArray{kind: ArrayInt, data: dst}
	case ArrayFloat:
		src := a.data.([]float64)
		dst := make([]float64, len(src))
		copy(dst, src)
		return &TypedArray{kind: ArrayFloat, data: dst}
	case ArrayObject:
		src := a.data.([]Value)
		dst := make([]Value, len(src))
		copy(dst, src)
		return &TypedArray{kind: ArrayObject, data: dst}
	}
	return nil
}

// --- Counted interface ---

func (a *TypedArray) Count() Value    { return Int(a.Len()) }
func (a *TypedArray) RawCount() int   { return a.Len() }
func (a *TypedArray) Empty() Collection { return NewObjectArray(0) }
func (a *TypedArray) Conj(v Value) Collection {
	// Conj on an array creates a new object-array with element appended
	n := a.Len()
	vals := make([]Value, n+1)
	for i := 0; i < n; i++ {
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

func (s *TypedArraySeq) Type() ValueType    { return ListType }
func (s *TypedArraySeq) Unbox() interface{} { return s }
func (s *TypedArraySeq) Meta() Value        { return NIL }
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

func (s *TypedArraySeq) Count() Value  { return Int(s.arr.Len() - s.i) }
func (s *TypedArraySeq) RawCount() int { return s.arr.Len() - s.i }
func (s *TypedArraySeq) Empty() Collection { return EmptyList }
func (s *TypedArraySeq) Conj(val Value) Collection {
	return s.Cons(val).(*List)
}

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
