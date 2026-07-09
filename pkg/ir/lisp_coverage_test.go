/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package ir_test

import (
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

func TestCoverageClassifierBuckets(t *testing.T) {
	ensureLoader()
	// Load the new ns explicitly (ensureLoader loads the IR set; coverage is new).
	runLispExpr(t, `(require 'ir.coverage)`)
	cases := map[string]string{
		"unresolved symbol set! in this context": ":missing-form/set!",
		"unresolved symbol frobnicate here":      ":unresolved/frobnicate",
		"unrecognized form somewhere":            ":build/unrecognized-form",
	}
	for in, want := range cases {
		expr := `(pr-str (ir.coverage/classify-error "` + in + `"))`
		got := runLispExpr(t, expr)
		if s, ok := got.(vm.String); !ok || string(s) != want {
			t.Fatalf("classify-error(%q) = %v, want %s", in, got, want)
		}
	}
}
