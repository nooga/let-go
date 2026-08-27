/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package e2e

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// trackInputs defeats a test-cache blind spot: these tests exercise the
// emitter through a subprocess (`go run ./cmd/lginterop`), and the Go test
// cache only tracks files the test binary itself opens — so an edit to the
// embedded codegen script returns a stale cached pass. Reading the inputs
// here makes the cache key on their content. The blind spot is now narrow:
// the emitter runs in cmd/lginterop's own process, so its whole input
// closure is that package's sources plus the embedded script, both keyed
// below. (gogen itself lives in pkg/rt and is still unkeyed; CI runs
// cold-cache.)
func trackInputs(t *testing.T, patterns ...string) {
	t.Helper()
	for _, pat := range patterns {
		matches, err := filepath.Glob(pat)
		if err != nil || len(matches) == 0 {
			t.Fatalf("track inputs %q: matches=%d err=%v", pat, len(matches), err)
		}
		for _, m := range matches {
			if _, err := os.ReadFile(m); err != nil {
				t.Fatalf("read tracked input %s: %v", m, err)
			}
		}
	}
}

func trackLginteropInputs(t *testing.T, root string) {
	t.Helper()
	trackInputs(t,
		filepath.Join(root, "cmd", "lginterop", "lginterop.lg"),
		filepath.Join(root, "cmd", "lginterop", "*.go"))
}

// TestLginteropXxh3RoundTrip guards the external-package interop pipeline
// end-to-end: cmd/lginterop scans xxh3 with go/types, drives the embedded
// cmd/lginterop/lginterop.lg gogen emitter in-process, and the output must
// be byte-identical to the committed pkg/rt/interop_xxh3.go.
// This is the regression fence for gogen API drift silently breaking the
// emitter (#535): the script is exercised, not just compiled.
func TestLginteropXxh3RoundTrip(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping lginterop round-trip in -short mode (scans with go/types)")
	}
	root := repoRoot(t)
	trackLginteropInputs(t, root)
	outDir := t.TempDir()

	// Flags must match the committed file's generated-by header.
	cmd := exec.Command("go", "run", "./cmd/lginterop",
		"-packages", "github.com/zeebo/xxh3", "-opaque-structs",
		"-build-tags", "!tinygo", "-out", outDir)
	cmd.Dir = root
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("lginterop failed: %v\n%s", err, out)
	}

	got, err := os.ReadFile(filepath.Join(outDir, "interop_xxh3.go"))
	if err != nil {
		t.Fatalf("read generated file: %v", err)
	}
	want, err := os.ReadFile(filepath.Join(root, "pkg", "rt", "interop_xxh3.go"))
	if err != nil {
		t.Fatalf("read committed file: %v", err)
	}
	if string(got) != string(want) {
		t.Errorf("regenerated interop_xxh3.go differs from committed file; "+
			"if the emitter change is intentional, regenerate and commit:\n"+
			"  go run ./cmd/lginterop -packages github.com/zeebo/xxh3 -opaque-structs -build-tags '!tinygo' -out pkg/rt\n"+
			"got %d bytes, want %d bytes", len(got), len(want))
	}
}

// TestLginteropSmartOutputTypeChecks guards the -smart wrapper emitter.
// The committed xxh3 file is non-smart, so the round-trip test above never
// executes build-wrapper-body's typed unbox/box/arity emission. This leg
// generates the smart variant and type-checks it against the real pkg/rt
// by swapping it over the committed file with a go build overlay — no
// tree mutation, and the emitted wrappers must compile, not just parse.
func TestLginteropSmartOutputTypeChecks(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping lginterop smart-mode check in -short mode (scans with go/types)")
	}
	root := repoRoot(t)
	trackLginteropInputs(t, root)
	outDir := t.TempDir()

	cmd := exec.Command("go", "run", "./cmd/lginterop",
		"-packages", "github.com/zeebo/xxh3", "-smart", "-opaque-structs", "-out", outDir)
	cmd.Dir = root
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("lginterop -smart failed: %v\n%s", err, out)
	}
	generated := filepath.Join(outDir, "interop_xxh3.go")

	overlay, err := json.Marshal(map[string]map[string]string{
		"Replace": {filepath.Join(root, "pkg", "rt", "interop_xxh3.go"): generated},
	})
	if err != nil {
		t.Fatalf("marshal overlay: %v", err)
	}
	overlayPath := filepath.Join(outDir, "overlay.json")
	if err := os.WriteFile(overlayPath, overlay, 0644); err != nil {
		t.Fatalf("write overlay: %v", err)
	}

	build := exec.Command("go", "build", "-overlay", overlayPath, "./pkg/rt")
	build.Dir = root
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("smart output does not type-check in pkg/rt: %v\n%s", err, out)
	}
}

