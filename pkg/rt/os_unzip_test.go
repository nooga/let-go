//go:build !tinygo

/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"archive/zip"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

// zipEntry describes one member of a zip built by writeZip.
type zipEntry struct {
	name string
	body string
	mode fs.FileMode // zero means "no unix mode recorded"
	dir  bool
}

// writeZip builds a zip at dir/name from entries and returns its path.
func writeZip(t *testing.T, dir, name string, entries []zipEntry) string {
	t.Helper()
	path := filepath.Join(dir, name)
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("create zip: %v", err)
	}
	defer f.Close()

	zw := zip.NewWriter(f)
	for _, e := range entries {
		hdr := &zip.FileHeader{Name: e.name, Method: zip.Deflate}
		if e.mode != 0 {
			hdr.SetMode(e.mode)
		}
		if e.dir && !strings.HasSuffix(hdr.Name, "/") {
			hdr.Name += "/"
		}
		w, err := zw.CreateHeader(hdr)
		if err != nil {
			t.Fatalf("create entry %q: %v", e.name, err)
		}
		if !e.dir {
			if _, err := w.Write([]byte(e.body)); err != nil {
				t.Fatalf("write entry %q: %v", e.name, err)
			}
		}
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("close zip: %v", err)
	}
	return path
}

// unzipFn resolves the registered os/unzip native.
func unzipFn(t *testing.T) vm.Fn {
	t.Helper()
	v := NS("os").Lookup(vm.Symbol("unzip"))
	if v == nil || v == vm.NIL {
		t.Fatal("os/unzip not found")
	}
	if vr, ok := v.(*vm.Var); ok {
		v = vr.Deref()
	}
	fn, ok := v.(vm.Fn)
	if !ok {
		t.Fatalf("os/unzip is not an Fn, got %T", v)
	}
	return fn
}

func callUnzip(t *testing.T, src, dest string) (vm.Value, error) {
	t.Helper()
	return unzipFn(t).Invoke([]vm.Value{vm.String(src), vm.String(dest)})
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(b)
}

func TestOsUnzipExtractsFilesAndDirs(t *testing.T) {
	tmp := t.TempDir()
	src := writeZip(t, tmp, "a.zip", []zipEntry{
		{name: "top.txt", body: "top"},
		{name: "nested/deep/leaf.txt", body: "leaf"},
		{name: "explicit", dir: true, mode: 0o755 | fs.ModeDir},
		{name: "explicit/inside.txt", body: "inside"},
	})
	dest := filepath.Join(tmp, "out")

	got, err := callUnzip(t, src, dest)
	if err != nil {
		t.Fatalf("unzip: %v", err)
	}
	if want := vm.String(dest); got != want {
		t.Errorf("return = %#v, want %#v", got, want)
	}

	for path, want := range map[string]string{
		filepath.Join(dest, "top.txt"):                    "top",
		filepath.Join(dest, "nested", "deep", "leaf.txt"): "leaf",
		filepath.Join(dest, "explicit", "inside.txt"):     "inside",
	} {
		if got := readFile(t, path); got != want {
			t.Errorf("%s = %q, want %q", path, got, want)
		}
	}

	// Implicit parent directories are created for nested entries.
	info, err := os.Stat(filepath.Join(dest, "nested", "deep"))
	if err != nil {
		t.Fatalf("stat nested dir: %v", err)
	}
	if !info.IsDir() {
		t.Errorf("nested/deep is not a directory")
	}
}

func TestOsUnzipCreatesMissingDestination(t *testing.T) {
	tmp := t.TempDir()
	src := writeZip(t, tmp, "a.zip", []zipEntry{{name: "f.txt", body: "x"}})
	dest := filepath.Join(tmp, "does", "not", "exist")

	if _, err := callUnzip(t, src, dest); err != nil {
		t.Fatalf("unzip: %v", err)
	}
	if got := readFile(t, filepath.Join(dest, "f.txt")); got != "x" {
		t.Errorf("f.txt = %q, want %q", got, "x")
	}
}

