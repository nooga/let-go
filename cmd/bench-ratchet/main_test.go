package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/perfdata"
)

func TestAggregateFromFileRetainsSamples(t *testing.T) {
	path := filepath.Join(t.TempDir(), "run.jsonl")
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	enc := json.NewEncoder(f)
	// Reduction semantics: first rep per benchmark is warmup (discarded),
	// remainder medians. Anchor [10,20] -> 20; BenchmarkA ns [30,60,40,50]
	// -> median(60,40,50) = 50; bytes [100,200,120,160] -> 160;
	// allocs [1,3,5,7] -> 5. Raw samples all retained.
	records := []StreamRecord{
		{Package: anchorPackage, Name: anchorName, Iterations: 100, NSPerOp: 10, CapturedAt: "2026-06-01T00:00:00Z"},
		{Package: anchorPackage, Name: anchorName, Iterations: 90, NSPerOp: 20, CapturedAt: "2026-06-01T00:00:01Z"},
		{Package: "pkg", Name: "BenchmarkA", Iterations: 50, NSPerOp: 30, BytesPerOp: 100, AllocsPerOp: 1, CapturedAt: "2026-06-01T00:00:02Z"},
		{Package: "pkg", Name: "BenchmarkA", Iterations: 60, NSPerOp: 60, BytesPerOp: 200, AllocsPerOp: 3, CapturedAt: "2026-06-01T00:00:03Z"},
		{Package: "pkg", Name: "BenchmarkA", Iterations: 55, NSPerOp: 40, BytesPerOp: 120, AllocsPerOp: 5, CapturedAt: "2026-06-01T00:00:04Z"},
		{Package: "pkg", Name: "BenchmarkA", Iterations: 58, NSPerOp: 50, BytesPerOp: 160, AllocsPerOp: 7, CapturedAt: "2026-06-01T00:00:05Z"},
	}
	for _, rec := range records {
		if err := enc.Encode(rec); err != nil {
			t.Fatal(err)
		}
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}

	baseline, err := aggregateFromFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(baseline.Anchor.Samples) != 2 {
		t.Fatalf("anchor samples = %d, want 2", len(baseline.Anchor.Samples))
	}
	entry := baseline.Benchmarks["pkg.BenchmarkA"]
	if entry.NSPerOp != 50 {
		t.Fatalf("ns/op = %v, want 50 (median of post-warmup reps)", entry.NSPerOp)
	}
	if entry.BytesPerOp != 160 {
		t.Fatalf("bytes/op = %d, want 160", entry.BytesPerOp)
	}
	if entry.AllocsPerOp != 5 {
		t.Fatalf("allocs/op = %d, want 5", entry.AllocsPerOp)
	}
	if len(entry.Samples) != 4 {
		t.Fatalf("samples = %d, want 4 (raw reps all retained)", len(entry.Samples))
	}
	// Per-sample ratios are raw values against the REDUCED anchor (20).
	if entry.Samples[0].RatioToAnchor != 1.5 {
		t.Fatalf("first sample ratio = %v, want 1.5", entry.Samples[0].RatioToAnchor)
	}
	if entry.Samples[1].RatioToAnchor != 3 {
		t.Fatalf("second sample ratio = %v, want 3", entry.Samples[1].RatioToAnchor)
	}
}

func TestReduceSamplesWarmupAndMedian(t *testing.T) {
	cases := []struct {
		name string
		in   []float64
		want float64
	}{
		{"empty", nil, 0},
		{"single passes through", []float64{7}, 7},
		{"two drops warmup", []float64{100, 40}, 40},
		{"three takes plain median, cold rep rejected", []float64{100, 40, 60}, 60},
		{"odd median after warmup", []float64{100, 60, 40, 50}, 50},
		{"even median after warmup", []float64{100, 10, 20, 30, 40}, 25},
		{"warmup spike ignored", []float64{9999, 10, 11, 12}, 11},
	}
	for _, c := range cases {
		if got := reduceSamples(c.in); got != c.want {
			t.Errorf("%s: reduceSamples(%v) = %v, want %v", c.name, c.in, got, c.want)
		}
	}
}

func TestSlugifyIsFilesystemSafeAndStable(t *testing.T) {
	cases := map[string]string{
		"arm64/Apple M1": "arm64-apple-m1",
		"arm64/Apple M2": "arm64-apple-m2",
		"arm64/Apple M3": "arm64-apple-m3",
		"amd64/Intel(R) Xeon(R) Platinum 8275CL CPU @ 3.00GHz": "amd64-intel-r-xeon-r-platinum-8275cl-cpu-3-00ghz",
		"  weird//name  ": "weird-name",
	}
	for in, want := range cases {
		if got := slugify(in); got != want {
			t.Errorf("slugify(%q) = %q, want %q", in, got, want)
		}
	}
}

// The whole point of the machine key in the snapshot filename: M1 and M2 both
// report GOARCH arm64, so a uname-based key would collide. Slugifying the
// <arch>/<CPUModel> partition key keeps them distinct.
func TestMachineKeyDistinguishesSameArchDifferentCPU(t *testing.T) {
	m1 := perfdata.MachineKey(Machine{OS: "darwin", Arch: "arm64", CPUModel: "Apple M1"})
	m2 := perfdata.MachineKey(Machine{OS: "darwin", Arch: "arm64", CPUModel: "Apple M2"})
	s1, s2 := slugify(m1), slugify(m2)
	if s1 == s2 {
		t.Fatalf("M1 and M2 slugify to the same key %q — snapshots would collide", s1)
	}
	if s1 != "arm64-apple-m1" || s2 != "arm64-apple-m2" {
		t.Fatalf("unexpected slugs: M1=%q M2=%q", s1, s2)
	}
}

func TestWriteBaselineWritesAtomicallyReadableJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "baseline.json")
	mach := Machine{OS: "darwin", Arch: "arm64", CPUModel: "Apple M3"}
	key := perfdata.MachineKey(mach)
	baseline := Baseline{
		Version: schemaVersion,
		Machines: map[string]MachineBaseline{
			key: {
				CapturedAt:    "2026-06-01T00:00:00Z",
				CapturedAtSHA: "abc",
				Machine:       mach,
				Benchmarks: map[string]BenchmarkEntry{
					"pkg.BenchmarkA": {NSPerOp: 1, RatioToAnchor: 2},
				},
			},
		},
	}
	if err := writeBaseline(path, baseline); err != nil {
		t.Fatal(err)
	}
	read, err := readBaseline(path)
	if err != nil {
		t.Fatal(err)
	}
	if read.Machines[key].Benchmarks["pkg.BenchmarkA"].RatioToAnchor != 2 {
		t.Fatalf("ratio = %v, want 2", read.Machines[key].Benchmarks["pkg.BenchmarkA"].RatioToAnchor)
	}
	matches, err := filepath.Glob(filepath.Join(filepath.Dir(path), ".*.tmp-*"))
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != 0 {
		t.Fatalf("left temporary files: %v", matches)
	}
}

func TestEffectiveCountPrefersJobOverride(t *testing.T) {
	if got := (captureJob{}).effectiveCount(3); got != 3 {
		t.Fatalf("zero override should fall back to default: got %d, want 3", got)
	}
	if got := (captureJob{count: 1}).effectiveCount(3); got != 1 {
		t.Fatalf("job override should win: got %d, want 1", got)
	}
}

