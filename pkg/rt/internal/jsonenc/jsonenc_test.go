/*
 * Copyright (c) 2026 Matt Parrett
 * SPDX-License-Identifier: MIT
 */

package jsonenc

import (
	"encoding/binary"
	"encoding/json"
	"math"
	"testing"
)

// assertParity checks Marshal against encoding/json, the oracle. Both must
// fail together or produce identical bytes.
func assertParity(t *testing.T, v any) {
	t.Helper()
	want, wantErr := json.Marshal(v)
	got, gotErr := Marshal(v)
	if (wantErr != nil) != (gotErr != nil) {
		t.Fatalf("error mismatch for %#v: encoding/json=%v jsonenc=%v", v, wantErr, gotErr)
	}
	if wantErr != nil {
		if wantErr.Error() != gotErr.Error() {
			t.Fatalf("error text for %#v: encoding/json=%q jsonenc=%q", v, wantErr, gotErr)
		}
		return
	}
	if string(want) != string(got) {
		t.Fatalf("output mismatch for %#v:\n encoding/json %s\n jsonenc       %s", v, want, got)
	}
}

func TestParityCases(t *testing.T) {
	cases := []any{
		nil, true, false, 0, -1, int64(math.MaxInt64), int64(math.MinInt64),
		"", "plain", "quote\" backslash\\ slash/", "\b\f\n\r\t\x00\x01\x1f\x7f",
		"<script>&amp;</script>", "\u2028\u2029", "\u00e9\u65e5\u672c\U0001F600", "\ufffd",
		"bad\xffutf8\xc3", "\xed\xa0\x80", // lone surrogate encoded as UTF-8 is invalid
		0.0, math.Copysign(0, -1), 1.0, 2.5, 0.1, 62.302, 1e-6, 1e-7, 1.5e-7,
		9.99999e-7, 1e20, 9.999999999999999e20, 1e21, -1e21, 1e100, 1e-100,
		5e-324, math.MaxFloat64, 12345678901234567890.0, 1.0 / 3,
		math.NaN(), math.Inf(1), math.Inf(-1),
		[]any{}, []any(nil), map[string]any{}, map[string]any(nil),
		[]any{1, "a", nil, true, 2.5, []any{map[string]any{"k": nil}}},
		map[string]any{"b": 1, "a": 2, "": 3, "B": 4, "\u00e9": 5, "<": 6, "a\x00": 7},
		map[string]any{"xsofy/stats": map[string]any{"hp": 12, "max-hp": 20, "depth": 3, "turn": 1234}},
		map[string]any{"bad": math.NaN()},
	}
	for _, c := range cases {
		assertParity(t, c)
	}
}

func TestUnsupportedType(t *testing.T) {
	if _, err := Marshal(struct{}{}); err == nil {
		t.Fatal("expected an error for a type outside the closed set")
	}
}

// valueFrom decodes fuzz bytes into a nested value from the closed set, so the
// fuzzer explores structure as well as leaf content.
func valueFrom(b []byte, depth int) (any, []byte) {
	if len(b) == 0 {
		return nil, b
	}
	tag, b := b[0], b[1:]
	take := func(n int) []byte {
		if n > len(b) {
			n = len(b)
		}
		out := b[:n]
		b = b[n:]
		return out
	}
	switch tag % 8 {
	case 0:
		return nil, b
	case 1:
		return tag&0x80 != 0, b
	case 2:
		n := 0
		if len(b) > 0 {
			n = int(b[0] % 32)
			b = b[1:]
		}
		return string(take(n)), b
	case 3:
		var buf [8]byte
		copy(buf[:], take(8))
		return int(int64(binary.LittleEndian.Uint64(buf[:]))), b
	case 4:
		var buf [8]byte
		copy(buf[:], take(8))
		return int64(binary.LittleEndian.Uint64(buf[:])), b
	case 5:
		var buf [8]byte
		copy(buf[:], take(8))
		return math.Float64frombits(binary.LittleEndian.Uint64(buf[:])), b
	case 6:
		if depth > 4 {
			return nil, b
		}
		n := 0
		if len(b) > 0 {
			n = int(b[0] % 5)
			b = b[1:]
		}
		arr := make([]any, 0, n)
		for i := 0; i < n; i++ {
			var e any
			e, b = valueFrom(b, depth+1)
			arr = append(arr, e)
		}
		return arr, b
	default:
		if depth > 4 {
			return nil, b
		}
		n := 0
		if len(b) > 0 {
			n = int(b[0] % 5)
			b = b[1:]
		}
		m := make(map[string]any, n)
		for i := 0; i < n; i++ {
			k := ""
			if len(b) > 0 {
				k = string(take(int(b[0]%8) + 1))
			}
			var e any
			e, b = valueFrom(b, depth+1)
			m[k] = e
		}
		return m, b
	}
}

func FuzzParity(f *testing.F) {
	f.Add([]byte{2, 5, '<', '&', 0xff, '\n', '"'})
	f.Add([]byte{5, 0, 0, 0, 0, 0, 0, 0xf0, 0x7f}) // +Inf
	f.Add([]byte{7, 3, 2, 'a', 3, 1, 2, 3, 4, 5, 6, 7, 8, 1, 'b', 6, 2, 1, 0})
	f.Fuzz(func(t *testing.T, b []byte) {
		v, _ := valueFrom(b, 0)
		assertParity(t, v)
	})
}

func FuzzStringParity(f *testing.F) {
	f.Add("\u2028<\xff>\x00")
	f.Fuzz(func(t *testing.T, s string) {
		assertParity(t, s)
		assertParity(t, map[string]any{s: s})
	})
}

func FuzzFloatParity(f *testing.F) {
	f.Add(1e21)
	f.Add(1e-7)
	f.Fuzz(func(t *testing.T, x float64) {
		assertParity(t, x)
	})
}
