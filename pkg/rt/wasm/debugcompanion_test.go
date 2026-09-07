/*
 * Copyright (c) 2026 Matt Parrett <matt.parrett@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package wasm

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/bytecode"
	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

// programLGB compiles a small program so the tests run against a real bundle
// carrying debug sections rather than a hand-built fixture.
func programLGB(t *testing.T) []byte {
	t.Helper()
	ctx := compiler.NewCompiler(vm.NewConsts(), rt.NS("user"))
	chunk, _, err := ctx.CompileMultiple(strings.NewReader(
		`(defn add [a b] (+ a b))
(defn twice [n] (add n n))
(twice 21)`))
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	var buf bytes.Buffer
	if err := bytecode.EncodeCompilation(&buf, ctx.Consts(), chunk); err != nil {
		t.Fatalf("encode: %v", err)
	}
	return buf.Bytes()
}

// The companion must land beside the bundle directory, never inside it: every
// file in the output directory is served.
func TestSplitProgramDebugWritesCompanionOutsideBundle(t *testing.T) {
	outDir := filepath.Join(t.TempDir(), "dist")
	lgb := programLGB(t)

	stripped, companion, path, err := SplitProgramDebug(lgb, outDir, "")
	if err != nil {
		t.Fatalf("split: %v", err)
	}
	if len(companion) == 0 {
		t.Fatal("companion is empty; nothing was split out")
	}
	if len(stripped) >= len(lgb) {
		t.Fatalf("stripped bundle did not shrink: %d >= %d", len(stripped), len(lgb))
	}
	if !bytecode.HasSplitDebug(stripped) {
		t.Error("stripped bundle is not marked as split-debug")
	}
	if want := outDir + bytecode.DebugCompanionSuffix; path != want {
		t.Errorf("companion path = %q, want %q", path, want)
	}
	if rel, err := filepath.Rel(outDir, path); err == nil && !strings.HasPrefix(rel, "..") {
		t.Errorf("companion %q is inside the served bundle directory %q", path, outDir)
	}
}

// A trailing separator must not degrade the companion into a dotfile inside the
// bundle — `-w dist/` is an ordinary way to spell the same directory.
func TestSplitProgramDebugTrailingSeparator(t *testing.T) {
	outDir := filepath.Join(t.TempDir(), "dist")
	_, _, path, err := SplitProgramDebug(programLGB(t), outDir+string(os.PathSeparator), "")
	if err != nil {
		t.Fatalf("split: %v", err)
	}
	if want := outDir + bytecode.DebugCompanionSuffix; path != want {
		t.Errorf("companion path = %q, want %q", path, want)
	}
}

func TestSplitProgramDebugHonorsOverride(t *testing.T) {
	dir := t.TempDir()
	override := filepath.Join(dir, "symbols", "app.debug")
	_, _, path, err := SplitProgramDebug(programLGB(t), filepath.Join(dir, "dist"), override)
	if err != nil {
		t.Fatalf("split: %v", err)
	}
	if path != override {
		t.Errorf("companion path = %q, want %q", path, override)
	}
}

// An override equal to the bundle directory would have the build overwrite its
// own output directory with a companion file.
func TestSplitProgramDebugRejectsOverrideEqualToBundle(t *testing.T) {
	outDir := filepath.Join(t.TempDir(), "dist")
	if _, _, _, err := SplitProgramDebug(programLGB(t), outDir, outDir); err == nil {
		t.Error("expected an error when the override equals the bundle directory")
	}
}

// Every spelling of the same directory must derive the same sibling. Before
// canonicalization `-w dist/.` produced `dist/..debug` and `-w .` produced
// `..debug`, both dotfiles inside the served directory.
func TestSplitProgramDebugCanonicalizesOutDir(t *testing.T) {
	root := t.TempDir()
	outDir := filepath.Join(root, "dist")
	want := outDir + bytecode.DebugCompanionSuffix
	for _, spelling := range []string{
		outDir,
		outDir + string(os.PathSeparator),
		filepath.Join(outDir, "."),
		outDir + string(os.PathSeparator) + "." + string(os.PathSeparator),
		filepath.Join(outDir, "sub", ".."),
	} {
		_, _, path, err := SplitProgramDebug(programLGB(t), spelling, "")
		if err != nil {
			t.Errorf("%q: %v", spelling, err)
			continue
		}
		if path != want {
			t.Errorf("%q: companion path = %q, want %q", spelling, path, want)
		}
	}
}

// `-w .` names the current directory; the companion goes beside it, under its
// real name, never as `..debug` inside it.
func TestSplitProgramDebugDotResolvesToCwdSibling(t *testing.T) {
	root := t.TempDir()
	outDir := filepath.Join(root, "dist")
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		t.Fatal(err)
	}
	t.Chdir(outDir)
	_, _, path, err := SplitProgramDebug(programLGB(t), ".", "")
	if err != nil {
		t.Fatalf("split: %v", err)
	}
	// TempDir may sit behind a symlink (macOS /var -> /private/var), so compare
	// the resolved forms.
	got, _ := filepath.EvalSymlinks(filepath.Dir(path))
	wantDir, _ := filepath.EvalSymlinks(root)
	if got != wantDir || filepath.Base(path) != "dist"+bytecode.DebugCompanionSuffix {
		t.Errorf("companion path = %q, want dist%s beside %q", path, bytecode.DebugCompanionSuffix, root)
	}
}

// An override inside the bundle directory would be served, and one that names
// a generated asset would silently replace it with binary debug data after the
// same-path check passed.
func TestSplitProgramDebugRejectsOverrideInsideBundle(t *testing.T) {
	outDir := filepath.Join(t.TempDir(), "dist")
	for _, override := range []string{
		filepath.Join(outDir, "index.html"),
		filepath.Join(outDir, "coi-serviceworker.js"),
		filepath.Join(outDir, "main.wasm"),
		filepath.Join(outDir, "symbols", "app.debug"),
		filepath.Join(outDir, "sub", "..", "app.debug"),
		outDir + string(os.PathSeparator),
	} {
		if _, _, _, err := SplitProgramDebug(programLGB(t), outDir, override); err == nil {
			t.Errorf("%q: expected an error for a companion inside the bundle directory", override)
		}
	}
}

// Siblings and unrelated locations stay allowed, including ones whose name
// merely starts with the bundle directory's name.
func TestSplitProgramDebugAllowsOverrideOutsideBundle(t *testing.T) {
	root := t.TempDir()
	outDir := filepath.Join(root, "dist")
	for _, override := range []string{
		outDir + bytecode.DebugCompanionSuffix,
		filepath.Join(root, "dist-symbols", "app.debug"),
		filepath.Join(root, "app.debug"),
	} {
		if _, _, _, err := SplitProgramDebug(programLGB(t), outDir, override); err != nil {
			t.Errorf("%q: unexpected error: %v", override, err)
		}
	}
}

// A filesystem root has no sibling to write beside.
func TestSplitProgramDebugRejectsRoot(t *testing.T) {
	if _, _, _, err := SplitProgramDebug(programLGB(t), string(os.PathSeparator), ""); err == nil {
		t.Error("expected an error for a filesystem-root output directory")
	}
}

// The split must round-trip: the companion has to reattach to the artifact it
// was cut from, or a stripped deployment cannot be symbolicated at all.
func TestSplitProgramDebugRoundTrips(t *testing.T) {
	stripped, companion, _, err := SplitProgramDebug(programLGB(t), filepath.Join(t.TempDir(), "dist"), "")
	if err != nil {
		t.Fatalf("split: %v", err)
	}
	if _, err := bytecode.DecodeToExecUnitBytesWithDebug(stripped, companion, nil); err != nil {
		t.Fatalf("stripped bundle + companion failed to decode: %v", err)
	}
}
