/*
 * Copyright (c) 2026 let-go contributors; see CONTRIBUTORS.
 * SPDX-License-Identifier: MIT
 */

package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestLgCompileEntryFrameOptIn pins the #628 ownership contract: frame
// emission is opt-in via --entry-frame. Without the flag, lg-compile writes
// only lowered packages (historical behavior for Gloat and other
// orchestrators). With the flag, exactly one package-main frame is emitted.
// Package lowering is identical in both modes.
func TestLgCompileEntryFrameOptIn(t *testing.T) {
	bin := buildLG(t)
	root := repoRoot(t)

	srcDir := t.TempDir()
	src := filepath.Join(srcDir, "app.lg")
	if err := os.WriteFile(src, []byte("(ns app)\n(defn -main [] 1)\n"), 0644); err != nil {
		t.Fatal(err)
	}

	runCompile := func(t *testing.T, outDir string, entryFrame bool) string {
		t.Helper()
		args := []string{"scripts/lg-compile"}
		if entryFrame {
			args = append(args, "--entry-frame")
		}
		args = append(args, outDir, "tmpmod", src)
		cmd := exec.Command(bin, args...)
		cmd.Dir = root
		cmd.Env = append(os.Environ(), "LG_SOURCE_PATHS=pkg/rt/core")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("lg-compile (entry-frame=%v): %v\n%s", entryFrame, err, out)
		}
		return string(out)
	}

	pkgRel := filepath.Join("app", "app.go")

	t.Run("no flag → no main.go", func(t *testing.T) {
		outDir := t.TempDir()
		log := runCompile(t, outDir, false)
		if _, err := os.Stat(filepath.Join(outDir, "main.go")); !os.IsNotExist(err) {
			t.Fatalf("main.go should be absent without --entry-frame; err=%v", err)
		}
		if _, err := os.Stat(filepath.Join(outDir, pkgRel)); err != nil {
			t.Fatalf("expected lowered package %s: %v", pkgRel, err)
		}
		if !strings.Contains(log, "fns lowered; native entry:") {
			t.Fatalf("expected summary line without frame emission; log:\n%s", log)
		}
		if strings.Contains(log, "native-entry frame ->") {
			t.Fatalf("must not report frame emission without --entry-frame; log:\n%s", log)
		}
	})

	t.Run("flag → exactly one frame", func(t *testing.T) {
		outDir := t.TempDir()
		log := runCompile(t, outDir, true)
		mainPath := filepath.Join(outDir, "main.go")
		body, err := os.ReadFile(mainPath)
		if err != nil {
			t.Fatalf("main.go missing with --entry-frame: %v\n%s", err, log)
		}
		src := string(body)
		if !strings.Contains(src, "package main") {
			t.Fatalf("frame missing package main:\n%s", src)
		}
		if c := strings.Count(src, "func main()"); c != 1 {
			t.Fatalf("want exactly one func main(), got %d", c)
		}
		if !strings.Contains(log, "native-entry frame ->") {
			t.Fatalf("expected frame emission log; got:\n%s", log)
		}
	})

	t.Run("package lowering identical in both modes", func(t *testing.T) {
		outA := t.TempDir()
		outB := t.TempDir()
		runCompile(t, outA, false)
		runCompile(t, outB, true)
		a, err := os.ReadFile(filepath.Join(outA, pkgRel))
		if err != nil {
			t.Fatal(err)
		}
		b, err := os.ReadFile(filepath.Join(outB, pkgRel))
		if err != nil {
			t.Fatal(err)
		}
		if string(a) != string(b) {
			t.Fatalf("lowered package differed with/without --entry-frame\n--- no flag ---\n%s\n--- flag ---\n%s", a, b)
		}
	})
}
