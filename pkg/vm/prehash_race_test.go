/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"sync"
	"testing"
)

// A constant that lowered Go holds in a package-level var is shared by every
// goroutine from the start. Prehash fills its hash cache in the initializer,
// so concurrent hashing only reads it; run with -race.
func TestPrehashedConstantsHashWithoutARace(t *testing.T) {
	vals := []Value{
		Prehash(NewArrayMap([]Value{Keyword("a"), Int(1), Keyword("b"), Int(2)})),
		Prehash(NewSet([]Value{Int(1), Int(2), Int(3)})),
		Prehash(NewList([]Value{Int(1), Int(2)})),
	}
	want := make([]uint32, len(vals))
	for i, v := range vals {
		want[i] = HashValue(v)
	}
	var wg sync.WaitGroup
	for g := 0; g < 8; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i, v := range vals {
				if h := HashValue(v); h != want[i] {
					t.Errorf("hash of %v changed: %d, want %d", v, h, want[i])
				}
			}
		}()
	}
	wg.Wait()
}
