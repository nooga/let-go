package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/nooga/let-go/pkg/perfdata"
)

// timelineName matches a perf-timeline snapshot filename:
// TIMESTAMP-SHORTSHA-MACHINESLUG.json.
//
// Anchoring the shape here rather than splitting on "-" matters because the
// machine slug is itself full of dashes: a positional split yields a
// plausible-looking but wrong SHA for any name that does not have exactly the
// expected leading fields, and nothing downstream can tell that apart from a
// real one. A non-matching name is skipped and reported instead.
var timelineName = regexp.MustCompile(`^(\d{8}T\d{6}Z)-([0-9a-f]{7,40})-(.+)\.json$`)

// timelineFile is one parsed snapshot filename.
type timelineFile struct {
	path      string
	sha       string
	machine   string // filename slug, e.g. "amd64-amd-epyc-7763-64-core-processor"
	timestamp string
}

// seedCandidate pairs a snapshot file with the machine profile read out of it.
type seedCandidate struct {
	file timelineFile
	mb   MachineBaseline
}

// seedOptions are the knobs on `seed-baseline`. Defaults are set in main().
type seedOptions struct {
	// window is how many recent snapshots per machine key participate. One
	// snapshot is one CI run, and one CI run is one sample.
	window int
	// coherenceTolerance rejects a candidate whose ratio_to_anchor values sit,
	// in median, further than this fraction off the rest of the window.
	coherenceTolerance float64
	// iterationTolerance reports a benchmark whose b.N spread across the
	// window exceeds this fraction.
	iterationTolerance float64
	// minIterations reports a benchmark whose median b.N falls below this.
	minIterations int64
	// archPrefix restricts seeding to one architecture (#651: amd64-only).
	archPrefix string
}

// Defaults for `seed-baseline`, shared by the flag definitions in main() and
// by defaultSeedOptions so the documented behaviour has one source.
const (
	defaultSeedWindow             = 5
	defaultSeedCoherenceTolerance = 0.05
	defaultSeedIterationTolerance = 0.10
	defaultSeedMinIterations      = 20
	defaultSeedArch               = "amd64"
)

// defaultSeedOptions is the configuration `seed-baseline` runs with when no
// flags override it.
func defaultSeedOptions() seedOptions {
	return seedOptions{
		window:             defaultSeedWindow,
		coherenceTolerance: defaultSeedCoherenceTolerance,
		iterationTolerance: defaultSeedIterationTolerance,
		minIterations:      defaultSeedMinIterations,
		archPrefix:         defaultSeedArch,
	}
}

// unstableBenchmarks are excluded from the seed by decision, not by
// measurement (#651): BenchmarkClojureTestSuite and its CompileAndRun sibling,
// under each of the bytecode / ir_bytecode / aot_native variants.
//
// This list is authoritative but not self-checking, so seedBaseline reports an
// entry that matched nothing. A filter that has silently stopped applying —
// after a rename, say — still reads as if it filtered.
var unstableBenchmarks = []string{
	"github.com/nooga/let-go/test.BenchmarkClojureTestSuite",
	"github.com/nooga/let-go/test.BenchmarkClojureTestSuiteCompileAndRun",
}

