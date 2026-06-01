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

func runSource(t *testing.T, src string) (vm.Value, error) {
	t.Helper()
	cp := vm.NewConsts()
	ns := rt.NS(rt.NameCoreNS)
	ctx := NewCompiler(cp, ns)
	chunk, _, err := ctx.CompileMultiple(strings.NewReader(src))
	if err != nil {
		return vm.NIL, err
	}
	return vm.NewFrame(chunk, nil).Run()
}

// set-field! mutates a deftype instance field in place; the new value is
// observable through the existing (.field instance) read path. This is the VM
// substrate for cljs-style ^:mutable deftype fields.
func TestSetField_MutateAndReadBack(t *testing.T) {
	src := `(do
              (def T (make-deftype "T" 'x 'y))
              (def t (make-deftype-instance T 1 2))
              (def before (.x t))
              (set-field! t 'x 99)
              (set-field! t 'y 7)
              (and (= before 1) (= (.x t) 99) (= (.y t) 7)))`
	val, err := runSource(t, src)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if val != vm.TRUE {
		t.Fatalf("expected TRUE (live mutate + read-back), got %v", val)
	}
}

// set-field! returns the assigned value.
func TestSetField_ReturnsValue(t *testing.T) {
	src := `(do
              (def T (make-deftype "T" 'x))
              (def t (make-deftype-instance T 1))
              (= 42 (set-field! t 'x 42)))`
	val, err := runSource(t, src)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if val != vm.TRUE {
		t.Fatalf("set-field! should return the assigned value, got %v", val)
	}
}

// set-field! on an unknown field is a runtime error.
func TestSetField_UnknownFieldErrors(t *testing.T) {
	src := `(do
              (def T (make-deftype "T" 'x))
              (def t (make-deftype-instance T 1))
              (set-field! t 'nope 5))`
	if _, err := runSource(t, src); err == nil {
		t.Fatal("expected error setting unknown field, got nil")
	}
}

// set-field! rejects a non-deftype-instance target.
func TestSetField_NonInstanceErrors(t *testing.T) {
	src := `(set-field! 5 'x 1)`
	if _, err := runSource(t, src); err == nil {
		t.Fatal("expected error on non-instance target, got nil")
	}
}
