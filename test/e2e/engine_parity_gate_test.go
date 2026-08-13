/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package e2e

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/genmanifest"
)

// Classification tests precede the harness implementation so the category and
// ratchet contracts remain the primary executable specification.
func TestClassifyCorpusUsesObservedResults(t *testing.T) {
	fixtures := []string{"full.lg", "partial.lg", "diverged.lg", "broken.lg"}
	observations := []fixtureObservation{
		observation("full.lg", succeeded("control"), succeeded("control"), succeeded("control"), succeeded("control")),
		observation("partial.lg", succeeded("control"), succeeded("control"), strictRejected("ir-compile-strict: unsupported"), succeeded("strict")),
		observation("diverged.lg", succeeded("control"), succeeded("control"), succeeded("bytecode"), succeeded("gogen")),
		observation("broken.lg", succeeded("control"), succeeded("control"), failed(2, "parse error at line 3"), succeeded("strict")),
	}

	census := classifyCorpus(fixtures, observations)
	if !census.Full["full.lg"] {
		t.Fatalf("full = %v; want full.lg", census.Full)
	}
	if !census.Partial["partial.lg"] {
		t.Fatalf("partial = %v; want strict rejection classified as partial", census.Partial)
	}
	if !census.Diverged["diverged.lg"] {
		t.Fatalf("diverged = %v; want successful strict output divergence", census.Diverged)
	}
	if len(census.Fatal) != 1 || !strings.Contains(census.Fatal[0], "parse error at line 3") {
		t.Fatalf("fatal = %v; want arbitrary strict failure with stderr provenance", census.Fatal)
	}
}

func TestStrictOutputsAreComparedToNonStrictControl(t *testing.T) {
	obs := observation("same-wrong.lg", succeeded("control"), succeeded("control"), succeeded("wrong"), succeeded("wrong"))
	census := classifyCorpus([]string{"same-wrong.lg"}, []fixtureObservation{obs})
	if !census.Diverged["same-wrong.lg"] {
		t.Fatalf("diverged = %v; want strict output change recorded even when engines agree with each other", census.Diverged)
	}
}

func TestStrictRejectionUsesHarnessOwnedExitCodeAndRetainsBothStreams(t *testing.T) {
	result := executionResult{
		Stdout:   "fixture output",
		Stderr:   "execution error: ir-compile-strict: unsupported form",
		ExitCode: strictRejectionExitCode,
	}
	if !result.strictRejection() {
		t.Fatal("dedicated strict rejection exit code was not recognized")
	}
	for _, want := range []string{result.Stdout, result.Stderr} {
		if !strings.Contains(result.evidence(), want) {
			t.Fatalf("evidence %q does not retain %q", result.evidence(), want)
		}
	}
}

func TestAmbientStrictWordsWithUnrelatedFailureAreFatal(t *testing.T) {
	forged := executionResult{
		Stdout:   "fixture says ir-compile-strict: but then fails",
		Stderr:   "boom",
		ExitCode: 1,
	}
	obs := observation("forged.lg", succeeded("control"), succeeded("control"), forged, succeeded("control"))
	census := classifyCorpus([]string{"forged.lg"}, []fixtureObservation{obs})
	if len(census.Fatal) != 1 || !strings.Contains(census.Fatal[0], "boom") {
		t.Fatalf("fatal = %v; want unrelated failure with preserved evidence", census.Fatal)
	}
	if census.Partial["forged.lg"] {
		t.Fatal("ambient fixture output forged a strict capability rejection")
	}
}

func TestPartialAndDivergedCategoriesMayOverlap(t *testing.T) {
	obs := observation(
		"both.lg",
		succeeded("control"),
		succeeded("control"),
		strictRejected("ir-compile-strict: unsupported"),
		succeeded("different"),
	)
	census := classifyCorpus([]string{"both.lg"}, []fixtureObservation{obs})
	if !census.Partial["both.lg"] || !census.Diverged["both.lg"] {
		t.Fatalf("partial=%v diverged=%v; want both measured facts recorded", census.Partial, census.Diverged)
	}
}

func TestNonStrictFailureIsFatalWithProvenance(t *testing.T) {
	obs := observation("bad.lg", failed(1, "reader exploded"), succeeded("unused"), strictRejected("ir-compile-strict"), succeeded("unused"))
	census := classifyCorpus([]string{"bad.lg"}, []fixtureObservation{obs})
	if len(census.Fatal) != 1 || !strings.Contains(census.Fatal[0], "reader exploded") {
		t.Fatalf("fatal = %v; want non-strict stderr", census.Fatal)
	}
	if census.Partial["bad.lg"] {
		t.Fatal("broken non-strict control was mislabeled as partial strict capability")
	}
}

func TestStrictSuccessWithoutEligibleSeamIsOnlyCapabilityEvidence(t *testing.T) {
	obs := observation("no-seam.lg", succeeded("same"), succeeded("same"), succeeded("same"), succeeded("same"))
	census := classifyCorpus([]string{"no-seam.lg"}, []fixtureObservation{obs})
	if !census.Full["no-seam.lg"] {
		t.Fatalf("full = %v; want observed full strict parity", census.Full)
	}
	if strings.Contains(strings.ToLower(census.Summary()), "genuine") || strings.Contains(strings.ToLower(census.Summary()), "native") {
		t.Fatalf("summary %q overclaims generated/native execution", census.Summary())
	}
}

func TestCorpusCompletenessRejectsMissingDuplicateAndUnexpectedResults(t *testing.T) {
	observations := []fixtureObservation{
		observation("a.lg", succeeded("x"), succeeded("x"), succeeded("x"), succeeded("x")),
		observation("a.lg", succeeded("x"), succeeded("x"), succeeded("x"), succeeded("x")),
		observation("extra.lg", succeeded("x"), succeeded("x"), succeeded("x"), succeeded("x")),
	}
	census := classifyCorpus([]string{"a.lg", "b.lg"}, observations)
	joined := strings.Join(census.Fatal, "\n")
	for _, want := range []string{"duplicate result for a.lg", "unexpected result for extra.lg", "missing result for b.lg"} {
		if !strings.Contains(joined, want) {
			t.Errorf("fatal errors %q do not contain %q", joined, want)
		}
	}
}

func TestRatchetsRejectNewAndStaleEntries(t *testing.T) {
	got := validateRatchet("partial", setOf("new.lg", "same.lg"), setOf("same.lg", "stale.lg"))
	joined := strings.Join(got, "\n")
	for _, want := range []string{"NEW partial: new.lg", "STALE partial: stale.lg"} {
		if !strings.Contains(joined, want) {
			t.Errorf("errors %q do not contain %q", joined, want)
		}
	}

	got = validateRatchet("strict divergence", setOf("new.lg"), setOf("stale.lg"))
	joined = strings.Join(got, "\n")
	for _, want := range []string{"NEW strict divergence: new.lg", "STALE strict divergence: stale.lg"} {
		if !strings.Contains(joined, want) {
			t.Errorf("errors %q do not contain %q", joined, want)
		}
	}
}

func TestBaselineBytesAreDeterministicAndSorted(t *testing.T) {
	first := baselineBytes(partialBaselineKind, setOf("z.lg", "a.lg"))
	second := baselineBytes(partialBaselineKind, setOf("a.lg", "z.lg"))
	if !bytes.Equal(first, second) {
		t.Fatalf("baseline bytes depend on map order:\nfirst:\n%s\nsecond:\n%s", first, second)
	}
	if strings.Index(string(first), "a.lg") > strings.Index(string(first), "z.lg") {
		t.Fatalf("baseline is not sorted:\n%s", first)
	}
}

func TestInvalidCensusCannotRewriteBaselines(t *testing.T) {
	root := t.TempDir()
	partialPath, divergencePath := baselinePaths(root)
	if err := os.MkdirAll(filepath.Dir(partialPath), 0o755); err != nil {
		t.Fatal(err)
	}
	before := map[string]string{
		ledgerPath(root): "ledger-before\n",
		partialPath:      "partial-before\n",
		divergencePath:   "divergence-before\n",
	}
	for path, data := range before {
		if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	invalid := corpusCensus{Partial: setOf("new.lg"), Diverged: setOf("new.lg"), Fatal: []string{"harness failed"}}
	ledger := parityLedger{Cases: []ledgerCase{{Fixture: "new.lg", Status: statusPartial, Reason: ledgerNoReason}}}
	if err := writeCensusBaselines(root, invalid, ledger); err == nil {
		t.Fatal("invalid census unexpectedly rewrote baselines")
	}
	for path, want := range before {
		got, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if string(got) != want {
			t.Fatalf("%s changed after invalid census: got %q, want %q", path, got, want)
		}
	}
}

func TestBaselineSetWriteLeavesAllThreeFilesOnInvalidTarget(t *testing.T) {
	root := t.TempDir()
	partialPath, divergencePath := baselinePaths(root)
	if err := os.MkdirAll(filepath.Dir(partialPath), 0o755); err != nil {
		t.Fatal(err)
	}
	unchanged := map[string]string{
		ledgerPath(root): "ledger-before\n",
		partialPath:      "partial-before\n",
	}
	for path, data := range unchanged {
		if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	// The last target of the set is unwritable, so the whole set must roll back.
	if err := os.Mkdir(divergencePath, 0o755); err != nil {
		t.Fatal(err)
	}
	census := corpusCensus{Complete: true, Partial: setOf("new.lg"), Diverged: setOf("also-new.lg")}
	ledger := parityLedger{Cases: []ledgerCase{{Fixture: "new.lg", Status: statusPartial, Reason: ledgerNoReason}}}
	if err := writeCensusBaselines(root, census, ledger); err == nil {
		t.Fatal("baseline set write unexpectedly accepted a directory target")
	}
	for path, want := range unchanged {
		got, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if string(got) != want {
			t.Fatalf("%s changed after a failed set write: got %q, want %q", path, got, want)
		}
	}
	info, err := os.Stat(divergencePath)
	if err != nil {
		t.Fatal(err)
	}
	if !info.IsDir() {
		t.Fatalf("divergence target changed after a failed set write: mode=%v", info.Mode())
	}
	entries, err := os.ReadDir(filepath.Dir(partialPath))
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".parity-") {
			t.Fatalf("failed set write left staging or backup debris: %s", entry.Name())
		}
	}
}

func TestValidRederiveWritesLedgerAndBothBaselines(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "test"), 0o755); err != nil {
		t.Fatal(err)
	}
	census := corpusCensus{Complete: true, Partial: setOf("z.lg", "a.lg"), Diverged: setOf("d.lg")}
	ledger := parityLedger{Cases: []ledgerCase{
		{Fixture: "a.lg", Status: statusPartial, Fallbacks: 1, Reason: "unsupported shape"},
		{Fixture: "d.lg", Status: statusDiverged, Fallbacks: 0, Reason: ledgerNoReason},
	}}
	if err := writeCensusBaselines(root, census, ledger); err != nil {
		t.Fatal(err)
	}
	partialPath, divergencePath := baselinePaths(root)
	for path, want := range map[string][]byte{
		ledgerPath(root): ledgerBytes(ledger),
		partialPath:      baselineBytes(partialBaselineKind, census.Partial),
		divergencePath:   baselineBytes(divergenceBaselineKind, census.Diverged),
	} {
		got, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(got, want) {
			t.Fatalf("%s bytes = %q; want %q", path, got, want)
		}
	}
}

