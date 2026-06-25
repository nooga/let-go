/*
 * Copyright (c) 2026 let-go contributors; see CONTRIBUTORS.
 * SPDX-License-Identifier: MIT
 */

package bundle

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"
)

func writeTrailerFile(t *testing.T, base, lgb, res, id []byte, kind string, lgbSizeField, resSizeField, idSizeField uint64) string {
	t.Helper()
	buf := append([]byte{}, base...)
	buf = append(buf, lgb...)
	buf = append(buf, res...)
	buf = append(buf, id...)
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
	case "lgb3":
		var tr [28]byte
		binary.LittleEndian.PutUint64(tr[0:8], lgbSizeField)
		binary.LittleEndian.PutUint64(tr[8:16], resSizeField)
		binary.LittleEndian.PutUint64(tr[16:24], idSizeField)
		copy(tr[24:], bundleMagicV3[:])
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
	return BaseBinarySize(f)
}

func TestPayloadFitsFile(t *testing.T) {
	const maxI64 = int64(^uint64(0) >> 1) // math.MaxInt64
	cases := []struct {
		name           string
		lgb, res, id   uint64
		trailer, total int64
		want           bool
	}{
		{"legacy fits", 7, 0, 0, 12, 30, true},
		{"v2 fits exactly", 3, 6, 0, 20, 29, true},
		{"v3 fits exactly", 3, 6, 4, 28, 41, true},
		{"v3 id exceeds remainder", 3, 6, 5, 28, 41, false},
		{"lgb exceeds file", 30, 0, 0, 12, 30, false},
		{"huge lgb", 0xFFFFFFFFFFFFFFFF, 0, 0, 20, 30, false},
		// Each size <= total, but lgb+res+id+trailer would overflow uint64 if summed.
		{"sum overflows uint64", uint64(maxI64), uint64(maxI64), uint64(maxI64), 28, maxI64, false},
	}
	for _, c := range cases {
		if got := payloadFitsFile(c.lgb, c.res, c.id, c.trailer, c.total); got != c.want {
			t.Errorf("%s: payloadFitsFile(%d,%d,%d,%d,%d) = %v, want %v",
				c.name, c.lgb, c.res, c.id, c.trailer, c.total, got, c.want)
		}
	}
}

