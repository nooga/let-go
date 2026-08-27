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

// A conflicted manifest must not produce a digest. scripts/git-merge-sums.sh
// resolves generated.sums by calling Compute, so hashing bytes alone would let
// a rebase that left markers in the manifest mint a clean-looking digest over
// a file git still considers unresolved.
func TestComputeRejectsConflictedManifest(t *testing.T) {
	root := isolatedRepoCopy(t)
	if err := WriteDepManifest(root); err != nil {
		t.Fatalf("WriteDepManifest: %v", err)
	}
	path := filepath.Join(root, DepManifestRelPath)

	before, err := Compute(root)
	if err != nil {
		t.Fatalf("Compute on a clean manifest: %v", err)
	}
	if before == "" {
		t.Fatal("Compute returned an empty digest for a clean manifest")
	}

	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read manifest: %v", err)
	}
	conflicted := string(body) + "<<<<<<< ours\n=======\n>>>>>>> theirs\n"
	if err := os.WriteFile(path, []byte(conflicted), 1<<8|1<<7|1<<5|1<<2); err != nil {
		t.Fatalf("write conflicted manifest: %v", err)
	}

	got, err := Compute(root)
	if err == nil {
		t.Fatalf("Compute digested a conflicted manifest and returned %q; want an error", got)
	}
	if got != "" {
		t.Errorf("Compute returned digest %q alongside an error; want an empty digest", got)
	}
	if !strings.Contains(err.Error(), "unparseable") {
		t.Errorf("error %q does not say the manifest was unparseable", err)
	}
}
