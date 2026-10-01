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

const closureCtorPrefix = "newClosure_"

// closureSites counts the closure creations in node, nested closures
// included: a call to a closure constructor (one allocation holding the
// NativeFn and its captures) or, for a literal the struct path does not
// cover, an rt.BoxNativeFn call. Each is one closure allocation site at run
// time. The constructors called are returned by name.
func closureSites(node ast.Node) (ctors []string, boxes int) {
	ast.Inspect(node, func(x ast.Node) bool {
		ce, ok := x.(*ast.CallExpr)
		if !ok {
			return true
		}
		if isSelector(ce.Fun, "rt", "BoxNativeFn") {
			boxes++
		}
		if id, ok := ce.Fun.(*ast.Ident); ok && strings.HasPrefix(id.Name, closureCtorPrefix) {
			ctors = append(ctors, id.Name)
		}
		return true
	})
	return ctors, boxes
}

// A capturing closure is created once, where its fn form sits, and every
// later use reads that one value. The bytecode VM allocates exactly one Fn per
// evaluation of the form; the lowered Go must match it, both so `identical?`
// agrees and so a use never costs an allocation. Each case names the number
// of closure forms in the source, which is the number of closure creations
// the lowered function may contain; each creation is a closure-struct
// constructor, never an rt.BoxNativeFn (a Go func plus a NativeFn: two
// allocations).
func TestCapturingClosureBoxedOncePerForm(t *testing.T) {
	ensureLoader()
	cases := []struct {
		name  string
		goFn  string
		form  string
		sites int
	}{
		{"let-bound local used twice", "SameLocalCap",
			`(defn same-local-cap [y] (let [f (fn [x] (+ x y))] (identical? f f)))`, 1},
		{"loop-carried local", "LoopCarriedSame",
			`(defn loop-carried-same [y]
			   (loop [f (fn [x] (+ x y)) i 0]
			     (if (< i 2) (recur f (inc i)) (identical? f f))))`, 1},
		{"captured by another closure", "CapturedByClosure",
			`(defn captured-by-closure [y]
			   (let [f (fn [x] (+ x y)) g (fn [] f)] (identical? f (g))))`, 2},
		{"stored twice in a collection", "StoredInCollection",
			`(defn stored-in-collection [y]
			   (let [f (fn [x] (+ x y)) v [f f]] (identical? (nth v 0) (nth v 1))))`, 1},
		{"variadic and mixed arities", "Variadics",
			`(defn variadics [y]
			   (let [f (fn [& xs] (cons y xs)) g (fn [a & xs] (cons a (cons y xs)))] [f g]))`, 2},
		{"multi-arity literal, one struct per arm", "MultiArity",
			`(defn multi-arity [y] (fn ([a] (+ a y)) ([a b] (+ a (+ b y)))))`, 2},
		{"nested literal, the inner created per outer call", "Nested",
			`(defn nested [y] (fn [a] (fn [b] (+ a (+ b y)))))`, 1},
		{"capture-free variadic literal stays inline", "FreeVariadic",
			`(defn free-variadic [] (let [f (fn [& xs] xs)] (f 1)))`, 1},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			source := runLispString(t, `(do
			  (create-ns 'closureonce)
			  (ir.passes.pipeline/lower-ns-to-go "closureonce" 'closureonce '[`+tc.form+`]))`)
			file := parseLoweredGo(t, source)
			fn, found := findFunc(file, func(name string) bool { return name == tc.goFn })
			if !found {
				t.Fatalf("%s was not lowered\n%s", tc.goFn, source)
			}
			ctors, boxes := closureSites(fn.Body)
			if boxes != 0 {
				t.Errorf("got %d rt.BoxNativeFn calls, want 0: every literal is a closure struct\n%s", boxes, source)
			}
			if len(ctors) != tc.sites {
				t.Fatalf("got %d closure constructions %v, want %d (one per closure form)\n%s", len(ctors), ctors, tc.sites, source)
			}
			for _, ctor := range ctors {
				assertClosureStruct(t, file, ctor, source)
			}
		})
	}
}

// assertClosureStruct checks the three decls behind a constructor: a struct
// type embedding vm.NativeFn, the constructor allocating exactly that struct
// once and returning its embedded NativeFn, and a CallNative method on it.
func assertClosureStruct(t *testing.T, file *ast.File, ctor, source string) {
	t.Helper()
	typeName := "closure_" + strings.TrimPrefix(ctor, closureCtorPrefix)
	var typeDecl *ast.TypeSpec
	var method, ctorDecl *ast.FuncDecl
	for _, d := range file.Decls {
		switch d := d.(type) {
		case *ast.GenDecl:
			for _, sp := range d.Specs {
				if ts, ok := sp.(*ast.TypeSpec); ok && ts.Name.Name == typeName {
					typeDecl = ts
				}
			}
		case *ast.FuncDecl:
			if d.Recv == nil && d.Name.Name == ctor {
				ctorDecl = d
			}
			if d.Recv != nil && d.Name.Name == "CallNative" && recvTypeName(d) == typeName {
				method = d
			}
		}
	}
	if typeDecl == nil || ctorDecl == nil || method == nil {
		t.Fatalf("%s: type %v, constructor %v, CallNative %v; want all three\n%s", ctor, typeDecl != nil, ctorDecl != nil, method != nil, source)
	}
	st, ok := typeDecl.Type.(*ast.StructType)
	if !ok || len(st.Fields.List) == 0 || !isSelector(st.Fields.List[0].Type, "vm", "NativeFn") || len(st.Fields.List[0].Names) != 0 {
		t.Errorf("%s: the struct must embed vm.NativeFn as its first field", typeName)
	}
	allocs := 0
	ast.Inspect(ctorDecl.Body, func(x ast.Node) bool {
		if u, ok := x.(*ast.UnaryExpr); ok {
			if _, lit := u.X.(*ast.CompositeLit); lit {
				allocs++
			}
		}
		if _, ok := x.(*ast.FuncLit); ok {
			t.Errorf("%s: the constructor allocates a Go closure", ctor)
		}
		return true
	})
	if allocs != 1 {
		t.Errorf("%s: %d composite-literal allocations, want exactly 1", ctor, allocs)
	}
}

func recvTypeName(d *ast.FuncDecl) string {
	if len(d.Recv.List) == 0 {
		return ""
	}
	if star, ok := d.Recv.List[0].Type.(*ast.StarExpr); ok {
		if id, ok := star.X.(*ast.Ident); ok {
			return id.Name
		}
	}
	return ""
}
