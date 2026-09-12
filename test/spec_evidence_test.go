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
		if hasEvidenceSkip(text) {
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

// hasEvidenceSkip reports whether a spec's frontmatter masthead carries the
// `evidence: skip` opt-out. Deliberately scoped to the masthead — the lines
// between the leading `---` and the next `---` — so a spec that *documents*
// the escape hatch inside a fenced example does not silently skip its own
// gate. A file with no masthead never opts out.
func hasEvidenceSkip(text []byte) bool {
	lines := strings.Split(string(text), "\n")
	if len(lines) == 0 || strings.TrimRight(lines[0], "\r") != "---" {
		return false
	}
	for _, line := range lines[1:] {
		line = strings.TrimRight(line, "\r")
		if line == "---" {
			return false
		}
		if line == "evidence: skip" {
			return true
		}
	}
	return false
}

func TestHasEvidenceSkip(t *testing.T) {
	cases := []struct {
		name string
		text string
		want bool
	}{
		{
			name: "masthead opt-out",
			text: "---\nstatus: draft\nevidence: skip\n---\n\n# S\n\n[R-x]\n",
			want: true,
		},
		{
			name: "only inside a fenced example after the masthead",
			text: "---\nstatus: draft\n---\n\n# S\n\nOpt out like this:\n\n```yaml\n---\nevidence: skip\n---\n```\n",
			want: false,
		},
		{
			name: "no masthead",
			text: "# S\n\nevidence: skip\n",
			want: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := hasEvidenceSkip([]byte(tc.text)); got != tc.want {
				t.Fatalf("hasEvidenceSkip = %v, want %v", got, tc.want)
			}
		})
	}
}
