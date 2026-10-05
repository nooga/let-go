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

func TestUsesCountWordBoundaries(t *testing.T) {
	ensureLoader()
	got := runLispExpr(t, `
	  (pr-str (mapv ir/uses-count
	    [[] [0 0] [1] [(bit-shift-left 1 63)] [-1]
	     [-1 1] [1 0 (bit-shift-left 1 63)] [-1 -1]]))`)
	want := "[0 0 1 1 64 65 2 128]"
	if s, ok := got.(vm.String); !ok || string(s) != want {
		t.Fatalf("use counts = %v, want %s", got, want)
	}
}
