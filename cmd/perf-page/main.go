package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/nooga/let-go/pkg/perfdata"
	"golang.org/x/perf/benchunit"
)

const modulePrefix = "github.com/nooga/let-go/"

// Baseline here is one MACHINE PROFILE's data — the page renders a single
// machine's numbers (anchor-relative ratios only normalize within a CPU model).
// MultiBaseline is the on-disk arch-partitioned file that loadBaseline reads
// and selects a profile from.
type Baseline = perfdata.MachineBaseline
type MultiBaseline = perfdata.Baseline
type Machine = perfdata.Machine
type Anchor = perfdata.Anchor
type BenchmarkEntry = perfdata.BenchmarkEntry
type BenchmarkSample = perfdata.BenchmarkSample

type BenchmarkRow struct {
	FullName     string
	Package      string
	Name         string
	NSPerOp      float64
	AllocsPerOp  float64
	BytesPerOp   float64
	Ratio        float64
	BestSinceSHA string
	BestSinceAt  string
	Delta        *float64
	BarWidth     float64
}

type ChangeRow struct {
	BenchmarkRow
	OldRatio float64
	NewRatio float64
}

type Snapshot struct {
	Name     string
	Baseline Baseline
	Captured time.Time
}

type Chart struct {
	Title    string
	Subtitle string
	Unit     string
	Series   []ChartSeries
	YMin     float64
	YMax     float64
	YMinText string
	YMaxText string

	// Reference line + regression budget band. The reference is the most
	// recent frozen release baseline; the budget band marks the ratchet's
	// regression ceiling (reference * (1+budget)) so "stay under the line"
	// reads off the chart.
	HasRef   bool
	RefY     float64
	RefLabel string

	ShowRefLine bool // draw the dashed line in-plot; false when the reference is off-scale

	HasBudget  bool
	BudgetY    float64 // y of the regression ceiling (top edge of the danger band)
	BudgetText string  // e.g. "±5%"
	RefLegend  string  // reference's legend entry, e.g. "v1.8.0 = 599M ↑" (↑/↓ marks off-scale)

	// Latest point's standing vs the reference (the status pill).
	Status      string
	StatusClass string // good | bad | flat | none

	XTicks []ChartXTick

	// Per-snapshot delta between the two series (series[1] vs series[0]),
	// e.g. "how much faster aot_native is than ir_bytecode" at each point.
	Deltas       []ChartDelta
	DeltaCaption string // footer caption, e.g. "Δ aot_native vs ir_bytecode"
}

type ChartXTick struct {
	X     float64
	Label string
}

type ChartDelta struct {
	X    float64
	Y    float64
	Text string
}

type ChartSeries struct {
	Label    string
	Color    string
	Path     string // line; subpaths restart (M) across missing snapshots so gaps don't bridge
	BandPath string // variance envelope (sample min..max), one closed subpath per contiguous run
	Points   []ChartPoint
}

type ChartPoint struct {
	X     float64
	Y     float64
	Index int
	Date  string
	SHA   string
	Value float64
	Text  string

	// Per-run sample envelope: Low/High are the spread's metric values,
	// LowY/HighY their pixel-y. HasBand is false when a snapshot carried <2
	// samples (nothing to spread), so the band collapses to the point.
	Low     float64
	High    float64
	LowY    float64
	HighY   float64
	HasBand bool
	Spread  string // the "low .. high" fragment shown in the point's hover tooltip
}

type Summary struct {
	BenchmarkCount int
	PackageCount   int
	ZeroAllocs     int
	Common         int
	New            int
	Missing        int
	Faster         int
	Slower         int
	MedianDelta    float64
}

type PageData struct {
	Title             string
	LogoDataURI       template.URL
	Current           Baseline
	ReferenceName     string
	Summary           Summary
	Timeline          []Snapshot
	Charts            []Chart
	Rows              []BenchmarkRow
	TopImprovements   []ChangeRow
	TopSlowdowns      []ChangeRow
	RecentlyTightened []BenchmarkRow

	// ExplorerURL is where the page fetches the tidy explorer data
	// ({date, cpu, bench, metric, value} rows) at load time — written as a
	// separate explorer.json rather than inlined, so the HTML stays small.
	ExplorerURL string
}

// explorerDatum is one tidy row: a single metric of one benchmark at one
// timeline snapshot, tagged by the CPU that produced it.
type explorerDatum struct {
	Date   string  `json:"date"`
	CPU    string  `json:"cpu"`
	Bench  string  `json:"bench"`
	Metric string  `json:"metric"`
	Value  float64 `json:"value"`
	// Lo/Hi are the min/max of this metric across the snapshot's retained
	// samples — the run-to-run spread the explorer shades as a band. With a
	// single sample (the end-to-end suite) Lo==Hi==Value (no band). This is an
	// honest spread, NOT a 95% CI (typically only ~3 samples).
	Lo float64 `json:"lo"`
	Hi float64 `json:"hi"`
	// Samples are the individual per-run measurements behind Value (the raw
	// gathered points), so the chart/sparkline can show each sample as an
	// observable dot rather than only the summarized mean+band.
	Samples []float64 `json:"samples,omitempty"`
}

// metricInfo carries a metric's display unit and direction so the explorer can
// label the axis and show whether an increase is good or bad. All current perf
// metrics are lower-is-better, but keeping it explicit (and per-metric) means a
// future higher-is-better metric (e.g. throughput) renders correctly.
type metricInfo struct {
	Unit          string `json:"unit"`
	LowerIsBetter bool   `json:"lower_is_better"`
}

// metricMeta maps each emitted metric to its unit + direction. The keys must
// match the metric strings used in explorerDatum rows.
var metricMeta = map[string]metricInfo{
	"ratio_to_anchor": {Unit: "× anchor", LowerIsBetter: true},
	"allocs_per_op":   {Unit: "allocs/op", LowerIsBetter: true},
	"bytes_per_op":    {Unit: "B/op", LowerIsBetter: true},
}

// explorerPayload is the shape of explorer.json: per-metric metadata plus the
// tidy rows.
type explorerPayload struct {
	Meta map[string]metricInfo `json:"meta"`
	Rows []explorerDatum       `json:"rows"`
}

// sig6 rounds to 6 significant figures — ample for a trend chart, and it keeps
// the embedded JSON small (raw float64 ratios stringify to ~17 digits).
func sig6(v float64) float64 {
	if v == 0 {
		return 0
	}
	d := math.Pow(10, 5-math.Floor(math.Log10(math.Abs(v))))
	return math.Round(v*d) / d
}

// buildExplorerData flattens the timeline into tidy rows for the Vega-Lite
// explorer. Each snapshot contributes, per benchmark, one row per non-zero
// metric (ratio_to_anchor / ns_per_op / allocs_per_op / bytes_per_op), tagged
// with the snapshot's CPU model — so the explorer plots a line per CPU.
func buildExplorerData(timeline []Snapshot) []byte {
	var rows []explorerDatum
	for _, s := range timeline {
		date := s.Baseline.CapturedAt
		cpu := s.Baseline.Machine.CPUModel
		if cpu == "" {
			cpu = "(unknown)"
		}
		for name, e := range s.Baseline.Benchmarks {
			bench := strings.TrimPrefix(name, modulePrefix)
			samples := e.Samples
			// add emits a row for one metric; sampleVal pulls that metric out of
			// a raw sample so we can shade the min/max run-to-run spread.
			add := func(metric string, v float64, sampleVal func(BenchmarkSample) float64) {
				if v == 0 {
					return
				}
				lo, hi, seen := v, v, false
				var pts []float64
				for _, sm := range samples {
					sv := sampleVal(sm)
					if sv == 0 {
						continue
					}
					pts = append(pts, sig6(sv))
					if !seen {
						lo, hi, seen = sv, sv, true
						continue
					}
					if sv < lo {
						lo = sv
					}
					if sv > hi {
						hi = sv
					}
				}
				// Drop samples that carry no spread: deterministic metrics
				// (allocs/op, bytes/op) repeat the identical value every run, so
				// their samples[] would just triple the payload with no signal.
				// The renderers fall back to the mean when samples is absent.
				distinct := false
				for _, p := range pts {
					if p != pts[0] {
						distinct = true
						break
					}
				}
				if !distinct {
					pts = nil
				}
				rows = append(rows, explorerDatum{
					Date: date, CPU: cpu, Bench: bench, Metric: metric,
					Value: sig6(v), Lo: sig6(lo), Hi: sig6(hi), Samples: pts,
				})
			}
			// ratio_to_anchor is the machine-portable timing metric (raw
			// ns_per_op is intentionally omitted — it's CPU-dependent and
			// redundant with the ratio); allocs/bytes are the deterministic
			// metrics this explorer was added to surface.
			add("ratio_to_anchor", e.RatioToAnchor, func(sm BenchmarkSample) float64 { return sm.RatioToAnchor })
			add("allocs_per_op", float64(e.AllocsPerOp), func(sm BenchmarkSample) float64 { return float64(sm.AllocsPerOp) })
			add("bytes_per_op", float64(e.BytesPerOp), func(sm BenchmarkSample) float64 { return float64(sm.BytesPerOp) })
		}
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Bench != rows[j].Bench {
			return rows[i].Bench < rows[j].Bench
		}
		if rows[i].Metric != rows[j].Metric {
			return rows[i].Metric < rows[j].Metric
		}
		return rows[i].Date < rows[j].Date
	})
	b, err := json.Marshal(explorerPayload{Meta: metricMeta, Rows: rows})
	if err != nil {
		return []byte(`{"meta":{},"rows":[]}`)
	}
	return b
}

func main() {
	var (
		baselinePath   = flag.String("baseline", "docs/perf/baseline.json", "current baseline JSON")
		historicalPath = flag.String("historical", "docs/perf/historical", "historical baseline directory")
		timelinePath   = flag.String("timeline", "docs/perf/timeline", "timeline snapshot directory")
		outPath        = flag.String("out", "docs/perf/index.html", "HTML output path")
		explorerOut    = flag.String("explorer-out", "", "path to write the explorer data JSON (default: explorer.json next to -out). The page fetches this file rather than inlining ~tens of thousands of rows.")
		explorerURL    = flag.String("explorer-url", "explorer.json", "URL the page fetches the explorer data from (relative to the page = same-origin; or an absolute https URL, e.g. raw.githubusercontent of the perf-data branch)")
		logoPath       = flag.String("logo", "meta/logo.svg", "logo SVG to embed")
		cpuFilter      = flag.String("cpu", "", "keep only timeline snapshots whose machine cpu_model contains this substring (CI runs land on ≥2 CPU tiers whose ratio_to_anchor doesn't normalize across them, so a mixed timeline zig-zags ~2x; filtering to one tier gives a clean series). Empty = all.")
		anchorName     = flag.String("anchor", "", "historical baseline to compare against, by file stem (e.g. 'v1.8.0'); empty = newest. Lets the page swap anchors — e.g. a fresh same-machine baseline for the modern suite vs the legacy v1.8.0 (Apple M3, pre-IR) reference.")
	)
	flag.Parse()

	current, err := loadBaseline(*baselinePath)
	if err != nil {
		die("load baseline: %v", err)
	}
	reference, referenceName, err := loadHistoricalAnchor(*historicalPath, *anchorName)
	if err != nil {
		die("load historical baseline: %v", err)
	}
	timeline, err := loadTimeline(*timelinePath, *historicalPath, current)
	if err != nil {
		die("load timeline: %v", err)
	}
	timeline = filterTimelineByCPU(timeline, *cpuFilter)
	logo, err := logoDataURI(*logoPath)
	if err != nil {
		die("load logo: %v", err)
	}

	page := buildPage(current, reference, referenceName, timeline, logo)
	page.ExplorerURL = *explorerURL
	html, err := renderPage(page)
	if err != nil {
		die("render page: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(*outPath), 0o755); err != nil {
		die("create output directory: %v", err)
	}
	if err := os.WriteFile(*outPath, html, 0o644); err != nil {
		die("write %s: %v", *outPath, err)
	}
	// Write the explorer data as a separate file the page fetches (keeps the
	// HTML small instead of inlining tens of thousands of rows).
	explorerPath := *explorerOut
	if explorerPath == "" {
		explorerPath = filepath.Join(filepath.Dir(*outPath), "explorer.json")
	}
	data := buildExplorerData(timeline)
	if err := os.WriteFile(explorerPath, data, 0o644); err != nil {
		die("write %s: %v", explorerPath, err)
	}
	fmt.Printf("wrote %s (%d benchmarks)\n  explorer data → %s (%d KB)\n",
		*outPath, len(current.Benchmarks), explorerPath, len(data)/1024)
}

func die(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "perf-page: "+format+"\n", args...)
	os.Exit(1)
}

