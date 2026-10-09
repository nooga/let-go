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

// TestCallerModule drives -w-module and `lg compile -module`: both
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

	// -import: the module holds a generated interop package for hash/crc32,
	// which registers the crc32 namespace from init(). The lg that compiles the
	// program does not have it, so the program requires it when -main runs; the
	// namespace resolves only if -import linked the package into the build.
	impApp := filepath.Join(fix, "imp.lg")
	writeFile(t, impApp, `(ns imp)
(defn -main []
  (require 'crc32)
  (println ((resolve 'crc32/ChecksumIEEE) (.getBytes "hello"))))
(when-not *compiling-aot* (-main))
`)
	const crcHello = "907060870\n"
	newInteropModule := func(t *testing.T) string {
		t.Helper()
		mod := newModule(t, true)
		if out, err := runCmd(ctx, t, "go", root, []string{"run", "./cmd/lginterop",
			"-packages", "hash/crc32", "-out-pkg", "interop", "-out", filepath.Join(mod, "interop")}); err != nil {
			t.Fatalf("lginterop: %v\n%s", err, out)
		}
		return mod
	}

	t.Run("-import links the package into lg compile and -w builds", func(t *testing.T) {
		mod := newInteropModule(t)
		before := readModFiles(t, mod)

		exe := filepath.Join(t.TempDir(), "imp.native")
		if out, err := runCmd(ctx, t, bin, fix, []string{"compile", "-o", exe, "-module", mod,
			"-import", "example.com/host/interop", impApp}, env...); err != nil {
			t.Fatalf("lg compile -import: %v\n%s", err, out)
		}
		stdout, stderr, exit, err := runNativeEntryBinary(ctx, exe)
		if err != nil || exit != 0 || string(stdout) != crcHello {
			t.Fatalf("native: err=%v exit=%d stdout=%q, want %q; stderr=%q", err, exit, stdout, crcHello, stderr)
		}

		// Without -import the same program must fail, or the check above
		// proves nothing about -import.
		if out, err := runCmd(ctx, t, bin, fix, []string{"compile", "-o", exe, "-module", mod, impApp}, env...); err != nil {
			t.Fatalf("lg compile without -import: %v\n%s", err, out)
		}
		if stdout, _, exit, _ := runNativeEntryBinary(ctx, exe); exit == 0 || string(stdout) == crcHello {
			t.Fatalf("without -import the program still resolved crc32: exit=%d stdout=%q", exit, stdout)
		}

		outDir := filepath.Join(t.TempDir(), "web")
		if out, err := runCmd(ctx, t, bin, fix, []string{"-w", outDir, "-w-shell", "none", "-w-wasm", "external",
			"-w-module", mod, "-import", "example.com/host/interop", impApp}, env...); err != nil {
			t.Fatalf("lg -w -import: %v\n%s", err, out)
		}
		if after := readModFiles(t, mod); !bytes.Equal(before[0], after[0]) || !bytes.Equal(before[1], after[1]) {
			t.Error("-import changed the caller's go.mod or go.sum")
		}
		node, err := exec.LookPath("node")
		if err != nil {
			t.Skip("node not on PATH; cannot run the wasm build")
		}
		goroot, err := runCmd(ctx, t, "go", root, []string{"env", "GOROOT"})
		if err != nil {
			t.Fatalf("go env GOROOT: %v", err)
		}
		// The generated main writes *out* through the _lgOutput JS global, so
		// the driver defines it and runs the module with Go's wasm_exec.js.
		driver := filepath.Join(t.TempDir(), "run.js")
		writeFile(t, driver, `const fs = require("fs");
require(process.argv[2]);
globalThis._lgOutput = (s) => process.stdout.write(s);
const go = new Go();
WebAssembly.instantiate(fs.readFileSync(process.argv[3]), go.importObject)
  .then((r) => go.run(r.instance))
  .catch((e) => { console.error(e); process.exit(1); });
`)
		cmd := exec.CommandContext(ctx, node, driver,
			filepath.Join(strings.TrimSpace(goroot), "lib", "wasm", "wasm_exec.js"),
			filepath.Join(outDir, "main.wasm"))
		var wout, werr bytes.Buffer
		cmd.Stdout, cmd.Stderr = &wout, &werr
		if err := cmd.Run(); err != nil || wout.String() != crcHello {
			t.Fatalf("wasm: err=%v stdout=%q, want %q; stderr=%q", err, wout.String(), crcHello, werr.String())
		}
	})

	t.Run("-import of a package the module does not provide is reported", func(t *testing.T) {
		const missing = "example.com/nowhere/interop"
		for _, cmd := range []string{"compile", "-w"} {
			mod := newModule(t, true)
			before := readModFiles(t, mod)
			args := []string{"compile", "-o", filepath.Join(t.TempDir(), "x"), "-module", mod, "-import", missing, lib, app}
			if cmd == "-w" {
				args = []string{"-w", filepath.Join(t.TempDir(), "web"), "-w-module", mod, "-import", missing, app}
			}
			out, err := runCmd(ctx, t, bin, fix, args, env...)
			if ee, ok := err.(*exec.ExitError); !ok || ee.ExitCode() != 1 || !strings.Contains(out, missing) {
				t.Errorf("lg %s: err=%v, want exit 1 naming %s\n%s", cmd, err, missing, out)
			}
			if after := readModFiles(t, mod); !bytes.Equal(before[0], after[0]) || !bytes.Equal(before[1], after[1]) {
				t.Errorf("lg %s: a failed -import build changed the caller's go.mod or go.sum", cmd)
			}
		}
	})

	t.Run("-import requires a caller module", func(t *testing.T) {
		for _, args := range [][]string{
			{"compile", "-import", "example.com/x", app},
			{"-w", t.TempDir(), "-import", "example.com/x", app},
		} {
			cmd := exec.CommandContext(ctx, bin, args...)
			cmd.Dir = fix
			out, err := cmd.CombinedOutput()
			if ee, ok := err.(*exec.ExitError); !ok || ee.ExitCode() != 2 || !strings.Contains(string(out), "-import requires") {
				t.Errorf("lg %v: err=%v, want exit 2 naming -import\n%s", args[0], err, out)
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

	t.Run("-w-module requires -w", func(t *testing.T) {
		cmd := exec.CommandContext(ctx, bin, "-w-module", t.TempDir(), app)
		cmd.Dir = fix
		out, err := cmd.CombinedOutput()
		if ee, ok := err.(*exec.ExitError); !ok || ee.ExitCode() != 2 || !strings.Contains(string(out), "-w-module requires -w") {
			t.Errorf("err=%v, want exit 2 naming -w\n%s", err, out)
		}
	})
}
