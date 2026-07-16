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
	"testing"
)

// TestLoweringDeterminism verifies that cold and hot lowering are byte-identical
// and that typeinfer never bails on core.
//
// This test compares lgbgen built with -tags bootstrap (cold, interpreted IR)
// against lgbgen built with -tags "bootstrap gogen_ir" (hot, native-optimized IR).
// The lowering encodes inferred types and generated code into Go source files,
// and for self-hosting to hold, both backends must produce identical output.
//
// The test also asserts that typeinfer never hit its work-unit guard
// (*typeinfer-max-drains*) during either pass. A bail would mean partial types
// (degraded, though still deterministic) lowering — the cap exists only as a
// backstop for pathological inputs, never for real core. The formerly
// standalone never-bails-on-core assertion is folded here since both passes
// already lower the full stdlib.
//
// This is the validation gate that Tasks 1 (position densify) and 2 (work-unit
// typeinfer guard) make byte-identical possible. If it fails, residual
// nondeterminism remains and must be fixed before proceeding.
//
// Skipped under testing.Short(); runs by default otherwise.
func TestLoweringDeterminism(t *testing.T) {
	if testing.Short() {
		t.Skip("lowering determinism harness runs lgbgen twice; run via `go test ./test/e2e/`")
	}

	root := repoRoot(t)

	cold := buildLgbgenTags(t, root, "bootstrap")

	// Materialize the gogen_ir wireup + native lowered tree in the checkout
	// before building hot (same artifacts `make lowered` produces; they are
	// gitignored and regenerated on demand). In a clean checkout they don't
	// exist, `-tags gogen_ir` selects no additional files, and "hot" compiles
	// identically to "cold" — the gate would compare cold against itself.
	materializeNativeLowering(t, root, cold)

	hot := buildLgbgenTags(t, root, "bootstrap gogen_ir")
	assertNativeLoweringLinked(t, root)

	base := t.TempDir()
	// The work-unit typeinfer guard (*typeinfer-max-drains*) is a defensive
	// backstop that must NEVER fire on real core — a bail would silently degrade
	// inferred types. Both passes run the full pipeline over the whole stdlib, so
	// asserting the bail line is absent here folds the former standalone
	// never-bails-on-core test onto lowerings the gate already performs.
	const bailMsg = "typeinfer: hit *typeinfer-max-drains*"

	coldDir, coldOut, err := generateLoweredTree(root, cold, filepath.Join(base, "cold"))
	if err != nil {
		t.Fatalf("cold lgbgen: %v", err)
	}
	if strings.Contains(coldOut, bailMsg) {
		t.Fatalf("typeinfer bailed during COLD lowering — *typeinfer-max-drains* cap hit on real core (should never happen). Raise the cap or investigate. saw: %q", bailMsg)
	}
	hotDir, hotOut, err := generateLoweredTree(root, hot, filepath.Join(base, "hot"))
	if err != nil {
		t.Fatalf("hot lgbgen: %v", err)
	}
	if strings.Contains(hotOut, bailMsg) {
		t.Fatalf("typeinfer bailed during HOT lowering — *typeinfer-max-drains* cap hit on real core. saw: %q", bailMsg)
	}

	diffs := compareDirectories(t, coldDir, hotDir)
	if len(diffs) > 0 {
		t.Errorf("cold vs hot lowering not byte-identical: %d files differ", len(diffs))
		for _, d := range diffs[:min(len(diffs), 5)] {
			t.Logf("  DIFF: %s", d)
		}
		if len(diffs) > 5 {
			t.Logf("  ... and %d more", len(diffs)-5)
		}
		t.FailNow()
	}
}

// buildLgbgenTags compiles an lgbgen binary with the specified tags into a
// temp path. Tags are space-separated (e.g., "bootstrap" or "bootstrap gogen_ir").
func buildLgbgenTags(t *testing.T, repoRoot, tags string) string {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "lgbgen")
	cmd := exec.Command("go", "build", "-tags", tags, "-o", bin, "./cmd/lgbgen")
	cmd.Dir = repoRoot
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("build lgbgen (-tags %q): %v\nstderr:\n%s", tags, err, stderr.String())
	}
	return bin
}

// materializeNativeLowering runs the cold lgbgen against the checkout's
// DEFAULT output locations (equivalent to `make lowered`): it writes the
// native lowered tree and the //go:build gogen_ir wireup files that the hot
// build needs to actually differ from cold. Both are gitignored on-demand
// artifacts, so this mutates no tracked state.
func materializeNativeLowering(t *testing.T, repoRoot, coldBin string) {
	t.Helper()
	cmd := exec.Command(coldBin, "--target=go")
	cmd.Dir = repoRoot
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("materialize native lowering (--target=go): %v\nstderr:\n%s", err, stderr.String())
	}
	// Same bail assertion as the cold/hot passes: a bail here would embed
	// degraded native passes in the hot binary. The cold pass would bail
	// identically and catch it, but fail at the source with the diagnostics.
	if strings.Contains(stdout.String(), "typeinfer: hit *typeinfer-max-drains*") {
		t.Fatalf("typeinfer bailed while materializing the native lowering:\n%s", stdout.String())
	}
}

// assertNativeLoweringLinked fails unless the gogen_ir-tagged lgbgen build
// actually depends on the native lowered tree. Without this assertion a
// missing wireup makes hot == cold and the byte-identity comparison vacuous.
func assertNativeLoweringLinked(t *testing.T, repoRoot string) {
	t.Helper()
	cmd := exec.Command("go", "list", "-deps", "-tags", "bootstrap gogen_ir", "./cmd/lgbgen")
	cmd.Dir = repoRoot
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("go list -deps (hot tags): %v\nstderr:\n%s", err, stderr.String())
	}
	if !strings.Contains(stdout.String(), "core_go_lowered") {
		t.Fatalf("hot lgbgen build does not link the native lowered tree " +
			"(no core_go_lowered package in go list -deps): the gogen_ir wireup " +
			"is missing, so the hot binary would be identical to cold and this " +
			"gate would pass without exercising the native self-hosting path")
	}
}

// generateLoweredTree runs lgbgen with the lowered tree (--target=go) and the
// gogen_ir wireup files (--code-dir) both directed under runDir, fully isolating
// the run from the real checkout: it never rewrites the tracked
// pkg/rt/core_go_lowered tree OR the wireup files (lg_gogen_ir.go, …) that
// `go test ./...` builds elsewhere (e.g. TestGogenAOTDiff's -tags gogen_ir
// build). The full isolation is also what lets two runs execute independently.
//
// Returns the lowered-tree dir, stdout (typeinfer diagnostics), and an error.
func generateLoweredTree(repoRoot, bin, runDir string) (string, string, error) {
	outDir := filepath.Join(runDir, "tree")
	codeDir := filepath.Join(runDir, "code")
	// cmd.Dir is the repo root so lgbgen can read the .lg sources, but the
	// lowered tree and wireup both go to absolute temp dirs under runDir.
	cmd := exec.Command(bin, "--target=go", "--code-dir", codeDir, outDir)
	cmd.Dir = repoRoot
	// Capture stdout (typeinfer diagnostics — e.g. the drain-cap bail line) and
	// stderr (timing summary / warnings). stdout is small diagnostic text; the
	// lowered tree itself is written to files under outDir, not to stdout.
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", "", fmt.Errorf("%v\nstderr:\n%s", err, stderr.String())
	}
	return outDir, stdout.String(), nil
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
