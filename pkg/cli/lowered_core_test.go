/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package cli

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/genmanifest"
)

// fakeLGBGen stands in for the go run of lgbgen: it writes a package per
// namespace in pkgs and counts its runs.
func fakeLGBGen(t *testing.T, runs *int, pkgs ...string) {
	t.Helper()
	prev := lgbgenRunner
	t.Cleanup(func() { lgbgenRunner = prev })
	lgbgenRunner = func(_, prefix, out string, _ bool) error {
		*runs++
		for _, p := range pkgs {
			dir := filepath.Join(out, p)
			if err := os.MkdirAll(dir, 0755); err != nil {
				return err
			}
			src := "package " + p + "\n// lowered under " + prefix + "\n"
			if err := os.WriteFile(filepath.Join(dir, p+".go"), []byte(src), 0644); err != nil {
				return err
			}
		}
		return genmanifest.WriteTreeManifest(out)
	}
}

func TestLowerCoreIntoWritesTreeAndImports(t *testing.T) {
	noLetGoCache(t)
	var runs int
	fakeLGBGen(t, &runs, "core", "walk")
	dir := t.TempDir()
	if err := lowerCoreInto(dir, "example.com/app/lgprogram", true); err != nil {
		t.Fatal(err)
	}
	if runs != 1 {
		t.Fatalf("lgbgen ran %d times, want 1", runs)
	}
	imports, err := os.ReadFile(filepath.Join(dir, loweredCoreImportsFile))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"package main",
		`_ "example.com/app/lgprogram/lgcore/core"`,
		`_ "example.com/app/lgprogram/lgcore/walk"`,
	} {
		if !strings.Contains(string(imports), want) {
			t.Errorf("imports file lacks %s:\n%s", want, imports)
		}
	}
}

// A build that asked for the lowered core must not come out without it.
func TestLowerCoreIntoFailsWhenAPackageIsMissing(t *testing.T) {
	noLetGoCache(t)
	var runs int
	fakeLGBGen(t, &runs, "core") // no walk
	dir := t.TempDir()
	err := lowerCoreInto(dir, "lgprogram", false)
	if err == nil || !strings.Contains(err.Error(), "walk") {
		t.Fatalf("err = %v, want one naming walk", err)
	}
	if _, statErr := os.Stat(filepath.Join(dir, loweredCoreImportsFile)); !os.IsNotExist(statErr) {
		t.Fatal("the imports file was written for an incomplete tree")
	}
	// Nothing half-installed either: a leftover lgcore/ would make the next
	// build refuse, taking it for a program package.
	if entries, _ := os.ReadDir(dir); len(entries) != 0 {
		t.Fatalf("a failed lowering left %d entries in the build directory", len(entries))
	}
}

func TestLowerCoreIntoReportsLGBGenFailure(t *testing.T) {
	noLetGoCache(t)
	prev := lgbgenRunner
	t.Cleanup(func() { lgbgenRunner = prev })
	lgbgenRunner = func(_, _, _ string, _ bool) error { return errors.New("lowering exploded") }
	dir := t.TempDir()
	err := lowerCoreInto(dir, "lgprogram", false)
	if err == nil || !strings.Contains(err.Error(), "lowering exploded") {
		t.Fatalf("err = %v", err)
	}
	if _, statErr := os.Stat(filepath.Join(dir, loweredCoreImportsFile)); !os.IsNotExist(statErr) {
		t.Fatal("the imports file was written after a failed lowering")
	}
}

// noLetGoCache makes lowerCoreInto lower directly, as a build against a
// let-go source directory does.
func noLetGoCache(t *testing.T) {
	t.Helper()
	prev := letgoResolver
	t.Cleanup(func() { letgoResolver = prev })
	letgoResolver = func(string, bool) (string, bool) { return "", false }
}

