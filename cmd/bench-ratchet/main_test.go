package main

import (
	"encoding/json"
	"os"
	"path/filepath"
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
	records := []StreamRecord{
		{Package: anchorPackage, Name: anchorName, Iterations: 100, NSPerOp: 10, CapturedAt: "2026-06-01T00:00:00Z"},
		{Package: anchorPackage, Name: anchorName, Iterations: 90, NSPerOp: 20, CapturedAt: "2026-06-01T00:00:01Z"},
		{Package: "pkg", Name: "BenchmarkA", Iterations: 50, NSPerOp: 30, BytesPerOp: 100, AllocsPerOp: 1, CapturedAt: "2026-06-01T00:00:02Z"},
		{Package: "pkg", Name: "BenchmarkA", Iterations: 60, NSPerOp: 60, BytesPerOp: 200, AllocsPerOp: 3, CapturedAt: "2026-06-01T00:00:03Z"},
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
	if entry.NSPerOp != 45 {
		t.Fatalf("ns/op = %v, want 45", entry.NSPerOp)
	}
	if entry.BytesPerOp != 150 {
		t.Fatalf("bytes/op = %d, want 150", entry.BytesPerOp)
	}
	if entry.AllocsPerOp != 2 {
		t.Fatalf("allocs/op = %d, want 2", entry.AllocsPerOp)
	}
	if len(entry.Samples) != 2 {
		t.Fatalf("samples = %d, want 2", len(entry.Samples))
	}
	if entry.Samples[0].RatioToAnchor != 2 {
		t.Fatalf("first sample ratio = %v, want 2", entry.Samples[0].RatioToAnchor)
	}
	if entry.Samples[1].RatioToAnchor != 4 {
		t.Fatalf("second sample ratio = %v, want 4", entry.Samples[1].RatioToAnchor)
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

// The full + fast profiles run the Clojure test suite once per execution mode
// (count=1): a single suite pass takes minutes and run-to-run variance is
// negligible, so 3 samples would just triple wall time. The cheap,
// benchtime-bounded vm/ir jobs keep the CLI default (count=0 → -count flag).
//
// There are exactly three suite modes, distinguished by the LG_SUITE_IR env
// toggle crossed with the gogen_ir build tag:
//   - bytecode    : *ir-compile* off, untagged
//   - ir_bytecode : *ir-compile* on  (LG_SUITE_IR=1), untagged
//   - aot_native  : *ir-compile* on  (LG_SUITE_IR=1), -tags gogen_ir
func TestSuiteJobsPinCountToOne(t *testing.T) {
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
		for _, j := range jobs {
			if j.pkg == suitePackage {
				if j.count != 1 {
					t.Errorf("full=%v: suite job [%s] count = %d, want 1", full, j.variant, j.count)
				}
				suite[j.variant] = j
			} else if j.count != 0 {
				t.Errorf("full=%v: non-suite job %s [%s] count = %d, want 0 (use CLI default)", full, j.pkg, j.variant, j.count)
			}
		}
		if len(suite) != 3 {
			t.Fatalf("full=%v: expected 3 suite variants (bytecode, ir_bytecode, aot_native), got %d: %v", full, len(suite), keysOf(suite))
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