func TestOsUnzipOverwritesExistingFile(t *testing.T) {
	tmp := t.TempDir()
	src := writeZip(t, tmp, "a.zip", []zipEntry{{name: "f.txt", body: "new"}})
	dest := filepath.Join(tmp, "out")
	if err := os.MkdirAll(dest, 0o755); err != nil {
		t.Fatalf("mkdir dest: %v", err)
	}
	// Pre-existing content is longer than the replacement — a truncating
	// write is required, not just an overwrite of the first bytes.
	if err := os.WriteFile(filepath.Join(dest, "f.txt"), []byte("stale stale stale"), 0o644); err != nil {
		t.Fatalf("seed file: %v", err)
	}

	if _, err := callUnzip(t, src, dest); err != nil {
		t.Fatalf("unzip: %v", err)
	}
	if got := readFile(t, filepath.Join(dest, "f.txt")); got != "new" {
		t.Errorf("f.txt = %q, want %q", got, "new")
	}
}

func TestOsUnzipAppliesEntryMode(t *testing.T) {
	tmp := t.TempDir()
	src := writeZip(t, tmp, "a.zip", []zipEntry{
		{name: "script.sh", body: "#!/bin/sh\n", mode: 0o755},
		{name: "plain.txt", body: "p"}, // no recorded mode -> default 0644
	})
	dest := filepath.Join(tmp, "out")
	if _, err := callUnzip(t, src, dest); err != nil {
		t.Fatalf("unzip: %v", err)
	}

	// Assert the bits that matter rather than the exact mode: os.OpenFile
	// applies the process umask, so an exact comparison would fail under an
	// unusual one.
	info, err := os.Stat(filepath.Join(dest, "script.sh"))
	if err != nil {
		t.Fatalf("stat script: %v", err)
	}
	if got := info.Mode().Perm(); got&0o100 == 0 {
		t.Errorf("script.sh mode = %v, want the owner-execute bit set", got)
	}

	info, err = os.Stat(filepath.Join(dest, "plain.txt"))
	if err != nil {
		t.Fatalf("stat plain: %v", err)
	}
	perm := info.Mode().Perm()
	if perm&0o111 != 0 {
		t.Errorf("plain.txt mode = %v, want no execute bits", perm)
	}
	// The FAT-style entry records no unix mode; archive/zip would synthesize
	// 0666, which is world-writable when the umask is permissive.
	if perm&0o022 != 0 {
		t.Errorf("plain.txt mode = %v, want no group/other write bits", perm)
	}
}

func TestOsUnzipRelativeDestination(t *testing.T) {
	tmp := t.TempDir()
	// Two entries under one nested directory: the second is extracted after
	// the directory already exists, which is where a containment check that
	// compares unnormalised paths goes wrong for a relative dest.
	src := writeZip(t, tmp, "a.zip", []zipEntry{
		{name: "sub/one.txt", body: "1"},
		{name: "sub/two.txt", body: "2"},
	})
	work := filepath.Join(tmp, "work")
	if err := os.MkdirAll(work, 0o755); err != nil {
		t.Fatalf("mkdir work: %v", err)
	}
	t.Chdir(work)

	if _, err := callUnzip(t, src, "."); err != nil {
		t.Fatalf("unzip into .: %v", err)
	}
	for name, want := range map[string]string{"one.txt": "1", "two.txt": "2"} {
		if got := readFile(t, filepath.Join(work, "sub", name)); got != want {
			t.Errorf("sub/%s = %q, want %q", name, got, want)
		}
	}
}

