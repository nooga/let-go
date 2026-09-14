package main

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/nooga/let-go/pkg/perfdata"
)

func TestSplitBenchmarkName(t *testing.T) {
	pkg, name := splitBenchmarkName("github.com/nooga/let-go/pkg/vm.BenchmarkFuncInvoke/Direct")
	if pkg != "pkg/vm" {
		t.Fatalf("package = %q, want pkg/vm", pkg)
	}
	if name != "BenchmarkFuncInvoke/Direct" {
		t.Fatalf("name = %q, want BenchmarkFuncInvoke/Direct", name)
	}

	pkg, name = splitBenchmarkName("github.com/nooga/let-go/pkg/ir.BenchmarkIRCompile [gogen_ir]")
	if pkg != "pkg/ir" {
		t.Fatalf("package = %q, want pkg/ir", pkg)
	}
	if name != "BenchmarkIRCompile [gogen_ir]" {
		t.Fatalf("name = %q, want BenchmarkIRCompile [gogen_ir]", name)
	}
}

func TestCompareWithHistorical(t *testing.T) {
	current := Baseline{Benchmarks: map[string]BenchmarkEntry{
		"pkg.BenchmarkA": {RatioToAnchor: 80},
		"pkg.BenchmarkB": {RatioToAnchor: 120},
		"pkg.BenchmarkC": {RatioToAnchor: 50},
	}}
	reference := Baseline{Benchmarks: map[string]BenchmarkEntry{
		"pkg.BenchmarkA": {RatioToAnchor: 100},
		"pkg.BenchmarkB": {RatioToAnchor: 100},
		"pkg.BenchmarkD": {RatioToAnchor: 10},
	}}

	changes, summary := compare(current, reference)
	if summary.Common != 2 {
		t.Fatalf("common = %d, want 2", summary.Common)
	}
	if summary.New != 1 {
		t.Fatalf("new = %d, want 1", summary.New)
	}
	if summary.Missing != 1 {
		t.Fatalf("missing = %d, want 1", summary.Missing)
	}
	if summary.Faster != 1 {
		t.Fatalf("faster = %d, want 1", summary.Faster)
	}
	if summary.Slower != 1 {
		t.Fatalf("slower = %d, want 1", summary.Slower)
	}
	if summary.MedianDelta != 0 {
		t.Fatalf("median delta = %v, want 0", summary.MedianDelta)
	}
	if got := changes["pkg.BenchmarkA"]; !near(got, -0.2) {
		t.Fatalf("BenchmarkA delta = %v, want -0.2", got)
	}
	if got := changes["pkg.BenchmarkB"]; !near(got, 0.2) {
		t.Fatalf("BenchmarkB delta = %v, want 0.2", got)
	}
}

func TestBuildCharts(t *testing.T) {
	timeline := []Snapshot{
		makeSnapshot("a", Baseline{
			CapturedAt:    "2026-06-01T00:00:00Z",
			CapturedAtSHA: "aaaaaaaaaaaa",
			Benchmarks: map[string]BenchmarkEntry{
				"github.com/nooga/let-go/test.BenchmarkClojureTestSuite [bytecode]": {RatioToAnchor: 100, AllocsPerOp: 10, BytesPerOp: 1000},
			},
		}),
		makeSnapshot("b", Baseline{
			CapturedAt:    "2026-06-02T00:00:00Z",
			CapturedAtSHA: "bbbbbbbbbbbb",
			Benchmarks: map[string]BenchmarkEntry{
				"github.com/nooga/let-go/test.BenchmarkClojureTestSuite [bytecode]": {RatioToAnchor: 80, AllocsPerOp: 9, BytesPerOp: 900},
			},
		}),
	}

	charts := buildCharts(timeline, Baseline{}, "", defaultBudgetFraction)
	if len(charts) != 3 {
		t.Fatalf("chart count = %d, want 3", len(charts))
	}
	if charts[0].Title != "End-to-end suite" {
		t.Fatalf("first chart = %q, want End-to-end suite", charts[0].Title)
	}
	if len(charts[0].Series) != 1 {
		t.Fatalf("series count = %d, want 1", len(charts[0].Series))
	}
	if len(charts[0].Series[0].Points) != 2 {
		t.Fatalf("point count = %d, want 2", len(charts[0].Series[0].Points))
	}
	if charts[0].Series[0].Path == "" {
		t.Fatal("expected SVG path")
	}
	if charts[0].Series[0].Points[1].X <= charts[0].Series[0].Points[0].X {
		t.Fatalf("x coordinates did not advance: %#v", charts[0].Series[0].Points)
	}
}

