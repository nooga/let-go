package test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestScriptWithNsDoesNotReExecute exercises the case where a script
// file declares a namespace whose name matches its filename. Without
// the fix, the compile-time (in-ns 'foo) early-detection in
// pkg/compiler/compiler.go would call rt.NS("foo"), which triggers
// the resolver to find foo.lg on disk and re-execute it inside the
// outer compile — producing duplicate side effects (println twice,
// state allocations twice).
//
// The fix swaps rt.NS for rt.LookupOrRegisterNSNoLoad in the
// compile-time hook, matching what the (in-ns) runtime function
// already does.
func TestScriptWithNsDoesNotReExecute(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "lg-ns-noreload-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// A minimal script: file 'foo.lg' declares (ns foo) and prints once.
	// The bug would cause it to print twice.
	scriptPath := filepath.Join(tmpDir, "foo.lg")
	script := "(ns foo)\n(println \"hello from foo\")\n"
	if err := os.WriteFile(scriptPath, []byte(script), 0644); err != nil {
		t.Fatal(err)
	}

	// Build the lg binary from the package we're testing.
	lgPath := filepath.Join(tmpDir, "lg-bin")
	build := exec.Command("go", "build", "-o", lgPath, "../")
	build.Dir = "." // tests run from the test/ directory
	var buildErr bytes.Buffer
	build.Stderr = &buildErr
	if err := build.Run(); err != nil {
		t.Fatalf("failed to build lg: %v\nstderr: %s", err, buildErr.String())
	}

	cmd := exec.Command(lgPath, scriptPath)
	cmd.Dir = tmpDir
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		t.Fatalf("lg failed: %v\noutput:\n%s", err, out.String())
	}

	output := out.String()
	count := strings.Count(output, "hello from foo")
	if count != 1 {
		t.Errorf("expected script body to execute exactly once (1 println), got %d executions\nfull output:\n%s",
			count, output)
	}
}