// The full + fast profiles run the Clojure test suite 4x per execution mode:
// the first rep is discarded as warmup by reduceSamples and the remaining 3
// median (the minimum for a meaningful median). The old count=1 pin assumed
// run-to-run variance was negligible; in practice single cold samples plus
// anchor drift manufactured phantom +38% regressions (2026-07-18), and the
// suite bench visibly climbs across in-process reps. The cheap,
// benchtime-bounded vm/ir jobs keep the CLI default (count=0 → -count flag).
//
// There are exactly three suite modes, distinguished by the LG_SUITE_IR env
// toggle crossed with the gogen_ir build tag:
//   - bytecode    : *ir-compile* off, untagged
//   - ir_bytecode : *ir-compile* on  (LG_SUITE_IR=1), untagged
//   - aot_native  : *ir-compile* on  (LG_SUITE_IR=1), -tags gogen_ir
func TestSuiteJobsPinCountToFour(t *testing.T) {
	hasEnv := func(env []string, want string) bool {
		for _, e := range env {
			if e == want {
				return true
			}
		}
		return false
	}
	for _, full := range []bool{true, false} {
		jobs, _, err := buildJobs("", "", full, false, nil)
		if err != nil {
			t.Fatalf("buildJobs(full=%v): %v", full, err)
		}
		suite := map[string]captureJob{}
		total := map[string]captureJob{}
		for _, j := range jobs {
			if j.pkg == suitePackage {
				if j.count != 4 {
					t.Errorf("full=%v: suite job [%s] count = %d, want 4", full, j.variant, j.count)
				}
				switch j.filter.String() {
				case suiteFilter:
					suite[j.variant] = j
				case suiteTotalFilter:
					total[j.variant] = j
				default:
					t.Errorf("full=%v: unexpected suite filter %q for variant %q", full, j.filter.String(), j.variant)
				}
			} else if j.count != 0 {
				t.Errorf("full=%v: non-suite job %s [%s] count = %d, want 0 (use CLI default)", full, j.pkg, j.variant, j.count)
			}
		}
		if len(suite) != 3 {
			t.Fatalf("full=%v: expected 3 suite variants (bytecode, ir_bytecode, aot_native), got %d: %v", full, len(suite), keysOf(suite))
		}
		if len(total) != 3 {
			t.Fatalf("full=%v: expected 3 total-suite variants (total_bytecode, total_ir_bytecode, total_aot_native), got %d: %v", full, len(total), keysOf(total))
		}
		// bytecode: no IR toggle, untagged.
		if j := suite["bytecode"]; len(j.env) != 0 || j.tags != "" {
			t.Errorf("full=%v: bytecode variant want no env/tags, got env=%v tags=%q", full, j.env, j.tags)
		}
		// ir_bytecode: IR on, still untagged (passes run as bytecode).
		if j := suite["ir_bytecode"]; !hasEnv(j.env, "LG_SUITE_IR=1") || j.tags != "" {
			t.Errorf("full=%v: ir_bytecode want LG_SUITE_IR=1 + untagged, got env=%v tags=%q", full, j.env, j.tags)
		}
		// aot_native: IR on AND gogen_ir tag (passes dispatch to native Go).
		if j := suite["aot_native"]; !hasEnv(j.env, "LG_SUITE_IR=1") || j.tags != "gogen_ir" {
			t.Errorf("full=%v: aot_native want LG_SUITE_IR=1 + -tags gogen_ir, got env=%v tags=%q", full, j.env, j.tags)
		}
		if j := total["total_bytecode"]; len(j.env) != 0 || j.tags != "" {
			t.Errorf("full=%v: total_bytecode want no env/tags, got env=%v tags=%q", full, j.env, j.tags)
		}
		if j := total["total_ir_bytecode"]; !hasEnv(j.env, "LG_SUITE_IR=1") || j.tags != "" {
			t.Errorf("full=%v: total_ir_bytecode want LG_SUITE_IR=1 + untagged, got env=%v tags=%q", full, j.env, j.tags)
		}
		if j := total["total_aot_native"]; !hasEnv(j.env, "LG_SUITE_IR=1") || j.tags != "gogen_ir" {
			t.Errorf("full=%v: total_aot_native want LG_SUITE_IR=1 + -tags gogen_ir, got env=%v tags=%q", full, j.env, j.tags)
		}
	}
}

// The pr-fast profile includes the anchor and the kept families, excludes the
// sub-nanosecond noise families, and sets its own tuning defaults.
func TestProfilePrFast(t *testing.T) {
	p, ok := profiles["pr-fast"]
	if !ok {
		t.Fatal("pr-fast profile missing")
	}
	if p.count == 0 || p.benchtime == "" || p.budget == 0 {
		t.Errorf("pr-fast should set count/benchtime/budget defaults, got %+v", p)
	}
	jobs, err := p.jobs("gogen_ir")
	if err != nil {
		t.Fatal(err)
	}
	if len(jobs) != 1 || jobs[0].pkg != anchorPackage {
		t.Fatalf("want one pkg/vm job, got %+v", jobs)
	}
	f := jobs[0].filter
	for _, keep := range []string{anchorName, "BenchmarkFrameDispatch", "BenchmarkMapAssoc"} {
		if !f.MatchString(keep) {
			t.Errorf("pr-fast should include %s", keep)
		}
	}
	for _, drop := range []string{"BenchmarkStackOps", "BenchmarkIsTruthy", "BenchmarkConsCreation"} {
		if f.MatchString(drop) {
			t.Errorf("pr-fast should exclude %s", drop)
		}
	}
}

