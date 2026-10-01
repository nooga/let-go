/*
 * Copyright (c) 2026 let-go contributors; see CONTRIBUTORS.
 * SPDX-License-Identifier: MIT
 */

package e2e

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

// TestLgCompileCommand drives `lg compile` (#596) end to end: lowering, the
// program.lgb child process, go.mod scaffolding, and go build, all from one
// invocation. The fixture is TestNativeEntryNamespaceReplay's, which runs the
// same pipeline step by step, so the two tests pin the same observable output.
func TestLgCompileCommand(t *testing.T) {
	if testing.Short() {
		t.Skip("builds lg and a native Go binary")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	root := repoRoot(t)
	// -buildvcs=false keeps this a dev build. From a clean checkout, which is
	// what CI has, build info would stamp a pseudo-version for the commit, and
	// the generated go.mod would try to fetch that commit from the module
	// proxy instead of using LETGO_SRC below.
	bin := filepath.Join(t.TempDir(), "lg")
	if out, err := runCmd(ctx, t, "go", root, []string{"build", "-buildvcs=false", "-o", bin, "."}); err != nil {
		t.Fatalf("build lg: %v\n%s", err, out)
	}
	fix := t.TempDir()
	lib := filepath.Join(fix, "libx", "core.lg")
	if err := os.MkdirAll(filepath.Dir(lib), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(lib, []byte(`(ns libx.core)
(defn helper [a] (+ a 41))
(when-not *compiling-aot* (println "libx: runtime init ran"))
`), 0644); err != nil {
		t.Fatal(err)
	}
	app := filepath.Join(fix, "app.lg")
	if err := os.WriteFile(app, []byte(`(ns app (:require [libx.core :as lib]))
(defn -main [] (println (lib/helper 1)))
(when-not *compiling-aot*
  (println "app: runtime init ran")
  (-main))
`), 0644); err != nil {
		t.Fatal(err)
	}
	// LETGO_SRC points the dev build's generated module at this checkout.
	env := []string{
		"LG_SOURCE_PATHS=" + fix + string(os.PathListSeparator) + filepath.Join(root, "pkg", "rt", "core"),
		"LETGO_SRC=" + root,
	}

	t.Run("builds and runs", func(t *testing.T) {
		work := t.TempDir()
		exe := filepath.Join(t.TempDir(), "app.native")
		if out, err := runCmd(ctx, t, bin, fix,
			[]string{"compile", "-o", exe, "-work", work, lib, app}, env...); err != nil {
			t.Fatalf("lg compile: %v\n%s", err, out)
		}
		stdout, stderr, exit, err := runNativeEntryBinary(ctx, exe)
		if err != nil || exit != 0 {
			t.Fatalf("run: err=%v exit=%d stdout=%q stderr=%q", err, exit, stdout, stderr)
		}
		want := []byte("libx: runtime init ran\napp: runtime init ran\n42\n")
		if !bytes.Equal(stdout, want) {
			t.Fatalf("stdout %q, want %q; stderr=%q", stdout, want, stderr)
		}
		// -work keeps the generated module where the caller asked for it.
		for _, f := range []string{"main.go", "program.lgb", "go.mod"} {
			if _, err := os.Stat(filepath.Join(work, f)); err != nil {
				t.Errorf("-work dir is missing %s: %v", f, err)
			}
		}
	})

	t.Run("usage errors exit 2", func(t *testing.T) {
		for _, args := range [][]string{{"compile"}, {"compile", "-bogus", app}, {"compile", "-o"}} {
			cmd := exec.CommandContext(ctx, bin, args...)
			cmd.Dir = fix
			out, err := cmd.CombinedOutput()
			if ee, ok := err.(*exec.ExitError); !ok || ee.ExitCode() != 2 {
				t.Errorf("lg %v: err=%v, want exit 2\n%s", args, err, out)
			}
		}
	})

	t.Run("a default output name that is a directory fails", func(t *testing.T) {
		// app.lg defaults to ./app; go build -o would write app/lgprogram.
		dir := filepath.Join(fix, "app")
		if err := os.Mkdir(dir, 0755); err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(dir)
		out, err := runCmd(ctx, t, bin, fix, []string{"compile", lib, app}, env...)
		if ee, ok := err.(*exec.ExitError); !ok || ee.ExitCode() != 1 {
			t.Fatalf("lg compile into a directory: err=%v, want exit 1\n%s", err, out)
		}
		if !bytes.Contains([]byte(out), []byte("is a directory")) {
			t.Fatalf("missing the directory diagnostic:\n%s", out)
		}
		if entries, _ := os.ReadDir(dir); len(entries) != 0 {
			t.Fatalf("lg compile wrote into %s: %v", dir, entries)
		}
	})

	t.Run("a program without an entry fails", func(t *testing.T) {
		out, err := runCmd(ctx, t, bin, fix, []string{"compile", lib}, env...)
		if ee, ok := err.(*exec.ExitError); !ok || ee.ExitCode() != 1 {
			t.Fatalf("lg compile on an entry-less file: err=%v, want exit 1\n%s", err, out)
		}
		if !bytes.Contains([]byte(out), []byte("no program entry")) {
			t.Fatalf("missing the no-entry diagnostic:\n%s", out)
		}
	})
}
