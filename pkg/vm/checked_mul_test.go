/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"math/big"
	"testing"
)

// TestCheckedMulInt checks checkedMulInt against a math/big oracle over a
// grid of adversarial operands: the fast-path guard boundary (mulGuardMin/
// Max and neighbors), the int extremes, and the historical minInt×-1 edge.
func TestCheckedMulInt(t *testing.T) {
	vals := []Int{
		0, 1, -1, 2, -2, 3, 255, 256, -255, -256,
		mulGuardMax, mulGuardMax - 1, mulGuardMax + 1,
		mulGuardMin, mulGuardMin + 1, mulGuardMin - 1,
		maxIntValue, maxIntValue - 1, minIntValue, minIntValue + 1,
		maxIntValue / 2, minIntValue / 2, mulGuardMax * 2, mulGuardMin * 2,
	}
	for _, a := range vals {
		for _, b := range vals {
			got, ok := checkedMulInt(a, b)
			ref := new(big.Int).Mul(big.NewInt(int64(a)), big.NewInt(int64(b)))
			fits := ref.IsInt64() && Int(ref.Int64()) >= minIntValue && Int(ref.Int64()) <= maxIntValue
			if fits != ok {
				t.Fatalf("checkedMulInt(%d, %d): ok=%v, oracle fits=%v", a, b, ok, fits)
			}
			if ok && int64(got) != ref.Int64() {
				t.Fatalf("checkedMulInt(%d, %d) = %d, oracle %s", a, b, got, ref)
			}
		}
	}
}

// BenchmarkCheckedMulInt measures the common-operand case (values well
// inside the fast-path guard) and the wide case that still takes the
// division-checked path.
func BenchmarkCheckedMulInt(b *testing.B) {
	b.Run("small", func(b *testing.B) {
		x, y := Int(31337), Int(-2711)
		for i := 0; i < b.N; i++ {
			r, ok := checkedMulInt(x, y)
			if !ok || r != x*y {
				b.Fatal("bad result")
			}
		}
	})
	b.Run("wide", func(b *testing.B) {
		// Outside the fast-path guard at any int width.
		x, y := (mulGuardMax+1)*8, Int(3)
		for i := 0; i < b.N; i++ {
			r, ok := checkedMulInt(x, y)
			if !ok || r != x*y {
				b.Fatal("bad result")
			}
		}
	})
}
