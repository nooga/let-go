/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"fmt"
	"reflect"
)

// ValueType represents a type of a Value
type ValueType interface {
	Value
	Name() string
	Box(any) (Value, error)
}

// Value is implemented by all LETGO values
type Value interface {
	fmt.Stringer
	Type() ValueType
	Unbox() any
}

// IMeta is implemented by values that support metadata.
type IMeta interface {
	Meta() Value          // returns the metadata map, or NIL
	WithMeta(Value) Value // returns a copy of this value with the given metadata
}

// Seq is implemented by all sequence-like values
type Seq interface {
	Value
	Cons(Value) Seq
	First() Value
	More() Seq
	Next() Seq
}

type Sequable interface {
	Value
	Seq() Seq
}

type Counted interface {
	Value
	RawCount() int
	Count() Value
}

// Collection is implemented by all collections
type Collection interface {
	Value
	Counted
	Empty() Collection
	Conj(Value) Collection
}

type Fn interface {
	Value
	Invoke([]Value) (Value, error)
	Arity() int
}

type Associative interface {
	Value
	Assoc(Value, Value) Associative
	Dissoc(Value) Associative
}

type Lookup interface {
	Value
	ValueAt(Value) Value
	ValueAtOr(Value, Value) Value
}

// KeywordLookup is an optional fast path for keyword-keyed lookup that avoids
// boxing the Keyword into a Value interface (type Keyword string boxes on
// interface conversion — a heap alloc). Keyword.Invoke prefers it when present.
type KeywordLookup interface {
	ValueAtKeyword(k Keyword) Value
	ValueAtKeywordOr(k Keyword, dflt Value) Value
}

// Indexed marks positional collections — those whose elements are addressed by
// a 0-based integer index, as opposed to maps/sets (key-addressable) and lazy
// seqs (sequential-only). `nth` dispatches on this to use an O(1)-ish indexed
// access without seq traversal. Implementations must satisfy: Nth(i) is valid
// for 0 <= i < RawCount(); callers bounds-check via RawCount() before calling.
//
// Crucially this excludes maps/sets (which are Lookup+Counted but NOT positional)
// and transient/chunk types that aren't seqable — so adding a new positional
// type is a matter of implementing this interface, not editing a type switch.
type Indexed interface {
	Value
	Nth(i int) Value
	RawCount() int
}

type Keyed interface {
	Value
	Contains(Value) Boolean
}

type Receiver interface {
	Value
	InvokeMethod(Symbol, []Value) (Value, error)
}

type Named interface {
	Value
	Name() Value
	Namespace() Value
}

type Reference interface {
	Deref() Value
}

// BlockingDeref is a reference whose value may not be available yet (promise,
// future). DerefTimeout blocks up to timeoutMs for the value, returning
// timeoutVal on timeout. Backs Clojure's 3-arg deref.
type BlockingDeref interface {
	DerefTimeout(timeoutMs int64, timeoutVal Value) Value
}

type theTypeType struct{}

var TypeType *theTypeType = &theTypeType{}

func (t *theTypeType) String() string  { return t.Name() }
func (t *theTypeType) Type() ValueType { return t }
func (t *theTypeType) Unbox() any      { return reflect.TypeFor[*theTypeType]() }

func (t *theTypeType) Name() string { return "let-go.lang.Type" }
func (t *theTypeType) Box(b any) (Value, error) {
	return NIL, NewTypeError(b, "can't be boxed as", t)
}

type theAnyType struct{}

var AnyType *theAnyType = &theAnyType{}

func (t *theAnyType) String() string  { return t.Name() }
func (t *theAnyType) Type() ValueType { return TypeType }
func (t *theAnyType) Unbox() any      { return reflect.TypeFor[*theAnyType]() }
func (t *theAnyType) Name() string    { return "java.lang.Object" }
func (t *theAnyType) Box(b any) (Value, error) {
	return NIL, NewTypeError(b, "can't be boxed as", t)
}

func ToLetGo(v any) (Value, error) {
	return BoxValue(reflect.ValueOf(v))
}

// MustBox boxes a Go value via BoxValue and panics on error. Intended for
// use at init time (e.g. namespace install functions) where boxing failure
// is a programmer error, not a runtime condition.
func MustBox(v any) Value {
	val, err := BoxValue(reflect.ValueOf(v))
	if err != nil {
		panic("vm.MustBox: " + err.Error())
	}
	return val
}

