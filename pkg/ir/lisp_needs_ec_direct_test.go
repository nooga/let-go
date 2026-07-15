/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package ir_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/rt"
)

// TestNeedsEcEmitsEcArg verifies that when a native module function has
// NeedsEC=true, the lowered Go call includes the ec argument as the first
// parameter: pkg.Fn(ec, args...) instead of pkg.Fn(args...).
func TestNeedsEcEmitsEcArg(t *testing.T) {
	ensureLoader()

	// Register a test native module with one function that needs ec, in clojure.core
	testPkg := "github.com/nooga/let-go/test/needs_ec_test"
	rt.RegisterNativeModule(&rt.NativeModule{
		GoPkg:     testPkg,
		Namespace: "clojure.core",
		Fns: map[string]rt.NativeDirectFn{
			"ec-test-echo": {
				GoIdent:    "EcTestEcho",
				Arity:      1,
				Variadic:   false,
				ParamSpecs: []string{"vm.Value"},
				ResultSpec: "vm.Value",
				NeedsError: true,
				NeedsEC:    true,
			},
		},
	})

	// Define the native function as a var in clojure.core so IR build can find it
	runLispExpr(t, `(def ec-test-echo (fn [x] x))`)

	// Build IR for a function that calls the ec-needing native function
	fn := buildLispIR(t, `(defn use-ec-test-echo [x] (ec-test-echo x))`)
	optimizeLispIR(t, fn)
	passVarCounter++
	varName := fmt.Sprintf("*needs-ec-fn-%d*", passVarCounter)
	rt.NS(rt.NameCoreNS).Def(varName, fn)

	// Lower with native registry seeded
	result := lowerWithNativeRegistry(t, varName)
	rendered := bindAndRenderGoDecl(t, result)

	// The call should include ec as the first argument
	if !strings.Contains(rendered, "EcTestEcho(ec,") && !strings.Contains(rendered, "EcTestEcho(ec )") {
		t.Fatalf("expected EcTestEcho to be called with ec as first arg (EcTestEcho(ec, ...)):\n--- go ---\n%s", rendered)
	}

	// Should not have the ec-free form
	if strings.Contains(rendered, "EcTestEcho(") && !strings.Contains(rendered, "EcTestEcho(ec") {
		t.Fatalf("expected call to include ec, but found EcTestEcho without ec:\n--- go ---\n%s", rendered)
	}
}
