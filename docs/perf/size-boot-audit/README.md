# let-go — binary size & boot-time audit (v1.7.4 → main)

An independent, controlled re-measurement of the size and startup deltas
between `v1.7.4` and `main` (`ed4ecc2`), prompted by a casual comparison
table that reported a "+40% startup / +24% binary" regression.

**Measured on** darwin/arm64, go1.26.3, single machine. Absolute figures are
machine-specific — **trust the ratios**, not the raw ms/MB. Every number here
is reproducible with the scripts in `scripts/` (see [Reproduce](#reproduce)).
Raw data in `data/`.

---

## TL;DR

- **The shipped binary grew +2.57 MB (+25%)** — 10.07 → 12.64 MB stripped. The
  original "18.2 MB / +24%" was an **unstripped** build; `-s -w` (what ships)
  is 12.64 MB. Growth is ~97% Go code, not the embedded core bundle.
- **Boot regressed ~7.4× in-process** (744 → 5499 µs) by v1.11.0, then was
  **fixed in v1.11.1** (5499 → 1145 µs) by making the `ir.*` pipeline lazy.
  `main` sits at +46% of the v1.7.4 baseline.
- **Boot tracks the core bundle, not the Go binary.** The regression was `ir.*`
  being eagerly baked into `core_compiled.lgb` and decoded on every boot.
- **It never showed on the perf page** because boot (`InitFromLGB`) wasn't a
  tracked metric until #355 — which landed the same day as the fix, one day
  after the peak. The dashboard's boot series effectively starts post-fix.
- **fib/tak: real but modest.** Controlled interleaved A/B shows main ~**+11%**
  slower on both (min-ratio) — a genuine per-call overhead regression, but
  smaller than the casual "+17–19%".
- **The ratchet gates in-proc boot, but on the noisy ns-ratio, not the
  deterministic allocs/B** its own comment relies on.
- **Available lever:** tagging `gogen` out of the default build saves **696 KB
  (5.4%)** — with a trade-off (see §1).
- **Bundle compression (`feat/lgb-compression-core`) trades boot for size:**
  DEFLATE shrinks the bundle 3.2× (238→75 KB) but costs **~1.4 ms coldstart**
  (in-proc boot ~doubles). Poor trade for a CLI — motivates random-access zip
  (decompress only what's touched; skip lazy sourceinfo). See §2.5.

---

## 1. Binary size

### Stripped vs unstripped (reconcile the headline)

|                          | v1.7.4   | main     | Δ         |
| ------------------------ | -------- | -------- | --------- |
| unstripped (dev build)   | 14.65 MB | 18.02 MB | +3.37 MB  |
| **stripped `-s -w`** (ships) | **10.07 MB** | **12.64 MB** | **+2.57 MB / +25%** |

The core bundle only grew 130 → 248 KB, so essentially all of the +2.57 MB is
**Go code**. The extra ~0.8 MB in the unstripped number is DWARF/debug that
doesn't ship.

### What grew (per-package, `nm-package-diff.py` on linux-ELF builds)

| package                                             | Δ            |
| --------------------------------------------------- | ------------ |
| `pkg/vm`                                             | +642 KB (~3×) |
| `pkg/rt`                                             | +435 KB (~2×) |
| `go/ast`+`parser`+`printer`+`token`+`scanner`+`doc` | +230 KB (new) |
| `chzyer/readline`                                   | +78 KB       |
| `runtime/pprof`                                      | +33 KB       |
| `zeebo/xxh3`                                         | +24 KB       |

`pkg/vm` + `pkg/rt` (~1.08 MB) is the register-VM / numeric-op / IR-execution /
type-inference / added-builtins mass — intended. The `go/ast` cluster is the
notable one: it's the Go source-generation toolchain, pulled in **unconditionally**
by `pkg/rt/gogen.go`.

### The gogen lever (measured)

`gogen.go` (`go/ast`+`format`+`parser`+`token`) has no build tag, so it ships in
every binary — including wasm. Tagging it `//go:build bootstrap` (patch:
`gogen-bootstrap-gate.patch`):

- **saves 696 KB (5.4%)** stripped, builds clean, normal programs boot & run;
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
bundle (694 → 228 KB) but Go code climbed back, so net stayed ~12.6 MB.

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
| v1.7.4  | 05-13 | 10.07 MB | 130 KB | 744 µs       | 5.39 ms | 7.95 ms  |
| v1.8.0  | 05-22 | 10.88 MB | 145 KB | 842 µs       | 5.20 ms | 6.08 ms  |
| v1.9.0  | 05-31 | 11.84 MB | 425 KB | 2852 µs      | 8.57 ms | 10.2 ms  |
| v1.10.0 | 06-08 | 12.15 MB | 528 KB | 3883 µs      | 9.73 ms | 10.72 ms |
| v1.11.0 | 06-28 | 12.64 MB | 694 KB | **5499 µs**  | 11.53 ms| 13.03 ms |
| v1.11.1 | 06-29 | 12.23 MB | 228 KB | **1145 µs**  | 5.90 ms | 7.00 ms  |
| main    | 07-15 | 12.64 MB | 242 KB | 1085 µs      | 5.93 ms | 6.90 ms  |

Boot tracks the **bundle** almost perfectly (694 KB → 5499 µs; 228 KB → 1145 µs),
not the Go binary. Root cause: the `ir.*` SSA pipeline (v1.9.0) was eagerly
bundled and decoded+evaluated on every boot; v1.11.1 made it lazy-loaded. Same
cause bloated the bundle *and* boot; the fix reclaimed the boot half.

### Why it didn't surface on the perf page

`BenchmarkInitFromLGB` was added to the ratchet by **#355 on 2026-06-29** — one
day after the v1.11.0 peak (06-28) and the same day as the v1.11.1 fix. Before
that the timeline suite was explicitly "execution-only" (#190): reduce/fib/tak,
no boot. The entire regression window (v1.9.0→v1.11.0) had **no boot metric**.
The sensor was installed after the spike it would have caught.

### 2.5 Compression & coldstart (`feat/lgb-compression-core`)

The compression branch adds an opt-in `lgbgen --compress` (DEFLATE the bundle
body, default off; `FlagCompressed` in the decoder). Same-branch A/B, compressed
vs uncompressed bundle (`data/compression-ab.csv`), isolates the compression
effect from branch drift:

| bundle       | bundle size | binary   | in-proc boot | e2e coldstart |
| ------------ | ----------- | -------- | ------------ | ------------- |
| uncompressed | 238 KB      | 12.48 MB | 1139 µs      | 6.8 ms        |
| DEFLATE      | **75 KB**   | 12.32 MB | **2171 µs**  | **8.2 ms**    |

- **Size win:** bundle 3.2× smaller; binary −163 KB (−1.3%).
- **Boot cost:** in-proc ~doubles (+1.0 ms); e2e coldstart +1.4 ms (plain 1.20×
  faster) — matching the `--compress` flag's own "~1.4 ms slower boot" note.

For a CLI, paying +1.4 ms coldstart on every invocation to save 160 KB is a poor
trade. The better direction (maintainer's) is **random-access zip**: store the
bundle as a zip and inflate only the entries actually decoded at boot. Since
sourceinfo is already lazy, boot would never inflate it — potentially landing
the size win *and* a boot speedup vs today's uncompressed decode. **Not yet
measured** — needs the random-access format; this table is whole-bundle DEFLATE
only.

---

## 3. Compute (fib / tak) — the controlled re-run

Interleaved hyperfine A/B, both stripped binaries, frozen workload
(`ab-compute.sh`, `data/ab-fib-tak.csv`):

| workload | v1.7.4 (min) | main (min) | mean-ratio        | min-ratio |
| -------- | ------------ | ---------- | ----------------- | --------- |
| fib 35   | 1.933 s      | 2.149 s    | 1.17× (+17%) ±0.10 | +11.1%   |
| tak      | 1.950 s      | 2.174 s    | 1.07× (+7%) ±0.08  | +11.5%   |

Both are reproducibly slower on `main` — a **genuine ~+11% per-call overhead
regression** (min-ratio is the stable read; the mean-ratios diverge only from
asymmetric run-to-run noise). Direction confirms the casual table; magnitude is
smaller than its "+17–19%". Consistent with the `pkg/vm` dispatch growth. Moves
opposite to `reduce` (−47%), i.e. the classic call-overhead-for-iteration trade.

---

## 4. The ratchet — what's actually guarded

- `InitFromLGB` **is** gated (since #355), in the `pr-fast` + `full` profiles.
- But `compareAndReport` flags a regression only on **`RatioToAnchor`** (an
  anchor-normalized **ns/op** ratio, 5% budget / 12% in pr-fast). It does **not**
  gate `AllocsPerOp` / `BytesPerOp`, despite the code comment at
  `cmd/bench-ratchet/main.go:104` selling the guard on their determinism. Those
  fields are captured (`-benchmem`) and used in baseline-tightening, not in the
  pass/fail check.
- **Connection to the A/B-median work:** the median A/B runner exists precisely
  because absolute time is noisy → compare interleaved anchor-relative medians.
  But that robust method lives in the **informational** perf-page path, while the
  **hard gate** uses a single-baseline ns-ratio (no interleaving, no median).
  Boot is guarded by the *less* robust of the two mechanisms the team built.

---

## 5. Recommendations

1. **Gate init on allocs/B**, not just the ns-ratio — makes the guard match its
   own comment and reliably catches the eager-bundling class (allocs scale with
   bundled content: main InitFromLGB is 1.0 MB / 8938 allocs per boot).
2. **Add a bundle-size + binary-size ratchet.** Boot is guarded now; the *size*
   half isn't a benchmark at all, so the +2.57 MB Go growth is ungated. `main`'s
   bundle already crept 228 → 242 KB.
3. **Pick one canonical "startup"** and label the rest. Wall-clock is honest but
   un-gateable (process-floor noise) → keep informational; quote in-proc
   `InitFromLGB` when describing what CI protects.
4. **Decide gogen (§1).** Option B keeps runtime ir-compile and still lands most
   of the 696 KB.
5. **Backfill the boot series** so the dashboard shows the real v1.9.0→v1.11.0
   spike (`data/release-compare.csv` is that history; #236 has a backfill
   dispatch).
6. **Prefer random-access zip over whole-bundle DEFLATE** for the core bundle
   (§2.5): whole-bundle compression costs +1.4 ms coldstart for ~160 KB;
   random-access + lazy sourceinfo could get the size win without the boot tax.
   Worth an A/B once the format exists.

---

## Reproduce

All scripts create their own throwaway worktree and never touch your checkout.

```
scripts/release-compare.sh                 # binary + bundle + in-proc boot per release
scripts/boot-e2e.sh                        # hyperfine wall-clock boot per release  (needs hyperfine)
scripts/binary-size-walk.sh v1.7.4 main 25 # stripped-size growth curve, sampled
scripts/ab-compute.sh v1.7.4 main benchmark/fib.clj   # interleaved A/B (needs hyperfine)
# per-package attribution (linux ELF for reliable nm sizes):
GOOS=linux GOARCH=arm64 go build -o /tmp/old . && go tool nm -size /tmp/old > /tmp/old.txt   # at v1.7.4
GOOS=linux GOARCH=arm64 go build -o /tmp/new . && go tool nm -size /tmp/new > /tmp/new.txt   # at main
scripts/nm-package-diff.py /tmp/old.txt /tmp/new.txt
```

## Files

- `scripts/` — the five tools + `lib.sh`
- `data/` — CSVs for every table above + hyperfine JSON
- `gogen-bootstrap-gate.patch` — the §1 lever (apply with `git apply`)

*Provenance note: this audit was prompted by a screenshot table from another
machine (single-run `time`); its startup/fib/tak magnitudes were not
reproducible as stated. The numbers here supersede it.*
