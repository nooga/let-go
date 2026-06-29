---
status: active
last-verified: 2026-06-05
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
| `docs/perf/index.html` | Static "Are we fast yet?" page generated from the committed baseline, historical snapshots, and timeline snapshots. |
| `cmd/perf-page/main.go` | Static page generator. It renders HTML only; it never runs benchmarks. |
| `docs/perf/.runs/*.jsonl` | Raw capture output. Gitignored; recreated on each run. |
| `.github/workflows/perf-timeline.yml` | Main-only CI job that records timeline snapshots and commits the regenerated static page. |
| `Makefile` targets | `bench-ratchet`, `bench-ratchet-update`, `bench-ratchet-show`, `perf-snapshot`, `perf-page`. |

## Scope

Default gate: calibration anchor, BenchmarkClojureTestSuite under both
bytecode and `gogen_ir`, and the targeted `pkg/ir` BenchmarkIRCompile
under both bytecode and `gogen_ir`.

Full profile (`-full`): the broad timeline/deep-dive profile. It runs
the full `pkg/vm` benchmark fleet under `-tags gogen_ir`, plus the
suite and IR compile benchmark under both bytecode and `gogen_ir`.

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
make perf-page               # refresh docs/perf/index.html from committed JSON
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
      "ratio_to_anchor": 41862.4
    }
  }
}
```

- `ns_per_op` and friends are kept for human eyeballing on same-machine
  drift checks. Comparisons across machines should ignore them.
- `ratio_to_anchor` = `ns_per_op / anchor.ns_per_op`. This is what the
  `check` mode actually compares.

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

## The baseline starts at v1.8.0

The active `docs/perf/baseline.json` is initialized as a copy of
`docs/perf/historical/v1.8.0.json`. This means `make bench-ratchet`
asks "how does the current code compare to the v1.8.0 release?",
not "to whatever main looked like yesterday." A regression bar
anchored to a release is meaningful; one anchored to yesterday's
main just tracks noise.

Two consequences:

- Some current benchmarks (e.g. `BenchmarkRatchetAnchor` itself,
  `BenchmarkIRPipelineCompile`) didn't exist at v1.8.0. They'll be
  flagged **NEW** in every check, which is informational, not a
  regression. Once they've matured into "things we want to gate on,"
  do a deliberate `update` to lift them into the bar.
- Some v1.8.0 benchmarks may have been renamed or removed. They'll
  appear as **MISSING**. Same treatment — `update` clears them
  out of the comparison once you confirm the removal was intentional.

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
| Bench in baseline but missing from current | **Kept in baseline** — a removed/renamed benchmark shouldn't release the bar. Use `-force` to drop it intentionally. |

`update` prints a summary of what tightened, what would-have-regressed-
but-was-pinned, what's new, and what's missing-but-kept. A pinned
entry is your signal to investigate: either there's an actual
regression hiding (run `make bench-ratchet` to see how far over),
or it's measurement noise that you can ratchet down with another
run.

### `-force`

`go run ./cmd/bench-ratchet -force update` bypasses the ratchet and
writes current numbers as-is, including any regressions. Use only
when:

- A regression has been investigated, discussed, and accepted as
  the new floor (rare; should have a paper trail in the commit).
- A benchmark was renamed or removed and you want to drop the
  historical entry.
- Initial seeding — though a missing `baseline.json` triggers a
  one-shot write without `-force`.

What NOT to do:

- Don't `-force` to silence a regression you can't explain. The
  whole point of the system is to surface regressions; force-updating
  defeats it. Investigate first.
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
