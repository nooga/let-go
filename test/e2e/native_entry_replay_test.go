/*
 * Copyright (c) 2026 let-go contributors; see CONTRIBUTORS.
 * SPDX-License-Identifier: MIT
 */

package e2e

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestNativeEntryNamespaceReplay pins both halves of the entry-frame contract:
// required namespaces initialize under ordinary runtime semantics, and the
// selected entry runs once. The guarded call is omitted from the bundle and
// the native frame makes it; a call the compiler can't omit runs during load,
// and the frame then skips its own call (#1072).
func TestNativeEntryNamespaceReplay(t *testing.T) {
	if testing.Short() {
		t.Skip("builds lg and a native-entry Go binary")
	}
	for _, c := range []struct{ name, call, want string }{
		{"guarded", "(when-not *compiling-aot*\n  (println \"app: runtime init ran\")\n  (-main))",
			"libx: runtime init ran\napp: runtime init ran\n42\n"},
		{"under when", "(when true (-main))",
			"libx: runtime init ran\n42\n"},
		{"through an alias", "(let [f -main] (mapv (fn [_] (f)) [1]))",
			"libx: runtime init ran\n42\n"},
	} {
		t.Run(c.name, func(t *testing.T) {
			stdout, stderr := runNativeEntryApp(t, "(ns app (:require [libx.core :as lib]))\n(defn -main [] (println (lib/helper 1)))\n"+c.call+"\n")
			if string(stdout) != c.want {
				t.Fatalf("stdout %q, want %q; stderr=%q", stdout, c.want, stderr)
			}
		})
	}
}

// runNativeEntryApp builds app source, which requires libx.core, into a
// native-entry binary and returns what the binary prints.
func runNativeEntryApp(t *testing.T, appSrc string) (stdout, stderr []byte) {
	t.Helper()
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
	if err := os.WriteFile(app, []byte(appSrc), 0644); err != nil {
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
	return stdout, stderr
}
