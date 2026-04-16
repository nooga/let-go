/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import "fmt"

// Comparator is a function that compares two Values, returning -1, 0, or 1.
type Comparator func(a, b Value) (int, error)

// DefaultCompare is the default comparator used by sorted collections.
// It handles nil, numbers, strings, keywords, symbols, booleans, and chars.
func DefaultCompare(a, b Value) (int, error) {
	if a == NIL && b == NIL {
		return 0, nil
	}
	if a == NIL {
		return -1, nil
	}
	if b == NIL {
		return 1, nil
	}
	switch va := a.(type) {
	case Int, Float, *BigInt, *Ratio, *BigDecimal:
		_ = va
		if r, err := NumLt(a, b); err == nil {
			if r {
				return -1, nil
			}
			if r2, _ := NumGt(a, b); r2 {
				return 1, nil
			}
			return 0, nil
		}
	case String:
		if vb, ok := b.(String); ok {
			switch {
			case string(va) < string(vb):
				return -1, nil
			case string(va) > string(vb):
				return 1, nil
			default:
				return 0, nil
			}
		}
	case Keyword:
		if vb, ok := b.(Keyword); ok {
			switch {
			case string(va) < string(vb):
				return -1, nil
			case string(va) > string(vb):
				return 1, nil
			default:
				return 0, nil
			}
		}
	case Symbol:
		if vb, ok := b.(Symbol); ok {
			switch {
			case string(va) < string(vb):
				return -1, nil
			case string(va) > string(vb):
				return 1, nil
			default:
				return 0, nil
			}
		}
	case Boolean:
		if vb, ok := b.(Boolean); ok {
			switch {
			case !bool(va) && bool(vb):
				return -1, nil
			case bool(va) && !bool(vb):
				return 1, nil
			default:
				return 0, nil
			}
		}
	case Char:
		if vb, ok := b.(Char); ok {
			switch {
			case rune(va) < rune(vb):
				return -1, nil
			case rune(va) > rune(vb):
				return 1, nil
			default:
				return 0, nil
			}
		}
	}
	return 0, fmt.Errorf("cannot compare %s and %s", a.Type(), b.Type())
}