func TestOsUnzipOverwriteAppliesEntryMode(t *testing.T) {
	tmp := t.TempDir()
	src := writeZip(t, tmp, "a.zip", []zipEntry{{name: "run.sh", body: "#!/bin/sh\n", mode: 0o755}})
	dest := filepath.Join(tmp, "out")
	if err := os.MkdirAll(dest, 0o755); err != nil {
		t.Fatalf("mkdir dest: %v", err)
	}
	// O_CREATE applies its perm only when it creates the file, so overwriting
	// in place would leave this one non-executable.
	if err := os.WriteFile(filepath.Join(dest, "run.sh"), []byte("old"), 0o644); err != nil {
		t.Fatalf("seed file: %v", err)
	}

	if _, err := callUnzip(t, src, dest); err != nil {
		t.Fatalf("unzip: %v", err)
	}
	info, err := os.Stat(filepath.Join(dest, "run.sh"))
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if got := info.Mode().Perm(); got&0o100 == 0 {
		t.Errorf("run.sh mode = %v, want the owner-execute bit set after overwrite", got)
	}
}

func TestOsUnzipAppliesDirectoryEntryMode(t *testing.T) {
	tmp := t.TempDir()
	src := writeZip(t, tmp, "a.zip", []zipEntry{
		{name: "private", dir: true, mode: 0o700 | fs.ModeDir},
		{name: "private/secret.txt", body: "s"},
		{name: "plain", dir: true},
	})
	dest := filepath.Join(tmp, "out")
	if _, err := callUnzip(t, src, dest); err != nil {
		t.Fatalf("unzip: %v", err)
	}

	info, err := os.Stat(filepath.Join(dest, "private"))
	if err != nil {
		t.Fatalf("stat private: %v", err)
	}
	// A 0700 entry must not widen to 0755 just because the umask is permissive.
	if got := info.Mode().Perm(); got&0o077 != 0 {
		t.Errorf("private mode = %v, want no group/other bits", got)
	}
	// The entry's own contents still extract — owner rwx is forced on.
	if got := readFile(t, filepath.Join(dest, "private", "secret.txt")); got != "s" {
		t.Errorf("secret.txt = %q, want %q", got, "s")
	}

	info, err = os.Stat(filepath.Join(dest, "plain"))
	if err != nil {
		t.Fatalf("stat plain: %v", err)
	}
	if !info.IsDir() {
		t.Errorf("plain is not a directory")
	}
}

func TestOsUnzipRejectsZipSlip(t *testing.T) {
	tmp := t.TempDir()
	src := writeZip(t, tmp, "evil.zip", []zipEntry{
		{name: "ok.txt", body: "ok"},
		{name: "../evil.txt", body: "pwned"},
	})
	dest := filepath.Join(tmp, "out")

	if _, err := callUnzip(t, src, dest); err == nil {
		t.Fatal("expected an error for a ../ entry, got nil")
	}
	if _, err := os.Stat(filepath.Join(tmp, "evil.txt")); !os.IsNotExist(err) {
		t.Errorf("escaped file was created outside dest (stat err = %v)", err)
	}
}

func TestOsUnzipRejectsAbsoluteEntryEscape(t *testing.T) {
	tmp := t.TempDir()
	outside := filepath.Join(tmp, "outside.txt")
	// A deep ../ chain: even after Join cleans it, the result must not escape.
	src := writeZip(t, tmp, "evil.zip", []zipEntry{
		{name: "../../../../../../../../etc/passwd-lg-test", body: "pwned"},
		{name: "sub/../../outside.txt", body: "pwned"},
	})
	dest := filepath.Join(tmp, "out")

	if _, err := callUnzip(t, src, dest); err == nil {
		t.Fatal("expected an error for escaping entries, got nil")
	}
	if _, err := os.Stat(outside); !os.IsNotExist(err) {
		t.Errorf("escaped file was created outside dest (stat err = %v)", err)
	}
}

