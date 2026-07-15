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

// buildLGRuntime builds the runtime-only binary (cmd/lg-runtime) into a temp
// dir and returns its path.
func buildLGRuntime(t *testing.T) string {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "lg-runtime")
	cmd := exec.Command("go", "build", "-o", bin, "./cmd/lg-runtime")
	cmd.Dir = repoRoot(t)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build lg-runtime: %v\n%s", err, out)
	}
	return bin
}

// TestRuntimeOnlyDepGraph is the mechanical check behind the runtime-only
// claim: cmd/lg-runtime's dependency graph must contain neither the compiler
// nor the resolver. If an import creeps back in (e.g. via pkg/rt), this fails
// before anyone inspects a binary.
func TestRuntimeOnlyDepGraph(t *testing.T) {
	cmd := exec.Command("go", "list", "-deps", "./cmd/lg-runtime")
	cmd.Dir = repoRoot(t)
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("go list -deps ./cmd/lg-runtime: %v", err)
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
			t.Errorf("cmd/lg-runtime dep graph must not contain %s", banned)
		}
	}
}

// TestRuntimeOnly covers the runtime-only binary end to end: it runs
// precompiled .lgb (plain and bundle-format) with the full CLI lifecycle
// (user args, baseline namespaces), runs as the base of a standalone bundle,
// and rejects source rather than compiling it — there is no path to source
// in this build.
func TestRuntimeOnly(t *testing.T) {
	full := buildLG(t)
	runtime := buildLGRuntime(t)

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

	runBin := func(t *testing.T, bin string, args ...string) (int, string) {
		t.Helper()
		out, err := exec.Command(bin, args...).CombinedOutput()
		if err != nil {
			var ee *exec.ExitError
			if !errors.As(err, &ee) {
				t.Fatalf("run %s %v: %v\n%s", bin, args, err, out)
			}
			return ee.ExitCode(), string(out)
		}
		return 0, string(out)
	}
	run := func(t *testing.T, args ...string) (int, string) {
		t.Helper()
		return runBin(t, runtime, args...)
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

	t.Run("publishes command-line args", func(t *testing.T) {
		lgb := compileLGB(t, `(prn *command-line-args*)`)
		code, out := run(t, lgb, "one", "two")
		if code != 0 || !strings.Contains(out, `("one" "two")`) {
			t.Fatalf(`want ("one" "two"), got %d:%s`, code, out)
		}
		// And nil when there are none — same contract as the full lg.
		code, out = run(t, lgb)
		if code != 0 || !strings.Contains(out, "nil") {
			t.Fatalf("want nil for no args, got %d:%s", code, out)
		}
	})

	t.Run("baseline namespace fns are loaded", func(t *testing.T) {
		// str-join lives in let-go.core, which is auto-refer'd but never
		// required — it only works if LoadCore replays the baselines eagerly.
		lgb := compileLGB(t, `(println (str-join "," ["a" "b" "c"]))`)
		code, out := run(t, lgb)
		if code != 0 || !strings.Contains(out, "a,b,c") {
			t.Fatalf("want a,b,c, got %d:\n%s", code, out)
		}
	})

	t.Run("runs as standalone bundle base", func(t *testing.T) {
		dir := t.TempDir()
		lg := filepath.Join(dir, "app.lg")
		app := filepath.Join(dir, "app")
		src := `
(ns app.main (:require [string]))
(println (string/upper-case "bundled") *command-line-args*)`
		if err := os.WriteFile(lg, []byte(src), runtimeOnlyTestFilePerm); err != nil {
			t.Fatal(err)
		}
		if out, err := exec.Command(full, "-b", app, "-bundle-base", runtime, lg).CombinedOutput(); err != nil {
			t.Fatalf("bundle with lg-runtime base: %v\n%s", err, out)
		}
		code, out := runBin(t, app, "x", "y")
		if code != 0 || !strings.Contains(out, "BUNDLED (x y)") {
			t.Fatalf("want BUNDLED (x y), got %d:%s", code, out)
		}
	})

	t.Run("bundle embeds resources", func(t *testing.T) {
		resDir := t.TempDir()
		if err := os.WriteFile(filepath.Join(resDir, "msg.txt"), []byte("embedded-resource"), runtimeOnlyTestFilePerm); err != nil {
			t.Fatal(err)
		}
		dir := t.TempDir()
		lg := filepath.Join(dir, "app.lg")
		app := filepath.Join(dir, "app")
		src := `(println (io/slurp (io/resource "msg.txt")))`
		if err := os.WriteFile(lg, []byte(src), runtimeOnlyTestFilePerm); err != nil {
			t.Fatal(err)
		}
		if out, err := exec.Command(full, "-b", app, "-resource-paths", resDir, "-bundle-base", runtime, lg).CombinedOutput(); err != nil {
			t.Fatalf("bundle with resources: %v\n%s", err, out)
		}
		// Run from a clean cwd with no resource files around — only the
		// embedded copy can satisfy io/resource.
		cmd := exec.Command(app)
		cmd.Dir = t.TempDir()
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("run bundle: %v\n%s", err, out)
		}
		if !strings.Contains(string(out), "embedded-resource") {
			t.Fatalf("want embedded resource contents, got: %q", out)
		}
	})

	t.Run("reads resources from LG_RESOURCE_PATHS", func(t *testing.T) {
		resDir := t.TempDir()
		if err := os.WriteFile(filepath.Join(resDir, "msg.txt"), []byte("env-resource"), runtimeOnlyTestFilePerm); err != nil {
			t.Fatal(err)
		}
		// The resource read is guarded to runtime: -c executes top-level forms
		// at compile time, and the compile step deliberately gets no resource
		// roots, so a pass proves lg-runtime's own env wiring served the read.
		lgb := compileLGB(t, `(when-not *compiling-aot* (println (io/slurp (io/resource "msg.txt"))))`)
		cmd := exec.Command(runtime, lgb)
		cmd.Env = append(os.Environ(), "LG_RESOURCE_PATHS="+resDir)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("run: %v\n%s", err, out)
		}
		if !strings.Contains(string(out), "env-resource") {
			t.Fatalf("want resource contents via LG_RESOURCE_PATHS, got: %q", out)
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

	t.Run("usage on no program", func(t *testing.T) {
		code, out := run(t)
		if code != 2 || !strings.Contains(out, "usage:") {
			t.Fatalf("want usage and exit 2, got %d:\n%s", code, out)
		}
	})
}
