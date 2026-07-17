/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"sync"
	"testing"
)

// Fresh types seen for the first time inside these tests, so the concurrent
// boxing below exercises the BoxedTypes insert path rather than a warm cache.
type boxRaceMapType struct{ n int }
type boxRaceHashType struct{ n int }

// TestBoxedTypesConcurrentBoxing pins down the BoxedTypes map race: many
// goroutines boxing the same not-yet-cached type at once hit lookup+insert
// concurrently (a plain-map data race, and a "concurrent map writes" panic).
// All callers must converge on the single shared *aBoxedType for the type.
func TestBoxedTypesConcurrentBoxing(t *testing.T) {
	const goroutines = 32
	var wg sync.WaitGroup
	start := make(chan struct{})
	results := make([]*Boxed, goroutines)
	for i := range goroutines {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-start
			results[i] = NewBoxed(boxRaceMapType{i})
		}(i)
	}
	close(start)
	wg.Wait()

	typ := results[0].typ
	for i, b := range results {
		if b.typ != typ {
			t.Fatalf("goroutine %d resolved a different *aBoxedType (%p vs %p)", i, b.typ, typ)
		}
	}
}

// TestBoxedConcurrentHash pins down the Hash() lazy-cache race: many goroutines
// racing to first-compute the cached hash of one shared Boxed. The expected
// value is derived independently so the cache is genuinely cold at spawn time.
func TestBoxedConcurrentHash(t *testing.T) {
	const goroutines = 32
	shared := NewBoxed(boxRaceHashType{7})
	want := hashString(shared.String()) // independent of the lazy cache

	var wg sync.WaitGroup
	start := make(chan struct{})
	got := make([]uint32, goroutines)
	for i := range goroutines {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-start
			got[i] = shared.Hash()
		}(i)
	}
	close(start)
	wg.Wait()

	for i, h := range got {
		if h != want {
			t.Fatalf("goroutine %d hash = %d, want %d", i, h, want)
		}
	}
}
