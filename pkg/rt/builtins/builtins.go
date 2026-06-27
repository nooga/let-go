/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

// Package builtins is the single-source home for direct-call-eligible
// clojure.core builtins: context-free functions with uniform vm.Value
// signatures. Each function here is meant to be BOTH the interpreter
// registration (lang.go boxes it) AND the AOT direct-call target (lang.go
// registers it as a native module, generated code calls builtins.X directly).
//
// It imports ONLY vm (never rt) so rt/lang.go can import it without a cycle —
// that layering is what lets lang.go register from these funcs instead of
// duplicating their bodies (the corefns/diff-test pattern this replaces).
package builtins

import (
	"fmt"

	"github.com/nooga/let-go/pkg/vm"
)

// seqOf mirrors the unexported pkg/rt/lang.go helper of the same name.
func seqOf(v vm.Value) (vm.Seq, error) {
	if v == vm.NIL || v == vm.EmptyList {
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

// Vector mirrors clojure.core/vector — `(vector & xs)`. Builtin is
// vm.NativeFnType.WrapNoErr(vm.NewArrayVector); same delegation as List.
func Vector(args ...vm.Value) (vm.Value, error) {
	return vm.NewArrayVector(args), nil
}

// Not mirrors clojure.core/not — `(not x)`.
func Not(v vm.Value) (vm.Value, error) {
	return vm.Boolean(!vm.IsTruthy(v)), nil
}

// Cons mirrors clojure.core/cons — `(cons x coll)`.
func Cons(elem, coll vm.Value) (vm.Value, error) {
	if coll == vm.NIL {
		return vm.NewCons(elem, nil), nil
	}
	seq, err := seqOf(coll)
	if err != nil {
		return vm.NIL, fmt.Errorf("cons expected Seq")
	}
	if seq == nil {
		return vm.NewCons(elem, nil), nil
	}
	return vm.NewCons(elem, seq), nil
}

// Contains mirrors clojure.core/contains? — `(contains? coll k)`.
func Contains(coll, k vm.Value) (vm.Value, error) {
	if coll == vm.NIL {
		return vm.FALSE, nil
	}
	if s, ok := coll.(vm.Keyed); ok {
		return s.Contains(k), nil
	}
	if idx, ok := k.(vm.Int); ok {
		i := int(idx)
		if c, ok := coll.(vm.Counted); ok {
			return vm.Boolean(i >= 0 && i < c.RawCount()), nil
		}
	}
	if s, ok := coll.(vm.String); ok {
		if idx, ok := k.(vm.Int); ok {
			return vm.Boolean(int(idx) >= 0 && int(idx) < len([]rune(string(s)))), nil
		}
		return vm.NIL, fmt.Errorf("contains? on a string requires an integer key, got %s", k.Type().Name())
	}
	return vm.FALSE, nil
}
