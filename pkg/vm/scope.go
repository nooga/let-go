/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// Scope is a node in the VM's goroutine supervision tree — a structured
// home for every goroutine the runtime spawns (futures, agents,
// core.async go-blocks, channel pipelines). It exists because Go has no
// goroutine registry of its own and cannot force-kill a goroutine: the
// runtime needs a way to observe, await, and cooperatively cancel the
// work it started.
//
// Each scope owns:
//   - a cancellable context derived from its PARENT's context, so
//     cancelling a scope cancels its whole subtree (children inherit the
//     cancel) while sibling scopes are unaffected;
//   - its OWN wait-group + live counter, so spawning into a scope touches
//     only that scope's bookkeeping — there is no single global counter
//     every goroutine contends on. The hot path (Go/Context) is
//     lock-free; the only mutex is the cold child-list, touched at scope
//     creation, not per spawn.
//
// This is the Go-native cousin of an Erlang supervision tree / a Trio
// nursery / Loom StructuredTaskScope. The ceiling is the same as Erlang's
// is not: Go can't reap a goroutine that ignores its context, so
// cancellation is cooperative (blocking ops select on the scope context).
//
// `Goroutines` is the process-wide ROOT scope; everything currently spawns
// there. Sub-scopes (Child) are opt-in and let a caller cancel/await one
// subtree independently — the per-goroutine cancellation the flat registry
// could not do.
type ctxState struct {
	ctx    context.Context
	cancel context.CancelFunc
}

type Scope struct {
	parent *Scope
	state  atomic.Pointer[ctxState] // own context generation (child of parent)
	wg     sync.WaitGroup           // this scope's direct goroutines
	live   atomic.Int64

	mu           sync.Mutex // guards children (cold path: scope creation only)
	children     map[*Scope]struct{}
	closeRestore func() // set by OpenChild, consumed by CloseScoped
}

// Goroutines is the process-wide root scope. Every VM spawn goes through
// a scope; blocking ops consult its Context for cancellation.
var Goroutines = newRootScope()

func newRootScope() *Scope {
	s := &Scope{}
	ctx, cancel := context.WithCancel(context.Background())
	s.state.Store(&ctxState{ctx: ctx, cancel: cancel})
	return s
}

// Child creates a sub-scope whose context is derived from this scope's
// current context. Cancelling this scope (or any ancestor) cancels the
// child; cancelling the child leaves this scope and its siblings running.
func (s *Scope) Child() *Scope {
	parentCtx := s.state.Load().ctx
	ctx, cancel := context.WithCancel(parentCtx)
	c := &Scope{parent: s}
	c.state.Store(&ctxState{ctx: ctx, cancel: cancel})
	s.mu.Lock()
	if s.children == nil {
		s.children = make(map[*Scope]struct{})
	}
	s.children[c] = struct{}{}
	s.mu.Unlock()
	return c
}

// removeChild drops c from this scope's child set (cold path: scope close).
func (s *Scope) removeChild(c *Scope) {
	s.mu.Lock()
	delete(s.children, c)
	s.mu.Unlock()
}

// Go runs fn in a tracked goroutine of THIS scope. fn receives the
// scope's current context; long-running work should select on ctx.Done()
// so Cancel/Shutdown/Drain can stop it.
func (s *Scope) Go(fn func(ctx context.Context)) {
	ctx := s.state.Load().ctx
	s.wg.Add(1)
	s.live.Add(1)
	// Only goroutines running under a NON-root scope need a gid registration:
	// a goroutine under the root resolves to Goroutines via the fast-path
	// fallback anyway, so registering it would pointlessly bump scopedLive and
	// defeat the scopedLive==0 guard for every other goroutine's channel ops
	// (futures, go-blocks, sleeps and the async pumps all spawn under root).
	register := s != Goroutines
	go func() {
		var gid int64
		if register {
			gid = goID()
			scopeByGID.Store(gid, s)
			scopedLive.Add(1)
		}
		defer func() {
			if register {
				scopeByGID.Delete(gid)
				scopedLive.Add(-1)
			}
			s.live.Add(-1)
			s.wg.Done()
		}()
		fn(ctx)
	}()
}

