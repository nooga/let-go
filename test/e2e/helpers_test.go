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

// lgBuild is one cached build of the lg binary for a tag set.
type lgBuild struct {
	once sync.Once
	bin  string
	err  error
}

var (
	lgDir    string
	lgMu     sync.Mutex
	lgBuilds = map[string]*lgBuild{}
)

// TestMain owns the directory the shared lg binaries live in, so they outlive
// any one test's t.TempDir.
func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "lg-e2e-")
	if err != nil {
		panic(err)
	}
	lgDir = dir
	code := m.Run()
	os.RemoveAll(dir)
	os.Exit(code)
}

// buildLG builds the default lg binary once per package run and returns its
// path. Building per test cost about half a second each, which is most of
// this package's wall time.
func buildLG(t *testing.T) string {
	t.Helper()
	return buildLGTags(t, "")
}

// buildLGTags builds lg with the given build tags ("" = default bytecode;
// "gogen_ir" = generated-core engine variant) once per package run. It does
// not AOT-compile fixtures. The build runs from the repo root so `go build .`
// targets the root main package. repoRoot is resolved before the Once: its
// t.Fatal would otherwise mark the Once done with neither a binary nor an
// error, and every later caller would get a path that was never written.
func buildLGTags(t *testing.T, tags string) string {
	t.Helper()
	root := repoRoot(t)
	lgMu.Lock()
	b := lgBuilds[tags]
	if b == nil {
		b = &lgBuild{}
		lgBuilds[tags] = b
	}
	lgMu.Unlock()
	b.once.Do(func() {
		name, args := "lg", []string{"build"}
		if tags != "" {
			name = "lg-" + tags
			args = append(args, "-tags", tags)
		}
		bin := filepath.Join(lgDir, name)
		cmd := exec.Command("go", append(args, "-o", bin, ".")...)
		cmd.Dir = root
		if out, err := cmd.CombinedOutput(); err != nil {
			b.err = fmt.Errorf("build lg (tags=%q): %v\n%s", tags, err, out)
			return
		}
		b.bin = bin
	})
	if b.err != nil {
		t.Fatal(b.err)
	}
	return b.bin
}