func TestLedgerBytesAreDeterministicSortedAndRecordTombstones(t *testing.T) {
	ledger := parityLedger{
		Cases: []ledgerCase{
			{Fixture: "z.lg", Status: statusPartial + "+" + statusDiverged, Fallbacks: 2, Reason: "unsupported shape"},
			{Fixture: "a.lg", Status: statusFull, Fallbacks: 0, Reason: ledgerNoReason},
		},
		Tombstones: []ledgerTombstone{
			{Fixture: "gone.lg", LastStatus: statusPartial, RemovedAt: "deadbeef", Reason: "fixture folded into seq.lg"},
		},
	}
	shuffled := parityLedger{
		Cases:      []ledgerCase{ledger.Cases[1], ledger.Cases[0]},
		Tombstones: ledger.Tombstones,
	}
	first := ledgerBytes(ledger)
	if !bytes.Equal(first, ledgerBytes(shuffled)) {
		t.Fatalf("ledger bytes depend on input order:\n%s\n---\n%s", first, ledgerBytes(shuffled))
	}
	text := string(first)
	for _, want := range []string{
		`CASE a.lg status=full fallbacks=0 reason="-"`,
		`CASE z.lg status=partial+diverged fallbacks=2 reason="unsupported shape"`,
		`TOMBSTONE gone.lg last_status=partial removed_at=deadbeef reason="fixture folded into seq.lg"`,
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("ledger does not contain %q:\n%s", want, text)
		}
	}
	if strings.Index(text, "CASE a.lg") > strings.Index(text, "CASE z.lg") {
		t.Fatalf("case lines are not sorted:\n%s", text)
	}
	if strings.Index(text, "CASE z.lg") > strings.Index(text, "TOMBSTONE gone.lg") {
		t.Fatalf("tombstones are not written after case lines:\n%s", text)
	}

	path := filepath.Join(t.TempDir(), "parity-ledger.txt")
	if err := os.WriteFile(path, first, 0o644); err != nil {
		t.Fatal(err)
	}
	parsed, err := parseLedger(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(ledgerBytes(parsed), first) {
		t.Fatalf("round-tripped ledger = %q; want %q", ledgerBytes(parsed), first)
	}
}

func TestLedgerStatusUsesCanonicalOverlapOrder(t *testing.T) {
	census := corpusCensus{Complete: true, Partial: setOf("both.lg"), Diverged: setOf("both.lg")}
	if got := ledgerStatus(census, "both.lg"); got != statusPartial+"+"+statusDiverged {
		t.Fatalf("ledgerStatus = %q; want canonical partial+diverged", got)
	}
	if got := ledgerStatus(corpusCensus{}, "missing.lg"); got != statusFatal {
		t.Fatalf("ledgerStatus = %q; want fatal for an unclassified fixture", got)
	}
}

func TestCensusLedgerCasesRecordMeasuredFallbackSitesAndReasons(t *testing.T) {
	census := corpusCensus{Complete: true, Full: setOf("full.lg"), Partial: setOf("partial.lg")}
	audits := map[string]fixtureAudit{"partial.lg": {Engines: []engineAudit{
		{Engine: "gogen_ir", Entries: []fallbackEntry{{Defn: "b", Reason: "zeta reason"}}},
		{Engine: "bytecode", Entries: []fallbackEntry{{Defn: "a", Reason: "alpha reason"}, {Defn: "b", Reason: "zeta reason"}}},
	}}}
	cases := censusLedgerCases([]string{"partial.lg", "full.lg"}, census, audits)
	want := []ledgerCase{
		{Fixture: "full.lg", Status: statusFull, Fallbacks: 0, Reason: ledgerNoReason},
		{Fixture: "partial.lg", Status: statusPartial, Fallbacks: 3, Reason: "alpha reason" + ledgerReasonSeparator + "zeta reason"},
	}
	if fmt.Sprint(cases) != fmt.Sprint(want) {
		t.Fatalf("cases = %v; want %v", cases, want)
	}
}

func TestLedgerDiffReportsAddedRemovedAndStatusChangedCases(t *testing.T) {
	committed := parityLedger{Cases: []ledgerCase{
		{Fixture: "kept.lg", Status: statusFull},
		{Fixture: "changed.lg", Status: statusFull},
		{Fixture: "removed.lg", Status: statusPartial},
	}}
	current := []ledgerCase{
		{Fixture: "kept.lg", Status: statusFull},
		{Fixture: "changed.lg", Status: statusPartial},
		{Fixture: "added.lg", Status: statusFull},
	}
	joined := strings.Join(diffLedger(current, committed), "\n")
	for _, want := range []string{
		"ADDED added.lg status=full",
		"REMOVED removed.lg status=partial",
		"STATUS-CHANGED changed.lg full -> partial",
	} {
		if !strings.Contains(joined, want) {
			t.Errorf("ledger diff %q does not contain %q", joined, want)
		}
	}
	if strings.Contains(joined, "kept.lg") {
		t.Errorf("ledger diff reported an unchanged case: %q", joined)
	}
}

func TestReappearingTombstonedFixtureIsAdded(t *testing.T) {
	committed := parityLedger{Tombstones: []ledgerTombstone{
		{Fixture: "back.lg", LastStatus: statusPartial, RemovedAt: "deadbeef", Reason: "renamed"},
	}}
	joined := strings.Join(diffLedger([]ledgerCase{{Fixture: "back.lg", Status: statusFull}}, committed), "\n")
	if !strings.Contains(joined, "ADDED back.lg status=full") {
		t.Fatalf("ledger diff %q does not report a reappearing tombstoned fixture as ADDED", joined)
	}
	if !strings.Contains(joined, "tombstone") {
		t.Fatalf("ledger diff %q does not tell the author to remove the tombstone deliberately", joined)
	}
}

func TestUnexplainedTombstoneFailsAndExplainedTombstonePasses(t *testing.T) {
	unexplained := parityLedger{Tombstones: []ledgerTombstone{
		{Fixture: "gone.lg", LastStatus: statusPartial, RemovedAt: "deadbeef", Reason: ledgerUnexplained},
	}}
	joined := strings.Join(diffLedger(nil, unexplained), "\n")
	if !strings.Contains(joined, "UNEXPLAINED-TOMBSTONE gone.lg") {
		t.Fatalf("ledger diff %q does not fail on an unexplained tombstone", joined)
	}

	explained := parityLedger{Tombstones: []ledgerTombstone{
		{Fixture: "gone.lg", LastStatus: statusPartial, RemovedAt: "deadbeef", Reason: "fixture merged into seq.lg"},
	}}
	if errs := diffLedger(nil, explained); len(errs) != 0 {
		t.Fatalf("explained tombstone failed the gate: %v", errs)
	}
}

func TestRederiveTombstonesMissingCasesWithEnvReason(t *testing.T) {
	previous := parityLedger{Cases: []ledgerCase{
		{Fixture: "kept.lg", Status: statusFull},
		{Fixture: "explained.lg", Status: statusPartial},
		{Fixture: "silent.lg", Status: statusDiverged},
	}}
	current := []ledgerCase{{Fixture: "kept.lg", Status: statusFull}}
	reasons := parseTombstoneReasons("explained.lg=folded into kept.lg;ignored")
	next, _ := rederiveLedger(current, previous, "cafebabe", reasons, false)
	want := []ledgerTombstone{
		{Fixture: "explained.lg", LastStatus: statusPartial, RemovedAt: "cafebabe", Reason: "folded into kept.lg"},
		{Fixture: "silent.lg", LastStatus: statusDiverged, RemovedAt: "cafebabe", Reason: ledgerUnexplained},
	}
	if fmt.Sprint(next.Tombstones) != fmt.Sprint(want) {
		t.Fatalf("tombstones = %v; want %v", next.Tombstones, want)
	}
	if fmt.Sprint(next.Cases) != fmt.Sprint(current) {
		t.Fatalf("cases = %v; want the measured census %v", next.Cases, current)
	}
}

func TestRederiveDropsTombstoneWhenFixtureIsMeasuredAgain(t *testing.T) {
	previous := parityLedger{Tombstones: []ledgerTombstone{
		{Fixture: "back.lg", LastStatus: statusPartial, RemovedAt: "deadbeef", Reason: "renamed"},
	}}
	next, _ := rederiveLedger([]ledgerCase{{Fixture: "back.lg", Status: statusFull}}, previous, "cafebabe", nil, false)
	if len(next.Tombstones) != 0 {
		t.Fatalf("tombstones = %v; want the tombstone dropped for a measured fixture", next.Tombstones)
	}
}

func TestLedgerCleanupDropsOnlyPriorRevisionTombstones(t *testing.T) {
	previous := parityLedger{Tombstones: []ledgerTombstone{
		{Fixture: "old.lg", LastStatus: statusPartial, RemovedAt: "deadbeef", Reason: "removed last release"},
		{Fixture: "new.lg", LastStatus: statusFull, RemovedAt: "cafebabe", Reason: "removed in this revision"},
	}}
	kept, note := rederiveLedger(nil, previous, "cafebabe", nil, false)
	if len(kept.Tombstones) != 2 {
		t.Fatalf("tombstones = %v; want both retained without cleanup", kept.Tombstones)
	}
	if note != "" {
		t.Fatalf("note = %q; want no note when cleanup was not requested", note)
	}
	cleaned, note := rederiveLedger(nil, previous, "cafebabe", nil, true)
	if len(cleaned.Tombstones) != 1 || cleaned.Tombstones[0].Fixture != "new.lg" {
		t.Fatalf("tombstones = %v; want only the current-revision tombstone retained", cleaned.Tombstones)
	}
	if note != "" {
		t.Fatalf("note = %q; want no note when cleanup ran against a resolvable revision", note)
	}
}

func TestLedgerCleanupWithUnknownRevisionDropsNothingAndReportsWhy(t *testing.T) {
	previous := parityLedger{Tombstones: []ledgerTombstone{
		{Fixture: "old.lg", LastStatus: statusPartial, RemovedAt: "deadbeef", Reason: "removed last release"},
		{Fixture: "murky.lg", LastStatus: statusFull, RemovedAt: ledgerUnknownCommit, Reason: "removed without provenance"},
	}}
	next, note := rederiveLedger(nil, previous, ledgerUnknownCommit, nil, true)
	if len(next.Tombstones) != 2 {
		t.Fatalf("tombstones = %v; want every tombstone retained when the revision is unresolvable", next.Tombstones)
	}
	if note == "" || !strings.Contains(note, ledgerCleanupEnv) || !strings.Contains(note, ledgerUnknownCommit) {
		t.Fatalf("note = %q; want an explanation naming %s and %s", note, ledgerCleanupEnv, ledgerUnknownCommit)
	}
}

func TestLedgerCleanupDropsUnknownRemovalWhenCurrentRevisionResolves(t *testing.T) {
	previous := parityLedger{Tombstones: []ledgerTombstone{
		{Fixture: "murky.lg", LastStatus: statusFull, RemovedAt: ledgerUnknownCommit, Reason: "removed without provenance"},
		{Fixture: "new.lg", LastStatus: statusPartial, RemovedAt: "cafebabe", Reason: "removed in this revision"},
	}}
	cleaned, note := rederiveLedger(nil, previous, "cafebabe", nil, true)
	if len(cleaned.Tombstones) != 1 || cleaned.Tombstones[0].Fixture != "new.lg" {
		t.Fatalf("tombstones = %v; want a recorded removed_at=%s tombstone dropped like any other foreign revision", cleaned.Tombstones, ledgerUnknownCommit)
	}
	if note != "" {
		t.Fatalf("note = %q; want no note when the current revision resolves", note)
	}
}

