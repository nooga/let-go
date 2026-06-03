package wasm

import (
	"flag"
	"os"
	"path/filepath"
	"testing"
)

var updateGolden = flag.Bool("update", false, "update golden files in testdata/")

// TestAssembleHTMLGolden pins the assembled HTML against
// testdata/assemble_golden.html. Any edit to host.html, lg-host.js, the
// markers, or the assembly logic that changes the bundle `lg -w` ships
// surfaces as a golden diff.
//
// This is a self-consistency pin, not a byte-identity guarantee against
// any prior implementation. Regenerate after intentional changes:
//
//	go test ./pkg/rt/wasm -update
func TestAssembleHTMLGolden(t *testing.T) {
	// Fixed, recognizable stubs. The real wasm_exec.js and gzipped WASM
	// blob change every build (non-deterministic Go toolchain output);
	// the test pins everything else.
	const wasmExecJS = "// stub wasm_exec.js for the golden test\nconsole.log('exec stub');\n"
	const wasmGzB64 = "STUBWASMBLOBB64=="

	got := AssembleHTML(wasmExecJS, wasmGzB64)
	goldenPath := filepath.Join("testdata", "assemble_golden.html")

	if *updateGolden {
		if err := os.MkdirAll("testdata", 0755); err != nil {
			t.Fatalf("mkdir testdata: %v", err)
		}
		if err := os.WriteFile(goldenPath, []byte(got), 0644); err != nil {
			t.Fatalf("writing golden: %v", err)
		}
		t.Logf("golden updated: %s (%d bytes)", goldenPath, len(got))
		return
	}

	golden, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("reading golden (run `go test ./pkg/rt/wasm -update` to create): %v", err)
	}
	if string(golden) != got {
		// Surface the size delta so the developer has a quick signal
		// before diving into a diff.
		t.Errorf("AssembleHTML output drift (golden=%d bytes, got=%d bytes).\n"+
			"Run `go test ./pkg/rt/wasm -update` to refresh after intentional changes.",
			len(golden), len(got))
	}
}

// TestMarkersGone protects against a different failure: the substitution
// could succeed structurally but leave stray markers if the source
// gains another copy of __WASM_EXEC_JS__ etc. End-to-end build still
// works in that case (the broken marker is just literal text in the
// JS), so the golden test alone wouldn't catch it cleanly.
func TestMarkersGone(t *testing.T) {
	got := AssembleHTML("anything", "whatever")
	for _, m := range []string{
		"__WASM_EXEC_JS__",
		"__WASM_GZ_B64__",
		"__LG_HOST_JS_BODY_PLACEHOLDER__",
	} {
		if contains(got, m) {
			t.Errorf("marker %q still present in assembled output", m)
		}
	}
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
