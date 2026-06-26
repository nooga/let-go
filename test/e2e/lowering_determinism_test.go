/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package e2e

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"testing"
)

// TestLoweringDeterminism verifies that the gogen_ir lowering process is
// deterministic: two runs of lgbgen produce identical byte-for-byte output.
//
// This is a property test that commits to the determinism invariant. The
// lowering encodes inferred types and generated code into Go source files,
// so any non-determinism in token positions, type inference, or code generation
// would produce different bytes across runs. This test guards against regressions.
//
// Skipped under testing.Short() since it runs lgbgen twice (expensive).
// Run explicitly via `make test` for full validation.
func TestLoweringDeterminism(t *testing.T) {
	if testing.Short() {
		t.Skip("lowering determinism harness runs lgbgen twice; run via `make test` or `go test ./test/e2e/`")
	}

	root := repoRoot(t)

	// Build lgbgen once and invoke the binary twice. `go run` would recompile
	// the tool on every call, and that compile dominates the wall clock for a
	// two-run test; building once keeps the harness inside the package test
	// timeout.
	bin := buildLgbgen(t, root)

	// Each run is fully output-isolated (its own --target tree and --code-dir
	// wireup under a private temp dir), so the two share no state and run
	// CONCURRENTLY — roughly halving wall clock vs sequential. Errors are
	// collected and reported from the test goroutine (t.Fatalf is unsafe from
	// a spawned goroutine).
	base := t.TempDir()
	dirs := make([]string, 2)
	errs := make([]error, 2)
	var wg sync.WaitGroup
	for i, label := range []string{"run1", "run2"} {
		wg.Add(1)
		go func(i int, label string) {
			defer wg.Done()
			dirs[i], errs[i] = generateLoweredTree(root, bin, filepath.Join(base, label))
		}(i, label)
	}
	wg.Wait()
	for i, err := range errs {
		if err != nil {
			t.Fatalf("lgbgen run %d: %v", i+1, err)
		}
	}

	// Compare all generated files recursively.
	diffs := compareDirectories(t, dirs[0], dirs[1])
	if len(diffs) > 0 {
		t.Errorf("Lowering is non-deterministic: %d files differ between runs", len(diffs))
		for _, diff := range diffs[:min(len(diffs), 5)] {
			t.Logf("  DIFF: %s", diff)
		}
		if len(diffs) > 5 {
			t.Logf("  ... and %d more", len(diffs)-5)
		}
		t.FailNow()
	}
}

// buildLgbgen compiles the bootstrap lgbgen tool once into a temp path so the
// determinism harness can invoke the binary twice without paying the Go
// compile on each run.
func buildLgbgen(t *testing.T, repoRoot string) string {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "lgbgen")

	cmd := exec.Command("go", "build", "-tags", "bootstrap", "-o", bin, "./cmd/lgbgen")
	cmd.Dir = repoRoot

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		t.Fatalf("build lgbgen: %v\nstderr:\n%s", err, stderr.String())
	}

	return bin
}

// generateLoweredTree runs lgbgen with the lowered tree (--target=go) and the
// gogen_ir wireup files (--code-dir) both directed under runDir, fully isolating
// the run from the real checkout: it never rewrites the tracked
// pkg/rt/core_go_lowered tree OR the wireup files (lg_gogen_ir.go, …) that
// `go test ./...` builds elsewhere (e.g. TestGogenAOTDiff's -tags gogen_ir
// build). The full isolation is also what lets two runs execute concurrently.
//
// Takes no *testing.T because it is called from goroutines, where t.Fatalf is
// unsafe; it returns the lowered-tree dir and an error instead.
func generateLoweredTree(repoRoot, bin, runDir string) (string, error) {
	outDir := filepath.Join(runDir, "tree")
	codeDir := filepath.Join(runDir, "code")

	// cmd.Dir is the repo root so lgbgen can read the .lg sources, but the
	// lowered tree and wireup both go to absolute temp dirs under runDir.
	// --source-paths includes pkg/rt/gogen so the ns->Go-package manglers
	// (gogen/ns->go-pkg) used by cross-package lowering resolve — mirrors the
	// Makefile's LGBGEN-SOURCE-PATHS that every --target=go invocation passes.
	sourcePaths := strings.Join([]string{"pkg/rt/core", "pkg/rt/gogen"}, string(os.PathListSeparator))
	cmd := exec.Command(bin, "--target=go", "--source-paths", sourcePaths, "--code-dir", codeDir, outDir)
	cmd.Dir = repoRoot

	// Capture stderr (timing summary / warnings) for the failure message.
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\nstderr:\n%s", err, stderr.String())
	}

	return outDir, nil
}

// compareDirectories recursively compares two directory trees and returns
// a list of relative paths that differ in content or presence.
func compareDirectories(t *testing.T, dir1, dir2 string) []string {
	t.Helper()
	var diffs []string

	err := filepath.Walk(dir1, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(dir1, path)
		other := filepath.Join(dir2, rel)

		if info.IsDir() {
			return nil // Walk handles recursion
		}

		// Compare file contents.
		b1, err1 := os.ReadFile(path)
		b2, err2 := os.ReadFile(other)

		if err1 != nil || err2 != nil || !bytes.Equal(b1, b2) {
			diffs = append(diffs, rel)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk dir1: %v", err)
	}

	// Also check for files in dir2 that don't exist in dir1.
	err = filepath.Walk(dir2, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(dir2, path)
		other := filepath.Join(dir1, rel)

		if info.IsDir() {
			return nil
		}

		if _, err := os.Stat(other); os.IsNotExist(err) {
			diffs = append(diffs, rel+" (missing in run1)")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk dir2: %v", err)
	}

	sort.Strings(diffs)
	return diffs
}

// min returns the minimum of two integers.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
