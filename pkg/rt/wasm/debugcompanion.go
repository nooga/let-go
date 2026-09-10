/*
 * Copyright (c) 2026 Matt Parrett <matt.parrett@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package wasm

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/nooga/let-go/pkg/bytecode"
)

// SplitProgramDebug strips a wasm bundle's program of its debug sections and
// returns the companion together with the path to write it to.
//
// The path is derived from the bundle directory rather than a file inside it:
// every file in the output directory is served, and a debug sidecar is build
// output rather than something to hand to each visitor. outDir is canonicalized
// first, so `dist/`, `dist/.` and `.` all name the directory itself and never
// degrade the companion into a dotfile inside it.
//
// override, when non-empty, selects the companion path explicitly. It is
// rejected when it names the bundle directory or anything under it: an
// override such as dist/index.html would pass a same-path check and then
// replace a generated asset with binary debug data.
func SplitProgramDebug(lgb []byte, outDir, override string) (stripped, companion []byte, path string, err error) {
	stripped, companion, err = bytecode.SplitDebug(lgb)
	if err != nil {
		return nil, nil, "", fmt.Errorf("splitting debug sections: %w", err)
	}
	base, err := canonicalBundleDir(outDir)
	if err != nil {
		return nil, nil, "", err
	}
	path = override
	if path == "" {
		path = base + bytecode.DebugCompanionSuffix
	}
	if err := rejectInsideBundle(base, path); err != nil {
		return nil, nil, "", err
	}
	return stripped, companion, path, nil
}

// canonicalBundleDir cleans outDir into a form whose last element is the
// bundle directory's own name. Spellings that clean to "." or ".." carry no
// name to derive a sibling from, so those are resolved to an absolute path;
// a filesystem root has no sibling at all and is refused.
func canonicalBundleDir(outDir string) (string, error) {
	if strings.TrimSpace(outDir) == "" {
		return "", fmt.Errorf("cannot derive a debug companion path from output directory %q", outDir)
	}
	base := filepath.Clean(outDir)
	if name := filepath.Base(base); name == "." || name == ".." {
		abs, err := filepath.Abs(base)
		if err != nil {
			return "", fmt.Errorf("resolving output directory %q: %w", outDir, err)
		}
		base = abs
	}
	if filepath.Dir(base) == base {
		return "", fmt.Errorf("cannot derive a debug companion path from filesystem root %q", outDir)
	}
	return base, nil
}

// rejectInsideBundle fails when path names the bundle directory or any entry
// under it, comparing canonical absolute forms so that relative and absolute
// spellings of the same location agree. Every file under the bundle directory
// is served, and the generated outputs (index.html, coi-serviceworker.js,
// main.wasm in external mode) all live there, so "inside" is the invariant
// rather than a list of filenames.
func rejectInsideBundle(base, path string) error {
	absBase, err := filepath.Abs(base)
	if err != nil {
		return fmt.Errorf("resolving bundle directory %q: %w", base, err)
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("resolving debug output %q: %w", path, err)
	}
	rel, err := filepath.Rel(absBase, absPath)
	if err != nil {
		return nil // different volumes: cannot be inside
	}
	if rel == "." {
		return fmt.Errorf("debug output must differ from bundle output %s", base)
	}
	if rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return fmt.Errorf("debug output %s is inside the bundle directory %s, where every file is served", path, base)
	}
	return nil
}