func TestParseBundleTrailer(t *testing.T) {
	base := []byte("BASEBINARY") // len 10

	t.Run("valid LGBX", func(t *testing.T) {
		lgb := []byte("LGBDATA")
		p := writeTrailerFile(t, base, lgb, nil, nil, "lgbx", uint64(len(lgb)), 0, 0)
		gotLgb, gotRes, gotID := ReadBundled(p)
		if string(gotLgb) != "LGBDATA" || gotRes != nil || gotID != "" {
			t.Fatalf("ReadBundled = (%q, %v, %q), want (LGBDATA, nil, \"\")", gotLgb, gotRes, gotID)
		}
		if bs, err := baseSize(t, p); err != nil || bs != int64(len(base)) {
			t.Fatalf("BaseBinarySize = (%d, %v), want (%d, nil)", bs, err, len(base))
		}
	})

	t.Run("valid LGB2", func(t *testing.T) {
		lgb := []byte("LGB")
		res := []byte("RESARC")
		p := writeTrailerFile(t, base, lgb, res, nil, "lgb2", uint64(len(lgb)), uint64(len(res)), 0)
		gotLgb, gotRes, gotID := ReadBundled(p)
		if string(gotLgb) != "LGB" || string(gotRes) != "RESARC" || gotID != "" {
			t.Fatalf("ReadBundled = (%q, %q, %q), want (LGB, RESARC, \"\")", gotLgb, gotRes, gotID)
		}
		if bs, err := baseSize(t, p); err != nil || bs != int64(len(base)) {
			t.Fatalf("BaseBinarySize = (%d, %v), want (%d, nil)", bs, err, len(base))
		}
	})

	t.Run("valid LGB3 with resources", func(t *testing.T) {
		lgb := []byte("LGB")
		res := []byte("RESARC")
		id := []byte("my-app")
		p := writeTrailerFile(t, base, lgb, res, id, "lgb3", uint64(len(lgb)), uint64(len(res)), uint64(len(id)))
		gotLgb, gotRes, gotID := ReadBundled(p)
		if string(gotLgb) != "LGB" || string(gotRes) != "RESARC" || gotID != "my-app" {
			t.Fatalf("ReadBundled = (%q, %q, %q), want (LGB, RESARC, my-app)", gotLgb, gotRes, gotID)
		}
		if bs, err := baseSize(t, p); err != nil || bs != int64(len(base)) {
			t.Fatalf("BaseBinarySize = (%d, %v), want (%d, nil)", bs, err, len(base))
		}
	})

	t.Run("valid LGB3 without resources", func(t *testing.T) {
		lgb := []byte("LGBDATA")
		id := []byte("solo")
		p := writeTrailerFile(t, base, lgb, nil, id, "lgb3", uint64(len(lgb)), 0, uint64(len(id)))
		gotLgb, gotRes, gotID := ReadBundled(p)
		if string(gotLgb) != "LGBDATA" || gotRes != nil || gotID != "solo" {
			t.Fatalf("ReadBundled = (%q, %v, %q), want (LGBDATA, nil, solo)", gotLgb, gotRes, gotID)
		}
	})

	t.Run("corrupt huge lgbSize does not panic", func(t *testing.T) {
		lgb := []byte("x")
		// lgbSize field claims a size far larger than the file.
		p := writeTrailerFile(t, base, lgb, nil, nil, "lgb2", 0xFFFFFFFFFFFFFFFF, 0, 0)
		gotLgb, gotRes, _ := ReadBundled(p)
		if gotLgb != nil || gotRes != nil {
			t.Fatalf("ReadBundled on corrupt trailer = (%q, %q), want (nil, nil)", gotLgb, gotRes)
		}
		if _, err := baseSize(t, p); err == nil {
			t.Fatalf("BaseBinarySize on corrupt trailer: expected error, got nil")
		}
	})

	t.Run("corrupt huge idSize does not panic", func(t *testing.T) {
		lgb := []byte("x")
		id := []byte("y")
		// idSize field claims a size far larger than the file.
		p := writeTrailerFile(t, base, lgb, nil, id, "lgb3", uint64(len(lgb)), 0, 0xFFFFFFFFFFFFFFFF)
		gotLgb, _, gotID := ReadBundled(p)
		if gotLgb != nil || gotID != "" {
			t.Fatalf("ReadBundled on corrupt v3 trailer = (%q, %q), want (nil, \"\")", gotLgb, gotID)
		}
		if _, err := baseSize(t, p); err == nil {
			t.Fatalf("BaseBinarySize on corrupt v3 trailer: expected error, got nil")
		}
	})

	t.Run("non-bundle file", func(t *testing.T) {
		junk := []byte("just some random bytes, definitely not a bundle trailer!!")
		p := writeTrailerFile(t, junk, nil, nil, nil, "none", 0, 0, 0)
		gotLgb, _, _ := ReadBundled(p)
		if gotLgb != nil {
			t.Fatalf("ReadBundled on non-bundle = %q, want nil", gotLgb)
		}
		if bs, err := baseSize(t, p); err != nil || bs != int64(len(junk)) {
			t.Fatalf("BaseBinarySize = (%d, %v), want (%d, nil)", bs, err, len(junk))
		}
	})
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

	files, err := CollectResources([]string{root}, "")
	if err != nil {
		t.Fatalf("CollectResources: %v", err)
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

	files, err := CollectResources([]string{root}, "") // must terminate
	if err != nil {
		t.Fatalf("CollectResources: %v", err)
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

	files, err := CollectResources([]string{root}, abs)
	if err != nil {
		t.Fatalf("CollectResources: %v", err)
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

	files, err := CollectResources([]string{root}, "")
	if err != nil {
		t.Fatalf("CollectResources: %v", err)
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
