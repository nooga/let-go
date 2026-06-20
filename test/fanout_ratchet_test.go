package test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// lgBin is the lg binary built once for all fanout-ratchet integration tests.
var lgBin string

// repoRoot is the absolute path to the repository root (tests run from test/).
var repoRoot string

func TestMain(m *testing.M) {
	tmp, err := os.MkdirTemp("", "lg-fanout-*")
	if err != nil {
		panic(err)
	}
	lgBin = filepath.Join(tmp, "lg-bin")
	build := exec.Command("go", "build", "-o", lgBin, "../")
	build.Dir = "."
	if out, err := build.CombinedOutput(); err != nil {
		panic("build lg: " + err.Error() + "\n" + string(out))
	}
	repoRoot, _ = filepath.Abs("..")
	code := m.Run()
	os.RemoveAll(tmp)
	os.Exit(code)
}

// runRatchet runs scripts/fanout-ratchet.lg with args from repo root; returns
// combined output and exit code.
func runRatchet(t *testing.T, args ...string) (string, int) {
	t.Helper()
	full := append([]string{"scripts/fanout-ratchet.lg"}, args...)
	cmd := exec.Command(lgBin, full...)
	cmd.Dir = repoRoot
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	code := 0
	if ee, ok := err.(*exec.ExitError); ok {
		code = ee.ExitCode()
	} else if err != nil {
		t.Fatalf("run ratchet: %v\noutput:\n%s", err, out.String())
	}
	return out.String(), code
}

// writeModule writes <dir>/<module>/<module>.go with exactly `size` bytes
// (size-1 'a's + trailing newline). file-size in the script == chars == bytes.
func writeModule(t *testing.T, treeDir, module string, size int) {
	t.Helper()
	d := filepath.Join(treeDir, module)
	if err := os.MkdirAll(d, 0755); err != nil {
		t.Fatal(err)
	}
	body := strings.Repeat("a", size-1) + "\n"
	if err := os.WriteFile(filepath.Join(d, module+".go"), []byte(body), 0644); err != nil {
		t.Fatal(err)
	}
}

// writeSource writes <dir>/<name>.lg with exactly `size` bytes — the denominator
// of the expansion ratio. The script sums :source-bytes over the --source-dir
// tree, so tests drive both numerator (--tree-dir Go bytes) and denominator
// (--source-dir .lg bytes) to make the ratio deterministic.
func writeSource(t *testing.T, srcDir, name string, size int) {
	t.Helper()
	if err := os.MkdirAll(srcDir, 0755); err != nil {
		t.Fatal(err)
	}
	body := strings.Repeat("a", size-1) + "\n"
	if err := os.WriteFile(filepath.Join(srcDir, name+".lg"), []byte(body), 0644); err != nil {
		t.Fatal(err)
	}
}

func writeBaseline(t *testing.T, path, edn string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(edn), 0644); err != nil {
		t.Fatal(err)
	}
}

func TestFanoutShow(t *testing.T) {
	tree := t.TempDir()
	writeModule(t, tree, "a", 1000)
	src := t.TempDir() // empty source dir keeps the run hermetic and fast
	out, code := runRatchet(t, "show", "--no-regen", "--no-wireup", "--tree-dir", tree, "--source-dir", src)
	if code != 0 {
		t.Fatalf("show exit=%d, want 0\n%s", code, out)
	}
	if !strings.Contains(out, ":total-bytes 1000") {
		t.Errorf("expected :total-bytes 1000 in output:\n%s", out)
	}
	if !strings.Contains(out, "\"a\"") {
		t.Errorf("expected module \"a\" in output:\n%s", out)
	}
}

// A missing tree (or source) directory must yield empty metrics, not an
// interpreter error — the portable os/ls walker has to match the old find
// wrapper's "empty vec if dir missing" contract.
func TestFanoutShowMissingTree(t *testing.T) {
	missingTree := filepath.Join(t.TempDir(), "does-not-exist")
	missingSrc := filepath.Join(t.TempDir(), "also-missing")
	out, code := runRatchet(t, "show", "--no-regen", "--no-wireup", "--tree-dir", missingTree, "--source-dir", missingSrc)
	if code != 0 {
		t.Fatalf("show exit=%d, want 0 (missing tree must not error)\n%s", code, out)
	}
	if strings.Contains(out, "ERROR") || strings.Contains(out, "error:") {
		t.Errorf("missing tree produced an error instead of empty metrics:\n%s", out)
	}
	if !strings.Contains(out, ":total-bytes 0") {
		t.Errorf("expected :total-bytes 0 for a missing tree:\n%s", out)
	}
}