func keysOf(m map[string]captureJob) []string {
	ks := make([]string, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	return ks
}

// TestJobsSelectScope pins the check-report scope predicate: baseline entries
// the active profile never selects are "out of scope" (skipped, counted once),
// while in-scope-but-absent entries stay MISSING. Guards the fast-gate report
// from drowning in MISSING rows for the full-profile pkg/vm fleet.
func TestJobsSelectScope(t *testing.T) {
	anchorRE := regexp.MustCompile(`^BenchmarkRatchetAnchor$`)
	suiteRE := regexp.MustCompile(`^BenchmarkClojureTestSuite$`)
	irRE := regexp.MustCompile(`^BenchmarkIRCompile$`)
	jobs := []captureJob{
		{pkg: "github.com/nooga/let-go/pkg/vm", filter: anchorRE},
		{pkg: "github.com/nooga/let-go/test", filter: suiteRE, variant: "bytecode"},
		{pkg: "github.com/nooga/let-go/pkg/ir", filter: irRE, variant: "gogen_ir"},
	}
	cases := []struct {
		name string
		want bool
	}{
		// anchor: variant-free job matches variant-free entry
		{"github.com/nooga/let-go/pkg/vm.BenchmarkRatchetAnchor", true},
		// full-fleet vm micro: same pkg, family not in any filter → out of scope
		{"github.com/nooga/let-go/pkg/vm.BenchmarkVectorCreation/PersistentVector", false},
		// suite under the selected variant, in scope
		{"github.com/nooga/let-go/test.BenchmarkClojureTestSuite [bytecode]", true},
		// suite under a variant this run doesn't measure → out of scope
		{"github.com/nooga/let-go/test.BenchmarkClojureTestSuite [aot_native]", false},
		// variant job must not claim the variant-free spelling
		{"github.com/nooga/let-go/pkg/ir.BenchmarkIRCompile", false},
		{"github.com/nooga/let-go/pkg/ir.BenchmarkIRCompile [gogen_ir]", true},
		// sub-benchmark matches on the family segment
		{"github.com/nooga/let-go/pkg/ir.BenchmarkIRCompile/warm [gogen_ir]", true},
		// unrelated package
		{"github.com/nooga/let-go/pkg/compiler.BenchmarkInitFromLGB [bytecode]", false},
	}
	for _, c := range cases {
		if got := jobsSelect(jobs, c.name); got != c.want {
			t.Errorf("jobsSelect(%q) = %v, want %v", c.name, got, c.want)
		}
	}
}

func TestSeedBaselineAmd64OnlyPreservesM3(t *testing.T) {
	// Create a temporary directory with mock timeline snapshots.
	tmpDir := t.TempDir()
	timelineDir := filepath.Join(tmpDir, "timeline")
	if err := os.MkdirAll(timelineDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Create mock baselines for amd64 and arm64 machines.
	createMockSnapshot := func(filename string, sha, timestamp, arch, machine string) {
		baseline := Baseline{
			Version: schemaVersion,
			Machines: map[string]MachineBaseline{
				perfdata.MachineKey(Machine{
					OS:        "linux",
					Arch:      arch,
					NumCPU:    16,
					CPUModel:  machine,
					GoVersion: "go1.26.4",
				}): {
					CapturedAt:    timestamp,
					CapturedAtSHA: sha,
					Machine: Machine{
						OS:        "linux",
						Arch:      arch,
						NumCPU:    16,
						CPUModel:  machine,
						GoVersion: "go1.26.4",
					},
					Anchor: AnchorRecord{
						Name:       anchorName,
						Package:    anchorPackage,
						NSPerOp:    1.5,
						Iterations: 1000000000,
						Samples: []BenchmarkSample{
							{Iterations: 1000000000, NSPerOp: 1.5, CapturedAt: timestamp},
						},
					},
					Benchmarks: map[string]BenchmarkEntry{
						"test.BenchmarkA": {
							NSPerOp:       100,
							RatioToAnchor: 66.67,
							Samples: []BenchmarkSample{
								{NSPerOp: 100, CapturedAt: timestamp},
							},
						},
						// Unstable benchmark that should be filtered out
						"github.com/nooga/let-go/test.BenchmarkClojureTestSuite": {
							NSPerOp:       5000,
							RatioToAnchor: 3333.0,
							Samples: []BenchmarkSample{
								{NSPerOp: 5000, CapturedAt: timestamp},
							},
						},
					},
				},
			},
		}
		data, err := json.Marshal(baseline)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(timelineDir, filename), data, 0o644); err != nil {
			t.Fatal(err)
		}
	}

	// Create amd64 snapshots (should be seeded): a full minimum window, since
	// a machine with fewer than -seed-min-window snapshots is skipped.
	createMockSnapshot("20260801T010134Z-b170a08eef47-amd64-amd-epyc-7763.json",
		"b170a08eef47", "2026-08-01T01:01:34Z", "amd64", "AMD EPYC 7763")
	createMockSnapshot("20260731T010134Z-a170a08eef47-amd64-amd-epyc-7763.json",
		"a170a08eef47", "2026-07-31T01:01:34Z", "amd64", "AMD EPYC 7763")
	createMockSnapshot("20260730T010134Z-9170a08eef47-amd64-amd-epyc-7763.json",
		"9170a08eef47", "2026-07-30T01:01:34Z", "amd64", "AMD EPYC 7763")

	// Create arm64 snapshots (should be ignored per #651)
	createMockSnapshot("20260801T010134Z-b170a08eef47-arm64-apple-m1-virtual.json",
		"b170a08eef47", "2026-08-01T01:01:34Z", "arm64", "Apple M1 (Virtual)")

	// Create existing baseline with M3 profile to be preserved
	existingBaseline := Baseline{
		Version: schemaVersion,
		Machines: map[string]MachineBaseline{
			perfdata.MachineKey(Machine{
				OS:        "darwin",
				Arch:      "arm64",
				NumCPU:    8,
				CPUModel:  "Apple M3",
				GoVersion: "go1.26.4",
			}): {
				CapturedAt:    "2026-06-01T00:00:00Z",
				CapturedAtSHA: "oldsha123",
				Machine: Machine{
					OS:        "darwin",
					Arch:      "arm64",
					NumCPU:    8,
					CPUModel:  "Apple M3",
					GoVersion: "go1.26.4",
				},
				Anchor: AnchorRecord{
					Name:       anchorName,
					Package:    anchorPackage,
					NSPerOp:    1.2,
					Iterations: 1000000000,
					Samples: []BenchmarkSample{
						{Iterations: 1000000000, NSPerOp: 1.2, CapturedAt: "2026-06-01T00:00:00Z"},
					},
				},
				Benchmarks: map[string]BenchmarkEntry{
					"test.BenchmarkA": {
						NSPerOp:       90,
						RatioToAnchor: 75.0,
						Samples: []BenchmarkSample{
							{NSPerOp: 90, CapturedAt: "2026-06-01T00:00:00Z"},
						},
					},
				},
			},
		},
	}
	existingData, err := json.Marshal(existingBaseline)
	if err != nil {
		t.Fatal(err)
	}
	baselineFile := filepath.Join(tmpDir, "baseline.json")
	if err := os.WriteFile(baselineFile, existingData, 0o644); err != nil {
		t.Fatal(err)
	}

	// Run seed-baseline.
	seedBaseline(baselineFile, timelineDir, defaultSeedOptions())

	// Verify the output.
	data, err := os.ReadFile(baselineFile)
	if err != nil {
		t.Fatal(err)
	}
	var result Baseline
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatal(err)
	}

	// Verify version and structure.
	if result.Version != schemaVersion {
		t.Errorf("version = %d, want %d", result.Version, schemaVersion)
	}

	// Should have 2 machines: amd64 (from perf-data) + M3 (preserved)
	if len(result.Machines) != 2 {
		t.Errorf("machines = %d, want 2 (amd64 + M3)", len(result.Machines))
	}

	// Verify amd64 machine is present and unstable benchmark is filtered
	hasAmd64 := false
	for key, mb := range result.Machines {
		if strings.Contains(key, "amd64") || strings.Contains(mb.Machine.CPUModel, "EPYC") {
			hasAmd64 = true
			// Verify unstable benchmark is filtered out
			if _, ok := mb.Benchmarks["github.com/nooga/let-go/test.BenchmarkClojureTestSuite"]; ok {
				t.Error("unstable BenchmarkClojureTestSuite should be filtered out")
			}
			// Verify stable benchmark is kept
			e, ok := mb.Benchmarks["test.BenchmarkA"]
			if !ok {
				t.Error("stable benchmark test.BenchmarkA should be kept")
			}
			// A seeded deterministic floor with no commit attached is a floor
			// the gate cannot compare like with like against.
			if e.BestSinceSHA == "" || e.BestSinceAt == "" {
				t.Errorf("seeded entry has no provenance: best_since_sha=%q best_since_at=%q",
					e.BestSinceSHA, e.BestSinceAt)
			}
			if e.BestSinceSHA != mb.CapturedAtSHA {
				t.Errorf("seeded provenance %q != profile captured_at_sha %q", e.BestSinceSHA, mb.CapturedAtSHA)
			}
			// The deterministic pair is measured by the seed window too, so it
			// carries its own date; without one the gate would rank the seeded
			// row below every stamped row.
			if e.AllocsSinceSHA != mb.CapturedAtSHA || e.AllocsSinceAt != mb.CapturedAt ||
				e.BytesSinceSHA != mb.CapturedAtSHA || e.BytesSinceAt != mb.CapturedAt {
				t.Errorf("seeded deterministic provenance = %q/%q, %q/%q, want %q/%q for both",
					e.AllocsSinceSHA, e.AllocsSinceAt, e.BytesSinceSHA, e.BytesSinceAt,
					mb.CapturedAtSHA, mb.CapturedAt)
			}
		}
	}
	if !hasAmd64 {
		t.Error("amd64 machine not found in merged baseline")
	}

	// Verify M3 profile is preserved
	hasM3 := false
	for key, mb := range result.Machines {
		if strings.Contains(key, "apple-m3") || strings.Contains(mb.Machine.CPUModel, "M3") {
			hasM3 = true
			// M3 should retain its old data
			if mb.CapturedAtSHA != "oldsha123" {
				t.Errorf("M3 captured_at_sha changed; want oldsha123, got %q", mb.CapturedAtSHA)
			}
		}
	}
	if !hasM3 {
		t.Error("M3 profile should be preserved")
	}
}

// TestForceRebaselineWritesOnlyItsOwnProfile pins that a forced acceptance is
// recorded where it was measured and nowhere else. The accepted numbers reach
// the gate by being the newest measurement of them, so copying them into other
// profiles would only record one machine's capture as those tiers' stored floor
// under a commit they never ran.
func TestForceRebaselineWritesOnlyItsOwnProfile(t *testing.T) {
	const benchmark = "pkg.BenchmarkAcceptedRegression"
	currentMachine := Machine{OS: "darwin", Arch: "arm64", CPUModel: "Apple M3", GoVersion: "go1.26.5"}
	otherMachine := Machine{OS: "linux", Arch: "amd64", CPUModel: "EPYC", GoVersion: "go1.26.5"}
	currentKey := perfdata.MachineKey(currentMachine)
	otherKey := perfdata.MachineKey(otherMachine)

	baseline := Baseline{Version: schemaVersion, Machines: map[string]MachineBaseline{
		currentKey: {
			CapturedAt: "2026-08-01T00:00:00Z", CapturedAtSHA: "oldsha",
			Machine:    currentMachine,
			Benchmarks: map[string]BenchmarkEntry{benchmark: {NSPerOp: 10, AllocsPerOp: 10, BytesPerOp: 100}},
		},
		otherKey: {
			CapturedAt: "2026-08-02T00:00:00Z", CapturedAtSHA: "othersha",
			Machine: otherMachine,
			Benchmarks: map[string]BenchmarkEntry{
				benchmark: {NSPerOp: 20, RatioToAnchor: 20, AllocsPerOp: 5, BytesPerOp: 50},
			},
		},
	}}
	current := MachineBaseline{
		CapturedAt:    "2026-08-18T21:00:00Z",
		CapturedAtSHA: "accepted-sha",
		Machine:       currentMachine,
		Benchmarks: map[string]BenchmarkEntry{
			benchmark: {NSPerOp: 15, RatioToAnchor: 15, AllocsPerOp: 12, BytesPerOp: 120},
		},
	}

	forceRebaseline(&baseline, currentKey, current, true)

	// Profile B is untouched — numbers and provenance both.
	gotOther := baseline.Machines[otherKey].Benchmarks[benchmark]
	if gotOther.AllocsPerOp != 5 || gotOther.BytesPerOp != 50 {
		t.Fatalf("other profile deterministic metrics were rewritten: %+v", gotOther)
	}
	if gotOther.BestSinceSHA != "" || gotOther.BestSinceAt != "" {
		t.Fatalf("other profile provenance was stamped: %+v", gotOther)
	}
	if gotOther.NSPerOp != 20 || gotOther.RatioToAnchor != 20 {
		t.Fatalf("other profile timing changed: %+v", gotOther)
	}

	// And the accepted row is nonetheless what the gate compares against, on B
	// as much as on A, because it is the newest measurement of this benchmark.
	bar := machineIndependentBar(baseline)[benchmark]
	if bar.AllocsPerOp != 12 || bar.BytesPerOp != 120 {
		t.Fatalf("deterministic bar = %d allocs/%d bytes, want the accepted 12/120", bar.AllocsPerOp, bar.BytesPerOp)
	}
	if bar.AllocsSinceSHA != "accepted-sha" {
		t.Fatalf("bar provenance = %q, want accepted-sha", bar.AllocsSinceSHA)
	}
	onOther := MachineBaseline{
		Machine:    otherMachine,
		Benchmarks: map[string]BenchmarkEntry{benchmark: {AllocsPerOp: 12, BytesPerOp: 120}},
	}
	if n := compareDeterministic(bar2map(bar, benchmark), onOther, allocBudget); n != 0 {
		t.Fatalf("check on the other machine reported %d regression(s) against the accepted row; want 0", n)
	}
}