// filterTimelineByCPU reports the CPU-tier mix of the timeline and, when sub is
// non-empty, keeps only snapshots whose machine cpu_model contains it. The mix
// matters because ratio_to_anchor does not normalize across CI CPU tiers (a
// trivial ~1ns anchor can't track each microarch's cache/memory/GC profile), so
// a mixed timeline zig-zags ~2x between tiers and a single point is unreadable.
// Filtering to one tier yields a clean, comparable series.
func filterTimelineByCPU(timeline []Snapshot, sub string) []Snapshot {
	counts := map[string]int{}
	for _, s := range timeline {
		cpu := s.Baseline.Machine.CPUModel
		if cpu == "" {
			cpu = "(unknown)"
		}
		counts[cpu]++
	}
	models := make([]string, 0, len(counts))
	for m := range counts {
		models = append(models, m)
	}
	sort.Slice(models, func(i, j int) bool { return counts[models[i]] > counts[models[j]] })
	fmt.Fprintf(os.Stderr, "timeline CPU tiers (%d snapshots):\n", len(timeline))
	for _, m := range models {
		fmt.Fprintf(os.Stderr, "  %3d  %s\n", counts[m], m)
	}
	if sub == "" {
		return timeline
	}
	want := strings.ToLower(sub)
	kept := make([]Snapshot, 0, len(timeline))
	for _, s := range timeline {
		if strings.Contains(strings.ToLower(s.Baseline.Machine.CPUModel), want) {
			kept = append(kept, s)
		}
	}
	fmt.Fprintf(os.Stderr, "cpu filter %q → kept %d of %d snapshots\n", sub, len(kept), len(timeline))
	return kept
}

// loadBaseline reads an arch-partitioned (v2) or legacy single-machine (v1)
// baseline file and returns ONE machine profile's data. With multiple profiles
// it currently selects the first by key (sorted) — the page is per-machine, and
// multi-machine rendering is a follow-up.
func loadBaseline(path string) (Baseline, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return Baseline{}, err
	}
	var probe struct {
		Version int `json:"version"`
	}
	if err := json.Unmarshal(b, &probe); err != nil {
		return Baseline{}, err
	}
	var multi MultiBaseline
	if probe.Version <= 1 {
		var v1 perfdata.LegacyBaselineV1
		if err := json.Unmarshal(b, &v1); err != nil {
			return Baseline{}, err
		}
		multi = MultiBaseline{
			Version:  2,
			Machines: map[string]Baseline{perfdata.MachineKey(v1.Machine): v1.ToMachineBaseline()},
		}
	} else if err := json.Unmarshal(b, &multi); err != nil {
		return Baseline{}, err
	}
	prof, ok := selectProfile(multi)
	if !ok || len(prof.Benchmarks) == 0 {
		return Baseline{}, fmt.Errorf("%s has no benchmarks", path)
	}
	return prof, nil
}

// selectProfile picks one machine profile from a multi-arch baseline (first by
// sorted key). Deterministic; multi-machine rendering can refine this later.
func selectProfile(b MultiBaseline) (Baseline, bool) {
	keys := make([]string, 0, len(b.Machines))
	for k := range b.Machines {
		keys = append(keys, k)
	}
	if len(keys) == 0 {
		return Baseline{}, false
	}
	sort.Strings(keys)
	return b.Machines[keys[0]], true
}

func loadLatestHistorical(dir string) (Baseline, string, error) {
	matches, err := filepath.Glob(filepath.Join(dir, "*.json"))
	if err != nil {
		return Baseline{}, "", err
	}
	if len(matches) == 0 {
		return Baseline{}, "", nil
	}

	type historical struct {
		path     string
		name     string
		baseline Baseline
		captured time.Time
	}
	var all []historical
	for _, match := range matches {
		baseline, err := loadBaseline(match)
		if err != nil {
			warnSkipBaseline(match, err)
			continue
		}
		captured, _ := time.Parse(time.RFC3339, baseline.CapturedAt)
		name := strings.TrimSuffix(filepath.Base(match), filepath.Ext(match))
		all = append(all, historical{
			path:     match,
			name:     name,
			baseline: baseline,
			captured: captured,
		})
	}
	if len(all) == 0 {
		return Baseline{}, "", nil
	}
	sort.Slice(all, func(i, j int) bool {
		if !all[i].captured.Equal(all[j].captured) {
			return all[i].captured.After(all[j].captured)
		}
		return all[i].path > all[j].path
	})
	return all[0].baseline, all[0].name, nil
}

// loadHistoricalAnchor selects the comparison baseline. With name empty it
// keeps the existing behavior (newest historical). With a name (file stem,
// e.g. "v1.8.0") it loads that specific baseline — so the page can anchor the
// modern suite to a fresh same-machine baseline while keeping v1.8.0 available
// as an explicit, labeled legacy reference rather than the silent default.
func loadHistoricalAnchor(dir, name string) (Baseline, string, error) {
	if name == "" {
		return loadLatestHistorical(dir)
	}
	path := filepath.Join(dir, name+".json")
	baseline, err := loadBaseline(path)
	if err != nil {
		return Baseline{}, "", fmt.Errorf("anchor %q (%s): %w", name, path, err)
	}
	return baseline, name, nil
}

func loadTimeline(timelineDir, historicalDir string, current Baseline) ([]Snapshot, error) {
	var snapshots []Snapshot
	matches, err := filepath.Glob(filepath.Join(timelineDir, "*.json"))
	if err != nil {
		return nil, err
	}
	for _, match := range matches {
		baseline, err := loadBaseline(match)
		if err != nil {
			warnSkipBaseline(match, err)
			continue
		}
		snapshots = append(snapshots, makeSnapshot(strings.TrimSuffix(filepath.Base(match), filepath.Ext(match)), baseline))
	}
	if len(snapshots) == 0 {
		historical, err := filepath.Glob(filepath.Join(historicalDir, "*.json"))
		if err != nil {
			return nil, err
		}
		for _, match := range historical {
			baseline, err := loadBaseline(match)
			if err != nil {
				warnSkipBaseline(match, err)
				continue
			}
			snapshots = append(snapshots, makeSnapshot(strings.TrimSuffix(filepath.Base(match), filepath.Ext(match)), baseline))
		}
		snapshots = append(snapshots, makeSnapshot("current-ratchet", current))
	}
	sort.Slice(snapshots, func(i, j int) bool {
		if !snapshots[i].Captured.Equal(snapshots[j].Captured) {
			return snapshots[i].Captured.Before(snapshots[j].Captured)
		}
		if snapshots[i].Baseline.CapturedAtSHA != snapshots[j].Baseline.CapturedAtSHA {
			return snapshots[i].Baseline.CapturedAtSHA < snapshots[j].Baseline.CapturedAtSHA
		}
		return snapshots[i].Name < snapshots[j].Name
	})
	return snapshots, nil
}

func warnSkipBaseline(path string, err error) {
	fmt.Fprintf(os.Stderr, "perf-page: warning: skipping %s: %v\n", path, err)
}

func makeSnapshot(name string, baseline Baseline) Snapshot {
	captured, _ := time.Parse(time.RFC3339, baseline.CapturedAt)
	return Snapshot{Name: name, Baseline: baseline, Captured: captured}
}

func logoDataURI(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	encoded := base64.StdEncoding.EncodeToString(b)
	return "data:image/svg+xml;base64," + encoded, nil
}

func buildPage(current, reference Baseline, referenceName string, timeline []Snapshot, logo string) PageData {
	rows := benchmarkRows(current)
	changes, summary := compare(current, reference)
	maxRatio := 0.0
	packageSet := map[string]struct{}{}
	for i := range rows {
		if rows[i].Ratio > maxRatio {
			maxRatio = rows[i].Ratio
		}
		packageSet[rows[i].Package] = struct{}{}
		if delta, ok := changes[rows[i].FullName]; ok {
			rows[i].Delta = &delta
		}
		if rows[i].AllocsPerOp == 0 {
			summary.ZeroAllocs++
		}
	}
	for i := range rows {
		rows[i].BarWidth = barWidth(rows[i].Ratio, maxRatio)
	}
	summary.BenchmarkCount = len(rows)
	summary.PackageCount = len(packageSet)

	recent := append([]BenchmarkRow(nil), rows...)
	sort.Slice(recent, func(i, j int) bool {
		ti, _ := time.Parse(time.RFC3339, recent[i].BestSinceAt)
		tj, _ := time.Parse(time.RFC3339, recent[j].BestSinceAt)
		if !ti.Equal(tj) {
			return ti.After(tj)
		}
		return recent[i].FullName < recent[j].FullName
	})
	if len(recent) > 8 {
		recent = recent[:8]
	}

	improvements, slowdowns := topChanges(current, reference, 8)
	return PageData{
		Title:             "Are we fast yet?",
		LogoDataURI:       template.URL(logo),
		Current:           current,
		ReferenceName:     referenceName,
		Summary:           summary,
		Timeline:          timeline,
		Charts:            buildCharts(timeline, reference, referenceName, defaultBudgetFraction),
		Rows:              rows,
		TopImprovements:   improvements,
		TopSlowdowns:      slowdowns,
		RecentlyTightened: recent,
	}
}

// defaultBudgetFraction mirrors bench-ratchet's regression budget (5%): the
// reference line plus this margin is the ceiling a benchmark must stay under.
const defaultBudgetFraction = 0.05

func buildCharts(timeline []Snapshot, reference Baseline, referenceName string, budget float64) []Chart {
	const (
		suite = "github.com/nooga/let-go/test.BenchmarkClojureTestSuite"
		ir    = "github.com/nooga/let-go/pkg/ir.BenchmarkIRCompile"
	)
	// Series carry candidate keys in priority order: the benchmark name has
	// grown a "[variant]" suffix over time (bytecode → ir_bytecode,
	// gogen_ir → aot_native). Matching the first key that exists keeps a
	// logical series continuous across the rename instead of silently
	// dropping to an empty chart.
	suiteSeries := []chartSeriesSpec{
		{label: "ir_bytecode", color: "#245c73", names: []string{suite + " [ir_bytecode]", suite + " [bytecode]"}},
		{label: "aot_native", color: "#167a48", names: []string{suite + " [aot_native]", suite + " [gogen_ir]"}},
	}
	specs := []struct {
		title    string
		subtitle string
		unit     string // used as-is for absolute charts; relative charts override it
		metric   func(BenchmarkEntry) float64
		sample   func(BenchmarkSample) float64
		format   func(float64) string
		series   []chartSeriesSpec
		refKeys  []string // keys to try in the reference baseline (release line)
		relative bool     // plot % vs reference (or first run) instead of raw values
	}{
		{
			title:    "End-to-end suite",
			subtitle: "Wall time relative to the release. Lower is better.",
			metric:   func(e BenchmarkEntry) float64 { return e.RatioToAnchor },
			sample:   func(s BenchmarkSample) float64 { return s.RatioToAnchor },
			format:   formatRatio,
			series:   suiteSeries,
			refKeys:  []string{suite, suite + " [ir_bytecode]", suite + " [bytecode]"},
			relative: true,
		},
		{
			title:    "IR compile",
			subtitle: "Compile time relative to the window start. Lower is better.",
			metric:   func(e BenchmarkEntry) float64 { return e.RatioToAnchor },
			sample:   func(s BenchmarkSample) float64 { return s.RatioToAnchor },
			format:   formatRatio,
			series: []chartSeriesSpec{
				{label: "bytecode", color: "#245c73", names: []string{ir + " [bytecode]"}},
				{label: "gogen_ir", color: "#167a48", names: []string{ir + " [gogen_ir]"}},
			},
			refKeys:  []string{ir, ir + " [bytecode]"},
			relative: true,
		},
		{
			title:    "Suite allocations",
			subtitle: "Allocations per op, both variants. Lower is better.",
			unit:     "allocs/op",
			metric:   func(e BenchmarkEntry) float64 { return float64(e.AllocsPerOp) },
			sample:   func(s BenchmarkSample) float64 { return float64(s.AllocsPerOp) },
			format:   formatCount,
			series:   suiteSeries,
		},
		{
			title:    "Suite memory",
			subtitle: "Heap bytes per op, both variants. Lower is better.",
			unit:     "B/op",
			metric:   func(e BenchmarkEntry) float64 { return float64(e.BytesPerOp) },
			sample:   func(s BenchmarkSample) float64 { return float64(s.BytesPerOp) },
			format:   formatBytes,
			series:   suiteSeries,
		},
	}

	charts := make([]Chart, 0, len(specs))
	for _, spec := range specs {
		refVal := 0.0
		if len(spec.refKeys) > 0 {
			refVal = lookupRef(reference, spec.refKeys, spec.metric)
		}
		refLabel := ""
		if refVal > 0 {
			refLabel = referenceName
		}
		unit := spec.unit
		if spec.relative {
			if refVal > 0 {
				unit = "% vs " + referenceName
			} else {
				unit = "% vs first run"
			}
		}
		chart := buildChart(timeline, spec.title, spec.subtitle, unit,
			spec.metric, spec.sample, spec.format, spec.series, refVal, refLabel, budget, spec.relative)
		if len(chart.Series) > 0 {
			charts = append(charts, chart)
		}
	}
	return charts
}