// seedBaseline derives docs/perf/baseline.json from the perf-data timeline.
//
// The numbers stored for a machine tier are medians over the last `window`
// snapshots for that tier, not the newest snapshot. One snapshot is one CI run:
// these distributions are a tight core with a one-sided slow tail, so a single
// run has a real chance of being a tail observation, and a baseline seeded from
// one is wrong in a direction nothing downstream can detect.
//
// Median rather than min across the window: `update` already takes a min over
// history when it ratchets, and seeding with a second minimum would stack two
// of them into a floor no clean run can reach.
//
// Identity (captured_at, captured_at_sha, machine) comes from the newest
// surviving snapshot; the numbers come from the window. The SHA therefore names
// the newest contributing run rather than the sole source of the numbers — the
// seed log prints the whole window so the distinction is visible.
func seedBaseline(baselinePath, perfDataDir string, opt seedOptions) {
	paths, err := filepath.Glob(filepath.Join(perfDataDir, "*.json"))
	if err != nil {
		die("list timeline files: %v", err)
	}
	if len(paths) == 0 {
		die("no timeline snapshots found in %s", perfDataDir)
	}

	fmt.Printf("bench-ratchet: seed-baseline from %s\n", perfDataDir)
	fmt.Printf("  window %d snapshots per machine key, coherence tolerance ±%.0f%%\n",
		opt.window, opt.coherenceTolerance*100)

	byKey, skipped := gatherCandidates(paths, opt)
	if len(skipped) > 0 {
		fmt.Printf("  skipped %d file(s) not matching TIMESTAMP-SHA-MACHINE.json: %s\n",
			len(skipped), strings.Join(truncate(skipped, 3), ", "))
	}
	if len(byKey) == 0 {
		die("no %s machine snapshots found in %s", opt.archPrefix, perfDataDir)
	}

	merged := Baseline{Version: schemaVersion, Machines: map[string]MachineBaseline{}}
	for _, key := range sortedKeys(byKey) {
		cands := byKey[key]
		// Newest first, then keep at most `window`.
		sort.Slice(cands, func(i, j int) bool { return cands[i].file.timestamp > cands[j].file.timestamp })
		if len(cands) > opt.window {
			cands = cands[:opt.window]
		}
		mb, ok := seedOneMachine(key, cands, opt)
		if !ok {
			continue
		}
		merged.Machines[key] = mb
	}

	// Preserve the local arm64/Apple M3 profile: it gates developer machines
	// and has no counterpart in CI (#651). The Arch guard matters as much as
	// the model string: without it an amd64 CPUModel containing "M3" would
	// overwrite the fresh seed this run just computed for its own tier.
	if data, err := os.ReadFile(baselinePath); err == nil {
		var existing Baseline
		if err := json.Unmarshal(data, &existing); err == nil {
			for key, mb := range existing.Machines {
				if mb.Machine.Arch != opt.archPrefix && strings.Contains(mb.Machine.CPUModel, "M3") {
					merged.Machines[key] = mb
					fmt.Printf("  preserved: %s (local, not CI-derived)\n", key)
				}
			}
		}
	}

	if len(merged.Machines) == 0 {
		die("no machine profiles survived seeding — refusing to write an empty baseline")
	}
	if err := writeBaseline(baselinePath, merged); err != nil {
		die("write baseline: %v", err)
	}
	fmt.Printf("  wrote baseline → %s (%d machine profiles)\n", baselinePath, len(merged.Machines))
}

// gatherCandidates reads every snapshot and groups the profiles inside them by
// the machine key derived from their CONTENT.
//
// Grouping on the content key rather than the filename slug is what keeps a
// single mis-named file from overwriting a good window. Under filename
// grouping such a file forms a group of its own, is reduced on its own, and is
// then stored under the key its content names; whichever group the outer loop
// reaches last wins, so a one-file reduction can replace the five-run
// reduction the correctly-named files just produced for that key. Grouped by
// content, the stray file is simply one more candidate in the window it
// belongs to.
//
// No snapshot in perf-data has this shape as of 2026-09-02 (checked across all
// 382 files; the branch is append-only, with no deletions or renames in its
// history). This is a guard on a silent, order-dependent failure, not a repair
// of live data.
//
// The filename is still reported: it is what the log prints and what a human
// greps for, so a divergence has to be visible even though it no longer decides
// anything.
func gatherCandidates(paths []string, opt seedOptions) (map[string][]seedCandidate, []string) {
	byKey := map[string][]seedCandidate{}
	var skipped []string
	for _, p := range paths {
		m := timelineName.FindStringSubmatch(filepath.Base(p))
		if m == nil {
			skipped = append(skipped, filepath.Base(p))
			continue
		}
		f := timelineFile{path: p, timestamp: m[1], sha: m[2], machine: m[3]}

		var b Baseline
		data, err := os.ReadFile(f.path)
		if err != nil {
			die("read timeline snapshot %s: %v", f.path, err)
		}
		if err := json.Unmarshal(data, &b); err != nil {
			die("parse timeline snapshot %s: %v", f.path, err)
		}
		for _, mb := range b.Machines {
			// The architecture filter runs on the content too. A filename that
			// promised this architecture and a profile that isn't it is the
			// surprising case, so that one is reported; a file named for
			// another architecture is skipped quietly, as before.
			if mb.Machine.Arch != opt.archPrefix {
				if strings.HasPrefix(f.machine, opt.archPrefix+"-") {
					fmt.Printf("  skipping %s profile inside %s\n",
						mb.Machine.Arch, filepath.Base(f.path))
				}
				continue
			}
			key := perfdata.MachineKey(mb.Machine)
			if got := slugify(key); got != f.machine {
				fmt.Printf("  WARNING: %s is named for %q but carries %q — seeding it into the %q window\n",
					filepath.Base(f.path), f.machine, got, got)
			}
			byKey[key] = append(byKey[key], seedCandidate{file: f, mb: mb})
		}
	}
	return byKey, skipped
}