func TestRederiveNeverEmitsTwoTombstonesForOneFixture(t *testing.T) {
	// A half-done hand edit can leave a fixture recorded as both a case and a
	// tombstone. Rederive must still produce a ledger its own parser accepts.
	previous := parityLedger{
		Cases:      []ledgerCase{{Fixture: "x.lg", Status: statusPartial, Reason: ledgerNoReason}},
		Tombstones: []ledgerTombstone{{Fixture: "x.lg", LastStatus: statusPartial, RemovedAt: "deadbeef", Reason: "removed earlier"}},
	}
	next, _ := rederiveLedger(nil, previous, "cafebabe", nil, false)
	if len(next.Tombstones) != 1 || next.Tombstones[0].Fixture != "x.lg" {
		t.Fatalf("tombstones = %v; want exactly one tombstone for x.lg", next.Tombstones)
	}

	path := filepath.Join(t.TempDir(), "parity-ledger.txt")
	data := ledgerBytes(next)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
	parsed, err := parseLedger(path)
	if err != nil {
		t.Fatalf("rederived ledger is rejected by its own parser: %v\n%s", err, data)
	}
	if !bytes.Equal(ledgerBytes(parsed), data) {
		t.Fatalf("round-tripped ledger = %q; want %q", ledgerBytes(parsed), data)
	}
}

func TestParseLedgerRejectsMalformedLedgers(t *testing.T) {
	for _, tc := range []struct {
		name string
		body string
		want string
	}{
		{"unrecognized line", `NOTE a.lg status=full fallbacks=0 reason="-"`, "unrecognized ledger line"},
		{"malformed case", `CASE a.lg reason="-"`, "malformed CASE line"},
		{"malformed tombstone", `TOMBSTONE a.lg reason="-"`, "malformed TOMBSTONE line"},
		{"wrong case field key", `CASE a.lg state=full fallbacks=0 reason="-"`, `expected status= field`},
		{"wrong tombstone field key", `TOMBSTONE a.lg status=full removed_at=x reason="-"`, `expected last_status= field`},
		{"non-numeric fallbacks", `CASE a.lg status=full fallbacks=many reason="-"`, "fallbacks:"},
		{"missing reason", `CASE a.lg status=full fallbacks=0`, "missing reason= field"},
		{"unquoted reason", `CASE a.lg status=full fallbacks=0 reason=-`, "unquote reason"},
		{
			"duplicate case",
			"CASE a.lg status=full fallbacks=0 reason=\"-\"\nCASE a.lg status=partial fallbacks=1 reason=\"-\"",
			`duplicate ledger entry "a.lg"`,
		},
		{
			"same fixture as case and tombstone",
			"CASE a.lg status=full fallbacks=0 reason=\"-\"\nTOMBSTONE a.lg last_status=full removed_at=deadbeef reason=\"gone\"",
			`duplicate ledger entry "a.lg"`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "parity-ledger.txt")
			if err := os.WriteFile(path, []byte("# header\n\n"+tc.body+"\n"), 0o644); err != nil {
				t.Fatal(err)
			}
			ledger, err := parseLedger(path)
			if err == nil {
				t.Fatalf("parseLedger accepted %q as %#v", tc.body, ledger)
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("error %q does not contain %q", err, tc.want)
			}
			if !strings.Contains(err.Error(), path) {
				t.Fatalf("error %q does not name the offending file %q", err, path)
			}
		})
	}
}

func TestFatalReasonsMatchOnlyClassifierMessageShapes(t *testing.T) {
	census := corpusCensus{Fatal: []string{
		"var.lg: non-strict control failed; bytecode={} gogen_ir={}",
		"missing result for var.lg",
		"missing result for extra-var.lg",
		"duplicate result for var.lg.bak",
		"unexpected result for other.lg",
		"novar.lg: bytecode strict run failed outside the strict-capability rejection path; x",
	}}
	got := fatalReasons(census, "var.lg")
	want := []string{
		"missing result for var.lg",
		"var.lg: non-strict control failed; bytecode={} gogen_ir={}",
	}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("fatalReasons = %v; want only the two classifier shapes for var.lg: %v", got, want)
	}
}

func TestRepoRevisionPrefersJJChangeIDThenGitCommitThenUnknown(t *testing.T) {
	type call struct {
		dir  string
		name string
		args []string
	}
	var calls []call
	runner := func(out map[string]string, fail map[string]bool) commandRunner {
		return func(dir string, name string, args ...string) (string, error) {
			calls = append(calls, call{dir, name, args})
			if fail[name] {
				return "", fmt.Errorf("%s unavailable", name)
			}
			return out[name], nil
		}
	}

	calls = nil
	jjOnly := runner(map[string]string{"jj": "zzzzqqqqwwww\n", "git": "deadbeef\n"}, nil)
	if got := resolveRepoRevision("/root", jjOnly); got != "zzzzqqqqwwww" {
		t.Fatalf("revision = %q; want the trimmed jj change id", got)
	}
	if len(calls) != 1 || calls[0].name != "jj" || calls[0].dir != "/root" {
		t.Fatalf("calls = %v; want exactly one jj call in the repository root", calls)
	}

	calls = nil
	gitFallback := runner(map[string]string{"git": "  deadbeef  \n"}, map[string]bool{"jj": true})
	if got := resolveRepoRevision("/root", gitFallback); got != "deadbeef" {
		t.Fatalf("revision = %q; want the git commit when jj fails", got)
	}
	if len(calls) != 2 || calls[1].name != "git" {
		t.Fatalf("calls = %v; want jj attempted before git", calls)
	}

	blankThenGit := runner(map[string]string{"jj": "   \n", "git": "deadbeef"}, nil)
	if got := resolveRepoRevision("/root", blankThenGit); got != "deadbeef" {
		t.Fatalf("revision = %q; want blank jj output to fall through to git", got)
	}

	multiline := runner(map[string]string{"jj": "one\ntwo\n", "git": "three\nfour"}, nil)
	if got := resolveRepoRevision("/root", multiline); got != ledgerUnknownCommit {
		t.Fatalf("revision = %q; want %s for multi-line output", got, ledgerUnknownCommit)
	}

	allFail := runner(nil, map[string]bool{"jj": true, "git": true})
	if got := resolveRepoRevision("/root", allFail); got != ledgerUnknownCommit {
		t.Fatalf("revision = %q; want %s when no version-control command succeeds", got, ledgerUnknownCommit)
	}

	panicking := func(string, string, ...string) (string, error) { panic("no version control here") }
	if got := resolveRepoRevision("/root", panicking); got != ledgerUnknownCommit {
		t.Fatalf("revision = %q; want %s instead of a panic escaping the gate", got, ledgerUnknownCommit)
	}
}

func TestFixtureWrappersShareSetupAndSequentialEvaluation(t *testing.T) {
	dir := t.TempDir()
	source := filepath.Join(dir, "fixture.lg")
	fixture := "(require '[example.lib :as lib])\n(lib/run)\n"
	if err := os.WriteFile(source, []byte(fixture), 0o644); err != nil {
		t.Fatal(err)
	}

	wrappers := map[fixtureMode]string{}
	for _, mode := range []fixtureMode{controlMode, strictMode} {
		dest := filepath.Join(dir, string(mode)+".lg")
		if err := createFixtureWrapper(dest, source, mode); err != nil {
			t.Fatal(err)
		}
		wrappedBytes, err := os.ReadFile(dest)
		if err != nil {
			t.Fatal(err)
		}
		wrappers[mode] = string(wrappedBytes)

		ordered := []string{
			"(require 'ir.passes.pipeline)",
			mode.compileBindings(),
			"(doseq [form (read-all-string ",
			"(eval form)",
		}
		previous := -1
		for _, shape := range ordered {
			index := strings.Index(wrappers[mode], shape)
			if index <= previous {
				t.Fatalf("%s wrapper does not contain %q in prerequisite/sequential order:\n%s", mode, shape, wrappers[mode])
			}
			previous = index
		}
		if strings.Contains(wrappers[mode], `(str "(do ")`) || strings.Contains(wrappers[mode], "(read-string") {
			t.Fatalf("%s wrapper still compiles all top-level forms before requires execute:\n%s", mode, wrappers[mode])
		}
	}

	strictWrapper := wrappers[strictMode]
	for _, shape := range []string{
		"(let [strict-fallback-log (atom [])]",
		"(reset! strict-fallback-log [])",
		"*ir-compile-verbose* true",
		"*ir-compile-fallback-log* strict-fallback-log",
		"(empty? (deref strict-fallback-log))",
		fmt.Sprintf("(os/exit %d)", strictRejectionExitCode),
	} {
		if !strings.Contains(strictWrapper, shape) {
			t.Fatalf("strict wrapper does not contain %q:\n%s", shape, strictWrapper)
		}
	}
	if strings.Contains(strictWrapper, "index-of") || strings.Contains(strictWrapper, `"ir-compile-strict:"`) {
		t.Fatalf("strict wrapper attributes rejection through exception text:\n%s", strictWrapper)
	}
	for _, shape := range []string{"strict-fallback-log", "*ir-compile-verbose*"} {
		if strings.Contains(wrappers[controlMode], shape) {
			t.Fatalf("control wrapper unexpectedly contains strict observation seam %q:\n%s", shape, wrappers[controlMode])
		}
	}
}

func TestParityReportBytesAreDeterministicAndSorted(t *testing.T) {
	census := corpusCensus{
		Complete: true,
		Full:     setOf("a-full.lg"),
		Partial:  setOf("z-partial.lg"),
		Diverged: setOf("m-diverged.lg"),
	}
	audits := map[string]fixtureAudit{
		"z-partial.lg": {Engines: []engineAudit{
			{Engine: "gogen_ir", Entries: []fallbackEntry{{Defn: "zeta", Reason: "unsupported interop"}, {Defn: "alpha", Reason: "unsupported interop"}}},
			{Engine: "bytecode", Entries: []fallbackEntry{{Defn: "alpha", Reason: "unsupported interop"}}},
		}},
		"m-diverged.lg": {Engines: []engineAudit{
			{Engine: "bytecode", Entries: []fallbackEntry{{Defn: "mid", Reason: "unsupported loop form"}}},
		}},
	}
	first := buildParityReport([]string{"z-partial.lg", "a-full.lg", "m-diverged.lg"}, census, audits)
	second := buildParityReport([]string{"m-diverged.lg", "z-partial.lg", "a-full.lg"}, census, audits)
	if !bytes.Equal(first, second) {
		t.Fatalf("report bytes depend on input or map order:\nfirst:\n%s\nsecond:\n%s", first, second)
	}
	text := string(first)
	if strings.Index(text, "FIXTURE a-full.lg") > strings.Index(text, "FIXTURE m-diverged.lg") ||
		strings.Index(text, "FIXTURE m-diverged.lg") > strings.Index(text, "FIXTURE z-partial.lg") {
		t.Fatalf("fixture lines are not sorted:\n%s", text)
	}
	bucketSection := text[strings.Index(text, "BUCKET-CENSUS"):]
	if strings.Index(bucketSection, `BUCKET count=3 reason="unsupported interop"`) >
		strings.Index(bucketSection, `BUCKET count=1 reason="unsupported loop form"`) {
		t.Fatalf("bucket census is not ordered by descending count then reason:\n%s", text)
	}
	if !strings.Contains(text, "BUCKET count=3 reason=\"unsupported interop\"") {
		t.Fatalf("bucket census does not report per-site counts:\n%s", text)
	}
	for _, want := range []string{
		"AGGREGATE full strict parity=1 partial strict capability=1 strict output divergence=1 fatal=0 total=3",
		"AGGREGATE audited_fixtures=2",
		"BUCKET-CENSUS distinct_reasons=2",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("report does not contain aggregate line %q:\n%s", want, text)
		}
	}
}

func TestFullParityFixtureReportsZeroFallbacksWithoutAuditRun(t *testing.T) {
	census := corpusCensus{Complete: true, Full: setOf("full.lg")}
	if got := fixturesNeedingAudit([]string{"full.lg"}, census); len(got) != 0 {
		t.Fatalf("fixturesNeedingAudit = %v; want no audit for full strict parity", got)
	}
	report := string(buildParityReport([]string{"full.lg"}, census, map[string]fixtureAudit{}))
	want := `FIXTURE full.lg category=full strict parity fallbacks=0 reason="strict success proves zero IR fallback"`
	if !strings.Contains(report, want) {
		t.Fatalf("report %q does not contain %q", report, want)
	}
	if strings.Contains(report, auditFailedPrefix) {
		t.Fatalf("full strict parity fixture reported a missing audit:\n%s", report)
	}
}

