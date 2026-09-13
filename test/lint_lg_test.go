package test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func TestLintDisplayNamesCoverBuiltInKinds(t *testing.T) {
	script, err := os.ReadFile(filepath.Join(repoRoot, "scripts", "lint.lg"))
	if err != nil {
		t.Fatal(err)
	}
	catalog, err := os.ReadFile(filepath.Join(repoRoot, "scripts", "lint-code-rules.edn"))
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	scriptsDir := filepath.Join(dir, "scripts")
	if err := os.Mkdir(scriptsDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(scriptsDir, "lint-code-rules.edn"), catalog, 0644); err != nil {
		t.Fatal(err)
	}
	// The probe runs after the normal driver, so it exercises the actual display
	// function and catalog loader without adding a test-only production flag.
	probe := `
(doseq [kind (concat [:commented-out-code :restatement :comment-divider
                      :duplicated-comment :comment-density-outlier
                      :devlog-comment :comment-churn]
                     (map :kind (load-code-rules)))]
  (println (str "DISPLAY\t" (name kind) "\t"
                (display-name kind (load-code-rules)))))
`
	if err := os.WriteFile(filepath.Join(scriptsDir, "lint.lg"), append(script, []byte(probe)...), 0644); err != nil {
		t.Fatal(err)
	}
	sourceDir := filepath.Join(dir, "empty-source")
	if err := os.Mkdir(sourceDir, 0755); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(lgBin, filepath.Join(scriptsDir, "lint.lg"), sourceDir)
	cmd.Dir = repoRoot
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("lint label probe: %v\n%s", err, out)
	}
	labels := map[string]string{}
	for _, line := range strings.Split(string(out), "\n") {
		if !strings.HasPrefix(line, "DISPLAY\t") {
			continue
		}
		parts := strings.SplitN(line, "\t", 3)
		if len(parts) != 3 {
			t.Fatalf("malformed label probe line %q", line)
		}
		labels[parts[1]] = parts[2]
	}
	wantKinds := map[string]bool{
		"commented-out-code": true, "restatement": true, "comment-divider": true,
		"duplicated-comment": true, "comment-density-outlier": true,
		"devlog-comment": true, "comment-churn": true,
	}
	for _, match := range regexp.MustCompile(`(?m)^\s*:kind\s+:([a-z][a-z-]*)`).FindAllStringSubmatch(string(catalog), -1) {
		wantKinds[match[1]] = true
	}
	for _, entry := range strings.Split(string(catalog), "\n\n") {
		if !regexp.MustCompile(`(?m)^\s*:kind\s+:`).MatchString(entry) {
			continue
		}
		if !regexp.MustCompile(`(?m)^\s*:label\s+"[^"\s][^"]*"`).MatchString(entry) {
			t.Errorf("catalog entry has no nonblank :label: %s", entry)
		}
	}
	if len(labels) != len(wantKinds) {
		t.Errorf("got %d distinct labels, want %d; labels=%v", len(labels), len(wantKinds), labels)
	}
	wantExact := map[string]string{
		"restatement":             "comment repeats code",
		"comment-divider":         "decorative section divider",
		"comment-density-outlier": "unusually comment-heavy definition",
		"devlog-comment":          "development-note phrase",
		"eta-expansion":           "argument-forwarding wrapper",
		"composable-accessor":     "nested sequence accessor",
	}
	seen := map[string]string{}
	for kind := range wantKinds {
		label := labels[kind]
		if strings.TrimSpace(label) == "" {
			t.Errorf("%s has no display label", kind)
		}
		if regexp.MustCompile(`\bR[1-6]\b`).MatchString(label) || strings.Contains(label, kind) || strings.Contains(label, "eta expansion") {
			t.Errorf("%s has a low-context failure: %q", kind, label)
		}
		if want, ok := wantExact[kind]; ok && label != want {
			t.Errorf("%s label = %q, want %q", kind, label, want)
		}
		if other, exists := seen[label]; exists && other != kind {
			t.Errorf("%s and %s share display label %q", kind, other, label)
		}
		seen[label] = kind
	}
}

func assertLintReaderText(t *testing.T, output string) {
	t.Helper()
	for _, forbidden := range []string{"see R1's section comment", "see the code-verbosity section comment"} {
		if strings.Contains(output, forbidden) {
			t.Errorf("report refers reader to source: %q\n%s", forbidden, output)
		}
	}
	if regexp.MustCompile(`\bR[1-6]\b`).MatchString(output) {
		t.Errorf("report exposes internal rule numbers:\n%s", output)
	}
	if regexp.MustCompile(`\[(?:commented-out-code|restatement|comment-divider|duplicated-comment|comment-density-outlier|devlog-comment|comment-churn)\]`).MatchString(output) {
		t.Errorf("report exposes an internal finding kind:\n%s", output)
	}
}

func TestLintReportLabelsOnEmptySource(t *testing.T) {
	dir := t.TempDir()
	out, err := exec.Command(lgBin, filepath.Join(repoRoot, "scripts", "lint.lg"), dir).CombinedOutput()
	if err != nil {
		t.Fatalf("lint empty source: %v\n%s", err, out)
	}
	text := string(out)
	assertLintReaderText(t, text)
	for _, want := range []string{
		"code left in a comment", "comment repeats code", "decorative section divider",
		"repeated comment block", "unusually comment-heavy definition: skipped",
		"development-note phrase", "comment additions relative to code: skipped",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("report missing reader label %q:\n%s", want, text)
		}
	}
}