// seedOneMachine reduces one machine key's window into a single profile.
// Returns ok=false if the tier could not be seeded.
func seedOneMachine(key string, cands []seedCandidate, opt seedOptions) (MachineBaseline, bool) {
	window := len(cands)

	cands, dropped := rejectIncoherent(cands, opt.coherenceTolerance)
	for _, d := range dropped {
		fmt.Printf("  %s: rejected %s — ratios sit %+.1f%% off the window across %d benchmarks (anchor %+.1f%%)\n",
			key, d.sha, d.ratioOffset*100, d.shared, d.anchorDev*100)
	}
	if len(cands) == 0 {
		fmt.Printf("  %s: no snapshot in the window agrees with the others — skipped\n", key)
		return MachineBaseline{}, false
	}

	// Newest surviving snapshot supplies identity.
	sort.Slice(cands, func(i, j int) bool { return cands[i].file.timestamp > cands[j].file.timestamp })
	newest := cands[0]

	anchorNs := medianOf(mapf(cands, func(c seedCandidate) float64 { return c.mb.Anchor.NSPerOp }))
	if anchorNs <= 0 {
		fmt.Printf("  %s: window anchor median is %.3f ns/op — skipped\n", key, anchorNs)
		return MachineBaseline{}, false
	}

	out := MachineBaseline{
		CapturedAt:    newest.mb.CapturedAt,
		CapturedAtSHA: newest.mb.CapturedAtSHA,
		Machine:       newest.mb.Machine,
		Anchor: AnchorRecord{
			Name:       newest.mb.Anchor.Name,
			Package:    newest.mb.Anchor.Package,
			NSPerOp:    anchorNs,
			Iterations: newest.mb.Anchor.Iterations,
		},
		Benchmarks: map[string]BenchmarkEntry{},
	}

	rep := reduceBenchmarks(cands, anchorNs, opt)
	out.Benchmarks = rep.entries

	fmt.Printf("  %s: %d benchmarks from %d/%d snapshots (anchor %.3f ns/op, newest %s)\n",
		key, len(out.Benchmarks), len(cands), window,
		anchorNs, newest.file.sha)
	rep.report(key, opt)

	return out, true
}

// coherence measures one candidate's agreement with the rest of its window.
type coherence struct {
	sha string
	// anchorDev is the anchor's deviation from the window median. Reported
	// only — see rejectIncoherent for why it is not a rejection criterion.
	anchorDev float64
	// ratioOffset is the median, across shared benchmarks, of this snapshot's
	// ratio_to_anchor against the window median for that benchmark.
	ratioOffset float64
	shared      int
}

// rejectIncoherent drops candidates whose NORMALIZED numbers disagree with the
// rest of the window.
//
// It deliberately does not gate on the anchor's absolute deviation, which was
// the obvious design and is wrong. Measured over the 24 most recent amd64
// snapshots in perf-data (2026-08-05): anchor deviation ranges -22.4%..+3.1%,
// while ratio_to_anchor for the same snapshots holds to a median 0.03% and a
// worst 1.75%. The -22.4% snapshot (a588a69d2759, EPYC 9V74) is uniformly
// 22.4% fast in raw ns/op across all 162 of its benchmarks and agrees with its
// window on every ratio to within 0.1% — the whole host was fast and the anchor
// divided that back out, which is what the anchor is for. Rejecting on anchor
// deviation would have discarded three good captures, two of which are the very
// snapshots this baseline is seeded from.
//
// What does need rejecting is the MIXED capture: the anchor caught the slow
// tail and the benchmarks did not, or vice versa. Then every ratio from that
// snapshot is uniformly wrong while its raw numbers look ordinary, and the
// stored floor is off by that amount permanently. A mixed capture shows up
// exactly here, as a whole-snapshot offset in normalized space, well clear of
// the 1.75% the corpus actually exhibits.
//
// Below three candidates there is no median worth testing against, so the
// window passes through unfiltered.
func rejectIncoherent(cands []seedCandidate, tolerance float64) ([]seedCandidate, []coherence) {
	if len(cands) < 3 || tolerance <= 0 {
		return cands, nil
	}
	anchorMed := medianOf(mapf(cands, func(c seedCandidate) float64 { return c.mb.Anchor.NSPerOp }))

	// Window median ratio per benchmark, over the candidates that carry it.
	perBench := map[string][]float64{}
	for _, c := range cands {
		for name, e := range c.mb.Benchmarks {
			if e.RatioToAnchor > 0 {
				perBench[name] = append(perBench[name], e.RatioToAnchor)
			}
		}
	}
	windowMed := map[string]float64{}
	for name, vals := range perBench {
		if len(vals) >= 3 {
			windowMed[name] = medianOf(vals)
		}
	}

	var kept []seedCandidate
	var dropped []coherence
	for _, c := range cands {
		var offsets []float64
		for name, e := range c.mb.Benchmarks {
			med, ok := windowMed[name]
			if !ok || med <= 0 || e.RatioToAnchor <= 0 {
				continue
			}
			offsets = append(offsets, e.RatioToAnchor/med-1)
		}
		co := coherence{sha: c.file.sha, shared: len(offsets)}
		if anchorMed > 0 {
			co.anchorDev = (c.mb.Anchor.NSPerOp - anchorMed) / anchorMed
		}
		// Too few shared benchmarks to judge: keep it rather than guess.
		if len(offsets) < 3 {
			kept = append(kept, c)
			continue
		}
		co.ratioOffset = medianOf(offsets)
		if co.ratioOffset > tolerance || co.ratioOffset < -tolerance {
			dropped = append(dropped, co)
			continue
		}
		kept = append(kept, c)
	}
	return kept, dropped
}

