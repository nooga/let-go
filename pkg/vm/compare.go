/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import "fmt"

// Comparator is a function that compares two Values, returning -1, 0, or 1.
type Comparator func(a, b Value) (int, error)

// ComparableFallback, when set, is consulted by DefaultCompare for values it
// cannot order natively. It dispatches on a's type only (mirroring Clojure's
// compare, which calls a.compareTo(b)) and returns (result, handled, err):
// handled is false when a's type has no Comparable implementation, in which
// case DefaultCompare reports its usual error. The runtime (package rt) wires
// this at init to route through the Comparable protocol; it lives here as a
// hook to avoid a vm->rt import cycle. Because compare, sort, and sorted-set
// all funnel through DefaultCompare, registering once makes all three respect
// Comparable — matching Clojure, where sort/sorted-set use compare.
var ComparableFallback func(a, b Value) (int, bool, error)

func isSeqComparable(v Value) bool {
	switch v.(type) {
	case ArrayVector, PersistentVector, *PersistentVector, MapEntry:
		return true
	}
	return false
}

func seqFrom(v Value) Seq {
	if s, ok := v.(Sequable); ok {
		return s.Seq()
	}
	if s, ok := v.(Seq); ok {
		return s
	}
	return nil
}

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
			return compareNamed(va.Namespace(), va.Name(), vb.Namespace(), vb.Name()), nil
		}
	case Symbol:
		if vb, ok := b.(Symbol); ok {
			return compareNamed(va.Namespace(), va.Name(), vb.Namespace(), vb.Name()), nil
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
	case *Instant:
		if vb, ok := b.(*Instant); ok {
			am, bm := va.t.UnixMilli(), vb.t.UnixMilli()
			switch {
			case am < bm:
				return -1, nil
			case am > bm:
				return 1, nil
			default:
				return 0, nil
			}
		}
	}
	// Vectors/sequential: lexicographic comparison
	if isSeqComparable(a) && isSeqComparable(b) {
		as := seqFrom(a)
		bs := seqFrom(b)
		for as != nil && bs != nil {
			c, err := DefaultCompare(as.First(), bs.First())
			if err != nil {
				return 0, err
			}
			if c != 0 {
				return c, nil
			}
			as = as.Next()
			bs = bs.Next()
		}
		if as == nil && bs == nil {
			return 0, nil
		}
		if as == nil {
			return -1, nil
		}
		return 1, nil
	}
	if ComparableFallback != nil {
		if c, handled, err := ComparableFallback(a, b); handled || err != nil {
			return c, err
		}
	}
	return 0, fmt.Errorf("cannot compare %s and %s", a.Type(), b.Type())
}

func compareNamed(aNs, aName, bNs, bName Value) int {
	if c := compareNamePart(aNs, bNs); c != 0 {
		return c
	}
	return compareNamePart(aName, bName)
}

func compareNamePart(a, b Value) int {
	as, bs := "", ""
	if a != NIL {
		as = string(a.(String))
	}
	if b != NIL {
		bs = string(b.(String))
	}
	switch {
	case as < bs:
		return -1
	case as > bs:
		return 1
	default:
		return 0
	}
}
