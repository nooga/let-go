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
