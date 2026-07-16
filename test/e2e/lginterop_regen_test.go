/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

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
