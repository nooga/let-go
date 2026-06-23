package rt

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFileStorageLogicalKeys(t *testing.T) {
	root := t.TempDir()
	store, err := NewFileStorage(root)
	if err != nil {
		t.Fatal(err)
	}
	keys := []string{
		"save:slot-1",
		"save:slot/with/slashes",
		"settings",
		"..",
		"",
	}
	for _, key := range keys {
		if err := store.Set(key, "value:"+key); err != nil {
			t.Fatalf("Set(%q): %v", key, err)
		}
	}
	for _, key := range keys {
		got, ok, err := store.Get(key)
		if err != nil {
			t.Fatalf("Get(%q): %v", key, err)
		}
		if !ok || got != "value:"+key {
			t.Fatalf("Get(%q) = %q, %v; want value", key, got, ok)
		}
	}
	got, err := store.Keys("save:")
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"save:slot-1", "save:slot/with/slashes"}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Fatalf("Keys(save:) = %#v, want %#v", got, want)
	}

	entries, err := os.ReadDir(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		name := e.Name()
		if strings.Contains(name, "/") || name == "." || name == ".." {
			t.Fatalf("unsafe filename for logical key: %q", name)
		}
		if filepath.Base(name) != name {
			t.Fatalf("filename escaped root: %q", name)
		}
	}
}

func TestFileStorageRemoveMissingIsNil(t *testing.T) {
	store, err := NewFileStorage(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Remove("missing"); err != nil {
		t.Fatalf("Remove missing: %v", err)
	}
}

// Constructing a store must not touch the filesystem; the root is created
// lazily on first Set so unused stores leave no empty dir behind.
func TestFileStorageLazyRoot(t *testing.T) {
	root := filepath.Join(t.TempDir(), "store")
	store, err := NewFileStorage(root)
	if err != nil {
		t.Fatal(err)
	}
	if _, statErr := os.Stat(root); !os.IsNotExist(statErr) {
		t.Fatalf("root created before any write: stat err = %v", statErr)
	}
	if _, ok, err := store.Get("k"); err != nil || ok {
		t.Fatalf("Get on absent root = ok %v, err %v; want absent, nil", ok, err)
	}
	if keys, err := store.Keys(""); err != nil || len(keys) != 0 {
		t.Fatalf("Keys on absent root = %#v, err %v; want empty", keys, err)
	}
	if err := store.Remove("k"); err != nil {
		t.Fatalf("Remove on absent root: %v", err)
	}
	if _, statErr := os.Stat(root); !os.IsNotExist(statErr) {
		t.Fatalf("read paths created the root: stat err = %v", statErr)
	}
	if err := store.Set("k", "v"); err != nil {
		t.Fatal(err)
	}
	if _, statErr := os.Stat(root); statErr != nil {
		t.Fatalf("Set did not create root: %v", statErr)
	}
}

// Set replaces atomically via a temp file; overwriting must return the latest
// value and must not leave the temp file visible to Keys.
func TestFileStorageOverwriteNoTempLeak(t *testing.T) {
	store, err := NewFileStorage(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Set("k", "v1"); err != nil {
		t.Fatal(err)
	}
	if err := store.Set("k", "v2"); err != nil {
		t.Fatal(err)
	}
	got, ok, err := store.Get("k")
	if err != nil || !ok || got != "v2" {
		t.Fatalf("Get = %q, %v, %v; want v2", got, ok, err)
	}
	keys, err := store.Keys("")
	if err != nil {
		t.Fatal(err)
	}
	if len(keys) != 1 || keys[0] != "k" {
		t.Fatalf("Keys = %#v, want [k] (temp files must not leak)", keys)
	}
}