func TestParityReportSurfacesAuditedDefnAndRawReason(t *testing.T) {
	stderr := "ignored preamble\n" + auditBeginMarker + "\n" +
		auditEntryMarker + "\t\"my-defn\"\t\"ir-compile: unsupported form\\n\\tat pass 3\"\n" +
		auditEndMarker + "\n"
	entries, err := parseAuditLog(stderr)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || entries[0].Defn != "my-defn" || entries[0].Reason != "ir-compile: unsupported form\n\tat pass 3" {
		t.Fatalf("parsed entries = %#v; want raw defn and reason text", entries)
	}

	census := corpusCensus{Complete: true, Partial: setOf("partial.lg")}
	audits := map[string]fixtureAudit{"partial.lg": {Engines: []engineAudit{{Engine: "bytecode", Entries: entries}}}}
	report := string(buildParityReport([]string{"partial.lg"}, census, audits))
	for _, want := range []string{
		"FIXTURE partial.lg category=partial strict capability distinct_reasons=1",
		`DETAIL partial.lg engine=bytecode defn=my-defn reason="ir-compile: unsupported form\n\tat pass 3"`,
		`BUCKET count=1 reason="ir-compile: unsupported form\n\tat pass 3"`,
		"BUCKET-SITE partial.lg engine=bytecode defn=my-defn",
	} {
		if !strings.Contains(report, want) {
			t.Fatalf("report does not contain %q:\n%s", want, report)
		}
	}
}

func TestParityReportDegradesWhenAuditFails(t *testing.T) {
	crashed := auditFromResult("gogen_ir", failed(1, "boom"))
	if crashed.Err == "" || len(crashed.Entries) != 0 {
		t.Fatalf("failed audit run = %#v; want a reported error and no entries", crashed)
	}
	if !strings.Contains(crashed.Err, "boom") {
		t.Fatalf("audit error %q loses run evidence", crashed.Err)
	}
	unmarked := auditFromResult("bytecode", succeeded(""))
	if unmarked.Err == "" {
		t.Fatal("audit output without markers was accepted as a measurement")
	}

	census := corpusCensus{Complete: true, Partial: setOf("partial.lg")}
	audits := map[string]fixtureAudit{"partial.lg": {Engines: []engineAudit{crashed, unmarked}}}
	report := string(buildParityReport([]string{"partial.lg"}, census, audits))
	if !strings.Contains(report, "FIXTURE partial.lg category=partial strict capability") {
		t.Fatalf("audit failure changed the reported classification:\n%s", report)
	}
	for _, want := range []string{
		"DETAIL partial.lg engine=bytecode " + auditFailedPrefix,
		"DETAIL partial.lg engine=gogen_ir " + auditFailedPrefix,
	} {
		if !strings.Contains(report, want) {
			t.Fatalf("report does not contain %q:\n%s", want, report)
		}
	}
	if strings.Contains(report, "BUCKET count=") {
		t.Fatalf("degraded audit invented lowering buckets:\n%s", report)
	}

	missing := string(buildParityReport([]string{"partial.lg"}, census, map[string]fixtureAudit{}))
	if !strings.Contains(missing, auditFailedPrefix+": no audit run recorded") {
		t.Fatalf("absent audit was silently reported as clean:\n%s", missing)
	}
}

func TestParityReportIsWrittenOnlyWhenEnvPathIsSet(t *testing.T) {
	report := buildParityReport([]string{"full.lg"}, corpusCensus{Complete: true, Full: setOf("full.lg")}, nil)
	if err := writeParityReport("", report); err != nil {
		t.Fatalf("unset %s must be a no-op: %v", parityReportEnv, err)
	}
	path := filepath.Join(t.TempDir(), "parity-report.txt")
	if err := writeParityReport(path, report); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, report) {
		t.Fatalf("written report = %q; want %q", got, report)
	}
}

func TestAuditWrapperMeasuresFallbacksWithoutStrictRejection(t *testing.T) {
	dir := t.TempDir()
	source := filepath.Join(dir, "fixture.lg")
	if err := os.WriteFile(source, []byte("(defn f [] 1)\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	dest := filepath.Join(dir, "audit.lg")
	if err := createFixtureWrapper(dest, source, auditMode); err != nil {
		t.Fatal(err)
	}
	wrappedBytes, err := os.ReadFile(dest)
	if err != nil {
		t.Fatal(err)
	}
	wrapper := string(wrappedBytes)
	for _, shape := range []string{
		"(let [audit-fallback-log (atom [])]",
		"*ir-compile* true",
		"*ir-compile-strict* false",
		"*ir-compile-verbose* true",
		"*ir-compile-fallback-log* audit-fallback-log",
		auditBeginMarker,
		auditEntryMarker,
		auditEndMarker,
		auditErrorMarker,
	} {
		if !strings.Contains(wrapper, shape) {
			t.Fatalf("audit wrapper does not contain %q:\n%s", shape, wrapper)
		}
	}
	if strings.Contains(wrapper, fmt.Sprintf("(os/exit %d)", strictRejectionExitCode)) {
		t.Fatalf("audit wrapper can emit the strict rejection exit code:\n%s", wrapper)
	}
	if depth := parenDepth(wrapper); depth != 0 {
		t.Fatalf("audit wrapper parens are unbalanced (depth=%d):\n%s", depth, wrapper)
	}
}

// parenDepth returns the final paren nesting depth of an .lg wrapper, skipping
// string literals, or -1 if it ever closes past zero.
func parenDepth(source string) int {
	depth := 0
	inString := false
	escaped := false
	for _, r := range source {
		if inString {
			switch {
			case escaped:
				escaped = false
			case r == '\\':
				escaped = true
			case r == '"':
				inString = false
			}
			continue
		}
		switch r {
		case '"':
			inString = true
		case '(':
			depth++
		case ')':
			depth--
			if depth < 0 {
				return -1
			}
		}
	}
	return depth
}

type executionResult struct {
	Stdout     string
	Stderr     string
	ExitCode   int
	StartError string
}

type fixtureObservation struct {
	Fixture      string
	ControlBC    executionResult
	ControlGogen executionResult
	StrictBC     executionResult
	StrictGogen  executionResult
}

type corpusCensus struct {
	Complete bool
	Full     map[string]bool
	Partial  map[string]bool
	Diverged map[string]bool
	Fatal    []string
}

const (
	parityRederiveEnv       = "LETGO_PARITY_REDERIVE"
	partialBaselineKind     = "partial"
	divergenceBaselineKind  = "strict-divergence"
	strictRejectionExitCode = 86
)

func succeeded(stdout string) executionResult {
	return executionResult{Stdout: stdout}
}

func failed(exitCode int, stderr string) executionResult {
	return executionResult{Stderr: stderr, ExitCode: exitCode}
}

func strictRejected(stderr string) executionResult {
	return executionResult{Stderr: stderr, ExitCode: strictRejectionExitCode}
}

func observation(name string, controlBC, controlGogen, strictBC, strictGogen executionResult) fixtureObservation {
	return fixtureObservation{name, controlBC, controlGogen, strictBC, strictGogen}
}

func (r executionResult) OK() bool {
	return r.StartError == "" && r.ExitCode == 0
}

func (r executionResult) strictRejection() bool {
	return r.StartError == "" && r.ExitCode == strictRejectionExitCode
}

func (r executionResult) evidence() string {
	return fmt.Sprintf("exit=%d start_error=%q stdout=%q stderr=%q", r.ExitCode, r.StartError, r.Stdout, r.Stderr)
}

func classifyCorpus(discovered []string, observations []fixtureObservation) corpusCensus {
	census := corpusCensus{Full: map[string]bool{}, Partial: map[string]bool{}, Diverged: map[string]bool{}}
	expected := setOf(discovered...)
	seen := map[string]bool{}
	for _, obs := range observations {
		if seen[obs.Fixture] {
			census.Fatal = append(census.Fatal, fmt.Sprintf("duplicate result for %s", obs.Fixture))
			continue
		}
		seen[obs.Fixture] = true
		if !expected[obs.Fixture] {
			census.Fatal = append(census.Fatal, fmt.Sprintf("unexpected result for %s", obs.Fixture))
			continue
		}
		classifyObservation(&census, obs)
	}
	for _, name := range discovered {
		if !seen[name] {
			census.Fatal = append(census.Fatal, fmt.Sprintf("missing result for %s", name))
		}
	}
	census.Complete = len(discovered) > 0 && len(seen) == len(expected)
	sort.Strings(census.Fatal)
	return census
}

func classifyObservation(census *corpusCensus, obs fixtureObservation) {
	if !obs.ControlBC.OK() || !obs.ControlGogen.OK() {
		census.Fatal = append(census.Fatal, fmt.Sprintf(
			"%s: non-strict control failed; bytecode={%s} gogen_ir={%s}",
			obs.Fixture, obs.ControlBC.evidence(), obs.ControlGogen.evidence()))
		return
	}
	if obs.ControlBC.Stdout != obs.ControlGogen.Stdout {
		census.Fatal = append(census.Fatal, fmt.Sprintf(
			"%s: non-strict controls diverged; bytecode={%s} gogen_ir={%s}",
			obs.Fixture, obs.ControlBC.evidence(), obs.ControlGogen.evidence()))
		return
	}

	strictRuns := []struct {
		name   string
		result executionResult
	}{{"bytecode", obs.StrictBC}, {"gogen_ir", obs.StrictGogen}}
	partial := false
	diverged := false
	controlOutput := obs.ControlBC.Stdout
	for _, run := range strictRuns {
		if run.result.OK() {
			if run.result.Stdout != controlOutput {
				diverged = true
			}
			continue
		}
		if run.result.strictRejection() {
			partial = true
			continue
		}
		census.Fatal = append(census.Fatal, fmt.Sprintf(
			"%s: %s strict run failed outside the strict-capability rejection path; %s",
			obs.Fixture, run.name, run.result.evidence()))
		return
	}
	if partial {
		census.Partial[obs.Fixture] = true
	}
	if diverged {
		census.Diverged[obs.Fixture] = true
	}
	if !partial && !diverged {
		census.Full[obs.Fixture] = true
	}
}

func (c corpusCensus) Valid() bool {
	return c.Complete && len(c.Fatal) == 0
}

func (c corpusCensus) Summary() string {
	classified := setOf()
	for name := range c.Full {
		classified[name] = true
	}
	for name := range c.Partial {
		classified[name] = true
	}
	for name := range c.Diverged {
		classified[name] = true
	}
	return fmt.Sprintf("full strict parity=%d partial strict capability=%d strict output divergence=%d fatal=%d total=%d",
		len(c.Full), len(c.Partial), len(c.Diverged), len(c.Fatal), len(classified))
}

func setOf(names ...string) map[string]bool {
	set := make(map[string]bool, len(names))
	for _, name := range names {
		set[name] = true
	}
	return set
}

func validateRatchet(kind string, current, baseline map[string]bool) []string {
	var errs []string
	for name := range current {
		if !baseline[name] {
			errs = append(errs, fmt.Sprintf("NEW %s: %s", kind, name))
		}
	}
	for name := range baseline {
		if !current[name] {
			errs = append(errs, fmt.Sprintf("STALE %s: %s; remove it (shrink-only ratchet)", kind, name))
		}
	}
	sort.Strings(errs)
	return errs
}

func baselinePaths(root string) (partial, divergence string) {
	return filepath.Join(root, "test/parity-partial.txt"), filepath.Join(root, "test/parity-divergence.txt")
}

func readBaseline(path string) (map[string]bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	entries := map[string]bool{}
	for lineNumber, line := range strings.Split(string(data), "\n") {
		entry := strings.TrimSpace(line)
		if entry == "" || strings.HasPrefix(entry, "#") {
			continue
		}
		if entries[entry] {
			return nil, fmt.Errorf("%s:%d: duplicate baseline entry %q", path, lineNumber+1, entry)
		}
		entries[entry] = true
	}
	return entries, nil
}

func baselineBytes(kind string, entries map[string]bool) []byte {
	names := make([]string, 0, len(entries))
	for name := range entries {
		names = append(names, name)
	}
	sort.Strings(names)
	var b strings.Builder
	switch kind {
	case partialBaselineKind:
		b.WriteString("# Machine-generated shrink-only baseline of fixtures with partial strict capability.\n")
		b.WriteString("# The target state is empty. Rederive only from a complete valid census.\n")
	case divergenceBaselineKind:
		b.WriteString("# Machine-generated shrink-only baseline of unattributed strict engine-output divergences.\n")
		b.WriteString("# The target state is empty. Rederive only from a complete valid census.\n")
	default:
		panic("unknown baseline kind: " + kind)
	}
	b.WriteString("# Partial and divergence are independent facts; a fixture may appear in both.\n")
	b.WriteString("# LETGO_PARITY_REDERIVE=1 make engine-parity-gate\n")
	// The blank separator belongs between the header and the entries, so it is
	// written only when there are entries. An empty baseline — the target state
	// for both kinds — otherwise ends on a blank line, which `git diff --check`
	// reports as whitespace damage on every rederive.
	if len(names) > 0 {
		b.WriteString("\n")
	}
	for _, name := range names {
		b.WriteString(name)
		b.WriteByte('\n')
	}
	return []byte(b.String())
}

// writeCensusBaselines rewrites the ledger and both shrink-only baselines as
// one validated, rollback-safe set. Either all three files reflect the same
// census or none of them changes.
func writeCensusBaselines(root string, census corpusCensus, ledger parityLedger) error {
	if !census.Valid() {
		return fmt.Errorf("refusing to rederive baselines from invalid census: %s: %s", census.Summary(), strings.Join(census.Fatal, "; "))
	}
	partialPath, divergencePath := baselinePaths(root)
	files := []baselineReplacement{
		{target: ledgerPath(root), data: ledgerBytes(ledger)},
		{target: partialPath, data: baselineBytes(partialBaselineKind, census.Partial)},
		{target: divergencePath, data: baselineBytes(divergenceBaselineKind, census.Diverged)},
	}
	return replaceBaselineSet(files)
}

type baselineReplacement struct {
	target string
	data   []byte
	staged string
	backup string
}

func replaceBaselineSet(files []baselineReplacement) error {
	for i := range files {
		staged, err := stageBaseline(files[i].target, files[i].data)
		if err != nil {
			removeStagedBaselines(files)
			return err
		}
		files[i].staged = staged
	}
	defer removeStagedBaselines(files)

	for _, file := range files {
		info, err := os.Lstat(file.target)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return fmt.Errorf("validate %s: %w", file.target, err)
		}
		if !info.Mode().IsRegular() {
			return fmt.Errorf("validate %s: target is not a regular file", file.target)
		}
	}

	for i := range files {
		if _, err := os.Lstat(files[i].target); os.IsNotExist(err) {
			continue
		} else if err != nil {
			rollbackErr := rollbackBaselineSet(files, 0)
			return joinBaselineErrors(fmt.Errorf("prepare backup for %s: %w", files[i].target, err), rollbackErr)
		}
		backup, err := reserveBackupPath(files[i].target)
		if err == nil {
			err = os.Rename(files[i].target, backup)
		}
		if err != nil {
			os.Remove(backup)
			rollbackErr := rollbackBaselineSet(files, 0)
			return joinBaselineErrors(fmt.Errorf("prepare backup for %s: %w", files[i].target, err), rollbackErr)
		}
		files[i].backup = backup
	}

	committed := 0
	for i := range files {
		if err := os.Rename(files[i].staged, files[i].target); err != nil {
			rollbackErr := rollbackBaselineSet(files, committed)
			return joinBaselineErrors(fmt.Errorf("replace %s: %w", files[i].target, err), rollbackErr)
		}
		files[i].staged = ""
		committed++
	}
	for i := range files {
		if files[i].backup != "" {
			// The set is committed. Backup cleanup is best-effort so a cleanup
			// failure cannot turn a successful replacement into a reported error.
			os.Remove(files[i].backup)
			files[i].backup = ""
		}
	}
	return nil
}

func stageBaseline(path string, data []byte) (string, error) {
	tmp, err := os.CreateTemp(filepath.Dir(path), ".parity-baseline-*")
	if err != nil {
		return "", fmt.Errorf("stage %s: %w", path, err)
	}
	tmpName := tmp.Name()
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		os.Remove(tmpName)
		return "", fmt.Errorf("write staged %s: %w", path, err)
	}
	if err := tmp.Chmod(0o644); err != nil {
		tmp.Close()
		os.Remove(tmpName)
		return "", fmt.Errorf("chmod staged %s: %w", path, err)
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmpName)
		return "", fmt.Errorf("close staged %s: %w", path, err)
	}
	return tmpName, nil
}

