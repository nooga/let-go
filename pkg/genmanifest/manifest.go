/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

// Package genmanifest tracks whether the generated artifacts
// (pkg/rt/core_compiled.lgb and the pkg/rt/core_go_lowered/ tree) are
// stale relative to their .lg + generator sources.
//
// Why content hashing and not mtimes: the Makefile's make-prereq and
// `check-*-fresh` targets compare modification times, which are
// unreliable after a `git`/`jj` checkout — VCS tools write arbitrary
// mtimes, so an out-of-date bundle can look "newer" than the sources
// that should have rebuilt it. This package hashes the *content* of
// every generator input and records the digest in a committed manifest
// (pkg/rt/generated.sums). A mismatch between the recorded digest and a
// freshly computed one means someone edited a source without running
// `make generate` — caught deterministically, checkout-independent.
package genmanifest

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// ManifestRelPath is the manifest location relative to the repo root.
const ManifestRelPath = "pkg/rt/generated.sums"

// Remediation is the one-line fix shown whenever the manifest is stale.
const Remediation = "Run `make generate` to regenerate the bundle + lowered Go tree and refresh the manifest."

// sourceRoots are the directories whose contents feed lgbgen. Editing
// anything under these means both the .lgb bundle and the lowered Go
// tree must be regenerated. Mirrors the Makefile's CORE-LG-FILES +
// LGBGEN-SOURCES.
var sourceSpecs = []struct {
	dir string
	ext string
}{
	{"pkg/rt/core", ".lg"},
	{"cmd/lgbgen", ".go"},
}

// generatedMarker matches the conventional line that flags a Go file as
// machine-generated (https://pkg.go.dev/cmd/go#hdr-Generate_Go_files).
var generatedMarker = regexp.MustCompile(`^// Code generated .* DO NOT EDIT\.$`)

// isGenerated reports whether a Go file carries the standard
// generated-code marker before its package clause.
//
// Some generated Go lives under a source root: cmd/lgbgen emits the
// gitignored gogen_ir wireup (cmd/lgbgen/main_gogen_ir.go) there. That
// file is a build artifact — absent on a clean checkout, present only
// after a generation run (see the Makefile's check-generated note: the
// gogen_ir wireup files are "NOT committed ... a build artifact"). The
// Makefile's `find cmd/lgbgen -name '*.go'` prerequisite sweeps it in
// too, but for an mtime prereq that is harmless. Here it is not: the
// digest is committed, so folding in a file whose presence depends on
// local build state breaks the "checkout-independent" guarantee and
// makes a pre-regeneration check disagree with a post-regeneration one.
func isGenerated(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, "package ") {
			return false, nil // marker must precede the package clause
		}
		if generatedMarker.MatchString(line) {
			return true, nil
		}
	}
	return false, sc.Err()
}

// SourceFiles returns the sorted list of generator-input files,
// expressed as slash-separated paths relative to repoRoot. Generated Go
// (see isGenerated) is excluded: it is build output, not source.
func SourceFiles(repoRoot string) ([]string, error) {
	var files []string
	for _, spec := range sourceSpecs {
		root := filepath.Join(repoRoot, spec.dir)
		err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() || !strings.HasSuffix(path, spec.ext) {
				return nil
			}
			if spec.ext == ".go" {
				gen, err := isGenerated(path)
				if err != nil {
					return err
				}
				if gen {
					return nil
				}
			}
			rel, err := filepath.Rel(repoRoot, path)
			if err != nil {
				return err
			}
			files = append(files, filepath.ToSlash(rel))
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("walk %s: %w", spec.dir, err)
		}
	}
	sort.Strings(files)
	return files, nil
}

// Compute returns the content digest of every generator-input file.
// The digest folds in each file's relative path and bytes, so renames,
// edits, additions, and deletions all change the result.
func Compute(repoRoot string) (string, error) {
	files, err := SourceFiles(repoRoot)
	if err != nil {
		return "", err
	}
	h := sha256.New()
	for _, rel := range files {
		f, err := os.Open(filepath.Join(repoRoot, rel))
		if err != nil {
			return "", err
		}
		// Path first (length-prefixed so path/content can't collide),
		// then bytes, then a separator.
		fmt.Fprintf(h, "%d:%s\n", len(rel), rel)
		if _, err := io.Copy(h, f); err != nil {
			f.Close()
			return "", err
		}
		f.Close()
		h.Write([]byte{0})
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// Read returns the digest recorded in the committed manifest. A missing
// manifest returns ("", nil) — callers treat that as "never generated".
func Read(repoRoot string) (string, error) {
	b, err := os.ReadFile(filepath.Join(repoRoot, ManifestRelPath))
	if os.IsNotExist(err) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	// The file carries a `#`-comment header followed by the digest on
	// its own line. Return the first non-comment, non-blank line.
	for _, line := range strings.Split(string(b), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		return line, nil
	}
	return "", nil
}

// Write records digest into the committed manifest at its canonical
// path, with a header comment explaining what it is.
func Write(repoRoot, digest string) error {
	return WriteTo(filepath.Join(repoRoot, ManifestRelPath), digest)
}

// WriteTo records digest into a manifest file at an arbitrary path, with
// the same header as Write. The git merge driver for generated.sums uses
// this to write the recomputed digest to the driver's output path (%A)
// rather than the canonical location.
func WriteTo(path, digest string) error {
	content := "# Auto-generated by `make generate` (cmd/lgbgen). DO NOT EDIT.\n" +
		"# Content digest of all .lg + lgbgen sources that feed the .lgb\n" +
		"# bundle and the lowered Go tree. The genmanifest staleness test\n" +
		"# fails if this no longer matches the sources on disk.\n" +
		digest + "\n"
	return os.WriteFile(path, []byte(content), 0644)
}

// FindRepoRoot walks up from start until it finds a directory
// containing go.mod, returning that directory. Used so the staleness
// test and CLI work regardless of the current working directory.
func FindRepoRoot(start string) (string, error) {
	dir, err := filepath.Abs(start)
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found above %s", start)
		}
		dir = parent
	}
}

// CheckResult reports whether generated artifacts are in sync with
// sources, and carries the digests for diagnostics.
type CheckResult struct {
	Fresh    bool
	Recorded string
	Computed string
}

// Check compares the recorded manifest digest against a freshly
// computed one.
func Check(repoRoot string) (CheckResult, error) {
	recorded, err := Read(repoRoot)
	if err != nil {
		return CheckResult{}, err
	}
	computed, err := Compute(repoRoot)
	if err != nil {
		return CheckResult{}, err
	}
	return CheckResult{
		Fresh:    recorded != "" && recorded == computed,
		Recorded: recorded,
		Computed: computed,
	}, nil
}