func TestOsUnzipRejectsAbsoluteEntryName(t *testing.T) {
	tmp := t.TempDir()
	src := writeZip(t, tmp, "evil.zip", []zipEntry{{name: "/etc/passwd-lg-test", body: "pwned"}})
	dest := filepath.Join(tmp, "out")

	// An absolute name is refused outright rather than quietly re-rooted
	// inside dest — the archive asked for something it cannot have.
	if _, err := callUnzip(t, src, dest); err == nil {
		t.Fatal("expected an error for an absolute entry name, got nil")
	}
	if _, err := os.Stat(filepath.Join(dest, "etc", "passwd-lg-test")); !os.IsNotExist(err) {
		t.Errorf("absolute entry was re-rooted into dest (stat err = %v)", err)
	}
}

func TestOsUnzipRejectsEscapeThroughPreexistingSymlink(t *testing.T) {
	tmp := t.TempDir()
	target := filepath.Join(tmp, "target")
	if err := os.MkdirAll(target, 0o755); err != nil {
		t.Fatalf("mkdir target: %v", err)
	}
	victim := filepath.Join(target, "victim.txt")
	if err := os.WriteFile(victim, []byte("untouched"), 0o644); err != nil {
		t.Fatalf("seed victim: %v", err)
	}

	dest := filepath.Join(tmp, "out")
	if err := os.MkdirAll(dest, 0o755); err != nil {
		t.Fatalf("mkdir dest: %v", err)
	}
	// A symlink that already lives inside dest and points outside it. The
	// lexical guard alone cannot see this — dest/link/victim.txt is
	// lexically contained.
	if err := os.Symlink(target, filepath.Join(dest, "link")); err != nil {
		t.Skipf("symlinks unsupported here: %v", err)
	}

	src := writeZip(t, tmp, "evil.zip", []zipEntry{
		{name: "link/victim.txt", body: "pwned"},
		{name: "link/deep/other.txt", body: "pwned"},
	})

	if _, err := callUnzip(t, src, dest); err == nil {
		t.Fatal("expected an error writing through a symlinked dir, got nil")
	}
	if got := readFile(t, victim); got != "untouched" {
		t.Errorf("victim.txt = %q, want %q", got, "untouched")
	}
	if _, err := os.Stat(filepath.Join(target, "deep")); !os.IsNotExist(err) {
		t.Errorf("directory created outside dest through symlink (stat err = %v)", err)
	}
}

func TestOsUnzipReplacesPreexistingSymlinkAtTarget(t *testing.T) {
	tmp := t.TempDir()
	victim := filepath.Join(tmp, "victim.txt")
	if err := os.WriteFile(victim, []byte("untouched"), 0o644); err != nil {
		t.Fatalf("seed victim: %v", err)
	}
	dest := filepath.Join(tmp, "out")
	if err := os.MkdirAll(dest, 0o755); err != nil {
		t.Fatalf("mkdir dest: %v", err)
	}
	// The escape the parent-directory guard cannot see: the symlink is the
	// target itself, so its parent is plain old dest.
	if err := os.Symlink(victim, filepath.Join(dest, "f.txt")); err != nil {
		t.Skipf("symlinks unsupported here: %v", err)
	}

	src := writeZip(t, tmp, "a.zip", []zipEntry{{name: "f.txt", body: "new"}})
	if _, err := callUnzip(t, src, dest); err != nil {
		t.Fatalf("unzip: %v", err)
	}
	if got := readFile(t, victim); got != "untouched" {
		t.Errorf("victim.txt = %q, want %q — wrote through the symlink", got, "untouched")
	}
	if got := readFile(t, filepath.Join(dest, "f.txt")); got != "new" {
		t.Errorf("f.txt = %q, want %q", got, "new")
	}
}

func TestOsUnzipSkipsSymlinkEntries(t *testing.T) {
	tmp := t.TempDir()
	src := writeZip(t, tmp, "a.zip", []zipEntry{
		{name: "real.txt", body: "real"},
		{name: "evil-link", body: "/etc/passwd", mode: 0o777 | fs.ModeSymlink},
	})
	dest := filepath.Join(tmp, "out")

	if _, err := callUnzip(t, src, dest); err != nil {
		t.Fatalf("unzip: %v", err)
	}
	if got := readFile(t, filepath.Join(dest, "real.txt")); got != "real" {
		t.Errorf("real.txt = %q, want %q", got, "real")
	}
	if _, err := os.Lstat(filepath.Join(dest, "evil-link")); !os.IsNotExist(err) {
		t.Errorf("symlink entry was materialised (lstat err = %v)", err)
	}
}

