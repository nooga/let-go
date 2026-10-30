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
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
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

var (
	lgOnce sync.Once
	lgBin  string
	lgErr  error
)

// TestMain owns the directory the shared lg binary lives in, so it outlives
// any one test's t.TempDir.
func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "lg-e2e-")
	if err != nil {
		panic(err)
	}
	lgBin = filepath.Join(dir, "lg")
	code := m.Run()
	os.RemoveAll(dir)
	os.Exit(code)
}

// buildLG builds the lg binary once per package run and returns its path.
// Building per test cost about half a second each, which is most of this
// package's wall time. The build runs from the repo root so `go build .`
// targets the root main package.
func buildLG(t *testing.T) string {
	t.Helper()
	lgOnce.Do(func() {
		cmd := exec.Command("go", "build", "-o", lgBin, ".")
		cmd.Dir = repoRoot(t)
		if out, err := cmd.CombinedOutput(); err != nil {
			lgErr = fmt.Errorf("build lg: %v\n%s", err, out)
		}
	})
	if lgErr != nil {
		t.Fatal(lgErr)
	}
	return lgBin
}
