/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

// Package builtins is the single-source home for direct-call-eligible
// clojure.core builtins: context-free functions with uniform vm.Value
// signatures. Each function here is meant to be BOTH the interpreter
// registration (lang.go boxes it) AND the AOT direct-call target (lang.go
// registers it as a native module, generated code calls builtins.X directly).
//
// It imports ONLY vm (never rt) so rt/lang.go can import it without a cycle —
// that layering is what lets lang.go register from these funcs instead of
// duplicating their bodies (the corefns/diff-test pattern this replaces).
package builtins

import (
	"fmt"

	"github.com/nooga/let-go/pkg/vm"
)

// seqOf mirrors the unexported pkg/rt/lang.go helper of the same name.
func seqOf(v vm.Value) (vm.Seq, error) {
	if v == vm.NIL || v == vm.EmptyList {
		return nil, nil
	}
	if _, isLazy := v.(*vm.LazySeq); !isLazy {
		if sq, ok := v.(vm.Sequable); ok {
			s := sq.Seq()
			if s == nil || s == vm.EmptyList {
				return nil, nil
			}
			return s, nil
		}
	}
	if s, ok := v.(vm.Seq); ok {
		return s, nil
	}
	return nil, fmt.Errorf("don't know how to create ISeq from %s", v.Type())
}

// Vector mirrors clojure.core/vector — `(vector & xs)`. Builtin is
// vm.NativeFnType.WrapNoErr(vm.NewArrayVector); same delegation as List.
func Vector(args ...vm.Value) (vm.Value, error) {
	return vm.NewArrayVector(args), nil
}

// Not mirrors clojure.core/not — `(not x)`.
func Not(v vm.Value) (vm.Value, error) {
	return vm.Boolean(!vm.IsTruthy(v)), nil
}

// Cons mirrors clojure.core/cons — `(cons x coll)`.
func Cons(elem, coll vm.Value) (vm.Value, error) {
	if coll == vm.NIL {
		return vm.NewCons(elem, nil), nil
	}
	seq, err := seqOf(coll)
	if err != nil {
		return vm.NIL, fmt.Errorf("cons expected Seq")
	}
	if seq == nil {
		return vm.NewCons(elem, nil), nil
	}
	return vm.NewCons(elem, seq), nil
}

// lg:native
// lg:name contains?
// Contains mirrors clojure.core/contains? — `(contains? coll k)`.
func Contains(coll, k vm.Value) (vm.Value, error) {
	if coll == vm.NIL {
		return vm.FALSE, nil
	}
	if s, ok := coll.(vm.Keyed); ok {
		return s.Contains(k), nil
	}
	if s, ok := coll.(vm.String); ok {
		if idx, ok := k.(vm.Int); ok {
			return vm.Boolean(int(idx) >= 0 && int(idx) < len([]rune(string(s)))), nil
		}
		return vm.NIL, fmt.Errorf("contains? on a string requires an integer key, got %s", k.Type().Name())
	}
	if _, isSeq := coll.(vm.Seq); isSeq {
		return vm.NIL, fmt.Errorf("contains? not supported on %s", coll.Type().Name())
	}
	indexed, ok := coll.(vm.Indexed)
	if !ok {
		return vm.NIL, fmt.Errorf("contains? not supported on %s", coll.Type().Name())
	}
	idx, ok := k.(vm.Int)
	if !ok {
		return vm.NIL, fmt.Errorf("contains? on an indexed value requires an integer key, got %s", k.Type().Name())
	}
	i := int(idx)
	return vm.Boolean(i >= 0 && i < indexed.RawCount()), nil
}

// lg:native
// HashSet mirrors clojure.core/hash-set — `(hash-set & xs)`.
func HashSet(args ...vm.Value) (vm.Value, error) {
	return vm.NewSet(args), nil
}

// lg:native
// lg:name array-map
// ArrayMap mirrors clojure.core/array-map — `(array-map & kvs)`.
func ArrayMap(args ...vm.Value) (vm.Value, error) {
	if len(args)%2 != 0 {
		return vm.NIL, fmt.Errorf("array-map requires an even number of arguments, got %d", len(args))
	}
	return vm.NewArrayMap(args), nil
}

// lg:native
// Vec mirrors clojure.core/vec — `(vec coll)`.
func Vec(coll vm.Value) (vm.Value, error) {
	if coll == vm.NIL || coll == vm.EmptyList {
		return vm.ArrayVector{}, nil
	}
	// Empty string → empty vector
	if s, ok := coll.(vm.String); ok && len(string(s)) == 0 {
		return vm.ArrayVector{}, nil
	}
	if v, ok := coll.(vm.ArrayVector); ok {
		return v, nil
	}
	if a, ok := coll.(*vm.TypedArray); ok && a.Kind() == vm.ArrayObject {
		return vm.ArrayVector(a.Unbox().([]vm.Value)), nil
	}
	seq, err := seqOf(coll)
	if err != nil {
		return vm.NIL, err
	}
	// Realize lazy seqs to check emptiness
	if ls, ok := seq.(*vm.LazySeq); ok {
		seq = ls.Seq()
	}
	if seq == nil || seq == vm.EmptyList {
		return vm.ArrayVector{}, nil
	}
	ret := []vm.Value{}
	for seq != nil {
		ret = append(ret, seq.First())
		seq = seq.Next()
	}
	return vm.NewArrayVector(ret), nil
}

