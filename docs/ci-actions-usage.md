---
status: active
last-verified: 2026-09-14
authoritative-for:
  - ci-actions-usage
human-verified:
---

# Where CI runner time goes

`scripts/ci-usage.sh` answers that from the Actions API. This doc records what
it reported for the week of **2026-09-07 to 2026-09-14** and what follows from
it. Numbers below are that window; re-run the script for current ones.

## Why a script and not the metrics page

The `/actions/metrics/usage` page has no REST equivalent for this repo. The
Actions usage-metrics API is org- and enterprise-scoped, and `nooga` is a user
account, so `repos/nooga/let-go/actions/metrics/usage` returns 404; the billing
endpoints need the owner's own token. The runs + jobs API is what is left.

It is also the better source, because it carries a field the metrics page never
shows: each job's **queue wait**, `created_at` to `started_at`. That is the
measurement that tells a shared-capacity problem apart from a self-inflicted
one, and the two have opposite fixes.

## The week

689 runs, 2,989 jobs, **9,957 runner-minutes** (166 hours). Two workflows are
97% of it:

| Workflow | Runner-minutes | Jobs |
|---|---:|---:|
| Perf Timeline | 6,112 (61%) | 56 |
| Go | 3,544 (36%) | 2,013 |
| CodeQL | 243 | 145 |
| Deploy wasm to Pages | 58 | 134 |

Split by runner: 6,233 minutes across 2,321 ubuntu jobs, and 3,723 minutes
across **28** macOS jobs. Those 28 jobs are all Perf Timeline.

Run counts on the metrics page overstate load considerably. 285 of the 689 runs
(41%) are Perf PR, Perf PR (repeat A/B), and Perf WASM, which between them
skipped **every** job in the window — they are label-gated and cost nothing.
Judge load by minutes, not runs.

## 1. Perf Timeline is the repo's CI budget

It triggers on every push to `main`, two matrix legs, and the legs cost
**~3.6 runner-hours per push** together. Per-week totals track commit volume
directly:

| Week of | Timeline runs | ≈ runner-hours |
|---|---:|---:|
| 2026-08-10 | 15 | 54 |
| 2026-08-17 | 11 | 39 |
| 2026-08-24 | 6 | 21 |
| 2026-08-31 | 27 | 98 |
| 2026-09-07 | 29 | 105 |

The workflow is behaving as designed. The question is whether per-commit granularity is worth 61% of all runner time,
given that the header comment already treats gaps as acceptable — "the explorer
plots by commit time, and backfill can densify any range on demand."

## 2. Both legs now run at their timeout ceiling

| Leg | `timeout-minutes` | Successes | p50 | Slowest success | Killed at the cap |
|---|---:|---:|---:|---:|---:|
| ubuntu-latest | 90 | 21 | 84 min | 89 min | **7 of 28 (25%)** |
| macos-14 | 150 | 27 | 133 min | 146 min | 1 of 28 |

All seven ubuntu cancellations ran exactly 90.0 minutes, so they are timeout
kills, not queue cancels. A quarter of the amd64 timeline is being dropped, and
each drop still costs a full 90 minutes. Together the two legs burned 780
minutes in the window producing no snapshot.

The slowest *successful* ubuntu leg finished with one minute of headroom. This
is the same failure #581 item 6 described for macOS, and that #573/#583 fixed by
raising the macOS ceiling to 150 — the ubuntu leg has since grown into its own.
Raising the ubuntu timeout or trimming its workload is the immediate fix; the
underlying question is that the leg's runtime is drifting upward and nothing
watches it.

## 3. The queue wait is self-serialization, not runner scarcity

Measured queue waits:

```
Perf Timeline / macos    median=127m   max=1147m (19h)
Perf Timeline / ubuntu   median=66m    max=631m
Go / ubuntu              median=0m     max=1m
CodeQL / ubuntu          median=0m     max=0m
Deploy wasm to Pages     median=0m     max=0m
```

Everything outside Perf Timeline starts within a minute, which rules out an
account-wide concurrency cap. The direct evidence is stronger than that: lining
the macOS legs up in start order, on all 19 legs that were already queued when
the previous leg finished, the next one started **6 to 18 seconds later**. A
scarce runner pool would show idle gaps there. There are none. The queue is saturated by `concurrency: {group:
perf-timeline-<os>, cancel-in-progress: false, queue: max}`, which is one server
per OS by construction.

That is fine while arrivals are spaced further apart than the 133-minute service
time, and it collapses when they are not. On 2026-09-11, 13 commits landed on
`main` between 17:35 and 23:12; the macOS queue did not drain until
**2026-09-12T22:49**, and the last leg in that burst waited 19 hours to start.
GitHub cancels jobs queued for 24 hours, which is the mechanism #581 item 6
originally suspected — it is reachable, just not from scarcity.

The serialization buys ordered snapshots and pays for them by turning a busy
day on `main` into a day-long backlog.

## 4. Does per-commit sampling buy resolution?

Section 1 leaves the cadence question open on cost alone. This answers it on
value, using the 474 snapshots on the `perf-data` branch (2026-06-04 to
2026-09-13).

The test: track how the median per-benchmark delta grows with the separation
between two snapshots. Data dominated by run-to-run noise stays flat as
separation grows. Data carrying real signal rises, because more code changed.

