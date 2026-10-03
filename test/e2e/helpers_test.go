/*
 * Copyright (c) 2026 let-go contributors; see CONTRIBUTORS.
 * SPDX-License-Identifier: MIT
 */

// Package e2e holds black-box end-to-end tests that build the lg binary and run
// it as a subprocess. They live here (not at the repo root) so the root stays
// free of test files. Because `go test` runs a package's tests with the working
// directory set to the package's own directory, these tests cannot rely on the
// repo root being the cwd — repoRoot resolves it explicitly, and buildLG / any
// run that touches a repo-relative path (scripts/, test/gold/, …) sets the
// subprocess cwd to it.
package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// repoRoot returns the repository root by walking up from this source file
// until it finds go.mod. Stable regardless of the test's working directory.
func repoRoot(t *testing.T) string {
	t.Helper()
	// go test runs a test in its package directory; the source path from
	// runtime.Caller is not usable, since -trimpath rewrites it.
	start, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	dir := start
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("repo root (go.mod) not found above " + start)
		}
		dir = parent
	}
}

// buildLG builds the lg binary once into a temp dir and returns its path. The
// build runs from the repo root so `go build .` targets the root main package.
func buildLG(t *testing.T) string {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "lg")
	cmd := exec.Command("go", "build", "-o", bin, ".")
	cmd.Dir = repoRoot(t)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build lg: %v\n%s", err, out)
	}
	return bin
}