func TestOsUnzipRejectsDirEntryOverExistingFile(t *testing.T) {
	tmp := t.TempDir()
	src := writeZip(t, tmp, "a.zip", []zipEntry{{name: "collide", dir: true}})
	dest := filepath.Join(tmp, "out")
	if err := os.MkdirAll(dest, 0o755); err != nil {
		t.Fatalf("mkdir dest: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dest, "collide"), []byte("i am a file"), 0o644); err != nil {
		t.Fatalf("seed file: %v", err)
	}

	// Reporting success here would leave the directory entry's contents with
	// nowhere to go.
	if _, err := callUnzip(t, src, dest); err == nil {
		t.Fatal("expected an error for a dir entry colliding with a file, got nil")
	}
}

// unsupportedMethod is a compression method the reader has no decompressor
// for, so f.Open() fails on the way out.
const unsupportedMethod uint16 = 99

func TestOsUnzipKeepsExistingFileWhenEntryCannotBeOpened(t *testing.T) {
	tmp := t.TempDir()
	// Written with a method registered only on the writer side.
	path := filepath.Join(tmp, "odd.zip")
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("create zip: %v", err)
	}
	zw := zip.NewWriter(f)
	zw.RegisterCompressor(unsupportedMethod, func(w io.Writer) (io.WriteCloser, error) {
		return nopWriteCloser{w}, nil
	})
	w, err := zw.CreateHeader(&zip.FileHeader{Name: "keep.txt", Method: unsupportedMethod})
	if err != nil {
		t.Fatalf("create entry: %v", err)
	}
	if _, err := w.Write([]byte("replacement")); err != nil {
		t.Fatalf("write entry: %v", err)
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("close zip: %v", err)
	}
	f.Close()

	dest := filepath.Join(tmp, "out")
	if err := os.MkdirAll(dest, 0o755); err != nil {
		t.Fatalf("mkdir dest: %v", err)
	}
	target := filepath.Join(dest, "keep.txt")
	if err := os.WriteFile(target, []byte("precious"), 0o644); err != nil {
		t.Fatalf("seed file: %v", err)
	}

	if _, err := callUnzip(t, path, dest); err == nil {
		t.Fatal("expected an error for an unreadable entry, got nil")
	}
	// The extraction failed; it must not have destroyed the old file first.
	if got := readFile(t, target); got != "precious" {
		t.Errorf("keep.txt = %q, want %q — deleted before the entry failed to open", got, "precious")
	}
}

type nopWriteCloser struct{ io.Writer }

func (nopWriteCloser) Close() error { return nil }

func TestOsUnzipArgErrors(t *testing.T) {
	tmp := t.TempDir()
	src := writeZip(t, tmp, "a.zip", []zipEntry{{name: "f.txt", body: "x"}})
	fn := unzipFn(t)

	cases := []struct {
		name string
		args []vm.Value
	}{
		{"no args", nil},
		{"one arg", []vm.Value{vm.String(src)}},
		{"three args", []vm.Value{vm.String(src), vm.String(tmp), vm.String(tmp)}},
		{"non-string path", []vm.Value{vm.Int(1), vm.String(tmp)}},
		{"non-string dest", []vm.Value{vm.String(src), vm.Int(1)}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if _, err := fn.Invoke(c.args); err == nil {
				t.Error("expected an error, got nil")
			}
		})
	}
}

func TestOsUnzipMissingArchive(t *testing.T) {
	tmp := t.TempDir()
	if _, err := callUnzip(t, filepath.Join(tmp, "nope.zip"), filepath.Join(tmp, "out")); err == nil {
		t.Fatal("expected an error for a missing archive, got nil")
	}
}