func bar2map(e BenchmarkEntry, name string) map[string]BenchmarkEntry {
	return map[string]BenchmarkEntry{name: e}
}

func TestForceRebaselineRetainsUnmeasuredEntries(t *testing.T) {
	const benchmark = "pkg.BenchmarkAcceptedRegression"
	const outsideScope = "pkg.BenchmarkOutsideScope"
	currentMachine := Machine{OS: "darwin", Arch: "arm64", CPUModel: "Apple M3", GoVersion: "go1.26.5"}
	otherMachine := Machine{OS: "linux", Arch: "amd64", CPUModel: "EPYC", GoVersion: "go1.26.5"}
	currentKey := perfdata.MachineKey(currentMachine)
	otherKey := perfdata.MachineKey(otherMachine)

	baseline := Baseline{Version: schemaVersion, Machines: map[string]MachineBaseline{
		currentKey: {
			Machine: currentMachine,
			Benchmarks: map[string]BenchmarkEntry{
				benchmark:    {NSPerOp: 10, AllocsPerOp: 10, BytesPerOp: 100},
				outsideScope: {NSPerOp: 11, AllocsPerOp: 4, BytesPerOp: 40},
			},
		},
		otherKey: {
			Machine: otherMachine,
			Benchmarks: map[string]BenchmarkEntry{
				benchmark:    {NSPerOp: 20, RatioToAnchor: 20, AllocsPerOp: 5, BytesPerOp: 50},
				outsideScope: {NSPerOp: 30, AllocsPerOp: 3, BytesPerOp: 30},
			},
		},
	}}
	current := MachineBaseline{
		CapturedAt:    "2026-08-18T21:00:00Z",
		CapturedAtSHA: "accepted-sha",
		Machine:       currentMachine,
		Benchmarks: map[string]BenchmarkEntry{
			benchmark: {NSPerOp: 15, RatioToAnchor: 15, AllocsPerOp: 12, BytesPerOp: 120},
		},
	}

	forceRebaseline(&baseline, currentKey, current, true)

	gotCurrent := baseline.Machines[currentKey].Benchmarks[benchmark]
	if gotCurrent.NSPerOp != 15 || gotCurrent.AllocsPerOp != 12 || gotCurrent.BytesPerOp != 120 {
		t.Fatalf("current profile was not replaced: %+v", gotCurrent)
	}
	// The benchmark this fast-gate run did not measure keeps its bar rather than
	// being dropped, so a rebaseline cannot erase full-profile history.
	if got := baseline.Machines[currentKey].Benchmarks[outsideScope]; got.NSPerOp != 11 || got.AllocsPerOp != 4 || got.BytesPerOp != 40 {
		t.Fatalf("unmeasured benchmark was not retained: %+v", got)
	}
	if got := baseline.Machines[otherKey].Benchmarks[outsideScope]; got.AllocsPerOp != 3 || got.BytesPerOp != 30 {
		t.Fatalf("other out-of-scope benchmark changed: %+v", got)
	}
}

// TestMachineIndependentBarUsesNewestMeasuredRow pins the rule that allocs/bytes
// are a property of the CODE at a commit, not of a machine: the deterministic
// reference must be the value measured at the newest commit that has the
// benchmark, never a minimum mixed across rows captured at different commits.
//
// The fixture is the shape that makes the difference visible — a tier whose
// runner stopped reporting, still carrying the lower numbers it measured before
// the code changed underneath it.
func TestMachineIndependentBarUsesNewestMeasuredRow(t *testing.T) {
	const benchmark = "pkg.BenchmarkInitFromLGB [bytecode]"
	baseline := Baseline{Version: schemaVersion, Machines: map[string]MachineBaseline{
		// Stale tier: low numbers, but measured at a commit long superseded.
		"amd64/Stale Xeon": {
			CapturedAt:    "2026-08-11T16:10:14Z",
			CapturedAtSHA: "f3ca5f9be1da",
			Machine:       Machine{OS: "linux", Arch: "amd64", CPUModel: "Stale Xeon"},
			Benchmarks: map[string]BenchmarkEntry{
				benchmark: {AllocsPerOp: 7041, BytesPerOp: 743356,
					AllocsSinceSHA: "f3ca5f9be1da", AllocsSinceAt: "2026-08-11T16:10:14Z",
					BytesSinceSHA: "f3ca5f9be1da", BytesSinceAt: "2026-08-11T16:10:14Z"},
			},
		},
		// Current tier: the newest row that carries this benchmark.
		"amd64/Fresh EPYC": {
			CapturedAt:    "2026-09-07T05:34:02Z",
			CapturedAtSHA: "477a5d36e25f",
			Machine:       Machine{OS: "linux", Arch: "amd64", CPUModel: "Fresh EPYC"},
			Benchmarks: map[string]BenchmarkEntry{
				benchmark: {AllocsPerOp: 13995, BytesPerOp: 1016128,
					AllocsSinceSHA: "477a5d36e25f", AllocsSinceAt: "2026-09-07T05:34:02Z",
					BytesSinceSHA: "477a5d36e25f", BytesSinceAt: "2026-09-07T05:34:02Z"},
			},
		},
		// A third tier, newer than the stale one but older than the newest.
		"arm64/Apple M3": {
			CapturedAt:    "2026-08-24T20:18:19Z",
			CapturedAtSHA: "dadb8e0b54b7",
			Machine:       Machine{OS: "darwin", Arch: "arm64", CPUModel: "Apple M3"},
			Benchmarks: map[string]BenchmarkEntry{
				benchmark: {AllocsPerOp: 26294, BytesPerOp: 1599971,
					BestSinceSHA: "dadb8e0b54b7", BestSinceAt: "2026-08-24T20:18:19Z",
					AllocsSinceSHA: "dadb8e0b54b7", AllocsSinceAt: "2026-08-24T20:18:19Z",
					BytesSinceSHA: "dadb8e0b54b7", BytesSinceAt: "2026-08-24T20:18:19Z"},
			},
		},
	}}

	bar := machineIndependentBar(baseline)[benchmark]
	if bar.AllocsPerOp != 13995 || bar.BytesPerOp != 1016128 {
		t.Fatalf("deterministic bar = %d allocs/%d bytes, want 13995/1016128 (the newest measured row)",
			bar.AllocsPerOp, bar.BytesPerOp)
	}
	if bar.AllocsSinceSHA != "477a5d36e25f" {
		t.Fatalf("bar provenance = %q, want 477a5d36e25f", bar.AllocsSinceSHA)
	}

	// The gate must now pass for code that is better than the newest reference,
	// even though it is far above the abandoned row.
	current := MachineBaseline{
		Machine:    Machine{OS: "darwin", Arch: "arm64", CPUModel: "Apple M3"},
		Benchmarks: map[string]BenchmarkEntry{benchmark: {AllocsPerOp: 11900, BytesPerOp: 900000}},
	}
	if n := compareDeterministic(machineIndependentBar(baseline), current, allocBudget); n != 0 {
		t.Fatalf("compareDeterministic reported %d regression(s) against the newest row; want 0", n)
	}

	// A genuine regression against the newest row must still be reported.
	worse := MachineBaseline{
		Machine:    current.Machine,
		Benchmarks: map[string]BenchmarkEntry{benchmark: {AllocsPerOp: 20000, BytesPerOp: 1016128}},
	}
	if n := compareDeterministic(machineIndependentBar(baseline), worse, allocBudget); n == 0 {
		t.Fatal("compareDeterministic reported no regression for 13995 -> 20000 allocs/op")
	}
}

