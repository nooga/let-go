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

// A capturing closure costs one allocation site in lowered Go, no matter how
// many times the enclosing fn reads it. The bytecode VM allocates one Fn when
// the fn form is evaluated; lowered Go must not allocate more than the
// bytecode it replaces, so a fn that binds a closure once and reads it twice
// must allocate exactly what one rt.BoxNativeFn costs.
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

func TestOneBoxPerClosureForm(t *testing.T) {
	ec, err := rt.BootCore()
	if err != nil {
		t.Fatal(err)
	}
	y := vm.Int(1)
	// The hand-written shape of the source: the closure is boxed once and the
	// one value is read twice. Whatever a box and the identical? call cost,
	// the lowered fn may cost no more.
	reference := testing.AllocsPerRun(200, func() {
		f := rt.BoxNativeFn(func(x vm.Value) vm.Value { return rt.AddValue(x, y) })
		if _, err := rt.CoreIdentical(f, f); err != nil {
			panic(err)
		}
	})
	got := testing.AllocsPerRun(200, func() {
		if _, err := SameLocalCap(ec, y); err != nil {
			panic(err)
		}
	})
	if got != reference {
		t.Fatalf("SameLocalCap allocates %v per call, want %v (one rt.BoxNativeFn plus the identical? call)", got, reference)
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
