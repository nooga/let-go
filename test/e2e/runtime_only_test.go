/*
 * Copyright (c) 2026 let-go contributors; see CONTRIBUTORS.
 * SPDX-License-Identifier: MIT
 */

package e2e

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const (
	runtimeOnlyUserRead  os.FileMode = 1 << 8
	runtimeOnlyUserWrite os.FileMode = 1 << 7
	runtimeOnlyGroupRead os.FileMode = 1 << 5
	runtimeOnlyOtherRead os.FileMode = 1 << 2

	runtimeOnlyTestFilePerm = runtimeOnlyUserRead | runtimeOnlyUserWrite | runtimeOnlyGroupRead | runtimeOnlyOtherRead
)

// buildLGRuntimeOnly builds the runtime-only lg binary (-tags runtime_only)
// into a temp dir and returns its path.
func buildLGRuntimeOnly(t *testing.T) string {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "lg-runtime")
	cmd := exec.Command("go", "build", "-tags", "runtime_only", "-o", bin, ".")
	cmd.Dir = repoRoot(t)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build lg (runtime_only): %v\n%s", err, out)
	}
	return bin
}

// TestRuntimeOnlyDepGraph is the mechanical check behind the runtime_only
// trust claim: the tagged build's dependency graph must contain neither the
// compiler nor the resolver. If an import creeps back in (e.g. via pkg/rt),
// this fails before anyone inspects a binary.
func TestRuntimeOnlyDepGraph(t *testing.T) {
	cmd := exec.Command("go", "list", "-deps", "-tags", "runtime_only", ".")
	cmd.Dir = repoRoot(t)
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("go list -deps -tags runtime_only: %v", err)
	}
	deps := string(out)
	if !strings.Contains(deps, "github.com/nooga/let-go/pkg/vm") {
		t.Fatalf("dep graph looks wrong (no pkg/vm); go list output:\n%s", deps)
	}
	for _, banned := range []string{
		"github.com/nooga/let-go/pkg/compiler",
		"github.com/nooga/let-go/pkg/resolver",
	} {
		if strings.Contains(deps, banned) {
			t.Errorf("runtime_only dep graph must not contain %s", banned)
		}
	}
}

// TestRuntimeOnly covers the runtime-only binary end to end: it runs
// precompiled .lgb (plain and bundle-format), and it rejects source rather
// than compiling it — there is no path to source in this build.
func TestRuntimeOnly(t *testing.T) {
	full := buildLG(t)
	runtime := buildLGRuntimeOnly(t)

	// compileLGB compiles src (a .lg body) to a .lgb with the full binary and
	// returns the .lgb path.
	compileLGB := func(t *testing.T, src string) string {
		t.Helper()
		dir := t.TempDir()
		lg := filepath.Join(dir, "app.lg")
		lgb := filepath.Join(dir, "app.lgb")
		if err := os.WriteFile(lg, []byte(src), runtimeOnlyTestFilePerm); err != nil {
			t.Fatal(err)
		}
		if out, err := exec.Command(full, "-c", lgb, lg).CombinedOutput(); err != nil {
			t.Fatalf("compile .lgb: %v\n%s", err, out)
		}
		return lgb
	}

	run := func(t *testing.T, args ...string) (int, string) {
		t.Helper()
		out, err := exec.Command(runtime, args...).CombinedOutput()
		if err != nil {
			var ee *exec.ExitError
			if !errors.As(err, &ee) {
				t.Fatalf("run lg-runtime %v: %v\n%s", args, err, out)
			}
			return ee.ExitCode(), string(out)
		}
		return 0, string(out)
	}

	t.Run("runs precompiled program", func(t *testing.T) {
		lgb := compileLGB(t, `(println (reduce + (range 10)))`)
		code, out := run(t, lgb)
		if code != 0 || !strings.Contains(out, "45") {
			t.Fatalf("want exit 0 with output 45, got %d:\n%s", code, out)
		}
	})

	t.Run("runs bundle with required namespace", func(t *testing.T) {
		lgb := compileLGB(t, `
(ns app.main (:require [string]))
(println (string/capitalize "bytecode"))`)
		code, out := run(t, lgb)
		if code != 0 || !strings.Contains(out, "Bytecode") {
			t.Fatalf("want exit 0 with output Bytecode, got %d:\n%s", code, out)
		}
	})

	t.Run("rejects source", func(t *testing.T) {
		lg := filepath.Join(t.TempDir(), "app.lg")
		if err := os.WriteFile(lg, []byte(`(println :should-never-run)`), runtimeOnlyTestFilePerm); err != nil {
			t.Fatal(err)
		}
		code, out := run(t, lg)
		if code == 0 {
			t.Fatalf("want nonzero exit for .lg source, got 0:\n%s", out)
		}
		if strings.Contains(out, "should-never-run") {
			t.Fatalf("source was executed:\n%s", out)
		}
	})

	t.Run("boots with no program", func(t *testing.T) {
		if code, out := run(t); code != 0 {
			t.Fatalf("want exit 0 for bare boot, got %d:\n%s", code, out)
		}
	})
}
