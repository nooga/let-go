/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"context"
	"testing"
	"time"
)

// TestRootSpawnsDoNotBumpScopedLive is the fail-fast regression guard for the
// goID() hot-path cost. The scope changes can only slow the channel-op hot path
// (<!/>!/alts!/sleep all call CurrentContext()) if the scopedLive==0 fast-path
// guard gets defeated and goID() — a ~1µs runtime.Stack parse — starts firing
// on every call.
//
// The subtle way that happens: futures, go-blocks, sleeps, and the async CSP
// pumps all spawn under the ROOT scope. If a root spawn bumped scopedLive, then
// any time one was live (i.e. constantly, in real programs) EVERY other
// goroutine's channel op would pay goID() — a large, suite-wide slowdown. Root
// goroutines resolve to root via the fallback anyway and must NOT register.
//
// This test pins that invariant: with many root-scoped goroutines live,
// scopedLive stays 0 and the unscoped CurrentContext() stays in the nanosecond
// regime. It runs in `go test ./pkg/vm`, so it guards every build cheaply —
// no multi-minute bench-ratchet needed to catch this class of regression.
func TestRootSpawnsDoNotBumpScopedLive(t *testing.T) {
	if got := scopedLive.Load(); got != 0 {
		t.Skipf("scopedLive not clean (%d) — another test left scoped goroutines live", got)
	}

	release := make(chan struct{})
	const rootGoroutines = 32
	for i := 0; i < rootGoroutines; i++ {
		Goroutines.Go(func(ctx context.Context) { <-release })
	}
	// Let them all start and register (if they wrongly did).
	deadline := time.Now().Add(time.Second)
	for Goroutines.Live() < rootGoroutines && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}

	// THE INVARIANT: root spawns do not bump the hot-path guard.
	if got := scopedLive.Load(); got != 0 {
		close(release)
		t.Fatalf("REGRESSION: %d root-scoped goroutines bumped scopedLive to %d — "+
			"the scopedLive==0 fast path is defeated and goID() now fires on every "+
			"channel op suite-wide. Root spawns must not register a gid.", rootGoroutines, got)
	}

	// With scopedLive==0 proven above, CurrentContext() provably does not call
	// goID() (it short-circuits to root on the atomic load). We log a rough
	// per-op cost for visibility but do NOT assert on wall-clock — it is
	// dominated by -race atomic instrumentation and would be flaky. The
	// structural invariant (scopedLive==0) is the real, deterministic guard.
	const n = 1_000_000
	var sink context.Context
	start := time.Now()
	for i := 0; i < n; i++ {
		sink = CurrentContext()
	}
	_ = sink
	t.Logf("unscoped CurrentContext()=%v/op with %d root goroutines live (informational)",
		time.Since(start)/n, rootGoroutines)

	close(release)
	Goroutines.Await(2 * time.Second)

	// Sanity the OTHER direction: a non-root sub-scope DOES register, so
	// sub-scope cancellation still works (the whole point of the feature).
	child := Goroutines.Child()
	defer Goroutines.removeChild(child)
	seen := make(chan int64, 1)
	child.Go(func(ctx context.Context) { seen <- scopedLive.Load() })
	if got := <-seen; got < 1 {
		t.Fatalf("a sub-scope goroutine must bump scopedLive (got %d) so its scope is discoverable", got)
	}
	child.Await(2 * time.Second)
}
