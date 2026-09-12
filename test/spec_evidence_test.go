package test

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

// specEvidence runs the spec-evidence CLI from the repository root with the
// helper namespaces on the source path, and returns its combined output.
func specEvidence(t *testing.T, args ...string) ([]byte, error) {
	t.Helper()
	cmd := exec.Command(lgBin, append([]string{"scripts/spec-evidence.lg"}, args...)...)
	cmd.Dir = ".."
	cmd.Env = append(os.Environ(), "LG_SOURCE_PATHS=scripts")
	out, err := cmd.CombinedOutput()
	return out, err
}

// TestSpecEvidence runs every gated spec through scripts/spec-evidence.lg
// against the lg binary TestMain built, one Go subtest per spec.
//
// Which specs are gated is NOT decided here: the `specs` subcommand is the
// single rule (carries `@R-` evidence, and no `evidence: skip` in the
// masthead), and the Makefile's spec-lint / spec-evidence / spec-oracle
// targets call the same subcommand. Re-deriving it in Go would let the gate
// and the make targets disagree about what they cover.
//
// See docs/specs/executable-evidence.md.
func TestSpecEvidence(t *testing.T) {
	listing, err := specEvidence(t, "specs", "docs/specs")
	if err != nil {
		t.Fatalf("listing gated specs: %v\n%s", err, listing)
	}
	var specs []string
	for _, line := range strings.Split(string(listing), "\n") {
		if line = strings.TrimSpace(line); line != "" {
			specs = append(specs, line)
		}
	}
	if len(specs) == 0 {
		t.Fatal("no gated spec under docs/specs; the gate would pass vacuously")
	}
	for _, spec := range specs {
		t.Run(baseName(spec), func(t *testing.T) {
			out, err := specEvidence(t, "run", spec)
			if err != nil {
				t.Fatalf("spec evidence failed for %s: %v\n%s", spec, err, out)
			}
			text := string(out)
			if !strings.Contains(text, "cases ok") {
				t.Fatalf("no summary line for %s:\n%s", spec, text)
			}
			// A spec whose blocks all vanished (say, every opener malformed)
			// would otherwise report a green 0/0. `run` refuses structural
			// errors outright, but a spec that tags a requirement and then
			// has nothing to execute is still not evidence.
			if strings.Contains(text, " 0/0 cases ok") {
				t.Fatalf("%s ran no cases:\n%s", spec, text)
			}
		})
	}
}

// baseName is filepath.Base for the forward-slash paths `specs` prints,
// independent of the host separator.
func baseName(path string) string {
	if i := strings.LastIndex(path, "/"); i >= 0 {
		return path[i+1:]
	}
	return path
}
