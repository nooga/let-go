/*
 * Copyright (c) 2026 let-go contributors; see CONTRIBUTORS.
 * SPDX-License-Identifier: MIT
 */

package main

import (
	"encoding/binary"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// writeTrailerFile builds base+lgb+res bytes followed by a trailer using the
// given raw size fields, so tests can craft both valid and corrupt trailers.
// kind is "lgbx", "lgb2", or "none".
func writeTrailerFile(t *testing.T, base, lgb, res []byte, kind string, lgbSizeField, resSizeField uint64) string {
	t.Helper()
	buf := append([]byte{}, base...)
	buf = append(buf, lgb...)
	buf = append(buf, res...)
	switch kind {
	case "lgbx":
		var tr [12]byte
		binary.LittleEndian.PutUint64(tr[:8], lgbSizeField)
		copy(tr[8:], bundleMagic[:])
		buf = append(buf, tr[:]...)
	case "lgb2":
		var tr [20]byte
		binary.LittleEndian.PutUint64(tr[0:8], lgbSizeField)
		binary.LittleEndian.PutUint64(tr[8:16], resSizeField)
		copy(tr[16:], bundleMagicV2[:])
		buf = append(buf, tr[:]...)
	case "none":
		// no trailer
	}
	p := filepath.Join(t.TempDir(), "bin")
	if err := os.WriteFile(p, buf, 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func baseSize(t *testing.T, path string) (int64, error) {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	return getBaseBinarySize(f)
}

func TestPayloadFitsFile(t *testing.T) {
	const maxI64 = int64(^uint64(0) >> 1) // math.MaxInt64
	cases := []struct {
		name           string
		lgb, res       uint64
		trailer, total int64
		want           bool
	}{
		{"legacy fits", 7, 0, 12, 30, true},
		{"v2 fits exactly", 3, 6, 20, 29, true},
		{"lgb exceeds file", 30, 0, 12, 30, false},
		{"huge lgb", 0xFFFFFFFFFFFFFFFF, 0, 20, 30, false},
		// Each size <= total, but lgb+res+trailer would overflow uint64 if summed.
		{"sum overflows uint64", uint64(maxI64), uint64(maxI64), 20, maxI64, false},
	}
	for _, c := range cases {
		if got := payloadFitsFile(c.lgb, c.res, c.trailer, c.total); got != c.want {
			t.Errorf("%s: payloadFitsFile(%d,%d,%d,%d) = %v, want %v",
				c.name, c.lgb, c.res, c.trailer, c.total, got, c.want)
		}
	}
}

func TestParseBundleTrailer(t *testing.T) {
	base := []byte("BASEBINARY") // len 10

	t.Run("valid LGBX", func(t *testing.T) {
		lgb := []byte("LGBDATA")
		p := writeTrailerFile(t, base, lgb, nil, "lgbx", uint64(len(lgb)), 0)
		gotLgb, gotRes := readBundledLGB(p)
		if string(gotLgb) != "LGBDATA" || gotRes != nil {
			t.Fatalf("readBundledLGB = (%q, %v), want (LGBDATA, nil)", gotLgb, gotRes)
		}
		if bs, err := baseSize(t, p); err != nil || bs != int64(len(base)) {
			t.Fatalf("getBaseBinarySize = (%d, %v), want (%d, nil)", bs, err, len(base))
		}
	})

	t.Run("valid LGB2", func(t *testing.T) {
		lgb := []byte("LGB")
		res := []byte("RESARC")
		p := writeTrailerFile(t, base, lgb, res, "lgb2", uint64(len(lgb)), uint64(len(res)))
		gotLgb, gotRes := readBundledLGB(p)
		if string(gotLgb) != "LGB" || string(gotRes) != "RESARC" {
			t.Fatalf("readBundledLGB = (%q, %q), want (LGB, RESARC)", gotLgb, gotRes)
		}
		if bs, err := baseSize(t, p); err != nil || bs != int64(len(base)) {
			t.Fatalf("getBaseBinarySize = (%d, %v), want (%d, nil)", bs, err, len(base))
		}
	})

	t.Run("corrupt huge lgbSize does not panic", func(t *testing.T) {
		lgb := []byte("x")
		// lgbSize field claims a size far larger than the file.
		p := writeTrailerFile(t, base, lgb, nil, "lgb2", 0xFFFFFFFFFFFFFFFF, 0)
		gotLgb, gotRes := readBundledLGB(p)
		if gotLgb != nil || gotRes != nil {
			t.Fatalf("readBundledLGB on corrupt trailer = (%q, %q), want (nil, nil)", gotLgb, gotRes)
		}
		if _, err := baseSize(t, p); err == nil {
			t.Fatalf("getBaseBinarySize on corrupt trailer: expected error, got nil")
		}
	})

	t.Run("non-bundle file", func(t *testing.T) {
		junk := []byte("just some random bytes, definitely not a bundle trailer!!")
		p := writeTrailerFile(t, junk, nil, nil, "none", 0, 0)
		gotLgb, _ := readBundledLGB(p)
		if gotLgb != nil {
			t.Fatalf("readBundledLGB on non-bundle = %q, want nil", gotLgb)
		}
		if bs, err := baseSize(t, p); err != nil || bs != int64(len(junk)) {
			t.Fatalf("getBaseBinarySize = (%d, %v), want (%d, nil)", bs, err, len(junk))
		}
	})
}

// TestResourceDevMode: running from source with -resource-paths, io/resource
// finds files on the filesystem and io/slurp reads them; a missing resource
// yields nil. (buildLG is defined in scope_e2e_test.go, same package.)
func TestResourceDevMode(t *testing.T) {
	bin := buildLG(t)
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "hello.txt"), []byte("hi-there"), 0644); err != nil {
		t.Fatal(err)
	}

	out, err := exec.Command(bin, "-resource-paths", dir, "-e",
		`(println (io/slurp (io/resource "hello.txt")))`).CombinedOutput()
	if err != nil {
		t.Fatalf("run: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "hi-there") {
		t.Fatalf("expected resource contents, got: %q", out)
	}

	out, err = exec.Command(bin, "-resource-paths", dir, "-e",
		`(println (if (io/resource "nope.txt") "FOUND" "MISSING"))`).CombinedOutput()
	if err != nil {
		t.Fatalf("run: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "MISSING") {
		t.Fatalf("expected MISSING for absent resource, got: %q", out)
	}
}

// TestResourceBundle: resources under -resource-paths are embedded into a -b
// standalone binary, and io/resource reads them at runtime even when the
// binary runs from a directory with no resource files present.
func TestResourceBundle(t *testing.T) {
	bin := buildLG(t)

	resDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(resDir, "msg.txt"), []byte("hello-resource"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(resDir, "sub"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(resDir, "sub", "n.txt"), []byte("nested-ok"), 0644); err != nil {
		t.Fatal(err)
	}

	prog := filepath.Join(t.TempDir(), "prog.lg")
	progSrc := `(println (io/slurp (io/resource "msg.txt")))` + "\n" +
		`(println (io/slurp (io/resource "sub/n.txt")))`
	if err := os.WriteFile(prog, []byte(progSrc), 0644); err != nil {
		t.Fatal(err)
	}

	outBin := filepath.Join(t.TempDir(), "app")
	if out, err := exec.Command(bin, "-b", outBin, "-resource-paths", resDir, prog).CombinedOutput(); err != nil {
		t.Fatalf("bundle: %v\n%s", err, out)
	}

	// Run from a clean cwd with no resource files around — only the embedded
	// copies can satisfy io/resource here.
	cmd := exec.Command(outBin)
	cmd.Dir = t.TempDir()
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("run bundle: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "hello-resource") {
		t.Fatalf("expected embedded resource contents, got: %q", out)
	}
	if !strings.Contains(string(out), "nested-ok") {
		t.Fatalf("expected nested embedded resource, got: %q", out)
	}
}

// TestCollectResourcesFollowsSymlinks: bundling matches the dev FS provider,
// which resolves names with os.Stat — so symlinks to files AND to directories
// are followed.
func TestCollectResourcesFollowsSymlinks(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "real.txt"), []byte("real-bytes"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(root, "real.txt"), filepath.Join(root, "link.txt")); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(root, "d"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "d", "inner.txt"), []byte("inner-bytes"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(root, "d"), filepath.Join(root, "dlink")); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}

	files, err := collectResources([]string{root}, "")
	if err != nil {
		t.Fatalf("collectResources: %v", err)
	}
	if string(files["real.txt"]) != "real-bytes" || string(files["d/inner.txt"]) != "inner-bytes" {
		t.Errorf("real files not embedded: %q / %q", files["real.txt"], files["d/inner.txt"])
	}
	if string(files["link.txt"]) != "real-bytes" {
		t.Errorf("symlink to file not embedded: got %q", files["link.txt"])
	}
	// A symlinked directory is followed and its contents embedded under the
	// symlink's name, matching what (io/resource "dlink/inner.txt") finds.
	if string(files["dlink/inner.txt"]) != "inner-bytes" {
		t.Errorf("symlink to directory not followed: got %q", files["dlink/inner.txt"])
	}
}

// TestCollectResourcesHandlesSymlinkCycle: a symlink that loops back to an
// ancestor must not cause infinite recursion.
func TestCollectResourcesHandlesSymlinkCycle(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "a.txt"), []byte("a"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "sub"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(root, filepath.Join(root, "sub", "loop")); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}

	files, err := collectResources([]string{root}, "") // must terminate
	if err != nil {
		t.Fatalf("collectResources: %v", err)
	}
	if string(files["a.txt"]) != "a" {
		t.Errorf("expected a.txt embedded, got %q", files["a.txt"])
	}
}

// TestCollectResourcesExcludesOutputBinary: a dst that lives inside a resource
// root must not be embedded into its own bundle.
func TestCollectResourcesExcludesOutputBinary(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "keep.txt"), []byte("keep"), 0644); err != nil {
		t.Fatal(err)
	}
	outBin := filepath.Join(root, "app")
	if err := os.WriteFile(outBin, []byte("BINARY"), 0755); err != nil {
		t.Fatal(err)
	}
	abs, _ := filepath.Abs(outBin)

	files, err := collectResources([]string{root}, abs)
	if err != nil {
		t.Fatalf("collectResources: %v", err)
	}
	if _, ok := files["app"]; ok {
		t.Errorf("output binary should be excluded from resources")
	}
	if string(files["keep.txt"]) != "keep" {
		t.Errorf("expected keep.txt embedded, got %q", files["keep.txt"])
	}
}

