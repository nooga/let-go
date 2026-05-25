/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package ir_test

import (
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

func TestBuildFnStarPreservesNestedFnTemplateInIR(t *testing.T) {
	ensureLoader()

	fn := buildLispIR(t, `(defn make-id [] (fn* [x] x))`)
	dump := lispDump(t, fn)

	if !strings.Contains(dump, ":kind :fn-template") {
		t.Fatalf("expected nested fn* const to preserve fn template\n--- dump ---\n%s", dump)
	}
}

func TestBuildFnStarPreservesNestedMultiFnTemplateInIR(t *testing.T) {
	ensureLoader()

	fn := buildLispIR(t, `(defn make-multi [] (fn* ([] :zero) ([x] x)))`)
	dump := lispDump(t, fn)

	if !strings.Contains(dump, ":kind :multi-fn-template") {
		t.Fatalf("expected nested multi-arity fn* const to preserve multi-fn template\n--- dump ---\n%s", dump)
	}
}

func TestNestedFnTemplateStillExecutesViaBytecodeLowering(t *testing.T) {
	ensureLoader()

	result := runLispExpr(t, `(((fn* [x] (fn* [y] (+ x y))) 1) 2)`)
	if got, ok := result.(vm.Int); !ok || got != 3 {
		t.Fatalf("expected nested closure execution to produce 3, got %T %v", result, result)
	}
}

func TestNestedMultiFnTemplateStillExecutesViaBytecodeLowering(t *testing.T) {
	ensureLoader()

	result := runLispExpr(t, `(((fn* [x] (fn* ([] x) ([y] y))) 1))`)
	if got, ok := result.(vm.Int); !ok || got != 1 {
		t.Fatalf("expected captured zero-arity branch to produce 1, got %T %v", result, result)
	}

	result = runLispExpr(t, `(((fn* [x] (fn* ([] x) ([y] (+ x y)))) 1) 2)`)
	if got, ok := result.(vm.Int); !ok || got != 3 {
		t.Fatalf("expected captured one-arity branch to produce 3, got %T %v", result, result)
	}
}
