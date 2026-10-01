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

// A constant collection literal is one shared value, so evaluating it in
// lowered Go allocates nothing: the empty vector is vm.EmptyVector, and any
// other constant literal is a package-level var. (type []) is the shape the
// core type predicates evaluate on every call.
func TestLoweredConstantLiteralAllocatesNothing(t *testing.T) {
	if testing.Short() {
		t.Skip("lowers and builds a Go module; run without -short")
	}
	ctx := context.Background()
	root := repoRoot(t)
	bin := buildLG(t)

	const module = "constalloc"
	outDir := t.TempDir()
	fixture := filepath.Join(outDir, "constalloc.lg")
	source := `(ns constalloc)

(defn type-of-empty [] (type []))
(defn mk-vec [] [1 2 3])
(defn mk-map [] {:a 1 :b 2})
(defn mk-set [] #{:a :b})
(defn mk-nested [] [1 {:k [2 #{3}]}])
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
	probe := `package constalloc

import (
	"testing"

	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

func TestNoAllocationPerConstantLiteral(t *testing.T) {
	ec, err := rt.BootCore()
	if err != nil {
		t.Fatal(err)
	}
	// A fn returning a constant lowers without an error result; (type [])
	// calls through a var and keeps one.
	cases := []struct {
		name string
		call func() (vm.Value, error)
	}{
		{"(type [])", func() (vm.Value, error) { return TypeOfEmpty(ec) }},
		{"vector literal", func() (vm.Value, error) { return MkVec(ec), nil }},
		{"map literal", func() (vm.Value, error) { return MkMap(ec), nil }},
		{"set literal", func() (vm.Value, error) { return MkSet(ec), nil }},
		{"nested literal", func() (vm.Value, error) { return MkNested(ec), nil }},
	}
	for _, tc := range cases {
		if _, err := tc.call(); err != nil {
			t.Fatal(err)
		}
		if n := testing.AllocsPerRun(200, func() {
			if _, err := tc.call(); err != nil {
				panic(err)
			}
		}); n != 0 {
			t.Errorf("%s: %v allocations per call, want 0", tc.name, n)
		}
		// Identity across evaluations is compared on the .lg side, where
		// identical? handles every value kind (test/gogen/const_identity.lg).
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
