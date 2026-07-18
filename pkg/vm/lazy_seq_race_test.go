/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
)

// Concurrent accessors racing the first realization must all observe the
// same realized seq, and the thunk must run exactly once. Exercises the
// lock-free done fast path against the mutex-guarded realization transition
// (run with -race).
func TestLazySeqConcurrentRealization(t *testing.T) {
	const goroutines = 16
	var thunkRuns atomic.Int32
	ls := NewLazySeq(thunkOf(func() Value {
		thunkRuns.Add(1)
		return NewArrayVector([]Value{Int(1), Int(2), Int(3)})
	}))

	var wg sync.WaitGroup
	firsts := make([]Value, goroutines)
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			firsts[i] = ls.First()
			// Walk the whole chain to hit the realized fast path repeatedly.
			n := 0
			for s := ls.Resolve(); s != nil; s = s.Next() {
				n++
			}
			if n != 3 {
				t.Errorf("goroutine %d: walked %d elements, want 3", i, n)
			}
		}(i)
	}
	wg.Wait()

	if got := thunkRuns.Load(); got != 1 {
		t.Fatalf("thunk ran %d times, want exactly 1", got)
	}
	for i, f := range firsts {
		if f != Int(1) {
			t.Fatalf("goroutine %d saw First()=%v, want 1", i, f)
		}
	}
	if !ls.IsRealized() {
		t.Fatal("IsRealized false after realization")
	}
}

// An empty realization (thunk yields an empty collection) must canonicalize
// to nil consistently under concurrent access.
func TestLazySeqConcurrentEmptyRealization(t *testing.T) {
	ls := NewLazySeq(thunkOf(func() Value { return NewArrayVector(nil) }))
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if s := ls.Resolve(); s != nil {
				t.Errorf("empty lazy seq resolved to non-nil %v", s)
			}
			if f := ls.First(); f != NIL {
				t.Errorf("First() on empty = %v, want NIL", f)
			}
		}()
	}
	wg.Wait()
}

// A thunk error must be cached and re-raised (as thrownPanic) on every
// access, including accesses after the failing one — the error path must
// never be bypassed by the realized fast path.
func TestLazySeqConcurrentErrorRealization(t *testing.T) {
	wantErr := errors.New("boom")
	failing, _ := NativeFnType.Wrap(func(_ []Value) (Value, error) {
		return NIL, wantErr
	})
	ls := NewLazySeq(failing.(Fn))

	catchOne := func() error {
		defer func() { recover() }()
		var caught error
		func() {
			defer func() {
				if r := recover(); r != nil {
					if tp, ok := r.(*thrownPanic); ok {
						caught = tp.err
					} else {
						t.Errorf("unexpected panic value %v", r)
					}
				}
			}()
			ls.First()
		}()
		return caught
	}

	var wg sync.WaitGroup
	errs := make([]error, 8)
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			errs[i] = catchOne()
		}(i)
	}
	wg.Wait()
	for i, err := range errs {
		if !errors.Is(err, wantErr) {
			t.Fatalf("goroutine %d caught %v, want %v", i, err, wantErr)
		}
	}
	// Repeated sequential access must still re-raise, not fast-path past it.
	if err := catchOne(); !errors.Is(err, wantErr) {
		t.Fatalf("post-race access caught %v, want %v", err, wantErr)
	}
}

// The other error branch: a thunk that realizes to a non-seq, non-Sequable
// value sets err inside seq() (not sval()) — that path must also cache and
// re-raise on every access and never be bypassed by the realized fast path.
func TestLazySeqConcurrentNonSeqRealization(t *testing.T) {
	ls := NewLazySeq(thunkOf(func() Value { return Int(42) }))

	catchOne := func() error {
		var caught error
		func() {
			defer func() {
				if r := recover(); r != nil {
					if tp, ok := r.(*thrownPanic); ok {
						caught = tp.err
					} else {
						t.Errorf("unexpected panic value %v", r)
					}
				}
			}()
			ls.First()
		}()
		return caught
	}

	var wg sync.WaitGroup
	errs := make([]error, 8)
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			errs[i] = catchOne()
		}(i)
	}
	wg.Wait()
	for i, err := range errs {
		if err == nil {
			t.Fatalf("goroutine %d: non-seq realization did not raise", i)
		}
	}
	// Repeated sequential access must still re-raise, not fast-path past it.
	if err := catchOne(); err == nil {
		t.Fatal("post-race access did not re-raise")
	}
}