// scripts/lint.lg flags comments that narrate the change rather than describe
// the code. Its skip markers are the part most likely to rot silently: a marker
// that can never match costs nothing at runtime and shows up only as a stray
// finding, which is how "//go:" survived review (nooga/let-go#836). The fixture
// below pins one line per marker plus the matching rules the phrase list leans
// on, so a change to the stripping or normalization has to move a line number.

// lintFixtureGo pairs each comment with the line it lands on; want-hit lines are
// listed in lintWantGo. Keep the two in sync when editing.
const lintFixtureGo = `package fixture

//go:generate run this pr adds a flag
func A() {}

// Code generated by tool. this pr adds a flag
func B() {}

//lgbgen: this pr adds a flag
func C() {}

// lint:ignore — quoting "this pr" is what the escape hatch is for
func D() {}

// Copyright 2026 — this pr adds a flag
func E() {}

// this pr adds a flag
func F() {}

// this process adds a flag
func G() {}

// the arms are not what they
// used to be before the split
func H() {}
`

// Line 18 is the plain devlog comment; 24 is the block that wraps a phrase
// across two lines and reports at the line it starts on.
var lintWantGo = []int{18, 24}

const lintFixtureLg = `;; this commit renames the arms
(def a 1)

;; lint:ignore — this commit is quoted on purpose
(def b 2)

;; the split is done; this is the current shape
(def c 3)
`

var lintWantLg = []int{1}

func TestLintLgSkipMarkersAndPhraseMatching(t *testing.T) {
	dir := t.TempDir()
	goFixture := filepath.Join(dir, "fixture.go")
	lgFixture := filepath.Join(dir, "fixture.lg")
	if err := os.WriteFile(goFixture, []byte(lintFixtureGo), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(lgFixture, []byte(lintFixtureLg), 0644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(lgBin, filepath.Join(repoRoot, "scripts", "lint.lg"), dir)
	cmd.Dir = repoRoot
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("lint.lg: %v\n%s", err, out)
	}
	assertLintReaderText(t, string(out))

	got := map[string][]int{}
	for _, line := range strings.Split(string(out), "\n") {
		if !strings.Contains(line, "[development-note phrase]") {
			continue
		}
		loc := strings.Fields(line)[0]
		i := strings.LastIndex(loc, ":")
		n, convErr := strconv.Atoi(loc[i+1:])
		if convErr != nil {
			t.Fatalf("unparsable finding location %q in:\n%s", loc, out)
		}
		got[filepath.Base(loc[:i])] = append(got[filepath.Base(loc[:i])], n)
	}

	for file, want := range map[string][]int{"fixture.go": lintWantGo, "fixture.lg": lintWantLg} {
		if !equalInts(got[file], want) {
			t.Errorf("%s: findings at lines %v, want %v\n%s", file, got[file], want, out)
		}
	}
}

// R1 — commented-out code (`.lg` only; Go is explicitly deferred, see
// scripts/lint.lg's own R1 section comment). A comment whose body reads as a
// balanced, delimited form is flagged; a bare identifier/literal and ordinary
// prose (even prose that happens to contain a parenthetical aside without a
// full balanced code shape) are not.
const lintFixtureR1 = `;; (inc counter)
(def a 1)

;; n
(def b 2)

;; increment the counter for retries
(def c 3)

;; 42
(def d 4)
`

var lintWantR1 = []int{1}

// A self-referential usage example — a comment that shows the calling
// convention of the very declaration it precedes (its parsed head symbol
// equals a token on the next source line) — must NOT flag: it is
// documentation, not dead code, even though it parses as a balanced form.
const lintFixtureR1ExampleDoc = `;; (frob widget & opts)
(defn frob [widget & opts]
  widget)
`

func TestLintR1CommentedOutCode(t *testing.T) {
	dir := t.TempDir()
	lgFixture := filepath.Join(dir, "fixture.lg")
	if err := os.WriteFile(lgFixture, []byte(lintFixtureR1), 0644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(lgBin, filepath.Join(repoRoot, "scripts", "lint.lg"), dir)
	cmd.Dir = repoRoot
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("lint.lg: %v\n%s", err, out)
	}
	assertLintReaderText(t, string(out))

	got := []int{}
	for _, line := range strings.Split(string(out), "\n") {
		if !strings.Contains(line, "[code left in a comment]") {
			continue
		}
		loc := strings.Fields(line)[0]
		// loc is "path:line-end"
		loc = strings.SplitN(loc, ":", 2)[1]
		loc = strings.SplitN(loc, "-", 2)[0]
		n, convErr := strconv.Atoi(loc)
		if convErr != nil {
			t.Fatalf("unparsable finding location %q in:\n%s", loc, out)
		}
		got = append(got, n)
	}

	if !equalInts(got, lintWantR1) {
		t.Errorf("commented-out-code findings at lines %v, want %v\n%s", got, lintWantR1, out)
	}

	// A self-referential usage-example doc comment must not flag.
	dir2 := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir2, "fixture.lg"), []byte(lintFixtureR1ExampleDoc), 0644); err != nil {
		t.Fatal(err)
	}
	out3, err := exec.Command(lgBin, filepath.Join(repoRoot, "scripts", "lint.lg"), dir2).CombinedOutput()
	if err != nil {
		t.Fatalf("lint.lg: %v\n%s", err, out3)
	}
	if strings.Contains(string(out3), "[code left in a comment]") {
		t.Errorf("a self-referential usage example must not flag as commented-out code:\n%s", out3)
	}

	// Go source must never be flagged for R1: the rule is deferred there.
	goFixture := filepath.Join(dir, "other.go")
	goSrc := "package fixture\n\n// x := computeSomething(a, b)\nfunc F() {}\n"
	if err := os.WriteFile(goFixture, []byte(goSrc), 0644); err != nil {
		t.Fatal(err)
	}
	out2, err := exec.Command(lgBin, filepath.Join(repoRoot, "scripts", "lint.lg"), dir).CombinedOutput()
	if err != nil {
		t.Fatalf("lint.lg: %v\n%s", err, out2)
	}
	for _, line := range strings.Split(string(out2), "\n") {
		if strings.Contains(line, "[code left in a comment]") && strings.Contains(line, "other.go") {
			t.Errorf("R1 must not fire on Go source (deferred): %s", line)
		}
	}
}