// TestMachineIndependentBarTieBreaksOnMin covers rows that share the newest
// provenance: with nothing to distinguish them by commit, the tightest wins, so
// the bar keeps ratcheting within one commit's worth of measurements.
func TestMachineIndependentBarTieBreaksOnMin(t *testing.T) {
	const benchmark = "pkg.BenchmarkTied"
	baseline := Baseline{Version: schemaVersion, Machines: map[string]MachineBaseline{
		"a": {
			CapturedAt: "2026-09-07T05:34:02Z", CapturedAtSHA: "same-sha",
			Benchmarks: map[string]BenchmarkEntry{benchmark: {AllocsPerOp: 100, BytesPerOp: 900}},
		},
		"b": {
			CapturedAt: "2026-09-07T05:34:02Z", CapturedAtSHA: "same-sha",
			Benchmarks: map[string]BenchmarkEntry{benchmark: {AllocsPerOp: 90, BytesPerOp: 1000}},
		},
	}}
	bar := machineIndependentBar(baseline)[benchmark]
	if bar.AllocsPerOp != 90 || bar.BytesPerOp != 900 {
		t.Fatalf("tied bar = %d allocs/%d bytes, want 90/900", bar.AllocsPerOp, bar.BytesPerOp)
	}
}

// TestMachineIndependentBarPrefersEntryProvenance checks that a per-entry
// best_since stamp outranks the profile's captured_at: within a profile,
// ratchetMerge keeps entries set at older commits alongside freshly measured
// ones, so the entry's own stamp is the more precise provenance.
func TestMachineIndependentBarPrefersEntryProvenance(t *testing.T) {
	const benchmark = "pkg.BenchmarkEntryStamp"
	baseline := Baseline{Version: schemaVersion, Machines: map[string]MachineBaseline{
		// Profile captured recently, but this entry's bar was set long ago.
		"a": {
			CapturedAt: "2026-09-07T05:34:02Z", CapturedAtSHA: "newprofile",
			Benchmarks: map[string]BenchmarkEntry{benchmark: {
				AllocsPerOp: 10, BytesPerOp: 100,
				BestSinceSHA: "oldsha", BestSinceAt: "2026-01-01T00:00:00Z",
				AllocsSinceSHA: "oldsha", AllocsSinceAt: "2026-01-01T00:00:00Z",
				BytesSinceSHA: "oldsha", BytesSinceAt: "2026-01-01T00:00:00Z",
			}},
		},
		"b": {
			CapturedAt: "2026-06-01T00:00:00Z", CapturedAtSHA: "midprofile",
			Benchmarks: map[string]BenchmarkEntry{benchmark: {
				AllocsPerOp: 50, BytesPerOp: 500,
				BestSinceSHA: "midsha", BestSinceAt: "2026-06-01T00:00:00Z",
				AllocsSinceSHA: "midsha", AllocsSinceAt: "2026-06-01T00:00:00Z",
				BytesSinceSHA: "midsha", BytesSinceAt: "2026-06-01T00:00:00Z",
			}},
		},
	}}
	bar := machineIndependentBar(baseline)[benchmark]
	if bar.AllocsPerOp != 50 || bar.AllocsSinceSHA != "midsha" {
		t.Fatalf("bar = %d allocs (%q), want 50 (midsha)", bar.AllocsPerOp, bar.AllocsSinceSHA)
	}
}

func TestFilterUnstableBenchmarksKeepsStableGateRows(t *testing.T) {
	stable := "github.com/nooga/let-go/pkg/ir.BenchmarkIRCompile [bytecode]"
	unstable := []string{
		"github.com/nooga/let-go/test.BenchmarkClojureTestSuite [bytecode]",
		"github.com/nooga/let-go/test.BenchmarkClojureTestSuite [ir_bytecode]",
		"github.com/nooga/let-go/test.BenchmarkClojureTestSuite [aot_native]",
		"github.com/nooga/let-go/test.BenchmarkClojureTestSuiteCompileAndRun [total_bytecode]",
		"github.com/nooga/let-go/test.BenchmarkClojureTestSuiteCompileAndRun [total_ir_bytecode]",
		"github.com/nooga/let-go/test.BenchmarkClojureTestSuiteCompileAndRun [total_aot_native]",
	}
	benchmarks := map[string]BenchmarkEntry{stable: {NSPerOp: 1}}
	for _, name := range unstable {
		benchmarks[name] = BenchmarkEntry{NSPerOp: 2}
	}
	filtered := filterUnstableBenchmarks(MachineBaseline{Benchmarks: benchmarks})
	if len(filtered.Benchmarks) != 1 {
		t.Fatalf("filtered benchmark count = %d, want 1", len(filtered.Benchmarks))
	}
	if _, ok := filtered.Benchmarks[stable]; !ok {
		t.Fatal("stable gate benchmark was filtered")
	}
}

// TestDeterministicProvenanceSurvivesTimingOnlyImprovement pins the rule that a
// row's deterministic stamp names the run that measured its allocs/bytes, not
// the run that last moved any metric. A timing win leaves allocs/bytes where
// they were, so it must leave their provenance where it was too — otherwise a
// pinned old floor would outrank a row that genuinely measured the benchmark
// later, and the gate would attribute the old numbers to the newer commit.
func TestDeterministicProvenanceSurvivesTimingOnlyImprovement(t *testing.T) {
	const benchmark = "pkg.BenchmarkPinned"

	january := MachineBaseline{
		CapturedAt: "2026-01-10T00:00:00Z", CapturedAtSHA: "jansha",
		Benchmarks: map[string]BenchmarkEntry{benchmark: {
			NSPerOp: 100, RatioToAnchor: 10, AllocsPerOp: 10, BytesPerOp: 100,
			BestSinceSHA: "jansha", BestSinceAt: "2026-01-10T00:00:00Z",
			AllocsSinceSHA: "jansha", AllocsSinceAt: "2026-01-10T00:00:00Z",
			BytesSinceSHA: "jansha", BytesSinceAt: "2026-01-10T00:00:00Z",
		}},
	}
	// March: faster, identical allocs/bytes.
	march := MachineBaseline{
		CapturedAt: "2026-03-10T00:00:00Z", CapturedAtSHA: "marsha",
		Benchmarks: map[string]BenchmarkEntry{benchmark: {
			NSPerOp: 50, RatioToAnchor: 5, AllocsPerOp: 10, BytesPerOp: 100,
		}},
	}
	merged, _ := ratchetMerge(january, march)
	got := merged.Benchmarks[benchmark]
	if got.BestSinceSHA != "marsha" {
		t.Fatalf("timing provenance = %q, want marsha", got.BestSinceSHA)
	}
	if got.AllocsSinceSHA != "jansha" || got.BytesSinceSHA != "jansha" ||
		got.AllocsSinceAt != "2026-01-10T00:00:00Z" || got.BytesSinceAt != "2026-01-10T00:00:00Z" {
		t.Fatalf("deterministic provenance = %q/%q, want jansha/2026-01-10T00:00:00Z for both",
			got.AllocsSinceSHA, got.BytesSinceSHA)
	}

	// February measured this benchmark for real; it is the newest row that did.
	baseline := Baseline{Version: schemaVersion, Machines: map[string]MachineBaseline{
		"a": merged,
		"b": {
			CapturedAt: "2026-02-10T00:00:00Z", CapturedAtSHA: "febsha",
			Benchmarks: map[string]BenchmarkEntry{benchmark: {
				AllocsPerOp: 15, BytesPerOp: 150,
				BestSinceSHA: "febsha", BestSinceAt: "2026-02-10T00:00:00Z",
				AllocsSinceSHA: "febsha", AllocsSinceAt: "2026-02-10T00:00:00Z",
				BytesSinceSHA: "febsha", BytesSinceAt: "2026-02-10T00:00:00Z",
			}},
		},
	}}
	bar := machineIndependentBar(baseline)[benchmark]
	if bar.AllocsPerOp != 15 || bar.BytesPerOp != 150 || bar.AllocsSinceSHA != "febsha" {
		t.Fatalf("bar = %d allocs/%d bytes (%q), want 15/150 (febsha)",
			bar.AllocsPerOp, bar.BytesPerOp, bar.AllocsSinceSHA)
	}
}

