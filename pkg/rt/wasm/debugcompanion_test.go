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