// lg:native
// Nth mirrors clojure.core/nth — `(nth coll i)`.
// Takes vm.Value arguments and handles type conversion internally,
// avoiding the type-coercion bottleneck of the primitive rt.Nth.
func Nth(coll, i vm.Value) (vm.Value, error) {
	// (nth nil i) => nil, matching rt.Nth / clojure.core/nth — nil is an empty
	// coll for nth. An empty list is out of bounds, like an empty vector.
	if coll == vm.NIL {
		return vm.NIL, nil
	}
	if coll == vm.EmptyList {
		return vm.NIL, fmt.Errorf("nth: index out of bounds")
	}

	// Convert index to int
	var idx int
	switch iv := i.(type) {
	case vm.Int:
		idx = int(iv)
	default:
		return vm.NIL, fmt.Errorf("nth: index must be numeric, got %s", i.Type().Name())
	}

	// Negative index
	if idx < 0 {
		return vm.NIL, fmt.Errorf("nth: negative index %d", idx)
	}

	// Handle strings
	if s, ok := coll.(vm.String); ok {
		runes := []rune(string(s))
		if idx >= len(runes) {
			return vm.NIL, fmt.Errorf("nth: index %d out of bounds for string", idx)
		}
		return vm.Char(runes[idx]), nil
	}

	// Handle indexed collections (vectors, arrays)
	if indexed, ok := coll.(vm.Indexed); ok {
		count := indexed.RawCount()
		if idx >= count {
			return vm.NIL, fmt.Errorf("nth: index %d out of bounds", idx)
		}
		return indexed.Nth(idx), nil
	}

	// Handle sequences. Realize a lazy seq at the head and each step (mirrors
	// rt.forceSeq): an EMPTY lazy seq is non-nil until forced, so First() on it
	// would wrongly yield nil instead of an out-of-bounds error.
	seq, err := seqOf(coll)
	if err != nil {
		return vm.NIL, fmt.Errorf("nth: cannot get sequence from %s", coll.Type().Name())
	}
	if ls, ok := seq.(*vm.LazySeq); ok {
		seq = ls.Seq()
	}
	if seq == nil {
		return vm.NIL, fmt.Errorf("nth: index %d out of bounds", idx)
	}
	for j := 0; j < idx; j++ {
		seq = seq.Next()
		if ls, ok := seq.(*vm.LazySeq); ok {
			seq = ls.Seq()
		}
		if seq == nil {
			return vm.NIL, fmt.Errorf("nth: index %d out of bounds", idx)
		}
	}
	return seq.First(), nil
}

// lg:native
// lg:name nth
// Nth3 mirrors clojure.core/nth — `(nth coll i notFound)`.
// Takes vm.Value arguments and handles type conversion internally.
func Nth3(coll, i, notFound vm.Value) (vm.Value, error) {
	// Handle nil/empty collection
	if coll == vm.NIL || coll == vm.EmptyList {
		return notFound, nil
	}

	// Convert index to int. A non-numeric index is an ERROR, not a not-found:
	// canonical nth coerces the index to int, and (nth coll bad not-found) still
	// throws on a bad index type — not-found is only for out-of-bounds.
	var idx int
	switch iv := i.(type) {
	case vm.Int:
		idx = int(iv)
	default:
		return vm.NIL, fmt.Errorf("nth: index must be numeric, got %s", i.Type().Name())
	}

	// Negative index
	if idx < 0 {
		return notFound, nil
	}

	// Handle strings
	if s, ok := coll.(vm.String); ok {
		runes := []rune(string(s))
		if idx >= len(runes) {
			return notFound, nil
		}
		return vm.Char(runes[idx]), nil
	}

	// Handle indexed collections (vectors, arrays)
	if indexed, ok := coll.(vm.Indexed); ok {
		count := indexed.RawCount()
		if idx >= count {
			return notFound, nil
		}
		return indexed.Nth(idx), nil
	}

	// Handle sequences. Realize a lazy seq at the head and each step so an EMPTY
	// lazy seq falls through to notFound instead of yielding a nil First().
	seq, err := seqOf(coll)
	if err != nil {
		return notFound, nil
	}
	if ls, ok := seq.(*vm.LazySeq); ok {
		seq = ls.Seq()
	}
	if seq == nil {
		return notFound, nil
	}
	for j := 0; j < idx; j++ {
		seq = seq.Next()
		if ls, ok := seq.(*vm.LazySeq); ok {
			seq = ls.Seq()
		}
		if seq == nil {
			return notFound, nil
		}
	}
	return seq.First(), nil
}

// lg:native
// lg:name seq?
// IsSeq mirrors clojure.core/seq? — true only for values that are actually
// vm.Seq. The explicit negative cases come first because several collection
// types satisfy vm.Seq structurally while Clojure reports seq? false for them.
// Single source of truth: rt's ns.Def("seq?") delegates here.
func IsSeq(v vm.Value) (vm.Value, error) {
	switch v.(type) {
	case *vm.Nil, vm.String, *vm.TypedArray,
		vm.ArrayVector, *vm.PersistentVector,
		*vm.PersistentMap, *vm.PersistentSet,
		*vm.SortedSet, *vm.SortedMap:
		return vm.FALSE, nil
	}
	_, ok := v.(vm.Seq)
	return vm.Boolean(ok), nil
}
