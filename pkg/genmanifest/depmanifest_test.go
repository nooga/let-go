package genmanifest

import (
	"os"
	"path/filepath"
	"testing"
)

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
		}
	}
}

func TestDepManifestRoundTrip(t *testing.T) {
	root, err := FindRepoRoot(".")
	if err != nil {
		t.Fatalf("repo root: %v", err)
	}
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
	for _, he := range got {
		if len(he.Sum) != 64 { // hex sha256
			t.Errorf("edge %v: bad sha %q", he.Edge, he.Sum)
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
	root, err := FindRepoRoot(".")
	if err != nil {
		t.Fatalf("repo root: %v", err)
	}
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
	// Mutate one .lg input, restore after: only the bundle+lowered go stale,
	// NOT the registrar (a .lg edit must not dirty the Go-annotation output).
	lg := filepath.Join(root, "pkg/rt/core/core.lg")
	orig, _ := os.ReadFile(lg)
	defer os.WriteFile(lg, orig, 0644)
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
