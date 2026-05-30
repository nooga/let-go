/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

// scopeByGID maps a goroutine id to the Scope that goroutine is currently
// running under. Absence means "unscoped" → the root Goroutines scope.
var scopeByGID sync.Map // int64 -> *Scope

// scopedLive counts the net number of live gid registrations. It is the
// hot-path guard: when zero, CurrentScope/CurrentContext skip the goID()
// parse entirely and return root. It self-clears to zero when all scoped
// work drains.
var scopedLive atomic.Int64

// goID parses the running goroutine's id from the first line of its stack:
//
//	"goroutine 123 [running]:\n..."  ->  123
//
// Returns 0 on any parse failure (callers treat 0 as "unscoped" → root).
func goID() int64 {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	s := buf[:n]
	const prefix = "goroutine "
	if len(s) < len(prefix) {
		return 0
	}
	s = s[len(prefix):]
	i := 0
	for i < len(s) && s[i] >= '0' && s[i] <= '9' {
		i++
	}
	id, err := strconv.ParseInt(string(s[:i]), 10, 64)
	if err != nil {
		return 0
	}
	return id
}

// CurrentScope returns the scope the calling goroutine runs under, or the
// root Goroutines scope if it is unscoped. Fast path: no goID() parse while
// no scopes are in play.
func CurrentScope() *Scope {
	if scopedLive.Load() == 0 {
		return Goroutines
	}
	if v, ok := scopeByGID.Load(goID()); ok {
		return v.(*Scope)
	}
	return Goroutines
}

// CurrentContext is the cancellation context blocking native ops select on.
func CurrentContext() context.Context { return CurrentScope().Context() }

// bind makes s the calling goroutine's current scope and returns a restore
// func that reinstates the previous mapping (or removes it). Used by
// with-scope for a synchronous dynamic extent.
func (s *Scope) bind() func() {
	gid := goID()
	prev, had := scopeByGID.Load(gid)
	scopeByGID.Store(gid, s)
	if !had {
		scopedLive.Add(1)
	}
	return func() {
		if had {
			scopeByGID.Store(gid, prev)
		} else {
			scopeByGID.Delete(gid)
			scopedLive.Add(-1)
		}
	}
}

// OpenChild creates a child of the current scope, makes it current on the
// calling goroutine, and returns it. Pair with CloseScoped. Used by the
// (scope-open) native behind the with-scope macro.
func OpenChild() *Scope {
	c := CurrentScope().Child()
	c.closeRestore = c.bind()
	return c
}

// CloseScoped tears a with-scope child down: cancel the subtree, drain up to
// timeout (warn on a straggler), restore the caller's previous current scope,
// and unparent the child so it is not retained.
func CloseScoped(c *Scope, timeout time.Duration) {
	c.Cancel()
	if !c.Await(timeout) {
		fmt.Fprintf(os.Stderr,
			"warning: with-scope drain timed out; %d goroutine(s) still live\n",
			c.LiveTree())
	}
	if c.closeRestore != nil {
		c.closeRestore()
		c.closeRestore = nil
	}
	if c.parent != nil {
		c.parent.removeChild(c)
	}
}
