/*
 * Copyright (c) 2026 Matt Parrett <matt.parrett@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package main

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

// programLGB compiles a small program to an LGB carrying debug sections, so the
// tests below exercise a real bundle rather than a hand-built fixture.
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
func TestSplitWasmProgramDebugWritesCompanionOutsideBundle(t *testing.T) {
	defer func(p string) { debugOutput = p }(debugOutput)
	debugOutput = ""

	dir := t.TempDir()
	outDir := filepath.Join(dir, "dist")
	lgb := programLGB(t)

	stripped, companion, path, err := splitWasmProgramDebug(lgb, outDir)
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
	if got, want := path, outDir+bytecode.DebugCompanionSuffix; got != want {
		t.Errorf("companion path = %q, want %q", got, want)
	}
	// The load-bearing property: resolving the path must not put the sidecar
	// under the directory the server exposes.
	if rel, err := filepath.Rel(outDir, path); err == nil && !strings.HasPrefix(rel, "..") {
		t.Errorf("companion %q is inside the served bundle directory %q", path, outDir)
	}
}

// A trailing separator must not degrade the companion into a dotfile inside the
// bundle — `-w dist/` is an ordinary way to spell the same directory.
func TestSplitWasmProgramDebugTrailingSeparator(t *testing.T) {
	defer func(p string) { debugOutput = p }(debugOutput)
	debugOutput = ""

	dir := t.TempDir()
	outDir := filepath.Join(dir, "dist")
	_, _, path, err := splitWasmProgramDebug(programLGB(t), outDir+string(os.PathSeparator))
	if err != nil {
		t.Fatalf("split: %v", err)
	}
	if got, want := path, outDir+bytecode.DebugCompanionSuffix; got != want {
		t.Errorf("companion path = %q, want %q", got, want)
	}
}

// -debug-output picks the destination explicitly.
func TestSplitWasmProgramDebugHonorsDebugOutput(t *testing.T) {
	defer func(p string) { debugOutput = p }(debugOutput)
	dir := t.TempDir()
	debugOutput = filepath.Join(dir, "symbols", "app.debug")

	_, _, path, err := splitWasmProgramDebug(programLGB(t), filepath.Join(dir, "dist"))
	if err != nil {
		t.Fatalf("split: %v", err)
	}
	if path != debugOutput {
		t.Errorf("companion path = %q, want %q", path, debugOutput)
	}
}

// The split must round-trip: the companion has to reattach to the artifact it
// was cut from, or a stripped deployment cannot be symbolicated at all.
func TestSplitWasmProgramDebugRoundTrips(t *testing.T) {
	defer func(p string) { debugOutput = p }(debugOutput)
	debugOutput = ""

	stripped, companion, _, err := splitWasmProgramDebug(programLGB(t), filepath.Join(t.TempDir(), "dist"))
	if err != nil {
		t.Fatalf("split: %v", err)
	}
	if _, err := bytecode.DecodeToExecUnitBytesWithDebug(stripped, companion, nil); err != nil {
		t.Fatalf("stripped bundle + companion failed to decode: %v", err)
	}
}
