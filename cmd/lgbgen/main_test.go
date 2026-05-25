package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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

	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}
