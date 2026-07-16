/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import "testing"

func persistentVectorIdentityFixture() PersistentVector {
	values := make([]Value, 40)
	for i := range values {
		values[i] = Int(i)
	}
	return NewPersistentVector(values).(PersistentVector)
}

func TestPersistentVectorCopiesPreserveIdentity(t *testing.T) {
	original := persistentVectorIdentityFixture()
	copy := original
	if !original.SameIdentity(copy) {
		t.Fatal("copy of a persistent vector should preserve identity")
	}

	distinct := persistentVectorIdentityFixture()
	if original.SameIdentity(distinct) {
		t.Fatal("separate constructors should allocate distinct identities")
	}
}

func TestPersistentVectorOperationsCreateIdentity(t *testing.T) {
	original := persistentVectorIdentityFixture()
	operations := []struct {
		name  string
		apply func(PersistentVector) PersistentVector
	}{
		{
			name: "with meta",
			apply: func(v PersistentVector) PersistentVector {
				return v.WithMeta(Map{Keyword("tag"): Keyword("value")}).(PersistentVector)
			},
		},
		{
			name: "empty",
			apply: func(v PersistentVector) PersistentVector {
				return v.Empty().(PersistentVector)
			},
		},
		{
			name: "conj",
			apply: func(v PersistentVector) PersistentVector {
				return v.Conj(Int(40)).(PersistentVector)
			},
		},
		{
			name: "pop",
			apply: func(v PersistentVector) PersistentVector {
				return v.Pop()
			},
		},
		{
			name: "assoc trie",
			apply: func(v PersistentVector) PersistentVector {
				return v.Assoc(Int(0), Int(100)).(PersistentVector)
			},
		},
		{
			name: "assoc tail",
			apply: func(v PersistentVector) PersistentVector {
				return v.Assoc(Int(39), Int(100)).(PersistentVector)
			},
		},
		{
			name: "assoc append",
			apply: func(v PersistentVector) PersistentVector {
				return v.Assoc(Int(40), Int(40)).(PersistentVector)
			},
		},
	}

	for _, operation := range operations {
		t.Run(operation.name, func(t *testing.T) {
			first := operation.apply(original)
			if original.SameIdentity(first) {
				t.Fatal("logically new vector reused input identity")
			}

			copy := first
			if !first.SameIdentity(copy) {
				t.Fatal("copy of operation result should preserve identity")
			}

			second := operation.apply(original)
			if first.SameIdentity(second) {
				t.Fatal("separate operation results should have distinct identities")
			}
		})
	}
}

func TestPersistentVectorEmptyFromZeroValueCreatesIdentity(t *testing.T) {
	first := (PersistentVector{}).Empty().(PersistentVector)
	copy := first
	if !first.SameIdentity(copy) {
		t.Fatal("empty vector copy should preserve identity")
	}

	second := (PersistentVector{}).Empty().(PersistentVector)
	if first.SameIdentity(second) {
		t.Fatal("separate empty operations should allocate distinct identities")
	}
}
