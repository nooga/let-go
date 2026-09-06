package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestGointeropFromDepsEdnAcceptsMetadataBearingVectors(t *testing.T) {
	dir := t.TempDir()
	deps := `{:gointerop ^:configured ["example.com/plain"]
             :gointerop-wrappers ^:configured ["example.com/wrapped"]}`
	if err := os.WriteFile(filepath.Join(dir, "deps.edn"), []byte(deps), 0o600); err != nil {
		t.Fatal(err)
	}

	got, err := gointeropFromDepsEdn(dir, false)
	if err != nil {
		t.Fatal(err)
	}
	want := []interopEntry{
		{pkg: "example.com/plain"},
		{pkg: "example.com/wrapped", smart: true},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("gointeropFromDepsEdn() = %#v, want %#v", got, want)
	}
}
