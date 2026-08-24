package genmanifest

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestManifestFieldRoundTripWhitespaceAndPercent(t *testing.T) {
	for _, value := range []string{
		"plain/path.go",
		"dir with spaces/embed file.txt",
		"tabs\tand\nnewlines%plus+signs",
	} {
		encoded := encodeManifestField(value)
		if strings.ContainsAny(encoded, " \t\r\n") {
			t.Fatalf("encoded field still contains whitespace: %q", encoded)
		}
		got, err := decodeManifestField(encoded)
		if err != nil {
			t.Fatalf("decode %q: %v", encoded, err)
		}
		if got != value {
			t.Fatalf("round trip %q: got %q", value, got)
		}
	}
}

func TestEdgesCoverExpectedOutputs(t *testing.T) {
	root, err := FindRepoRoot(".")
	if err != nil {
		t.Fatalf("repo root: %v", err)
	}
	edges, err := Edges(root)
	if err != nil {
		t.Fatalf("Edges: %v", err)
	}
	byOutput := map[string]bool{}
	for _, e := range edges {
		byOutput[e.Output] = true
	}
	for _, want := range []string{
		"pkg/rt/zz_primitives_generated.go",
		"pkg/rt/corefns/zz_primitives_generated.go",
		"pkg/rt/core_compiled.lgb",
		"pkg/rt/core_go_lowered/",
		"pkg/ir/op_generated.go",
		"pkg/rt/ir_bridge_generated.go",
		"pkg/rt/core/ir/data/generated.lg",
	} {
		if !byOutput[want] {
			t.Errorf("no edges for expected output %q", want)
		}
	}
}

func TestEdgesExcludeGeneratedAndTestInputs(t *testing.T) {
	root, err := FindRepoRoot(".")
	if err != nil {
		t.Fatalf("repo root: %v", err)
	}
	edges, err := Edges(root)
	if err != nil {
		t.Fatalf("Edges: %v", err)
	}
	for _, e := range edges {
		if e.Kind == "input" && e.Output == "pkg/rt/zz_primitives_generated.go" {
			if e.Input == "pkg/rt/zz_primitives_generated.go" {
				t.Errorf("registrar input sweep must exclude the generated output itself")
			}
			if len(e.Input) > 8 && e.Input[len(e.Input)-8:] == "_test.go" {
				t.Errorf("registrar input sweep must exclude _test.go: %s", e.Input)
			}
			if strings.HasPrefix(e.Input, "pkg/rt/core_go_lowered/") {
				t.Errorf("registrar input sweep must exclude generated directory outputs: %s", e.Input)
			}
		}
	}
}

func TestSweepPrunesDeclaredDirectoryOutputsBeforeOpeningFiles(t *testing.T) {
	root := t.TempDir()
	lowered := filepath.Join(root, "pkg", "rt", "core_go_lowered")
	if err := os.MkdirAll(lowered, 0o755); err != nil {
		t.Fatal(err)
	}
	// A broken .go symlink makes any attempted isGenerated/read fail. The
	// declared output directory must be pruned before its files are inspected.
	if err := os.Symlink("missing.go", filepath.Join(lowered, "unstable.go")); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
	files, err := sweepFiles(root, sweep{"pkg/rt", ".go"}, true)
	if err != nil {
		t.Fatalf("declared output directory was traversed: %v", err)
	}
	for _, file := range files {
		if strings.HasPrefix(file, "pkg/rt/core_go_lowered/") {
			t.Fatalf("declared output leaked into sweep: %s", file)
		}
	}
}

