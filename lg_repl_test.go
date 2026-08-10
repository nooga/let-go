//go:build !plan9 && !js && !wasip1

/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package main

import (
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

func testCompleter(t *testing.T) *completer {
	t.Helper()
	ns := rt.NS("user")
	if ns == nil {
		t.Fatal("user namespace not found")
	}
	return &completer{ctx: compiler.NewCompiler(vm.NewConsts(), ns)}
}

// Tab completion walks the refer graph, which is cyclic (clojure.core requires
// let-go.types, RegisterNS refers clojure.core back into it). Completing at all
// is the assertion here — before the visited set this overflowed the stack.
func TestCompleterDoReturnsUntypedRemainder(t *testing.T) {
	c := testCompleter(t)

	line := []rune("(ns-unm")
	candidates, length := c.Do(line, len(line))
	if len(candidates) == 0 {
		t.Fatal("no candidates for ns-unm")
	}
	if length != len("ns-unm") {
		t.Fatalf("length = %d, want %d", length, len("ns-unm"))
	}

	// readline inserts the candidate at the cursor, so each one must be only the
	// part not typed yet — a whole symbol here would render as "ns-unmns-unmap".
	found := false
	for _, cand := range candidates {
		if strings.HasPrefix(string(cand), "ns-unm") {
			t.Fatalf("candidate %q repeats the typed prefix", string(cand))
		}
		if string(cand) == "ap " {
			found = true
		}
	}
	if !found {
		t.Fatalf("ns-unmap missing from candidates: %q", candidates)
	}
}

func TestCompleterDoUnknownNamespace(t *testing.T) {
	c := testCompleter(t)
	line := []rune("nosuchns/x")
	if candidates, _ := c.Do(line, len(line)); len(candidates) != 0 {
		t.Fatalf("want no candidates for an unknown namespace, got %q", candidates)
	}
}
