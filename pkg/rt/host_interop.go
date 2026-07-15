/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"fmt"

	"github.com/nooga/let-go/pkg/vm"
)

// This file provides generic JVM collection-interface interop for `.`-form
// method calls on let-go collections: the ILookup (.valAt), Iterable (.iterator
// + the returned Iterator's .hasNext/.next), Associative (.assoc),
// IPersistentCollection (.cons), Indexed (.nth), Counted (.count), and
// Object/Long (.hashCode/.longValue) methods that Clojure libraries reach on
// their :clj branches. metosin/malli uses .valAt (map validator), .iterator
// (malli.impl.util/-vmap), and the regex cache's .hashCode/.longValue.
//
// It dispatches on the value's let-go interface (Lookup/Sequable/Associative/
// Collection/Indexed/Counted), so it works for every collection type, not just
// the ones malli happens to touch.

// hostNormSeq returns v's seq as a Value, or NIL when empty/exhausted. Uses the
// Go-level seqOf so it needs no bootstrapped runtime (works in unit tests too).
// seqOf passes a LazySeq through without an emptiness check, so an empty
// (lazy-seq nil) would look non-empty — SeqIsEmpty forces and normalizes it.
func hostNormSeq(v vm.Value) vm.Value {
	if v == nil {
		return vm.NIL
	}
	s, err := seqOf(v)
	if err != nil || s == nil {
		return vm.NIL
	}
	// seqOf passes a LazySeq through unrealized, so an empty (lazy-seq nil) would
	// look non-empty. Resolve() forces the head and returns nil when it's empty.
	if ls, ok := s.(*vm.LazySeq); ok {
		s = ls.Resolve()
	}
	if vm.SeqIsEmpty(s) {
		return vm.NIL
	}
	return s
}

// hostHashOf returns a stable hash of v. A weak hash is still correct for
// malli's regex cache: it also compares candidate keys with `=`, so collisions
// only cost time.
func hostHashOf(v vm.Value) vm.Int {
	return vm.Int(int64(vm.HashValue(v)))
}

// theHostIteratorType is the ValueType for seqIterator values.
type theHostIteratorType struct{}

func (t *theHostIteratorType) String() string     { return t.Name() }
func (t *theHostIteratorType) Type() vm.ValueType { return vm.TypeType }
func (t *theHostIteratorType) Unbox() any         { return nil }
func (t *theHostIteratorType) Name() string       { return "java.util.Iterator" }
func (t *theHostIteratorType) Box(any) (vm.Value, error) {
	return vm.NIL, nil
}

var hostIteratorType = &theHostIteratorType{}

// seqIterator is a java.util.Iterator over a let-go seq. `cur` holds the current
// (non-empty) seq, or NIL once exhausted.
type seqIterator struct{ cur vm.Value }

func (it *seqIterator) Type() vm.ValueType { return hostIteratorType }
func (it *seqIterator) Unbox() any         { return it }
func (it *seqIterator) String() string     { return "#<java.util.Iterator>" }

func (it *seqIterator) InvokeMethod(name vm.Symbol, args []vm.Value) (vm.Value, error) {
	switch string(name) {
	case "hasNext":
		if len(args) == 0 {
			return vm.Boolean(it.cur != vm.NIL), nil
		}
	case "next":
		if len(args) == 0 {
			if it.cur == vm.NIL {
				return vm.NIL, fmt.Errorf("java.util.Iterator: next past end")
			}
			s, ok := it.cur.(vm.Seq)
			if !ok {
				return vm.NIL, fmt.Errorf("java.util.Iterator: not a seq")
			}
			v := s.First()
			it.cur = hostNormSeq(s.Next())
			return v, nil
		}
	}
	return vm.NIL, fmt.Errorf("java.util.Iterator has no method .%s/%d", name, len(args))
}

// hostCollectionMethod maps a JVM collection-interface `.`-method to the
// matching let-go collection operation. Returns (value, handled, error);
// handled is false when the method/receiver isn't one it covers, so the caller
// falls through to the other fallbacks.
func hostCollectionMethod(rec vm.Value, name vm.Symbol, args []vm.Value) (vm.Value, bool, error) {
	if rec == nil {
		return vm.NIL, false, nil
	}
	// Scope this shim to real collection types. vm.String satisfies Counted/
	// Indexed/Collection but with byte semantics (String.Count is len in bytes),
	// while core count/nth are rune-based — and java.lang.String isn't
	// clojure.lang.Counted anyway. Letting strings fall through keeps
	// (.count "café") = 4 (core), not 5 (bytes). Same carve-out indexed? uses.
	if _, ok := rec.(vm.String); ok {
		return vm.NIL, false, nil
	}
	switch string(name) {
	case "valAt": // clojure.lang.ILookup
		if l, ok := rec.(vm.Lookup); ok {
			switch len(args) {
			case 1:
				return l.ValueAt(args[0]), true, nil
			case 2:
				return l.ValueAtOr(args[0], args[1]), true, nil
			}
		}
	case "iterator": // java.lang.Iterable
		if len(args) == 0 {
			if _, ok := rec.(vm.Sequable); ok {
				return &seqIterator{cur: hostNormSeq(rec)}, true, nil
			}
			if _, ok := rec.(vm.Seq); ok {
				return &seqIterator{cur: hostNormSeq(rec)}, true, nil
			}
		}
	case "assoc": // clojure.lang.Associative
		if a, ok := rec.(vm.Associative); ok && len(args) == 2 {
			return a.Assoc(args[0], args[1]), true, nil
		}
	case "cons": // clojure.lang.IPersistentCollection (append/prepend per type)
		if c, ok := rec.(vm.Collection); ok && len(args) == 1 {
			return c.Conj(args[0]), true, nil
		}
	case "nth": // clojure.lang.Indexed: (nth i) or (nth i notFound)
		if idx, ok := rec.(vm.Indexed); ok && (len(args) == 1 || len(args) == 2) {
			if i, ok := args[0].(vm.Int); ok {
				if int(i) >= 0 && int(i) < idx.RawCount() {
					return idx.Nth(int(i)), true, nil
				}
				if len(args) == 2 {
					return args[1], true, nil
				}
			}
		}
	case "count": // clojure.lang.Counted
		if c, ok := rec.(vm.Counted); ok && len(args) == 0 {
			return c.Count(), true, nil
		}
	case "hashCode": // java.lang.Object
		if len(args) == 0 {
			return hostHashOf(rec), true, nil
		}
	case "longValue": // java.lang.Long
		if i, ok := rec.(vm.Int); ok && len(args) == 0 {
			return i, true, nil
		}
	}
	return vm.NIL, false, nil
}