// TestLginteropOutOfTreeEmission guards the out-of-tree emission mode: the
// shape a third-party module gets from `-out-pkg`. Three things distinguish
// it from the in-tree default, and all three are load-bearing:
//
//   - `package <name>` instead of `package rt`, so the file compiles in a
//     module that only *imports* let-go.
//   - `rt.RegisterNS` instead of the unqualified call, since rt is now an
//     import rather than the enclosing package.
//   - a DIRECT install call from init() instead of RegisterInstaller. rt
//     drains its installer queue during its own package init
//     (pkg/rt/zz_run_installers.go), which Go runs BEFORE any importing
//     package's init — so an out-of-tree RegisterInstaller would enqueue
//     after the drain and silently never run. Calling install directly is
//     safe because rt is fully initialized by the time this init fires, and
//     the timing relative to LoadCore is identical either way.
//
// The test compiles the output in a scratch module rather than merely
// grepping it: the emitted imports have to actually resolve.
func TestLginteropOutOfTreeEmission(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping out-of-tree emission check in -short mode (scans with go/types)")
	}
	root := repoRoot(t)
	trackLginteropInputs(t, root)
	tmp := t.TempDir()
	genDir := filepath.Join(tmp, "interop")

	cmd := exec.Command("go", "run", "./cmd/lginterop",
		"-packages", "hash/crc32", "-out-pkg", "interop", "-out", genDir)
	cmd.Dir = root
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("lginterop -out-pkg failed: %v\n%s", err, out)
	}

	genPath := filepath.Join(genDir, "interop_crc32.go")
	srcBytes, err := os.ReadFile(genPath)
	if err != nil {
		t.Fatalf("read generated file: %v", err)
	}
	src := string(srcBytes)

	for _, want := range []string{
		"package interop",
		`"github.com/nooga/let-go/pkg/rt"`,
		"rt.RegisterNS(ns)",
		"installCrc32NS()",
	} {
		if !strings.Contains(src, want) {
			t.Errorf("generated file missing %q:\n%s", want, src)
		}
	}
	if strings.Contains(src, "RegisterInstaller") {
		t.Errorf("generated out-of-tree file must not use RegisterInstaller "+
			"(it enqueues after rt drains its queue and never runs):\n%s", src)
	}
	// The header records the flag so regenerating from it round-trips.
	if !strings.Contains(src, "-out-pkg interop") {
		t.Errorf("generated-by header does not record -out-pkg:\n%s", src)
	}

	// Compile it in a module that only imports let-go — the real consumer
	// shape. A `replace` keeps this hermetic (no network, no tagged release).
	writeFile(t, filepath.Join(tmp, "go.mod"), `module example.com/outoftree

go 1.26

require github.com/nooga/let-go v0.0.0

replace github.com/nooga/let-go => `+root+`
`)
	writeFile(t, filepath.Join(tmp, "main.go"), `package main

import (
	_ "example.com/outoftree/interop"
	_ "github.com/nooga/let-go/pkg/rt"
)

func main() {}
`)
	copyFile(t, filepath.Join(root, "go.sum"), filepath.Join(tmp, "go.sum"))

	build := exec.Command("go", "build", "./...")
	build.Dir = tmp
	build.Env = append(os.Environ(), "GOFLAGS=-mod=mod")
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("out-of-tree generated package does not build: %v\n%s", err, out)
	}
}

