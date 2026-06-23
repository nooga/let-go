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