// lookupRef returns the first reference benchmark found among keys, as the
// chart metric. Zero means no reference (e.g. IRCompile predates the release).
func lookupRef(reference Baseline, keys []string, metric func(BenchmarkEntry) float64) float64 {
	for _, k := range keys {
		if e, ok := reference.Benchmarks[k]; ok {
			if v := metric(e); v > 0 {
				return v
			}
		}
	}
	return 0
}

type chartSeriesSpec struct {
	label string
	color string
	names []string // candidate keys, first match wins (tolerates variant renames)
}

func buildChart(timeline []Snapshot, title, subtitle, unit string,
	metric func(BenchmarkEntry) float64, sampleMetric func(BenchmarkSample) float64,
	format func(float64) string, specs []chartSeriesSpec,
	refVal float64, refLabel string, budget float64, relative bool) Chart {
	const (
		left   = 46.0
		right  = 18.0
		top    = 22.0
		bottom = 34.0
		width  = 520.0
		height = 210.0
	)
	plotW := width - left - right
	plotH := height - top - bottom

	// Pass 1: collect raw metric values. Value/Low/High hold raw numbers until
	// the display transform below; range tracking waits until the relative
	// basis is known.
	series := make([]ChartSeries, 0, len(specs))
	oldestRaw, oldestIdx := 0.0, math.MaxInt
	// Latest plotted value per series, for a worst-case status across all
	// series rather than whichever series sorts first (see chartStatus).
	var latest []seriesLatest
	for _, spec := range specs {
		var pts []ChartPoint
		for i, snap := range timeline {
			entry, ok := lookupEntry(snap.Baseline.Benchmarks, spec.names)
			if !ok {
				continue
			}
			value := metric(entry)
			if value <= 0 {
				continue
			}
			lo, hi, hasBand := sampleSpread(entry.Samples, sampleMetric)
			if !hasBand {
				lo, hi = value, value
			}
			pts = append(pts, ChartPoint{
				Index: i,
				Date:  formatDate(snap.Baseline.CapturedAt),
				SHA:   shortSHA(snap.Baseline.CapturedAtSHA),
				Value: value, Low: lo, High: hi, HasBand: hasBand,
			})
			if i < oldestIdx {
				oldestIdx, oldestRaw = i, value
			}
		}
		if len(pts) == 0 {
			continue
		}
		// pts are appended in ascending timeline order, so the last is latest.
		latest = append(latest, seriesLatest{label: spec.label, value: pts[len(pts)-1].Value})
		series = append(series, ChartSeries{Label: spec.label, Color: spec.color, Points: pts})
	}
	if len(series) == 0 {
		return Chart{}
	}

	// Inter-series gap (series[1] vs series[0]) per shared snapshot, captured
	// from RAW values before the display transform below overwrites them.
	type rawDelta struct {
		index int
		delta float64
	}
	var deltasRaw []rawDelta
	if len(series) == 2 {
		base := make(map[int]float64, len(series[0].Points))
		for _, p := range series[0].Points {
			base[p.Index] = p.Value
		}
		for _, p := range series[1].Points {
			if v0, ok := base[p.Index]; ok && v0 > 0 {
				deltasRaw = append(deltasRaw, rawDelta{p.Index, p.Value/v0 - 1})
			}
		}
	}

	// Relative mode plots every point as a fraction of a basis — the release
	// reference when present, else the oldest run in the window. This keeps the
	// baseline at a fixed on-scale coordinate (0) and turns the huge raw anchor
	// ratios into readable percentages.
	basis := 0.0
	if relative {
		if refVal > 0 {
			basis = refVal
		} else {
			basis = oldestRaw
		}
		if basis <= 0 {
			relative = false
		}
	}
	disp := func(raw float64) float64 {
		if relative {
			return raw/basis - 1
		}
		return raw
	}
	fmtv := func(raw float64) string {
		if relative {
			return formatPct(raw/basis - 1)
		}
		return format(raw)
	}

	yMin, yMax := math.Inf(1), math.Inf(-1)
	note := func(v float64) {
		if math.IsInf(v, 0) || math.IsNaN(v) {
			return
		}
		if v < yMin {
			yMin = v
		}
		if v > yMax {
			yMax = v
		}
	}

	hasRef := refVal > 0
	hasBudget := hasRef && budget > 0

	// Pass 2: raw → display, format, and track the y-range.
	for si := range series {
		for pi := range series[si].Points {
			p := &series[si].Points[pi]
			rawV, rawLo, rawHi := p.Value, p.Low, p.High
			p.Text = fmtv(rawV)
			if p.HasBand {
				p.Spread = fmtv(rawLo) + " .. " + fmtv(rawHi)
			}
			p.Value, p.Low, p.High = disp(rawV), disp(rawLo), disp(rawHi)
			note(p.Value)
			note(p.Low)
			note(p.High)
		}
	}

	if yMin == yMax {
		yMin *= 0.95
		yMax *= 1.05
		if yMin == yMax {
			yMin, yMax = 0, 1
		}
	}
	pad := (yMax - yMin) * 0.08
	yMin -= pad
	yMax += pad

	yOf := func(v float64) float64 { return top + ((yMax-v)/(yMax-yMin))*plotH }
	denom := float64(maxInt(len(timeline)-1, 1))
	xOf := func(index int) float64 { return left + (float64(index)/denom)*plotW }

	tickSet := map[int]struct{}{}
	for si := range series {
		for pi := range series[si].Points {
			p := &series[si].Points[pi]
			p.X = xOf(p.Index)
			p.Y = yOf(p.Value)
			p.LowY = yOf(p.Low)
			p.HighY = yOf(p.High)
			tickSet[p.Index] = struct{}{}
		}
		series[si].Path = brokenLinePath(series[si].Points)
		series[si].BandPath = bandPath(series[si].Points)
	}

	yMinText, yMaxText, axisUnit := axisLabels(yMin, yMax, unit, relative, format)
	chart := Chart{
		Title:    title,
		Subtitle: subtitle,
		Unit:     axisUnit,
		Series:   series,
		YMin:     yMin,
		YMax:     yMax,
		YMinText: yMinText,
		YMaxText: yMaxText,
		XTicks:   buildXTicks(timeline, tickSet, xOf),
	}

	if len(deltasRaw) > 0 {
		yAt := func(si, index int) (float64, bool) {
			for _, p := range series[si].Points {
				if p.Index == index {
					return p.Y, true
				}
			}
			return 0, false
		}
		// Thin to a readable number of labels on dense timelines (always keep
		// the latest); the lines themselves still show every point. A min
		// horizontal gap then prevents adjacent labels (and the strided run
		// meeting the forced last) from overprinting.
		const maxGapLabels = 6
		const minGapPx = 56.0
		stride := 1
		if len(deltasRaw) > maxGapLabels {
			stride = (len(deltasRaw) + maxGapLabels - 1) / maxGapLabels
		}
		for n, g := range deltasRaw {
			if n%stride != 0 && n != len(deltasRaw)-1 {
				continue
			}
			x := xOf(g.index)
			if k := len(chart.Deltas); k > 0 && x-chart.Deltas[k-1].X < minGapPx {
				chart.Deltas = chart.Deltas[:k-1] // drop the crowded predecessor, keep the later one
			}
			y0, _ := yAt(0, g.index)
			y1, _ := yAt(1, g.index)
			upper, lower := math.Min(y0, y1), math.Max(y0, y1)
			labelY := upper - 9 // above the higher point, with clearance
			if upper < top+16 {
				labelY = lower + 15 // too close to the top axis — drop below the pair
			}
			chart.Deltas = append(chart.Deltas, ChartDelta{X: x, Y: labelY, Text: formatPct(g.delta)})
		}
		chart.DeltaCaption = "Δ " + series[1].Label + " vs " + series[0].Label
	}

	switch {
	case hasRef:
		// The reference and budget are evaluated in display space. We never
		// stretch the axis to include them: when off-scale (the usual case in
		// relative mode, where 0% sits above an all-negative data range) they
		// show as a legend marker; when the data rises within range of the
		// ceiling, the line and budget band draw in place automatically.
		refDisp := disp(refVal)
		chart.HasRef = true
		chart.RefLabel = refLabel
		// Keep the absolute magnitude of the 0 baseline visible; the axis unit
		// ("% vs <ref>") already says it's the 0% line, so the chip stays terse.
		chart.RefLegend = refLabel + " = " + format(refVal)
		switch {
		case refDisp >= yMin && refDisp <= yMax:
			chart.ShowRefLine = true
			chart.RefY = yOf(refDisp)
			if hasBudget {
				chart.HasBudget = true
				chart.BudgetY = math.Max(yOf(disp(refVal*(1+budget))), top)
				chart.BudgetText = fmt.Sprintf("±%.0f%%", budget*100)
			}
		case refDisp > yMax:
			chart.RefLegend += " ↑" // 0% baseline is above the plotted range
		default:
			chart.RefLegend += " ↓"
		}
	case relative:
		// No release reference: the baseline is the first run in the window,
		// which is a real plotted point at 0%, so the line stays in range.
		chart.HasRef = true
		chart.ShowRefLine = true
		chart.RefY = yOf(0)
		chart.RefLabel = "first run"
		chart.RefLegend = "first run = " + format(basis)
	}

	chart.Status, chart.StatusClass = chartStatus(latest, refVal, refLabel, budget)
	return chart
}

// axisLabels formats the y-axis bound labels and resolves the axis unit. For
// absolute charts it picks ONE scale from the larger bound and emits bare
// numbers (e.g. "1.19" / "0.538" with the unit "GiB/op" carried in the meta
// line) — long per-label strings like "1.185 GiB" overflow the narrow left
// margin and pick inconsistent units (537.9 MiB vs 1.185 GiB) between bounds.
func axisLabels(yMin, yMax float64, unit string, relative bool, format func(float64) string) (minText, maxText, axisUnit string) {
	switch {
	case relative:
		return formatPct(yMin), formatPct(yMax), unit
	case unit == "B/op":
		f, u := pickBytesScale(yMax)
		return trimAxisNum(yMin / f), trimAxisNum(yMax / f), u + "/op"
	case strings.Contains(unit, "allocs"):
		f, u := pickCountScale(yMax)
		return trimAxisNum(yMin / f), trimAxisNum(yMax / f), strings.TrimSpace(u + " allocs/op")
	default:
		return format(yMin), format(yMax), unit
	}
}

func pickBytesScale(v float64) (float64, string) {
	switch {
	case v >= 1<<30:
		return 1 << 30, "GiB"
	case v >= 1<<20:
		return 1 << 20, "MiB"
	case v >= 1<<10:
		return 1 << 10, "KiB"
	default:
		return 1, "B"
	}
}

func pickCountScale(v float64) (float64, string) {
	switch {
	case v >= 1e9:
		return 1e9, "G"
	case v >= 1e6:
		return 1e6, "M"
	case v >= 1e3:
		return 1e3, "k"
	default:
		return 1, ""
	}
}

// trimAxisNum renders ~3 significant figures, trailing zeros trimmed.
func trimAxisNum(v float64) string { return fmt.Sprintf("%.3g", v) }

func lookupEntry(benchmarks map[string]BenchmarkEntry, names []string) (BenchmarkEntry, bool) {
	for _, n := range names {
		if e, ok := benchmarks[n]; ok {
			return e, true
		}
	}
	return BenchmarkEntry{}, false
}