// seedReport carries what reduceBenchmarks noticed but did not act on.
type seedReport struct {
	entries map[string]BenchmarkEntry
	// thin lists benchmarks that appeared in too few snapshots to reduce.
	thin []string
	// moved lists benchmarks whose b.N spread exceeded the tolerance.
	moved []string
	// lowN lists benchmarks whose median b.N fell below the floor.
	lowN []string
	// unmatched lists unstableBenchmarks entries that filtered nothing.
	unmatched []string
}

// reduceBenchmarks medians each benchmark across the surviving window.
//
// The reduction happens in RATIO space, not raw ns/op, and ns_per_op is derived
// back from the reduced ratio and the window's anchor. Measured over the recent
// amd64 corpus (2026-08-05): a snapshot's raw ns/op can sit 22% off its window
// because the host was fast that day, while its ratio_to_anchor holds to 1.75%
// worst case. Reducing the quantity that carries host speed, then dividing,
// imports that speed into the stored floor; reducing the normalized quantity
// does not. The two agree whenever the window is drawn from one host and differ
// exactly when it is not, which is the case a shared runner pool guarantees.
//
// Deriving ns_per_op from the same anchor that is stored also keeps
// entry.ratio_to_anchor == entry.ns_per_op / anchor.ns_per_op true by
// construction, which is what `check` relies on.
func reduceBenchmarks(cands []seedCandidate, anchorNs float64, opt seedOptions) seedReport {
	rep := seedReport{entries: map[string]BenchmarkEntry{}}

	type acc struct {
		ratios, bytes, allocs []float64
		iters                 []float64
	}
	byName := map[string]*acc{}
	matched := map[string]bool{}
	for _, c := range cands {
		for name, e := range c.mb.Benchmarks {
			if markUnstable(name, matched) {
				continue
			}
			a := byName[name]
			if a == nil {
				a = &acc{}
				byName[name] = a
			}
			// Prefer the recorded ratio; fall back to deriving it against the
			// snapshot's OWN anchor, which is the normalization it was captured
			// under.
			r := e.RatioToAnchor
			if r <= 0 && c.mb.Anchor.NSPerOp > 0 {
				r = e.NSPerOp / c.mb.Anchor.NSPerOp
			}
			if r <= 0 {
				continue
			}
			a.ratios = append(a.ratios, r)
			a.bytes = append(a.bytes, float64(e.BytesPerOp))
			a.allocs = append(a.allocs, float64(e.AllocsPerOp))
			if n := maxIterations(e); n > 0 {
				a.iters = append(a.iters, float64(n))
			}
		}
	}
	for _, prefix := range unstableBenchmarks {
		if !matched[prefix] {
			rep.unmatched = append(rep.unmatched, prefix)
		}
	}

	// A benchmark must appear in at least half the surviving window to be
	// seeded. One that does not is either newly added or intermittent, and
	// seeding it from a single observation is the failure this whole path
	// exists to avoid; `check` reports it as NEW instead, which is honest.
	quorum := (len(cands) + 1) / 2

	for _, name := range sortedKeys(byName) {
		a := byName[name]
		if len(a.ratios) < quorum {
			rep.thin = append(rep.thin, fmt.Sprintf("%s (%d/%d)", name, len(a.ratios), len(cands)))
			continue
		}
		ratio := medianOf(a.ratios)
		rep.entries[name] = BenchmarkEntry{
			NSPerOp:       ratio * anchorNs,
			BytesPerOp:    int64(medianOf(a.bytes)),
			AllocsPerOp:   int64(medianOf(a.allocs)),
			RatioToAnchor: ratio,
		}

		// b.N is an OUTPUT of Go's timing loop (N ≈ benchtime / per-op cost),
		// so it moves whenever per-op cost moves. Where iterations share state,
		// ns/op is genuinely N-dependent, which makes two snapshots taken at
		// different N two different protocols rather than two samples. Report
		// it: acting on it would mean silently shrinking the gate.
		if len(a.iters) > 1 {
			lo, hi := minMax(a.iters)
			if lo > 0 && opt.iterationTolerance > 0 && (hi/lo-1) > opt.iterationTolerance {
				rep.moved = append(rep.moved, fmt.Sprintf("%s (b.N %.0f→%.0f, %+.0f%%)", name, lo, hi, (hi/lo-1)*100))
			}
			if med := medianOf(a.iters); opt.minIterations > 0 && med < float64(opt.minIterations) {
				rep.lowN = append(rep.lowN, fmt.Sprintf("%s (b.N %.0f)", name, med))
			}
		}
	}
	return rep
}

