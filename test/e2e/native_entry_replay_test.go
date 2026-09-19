/*
 * Copyright (c) 2026 let-go contributors; see CONTRIBUTORS.
 * SPDX-License-Identifier: MIT
 */

package e2e

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestNativeEntryNamespaceReplay pins both halves of the entry-frame contract:
// required namespaces initialize under ordinary runtime semantics, and the
// selected top-level entry invocation executes only in the native frame.
func TestNativeEntryNamespaceReplay(t *testing.T) {
	if testing.Short() {
		t.Skip("builds lg and a native-entry Go binary")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	root := repoRoot(t)
	bin := buildLG(t)
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
	outDir := t.TempDir()
	sourcePaths := "LG_SOURCE_PATHS=" + fix + string(os.PathListSeparator) + filepath.Join(root, "pkg", "rt", "core")
	if out, err := runCmd(ctx, t, bin, root,
		[]string{"scripts/lg-compile", "--entry-frame", outDir, "tmpmod", lib, app}, sourcePaths); err != nil {
		t.Fatalf("lg-compile --entry-frame: %v\n%s", err, out)
	}
	lgb := filepath.Join(outDir, "program.lgb")
	if out, err := runCmd(ctx, t, bin, root,
		[]string{"-c", lgb, "-entry-frame-entry", "app/-main", app}, sourcePaths); err != nil {
		t.Fatalf("lg -c -entry-frame-entry: %v\n%s", err, out)
	}
	mod := "module tmpmod\n\ngo 1.23\n\nrequire github.com/nooga/let-go v0.0.0\nreplace github.com/nooga/let-go => " + root + "\n"
	if err := os.WriteFile(filepath.Join(outDir, "go.mod"), []byte(mod), 0644); err != nil {
		t.Fatal(err)
	}
	if out, err := runCmd(ctx, t, "go", outDir, []string{"mod", "tidy"}); err != nil {
		t.Fatalf("go mod tidy: %v\n%s", err, out)
	}
	exe := filepath.Join(outDir, "app.native")
	if out, err := runCmd(ctx, t, "go", outDir, []string{"build", "-o", exe, "."}); err != nil {
		t.Fatalf("go build: %v\n%s", err, out)
	}
	stdout, stderr, exit, err := runNativeEntryBinary(ctx, exe)
	if err != nil || exit != 0 {
		t.Fatalf("run: err=%v exit=%d stdout=%q stderr=%q", err, exit, stdout, stderr)
	}
	want := []byte("libx: runtime init ran\napp: runtime init ran\n42\n")
	if !bytes.Equal(stdout, want) {
		t.Fatalf("stdout %q, want %q; stderr=%q", stdout, want, stderr)
	}
}
