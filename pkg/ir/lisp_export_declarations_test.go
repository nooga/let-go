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

func TestExportDiscoverySeparatesDeclarations(t *testing.T) {
	ensureLoader()
	var source vm.Value
	out := captureLetGoOut(t, func() {
		source = runLispExpr(t, `
		  (create-ns 'exportdecls)
		  (first (ir.passes.pipeline/lower-all-ns-to-go
		    [["exportdecls" 'exportdecls
		      '[(defprotocol Marker (value [this])) (defn keep-value [x] x)]
		      "example/exportdecls"]]))`)
	})
	if strings.Contains(out, "WARNING:") {
		t.Errorf("valid declaration produced a lowering warning: %s", out)
	}
	s, ok := source.(vm.String)
	if !ok || !strings.Contains(string(s), "type Marker interface") || !strings.Contains(string(s), "func KeepValue(") {
		t.Fatalf("declaration or function missing from emitted Go: %v", source)
	}
}

func TestExportDiscoveryReportsLoweringErrors(t *testing.T) {
	ensureLoader()
	out := captureLetGoOut(t, func() {
		runLispExpr(t, `
		  (create-ns 'exporterrors)
		  (ir.passes.pipeline/lower-all-ns-to-go
		    [["exporterrors" 'exporterrors
		      '[(defn unresolved [x] (missing-export-callee x))]
		      "example/exporterrors"]])`)
	})
	if !strings.Contains(out, "WARNING: lowering failed for unresolved") {
		t.Fatalf("missing diagnostic for invalid function: %s", out)
	}
}