// R2 — restatement: the comment's content words overlap heavily with the
// identifiers of the form immediately following it, and it adds nothing that
// isn't already spelled out by that form. A comment giving a reason, even one
// naming the same identifiers, is not flagged because it also carries extra
// content words.
const lintFixtureR2 = `;; set counter
(set counter val)

;; set counter because writer holds the exclusive lock
(set counter val)

// --- Func roundtrip ---
func TestFuncRoundtrip(t *testing.T) {}
`

// R2 also classifies a restatement-shaped comment separately when its text is
// a section-divider (bracketed by runs of punctuation like dashes/equals,
// with no sentence structure): that's conventional house style, not slop, so
// it must NOT be reported as [comment repeats code] — but it must still show up under
// its own [decorative section divider] kind so it stays visible.
func TestLintR2Restatement(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "fixture.lg"), []byte(lintFixtureR2), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "fixture.go"), []byte("package fixture\n\n"+lintFixtureR2[strings.Index(lintFixtureR2, "// ---"):]), 0644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(lgBin, filepath.Join(repoRoot, "scripts", "lint.lg"), dir)
	cmd.Dir = repoRoot
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("lint.lg: %v\n%s", err, out)
	}
	assertLintReaderText(t, string(out))
	outStr := string(out)

	got := []int{}
	for _, line := range strings.Split(outStr, "\n") {
		if !strings.Contains(line, "[comment repeats code]") {
			continue
		}
		loc := strings.Fields(line)[0]
		loc = strings.SplitN(loc, ":", 2)[1]
		loc = strings.SplitN(loc, "-", 2)[0]
		n, convErr := strconv.Atoi(loc)
		if convErr != nil {
			t.Fatalf("unparsable finding location %q in:\n%s", loc, outStr)
		}
		got = append(got, n)
	}

	want := []int{1}
	if !equalInts(got, want) {
		t.Errorf("restatement findings at lines %v, want %v\n%s", got, want, outStr)
	}

	if !strings.Contains(outStr, "[decorative section divider]") {
		t.Errorf("section-divider comment should be reported under its own kind:\n%s", outStr)
	}
	if strings.Contains(outStr, "section-divider,") {
		t.Errorf("divider evidence repeats the internal kind name:\n%s", outStr)
	}
	for _, line := range strings.Split(outStr, "\n") {
		if strings.Contains(line, "[comment repeats code]") && strings.Contains(line, "Func roundtrip") {
			t.Errorf("section-divider must not be reported as [comment repeats code]: %s", line)
		}
	}
}

