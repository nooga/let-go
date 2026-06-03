/*
 * Copyright (c) 2026 let-go contributors; see CONTRIBUTORS.
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"io"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

func TestNormalizeResourceName(t *testing.T) {
	cases := []struct {
		in   string
		want string
		ok   bool
	}{
		{"index.html", "index.html", true},
		{"/index.html", "index.html", true},
		{"./templates/x.html", "templates/x.html", true},
		{"a//b", "a/b", true},
		{"templates/../static/x.css", "static/x.css", true},
		{"//index.html", "index.html", true}, // repeated leading slashes collapse
		{"", "", false},
		{".", "", false},
		{"../secret", "", false},
		{"/../secret", "", false},
		{"//a/../../x", "", false}, // leading slashes must not bypass the .. check
		{"a/../../escape", "", false},
	}
	for _, c := range cases {
		got, ok := NormalizeResourceName(c.in)
		if got != c.want || ok != c.ok {
			t.Errorf("NormalizeResourceName(%q) = (%q, %v), want (%q, %v)",
				c.in, got, ok, c.want, c.ok)
		}
	}
}

func TestResourceArchiveRoundTrip(t *testing.T) {
	files := map[string][]byte{
		"index.html":         []byte("<!doctype html><h1>hi</h1>"),
		"css/app.css":        []byte("body{margin:0}"),
		"data/blob.bin":      {0x00, 0x01, 0xff, 0xfe, 0x80, 0x7f}, // non-UTF-8
		"empty.txt":          {},
		"nested/deep/x.json": []byte(`{"k":1}`),
	}

	blob, err := EncodeResourceArchive(files)
	if err != nil {
		t.Fatalf("EncodeResourceArchive: %v", err)
	}
	if len(blob) == 0 {
		t.Fatalf("EncodeResourceArchive returned empty blob")
	}

	got, err := DecodeResourceArchive(blob)
	if err != nil {
		t.Fatalf("DecodeResourceArchive: %v", err)
	}
	if !reflect.DeepEqual(got, files) {
		t.Fatalf("round-trip mismatch:\n got  %#v\n want %#v", got, files)
	}
}

func TestDecodeResourceArchiveRejectsGarbage(t *testing.T) {
	if _, err := DecodeResourceArchive([]byte("not a gzip stream")); err == nil {
		t.Fatalf("expected error decoding garbage, got nil")
	}
}

// openString reads a provider entry fully into a string (or fails the test if
// the entry is missing).
func openString(t *testing.T, p ResourceProvider, name string) (string, bool) {
	t.Helper()
	rc, ok := p.Open(name)
	if !ok {
		return "", false
	}
	defer rc.Close()
	data, err := io.ReadAll(rc)
	if err != nil {
		t.Fatalf("read %q: %v", name, err)
	}
	return string(data), true
}

func TestFSResourceProviderFirstRootWins(t *testing.T) {
	root1 := t.TempDir()
	root2 := t.TempDir()
	if err := os.WriteFile(filepath.Join(root1, "x.txt"), []byte("from-root1"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root2, "x.txt"), []byte("from-root2"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root2, "only2.txt"), []byte("only2"), 0644); err != nil {
		t.Fatal(err)
	}
	p := NewFSResourceProvider([]string{root1, root2})

	if got, ok := openString(t, p, "x.txt"); !ok || got != "from-root1" {
		t.Errorf("x.txt = (%q,%v), want (from-root1,true)", got, ok)
	}
	if got, ok := openString(t, p, "only2.txt"); !ok || got != "only2" {
		t.Errorf("only2.txt = (%q,%v), want (only2,true)", got, ok)
	}
	// leading "./" and "/" normalize to the same entry
	if got, ok := openString(t, p, "/x.txt"); !ok || got != "from-root1" {
		t.Errorf("/x.txt = (%q,%v), want (from-root1,true)", got, ok)
	}
}

func TestFSResourceProviderMissingAndTraversal(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "ok.txt"), []byte("ok"), 0644); err != nil {
		t.Fatal(err)
	}
	// a sub-directory must not resolve as a readable resource
	if err := os.Mkdir(filepath.Join(root, "adir"), 0755); err != nil {
		t.Fatal(err)
	}
	p := NewFSResourceProvider([]string{root})

	if _, ok := p.Open("missing.txt"); ok {
		t.Errorf("missing.txt: expected not found")
	}
	if _, ok := p.Open("../" + filepath.Base(root) + "/ok.txt"); ok {
		t.Errorf("traversal: expected not found")
	}
	if _, ok := p.Open("adir"); ok {
		t.Errorf("directory: expected not found")
	}
	// FIFO/non-regular rejection is exercised in resources_unix_test.go
	// (syscall.Mkfifo is not available on every supported platform).
}

func TestFSResourceProviderRejectsSymlinkEscape(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	if err := os.WriteFile(filepath.Join(outside, "secret.txt"), []byte("SECRET"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(root, "up")); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	// An in-root symlink must still resolve (containment, not a blanket ban).
	if err := os.WriteFile(filepath.Join(root, "real.txt"), []byte("ok"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(root, "real.txt"), filepath.Join(root, "link.txt")); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}

	p := NewFSResourceProvider([]string{root})
	if _, ok := p.Open("up/secret.txt"); ok {
		t.Error("symlink escaping the root must be denied")
	}
	if got, ok := openString(t, p, "link.txt"); !ok || got != "ok" {
		t.Errorf("in-root symlink should resolve: got (%q,%v)", got, ok)
	}
}

func TestWithinRoot(t *testing.T) {
	sep := string(os.PathSeparator)
	root := sep + "a" + sep + "b"
	cases := []struct {
		p    string
		want bool
	}{
		{root, true},
		{root + sep + "c", true},
		{root + sep + "c" + sep + "d.txt", true},
		{sep + "a" + sep + "bc", false}, // sibling sharing a name prefix
		{sep + "a", false},              // ancestor
		{sep + "x", false},
	}
	for _, c := range cases {
		if got := WithinRoot(root, c.p); got != c.want {
			t.Errorf("WithinRoot(%q, %q) = %v, want %v", root, c.p, got, c.want)
		}
	}
}

func TestEmbeddedResourceProvider(t *testing.T) {
	files := map[string][]byte{
		"a.txt":     []byte("alpha"),
		"d/bin.dat": {0x00, 0xff},
	}
	p := NewEmbeddedResourceProvider(files)

	if got, ok := openString(t, p, "a.txt"); !ok || got != "alpha" {
		t.Errorf("a.txt = (%q,%v), want (alpha,true)", got, ok)
	}
	if got, ok := openString(t, p, "./d/bin.dat"); !ok || got != "\x00\xff" {
		t.Errorf("d/bin.dat = (%q,%v)", got, ok)
	}
	if _, ok := p.Open("nope"); ok {
		t.Errorf("nope: expected not found")
	}
}

func TestResourceProviderSeam(t *testing.T) {
	orig := GetResourceProvider()
	defer SetResourceProvider(orig)

	p := NewEmbeddedResourceProvider(map[string][]byte{"k": []byte("v")})
	SetResourceProvider(p)
	if GetResourceProvider() != p {
		t.Fatalf("GetResourceProvider did not return the set provider")
	}
}

func TestIOResource(t *testing.T) {
	orig := GetResourceProvider()
	defer SetResourceProvider(orig)
	SetResourceProvider(NewEmbeddedResourceProvider(map[string][]byte{
		"templates/index.html": []byte("<h1>hi</h1>"),
	}))

	ions := NS("io")
	resVar := ions.LookupLocal(vm.Symbol("resource"))
	if resVar == nil {
		t.Fatal("io/resource is not defined")
	}

	// Present resource → non-nil handle; slurp returns the contents.
	res, err := resVar.Invoke([]vm.Value{vm.String("templates/index.html")})
	if err != nil {
		t.Fatalf("io/resource: %v", err)
	}
	if res == vm.NIL {
		t.Fatalf("io/resource returned nil for an existing resource")
	}
	slurpVar := ions.LookupLocal(vm.Symbol("slurp"))
	content, err := slurpVar.Invoke([]vm.Value{res})
	if err != nil {
		t.Fatalf("io/slurp: %v", err)
	}
	if got, ok := content.(vm.String); !ok || string(got) != "<h1>hi</h1>" {
		t.Fatalf("io/slurp = %#v, want \"<h1>hi</h1>\"", content)
	}

	// Composes with io/reader too.
	if rv := ions.LookupLocal(vm.Symbol("reader")); rv != nil {
		res2, _ := resVar.Invoke([]vm.Value{vm.String("templates/index.html")})
		if _, err := rv.Invoke([]vm.Value{res2}); err != nil {
			t.Fatalf("io/reader on resource: %v", err)
		}
	}

	// Missing resource → nil (no error).
	missing, err := resVar.Invoke([]vm.Value{vm.String("nope.txt")})
	if err != nil {
		t.Fatalf("io/resource(missing): %v", err)
	}
	if missing != vm.NIL {
		t.Fatalf("io/resource(missing) = %#v, want nil", missing)
	}
}
