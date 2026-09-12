//go:build !tinygo || wasm

/*
 * Copyright (c) 2026 Matt Parrett
 * SPDX-License-Identifier: MIT
 *
 * Non-reflect xxh3 bindings, mirroring the four functions the reflect-boxed
 * `xxh3` namespace (interop_xxh3.go) exposes for seed hashing. That file is
 * lginterop-generated and boxes through reflect, which TinyGo cannot call, so
 * TinyGo builds had no `xxh3` namespace at all. These adapters use the same
 * Wrap path as hash_murmur3.go and are compiled on every build so the stock
 * test suite exercises them; interop_xxh3_tinygo.go registers them as the
 * `xxh3` namespace under the tinygo tag, where the generated file is absent.
 *
 * Results are boxed the way the reflect path boxes a uint64: reinterpreted
 * as int64, no masking, so they are bit-identical to the generated bindings.
 *
 * Build constraint: stock Go always, TinyGo only on wasm. zeebo/xxh3's
 * large-input path is assembly on arm64/amd64 and TinyGo cannot link it, so a
 * native TinyGo build must not import the package at all.
 */

package rt

import (
	"fmt"

	"github.com/nooga/let-go/pkg/vm"
	"github.com/zeebo/xxh3"
)

func boxXxh3(h uint64) vm.Value { return vm.Int(int64(h)) }

func xxh3Hash(vs []vm.Value) (vm.Value, error) {
	if len(vs) != 1 {
		return vm.NIL, fmt.Errorf("xxh3/Hash expects 1 arg, got %d", len(vs))
	}
	b, ok := asBytes(vs[0])
	if !ok {
		return vm.NIL, vm.NewTypeError(vs[0], "is not a byte-array/String for", vm.NativeFnType)
	}
	return boxXxh3(xxh3.Hash(b)), nil
}

func xxh3HashSeed(vs []vm.Value) (vm.Value, error) {
	if len(vs) != 2 {
		return vm.NIL, fmt.Errorf("xxh3/HashSeed expects 2 args, got %d", len(vs))
	}
	b, ok := asBytes(vs[0])
	if !ok {
		return vm.NIL, vm.NewTypeError(vs[0], "is not a byte-array/String for", vm.NativeFnType)
	}
	seed, ok := hashSeedArg(vs[1])
	if !ok {
		return vm.NIL, vm.NewTypeError(vs[1], "is not an integer seed for", vm.NativeFnType)
	}
	return boxXxh3(xxh3.HashSeed(b, seed)), nil
}

func xxh3HashString(vs []vm.Value) (vm.Value, error) {
	if len(vs) != 1 {
		return vm.NIL, fmt.Errorf("xxh3/HashString expects 1 arg, got %d", len(vs))
	}
	s, ok := vs[0].(vm.String)
	if !ok {
		return vm.NIL, vm.NewTypeError(vs[0], "is not a String for", vm.NativeFnType)
	}
	return boxXxh3(xxh3.HashString(string(s))), nil
}

func xxh3HashStringSeed(vs []vm.Value) (vm.Value, error) {
	if len(vs) != 2 {
		return vm.NIL, fmt.Errorf("xxh3/HashStringSeed expects 2 args, got %d", len(vs))
	}
	s, ok := vs[0].(vm.String)
	if !ok {
		return vm.NIL, vm.NewTypeError(vs[0], "is not a String for", vm.NativeFnType)
	}
	seed, ok := hashSeedArg(vs[1])
	if !ok {
		return vm.NIL, vm.NewTypeError(vs[1], "is not an integer seed for", vm.NativeFnType)
	}
	return boxXxh3(xxh3.HashStringSeed(string(s), seed)), nil
}