// R3 — duplicated comment: the same normalized comment text appears in three
// or more places across the corpus. Twice is not enough; a licence header
// (already a skip-marker) is never counted even if repeated many times.
func TestLintR3DuplicatedComment(t *testing.T) {
	dir := t.TempDir()
	files := map[string]string{
		"a.lg": ";; Copyright 2026 Example — do not remove\n\n;; retries three times before giving up\n(def a 1)\n",
		"b.lg": ";; Copyright 2026 Example — do not remove\n\n;; retries three times before giving up\n(def b 2)\n",
		"c.lg": ";; Copyright 2026 Example — do not remove\n\n;; retries three times before giving up\n(def c 3)\n",
		"d.lg": ";; Copyright 2026 Example — do not remove\n\n;; seen only twice, should not flag\n(def d 4)\n",
		"e.lg": ";; Copyright 2026 Example — do not remove\n\n;; seen only twice, should not flag\n(def e 5)\n",
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	out, err := exec.Command(lgBin, filepath.Join(repoRoot, "scripts", "lint.lg"), dir).CombinedOutput()
	if err != nil {
		t.Fatalf("lint.lg: %v\n%s", err, out)
	}
	assertLintReaderText(t, string(out))

	dupCount, copyrightHit, twiceHit := 0, false, false
	for _, line := range strings.Split(string(out), "\n") {
		if !strings.Contains(line, "[repeated comment block]") {
			continue
		}
		dupCount++
		if strings.Contains(line, "copyright") || strings.Contains(strings.ToLower(line), "do not remove") {
			copyrightHit = true
		}
		if strings.Contains(line, "seen only twice") {
			twiceHit = true
		}
	}

	if dupCount != 3 {
		t.Errorf("duplicated-comment findings = %d, want 3 (one per occurrence of the 3x comment)\n%s", dupCount, out)
	}
	if copyrightHit {
		t.Errorf("licence header must never be flagged as duplicated:\n%s", out)
	}
	if twiceHit {
		t.Errorf("a comment seen only twice must not flag:\n%s", out)
	}
}

// R3 must not flag idiomatic Go interface-implementation boilerplate: a
// one-line doc comment that names the method it precedes (the standard
// "MethodName implements Interface." convention), repeated once per type
// implementing that method. A genuinely duplicated, non-boilerplate comment
// in the same corpus must still flag.
func TestLintR3ExcludesInterfaceBoilerplate(t *testing.T) {
	dir := t.TempDir()
	mk := func(typeName string) string {
		return fmt.Sprintf("package fixture\n\n// Meta implements IMeta.\nfunc (t *%s) Meta() {}\n\n// retry on transient failure before giving up\nfunc (t *%s) Other() {}\n",
			typeName, typeName)
	}
	for _, name := range []string{"A", "B", "C"} {
		if err := os.WriteFile(filepath.Join(dir, name+".go"), []byte(mk(name)), 0644); err != nil {
			t.Fatal(err)
		}
	}

	out, err := exec.Command(lgBin, filepath.Join(repoRoot, "scripts", "lint.lg"), dir).CombinedOutput()
	if err != nil {
		t.Fatalf("lint.lg: %v\n%s", err, out)
	}
	assertLintReaderText(t, string(out))
	outStr := string(out)

	if strings.Contains(outStr, "Meta implements IMeta") {
		t.Errorf("interface-implementation boilerplate must be excluded from duplicated-comment:\n%s", outStr)
	}
	if !strings.Contains(outStr, "retry on transient failure") {
		t.Errorf("a genuine non-boilerplate duplicate must still be flagged:\n%s", outStr)
	}
}

// R4 — comment density outlier: per definition, INTERIOR comment-lines /
// interior-total-lines (the leading doc-comment block is excluded from both
// the numerator and the denominator — a doc comment on an exported
// declaration is documentation, not density abuse; see scripts/lint.lg's R4
// section comment for why "exclude leading docs" was chosen over
// "unexported only"). Flagged when it exceeds the 95th percentile computed
// over the corpus itself.
//
// This fixture builds 20 definitions with a strictly increasing,
// hand-computable ratio: each definition has a fixed 3-line leading doc
// comment (must NOT affect the ratio) plus Ci INTERIOR comment lines before
// a 1-line body, Ci from 3 to 22, so every definition clears the
// min-definition-lines guard on its interior span alone. With nearest-rank
// on 20 values, rank = ceil(0.95*20) = 19, i.e. the 19th smallest (Ci=21,
// ratio 21/23=0.9130). Only the strict maximum (Ci=22, ratio 22/24=0.9167)
// beats that threshold, so exactly one definition must be flagged.
func buildLintFixtureR4() (string, int) {
	var b strings.Builder
	lastDeclLine := 0
	line := 1
	for ci := 3; ci <= 22; ci++ {
		if ci > 3 {
			b.WriteString("\n")
			line++
		}
		b.WriteString(";; doc line 1\n;; doc line 2\n;; doc line 3\n")
		line += 3
		lastDeclLine = line
		fmt.Fprintf(&b, "(defn f%d [x]\n", ci)
		line++
		for c := 1; c <= ci; c++ {
			fmt.Fprintf(&b, "  ;; c%d-%d\n", ci, c)
			line++
		}
		b.WriteString("  x)\n")
		line++
	}
	return b.String(), lastDeclLine
}

func TestLintR4DensityOutlier(t *testing.T) {
	dir := t.TempDir()
	src, lastDeclLine := buildLintFixtureR4()
	lgFixture := filepath.Join(dir, "fixture.lg")
	if err := os.WriteFile(lgFixture, []byte(src), 0644); err != nil {
		t.Fatal(err)
	}

	out, err := exec.Command(lgBin, filepath.Join(repoRoot, "scripts", "lint.lg"), dir).CombinedOutput()
	if err != nil {
		t.Fatalf("lint.lg: %v\n%s", err, out)
	}
	assertLintReaderText(t, string(out))
	if !strings.Contains(string(out), "unusually comment-heavy definition finding(s) (95th-percentile threshold=") {
		t.Errorf("calibrated density summary lacks reader label or threshold:\n%s", out)
	}

	got := []int{}
	for _, l := range strings.Split(string(out), "\n") {
		if !strings.Contains(l, "[unusually comment-heavy definition]") {
			continue
		}
		loc := strings.Fields(l)[0]
		loc = strings.SplitN(loc, ":", 2)[1]
		loc = strings.SplitN(loc, "-", 2)[0]
		n, convErr := strconv.Atoi(loc)
		if convErr != nil {
			t.Fatalf("unparsable finding location %q in:\n%s", loc, out)
		}
		got = append(got, n)
	}

	want := []int{lastDeclLine}
	if !equalInts(got, want) {
		t.Errorf("comment-density-outlier findings at lines %v, want %v (the Ci=22 definition's decl line, doc excluded)\n%s", got, want, out)
	}

	// The leading doc comment itself must never be part of the flagged
	// range: the reported :line must be the declaration line, not the doc's.
	for _, l := range strings.Split(string(out), "\n") {
		if strings.Contains(l, "[unusually comment-heavy definition]") && strings.Contains(l, "doc line") {
			t.Errorf("leading doc comment text leaked into a density-outlier finding: %s", l)
		}
	}
}

func TestLintR4GoFunctionBodies(t *testing.T) {
	dir := t.TempDir()
	var src strings.Builder
	src.WriteString("package fixture\n\n")
	for i := 0; i < 19; i++ {
		fmt.Fprintf(&src, "func F%d() {\n\t// ordinary interior comment\n\tx := 1\n\t_ = x\n}\n\n", i)
	}
	lastDeclLine := strings.Count(src.String(), "\n") + 1
	src.WriteString("func F19(\n\tx int,\n) int {\n")
	for i := 0; i < 4; i++ {
		fmt.Fprintf(&src, "\t// dense interior comment %d\n", i)
	}
	src.WriteString("\t_ = x\n\treturn x\n}\n")
	path := filepath.Join(dir, "fixture.go")
	if err := os.WriteFile(path, []byte(src.String()), 0644); err != nil {
		t.Fatal(err)
	}

	out, err := exec.Command(lgBin, filepath.Join(repoRoot, "scripts", "lint.lg"), dir).CombinedOutput()
	if err != nil {
		t.Fatalf("lint.lg: %v\n%s", err, out)
	}
	if strings.Contains(string(out), "only 0 definition(s) found") {
		t.Fatalf("R4 did not measure Go functions:\n%s", out)
	}
	want := fmt.Sprintf("%s:%d-", path, lastDeclLine)
	var findings []string
	for _, line := range strings.Split(string(out), "\n") {
		if strings.Contains(line, "[unusually comment-heavy definition]") {
			findings = append(findings, line)
		}
	}
	if len(findings) != 1 || !strings.HasPrefix(findings[0], want) {
		t.Errorf("R4 findings = %v, want only dense Go function at %s:\n%s", findings, want, out)
	}
}

func TestLintR4MalformedGoFailsLoudly(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "broken.go")
	if err := os.WriteFile(path, []byte("package p\nfunc Broken( {\n"), 0644); err != nil {
		t.Fatal(err)
	}
	out, err := exec.Command(lgBin, filepath.Join(repoRoot, "scripts", "lint.lg"), "--edn", path).CombinedOutput()
	if err == nil {
		t.Fatalf("malformed Go must fail lint explicitly:\n%s", out)
	}
	if !strings.Contains(string(out), "cannot parse Go file") {
		t.Errorf("missing parse error in lint output:\n%s", out)
	}
}

