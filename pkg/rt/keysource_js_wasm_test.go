//go:build js && wasm

/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"sort"
	"syscall/js"
	"testing"
	"time"
)

// installKeyRing sets up the globals HostKeySource reads, over a fresh
// SharedArrayBuffer, and returns a producer that pushes one key the way
// lg-host-core.js does: write the slot, then advance writeIdx.
func installKeyRing(t *testing.T) (push func(string), readIdx func() int) {
	t.Helper()
	sab := js.Global().Get("SharedArrayBuffer").New(keyOffset + keyCapacity*keyMaxLen)
	i32 := js.Global().Get("Int32Array").New(sab)
	u8 := js.Global().Get("Uint8Array").New(sab)
	flush := js.FuncOf(func(js.Value, []js.Value) any { return nil })
	js.Global().Set("_lgKeyInt32", i32)
	js.Global().Set("_lgKeyUint8", u8)
	js.Global().Set("_lgFlush", flush)
	t.Cleanup(func() {
		js.Global().Delete("_lgKeyInt32")
		js.Global().Delete("_lgKeyUint8")
		js.Global().Delete("_lgFlush")
		flush.Release()
	})
	atomics := js.Global().Get("Atomics")
	push = func(k string) {
		w := atomics.Call("load", i32, writeIdxCell).Int()
		slot := w % keyCapacity
		u8.SetIndex(keyLenOffset+slot, len(k))
		for i := 0; i < len(k); i++ {
			u8.SetIndex(keyOffset+slot*keyMaxLen+i, k[i])
		}
		atomics.Call("store", i32, writeIdxCell, w+1)
	}
	readIdx = func() int { return atomics.Call("load", i32, readIdxCell).Int() }
	return push, readIdx
}

func TestHostReadKeyLetsOtherGoroutinesRun(t *testing.T) {
	push, _ := installKeyRing(t)
	ticks := make(chan struct{}, 100)
	done := make(chan string)
	go func() {
		k, _ := HostKeySource{}.ReadKey()
		done <- k
	}()
	go func() {
		for i := 0; i < 5; i++ {
			time.Sleep(5 * time.Millisecond)
			ticks <- struct{}{}
		}
	}()
	for i := 0; i < 5; i++ {
		select {
		case <-ticks:
		case <-time.After(2 * time.Second):
			t.Fatal("a goroutine stalled while ReadKey waited")
		}
	}
	push("x")
	if k := <-done; k != "x" {
		t.Fatalf("ReadKey = %q, want %q", k, "x")
	}
}

// Two readers waiting on an empty ring must each get a different key. Each
// used to hold the readIdx it saw on entry, so both consumed the first key.
func TestHostReadKeyConcurrentReadersGetDistinctKeys(t *testing.T) {
	push, readIdx := installKeyRing(t)
	got := make(chan string, 2)
	for i := 0; i < 2; i++ {
		go func() {
			k, _ := HostKeySource{}.ReadKey()
			got <- k
		}()
	}
	time.Sleep(20 * time.Millisecond)
	push("x")
	time.Sleep(30 * time.Millisecond)
	push("y")
	keys := []string{<-got, <-got}
	sort.Strings(keys)
	if keys[0] != "x" || keys[1] != "y" {
		t.Fatalf("readers got %q, want [x y]", keys)
	}
	if r := readIdx(); r != 2 {
		t.Fatalf("readIdx = %d, want 2", r)
	}
}
