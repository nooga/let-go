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

	// Create amd64 snapshots (should be seeded)
	createMockSnapshot("20260801T010134Z-b170a08eef47-amd64-amd-epyc-7763.json",
		"b170a08eef47", "2026-08-01T01:01:34Z", "amd64", "AMD EPYC 7763")

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
	seedBaseline(baselineFile, timelineDir, "unused-release-sha")

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
			if _, ok := mb.Benchmarks["test.BenchmarkA"]; !ok {
				t.Error("stable benchmark test.BenchmarkA should be kept")
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

func TestForceRebaselinePropagatesDeterministicBarAcrossProfiles(t *testing.T) {
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

	forceRebaseline(&baseline, currentKey, current)

	gotCurrent := baseline.Machines[currentKey].Benchmarks[benchmark]
	if gotCurrent.NSPerOp != 15 || gotCurrent.AllocsPerOp != 12 || gotCurrent.BytesPerOp != 120 {
		t.Fatalf("current profile was not replaced: %+v", gotCurrent)
	}
	gotOther := baseline.Machines[otherKey].Benchmarks[benchmark]
	if gotOther.NSPerOp != 20 || gotOther.RatioToAnchor != 20 {
		t.Fatalf("other profile timing changed: %+v", gotOther)
	}
	if gotOther.AllocsPerOp != 12 || gotOther.BytesPerOp != 120 {
		t.Fatalf("other deterministic bar = %d allocs/%d bytes, want 12/120", gotOther.AllocsPerOp, gotOther.BytesPerOp)
	}
	if gotOther.BestSinceSHA != "accepted-sha" || gotOther.BestSinceAt != current.CapturedAt {
		t.Fatalf("other deterministic provenance was not updated: %+v", gotOther)
	}
	if got := baseline.Machines[currentKey].Benchmarks[outsideScope]; got.NSPerOp != 11 || got.AllocsPerOp != 4 || got.BytesPerOp != 40 {
		t.Fatalf("current out-of-scope benchmark changed: %+v", got)
	}
	if got := baseline.Machines[otherKey].Benchmarks[outsideScope]; got.AllocsPerOp != 3 || got.BytesPerOp != 30 {
		t.Fatalf("other out-of-scope benchmark changed: %+v", got)
	}
	bar := machineIndependentBar(baseline)[benchmark]
	if bar.AllocsPerOp != 12 || bar.BytesPerOp != 120 {
		t.Fatalf("global deterministic bar = %d allocs/%d bytes, want 12/120", bar.AllocsPerOp, bar.BytesPerOp)
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
