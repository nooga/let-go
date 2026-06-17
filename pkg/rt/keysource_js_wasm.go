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
 * SAB protocol — SPSC ring buffer between the main thread (producer) and
 * this worker (consumer). Layout MUST stay in sync with the producer in
 * pkg/rt/wasm/lg-host-core.js; the consts below mirror the JS side. The
 * producer writes slot [writeIdx % keyCapacity] then advances writeIdx; the
 * consumer waits on writeIdx, reads slot [readIdx % keyCapacity], then
 * advances readIdx. The producer never blocks the main thread — when the
 * ring is full it drops.
 *
 * Int32Array view (used cells):
 *
 *	[0]  readIdx   — consumer increments after reading a slot
 *	[1]  writeIdx  — producer increments after writing a slot
 *
 * Uint8Array view (full SAB):
 *
 *	bytes 64..71   slot lengths (1 byte per slot, 8 slots)
 *	bytes 72..199  slot keys    (16 bytes per slot, 8 slots)
 */

package rt

import (
	"fmt"
	"syscall/js"
)

const (
	keyCapacity  = 8  // ring slots
	keyMaxLen    = 16 // bytes per key
	keyLenOffset = 64 // byte offset of the lengths region
	keyOffset    = 72 // byte offset of the keys region
	readIdxCell  = 0  // Int32 index of the read pointer
	writeIdxCell = 1  // Int32 index of the write pointer
)

// HostKeySource reads keys from the bundle's SharedArrayBuffer ring.
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

	// Park while the ring is empty (writeIdx == readIdx). Atomics.wait
	// returns ok / timed-out / not-equal; in every case we re-check the
	// predicate, so spurious or early wakes just loop and re-wait. Passing
	// the observed writeIdx as the expected value means a producer that
	// advanced before we hit wait returns immediately with "not-equal".
	r := jsLoadInt(atomics, keyInt32, readIdxCell)
	for {
		w := jsLoadInt(atomics, keyInt32, writeIdxCell)
		if w != r {
			break
		}
		atomics.Call("wait", keyInt32, writeIdxCell, w)
	}

	slot := r % keyCapacity
	keyLen := int(byte(keyUint8.Index(keyLenOffset + slot).Int()))
	if keyLen <= 0 || keyLen > keyMaxLen {
		// Drain the bad slot and bail. Never block on it.
		atomics.Call("store", keyInt32, readIdxCell, r+1)
		return "", nil
	}

	// Critical ordering: copy slot bytes into a local Go slice BEFORE
	// advancing readIdx. If the ring is full, advancing first would invite
	// the producer to overwrite slot[r%capacity] while we're still reading.
	keyBytes := make([]byte, keyLen)
	keyBase := keyOffset + slot*keyMaxLen
	for i := 0; i < keyLen; i++ {
		keyBytes[i] = byte(keyUint8.Index(keyBase + i).Int())
	}

	// Tier-1 lazy coalesce: peek r+1, r+2, ... while the next slot matches
	// both length and bytes, and drain past all matches in a single readIdx
	// advance. A held key that outruns the game loop becomes one turn
	// instead of a backlog of queued turns.
	finalR := r
	for {
		nextR := finalR + 1
		w := jsLoadInt(atomics, keyInt32, writeIdxCell)
		if nextR >= w {
			break
		}
		nextSlot := nextR % keyCapacity
		if int(byte(keyUint8.Index(keyLenOffset+nextSlot).Int())) != keyLen {
			break
		}
		matches := true
		nextBase := keyOffset + nextSlot*keyMaxLen
		for i := 0; i < keyLen; i++ {
			if byte(keyUint8.Index(nextBase+i).Int()) != keyBytes[i] {
				matches = false
				break
			}
		}
		if !matches {
			break
		}
		finalR = nextR
	}

	atomics.Call("store", keyInt32, readIdxCell, finalR+1)
	return string(keyBytes), nil
}

func (HostKeySource) KeyPending() bool {
	keyInt32 := js.Global().Get("_lgKeyInt32")
	atomics := js.Global().Get("Atomics")
	if keyInt32.IsUndefined() || atomics.IsUndefined() {
		return false
	}
	// Non-consuming: true iff at least one slot is queued in the ring.
	r := jsLoadInt(atomics, keyInt32, readIdxCell)
	w := jsLoadInt(atomics, keyInt32, writeIdxCell)
	return w > r
}
