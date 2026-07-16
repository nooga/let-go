/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"math"
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

func invokeIdentical(t *testing.T, left, right vm.Value) bool {
	t.Helper()

	identicalVar := LookupCoreVar("identical?")
	if identicalVar == nil {
		t.Fatal("core/identical? not found")
	}
	identicalFn, ok := identicalVar.Deref().(vm.Fn)
	if !ok {
		t.Fatal("core/identical? is not an Fn")
	}

	got, err := identicalFn.Invoke([]vm.Value{left, right})
	if err != nil {
		t.Fatalf("identical? returned error: %v", err)
	}
	result, ok := got.(vm.Boolean)
	if !ok {
		t.Fatalf("identical? returned %T, want vm.Boolean", got)
	}
	return bool(result)
}

func TestIdenticalUsesRepresentationForNonComparableValues(t *testing.T) {
	tests := []struct {
		name string
		make func() vm.Value
	}{
		{
			name: "array vector",
			make: func() vm.Value {
				return vm.ArrayVector{vm.Int(1), vm.Int(2)}
			},
		},
		{
			name: "map",
			make: func() vm.Value {
				return vm.Map{vm.Keyword("key"): vm.Int(1)}
			},
		},
		{
			name: "set",
			make: func() vm.Value {
				return vm.Set{vm.Int(1): {}}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			left := test.make()
			if !invokeIdentical(t, left, left) {
				t.Fatal("same representation should be identical")
			}

			right := test.make()
			if invokeIdentical(t, left, right) {
				t.Fatal("distinct equal representations should not be identical")
			}
		})
	}
}

func TestIdenticalArrayVectorRequiresMatchingSliceHeader(t *testing.T) {
	backing := vm.ArrayVector{vm.Int(1), vm.Int(2), vm.Int(3)}
	left := backing[:2:2]
	same := left
	if !invokeIdentical(t, left, same) {
		t.Fatal("copied slice header should be identical")
	}

	differentCapacity := backing[:2:3]
	if invokeIdentical(t, left, differentCapacity) {
		t.Fatal("slices sharing a first element but with different capacities should not be identical")
	}

	differentLength := backing[:1:2]
	if invokeIdentical(t, left, differentLength) {
		t.Fatal("slices sharing a first element but with different lengths should not be identical")
	}
}

func TestRawMapAndSetPointersAreStableAndDistinct(t *testing.T) {
	leftMap := vm.Map{vm.Keyword("key"): vm.Int(1)}
	rightMap := vm.Map{vm.Keyword("key"): vm.Int(1)}
	leftMapPtr := mapPtr(leftMap)
	if leftMapPtr == 0 {
		t.Fatal("non-nil map identity should be non-zero")
	}
	if got := mapPtr(leftMap); got != leftMapPtr {
		t.Fatalf("map identity changed: got %#x, want %#x", got, leftMapPtr)
	}
	if rightMapPtr := mapPtr(rightMap); rightMapPtr == leftMapPtr {
		t.Fatalf("distinct maps share identity %#x", leftMapPtr)
	}

	leftSet := vm.Set{vm.Int(1): {}}
	rightSet := vm.Set{vm.Int(1): {}}
	leftSetPtr := setPtr(leftSet)
	if leftSetPtr == 0 {
		t.Fatal("non-nil set identity should be non-zero")
	}
	if got := setPtr(leftSet); got != leftSetPtr {
		t.Fatalf("set identity changed: got %#x, want %#x", got, leftSetPtr)
	}
	if rightSetPtr := setPtr(rightSet); rightSetPtr == leftSetPtr {
		t.Fatalf("distinct sets share identity %#x", leftSetPtr)
	}
}

func makePersistentVectorIdentityFixture() vm.PersistentVector {
	values := make([]vm.Value, 40)
	for i := range values {
		values[i] = vm.Int(i)
	}
	return vm.NewPersistentVector(values).(vm.PersistentVector)
}

func TestIdenticalPersistentVectorUsesObjectIdentity(t *testing.T) {
	left := makePersistentVectorIdentityFixture()
	if !invokeIdentical(t, left, left) {
		t.Fatal("same persistent vector binding should be identical")
	}

	right := makePersistentVectorIdentityFixture()
	if invokeIdentical(t, left, right) {
		t.Fatal("separately constructed equal persistent vectors should not be identical")
	}
}

func TestPersistentVectorValuePointerIsStableAndDistinct(t *testing.T) {
	left := makePersistentVectorIdentityFixture()
	leftPtr, ok := valuePtr(left)
	if !ok || leftPtr == 0 {
		t.Fatal("persistent vector should expose a non-zero identity key")
	}
	if got, ok := valuePtr(left); !ok || got != leftPtr {
		t.Fatalf("persistent vector identity changed: got %#x, want %#x", got, leftPtr)
	}

	right := makePersistentVectorIdentityFixture()
	rightPtr, ok := valuePtr(right)
	if !ok || rightPtr == 0 {
		t.Fatal("distinct persistent vector should expose a non-zero identity key")
	}
	if rightPtr == leftPtr {
		t.Fatalf("distinct persistent vectors share identity %#x", leftPtr)
	}
}

func TestIdenticalPreservesComparableAndPointerBehavior(t *testing.T) {
	if !invokeIdentical(t, vm.Keyword("same"), vm.Keyword("same")) {
		t.Fatal("equal comparable singleton values should be identical")
	}
	if invokeIdentical(t, vm.Keyword("left"), vm.Keyword("right")) {
		t.Fatal("different comparable values should not be identical")
	}

	nan := vm.Float(math.NaN())
	if invokeIdentical(t, nan, nan) {
		t.Fatal("NaN should retain Go equality behavior")
	}

	atom := vm.NewAtom(vm.Int(1))
	if !invokeIdentical(t, atom, atom) {
		t.Fatal("same atom pointer should be identical")
	}
	if invokeIdentical(t, atom, vm.NewAtom(vm.Int(1))) {
		t.Fatal("distinct atoms should not be identical")
	}

	persistentMap := vm.NewPersistentMap([]vm.Value{vm.Keyword("key"), vm.Int(1)})
	if !invokeIdentical(t, persistentMap, persistentMap) {
		t.Fatal("same persistent map pointer should be identical")
	}
	if invokeIdentical(t, persistentMap, vm.NewPersistentMap([]vm.Value{vm.Keyword("key"), vm.Int(1)})) {
		t.Fatal("distinct equal persistent maps should not be identical")
	}

	persistentSet := vm.NewPersistentSet([]vm.Value{vm.Int(1)})
	if !invokeIdentical(t, persistentSet, persistentSet) {
		t.Fatal("same persistent set pointer should be identical")
	}
	if invokeIdentical(t, persistentSet, vm.NewPersistentSet([]vm.Value{vm.Int(1)})) {
		t.Fatal("distinct equal persistent sets should not be identical")
	}
}
