## Benchmark Results

### Methodology

All benchmarks use [hyperfine](https://github.com/sharkdp/hyperfine) with 3 warmup runs
and 10 timed runs per benchmark. Times shown are mean ± σ wall-clock time. Peak memory is
measured via `/usr/bin/time -l` (median of 3 runs).

**let-go** and **let-go AOT** run the identical driver + fixture (`benchmark/aot/`): the
same program, through the same binary entry point — the only difference is dispatch. The
plain build executes bytecode on the VM; the AOT build (produced by
`benchmark/aot/build.lg`) carries the same namespaces IR-lowered to native Go, registered
as overrides, so the same `(bench.<name>/run)` call dispatches to compiled Go. Wall-clock
includes binary startup and namespace load for both legs. babashka and Clojure JVM run the
plain `.clj` scripts, which contain the same workloads as the fixtures.

Clojure JVM times include full JVM startup (~350-500ms) which dominates short benchmarks.

**System:** Darwin arm64, Apple M1 Pro

**Runtimes:**

| | let-go | let-go AOT | babashka | clojure JVM |
|---|---|---|---|---|
| **Version** | — | — | babashka v1.12.217 | Clojure CLI version 1.12.4.1618 |
| **Platform** | Go bytecode VM | IR-lowered native Go | GraalVM native | JVM (HotSpot) |
| **Binary/runtime size** | **13M** | 18M | 68M | 304M |

### Startup Time

| Runtime | Time |
|---|---|
| let-go | 11.1ms ± 0.8ms (1.0x) |
| **let-go AOT** | **10.7ms ± 1.1ms** (1.0x) |
| babashka | 20.4ms ± 1.3ms (1.8x) |
| clojure JVM | 0.364s ± 0.010s (32.8x) |

### Peak Memory Usage (RSS)

| Workload | let-go | let-go AOT | babashka | clojure JVM |
|---|---|---|---|---|
| startup (nil) | **15.2MB** (1.0x) | **15.2MB** (1.0x) | 27.0MB (1.8x) | 97.7MB (6.4x) |
| fib(35) | 15.6MB (1.0x) | **15.5MB** (1.0x) | 77.4MB (5.0x) | 119.0MB (7.6x) |
| reduce 1M | 21.7MB (1.0x) | **21.4MB** (1.0x) | 59.2MB (2.7x) | 116.8MB (5.4x) |

### Performance

| Benchmark | let-go | let-go AOT | babashka | clojure JVM |
|---|---|---|---|---|
| fib | 2.416s ± 0.012s (1.0x) | **0.110s ± 0.001s** (0.0x) | 1.943s ± 0.049s (0.8x) | 0.567s ± 0.005s (0.2x) |
| loop-recur | 68.8ms ± 1.0ms (1.0x) | **10.7ms ± 0.9ms** (0.2x) | 64.4ms ± 1.3ms (0.9x) | 0.475s ± 0.004s (6.9x) |
| map-filter | 11.3ms ± 0.8ms (1.0x) | **10.4ms ± 0.7ms** (0.9x) | 19.2ms ± 0.9ms (1.7x) | 0.387s ± 0.011s (34.3x) |
| persistent-map | 22.2ms ± 0.7ms (1.0x) | 22.7ms ± 1.4ms (1.0x) | **21.7ms ± 1.2ms** (1.0x) | 0.517s ± 0.011s (23.3x) |
| reduce | 39.3ms ± 1.1ms (1.0x) | 66.6ms ± 1.3ms (1.7x) | **36.0ms ± 1.0ms** (0.9x) | 0.390s ± 0.015s (9.9x) |
| tak | 2.397s ± 0.022s (1.0x) | **93.7ms ± 0.6ms** (0.0x) | 1.908s ± 0.037s (0.8x) | 0.613s ± 0.046s (0.3x) |
| transducers | 46.5ms ± 0.9ms (1.0x) | 43.0ms ± 1.0ms (0.9x) | **20.3ms ± 0.5ms** (0.4x) | 0.373s ± 0.011s (8.0x) |

