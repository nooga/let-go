---
status: active
last-verified: 2026-08-05
authoritative-for:
  - benchmark-ratchet
human-verified:
---

# Benchmark Ratchet

A small system for catching perf regressions in CI without requiring
CI to run on the same hardware as developers. The core idea is
**anchor-relative measurement**: every benchmark is reported as a
multiple of a tight CPU loop that has no allocations and no
project-specific code. That ratio is roughly stable across machines
(within a CPU family), so a baseline captured on one machine still
flags regressions when checked on another.

## Architecture

Two phases, separately invokable:

1. **capture** — runs `go test -bench` per package. As each benchmark
   line is parsed it's appended (and fsync'd) as one JSON object to a
   `.jsonl` file. If a later benchmark hangs or panics, earlier
   results are already on disk; a per-package `-timeout` keeps any
   single test from blocking the whole sweep forever.
2. **aggregate** — reads one or more `.jsonl` files, normalizes
   against the anchor, emits the consolidated `baseline.json`.
3. **snapshot** — captures and aggregates like `show`, then writes
   the current run as an immutable JSON snapshot. This is for timeline
   graphs and does not ratchet or update `baseline.json`.

The one-shot `check` / `update` / `show` modes are just wrappers:
capture-then-aggregate-then-(compare|write|print).

## Components

| File | Purpose |
|---|---|
| `cmd/bench-ratchet/main.go` | The tool. Streams raw `.jsonl`, aggregates, compares. |
| `pkg/vm/bench_ratchet_anchor_test.go` | `BenchmarkRatchetAnchor` — frozen calibration loop, lives in pkg/vm so it's automatically in default scope. |
| `test/zz_bench_test.go` | `BenchmarkClojureTestSuite` — end-to-end wall time of the full clojure-test-suite (jank) corpus. Catches compile + run regressions that pkg/vm micro-benches don't see. |
| `docs/perf/baseline.json` | The current committed baseline. |
| `docs/perf/historical/*.json` | Frozen historical snapshots (e.g. `v1.8.0.json`). Captured against the same anchor so any current run can `-baseline docs/perf/historical/v1.8.0.json check` and see "how much have we drifted since release N." |
| `docs/perf/timeline/*.json` | Append-only full perf snapshots captured on pushes to `main`. These drive trend charts and record actual runs over time. |
| `docs/perf/index.html` | Static "Are we fast yet?" page, generated on demand by `make perf-page` from the committed baseline and historical snapshots. Gitignored: the deployed page is built fresh by `pages.yml`, which also pulls timeline snapshots from the `perf-data` branch. |
| `cmd/perf-page/main.go` | Static page generator. It renders HTML only; it never runs benchmarks. |
| `docs/perf/.runs/*.jsonl` | Raw capture output. Gitignored; recreated on each run. |
| `.github/workflows/perf-timeline.yml` | Main-only CI job that records timeline snapshots and pushes them to the `perf-data` branch. It does not touch the static page. |
| `Makefile` targets | `bench-ratchet`, `bench-ratchet-update`, `bench-ratchet-show`, `perf-snapshot`, `perf-page`. |

## Scope

Default gate: calibration anchor, BenchmarkClojureTestSuite under both
bytecode and `gogen_ir`, and the targeted `pkg/ir` BenchmarkIRCompile
under both bytecode and `gogen_ir`.

Full profile (`-full`): the broad timeline/deep-dive profile. It runs
the full `pkg/vm` benchmark fleet under `-tags gogen_ir`, plus the
suite and IR compile benchmark under both bytecode and `gogen_ir`.

The six b.N=1 `BenchmarkClojureTestSuite*` variants remain observational:
they run and appear in reports, but `update`, forced rebaseline, and
`seed-baseline` filter them from the active ratchet because repeated samples
mutate shared runtime state and are too noisy for a blocking timing bar.

