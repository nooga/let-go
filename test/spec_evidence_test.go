package test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestSpecEvidence runs every docs/specs/*.md that carries tagged evidence
// through scripts/spec-evidence.lg against the lg binary TestMain built. One
// Go subtest per spec; a spec opts out with `evidence: skip` in its masthead,
// the escape hatch for a spec whose engine dependencies are not on this
// branch yet. See docs/specs/executable-evidence.md.
func TestSpecEvidence(t *testing.T) {
	specs, err := filepath.Glob("../docs/specs/*.md")
	if err != nil {
		t.Fatalf("glob docs/specs: %v", err)
	}
	if len(specs) == 0 {
		t.Fatal("no docs/specs/*.md found")
	}
	ran := 0
	for _, spec := range specs {
		text, err := os.ReadFile(spec)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Contains(text, []byte("@R-")) {
			continue
		}
		if bytes.Contains(text, []byte("\nevidence: skip\n")) {
			continue
		}
		rel, err := filepath.Rel("..", spec)
		if err != nil {
			t.Fatal(err)
		}
		ran++
		t.Run(filepath.Base(spec), func(t *testing.T) {
			cmd := exec.Command(lgBin, "scripts/spec-evidence.lg", "run", rel)
			cmd.Dir = ".."
			cmd.Env = append(os.Environ(), "LG_SOURCE_PATHS=scripts")
			out, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("spec evidence failed for %s: %v\n%s", rel, err, out)
			}
			if !strings.Contains(string(out), "cases ok") {
				t.Fatalf("no summary line for %s:\n%s", rel, out)
			}
		})
	}
	if ran == 0 {
		t.Fatal("no spec carried @R- evidence; the gate would pass vacuously")
	}
}
