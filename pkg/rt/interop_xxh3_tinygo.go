//go:build tinygo

/*
 * Copyright (c) 2026 Matt Parrett
 * SPDX-License-Identifier: MIT
 *
 * Pure-Go hash backing the `xxh3` namespace under TinyGo (native AND wasm).
 *
 * WHY NOT REAL xxh3: two independent walls. (1) The lginterop-generated
 * interop_xxh3.go (//go:build !tinygo) boxes typed Go funcs via reflection,
 * which TinyGo can't call. (2) zeebo/xxh3 is assembly-only on arm64/amd64
 * (accumNEON / accumSSE, no purego mode even in v1.1.0, the latest) — TinyGo
 * can't link Go assembly, so a native TinyGo build won't link it at all. Only
 * "other" arches (wasm) get xxh3's pure-Go path, and even there we'd rather not
 * carry the dep. So under TinyGo we back the `xxh3` namespace with a small
 * pure-Go hash instead: FNV-1a to absorb the bytes, then a splitmix64/Murmur3
 * finalizer for avalanche.
 *
 * NOT BIT-COMPATIBLE with real xxh3. Cross-platform determinism parity is
 * deliberately set aside — a TinyGo build is its own determinism domain (its
 * worlds won't match a stock-let-go build's). The result is masked to 31 bits
 * so it lands in let-go's Int identically on native (64-bit int) and wasm
 * (32-bit int) and is always non-negative — i.e. consistent across the whole
 * TinyGo domain. Good enough for game-seed determinism (xsofy.hash →
 * xxh3/HashSeed); it is not for security or cross-implementation reproducibility.
 * The namespace keeps the `xxh3` name so callers (xsofy) build unchanged.
 */

package rt

import "github.com/nooga/let-go/pkg/vm"

// fallbackHash64 mixes bytes into a seed. FNV-1a absorbs the input; a
// splitmix64/Murmur3 finalizer then diffuses it so small input changes avalanche
// across the whole word. Pure uint64 — no assembly, no platform int width.
func fallbackHash64(b []byte, seed uint64) uint64 {
	const (
		fnvOffset uint64 = 0xcbf29ce484222325
		fnvPrime  uint64 = 0x100000001b3
	)
	h := seed ^ fnvOffset
	for _, c := range b {
		h ^= uint64(c)
		h *= fnvPrime
	}
	h ^= h >> 33
	h *= 0xff51afd7ed558ccd
	h ^= h >> 33
	h *= 0xc4ceb9fe1a85ec53
	h ^= h >> 33
	return h
}

// boxHash masks to 31 bits so the value fits let-go's Int identically on native
// (64-bit int) and wasm (32-bit int) TinyGo, and is always non-negative.
func boxHash(h uint64) vm.Value { return vm.Int(int64(h & 0x7fffffff)) }

// hashSeedArg coerces a let-go integer seed to uint64.
func hashSeedArg(v vm.Value) (uint64, bool) {
	if n, ok := v.(vm.Int); ok {
		return uint64(int64(n)), true
	}
	return 0, false
}

func installXxh3NS() {
	ns := vm.NewNamespace("xxh3")

	hash, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		b, ok := asBytes(vs[0])
		if !ok {
			return vm.NIL, vm.NewTypeError(vs[0], "is not a byte-array/String for", vm.NativeFnType)
		}
		return boxHash(fallbackHash64(b, 0)), nil
	})
	ns.Def("Hash", hash)

	hashSeed, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		b, ok := asBytes(vs[0])
		if !ok {
			return vm.NIL, vm.NewTypeError(vs[0], "is not a byte-array/String for", vm.NativeFnType)
		}
		seed, ok := hashSeedArg(vs[1])
		if !ok {
			return vm.NIL, vm.NewTypeError(vs[1], "is not an integer seed for", vm.NativeFnType)
		}
		return boxHash(fallbackHash64(b, seed)), nil
	})
	ns.Def("HashSeed", hashSeed)

	hashString, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		s, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, vm.NewTypeError(vs[0], "is not a String for", vm.NativeFnType)
		}
		return boxHash(fallbackHash64([]byte(string(s)), 0)), nil
	})
	ns.Def("HashString", hashString)

	hashStringSeed, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		s, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, vm.NewTypeError(vs[0], "is not a String for", vm.NativeFnType)
		}
		seed, ok := hashSeedArg(vs[1])
		if !ok {
			return vm.NIL, vm.NewTypeError(vs[1], "is not an integer seed for", vm.NativeFnType)
		}
		return boxHash(fallbackHash64([]byte(string(s)), seed)), nil
	})
	ns.Def("HashStringSeed", hashStringSeed)

	RegisterNS(ns)
}

func init() { RegisterInstaller(installXxh3NS) }
