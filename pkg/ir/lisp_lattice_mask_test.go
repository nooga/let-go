/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package ir_test

import (
	"os"
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

func TestLatticeMaskCanonicalTypes(t *testing.T) {
	ensureLoader()
	// Frozen canonical outputs cover the complete ten-bit domain, including
	// bool subsumption, numeric unions, member ordering, and any/unknown precedence.
	want, err := os.ReadFile("../../test/fixtures/ir-mask-types.edn")
	if err != nil {
		t.Fatal(err)
	}
	got := runLispExpr(t, `(pr-str (mapv ir.lattice/mask->type (range 1024)))`)
	if s, ok := got.(vm.String); !ok || string(s) != strings.TrimSpace(string(want)) {
		t.Fatal("canonical type masks differ from the exhaustive baseline")
	}
}
