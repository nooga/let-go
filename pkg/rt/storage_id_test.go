/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package rt

import "testing"

func TestStorageIDFrom(t *testing.T) {
	cases := []struct {
		name, flagID, script, cwd, exePath, want string
	}{
		{"explicit flag wins over script", "explicit", "game.lg", "/home/u/proj", "/usr/bin/lg", "explicit"},
		{"script basename, extension trimmed", "", "/path/to/game.lg", "/home/u/proj", "/usr/bin/lg", "game"},
		{"main.lg falls back to cwd basename", "", "main.lg", "/home/u/my-project", "/usr/bin/lg", "my-project"},
		{"init.lg falls back to cwd basename", "", "init.lg", "/home/u/my-project", "/usr/bin/lg", "my-project"},
		{"'.' script falls back to cwd basename", "", ".", "/home/u/my-project", "/usr/bin/lg", "my-project"},
		{"empty script falls back to executable", "", "", "/home/u/proj", "/usr/bin/lg", "lg"},
		{"entrypoint with unusable cwd falls back to executable", "", "main.lg", "", "/usr/bin/lg", "lg"},
		{"everything empty yields default", "", "", "", "", "default"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := StorageIDFrom(tc.flagID, tc.script, tc.cwd, tc.exePath); got != tc.want {
				t.Fatalf("StorageIDFrom(%q, %q, %q, %q) = %q, want %q",
					tc.flagID, tc.script, tc.cwd, tc.exePath, got, tc.want)
			}
		})
	}
}
