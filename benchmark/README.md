# Benchmarks

Compares the let-go bytecode VM against let-go AOT (the same code IR-lowered
to native Go), babashka, and Clojure JVM.

The VM and AOT legs run the **identical driver + fixture** — only dispatch
differs. The plain build executes bytecode; the AOT build carries the same
namespaces lowered to native Go (via the self-hosted IR pipeline) and
registered as overrides, so the same `(bench.<name>/run)` call lands in
compiled Go.

## Prerequisites

Required:

- Go (for building let-go)
- [hyperfine](https://github.com/sharkdp/hyperfine) — benchmark runner
  (`brew install hyperfine`)
- python3 — results formatting

Optional (the script auto-detects what's available and skips the rest):

| Runtime | Install |
|---|---|
| [babashka](https://babashka.org/) | `brew install borkdude/brew/babashka` |
| [Clojure JVM](https://clojure.org/) | `brew install clojure` |

## Running

From the repo root:

```bash
bash benchmark/run.sh
```

This:

1. builds the plain let-go binary (`go build -ldflags="-s -w"`),
2. builds the AOT binary — `benchmark/aot/build.lg` lowers every fixture in
   `benchmark/aot/src/bench/*.lg` to native Go, wires the lowered packages
   into one `lg_bench`-tagged build, and **fails if any fixture silently
   falls back to bytecode**,
3. measures startup, peak RSS, and all seven workloads with hyperfine
   (3 warmups, 10 timed runs each), strictly sequentially,
4. regenerates `benchmark/results.md`.

A full run takes several minutes. Don't run anything CPU-heavy alongside it —
the benchmarks compete for cores and the numbers will be garbage.

### Iterating on a single benchmark

Pass benchmark names (basename of the `.clj` file) as arguments:

```bash
bash benchmark/run.sh fib
bash benchmark/run.sh fib tak
```

Filter mode uses more runs (10 warmups, 50 timed) for tighter σ, skips the
startup/memory sections, prints to stdout only, and does **not** rewrite
`results.md`.

## Layout

```
benchmark/
  *.clj                 # the workloads as plain Clojure scripts (bb / clj run these)
  aot/
    src/bench/*.lg      # the same workloads as fixtures: (ns bench.<name>) + 0-arg run
    drivers/*.lg        # per-benchmark driver: (require 'bench.<name>) (bench.<name>/run)
    build.lg            # lowers fixtures, builds the tagged letgo-aot binary, verifies
                        # native dispatch, cleans up everything it generated
  run.sh                # orchestrates builds + hyperfine + results.md generation
  results.md            # generated — do not edit by hand
```

## Adding a benchmark

1. Add the workload as a plain script: `benchmark/<name>.clj`.
2. Add the fixture: `benchmark/aot/src/bench/<name>.lg` with
   `(ns bench.<name>)` and a 0-arg `run` containing the same workload.
3. Add the driver: `benchmark/aot/drivers/<name>.lg` with
   `(require 'bench.<name>)` followed by `(bench.<name>/run)`.
4. `bash benchmark/run.sh <name>` to smoke-test it in filter mode.

Fixture caveat: a defn that lowers to a *fully typed* native signature (e.g.
a 0-arg fn inferred to return `int`) does not currently get an override
registration, so `build.lg`'s native-dispatch check will fail it. Keep the
hot code in its own defn and expose it through a boxed 0-arg `run` wrapper —
see `benchmark/aot/src/bench/loop-recur.lg` for the pattern.

## Methodology notes

- Times are mean ± σ wall-clock from hyperfine; peak memory is the median of
  3 runs of `/usr/bin/time -l`.
- Wall-clock includes binary startup and namespace load for both let-go legs.
- babashka and Clojure JVM run the plain `.clj` scripts (same workloads);
  Clojure JVM times include full JVM startup, which dominates the short
  benchmarks.
- Results in `results.md` were captured on an otherwise idle machine noted
  in its **System** line; absolute numbers vary by host — the ratios are the
  point.