> **Reading the `gogen_ir` / `aot_native` numbers.** These exercise the
> natively-lowered IR passes (the dispatch is guarded by
> `TestIRBenchDispatchesNativeUnderTag`), but the lowered Go currently
> *boxes* most cross-namespace and `clojure.core` calls — runtime
> `LookupVar`/`Deref` + a fresh `[]vm.Value` arg slice per call — rather
> than emitting direct calls. So the `gogen_ir` variant can allocate
> **more** than `bytecode` (e.g. IRCompile ~1.7×) and is only marginally
> faster. Treat these as "native dispatch works, codegen not yet
> optimized" — not "native is fast." Direct-call lowering (cross-ns +
> IFn-typed + `clojure.core`) is the lever that would push them below the
> bytecode line.

Deliberately out of scope: `pkg/compiler` (compile time — measured by
parity scripts) and `pkg/bytecode` (decode time — measured at `make
build`). Keeping the ratchet narrow means a slow compiler change
doesn't trip a runtime gate AND each in-scope bench gets more of the
sample budget.

Override via `-packages "github.com/nooga/let-go/pkg/X github.com/.../Y"`.

## Build tags

Default: `-tags gogen_ir`. This compiles the lowered-to-Go VM
(`pkg/rt/core_go_lowered/**/*.go`) into the test binary alongside the
bytecode VM, so any runtime that consults the native-direct registry
sees the lowered path available. At releases that pre-date the
lowered-Go work (e.g. v1.8.0), the build tag matches no files and the
flag is a no-op — no special-casing needed.

Override via `-tags ""` (vanilla build) or `-tags "foo bar"` (custom).

## Usage

One-shot (Makefile):

```sh
make bench-ratchet           # check current vs baseline (CI mode)
make bench-ratchet-update    # overwrite baseline with current numbers
make bench-ratchet-show      # capture & print, write nothing
make perf-page               # render docs/perf/index.html locally from committed JSON
make perf-snapshot           # full capture into docs/perf/timeline/<ts>-<sha>.json
```

Explicit two-phase (useful when you want progress visibility, or are
capturing on one machine and aggregating on another):

```sh
# Stream raw results
go run ./cmd/bench-ratchet -out /tmp/run.jsonl capture
# Aggregate to a baseline
go run ./cmd/bench-ratchet -in /tmp/run.jsonl -baseline /tmp/b.json aggregate
# Aggregate to an immutable timeline snapshot
go run ./cmd/bench-ratchet -in /tmp/run.jsonl -baseline docs/perf/timeline/<ts>-<sha>.json snapshot
```

Finer control on either phase:

```sh
go run ./cmd/bench-ratchet -budget 0.10 check                # 10% budget
go run ./cmd/bench-ratchet -filter '^BenchmarkIR' check      # subset
go run ./cmd/bench-ratchet -count 3 -benchtime 2s update     # rigorous
go run ./cmd/bench-ratchet -baseline docs/perf/historical/v1.8.0.json check
                                                             # vs v1.8.0
```

### Seeding from CI

The active baseline is seeded from CI timeline snapshots via the `seed-baseline`
command. It reduces a WINDOW of recent snapshots per machine key — not the
newest one — and merges the result with any existing local M3 profile:

```sh
# Fetch the perf-data branch containing timeline snapshots
git fetch origin perf-data

# Seed the baseline from amd64 profiles (preserving local M3)
bench-ratchet -perf-data-dir <perf-data-root>/timeline \
  -baseline docs/perf/baseline.json \
  seed-baseline
```

