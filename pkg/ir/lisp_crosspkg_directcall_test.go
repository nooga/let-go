/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package ir_test

// Direct cross-PACKAGE calls in AOT: a call to a fn in ANOTHER lowered package
// (registered in *cross-package-registry* with a :go-pkg + exported :go-name)
// lowers to a direct `pkg.Fn(ec, …)` call — NOT a runtime CachedVarFn trampoline
// — and the callee's Go package is imported. *export-lowered-fns* makes the
// callee name exported. Both default off, so core lowering is byte-identical.
//
// Anti-stub: asserts the rendered Go contains the qualified call AND the import
// (neither appears on the trampoline path), and that NO CachedVarFn trampoline
// to the registered callee is emitted.

import (
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

func TestCrossPackageDirectCallLowering(t *testing.T) {
	ensureLoader()

	// Intern lib/greet so (lib/greet x) resolves during the app build.
	// (ns …) (not bare in-ns) so clojure.core is referred.
	runLispExpr(t, `(ns xplib)`)
	runLispExpr(t, `(def greet (clojure.core/fn [x] x))`)
	runLispExpr(t, `(ns xpapp)`)

	v := runLispExpr(t, `(clojure.core/binding
	  [ir.lower-go/*cross-package-registry*
	     {(ir.lower-go/registry-key (quote xplib) "greet" 1)
	      {:go-name "Greet" :arity 1 :needs-error? true
	       :param-specs ["vm.Value"] :result-spec "vm.Value"
	       :native? false :go-pkg "example.com/m/xplib"}}
	   ir.lower-go/*export-lowered-fns* true]
	  (ir.passes.pipeline/lower-ns-to-go "xpapp" (quote xpapp)
	    (quote [(defn run [x] (xplib/greet x))])))`)

	s, ok := v.(vm.String)
	if !ok {
		t.Fatalf("expected rendered Go string, got %T: %v", v, v)
	}
	g := string(s)

	for _, want := range []string{
		"xplib.Greet(ec",        // direct cross-package call to the exported callee
		`"example.com/m/xplib"`, // the callee package is imported
		"func Run(ec",           // the caller itself is exported (*export-lowered-fns*)
	} {
		if !strings.Contains(g, want) {
			t.Fatalf("cross-package direct lowering missing %q:\n--- rendered Go ---\n%s", want, g)
		}
	}
	// Must NOT trampoline to the registered callee.
	if strings.Contains(g, `CachedVarFn(&__v_xplib_greet`) || strings.Contains(g, `"xplib", "greet"`) {
		t.Fatalf("registered cross-package callee should be a DIRECT call, not a trampoline:\n%s", g)
	}
}