// isolatedRepoCopy builds a tempdir holding exactly the input files the
// dependency manifest reads — every edge's input, at its real relative path —
// so the write-path tests exercise WriteDepManifest/StaleOutputs against a
// private tree instead of rewriting the tracked pkg/rt/generated.manifest
// (which would invalidate generated.sums and fail TestGeneratedArtifactsAreFresh
// when the package runs). Only inputs are needed: WriteDepManifest hashes edge
// inputs, never outputs. Placeholder outputs are created so output-existence
// checks pass during staleness tests.
func isolatedRepoCopy(t *testing.T) string {
	t.Helper()
	realRoot, err := FindRepoRoot(".")
	if err != nil {
		t.Fatalf("repo root: %v", err)
	}
	realEdges, err := Edges(realRoot)
	if err != nil {
		t.Fatalf("edges: %v", err)
	}
	dst := t.TempDir()
	for _, moduleFile := range []string{"go.mod", "go.sum"} {
		data, err := os.ReadFile(filepath.Join(realRoot, moduleFile))
		if err != nil {
			t.Fatalf("read %s: %v", moduleFile, err)
		}
		if err := os.WriteFile(filepath.Join(dst, moduleFile), data, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	// Copy all input files needed for dynamic dependency discovery.
	for _, e := range realEdges {
		data, err := os.ReadFile(filepath.Join(realRoot, e.Input))
		if err != nil {
			t.Fatalf("read %s: %v", e.Input, err)
		}
		out := filepath.Join(dst, e.Input)
		if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(out, data, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	edges, err := Edges(dst)
	if err != nil {
		t.Fatalf("temp edges: %v", err)
	}
	// Create placeholder outputs so output-existence checks pass
	outputs := map[string]bool{}
	for _, e := range edges {
		outputs[e.Output] = true
	}
	for out := range outputs {
		outPath := filepath.Join(dst, out)
		if strings.HasSuffix(out, "/") {
			// Directory output: create the directory with a placeholder file
			if err := os.MkdirAll(outPath, 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(outPath, ".placeholder"), []byte(""), 0o644); err != nil {
				t.Fatal(err)
			}
		} else {
			// Preserve generated outputs already copied because another edge uses
			// them as generator inputs.
			if _, err := os.Stat(outPath); err == nil {
				continue
			} else if !os.IsNotExist(err) {
				t.Fatal(err)
			}
			// Readiness validation is format-aware, so use the real generated
			// artifact rather than a vacuous empty placeholder.
			data, err := os.ReadFile(filepath.Join(realRoot, out))
			if err != nil {
				t.Fatalf("read generated output %s: %v", out, err)
			}
			if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(outPath, data, 0o644); err != nil {
				t.Fatal(err)
			}
		}
	}
	if err := WriteTreeManifest(filepath.Join(dst, "pkg/rt/core_go_lowered")); err != nil {
		t.Fatal(err)
	}
	return dst
}

func TestDepManifestRoundTrip(t *testing.T) {
	root := isolatedRepoCopy(t)
	if err := WriteDepManifest(root); err != nil {
		t.Fatalf("write: %v", err)
	}
	got, err := ReadDepManifest(root)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if len(got) == 0 {
		t.Fatal("manifest has no edges")
	}
	outputRecords := map[string]bool{}
	for _, he := range got {
		if len(he.Sum) != 64 { // hex sha256
			t.Errorf("edge %v: bad sha %q", he.Edge, he.Sum)
		}
		if he.Kind == "output" {
			if he.Input != he.Output {
				t.Errorf("output readiness record must hash itself: %+v", he.Edge)
			}
			outputRecords[he.Output] = true
		}
	}
	for _, spec := range outputSpecs {
		if !strings.HasSuffix(spec.output, "/") && !outputRecords[spec.output] {
			t.Errorf("missing output readiness record for %s", spec.output)
		}
	}
	// Deterministic: a second write is byte-identical.
	before, _ := os.ReadFile(filepath.Join(root, DepManifestRelPath))
	if err := WriteDepManifest(root); err != nil {
		t.Fatalf("rewrite: %v", err)
	}
	after, _ := os.ReadFile(filepath.Join(root, DepManifestRelPath))
	if string(before) != string(after) {
		t.Error("manifest is not deterministic across writes")
	}
}

func TestStaleOutputsDetectsChangedInput(t *testing.T) {
	root := isolatedRepoCopy(t)
	if err := WriteDepManifest(root); err != nil {
		t.Fatalf("write: %v", err)
	}
	// Fresh manifest → nothing stale.
	stale, err := StaleOutputs(root)
	if err != nil {
		t.Fatalf("stale: %v", err)
	}
	if len(stale) != 0 {
		t.Fatalf("expected nothing stale right after write, got %v", stale)
	}
	// Mutate one .lg input (in the isolated copy): only the bundle+lowered go
	// stale, NOT the registrar (a .lg edit must not dirty the Go-annotation
	// output).
	lg := filepath.Join(root, "pkg/rt/core/core.lg")
	orig, _ := os.ReadFile(lg)
	if err := os.WriteFile(lg, append(append([]byte{}, orig...), []byte("\n;; probe\n")...), 0644); err != nil {
		t.Fatal(err)
	}
	stale, err = StaleOutputs(root)
	if err != nil {
		t.Fatalf("stale: %v", err)
	}
	set := map[string]bool{}
	for _, s := range stale {
		set[s] = true
	}
	if !set["pkg/rt/core_compiled.lgb"] {
		t.Error("core.lg edit should mark the bundle stale")
	}
	if set["pkg/rt/zz_primitives_generated.go"] {
		t.Error("a .lg edit must NOT mark the Go registrar stale")
	}
}

func TestMissingOutputDetected(t *testing.T) {
	root := isolatedRepoCopy(t)
	if err := WriteDepManifest(root); err != nil {
		t.Fatalf("write: %v", err)
	}
	// Fresh manifest → nothing stale.
	stale, err := StaleOutputs(root)
	if err != nil {
		t.Fatalf("stale: %v", err)
	}
	if len(stale) != 0 {
		t.Fatalf("expected nothing stale right after write, got %v", stale)
	}
	// Remove a file output and verify it's detected as stale
	if err := os.Remove(filepath.Join(root, "pkg/rt/zz_primitives_generated.go")); err != nil {
		t.Fatal(err)
	}
	stale, err = StaleOutputs(root)
	if err != nil {
		t.Fatalf("stale after removal: %v", err)
	}
	found := false
	for _, s := range stale {
		if s == "pkg/rt/zz_primitives_generated.go" {
			found = true
			break
		}
	}
	if !found {
		t.Error("missing file output should be detected as stale")
	}
	// Remove a directory output and verify it's detected as stale
	dirOut := filepath.Join(root, "pkg/rt/core_go_lowered")
	if err := os.RemoveAll(dirOut); err != nil {
		t.Fatal(err)
	}
	stale, err = StaleOutputs(root)
	if err != nil {
		t.Fatalf("stale after dir removal: %v", err)
	}
	found = false
	for _, s := range stale {
		if s == "pkg/rt/core_go_lowered/" {
			found = true
			break
		}
	}
	if !found {
		t.Error("missing directory output should be detected as stale")
	}
}

func TestIncompleteFileOutputsDetected(t *testing.T) {
	root := isolatedRepoCopy(t)
	if err := WriteDepManifest(root); err != nil {
		t.Fatalf("write: %v", err)
	}

	assertStale := func(t *testing.T, want string) {
		t.Helper()
		stale, err := StaleOutputs(root)
		if err != nil {
			t.Fatalf("stale: %v", err)
		}
		for _, output := range stale {
			if output == want {
				return
			}
		}
		t.Fatalf("incomplete output %s must be stale, got %v", want, stale)
	}

	for _, spec := range outputSpecs {
		if strings.HasSuffix(spec.output, "/") {
			continue
		}
		path := filepath.Join(root, spec.output)
		original, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", spec.output, err)
		}

		t.Run(spec.output+"/empty", func(t *testing.T) {
			if err := os.WriteFile(path, nil, 0o644); err != nil {
				t.Fatal(err)
			}
			assertStale(t, spec.output)
			if err := os.WriteFile(path, original, 0o644); err != nil {
				t.Fatal(err)
			}
		})

		var truncated []byte
		switch filepath.Ext(spec.output) {
		case ".lgb":
			truncated = []byte{'L', 'G', 'B', 0x01}
		case ".go":
			truncated = []byte("// Code generated by test. DO NOT EDIT.\n")
		case ".lg":
			truncated = []byte(";; Code generated by test. DO NOT EDIT.\n(\n")
		default:
			t.Fatalf("declared file output %s has no readiness contract", spec.output)
		}
		t.Run(spec.output+"/truncated", func(t *testing.T) {
			if err := os.WriteFile(path, truncated, 0o644); err != nil {
				t.Fatal(err)
			}
			assertStale(t, spec.output)
			if err := os.WriteFile(path, original, 0o644); err != nil {
				t.Fatal(err)
			}
		})

		t.Run(spec.output+"/directory", func(t *testing.T) {
			if err := os.Remove(path); err != nil {
				t.Fatal(err)
			}
			if err := os.Mkdir(path, 0o755); err != nil {
				t.Fatal(err)
			}
			assertStale(t, spec.output)
			if err := os.RemoveAll(path); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, original, 0o644); err != nil {
				t.Fatal(err)
			}
		})
	}

	const invalidGoOutput = "pkg/rt/zz_primitives_generated.go"
	invalidGoPath := filepath.Join(root, invalidGoOutput)
	originalGo, err := os.ReadFile(invalidGoPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(invalidGoPath, []byte(`// Code generated by test. DO NOT EDIT.
package rt
import _ "example.invalid/missing"
var Generated = true
`), 0o644); err != nil {
		t.Fatal(err)
	}
	assertStale(t, invalidGoOutput)
	if err := os.WriteFile(invalidGoPath, originalGo, 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestCheckDepManifestIgnoresMissingOutputs(t *testing.T) {
	root := isolatedRepoCopy(t)
	if err := WriteDepManifest(root); err != nil {
		t.Fatalf("write: %v", err)
	}
	if err := os.RemoveAll(filepath.Join(root, "pkg/rt/core_go_lowered")); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(filepath.Join(root, "pkg/rt/core_compiled.lgb")); err != nil {
		t.Fatal(err)
	}

	if err := CheckDepManifest(root); err != nil {
		t.Fatalf("output state must not make input hashes stale: %v", err)
	}
}

func TestStaleOutputsDetectsTornLoweredTree(t *testing.T) {
	root := isolatedRepoCopy(t)
	if err := WriteDepManifest(root); err != nil {
		t.Fatalf("write: %v", err)
	}
	lowered := filepath.Join(root, "pkg/rt/core_go_lowered")
	if err := os.Remove(filepath.Join(lowered, ".placeholder")); err != nil {
		t.Fatal(err)
	}

	stale, err := StaleOutputs(root)
	if err != nil {
		t.Fatalf("stale: %v", err)
	}
	for _, output := range stale {
		if output == "pkg/rt/core_go_lowered/" {
			return
		}
	}
	t.Fatalf("torn lowered tree must be stale, got %v", stale)
}

// generated.manifest has no merge driver and is one sorted record per line, so a
// clean text merge can land both sides' record for a single edge. The parse must
// refuse that, and Compute with it — Compute is the path scripts/git-merge-sums.sh
// calls, so a duplicated edge would otherwise mint a digest over a manifest that
// records two hashes for one input.
func TestDepManifestRejectsDuplicateEdge(t *testing.T) {
	root := isolatedRepoCopy(t)
	if err := WriteDepManifest(root); err != nil {
		t.Fatalf("write: %v", err)
	}
	path := filepath.Join(root, DepManifestRelPath)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var record string
	for _, line := range strings.Split(string(data), "\n") {
		if line != "" && !strings.HasPrefix(line, "#") {
			record = line
			break
		}
	}
	if record == "" {
		t.Fatal("manifest has no records to duplicate")
	}
	if err := os.WriteFile(path, append(data, []byte(record+"\n")...), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := ReadDepManifest(root); err == nil {
		t.Fatal("ReadDepManifest accepted a duplicated edge")
	}
	if _, err := Compute(root); err == nil {
		t.Fatal("Compute digested a manifest carrying a duplicated edge")
	}
}

// A manifest that merged wrong carries more than one duplicated edge, and the
// first one found says nothing about how far the damage runs. Every repeat is
// named so the report itself distinguishes a one-line slip from a bad merge.
func TestDepManifestReportsEveryDuplicateEdge(t *testing.T) {
	root := isolatedRepoCopy(t)
	if err := WriteDepManifest(root); err != nil {
		t.Fatalf("write: %v", err)
	}
	path := filepath.Join(root, DepManifestRelPath)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var records []string
	for _, line := range strings.Split(string(data), "\n") {
		if line != "" && !strings.HasPrefix(line, "#") {
			records = append(records, line)
			if len(records) == 2 {
				break
			}
		}
	}
	if len(records) < 2 {
		t.Fatalf("manifest has %d records, need 2 to duplicate", len(records))
	}
	// The second record is repeated twice, so the count carries past a single repeat.
	dupes := records[0] + "\n" + records[1] + "\n" + records[1] + "\n"
	if err := os.WriteFile(path, append(data, []byte(dupes)...), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err = ReadDepManifest(root)
	if err == nil {
		t.Fatal("ReadDepManifest accepted duplicated edges")
	}
	msg := err.Error()
	for i, record := range records {
		edge := strings.Fields(record)[:3]
		for _, field := range edge {
			decoded, decErr := decodeManifestField(field)
			if decErr != nil {
				t.Fatalf("decode %q: %v", field, decErr)
			}
			if !strings.Contains(msg, decoded) {
				t.Fatalf("duplicate report omits record %d field %q: %s", i, decoded, msg)
			}
		}
	}
	if !strings.Contains(msg, "2 duplicate manifest entries") {
		t.Fatalf("duplicate report does not count both edges: %s", msg)
	}
	if !strings.Contains(msg, "(3 occurrences)") {
		t.Fatalf("duplicate report does not count the thrice-recorded edge: %s", msg)
	}
}