// Ratio baseline: go=2000 over src=1000 → expansion 2000 (go bytes per 1000
// source bytes). The +5% band admits ratios up to 2100.
const fixtureBaseline = `{:total-bytes 2000
 :source-bytes 1000
 :total-loc 2
 :files 2
 :modules {
            "a" {:bytes 1000 :loc 1 :files 1}
            "b" {:bytes 1000 :loc 1 :files 1}
           }}
`

// Current go=2050 over src=1000 → ratio 2050, within the 2100 band → OK.
func TestFanoutCheckWithinBand(t *testing.T) {
	tree := t.TempDir()
	writeModule(t, tree, "a", 1000)
	writeModule(t, tree, "b", 1050) // cur-go 2050
	src := t.TempDir()
	writeSource(t, src, "core", 1000) // cur-src 1000 → ratio 2050 <= 2100
	bl := filepath.Join(t.TempDir(), "base.edn")
	writeBaseline(t, bl, fixtureBaseline)
	out, code := runRatchet(t, "check", "--no-regen", "--no-wireup", "--tree-dir", tree, "--source-dir", src, "--baseline", bl)
	if code != 0 {
		t.Fatalf("check exit=%d, want 0 (within band)\n%s", code, out)
	}
	if !strings.Contains(out, "OK") {
		t.Errorf("expected OK in output:\n%s", out)
	}
}

// Current go=2200 over src=1000 → ratio 2200, beyond the 2100 band → REGRESSION.
func TestFanoutCheckOverBand(t *testing.T) {
	tree := t.TempDir()
	writeModule(t, tree, "a", 1000)
	writeModule(t, tree, "b", 1200) // cur-go 2200
	src := t.TempDir()
	writeSource(t, src, "core", 1000) // cur-src 1000 → ratio 2200 > 2100
	bl := filepath.Join(t.TempDir(), "base.edn")
	writeBaseline(t, bl, fixtureBaseline)
	out, code := runRatchet(t, "check", "--no-regen", "--no-wireup", "--tree-dir", tree, "--source-dir", src, "--baseline", bl)
	if code != 1 {
		t.Fatalf("check exit=%d, want 1 (over band)\n%s", code, out)
	}
	if !strings.Contains(out, "REGRESSION") {
		t.Errorf("expected REGRESSION in output:\n%s", out)
	}
}

// The headline property of the ratio gate: coverage growth — lowering a brand
// new module — adds to BOTH the Go output and the .lg source, so the ratio
// stays flat and the gate does NOT trip, even though absolute Go bytes double.
// (Under the old byte-band gate this was "new modules are exempt"; under the
// ratio gate it generalizes to "proportional growth is exempt".)
const fixtureBaselineSingle = `{:total-bytes 2000
 :source-bytes 1000
 :total-loc 1
 :files 1
 :modules {
            "a" {:bytes 2000 :loc 1 :files 1}
           }}
`

func TestFanoutCheckCoverageGrowthNeutral(t *testing.T) {
	tree := t.TempDir()
	writeModule(t, tree, "a", 2000)
	writeModule(t, tree, "c", 2000) // brand new module, doubles Go bytes to 4000
	src := t.TempDir()
	writeSource(t, src, "a", 1000)
	writeSource(t, src, "c", 1000) // source also doubles to 2000 → ratio stays 2000
	bl := filepath.Join(t.TempDir(), "base.edn")
	writeBaseline(t, bl, fixtureBaselineSingle)
	out, code := runRatchet(t, "check", "--no-regen", "--no-wireup", "--tree-dir", tree, "--source-dir", src, "--baseline", bl)
	if code != 0 {
		t.Fatalf("check exit=%d, want 0 (proportional growth exempt)\n%s", code, out)
	}
	if !strings.Contains(out, "new modules") || !strings.Contains(out, "c") {
		t.Errorf("expected new module c reported:\n%s", out)
	}
}