| Profile | lag 1 | lag 2 | lag 5 | lag 10 | lag 20 |
|---|---:|---:|---:|---:|---:|
| amd64 EPYC 7763 | 1.11% | 1.28% | 1.56% | 1.83% | 2.25% |
| amd64 EPYC 9V74 | 1.24% | 1.31% | 1.69% | 1.92% | 2.71% |
| arm64 Apple M1 | 11.01% | 11.28% | 11.44% | 11.97% | 12.03% |

**The arm64 leg is noise-dominated.** Comparing commits 20 apart tells you 9%
more than comparing adjacent ones, against 103% on amd64. Its resolution is
set by run-to-run variance rather than by the code under test, so per-commit
sampling adds little a weekly cadence would not also give. That leg is the
most expensive thing the repo runs: 133 minutes per leg, 37% of all runner
time in the section-1 window, and the sole cause of the queue backlog in
section 3.

The amd64 lanes carry real signal, but at lag 1 the delta sits on the noise
floor; separation of 10 to 20 commits is where a change clearly clears it.
That is an argument for sampling coarsely and backfilling on detection, which
is the workflow header's own stated design.

Two structural findings came out of the same data.

**The amd64 lane is six CPU models, not one.** Snapshots since 2026-06-04
split as EPYC 7763 (204), EPYC 9V74 (100), Xeon 8370C (21), Xeon 8573C (11),
Xeon 6973P-C (8), EPYC 9V45 (2). The ratchet partitions by machine profile, so
these are six series, not one. Pushing per commit therefore does not produce a
per-commit series in any comparable lane: roughly 60% of snapshots land in the
top profile and the rest scatter into series too sparse to read.

**No snapshot has an A/A control.** Across all seven profiles, every snapshot
carries a distinct `captured_at_sha`: 474 snapshots, 474 distinct commits, zero
repeats. The timeline therefore cannot separate a regression from noise using
its own data, since it never measures one commit twice on one profile.
`perf-pr-repeat.yml` exists for exactly that measurement and skipped every job
in the section-1 window.

## 5. What runs on every PR push

`go.yml` carries all nine required status checks. Over 465 runs
(2026-08-15 to 2026-09-14) it cost 9,222 runner-minutes.

**One test step is 31% of that.** Sampled across 54 build jobs, "Expensive
lowering e2e" is 74% of the `build` job, or about 2,850 minutes a month. It
also grew 32% when `TestCustomMain` joined it, from 352s to 465s average.

Fire rates over the same 465 runs:

| Job | Failures | Minutes | Required |
|---|---:|---:|---|
| build | 43 | 3,841 | yes |
| gogen-diff | 0 | 1,481 | yes |
| generated-artifacts | 34 | 1,278 | yes |
| tinygo-wasi-build | 1 | 1,039 | yes |
| lint | 4 | 385 | yes |
| race | 1 | 374 | yes |
| wasip1-build | 0 | 230 | yes |
| no-http-build | 0 | 228 | yes |
| gold-differential | 0 | 152 | no |
| default-deps | 1 | 78 | no |
| test-location | 1 | 39 | yes |
| docs-status | 0 | 50 | no |
| docs-frontmatter | 1 | 47 | no |

A gate that has not fired is not thereby useless, and two of these are worth
naming. `gogen-diff` costs 25 runner-hours a month on a clean record, but it
is a cross-engine differential plus parity ledger, not a duplicate of
`generated-artifacts` — the class of failure it guards against is rare and
severe. `tinygo-wasi-build` is 17 runner-hours; its toolchain install is only
9 seconds, so caching does not help, and the 107-second TinyGo build is the
cost. For both, path-gating preserves the gate; removal does not.

**Path filtering is worth less than it looks.** `go.yml` has no `paths` filter,
so a docs-only PR runs the full matrix. Measured, that was 21 of 465 runs and
428 runner-minutes a month, 4.6% of the workflow. Cheap to add, small payoff.

## Ranked by payoff

1. arm64 timeline cadence: about 37% of all runner time, at close to no loss
   of information (section 4).
2. amd64 timeline cadence: most of another 24%, at a real but bounded
   resolution cost.
3. The "Expensive lowering e2e" step: 47 runner-hours a month, and the open
   question is whether it needs to run on every push.
4. `gogen-diff` path-gating: 25 runner-hours a month, gate preserved.
5. A docs path filter on `go.yml`: 7 runner-hours a month.

## What this leaves open

Trade-offs for the team, not conclusions:

- **Cadence.** A scheduled or coalescing trigger instead of `on: push` would cut
  the largest line item substantially and drain the queue, at the cost of
  timeline density that backfill is designed to restore. #581 item 6's step 2
  proposed exactly this and is still open.
- **The ubuntu ceiling.** Needs a decision now, independent of cadence: 25% of
  amd64 snapshots are being lost at full price.
- **Drift watch.** Both legs reached their ceilings by growing into them. Some
  check on leg runtime would catch the next one before it starts dropping data.
- **What the legs measure.** Neither the per-commit cadence nor the workload
  size has been re-justified against what the timeline is actually used for.
  That analysis is not in this doc.

Related: #581 (CI audit, item 6 is the ancestor of this), #573 and #583 (the
macOS timeout fix), #693 (label-opt-in benchmark lane), #752 (pages.yml
redeploys on any Timeline conclusion).
