/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 *
 * SHA-1 and SHA-256 digests as `hash/sha1` and `hash/sha256`.
 *
 * No build tags: crypto/sha1 and crypto/sha256 are pure Go, so tinygo/wasm
 * builds get these too (unlike xxh3, which needs the murmur3 fallback).
 *
 * The installer seeds the `hash` namespace from Go; pkg/rt/core/hash.lg
 * later declares `(ns hash)` and EXTENDS this namespace rather than
 * replacing it (LookupOrRegisterNS returns the registered instance).
 *
 * let-go strings are raw Go byte strings, so these digest the input's bytes
 * — a jar slurped from disk hashes to the same hex `sha1sum` prints.
 */

package rt

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/nooga/let-go/pkg/vm"
)

func init() { RegisterInstaller(installHashShaNS) }

// shaArg validates the single String argument shared by both digests.
func shaArg(name string, vs []vm.Value) (string, error) {
	if len(vs) != 1 {
		return "", fmt.Errorf("%s expects 1 arg", name)
	}
	s, ok := vs[0].(vm.String)
	if !ok {
		return "", fmt.Errorf("%s expected String", name)
	}
	return string(s), nil
}

func installHashShaNS() {
	ns := vm.NewNamespace("hash")

	// hash/sha1 — (hash/sha1 s) → lowercase hex SHA-1 of s's bytes.
	ns.Def("sha1", mustWrap(func(vs []vm.Value) (vm.Value, error) {
		s, err := shaArg("hash/sha1", vs)
		if err != nil {
			return vm.NIL, err
		}
		sum := sha1.Sum([]byte(s))
		return vm.String(hex.EncodeToString(sum[:])), nil
	}))

	// hash/sha256 — (hash/sha256 s) → lowercase hex SHA-256 of s's bytes.
	ns.Def("sha256", mustWrap(func(vs []vm.Value) (vm.Value, error) {
		s, err := shaArg("hash/sha256", vs)
		if err != nil {
			return vm.NIL, err
		}
		sum := sha256.Sum256([]byte(s))
		return vm.String(hex.EncodeToString(sum[:])), nil
	}))

	RegisterNS(ns)
}
