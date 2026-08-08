//go:build !tinygo

/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

func TestOsAbsolutePathResolvesAgainstCwd(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	got, err := callOsFn(t, "absolute-path", "relative/leaf.txt")
	if err != nil {
		t.Fatalf("absolute-path: %v", err)
	}
	if want := vm.String(filepath.Join(cwd, "relative/leaf.txt")); got != want {
		t.Errorf("absolute-path = %v, want %v", got, want)
	}
}

func TestOsAbsolutePathCleansTraversal(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	got, err := callOsFn(t, "absolute-path", "a/b/../c")
	if err != nil {
		t.Fatalf("absolute-path: %v", err)
	}
	if want := vm.String(filepath.Join(cwd, "a/c")); got != want {
		t.Errorf("absolute-path = %v, want %v", got, want)
	}
}

func TestOsAbsolutePathLeavesAnAbsolutePathAlone(t *testing.T) {
	tmp := t.TempDir()
	got, err := callOsFn(t, "absolute-path", tmp)
	if err != nil {
		t.Fatalf("absolute-path: %v", err)
	}
	if want := vm.String(tmp); got != want {
		t.Errorf("absolute-path = %v, want %v", got, want)
	}
}

// absolute-path is lexical, so it must not require the path to exist.
func TestOsAbsolutePathAcceptsAMissingPath(t *testing.T) {
	tmp := t.TempDir()
	if _, err := callOsFn(t, "absolute-path", filepath.Join(tmp, "not-created-yet")); err != nil {
		t.Errorf("absolute-path on a missing path returned %v, want nil", err)
	}
}

func TestOsCanonicalPathResolvesASymlink(t *testing.T) {
	tmp := t.TempDir()
	target := filepath.Join(tmp, "target.txt")
	writeFileAt(t, target, "x")
	link := filepath.Join(tmp, "link.txt")
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("symlink unsupported here: %v", err)
	}

	got, err := callOsFn(t, "canonical-path", link)
	if err != nil {
		t.Fatalf("canonical-path: %v", err)
	}
	// t.TempDir can itself sit under a symlink (/var → /private/var on
	// darwin), so compare against the resolved target rather than target.
	wantReal, err := filepath.EvalSymlinks(target)
	if err != nil {
		t.Fatalf("evalsymlinks: %v", err)
	}
	if want := vm.String(wantReal); got != want {
		t.Errorf("canonical-path = %v, want %v", got, want)
	}
}

// Two names for one file must canonicalize to one string — the property
// that makes this the form to compare or key on.
func TestOsCanonicalPathCollapsesTwoNamesForOneFile(t *testing.T) {
	tmp := t.TempDir()
	target := filepath.Join(tmp, "real.txt")
	writeFileAt(t, target, "x")
	link := filepath.Join(tmp, "alias.txt")
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("symlink unsupported here: %v", err)
	}

	viaLink, err := callOsFn(t, "canonical-path", link)
	if err != nil {
		t.Fatalf("canonical-path via link: %v", err)
	}
	viaReal, err := callOsFn(t, "canonical-path", target)
	if err != nil {
		t.Fatalf("canonical-path via target: %v", err)
	}
	if viaLink != viaReal {
		t.Errorf("link canonicalized to %v, target to %v; want one string", viaLink, viaReal)
	}
}

func TestOsCanonicalPathAcceptsARelativePath(t *testing.T) {
	got, err := callOsFn(t, "canonical-path", ".")
	if err != nil {
		t.Fatalf("canonical-path: %v", err)
	}
	if s, ok := got.(vm.String); !ok || !strings.HasPrefix(string(s), "/") {
		t.Errorf("canonical-path(\".\") = %v, want an absolute path", got)
	}
}

// The documented divergence from absolute-path: canonicalizing requires
// looking at the filesystem, so a missing path has no answer.
func TestOsCanonicalPathReportsAMissingPath(t *testing.T) {
	tmp := t.TempDir()
	if _, err := callOsFn(t, "canonical-path", filepath.Join(tmp, "absent")); err == nil {
		t.Error("canonical-path on a missing path succeeded, want an error")
	}
}

func TestOsPathFnsRejectBadArgs(t *testing.T) {
	for _, name := range []string{"absolute-path", "canonical-path"} {
		if _, err := osFn(t, name).Invoke(nil); err == nil {
			t.Errorf("os/%s with no args succeeded, want an error", name)
		}
		if _, err := osFn(t, name).Invoke([]vm.Value{vm.Int(1)}); err == nil {
			t.Errorf("os/%s with a non-String path succeeded, want an error", name)
		}
	}
}
