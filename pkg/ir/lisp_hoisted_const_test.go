/*
 * Copyright (c) 2026 let-go contributors; see CONTRIBUTORS.
 * SPDX-License-Identifier: MIT
 */

package ir_test

import (
	"go/ast"
	"go/token"
	"strings"
	"testing"
)

// A constant collection literal lowers to one package-level var that the
// function reads, never to a constructor call inside the function: the value
// is shared by every evaluation, as the bytecode VM and Clojure share it. Each
// literal is its own var, so two literals of equal content in two fns are two
// values, as they are two constants on the VM. The empty vector is
// vm.EmptyVector itself. A literal with a non-constant element, or with
// metadata, is still built where it is evaluated.
func TestConstantCollectionLiteralsAreHoisted(t *testing.T) {
	ensureLoader()
	source := runLispString(t, `(do
	  (create-ns 'hoistconst)
	  (ir.passes.pipeline/lower-ns-to-go "hoistconst" 'hoistconst
	    '[(defn mk-vec [] [1 2 3])
	      (defn mk-vec-again [] [1 2 3])
	      (defn mk-map [] {:a 1})
	      (defn mk-set [] #{:a})
	      (defn mk-nested [] [1 {:k [2]}])
	      (defn mk-empty [] [])
	      (defn mk-fresh [x] [x 1])
	      (defn mk-meta [] ^{:t 1} [1 2])]))`)
	file := parseLoweredGo(t, source)

	hoisted := map[string]bool{}
	for _, d := range file.Decls {
		gd, ok := d.(*ast.GenDecl)
		if !ok || gd.Tok != token.VAR {
			continue
		}
		for _, sp := range gd.Specs {
			if vs, ok := sp.(*ast.ValueSpec); ok {
				for _, n := range vs.Names {
					if strings.HasPrefix(n.Name, "__k_") {
						if len(vs.Values) != 1 {
							t.Fatalf("%s has no initializer\n%s", n.Name, source)
						}
						hoisted[n.Name] = true
					}
				}
			}
		}
	}

	// ctors are the collection constructions in fn's body (vm constructors
	// and the vector/array-map/hash-set builtins the lowering direct-calls);
	// vars are the distinct hoisted constants it reads.
	ctorCalls := func(fn *ast.FuncDecl) (ctors []string, vars []string) {
		seen := map[string]bool{}
		ast.Inspect(fn.Body, func(x ast.Node) bool {
			switch n := x.(type) {
			case *ast.CallExpr:
				for _, c := range []string{"NewArrayVector", "NewPersistentMap", "NewSet", "NewList"} {
					if isSelector(n.Fun, "vm", c) {
						ctors = append(ctors, c)
					}
				}
				for _, c := range []string{"Vector", "ArrayMap", "HashSet"} {
					if isSelector(n.Fun, "builtins", c) {
						ctors = append(ctors, c)
					}
				}
			case *ast.Ident:
				if strings.HasPrefix(n.Name, "__k_") && !seen[n.Name] {
					seen[n.Name] = true
					vars = append(vars, n.Name)
				}
			}
			return true
		})
		return
	}
	fnNamed := func(name string) *ast.FuncDecl {
		fn, ok := findFunc(file, func(n string) bool { return n == name })
		if !ok {
			t.Fatalf("%s was not lowered\n%s", name, source)
		}
		return fn
	}

	for _, name := range []string{"MkVec", "MkMap", "MkSet", "MkNested"} {
		ctors, vars := ctorCalls(fnNamed(name))
		if len(ctors) != 0 || len(vars) != 1 || !hoisted[vars[0]] {
			t.Errorf("%s: constructors %v, hoisted reads %v; want no constructor and one read of a hoisted var\n%s", name, ctors, vars, source)
		}
	}
	_, v1 := ctorCalls(fnNamed("MkVec"))
	_, v2 := ctorCalls(fnNamed("MkVecAgain"))
	if len(v1) != 1 || len(v2) != 1 || v1[0] == v2[0] {
		t.Errorf("equal literals in two fns must be two vars: %v vs %v", v1, v2)
	}
	if ctors, vars := ctorCalls(fnNamed("MkEmpty")); len(ctors) != 0 || len(vars) != 0 || !strings.Contains(source, "vm.EmptyVector") {
		t.Errorf("MkEmpty: want vm.EmptyVector, got constructors %v, vars %v\n%s", ctors, vars, source)
	}
	if ctors, vars := ctorCalls(fnNamed("MkFresh")); len(ctors) == 0 || len(vars) != 0 {
		t.Errorf("MkFresh: constructors %v, hoisted reads %v; want a construction and no hoisted read\n%s", ctors, vars, source)
	}
	// A literal with metadata applies with-meta on every evaluation to its
	// (hoisted) constant parts, so the result is fresh each time.
	ctors, vars := ctorCalls(fnNamed("MkMeta"))
	if len(ctors) != 0 || len(vars) != 2 || !strings.Contains(source, "with_meta") && !strings.Contains(source, "WithMeta") {
		t.Errorf("MkMeta: constructors %v, hoisted reads %v; want with-meta applied per call to two hoisted constants\n%s", ctors, vars, source)
	}
}
