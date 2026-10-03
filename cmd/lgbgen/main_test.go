package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestLGBGenUsesSourceBootstrap(t *testing.T) {
	root := repoRoot(t)
	out := filepath.Join(t.TempDir(), "core_compiled.lgb")

	cmd := exec.Command("go", "run", "-tags", "bootstrap", "./cmd/lgbgen", out)
	cmd.Dir = root
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("lgbgen source bootstrap failed: %v\n%s", err, output)
	}

	info, err := os.Stat(out)
	if err != nil {
		t.Fatalf("stat generated lgb: %v", err)
	}
	if info.Size() == 0 {
		t.Fatal("generated lgb is empty")
	}
}

func repoRoot(t *testing.T) string {
	t.Helper()

	// go test runs a test in its package directory, cmd/lgbgen; the source
	// path from runtime.Caller is not usable, since -trimpath rewrites it.
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	return filepath.Clean(filepath.Join(wd, "..", ".."))
}
