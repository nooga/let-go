/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package e2e

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/nooga/let-go/internal/gofragment"
)

const (
	// jankRunnerGoFn is the lowered Go function for the adapter's
	// `run-selected-test`; it is what selected_test.go calls and whose
	// structured result the generated test asserts.
	jankRunnerGoFn = "RunSelectedTest"
	// jankGeneratedPackage is the Go package (and module) name of the
	// generated tree. It is generated into t.TempDir(), never into the
	// repository.
	jankGeneratedPackage = "jankdirectabi"
	// jankSemanticWrap is the semantic mutant: the original returned
	// expression is still evaluated (so generated temporaries stay used and
	// the mutant compiles) but the produced value becomes vm.FALSE, which the
	// generated test must reject. semanticWrapHole is defined by the
	// native-entry gate in this package.
	jankSemanticWrap = "func() vm.Value { _ = " + semanticWrapHole + "; return vm.FALSE }()"
)

// TestJankSuiteDirectABIGeneratedGo lowers one pinned jank deftest
// (`clojure.core_test/identical?`) into a real Go test package and requires four
// pieces of evidence: an inline Go AST oracle match, that oracle's falsifier, a
// passing `go test` of the generated package, and the death of an AST-located
// semantic mutant.
//
// The generated package is written into t.TempDir() as its own Go module that
// `replace`s github.com/nooga/let-go with this checkout — exactly like the
// native-entry gate — so no generated code ever lands in the repository tree.
func TestJankSuiteDirectABIGeneratedGo(t *testing.T) {
	if testing.Short() {
		t.Skip("generated jank harness builds and runs a Go test package")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	root := repoRoot(t)
	fixture := filepath.Join(root, "test/clojure-test-suite/test/clojure/core_test/identical_qmark.cljc")
	fixtureSource, err := os.ReadFile(fixture)
	if err != nil {
		// Never a skip: a missing capability must be loud. This test is part
		// of `make native-entry-gate`, and a silent skip would turn the
		// direct-ABI claim it proves into an unmeasured assumption.
		t.Fatalf("pinned jank fixture %s is unreadable: %v\n"+
			"test/clojure-test-suite is a git submodule; jj does not manage submodules and\n"+
			"git worktrees / jj workspaces do not materialize them. Remedies:\n"+
			"  1. in the PRIMARY git worktree: git submodule update --init test/clojure-test-suite\n"+
			"  2. in a jj workspace or secondary worktree, symlink it from the primary worktree:\n"+
			"     scripts/link-clojure-test-suite.sh <workspace-path>\n"+
			"     (equivalently: ln -s <primary>/test/clojure-test-suite <workspace>/test/clojure-test-suite)\n"+
			"jj does not report that symlink as a working-copy change, so it does not pollute commits.",
			fixture, err)
	}
	if !strings.Contains(string(fixtureSource), "(deftest test-identical?") {
		t.Fatalf("pinned identical? fixture no longer contains test-identical?")
	}

	generatedDir := t.TempDir()

	lg := buildLG(t)
	sourcePaths := strings.Join([]string{
		filepath.Join(root, "test/compat"),
		filepath.Join(root, "test/clojure-test-suite/test"),
		root,
	}, string(os.PathListSeparator))
	generate := exec.CommandContext(
		ctx,
		lg,
		"-source-paths", sourcePaths,
		"test/tools/jank-go-harness.lg",
		"--fixture", fixture,
		"--test", "test-identical?",
		"--package", jankGeneratedPackage,
		"--out-dir", generatedDir,
	)
	generate.Dir = root
	if out, err := generate.CombinedOutput(); err != nil {
		t.Fatalf("generate pinned jank Go harness: %v\n%s", err, out)
	}

	generatedFile := filepath.Join(generatedDir, "selected.go")
	generated, err := os.ReadFile(generatedFile)
	if err != nil {
		t.Fatalf("read generated selected Go: %v", err)
	}
	oracle := gofragment.GoMatchRequest{
		Kind:     gofragment.GoExpression,
		Target:   "TestIdenticalQmark",
		Expected: `ec.Invoke(rt.CachedVarFn(&__v_clojure_core_identical_QMARK_, "clojure.core", "identical?"), []vm.Value{x, y})`,
	}
	if err := gofragment.MatchGeneratedGoFragment(oracle, string(generated)); err != nil {
		t.Fatalf("inline Go AST oracle mismatch: %v", err)
	}

	badOracle := oracle
	badOracle.Expected = strings.Replace(oracle.Expected, `"identical?"`, `"not-identical?"`, 1)
	if err := gofragment.MatchGeneratedGoFragment(badOracle, string(generated)); err == nil {
		t.Fatal("structural falsifiability: mismatched inline Go oracle unexpectedly passed")
	}

	prepareJankGeneratedModule(ctx, t, root, generatedDir)
	if out, err := runGeneratedGoTest(ctx, t, generatedDir); err != nil {
		t.Fatalf("generated pinned jank Go test failed: %v\n%s", err, out)
	}

	// Semantic falsifiability. The mutation site is located structurally with
	// go/ast (shared helper, same mechanism as the native-entry gate): the last
	// returned expression of the lowered runner is wrapped so the original
	// expression is still evaluated — keeping the mutant compilable — while the
	// produced value becomes vm.FALSE. No gensym-numbered temporary such as
	// `v78` ever appears in this test's source.
	mutated := wrapLastReturnResult(t, generated, jankRunnerGoFn, jankSemanticWrap)
	if bytes.Equal(mutated, generated) {
		t.Fatalf("semantic mutation of %s produced an identical file", jankRunnerGoFn)
	}
	if err := os.WriteFile(generatedFile, mutated, 0o644); err != nil {
		t.Fatal(err)
	}
	if out, err := runGeneratedGoTest(ctx, t, generatedDir); err == nil {
		t.Fatalf("generated-code semantic mutant unexpectedly passed:\n%s", out)
	}
}

// prepareJankGeneratedModule turns the generated directory into a standalone Go
// module resolving github.com/nooga/let-go to this checkout, so the generated
// package can be built and run without ever being written into the repository.
func prepareJankGeneratedModule(ctx context.Context, t *testing.T, root, dir string) {
	t.Helper()
	mod := "module " + jankGeneratedPackage + "\n\ngo 1.23\n\n" +
		"require github.com/nooga/let-go v0.0.0\n" +
		"replace github.com/nooga/let-go => " + root + "\n"
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(mod), 0o644); err != nil {
		t.Fatal(err)
	}
	if out, err := runCmd(ctx, t, "go", dir, []string{"mod", "tidy"}); err != nil {
		t.Fatalf("go mod tidy in generated module: %v\n%s", err, out)
	}
}

func runGeneratedGoTest(ctx context.Context, t *testing.T, dir string) (string, error) {
	t.Helper()
	return runCmd(ctx, t, "go", dir, []string{"test", ".", "-count=1"})
}