// TestCollectResourcesRejectsSymlinkEscape: a symlink under a resource root
// that points outside the root must not pull external files into the bundle.
func TestCollectResourcesRejectsSymlinkEscape(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	if err := os.WriteFile(filepath.Join(outside, "secret.txt"), []byte("SECRET"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "keep.txt"), []byte("keep"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(root, "up")); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}

	files, err := collectResources([]string{root}, "")
	if err != nil {
		t.Fatalf("collectResources: %v", err)
	}
	if string(files["keep.txt"]) != "keep" {
		t.Errorf("in-root file should be embedded, got %q", files["keep.txt"])
	}
	if _, ok := files["up/secret.txt"]; ok {
		t.Errorf("escaping symlink must not be embedded")
	}
}

// TestResourcePathsEmptyFlagOverridesEnv: -resource-paths "" clears the
// LG_RESOURCE_PATHS fallback; an unset flag still honors the env var.
func TestResourcePathsEmptyFlagOverridesEnv(t *testing.T) {
	bin := buildLG(t)
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "f.txt"), []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}
	check := `(println (if (io/resource "f.txt") "FOUND" "MISSING"))`

	// env set, no flag → env honored → FOUND
	cmd := exec.Command(bin, "-e", check)
	cmd.Env = append(os.Environ(), "LG_RESOURCE_PATHS="+dir)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("env run: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "FOUND") {
		t.Fatalf("env fallback: expected FOUND, got %q", out)
	}

	// env set, explicit empty flag → env cleared → MISSING
	cmd = exec.Command(bin, "-resource-paths", "", "-e", check)
	cmd.Env = append(os.Environ(), "LG_RESOURCE_PATHS="+dir)
	out, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("empty-flag run: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "MISSING") {
		t.Fatalf("empty flag should clear env: expected MISSING, got %q", out)
	}
}

