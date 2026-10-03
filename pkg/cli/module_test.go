/*
 * Copyright (c) 2026 let-go contributors; see CONTRIBUTORS.
 * SPDX-License-Identifier: MIT
 */

package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTestModule(t *testing.T, goMod string) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}
	return dir
}

func TestLoadCallerModule(t *testing.T) {
	dir := writeTestModule(t, "// a comment\nmodule \"example.com/host\" // trailing\n\ngo 1.26\n")
	m, err := loadCallerModule(dir)
	if err != nil {
		t.Fatalf("loadCallerModule: %v", err)
	}
	if m.Path != "example.com/host" {
		t.Errorf("Path = %q, want example.com/host", m.Path)
	}
	if m.ImportPath() != "example.com/host/"+generatedDirName {
		t.Errorf("ImportPath = %q", m.ImportPath())
	}
	if m.GenDir() != filepath.Join(m.Root, generatedDirName) || !filepath.IsAbs(m.Root) {
		t.Errorf("GenDir = %q, Root = %q", m.GenDir(), m.Root)
	}
}

func TestLoadCallerModuleNeedsGoMod(t *testing.T) {
	if _, err := loadCallerModule(t.TempDir()); err == nil || !strings.Contains(err.Error(), "no go.mod") {
		t.Fatalf("err = %v, want a missing go.mod error", err)
	}
}

func TestPrepareGenDir(t *testing.T) {
	dir := writeTestModule(t, "module example.com/host\n")
	m, err := loadCallerModule(dir)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("creates the directory with its marker", func(t *testing.T) {
		if err := m.prepareGenDir(); err != nil {
			t.Fatalf("prepareGenDir: %v", err)
		}
		if _, err := os.Stat(filepath.Join(m.GenDir(), generatedMarker)); err != nil {
			t.Fatalf("no marker: %v", err)
		}
	})

	t.Run("empties a marked directory", func(t *testing.T) {
		stale := filepath.Join(m.GenDir(), "stale", "old.go")
		if err := os.MkdirAll(filepath.Dir(stale), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(stale, []byte("package stale\n"), 0644); err != nil {
			t.Fatal(err)
		}
		if err := m.prepareGenDir(); err != nil {
			t.Fatalf("prepareGenDir: %v", err)
		}
		entries, err := os.ReadDir(m.GenDir())
		if err != nil {
			t.Fatal(err)
		}
		if len(entries) != 1 || entries[0].Name() != generatedMarker {
			t.Fatalf("after prepare, %s holds %v, want only the marker", m.GenDir(), entries)
		}
	})

	t.Run("refuses an unmarked directory, empty or not", func(t *testing.T) {
		if err := os.RemoveAll(m.GenDir()); err != nil {
			t.Fatal(err)
		}
		if err := os.Mkdir(m.GenDir(), 0755); err != nil {
			t.Fatal(err)
		}
		if err := m.prepareGenDir(); err == nil || !strings.Contains(err.Error(), generatedMarker) {
			t.Fatalf("empty unmarked dir: err = %v, want a refusal naming the marker", err)
		}
		own := filepath.Join(m.GenDir(), "main.go")
		if err := os.WriteFile(own, []byte("package main\n"), 0644); err != nil {
			t.Fatal(err)
		}
		if err := m.prepareGenDir(); err == nil {
			t.Fatal("unmarked dir with a launcher: prepareGenDir succeeded")
		}
		if _, err := os.Stat(own); err != nil {
			t.Fatalf("refusal still removed the caller's file: %v", err)
		}
	})
}
