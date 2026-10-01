/*
 * Copyright (c) 2026 Matt Parrett
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"math"
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

// Pods decode JSON responses through JSONDecodeValue; an id echoed back to a
// peer has to survive exactly.
func TestJSONDecodeValueExactID(t *testing.T) {
	v, err := JSONDecodeValue(`{"id": 9007199254740993, "id": 9223372036854775807}`)
	if err != nil {
		t.Fatal(err)
	}
	got := v.(vm.Lookup).ValueAt(vm.Keyword("id"))
	if got != vm.Int(math.MaxInt64) {
		t.Fatalf("id = %v, want the last duplicate %d", got, int64(math.MaxInt64))
	}
}

// float64(int64(f)) is platform-defined outside int64 range, so the whole-float
// conversion has to bound f itself rather than round-trip it.
func TestFloatToValueInt64Bounds(t *testing.T) {
	cases := []struct {
		f    float64
		want vm.Value
	}{
		{3.0, vm.Int(3)},
		{-(1 << 63), vm.Int(math.MinInt64)},
		{1 << 63, vm.Float(1 << 63)},
		{-(1 << 64), vm.Float(-(1 << 64))},
		{1.5, vm.Float(1.5)},
	}
	for _, c := range cases {
		if got := floatToValue(c.f); got != c.want {
			t.Errorf("floatToValue(%g) = %v (%T), want %v (%T)", c.f, got, got, c.want, c.want)
		}
	}
}

func TestDecodeJSONRejectsTrailingData(t *testing.T) {
	for _, s := range []string{"1 2", `{"a":1}x`, "[] []"} {
		if _, err := decodeJSON(s); err == nil {
			t.Errorf("decodeJSON(%q) accepted trailing data", s)
		}
	}
	if _, err := decodeJSON(" 1 \n\t"); err != nil {
		t.Errorf("surrounding whitespace rejected: %v", err)
	}
}