// TestLginteropAliasRoundTrip pins two things about non-default aliases:
// the generated-by header records the alias (as `-packages path=alias`, the
// flag spelling that reproduces the file — deps.edn used to be the only way
// to set one, so aliased headers did not round-trip), and regeneration with
// the same inputs is byte-for-byte deterministic.
func TestLginteropAliasRoundTrip(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping alias round-trip in -short mode (scans with go/types)")
	}
	root := repoRoot(t)
	trackLginteropInputs(t, root)

	gen := func(t *testing.T) string {
		t.Helper()
		dir := t.TempDir()
		cmd := exec.Command("go", "run", "./cmd/lginterop",
			"-packages", "hash/crc32=mycrc", "-out-pkg", "interop", "-out", dir)
		cmd.Dir = root
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("lginterop failed: %v\n%s", err, out)
		}
		return readFile(t, filepath.Join(dir, "interop_mycrc.go"))
	}

	first := gen(t)
	if !strings.Contains(first, "-packages hash/crc32=mycrc") {
		t.Errorf("generated-by header does not record the non-default alias:\n%s",
			strings.SplitN(first, "\n", 2)[0])
	}
	if second := gen(t); first != second {
		t.Errorf("regeneration with identical inputs is not byte-identical "+
			"(%d vs %d bytes)", len(first), len(second))
	}
}

// TestLginteropSkipAccounting pins the success report for a scanned package
// with nothing to bind: it must count as skipped — not as generated output
// that was never written — and must not fail the run.
func TestLginteropSkipAccounting(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping skip-accounting check in -short mode (scans with go/types)")
	}
	root := repoRoot(t)
	trackLginteropInputs(t, root)
	dir := t.TempDir()

	cmd := exec.Command("go", "run", "./cmd/lginterop",
		"-packages", "github.com/nooga/let-go/internal/lginteroptest/empty,hash/crc32",
		"-out-pkg", "interop", "-out", dir)
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("a skipped package must not fail the run: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "generated 1/2") {
		t.Errorf("expected a 1/2 generation report:\n%s", out)
	}
	if !strings.Contains(string(out), "1 skipped: no eligible exports") {
		t.Errorf("skip not accounted for in the summary:\n%s", out)
	}
	if _, err := os.Stat(filepath.Join(dir, "interop_crc32.go")); err != nil {
		t.Errorf("eligible package's output missing: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "interop_empty.go")); err == nil {
		t.Errorf("skipped package unexpectedly produced a file")
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func copyFile(t *testing.T, src, dst string) {
	t.Helper()
	b, err := os.ReadFile(src)
	if err != nil {
		t.Fatalf("read %s: %v", src, err)
	}
	writeFile(t, dst, string(b))
}

// TestLginteropOutOfTreeRtAliasCollision guards the one import the
// out-of-tree mode adds unconditionally. A wrapped package whose own alias is
// `rt` — any import path ending in /rt — would land in the same import block
// as github.com/nooga/let-go/pkg/rt and the file would not compile
// ("rt redeclared"). The emitter must give the runtime a distinct alias.
//
// The equivalent collision for `vm` (always imported) and `fmt` (smart mode)
// is handled differently: those aliases are rejected up front (validateEntries)
// rather than aliased away, since no Go package path is forced to end in /vm
// the way rt-collisions arise. See TestValidateEntries in cmd/lginterop.
//
// deps.edn's alias form is the cheap way to force the collision: no stdlib
// package is named rt, so ask for hash/crc32 under that alias instead.
func TestLginteropOutOfTreeRtAliasCollision(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping rt-alias collision check in -short mode (scans with go/types)")
	}
	root := repoRoot(t)
	trackLginteropInputs(t, root)
	tmp := t.TempDir()
	genDir := filepath.Join(tmp, "interop")

	depsDir := filepath.Join(tmp, "deps")
	if err := os.MkdirAll(depsDir, 0755); err != nil {
		t.Fatalf("mkdir deps dir: %v", err)
	}
	writeFile(t, filepath.Join(depsDir, "deps.edn"), `{:gointerop [{"hash/crc32" "rt"}]}`)

	cmd := exec.Command("go", "run", "./cmd/lginterop",
		"-dir", depsDir, "-out-pkg", "interop", "-out", genDir)
	cmd.Dir = root
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("lginterop failed: %v\n%s", err, out)
	}

	src := readFile(t, filepath.Join(genDir, "interop_rt.go"))
	if !strings.Contains(src, `letgort "github.com/nooga/let-go/pkg/rt"`) {
		t.Errorf("runtime import not aliased away from the colliding package alias:\n%s", src)
	}
	if !strings.Contains(src, "letgort.RegisterNS(ns)") {
		t.Errorf("registration does not use the aliased runtime import:\n%s", src)
	}

	// Compiling is the real assertion: "rt redeclared" is a compile error.
	writeFile(t, filepath.Join(tmp, "go.mod"), `module example.com/rtcollision

go 1.26

require github.com/nooga/let-go v0.0.0

replace github.com/nooga/let-go => `+root+`
`)
	writeFile(t, filepath.Join(tmp, "main.go"), `package main

import _ "example.com/rtcollision/interop"

func main() {}
`)
	copyFile(t, filepath.Join(root, "go.sum"), filepath.Join(tmp, "go.sum"))

	build := exec.Command("go", "build", "./...")
	build.Dir = tmp
	build.Env = append(os.Environ(), "GOFLAGS=-mod=mod")
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("generated package with an rt-aliased target does not build: %v\n%s", err, out)
	}
}

