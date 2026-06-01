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

func runDT(t *testing.T, src string) vm.Value {
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

// 3-arg deref: (deref blocking-ref timeout-ms timeout-val) blocks up to ms for
// a promise/future, returning timeout-val on timeout.
func TestDerefTimeout(t *testing.T) {
	// Undelivered promise → timeout value.
	if v := runDT(t, `(deref (promise) 20 :timed-out)`); v.String() != ":timed-out" {
		t.Fatalf("undelivered: want :timed-out, got %v", v)
	}
	// Already-delivered promise → value (immediately).
	if v := runDT(t, `(let [p (promise)] (deliver p 42) (deref p 1000 :to))`); v.String() != "42" {
		t.Fatalf("delivered: want 42, got %v", v)
	}
	// future delivers within the timeout window.
	if v := runDT(t, `(deref (future 7) 2000 :to)`); v.String() != "7" {
		t.Fatalf("future: want 7, got %v", v)
	}
	// 1-arg deref still works.
	if v := runDT(t, `(let [p (promise)] (deliver p 5) (deref p))`); v.String() != "5" {
		t.Fatalf("1-arg deref: want 5, got %v", v)
	}
}
