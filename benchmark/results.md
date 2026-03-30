## Benchmark Results

### Methodology

All benchmarks use [hyperfine](https://github.com/sharkdp/hyperfine) with 3 warmup runs
and 10 timed runs per benchmark. Times shown are mean ± σ wall-clock time. Each benchmark
file is valid Clojure that runs unmodified on all runtimes. Peak memory is measured
via `/usr/bin/time -l` (median of 3 runs). Clojure JVM times include full JVM startup
(~350-500ms) which dominates short benchmarks. Joker is skipped for benchmarks that
would exceed reasonable time limits or use unsupported features (transducers).

**System:** Darwin arm64, Apple M1 Pro

**Runtimes:**

| | let-go | babashka | joker | clojure JVM |
|---|---|---|---|---|
| **Version** | — | babashka v1.12.217 | joker v1.7.1 | Clojure CLI version 1.12.4.1618 |
| **Platform** | Go bytecode VM | GraalVM native | Go tree-walk interpreter | JVM (HotSpot) |
| **Binary/runtime size** | **9.6M** | 68M | 26M | 304M |

### Startup Time

| Runtime | Time |
|---|---|
| **let-go** | **7.0ms ± 0.8ms** (1.0x) |
| babashka | 24.6ms ± 10.1ms (3.5x) |
| joker | 11.4ms ± 0.7ms (1.6x) |
| clojure JVM | 0.333s ± 0.016s (47.7x) |

### Peak Memory Usage (RSS)

| Workload | let-go | babashka | joker | clojure JVM |
|---|---|---|---|---|
| startup (nil) | **12.8MB** (1.0x) | 26.7MB (2.1x) | 21.2MB (1.7x) | 92.9MB (7.3x) |
| fib(35) | **13.8MB** (1.0x) | 77.1MB (5.6x) | 33.3MB (2.4x) | 111.7MB (8.1x) |
| reduce 1M | **19.9MB** (1.0x) | 58.9MB (3.0x) | 33.0MB (1.7x) | 112.2MB (5.6x) |

### Performance

| Benchmark | let-go | babashka | joker | clojure JVM |
|---|---|---|---|---|
| fib | 2.013s ± 0.014s (1.0x) | 1.920s ± 0.026s (1.0x) | 20.161s ± 0.133s (10.0x) | **0.569s ± 0.033s** (0.3x) |
| loop-recur | **58.0ms ± 1.1ms** (1.0x) | 64.7ms ± 2.3ms (1.1x) | 0.717s ± 0.016s (12.4x) | 0.466s ± 0.010s (8.0x) |
| map-filter | **6.9ms ± 0.6ms** (1.0x) | 20.6ms ± 3.7ms (3.0x) | 12.8ms ± 0.7ms (1.9x) | 0.358s ± 0.006s (51.6x) |
| persistent-map | **19.2ms ± 0.9ms** (1.0x) | 21.7ms ± 0.8ms (1.1x) | 49.1ms ± 1.4ms (2.6x) | 0.490s ± 0.009s (25.5x) |
| reduce | 74.0ms ± 1.6ms (1.0x) | **36.5ms ± 3.4ms** (0.5x) | 2.522s ± 0.037s (34.1x) | 0.357s ± 0.006s (4.8x) |
| tak | 2.063s ± 0.008s (1.0x) | 1.958s ± 0.052s (0.9x) | — | **0.589s ± 0.015s** (0.3x) |
| transducers | **6.1ms ± 0.5ms** (1.0x) | 20.6ms ± 2.6ms (3.4x) | — | 0.361s ± 0.012s (59.5x) |