// TestDeterministicProvenanceMovesWhenAllocsChange is the other half: a run that
// actually lowers allocs/bytes re-dates them.
func TestDeterministicProvenanceMovesWhenAllocsChange(t *testing.T) {
	const benchmark = "pkg.BenchmarkTightened"
	january := MachineBaseline{
		CapturedAt: "2026-01-10T00:00:00Z", CapturedAtSHA: "jansha",
		Benchmarks: map[string]BenchmarkEntry{benchmark: {
			NSPerOp: 100, AllocsPerOp: 10, BytesPerOp: 100,
			AllocsSinceSHA: "jansha", AllocsSinceAt: "2026-01-10T00:00:00Z",
			BytesSinceSHA: "jansha", BytesSinceAt: "2026-01-10T00:00:00Z",
		}},
	}
	march := MachineBaseline{
		CapturedAt: "2026-03-10T00:00:00Z", CapturedAtSHA: "marsha",
		Benchmarks: map[string]BenchmarkEntry{benchmark: {
			NSPerOp: 100, AllocsPerOp: 8, BytesPerOp: 100,
		}},
	}
	merged, _ := ratchetMerge(january, march)
	got := merged.Benchmarks[benchmark]
	if got.AllocsPerOp != 8 {
		t.Fatalf("merged allocs = %d, want 8", got.AllocsPerOp)
	}
	if got.AllocsSinceSHA != "marsha" || got.AllocsSinceAt != "2026-03-10T00:00:00Z" {
		t.Fatalf("allocs provenance = %q/%q, want marsha/2026-03-10T00:00:00Z",
			got.AllocsSinceSHA, got.AllocsSinceAt)
	}
}

// TestDeterministicProvenanceUnknownSortsOldest pins that a RATCHETED row
// predating the deterministic stamp claims no measurement date: best_since says
// only that some metric moved at that commit, which may have been the timing,
// so the row may be the reference only when nothing better exists.
func TestDeterministicProvenanceUnknownSortsOldest(t *testing.T) {
	const benchmark = "pkg.BenchmarkUnstamped"
	baseline := Baseline{Version: schemaVersion, Machines: map[string]MachineBaseline{
		"a": {
			CapturedAt: "2026-09-07T05:34:02Z", CapturedAtSHA: "newprofile",
			Benchmarks: map[string]BenchmarkEntry{benchmark: {
				AllocsPerOp: 10, BytesPerOp: 100,
				BestSinceSHA: "newprofile", BestSinceAt: "2026-09-07T05:34:02Z",
			}},
		},
		"b": {
			CapturedAt: "2026-06-01T00:00:00Z", CapturedAtSHA: "midprofile",
			Benchmarks: map[string]BenchmarkEntry{benchmark: {
				AllocsPerOp: 50, BytesPerOp: 500,
				AllocsSinceSHA: "midsha", AllocsSinceAt: "2026-06-01T00:00:00Z",
				BytesSinceSHA: "midsha", BytesSinceAt: "2026-06-01T00:00:00Z",
			}},
		},
	}}
	bar := machineIndependentBar(baseline)[benchmark]
	if bar.AllocsPerOp != 50 || bar.AllocsSinceSHA != "midsha" {
		t.Fatalf("bar = %d allocs (%q), want 50 (midsha)", bar.AllocsPerOp, bar.AllocsSinceSHA)
	}
}

// TestMachineIndependentBarTiesOnSHANotTimestamp pins that "same commit" is
// decided by the commit, not the clock: two captures of one code state tie
// however far apart they ran, and the tightest of them is the floor.
func TestMachineIndependentBarTiesOnSHANotTimestamp(t *testing.T) {
	const benchmark = "pkg.BenchmarkSameCommit"
	baseline := Baseline{Version: schemaVersion, Machines: map[string]MachineBaseline{
		"a": {Benchmarks: map[string]BenchmarkEntry{benchmark: {
			AllocsPerOp: 10, BytesPerOp: 100,
			AllocsSinceSHA: "same-sha", AllocsSinceAt: "2026-09-07T05:34:02Z",
			BytesSinceSHA: "same-sha", BytesSinceAt: "2026-09-07T05:34:02Z",
		}}},
		"b": {Benchmarks: map[string]BenchmarkEntry{benchmark: {
			AllocsPerOp: 20, BytesPerOp: 200,
			AllocsSinceSHA: "same-sha", AllocsSinceAt: "2026-09-08T11:00:00Z",
			BytesSinceSHA: "same-sha", BytesSinceAt: "2026-09-08T11:00:00Z",
		}}},
	}}
	bar := machineIndependentBar(baseline)[benchmark]
	if bar.AllocsPerOp != 10 || bar.BytesPerOp != 100 {
		t.Fatalf("bar = %d allocs/%d bytes, want the same-commit minimum 10/100",
			bar.AllocsPerOp, bar.BytesPerOp)
	}
	if bar.AllocsSinceSHA != "same-sha" {
		t.Fatalf("bar provenance = %q, want same-sha", bar.AllocsSinceSHA)
	}
}

// TestMachineIndependentBarDoesNotMergeDistinctSHAs is the converse: rows from
// different commits describe different code, so they are ranked rather than
// mixed, even when their timestamps collide.
func TestMachineIndependentBarDoesNotMergeDistinctSHAs(t *testing.T) {
	const benchmark = "pkg.BenchmarkDistinctCommits"
	const sameSecond = "2026-09-07T05:34:02Z"
	baseline := Baseline{Version: schemaVersion, Machines: map[string]MachineBaseline{
		"a": {Benchmarks: map[string]BenchmarkEntry{benchmark: {
			AllocsPerOp: 10, BytesPerOp: 100,
			AllocsSinceSHA: "aaaa1111", AllocsSinceAt: sameSecond,
			BytesSinceSHA: "aaaa1111", BytesSinceAt: sameSecond,
		}}},
		"b": {Benchmarks: map[string]BenchmarkEntry{benchmark: {
			AllocsPerOp: 20, BytesPerOp: 200,
			AllocsSinceSHA: "bbbb2222", AllocsSinceAt: sameSecond,
			BytesSinceSHA: "bbbb2222", BytesSinceAt: sameSecond,
		}}},
	}}
	bar := machineIndependentBar(baseline)[benchmark]
	// One row wins whole; a 10/200 or 10/100-from-two-commits mix would mean the
	// bar describes code that never existed.
	if !(bar.AllocsPerOp == 10 && bar.BytesPerOp == 100 && bar.AllocsSinceSHA == "aaaa1111") &&
		!(bar.AllocsPerOp == 20 && bar.BytesPerOp == 200 && bar.AllocsSinceSHA == "bbbb2222") {
		t.Fatalf("bar = %d allocs/%d bytes (%q), want one row intact",
			bar.AllocsPerOp, bar.BytesPerOp, bar.AllocsSinceSHA)
	}
}

