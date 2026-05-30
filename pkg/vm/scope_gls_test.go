/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestGoIDIsNonZeroAndDistinct(t *testing.T) {
	a := goID()
	if a == 0 {
		t.Fatal("goID returned 0 for the running goroutine")
	}
	var b int64
	done := make(chan struct{})
	go func() { b = goID(); close(done) }()
	<-done
	if b == 0 || a == b {
		t.Fatalf("expected distinct non-zero goids, got a=%d b=%d", a, b)
	}
}

func TestCurrentScopeDefaultsToRootWhenUnscoped(t *testing.T) {
	if scopedLive.Load() != 0 {
		t.Skipf("scopedLive not clean: %d", scopedLive.Load())
	}
	if CurrentScope() != Goroutines {
		t.Fatal("unscoped goroutine should see the root scope")
	}
	if CurrentContext() != Goroutines.Context() {
		t.Fatal("unscoped CurrentContext should be the root context")
	}
}

func TestBindRestoreNesting(t *testing.T) {
	c1 := Goroutines.Child()
	c2 := Goroutines.Child()
	r1 := c1.bind()
	if CurrentScope() != c1 {
		t.Fatal("after binding c1, current should be c1")
	}
	r2 := c2.bind()
	if CurrentScope() != c2 {
		t.Fatal("after binding c2, current should be c2")
	}
	r2()
	if CurrentScope() != c1 {
		t.Fatal("after restoring c2, current should be c1 again")
	}
	r1()
	if CurrentScope() != Goroutines {
		t.Fatal("after restoring c1, current should be root")
	}
	if scopedLive.Load() != 0 {
		t.Fatalf("scopedLive should return to 0, got %d", scopedLive.Load())
	}
}

func TestConcurrentGoIDNoRace(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 64; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); _ = goID() }()
	}
	wg.Wait()
}

func TestScopeGoRegistersItsScope(t *testing.T) {
	child := Goroutines.Child()
	defer Goroutines.removeChild(child)
	seen := make(chan *Scope, 1)
	child.Go(func(ctx context.Context) { seen <- CurrentScope() })
	got := <-seen
	if got != child {
		t.Fatalf("goroutine spawned via child.Go should see child as current, got %v", got)
	}
	child.Await(time.Second)
	if scopedLive.Load() != 0 {
		t.Fatalf("scopedLive should return to 0 after the goroutine exits, got %d", scopedLive.Load())
	}
}

func TestSubScopeCancelIsIndependent(t *testing.T) {
	a := Goroutines.Child()
	b := Goroutines.Child()
	defer Goroutines.removeChild(a)
	defer Goroutines.removeChild(b)
	ch := make(chan Value) // never fed → recv parks
	park := func(ctx context.Context) {
		select {
		case <-ch:
		case <-CurrentContext().Done(): // gid lookup resolves to a or b
		}
	}
	a.Go(park)
	b.Go(park)
	waitLiveScoped(t, a, 1)
	waitLiveScoped(t, b, 1)

	a.Cancel() // cancel ONLY a's subtree
	if !a.Await(2 * time.Second) {
		t.Fatal("a's parked goroutine should unblock on a.Cancel")
	}
	if b.Live() != 1 {
		t.Fatalf("b must remain live after cancelling a, got %d", b.Live())
	}
	b.Cancel()
	b.Await(2 * time.Second)
}

func waitLiveScoped(t *testing.T, s *Scope, want int) {
	t.Helper()
	deadline := time.Now().Add(time.Second)
	for s.Live() != want && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}
	if s.Live() != want {
		t.Fatalf("expected %d live, got %d", want, s.Live())
	}
}

func TestCloseScopedUnparentsNoGrowth(t *testing.T) {
	before := childCount(Goroutines)
	for i := 0; i < 20; i++ {
		c := OpenChild()
		CloseScoped(c, time.Second)
	}
	if after := childCount(Goroutines); after != before {
		t.Fatalf("root children grew: before=%d after=%d", before, after)
	}
}

func childCount(s *Scope) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.children)
}

func TestScopeIsAValue(t *testing.T) {
	var v Value = Goroutines.Child()
	if v.Type() != ScopeType {
		t.Fatalf("scope Type should be ScopeType, got %v", v.Type())
	}
	if v.String() == "" {
		t.Fatal("scope String should be non-empty")
	}
}
