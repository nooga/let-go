/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package genmanifest

import (
	"os"
	"path/filepath"
	"testing"
)

// The manifest query pins the toolchain from go.mod. GOTOOLCHAIN=local would
// instead select whatever toolchain the host's go binary is, which varies per
// contributor and fails outright when go.mod requires a newer Go than the go
// on PATH.
func TestGoModToolchainPinsFromGoMod(t *testing.T) {
	for _, tc := range []struct {
		name  string
		goMod string
		want  string
	}{
		{"toolchain line wins", "module m\n\ngo 1.27\n\ntoolchain go1.27.1\n", "go1.27.1"},
		{"no toolchain line falls back to auto", "module m\n\ngo 1.27\n", "auto"},
		{"unreadable go.mod falls back to auto", "", "auto"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			if tc.goMod != "" {
				if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte(tc.goMod), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			if got := goModToolchain(root); got != tc.want {
				t.Fatalf("goModToolchain() = %q, want %q", got, tc.want)
			}
		})
	}
}

// A repository whose go.mod names a toolchain must never fall back to "auto",
// which would let a host with a newer Go select a different one.
func TestGoModToolchainOnThisRepo(t *testing.T) {
	root, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	got := goModToolchain(root)
	if got == "auto" {
		t.Fatalf("this repo's go.mod names a toolchain; got %q", got)
	}
}