// TestForceRebaselineKeepsDeterministicUnlessImproved pins that -force is a
// timing instrument: it replaces this machine's wall-clock numbers, but an
// allocs/bytes regression is a code fact no local recapture can accept, so the
// stored value and its provenance stay and the rejection is reported.
func TestForceRebaselineKeepsDeterministicUnlessImproved(t *testing.T) {
	const regressed = "pkg.BenchmarkRegressedBytes"
	const improved = "pkg.BenchmarkImprovedBytes"
	machine := Machine{OS: "darwin", Arch: "arm64", CPUModel: "Apple M3", GoVersion: "go1.26.5"}
	key := perfdata.MachineKey(machine)

	baseline := Baseline{Version: schemaVersion, Machines: map[string]MachineBaseline{key: {
		CapturedAt: "2026-08-01T00:00:00Z", CapturedAtSHA: "oldsha",
		Machine: machine,
		Benchmarks: map[string]BenchmarkEntry{
			regressed: {NSPerOp: 100, AllocsPerOp: 10, BytesPerOp: 100,
				AllocsSinceSHA: "oldsha", AllocsSinceAt: "2026-08-01T00:00:00Z",
				BytesSinceSHA: "oldsha", BytesSinceAt: "2026-08-01T00:00:00Z"},
			improved: {NSPerOp: 200, AllocsPerOp: 20, BytesPerOp: 200,
				AllocsSinceSHA: "oldsha", AllocsSinceAt: "2026-08-01T00:00:00Z",
				BytesSinceSHA: "oldsha", BytesSinceAt: "2026-08-01T00:00:00Z"},
		},
	}}}
	current := MachineBaseline{
		CapturedAt: "2026-09-18T17:13:27Z", CapturedAtSHA: "newsha",
		Machine: machine,
		Benchmarks: map[string]BenchmarkEntry{
			regressed: {NSPerOp: 90, AllocsPerOp: 12, BytesPerOp: 150},
			improved:  {NSPerOp: 180, AllocsPerOp: 15, BytesPerOp: 150},
		},
	}

	rejected := forceRebaseline(&baseline, key, current, false)

	got := baseline.Machines[key].Benchmarks[regressed]
	if got.NSPerOp != 90 {
		t.Fatalf("forced timing was not written: ns_per_op = %v, want 90", got.NSPerOp)
	}
	if got.AllocsPerOp != 10 || got.BytesPerOp != 100 {
		t.Fatalf("regression was accepted: %d allocs/%d bytes, want the stored 10/100",
			got.AllocsPerOp, got.BytesPerOp)
	}
	if got.AllocsSinceSHA != "oldsha" || got.BytesSinceSHA != "oldsha" {
		t.Fatalf("kept values were re-stamped: %q/%q", got.AllocsSinceSHA, got.BytesSinceSHA)
	}

	gotImproved := baseline.Machines[key].Benchmarks[improved]
	if gotImproved.AllocsPerOp != 15 || gotImproved.BytesPerOp != 150 {
		t.Fatalf("improvement was not adopted: %d allocs/%d bytes, want 15/150",
			gotImproved.AllocsPerOp, gotImproved.BytesPerOp)
	}
	if gotImproved.AllocsSinceSHA != "newsha" || gotImproved.BytesSinceSHA != "newsha" ||
		gotImproved.AllocsSinceAt != "2026-09-18T17:13:27Z" {
		t.Fatalf("improvement was not stamped: %q/%q",
			gotImproved.AllocsSinceSHA, gotImproved.BytesSinceSHA)
	}

	if len(rejected) != 2 {
		t.Fatalf("rejected = %+v, want the allocs and bytes regressions of %s", rejected, regressed)
	}
	for _, r := range rejected {
		if r.Name != regressed {
			t.Fatalf("rejected names %q, want %q", r.Name, regressed)
		}
	}
}

// TestForceRebaselineAcceptsDeterministicWhenAsked pins the explicit escape
// hatch: an operator who has justified an allocation regression can record it,
// and the new numbers then carry this run's provenance.
func TestForceRebaselineAcceptsDeterministicWhenAsked(t *testing.T) {
	const benchmark = "pkg.BenchmarkAcceptedRegression"
	machine := Machine{OS: "darwin", Arch: "arm64", CPUModel: "Apple M3", GoVersion: "go1.26.5"}
	key := perfdata.MachineKey(machine)
	baseline := Baseline{Version: schemaVersion, Machines: map[string]MachineBaseline{key: {
		CapturedAt: "2026-08-01T00:00:00Z", CapturedAtSHA: "oldsha",
		Machine:    machine,
		Benchmarks: map[string]BenchmarkEntry{benchmark: {NSPerOp: 10, AllocsPerOp: 10, BytesPerOp: 100}},
	}}}
	current := MachineBaseline{
		CapturedAt: "2026-08-18T21:00:00Z", CapturedAtSHA: "accepted-sha",
		Machine:    machine,
		Benchmarks: map[string]BenchmarkEntry{benchmark: {NSPerOp: 15, AllocsPerOp: 12, BytesPerOp: 120}},
	}
	if rejected := forceRebaseline(&baseline, key, current, true); len(rejected) != 0 {
		t.Fatalf("rejected = %+v, want none under -accept-deterministic", rejected)
	}
	got := baseline.Machines[key].Benchmarks[benchmark]
	if got.AllocsPerOp != 12 || got.BytesPerOp != 120 {
		t.Fatalf("accepted numbers not written: %d allocs/%d bytes", got.AllocsPerOp, got.BytesPerOp)
	}
	if got.AllocsSinceSHA != "accepted-sha" || got.BytesSinceSHA != "accepted-sha" {
		t.Fatalf("accepted numbers not stamped: %q/%q", got.AllocsSinceSHA, got.BytesSinceSHA)
	}
}

// TestDeterministicProvenanceFallsBackToProfileForNeverRatchetedRow pins the one
// case where the profile's own capture identity IS the deterministic
// provenance: an entry that carries no best_since stamp has never been
// ratcheted, so seed or capture wrote every one of its numbers in the run the
// profile records. Reading them as undated would throw away a fact the file
// states, and hand the bar back to the fleet-wide minimum the newest-row rule
// exists to replace.
func TestDeterministicProvenanceFallsBackToProfileForNeverRatchetedRow(t *testing.T) {
	const benchmark = "pkg.BenchmarkSeeded"
	baseline := Baseline{Version: schemaVersion, Machines: map[string]MachineBaseline{
		// Never ratcheted: no best_since, no deterministic stamp.
		"a": {
			CapturedAt: "2026-09-07T05:34:02Z", CapturedAtSHA: "477a5d36e25f",
			Benchmarks: map[string]BenchmarkEntry{benchmark: {AllocsPerOp: 13995, BytesPerOp: 1016128}},
		},
		// An older tier whose lower numbers describe superseded code.
		"b": {
			CapturedAt: "2026-08-11T16:10:14Z", CapturedAtSHA: "f3ca5f9be1da",
			Benchmarks: map[string]BenchmarkEntry{benchmark: {AllocsPerOp: 7041, BytesPerOp: 743356}},
		},
	}}
	bar := machineIndependentBar(baseline)[benchmark]
	if bar.AllocsPerOp != 13995 || bar.BytesPerOp != 1016128 {
		t.Fatalf("bar = %d allocs/%d bytes, want the newest profile's 13995/1016128",
			bar.AllocsPerOp, bar.BytesPerOp)
	}
	if bar.AllocsSinceSHA != "477a5d36e25f" {
		t.Fatalf("bar provenance = %q, want 477a5d36e25f", bar.AllocsSinceSHA)
	}
}

// TestForceRebaselineRecordsProvenanceOfKeptValues pins that a kept
// deterministic value keeps its provenance in a form the next reader can still
// see. The stored row's date may be implicit in the profile it sat in; a forced
// write stamps the profile with this run, so the date has to be written onto
// the entry or it is lost.
func TestForceRebaselineRecordsProvenanceOfKeptValues(t *testing.T) {
	const benchmark = "pkg.BenchmarkKeptProvenance"
	machine := Machine{OS: "darwin", Arch: "arm64", CPUModel: "Apple M3", GoVersion: "go1.26.5"}
	key := perfdata.MachineKey(machine)
	baseline := Baseline{Version: schemaVersion, Machines: map[string]MachineBaseline{key: {
		CapturedAt: "2026-08-01T00:00:00Z", CapturedAtSHA: "seedsha",
		Machine: machine,
		// Never ratcheted: its numbers are the profile's own capture.
		Benchmarks: map[string]BenchmarkEntry{benchmark: {NSPerOp: 100, AllocsPerOp: 10, BytesPerOp: 100}},
	}}}
	current := MachineBaseline{
		CapturedAt: "2026-09-18T17:13:27Z", CapturedAtSHA: "newsha",
		Machine:    machine,
		Benchmarks: map[string]BenchmarkEntry{benchmark: {NSPerOp: 90, AllocsPerOp: 12, BytesPerOp: 150}},
	}
	if rejected := forceRebaseline(&baseline, key, current, false); len(rejected) != 2 {
		t.Fatalf("rejected = %+v, want the allocs and bytes regressions", rejected)
	}
	got := baseline.Machines[key].Benchmarks[benchmark]
	if got.AllocsPerOp != 10 || got.BytesPerOp != 100 {
		t.Fatalf("kept values = %d allocs/%d bytes, want 10/100", got.AllocsPerOp, got.BytesPerOp)
	}
	if got.AllocsSinceSHA != "seedsha" || got.BytesSinceSHA != "seedsha" ||
		got.AllocsSinceAt != "2026-08-01T00:00:00Z" {
		t.Fatalf("kept provenance = %q/%q, want seedsha for both",
			got.AllocsSinceSHA, got.BytesSinceSHA)
	}
}