// A legacy baseline written before the expansion gate has no :source-bytes. The
// ratio is then uncomputable, so check must treat it as in-band (exit 0) and
// say so — never crash multiplying nil.
const fixtureBaselineLegacy = `{:total-bytes 2000
 :total-loc 2
 :files 2
 :modules {
            "a" {:bytes 1000 :loc 1 :files 1}
            "b" {:bytes 1000 :loc 1 :files 1}
           }}
`

func TestFanoutCheckLegacyBaseline(t *testing.T) {
	tree := t.TempDir()
	writeModule(t, tree, "a", 1000)
	writeModule(t, tree, "b", 5000) // would be way over band if it were comparable
	src := t.TempDir()
	writeSource(t, src, "core", 1000)
	bl := filepath.Join(t.TempDir(), "base.edn")
	writeBaseline(t, bl, fixtureBaselineLegacy)
	out, code := runRatchet(t, "check", "--no-regen", "--no-wireup", "--tree-dir", tree, "--source-dir", src, "--baseline", bl)
	if code != 0 {
		t.Fatalf("check exit=%d, want 0 (legacy baseline uncomparable, not a regression)\n%s", code, out)
	}
	if !strings.Contains(out, "source-bytes") {
		t.Errorf("expected a notice that the baseline lacks source-bytes:\n%s", out)
	}
}

func TestFanoutUpdateTightens(t *testing.T) {
	tree := t.TempDir()
	writeModule(t, tree, "a", 800) // shrank from 1000
	writeModule(t, tree, "b", 1000)
	src := t.TempDir()
	bl := filepath.Join(t.TempDir(), "base.edn")
	writeBaseline(t, bl, fixtureBaseline)
	_, code := runRatchet(t, "update", "--no-regen", "--no-wireup", "--tree-dir", tree, "--source-dir", src, "--baseline", bl)
	if code != 0 {
		t.Fatalf("update exit=%d, want 0", code)
	}
	got, _ := os.ReadFile(bl)
	s := string(got)
	if !strings.Contains(s, ":total-bytes 1800") {
		t.Errorf("expected total-bytes tightened to 1800:\n%s", s)
	}
	if !strings.Contains(s, "\"a\" {:bytes 800") {
		t.Errorf("expected a tightened to 800:\n%s", s)
	}
}

func TestFanoutUpdateFoldsNewModule(t *testing.T) {
	tree := t.TempDir()
	writeModule(t, tree, "a", 1000)
	writeModule(t, tree, "b", 1000)
	writeModule(t, tree, "c", 500) // new module
	src := t.TempDir()
	bl := filepath.Join(t.TempDir(), "base.edn")
	writeBaseline(t, bl, fixtureBaseline)
	_, code := runRatchet(t, "update", "--no-regen", "--no-wireup", "--tree-dir", tree, "--source-dir", src, "--baseline", bl)
	if code != 0 {
		t.Fatalf("update exit=%d, want 0", code)
	}
	got, _ := os.ReadFile(bl)
	s := string(got)
	if !strings.Contains(s, "\"c\" {:bytes 500") {
		t.Errorf("expected new module c folded in:\n%s", s)
	}
	if !strings.Contains(s, ":total-bytes 2500") {
		t.Errorf("expected total-bytes 2500 (1000+1000+500):\n%s", s)
	}
}

func TestFanoutUpdateDeterministic(t *testing.T) {
	tree := t.TempDir()
	writeModule(t, tree, "b", 1000)
	writeModule(t, tree, "a", 1000)
	writeModule(t, tree, "c", 1000)
	src := t.TempDir()
	bl := filepath.Join(t.TempDir(), "base.edn")
	writeBaseline(t, bl, fixtureBaseline)
	runRatchet(t, "update", "--no-regen", "--no-wireup", "--tree-dir", tree, "--source-dir", src, "--baseline", bl)
	first, _ := os.ReadFile(bl)
	runRatchet(t, "update", "--no-regen", "--no-wireup", "--tree-dir", tree, "--source-dir", src, "--baseline", bl)
	second, _ := os.ReadFile(bl)
	if string(first) != string(second) {
		t.Errorf("update output not byte-stable:\nfirst:\n%s\nsecond:\n%s", first, second)
	}
	s := string(first)
	ia, ib, ic := strings.Index(s, "\"a\""), strings.Index(s, "\"b\""), strings.Index(s, "\"c\"")
	if !(ia < ib && ib < ic) {
		t.Errorf("modules not sorted a<b<c:\n%s", s)
	}
}
