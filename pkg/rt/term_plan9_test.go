//go:build plan9

/*
 * Copyright (c) 2026 let-go contributors; see CONTRIBUTORS.
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"sync/atomic"
	"testing"
	"time"
)

// TestPlan9KeySourceCatchesResizeBeforeFirstRead covers the startup window
// where a terminal publishes corrected geometry after the root key source is
// constructed but before an application reaches its first blocking read-key.
// The constructor's generation baseline must make that change observable.
func TestPlan9KeySourceCatchesResizeBeforeFirstRead(t *testing.T) {
	var generation atomic.Value
	generation.Store("1")
	readWinch := func() (string, bool) {
		return generation.Load().(string), true
	}

	ch := make(chan []byte)
	ks := newPlan9KeySourceWithWinch(
		&chunkReader{ch: ch},
		readWinch,
		time.Millisecond,
	)

	// The correction lands after construction but before ReadKey starts.
	generation.Store("2")
	result := make(chan string, 1)
	go func() {
		k, _ := ks.ReadKey()
		result <- k
	}()

	select {
	case got := <-result:
		if got != terminalWakeKey {
			t.Fatalf("startup resize key = %q, want %q", got, terminalWakeKey)
		}
	case <-time.After(time.Second):
		t.Fatal("resize before first ReadKey was missed")
	}
	close(ch)
}
