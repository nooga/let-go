/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package ir_test

// Behavior test for AOT direct cross-PACKAGE calls: compile a two-namespace
// program (lib + app, where app calls lib/greet) to native Go via
// scripts/lg-compile, then `go build` the result. This proves the WHOLE
// contract by construction — the emitted cross-package call qualifies as
// `alias.Fn(ec,…)`, the callee package is imported under that SAME alias, the
// callee is exported, and (because Go's own compiler checks it) there are no
// undefined symbols, unused imports, or alias collisions. Far stronger than
// string-matching the rendered Go, and immune to alias-scheme changes.
//
// Integration test: shells out to `go run . scripts/lg-compile …` and `go
// build`. Skipped under -short.

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestCrossPackageDirectCallCompiles(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test (lg-compile + go build); skipped under -short")
	}

	// repoRoot = two levels up from pkg/ir.
	_, thisFile, _, _ := runtime.Caller(0)
	repoRoot, err := filepath.Abs(filepath.Join(filepath.Dir(thisFile), "..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	src := filepath.Join(repoRoot, "examples", "aot", "cross-package", "src")
	if _, err := os.Stat(filepath.Join(src, "app.lg")); err != nil {
		t.Skipf("cross-package fixtures not present: %v", err)
	}

	work := t.TempDir()
	outDir := filepath.Join(work, "gen")
	const modPath = "crosspkgtest.example/m"

	// 1. Lower lib + app to native Go via lg-compile (same driver users run).
	gen := exec.Command("go", "run", ".", "scripts/lg-compile", outDir, modPath,
		filepath.Join(src, "lib.lg"), filepath.Join(src, "app.lg"))
	gen.Dir = repoRoot
	gen.Env = append(os.Environ(),
		"LG_SOURCE_PATHS="+src+string(os.PathListSeparator)+filepath.Join(repoRoot, "pkg", "rt", "gogen"))
	if out, err := gen.CombinedOutput(); err != nil {
		t.Fatalf("lg-compile failed: %v\n%s", err, out)
	} else if strings.Contains(string(out), "EMIT-FAIL") {
		t.Fatalf("lg-compile reported EMIT-FAIL:\n%s", out)
	}

	// 2. Make the generated tree a module that resolves the let-go imports.
	gomod := "module " + modPath + "\n\ngo 1.22\n\n" +
		"require github.com/nooga/let-go v0.0.0\n\n" +
		"replace github.com/nooga/let-go => " + repoRoot + "\n"
	if err := os.WriteFile(filepath.Join(outDir, "go.mod"), []byte(gomod), 0o644); err != nil {
		t.Fatal(err)
	}
	tidy := exec.Command("go", "mod", "tidy")
	tidy.Dir = outDir
	tidy.Env = os.Environ()
	if out, err := tidy.CombinedOutput(); err != nil {
		t.Fatalf("go mod tidy failed: %v\n%s", err, out)
	}

	// 3. go build — the real behavior assertion. A bare call, unqualified
	//    callee, alias collision, or unused import all fail here.
	build := exec.Command("go", "build", "./...")
	build.Dir = outDir
	build.Env = os.Environ()
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("generated cross-package Go failed to compile:\n%s", out)
	}
}