// Live reports how many goroutines are running directly in this scope
// (not counting sub-scopes). Use LiveTree for the whole subtree.
func (s *Scope) Live() int { return int(s.live.Load()) }

// LiveTree reports the live goroutine count of this scope and all
// descendants.
func (s *Scope) LiveTree() int {
	n := int(s.live.Load())
	s.mu.Lock()
	kids := make([]*Scope, 0, len(s.children))
	for c := range s.children {
		kids = append(kids, c)
	}
	s.mu.Unlock()
	for _, c := range kids {
		n += c.LiveTree()
	}
	return n
}

// Context returns the scope's current cancellation context. Blocking
// native ops (sleep, channel send/recv, alts!) select on its Done() so
// they unblock when the scope is cancelled. Lock-free (one atomic load),
// cheap to call on every channel op.
func (s *Scope) Context() context.Context {
	return s.state.Load().ctx
}

// Cancel cancels this scope's context — and, by context inheritance, its
// whole subtree — without installing a fresh one. Terminal: the scope
// stops accepting new cancellable work. Siblings are unaffected. Does not
// wait; pair with Await (or use Shutdown).
func (s *Scope) Cancel() {
	s.state.Load().cancel()
}

// CancelAll cancels the current context generation and installs a fresh
// one derived from the parent, so the scope can keep accepting spawns
// after a drain. This is the "drain and continue" form used by the root
// between benchmark iterations / at soft shutdown points. Lock-free: the
// fresh generation is swapped in atomically and we cancel whatever we
// displaced, so concurrent CancelAlls are safe.
func (s *Scope) CancelAll() {
	parentCtx := context.Background()
	if s.parent != nil {
		parentCtx = s.parent.state.Load().ctx
	}
	ctx, cancel := context.WithCancel(parentCtx)
	old := s.state.Swap(&ctxState{ctx: ctx, cancel: cancel})
	old.cancel()
}

// Await blocks until every goroutine in this scope AND its subtree has
// exited, or timeout elapses. Returns true if fully drained, false on
// timeout. A non-positive timeout waits indefinitely.
func (s *Scope) Await(timeout time.Duration) bool {
	done := make(chan struct{})
	go func() {
		s.awaitTree()
		close(done)
	}()
	if timeout <= 0 {
		<-done
		return true
	}
	t := time.NewTimer(timeout)
	defer t.Stop()
	select {
	case <-done:
		return true
	case <-t.C:
		return false
	}
}

func (s *Scope) awaitTree() {
	s.wg.Wait()
	s.mu.Lock()
	kids := make([]*Scope, 0, len(s.children))
	for c := range s.children {
		kids = append(kids, c)
	}
	s.mu.Unlock()
	for _, c := range kids {
		c.awaitTree()
	}
}

// Shutdown cancels this scope (terminal, cascades to the subtree) and
// waits up to timeout for it to drain. The structured-shutdown form for
// a sub-scope; for the root's drain-and-continue use Drain.
func (s *Scope) Shutdown(timeout time.Duration) bool {
	s.Cancel()
	return s.Await(timeout)
}

// Drain cancels the scope's current generation and reinstalls a fresh one
// (so the scope keeps working), then waits for the displaced work to
// exit. Used by the root between bench iterations and at soft shutdown.
func (s *Scope) Drain(timeout time.Duration) bool {
	s.CancelAll()
	return s.Await(timeout)
}

// Scope is a Value so it can be handed to Lisp as an opaque handle.
func (s *Scope) Type() ValueType { return ScopeType }
func (s *Scope) Unbox() any      { return s }
func (s *Scope) String() string  { return fmt.Sprintf("#<scope live=%d>", s.LiveTree()) }

type theScopeType struct{}

func (t *theScopeType) String() string  { return t.Name() }
func (t *theScopeType) Type() ValueType { return TypeType }
func (t *theScopeType) Unbox() any      { return nil }
func (t *theScopeType) Name() string    { return "let-go.lang.Scope" }
func (t *theScopeType) Box(bare any) (Value, error) {
	return NIL, NewTypeError(bare, "can't be boxed as", t)
}

// ScopeType is the Value type of a Scope handle.
var ScopeType *theScopeType = &theScopeType{}