func sampleSpread(samples []BenchmarkSample, metric func(BenchmarkSample) float64) (lo, hi float64, ok bool) {
	lo, hi = math.Inf(1), math.Inf(-1)
	n := 0
	for _, s := range samples {
		v := metric(s)
		if v <= 0 {
			continue
		}
		n++
		if v < lo {
			lo = v
		}
		if v > hi {
			hi = v
		}
	}
	if n < 2 || lo == hi {
		return 0, 0, false
	}
	return lo, hi, true
}

// brokenLinePath restarts the subpath (M) wherever consecutive plotted points
// skip a snapshot, so a gap in the data reads as a gap instead of a straight
// line bridging across missing runs.
func brokenLinePath(points []ChartPoint) string {
	parts := make([]string, 0, len(points))
	for i := range points {
		cmd := "L"
		if i == 0 || points[i].Index != points[i-1].Index+1 {
			cmd = "M"
		}
		parts = append(parts, fmt.Sprintf("%s %.2f %.2f", cmd, points[i].X, points[i].Y))
	}
	return strings.Join(parts, " ")
}

// bandPath builds the sample min..max envelope as one closed polygon per
// contiguous run of band-bearing points (top edge left-to-right, bottom edge
// back). Runs break on the same gaps as the line.
func bandPath(points []ChartPoint) string {
	var b strings.Builder
	i := 0
	for i < len(points) {
		if !points[i].HasBand {
			i++
			continue
		}
		j := i
		for j+1 < len(points) && points[j+1].HasBand && points[j+1].Index == points[j].Index+1 {
			j++
		}
		if j == i {
			i++
			continue // a lone banded point has no width to fill
		}
		for k := i; k <= j; k++ {
			cmd := "L"
			if k == i {
				cmd = "M"
			}
			fmt.Fprintf(&b, "%s %.2f %.2f ", cmd, points[k].X, points[k].HighY)
		}
		for k := j; k >= i; k-- {
			fmt.Fprintf(&b, "L %.2f %.2f ", points[k].X, points[k].LowY)
		}
		b.WriteString("Z ")
		i = j + 1
	}
	return strings.TrimSpace(b.String())
}

// buildXTicks labels snapshot positions with a short capture date ("Jun 04"),
// thinned to at most 8 so a long timeline doesn't crowd the axis. The exact
// timestamp + SHA stay in each point's hover tooltip — dates are low-entropy
// and scannable as axis labels; SHAs are not.
func buildXTicks(timeline []Snapshot, present map[int]struct{}, xOf func(int) float64) []ChartXTick {
	idx := make([]int, 0, len(present))
	for i := range present {
		idx = append(idx, i)
	}
	sort.Ints(idx)
	if len(idx) == 0 {
		return nil
	}
	stride := 1
	if len(idx) > 8 {
		stride = (len(idx) + 7) / 8
	}
	// A "Jun 04" label at ~8.5px needs ~44px of clearance so adjacent ticks
	// (and the strided run meeting the forced last) don't collide.
	const minGapPx = 44.0
	last := idx[len(idx)-1]
	ticks := make([]ChartXTick, 0, 8)
	add := func(i int) {
		x := xOf(i)
		label := tickDate(timeline[i].Baseline.CapturedAt)
		// Drop the previous tick if the new one would crowd it, or if it
		// repeats the same date (e.g. two same-day snapshots); keep the later.
		if n := len(ticks); n > 0 && (x-ticks[n-1].X < minGapPx || ticks[n-1].Label == label) {
			ticks = ticks[:n-1]
		}
		ticks = append(ticks, ChartXTick{X: x, Label: label})
	}
	for n, i := range idx {
		if i == last || n%stride == 0 {
			add(i)
		}
	}
	return ticks
}

// tickDate renders a capture timestamp as a short, scannable axis label.
func tickDate(captured string) string {
	t, err := time.Parse(time.RFC3339, captured)
	if err != nil {
		return ""
	}
	return t.UTC().Format("Jan 02")
}

// seriesLatest is a series' most recent plotted value, kept so chartStatus can
// pick the worst across series instead of trusting whichever sorts first.
type seriesLatest struct {
	label string
	value float64
}

// chartStatus summarizes the latest points against the reference line. With
// multiple series it reports the worst (highest, since lower is better) latest
// value and names the series, so the single status pill can't silently hide a
// regression in a non-first series.
func chartStatus(latest []seriesLatest, refVal float64, refLabel string, budget float64) (string, string) {
	if refVal <= 0 {
		return "no release reference", "none"
	}
	if len(latest) == 0 {
		return "no recent data", "none"
	}
	worst := latest[0]
	for _, s := range latest[1:] {
		if s.value > worst.value {
			worst = s
		}
	}
	delta := worst.value/refVal - 1
	pctText := fmt.Sprintf("%.0f%%", math.Abs(delta)*100)
	suffix := ""
	if len(latest) > 1 {
		suffix = " (" + worst.label + ")"
	}
	switch {
	case delta > budget:
		return pctText + " over " + refLabel + suffix, "bad"
	case delta < -budget:
		return pctText + " under " + refLabel + suffix, "good"
	default:
		return "within budget vs " + refLabel, "flat"
	}
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func benchmarkRows(baseline Baseline) []BenchmarkRow {
	rows := make([]BenchmarkRow, 0, len(baseline.Benchmarks))
	for fullName, entry := range baseline.Benchmarks {
		pkg, name := splitBenchmarkName(fullName)
		rows = append(rows, BenchmarkRow{
			FullName:     fullName,
			Package:      pkg,
			Name:         name,
			NSPerOp:      entry.NSPerOp,
			AllocsPerOp:  float64(entry.AllocsPerOp),
			BytesPerOp:   float64(entry.BytesPerOp),
			Ratio:        entry.RatioToAnchor,
			BestSinceSHA: entry.BestSinceSHA,
			BestSinceAt:  entry.BestSinceAt,
		})
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Package != rows[j].Package {
			return rows[i].Package < rows[j].Package
		}
		return rows[i].Name < rows[j].Name
	})
	return rows
}

func splitBenchmarkName(fullName string) (string, string) {
	name := strings.TrimPrefix(fullName, modulePrefix)
	idx := strings.Index(name, ".Benchmark")
	if idx < 0 {
		return "", name
	}
	return name[:idx], name[idx+1:]
}

func compare(current, reference Baseline) (map[string]float64, Summary) {
	changes := make(map[string]float64)
	if len(reference.Benchmarks) == 0 {
		return changes, Summary{New: len(current.Benchmarks)}
	}
	var deltas []float64
	var summary Summary
	for name, cur := range current.Benchmarks {
		ref, ok := reference.Benchmarks[name]
		if !ok {
			summary.New++
			continue
		}
		if ref.RatioToAnchor == 0 {
			continue
		}
		delta := cur.RatioToAnchor/ref.RatioToAnchor - 1
		changes[name] = delta
		deltas = append(deltas, delta)
		summary.Common++
		if delta <= -0.05 {
			summary.Faster++
		} else if delta >= 0.05 {
			summary.Slower++
		}
	}
	for name := range reference.Benchmarks {
		if _, ok := current.Benchmarks[name]; !ok {
			summary.Missing++
		}
	}
	summary.MedianDelta = median(deltas)
	return changes, summary
}

func topChanges(current, reference Baseline, limit int) ([]ChangeRow, []ChangeRow) {
	if len(reference.Benchmarks) == 0 {
		return nil, nil
	}
	var rows []ChangeRow
	for _, row := range benchmarkRows(current) {
		cur := current.Benchmarks[row.FullName]
		ref, ok := reference.Benchmarks[row.FullName]
		if !ok || ref.RatioToAnchor == 0 {
			continue
		}
		delta := cur.RatioToAnchor/ref.RatioToAnchor - 1
		row.Delta = &delta
		rows = append(rows, ChangeRow{
			BenchmarkRow: row,
			OldRatio:     ref.RatioToAnchor,
			NewRatio:     cur.RatioToAnchor,
		})
	}

	improvements := append([]ChangeRow(nil), rows...)
	sort.Slice(improvements, func(i, j int) bool {
		return *improvements[i].Delta < *improvements[j].Delta
	})
	slowdowns := append([]ChangeRow(nil), rows...)
	sort.Slice(slowdowns, func(i, j int) bool {
		return *slowdowns[i].Delta > *slowdowns[j].Delta
	})

	return takeSignificant(improvements, limit, true), takeSignificant(slowdowns, limit, false)
}

func takeSignificant(rows []ChangeRow, limit int, faster bool) []ChangeRow {
	out := make([]ChangeRow, 0, limit)
	for _, row := range rows {
		if row.Delta == nil {
			continue
		}
		if faster && *row.Delta >= -0.01 {
			continue
		}
		if !faster && *row.Delta <= 0.01 {
			continue
		}
		out = append(out, row)
		if len(out) == limit {
			return out
		}
	}
	return out
}

func median(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	values = append([]float64(nil), values...)
	sort.Float64s(values)
	mid := len(values) / 2
	if len(values)%2 == 1 {
		return values[mid]
	}
	return (values[mid-1] + values[mid]) / 2
}

func barWidth(ratio, maxRatio float64) float64 {
	if maxRatio <= 0 || ratio <= 0 {
		return 0
	}
	return math.Log1p(ratio) / math.Log1p(maxRatio) * 100
}

func renderPage(page PageData) ([]byte, error) {
	funcs := template.FuncMap{
		"date":       formatDate,
		"shortSHA":   shortSHA,
		"ns":         formatNS,
		"bytes":      formatBytes,
		"count":      formatCount,
		"ratio":      formatRatio,
		"pct":        formatPct,
		"deltaClass": deltaClass,
		"deltaText":  deltaText,
		"bar":        formatBar,
		"sub":        func(a, b float64) float64 { return a - b },
	}
	tmpl, err := template.New("page").Funcs(funcs).Parse(pageTemplate)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, page); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func formatDate(value string) string {
	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return value
	}
	return t.UTC().Format("2006-01-02 15:04 UTC")
}

func shortSHA(value string) string {
	if len(value) <= 12 {
		return value
	}
	return value[:12]
}

func formatNS(value float64) string {
	return formatScaledUnit(value, "ns")
}

func formatBytes(value float64) string {
	return formatScaledUnit(value, "B")
}

func formatCount(value float64) string {
	number, prefix := splitBenchunitScale(benchunit.Scale(value, benchunit.Decimal))
	return trimScaledNumber(number) + prefix
}

func formatScaledUnit(value float64, unit string) string {
	tidiedValue, tidiedUnit := benchunit.Tidy(value, unit)
	number, prefix := splitBenchunitScale(benchunit.Scale(tidiedValue, benchunit.ClassOf(tidiedUnit)))
	return trimScaledNumber(number) + " " + displayUnit(prefix, tidiedUnit)
}

func splitBenchunitScale(value string) (string, string) {
	idx := len(value)
	for idx > 0 {
		r, size := utf8.DecodeLastRuneInString(value[:idx])
		if r == '.' || r == '-' || r == '+' || unicode.IsDigit(r) {
			break
		}
		idx -= size
	}
	return value[:idx], value[idx:]
}

func trimScaledNumber(value string) string {
	if !strings.Contains(value, ".") {
		return value
	}
	value = strings.TrimRight(value, "0")
	value = strings.TrimRight(value, ".")
	if value == "-0" {
		return "0"
	}
	return value
}

func displayUnit(prefix, unit string) string {
	switch unit {
	case "sec":
		switch prefix {
		case "n":
			return "ns"
		case "µ":
			return "us"
		case "m":
			return "ms"
		case "":
			return "s"
		default:
			return prefix + "s"
		}
	case "B":
		return prefix + "B"
	default:
		return prefix + unit
	}
}

func formatRatio(value float64) string {
	if math.Abs(value) >= 1_000 {
		return formatCompact(value)
	}
	if value >= 10 {
		return fmt.Sprintf("%.1f", value)
	}
	return fmt.Sprintf("%.2f", value)
}

func formatCompact(value float64) string {
	abs := math.Abs(value)
	switch {
	case abs >= 1_000_000_000:
		return fmt.Sprintf("%.2fB", value/1_000_000_000)
	case abs >= 1_000_000:
		return fmt.Sprintf("%.2fM", value/1_000_000)
	case abs >= 1_000:
		return fmt.Sprintf("%.1fk", value/1_000)
	default:
		return fmt.Sprintf("%.0f", value)
	}
}

func formatPct(value float64) string {
	return fmt.Sprintf("%+.1f%%", value*100)
}

