package genmanifest

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func generateScript(t *testing.T) string {
	t.Helper()
	root, err := FindRepoRoot(".")
	if err != nil {
		t.Fatalf("repo root: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(root, "scripts/generate.lg"))
	if err != nil {
		t.Fatalf("read generate.lg: %v", err)
	}
	return string(data)
}

func TestGenerateScriptRequeriesStalenessAfterIRGeneration(t *testing.T) {
	script := generateScript(t)
	upstream := strings.Index(script, `(gen-lisp "pkg/ir/ir_data.lg"`)
	lgbgen := strings.Index(script, `(if (stale? "pkg/rt/core_compiled.lgb"`)
	if upstream < 0 || lgbgen < 0 || upstream >= lgbgen {
		t.Fatal("could not locate IR generation followed by the lgbgen decision")
	}
	between := script[upstream:lgbgen]
	if !strings.Contains(between, `"-needs-generation"`) {
		t.Fatal("staleness must be requeried after IR generation and before the lgbgen decision")
	}
}

func TestLGBGenOutputsTrackImplementationDependencies(t *testing.T) {
	root, err := FindRepoRoot(".")
	if err != nil {
		t.Fatalf("repo root: %v", err)
	}
	edges, err := Edges(root)
	if err != nil {
		t.Fatalf("Edges: %v", err)
	}

	outputs := []string{"pkg/rt/core_compiled.lgb", "pkg/rt/core_go_lowered/"}
	dirs := []string{"cmd/lgbgen/", "pkg/compiler/", "pkg/bytecode/", "pkg/ir/", "pkg/resolver/", "pkg/vm/", "pkg/rt/"}
	for _, output := range outputs {
		for _, dir := range dirs {
			found := false
			for _, edge := range edges {
				if edge.Output == output && edge.Kind == "generator" && strings.HasPrefix(edge.Input, dir) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("%s has no generator dependency under %s", output, dir)
			}
		}
	}
}

func TestIROutputsAttributeGenerateScript(t *testing.T) {
	root, err := FindRepoRoot(".")
	if err != nil {
		t.Fatalf("repo root: %v", err)
	}
	edges, err := Edges(root)
	if err != nil {
		t.Fatalf("Edges: %v", err)
	}

	outputs := map[string]bool{
		"pkg/ir/op_generated.go":           true,
		"pkg/rt/ir_bridge_generated.go":    true,
		"pkg/rt/core/ir/data/generated.lg": true,
	}
	seenScript := map[string]bool{}
	seenLGBGen := map[string]bool{}
	for _, edge := range edges {
		if !outputs[edge.Output] || edge.Kind != "generator" {
			continue
		}
		if edge.Input == "scripts/generate.lg" {
			seenScript[edge.Output] = true
		}
		if strings.HasPrefix(edge.Input, "cmd/lgbgen/") {
			seenLGBGen[edge.Output] = true
		}
	}
	for output := range outputs {
		if !seenScript[output] {
			t.Errorf("%s lacks scripts/generate.lg as its generator", output)
		}
		if seenLGBGen[output] {
			t.Errorf("%s incorrectly attributes generation to cmd/lgbgen", output)
		}
	}
}

func TestGenerateScriptRejectsFailedStalenessQuery(t *testing.T) {
	script := generateScript(t)
	start := strings.Index(script, ";; --- selective regen:")
	end := strings.Index(script, ";; --- provenance ---")
	if start < 0 || end < 0 || start >= end {
		t.Fatal("could not locate selective regeneration section")
	}
	query := script[start:end]
	for _, required := range []string{"(:exit", "(:err", "die"} {
		if !strings.Contains(query, required) {
			t.Errorf("staleness query must inspect/handle %q before treating stdout as a stale set", required)
		}
	}
}
