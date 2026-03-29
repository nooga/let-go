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
| **Binary/runtime size** | **9.4M** | 68M | 26M | 304M |

### Startup Time

| Runtime | Time |
|---|---|
| **let-go** | **6.2ms ± 0.6ms** (1.0x) |
| babashka | 19.5ms ± 0.9ms (3.1x) |
| joker | 12.2ms ± 1.4ms (2.0x) |
| clojure JVM | 0.353s ± 0.020s (56.8x) |

### Peak Memory Usage (RSS)

| Workload | let-go | babashka | joker | clojure JVM |
|---|---|---|---|---|
| startup (nil) | **12.8MB** (1.0x) | 26.7MB (2.1x) | 21.2MB (1.7x) | 97.6MB (7.6x) |
| fib(35) | **13.6MB** (1.0x) | 77.1MB (5.7x) | 33.3MB (2.4x) | 118.0MB (8.7x) |
| reduce 1M | **19.5MB** (1.0x) | 59.0MB (3.0x) | 33.0MB (1.7x) | 116.7MB (6.0x) |

### Performance

| Benchmark | let-go | babashka | joker | clojure JVM |
|---|---|---|---|---|
| fib | 2.104s ± 0.067s (1.0x) | 1.950s ± 0.025s (0.9x) | 20.597s ± 0.074s (9.8x) | **0.588s ± 0.025s** (0.3x) |
| loop-recur | **60.7ms ± 4.7ms** (1.0x) | 65.8ms ± 1.4ms (1.1x) | 0.733s ± 0.016s (12.1x) | 0.481s ± 0.007s (7.9x) |
| map-filter | **6.0ms ± 0.3ms** (1.0x) | 19.2ms ± 0.9ms (3.2x) | 12.2ms ± 0.9ms (2.0x) | 0.376s ± 0.007s (62.5x) |
| persistent-map | **19.0ms ± 1.0ms** (1.0x) | 22.6ms ± 0.9ms (1.2x) | 48.5ms ± 1.3ms (2.6x) | 0.509s ± 0.022s (26.8x) |
| reduce | 76.4ms ± 1.8ms (1.0x) | **35.6ms ± 1.0ms** (0.5x) | 2.548s ± 0.022s (33.4x) | 0.371s ± 0.006s (4.9x) |
| tak | 2.094s ± 0.032s (1.0x) | 1.948s ± 0.030s (0.9x) | — | **0.597s ± 0.019s** (0.3x) |
| transducers | **7.1ms ± 0.6ms** (1.0x) | 20.4ms ± 1.1ms (2.9x) | — | 0.377s ± 0.015s (53.0x) |

