/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package ir_test

import (
	"regexp"
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

// Part 3 of the lambda-lift native-safety work: the MAIN fn of a lambda-lift
// multi-fn-template must carry its full lowering result through the flatten,
// exactly like its lifted siblings (Part 2). Otherwise the main is stripped to
// a bare decl — no :fn-name/:override-eligible? — so it is (a) never registered
// as a native override, and (b) invisible to the direct-call registry, forcing
// every intra-ns caller through the rt.CachedVarFn bytecode trampoline even
// though a native decl for it sits in the same package (self-trampolining).
func TestLambdaLiftMainIsRegisteredAndDirectCallable(t *testing.T) {
	ensureLoader()

	// Define the fns in their namespace first (interns liftmain/cx so use-cx's
	// (cx x) resolves during build), then lower — mirroring how the bootstrap
	// lowers a fully-loaded namespace.
	v := runLispExpr(t, `
      (ns liftmain)
      (defn cx [p] (identity (fn* [q] q)))
      (defn use-cx [x] (cx x))
      (ns user)
      (ir.passes.pipeline/lower-ns-to-go "liftmain" (quote liftmain)
        [(quote (defn cx [p] (identity (fn* [q] q))))
         (quote (defn use-cx [x] (cx x)))])`)
	s, ok := v.(vm.String)
	if !ok {
		t.Fatalf("expected rendered Go source string, got %T", v)
	}
	src := string(s)

	// Sanity: main + sibling decls exist, sibling is registered (Parts 1+2).
	if !strings.Contains(src, "func Cx(") || !strings.Contains(src, "func Cx__lifted0(") {
		t.Fatalf("expected both main and sibling decls:\n%s", src)
	}
	if !regexp.MustCompile(`"cx__lifted0":\s*__gogen_wrap`).MatchString(src) {
		t.Fatalf("sibling cx__lifted0 not registered (Part 2 regressed?):\n%s", src)
	}

	// Part 3a: the main itself is registered as a native override.
	if !regexp.MustCompile(`"cx":\s*__gogen_wrap`).MatchString(src) {
		t.Fatalf("main cx is NOT registered in the override init map:\n%s", src)
	}

	// Part 3b: an intra-ns caller direct-calls the native main — no trampoline.
	if !regexp.MustCompile(`=\s*Cx\(ec,`).MatchString(src) {
		t.Fatalf("use-cx does not direct-call Cx natively:\n%s", src)
	}
	if strings.Contains(src, `"liftmain", "cx")`) {
		t.Fatalf("use-cx still trampolines to cx via var lookup:\n%s", src)
	}
}
