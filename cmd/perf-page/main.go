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

type Baseline = perfdata.Baseline
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
}

func main() {
	var (
		baselinePath   = flag.String("baseline", "docs/perf/baseline.json", "current baseline JSON")
		historicalPath = flag.String("historical", "docs/perf/historical", "historical baseline directory")
		timelinePath   = flag.String("timeline", "docs/perf/timeline", "timeline snapshot directory")
		outPath        = flag.String("out", "docs/perf/index.html", "HTML output path")
		logoPath       = flag.String("logo", "meta/logo.svg", "logo SVG to embed")
		cpuFilter      = flag.String("cpu", "", "keep only timeline snapshots whose machine cpu_model contains this substring (CI runs land on ≥2 CPU tiers whose ratio_to_anchor doesn't normalize across them, so a mixed timeline zig-zags ~2x; filtering to one tier gives a clean series). Empty = all.")
	)
	flag.Parse()

	current, err := loadBaseline(*baselinePath)
	if err != nil {
		die("load baseline: %v", err)
	}
	reference, referenceName, err := loadLatestHistorical(*historicalPath)
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
	fmt.Printf("wrote %s (%d benchmarks)\n", *outPath, len(current.Benchmarks))
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

func loadBaseline(path string) (Baseline, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return Baseline{}, err
	}
	var baseline Baseline
	if err := json.Unmarshal(b, &baseline); err != nil {
		return Baseline{}, err
	}
	if len(baseline.Benchmarks) == 0 {
		return Baseline{}, fmt.Errorf("%s has no benchmarks", path)
	}
	return baseline, nil
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
</body>
</html>
`