func BoxValue(v reflect.Value) (Value, error) {
	if !v.IsValid() {
		return NIL, NewTypeError(v, "can't be boxed", nil)
	}
	// A reflect.Value carrying a nil-typed interface (e.g. a Go fn returned
	// `(vm.Value)(nil)` instead of vm.NIL) reaches the default branch of the
	// switch below and panics inside NewBoxed → reflect.TypeOf(nil).Name().
	// Treat all nil-interface returns as let-go NIL up front.
	if v.IsValid() && v.Kind() == reflect.Interface && v.IsNil() {
		return NIL, nil
	}
	if v.CanInterface() {
		rv, ok := v.Interface().(Value)
		if ok {
			if rv == nil {
				return NIL, nil
			}
			return rv, nil
		}
	}
	switch v.Type().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return Int(v.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return Int(v.Uint()), nil
	case reflect.String:
		return StringType.Box(v.Interface())
	case reflect.Float32, reflect.Float64:
		return Float(v.Float()), nil
	case reflect.Bool:
		return BooleanType.Box(v.Interface())
	case reflect.Func:
		return NativeFnType.Box(v.Interface())
	case reflect.Struct:
		if v.CanInterface() {
			if m := LookupStructMapping(v.Type()); m != nil {
				return m.StructToRecord(v.Interface()), nil
			}
		}
		if v.CanInterface() {
			return NewBoxed(v.Interface()), nil
		}
		return NIL, NewTypeError(v, "is not boxable", nil)
	case reflect.Pointer:
		if v.IsNil() {
			return NIL, nil
		}
		// Check if pointed-to struct has a registered mapping
		if v.Elem().Kind() == reflect.Struct && v.CanInterface() {
			if m := LookupStructMapping(v.Elem().Type()); m != nil {
				return m.StructToRecord(v.Interface()), nil
			}
		}
		// Wrap non-nil, non-Value pointers as opaque boxed values
		if v.CanInterface() {
			return NewBoxed(v.Interface()), nil
		}
		return NIL, NewTypeError(v, "is not boxable", nil)
	case reflect.Slice, reflect.Array:
		if v.IsNil() {
			return NIL, nil
		}
		switch v.Type().Elem().Kind() {
		case reflect.Uint8:
			return String(v.Bytes()), nil
		case reflect.Int64:
			src := v.Interface().([]int64)
			dst := make([]int64, len(src))
			copy(dst, src)
			return NewIntArrayFrom(dst), nil
		case reflect.Float64:
			src := v.Interface().([]float64)
			dst := make([]float64, len(src))
			copy(dst, src)
			return NewFloatArrayFrom(dst), nil
		}
		in := make([]Value, v.Len())
		for i := 0; i < v.Len(); i++ {
			mv, err := boxByDynamicType(v.Index(i))
			if err != nil {
				return NIL, NewTypeError(v.Index(i), "can't be boxed", nil).Wrap(err)
			}
			in[i] = mv
		}
		return ArrayVector(in), nil
	case reflect.Map:
		if v.IsNil() {
			return NIL, nil
		}
		result := EmptyPersistentMap
		iter := v.MapRange()
		for iter.Next() {
			// Keys and values are boxed by their DYNAMIC type for the same
			// reason slice elements are: the map's STATIC element type of a
			// map[string]any is Interface, so every value would otherwise miss
			// the scalar fast paths and land as an opaque Boxed. map[string]any
			// is the natural Go shape for a JSON object or an options bag,
			// which is what a Wails v3 binding hands back.
			k, err := boxByDynamicType(iter.Key())
			if err != nil {
				return NIL, err
			}
			val, err := boxByDynamicType(iter.Value())
			if err != nil {
				return NIL, err
			}
			result = result.Assoc(k, val).(*PersistentMap)
		}
		return result, nil
	case reflect.Chan:
		if v.IsNil() {
			return NIL, nil
		}
		return ChanType.Box(v.Interface())
	default:
		if v.CanInterface() {
			return NewBoxed(v.Interface()), nil
		}
		return NIL, NewTypeError(v, "is not boxable", nil)
	}
}

