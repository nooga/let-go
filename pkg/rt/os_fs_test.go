//go:build !tinygo

/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

// osFn resolves a registered native from the os namespace.
func osFn(t *testing.T, name string) vm.Fn {
	t.Helper()
	v := NS("os").Lookup(vm.Symbol(name))
	if v == nil || v == vm.NIL {
		t.Fatalf("os/%s not found", name)
	}
	if vr, ok := v.(*vm.Var); ok {
		v = vr.Deref()
	}
	fn, ok := v.(vm.Fn)
	if !ok {
		t.Fatalf("os/%s is not an Fn, got %T", name, v)
	}
	return fn
}

func callOsFn(t *testing.T, name string, args ...string) (vm.Value, error) {
	t.Helper()
	vs := make([]vm.Value, len(args))
	for i, a := range args {
		vs[i] = vm.String(a)
	}
	return osFn(t, name).Invoke(vs)
}

func writeFileAt(t *testing.T, path, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir for %s: %v", path, err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func TestOsRenameMovesAFile(t *testing.T) {
	tmp := t.TempDir()
	from := filepath.Join(tmp, "before.txt")
	to := filepath.Join(tmp, "after.txt")
	writeFileAt(t, from, "payload")

	got, err := callOsFn(t, "rename", from, to)
	if err != nil {
		t.Fatalf("rename: %v", err)
	}
	if want := vm.String(to); got != want {
		t.Errorf("rename returned %v, want %v", got, want)
	}
	if body := readFile(t, to); body != "payload" {
		t.Errorf("destination holds %q, want %q", body, "payload")
	}
	if _, err := os.Stat(from); !os.IsNotExist(err) {
		t.Errorf("source still present after rename, stat err = %v", err)
	}
}

// The publish-by-rename pattern is the reason this native exists. This
// checks the end state only — a torn intermediate is not observable from
// here — so it pins that the destination ends up holding the new content
// with no leftover source.
func TestOsRenameReplacesAnExistingFile(t *testing.T) {
	tmp := t.TempDir()
	staged := filepath.Join(tmp, "staged")
	live := filepath.Join(tmp, "live")
	writeFileAt(t, staged, "new")
	writeFileAt(t, live, "old")

	if _, err := callOsFn(t, "rename", staged, live); err != nil {
		t.Fatalf("rename over existing: %v", err)
	}
	if body := readFile(t, live); body != "new" {
		t.Errorf("live holds %q, want %q", body, "new")
	}
}

func TestOsRenameMovesADirectoryWithContents(t *testing.T) {
	tmp := t.TempDir()
	from := filepath.Join(tmp, "tree")
	to := filepath.Join(tmp, "moved")
	writeFileAt(t, filepath.Join(from, "nested", "leaf.txt"), "leaf")

	if _, err := callOsFn(t, "rename", from, to); err != nil {
		t.Fatalf("rename dir: %v", err)
	}
	if body := readFile(t, filepath.Join(to, "nested", "leaf.txt")); body != "leaf" {
		t.Errorf("moved leaf holds %q, want %q", body, "leaf")
	}
}

func TestOsRenameReportsAMissingSource(t *testing.T) {
	tmp := t.TempDir()
	_, err := callOsFn(t, "rename", filepath.Join(tmp, "absent"), filepath.Join(tmp, "dest"))
	if err == nil {
		t.Fatal("renaming a missing source succeeded, want an error")
	}
}

func TestOsRenameRejectsBadArgs(t *testing.T) {
	tmp := t.TempDir()
	if _, err := osFn(t, "rename").Invoke(nil); err == nil {
		t.Error("rename with no args succeeded, want an error")
	}
	if _, err := callOsFn(t, "rename", filepath.Join(tmp, "a")); err == nil {
		t.Error("rename with 1 arg succeeded, want an error")
	}
	if _, err := osFn(t, "rename").Invoke([]vm.Value{vm.Int(1), vm.String(tmp)}); err == nil {
		t.Error("rename with a non-String source succeeded, want an error")
	}
	if _, err := osFn(t, "rename").Invoke([]vm.Value{vm.String(tmp), vm.Int(1)}); err == nil {
		t.Error("rename with a non-String destination succeeded, want an error")
	}
}

func TestOsDeleteTreeRemovesRecursively(t *testing.T) {
	tmp := t.TempDir()
	root := filepath.Join(tmp, "root")
	writeFileAt(t, filepath.Join(root, "top.txt"), "top")
	writeFileAt(t, filepath.Join(root, "a", "b", "deep.txt"), "deep")

	got, err := callOsFn(t, "delete-tree", root)
	if err != nil {
		t.Fatalf("delete-tree: %v", err)
	}
	if got != vm.NIL {
		t.Errorf("delete-tree returned %v, want nil", got)
	}
	if _, err := os.Stat(root); !os.IsNotExist(err) {
		t.Errorf("tree still present, stat err = %v", err)
	}
}