// --edn / --gate: machine-readable output and opt-in exit codes.
func TestLintEdnAndGate(t *testing.T) {
	dir := t.TempDir()
	// One R1 (commented-out code) finding and nothing else interesting.
	src := ";; (inc counter)\n(def a 1)\n"
	if err := os.WriteFile(filepath.Join(dir, "fixture.lg"), []byte(src), 0644); err != nil {
		t.Fatal(err)
	}

	// --edn: stdout must be exactly one EDN vector of finding maps.
	ednOut, err := exec.Command(lgBin, filepath.Join(repoRoot, "scripts", "lint.lg"), "--edn", dir).CombinedOutput()
	if err != nil {
		t.Fatalf("lint.lg --edn: %v\n%s", err, ednOut)
	}
	trimmed := strings.TrimSpace(string(ednOut))
	if !strings.HasPrefix(trimmed, "[") || !strings.HasSuffix(trimmed, "]") {
		t.Fatalf("--edn output is not a single EDN vector:\n%s", trimmed)
	}
	if !strings.Contains(trimmed, ":commented-out-code") {
		t.Errorf("--edn output missing the R1 finding:\n%s", trimmed)
	}
	if !strings.Contains(trimmed, ":line") || !strings.Contains(trimmed, ":end") {
		t.Errorf("--edn findings must carry :line and :end:\n%s", trimmed)
	}

	// --gate on a kind with findings exits non-zero.
	gateCmd := exec.Command(lgBin, filepath.Join(repoRoot, "scripts", "lint.lg"),
		"--gate", "commented-out-code", dir)
	gateOut, gateErr := gateCmd.CombinedOutput()
	if gateErr == nil {
		t.Fatalf("--gate commented-out-code should exit non-zero when R1 has findings:\n%s", gateOut)
	}

	// --gate on a kind with no findings exits zero.
	cleanCmd := exec.Command(lgBin, filepath.Join(repoRoot, "scripts", "lint.lg"),
		"--gate", "duplicated-comment", dir)
	cleanOut, cleanErr := cleanCmd.CombinedOutput()
	if cleanErr != nil {
		t.Fatalf("--gate duplicated-comment should exit zero (no such findings): %v\n%s", cleanErr, cleanOut)
	}

	// Default mode (no --gate) still exits zero even though R1 has findings.
	defaultCmd := exec.Command(lgBin, filepath.Join(repoRoot, "scripts", "lint.lg"), dir)
	defaultOut, defaultErr := defaultCmd.CombinedOutput()
	if defaultErr != nil {
		t.Fatalf("default mode must exit zero regardless of findings: %v\n%s", defaultErr, defaultOut)
	}

	// Gating on the R5 heuristic kind must never fail the run, even though
	// it is not a real gate target.
	r5GateCmd := exec.Command(lgBin, filepath.Join(repoRoot, "scripts", "lint.lg"),
		"--gate", "devlog-comment", dir)
	r5Out, r5Err := r5GateCmd.CombinedOutput()
	if r5Err != nil {
		t.Fatalf("--gate devlog-comment must not gate (R5 is heuristic-only): %v\n%s", r5Err, r5Out)
	}
	assertLintReaderText(t, string(r5Out))
	if !strings.Contains(string(r5Out), "--gate ignored development-note phrase") ||
		!strings.Contains(string(r5Out), "development-note phrases use a heuristic") {
		t.Errorf("ignored phrase gate lacks a readable explanation:\n%s", r5Out)
	}
}

// R6 — comment churn for a revision range: comment-lines-added / code-lines-
// added within the range, compared against the corpus's own current
// comment/code ratio (the "baseline" — a snapshot proxy for typical density,
// not a full historical per-commit walk; see scripts/lint.lg's R6 section
// comment and the calibration report for why).
func gitRun(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(),
		"GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t.com",
		"GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t.com")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return strings.TrimSpace(string(out))
}