func near(got, want float64) bool {
	diff := got - want
	return diff < 0.0000001 && diff > -0.0000001
}

func TestFormatNS(t *testing.T) {
	tests := map[float64]string{
		12.345:         "12.35 ns",
		12_345:         "12.35 us",
		12_345_000:     "12.35 ms",
		12_345_000_000: "12.35 s",
	}
	for input, want := range tests {
		if got := formatNS(input); got != want {
			t.Fatalf("formatNS(%v) = %q, want %q", input, got, want)
		}
	}
}

func TestLoadTimelineSkipsCorruptSnapshots(t *testing.T) {
	dir := t.TempDir()
	timelineDir := filepath.Join(dir, "timeline")
	historicalDir := filepath.Join(dir, "historical")
	if err := os.MkdirAll(timelineDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(historicalDir, 0o755); err != nil {
		t.Fatal(err)
	}
	good := `{"version":1,"captured_at":"2026-06-01T00:00:00Z","captured_at_sha":"abc","benchmarks":{"pkg.BenchmarkA":{"ns_per_op":1,"ratio_to_anchor":2}}}`
	if err := os.WriteFile(filepath.Join(timelineDir, "good.json"), []byte(good), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(timelineDir, "bad.json"), []byte(`{"version":`), 0o644); err != nil {
		t.Fatal(err)
	}

	var stderr strings.Builder
	oldStderr := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stderr = w
	snapshots, err := loadTimeline(timelineDir, historicalDir, Baseline{})
	_ = w.Close()
	os.Stderr = oldStderr
	if _, copyErr := io.Copy(&stderr, r); copyErr != nil {
		t.Fatal(copyErr)
	}
	if err != nil {
		t.Fatalf("loadTimeline returned error: %v", err)
	}
	if len(snapshots) != 1 {
		t.Fatalf("snapshot count = %d, want 1", len(snapshots))
	}
	if !strings.Contains(stderr.String(), "skipping") {
		t.Fatalf("stderr = %q, want skip warning", stderr.String())
	}
}

func TestFormatRatioCompactsLargeValues(t *testing.T) {
	tests := map[float64]string{
		1_124_520_183: "1.12B",
		8_519_621:     "8.52M",
		66_309:        "66.3k",
		16.209:        "16.2",
		4.688:         "4.69",
	}
	for input, want := range tests {
		if got := formatRatio(input); got != want {
			t.Fatalf("formatRatio(%v) = %q, want %q", input, got, want)
		}
	}
}

// Timeline-explorer regression tests. Each of these pins a failure that
// reached review on #880: a null slice the page reads .length on, a timestamp
// only Chromium would parse, a tier reachable in the data but not in the
// filter, and a link that assumed one directory layout.

func viewerTestSnapshot(at, cpu string, ratio float64) Snapshot {
	return Snapshot{Baseline: Baseline{
		CapturedAt: at,
		Machine:    Machine{Arch: "amd64", CPUModel: cpu},
		Anchor:     perfdata.Anchor{NSPerOp: 1},
		Benchmarks: map[string]BenchmarkEntry{
			"github.com/nooga/let-go/pkg/ir.BenchmarkIRCompile [bytecode]": {
				RatioToAnchor: ratio, AllocsPerOp: 10, BytesPerOp: 100,
			},
		},
	}}
}

func TestBuildViewerDataEmitsNoNilPointSlices(t *testing.T) {
	// A nil slice marshals as null, and the page reads .length on it, so an
	// empty series must not survive into the payload at all.
	data := buildViewerData([]Snapshot{viewerTestSnapshot("2026-06-04T22:25:22Z", "AMD EPYC 7763", 2)})
	blob, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if bytes.Contains(blob, []byte(`"pts":null`)) {
		t.Error(`payload contains "pts":null; the page does s.pts.length on it`)
	}
	for _, ch := range data.Charts {
		if len(ch.Series) == 0 {
			t.Errorf("chart %q has no series and should have been dropped", ch.Title)
		}
		for _, s := range ch.Series {
			if len(s.Points) == 0 {
				t.Errorf("chart %q series %q is empty and should have been dropped", ch.Title, s.Label)
			}
		}
	}
}

func TestBuildViewerDataKeepsTimestampsParseable(t *testing.T) {
	// Date.parse only guarantees RFC3339. formatDate's "2006-01-02 15:04 UTC"
	// is accepted by Chromium and rejected by JavaScriptCore, so emitting it
	// renders an empty page in Safari.
	const want = "2026-06-04T22:25:22Z"
	data := buildViewerData([]Snapshot{
		viewerTestSnapshot(want, "AMD EPYC 7763", 2),
		viewerTestSnapshot("2026-06-05T10:00:00Z", "AMD EPYC 7763", 3),
	})
	for _, ch := range data.Charts {
		for _, s := range ch.Series {
			for _, p := range s.Points {
				if _, err := time.Parse(time.RFC3339, p.Date); err != nil {
					t.Fatalf("chart %q series %q: %q is not RFC3339", ch.Title, s.Label, p.Date)
				}
			}
		}
	}
}

func TestBuildViewerDataListsEveryPlottedTier(t *testing.T) {
	// A tier that appears only in the geomean must still be selectable, or
	// "All (n)" plots more tiers than it counts.
	//
	// The fixture has to earn that: both snapshots share an UNCHARTED
	// benchmark, so the geomean basket is non-empty, while only the first also
	// carries a charted one. Apple M3 is therefore reachable through the
	// geomean and through nothing else — which is the case that a CPU list
	// built from the chart specs alone silently drops.
	// Two of them: viewerGeomean needs a basket of at least two benchmarks
	// before it will emit a line at all.
	uncharted := []string{
		"github.com/nooga/let-go/pkg/vm.BenchmarkSomethingElse",
		"github.com/nooga/let-go/pkg/vm.BenchmarkSomethingElser",
	}
	withBasket := func(at, cpu string, charted bool) Snapshot {
		snap := viewerTestSnapshot(at, cpu, 2)
		for i, name := range uncharted {
			snap.Baseline.Benchmarks[name] = BenchmarkEntry{
				RatioToAnchor: float64(5 + i), AllocsPerOp: 1, BytesPerOp: 8,
			}
		}
		if !charted {
			delete(snap.Baseline.Benchmarks, "github.com/nooga/let-go/pkg/ir.BenchmarkIRCompile [bytecode]")
		}
		return snap
	}
	data := buildViewerData([]Snapshot{
		withBasket("2026-06-04T22:25:22Z", "AMD EPYC 7763", true),
		withBasket("2026-06-05T10:00:00Z", "Apple M3", false),
	})
	if len(data.Charts) == 0 {
		t.Fatal("fixture produced no charts; the geomean basket is empty")
	}
	if data.Charts[0].Title[:7] != "Overall" {
		t.Fatalf("expected the geomean chart first, got %q", data.Charts[0].Title)
	}
	listed := map[string]bool{}
	for _, c := range data.CPUs {
		listed[c] = true
	}
	for _, ch := range data.Charts {
		for _, s := range ch.Series {
			for _, p := range s.Points {
				if !listed[p.CPU] {
					t.Errorf("chart %q plots tier %q, which is not in CPUs %v", ch.Title, p.CPU, data.CPUs)
				}
			}
		}
	}
}

func TestRelLinkNamesTheFileNotTheDirectory(t *testing.T) {
	for _, tc := range []struct{ from, to, want string }{
		{"out/perf/index.html", "out/perf/explore/index.html", "explore/index.html"},
		{"out/perf/explore/index.html", "out/perf/index.html", "../index.html"},
		// A flat layout has no index.html to fall back on, so a bare "./" here
		// would link the summary to itself.
		{"out/summary.html", "out/timeline.html", "timeline.html"},
		{"out/timeline.html", "out/summary.html", "summary.html"},
	} {
		if got := relLink(tc.from, tc.to); got != tc.want {
			t.Errorf("relLink(%q, %q) = %q, want %q", tc.from, tc.to, got, tc.want)
		}
	}
}