// lgbgen runs with the module directory as its working directory, so the
// out dir it is handed has to be absolute even when -work is relative.
func TestLowerCoreIntoResolvesARelativeDir(t *testing.T) {
	noLetGoCache(t)
	prev := lgbgenRunner
	t.Cleanup(func() { lgbgenRunner = prev })
	var gotModuleDir, gotOut string
	lgbgenRunner = func(moduleDir, _, out string, _ bool) error {
		gotModuleDir, gotOut = moduleDir, out
		for _, p := range []string{"core", "walk"} {
			if err := os.MkdirAll(filepath.Join(out, p), 0755); err != nil {
				return err
			}
		}
		return nil
	}
	t.Chdir(t.TempDir())
	if err := lowerCoreInto("work", "lgprogram", false); err != nil {
		t.Fatal(err)
	}
	if !filepath.IsAbs(gotModuleDir) || !filepath.IsAbs(gotOut) {
		t.Fatalf("lgbgen got module dir %q and out %q; both must be absolute", gotModuleDir, gotOut)
	}
	if _, err := os.Stat(filepath.Join("work", "lgcore", "core")); err != nil {
		t.Fatalf("the tree is not under work/lgcore: %v", err)
	}
}

// A kept -work directory reused without the flag must not go on linking the
// core, and clearing must not take a program's own lgcore.* package with it.
func TestClearLoweredCore(t *testing.T) {
	noLetGoCache(t)
	var runs int
	fakeLGBGen(t, &runs, "core", "walk")
	dir := t.TempDir()
	if err := lowerCoreInto(dir, "lgprogram", false); err != nil {
		t.Fatal(err)
	}
	// An interrupted build's staging directory goes too.
	if err := os.MkdirAll(filepath.Join(dir, ".lgcore-stage-123", "tree", "core"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := clearLoweredCore(dir); err != nil {
		t.Fatal(err)
	}
	for _, gone := range []string{loweredCoreImportsFile, loweredCoreDirName, ".lgcore-stage-123"} {
		if _, err := os.Stat(filepath.Join(dir, gone)); !os.IsNotExist(err) {
			t.Errorf("%s survived clearLoweredCore", gone)
		}
	}

	prog := filepath.Join(dir, loweredCoreDirName, "core", "core.go")
	if err := os.MkdirAll(filepath.Dir(prog), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(prog, []byte("package core\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := clearLoweredCore(dir); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(prog); err != nil {
		t.Fatalf("clearLoweredCore removed a program package: %v", err)
	}
	if err := clearLoweredCore(filepath.Join(dir, "never-created")); err != nil {
		t.Fatalf("clearing a directory that does not exist: %v", err)
	}
}

// A program namespace lgcore.* lowers to the directory the core needs. The
// build has to stop, not replace the program's package.
func TestLowerCoreIntoRefusesAProgramPackageInItsWay(t *testing.T) {
	noLetGoCache(t)
	var runs int
	fakeLGBGen(t, &runs, "core", "walk")
	dir := t.TempDir()
	prog := filepath.Join(dir, loweredCoreDirName, "core", "core.go")
	if err := os.MkdirAll(filepath.Dir(prog), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(prog, []byte("package core // the program's\n"), 0644); err != nil {
		t.Fatal(err)
	}
	err := lowerCoreInto(dir, "lgprogram", false)
	if err == nil || !strings.Contains(err.Error(), "a program package is already there") {
		t.Fatalf("err = %v, want a refusal naming the program's package", err)
	}
	if runs != 0 {
		t.Fatalf("lgbgen ran %d times before the refusal", runs)
	}
	if src, _ := os.ReadFile(prog); !strings.Contains(string(src), "the program's") {
		t.Fatal("the program's package was replaced")
	}
}

// The cache is keyed by the let-go the module being built links, which with
// -module is the caller's go.mod's choice and not this lg's own version.
func TestLoweredCoreCacheFollowsTheModulesLetGo(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("XDG_CACHE_HOME", t.TempDir())
	prev := letgoResolver
	t.Cleanup(func() { letgoResolver = prev })
	ref := "github.com/nooga/let-go@v1.14.0"
	letgoResolver = func(string, bool) (string, bool) { return ref, true }
	var runs int
	fakeLGBGen(t, &runs, "core", "walk")

	build := func() {
		t.Helper()
		if err := lowerCoreInto(t.TempDir(), "lgprogram", false); err != nil {
			t.Fatal(err)
		}
	}
	build()
	build()
	if runs != 1 {
		t.Fatalf("two builds against one let-go version ran lgbgen %d times, want 1", runs)
	}
	ref = "github.com/nooga/let-go@v1.14.1"
	build()
	if runs != 2 {
		t.Fatalf("a build against another let-go version reused the cache (lgbgen runs = %d)", runs)
	}
	// A cache entry that lost files (a cache cleaner, an interrupted copy) no
	// longer matches its manifest and is lowered again, not reused.
	cache, ok := loweredCoreCacheDir("", "lgprogram/lgcore", false)
	if !ok {
		t.Fatal("no cache entry for a resolved let-go version")
	}
	if err := os.Remove(filepath.Join(cache, "walk", "walk.go")); err != nil {
		t.Fatal(err)
	}
	build()
	if runs != 3 {
		t.Fatalf("a damaged cache entry was reused (lgbgen runs = %d, want 3)", runs)
	}
	build()
	if runs != 3 {
		t.Fatalf("the repaired cache entry was not reused (lgbgen runs = %d, want 3)", runs)
	}
	letgoResolver = func(string, bool) (string, bool) { return "", false }
	build()
	build()
	if runs != 5 {
		t.Fatalf("builds against a let-go source directory must lower every time (lgbgen runs = %d, want 5)", runs)
	}
}

// lgbgen runs on this machine whatever the program is being built for.
func TestHostGoEnvOverridesACrossCompileTarget(t *testing.T) {
	env := hostGoEnv([]string{"GOOS=plan9", "GOARCH=386", "HOME=/h"})
	last := map[string]string{}
	for _, kv := range env {
		k, v, _ := strings.Cut(kv, "=")
		last[k] = v // exec.Cmd uses the last value of a duplicated key
	}
	if last["GOOS"] != runtime.GOOS || last["GOARCH"] != runtime.GOARCH || last["HOME"] != "/h" {
		t.Fatalf("env = %v; want GOOS=%s GOARCH=%s and the rest kept", env, runtime.GOOS, runtime.GOARCH)
	}
}

func TestLetGoRefFromGoList(t *testing.T) {
	for _, tc := range []struct {
		name, json, want string
		ok               bool
	}{
		{"a required version", `{"Path":"github.com/nooga/let-go","Version":"v1.14.0"}`, "github.com/nooga/let-go@v1.14.0", true},
		{"a pseudo-version", `{"Path":"github.com/nooga/let-go","Version":"v1.13.1-0.20261113224131-c0a93f0dbac1"}`, "github.com/nooga/let-go@v1.13.1-0.20261113224131-c0a93f0dbac1", true},
		{"replaced by a fork's version", `{"Path":"github.com/nooga/let-go","Version":"v1.14.0","Replace":{"Path":"github.com/x/let-go","Version":"v0.0.1"}}`, "github.com/x/let-go@v0.0.1", true},
		{"replaced by a directory", `{"Path":"github.com/nooga/let-go","Version":"v1.14.0","Replace":{"Path":"../let-go","Dir":"/src/let-go"}}`, "", false},
		{"the main module itself", `{"Path":"github.com/nooga/let-go","Main":true,"Dir":"/src/let-go"}`, "", false},
		{"not json", `go: errors`, "", false},
	} {
		got, ok := letgoRefFrom([]byte(tc.json))
		if got != tc.want || ok != tc.ok {
			t.Errorf("%s: got %q, %v; want %q, %v", tc.name, got, ok, tc.want, tc.ok)
		}
	}
	a := loweredCoreCacheKey("github.com/nooga/let-go@v1.14.0", "lgprogram/lgcore")
	if b := loweredCoreCacheKey("github.com/nooga/let-go@v1.14.0", "example.com/app/lgprogram/lgcore"); a == b {
		t.Error("two import prefixes share a cache entry")
	}
}