// TestLginteropRejectsInvalidOutPkg checks the flag refuses names that would
// emit unparseable Go rather than reporting a successful generation.
func TestLginteropRejectsInvalidOutPkg(t *testing.T) {
	root := repoRoot(t)
	for _, bad := range []string{"rt", "main", "package", "my-pkg", "1interop"} {
		t.Run(bad, func(t *testing.T) {
			cmd := exec.Command("go", "run", "./cmd/lginterop",
				"-packages", "hash/crc32", "-out-pkg", bad, "-out", t.TempDir())
			cmd.Dir = root
			out, err := cmd.CombinedOutput()
			if err == nil {
				t.Fatalf("-out-pkg %q was accepted; output:\n%s", bad, out)
			}
			if !strings.Contains(string(out), "lginterop:") {
				t.Errorf("-out-pkg %q failed without a diagnostic:\n%s", bad, out)
			}
		})
	}
}

// TestLginteropFailsLoudly pins the exit status. Per-package errors only log
// and continue, so a run that generated nothing used to still exit 0 — a
// build pipeline driving this would see success and a missing file.
func TestLginteropFailsLoudly(t *testing.T) {
	root := repoRoot(t)
	out := t.TempDir()
	cmd := exec.Command("go", "run", "./cmd/lginterop",
		"-packages", "example.invalid/does/not/exist", "-out", out)
	cmd.Dir = root
	combined, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("scanning a nonexistent package exited 0; output:\n%s", combined)
	}
	if !strings.Contains(string(combined), "generated 0/1") {
		t.Errorf("expected a 0/1 generation report:\n%s", combined)
	}
	// The import failure should name the module-context requirement.
	if !strings.Contains(string(combined), "go get") {
		t.Errorf("import failure lacks the module-context hint:\n%s", combined)
	}
}

// TestLginteropRegisterStructOnlyForStructs pins which named types get
// vm.RegisterStruct. ex-struct? used to test the export vector's LENGTH, but
// every :type shape serializes to 6 elements — the non-struct branch writes
// `nil []` where the struct branch writes `:struct [fields]` — so it answered
// true for every named type. That emitted vm.RegisterStruct for named arrays
// like crc32.Table (`type Table [256]uint32`), which panics at init with
// "expected struct, got array". The in-tree xxh3 golden never caught it
// because it is generated with -opaque-structs, which skips RegisterStruct.
//
// Both directions matter, so assert both: a real struct still registers.
func TestLginteropRegisterStructOnlyForStructs(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping RegisterStruct check in -short mode (scans with go/types)")
	}
	root := repoRoot(t)
	trackLginteropInputs(t, root)

	gen := func(t *testing.T, pkg, file string) string {
		t.Helper()
		dir := t.TempDir()
		cmd := exec.Command("go", "run", "./cmd/lginterop",
			"-packages", pkg, "-out-pkg", "interop", "-out", dir)
		cmd.Dir = root
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("lginterop %s failed: %v\n%s", pkg, err, out)
		}
		return readFile(t, filepath.Join(dir, file))
	}

	// hash/crc32's only exported type is `type Table [256]uint32`.
	if src := gen(t, "hash/crc32", "interop_crc32.go"); strings.Contains(src, "RegisterStruct") {
		t.Errorf("named array type must not be registered as a struct:\n%s", src)
	}
	// image/color is all structs — these must keep registering.
	src := gen(t, "image/color", "interop_color.go")
	if !strings.Contains(src, `vm.RegisterStruct[color.RGBA]("color/RGBA")`) {
		t.Errorf("genuine struct type lost its RegisterStruct:\n%s", src)
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(b)
}