func reserveBackupPath(target string) (string, error) {
	tmp, err := os.CreateTemp(filepath.Dir(target), ".parity-backup-*")
	if err != nil {
		return "", err
	}
	name := tmp.Name()
	if err := tmp.Close(); err != nil {
		os.Remove(name)
		return "", err
	}
	if err := os.Remove(name); err != nil {
		return "", err
	}
	return name, nil
}

func rollbackBaselineSet(files []baselineReplacement, committed int) error {
	var errs []string
	for i := 0; i < committed; i++ {
		if err := os.Remove(files[i].target); err != nil && !os.IsNotExist(err) {
			errs = append(errs, fmt.Sprintf("remove replacement %s: %v", files[i].target, err))
		}
	}
	for i := range files {
		if files[i].backup == "" {
			continue
		}
		if err := os.Rename(files[i].backup, files[i].target); err != nil {
			errs = append(errs, fmt.Sprintf("restore %s: %v", files[i].target, err))
		}
		files[i].backup = ""
	}
	if len(errs) != 0 {
		return fmt.Errorf("rollback baseline set: %s", strings.Join(errs, "; "))
	}
	return nil
}

func removeStagedBaselines(files []baselineReplacement) {
	for _, file := range files {
		if file.staged != "" {
			os.Remove(file.staged)
		}
	}
}

func joinBaselineErrors(primary, rollback error) error {
	if rollback == nil {
		return primary
	}
	return fmt.Errorf("%v; %w", primary, rollback)
}

// ---------------------------------------------------------------------------
// Case ledger
//
// test/parity-ledger.txt is a committed, diffable record of every measured
// case plus tombstones for cases that were removed. Unlike the shrink-only
// partial/divergence baselines, the ledger records full state: a reviewer can
// read a commit diff and see exactly which cases were added, removed, or
// changed status. It never changes classification.
// ---------------------------------------------------------------------------

const (
	ledgerTombstoneReasonEnv = "LETGO_PARITY_TOMBSTONE_REASON"
	ledgerCleanupEnv         = "LETGO_PARITY_LEDGER_CLEANUP"
	ledgerUnexplained        = "UNEXPLAINED"
	ledgerUnknownCommit      = "unknown"
	ledgerNoReason           = "-"
	ledgerReasonSeparator    = " | "
	statusFull               = "full"
	statusPartial            = "partial"
	statusDiverged           = "diverged"
	statusFatal              = "fatal"
)

// ledgerCase is one measured fixture as recorded in the committed ledger.
type ledgerCase struct {
	Fixture   string
	Status    string
	Fallbacks int
	Reason    string
}

// ledgerTombstone records a fixture that used to be in the ledger and is no
// longer discovered by the census.
type ledgerTombstone struct {
	Fixture    string
	LastStatus string
	RemovedAt  string
	Reason     string
}

// parityLedger is the parsed contents of test/parity-ledger.txt.
type parityLedger struct {
	Cases      []ledgerCase
	Tombstones []ledgerTombstone
}

func ledgerPath(root string) string {
	return filepath.Join(root, "test/parity-ledger.txt")
}

// ledgerStatus renders the canonical status token for a fixture. Partial and
// divergence are independent measured facts, so overlaps join in a fixed
// order: always partial+diverged, never diverged+partial.
func ledgerStatus(census corpusCensus, fixture string) string {
	var parts []string
	if census.Full[fixture] {
		parts = append(parts, statusFull)
	}
	if census.Partial[fixture] {
		parts = append(parts, statusPartial)
	}
	if census.Diverged[fixture] {
		parts = append(parts, statusDiverged)
	}
	if len(parts) == 0 {
		return statusFatal
	}
	return strings.Join(parts, "+")
}

// ledgerFallbacksAndReason summarizes the audit measurement for one fixture:
// the number of measured fallback sites and the distinct reasons behind them.
// A full strict parity fixture is never audited; strict success already proves
// zero fallback.
func ledgerFallbacksAndReason(census corpusCensus, fixture string, audits map[string]fixtureAudit) (int, string) {
	if !needsAudit(census, fixture) {
		return 0, ledgerNoReason
	}
	audit, measured := audits[fixture]
	if !measured {
		return 0, auditFailedPrefix + ": no audit run recorded"
	}
	sites := 0
	reasons := map[string]bool{}
	for _, engine := range audit.Engines {
		if engine.Err != "" {
			reasons[engine.Err] = true
			continue
		}
		for _, entry := range engine.Entries {
			sites++
			reasons[entry.Reason] = true
		}
	}
	if len(reasons) == 0 {
		return sites, ledgerNoReason
	}
	return sites, strings.Join(sortedKeys(reasons), ledgerReasonSeparator)
}

// censusLedgerCases renders the measured census as ledger cases.
func censusLedgerCases(names []string, census corpusCensus, audits map[string]fixtureAudit) []ledgerCase {
	ordered := append([]string(nil), names...)
	sort.Strings(ordered)
	cases := make([]ledgerCase, 0, len(ordered))
	for _, name := range ordered {
		fallbacks, reason := ledgerFallbacksAndReason(census, name, audits)
		cases = append(cases, ledgerCase{Fixture: name, Status: ledgerStatus(census, name), Fallbacks: fallbacks, Reason: reason})
	}
	return cases
}