func TestOsDeleteTreeRemovesASingleFile(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "lone.txt")
	writeFileAt(t, path, "x")

	if _, err := callOsFn(t, "delete-tree", path); err != nil {
		t.Fatalf("delete-tree on a file: %v", err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Errorf("file still present, stat err = %v", err)
	}
}

// Absent-is-success is the documented divergence from delete-file, which
// fails on a missing path. Pinned so it can't drift back silently.
func TestOsDeleteTreeSucceedsWhenAlreadyAbsent(t *testing.T) {
	tmp := t.TempDir()
	if _, err := callOsFn(t, "delete-tree", filepath.Join(tmp, "never-existed")); err != nil {
		t.Errorf("delete-tree on an absent path returned %v, want nil", err)
	}
}

// The containment property worth pinning: a symlink inside the doomed tree
// pointing out of it must be unlinked, not followed. Two unrelated sibling
// directories would prove nothing, since no implementation could conflate
// them.
func TestOsDeleteTreeUnlinksSymlinksRatherThanFollowingThem(t *testing.T) {
	tmp := t.TempDir()
	doomed := filepath.Join(tmp, "doomed")
	outside := filepath.Join(tmp, "outside")
	writeFileAt(t, filepath.Join(doomed, "x.txt"), "x")
	writeFileAt(t, filepath.Join(outside, "keep.txt"), "keep")
	if err := os.Symlink(outside, filepath.Join(doomed, "escape")); err != nil {
		t.Skipf("symlink unsupported here: %v", err)
	}

	if _, err := callOsFn(t, "delete-tree", doomed); err != nil {
		t.Fatalf("delete-tree: %v", err)
	}
	if _, err := os.Stat(doomed); !os.IsNotExist(err) {
		t.Errorf("doomed tree survived, stat err = %v", err)
	}
	if body := readFile(t, filepath.Join(outside, "keep.txt")); body != "keep" {
		t.Errorf("the symlink was followed out of the tree; outside holds %q", body)
	}
}

// The corollary, and the surprising half: when the path itself is a symlink
// to a directory, only the link goes. Documented on the native because
// "removes path and everything beneath it" does not lead a caller to expect
// it.
func TestOsDeleteTreeOnASymlinkRemovesOnlyTheLink(t *testing.T) {
	tmp := t.TempDir()
	target := filepath.Join(tmp, "target")
	writeFileAt(t, filepath.Join(target, "t.txt"), "t")
	link := filepath.Join(tmp, "linkdir")
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("symlink unsupported here: %v", err)
	}

	if _, err := callOsFn(t, "delete-tree", link); err != nil {
		t.Fatalf("delete-tree: %v", err)
	}
	if _, err := os.Lstat(link); !os.IsNotExist(err) {
		t.Errorf("link survived, lstat err = %v", err)
	}
	if body := readFile(t, filepath.Join(target, "t.txt")); body != "t" {
		t.Errorf("the link's target lost its contents; holds %q", body)
	}
}

// The empty-path guard exists to catch an unset variable. "/" is the same
// mistake with a worse outcome, and in lg (str nil) is "", so "$root/$name"
// with root unset yields "/name" rather than "".
func TestOsDeleteTreeRefusesTheFilesystemRoot(t *testing.T) {
	for _, p := range []string{"/", "//", "/.."} {
		if _, err := callOsFn(t, "delete-tree", p); err == nil {
			t.Errorf("delete-tree(%q) succeeded, want an error", p)
		}
	}
}

// RemoveAll("") is a silent no-op in Go. Returning an error instead turns
// an unset path variable into a failure the caller can see.
func TestOsDeleteTreeRefusesAnEmptyPath(t *testing.T) {
	if _, err := callOsFn(t, "delete-tree", ""); err == nil {
		t.Error("delete-tree on an empty path succeeded, want an error")
	}
}

func TestOsDeleteTreeRejectsBadArgs(t *testing.T) {
	if _, err := osFn(t, "delete-tree").Invoke(nil); err == nil {
		t.Error("delete-tree with no args succeeded, want an error")
	}
	if _, err := osFn(t, "delete-tree").Invoke([]vm.Value{vm.Int(1)}); err == nil {
		t.Error("delete-tree with a non-String path succeeded, want an error")
	}
}
