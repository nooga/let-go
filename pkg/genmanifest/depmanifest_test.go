package genmanifest

import (
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