// ledgerBytes renders a deterministic, sorted ledger file.
func ledgerBytes(ledger parityLedger) []byte {
	cases := append([]ledgerCase(nil), ledger.Cases...)
	sort.Slice(cases, func(i, j int) bool { return cases[i].Fixture < cases[j].Fixture })
	tombstones := append([]ledgerTombstone(nil), ledger.Tombstones...)
	sort.Slice(tombstones, func(i, j int) bool { return tombstones[i].Fixture < tombstones[j].Fixture })

	var b strings.Builder
	b.WriteString("# Machine-generated case ledger for the engine parity gate.\n")
	b.WriteString("# generated-by: TestEngineParityGate (test/e2e/engine_parity_gate_test.go)\n")
	b.WriteString("# rederive: LETGO_PARITY_REDERIVE=1 make engine-parity-gate\n")
	b.WriteString("# This file is a record AND a contract, not a shrink-only ratchet: it lists\n")
	b.WriteString("# every measured case so a commit diff shows exactly what changed. The gate\n")
	b.WriteString("# fails on ADDED, REMOVED, or STATUS-CHANGED cases and on any tombstone whose\n")
	b.WriteString("# reason is " + ledgerUnexplained + ".\n")
	b.WriteString("# status: full | partial | diverged | fatal; overlapping facts join with '+'\n")
	b.WriteString("# in the canonical order partial+diverged.\n")
	b.WriteString("# fallbacks: measured audit fallback site count (0 for full strict parity).\n")
	b.WriteString("# reason: distinct measured reasons joined with '" + ledgerReasonSeparator + "'; \"" + ledgerNoReason + "\" when none.\n")
	b.WriteString("# Tombstone reasons come from " + ledgerTombstoneReasonEnv + "=\"fixture=why;fixture2=why2\".\n")
	b.WriteString("# removed_at: the revision that removed the fixture (jj change id, else git commit,\n")
	b.WriteString("# else \"" + ledgerUnknownCommit + "\").\n")
	b.WriteString("# " + ledgerCleanupEnv + "=1 during rederive drops tombstones inherited from earlier revisions;\n")
	b.WriteString("# it does nothing when the current revision is \"" + ledgerUnknownCommit + "\".\n\n")
	for _, c := range cases {
		b.WriteString(fmt.Sprintf("CASE %s status=%s fallbacks=%d reason=%s\n", c.Fixture, c.Status, c.Fallbacks, strconv.Quote(c.Reason)))
	}
	for _, t := range tombstones {
		b.WriteString(fmt.Sprintf("TOMBSTONE %s last_status=%s removed_at=%s reason=%s\n",
			t.Fixture, t.LastStatus, t.RemovedAt, strconv.Quote(t.Reason)))
	}
	return []byte(b.String())
}

// ledgerField extracts a `key=value` token from a ledger line prefix section.
func ledgerField(fields []string, index int, key string) (string, error) {
	if index >= len(fields) {
		return "", fmt.Errorf("missing %s= field", key)
	}
	if !strings.HasPrefix(fields[index], key+"=") {
		return "", fmt.Errorf("expected %s= field, got %q", key, fields[index])
	}
	return strings.TrimPrefix(fields[index], key+"="), nil
}

// ledgerQuotedReason unquotes the trailing `reason="..."` field of a line.
func ledgerQuotedReason(line string) (string, error) {
	index := strings.Index(line, "reason=")
	if index < 0 {
		return "", fmt.Errorf("missing reason= field")
	}
	reason, err := strconv.Unquote(strings.TrimSpace(line[index+len("reason="):]))
	if err != nil {
		return "", fmt.Errorf("unquote reason: %w", err)
	}
	return reason, nil
}

// parseLedger reads a committed ledger. A malformed ledger is an error, never
// a silently empty ledger. Fixture identity is global across kinds: one
// fixture may appear once, as a CASE or as a TOMBSTONE, never both.
func parseLedger(path string) (parityLedger, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return parityLedger{}, fmt.Errorf("read %s: %w", path, err)
	}
	var ledger parityLedger
	seen := map[string]bool{}
	for lineNumber, raw := range strings.Split(string(data), "\n") {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		head := line
		if index := strings.Index(line, "reason="); index >= 0 {
			head = line[:index]
		}
		fields := strings.Fields(head)
		if len(fields) == 0 || (fields[0] != "CASE" && fields[0] != "TOMBSTONE") {
			return parityLedger{}, fmt.Errorf("%s:%d: unrecognized ledger line %q", path, lineNumber+1, line)
		}
		reason, err := ledgerQuotedReason(line)
		if err != nil {
			return parityLedger{}, fmt.Errorf("%s:%d: %w", path, lineNumber+1, err)
		}
		switch fields[0] {
		case "CASE":
			if len(fields) < 4 {
				return parityLedger{}, fmt.Errorf("%s:%d: malformed CASE line %q", path, lineNumber+1, line)
			}
			status, err := ledgerField(fields, 2, "status")
			if err != nil {
				return parityLedger{}, fmt.Errorf("%s:%d: %w", path, lineNumber+1, err)
			}
			rawFallbacks, err := ledgerField(fields, 3, "fallbacks")
			if err != nil {
				return parityLedger{}, fmt.Errorf("%s:%d: %w", path, lineNumber+1, err)
			}
			fallbacks, err := strconv.Atoi(rawFallbacks)
			if err != nil {
				return parityLedger{}, fmt.Errorf("%s:%d: fallbacks: %w", path, lineNumber+1, err)
			}
			if seen[fields[1]] {
				return parityLedger{}, fmt.Errorf("%s:%d: duplicate ledger entry %q", path, lineNumber+1, fields[1])
			}
			seen[fields[1]] = true
			ledger.Cases = append(ledger.Cases, ledgerCase{Fixture: fields[1], Status: status, Fallbacks: fallbacks, Reason: reason})
		case "TOMBSTONE":
			if len(fields) < 4 {
				return parityLedger{}, fmt.Errorf("%s:%d: malformed TOMBSTONE line %q", path, lineNumber+1, line)
			}
			lastStatus, err := ledgerField(fields, 2, "last_status")
			if err != nil {
				return parityLedger{}, fmt.Errorf("%s:%d: %w", path, lineNumber+1, err)
			}
			removedAt, err := ledgerField(fields, 3, "removed_at")
			if err != nil {
				return parityLedger{}, fmt.Errorf("%s:%d: %w", path, lineNumber+1, err)
			}
			if seen[fields[1]] {
				return parityLedger{}, fmt.Errorf("%s:%d: duplicate ledger entry %q", path, lineNumber+1, fields[1])
			}
			seen[fields[1]] = true
			ledger.Tombstones = append(ledger.Tombstones, ledgerTombstone{
				Fixture: fields[1], LastStatus: lastStatus, RemovedAt: removedAt, Reason: reason})
		}
	}
	return ledger, nil
}

// parseTombstoneReasons parses LETGO_PARITY_TOMBSTONE_REASON, which uses the
// form "fixture=why;fixture2=why2".
func parseTombstoneReasons(spec string) map[string]string {
	reasons := map[string]string{}
	for _, clause := range strings.Split(spec, ";") {
		clause = strings.TrimSpace(clause)
		if clause == "" {
			continue
		}
		parts := strings.SplitN(clause, "=", 2)
		if len(parts) != 2 {
			continue
		}
		fixture := strings.TrimSpace(parts[0])
		why := strings.TrimSpace(parts[1])
		if fixture == "" || why == "" {
			continue
		}
		reasons[fixture] = why
	}
	return reasons
}

// commandRunner runs a read-only version-control command in dir and returns
// its raw stdout. It exists as a seam so revision resolution is unit-testable
// without shelling out.
type commandRunner func(dir string, name string, args ...string) (string, error)

// execCommandRunner is the production commandRunner.
func execCommandRunner(dir string, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return stdout.String(), nil
}

// revisionCommands are tried in order. jj comes first because this repository
// is driven by jj and a jj workspace need not contain a .git directory; the jj
// change id is also stable across amends, which is the correct identity for
// "the work of this revision". git is the fallback for plain git checkouts.
var revisionCommands = [][]string{
	{"jj", "log", "--no-graph", "-r", "@", "-T", "change_id"},
	{"git", "rev-parse", "HEAD"},
}

// singleLineRevision accepts only a non-empty, single-line revision token.
func singleLineRevision(out string) (string, bool) {
	revision := strings.TrimSpace(out)
	if revision == "" || strings.ContainsAny(revision, "\n\r") {
		return "", false
	}
	return revision, true
}

// resolveRepoRevision returns the current revision (jj change id, else git
// commit), or "unknown". A failing, missing, or garbled command is never fatal.
func resolveRepoRevision(root string, run commandRunner) (revision string) {
	defer func() {
		if recover() != nil {
			revision = ledgerUnknownCommit
		}
	}()
	for _, command := range revisionCommands {
		out, err := run(root, command[0], command[1:]...)
		if err != nil {
			continue
		}
		if revision, ok := singleLineRevision(out); ok {
			return revision
		}
	}
	return ledgerUnknownCommit
}

// repoRevision returns the current repository revision, or "unknown" when it
// cannot be determined. The ledger records provenance; it never fails a run
// because version-control metadata is unavailable.
func repoRevision(root string) string {
	return resolveRepoRevision(root, execCommandRunner)
}

// ledgerCleanupUnresolvedRevision explains why a requested cleanup was skipped.
const ledgerCleanupUnresolvedRevision = ledgerCleanupEnv + "=1 ignored: cleanup requires a resolvable revision " +
	"(jj change id or git commit) and this run resolved " + ledgerUnknownCommit +
	"; no tombstones were dropped"

// rederiveLedger builds the next ledger from the measured cases and the
// previously committed ledger. Fixtures that disappeared become tombstones; a
// tombstoned fixture that is measured again loses its tombstone. The returned
// note is non-empty only when a requested cleanup was refused.
//
// Cleanup drops every recorded tombstone whose removed_at differs from the
// current revision, including one recorded as "unknown", and keeps only this
// revision's. It is therefore gated on a resolvable current revision: when the
// current revision is "unknown" the comparison would be meaningless, so
// cleanup is refused and nothing is dropped.
func rederiveLedger(current []ledgerCase, previous parityLedger, revision string, reasons map[string]string, cleanup bool) (parityLedger, string) {
	note := ""
	if cleanup && revision == ledgerUnknownCommit {
		cleanup = false
		note = ledgerCleanupUnresolvedRevision
	}
	measured := make(map[string]bool, len(current))
	for _, c := range current {
		measured[c.Fixture] = true
	}
	next := parityLedger{Cases: append([]ledgerCase(nil), current...)}
	// One fixture yields at most one tombstone, even if a hand-edited previous
	// ledger recorded it as both a case and a tombstone.
	tombstoned := map[string]bool{}
	for _, tombstone := range previous.Tombstones {
		if measured[tombstone.Fixture] || tombstoned[tombstone.Fixture] {
			continue
		}
		if cleanup && tombstone.RemovedAt != revision {
			continue
		}
		tombstoned[tombstone.Fixture] = true
		next.Tombstones = append(next.Tombstones, tombstone)
	}
	for _, previousCase := range previous.Cases {
		if measured[previousCase.Fixture] || tombstoned[previousCase.Fixture] {
			continue
		}
		tombstoned[previousCase.Fixture] = true
		reason := ledgerUnexplained
		if why, ok := reasons[previousCase.Fixture]; ok {
			reason = why
		}
		next.Tombstones = append(next.Tombstones, ledgerTombstone{
			Fixture:    previousCase.Fixture,
			LastStatus: previousCase.Status,
			RemovedAt:  revision,
			Reason:     reason,
		})
	}
	sort.Slice(next.Cases, func(i, j int) bool { return next.Cases[i].Fixture < next.Cases[j].Fixture })
	sort.Slice(next.Tombstones, func(i, j int) bool { return next.Tombstones[i].Fixture < next.Tombstones[j].Fixture })
	return next, note
}

// diffLedger compares the measured cases against the committed ledger and
// returns a precise, sorted, per-case list of contract violations.
func diffLedger(current []ledgerCase, committed parityLedger) []string {
	committedCases := make(map[string]ledgerCase, len(committed.Cases))
	for _, c := range committed.Cases {
		committedCases[c.Fixture] = c
	}
	tombstoned := make(map[string]ledgerTombstone, len(committed.Tombstones))
	for _, t := range committed.Tombstones {
		tombstoned[t.Fixture] = t
	}

	var errs []string
	measured := make(map[string]bool, len(current))
	for _, c := range current {
		measured[c.Fixture] = true
		previous, ok := committedCases[c.Fixture]
		if !ok {
			line := fmt.Sprintf("ADDED %s status=%s", c.Fixture, c.Status)
			if _, buried := tombstoned[c.Fixture]; buried {
				line += "; the ledger tombstones this fixture, so remove the tombstone deliberately"
			}
			errs = append(errs, line+"; rederive the ledger with "+parityRederiveEnv+"=1")
			continue
		}
		if previous.Status != c.Status {
			errs = append(errs, fmt.Sprintf("STATUS-CHANGED %s %s -> %s; rederive the ledger with %s=1",
				c.Fixture, previous.Status, c.Status, parityRederiveEnv))
		}
	}
	for _, c := range committed.Cases {
		if !measured[c.Fixture] {
			errs = append(errs, fmt.Sprintf("REMOVED %s status=%s; rederive the ledger with %s=1 and explain it via %s",
				c.Fixture, c.Status, parityRederiveEnv, ledgerTombstoneReasonEnv))
		}
	}
	for _, t := range committed.Tombstones {
		if t.Reason == ledgerUnexplained {
			errs = append(errs, fmt.Sprintf("UNEXPLAINED-TOMBSTONE %s last_status=%s removed_at=%s; explain it via %s",
				t.Fixture, t.LastStatus, t.RemovedAt, ledgerTombstoneReasonEnv))
		}
	}
	sort.Strings(errs)
	return errs
}

