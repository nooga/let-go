---
status: active
last-verified: 2026-09-17
human-verified:
---

# let-go — binary size & boot-time audit (v1.7.4 → main @ ed4ecc2)

An independent, controlled re-measurement of the size and startup deltas
between `v1.7.4` and `main` (`ed4ecc2`), prompted by a casual comparison
table that reported a "+40% startup / +24% binary" regression.

Size and boot were measured 2026-07-17 on darwin/arm64, go1.26.3, single
machine; compute was re-run 2026-09-17 as noted in §3. Absolute figures are
machine-specific — **trust the ratios**, not the raw times and sizes. Every
number here is reproducible with the scripts in `scripts/` (see
[Reproduce](#reproduce)). Raw data is in `data/`.

---

## TL;DR

- **The shipped binary grew +2.57 MiB (+25%)** — 10.07 → 12.64 MiB stripped.
  The original "18.2 MB / +24%" was an **unstripped** build; `-s -w` (what
  ships) is 12.64 MiB. About 96% of the growth is outside the embedded core
  bundle; package attribution points to Go code.
- **Boot regressed ~7.4× in-process** (744 → 5499 µs) by v1.11.0, then was
  **fixed in v1.11.1** (5499 → 1145 µs) by making the `ir.*` pipeline lazy.
  `ed4ecc2` sits at +46% of the v1.7.4 baseline.
- **Boot tracks the core bundle, not the Go binary.** The regression was `ir.*`
  being eagerly baked into `core_compiled.lgb` and decoded on every boot.
- **It never showed on the perf page** because boot (`InitFromLGB`) wasn't a
  tracked metric until #355 — which landed the same day as the fix, one day
  after the peak. The dashboard's boot series effectively starts post-fix.
- **fib/tak are slower.** A corrected, paired A/B run alternating `AB`/`BA`
  order shows `ed4ecc2` **+12.8%** on fib and **+17.6%** on tak by paired median.
- **The audited ratchet gated boot timing but not allocations.** That gap was
  subsequently fixed by #780; current `main` gates deterministic allocs/B too.
- **Available lever:** tagging `gogen` out of the default build saves **696 KiB
  (5.4%)** — with a trade-off (see §1).

---

## 1. Binary size

### Stripped vs unstripped (reconcile the headline)

|                          | v1.7.4   | `ed4ecc2` | Δ         |
| ------------------------ | -------- | -------- | --------- |
| unstripped (dev build)   | 14.65 MiB | 18.02 MiB | +3.37 MiB  |
| **stripped `-s -w`** (ships) | **10.07 MiB** | **12.64 MiB** | **+2.57 MiB / +25%** |

The core bundle grew 130 → 242 KiB, only about 4.3% of the +2.57 MiB. Package
attribution below points to Go code for most of the remaining growth. The extra
~0.8 MiB in the unstripped number is DWARF/debug that doesn't ship.

### What grew (per-package, `nm-package-diff.py` on linux-ELF builds)

| package                                             | Δ            |
| --------------------------------------------------- | ------------ |
| `pkg/vm`                                             | +642 KiB (~3×) |
| `pkg/rt`                                             | +435 KiB (~2×) |
| `go/ast`+`parser`+`printer`+`token`+`scanner`+`doc` | +230 KiB (new) |
| `chzyer/readline`                                   | +78 KiB       |
| `runtime/pprof`                                      | +33 KiB       |
| `zeebo/xxh3`                                         | +24 KiB       |

`pkg/vm` + `pkg/rt` (~1.08 MiB) is the register-VM / numeric-op / IR-execution /
type-inference / added-builtins mass — intended. The `go/ast` cluster is the
notable one: it's the Go source-generation toolchain, pulled in **unconditionally**
by `pkg/rt/gogen.go`.

### The gogen lever (measured)

`gogen.go` (`go/ast`+`format`+`parser`+`token`) has no build tag, so it ships in
every binary — including wasm. Tagging it `//go:build bootstrap` (patch:
`gogen-bootstrap-gate.patch`):

- **saves 696 KiB (5.4%)** stripped, builds clean, normal programs boot & run;
- **but breaks runtime `*ir-compile*`** — its lazy `(require 'ir.passes.pipeline)`
  transitively needs `ir.lower-go` → `gogen`.

So it's a **trade-off, not a freebie**. Two options:
- **A (quick):** accept the tag; drops runtime `*ir-compile*` from the stock
  binary (defensible if it's treated as experimental — the daily/xsofy AOT path
  uses `-tags bootstrap` codegen, not the runtime optimizer).
- **B (full win):** decouple `ir.lower-go` from `pipeline.lg`'s unconditional
  `:require` so only the `:go` (AOT-to-native) target pulls gogen; then gate
  gogen+lower-go behind `bootstrap`. Keeps runtime bytecode ir-compile.

### Growth curve

Sampled per-commit walk (`binary-size-walk.sh`, `data/binary-size-walk.csv`):
the biggest single step is the **SSA IR / pass-pipeline landing** (late May,
`896316d8`); steady feature accretion through June; a July cleanup halved the
bundle (694 → 228 KiB) but Go code climbed back, so net stayed ~12.6 MiB.

---

## 2. Boot time

### "Startup" means three different things

| "startup"                          | measures                                | magnitude | role              |
| ---------------------------------- | --------------------------------------- | --------- | ----------------- |
| `benchmark/run.sh` (hyperfine)     | real process launch (`lg -e nil`)       | ~7–10 ms  | comparison table  |
| `bench-ratchet` **InitFromLGB**    | in-proc decode-bundle + eval, ns/op     | ~1 ms     | **the CI gate**   |
| `InitFromSource`                   | compile core from `.lg` source          | ~13 ms    | not a shipped path |

The wall-clock number carries a ~5 ms fixed OS-spawn + Go-runtime floor; the
in-proc number is the ~9×-smaller variable part that actually regressed. The
casual table's "9.9 ms" was a single-run wall-clock `time` on another machine —
untrusted by construction (its magnitude is wall-clock, its direction is noise).

### Per-release (`release-compare.sh` + `boot-e2e.sh`)

| release | date  | binary   | bundle | in-proc boot | e2e min | e2e mean |
| ------- | ----- | -------- | ------ | ------------ | ------- | -------- |
| v1.7.4  | 05-13 | 10.07 MiB | 130 KiB | 744 µs       | 5.39 ms | 7.95 ms  |
| v1.8.0  | 05-22 | 10.88 MiB | 145 KiB | 842 µs       | 5.20 ms | 6.08 ms  |
| v1.9.0  | 05-31 | 11.84 MiB | 425 KiB | 2852 µs      | 8.57 ms | 10.2 ms  |
| v1.10.0 | 06-08 | 12.15 MiB | 528 KiB | 3883 µs      | 9.73 ms | 10.72 ms |
| v1.11.0 | 06-28 | 12.64 MiB | 694 KiB | **5499 µs**  | 11.53 ms| 13.03 ms |
| v1.11.1 | 06-29 | 12.23 MiB | 228 KiB | **1145 µs**  | 5.90 ms | 7.00 ms  |
| `ed4ecc2` | 07-15 | 12.64 MiB | 242 KiB | 1085 µs    | 5.93 ms | 6.90 ms  |

Boot tracks the **bundle** almost perfectly (694 KiB → 5499 µs; 228 KiB → 1145 µs),
not the Go binary. Root cause: the `ir.*` SSA pipeline (v1.9.0) was eagerly
bundled and decoded+evaluated on every boot; v1.11.1 made it lazy-loaded. Same
cause bloated the bundle *and* boot; the fix reclaimed the boot half.

### Why it didn't surface on the perf page

`BenchmarkInitFromLGB` was added to the ratchet by **#355 on 2026-06-29** — one
day after the v1.11.0 peak (06-28) and the same day as the v1.11.1 fix. Before
that the timeline suite was explicitly "execution-only" (#190): reduce/fib/tak,
no boot. The entire regression window (v1.9.0→v1.11.0) had **no boot metric**.
The sensor was installed after the spike it would have caught.

## 3. Compute (fib / tak) — the controlled re-run

Balanced paired A/B, both stripped binaries, frozen workload
(`ab-compute.sh`, `data/ab-fib-tak.csv`). Each of 15 rounds measures both refs,
alternating order `AB`, `BA`, … so time-dependent machine drift is balanced.
This section was re-run 2026-09-17 on darwin/arm64 with go1.26.5 after correcting
the original block-ordered Hyperfine invocation:

| workload | v1.7.4 median | `ed4ecc2` median | paired B/A mean ± SD | paired B/A median |
| -------- | ------------- | ----------------- | -------------------- | ----------------- |
| fib 35   | 3.032 s       | 3.249 s           | 1.12× ±0.18          | **1.128× (+12.8%)** |
| tak      | 2.135 s       | 2.512 s           | 1.17× ±0.15          | **1.176× (+17.6%)** |

Both workloads are slower at `ed4ecc2`. The paired medians put fib below the
casual table's "+17–19%" claim and tak within it. The paired raw ratios and exact
round order are stored in `data/ab-fib.json` and `data/ab-tak.json`.

---

## 4. The ratchet — audited state and current state

- At the audited ref (`ed4ecc2`), `InitFromLGB` was gated (since #355) in the
  `pr-fast` + `full` profiles, but pass/fail used only the anchor-normalized
  **ns/op** ratio. `AllocsPerOp` and `BytesPerOp` were captured but not gated.
- **Current `main`: fixed by #780 on 2026-09-02.** `bench-ratchet check` now
  builds a machine-independent global-min bar and calls `compareDeterministic`
  with a 2% budget for allocs/op and bytes/op, even when the machine has no
  timing profile. The timing gate remains anchor-normalized and machine-scoped.

---

## 5. Recommendations

1. **Shipped in #780 (2026-09-02): gate init on allocs/B**, not just the
   ns-ratio. This now catches the eager-bundling class with deterministic
   metrics as well as timing.
2. **Add a bundle-size + binary-size ratchet.** Boot is guarded now; the *size*
   half isn't a benchmark at all, so the +2.57 MiB Go growth is ungated. By the
   audited ref, the bundle had already crept 228 → 242 KiB.
3. **Pick one canonical "startup"** and label the rest. Wall-clock is honest but
   un-gateable (process-floor noise) → keep informational; quote in-proc
   `InitFromLGB` when describing what CI protects.
4. **Decide gogen (§1).** Option B keeps runtime ir-compile and still lands most
   of the 696 KiB.
5. **Backfill the boot series** so the dashboard shows the real v1.9.0→v1.11.0
   spike (`data/release-compare.csv` is that history; #236 has a backfill
   dispatch).

---

## Reproduce

All scripts create their own throwaway worktree and never touch your checkout.

```
scripts/release-compare.sh                 # binary + bundle + in-proc boot per release
scripts/boot-e2e.sh                        # hyperfine wall-clock boot per release  (needs hyperfine)
scripts/binary-size-walk.sh v1.7.4 ed4ecc215 25 # stripped-size growth curve, sampled
scripts/ab-compute.sh v1.7.4 ed4ecc215 benchmark/fib.clj # paired AB/BA A/B
# per-package attribution (linux ELF for reliable nm sizes):
GOOS=linux GOARCH=arm64 go build -o /tmp/old . && go tool nm -size /tmp/old > /tmp/old.txt   # at v1.7.4
GOOS=linux GOARCH=arm64 go build -o /tmp/new . && go tool nm -size /tmp/new > /tmp/new.txt   # at ed4ecc215
scripts/nm-package-diff.py /tmp/old.txt /tmp/new.txt
```

## Files

- `scripts/` — the six tools + `lib.sh`
- `data/` — CSVs for every table above + paired A/B JSON
- `gogen-bootstrap-gate.patch` — the §1 lever (apply with `git apply`)

*Provenance note: this audit was prompted by a screenshot table from another
machine (single-run `time`); its startup/fib/tak magnitudes were not
reproducible as stated. The numbers here supersede it.*