func TestLintR6CommentChurn(t *testing.T) {
	dir := t.TempDir()
	gitRun(t, dir, "init", "-q")
	gitRun(t, dir, "config", "commit.gpgsign", "false")

	// C1: baseline — 1 comment line, 9 code lines (package decl + 8 funcs).
	c1 := "package fixture\n\n// base comment\nfunc A() {}\nfunc B() {}\nfunc C() {}\nfunc D() {}\nfunc E() {}\nfunc F() {}\nfunc G() {}\nfunc H() {}\n"
	if err := os.WriteFile(filepath.Join(dir, "a.go"), []byte(c1), 0644); err != nil {
		t.Fatal(err)
	}
	gitRun(t, dir, "add", "a.go")
	gitRun(t, dir, "commit", "-q", "-m", "c1")
	c1sha := gitRun(t, dir, "rev-parse", "HEAD")

	// C2: adds 4 comment lines and 1 code line — a much higher ratio (4.0)
	// than what the FINAL corpus (C1+C2 combined: 5 comment / 10 code =
	// 0.5) will show as its own baseline. Expected multiple = 4.0/0.5 = 8.
	c2 := c1 + "\n// new c1\n// new c2\n// new c3\n// new c4\nfunc I() {}\n"
	if err := os.WriteFile(filepath.Join(dir, "a.go"), []byte(c2), 0644); err != nil {
		t.Fatal(err)
	}
	gitRun(t, dir, "add", "a.go")
	gitRun(t, dir, "commit", "-q", "-m", "c2")
	c2sha := gitRun(t, dir, "rev-parse", "HEAD")

	cmd := exec.Command(lgBin, filepath.Join(repoRoot, "scripts", "lint.lg"),
		"--churn", c1sha+".."+c2sha, ".")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("lint.lg --churn: %v\n%s", err, out)
	}

	found := false
	for _, line := range strings.Split(string(out), "\n") {
		if !strings.Contains(line, "[comment additions relative to code]") {
			continue
		}
		found = true
		if !strings.Contains(line, "range-ratio=4") {
			t.Errorf("expected range-ratio=4, got: %s", line)
		}
		if !strings.Contains(line, "baseline=0.5") {
			t.Errorf("expected baseline=0.5, got: %s", line)
		}
		if !strings.Contains(line, "multiple=8") {
			t.Errorf("expected multiple=8, got: %s", line)
		}
	}
	if !found {
		t.Fatalf("no [comment additions relative to code] finding in output:\n%s", out)
	}
	assertLintReaderText(t, string(out))

	// A range adding comments but no code cannot yield a comparable ratio.
	if err := os.WriteFile(filepath.Join(dir, "a.go"), []byte(c2+"// comment only\n"), 0644); err != nil {
		t.Fatal(err)
	}
	gitRun(t, dir, "add", "a.go")
	gitRun(t, dir, "commit", "-q", "-m", "c3")
	c3sha := gitRun(t, dir, "rev-parse", "HEAD")
	emptyCmd := exec.Command(lgBin, filepath.Join(repoRoot, "scripts", "lint.lg"), "--churn", c2sha+".."+c3sha, ".")
	emptyCmd.Dir = dir
	emptyOut, err := emptyCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("comment-only range: %v\n%s", err, emptyOut)
	}
	assertLintReaderText(t, string(emptyOut))
	if !strings.Contains(string(emptyOut), "comment additions relative to code: no comparable measurement") ||
		strings.Contains(string(emptyOut), "[comment additions relative to code]") {
		t.Errorf("comment-only range should explain why it has no finding:\n%s", emptyOut)
	}

	// Without --churn, R6 must not run at all (no default range makes sense).
	noChurnCmd := exec.Command(lgBin, filepath.Join(repoRoot, "scripts", "lint.lg"), ".")
	noChurnCmd.Dir = dir
	noChurnOut, err := noChurnCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("lint.lg (no --churn): %v\n%s", err, noChurnOut)
	}
	if strings.Contains(string(noChurnOut), "[comment additions relative to code]") {
		t.Errorf("R6 must not fire without --churn:\n%s", noChurnOut)
	}
}

func TestLintR6ZeroCommentRangeIsDiagnostic(t *testing.T) {
	dir := t.TempDir()
	gitRun(t, dir, "init", "-q")
	gitRun(t, dir, "config", "commit.gpgsign", "false")
	path := filepath.Join(dir, "a.go")
	base := "package p\n// baseline\nfunc A() {}\n"
	if err := os.WriteFile(path, []byte(base), 0644); err != nil {
		t.Fatal(err)
	}
	gitRun(t, dir, "add", "a.go")
	gitRun(t, dir, "commit", "-q", "-m", "c1")
	c1 := gitRun(t, dir, "rev-parse", "HEAD")
	if err := os.WriteFile(path, []byte(base+"func B() {}\n"), 0644); err != nil {
		t.Fatal(err)
	}
	gitRun(t, dir, "add", "a.go")
	gitRun(t, dir, "commit", "-q", "-m", "c2")
	c2 := gitRun(t, dir, "rev-parse", "HEAD")
	rangeArg := c1 + ".." + c2

	cmd := exec.Command(lgBin, filepath.Join(repoRoot, "scripts", "lint.lg"), "--churn", rangeArg, ".")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("zero-comment range: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "[comment additions relative to code]") || !strings.Contains(string(out), "multiple=0") {
		t.Errorf("zero-comment range must report a 0x multiple:\n%s", out)
	}

	gate := exec.Command(lgBin, filepath.Join(repoRoot, "scripts", "lint.lg"),
		"--gate", "comment-churn", "--churn", rangeArg, ".")
	gate.Dir = dir
	gateOut, err := gate.CombinedOutput()
	if err != nil {
		t.Fatalf("R6 is a measurement, not a gate: %v\n%s", err, gateOut)
	}
	if !strings.Contains(string(gateOut), "comment additions relative to code is report-only") {
		t.Errorf("ignored R6 gate must explain why:\n%s", gateOut)
	}
	assertLintReaderText(t, string(gateOut))

	noChurn := exec.Command(lgBin, filepath.Join(repoRoot, "scripts", "lint.lg"),
		"--gate", "comment-churn", ".")
	noChurn.Dir = dir
	noChurnOut, err := noChurn.CombinedOutput()
	if err != nil {
		t.Fatalf("R6 without --churn must not gate: %v\n%s", err, noChurnOut)
	}
	if !strings.Contains(string(noChurnOut), "comment additions relative to code is report-only") ||
		!strings.Contains(string(noChurnOut), "comment additions relative to code: skipped") {
		t.Errorf("missing R6 ignored/skipped notes:\n%s", noChurnOut)
	}
}