// TestResourceBundleRelativePathSurvivesChdir: a relative -resource-paths is
// resolved against the cwd at launch, even if the program changes cwd during
// AOT compilation (top-level forms run at compile time). Regression test for
// resolving resource roots after user code instead of before.
func TestResourceBundleRelativePathSurvivesChdir(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("syscall/chdir is linux-only")
	}
	bin := buildLG(t)

	work := t.TempDir()
	if err := os.MkdirAll(filepath.Join(work, "res"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(work, "res", "msg.txt"), []byte("rel-resource"), 0644); err != nil {
		t.Fatal(err)
	}
	elsewhere := t.TempDir()

	// The chdir runs at both compile and run time; the resource read is guarded
	// to runtime only (at compile time the cwd has already changed away).
	prog := filepath.Join(work, "prog.lg")
	progSrc := `(syscall/chdir "` + elsewhere + `")` + "\n" +
		`(when-not *compiling-aot* (println (io/slurp (io/resource "msg.txt"))))`
	if err := os.WriteFile(prog, []byte(progSrc), 0644); err != nil {
		t.Fatal(err)
	}

	outBin := filepath.Join(t.TempDir(), "app")
	cmd := exec.Command(bin, "-b", outBin, "-resource-paths", "res", "prog.lg")
	cmd.Dir = work // relative "res" must resolve against here, not `elsewhere`
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("bundle: %v\n%s", err, out)
	}

	run := exec.Command(outBin)
	run.Dir = t.TempDir()
	out, err := run.CombinedOutput()
	if err != nil {
		t.Fatalf("run bundle: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "rel-resource") {
		t.Fatalf("relative resource path lost after compile-time chdir: %q", out)
	}
}

// TestLegacyBundleStillRuns: a -b bundle with no resources keeps working
// (exercises the resSize==0 / legacy-trailer path).
func TestLegacyBundleStillRuns(t *testing.T) {
	bin := buildLG(t)
	prog := filepath.Join(t.TempDir(), "prog.lg")
	if err := os.WriteFile(prog, []byte(`(println "no-resources-here")`), 0644); err != nil {
		t.Fatal(err)
	}
	outBin := filepath.Join(t.TempDir(), "app2")
	if out, err := exec.Command(bin, "-b", outBin, prog).CombinedOutput(); err != nil {
		t.Fatalf("bundle: %v\n%s", err, out)
	}
	cmd := exec.Command(outBin)
	cmd.Dir = t.TempDir()
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("run bundle: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "no-resources-here") {
		t.Fatalf("expected program output, got: %q", out)
	}
}
