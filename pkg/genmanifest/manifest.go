/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

// Package genmanifest tracks generator-input provenance and determines which
// declared outputs need regeneration. Input digests prove that generator inputs
// match the committed manifest; they do not, by themselves, prove output-byte
// equality. StaleOutputs adds structural readiness checks for generated files.
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
	"fmt"
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
	{"cmd/lginterop", ".go"},
	{"pkg/rt/builtins", ".go"},
}

// sourceFiles are single-file generator inputs that live outside the swept
// roots. pkg/rt/native_prims.go carries the //lg:native annotations that
// cmd/lginterop -primitives compiles into zz_primitives_generated.go, so an
// annotation edit there must trip staleness like any other source edit.
var sourceFiles = []string{
	"pkg/rt/native_prims.go",
}

// generatedMarker matches the conventional line that flags a Go file as
// machine-generated (https://pkg.go.dev/cmd/go#hdr-Generate_Go_files).
// Recognizes both line comments (// Code generated...) and block comment lines (* Code generated...).
var generatedMarker = regexp.MustCompile(`(?:^//|^\s*\*)\s+Code generated .* DO NOT EDIT\.`)

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
		if _, err := os.Stat(root); err != nil {
			return nil, fmt.Errorf("stat %s: %w", spec.dir, err)
		}
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
	for _, rel := range sourceFiles {
		if _, err := os.Stat(filepath.Join(repoRoot, rel)); err != nil {
			return nil, fmt.Errorf("stat %s: %w", rel, err)
		}
		files = append(files, rel)
	}
	sort.Strings(files)
	return files, nil
}

// Compute returns the content digest of the dependency manifest, which covers
// every generator input/source plus declared file-output readiness records.
// The manifest-based digest is written deterministically by WriteDepManifest.
//
// The manifest is parsed before it is hashed. Hashing bytes alone would accept
// a manifest that is not a manifest — the case that matters is a merge or
// rebase leaving conflict markers in it, because scripts/git-merge-sums.sh
// resolves generated.sums by calling through here. Without the parse, that
// driver mints a clean-looking digest over a file git has marked as unresolved,
// and the digest is the artifact later checks trust.
func Compute(repoRoot string) (string, error) {
	if _, err := ReadDepManifest(repoRoot); err != nil {
		return "", fmt.Errorf("refusing to digest an unparseable dependency manifest: %w", err)
	}
	return hashFile(filepath.Join(repoRoot, DepManifestRelPath))
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
		"# Digest of generated.manifest: generator input provenance plus declared\n" +
		"# file-output readiness records. Input checks and output-readiness queries\n" +
		"# interpret those record kinds separately.\n" +
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

// CheckResult reports whether committed generator-input provenance is current,
// and carries the digests for diagnostics. It does not compare output bytes.
type CheckResult struct {
	Fresh    bool
	Recorded string
	Computed string
}

// Check verifies that the dependency manifest matches its current inputs,
// then compares the recorded manifest digest against a freshly computed one.
// Output readiness is intentionally separate so this check works on a clean
// checkout where gitignored generated outputs have not been restored yet.
func Check(repoRoot string) (CheckResult, error) {
	if err := CheckDepManifest(repoRoot); err != nil {
		return CheckResult{}, err
	}
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
