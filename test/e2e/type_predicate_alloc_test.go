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

// vector? compares type objects, so answering it for any value costs no
// allocation in lowered Go.
func TestLoweredVectorPredicateAllocatesNothing(t *testing.T) {
	if testing.Short() {
		t.Skip("lowers and builds a Go module; run without -short")
	}
	ctx := context.Background()
	root := repoRoot(t)
	bin := buildLG(t)

	const module = "predalloc"
	outDir := t.TempDir()
	fixture := filepath.Join(outDir, "predalloc.lg")
	source := `(ns predalloc)

(defn is-vec [x] (vector? x))
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
	probe := `package predalloc

import (
	"testing"

	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

func TestNoAllocationPerVectorPredicate(t *testing.T) {
	ec, err := rt.BootCore()
	if err != nil {
		t.Fatal(err)
	}
	vals := []struct {
		name string
		v    vm.Value
		want vm.Value
	}{
		{"vector", vm.NewArrayVector([]vm.Value{vm.Int(1)}), vm.TRUE},
		{"empty vector", vm.EmptyVector, vm.TRUE},
		{"nil", vm.NIL, vm.FALSE},
		{"keyword", vm.Keyword("k"), vm.FALSE},
		{"string", vm.String("s"), vm.FALSE},
	}
	for _, tc := range vals {
		got, err := IsVec(ec, tc.v)
		if err != nil || got != tc.want {
			t.Fatalf("%s: got %v, %v; want %v", tc.name, got, err, tc.want)
		}
		if n := testing.AllocsPerRun(200, func() {
			if _, err := IsVec(ec, tc.v); err != nil {
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
