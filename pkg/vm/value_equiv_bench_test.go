/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package vm

import "testing"

// Direct benchmarks for valueEquiv's dominant map-key type pairs. The
// map-level benchmarks (MapAssoc/*) exercise these through the trie; these
// isolate the comparison itself.
func BenchmarkValueEquiv(b *testing.B) {
	kw1, kw2 := Keyword("player-position"), Keyword("player-position")
	kwOther := Keyword("world-seed")
	i1, i2, iOther := Int(1048576), Int(1048576), Int(99)
	f := Float(1048576)

	b.Run("Keyword-equal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if !valueEquiv(kw1, kw2) {
				b.Fatal("expected equal")
			}
		}
	})
	b.Run("Keyword-unequal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if valueEquiv(kw1, kwOther) {
				b.Fatal("expected unequal")
			}
		}
	})
	b.Run("Int-equal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if !valueEquiv(i1, i2) {
				b.Fatal("expected equal")
			}
		}
	})
	b.Run("Int-unequal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if valueEquiv(i1, iOther) {
				b.Fatal("expected unequal")
			}
		}
	})
	// Mixed numeric must keep falling through to the NumEq path, which
	// follows Clojure = semantics: Int and Float are never map-key
	// equivalent (verified identical on the pre-change code). The Int fast
	// path above must not short-circuit this pair.
	b.Run("Int-vs-Float", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if valueEquiv(i1, f) {
				b.Fatal("Int and Float must not be map-key equivalent")
			}
		}
	})
}
