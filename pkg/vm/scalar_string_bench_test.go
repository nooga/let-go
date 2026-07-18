/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package vm

import "testing"

// String() on the scalar types sits under str/pr-str and any render path
// that stringifies values.
func BenchmarkScalarString(b *testing.B) {
	b.Run("Int", func(b *testing.B) {
		v := Int(1048576)
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			if len(v.String()) == 0 {
				b.Fatal("empty")
			}
		}
	})
	b.Run("Keyword", func(b *testing.B) {
		v := Keyword("player-position")
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			if len(v.String()) == 0 {
				b.Fatal("empty")
			}
		}
	})
}
