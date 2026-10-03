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
	"strings"
	"testing"
	"time"
)

// TestCallerModule drives -w-module and `lg compile -module` (#998): both
// commands build inside a Go module the caller prepared, write only the
// directory let-go owns there, and leave the caller's go.mod and go.sum as
// they were. The module reaches let-go through a directory replace, and the
// builds run with GOPROXY=off and GOFLAGS=-mod=mod: the first proves nothing
// was fetched, the second that -mod=readonly still wins over the environment.
func TestCallerModule(t *testing.T) {
	if testing.Short() {
		t.Skip("builds lg, a native Go binary, and a wasm module")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	root := repoRoot(t)
	// -buildvcs=false for the reason TestLgCompileCommand gives.
	bin := filepath.Join(t.TempDir(), "lg")
	if out, err := runCmd(ctx, t, "go", root, []string{"build", "-buildvcs=false", "-o", bin, "."}); err != nil {
		t.Fatalf("build lg: %v\n%s", err, out)
	}

	fix := t.TempDir()
	lib := filepath.Join(fix, "libx", "core.lg")
	if err := os.MkdirAll(filepath.Dir(lib), 0755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, lib, `(ns libx.core)
(defn helper [a] (+ a 41))
`)
	app := filepath.Join(fix, "app.lg")
	writeFile(t, app, `(ns app (:require [libx.core :as lib]))
(defn -main [] (println (lib/helper 1)))
(when-not *compiling-aot* (-main))
`)
	env := []string{
		"LG_SOURCE_PATHS=" + fix + string(os.PathListSeparator) + filepath.Join(root, "pkg", "rt", "core"),
		"GOPROXY=off",
		"GOFLAGS=-mod=mod",
	}

	// newModule requires let-go through a directory replace. A complete one
	// also carries let-go's own requires, as gomod's dev-build module does; an
	// incomplete one is what -mod=mod would quietly rewrite offline.
	letgoMod, err := os.ReadFile(filepath.Join(root, "go.mod"))
	if err != nil {
		t.Fatal(err)
	}
	newModule := func(t *testing.T, complete bool) string {
		t.Helper()
		dir := t.TempDir()
		goMod := "module example.com/host\n\ngo 1.26\n"
		if complete {
			goMod = strings.Replace(string(letgoMod), "module github.com/nooga/let-go", "module example.com/host", 1)
		}
		goMod += "\nrequire github.com/nooga/let-go v0.0.0\n\nreplace github.com/nooga/let-go => " + root + "\n"
		writeFile(t, filepath.Join(dir, "go.mod"), goMod)
		copyFile(t, filepath.Join(root, "go.sum"), filepath.Join(dir, "go.sum"))
		return dir
	}
	readModFiles := func(t *testing.T, dir string) [2][]byte {
		t.Helper()
		var got [2][]byte
		for i, f := range []string{"go.mod", "go.sum"} {
			b, err := os.ReadFile(filepath.Join(dir, f))
			if err != nil {
				t.Fatal(err)
			}
			got[i] = b
		}
		return got
	}

	t.Run("lg compile builds in the module and imports under its path", func(t *testing.T) {
		mod := newModule(t, true)
		before := readModFiles(t, mod)
		exe := filepath.Join(t.TempDir(), "app.native")
		// Twice: the second build replaces the first's generated directory.
		for range 2 {
			if out, err := runCmd(ctx, t, bin, fix,
				[]string{"compile", "-o", exe, "-module", mod, lib, app}, env...); err != nil {
				t.Fatalf("lg compile -module: %v\n%s", err, out)
			}
		}
		stdout, stderr, exit, err := runNativeEntryBinary(ctx, exe)
		if err != nil || exit != 0 {
			t.Fatalf("run: err=%v exit=%d stdout=%q stderr=%q", err, exit, stdout, stderr)
		}
		if want := []byte("42\n"); !bytes.Equal(stdout, want) {
			t.Fatalf("stdout %q, want %q; stderr=%q", stdout, want, stderr)
		}
		if after := readModFiles(t, mod); !bytes.Equal(before[0], after[0]) || !bytes.Equal(before[1], after[1]) {
			t.Error("lg compile -module changed the caller's go.mod or go.sum")
		}
		gen := filepath.Join(mod, "lgprogram")
		if _, err := os.Stat(filepath.Join(gen, ".lg-generated")); err != nil {
			t.Errorf("generated directory has no marker: %v", err)
		}
		// The frame imports the lowered app package by the module's path,
		// not the hardcoded lgprogram module name.
		mainGo, err := os.ReadFile(filepath.Join(gen, "main.go"))
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Contains(mainGo, []byte(`"example.com/host/lgprogram/`)) {
			t.Errorf("main.go does not import under example.com/host/lgprogram:\n%s", mainGo)
		}
		entries, _ := os.ReadDir(mod)
		for _, e := range entries {
			switch e.Name() {
			case "go.mod", "go.sum", "lgprogram":
			default:
				t.Errorf("lg compile -module wrote %s outside lgprogram", e.Name())
			}
		}
	})

	t.Run("-w builds in the module", func(t *testing.T) {
		mod := newModule(t, true)
		before := readModFiles(t, mod)
		outDir := filepath.Join(t.TempDir(), "web")
		script := filepath.Join(fix, "web.lg")
		writeFile(t, script, "(println \"hi\")\n")
		if out, err := runCmd(ctx, t, bin, fix,
			[]string{"-w", outDir, "-w-shell", "none", "-w-module", mod, script}, env...); err != nil {
			t.Fatalf("lg -w -w-module: %v\n%s", err, out)
		}
		if _, err := os.Stat(filepath.Join(outDir, "index.html")); err != nil {
			t.Errorf("no wasm output: %v", err)
		}
		for _, f := range []string{"main.go", "program.lgb", ".lg-generated"} {
			if _, err := os.Stat(filepath.Join(mod, "lgprogram", f)); err != nil {
				t.Errorf("lgprogram is missing %s: %v", f, err)
			}
		}
		if after := readModFiles(t, mod); !bytes.Equal(before[0], after[0]) || !bytes.Equal(before[1], after[1]) {
			t.Error("-w-module changed the caller's go.mod or go.sum")
		}
	})

	t.Run("an unmarked lgprogram is refused", func(t *testing.T) {
		mod := newModule(t, true)
		own := filepath.Join(mod, "lgprogram", "main.go")
		if err := os.Mkdir(filepath.Dir(own), 0755); err != nil {
			t.Fatal(err)
		}
		writeFile(t, own, "package main\n\nfunc main() {}\n")
		for _, args := range [][]string{
			{"compile", "-o", filepath.Join(t.TempDir(), "x"), "-module", mod, lib, app},
			{"-w", filepath.Join(t.TempDir(), "web"), "-w-module", mod, app},
		} {
			out, err := runCmd(ctx, t, bin, fix, args, env...)
			if err == nil || !strings.Contains(out, ".lg-generated") {
				t.Errorf("lg %v: err=%v, want a refusal naming the marker\n%s", args[0], err, out)
			}
		}
		if b, err := os.ReadFile(own); err != nil || !bytes.Contains(b, []byte("func main")) {
			t.Fatalf("the caller's lgprogram/main.go was touched: %v", err)
		}
	})

	t.Run("missing requires are reported, not added", func(t *testing.T) {
		for _, cmd := range []string{"compile", "-w"} {
			mod := newModule(t, false)
			before := readModFiles(t, mod)
			args := []string{"compile", "-o", filepath.Join(t.TempDir(), "x"), "-module", mod, lib, app}
			if cmd == "-w" {
				args = []string{"-w", filepath.Join(t.TempDir(), "web"), "-w-module", mod, app}
			}
			out, err := runCmd(ctx, t, bin, fix, args, env...)
			if ee, ok := err.(*exec.ExitError); !ok || ee.ExitCode() != 1 {
				t.Errorf("lg %s on an untidy module: err=%v, want exit 1\n%s", cmd, err, out)
				continue
			}
			if !strings.Contains(out, "-mod=readonly") {
				t.Errorf("lg %s: error does not say the build was read-only:\n%s", cmd, out)
			}
			if after := readModFiles(t, mod); !bytes.Equal(before[0], after[0]) || !bytes.Equal(before[1], after[1]) {
				t.Errorf("lg %s: a failed build changed the caller's go.mod or go.sum", cmd)
			}
		}
	})

	t.Run("-work and -module are exclusive", func(t *testing.T) {
		cmd := exec.CommandContext(ctx, bin, "compile", "-work", t.TempDir(), "-module", t.TempDir(), app)
		cmd.Dir = fix
		out, err := cmd.CombinedOutput()
		if ee, ok := err.(*exec.ExitError); !ok || ee.ExitCode() != 2 {
			t.Errorf("err=%v, want exit 2\n%s", err, out)
		}
	})
}
