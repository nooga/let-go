/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package compiler

import (
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

func runCAS(t *testing.T, src string) vm.Value {
	t.Helper()
	cp := vm.NewConsts()
	ns := rt.NS(rt.NameCoreNS)
	ctx := NewCompiler(cp, ns)
	ch, _, err := ctx.CompileMultiple(strings.NewReader(src))
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	v, err := vm.NewFrame(ch, nil).Run()
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	return v
}

// compare-and-set! sets the atom to new only when the current value is
// identical to old, returning a boolean (Clojure identity semantics).
func TestCompareAndSet(t *testing.T) {
	// match → set, true
	if v := runCAS(t, `(let [a (atom 1)] [(compare-and-set! a 1 2) (deref a)])`); v.String() != "[true 2]" {
		t.Fatalf("match: want [true 2], got %v", v)
	}
	// mismatch → no set, false
	if v := runCAS(t, `(let [a (atom 1)] [(compare-and-set! a 99 2) (deref a)])`); v.String() != "[false 1]" {
		t.Fatalf("mismatch: want [false 1], got %v", v)
	}
	// nil old value (malli tie-the-knot pattern): first wins, second fails
	if v := runCAS(t, `(let [a (atom nil)]
                          [(compare-and-set! a nil 5) (deref a)
                           (compare-and-set! a nil 9) (deref a)])`); v.String() != "[true 5 false 5]" {
		t.Fatalf("nil-old: want [true 5 false 5], got %v", v)
	}
}
