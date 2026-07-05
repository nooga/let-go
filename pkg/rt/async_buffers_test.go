/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package rt

// Parity tests for the core.async buffer policies and put/close semantics,
// written against clojure.core.async's documented behavior:
//   dropping-buffer — "When full, puts will complete but val will be
//                      dropped (no transfer)."
//   sliding-buffer  — "When full, puts will complete and be buffered, but
//                      the oldest elements in the buffer will be dropped."
//   >!!             — "Returns true unless the channel is already closed."
//   offer!          — true if the put succeeds immediately, else nil.
//   close!          — no-op on an already-closed channel.

import (
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

func mkChanWith(t *testing.T, bufFn string, n int) vm.Chan {
	t.Helper()
	buf := invoke(t, asyncFn(t, bufFn), vm.Int(int64(n)))
	ch, ok := invoke(t, asyncFn(t, "chan"), buf).(vm.Chan)
	if !ok {
		t.Fatalf("chan with %s did not return a Chan", bufFn)
	}
	return ch
}

func TestDroppingBufferDropsNewestWhenFull(t *testing.T) {
	put := asyncFn(t, ">!")
	poll := asyncFn(t, "poll!")
	ch := mkChanWith(t, "dropping-buffer", 2)

	for i := 1; i <= 3; i++ {
		if got := invoke(t, put, ch, vm.Int(int64(i))); got != vm.TRUE {
			t.Fatalf("put %d on dropping chan: expected true (put completes), got %v", i, got)
		}
	}
	// Buffer holds the FIRST two; the third (newest) was dropped.
	if got := invoke(t, poll, ch); got != vm.Int(1) {
		t.Fatalf("expected 1, got %v", got)
	}
	if got := invoke(t, poll, ch); got != vm.Int(2) {
		t.Fatalf("expected 2, got %v", got)
	}
	if got := invoke(t, poll, ch); got != vm.NIL {
		t.Fatalf("expected empty (nil), got %v", got)
	}
}

func TestSlidingBufferDropsOldestWhenFull(t *testing.T) {
	put := asyncFn(t, ">!")
	poll := asyncFn(t, "poll!")
	ch := mkChanWith(t, "sliding-buffer", 2)

	for i := 1; i <= 3; i++ {
		if got := invoke(t, put, ch, vm.Int(int64(i))); got != vm.TRUE {
			t.Fatalf("put %d on sliding chan: expected true (put completes), got %v", i, got)
		}
	}
	// Buffer holds the LAST two; the first (oldest) was evicted.
	if got := invoke(t, poll, ch); got != vm.Int(2) {
		t.Fatalf("expected 2, got %v", got)
	}
	if got := invoke(t, poll, ch); got != vm.Int(3) {
		t.Fatalf("expected 3, got %v", got)
	}
	if got := invoke(t, poll, ch); got != vm.NIL {
		t.Fatalf("expected empty (nil), got %v", got)
	}
}

func TestFixedBufferMarkerMatchesPlainSizedChan(t *testing.T) {
	offer := asyncFn(t, "offer!")
	ch := mkChanWith(t, "buffer", 1)
	if got := invoke(t, offer, ch, vm.Int(1)); got != vm.TRUE {
		t.Fatalf("first offer on (chan (buffer 1)): expected true, got %v", got)
	}
	// Full fixed buffer: offer! cannot complete → nil (not false)
	if got := invoke(t, offer, ch, vm.Int(2)); got != vm.NIL {
		t.Fatalf("offer on full fixed chan: expected nil, got %v", got)
	}
}

func TestPutOnClosedChannelReturnsFalse(t *testing.T) {
	put := asyncFn(t, ">!")
	closef := asyncFn(t, "close!")
	ch := invoke(t, asyncFn(t, "chan"), vm.Int(1)).(vm.Chan)
	invoke(t, closef, ch)
	// close! twice: no-op, must not panic
	invoke(t, closef, ch)
	if got := invoke(t, put, ch, vm.Int(1)); got != vm.FALSE {
		t.Fatalf(">! on closed chan: expected false, got %v", got)
	}
}

func TestOfferOnClosedChannelReturnsNil(t *testing.T) {
	offer := asyncFn(t, "offer!")
	closef := asyncFn(t, "close!")
	ch := invoke(t, asyncFn(t, "chan"), vm.Int(1)).(vm.Chan)
	invoke(t, closef, ch)
	if got := invoke(t, offer, ch, vm.Int(1)); got != vm.NIL {
		t.Fatalf("offer! on closed chan: expected nil, got %v", got)
	}
}

func TestDroppingSlidingPutsNeverPark(t *testing.T) {
	// A put on a full dropping/sliding channel must return immediately —
	// run many puts synchronously; if any parked this test would hang and
	// the suite's timeout catches it.
	put := asyncFn(t, ">!")
	for _, bufFn := range []string{"dropping-buffer", "sliding-buffer"} {
		ch := mkChanWith(t, bufFn, 1)
		for i := 0; i < 100; i++ {
			if got := invoke(t, put, ch, vm.Int(int64(i))); got != vm.TRUE {
				t.Fatalf("%s put %d: expected true, got %v", bufFn, i, got)
			}
		}
	}
}