func deltaClass(value *float64) string {
	if value == nil || math.Abs(*value) < 0.0005 {
		return "muted" // not in reference, or unchanged — don't color it
	}
	if *value <= -0.05 {
		return "good"
	}
	if *value >= 0.05 {
		return "bad"
	}
	return "flat"
}

func deltaText(value *float64) string {
	if value == nil {
		return "new"
	}
	if math.Abs(*value) < 0.0005 {
		return "—" // unchanged vs reference; "+0.0%" reads like a regression
	}
	return formatPct(*value)
}

func formatBar(value float64) string {
	return fmt.Sprintf("%.2f", value)
}

const pageTemplate = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{{.Title}} - let-go perf</title>
  <style>
    :root {
      color-scheme: light;
      --bg: #f7f7f4;
      --paper: #ffffff;
      --ink: #171717;
      --muted: #65645f;
      --line: #deddd6;
      --soft: #eeede7;
      --green: #167a48;
      --green-bg: #e4f3ea;
      --red: #aa2e2e;
      --red-bg: #f8e3e0;
      --amber: #8b5e12;
      --amber-bg: #f4ead2;
      --accent: #245c73;
      --shadow: 0 18px 50px rgba(23, 23, 23, 0.08);
    }
    * { box-sizing: border-box; }
    body {
      margin: 0;
      background: var(--bg);
      color: var(--ink);
      font: 15px/1.5 ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
      letter-spacing: 0;
    }
    .wrap {
      width: min(1180px, calc(100% - 32px));
      margin: 0 auto;
    }
    header {
      padding: 34px 0 24px;
      border-bottom: 1px solid var(--line);
      background:
        linear-gradient(90deg, rgba(36, 92, 115, 0.12), transparent 42%),
        linear-gradient(180deg, #ffffff, var(--bg));
    }
    .topline {
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 18px;
      margin-bottom: 28px;
    }
    .brand {
      display: flex;
      align-items: center;
      gap: 12px;
      min-width: 0;
      font-weight: 750;
      letter-spacing: 0;
    }
    .brand img {
      width: 42px;
      height: 42px;
      display: block;
    }
    .brand span {
      white-space: nowrap;
    }
    .links {
      display: flex;
      gap: 10px;
      flex-wrap: wrap;
      justify-content: flex-end;
    }
    .links a {
      color: var(--ink);
      text-decoration: none;
      border: 1px solid var(--line);
      background: rgba(255, 255, 255, 0.7);
      border-radius: 7px;
      padding: 7px 10px;
      font-size: 13px;
      font-weight: 650;
    }
    h1 {
      margin: 0;
      font-size: clamp(42px, 7vw, 86px);
      line-height: 0.94;
      letter-spacing: 0;
      max-width: 820px;
    }
    .lede {
      margin: 18px 0 0;
      max-width: 760px;
      color: var(--muted);
      font-size: 18px;
    }
    .meta {
      margin-top: 20px;
      display: flex;
      flex-wrap: wrap;
      gap: 8px;
    }
    .chip {
      display: inline-flex;
      align-items: center;
      min-height: 30px;
      padding: 5px 9px;
      border-radius: 7px;
      border: 1px solid var(--line);
      background: rgba(255, 255, 255, 0.72);
      color: var(--muted);
      font-size: 13px;
      font-weight: 620;
    }
    main { padding: 28px 0 54px; }
    .cards {
      display: grid;
      grid-template-columns: repeat(4, minmax(0, 1fr));
      gap: 12px;
      margin-bottom: 30px;
    }
    .card {
      background: var(--paper);
      border: 1px solid var(--line);
      border-radius: 8px;
      box-shadow: var(--shadow);
      padding: 16px;
      min-height: 116px;
    }
    .label {
      margin: 0 0 8px;
      color: var(--muted);
      font-size: 12px;
      font-weight: 760;
      text-transform: uppercase;
      letter-spacing: 0.08em;
    }
    .value {
      margin: 0;
      font-size: 32px;
      line-height: 1;
      font-weight: 820;
      letter-spacing: 0;
    }
    .note {
      margin: 9px 0 0;
      color: var(--muted);
      font-size: 13px;
    }
    section {
      margin-top: 34px;
    }
    .section-head {
      display: flex;
      align-items: end;
      justify-content: space-between;
      gap: 20px;
      margin-bottom: 12px;
    }
    h2 {
      margin: 0;
      font-size: 24px;
      line-height: 1.15;
      letter-spacing: 0;
    }
    .section-head p {
      margin: 0;
      color: var(--muted);
      font-size: 13px;
    }
    .grid-2 {
      display: grid;
      grid-template-columns: repeat(2, minmax(0, 1fr));
      gap: 12px;
    }
    .chart-grid {
      display: grid;
      grid-template-columns: repeat(2, minmax(0, 1fr));
      gap: 12px;
    }
    .chart {
      background: var(--paper);
      border: 1px solid var(--line);
      border-radius: 8px;
      box-shadow: var(--shadow);
      padding: 14px;
      min-width: 0;
    }
    .chart h3 {
      margin: 0;
      font-size: 16px;
      line-height: 1.2;
      letter-spacing: 0;
    }
    .chart p {
      margin: 5px 0 12px;
      color: var(--muted);
      font-size: 12px;
    }
    /* Axis-unit + delta-descriptor line, between subtitle and chart. */
    /* Footer row under the chart: axis-unit + Δ descriptor on the left,
       series/reference legend on the right. */
    .chart-foot {
      display: flex;
      justify-content: space-between;
      align-items: baseline;
      gap: 6px 16px;
      flex-wrap: wrap;
      margin-top: 10px;
    }
    .chart-meta {
      margin: 0;
      color: var(--muted);
      font-size: 11.5px;
      font-style: italic;
      font-variant-numeric: tabular-nums;
    }
    .chart-unit { font-weight: 650; }
    .gap-delta { opacity: 0.9; }
    .chart svg {
      width: 100%;
      height: auto;
      display: block;
      overflow: visible;
    }
    .axis {
      stroke: var(--line);
      stroke-width: 1;
    }
    .axis-label {
      fill: var(--muted);
      font-size: 11px;
      font-variant-numeric: tabular-nums;
    }
    .chart-line {
      fill: none;
      stroke-width: 2.5;
      stroke-linecap: round;
      stroke-linejoin: round;
    }
    .point {
      stroke: var(--paper);
      stroke-width: 1.6;
    }
    .chart-band {
      opacity: 0.16;
      stroke: none;
    }
    .ref-line {
      stroke: var(--ink);
      stroke-width: 1;
      stroke-dasharray: 4 3;
      opacity: 0.5;
    }
    .budget-band {
      fill: var(--red);
      opacity: 0.08;
    }
    .tick {
      stroke: var(--line);
      stroke-width: 1;
    }
    .tick-label {
      fill: var(--muted);
      font-size: 9.5px;
      text-anchor: middle;
      font-variant-numeric: tabular-nums;
    }
    .gap-label {
      fill: var(--muted);
      font-size: 8.5px;
      font-weight: 700;
      text-anchor: middle;
      font-variant-numeric: tabular-nums;
    }
    .chart-head {
      display: flex;
      align-items: baseline;
      justify-content: space-between;
      gap: 10px;
    }
    .status {
      font-size: 11px;
      font-weight: 720;
      padding: 2px 8px;
      border-radius: 999px;
      white-space: nowrap;
    }
    .status.good { color: var(--green); background: var(--green-bg); }
    .status.bad { color: var(--red); background: var(--red-bg); }
    .status.flat { color: var(--amber); background: var(--amber-bg); }
    .legend {
      display: flex;
      column-gap: 14px;
      row-gap: 5px;
      flex-wrap: wrap;
      align-items: center;
      justify-content: flex-end;
      margin-left: auto; /* keep the legend right-aligned even when it wraps below the meta */
      color: var(--muted);
      font-size: 11.5px;
      font-weight: 650;
    }
    .legend span {
      display: inline-flex;
      align-items: center;
      gap: 6px;
    }
    .swatch {
      width: 16px;
      height: 3px;
      border-radius: 999px;
      background: var(--series);
      display: inline-block;
    }
    .swatch.dash {
      background: none;
      border-top: 1px dashed var(--ink);
      opacity: 0.6;
      height: 0;
    }
    .swatch.budget {
      height: 10px;
      border-radius: 2px;
      background: var(--red);
      opacity: 0.18;
    }
    table {
      width: 100%;
      border-collapse: collapse;
      background: var(--paper);
      border: 1px solid var(--line);
      border-radius: 8px;
      overflow: hidden;
      box-shadow: var(--shadow);
    }
    th, td {
      padding: 10px 12px;
      border-bottom: 1px solid var(--soft);
      text-align: left;
      vertical-align: middle;
    }
    th {
      color: var(--muted);
      font-size: 12px;
      font-weight: 780;
      text-transform: uppercase;
      letter-spacing: 0.08em;
      background: #fbfbf9;
    }
    tr:last-child td { border-bottom: 0; }
    .bench {
      min-width: 260px;
      font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
      font-size: 13px;
      line-height: 1.35;
      overflow-wrap: anywhere;
    }
    .pkg {
      display: block;
      color: var(--muted);
      font-family: ui-sans-serif, system-ui, sans-serif;
      font-size: 12px;
      margin-bottom: 2px;
    }
    .num {
      white-space: nowrap;
      font-variant-numeric: tabular-nums;
    }
    .delta {
      display: inline-flex;
      justify-content: center;
      min-width: 72px;
      padding: 4px 7px;
      border-radius: 6px;
      font-weight: 760;
      font-variant-numeric: tabular-nums;
    }
    .good { color: var(--green); background: var(--green-bg); }
    .bad { color: var(--red); background: var(--red-bg); }
    .flat { color: var(--amber); background: var(--amber-bg); }
    .muted { color: var(--muted); background: var(--soft); }
    .barcell {
      min-width: 170px;
    }
    .bartrack {
      height: 9px;
      border-radius: 999px;
      background: var(--soft);
      overflow: hidden;
    }
    .barfill {
      display: block;
      height: 100%;
      width: calc(var(--w) * 1%);
      min-width: 3px;
      background: linear-gradient(90deg, var(--accent), #6a9c7d);
      border-radius: inherit;
    }
    .empty {
      padding: 18px;
      border: 1px solid var(--line);
      border-radius: 8px;
      background: var(--paper);
      color: var(--muted);
    }
    footer {
      padding: 22px 0 38px;
      color: var(--muted);
      border-top: 1px solid var(--line);
      font-size: 13px;
    }
    @media (max-width: 860px) {
      .cards, .grid-2, .chart-grid { grid-template-columns: 1fr; }
      .section-head { display: block; }
      .section-head p { margin-top: 6px; }
      table { display: block; overflow-x: auto; }
      .topline { align-items: flex-start; }
    }
    @media (max-width: 560px) {
      .wrap { width: min(100% - 22px, 1180px); }
      header { padding-top: 20px; }
      .topline { flex-direction: column; }
      .links { justify-content: flex-start; }
      h1 { font-size: 42px; }
      .lede { font-size: 16px; }
      th, td { padding: 9px; }
    }
    .explorer-controls { display: flex; gap: 1rem; flex-wrap: wrap; align-items: center; margin-bottom: 0.75rem; }
    .explorer-controls label { font-size: 0.8rem; color: var(--muted); display: inline-flex; align-items: center; gap: 0.4rem; }
    .explorer-controls select { font: inherit; font-size: 0.8rem; padding: 0.25rem 0.45rem; border: 1px solid rgba(0,0,0,0.15); border-radius: 6px; background: var(--paper); color: var(--ink); max-width: 420px; }
    .explorer-controls input[type="search"] { font: inherit; font-size: 0.8rem; padding: 0.25rem 0.45rem; border: 1px solid rgba(0,0,0,0.15); border-radius: 6px; background: var(--paper); color: var(--ink); min-width: 180px; }
    .explorer-controls label.cb { gap: 0.35rem; cursor: pointer; }
    .explorer-controls label.cb input { cursor: pointer; }
    .spark-count { font-size: 0.78rem; color: var(--muted); margin-left: auto; font-variant-numeric: tabular-nums; }
    .explorer-chart { width: 100%; overflow-x: auto; }
    .explorer-chart figure { margin: 0; }
    .spark-table-wrap { max-height: 620px; overflow: auto; border: 1px solid rgba(0,0,0,0.08); border-radius: 10px; background: var(--paper); }
    table.spark-table { width: 100%; border-collapse: collapse; font-size: 0.82rem; }
    table.spark-table th { position: sticky; top: 0; background: var(--paper); text-align: left; padding: 0.5rem 0.75rem; color: var(--muted); font-weight: 640; border-bottom: 1px solid rgba(0,0,0,0.1); z-index: 1; white-space: nowrap; }
    table.spark-table td { padding: 0.25rem 0.75rem; border-bottom: 1px solid rgba(0,0,0,0.045); white-space: nowrap; }
    table.spark-table td.spark { padding: 0.1rem 0.5rem; line-height: 0; }
    table.spark-table th.num, table.spark-table td.num { text-align: right; font-variant-numeric: tabular-nums; }
    table.spark-table td.num { color: var(--muted); }
    table.spark-table tr:hover td { background: rgba(36,92,115,0.045); }
    table.spark-table .bench { color: var(--ink); max-width: 300px; overflow: hidden; text-overflow: ellipsis; }
    table.spark-table .delta { min-width: 0; }
    table.spark-table td.scales { white-space: normal; text-align: right; }
    .scalepill { display: inline-flex; gap: 0.25rem; align-items: baseline; margin-left: 0.4rem; font-size: 0.72rem; padding: 1px 5px; border-radius: 5px; font-variant-numeric: tabular-nums; }
    .scalepill b { font-weight: 700; opacity: 0.65; }
    .scalepill.good { color: var(--green); background: var(--green-bg); }
    .scalepill.bad { color: var(--red); background: var(--red-bg); }
    .spark-tip { position: fixed; z-index: 50; pointer-events: none; background: var(--ink); color: #fff; padding: 6px 9px; border-radius: 7px; font-size: 0.72rem; line-height: 1.4; box-shadow: 0 4px 16px rgba(0,0,0,0.28); max-width: 280px; opacity: 0; transition: opacity 0.08s; }
    .spark-tip .tip-h { font-weight: 700; }
    .spark-tip .tip-sub { opacity: 0.65; margin-bottom: 3px; }
    .spark-tip .tip-r { display: flex; gap: 10px; align-items: baseline; font-variant-numeric: tabular-nums; }
    .spark-tip .tip-r b { min-width: 26px; opacity: 0.65; font-weight: 700; }
    .spark-tip .tip-r span:nth-of-type(1) { margin-left: auto; }
    .spark-tip .tip-r .good { color: #7fdca4; }
    .spark-tip .tip-r .bad { color: #f3a3a3; }
    .spark-tip .tip-r.muted { opacity: 0.6; }
  </style>
</head>
<body>
  <header>
    <div class="wrap">
      <div class="topline">
        <div class="brand">
          {{if .LogoDataURI}}<img alt="" src="{{.LogoDataURI}}">{{end}}
          <span>let-go perf</span>
        </div>
        <nav class="links" aria-label="Links">
          <a href="../">WASM repl</a>
          <a href="https://github.com/nooga/let-go">GitHub</a>
          <a href="https://github.com/nooga/let-go/blob/main/docs/perf/ratchet.md">Ratchet docs</a>
        </nav>
      </div>
      <h1>{{.Title}}</h1>
      <p class="lede">Committed benchmark ratchet data, rendered as a static page. Ratios are normalized to {{.Current.Anchor.Name}}, so lower is better and cross-machine drift is less noisy.</p>
      <div class="meta">
        <span class="chip">captured {{date .Current.CapturedAt}}</span>
        <span class="chip">sha {{shortSHA .Current.CapturedAtSHA}}</span>
        <span class="chip">{{.Current.Machine.CPUModel}}</span>
        <span class="chip">{{.Current.Machine.GoVersion}}</span>
      </div>
    </div>
  </header>

  <main class="wrap">
    <div class="cards" aria-label="Summary">
      <article class="card">
        <p class="label">Tracked benches</p>
        <p class="value">{{.Summary.BenchmarkCount}}</p>
        <p class="note">{{.Summary.PackageCount}} packages in the current baseline</p>
      </article>
      <article class="card">
        <p class="label">Zero allocs</p>
        <p class="value">{{.Summary.ZeroAllocs}}</p>
        <p class="note">benchmarks currently at 0 allocs/op</p>
      </article>
      <article class="card">
        <p class="label">Since {{.ReferenceName}}</p>
        <p class="value">{{pct .Summary.MedianDelta}}</p>
        <p class="note">median anchor-relative delta across {{.Summary.Common}} shared benches</p>
      </article>
      <article class="card">
        <p class="label">Movement</p>
        <p class="value">{{.Summary.Faster}} / {{.Summary.Slower}}</p>
        <p class="note">faster / slower by at least 5 percent</p>
      </article>
    </div>

    <section>
      <div class="section-head">
        <h2>Timeline</h2>
        <p>{{len .Timeline}} snapshot(s). CI snapshots graph real runs; seed points use committed historical/current JSON until the timeline fills in.</p>
      </div>
      {{if .Charts}}
      <div class="chart-grid">
        {{range .Charts}}
        <article class="chart">
          <div class="chart-head">
            <h3>{{.Title}}</h3>
            {{if ne .StatusClass "none"}}<span class="status {{.StatusClass}}">{{.Status}}</span>{{end}}
          </div>
          <p>{{.Subtitle}}</p>
          <svg viewBox="0 0 520 210" role="img" aria-label="{{.Title}} trend chart">
            {{if .HasBudget}}<rect class="budget-band" x="46" y="{{printf "%.2f" .BudgetY}}" width="456" height="{{printf "%.2f" (sub .RefY .BudgetY)}}"></rect>{{end}}
            {{if .ShowRefLine}}<line class="ref-line" x1="46" y1="{{printf "%.2f" .RefY}}" x2="502" y2="{{printf "%.2f" .RefY}}"></line>{{end}}
            <line class="axis" x1="46" y1="22" x2="46" y2="176"></line>
            <line class="axis" x1="46" y1="176" x2="502" y2="176"></line>
            <text class="axis-label" x="42" y="26" text-anchor="end">{{.YMaxText}}</text>
            <text class="axis-label" x="42" y="173" text-anchor="end">{{.YMinText}}</text>
            {{range .XTicks}}
            <line class="tick" x1="{{printf "%.2f" .X}}" y1="176" x2="{{printf "%.2f" .X}}" y2="179"></line>
            <text class="tick-label" x="{{printf "%.2f" .X}}" y="188">{{.Label}}</text>
            {{end}}
            {{range .Series}}
            {{$color := .Color}}
            {{if .BandPath}}<path class="chart-band" fill="{{$color}}" d="{{.BandPath}}"></path>{{end}}
            <path class="chart-line" stroke="{{$color}}" d="{{.Path}}"></path>
            {{range .Points}}
            <circle class="point" fill="{{$color}}" cx="{{printf "%.2f" .X}}" cy="{{printf "%.2f" .Y}}" r="3.2">
              <title>{{.Date}} @ {{.SHA}}: {{.Text}}{{if .HasBand}} ({{.Spread}}){{end}}</title>
            </circle>
            {{end}}
            {{end}}
            {{range .Deltas}}
            <text class="gap-label" x="{{printf "%.2f" .X}}" y="{{printf "%.2f" .Y}}">{{.Text}}</text>
            {{end}}
          </svg>
          <div class="chart-foot">
            <p class="chart-meta"><span class="chart-unit">{{.Unit}}</span>{{if .DeltaCaption}} · <span class="gap-delta">{{.DeltaCaption}}</span>{{end}}</p>
            <div class="legend">
              {{range .Series}}<span><i class="swatch" style="--series: {{.Color}}"></i>{{.Label}}</span>{{end}}
              {{if .HasRef}}<span><i class="swatch dash"></i>{{.RefLegend}}</span>{{end}}
              {{if .HasBudget}}<span><i class="swatch budget"></i>{{.BudgetText}} budget</span>{{end}}
            </div>
          </div>
        </article>
        {{end}}
      </div>
      {{else}}
      <div class="empty">No timeline-capable benchmarks found yet.</div>
      {{end}}
    </section>

    <section>
      <div class="section-head">
        <h2>Explore metrics over time</h2>
        <p>Pick any benchmark and metric; each line is a CPU model, the shaded band is the per-run min/max spread (not a 95% CI — typically ~3 samples) and every individual gathered sample is plotted as a dot. Hover for values. The anchor-relative ratio only normalizes within a CPU, so compare trends per CPU rather than absolute levels across them.</p>
      </div>
      <div id="perf-explorer" style="width:100%;min-height:420px"></div>
    </section>

    <section>
      <div class="section-head">
        <h2>Trend sparklines (first → last)</h2>
        <p>One row per benchmark × CPU. Benchmarks that differ only by a scaling factor (…/10, /100, /1000) are merged into a single row: their lines are overlaid and each is indexed to its own first value (% change from a shared 0% baseline) so you can compare how each scale moved — darker line = larger scale — with the per-scale Δ shown at right. Un-scaled benchmarks show one absolute sparkline (every sample a faint dot; hollow ring = first snapshot, filled = last) plus first/last/Δ. Slope reads direction; green = improvement, red = regression. Sorted by the largest endpoint change — click Δ to flip; use the filters to hide single-point or unchanged series.</p>
      </div>
      <div id="perf-sparklines" style="width:100%"></div>
    </section>

    <section>
      <div class="section-head">
        <h2>Largest changes</h2>
        <p>Current ratchet compared with {{.ReferenceName}}.</p>
      </div>
      <div class="grid-2">
        {{if .TopImprovements}}
        <table>
          <thead><tr><th>Improved</th><th>Delta</th><th>Now</th></tr></thead>
          <tbody>
            {{range .TopImprovements}}
            <tr>
              <td class="bench"><span class="pkg">{{.Package}}</span>{{.Name}}</td>
              <td><span class="delta good">{{deltaText .Delta}}</span></td>
              <td class="num">{{ratio .NewRatio}} anchors</td>
            </tr>
            {{end}}
          </tbody>
        </table>
        {{else}}<div class="empty">No historical improvements found.</div>{{end}}

        {{if .TopSlowdowns}}
        <table>
          <thead><tr><th>Slower</th><th>Delta</th><th>Now</th></tr></thead>
          <tbody>
            {{range .TopSlowdowns}}
            <tr>
              <td class="bench"><span class="pkg">{{.Package}}</span>{{.Name}}</td>
              <td><span class="delta bad">{{deltaText .Delta}}</span></td>
              <td class="num">{{ratio .NewRatio}} anchors</td>
            </tr>
            {{end}}
          </tbody>
        </table>
        {{else}}<div class="empty">No historical slowdowns found.</div>{{end}}
      </div>
    </section>

    <section>
      <div class="section-head">
        <h2>Recently tightened</h2>
        <p>Most recently lowered ratchet bars. × anchor normalizes wall time across machines; the last column is the change vs {{.ReferenceName}}.</p>
      </div>
      <table>
        <thead><tr><th>Benchmark</th><th>Bar set</th><th>× anchor</th><th>Wall</th><th>Allocs</th><th>vs {{.ReferenceName}}</th></tr></thead>
        <tbody>
          {{range .RecentlyTightened}}
          <tr>
            <td class="bench"><span class="pkg">{{.Package}}</span>{{.Name}}</td>
            <td class="num">{{date .BestSinceAt}} @ {{shortSHA .BestSinceSHA}}</td>
            <td class="num">{{ratio .Ratio}}</td>
            <td class="num">{{ns .NSPerOp}}</td>
            <td class="num">{{count .AllocsPerOp}}</td>
            <td class="num {{deltaClass .Delta}}">{{deltaText .Delta}}</td>
          </tr>
          {{end}}
        </tbody>
      </table>
    </section>

    <section>
      <div class="section-head">
        <h2>Current baseline</h2>
        <p>Sorted by package and benchmark. Lower anchor ratio is faster.</p>
      </div>
      <table>
        <thead>
          <tr>
            <th>Benchmark</th>
            <th>Ratio</th>
            <th>Wall</th>
            <th>Alloc</th>
            <th>Bytes</th>
            <th>Delta</th>
            <th>Scale</th>
          </tr>
        </thead>
        <tbody>
          {{range .Rows}}
          <tr>
            <td class="bench"><span class="pkg">{{.Package}}</span>{{.Name}}</td>
            <td class="num">{{ratio .Ratio}}</td>
            <td class="num">{{ns .NSPerOp}}</td>
            <td class="num">{{count .AllocsPerOp}}</td>
            <td class="num">{{bytes .BytesPerOp}}</td>
            <td><span class="delta {{deltaClass .Delta}}">{{deltaText .Delta}}</span></td>
            <td class="barcell"><div class="bartrack"><span class="barfill" style="--w: {{bar .BarWidth}}"></span></div></td>
          </tr>
          {{end}}
        </tbody>
      </table>
    </section>
  </main>

  <footer>
    <div class="wrap">Source data: docs/perf/baseline.json, docs/perf/historical/*.json, and docs/perf/timeline/*.json. Page generation does not run benchmarks.</div>
  </footer>

  <!-- One declarative d3-based stack drives every chart on the page. Observable
       Plot (the grammar-of-graphics layer from the d3 authors) renders the
       marks; it externalizes d3, so d3 MUST load first and Plot reuses that same
       global d3 (no duplicate copy — that's why Plot's bundle is small). We also
       use d3-array directly to shape the timeline (group/extent/endpoints).
       Both pinned + SRI-verified. This replaced the vega/vega-lite/vega-embed
       trio: Plot is lighter, consistent with our d3 data layer, and has the
       exact marks this dashboard needs (faceting, link/dumbbell, tip). -->
  <script src="https://cdn.jsdelivr.net/npm/d3@7.9.0/dist/d3.min.js" integrity="sha384-CjloA8y00+1SDAUkjs099PVfnY2KmDC2BZnws9kh8D/lX1s46w6EPhpXdqMfjK6i" crossorigin="anonymous"></script>
  <script src="https://cdn.jsdelivr.net/npm/@observablehq/plot@0.6.16/dist/plot.umd.min.js" integrity="sha384-XzQ+KW4LBWHv6FLjiPA1vjDz//oGlY8We4XROgpir5jHgg2+Qvo/6fT5TjbnqyYi" crossorigin="anonymous"></script>
  <script>
    // Shared helpers for the two Observable Plot views below.
    const PERF = (function () {
      const rootStyle = getComputedStyle(document.documentElement);
      const cssVar = function (name, fallback) {
        const v = rootStyle.getPropertyValue(name);
        return (v && v.trim()) || fallback;
      };
      return {
        GREEN: cssVar("--green", "#167a48"),
        RED: cssVar("--red", "#aa2e2e"),
        INK: cssVar("--ink", "#171717"),
        PAPER: cssVar("--paper", "#ffffff"),
        // CPU model strings are long; collapse to a short tag for axis labels.
        shortCPU: function (s) {
          if (!s) return s;
          const t = s.replace(/\(R\)|\(TM\)|Processor|Platinum|\d+-Core|Core|CPU.*$/gi, " ");
          const m = t.match(/(EPYC\s+\w+|Xeon[\s\w]*?\d{3,}\w*|Apple\s+M\w+|Ryzen\s+\w+)/i);
          return (m ? m[1] : t).replace(/\s+/g, " ").trim();
        },
        metricLabel: function (meta, m) {
          const mi = meta[m];
          return mi ? m + "  (" + mi.unit + ", " + (mi.lower_is_better ? "↓ lower is better" : "↑ higher is better") + ")" : m;
        },
        fmt: function (v) {
          if (v == null || isNaN(v)) return "—";
          const a = Math.abs(v);
          if (a >= 1000) return Math.round(v).toLocaleString();
          if (a >= 10) return v.toFixed(1);
          return v.toPrecision(3);
        },
        // Build the bench + metric <select> controls; calls onChange() on input.
        controls: function (host, opts) {
          const bar = d3.select(host).append("div").attr("class", "explorer-controls");
          if (opts.benches) {
            const bSel = bar.append("label").text("Benchmark ").append("select");
            bSel.selectAll("option").data(opts.benches).join("option")
              .attr("value", function (d) { return d; }).text(function (d) { return d; })
              .property("selected", function (d) { return d === opts.bench(); });
            bSel.on("change", function () { opts.setBench(this.value); opts.onChange(); });
          }
          const mSel = bar.append("label").text("Metric ").append("select");
          mSel.selectAll("option").data(opts.metrics).join("option")
            .attr("value", function (d) { return d; }).text(opts.metricText)
            .property("selected", function (d) { return d === opts.metric(); });
          mSel.on("change", function () { opts.setMetric(this.value); opts.onChange(); });
        }
      };
    })();
  </script>
  <script>
    // Explorer (Observable Plot): pick a benchmark + metric, see one line per
    // CPU with its min/max band AND every individual gathered sample as a dot.
    (function () {
      const host = document.getElementById("perf-explorer");
      if (!host) return;
      fetch({{.ExplorerURL}})
        .then(function (r) { if (!r.ok) throw new Error("HTTP " + r.status); return r.json(); })
        .then(render)
        .catch(function (e) { host.innerHTML = '<div class="empty">Could not load explorer data (' + e + ').</div>'; });

      function render(payload) {
        const rows = (payload && payload.rows) || [];
        const meta = (payload && payload.meta) || {};
        if (!window.Plot || !window.d3 || !rows.length) {
          host.innerHTML = '<div class="empty">No timeline data to explore yet.</div>';
          return;
        }
        rows.forEach(function (r) { r._t = new Date(r.date); });
        const benches = Array.from(new Set(rows.map(function (r) { return r.bench; }))).sort();
        const metrics = Array.from(new Set(rows.map(function (r) { return r.metric; }))).sort();
        let curBench = benches.find(function (b) { return b.indexOf("ClojureTestSuite [aot_native]") >= 0; }) || benches[0];
        let curMetric = metrics.indexOf("ratio_to_anchor") >= 0 ? "ratio_to_anchor" : metrics[0];

        host.innerHTML = "";
        PERF.controls(host, {
          benches: benches, metrics: metrics,
          bench: function () { return curBench; }, setBench: function (v) { curBench = v; },
          metric: function () { return curMetric; }, setMetric: function (v) { curMetric = v; },
          metricText: function (m) { return PERF.metricLabel(meta, m); },
          onChange: draw
        });
        const chart = d3.select(host).append("div").attr("class", "explorer-chart");

        function draw() {
          const mi = meta[curMetric] || { unit: curMetric, lower_is_better: true };
          const sub = rows.filter(function (r) { return r.bench === curBench && r.metric === curMetric; })
            .map(function (r) { return { _t: r._t, cpu: PERF.shortCPU(r.cpu), value: r.value, lo: r.lo, hi: r.hi, samples: r.samples }; });
          chart.selectAll("*").remove();
          if (!sub.length) { chart.append("div").attr("class", "empty").text("No data for this selection."); return; }
          // Flatten every retained sample into its own observable point.
          const samples = [];
          sub.forEach(function (r) {
            (r.samples && r.samples.length ? r.samples : [r.value]).forEach(function (v) {
              samples.push({ _t: r._t, cpu: r.cpu, sample: v });
            });
          });
          const width = chart.node().clientWidth || 720;
          const fig = Plot.plot({
            width: width, height: 380, marginLeft: 60, marginBottom: 34,
            x: { label: null, type: "utc" },
            y: { grid: true, label: mi.unit + "  (" + (mi.lower_is_better ? "↓ better" : "↑ better") + ")", tickFormat: "~s" },
            color: { legend: true, label: "CPU" },
            marks: [
              Plot.areaY(sub, { x: "_t", y1: "lo", y2: "hi", fill: "cpu", fillOpacity: 0.1, curve: "monotone-x" }),
              Plot.dot(samples, { x: "_t", y: "sample", stroke: "cpu", r: 1.7, strokeOpacity: 0.45 }),
              Plot.lineY(sub, { x: "_t", y: "value", stroke: "cpu", strokeWidth: 1.7, curve: "monotone-x" }),
              Plot.dot(sub, { x: "_t", y: "value", fill: "cpu", r: 2.6, stroke: "white", strokeWidth: 0.6, tip: true })
            ]
          });
          chart.node().append(fig);
        }
        draw();
        let raf;
        window.addEventListener("resize", function () { cancelAnimationFrame(raf); raf = requestAnimationFrame(draw); });
      }
    })();
  </script>
  <script>
    // Trend sparklines (canvas + d3). One row per benchmark × CPU; the cell is
    // a real sparkline of the whole series (so you see the SHAPE, not just two
    // endpoints), each scaled to its own y-range, with every gathered sample as
    // a faint dot and hollow/filled rings marking first/last. Canvas, not SVG or
    // Plot: hundreds of rows × hundreds of sample dots would be tens of thousands
    // of nodes, and Plot can't give each row an independent y-scale (faceted
    // scales are shared) — exactly the dense, per-row-scaled case canvas owns.
    (function () {
      const host = document.getElementById("perf-sparklines");
      if (!host) return;
      fetch({{.ExplorerURL}})
        .then(function (r) { if (!r.ok) throw new Error("HTTP " + r.status); return r.json(); })
        .then(render)
        .catch(function (e) { host.innerHTML = '<div class="empty">Could not load sparkline data (' + e + ').</div>'; });

      function render(payload) {
        const rows = (payload && payload.rows) || [];
        const meta = (payload && payload.meta) || {};
        if (!window.d3 || !rows.length) {
          host.innerHTML = '<div class="empty">No timeline data for sparklines yet.</div>';
          return;
        }
        rows.forEach(function (r) { r._t = +new Date(r.date); });
        const metrics = Array.from(new Set(rows.map(function (r) { return r.metric; }))).sort();
        let curMetric = metrics.indexOf("ratio_to_anchor") >= 0 ? "ratio_to_anchor" : metrics[0];
        let sortDesc = true;   // by |Δ| descending
        let hideSingle = true; // drop series with a single snapshot (no trend, Δ=0)
        let hideFlat = true;   // drop series whose endpoints didn't move (|Δ| < 0.1%)
        let nameFilter = "";   // substring match on "bench · cpu"

        host.innerHTML = "";
        const bar = d3.select(host).append("div").attr("class", "explorer-controls");
        const mSel = bar.append("label").text("Metric ").append("select");
        mSel.selectAll("option").data(metrics).join("option")
          .attr("value", function (d) { return d; }).text(function (m) { return PERF.metricLabel(meta, m); })
          .property("selected", function (d) { return d === curMetric; });
        mSel.on("change", function () { curMetric = this.value; draw(); });
        bar.append("label").text("Filter ").append("input")
          .attr("type", "search").attr("placeholder", "benchmark or CPU…")
          .on("input", function () { nameFilter = this.value.trim().toLowerCase(); draw(); });
        const cbSingle = bar.append("label").attr("class", "cb");
        cbSingle.append("input").attr("type", "checkbox").property("checked", hideSingle)
          .on("change", function () { hideSingle = this.checked; draw(); });
        cbSingle.append("span").text("hide single-point");
        const cbFlat = bar.append("label").attr("class", "cb");
        cbFlat.append("input").attr("type", "checkbox").property("checked", hideFlat)
          .on("change", function () { hideFlat = this.checked; draw(); });
        cbFlat.append("span").text("hide unchanged (Δ≈0)");
        const count = bar.append("span").attr("class", "spark-count");
        const wrap = d3.select(host).append("div").attr("class", "spark-table-wrap");

        // Hover tooltip: the actual value (+ % change vs each line's first) at
        // the snapshot under the cursor — the numbers the overlaid lines can't
        // show on their own.
        const tip = d3.select(host).append("div").attr("class", "spark-tip").style("opacity", 0);
        function esc(s) { return String(s).replace(/[&<>]/g, function (c) { return c === "&" ? "&amp;" : c === "<" ? "&lt;" : "&gt;"; }); }
        function hideTip() { tip.style("opacity", 0); }
        function rowAt(se, t0) { for (let i = 0; i < se.s.length; i++) { if (se.s[i]._t === t0) return se.s[i]; } return null; }
        function showTip(ev, d, t0) {
          const lib = (meta[curMetric] || {}).lower_is_better;
          let dateStr = "";
          for (let i = 0; i < d.series.length; i++) { const r = rowAt(d.series[i], t0); if (r && r.date) { dateStr = String(r.date).slice(0, 10); break; } }
          let html = '<div class="tip-h">' + esc(d.base) + '</div><div class="tip-sub">' + esc(d.cpu) + (dateStr ? " · " + dateStr : "") + '</div>';
          d.series.forEach(function (se) {
            const r = rowAt(se, t0); if (!r) return;
            const pct = se.first ? (r.value / se.first - 1) * 100 : 0;
            const cls = pct === 0 ? "" : ((lib ? pct < 0 : pct > 0) ? "good" : "bad");
            const lab = d.grouped ? '<b>' + fmtScale(se.scale) + '</b>' : '<b></b>';
            html += '<div class="tip-r">' + lab + '<span>' + PERF.fmt(r.value) + '</span><span class="' + cls + '">' + deltaStr(pct) + '</span></div>';
            if (!d.grouped && r.lo !== r.hi) { html += '<div class="tip-r muted"><b></b><span>spread</span><span>' + PERF.fmt(r.lo) + "–" + PERF.fmt(r.hi) + '</span></div>'; }
          });
          tip.html(html).style("opacity", 1);
          const tw = tip.node().offsetWidth, th = tip.node().offsetHeight, gap = 14;
          let lx = ev.clientX + gap, ly = ev.clientY + gap;
          if (lx + tw > window.innerWidth - 8) lx = ev.clientX - tw - gap;
          if (ly + th > window.innerHeight - 8) ly = ev.clientY - th - gap;
          tip.style("left", lx + "px").style("top", ly + "px");
        }

        // Split a benchmark name into its base and scaling factor: a trailing
        // /<digits> size (before any [mode] suffix). MapAssoc/HAMT-Assoc/01000
        // [bytecode] → base "MapAssoc/HAMT-Assoc [bytecode]", scale 1000.
        function parseBench(bench) {
          let mode = "";
          const mm = bench.match(/(\s*\[[^\]]*\])$/);
          let core = bench;
          if (mm) { mode = mm[1]; core = bench.slice(0, bench.length - mm[1].length); }
          const sm = core.match(/^(.*)\/0*(\d+)$/);
          if (sm) { return { base: sm[1] + mode, scale: +sm[2] }; }
          return { base: bench, scale: null };
        }
        function fmtScale(n) {
          if (n == null) return "";
          if (n >= 1000000) return (n / 1000000) + "M";
          if (n >= 1000) return (n / 1000) + "k";
          return "" + n;
        }
        function deltaStr(p) { return (p >= 0 ? "+" : "") + p.toFixed(1) + "%"; }

        function draw() {
          const mi = meta[curMetric] || { unit: curMetric, lower_is_better: true };
          const sub = rows.filter(function (r) { return r.metric === curMetric; });
          // Group by (base benchmark, CPU); scaling-factor variants collapse into
          // one row carrying a SET of per-scale series.
          const byKey = d3.group(sub, function (d) {
            const p = parseBench(d.bench); d.__base = p.base; d.__scale = p.scale;
            return p.base + " @@ " + d.cpu;
          });
          let recs = [];
          byKey.forEach(function (items) {
            const base = items[0].__base, cpu = items[0].cpu;
            const byScale = d3.group(items, function (d) { return d.__scale; });
            let series = [];
            byScale.forEach(function (rs, scale) {
              const s = rs.slice().sort(function (a, b) { return a._t - b._t; });
              if (!s.length) return;
              const first = s[0].value, last = s[s.length - 1].value;
              const deltaPct = first ? (last - first) / first * 100 : 0;
              const good = mi.lower_is_better ? deltaPct < 0 : deltaPct > 0;
              series.push({ scale: scale, s: s, first: first, last: last, deltaPct: deltaPct, good: good, n: s.length });
            });
            series.sort(function (a, b) { return (a.scale || 0) - (b.scale || 0); });
            recs.push({
              base: base, cpu: PERF.shortCPU(cpu), series: series, grouped: series.length > 1,
              n: d3.max(series, function (d) { return d.n; }) || 1,
              sortDelta: d3.max(series, function (d) { return Math.abs(d.deltaPct); }) || 0
            });
          });
          const total = recs.length;
          // Dynamic filtering: most rows are single-snapshot or flat (Δ≈0).
          if (hideSingle) { recs = recs.filter(function (r) { return r.n >= 2; }); }
          if (hideFlat) { recs = recs.filter(function (r) { return r.sortDelta >= 0.1; }); }
          if (nameFilter) { recs = recs.filter(function (r) { return (r.base + " · " + r.cpu).toLowerCase().indexOf(nameFilter) >= 0; }); }
          recs.sort(function (a, b) { return (sortDesc ? -1 : 1) * (a.sortDelta - b.sortDelta); });
          count.text(recs.length + " of " + total + " series");

          wrap.selectAll("*").remove();
          if (!recs.length) { wrap.append("div").attr("class", "empty").text("No series match the current filters."); return; }

          const table = wrap.append("table").attr("class", "spark-table");
          const htr = table.append("thead").append("tr");
          htr.append("th").text("Benchmark");
          htr.append("th").text("CPU");
          htr.append("th").text("Trend (scale variants overlaid, indexed to % change)");
          htr.append("th").attr("class", "num").text("First");
          htr.append("th").attr("class", "num").text("Last");
          htr.append("th").attr("class", "num").style("cursor", "pointer").text("Δ " + (sortDesc ? "▼" : "▲"))
            .on("click", function () { sortDesc = !sortDesc; draw(); });

          const tb = table.append("tbody");
          tb.selectAll("tr").data(recs).join("tr").each(function (d) {
            const row = d3.select(this);
            row.append("td").attr("class", "bench").attr("title", d.base).text(d.base);
            row.append("td").text(d.cpu);
            row.append("td").attr("class", "spark").each(function () { drawSpark(this, d); });
            if (d.grouped) {
              // Scale family: per-scale Δ pills span the First/Last/Δ columns.
              const cell = row.append("td").attr("colspan", 3).attr("class", "scales");
              d.series.forEach(function (se) {
                const pill = cell.append("span").attr("class", "scalepill " + (se.good ? "good" : "bad"));
                pill.append("b").text(fmtScale(se.scale));
                pill.append("span").text(deltaStr(se.deltaPct));
              });
            } else {
              const se = d.series[0];
              row.append("td").attr("class", "num").text(PERF.fmt(se.first));
              row.append("td").attr("class", "num").text(PERF.fmt(se.last));
              row.append("td").attr("class", "num").append("span")
                .attr("class", "delta " + (se.good ? "good" : "bad")).text(deltaStr(se.deltaPct));
            }
          });
        }

        // Build the sparkline canvas plus a render(hoverTime) closure and the
        // list of snapshot times, then wire mousemove → crosshair + tooltip so
        // you can read the actual value (and % change) at any snapshot.
        function drawSpark(td, d) {
          const W = 160, H = 30, pad = 4;
          const dpr = window.devicePixelRatio || 1;
          const canvas = d3.select(td).append("canvas")
            .attr("width", Math.round(W * dpr)).attr("height", Math.round(H * dpr))
            .style("width", W + "px").style("height", H + "px").node();
          const ctx = canvas.getContext("2d");
          ctx.scale(dpr, dpr); // CSS pixels, crisp on HiDPI
          let x, times, render;

          if (d.grouped) {
            // One line per scale, each indexed to its own first value (% change)
            // so they share a 0% baseline; color ramps light→dark with scale.
            const allT = [];
            d.series.forEach(function (se) {
              se.idx = se.s.map(function (p) { return se.first ? (p.value / se.first - 1) * 100 : 0; });
              se.s.forEach(function (p) { allT.push(p._t); });
            });
            x = d3.scaleLinear().domain(d3.extent(allT)).range([pad, W - pad]);
            let lo = d3.min(d.series, function (se) { return d3.min(se.idx); });
            let hi = d3.max(d.series, function (se) { return d3.max(se.idx); });
            if (lo === hi) { lo -= 1; hi += 1; }
            const y = d3.scaleLinear().domain([lo, hi]).range([H - pad, pad]);
            times = Array.from(new Set(allT)).sort(function (a, b) { return a - b; });
            const n = d.series.length;
            render = function (ht) {
              ctx.clearRect(0, 0, W, H);
              ctx.beginPath(); ctx.moveTo(pad, y(0)); ctx.lineTo(W - pad, y(0));
              ctx.strokeStyle = "rgba(0,0,0,0.18)"; ctx.lineWidth = 0.5; ctx.stroke();
              if (ht != null) { ctx.beginPath(); ctx.moveTo(x(ht), pad); ctx.lineTo(x(ht), H - pad); ctx.strokeStyle = "rgba(0,0,0,0.3)"; ctx.lineWidth = 0.5; ctx.stroke(); }
              d.series.forEach(function (se, i) {
                const col = d3.interpolateBlues(0.45 + 0.5 * (n > 1 ? i / (n - 1) : 0));
                ctx.beginPath();
                d3.line().x(function (v, j) { return x(se.s[j]._t); }).y(function (v) { return y(v); }).curve(d3.curveMonotoneX).context(ctx)(se.idx);
                ctx.lineWidth = 1.2; ctx.strokeStyle = col; ctx.stroke();
                const last = se.s.length - 1;
                ctx.beginPath(); ctx.arc(x(se.s[0]._t), y(se.idx[0]), 2, 0, 2 * Math.PI); ctx.fillStyle = PERF.PAPER; ctx.fill(); ctx.lineWidth = 1; ctx.strokeStyle = col; ctx.stroke();
                ctx.beginPath(); ctx.arc(x(se.s[last]._t), y(se.idx[last]), 2, 0, 2 * Math.PI); ctx.fillStyle = col; ctx.fill();
                if (ht != null) { for (let j = 0; j < se.s.length; j++) { if (se.s[j]._t === ht) { ctx.beginPath(); ctx.arc(x(se.s[j]._t), y(se.idx[j]), 2.6, 0, 2 * Math.PI); ctx.fillStyle = col; ctx.fill(); break; } } }
              });
            };
          } else {
            // Absolute trend + every sample dot + first/last rings.
            const se = d.series[0], col = se.good ? PERF.GREEN : PERF.RED;
            x = d3.scaleLinear().domain(d3.extent(se.s, function (p) { return p._t; })).range([pad, W - pad]);
            let lo = d3.min(se.s, function (p) { return Math.min(p.lo, d3.min(p.samples && p.samples.length ? p.samples : [p.value])); });
            let hi = d3.max(se.s, function (p) { return Math.max(p.hi, d3.max(p.samples && p.samples.length ? p.samples : [p.value])); });
            if (lo === hi) { lo -= 1; hi += 1; }
            const y = d3.scaleLinear().domain([lo, hi]).range([H - pad, pad]);
            times = se.s.map(function (p) { return p._t; });
            render = function (ht) {
              ctx.clearRect(0, 0, W, H);
              if (ht != null) { ctx.beginPath(); ctx.moveTo(x(ht), pad); ctx.lineTo(x(ht), H - pad); ctx.strokeStyle = "rgba(0,0,0,0.3)"; ctx.lineWidth = 0.5; ctx.stroke(); }
              ctx.beginPath();
              d3.line().x(function (p) { return x(p._t); }).y(function (p) { return y(p.value); }).curve(d3.curveMonotoneX).context(ctx)(se.s);
              ctx.lineWidth = 1.3; ctx.strokeStyle = col; ctx.stroke();
              ctx.fillStyle = col; ctx.globalAlpha = 0.5;
              se.s.forEach(function (p) {
                const vs = p.samples && p.samples.length ? p.samples : [p.value];
                for (let i = 0; i < vs.length; i++) { ctx.beginPath(); ctx.arc(x(p._t), y(vs[i]), 1.3, 0, 2 * Math.PI); ctx.fill(); }
              });
              ctx.globalAlpha = 1;
              const f = se.s[0], l = se.s[se.s.length - 1];
              ctx.beginPath(); ctx.arc(x(f._t), y(f.value), 2.4, 0, 2 * Math.PI); ctx.fillStyle = PERF.PAPER; ctx.fill(); ctx.lineWidth = 1.1; ctx.strokeStyle = col; ctx.stroke();
              ctx.beginPath(); ctx.arc(x(l._t), y(l.value), 2.4, 0, 2 * Math.PI); ctx.fillStyle = col; ctx.fill();
              if (ht != null) { for (let j = 0; j < se.s.length; j++) { if (se.s[j]._t === ht) { ctx.beginPath(); ctx.arc(x(se.s[j]._t), y(se.s[j].value), 2.8, 0, 2 * Math.PI); ctx.fillStyle = col; ctx.fill(); break; } } }
            };
          }
          render(null);
          d3.select(canvas).style("cursor", "crosshair")
            .on("mousemove", function (ev) {
              const px = d3.pointer(ev, canvas)[0];
              const t = x.invert(px);
              let t0 = times[0], best = Infinity;
              for (let i = 0; i < times.length; i++) { const dd = Math.abs(times[i] - t); if (dd < best) { best = dd; t0 = times[i]; } }
              render(t0); showTip(ev, d, t0);
            })
            .on("mouseleave", function () { render(null); hideTip(); });
        }

        draw();
      }
    })();
  </script>
</body>
</html>
`
