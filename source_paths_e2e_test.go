/*
 * Copyright (c) 2026 let-go contributors; see CONTRIBUTORS.
 * SPDX-License-Identifier: MIT
 */

package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestSourcePathsControlSearchPath: an explicit -source-paths / LG_SOURCE_PATHS
// fully defines the namespace search path — the current directory is no longer
// searched implicitly. A namespace under the given root resolves; one sitting
// only in the cwd does not, unless "." is listed explicitly. A present-but-empty
// env var means "no paths" (symmetric with an empty flag).
//
// Note: `lg -e` exits 0 even when a require fails, so these subtests assert on
// output content (sentinel tokens / the load-failure message), not exit code.
// (buildLG is defined in scope_e2e_test.go, same package.)
func TestSourcePathsControlSearchPath(t *testing.T) {
	bin := buildLG(t)

	libDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(libDir, "mylib.lg"),
		[]byte(`(ns mylib) (def y "MYLIB_OK")`), 0644); err != nil {
		t.Fatal(err)
	}
	cwdDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(cwdDir, "cwdlib.lg"),
		[]byte(`(ns cwdlib) (def x "CWDLIB_OK")`), 0644); err != nil {
		t.Fatal(err)
	}

	sep := string(os.PathListSeparator)
	const loadFail = "unable to load namespace cwdlib"

	// cleanEnv returns the ambient environment with any LG_SOURCE_PATHS removed,
	// so the test controls that variable's presence/absence per subtest and is
	// not affected by an inherited value.
	cleanEnv := func() []string {
		out := os.Environ()[:0:0]
		for _, kv := range os.Environ() {
			if strings.HasPrefix(kv, "LG_SOURCE_PATHS=") {
				continue
			}
			out = append(out, kv)
		}
		return out
	}

	// run executes lg from cwdDir with the given environment and returns the
	// combined output. lg exits 0 on a failed require, so the error is ignored
	// and assertions are made on output content.
	run := func(t *testing.T, env []string, args ...string) string {
		t.Helper()
		cmd := exec.Command(bin, args...)
		cmd.Dir = cwdDir
		cmd.Env = env
		out, _ := cmd.CombinedOutput()
		return string(out)
	}

	t.Run("lib dir resolves", func(t *testing.T) {
		out := run(t, cleanEnv(), "-source-paths", libDir,
			"-e", `(require '[mylib :as m]) (println m/y)`)
		if !strings.Contains(out, "MYLIB_OK") {
			t.Fatalf("expected mylib to resolve from -source-paths, got:\n%s", out)
		}
	})

	t.Run("cwd not searched with explicit source-paths", func(t *testing.T) {
		out := run(t, cleanEnv(), "-source-paths", libDir,
			"-e", `(require '[cwdlib :as c]) (println c/x)`)
		if strings.Contains(out, "CWDLIB_OK") {
			t.Fatalf("cwd should not be searched, but cwdlib resolved:\n%s", out)
		}
		if !strings.Contains(out, loadFail) {
			t.Fatalf("expected %q for cwdlib, got:\n%s", loadFail, out)
		}
	})

	t.Run("explicit dot opts cwd back in", func(t *testing.T) {
		out := run(t, cleanEnv(), "-source-paths", "."+sep+libDir,
			"-e", `(require '[cwdlib :as c]) (println c/x)`)
		if !strings.Contains(out, "CWDLIB_OK") {
			t.Fatalf("expected cwdlib to resolve with '.' on the path, got:\n%s", out)
		}
	})

	t.Run("env var also drops implicit cwd", func(t *testing.T) {
		env := append(cleanEnv(), "LG_SOURCE_PATHS="+libDir)
		out := run(t, env, "-e", `(require '[cwdlib :as c]) (println c/x)`)
		if strings.Contains(out, "CWDLIB_OK") {
			t.Fatalf("env source-paths should not search cwd, but cwdlib resolved:\n%s", out)
		}
		if !strings.Contains(out, loadFail) {
			t.Fatalf("expected %q for cwdlib, got:\n%s", loadFail, out)
		}
	})

	t.Run("present-but-empty env means no paths", func(t *testing.T) {
		env := append(cleanEnv(), "LG_SOURCE_PATHS=")
		out := run(t, env, "-e", `(require '[cwdlib :as c]) (println c/x)`)
		if strings.Contains(out, "CWDLIB_OK") {
			t.Fatalf("empty env should yield no paths, but cwdlib resolved:\n%s", out)
		}
		if !strings.Contains(out, loadFail) {
			t.Fatalf("expected %q for cwdlib, got:\n%s", loadFail, out)
		}
	})
}
