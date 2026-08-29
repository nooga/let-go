/*
 * Copyright (c) 2026 let-go contributors
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

// TestCustomMain builds the thing this whole change exists for: a third-party
// binary that IS lg — full CLI via pkg/cli — wrapped around interop generated
// for a package let-go has never heard of.
//
// It has to be the real lg and not an lg-runtime-style bytecode host, because
// the custom binary serves two roles at once: it runs scripts (so the native
// namespace must resolve), and it is the compiling side of `lg -b` (top-level
// forms execute at AOT time, so the natives must be linked THERE too). Leg (b)
// covers that second role by using the custom binary as its own bundle base.
//
// Hermetic: a replace directive points at this worktree and the target is a
// stdlib package, so nothing is fetched.
func TestCustomMain(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping custom-main e2e in -short mode (builds a Go binary)")
	}
	root := repoRoot(t)
	// Everything this test exercises runs in child processes, so the outer
	// test cache cannot see it: an edit to the emitter or to pkg/cli would
	// otherwise return a stale cached pass. Read the inputs here to key the
	// cache on their content — same blind spot trackLginteropInputs covers.
	trackLginteropInputs(t, root)
	trackInputs(t,
		filepath.Join(root, "pkg", "cli", "*.go"),
		filepath.Join(root, "lg.go"))
	tmp := t.TempDir()

	// 1. Generate out-of-tree interop for a stdlib package.
	gen := exec.Command("go", "run", "./cmd/lginterop",
		"-packages", "hash/crc32", "-out-pkg", "interop",
		"-out", filepath.Join(tmp, "interop"))
	gen.Dir = root
	if out, err := gen.CombinedOutput(); err != nil {
		t.Fatalf("lginterop failed: %v\n%s", err, out)
	}

	// 2. A custom main: pkg/cli plus a blank import of the generated package.
	//    This is verbatim the pattern docs/guide/custom-lg.md documents.
	writeFile(t, filepath.Join(tmp, "go.mod"), `module example.com/customlg

go 1.26

require github.com/nooga/let-go v0.0.0

replace github.com/nooga/let-go => `+root+`
`)
	writeFile(t, filepath.Join(tmp, "main.go"), `package main

import (
	"os"

	_ "example.com/customlg/interop"
	"github.com/nooga/let-go/pkg/cli"
)

func main() { os.Exit(cli.Main("host-v9", "none")) }
`)
	copyFile(t, filepath.Join(root, "go.sum"), filepath.Join(tmp, "go.sum"))

	customLG := filepath.Join(tmp, "customlg")
	build := exec.Command("go", "build", "-o", customLG, ".")
	build.Dir = tmp
	build.Env = append(os.Environ(), "GOFLAGS=-mod=mod")
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("custom main does not build: %v\n%s", err, out)
	}

	// crc32.ChecksumIEEE("hello") — the value is fixed, so a wrong result is a
	// real failure rather than a tautology against the same code path.
	const script = `(require '[crc32])
(println (crc32/ChecksumIEEE (.getBytes "hello")))
`
	scriptPath := filepath.Join(tmp, "script.lg")
	writeFile(t, scriptPath, script)
	const want = "907060870"

	t.Run("runs a script against the generated namespace", func(t *testing.T) {
		cmd := exec.Command(customLG, scriptPath)
		cmd.Dir = tmp
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("custom lg failed to run script: %v\n%s", err, out)
		}
		if !strings.Contains(string(out), want) {
			t.Errorf("got %q, want it to contain %q", out, want)
		}
	})

	// The host stamp ("host-v9") is CLI metadata: -v may show it, but the
	// runtime's let-go.version must describe let-go itself — here resolved
	// through the directory replace, so "dev". A custom host leaking its own
	// version into System/getProperty would break runtime feature checks.
	t.Run("host metadata does not overwrite the runtime identity", func(t *testing.T) {
		verScript := filepath.Join(tmp, "version.lg")
		writeFile(t, verScript, `(println (System/getProperty "let-go.version"))`+"\n")
		cmd := exec.Command(customLG, verScript)
		cmd.Dir = tmp
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("custom lg failed to run version script: %v\n%s", err, out)
		}
		if strings.Contains(string(out), "host-v9") {
			t.Errorf("let-go.version leaked the host stamp: %q", out)
		}
		if !strings.Contains(string(out), "dev") {
			t.Errorf("let-go.version = %q, want the replace-resolved \"dev\"", out)
		}
	})

	// -w must resolve let-go through the local `replace`, not the proxy. The
	// replacement path lives in the custom binary's build info; without
	// carrying it into the generated module, gomod falls back to @latest,
	// which fails offline and silently builds a different let-go online.
	// GOPROXY=off is the assertion: any proxy access fails the test.
	t.Run("-w resolves let-go through the local replace, offline", func(t *testing.T) {
		if os.Getenv("LETGO_SRC") != "" {
			t.Skip("LETGO_SRC set: it would supply the source dir independently")
		}
		outDir := filepath.Join(tmp, "wasmout")
		cmd := exec.Command(customLG, "-w", outDir, "-w-shell", "none", scriptPath)
		cmd.Dir = tmp
		cmd.Env = append(os.Environ(), "GOPROXY=off", "GOFLAGS=-mod=mod")
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("custom lg -w failed with GOPROXY=off: %v\n%s", err, out)
		}
		if _, err := os.Stat(filepath.Join(outDir, "index.html")); err != nil {
			t.Errorf("no wasm output produced: %v", err)
		}
	})

	// The custom binary is its own bundle base: -b copies the running
	// executable and appends the bytecode, so the standalone binary inherits
	// the generated natives. A stock lg as base would produce a binary that
	// cannot resolve crc32.
	t.Run("bundles with itself as bundle base", func(t *testing.T) {
		bundled := filepath.Join(tmp, "bundled")
		bundle := exec.Command(customLG, "-b", bundled, scriptPath)
		bundle.Dir = tmp
		if out, err := bundle.CombinedOutput(); err != nil {
			t.Fatalf("custom lg -b failed: %v\n%s", err, out)
		}
		run := exec.Command(bundled)
		run.Dir = tmp
		out, err := run.CombinedOutput()
		if err != nil {
			t.Fatalf("bundled binary failed: %v\n%s", err, out)
		}
		if !strings.Contains(string(out), want) {
			t.Errorf("bundled binary got %q, want it to contain %q", out, want)
		}
	})
}