// TestDeterministicProvenanceIsPerMetric pins that allocs and bytes are dated
// independently. A run that lowers one and regresses the other stores a pair
// neither run measured — the new allocs beside the old bytes — so one stamp
// over both necessarily lies about one of them, and dating the old bytes to the
// new commit resurrects them as a current floor.
func TestDeterministicProvenanceIsPerMetric(t *testing.T) {
	const benchmark = "pkg.BenchmarkMixedChange"
	january := MachineBaseline{
		CapturedAt: "2026-01-10T00:00:00Z", CapturedAtSHA: "jansha",
		Benchmarks: map[string]BenchmarkEntry{benchmark: {
			NSPerOp: 100, AllocsPerOp: 10, BytesPerOp: 100,
			AllocsSinceSHA: "jansha", AllocsSinceAt: "2026-01-10T00:00:00Z",
			BytesSinceSHA: "jansha", BytesSinceAt: "2026-01-10T00:00:00Z",
		}},
	}
	// March: allocs 10 -> 9, bytes 100 -> 200. The ratchet keeps 100 bytes.
	march := MachineBaseline{
		CapturedAt: "2026-03-10T00:00:00Z", CapturedAtSHA: "marsha",
		Benchmarks: map[string]BenchmarkEntry{benchmark: {
			NSPerOp: 100, AllocsPerOp: 9, BytesPerOp: 200,
		}},
	}
	merged, _ := ratchetMerge(january, march)
	got := merged.Benchmarks[benchmark]
	if got.AllocsPerOp != 9 || got.BytesPerOp != 100 {
		t.Fatalf("merged = %d allocs/%d bytes, want 9/100", got.AllocsPerOp, got.BytesPerOp)
	}
	if got.AllocsSinceSHA != "marsha" {
		t.Fatalf("allocs provenance = %q, want marsha", got.AllocsSinceSHA)
	}
	if got.BytesSinceSHA != "jansha" {
		t.Fatalf("bytes provenance = %q, want jansha (March never measured 100 bytes)", got.BytesSinceSHA)
	}

	// February measured the benchmark for real; it is the newest row that has
	// bytes, while March is the newest that has allocs.
	baseline := Baseline{Version: schemaVersion, Machines: map[string]MachineBaseline{
		"a": merged,
		"b": {
			CapturedAt: "2026-02-10T00:00:00Z", CapturedAtSHA: "febsha",
			Benchmarks: map[string]BenchmarkEntry{benchmark: {
				AllocsPerOp: 15, BytesPerOp: 150,
				AllocsSinceSHA: "febsha", AllocsSinceAt: "2026-02-10T00:00:00Z",
				BytesSinceSHA: "febsha", BytesSinceAt: "2026-02-10T00:00:00Z",
			}},
		},
	}}
	bar := machineIndependentBar(baseline)[benchmark]
	if bar.BytesPerOp != 150 || bar.BytesSinceSHA != "febsha" {
		t.Fatalf("bytes bar = %d (%q), want 150 (febsha)", bar.BytesPerOp, bar.BytesSinceSHA)
	}
	if bar.AllocsPerOp != 9 || bar.AllocsSinceSHA != "marsha" {
		t.Fatalf("allocs bar = %d (%q), want 9 (marsha)", bar.AllocsPerOp, bar.AllocsSinceSHA)
	}
}

// TestForcedKeepStampsAdoptedAllocsSeparately is the inverse, as it occurs on a
// forced timing recapture: bytes are kept and allocs adopted, so the adopted
// allocs must carry this run's date or the gate keeps selecting a higher bar
// from an older tier and misses later regressions under it.
func TestForcedKeepStampsAdoptedAllocsSeparately(t *testing.T) {
	const benchmark = "pkg.BenchmarkIRCompile [bytecode]"
	machine := Machine{OS: "darwin", Arch: "arm64", CPUModel: "Apple M3", GoVersion: "go1.26.5"}
	key := perfdata.MachineKey(machine)
	baseline := Baseline{Version: schemaVersion, Machines: map[string]MachineBaseline{
		key: {
			CapturedAt: "2026-08-24T20:18:19Z", CapturedAtSHA: "dadb8e0b54b7",
			Machine: machine,
			Benchmarks: map[string]BenchmarkEntry{benchmark: {
				NSPerOp: 15137978, AllocsPerOp: 114255, BytesPerOp: 4963229,
				BestSinceSHA: "dadb8e0b54b7", BestSinceAt: "2026-08-24T20:18:19Z",
			}},
		},
		"amd64/EPYC": {
			CapturedAt: "2026-09-07T05:34:02Z", CapturedAtSHA: "477a5d36e25f",
			Machine: Machine{OS: "linux", Arch: "amd64", CPUModel: "EPYC"},
			Benchmarks: map[string]BenchmarkEntry{benchmark: {
				AllocsPerOp: 113777, BytesPerOp: 4941030,
			}},
		},
	}}
	current := MachineBaseline{
		CapturedAt: "2026-09-18T17:13:27Z", CapturedAtSHA: "3bbde90a0513",
		Machine: machine,
		Benchmarks: map[string]BenchmarkEntry{benchmark: {
			NSPerOp: 17295854, AllocsPerOp: 105814, BytesPerOp: 5232922,
		}},
	}
	rejected := forceRebaseline(&baseline, key, current, false)
	if len(rejected) != 1 || rejected[0].Metric != "bytes/op" {
		t.Fatalf("rejected = %+v, want only the bytes regression", rejected)
	}
	got := baseline.Machines[key].Benchmarks[benchmark]
	if got.AllocsPerOp != 105814 || got.AllocsSinceSHA != "3bbde90a0513" {
		t.Fatalf("adopted allocs = %d (%q), want 105814 (3bbde90a0513)", got.AllocsPerOp, got.AllocsSinceSHA)
	}
	// The stored M3 row had been ratcheted, so when its bytes were measured is
	// not knowable: the kept value stays undated rather than borrowing a date.
	if got.BytesPerOp != 4963229 || got.BytesSinceSHA != "" {
		t.Fatalf("kept bytes = %d (%q), want 4963229 undated", got.BytesPerOp, got.BytesSinceSHA)
	}
	// The EPYC row was never ratcheted, so it dates its own bytes and supplies
	// the bytes bar; the M3's freshly measured allocs supply the allocs bar.
	if b := machineIndependentBar(baseline)[benchmark]; b.BytesPerOp != 4941030 || b.BytesSinceSHA != "477a5d36e25f" {
		t.Fatalf("bytes bar = %d (%q), want 4941030 (477a5d36e25f)", b.BytesPerOp, b.BytesSinceSHA)
	}

	bar := machineIndependentBar(baseline)[benchmark]
	if bar.AllocsPerOp != 105814 {
		t.Fatalf("allocs bar = %d, want the M3's freshly measured 105814", bar.AllocsPerOp)
	}
	worse := MachineBaseline{
		Machine:    machine,
		Benchmarks: map[string]BenchmarkEntry{benchmark: {AllocsPerOp: 110000, BytesPerOp: 4963229}},
	}
	if n := compareDeterministic(bar2map(bar, benchmark), worse, allocBudget); n == 0 {
		t.Fatal("110000 allocs against a 105814 bar was not reported as a regression")
	}
}

// TestMachineIndependentBarMergedGroupTakesNewestTimestamp pins that a merged
// same-commit group is ranked by its NEWEST measurement, not by whichever
// profile happened to be read first: the group's claim on being current rests
// on its latest capture.
func TestMachineIndependentBarMergedGroupTakesNewestTimestamp(t *testing.T) {
	const benchmark = "pkg.BenchmarkGroupAt"
	baseline := Baseline{Version: schemaVersion, Machines: map[string]MachineBaseline{
		"a": {Benchmarks: map[string]BenchmarkEntry{benchmark: {
			AllocsPerOp: 10, BytesPerOp: 100,
			AllocsSinceSHA: "groupsha", AllocsSinceAt: "2026-01-10T00:00:00Z",
			BytesSinceSHA: "groupsha", BytesSinceAt: "2026-01-10T00:00:00Z",
		}}},
		"b": {Benchmarks: map[string]BenchmarkEntry{benchmark: {
			AllocsPerOp: 20, BytesPerOp: 200,
			AllocsSinceSHA: "groupsha", AllocsSinceAt: "2026-03-10T00:00:00Z",
			BytesSinceSHA: "groupsha", BytesSinceAt: "2026-03-10T00:00:00Z",
		}}},
		"c": {Benchmarks: map[string]BenchmarkEntry{benchmark: {
			AllocsPerOp: 50, BytesPerOp: 500,
			AllocsSinceSHA: "febsha", AllocsSinceAt: "2026-02-10T00:00:00Z",
			BytesSinceSHA: "febsha", BytesSinceAt: "2026-02-10T00:00:00Z",
		}}},
	}}
	bar := machineIndependentBar(baseline)[benchmark]
	if bar.AllocsPerOp != 10 || bar.BytesPerOp != 100 || bar.AllocsSinceSHA != "groupsha" {
		t.Fatalf("bar = %d allocs/%d bytes (%q), want 10/100 (groupsha, the March group)",
			bar.AllocsPerOp, bar.BytesPerOp, bar.AllocsSinceSHA)
	}
}