func TestLintR6ZeroBaselineIsUndefined(t *testing.T) {
	dir := t.TempDir()
	gitRun(t, dir, "init", "-q")
	gitRun(t, dir, "config", "commit.gpgsign", "false")
	path := filepath.Join(dir, "a.go")
	base := "package p\nfunc A() {}\n"
	if err := os.WriteFile(path, []byte(base), 0644); err != nil {
		t.Fatal(err)
	}
	gitRun(t, dir, "add", "a.go")
	gitRun(t, dir, "commit", "-q", "-m", "c1")
	c1 := gitRun(t, dir, "rev-parse", "HEAD")
	if err := os.WriteFile(path, []byte(base+"func B() {}\n"), 0644); err != nil {
		t.Fatal(err)
	}
	gitRun(t, dir, "add", "a.go")
	gitRun(t, dir, "commit", "-q", "-m", "c2")
	c2 := gitRun(t, dir, "rev-parse", "HEAD")
	rangeArg := c1 + ".." + c2

	for _, tc := range []struct {
		name string
		args []string
		want string
	}{
		{"text", []string{"--churn", rangeArg, "."}, "multiple=undefined"},
		{"edn", []string{"--edn", "--churn", rangeArg, "."}, ":measure nil"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cmd := exec.Command(lgBin, append([]string{filepath.Join(repoRoot, "scripts", "lint.lg")}, tc.args...)...)
			cmd.Dir = dir
			out, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("zero baseline: %v\n%s", err, out)
			}
			if !strings.Contains(string(out), tc.want) {
				t.Errorf("missing %q for zero baseline:\n%s", tc.want, out)
			}
			if tc.name == "text" && !strings.Contains(string(out), "relative multiple undefined because the corpus has no comment lines") {
				t.Errorf("undefined multiple lacks its cause:\n%s", out)
			}
		})
	}
}

// Code-verbosity: a generic pattern/replace matcher over parsed .lg forms,
// catalog-driven from scripts/lint-code-rules.edn (see that file, and the
// "Code-verbosity" section comment in scripts/lint.lg, for the matcher
// semantics: a "?name" metavariable, a "?&name" rest-metavariable, and
// atoms-eliminable = atom-count(matched) - atom-count(replacement)).
//
// One fixture exercises the whole shipped catalog at once, each rule paired
// with a near-miss that must NOT fire (the near-miss is what proves the
// matcher is precise, not just present).
// redundant-if-true-false ("(if ?c true false)" -> "?c") and
// let-empty-bindings ("(let [] ?e)" -> "?e") are NOT included below: both
// were built, tested, found to fire zero times on the real corpus, and
// dropped from the shipped catalog per the calibration duty (see the
// report and scripts/lint-code-rules.edn's own comments) — so this fixture,
// which is checked against the shipped catalog, only exercises what ships.
const lintFixtureCodeVerbosity = `(defn b1 [c] (if c false true))
(defn b2 [c] (if c false false))

(defn c1 [x] (= true x))
(defn c2 [x] (= false x))

(defn d1 [a b] (not (= a b)))
(defn d2 [a b] (not (and a b)))

(defn e1 [x] (not (empty? x)))
(defn e2 [x] (not (nil? x)))

(def f1 (fn [a b] (helper a b)))
(def f2 (fn [a b] (helper b a)))

(defn g1 [] (let [x 5] x))
(defn g2 [] (let [x 5] (inc x)))

(defn i1 [] (do (single-thing)))
(defn i2 [] (do (one) (two)))

(defn j1 [x] (first (first x)))
(defn j2 [x] (first (rest x)))

(defn k1 [sep coll] (apply str (interpose sep coll)))
(defn k2 [sep coll] (apply str coll))

(defn l1 [c1 c2 c3] (if c1 :a (if c2 :b (if c3 :c :d))))
(defn l2 [c1 c2] (if c1 :a (if c2 :b :c)))
`

// wantCodeVerbosity maps each fixture line to the :kind that must (line
// suffix "1") or must not (line suffix "2") fire there.
var codeVerbosityPositive = map[int]string{
	1:  "redundant-if-false-true",
	4:  "redundant-eq-true",
	7:  "redundant-not-eq",
	10: "redundant-not-empty",
	13: "eta-expansion",
	16: "let-identity",
	19: "do-single-form",
	22: "composable-accessor",
	25: "apply-str-interpose",
	28: "nested-if-chain",
}

