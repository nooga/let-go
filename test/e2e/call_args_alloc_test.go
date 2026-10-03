/*
 * Copyright (c) 2026 let-go contributors; see CONTRIBUTORS.
 * SPDX-License-Identifier: MIT
 */

package e2e

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

// A lowered call through a value or a var passes its arguments the way the
// bytecode VM does: borrowed by the callee for the call, never heap-allocated
// at the site. The bytecode VM lends a callee a window of the caller's operand
// stack; lowered Go used to build a []vm.Value on the heap for every call. So
// a lowered fn whose whole work is one call of one to three arguments, to a
// callee that itself allocates nothing, must allocate nothing per call — for
// a boxed Go func, a bytecode fn, a keyword, and a var bound to a lowered fn
// (whose own body calls (type x) through a var the same way).
//
// Lowered Go only runs as a built binary, so the fixture is lowered with the
// production path (scripts/lg-compile), completed into its own Go module
// against this checkout, and measured by a Go test inside that module.
func TestLoweredCallAllocatesNothingForArgs(t *testing.T) {
	if testing.Short() {
		t.Skip("lowers and builds a Go module; run without -short")
	}
	ctx := context.Background()
	root := repoRoot(t)
	bin := buildLG(t)

	const module = "callalloc"
	outDir := t.TempDir()
	fixture := filepath.Join(outDir, "callalloc.lg")
	source := `(ns callalloc)

(defn call1 [f x] (f x))
(defn call2 [f x y] (f x y))
(defn call3 [f x y z] (f x y z))
(defn keyword-call [m] (:k m))
(defn keyword-call-default [m] (:z m 7))

(defn via-var [x] (keyword? x))
`
	if err := os.WriteFile(fixture, []byte(source), 0o644); err != nil {
		t.Fatal(err)
	}
	if out, err := runCmd(ctx, t, bin, root, []string{"scripts/lg-compile", outDir, module, fixture}); err != nil {
		t.Fatalf("lg-compile: %v\n%s", err, out)
	}

	mod := "module " + module + "\n\ngo " + rootGoDirective(t, root) + "\n\n" +
		"require github.com/nooga/let-go v0.0.0\n" +
		"replace github.com/nooga/let-go => " + root + "\n"
	if err := os.WriteFile(filepath.Join(outDir, "go.mod"), []byte(mod), 0o644); err != nil {
		t.Fatal(err)
	}
	probe := `package callalloc

import (
	"testing"

	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

func TestNoArgAllocationPerLoweredCall(t *testing.T) {
	ec, err := rt.BootCore()
	if err != nil {
		t.Fatal(err)
	}
	x, y, z := vm.Int(1), vm.Int(2), vm.Int(3)
	// Callees that allocate nothing themselves, one per kind lowered code
	// dispatches to dynamically.
	boxed1 := rt.BoxNativeFn(func(a vm.Value) vm.Value { return a })
	boxed2 := rt.BoxNativeFn(func(a, b vm.Value) vm.Value { return b })
	boxed3 := rt.BoxNativeFn(func(a, b, c vm.Value) vm.Value { return c })
	closure2, err := compiler.Eval("(let [k 0] (fn [a b] b))")
	if err != nil {
		t.Fatal(err)
	}
	bytecode2, err := compiler.Eval("(fn [a b] b)")
	if err != nil {
		t.Fatal(err)
	}
	bytecode1, err := compiler.Eval("(fn [a] a)")
	if err != nil {
		t.Fatal(err)
	}
	m := vm.NewPersistentMap([]vm.Value{vm.Keyword("k"), x})
	// keyword? is a lowered clojure.core fn reached through its var (a
	// cross-namespace trampoline), whose own body calls (type x) the same way.
	viaVar := func() (vm.Value, error) { return ViaVar(ec, vm.Keyword("k")) }

	cases := []struct {
		name string
		call func() (vm.Value, error)
		want vm.Value
	}{
		{"boxed Go func, 1 arg", func() (vm.Value, error) { return Call1(ec, boxed1, x) }, x},
		{"boxed Go func, 2 args", func() (vm.Value, error) { return Call2(ec, boxed2, x, y) }, y},
		{"boxed Go func, 3 args", func() (vm.Value, error) { return Call3(ec, boxed3, x, y, z) }, z},
		{"bytecode fn, 1 arg", func() (vm.Value, error) { return Call1(ec, bytecode1, x) }, x},
		{"bytecode fn, 2 args", func() (vm.Value, error) { return Call2(ec, bytecode2, x, y) }, y},
		{"bytecode closure, 2 args", func() (vm.Value, error) { return Call2(ec, closure2, x, y) }, y},
		{"keyword", func() (vm.Value, error) { return KeywordCall(ec, m) }, x},
		{"keyword with default", func() (vm.Value, error) { return KeywordCallDefault(ec, m) }, vm.Int(7)},
		{"var bound to a lowered fn", viaVar, vm.TRUE},
	}
	for _, tc := range cases {
		got, err := tc.call()
		if err != nil || !vm.ValueEquals(got, tc.want) {
			t.Fatalf("%s: got %v, %v; want %v", tc.name, got, err, tc.want)
		}
		// Warm the frame pool once, so a bytecode callee's pooled frame is
		// not counted; the VM pays that once as well.
		_, _ = tc.call()
		if n := testing.AllocsPerRun(200, func() {
			if _, err := tc.call(); err != nil {
				panic(err)
			}
		}); n != 0 {
			t.Errorf("%s: %v allocations per call, want 0", tc.name, n)
		}
	}
}
`
	pkgDir := filepath.Join(outDir, module)
	if err := os.WriteFile(filepath.Join(pkgDir, "alloc_test.go"), []byte(probe), 0o644); err != nil {
		t.Fatal(err)
	}
	if out, err := runCmd(ctx, t, "go", outDir, []string{"mod", "tidy"}); err != nil {
		t.Fatalf("go mod tidy: %v\n%s", err, out)
	}
	if out, err := runCmd(ctx, t, "go", outDir, []string{"test", "./" + module, "-count=1"}); err != nil {
		t.Fatalf("allocation probe failed: %v\n%s", err, out)
	}
}
