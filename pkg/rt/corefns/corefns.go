/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

// Package corefns is a placeholder while the lginterop-driven direct-call
// path is being wired up. It currently exports just `Seq` so the dynamic
// native-direct registry plumbing (registry, pipeline-seeding,
// lower-go emit) has at least one entry to exercise.
//
// The long-term shape is: corefns hosts real Go top-level funcs with
// typed Go signatures (not the vm.Value-only shape used here). The
// metadata the registry holds is whatever lginterop's go/types pass
// extracts — typed params, typed results, variadic, &c. lower-go reads
// the typed signature and emits per-arg coercion + direct call, mirroring
// what vm.NativeFnType.Box already does reflectively at the boundary.
package corefns

import (
	"fmt"

	"github.com/nooga/let-go/pkg/vm"
)

// Seq returns a seq of coll, or nil for empty/nil colls. Mirrors the
// closure body previously inlined in pkg/rt/lang.go around the `seq`
// registration.
//
//lg:native
//lg:ns clojure.core
func Seq(v vm.Value) (vm.Value, error) {
	if v == vm.NIL {
		return vm.NIL, nil
	}
	switch v.(type) {
	case *vm.Cons, *vm.LazySeq:
		// fall through
	default:
		if coll, ok := v.(vm.Collection); ok {
			if coll.RawCount() == 0 {
				return vm.NIL, nil
			}
		}
	}
	var n vm.Seq
	if sqbl, ok := v.(vm.Sequable); ok {
		n = sqbl.Seq()
	} else if s, ok := v.(vm.Seq); ok {
		n = s
	} else {
		return vm.NIL, fmt.Errorf("seq expected Seqable")
	}
	if n == nil || n == vm.EmptyList {
		return vm.NIL, nil
	}
	return n, nil
}

// seqOf mirrors the unexported pkg/rt/lang.go helper of the same name —
// corefns cannot import rt (rt imports corefns), so the logic is duplicated.
// The corefns_diff_test guards against drift from the original.
func seqOf(v vm.Value) (vm.Seq, error) {
	if v == vm.NIL {
		return nil, nil
	}
	if v == vm.EmptyList {
		return nil, nil
	}
	if _, isLazy := v.(*vm.LazySeq); !isLazy {
		if sq, ok := v.(vm.Sequable); ok {
			s := sq.Seq()
			if s == nil || s == vm.EmptyList {
				return nil, nil
			}
			return s, nil
		}
	}
	if s, ok := v.(vm.Seq); ok {
		return s, nil
	}
	return nil, fmt.Errorf("don't know how to create ISeq from %s", v.Type())
}

// First mirrors the clojure.core/first builtin (pkg/rt/lang.go).
//
//lg:native
//lg:ns clojure.core
func First(v vm.Value) (vm.Value, error) {
	if v == vm.NIL {
		return vm.NIL, nil
	}
	if seq, ok := v.(vm.Seq); ok {
		return seq.First(), nil
	}
	if sq, ok := v.(vm.Sequable); ok {
		s := sq.Seq()
		if s == nil || s == vm.EmptyList {
			return vm.NIL, nil
		}
		return s.First(), nil
	}
	return vm.NIL, fmt.Errorf("first expected Seq")
}

// Second mirrors the clojure.core/second builtin (pkg/rt/lang.go).
//
//lg:native
//lg:ns clojure.core
func Second(v vm.Value) (vm.Value, error) {
	if v == vm.NIL {
		return vm.NIL, nil
	}
	seq, err := seqOf(v)
	if err != nil {
		return vm.NIL, fmt.Errorf("second expected Seq")
	}
	if seq == nil {
		return vm.NIL, nil
	}
	n := seq.Next()
	if n == nil {
		return vm.NIL, nil
	}
	return n.First(), nil
}

// Next mirrors the clojure.core/next builtin (pkg/rt/lang.go).
//
//lg:native
//lg:ns clojure.core
func Next(v vm.Value) (vm.Value, error) {
	if v == vm.NIL {
		return vm.NIL, nil
	}
	seq, err := seqOf(v)
	if err != nil {
		return vm.NIL, fmt.Errorf("next expected Seq")
	}
	if seq == nil {
		return vm.NIL, nil
	}
	n := seq.Next()
	if n == nil {
		return vm.NIL, nil
	}
	return n, nil
}

// Rest mirrors the clojure.core/rest builtin (pkg/rt/lang.go).
//
//lg:native
//lg:ns clojure.core
func Rest(v vm.Value) (vm.Value, error) {
	if v == vm.NIL {
		return vm.EmptyList, nil
	}
	s, err := seqOf(v)
	if err != nil {
		return vm.NIL, fmt.Errorf("rest expected Seq")
	}
	if s == nil {
		return vm.EmptyList, nil
	}
	return s.More(), nil
}

// Count mirrors the clojure.core/count builtin (pkg/rt/lang.go).
//
//lg:native
//lg:ns clojure.core
func Count(v vm.Value) (vm.Value, error) {
	if v == vm.NIL {
		return vm.MakeInt(0), nil
	}
	if s, ok := v.(vm.String); ok {
		return vm.MakeInt(len([]rune(string(s)))), nil
	}
	seq, ok := v.(vm.Counted)
	if !ok {
		return vm.NIL, fmt.Errorf("count expected Counted")
	}
	return seq.Count(), nil
}