func TestLintCodeVerbosityCatalog(t *testing.T) {
	dir := t.TempDir()
	fixture := filepath.Join(dir, "fixture.lg")
	if err := os.WriteFile(fixture, []byte(lintFixtureCodeVerbosity), 0644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(lgBin, filepath.Join(repoRoot, "scripts", "lint.lg"), "--edn", dir)
	cmd.Dir = repoRoot
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("lint.lg --edn: %v\n%s", err, out)
	}

	// One finding per positive line, none on the near-miss (line+1) lines.
	gotLines := map[int]string{}
	for line, kind := range findKindsPerLine(t, string(out)) {
		gotLines[line] = kind
	}
	for line, wantKind := range codeVerbosityPositive {
		if got, ok := gotLines[line]; !ok {
			t.Errorf("line %d: expected [%s] to fire, nothing did\n%s", line, wantKind, out)
		} else if got != wantKind {
			t.Errorf("line %d: got kind %q, want %q", line, got, wantKind)
		}
		if nearMiss, ok := gotLines[line+1]; ok {
			t.Errorf("near-miss at line %d must not fire, but got kind %q", line+1, nearMiss)
		}
	}
}

func TestLintCodeVerbosityReportText(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "fixture.lg"), []byte("(defn f [x] (first (first x)))\n"), 0644); err != nil {
		t.Fatal(err)
	}
	out, err := exec.Command(lgBin, filepath.Join(repoRoot, "scripts", "lint.lg"), dir).CombinedOutput()
	if err != nil {
		t.Fatalf("code pattern report: %v\n%s", err, out)
	}
	assertLintReaderText(t, string(out))
	if !strings.Contains(string(out), "[nested sequence accessor]") ||
		!strings.Contains(string(out), "ffirst says the same thing directly") ||
		strings.Contains(string(out), "composable-accessor:") ||
		!strings.Contains(string(out), "code pattern(s) that can be simplified") {
		t.Errorf("code pattern report lacks readable label or evidence:\n%s", out)
	}
}

// findKindsPerLine extracts {startLine: kind} from --edn output by scanning
// for `:line N, ... :kind :K` pairs — a small hand parser since we don't
// want a real EDN reader in the test just to check shape.
func findKindsPerLine(t *testing.T, ednOut string) map[int]string {
	t.Helper()
	out := map[int]string{}
	entries := strings.Split(ednOut, "{:file")
	for _, e := range entries[1:] {
		lineIdx := strings.Index(e, ":line ")
		kindIdx := strings.Index(e, ":kind :")
		if lineIdx < 0 || kindIdx < 0 {
			continue
		}
		lineStr := strings.Fields(e[lineIdx+len(":line "):])[0]
		lineStr = strings.TrimSuffix(lineStr, ",")
		n, err := strconv.Atoi(lineStr)
		if err != nil {
			continue
		}
		kindStr := strings.Fields(e[kindIdx+len(":kind :"):])[0]
		kindStr = strings.TrimSuffix(kindStr, ",")
		if _, exists := out[n]; !exists {
			out[n] = kindStr
		}
	}
	return out
}

// A rule added purely as catalog DATA — no change to lint.lg — must be
// found and applied. This is the property the data-driven design exists to
// provide. The catalog is resolved next to lint.lg's OWN invoked path, so
// this test invokes a copy of lint.lg placed alongside a fabricated
// catalog, proving the substitution needs zero changes to lint.lg's code —
// only a different scripts/lint-code-rules.edn.
func TestLintCodeVerbosityCatalogIsData(t *testing.T) {
	dir := t.TempDir()
	scriptsDir := filepath.Join(dir, "scripts")
	if err := os.MkdirAll(scriptsDir, 0755); err != nil {
		t.Fatal(err)
	}
	realLint, err := os.ReadFile(filepath.Join(repoRoot, "scripts", "lint.lg"))
	if err != nil {
		t.Fatal(err)
	}
	probe := `
(println (str "DISPLAY\t" (display-name :fabricated-test-rule (load-code-rules))))
`
	if err := os.WriteFile(filepath.Join(scriptsDir, "lint.lg"), append(realLint, []byte(probe)...), 0644); err != nil {
		t.Fatal(err)
	}
	customCatalog := `{:name :fabricated-test-rule
 :kind :fabricated-test-rule
 :pattern (triple ?x)
 :replace (* 3 ?x)
 :why "test-only pattern proving the catalog is pure data"}
`
	if err := os.WriteFile(filepath.Join(scriptsDir, "lint-code-rules.edn"), []byte(customCatalog), 0644); err != nil {
		t.Fatal(err)
	}
	fixtureDir := filepath.Join(dir, "src")
	if err := os.MkdirAll(fixtureDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(fixtureDir, "f.lg"), []byte("(defn f [n] (triple n))\n"), 0644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(lgBin, filepath.Join(scriptsDir, "lint.lg"), fixtureDir)
	cmd.Dir = repoRoot
	out, cmdErr := cmd.CombinedOutput()
	if cmdErr != nil {
		t.Fatalf("lint.lg: %v\n%s", cmdErr, out)
	}
	if !strings.Contains(string(out), "[fabricated test rule]") {
		t.Errorf("a catalog-only rule addition (no lint.lg code change) was not found:\n%s", out)
	}
	if !strings.Contains(string(out), "DISPLAY\tfabricated test rule") {
		t.Errorf("unlabeled catalog kind did not get a readable fallback:\n%s", out)
	}
}

func equalInts(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
