package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/perfdata"
)

// seedFixture describes one mock timeline snapshot.
type seedFixture struct {
	stamp    string // filename timestamp, e.g. "20260801T010134Z"
	sha      string
	anchorNs float64
	// benches maps benchmark name to its ns/op in this snapshot.
	benches map[string]float64
	// iters maps benchmark name to its b.N in this snapshot (optional).
	iters map[string]int64
	arch  string
	model string
}

func (f seedFixture) machine() perfdata.Machine {
	arch, model := f.arch, f.model
	if arch == "" {
		arch = "amd64"
	}
	if model == "" {
		model = "AMD EPYC 7763"
	}
	return perfdata.Machine{OS: "linux", Arch: arch, NumCPU: 16, CPUModel: model, GoVersion: "go1.26.4"}
}

// writeSeedFixtures materialises snapshots into a fresh timeline dir and
// returns (timelineDir, baselinePath).
func writeSeedFixtures(t *testing.T, fixtures []seedFixture) (string, string) {
	t.Helper()
	tmp := t.TempDir()
	timeline := filepath.Join(tmp, "timeline")
	if err := os.MkdirAll(timeline, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, f := range fixtures {
		m := f.machine()
		entries := map[string]BenchmarkEntry{}
		for name, ns := range f.benches {
			e := BenchmarkEntry{NSPerOp: ns, RatioToAnchor: ns / f.anchorNs}
			if n, ok := f.iters[name]; ok {
				e.Samples = []BenchmarkSample{{Iterations: n, NSPerOp: ns}}
			}
			entries[name] = e
		}
		b := Baseline{
			Version: schemaVersion,
			Machines: map[string]MachineBaseline{
				perfdata.MachineKey(m): {
					CapturedAt:    f.stamp,
					CapturedAtSHA: f.sha,
					Machine:       m,
					Anchor: AnchorRecord{
						Name: anchorName, Package: anchorPackage,
						NSPerOp: f.anchorNs, Iterations: 1000000000,
					},
					Benchmarks: entries,
				},
			},
		}
		data, err := json.Marshal(b)
		if err != nil {
			t.Fatal(err)
		}
		name := fmt.Sprintf("%s-%s-%s.json", f.stamp, f.sha, slugify(perfdata.MachineKey(m)))
		if err := os.WriteFile(filepath.Join(timeline, name), data, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return timeline, filepath.Join(tmp, "baseline.json")
}

// runSeed executes seedBaseline with stdout captured, returning the parsed
// baseline and everything printed.
func runSeed(t *testing.T, timeline, baselinePath string, opt seedOptions) (Baseline, string) {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w
	done := make(chan string)
	go func() {
		b, _ := io.ReadAll(r)
		done <- string(b)
	}()
	func() {
		defer func() {
			w.Close()
			os.Stdout = old
		}()
		seedBaseline(baselinePath, timeline, opt)
	}()
	out := <-done

	data, err := os.ReadFile(baselinePath)
	if err != nil {
		t.Fatalf("read seeded baseline: %v", err)
	}
	var got Baseline
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("parse seeded baseline: %v", err)
	}
	return got, out
}

func onlyMachine(t *testing.T, b Baseline) MachineBaseline {
	t.Helper()
	if len(b.Machines) != 1 {
		t.Fatalf("machines = %d, want 1", len(b.Machines))
	}
	for _, mb := range b.Machines {
		return mb
	}
	return MachineBaseline{}
}

// The point of the window: the newest snapshot does not decide the baseline.
// Five runs of the same benchmark, the newest of which is a slow observation
// inside the anchor tolerance — the stored value must be the window median.
func TestSeedBaselineReducesWindowNotNewest(t *testing.T) {
	ns := []float64{112, 100, 101, 99, 100} // newest first once sorted
	var fx []seedFixture
	for i, v := range ns {
		fx = append(fx, seedFixture{
			stamp:    fmt.Sprintf("2026080%dT010000Z", 5-i),
			sha:      fmt.Sprintf("%012d", i),
			anchorNs: 1.5,
			benches:  map[string]float64{"test.BenchmarkA": v},
		})
	}
	timeline, baselinePath := writeSeedFixtures(t, fx)
	got, _ := runSeed(t, timeline, baselinePath, defaultSeedOptions())

	mb := onlyMachine(t, got)
	e, ok := mb.Benchmarks["test.BenchmarkA"]
	if !ok {
		t.Fatal("test.BenchmarkA missing from seeded baseline")
	}
	// median{112,100,101,99,100} = 100; the newest sample (112) is 12% high.
	if e.NSPerOp != 100 {
		t.Errorf("ns_per_op = %v, want 100 (window median, not the newest sample)", e.NSPerOp)
	}
	// Identity still comes from the newest surviving snapshot.
	if mb.CapturedAtSHA != "000000000000" {
		t.Errorf("captured_at_sha = %q, want the newest snapshot's SHA", mb.CapturedAtSHA)
	}
}

// The negative result, kept as a test: a snapshot can sit far off its window's
// ANCHOR and still be a perfectly good capture. Measured over the recent amd64
// corpus, one snapshot ran 22.4% fast in raw ns/op across all 162 of its
// benchmarks and agreed with its window on every ratio to within 0.1%. Gating
// on anchor deviation would discard it — and two more like it.
func TestSeedBaselineKeepsUniformlyFastSnapshot(t *testing.T) {
	// Same ratio (66.67) everywhere; the fast host moves anchor and ns/op together.
	fx := []seedFixture{
		{stamp: "20260805T010000Z", sha: "aaaaaaaaaaaa", anchorNs: 1.1, benches: map[string]float64{"test.BenchmarkA": 73.337}},
		{stamp: "20260804T010000Z", sha: "bbbbbbbbbbbb", anchorNs: 1.5, benches: map[string]float64{"test.BenchmarkA": 100}},
		{stamp: "20260803T010000Z", sha: "cccccccccccc", anchorNs: 1.5, benches: map[string]float64{"test.BenchmarkA": 100}},
		{stamp: "20260802T010000Z", sha: "dddddddddddd", anchorNs: 1.52, benches: map[string]float64{"test.BenchmarkA": 101.333}},
		{stamp: "20260801T010000Z", sha: "eeeeeeeeeeee", anchorNs: 1.48, benches: map[string]float64{"test.BenchmarkA": 98.667}},
	}
	timeline, baselinePath := writeSeedFixtures(t, fx)
	got, out := runSeed(t, timeline, baselinePath, defaultSeedOptions())

	if strings.Contains(out, "rejected") {
		t.Errorf("a uniformly-fast host is not a bad capture; output was:\n%s", out)
	}
	mb := onlyMachine(t, got)
	if mb.CapturedAtSHA != "aaaaaaaaaaaa" {
		t.Errorf("captured_at_sha = %q, want the newest snapshot kept", mb.CapturedAtSHA)
	}
	if got := mb.Benchmarks["test.BenchmarkA"].RatioToAnchor; math.Abs(got-66.667) > 0.01 {
		t.Errorf("ratio_to_anchor = %v, want ~66.667 (unaffected by host speed)", got)
	}
}

// The case that does need rejecting: the anchor moved and the benchmarks did
// not, so every ratio from this snapshot is uniformly wrong while its raw
// numbers look ordinary.
func TestSeedBaselineRejectsMixedCapture(t *testing.T) {
	fx := []seedFixture{
		// anchor 22% fast, benchmark unchanged -> ratio 28% high.
		{stamp: "20260805T010000Z", sha: "aaaaaaaaaaaa", anchorNs: 1.17, benches: map[string]float64{"test.BenchmarkA": 100, "test.BenchmarkB": 200, "test.BenchmarkC": 300}},
		{stamp: "20260804T010000Z", sha: "bbbbbbbbbbbb", anchorNs: 1.5, benches: map[string]float64{"test.BenchmarkA": 100, "test.BenchmarkB": 200, "test.BenchmarkC": 300}},
		{stamp: "20260803T010000Z", sha: "cccccccccccc", anchorNs: 1.5, benches: map[string]float64{"test.BenchmarkA": 100, "test.BenchmarkB": 200, "test.BenchmarkC": 300}},
		{stamp: "20260802T010000Z", sha: "dddddddddddd", anchorNs: 1.5, benches: map[string]float64{"test.BenchmarkA": 100, "test.BenchmarkB": 200, "test.BenchmarkC": 300}},
		{stamp: "20260801T010000Z", sha: "eeeeeeeeeeee", anchorNs: 1.5, benches: map[string]float64{"test.BenchmarkA": 100, "test.BenchmarkB": 200, "test.BenchmarkC": 300}},
	}
	timeline, baselinePath := writeSeedFixtures(t, fx)
	got, out := runSeed(t, timeline, baselinePath, defaultSeedOptions())

	if !strings.Contains(out, "rejected aaaaaaaaaaaa") {
		t.Errorf("want the mixed capture rejected; output was:\n%s", out)
	}
	mb := onlyMachine(t, got)
	// Identity must come from the newest SURVIVING snapshot.
	if mb.CapturedAtSHA != "bbbbbbbbbbbb" {
		t.Errorf("captured_at_sha = %q, want bbbbbbbbbbbb (newest survivor)", mb.CapturedAtSHA)
	}
	if got := mb.Benchmarks["test.BenchmarkA"].RatioToAnchor; math.Abs(got-66.667) > 0.01 {
		t.Errorf("ratio_to_anchor = %v, want ~66.667 (the mixed snapshot excluded)", got)
	}
}

// ratio_to_anchor is what `check` compares, so it must agree with the ns/op and
// anchor actually stored beside it — not be reduced independently.
func TestSeedBaselineRatioAgreesWithStoredValues(t *testing.T) {
	fx := []seedFixture{
		{stamp: "20260805T010000Z", sha: "aaaaaaaaaaaa", anchorNs: 1.6, benches: map[string]float64{"test.BenchmarkA": 90, "test.BenchmarkB": 300}},
		{stamp: "20260804T010000Z", sha: "bbbbbbbbbbbb", anchorNs: 1.5, benches: map[string]float64{"test.BenchmarkA": 100, "test.BenchmarkB": 310}},
		{stamp: "20260803T010000Z", sha: "cccccccccccc", anchorNs: 1.4, benches: map[string]float64{"test.BenchmarkA": 110, "test.BenchmarkB": 290}},
	}
	timeline, baselinePath := writeSeedFixtures(t, fx)
	got, _ := runSeed(t, timeline, baselinePath, defaultSeedOptions())

	mb := onlyMachine(t, got)
	for name, e := range mb.Benchmarks {
		want := e.NSPerOp / mb.Anchor.NSPerOp
		if math.Abs(e.RatioToAnchor-want) > 1e-9 {
			t.Errorf("%s: ratio_to_anchor = %v, want %v (ns_per_op / anchor)", name, e.RatioToAnchor, want)
		}
	}
}

// b.N is an output of the timing loop, so it moves when per-op cost moves.
// Reducing across snapshots taken at different N mixes two protocols; the seed
// reports that rather than silently dropping the benchmark.
func TestSeedBaselineReportsMovedIterations(t *testing.T) {
	fx := []seedFixture{
		{stamp: "20260805T010000Z", sha: "aaaaaaaaaaaa", anchorNs: 1.5,
			benches: map[string]float64{"test.BenchmarkA": 100, "test.BenchmarkSteady": 50},
			iters:   map[string]int64{"test.BenchmarkA": 400, "test.BenchmarkSteady": 1000}},
		{stamp: "20260804T010000Z", sha: "bbbbbbbbbbbb", anchorNs: 1.5,
			benches: map[string]float64{"test.BenchmarkA": 100, "test.BenchmarkSteady": 50},
			iters:   map[string]int64{"test.BenchmarkA": 700, "test.BenchmarkSteady": 1010}},
		{stamp: "20260803T010000Z", sha: "cccccccccccc", anchorNs: 1.5,
			benches: map[string]float64{"test.BenchmarkA": 100, "test.BenchmarkSteady": 50},
			iters:   map[string]int64{"test.BenchmarkA": 900, "test.BenchmarkSteady": 990}},
	}
	timeline, baselinePath := writeSeedFixtures(t, fx)
	got, out := runSeed(t, timeline, baselinePath, defaultSeedOptions())

	if !strings.Contains(out, "b.N moved") || !strings.Contains(out, "test.BenchmarkA") {
		t.Errorf("want a b.N movement warning naming test.BenchmarkA; output was:\n%s", out)
	}
	if strings.Contains(out, "test.BenchmarkSteady (b.N") {
		t.Errorf("test.BenchmarkSteady moved 2%% and should not be flagged; output was:\n%s", out)
	}
	// Reported, not excluded — dropping it would shrink the gate silently.
	if _, ok := onlyMachine(t, got).Benchmarks["test.BenchmarkA"]; !ok {
		t.Error("test.BenchmarkA should still be seeded; the movement is reported, not acted on")
	}
}

// The exclusion list is a decision, not a measurement, so it cannot notice a
// rename on its own. Seeding says when an entry matched nothing.
func TestSeedBaselineWarnsWhenExclusionMatchesNothing(t *testing.T) {
	fx := []seedFixture{
		{stamp: "20260805T010000Z", sha: "aaaaaaaaaaaa", anchorNs: 1.5,
			benches: map[string]float64{"test.BenchmarkA": 100}},
	}
	timeline, baselinePath := writeSeedFixtures(t, fx)
	_, out := runSeed(t, timeline, baselinePath, defaultSeedOptions())

	if !strings.Contains(out, "matched nothing") {
		t.Errorf("want a stale-exclusion warning; output was:\n%s", out)
	}
	for _, prefix := range unstableBenchmarks {
		if !strings.Contains(out, prefix) {
			t.Errorf("want %q named in the warning; output was:\n%s", prefix, out)
		}
	}
}

// A benchmark present in only a minority of the window has too few samples to
// reduce. `check` reporting it as NEW is more honest than seeding it from one
// observation.
func TestSeedBaselineSkipsBenchmarksBelowQuorum(t *testing.T) {
	fx := []seedFixture{
		{stamp: "20260805T010000Z", sha: "aaaaaaaaaaaa", anchorNs: 1.5,
			benches: map[string]float64{"test.BenchmarkA": 100, "test.BenchmarkNew": 42}},
		{stamp: "20260804T010000Z", sha: "bbbbbbbbbbbb", anchorNs: 1.5, benches: map[string]float64{"test.BenchmarkA": 100}},
		{stamp: "20260803T010000Z", sha: "cccccccccccc", anchorNs: 1.5, benches: map[string]float64{"test.BenchmarkA": 100}},
		{stamp: "20260802T010000Z", sha: "dddddddddddd", anchorNs: 1.5, benches: map[string]float64{"test.BenchmarkA": 100}},
		{stamp: "20260801T010000Z", sha: "eeeeeeeeeeee", anchorNs: 1.5, benches: map[string]float64{"test.BenchmarkA": 100}},
	}
	timeline, baselinePath := writeSeedFixtures(t, fx)
	got, out := runSeed(t, timeline, baselinePath, defaultSeedOptions())

	mb := onlyMachine(t, got)
	if _, ok := mb.Benchmarks["test.BenchmarkNew"]; ok {
		t.Error("test.BenchmarkNew appeared in 1 of 5 snapshots and should not be seeded")
	}
	if _, ok := mb.Benchmarks["test.BenchmarkA"]; !ok {
		t.Error("test.BenchmarkA appeared in all 5 snapshots and should be seeded")
	}
	if !strings.Contains(out, "below quorum") {
		t.Errorf("want the skip reported; output was:\n%s", out)
	}
}

// A positional split on "-" yields a plausible-looking wrong SHA for a name
// that does not match the expected shape. Skip and say so instead.
func TestSeedBaselineSkipsMalformedFilenames(t *testing.T) {
	fx := []seedFixture{
		{stamp: "20260805T010000Z", sha: "aaaaaaaaaaaa", anchorNs: 1.5, benches: map[string]float64{"test.BenchmarkA": 100}},
	}
	timeline, baselinePath := writeSeedFixtures(t, fx)
	if err := os.WriteFile(filepath.Join(timeline, "summary.json"), []byte(`{"version":2}`), 0o644); err != nil {
		t.Fatal(err)
	}
	_, out := runSeed(t, timeline, baselinePath, defaultSeedOptions())

	if !strings.Contains(out, "summary.json") {
		t.Errorf("want summary.json reported as skipped; output was:\n%s", out)
	}
}

// The architecture filter must hold against the file's CONTENT. A snapshot
// named for one machine while carrying another is the divergence that put an
// EPYC profile under an Intel key once already.
func TestSeedBaselineWarnsOnSlugContentMismatch(t *testing.T) {
	fx := []seedFixture{
		{stamp: "20260805T010000Z", sha: "aaaaaaaaaaaa", anchorNs: 1.5, benches: map[string]float64{"test.BenchmarkA": 100}},
	}
	timeline, baselinePath := writeSeedFixtures(t, fx)
	// Rename the file so its slug claims a different CPU than it carries.
	old := filepath.Join(timeline, "20260805T010000Z-aaaaaaaaaaaa-amd64-amd-epyc-7763.json")
	renamed := filepath.Join(timeline, "20260805T010000Z-aaaaaaaaaaaa-amd64-intel-r-xeon-r-platinum-8573c.json")
	if err := os.Rename(old, renamed); err != nil {
		t.Fatal(err)
	}
	got, out := runSeed(t, timeline, baselinePath, defaultSeedOptions())

	if !strings.Contains(out, "is named for") {
		t.Errorf("want a slug/content mismatch warning; output was:\n%s", out)
	}
	// It is stored under the profile it carries, not the one it is named for.
	if _, ok := got.Machines["amd64/AMD EPYC 7763"]; !ok {
		t.Errorf("want the profile keyed by its content; got keys %v", got.Machines)
	}
}

// Overlapping prefixes: the shorter exclusion must not shadow the longer one
// into looking stale.
func TestSeedBaselineMarksEveryMatchingExclusion(t *testing.T) {
	fx := []seedFixture{
		{stamp: "20260805T010000Z", sha: "aaaaaaaaaaaa", anchorNs: 1.5, benches: map[string]float64{
			"test.BenchmarkA": 100,
			"github.com/nooga/let-go/test.BenchmarkClojureTestSuite [bytecode]":                    5000,
			"github.com/nooga/let-go/test.BenchmarkClojureTestSuiteCompileAndRun [total_bytecode]": 9000,
		}},
	}
	timeline, baselinePath := writeSeedFixtures(t, fx)
	got, out := runSeed(t, timeline, baselinePath, defaultSeedOptions())

	if strings.Contains(out, "matched nothing") {
		t.Errorf("both exclusions are present and should be marked matched; output was:\n%s", out)
	}
	for name := range onlyMachine(t, got).Benchmarks {
		if strings.Contains(name, "ClojureTestSuite") {
			t.Errorf("%s should have been excluded from the seed", name)
		}
	}
}

func TestMedianOf(t *testing.T) {
	cases := []struct {
		in   []float64
		want float64
	}{
		{nil, 0},
		{[]float64{5}, 5},
		{[]float64{4, 2}, 3},
		{[]float64{9, 1, 5}, 5},
		{[]float64{4, 1, 3, 2}, 2.5},
		// Unlike reduceSamples, no rep is discarded: order must not matter.
		{[]float64{100, 1, 1, 1, 1}, 1},
	}
	for _, c := range cases {
		if got := medianOf(c.in); got != c.want {
			t.Errorf("medianOf(%v) = %v, want %v", c.in, got, c.want)
		}
	}
}

// The M3 preservation carry-over keys on Arch as well as the model string.
// Without the Arch guard an amd64 CPUModel containing "M3" matches too, and the
// stale carried-over entry overwrites the fresh seed this very run computed for
// that tier — a silent revert to whatever the old baseline held.
func TestSeedBaselinePreservesOnlyOffArchM3(t *testing.T) {
	const amdM3 = "AMD EPYC M3000" // an amd64 model string that contains "M3"

	var fx []seedFixture
	for i := 0; i < 5; i++ {
		fx = append(fx, seedFixture{
			stamp:    fmt.Sprintf("2026080%dT010000Z", 5-i),
			sha:      fmt.Sprintf("%012d", i),
			anchorNs: 1.5,
			benches:  map[string]float64{"test.BenchmarkA": 100},
			arch:     "amd64",
			model:    amdM3,
		})
	}
	timeline, baselinePath := writeSeedFixtures(t, fx)

	// A prior baseline holding both tiers: the local arm64 M3 that seeding
	// cannot derive, and a stale entry for the very amd64 tier being seeded.
	local := perfdata.Machine{OS: "darwin", Arch: "arm64", NumCPU: 8, CPUModel: "Apple M3", GoVersion: "go1.26.4"}
	seeded := perfdata.Machine{OS: "linux", Arch: "amd64", NumCPU: 16, CPUModel: amdM3, GoVersion: "go1.26.4"}
	localKey, seededKey := perfdata.MachineKey(local), perfdata.MachineKey(seeded)
	prior := Baseline{
		Version: schemaVersion,
		Machines: map[string]MachineBaseline{
			localKey: {
				CapturedAt: "20260701T010000Z", CapturedAtSHA: "1111111111ff", Machine: local,
				Anchor:     AnchorRecord{Name: anchorName, Package: anchorPackage, NSPerOp: 1.5, Iterations: 1000000000},
				Benchmarks: map[string]BenchmarkEntry{"test.BenchmarkA": {NSPerOp: 42, RatioToAnchor: 28}},
			},
			seededKey: {
				CapturedAt: "20260701T010000Z", CapturedAtSHA: "2222222222ff", Machine: seeded,
				Anchor:     AnchorRecord{Name: anchorName, Package: anchorPackage, NSPerOp: 1.5, Iterations: 1000000000},
				Benchmarks: map[string]BenchmarkEntry{"test.BenchmarkA": {NSPerOp: 999, RatioToAnchor: 666}},
			},
		},
	}
	data, err := json.Marshal(prior)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(baselinePath, data, 0o644); err != nil {
		t.Fatal(err)
	}

	got, out := runSeed(t, timeline, baselinePath, defaultSeedOptions())

	if len(got.Machines) != 2 {
		t.Fatalf("machines = %d, want 2 (seeded amd64 + preserved arm64)", len(got.Machines))
	}
	if e := got.Machines[seededKey].Benchmarks["test.BenchmarkA"]; e.NSPerOp != 100 {
		t.Errorf("seeded amd64 ns_per_op = %v, want 100 — the stale 999 was carried over on a model-string match", e.NSPerOp)
	}
	if e := got.Machines[localKey].Benchmarks["test.BenchmarkA"]; e.NSPerOp != 42 {
		t.Errorf("preserved arm64 ns_per_op = %v, want 42 (the off-arch M3 profile must survive)", e.NSPerOp)
	}
	if strings.Count(out, "preserved:") != 1 {
		t.Errorf("preserved lines = %d, want exactly 1:\n%s", strings.Count(out, "preserved:"), out)
	}
}
