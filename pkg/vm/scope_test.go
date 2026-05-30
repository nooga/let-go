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

func TestScopeTracksLiveAndAwait(t *testing.T) {
	s := newRootScope()
	release := make(chan struct{})

	for i := 0; i < 3; i++ {
		s.Go(func(ctx context.Context) { <-release })
	}
	deadline := time.Now().Add(time.Second)
	for s.Live() != 3 && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}
	if got := s.Live(); got != 3 {
		t.Fatalf("expected 3 live goroutines, got %d", got)
	}
	if s.Await(50 * time.Millisecond) {
		t.Fatal("Await should have timed out while goroutines blocked")
	}
	close(release)
	if !s.Await(2 * time.Second) {
		t.Fatal("Await should have drained after release")
	}
	if got := s.Live(); got != 0 {
		t.Fatalf("expected 0 live after drain, got %d", got)
	}
}

func TestScopeDrainReinstallsContext(t *testing.T) {
	s := newRootScope()
	for i := 0; i < 4; i++ {
		s.Go(func(ctx context.Context) {
			select {
			case <-ctx.Done():
			case <-time.After(30 * time.Second):
			}
		})
	}
	if !s.Drain(2 * time.Second) {
		t.Fatal("Drain should have cancelled+drained ctx-aware goroutines fast")
	}
	if got := s.Live(); got != 0 {
		t.Fatalf("expected 0 live after drain, got %d", got)
	}
	// Drain reinstalls a fresh context: new spawns are NOT born cancelled.
	select {
	case <-s.Context().Done():
		t.Fatal("fresh context after drain should not be cancelled")
	default:
	}
}

// blockOnCtx spawns n goroutines in scope that park until their scope
// context is cancelled (or a long fallback).
func blockOnCtx(s *Scope, n int) {
	for i := 0; i < n; i++ {
		s.Go(func(ctx context.Context) {
			select {
			case <-ctx.Done():
			case <-time.After(30 * time.Second):
			}
		})
	}
}

func waitLive(s *Scope, want int) bool {
	deadline := time.Now().Add(time.Second)
	for s.Live() != want && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}
	return s.Live() == want
}

// TestScopeSubtreeCancelLeavesSiblingsRunning is the headline of the
// prototype: cancelling one sub-scope drains ITS goroutines while a
// sibling sub-scope keeps running untouched — the per-subtree
// cancellation the flat registry could not do.
func TestScopeSubtreeCancelLeavesSiblingsRunning(t *testing.T) {
	root := newRootScope()
	a := root.Child()
	b := root.Child()

	blockOnCtx(a, 3)
	blockOnCtx(b, 2)
	if !waitLive(a, 3) || !waitLive(b, 2) {
		t.Fatalf("expected a=3 b=2 live, got a=%d b=%d", a.Live(), b.Live())
	}

	// Shut down only scope a.
	if !a.Shutdown(2 * time.Second) {
		t.Fatal("Shutdown(a) should have drained a's goroutines")
	}
	if a.Live() != 0 {
		t.Fatalf("expected a drained, got %d", a.Live())
	}
	// Sibling b is untouched.
	if b.Live() != 2 {
		t.Fatalf("sibling b should be unaffected, got %d live", b.Live())
	}

	if !b.Shutdown(2 * time.Second) {
		t.Fatal("Shutdown(b) should have drained b")
	}
	if root.LiveTree() != 0 {
		t.Fatalf("expected whole tree drained, got %d", root.LiveTree())
	}
}

// TestScopeParentCancelCascades: cancelling a parent cancels and drains
// the whole subtree (children inherit the context).
func TestScopeParentCancelCascades(t *testing.T) {
	root := newRootScope()
	parent := root.Child()
	child := parent.Child()
	grandchild := child.Child()

	blockOnCtx(parent, 1)
	blockOnCtx(child, 2)
	blockOnCtx(grandchild, 1)

	deadline := time.Now().Add(time.Second)
	for root.LiveTree() != 4 && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}
	if got := root.LiveTree(); got != 4 {
		t.Fatalf("expected 4 live in subtree, got %d", got)
	}

	// Cancel the top of the subtree; the whole subtree must drain.
	if !parent.Shutdown(2 * time.Second) {
		t.Fatal("parent.Shutdown should have drained the subtree")
	}
	if got := parent.LiveTree(); got != 0 {
		t.Fatalf("expected subtree drained, got %d", got)
	}
}

// TestScopeChildContextDerivesFromParent: a child's context is cancelled
// when the parent is, but a child cancel does not cancel the parent.
func TestScopeChildContextDerivesFromParent(t *testing.T) {
	root := newRootScope()
	a := root.Child()
	b := root.Child()

	// Cancelling a must not cancel b or root.
	a.Cancel()
	select {
	case <-a.Context().Done():
	default:
		t.Fatal("a.Context should be cancelled after a.Cancel")
	}
	select {
	case <-b.Context().Done():
		t.Fatal("b.Context must NOT be cancelled by a.Cancel")
	default:
	}
	select {
	case <-root.Context().Done():
		t.Fatal("root.Context must NOT be cancelled by a.Cancel")
	default:
	}
}
