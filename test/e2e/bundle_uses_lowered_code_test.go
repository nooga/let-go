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

// An lg built with a lowered package linked in runs that package's Go code,
// and so does an `lg -b` bundle it makes (#991): the bundle decodes its var
// references before the lowered package's overrides are applied, so an
// override must take the existing Var, or the bundle keeps calling bytecode.
func TestBundleUsesLinkedLoweredCode(t *testing.T) {
	if testing.Short() {
		t.Skip("lowers a namespace and builds a custom lg; run without -short")
	}
	root := repoRoot(t)
	lg := buildLG(t)
	tmp := t.TempDir()
	src := filepath.Join(tmp, "src")
	if err := os.MkdirAll(src, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(src, "lib.lg"), "(ns lib)\n(defn g [x] (str x \"!\"))\n")
	app := filepath.Join(tmp, "app.lg")
	writeFile(t, app, "(ns app (:require [lib]))\n(println (str lib/g) (lib/g 1))\n")

	// The lowered package, and an lg that links it: pkg/cli plus a blank
	// import, as docs/guide/custom-lg.md documents.
	host := filepath.Join(tmp, "host")
	compile := exec.Command(lg, "scripts/lg-compile", host, "aothost", filepath.Join(src, "lib.lg"))
	compile.Dir = root
	if out, err := compile.CombinedOutput(); err != nil {
		t.Fatalf("lg-compile: %v\n%s", err, out)
	}
	writeFile(t, filepath.Join(host, "go.mod"), "module aothost\n\ngo "+rootGoDirective(t, root)+
		"\n\nrequire github.com/nooga/let-go v0.0.0\n\nreplace github.com/nooga/let-go => "+root+"\n")
	writeFile(t, filepath.Join(host, "main.go"), `package main

import (
	"os"

	_ "aothost/lib"
	"github.com/nooga/let-go/pkg/cli"
)

func main() { os.Exit(cli.Main("host", "none")) }
`)
	copyFile(t, filepath.Join(root, "go.sum"), filepath.Join(host, "go.sum"))
	customLG := filepath.Join(host, "customlg")
	build := exec.Command("go", "build", "-o", customLG, ".")
	build.Dir = host
	build.Env = append(os.Environ(), "GOFLAGS=-mod=mod")
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("custom lg does not build: %v\n%s", err, out)
	}

	env := append(os.Environ(), "LG_SOURCE_PATHS="+src)
	const want = "<native-fn g"
	run := func(name string, args ...string) string {
		t.Helper()
		cmd := exec.Command(name, args...)
		cmd.Dir = tmp
		cmd.Env = env
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("%s %v: %v\n%s", name, args, err, out)
		}
		return string(out)
	}
	if out := run(customLG, app); !strings.Contains(out, want) || !strings.Contains(out, "1!") {
		t.Fatalf("custom lg runs %q, want lib/g's lowered code (%s …) returning 1!", out, want)
	}
	bundled := filepath.Join(tmp, "app.bin")
	run(customLG, "-b", bundled, app)
	if out := run(bundled); !strings.Contains(out, want) || !strings.Contains(out, "1!") {
		t.Fatalf("bundle runs %q, want lib/g's lowered code (%s …) returning 1!", out, want)
	}
}
