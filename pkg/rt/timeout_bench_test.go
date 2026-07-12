/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"context"
	"testing"
	"time"

	"github.com/nooga/let-go/pkg/vm"
)

// These benchmarks quantify what the shared timer daemon buys per (timeout ms)
// call: a channel plus a queue entry, versus the old shape's goroutine spawn,
// goroutine stack, timer allocation, and closure per call. The old shape is
// preserved below as the baseline so the pair stays comparable in one run:
//
//	go test -bench Timeout -benchmem ./pkg/rt/
//
// A 1ms deadline keeps the daemon's queue short while the loop runs (entries
// drain about as fast as they are created), so the daemon numbers include its
// steady-state rescan cost rather than an ever-growing queue.

// timeoutSink defeats dead-code elimination of the returned channel.
var timeoutSink vm.Chan

// oldTimeoutChan is the pre-daemon implementation: one goroutine per call.
func oldTimeoutChan(ms int) vm.Chan {
	ch := make(vm.Chan)
	vm.Goroutines.Go(func(ctx context.Context) {
		t := time.NewTimer(time.Duration(ms) * time.Millisecond)
		defer t.Stop()
		select {
		case <-t.C:
		case <-ctx.Done():
		}
		close(ch)
	})
	return ch
}

func benchTimeouts(b *testing.B, mk func(int) vm.Chan) {
	b.ReportAllocs()
	var last vm.Chan
	for i := 0; i < b.N; i++ {
		last = mk(1)
	}
	timeoutSink = last
	b.StopTimer()
	<-last // drain: don't leak this run's pending work into the next benchmark
}

func BenchmarkTimeoutDaemon(b *testing.B) {
	benchTimeouts(b, timeoutChan)
}

func BenchmarkTimeoutGoroutinePerCall(b *testing.B) {
	benchTimeouts(b, oldTimeoutChan)
}
