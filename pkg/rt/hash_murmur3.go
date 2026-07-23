/*
 * Copyright (c) 2026 Matt Parrett
 * SPDX-License-Identifier: MIT
 *
 * MurmurHash3 (x64 128-bit variant, low word) as the `murmur3` namespace.
 *
 * This is the portable seed-hash for builds that can't have real xxh3:
 * TinyGo can neither call the reflect-boxed bindings in interop_xxh3.go nor
 * link zeebo/xxh3's arm64/amd64 assembly (accumNEON / accumSSE — no purego
 * mode even in v1.1.0, the latest). TinyGo builds therefore have no `xxh3`
 * namespace at all; this one is registered on ALL builds — not just TinyGo —
 * so the implementation is exercised by the stock test suite (the
 * tinygo-tagged files are not), and a program that calls murmur3/* gets
 * identical values on stock Go, TinyGo native, and TinyGo wasm.
 *
 * Results are masked to 31 bits so they land in let-go's Int identically on
 * 64-bit (native) and 32-bit (TinyGo wasm) int widths and are always
 * non-negative. Deterministic and portable; not for security.
 *
 * Algorithm: Austin Appleby's public-domain MurmurHash3_x64_128; we return
 * the low word (h1). The seed is accepted as a 64-bit value with
 * h1 = h2 = seed — bit-identical to the reference for seeds in uint32
 * range, and to twmb/murmur3's SeedSum128(seed, seed) beyond it (verified
 * against both in hash_murmur3_test.go).
 */

package rt

import (
	"encoding/binary"
	"math/bits"

	"github.com/nooga/let-go/pkg/vm"
)

func fmix64(k uint64) uint64 {
	k ^= k >> 33
	k *= 0xff51afd7ed558ccd
	k ^= k >> 33
	k *= 0xc4ceb9fe1a85ec53
	k ^= k >> 33
	return k
}

// murmur3Sum64 is MurmurHash3_x64_128 returning the low word (h1).
func murmur3Sum64(data []byte, seed uint64) uint64 {
	const (
		c1 uint64 = 0x87c37b91114253d5
		c2 uint64 = 0x4cf5ad432745937f
	)
	h1, h2 := seed, seed
	n := len(data)

	nblocks := n / 16
	for i := 0; i < nblocks; i++ {
		k1 := binary.LittleEndian.Uint64(data[i*16:])
		k2 := binary.LittleEndian.Uint64(data[i*16+8:])

		k1 *= c1
		k1 = bits.RotateLeft64(k1, 31)
		k1 *= c2
		h1 ^= k1

		h1 = bits.RotateLeft64(h1, 27)
		h1 += h2
		h1 = h1*5 + 0x52dce729

		k2 *= c2
		k2 = bits.RotateLeft64(k2, 33)
		k2 *= c1
		h2 ^= k2

		h2 = bits.RotateLeft64(h2, 31)
		h2 += h1
		h2 = h2*5 + 0x38495ab5
	}

	tail := data[nblocks*16:]
	var k1, k2 uint64
	switch len(tail) {
	case 15:
		k2 ^= uint64(tail[14]) << 48
		fallthrough
	case 14:
		k2 ^= uint64(tail[13]) << 40
		fallthrough
	case 13:
		k2 ^= uint64(tail[12]) << 32
		fallthrough
	case 12:
		k2 ^= uint64(tail[11]) << 24
		fallthrough
	case 11:
		k2 ^= uint64(tail[10]) << 16
		fallthrough
	case 10:
		k2 ^= uint64(tail[9]) << 8
		fallthrough
	case 9:
		k2 ^= uint64(tail[8])
		k2 *= c2
		k2 = bits.RotateLeft64(k2, 33)
		k2 *= c1
		h2 ^= k2
		fallthrough
	case 8:
		k1 ^= uint64(tail[7]) << 56
		fallthrough
	case 7:
		k1 ^= uint64(tail[6]) << 48
		fallthrough
	case 6:
		k1 ^= uint64(tail[5]) << 40
		fallthrough
	case 5:
		k1 ^= uint64(tail[4]) << 32
		fallthrough
	case 4:
		k1 ^= uint64(tail[3]) << 24
		fallthrough
	case 3:
		k1 ^= uint64(tail[2]) << 16
		fallthrough
	case 2:
		k1 ^= uint64(tail[1]) << 8
		fallthrough
	case 1:
		k1 ^= uint64(tail[0])
		k1 *= c1
		k1 = bits.RotateLeft64(k1, 31)
		k1 *= c2
		h1 ^= k1
	}

	h1 ^= uint64(n)
	h2 ^= uint64(n)
	h1 += h2
	h2 += h1
	h1 = fmix64(h1)
	h2 = fmix64(h2)
	h1 += h2
	// h2 += h1 completes the reference finalizer; we only return h1.
	return h1
}

// boxHash masks to 31 bits so the value fits let-go's Int identically on
// 64-bit (native) and 32-bit (TinyGo wasm) int widths, and is always
// non-negative.
func boxHash(h uint64) vm.Value { return vm.Int(int64(h & 0x7fffffff)) }

// hashSeedArg coerces a let-go integer seed to uint64. BigInt is accepted
// because on 32-bit-int builds (TinyGo wasm) an integer literal above 2^31-1
// reads as BigInt — the same seed source must work at both int widths.
func hashSeedArg(v vm.Value) (uint64, bool) {
	switch n := v.(type) {
	case vm.Int:
		return uint64(int64(n)), true
	case *vm.BigInt:
		if i, ok := n.ToInt64(); ok {
			return uint64(i), true
		}
	}
	return 0, false
}

func installMurmur3NS() {
	ns := vm.NewNamespace("murmur3")

	ns.Def("Hash", mustWrap(func(vs []vm.Value) (vm.Value, error) {
		b, ok := asBytes(vs[0])
		if !ok {
			return vm.NIL, vm.NewTypeError(vs[0], "is not a byte-array/String for", vm.NativeFnType)
		}
		return boxHash(murmur3Sum64(b, 0)), nil
	}))

	ns.Def("HashSeed", mustWrap(func(vs []vm.Value) (vm.Value, error) {
		b, ok := asBytes(vs[0])
		if !ok {
			return vm.NIL, vm.NewTypeError(vs[0], "is not a byte-array/String for", vm.NativeFnType)
		}
		seed, ok := hashSeedArg(vs[1])
		if !ok {
			return vm.NIL, vm.NewTypeError(vs[1], "is not an integer seed for", vm.NativeFnType)
		}
		return boxHash(murmur3Sum64(b, seed)), nil
	}))

	ns.Def("HashString", mustWrap(func(vs []vm.Value) (vm.Value, error) {
		s, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, vm.NewTypeError(vs[0], "is not a String for", vm.NativeFnType)
		}
		return boxHash(murmur3Sum64([]byte(string(s)), 0)), nil
	}))

	ns.Def("HashStringSeed", mustWrap(func(vs []vm.Value) (vm.Value, error) {
		s, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, vm.NewTypeError(vs[0], "is not a String for", vm.NativeFnType)
		}
		seed, ok := hashSeedArg(vs[1])
		if !ok {
			return vm.NIL, vm.NewTypeError(vs[1], "is not an integer seed for", vm.NativeFnType)
		}
		return boxHash(murmur3Sum64([]byte(string(s)), seed)), nil
	}))

	RegisterNS(ns)
}

func init() { RegisterInstaller(installMurmur3NS) }