The command:
- Scans the timeline directory for snapshot files named `TIMESTAMP-SHORTSHA-MACHINE.json`, reporting any name it cannot parse
- Filters to one architecture (`-seed-arch`, default amd64 per #651), on the file's CONTENT as well as its name
- Groups every profile it reads by the machine key derived from the file's CONTENT, reporting a filename that names a different machine than it carries
- Takes the newest `-seed-window` snapshots (default 5) per machine key
- Skips a machine key with fewer than `-seed-min-window` snapshots (default 3): a tier that has only just started reporting has no window to disagree with, so it stays ungated until the runs accrue
- Rejects a snapshot whose `ratio_to_anchor` values sit more than `-seed-coherence-tolerance` (default 5%) off the rest of its window
- Medians each benchmark across the survivors **in ratio space**, deriving `ns_per_op` back from the window's anchor
- Skips a benchmark present in fewer than half the surviving snapshots, rather than seeding it from one observation
- Reports `b.N` movement across the window and any exclusion-list entry that matched nothing
- Preserves any existing arm64/Apple M3 profile for local developer gating
- Excludes the six BenchmarkClojureTestSuite* variants (too noisy to ratchet, per #651)

`captured_at_sha` names the newest *surviving* snapshot in the window, which is
the identity of the profile rather than the sole source of its numbers. The seed
log prints the window size and how many snapshots contributed. Every seeded
entry is stamped with that same identity, as its `best_since_sha` /
`best_since_at` and as its `allocs_since_*` / `bytes_since_*` stamps:
a deterministic floor with no commit attached is one the
[deterministic gate](#the-deterministic-gate-allocsop-bytesop) cannot place in
time, and a tier that later stops reporting would leave an unattributable row
behind.

Future work (#597, separate) will backfill per-tier v1.8.0 release-reference
snapshots.

### Why a window, and why not gate on the anchor

**One snapshot is one CI run, and one CI run is one sample.** Seeded from the
newest snapshot versus a median of five, on the same corpus (2026-08-05): 22.6%
of the 758 (tier, benchmark) floors differ by more than the 5% regression
budget, and 11.3% by more than 10%. Part of that is real code movement across
the window and part is sampling — but either way, seeding from one snapshot sets
a fifth of the gate's thresholds from a single observation of it.

**The reduction happens in ratio space** because raw `ns_per_op` carries host
speed and `ratio_to_anchor` does not. Over the 24 most recent amd64 snapshots,
anchor deviation from the tier median ranges −22.4%..+3.1% while
`ratio_to_anchor` holds to a median 0.03% and a worst 1.75%.

**The gate is on ratio coherence, not on anchor drift** — which is the obvious
design and is wrong. Snapshot `a588a69d2759` (EPYC 9V74) sits 22.4% off its
window's anchor and is uniformly 22.4% fast in raw `ns/op` across all 162 of its
benchmarks, agreeing with its window on every ratio to within 0.1%: the host was
fast and the anchor divided that back out, which is what the anchor is for.
Gating on anchor deviation would discard it and two more like it, two of which
are snapshots this baseline is seeded from.

What does need rejecting is the *mixed* capture — anchor caught the slow tail,
benchmarks did not — where every ratio is uniformly wrong while the raw numbers
look ordinary. That shows up as a whole-snapshot offset in normalized space,
well clear of the 1.75% the corpus exhibits.

This approach makes the baseline auditable (provenance is in git history), CI-sourced
(no local machine capture noise), and reproducible across runs (idempotent).

## Streaming visibility

`capture` prints per-package progress to stderr:

```
  [1/5] github.com/nooga/let-go/pkg/api ... 7 records
  [2/5] github.com/nooga/let-go/pkg/bytecode ... 4 records
  [3/5] github.com/nooga/let-go/pkg/compiler ... 24 records
  ...
```

And `tail -F docs/perf/.runs/<sha>-<ts>.jsonl` shows individual
benchmark results as they land. Each line is a self-contained
`StreamRecord`:

```jsonc
{"package":"github.com/nooga/let-go/pkg/vm",
 "name":"BenchmarkIsTruthy/int",
 "iterations":427510606,
 "ns_per_op":0.566,
 "bytes_per_op":0,
 "allocs_per_op":0,
 "captured_at":"2026-05-29T18:38:24Z"}
```

If the sweep crashes or is interrupted, the `.jsonl` survives. You
can pick up where you left off by capturing the missing packages
into a second `.jsonl` and aggregating both.

## What the baseline records

```jsonc
{
  "version": 1,
  "captured_at": "2026-05-29T...",
  "captured_at_sha": "<short-sha>",
  "machine": {
    "os": "darwin", "arch": "arm64", "num_cpu": 8,
    "cpu_model": "Apple M3", "go_version": "go1.26.3"
  },
  "anchor": {
    "name": "BenchmarkRatchetAnchor",
    "package": "github.com/nooga/let-go/pkg/api",
    "ns_per_op": 1.09,
    "iterations": 1000000000
  },
  "benchmarks": {
    "github.com/nooga/let-go/pkg/api.BenchmarkIRPipelineCompile": {
      "ns_per_op": 45590.0,
      "allocs_per_op": 720,
      "bytes_per_op": 41200,
      "ratio_to_anchor": 41862.4,
      "best_since_sha": "<short-sha>",
      "best_since_at": "2026-09-07T05:34:02Z",
      "allocs_since_sha": "<short-sha>",
      "allocs_since_at": "2026-09-07T05:34:02Z",
      "bytes_since_sha": "<short-sha>",
      "bytes_since_at": "2026-09-07T05:34:02Z"
    }
  }
}
```

- `ns_per_op` and friends are kept for human eyeballing on same-machine
  drift checks. Comparisons across machines should ignore them.
- `ratio_to_anchor` = `ns_per_op / anchor.ns_per_op`. This is what the
  `check` mode actually compares.
- `best_since_sha` / `best_since_at` date the entry as a whole: the run that
  last moved **any** of its metrics toward better. They are what the timing
  report's "best since" column names.
- `allocs_since_sha` / `allocs_since_at` date `allocs_per_op`, and
  `bytes_since_sha` / `bytes_since_at` date `bytes_per_op` — each the run that
  most recently **confirmed** that metric (measured it equal to or better than
  the stored value), and so the newest run that measured what is stored. A run
  that measures a metric worse leaves its value and stamp untouched — the
  ratchet only tightens — but a run that merely re-measures the same value
  still restamps, because it is current evidence for the code state, and the
  gate ranks by newest evidence. They exist because `best_since_*` also moves
  on a timing-only win, which would date
  a pinned allocation floor to a commit that never measured it. The two metrics
  are dated **separately** because the ratchet takes each one's minimum
  separately: a run that lowers allocs while regressing bytes leaves this run's
  allocs stored beside an older run's bytes, and one stamp over both would date
  the kept metric to a commit that never produced it. The
  [deterministic gate](#the-deterministic-gate-allocsop-bytesop) selects its
  reference per metric by these stamps, falling back to the profile's capture
  identity only for a row that was never ratcheted (see that section for the two
  cases).

## How the check works

For each benchmark in the baseline:

```
delta = (current.ratio_to_anchor / baseline.ratio_to_anchor) - 1
```

- `delta > +budget` → **REGRESSION** (counted toward non-zero exit).
- `delta < -budget` → **IMPROVED** (informational; safe to ignore or
  ratchet down by running `update`).
- `|delta| ≤ budget` → ok.

Benchmarks in the baseline that don't appear in the current run are
flagged **MISSING** (likely renamed or removed). Benchmarks in the
current run that don't appear in the baseline are flagged **NEW**.

The default budget is **5%**. Raise it for noisy benchmarks via
`-budget`. Lower it once you've improved benchmark stability (e.g.
`-benchtime 5s -count 5`).

### The deterministic gate (allocs/op, bytes/op)

`allocs/op` and `bytes/op` carry no CPU-dependent noise, so they are gated
separately from timing, at a tight 2% budget, and on every machine — including
one with no timing profile of its own.

They are a property of **the code at a commit**, not of a machine. That is what
makes them portable across profiles, and it is also why the reference is **not**
a minimum across profiles. Each profile is captured at whatever commit its
machine last ran at, and a tier keeps its last numbers indefinitely once its
runner stops reporting — so the profiles present at any moment describe several
different code states. A minimum over them answers "the least anyone has ever
measured", which is a fact about the fleet's history rather than about any one
commit, and gates current code against whichever code state happened to allocate
least.

So for each benchmark the reference is **the value with the newest provenance
among the profiles that carry it** — the most recent commit anyone measured it
at — with rows from the **same commit** reduced to their minimum, so the bar
still ratchets across repeated measurements of one code state.

`allocs_per_op` and `bytes_per_op` are selected **separately**, since they are
dated separately: a profile can hold a freshly measured allocs figure beside a
bytes figure the ratchet kept from an older run, and taking both from whichever
row won on one of them would either throw away the fresh measurement or adopt
the stale one.

Provenance is the metric's own `allocs_since_*` / `bytes_since_*` stamp when it
has one. A value predating those fields is read in one of two ways:

- **No `best_since_*` either** — the entry has never been ratcheted, so seed or
  capture wrote every one of its numbers in the run the profile records. The
  profile's `captured_at_sha` / `captured_at` is then their exact provenance and
  is used.
- **`best_since_*` set, deterministic stamp absent** — the entry HAS been
  ratcheted, and `best_since_*` moves whenever any metric improves, timing
  included, so on a row whose allocs/bytes were pinned from an older run it
  names a commit that never measured them. The profile's `captured_at` is no
  better: it dates the run that wrote the file, not the run that set a bar the
  ratchet carried forward. Such a value has **unknown** provenance: it sorts
  oldest, ties with the other undated values at their minimum, and is displaced
  by any row that can name its commit.

The reported regression line names the commit the reference came from (`—` when
unknown), so a surprising bar can be traced to the run that set it.

Same-commit is decided by the SHA, not by the clock. Two machines run one commit
whenever their queues allow, so equal timestamps neither identify a shared code
state nor are needed to; conversely two different commits that happen to be
stamped in the same second are ranked rather than mixed, since a bar assembled
from two commits describes code that never existed. A merged same-commit group
then carries the **newest** of its members' timestamps: its claim on describing
current code rests on its latest capture, not on whichever profile happened to be
read first. The timestamp orders groups across commits (a SHA cannot be ordered
without the repository, and baselines are read on machines that lack the
history), and breaks a remaining tie by SHA so the selection does not depend on
map order.

## Running the check on a PR (the `perf` label)

CI runs the A/B on a pull request only when the PR carries the **`perf`**
label — `.github/workflows/perf-pr.yml`. Add the label to run it; it
re-runs on each push while the label is on; remove it to stop. An
unlabeled PR skips the job instantly, with no runner allocated, so the
label is a deliberate opt-in for the multi-minute sweep.

Guards and escape hatches:

- **Path filter.** The job only fires for changes under `pkg/**` or
  `cmd/**` (plus the workflow file itself), so a labeled doc-only PR
  still skips.
- **Manual dispatch.** Run against any PR by number without touching
  labels via the workflow's `workflow_dispatch` input `pr` (needs Actions
  access).

The PR job builds with the same `-tags gogen_ir` default as local runs
(see [Build tags](#build-tags)), so it measures the native-lowered path.
It runs the `pr-fast` profile — the stable `pkg/vm` micro-benchmark
families (frame dispatch, allocation, invoke, arity, record ops) plus the
calibration anchor, sized for quick two-pass feedback. The end-to-end
suite and the IR-compile bench are not in the PR gate; they run in the
fuller profiles on `main`. None of the `pr-fast` families is a
numeric/float kernel, so a change that only helps tight numeric code (e.g.
native float unboxing) reads flat here until such a benchmark is added —
see [Adding a new benchmark](#adding-a-new-benchmark).

## Cross-machine sanity

When `cpu_model` or `go_version` differs between baseline and
current, the tool prints a warning. The ratio-based comparison is
still meaningful — that's the whole point of the anchor — but with
caveats:

- Memory-bound benchmarks (large map/seq sweeps) ratio against an
  L1-cache-resident anchor will look slower on machines with lower
  memory bandwidth, even when CPU is fine. Treat large
  anchor-ratio deltas on memory-heavy benches with suspicion before
  blaming code.
- A different `go_version` can shift the anchor itself if the
  compiler grows new optimizations. The anchor is designed to
  resist this (PCG constants, escape via `runtime.KeepAlive`), but
  it's not impossible.

The anchor's absolute `ns_per_op` between baseline and current is
printed at the top of every `check` report. If it has drifted a lot,
the ratio comparison is on shakier ground.

## Baseline seeding and machine-key selection

The active `docs/perf/baseline.json` is seeded from CI timeline snapshots and
gates against the incremental (newest) snapshot for real-time drift tracking:

- **amd64 profiles**: Seeded from a median over the most recent snapshots per
  amd64 machine key (e.g. AMD EPYC 7763, 9V74) to survive runner rotation. The
  gate budget is 5%; no null control has been run on these tiers, so the floor
  beneath that budget — the gap between two builds that cannot differ — is not
  yet measured here.

- **arm64/Apple M3**: Preserved from any existing local baseline, allowing M3
  developers to gate against a machine-specific baseline without CI noise.

- **arm64/Apple M1 (Virtual)**: Deliberately excluded (per #651) because CI
  runner noise (27.2% of entries) is too high for reliable ratcheting; reverts
  to deterministic-only gating at `main.go:452-459`.

The `captured_at_sha` field in each machine entry points to the incremental SHA,
representing the "current" state for `make bench-ratchet check` comparisons.

The `bench-ratchet seed-baseline` command (§ [Seeding from CI](#seeding-from-ci))
reduces a window of timeline snapshots from the perf-data branch, making the
baseline reproducible, auditable, and independent of local machine captures.

## Current baseline: amd64-seeded with M3 fallback

The active `docs/perf/baseline.json` is seeded from CI timeline snapshots via
`seed-baseline`, with amd64 as the primary machine tier and the existing local
M3 profile (arm64/Apple M3) preserved for local developer gating. This means
`make bench-ratchet` gates against a median of the most recent successful amd64
builds in CI, providing a durable, reproducible baseline free of local machine
noise.

A future release-reference baseline (#597, separate) will backfill v1.8.0
reference snapshots for each machine tier to answer "how do we compare to the
last release?"; that work will coexist with the current drift-tracking baseline.

## Updating the baseline (the ratchet)

**The ratchet only moves one direction: tighter.** `make
bench-ratchet-update` reads the existing baseline, runs the current
benchmarks, and writes a merged baseline where each `(benchmark,
metric)` is the MIN of (existing baseline, current). So:

| Outcome | What happens |
|---|---|
| Bench got faster than baseline | New (faster) value adopted. |
| Bench got more allocs but the same speed | Speed adopted (it didn't change); allocs **pinned at baseline** — the new (higher) alloc count is rejected. |
| Bench is brand new (not in baseline) | Adopted as-is. There's no prior bar to ratchet against. |
| Bench in baseline but missing from current | **Kept in baseline** — a removed/renamed benchmark shouldn't release the bar accidentally. |

`update` prints a summary of what tightened, what would-have-regressed-
but-was-pinned, what's new, and what's missing-but-kept. A pinned
entry is your signal to investigate: either there's an actual
regression hiding (run `make bench-ratchet` to see how far over),
or it's measurement noise that you can ratchet down with another
run.

### `-force`

`go run ./cmd/bench-ratchet -force update` bypasses the ratchet for the current
machine's **timing**: it writes this run's `ns_per_op` / `ratio_to_anchor` as-is,
regressions included, stamped with this run's commit, and writes no other
profile. Unmeasured entries are retained, so a fast-gate rebaseline cannot erase
full-profile history.

`allocs/op` and `bytes/op` are not part of that bargain. Timing is a property of
the host, so a machine may re-declare its own; the deterministic pair is a
property of the code at a commit and every profile is gated against it, so a
local recapture that allocates more is a regression everywhere rather than a
number this host gets to reset. A forced update therefore keeps any stored
deterministic value it measured worse, along with its
`allocs_since_*` or `bytes_since_*` stamp — each metric keeps or adopts on its
own — and prints what it declined:

```
  NOT ACCEPTED (deterministic regression; stored bar kept — re-run with -accept-deterministic to record it):
    ! pkg/ir.BenchmarkIRCompile [bytecode]   bytes/op  kept 4941030 (since 477a5d36e25f), measured 5232922
```

A forced update is gated on **both** bars, not just this machine's. A value can
be an improvement over a stale local row and still be a regression against the
newest row any profile carries; since the gate selects by newest provenance,
adopting it would stamp it as the newest evidence and raise the bar everyone is
measured against, so the regression it represents would stop being reported.
Either violation keeps the stored value and its date, and the rejection names
whichever bar was exceeded.

An improvement needs no ceremony: it is lower, so it ratchets and takes this
run's stamp. So does a value measured **equal** to the stored one: a timing
recapture that re-measures the same allocations is current evidence for them,
and dating it keeps a legacy row from losing to an older, worse profile. To record a regression deliberately, add
`-accept-deterministic`, which writes the measured pair and dates it to this run.
Accepted values reach the
[deterministic gate](#the-deterministic-gate-allocsop-bytesop) by being the
newest measurement of them, so no other tier has to be told; copying them into
the other profiles would record one machine's local capture as those tiers'
stored floor under a commit they never ran.

Use only when:

- A regression has been investigated, discussed, and accepted as
  the new floor (rare; should have a paper trail in the commit).
- Initial seeding — though a missing `baseline.json` triggers a
  one-shot write without `-force`.

What NOT to do:

- Don't `-force` to silence a regression you can't explain. The
  whole point of the system is to surface regressions; force-updating
  defeats it. Investigate first.
- Don't reach for `-accept-deterministic` to clear the "NOT ACCEPTED" lines. An
  allocation regression the gate reports on every machine needs the paper trail,
  not the flag.
- Don't bundle a baseline update with the change being measured —
  split them so the reviewer sees which numbers moved and why.

## Historical baselines

`docs/perf/historical/vX.Y.Z.json` files are **frozen**. They are
captured once at release time and never updated. The active
`docs/perf/baseline.json` *starts* as a copy of the most recent
release's historical baseline, then ratchets tighter as perf work
lands.

Compare against an older release directly:

```sh
go run ./cmd/bench-ratchet -baseline docs/perf/historical/v1.8.0.json check
```

This tells you "how much have we drifted since v1.8.0?", with no
ratchet logic involved — just a direct check.

## Adding a new benchmark

Just add `func BenchmarkX(b *testing.B)` to any `*_test.go` under
`pkg/...`. The ratchet auto-discovers any package that has a
`Benchmark*` function. On the next `update`, the new benchmark
appears in `benchmarks` and gets tracked.

## Why an anchor, not benchstat?

`benchstat` does same-machine comparison really well, but its mental
model is "two runs, compare them statistically." That doesn't survive
crossing the machine boundary. The anchor approach gives a single
portable number per benchmark, which slots into CI without the
two-run-on-the-same-machine constraint. For deeper same-machine
analysis when investigating a regression, run `benchstat` against
the raw output — they compose.

## What's NOT in scope (yet)

- Per-benchmark budget overrides (some benchmarks are inherently
  noisier than others — useful, but not implemented).
- Historical baselines / regression graphs over time.
- Auto-comment posting on PRs.

All three are easy follow-ups when the bare-bones gate gets used
enough to justify them.
