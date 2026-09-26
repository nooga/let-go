/*
 * Copyright (c) 2026 let-go contributors; see CONTRIBUTORS.
 * SPDX-License-Identifier: MIT
 */

package ir_test

import (
	"go/ast"
	"strings"
	"testing"
)

func TestGuardedNativeCallSharesClosureArgument(t *testing.T) {
	ensureLoader()
	source := runLispString(t, `(do
      (require 'ir.passes.pipeline)
      (create-ns 'closureguard)
      (ir.passes.pipeline/lower-ns-to-go "closureguard" 'closureguard
        '[(defn store-thunk [a x] (reset! a (fn [] (fn [] x))))]))`)
	fn, found := findFunc(parseLoweredGo(t, source), func(name string) bool { return name == "StoreThunk" })
	if !found {
		t.Fatal("StoreThunk was not lowered")
	}
	closures := 0
	ast.Inspect(fn.Body, func(node ast.Node) bool {
		if _, ok := node.(*ast.FuncLit); ok {
			closures++
		}
		return true
	})
	if closures != 2 {
		t.Fatalf("got %d closure bodies, want 2 (one per source closure)\n%s", closures, source)
	}
	for _, call := range []string{"rt.CoreReset(", "rt.NativePrimsIntact()", "ec.Invoke("} {
		if !strings.Contains(source, call) {
			t.Errorf("missing %s in generated source", call)
		}
	}
}

func TestSharedClosureTempAvoidsSourceLocal(t *testing.T) {
	ensureLoader()
	source := runLispString(t, `(do
      (create-ns 'closurecollision)
      (ir.passes.pipeline/lower-ns-to-go "closurecollision" 'closurecollision
        '[(defn store-thunk [a x]
            (let [ca8_1 (identity x)]
              (reset! a (fn [] (fn [] ca8_1)))))]))`)
	fn, found := findFunc(parseLoweredGo(t, source), func(name string) bool { return name == "StoreThunk" })
	if !found {
		t.Fatal("StoreThunk was not lowered")
	}
	names := map[string]bool{}
	ast.Inspect(fn.Body, func(node ast.Node) bool {
		if _, nested := node.(*ast.FuncLit); nested {
			return false
		}
		if spec, ok := node.(*ast.ValueSpec); ok {
			for _, name := range spec.Names {
				if names[name.Name] {
					t.Fatalf("duplicate local %s\n%s", name.Name, source)
				}
				names[name.Name] = true
			}
		}
		return true
	})
}

func TestGuardedArgumentsPreserveLoweringBindings(t *testing.T) {
	ensureLoader()
	got := runLispString(t, `(let [emit (deref (get (ns-interns 'ir.lower-go) 'shared-arg-exprs))]
      (binding [ir.lower-go/*closure-arg-prefix* "captured_"]
      (gogen/render
        (first (emit {} [0]
          (fn [_] (gogen/ident ir.lower-go/*closure-arg-prefix*)))))))`)
	if strings.TrimSpace(got) != "captured_" {
		t.Fatalf("argument emitter lost its lowering binding: %q", got)
	}
}
