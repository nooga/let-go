//go:build bootstrap

package main

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestWriteFileAtomicallyPreservesDestinationOnFailure(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "core_compiled.lgb")
	if err := os.WriteFile(path, []byte("original"), 0o640); err != nil {
		t.Fatal(err)
	}

	wantErr := errors.New("encode failed")
	_, err := writeFileAtomically(path, func(w io.Writer) error {
		if _, err := io.WriteString(w, "partial"); err != nil {
			return err
		}
		return wantErr
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("writeFileAtomically() error = %v, want %v", err, wantErr)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "original" {
		t.Fatalf("destination = %q after failed write, want original", got)
	}
	assertNoAtomicTemps(t, dir)
}

func TestWriteFileAtomicallyReplacesDestination(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "core_compiled.lgb")
	if err := os.WriteFile(path, []byte("old"), 0o640); err != nil {
		t.Fatal(err)
	}

	size, err := writeFileAtomically(path, func(w io.Writer) error {
		_, err := io.WriteString(w, "replacement")
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	if size != int64(len("replacement")) {
		t.Fatalf("size = %d, want %d", size, len("replacement"))
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "replacement" {
		t.Fatalf("destination = %q, want replacement", got)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if gotMode := info.Mode().Perm(); gotMode != 0o640 {
		t.Fatalf("destination mode = %o, want 640", gotMode)
	}
	assertNoAtomicTemps(t, dir)
}

func assertNoAtomicTemps(t *testing.T, dir string) {
	t.Helper()
	matches, err := filepath.Glob(filepath.Join(dir, ".core_compiled.lgb.tmp-*"))
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != 0 {
		t.Fatalf("temporary files left behind: %v", matches)
	}
}

// The size lgbgen prints must come from the file it names. Reading it off the
// temporary file gives the same number today but stops being checkable the
// moment the two can differ, so pin the destination as the source.
func TestWriteFileAtomicallyReportsDestinationSize(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "core_compiled.lgb")
	content := []byte("bundle bytes, more than the original")
	if err := os.WriteFile(path, []byte("short"), 0o644); err != nil {
		t.Fatal(err)
	}

	size, err := writeFileAtomically(path, func(w io.Writer) error {
		_, err := w.Write(content)
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	if size != int64(len(content)) {
		t.Errorf("reported size = %d, want %d", size, len(content))
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if size != info.Size() {
		t.Errorf("reported size = %d, destination on disk = %d", size, info.Size())
	}
	assertNoAtomicTemps(t, dir)
}
