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

// Creating a closure in lowered Go costs one allocation, the same as the one
// Fn the bytecode VM allocates when the fn form is evaluated, and a literal
// that closes over nothing costs none, as on the VM, where it is a constant:
// lowered Go must not allocate more than the bytecode it replaces. Each fixture fn below
// evaluates one fn form of a given shape and reads the value through an
// `identical?` call, so whatever that call costs is measured separately and
// subtracted; what remains is the closure's own cost. A multi-arity literal
// is one object per arm plus the combining rt.MakeNativeMultiArity, also
// measured on its own.
//
// Lowered Go only runs as a built binary, so the fixture is lowered with the
// production path (scripts/lg-compile), completed into its own Go module
// against this checkout, and measured by a Go test inside that module.
func TestLoweredClosureAllocatesOncePerForm(t *testing.T) {
	if testing.Short() {
		t.Skip("lowers and builds a Go module; run without -short")
	}
	ctx := context.Background()
	root := repoRoot(t)
	bin := buildLG(t)

	const module = "closurealloc"
	outDir := t.TempDir()
	fixture := filepath.Join(outDir, "closurealloc.lg")
	source := `(ns closurealloc)

(defn same-local-cap [y]
  (let [f (fn [x] (+ x y))]
    (identical? f f)))

(defn cap-fixed0 [y]
  (let [f (fn [] y)]
    (identical? f f)))

(defn cap-fixed2 [y]
  (let [f (fn [a b] (+ a (+ b y)))]
    (identical? f f)))

(defn cap-variadic [y]
  (let [f (fn [& xs] (cons y xs))]
    (identical? f f)))

(defn cap-mixed [y]
  (let [f (fn [a & xs] (cons a (cons y xs)))]
    (identical? f f)))

(defn cap-multi [y]
  (let [f (fn ([a] (+ a y)) ([a b] (+ a (+ b y))))]
    (identical? f f)))

(defn cap-nested [y]
  (let [f (fn [a] (fn [b] (+ a (+ b y))))]
    (identical? f f)))

(defn free-variadic []
  (let [f (fn [& xs] xs)]
    (identical? f nil)))
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
	probe := `package closurealloc

import (
	"testing"

	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

func TestOneAllocationPerClosureForm(t *testing.T) {
	ec, err := rt.BootCore()
	if err != nil {
		t.Fatal(err)
	}
	y := vm.Int(1)
	// Costs the fixtures pay besides creating the closure, measured on
	// closures that already exist so the creation itself never counts.
	f1 := rt.BoxNativeFn(func(x vm.Value) vm.Value { return rt.AddValue(x, y) })
	f2 := rt.BoxNativeFn(func(a, b vm.Value) vm.Value { return rt.AddValue(a, b) })
	identicalCost := testing.AllocsPerRun(200, func() {
		if _, err := rt.CoreIdentical(f1, f1); err != nil {
			panic(err)
		}
	})
	multiArityCost := testing.AllocsPerRun(200, func() {
		rt.MakeNativeMultiArity([]vm.Value{f1, f2})
	})
	cases := []struct {
		name     string
		fn       func(*vm.ExecContext, vm.Value) (vm.Value, error)
		closures float64
	}{
		{"SameLocalCap", SameLocalCap, 1},
		{"CapFixed0", CapFixed0, 1},
		{"CapFixed2", CapFixed2, 1},
		{"CapVariadic", CapVariadic, 1},
		{"CapMixed", CapMixed, 1},
		{"CapMulti", CapMulti, 2 + multiArityCost},
		{"CapNested", CapNested, 1},
		// Closes over nothing: one shared fn, as the VM's fn constant.
		{"FreeVariadic", func(ec *vm.ExecContext, _ vm.Value) (vm.Value, error) { return FreeVariadic(ec) }, 0},
	}
	for _, tc := range cases {
		got := testing.AllocsPerRun(200, func() {
			if _, err := tc.fn(ec, y); err != nil {
				panic(err)
			}
		})
		if want := tc.closures + identicalCost; got != want {
			t.Errorf("%s allocates %v per call, want %v (%v for the closure forms plus %v for identical?)", tc.name, got, want, tc.closures, identicalCost)
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
