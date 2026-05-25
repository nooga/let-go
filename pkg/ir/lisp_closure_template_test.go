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

func TestNestedFnTemplateStillExecutesViaBytecodeLowering(t *testing.T) {
	ensureLoader()

	result := runLispExpr(t, `(((fn* [x] (fn* [y] (+ x y))) 1) 2)`)
	if got, ok := result.(vm.Int); !ok || got != 3 {
		t.Fatalf("expected nested closure execution to produce 3, got %T %v", result, result)
	}
}
