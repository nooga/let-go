/*
 * Copyright (c) 2026 let-go contributors; see CONTRIBUTORS.
 * SPDX-License-Identifier: MIT
 */

package gomod

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestReleaseRef pins the released-vs-dev decision that routes Generate. A
// released binary must pin its require to its own version (nooga/let-go#594:
// an unpinned require made released `lg -w` builds resolve a go.sum that did
// not match the runtime they were emitted by); "dev" and an empty version must
// not, so they reach the local-source path instead.
func TestReleaseRef(t *testing.T) {
	for _, tc := range []struct {
		version string
		wantRef string
		wantOK  bool
	}{
		{"1.12.2", "github.com/nooga/let-go@v1.12.2", true},
		{"0.1.0", "github.com/nooga/let-go@v0.1.0", true},
		{"dev", "", false},
		{"", "", false},
	} {
		ref, ok := ReleaseRef(tc.version)
		if ok != tc.wantOK || ref != tc.wantRef {
			t.Errorf("ReleaseRef(%q) = (%q, %v), want (%q, %v)", tc.version, ref, ok, tc.wantRef, tc.wantOK)
		}
	}
}

// writeFakeCheckout lays down the minimum that reads as a let-go source tree:
// a go.mod naming the module, plus a go.sum to be carried across.
func writeFakeCheckout(t *testing.T, goSum string) string {
	t.Helper()
	dir := t.TempDir()
	mod := "module " + ModulePath + "\n\ngo 1.26\n\nrequire example.com/dep v1.2.3\n"
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(mod), 0644); err != nil {
		t.Fatal(err)
	}
	if goSum != "" {
		if err := os.WriteFile(filepath.Join(dir, "go.sum"), []byte(goSum), 0644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

// TestGenerateLocalSource covers the dev-build path end to end through
// Generate: LETGO_SRC points at a checkout, so no network or `go get` is
// involved. It asserts the two properties the generated module must have —
// the caller's module name (not let-go's, and not a baked-in one), and a
// replace pointing back at the checkout.
func TestGenerateLocalSource(t *testing.T) {
	src := writeFakeCheckout(t, "example.com/dep v1.2.3 h1:abc=\n")
	t.Setenv("LETGO_SRC", src)

	files, err := Generate(t.TempDir(), "some-app", "dev")
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}

	if !strings.HasPrefix(files.Mod, "module some-app\n") {
		t.Errorf("module name not substituted; go.mod starts:\n%s", firstLine(files.Mod))
	}
	if strings.Contains(files.Mod, "module "+ModulePath) {
		t.Errorf("let-go's own module line survived into the generated go.mod:\n%s", files.Mod)
	}
	if want := "replace " + ModulePath + " => " + src; !strings.Contains(files.Mod, want) {
		t.Errorf("missing %q in:\n%s", want, files.Mod)
	}
	if !strings.Contains(files.Mod, "require "+ModulePath+" v0.0.0") {
		t.Errorf("missing require of let-go in:\n%s", files.Mod)
	}
	if string(files.Sum) != "example.com/dep v1.2.3 h1:abc=\n" {
		t.Errorf("go.sum not carried across, got %q", files.Sum)
	}
}

// TestGenerateModuleNameIsCallerSupplied is the point of the extraction: the
// module name was baked in when the WASM build was the only caller, so two
// callers must be able to get two different names from one implementation.
func TestGenerateModuleNameIsCallerSupplied(t *testing.T) {
	src := writeFakeCheckout(t, "")
	t.Setenv("LETGO_SRC", src)

	for _, name := range []string{"lg-wasm-app", "aotmod", "example.com/native/prog"} {
		files, err := Generate(t.TempDir(), name, "dev")
		if err != nil {
			t.Fatalf("Generate(%q): %v", name, err)
		}
		if !strings.HasPrefix(files.Mod, "module "+name+"\n") {
			t.Errorf("module %q not honored; go.mod starts:\n%s", name, firstLine(files.Mod))
		}
	}
}

// TestFindLetGoSourceDirPrefersEnv pins LETGO_SRC as the override: it wins over
// the walk-up, which matters because the walk-up finds whatever checkout the
// process happens to be running inside.
func TestFindLetGoSourceDirPrefersEnv(t *testing.T) {
	src := writeFakeCheckout(t, "")
	t.Setenv("LETGO_SRC", src)

	got, err := FindLetGoSourceDir()
	if err != nil {
		t.Fatalf("FindLetGoSourceDir: %v", err)
	}
	if got != src {
		t.Errorf("FindLetGoSourceDir() = %q, want %q", got, src)
	}
}

// TestGoToolPath checks the returned tool is usable: either an existing path
// under GOROOT, or the bare "go" fallback for PATH lookup.
func TestGoToolPath(t *testing.T) {
	got := GoToolPath()
	if got == "go" {
		return
	}
	if _, err := os.Stat(got); err != nil {
		t.Errorf("GoToolPath() = %q, which does not exist: %v", got, err)
	}
}

func firstLine(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return s[:i]
	}
	return s
}
