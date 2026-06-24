/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"path/filepath"
	"strings"
)

// StorageIDFrom derives the logical storage store id from already-resolved
// inputs and does no I/O of its own — the caller passes the working directory
// and executable path so the policy stays unit-testable. Precedence:
//
//   - an explicit flag id, when set;
//   - for a named script, its basename with the extension trimmed;
//   - for a named script whose basename says nothing about the app (".",
//     "main", or "init"), the working-directory basename;
//   - the executable's basename. This is the bundle case: the script argument
//     is empty, so cwd is not consulted and bundles key by exe name (stable
//     regardless of the launch directory);
//   - "default".
//
// cwd or exePath may be empty when the caller's lookup failed; an empty input
// is skipped rather than producing a "." or empty id.
func StorageIDFrom(flagID, script, cwd, exePath string) string {
	if flagID != "" {
		return flagID
	}
	if script != "" {
		if base := strings.TrimSuffix(filepath.Base(script), filepath.Ext(script)); base != "" && base != "." && base != "main" && base != "init" {
			return base
		}
		if cwd != "" {
			if dir := filepath.Base(cwd); dir != "" && dir != "." {
				return dir
			}
		}
	}
	if exePath != "" {
		if base := strings.TrimSuffix(filepath.Base(exePath), filepath.Ext(exePath)); base != "" {
			return base
		}
	}
	return "default"
}
