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
| let-go | 9.1ms ± 1.2ms (1.0x) |
| **let-go AOT** | **8.7ms ± 0.8ms** (1.0x) |
| babashka | 20.9ms ± 1.9ms (2.3x) |
| clojure JVM | 0.362s ± 0.018s (39.8x) |

### Peak Memory Usage (RSS)

| Workload | let-go | let-go AOT | babashka | clojure JVM |
|---|---|---|---|---|
| startup (nil) | **15.0MB** (1.0x) | 15.2MB (1.0x) | 27.0MB (1.8x) | 103.4MB (6.9x) |
| fib(35) | 16.0MB (1.0x) | **15.6MB** (1.0x) | 77.4MB (4.8x) | 121.8MB (7.6x) |
| reduce 1M | **21.0MB** (1.0x) | 22.0MB (1.0x) | 59.2MB (2.8x) | 121.9MB (5.8x) |

### Performance

| Benchmark | let-go | let-go AOT | babashka | clojure JVM |
|---|---|---|---|---|
| fib | 2.424s ± 0.028s (1.0x) | **0.113s ± 0.001s** (0.0x) | 1.935s ± 0.021s (0.8x) | 0.603s ± 0.018s (0.2x) |
| loop-recur | 70.9ms ± 0.9ms (1.0x) | **10.6ms ± 0.8ms** (0.1x) | 69.0ms ± 1.9ms (1.0x) | 0.501s ± 0.015s (7.1x) |
| map-filter | **10.7ms ± 1.7ms** (1.0x) | 12.4ms ± 1.7ms (1.2x) | 22.7ms ± 3.6ms (2.1x) | 0.398s ± 0.015s (37.1x) |
| persistent-map | **23.0ms ± 1.3ms** (1.0x) | 24.0ms ± 1.1ms (1.0x) | 29.3ms ± 6.1ms (1.3x) | 0.534s ± 0.017s (23.2x) |
| reduce | 41.0ms ± 1.5ms (1.0x) | 68.2ms ± 3.5ms (1.7x) | **37.6ms ± 2.3ms** (0.9x) | 0.376s ± 0.022s (9.2x) |
| tak | 2.432s ± 0.025s (1.0x) | **94.4ms ± 1.1ms** (0.0x) | 1.939s ± 0.022s (0.8x) | 0.649s ± 0.011s (0.3x) |
| transducers | 60.9ms ± 21.3ms (1.0x) | 49.4ms ± 1.6ms (0.8x) | **30.8ms ± 2.4ms** (0.5x) | 0.431s ± 0.014s (7.1x) |

