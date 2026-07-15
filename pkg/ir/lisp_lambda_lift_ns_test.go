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

	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

// compileFormToGoInNs is compileFormToGo but with *ns* bound to a non-core
// namespace via compile-form's caller-ns arg. This exposes ns-sensitivity in
// lowering: a lambda-lifted sibling lives in the SAME ns as its outer fn, so a
// reference to it must resolve against that ns — not the hardcoded "core"
// default that load-var-deref-expr applies to bare (unqualified) symbols.
func compileFormToGoInNs(t *testing.T, nsName, src string) string {
	t.Helper()
	passVarCounter++
	formName := fmt.Sprintf("*compile-form-ns-%d*", passVarCounter)

	coreNS := rt.NS(rt.NameCoreNS)

	consts := vm.NewConsts()
	c := compiler.NewCompiler(consts, coreNS)
	c.SetSource("compile-form-to-go-ns-parse")
	_, parsed, err := c.CompileMultiple(strings.NewReader(fmt.Sprintf(`(quote %s)`, src)))
	if err != nil {
		t.Fatalf("parse form: %v", err)
	}
	coreNS.Def(formName, parsed)

	consts = vm.NewConsts()
	c = compiler.NewCompiler(consts, coreNS)
	c.SetSource("compile-form-to-go-ns")
	compileExpr := fmt.Sprintf(`(do (create-ns (quote %s))
                                    (binding [ir.passes.pipeline/*target* :go]
                                      (ir.passes.pipeline/compile-form %s (the-ns (quote %s)))))`,
		nsName, formName, nsName)
	_, result, err := c.CompileMultiple(strings.NewReader(compileExpr))
	if err != nil {
		t.Fatalf("compile-form with *target* :go in ns %s: %v", nsName, err)
	}

	declVarName := fmt.Sprintf("*compile-result-ns-%d*", passVarCounter)
	coreNS.Def(declVarName, result)

	renderExpr := fmt.Sprintf(`
(let [decls (if (= :multi-fn-template (:kind %s))
              (mapv (fn* [entry] (:fn entry)) (:fns %s))
              (if (= :lowered (:status %s))
                [(:decl %s)]
                []))]
  (apply str (mapv (fn* [d] (gogen/render d)) decls)))`, declVarName, declVarName, declVarName, declVarName)

	consts = vm.NewConsts()
	c = compiler.NewCompiler(consts, coreNS)
	c.SetSource("render-decls-ns")
	_, rendered, err := c.CompileMultiple(strings.NewReader(renderExpr))
	if err != nil {
		t.Fatalf("render decls: %v", err)
	}
	if s, ok := rendered.(vm.String); ok {
		return string(s)
	}
	t.Fatalf("expected gogen/render to return string, got %T", rendered)
	return ""
}

// A lambda-lifted sibling is interned into the SAME namespace as the outer fn.
// When the outer fn references it as a value (here, passed to identity), the
// reference lowers to rt.LookupVar(<ns>, "<name>__lifted0"). The ns MUST be the
// outer fn's ns, not the "core" default. Before the fix, load-var-deref-expr
// saw a bare (unqualified) symbol and emitted LookupVar("core", ...), which
// dangles at runtime in any non-core ns (the observed AnalyzeNsForms nil-deref).
func TestLambdaLiftSiblingRefUsesOuterNs(t *testing.T) {
	ensureLoader()
	out := compileFormToGoInNs(t, "lifttest.pkg", `(defn cx [p] (identity (fn* [q] q)))`)

	if !strings.Contains(out, `"cx__lifted0"`) {
		t.Fatalf("expected a var reference to the lifted sibling cx__lifted0:\n%s", out)
	}
	if strings.Contains(out, `LookupVar("core", "cx__lifted0")`) {
		t.Fatalf("lifted sibling ref wrongly defaulted to \"core\" ns:\n%s", out)
	}
	if !strings.Contains(out, `LookupVar("lifttest.pkg", "cx__lifted0")`) {
		t.Fatalf("expected LookupVar(\"lifttest.pkg\", \"cx__lifted0\"), got:\n%s", out)
	}
}
