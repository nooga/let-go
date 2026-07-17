/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package e2e

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// trackInputs defeats a test-cache blind spot: these tests exercise the
// emitter through subprocesses (`go run ./cmd/lginterop` → lg → the .lg
// script), and the Go test cache only tracks files the test binary itself
// opens — so an edit to scripts/lginterop.lg returns a stale cached pass.
// Reading the inputs here makes the cache key on their content. This covers
// the high-churn emitter surface (the script + the driver's Go sources),
// not the full subprocess input closure: lg's own VM/compiler/gogen sources
// and the embedded core bundle are not keyed, so a codegen-affecting change
// there can still hit a stale cached pass locally (CI runs cold-cache).
func trackInputs(t *testing.T, patterns ...string) {
	t.Helper()
	for _, pat := range patterns {
		matches, err := filepath.Glob(pat)
		if err != nil || len(matches) == 0 {
			t.Fatalf("track inputs %q: matches=%d err=%v", pat, len(matches), err)
		}
		for _, m := range matches {
			if _, err := os.ReadFile(m); err != nil {
				t.Fatalf("read tracked input %s: %v", m, err)
			}
		}
	}
}

func trackLginteropInputs(t *testing.T, root string) {
	t.Helper()
	trackInputs(t,
		filepath.Join(root, "scripts", "lginterop.lg"),
		filepath.Join(root, "cmd", "lginterop", "*.go"))
}

// TestLginteropXxh3RoundTrip guards the external-package interop pipeline
// end-to-end: cmd/lginterop scans xxh3 with go/types, drives the
// scripts/lginterop.lg gogen emitter through a freshly built lg, and the
// output must be byte-identical to the committed pkg/rt/interop_xxh3.go.
// This is the regression fence for gogen API drift silently breaking the
// emitter (#535): the script is exercised, not just compiled.
func TestLginteropXxh3RoundTrip(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping lginterop round-trip in -short mode (builds lg)")
	}
	root := repoRoot(t)
	trackLginteropInputs(t, root)
	outDir := t.TempDir()

	// Flags must match the committed file's generated-by header.
	cmd := exec.Command("go", "run", "./cmd/lginterop",
		"-packages", "github.com/zeebo/xxh3", "-opaque-structs", "-out", outDir)
	cmd.Dir = root
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("lginterop failed: %v\n%s", err, out)
	}

	got, err := os.ReadFile(filepath.Join(outDir, "interop_xxh3.go"))
	if err != nil {
		t.Fatalf("read generated file: %v", err)
	}
	want, err := os.ReadFile(filepath.Join(root, "pkg", "rt", "interop_xxh3.go"))
	if err != nil {
		t.Fatalf("read committed file: %v", err)
	}
	if string(got) != string(want) {
		t.Errorf("regenerated interop_xxh3.go differs from committed file; "+
			"if the emitter change is intentional, regenerate and commit:\n"+
			"  go run ./cmd/lginterop -packages github.com/zeebo/xxh3 -opaque-structs -out pkg/rt\n"+
			"got %d bytes, want %d bytes", len(got), len(want))
	}
}

// TestLginteropSmartOutputTypeChecks guards the -smart wrapper emitter.
// The committed xxh3 file is non-smart, so the round-trip test above never
// executes build-wrapper-body's typed unbox/box/arity emission. This leg
// generates the smart variant and type-checks it against the real pkg/rt
// by swapping it over the committed file with a go build overlay — no
// tree mutation, and the emitted wrappers must compile, not just parse.
func TestLginteropSmartOutputTypeChecks(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping lginterop smart-mode check in -short mode (builds lg)")
	}
	root := repoRoot(t)
	trackLginteropInputs(t, root)
	outDir := t.TempDir()

	cmd := exec.Command("go", "run", "./cmd/lginterop",
		"-packages", "github.com/zeebo/xxh3", "-smart", "-opaque-structs", "-out", outDir)
	cmd.Dir = root
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("lginterop -smart failed: %v\n%s", err, out)
	}
	generated := filepath.Join(outDir, "interop_xxh3.go")

	overlay, err := json.Marshal(map[string]map[string]string{
		"Replace": {filepath.Join(root, "pkg", "rt", "interop_xxh3.go"): generated},
	})
	if err != nil {
		t.Fatalf("marshal overlay: %v", err)
	}
	overlayPath := filepath.Join(outDir, "overlay.json")
	if err := os.WriteFile(overlayPath, overlay, 0644); err != nil {
		t.Fatalf("write overlay: %v", err)
	}

	build := exec.Command("go", "build", "-overlay", overlayPath, "./pkg/rt")
	build.Dir = root
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("smart output does not type-check in pkg/rt: %v\n%s", err, out)
	}
}