// ---------------------------------------------------------------------------
// Reporting
//
// Report CONTENT is informational only: it never feeds classification, the
// partial/divergence ratchets, or rederive behavior. The one exit-status
// coupling is deliberate: when LETGO_PARITY_REPORT explicitly requests a
// report file, a failed write is reported as a test error rather than
// silently dropped.
// ---------------------------------------------------------------------------

const (
	parityReportEnv   = "LETGO_PARITY_REPORT"
	auditBeginMarker  = "LETGO-PARITY-AUDIT-BEGIN"
	auditEndMarker    = "LETGO-PARITY-AUDIT-END"
	auditEntryMarker  = "LETGO-PARITY-AUDIT-ENTRY"
	auditErrorMarker  = "LETGO-PARITY-AUDIT-ERROR"
	fullCategory      = "full strict parity"
	partialCategory   = "partial strict capability"
	divergedCategory  = "strict output divergence"
	fatalCategory     = "fatal"
	strictProvenNote  = "strict success proves zero IR fallback"
	auditFailedPrefix = "audit unavailable"
)

// fallbackEntry is one measured [defn, reason] pair recorded by the hybrid
// fallback path in clojure.core/defn. The reason text is kept raw.
type fallbackEntry struct {
	Defn   string
	Reason string
}

// engineAudit is the audit outcome for one fixture on one engine. Err is set
// when the audit run itself could not be measured; it degrades the report for
// that fixture and nothing else.
type engineAudit struct {
	Engine  string
	Entries []fallbackEntry
	Err     string
}

// fixtureAudit collects the per-engine audit outcomes for a single fixture.
type fixtureAudit struct {
	Engines []engineAudit
}

// needsAudit reports whether a fixture requires an extra audit run. A full
// strict parity fixture does not: strict success already proves zero fallback.
func needsAudit(census corpusCensus, fixture string) bool {
	return !census.Full[fixture]
}

// fixturesNeedingAudit returns the sorted subset of the corpus that needs an
// audit run to explain its category.
func fixturesNeedingAudit(names []string, census corpusCensus) []string {
	var out []string
	for _, name := range names {
		if needsAudit(census, name) {
			out = append(out, name)
		}
	}
	sort.Strings(out)
	return out
}

// parseAuditLog extracts the measured fallback pairs from a marked audit
// stderr block. A missing or unterminated block is an audit failure, not a
// fixture failure.
func parseAuditLog(stderr string) ([]fallbackEntry, error) {
	inBlock := false
	closed := false
	var entries []fallbackEntry
	for _, line := range strings.Split(stderr, "\n") {
		line = strings.TrimRight(line, "\r")
		switch {
		case line == auditBeginMarker:
			inBlock = true
			continue
		case line == auditEndMarker:
			if !inBlock {
				return nil, fmt.Errorf("audit end marker without begin marker")
			}
			inBlock = false
			closed = true
			continue
		}
		if !inBlock || !strings.HasPrefix(line, auditEntryMarker+"\t") {
			continue
		}
		fields := strings.Split(line, "\t")
		if len(fields) != 3 {
			return nil, fmt.Errorf("malformed audit entry %q", line)
		}
		defn, err := strconv.Unquote(fields[1])
		if err != nil {
			return nil, fmt.Errorf("unquote audit defn %q: %w", fields[1], err)
		}
		reason, err := strconv.Unquote(fields[2])
		if err != nil {
			return nil, fmt.Errorf("unquote audit reason %q: %w", fields[2], err)
		}
		entries = append(entries, fallbackEntry{Defn: defn, Reason: reason})
	}
	if !closed {
		return nil, fmt.Errorf("audit log block missing %s/%s markers", auditBeginMarker, auditEndMarker)
	}
	return entries, nil
}

// auditFromResult converts a raw audit execution into a reportable outcome.
func auditFromResult(engine string, result executionResult) engineAudit {
	if !result.OK() {
		return engineAudit{Engine: engine, Err: fmt.Sprintf("%s: audit run failed; %s", auditFailedPrefix, result.evidence())}
	}
	entries, err := parseAuditLog(result.Stderr)
	if err != nil {
		return engineAudit{Engine: engine, Err: fmt.Sprintf("%s: %v", auditFailedPrefix, err)}
	}
	return engineAudit{Engine: engine, Entries: entries}
}

// fixtureCategory names every measured category a fixture belongs to. Partial
// and divergence are independent facts, so a fixture may report both.
func fixtureCategory(census corpusCensus, fixture string) string {
	var parts []string
	if census.Full[fixture] {
		parts = append(parts, fullCategory)
	}
	if census.Partial[fixture] {
		parts = append(parts, partialCategory)
	}
	if census.Diverged[fixture] {
		parts = append(parts, divergedCategory)
	}
	if len(parts) == 0 {
		return fatalCategory
	}
	return strings.Join(parts, "+")
}

// corpusFatalPrefixes are the whole-corpus fatal messages classifyCorpus emits;
// each is a prefix of exactly "<prefix><fixture>".
var corpusFatalPrefixes = []string{"duplicate result for ", "unexpected result for ", "missing result for "}

// fatalReasons returns the fatal messages attributed to one fixture, sorted.
// It matches only the two message shapes the classifier actually produces:
// classifyObservation's "<fixture>: ..." per-fixture messages and
// classifyCorpus's whole-corpus "<kind> result for <fixture>" messages. Loose
// substring matching would misattribute (e.g. "extra-var.lg" to "var.lg").
func fatalReasons(census corpusCensus, fixture string) []string {
	var out []string
	for _, fatal := range census.Fatal {
		matched := strings.HasPrefix(fatal, fixture+": ")
		for _, prefix := range corpusFatalPrefixes {
			if fatal == prefix+fixture {
				matched = true
			}
		}
		if matched {
			out = append(out, fatal)
		}
	}
	sort.Strings(out)
	return out
}

type reasonBucketKey struct {
	Reason string
	Site   string
}

