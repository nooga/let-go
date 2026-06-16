//go:build js && wasm

/*
 * Copyright (c) 2026 Matt Parrett
 * SPDX-License-Identifier: MIT
 *
 * HostKeySource is the WASM key source: term/read-key and key-pending? read
 * keystrokes from the bundle's SharedArrayBuffer instead of poking the
 * _lgKeyInt32 / Atomics globals inline. The dual of HostWriter / HostEmitter
 * — installed at the *keys* root by term_wasm.go's install, resolving the JS
 * globals per call so SAB setup order doesn't matter.
 *
 * SAB layout (Int32Array): [0] key-ready flag, [1] byte count; Uint8Array at
 * offset 8, length 16, holds the key bytes. See pkg/rt/wasm/lg-host.js.
 */

package rt

import (
	"fmt"
	"syscall/js"
)

// HostKeySource reads keys from the bundle's SharedArrayBuffer.
type HostKeySource struct{}

// NewHostKeySource returns a HostKeySource for the *keys* root binding.
func NewHostKeySource() *HostKeySource { return &HostKeySource{} }

func (HostKeySource) ReadKey() (string, error) {
	atomics := js.Global().Get("Atomics")
	keyInt32 := js.Global().Get("_lgKeyInt32")
	keyUint8 := js.Global().Get("_lgKeyUint8")
	if keyInt32.IsUndefined() || keyUint8.IsUndefined() {
		return "", fmt.Errorf("read-key: terminal input not available (no SharedArrayBuffer)")
	}

	// Flush output before parking on the key wait.
	js.Global().Call("_lgFlush")

	// Block until a key is ready.
	atomics.Call("wait", keyInt32, 0, 0)

	keyLen := jsLoadInt(atomics, keyInt32, 1)
	if keyLen <= 0 || keyLen > 16 {
		atomics.Call("store", keyInt32, 0, 0)
		return "", nil
	}

	keyBytes := make([]byte, keyLen)
	for i := 0; i < keyLen; i++ {
		keyBytes[i] = byte(keyUint8.Index(i).Int())
	}
	atomics.Call("store", keyInt32, 0, 0)
	return string(keyBytes), nil
}

func (HostKeySource) KeyPending() bool {
	keyInt32 := js.Global().Get("_lgKeyInt32")
	atomics := js.Global().Get("Atomics")
	if keyInt32.IsUndefined() || atomics.IsUndefined() {
		return false
	}
	v := atomics.Call("load", keyInt32, 0)
	if v.IsUndefined() || v.IsNull() {
		return false
	}
	return v.Int() != 0
}
