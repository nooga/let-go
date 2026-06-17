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

	// Pin every bundle shape across both axes: shell (xterm | none) and
	// payload delivery (inline | external).
	for _, tc := range []struct {
		name     string
		shell    bool
		external bool
		golden   string
	}{
		{"xterm shell, inline wasm", true, false, "assemble_golden.html"},
		{"shell-less core, inline wasm", false, false, "assemble_golden_shellless.html"},
		{"xterm shell, external wasm", true, true, "assemble_golden_external.html"},
		{"shell-less core, external wasm", false, true, "assemble_golden_shellless_external.html"},
	} {
		got := AssembleHTML(wasmExecJS, wasmGzB64, tc.shell, tc.external)
		goldenPath := filepath.Join("testdata", tc.golden)

		if *updateGolden {
			if err := os.MkdirAll("testdata", 0755); err != nil {
				t.Fatalf("mkdir testdata: %v", err)
			}
			if err := os.WriteFile(goldenPath, []byte(got), 0644); err != nil {
				t.Fatalf("writing golden: %v", err)
			}
			t.Logf("golden updated: %s (%d bytes)", goldenPath, len(got))
			continue
		}

		golden, err := os.ReadFile(goldenPath)
		if err != nil {
			t.Fatalf("%s: reading golden (run `go test ./pkg/rt/wasm -update` to create): %v", tc.name, err)
		}
		if string(golden) != got {
			t.Errorf("%s: AssembleHTML output drift (golden=%d bytes, got=%d bytes).\n"+
				"Run `go test ./pkg/rt/wasm -update` to refresh after intentional changes.",
				tc.name, len(golden), len(got))
		}
	}
}

// TestMarkersGone protects against a different failure: the substitution
// could succeed structurally but leave stray markers if the source
// gains another copy of __WASM_EXEC_JS__ etc. End-to-end build still
// works in that case (the broken marker is just literal text in the
// JS), so the golden test alone wouldn't catch it cleanly.
func TestMarkersGone(t *testing.T) {
	for _, shell := range []bool{true, false} {
		for _, external := range []bool{true, false} {
			got := AssembleHTML("anything", "whatever", shell, external)
			for _, m := range []string{
				"__WASM_EXEC_JS__",
				"__WASM_MODE__",
				"__WASM_GZ_B64__",
				"__LG_HOST_JS_BODY_PLACEHOLDER__",
				"__LG_XTERM_CSS__",
				"__LG_XTERM_JS__",
			} {
				if contains(got, m) {
					t.Errorf("shell=%v external=%v: marker %q still present in assembled output", shell, external, m)
				}
			}
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
