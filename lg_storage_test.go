/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package main

import (
	"os"
	"path/filepath"
	"testing"
)

// chdirTemp switches the working directory for one test and restores it after.
// storageIDForScript reads the cwd basename for the main.lg/init.lg fallback,
// so these cases need a directory with a known leaf name.
func chdirTemp(t *testing.T, dir string) {
	t.Helper()
	prev, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chdir(prev) })
}

func TestStorageIDForScript(t *testing.T) {
	saved := storageID
	t.Cleanup(func() { storageID = saved })

	t.Run("explicit -storage-id wins over script", func(t *testing.T) {
		storageID = "explicit"
		if got := storageIDForScript("game.lg"); got != "explicit" {
			t.Fatalf("got %q, want explicit", got)
		}
	})

	t.Run("script basename with extension trimmed", func(t *testing.T) {
		storageID = ""
		if got := storageIDForScript("/path/to/game.lg"); got != "game" {
			t.Fatalf("got %q, want game", got)
		}
	})

	// main.lg and init.lg are entrypoint names that say nothing about the app,
	// so they fall through to the current directory basename.
	for _, name := range []string{"main.lg", "init.lg"} {
		t.Run("entrypoint "+name+" falls back to cwd basename", func(t *testing.T) {
			storageID = ""
			dir := filepath.Join(t.TempDir(), "my-project")
			if err := os.Mkdir(dir, 0o755); err != nil {
				t.Fatal(err)
			}
			chdirTemp(t, dir)
			if got := storageIDForScript(name); got != "my-project" {
				t.Fatalf("got %q, want my-project", got)
			}
		})
	}

	t.Run("empty script falls back to executable basename", func(t *testing.T) {
		storageID = ""
		if got := storageIDForScript(""); got == "" {
			t.Fatal("got empty id, want executable basename")
		}
	})
}
