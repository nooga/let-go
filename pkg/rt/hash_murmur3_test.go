/*
 * Copyright (c) 2026 Matt Parrett
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

// Vectors verified against two independent references: python mmh3 (wraps
// Appleby's reference C++) and twmb/murmur3. Both agree on every case; the
// 64-bit seed case is twmb SeedSum128(seed, seed), the h1=h2=seed extension
// (mmh3 seeds are uint32-only).
func TestMurmur3Sum64ReferenceVectors(t *testing.T) {
	cases := []struct {
		in   string
		seed uint64
		want uint64
	}{
		{"", 0, 0x0000000000000000},
		{"hello", 0, 0xcbd8a7b341bd9b02},
		{"hello, world", 0, 0x342fac623a5ebc8e},
		{"The quick brown fox jumps over the lazy dog", 0, 0xe34bbc7bbc071b6c},
		{"xsofy", 0xdeadbeef, 0x2470512280477ca6},
		{"seed-hash", 0x1122334455667788, 0xaba3dc1297abb8d7},
		// block-boundary coverage: exactly one 16-byte block, then +1 tail byte
		{"0123456789abcdef", 0, 0x4be06d94cf4ad1a7},
		{"0123456789abcdefX", 7, 0x58550bbc8514f46d},
	}
	for _, c := range cases {
		if got := murmur3Sum64([]byte(c.in), c.seed); got != c.want {
			t.Errorf("murmur3Sum64(%q, %#x) = %#016x, want %#016x", c.in, c.seed, got, c.want)
		}
	}
}

func TestMurmur3BoxHashMasksTo31Bits(t *testing.T) {
	// The boxed value must be identical on 64-bit and 32-bit int platforms
	// and non-negative, so the namespace-level result is the low 31 bits.
	h := murmur3Sum64([]byte("hello"), 0) // 0xcbd8a7b341bd9b02
	want := vm.Int(h & 0x7fffffff)        // 0x41bd9b02
	if got := boxHash(h); got != want {
		t.Errorf("boxHash = %v, want %v", got, want)
	}
}