// buildParityReport renders a deterministic, sorted coverage report over the
// whole discovered corpus.
func buildParityReport(names []string, census corpusCensus, audits map[string]fixtureAudit) []byte {
	ordered := append([]string(nil), names...)
	sort.Strings(ordered)

	buckets := map[string]map[string]bool{}
	var b strings.Builder
	b.WriteString("# let-go engine parity coverage report\n")
	b.WriteString("# Informational only: this content does not affect classification, ratchets, or rederive.\n")
	b.WriteString("# Fixture lines: <category> then the measured reason for every non-full fixture.\n\n")

	for _, name := range ordered {
		category := fixtureCategory(census, name)
		b.WriteString(fmt.Sprintf("FIXTURE %s category=%s", name, category))
		if !needsAudit(census, name) {
			b.WriteString(" fallbacks=0 reason=" + strconv.Quote(strictProvenNote) + "\n")
			continue
		}

		audit, measured := audits[name]
		var details []string
		reasons := map[string]bool{}
		if !measured {
			details = append(details, fmt.Sprintf("  DETAIL %s: %s: no audit run recorded", name, auditFailedPrefix))
			reasons[auditFailedPrefix+": no audit run recorded"] = true
		}
		engines := append([]engineAudit(nil), audit.Engines...)
		sort.Slice(engines, func(i, j int) bool { return engines[i].Engine < engines[j].Engine })
		for _, engine := range engines {
			if engine.Err != "" {
				details = append(details, fmt.Sprintf("  DETAIL %s engine=%s %s", name, engine.Engine, engine.Err))
				reasons[engine.Err] = true
				continue
			}
			if len(engine.Entries) == 0 {
				details = append(details, fmt.Sprintf("  DETAIL %s engine=%s no IR fallback recorded", name, engine.Engine))
				continue
			}
			entries := append([]fallbackEntry(nil), engine.Entries...)
			sort.Slice(entries, func(i, j int) bool {
				if entries[i].Defn != entries[j].Defn {
					return entries[i].Defn < entries[j].Defn
				}
				return entries[i].Reason < entries[j].Reason
			})
			for _, entry := range entries {
				details = append(details, fmt.Sprintf("  DETAIL %s engine=%s defn=%s reason=%s",
					name, engine.Engine, entry.Defn, strconv.Quote(entry.Reason)))
				reasons[entry.Reason] = true
				site := fmt.Sprintf("%s engine=%s defn=%s", name, engine.Engine, entry.Defn)
				if buckets[entry.Reason] == nil {
					buckets[entry.Reason] = map[string]bool{}
				}
				buckets[entry.Reason][site] = true
			}
		}
		for _, fatal := range fatalReasons(census, name) {
			details = append(details, "  DETAIL "+name+" fatal="+strconv.Quote(fatal))
			reasons[fatal] = true
		}

		distinct := sortedKeys(reasons)
		if len(distinct) == 0 {
			b.WriteString(" fallbacks=0 reason=" + strconv.Quote("no measured IR fallback; see engine detail") + "\n")
		} else {
			quoted := make([]string, len(distinct))
			for i, reason := range distinct {
				quoted[i] = strconv.Quote(reason)
			}
			b.WriteString(fmt.Sprintf(" distinct_reasons=%d reason=%s\n", len(distinct), strings.Join(quoted, " ")))
		}
		sort.Strings(details)
		for _, detail := range details {
			b.WriteString(detail + "\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("AGGREGATE %s\n", census.Summary()))
	b.WriteString(fmt.Sprintf("AGGREGATE audited_fixtures=%d\n", len(fixturesNeedingAudit(ordered, census))))
	b.WriteString(fmt.Sprintf("BUCKET-CENSUS distinct_reasons=%d\n", len(buckets)))
	// Triage order: largest gap first, then reason text for determinism.
	bucketReasons := sortedKeys(mapKeysAsSet(buckets))
	sort.SliceStable(bucketReasons, func(i, j int) bool {
		return len(buckets[bucketReasons[i]]) > len(buckets[bucketReasons[j]])
	})
	for _, reason := range bucketReasons {
		sites := sortedKeys(buckets[reason])
		b.WriteString(fmt.Sprintf("BUCKET count=%d reason=%s\n", len(sites), strconv.Quote(reason)))
		for _, site := range sites {
			b.WriteString("  BUCKET-SITE " + site + "\n")
		}
	}
	return []byte(b.String())
}

func mapKeysAsSet(buckets map[string]map[string]bool) map[string]bool {
	set := make(map[string]bool, len(buckets))
	for key := range buckets {
		set[key] = true
	}
	return set
}

func sortedKeys(set map[string]bool) []string {
	keys := make([]string, 0, len(set))
	for key := range set {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

// writeParityReport persists the report when LETGO_PARITY_REPORT names a path.
// The report is never committed and never becomes a baseline. An unset path is
// a no-op; a requested write that fails is returned to the caller, which
// reports it.
func writeParityReport(path string, data []byte) error {
	if path == "" {
		return nil
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write parity report %s: %w", path, err)
	}
	return nil
}

type fixtureMode string

const (
	controlMode fixtureMode = "control"
	strictMode  fixtureMode = "strict"
	// auditMode is reporting-only: it never classifies a fixture. It reruns a
	// non-full fixture with IR compilation on and strict OFF so the hybrid
	// fallback path records why each defn could not be lowered.
	auditMode fixtureMode = "audit"
)

func (mode fixtureMode) compileBindings() string {
	switch mode {
	case controlMode:
		return "(binding [*ir-compile* false *ir-compile-strict* false]"
	case auditMode:
		return `(let [audit-fallback-log (atom [])]
  (reset! audit-fallback-log [])
  (binding [*ir-compile* true
            *ir-compile-strict* false
            *ir-compile-verbose* true
            *ir-compile-fallback-log* audit-fallback-log]`
	case strictMode:
		return `(let [strict-fallback-log (atom [])]
  (reset! strict-fallback-log [])
  (binding [*ir-compile* true
            *ir-compile-strict* true
            *ir-compile-verbose* true
            *ir-compile-fallback-log* strict-fallback-log]`
	default:
		panic("unknown fixture mode: " + mode)
	}
}

func (mode fixtureMode) catchAndClose() string {
	switch mode {
	case controlMode:
		return `(catch e
      (throw e))))`
	case auditMode:
		// An audit failure is reported, never fatal: the marked block is emitted
		// after the try so a mid-fixture throw still yields whatever fallbacks
		// were recorded before it.
		return fmt.Sprintf(`(catch e
      (binding [*out* *err*]
        (println (str %q " " e)))))
  (binding [*out* *err*]
    (println %q)
    (doseq [entry (deref audit-fallback-log)]
      (println (str %q "\t"
                    (pr-str (str (first entry)))
                    "\t"
                    (pr-str (str (second entry))))))
    (println %q))))`, auditErrorMarker, auditBeginMarker, auditEntryMarker, auditEndMarker)
	case strictMode:
		return fmt.Sprintf(`(catch e
      (if (empty? (deref strict-fallback-log))
        (throw e)
        (do
          (binding [*out* *err*] (println (str e)))
          (os/exit %d)))))))`, strictRejectionExitCode)
	default:
		panic("unknown fixture mode: " + mode)
	}
}

func createFixtureWrapper(dest, source string, mode fixtureMode) error {
	content, err := os.ReadFile(source)
	if err != nil {
		return fmt.Errorf("read %s: %w", source, err)
	}
	wrapper := fmt.Sprintf(`(require 'ir.passes.pipeline)
%s
  (try
    (doseq [form (read-all-string %q)]
      (eval form))
    %s
`, mode.compileBindings(), string(content), mode.catchAndClose())
	if err := os.WriteFile(dest, []byte(wrapper), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", dest, err)
	}
	return nil
}

func runFixtureDetailed(bin, script string) executionResult {
	cmd := exec.Command(bin, script)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	result := executionResult{Stdout: normalizeOutput(stdout.String()), Stderr: stderr.String()}
	if err == nil {
		return result
	}
	if exitError, ok := err.(*exec.ExitError); ok {
		result.ExitCode = exitError.ExitCode()
	} else {
		result.ExitCode = -1
		result.StartError = err.Error()
	}
	return result
}

func discoverFixtureNames(dir string) ([]string, []string, error) {
	paths, err := filepath.Glob(filepath.Join(dir, "*.lg"))
	if err != nil {
		return nil, nil, err
	}
	if len(paths) == 0 {
		return nil, nil, fmt.Errorf("no fixtures in %s", dir)
	}
	sort.Strings(paths)
	names := make([]string, len(paths))
	for i, path := range paths {
		names[i] = filepath.Base(path)
	}
	return names, paths, nil
}

// TestEngineParityGate measures normalized engine-output parity and strict
// compilation capability over every test/gold-aot fixture. Strict success is a
// capability signal only: the wrapper does not prove that fixture source ran as
// generated Go or entered a native lowering seam.
func TestEngineParityGate(t *testing.T) {
	if testing.Short() {
		t.Skip("engine parity gate builds let-go twice; run via `make engine-parity-gate`")
	}
	root := repoRoot(t)
	if err := genmanifest.CheckTreeManifest(filepath.Join(root, "pkg/rt/core_go_lowered")); err != nil {
		t.Fatalf("generated lowered tree is not valid: %v\nrun `make generate` (or `make lowered`) and retry", err)
	}

	bc := buildLGTags(t, "")
	gogen := buildLGTags(t, "gogen_ir")
	fixtureDir := filepath.Join(root, "test/gold-aot")
	names, paths, err := discoverFixtureNames(fixtureDir)
	if err != nil {
		t.Fatal(err)
	}
	wrapDir := t.TempDir()

	// Falsify message-based ownership with both built engines without adding the
	// synthetic fixture to the discovered gold-fixture census.
	const ownershipFixture = "strict-rejection-ownership.lg"
	ownershipSource := filepath.Join(wrapDir, ownershipFixture)
	if err := os.WriteFile(ownershipSource, []byte(`(when *ir-compile-strict*
  (throw "user error: ir-compile-strict: boom"))
`), 0o644); err != nil {
		t.Fatal(err)
	}
	ownershipControl := filepath.Join(wrapDir, ownershipFixture+".control.lg")
	ownershipStrict := filepath.Join(wrapDir, ownershipFixture+".strict.lg")
	if err := createFixtureWrapper(ownershipControl, ownershipSource, controlMode); err != nil {
		t.Fatal(err)
	}
	if err := createFixtureWrapper(ownershipStrict, ownershipSource, strictMode); err != nil {
		t.Fatal(err)
	}
	ownershipObservation := observation(
		ownershipFixture,
		runFixtureDetailed(bc, ownershipControl),
		runFixtureDetailed(gogen, ownershipControl),
		runFixtureDetailed(bc, ownershipStrict),
		runFixtureDetailed(gogen, ownershipStrict),
	)
	for engine, result := range map[string]executionResult{
		"bytecode": ownershipObservation.StrictBC,
		"gogen_ir": ownershipObservation.StrictGogen,
	} {
		if result.strictRejection() || result.OK() || !strings.Contains(result.evidence(), "user error: ir-compile-strict: boom") {
			t.Fatalf("%s synthetic runtime throw was not preserved as a fatal result: %s", engine, result.evidence())
		}
	}
	ownershipCensus := classifyCorpus([]string{ownershipFixture}, []fixtureObservation{ownershipObservation})
	if len(ownershipCensus.Fatal) == 0 || !strings.Contains(ownershipCensus.Fatal[0], "user error: ir-compile-strict: boom") {
		t.Fatalf("synthetic runtime throw classification = %s fatal=%v; want Fatal with preserved output evidence", ownershipCensus.Summary(), ownershipCensus.Fatal)
	}
	if ownershipCensus.Partial[ownershipFixture] {
		t.Fatalf("synthetic runtime throw classified as partial: %v", ownershipCensus.Partial)
	}

	observations := make([]fixtureObservation, 0, len(paths))
	sourceByName := make(map[string]string, len(paths))
	for i, path := range paths {
		name := names[i]
		sourceByName[name] = path
		controlPath := filepath.Join(wrapDir, name+".control.lg")
		strictPath := filepath.Join(wrapDir, name+".strict.lg")
		if err := createFixtureWrapper(controlPath, path, controlMode); err != nil {
			t.Fatalf("create control wrapper for %s: %v", name, err)
		}
		if err := createFixtureWrapper(strictPath, path, strictMode); err != nil {
			t.Fatalf("create strict wrapper for %s: %v", name, err)
		}
		observations = append(observations, observation(
			name,
			runFixtureDetailed(bc, controlPath),
			runFixtureDetailed(gogen, controlPath),
			runFixtureDetailed(bc, strictPath),
			runFixtureDetailed(gogen, strictPath),
		))
	}

	census := classifyCorpus(names, observations)

	// Reporting stage: measure and print why every non-full fixture is not at
	// full strict parity. This runs after classification and changes nothing
	// about it; audit failures degrade the report only.
	audits := map[string]fixtureAudit{}
	for _, name := range fixturesNeedingAudit(names, census) {
		source, ok := sourceByName[name]
		if !ok {
			audits[name] = fixtureAudit{Engines: []engineAudit{{
				Engine: "harness",
				Err:    fmt.Sprintf("%s: no discovered source for %s", auditFailedPrefix, name),
			}}}
			continue
		}
		auditPath := filepath.Join(wrapDir, name+".audit.lg")
		if err := createFixtureWrapper(auditPath, source, auditMode); err != nil {
			audits[name] = fixtureAudit{Engines: []engineAudit{{
				Engine: "harness",
				Err:    fmt.Sprintf("%s: create audit wrapper: %v", auditFailedPrefix, err),
			}}}
			continue
		}
		audits[name] = fixtureAudit{Engines: []engineAudit{
			auditFromResult("bytecode", runFixtureDetailed(bc, auditPath)),
			auditFromResult("gogen_ir", runFixtureDetailed(gogen, auditPath)),
		}}
	}
	report := buildParityReport(names, census, audits)
	t.Logf("parity coverage report:\n%s", report)
	if err := writeParityReport(os.Getenv(parityReportEnv), report); err != nil {
		t.Error(err)
	}

	for _, fatal := range census.Fatal {
		t.Error("FATAL: " + fatal)
	}
	if !census.Valid() {
		t.Fatalf("invalid corpus census; baselines left unchanged: %s", census.Summary())
	}
	for _, obs := range observations {
		if census.Full[obs.Fixture] {
			t.Logf("FULL STRICT PARITY %s: output=%q", obs.Fixture, obs.StrictBC.Stdout)
		}
		if census.Partial[obs.Fixture] {
			t.Logf("PARTIAL STRICT CAPABILITY %s: bytecode={%s} gogen_ir={%s}", obs.Fixture, obs.StrictBC.evidence(), obs.StrictGogen.evidence())
		}
		if census.Diverged[obs.Fixture] {
			t.Logf("STRICT OUTPUT DIVERGENCE %s: bytecode={%s} gogen_ir={%s}", obs.Fixture, obs.StrictBC.evidence(), obs.StrictGogen.evidence())
		}
	}

	partialPath, divergencePath := baselinePaths(root)
	measuredCases := censusLedgerCases(names, census, audits)
	if os.Getenv(parityRederiveEnv) == "1" {
		previous, err := parseLedger(ledgerPath(root))
		if err != nil {
			if !os.IsNotExist(errors.Unwrap(err)) {
				t.Fatal(err)
			}
			previous = parityLedger{}
		}
		next, cleanupNote := rederiveLedger(
			measuredCases,
			previous,
			repoRevision(root),
			parseTombstoneReasons(os.Getenv(ledgerTombstoneReasonEnv)),
			os.Getenv(ledgerCleanupEnv) == "1",
		)
		if cleanupNote != "" {
			t.Log(cleanupNote)
		}
		if err := writeCensusBaselines(root, census, next); err != nil {
			t.Fatal(err)
		}
		t.Logf("rederived %s, %s and %s from complete %d-fixture census",
			ledgerPath(root), partialPath, divergencePath, len(names))
		return
	}
	committedLedger, err := parseLedger(ledgerPath(root))
	if err != nil {
		t.Fatal(err)
	}
	for _, ledgerError := range diffLedger(measuredCases, committedLedger) {
		t.Error(ledgerError)
	}
	partialBaseline, err := readBaseline(partialPath)
	if err != nil {
		t.Fatal(err)
	}
	divergenceBaseline, err := readBaseline(divergencePath)
	if err != nil {
		t.Fatal(err)
	}
	for _, ratchetError := range append(
		validateRatchet("partial", census.Partial, partialBaseline),
		validateRatchet("strict divergence", census.Diverged, divergenceBaseline)...,
	) {
		t.Error(ratchetError)
	}
	t.Log(census.Summary())
}