// report prints what the reduction noticed. None of it changes the baseline —
// it changes whether someone reading the seed knows what is in it.
func (r seedReport) report(slug string, opt seedOptions) {
	if len(r.thin) > 0 {
		fmt.Printf("    %d benchmark(s) below quorum, not seeded: %s\n",
			len(r.thin), strings.Join(truncate(r.thin, 3), ", "))
	}
	if len(r.moved) > 0 {
		fmt.Printf("    b.N moved over %.0f%% across the window for %d benchmark(s): %s\n",
			opt.iterationTolerance*100, len(r.moved), strings.Join(truncate(r.moved, 3), ", "))
	}
	if len(r.lowN) > 0 {
		fmt.Printf("    NOTE: %d benchmark(s) under b.N=%d, seeded anyway: %s\n",
			len(r.lowN), opt.minIterations, strings.Join(truncate(r.lowN, 3), ", "))
	}
	if len(r.unmatched) > 0 {
		fmt.Printf("    WARNING: %d unstable-benchmark prefix(es) matched nothing in %s: %s\n",
			len(r.unmatched), slug, strings.Join(r.unmatched, ", "))
		fmt.Printf("      the exclusion list may have gone stale against a rename.\n")
	}
}

// markUnstable records EVERY exclusion prefix that matches name and reports
// whether any did.
//
// Every prefix, not just the first: "…BenchmarkClojureTestSuite" is itself a
// prefix of "…BenchmarkClojureTestSuiteCompileAndRun", so returning on the
// first hit leaves the longer entry permanently unmarked and the staleness
// check then reports a live exclusion as matching nothing on every run.
func markUnstable(name string, matched map[string]bool) bool {
	any := false
	for _, prefix := range unstableBenchmarks {
		if strings.HasPrefix(name, prefix) {
			matched[prefix] = true
			any = true
		}
	}
	return any
}

// maxIterations returns the largest b.N recorded for an entry, preferring the
// raw samples when the snapshot retained them.
func maxIterations(e BenchmarkEntry) int64 {
	var n int64
	for _, s := range e.Samples {
		if s.Iterations > n {
			n = s.Iterations
		}
	}
	return n
}

// medianOf is a plain median. Unlike reduceSamples it does NOT discard a
// warmup rep: that adjustment exists for chronological repetitions inside one
// `go test` process, and snapshots from separate CI runs have no such ordering
// — the first is not systematically the slowest.
func medianOf(vals []float64) float64 {
	if len(vals) == 0 {
		return 0
	}
	s := append([]float64(nil), vals...)
	sort.Float64s(s)
	if n := len(s); n%2 == 1 {
		return s[n/2]
	}
	n := len(s)
	return (s[n/2-1] + s[n/2]) / 2
}

func minMax(vals []float64) (float64, float64) {
	lo, hi := vals[0], vals[0]
	for _, v := range vals[1:] {
		if v < lo {
			lo = v
		}
		if v > hi {
			hi = v
		}
	}
	return lo, hi
}

func mapf[T any](in []T, f func(T) float64) []float64 {
	out := make([]float64, 0, len(in))
	for _, v := range in {
		out = append(out, f(v))
	}
	return out
}

func sortedKeys[V any](m map[string]V) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func truncate(in []string, n int) []string {
	if len(in) <= n {
		return in
	}
	return append(append([]string(nil), in[:n]...), fmt.Sprintf("… +%d more", len(in)-n))
}
