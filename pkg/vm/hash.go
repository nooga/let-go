/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import "unsafe"

const (
	fnvOffset32 = uint32(2166136261)
	fnvPrime32  = uint32(16777619)
)

// hashValue computes a 32-bit hash for any Value type.
// Uses FNV-1a for byte data and Murmur3 finalizer for integers.
// Inlined for zero allocations on primitive types.
func hashValue(v Value) uint32 {
	switch x := v.(type) {
	case Int:
		return hashUint64(uint64(x))
	case Float:
		// Use bit pattern for exact float hashing
		f := float64(x)
		return hashUint64(*(*uint64)(unsafe.Pointer(&f)))
	case String:
		return hashBytes([]byte(string(x)))
	case Keyword:
		return hashBytes([]byte(string(x))) ^ 0x9e3779b9
	case Symbol:
		return hashBytes([]byte(string(x))) ^ 0x517cc1b7
	case Boolean:
		if bool(x) {
			return 1
		}
		return 0
	case Char:
		return hashUint64(uint64(x))
	case *Nil:
		return 0
	default:
		return hashBytes([]byte(v.String()))
	}
}

func hashBytes(b []byte) uint32 {
	h := fnvOffset32
	for _, c := range b {
		h ^= uint32(c)
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

// valueEquiv tests if two Values are equivalent for map key purposes.
// Uses == for primitive types, falls back to Equals interface for compound types.
func valueEquiv(a, b Value) bool {
	// Use a recover-free path: check comparability before using ==
	if isComparable(a) && isComparable(b) {
		if a == b {
			return true
		}
	}
	// Numeric cross-type: Int(1) == Float(1.0) for map lookup
	if IsNumber(a) && IsNumber(b) {
		return NumEq(a, b)
	}
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
		return false
	}
}
