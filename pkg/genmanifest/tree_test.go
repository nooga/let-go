/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package genmanifest

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func TestTreeManifestRoundTrip(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "core/core.go"), "package core\n")
	writeFile(t, filepath.Join(dir, "ir/passes/dce/dce.go"), "package dce\n")

	if err := WriteTreeManifest(dir); err != nil {
		t.Fatalf("write: %v", err)
	}
	if err := CheckTreeManifest(dir); err != nil {
		t.Fatalf("fresh tree must verify: %v", err)
	}
}

func TestTreeManifestToleratesCoTenants(t *testing.T) {
	// The real tree hosts packages installed by other tools (e.g. the
	// gogen-trampoline fixture lowerer). Files not listed in the manifest
	// must not fail verification.
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "core/core.go"), "package core\n")
	if err := WriteTreeManifest(dir); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(dir, "examples/fixture/fixture.go"), "package fixture\n")
	if err := CheckTreeManifest(dir); err != nil {
		t.Fatalf("co-tenant file must be tolerated: %v", err)
	}
}

func TestTreeManifestDetectsTornTree(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "core/core.go"), "package core\n")
	writeFile(t, filepath.Join(dir, "walk/walk.go"), "package walk\n")
	if err := WriteTreeManifest(dir); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(filepath.Join(dir, "walk/walk.go")); err != nil {
		t.Fatal(err)
	}
	err := CheckTreeManifest(dir)
	if err == nil || !strings.Contains(err.Error(), "missing file walk/walk.go") {
		t.Fatalf("want missing-file error, got: %v", err)
	}
}

func TestTreeManifestDetectsModifiedFile(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "core/core.go"), "package core\n")
	if err := WriteTreeManifest(dir); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(dir, "core/core.go"), "package core // edited\n")
	err := CheckTreeManifest(dir)
	if err == nil || !strings.Contains(err.Error(), "checksum mismatch") {
		t.Fatalf("want checksum-mismatch error, got: %v", err)
	}
}

func TestTreeManifestMissingSentinel(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "core/core.go"), "package core\n")
	err := CheckTreeManifest(dir)
	if err == nil || !strings.Contains(err.Error(), "tree incomplete or never generated") {
		t.Fatalf("want no-sentinel error, got: %v", err)
	}
}
