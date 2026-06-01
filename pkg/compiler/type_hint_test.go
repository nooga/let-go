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

func runTH(t *testing.T, src string) vm.Value {
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

// A type hint naming a host class (`^Iterable`, `^String`, …) must not require
// that symbol to resolve as a var — let-go ignores type hints, but they appear
// throughout real Clojure source.
func TestTypeHint_UnresolvedTagIsIgnored(t *testing.T) {
	// Hint on an expression whose tag (Iterable) is undefined.
	if v := runTH(t, `(count ^Iterable [1 2 3])`); v.String() != "3" {
		t.Fatalf("^Iterable hint: want 3, got %v", v)
	}
	// Hint on a fn parameter.
	if v := runTH(t, `((fn [^String s] (str s "!")) "hi")`); v.String() != `"hi!"` {
		t.Fatalf("^String param hint: want \"hi!\", got %v", v)
	}
	// Hint on a let binding name (unwrapped, like fn params).
	if v := runTH(t, `(let [^CharSequence x "ok"] x)`); v.String() != `"ok"` {
		t.Fatalf("^CharSequence let-name hint: want \"ok\", got %v", v)
	}
	// Hint on a let binding value (expression position).
	if v := runTH(t, `(let [x ^Iterable [1 2]] (count x))`); v.String() != "2" {
		t.Fatalf("^Iterable let-value hint: want 2, got %v", v)
	}
}

// Expression-position type hints are PRESERVED as :tag metadata (the symbol
// datum), so the IR/typeinfer can use them — they are not discarded.
func TestTypeHint_PreservedAsTagMetadata(t *testing.T) {
	if v := runTH(t, `(:tag (meta ^Foo [1 2 3]))`); v.String() != "Foo" {
		t.Fatalf("expected :tag metadata Foo to be preserved, got %v", v)
	}
}
