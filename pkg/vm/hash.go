/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

const (
	fnvOffset32 = uint32(2166136261)
	fnvPrime32  = uint32(16777619)
)

// Hashable is implemented by types that cache their hash for fast map lookups.
type Hashable interface {
	Hash() uint32
}

// hashValue computes a 32-bit hash for any Value type.
// Checks for Hashable first (cached hash), then falls back to computing.
func hashValue(v Value) uint32 {
	if h, ok := v.(Hashable); ok {
		return h.Hash()
	}
	return computeHash(v)
}

// computeHash is the fallback for types that don't implement Hashable.
func computeHash(v Value) uint32 {
	if v == NIL {
		return 0
	}
	return hashBytes([]byte(v.String()))
}

func hashBytes(b []byte) uint32 {
	h := fnvOffset32
	for _, c := range b {
		h ^= uint32(c)
		h *= fnvPrime32
	}
	return h
}

// hashString hashes a string without allocating a []byte copy.
func hashString(s string) uint32 {
	h := fnvOffset32
	for i := 0; i < len(s); i++ {
		h ^= uint32(s[i])
		h *= fnvPrime32
	}
	return h
}

func hashUint64(v uint64) uint32 {
	// Murmur3 finalizer
	v ^= v >> 33
	v *= 0xff51afd7ed558ccd
	v ^= v >> 33
	v *= 0xc4ceb9fe1a85ec53
	v ^= v >> 33
	return uint32(v)
}

// --- Hash combining (Murmur3-style, matching Clojure's approach) ---

// hashOrdered computes a hash for an ordered collection (vector, list).
// Matches Clojure's Murmur3.hashOrdered.
func hashOrdered(seq Seq) uint32 {
	h := uint32(1)
	for s := seq; s != nil; s = s.Next() {
		h = 31*h + hashValue(s.First())
	}
	return mixFinish(h)
}

// hashUnordered computes a hash for an unordered collection (map, set).
// Matches Clojure's Murmur3.hashUnordered — order-independent via XOR+addition.
func hashUnordered(seq Seq) uint32 {
	var h uint32
	for s := seq; s != nil; s = s.Next() {
		h += hashValue(s.First())
	}
	return mixFinish(h)
}

// mixFinish is Murmur3's fmix32.
func mixFinish(h uint32) uint32 {
	h ^= h >> 16
	h *= 0x85ebca6b
	h ^= h >> 13
	h *= 0xc2b2ae35
	h ^= h >> 16
	return h
}

// valueEquiv tests if two Values are equivalent for map key purposes.
// Uses hash as a fast negative check, then structural comparison.
func valueEquiv(a, b Value) bool {
	// Fast path: pointer/value identity
	if isComparable(a) && isComparable(b) {
		if a == b {
			return true
		}
	}
	// Numeric cross-type: Int(1) == Float(1.0)
	if IsNumber(a) && IsNumber(b) {
		return NumEq(a, b)
	}
	// Fast negative: if both are Hashable and hashes differ, not equal
	ha, aOk := a.(Hashable)
	hb, bOk := b.(Hashable)
	if aOk && bOk {
		if ha.Hash() != hb.Hash() {
			return false
		}
	}
	// Structural comparison
	if eq, ok := a.(interface{ Equals(Value) bool }); ok {
		return eq.Equals(b)
	}
	return false
}

// isComparable returns true if the Value can be safely compared with ==.
func isComparable(v Value) bool {
	switch v.(type) {
	case Int, Float, String, Keyword, Symbol, Boolean, Char, *Nil, *Var, *Namespace:
		return true
	default:
		// Type objects (singletons) are pointer-comparable
		if _, ok := v.(ValueType); ok {
			return true
		}
		return false
	}
}
