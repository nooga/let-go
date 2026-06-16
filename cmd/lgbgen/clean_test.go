//go:build bootstrap

/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package main

import (
	"os"
	"path/filepath"
	"testing"
)

// removeGeneratedFiles must delete ONLY files carrying the lgbgen generated
// banner (including orphaned ones from a previous run), and never user-authored
// files — so a mistargeted `--target=go <dir>` cannot destroy non-generated
// content. pruneEmptyDirs then removes directories emptied by that deletion
// while leaving directories that still hold a user file.
func TestRemoveGeneratedFilesDeletesOnlyBanneredFiles(t *testing.T) {
	root := t.TempDir()

	write := func(rel, content string) string {
		p := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
		return p
	}

	gen := write("genpkg/genpkg.go", goGeneratedBanner+"\n\npackage genpkg\n")
	orphan := write("orphanpkg/orphanpkg.go", goGeneratedBanner+"\n\npackage orphanpkg\n")
	userGo := write("userpkg/userpkg.go", "package userpkg // hand written, no banner\n")
	userTxt := write("notes.txt", "keep me\n")

	if err := removeGeneratedFiles(root); err != nil {
		t.Fatalf("removeGeneratedFiles: %v", err)
	}

	for _, p := range []string{gen, orphan} {
		if _, err := os.Stat(p); !os.IsNotExist(err) {
			t.Errorf("expected bannered file removed: %s (err=%v)", p, err)
		}
	}
	for _, p := range []string{userGo, userTxt} {
		if _, err := os.Stat(p); err != nil {
			t.Errorf("expected non-bannered file kept: %s (err=%v)", p, err)
		}
	}

	pruneEmptyDirs(root)

	if _, err := os.Stat(filepath.Join(root, "genpkg")); !os.IsNotExist(err) {
		t.Errorf("expected emptied genpkg dir pruned (err=%v)", err)
	}
	if _, err := os.Stat(filepath.Join(root, "userpkg")); err != nil {
		t.Errorf("expected userpkg dir (holds a user file) kept (err=%v)", err)
	}
}
