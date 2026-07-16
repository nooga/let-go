/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package test

import (
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

func TestBigdecPreservesBigIntPrecision(t *testing.T) {
	got, err := evalCoreCompat(`(= 123456789012345678901234567890.0M
	                              (bigdec 123456789012345678901234567890N))`)
	if err != nil {
		t.Fatal(err)
	}
	if got != vm.TRUE {
		t.Fatalf("bigdec lost BigInt precision: got %s", got)
	}
}

func TestBigdecRejectsNonFiniteFloats(t *testing.T) {
	for _, input := range []string{"##Inf", "##-Inf", "##NaN"} {
		t.Run(input, func(t *testing.T) {
			if _, err := evalCoreCompat(`(bigdec ` + input + `)`); err == nil {
				t.Fatalf("bigdec accepted %s", input)
			}
		})
	}
}

func TestBigDecimalSignedZeroEqualityAfterMultiplication(t *testing.T) {
	got, err := evalCoreCompat(`[(= 0.0M (*' -1 0.0M))
	                             (= 0.0M (*' 0.0M -1))
	                             (= 0.0M (bigdec -0.0))
	                             (pr-str (*' -1 0.0M))]`)
	if err != nil {
		t.Fatal(err)
	}
	if got.String() != `[true true true "0.0M"]` {
		t.Fatalf("signed BigDecimal zero must compare equal: got %s", got)
	}
}