// boxByDynamicType boxes a collection element by its DYNAMIC type when
// dynamicBoxSafe admits it, rather than by the static type the collection
// declares. For a []any or a map[string]any that static kind is Interface, so
// every element would otherwise miss the string/int/float paths and land as an
// opaque Boxed — printing as <go.string Ada> and comparing equal to nothing in
// let-go. Those are the natural Go types for a row of unknown column types and
// for a JSON object, so that made the documented database/sql and Wails
// wrapper shapes unusable.
//
// The IsNil test is load-bearing, not defensive: Elem() on a nil interface
// yields the zero reflect.Value, and BoxValue reports an invalid value as an
// error rather than NIL. Leaving a nil element wrapped lets the nil-interface
// guard at the top of BoxValue return NIL — which is what a SQL NULL has to
// become.
func boxByDynamicType(v reflect.Value) (Value, error) {
	e := v
	if e.Kind() == reflect.Interface && !e.IsNil() {
		if d := e.Elem(); dynamicBoxSafe(d.Type()) {
			e = d
		}
	}
	bv, err := BoxValue(e)
	if err != nil && e != v {
		// The dynamic type turned out to be one BoxValue cannot handle — a
		// defined scalar like json.Number or MyBool, whose XType.Box accepts
		// only the exact predeclared type. Boxing the element wrapped yields
		// the opaque Boxed it produced before this unwrapping existed, so an
		// element we cannot improve is preserved rather than failing the whole
		// collection.
		return BoxValue(v)
	}
	return bv, err
}

// dynamicBoxSafe reports whether an interface element's dynamic type is one we
// unwrap and box directly. It is an ALLOWLIST, deliberately: BoxValue reaches
// several paths that assume an exact predeclared type, and some of them panic
// rather than erroring, so an error fallback cannot make an opt-out list safe.
//
//   - the slice/array case calls IsNil, invalid for an array;
//   - the []byte/[]int64/[]float64 fast paths type-assert the exact slice
//     type, so a defined type over the same underlying slice explodes;
//   - ChanType.Box spawns a goroutine that calls Recv, which for a send-only
//     channel panics in ANOTHER goroutine and takes the process down —
//     unrecoverable at this call site;
//   - and composites reach all of the above recursively, so a top-level-only
//     check is not enough ([][1]int panics on its element).
//
// These all predate dynamic boxing — passing such a value straight to BoxValue
// has always behaved this way — but unwrapping would newly route []any
// elements into them. Listing what we positively want keeps everything else
// boxing exactly as it did before: as an opaque Boxed.
//
// []any is included and is what makes nesting work: its elements are
// themselves interfaces, so they recurse through this same allowlist.
//
// map[string]any is included on the same terms and for the same reason: it is
// the natural Go shape for a JSON object or an options bag, its values are
// interfaces that recurse through this allowlist, and BoxValue's map branch
// has no exact-type fast path to explode on. Only a string key and the empty
// interface qualify — any other key type would need its own boxing story, and
// a non-empty interface element can hold anything.
func dynamicBoxSafe(t reflect.Type) bool {
	// A defined type (json.Number, MyBool, MyInts) is never unwrapped: the
	// string and bool branches call XType.Box, which accepts only the exact
	// predeclared type, and the slice fast paths assert it.
	if t.PkgPath() != "" {
		return false
	}
	switch t.Kind() {
	case reflect.String, reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	case reflect.Slice:
		switch t.Elem().Kind() {
		case reflect.Uint8, reflect.Int64, reflect.Float64:
			// Exactly []byte / []int64 / []float64, whose fast paths assert
			// precisely those types.
			return t.Elem().PkgPath() == ""
		case reflect.Interface:
			// []any — elements recurse through this allowlist.
			return t.Elem().NumMethod() == 0
		}
	case reflect.Map:
		// Exactly map[string]any — values recurse through this allowlist. The
		// key's PkgPath test rejects a defined string type, which the string
		// branch's StringType.Box would not accept.
		return t.Key().Kind() == reflect.String && t.Key().PkgPath() == "" &&
			t.Elem().Kind() == reflect.Interface && t.Elem().NumMethod() == 0
	}
	return false
}

func IsTruthy(v Value) bool {
	return v != NIL && v != FALSE
}
