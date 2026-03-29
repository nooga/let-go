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
| **let-go** | **9.2ms ± 6.4ms** (1.0x) |
| babashka | 21.3ms ± 2.9ms (2.3x) |
| joker | 11.8ms ± 0.7ms (1.3x) |
| clojure JVM | 0.346s ± 0.015s (37.6x) |

### Peak Memory Usage (RSS)

| Workload | let-go | babashka | joker | clojure JVM |
|---|---|---|---|---|
| startup (nil) | **12.8MB** (1.0x) | 26.7MB (2.1x) | 21.2MB (1.7x) | 92.4MB (7.2x) |
| fib(35) | **13.7MB** (1.0x) | 77.1MB (5.6x) | 32.9MB (2.4x) | 112.4MB (8.2x) |
| reduce 1M | **20.2MB** (1.0x) | 58.9MB (2.9x) | 33.0MB (1.6x) | 112.2MB (5.6x) |

### Performance

| Benchmark | let-go | babashka | joker | clojure JVM |
|---|---|---|---|---|
| fib | 2.014s ± 0.019s (1.0x) | 1.938s ± 0.014s (1.0x) | 20.512s ± 0.107s (10.2x) | **0.579s ± 0.017s** (0.3x) |
| loop-recur | **58.2ms ± 0.5ms** (1.0x) | 65.8ms ± 1.1ms (1.1x) | 0.729s ± 0.013s (12.5x) | 0.486s ± 0.016s (8.4x) |
| map-filter | **7.1ms ± 0.4ms** (1.0x) | 20.5ms ± 0.8ms (2.9x) | 13.8ms ± 0.9ms (1.9x) | 0.376s ± 0.009s (53.0x) |
| persistent-map | **18.8ms ± 0.6ms** (1.0x) | 22.2ms ± 1.0ms (1.2x) | 53.1ms ± 12.3ms (2.8x) | 0.504s ± 0.006s (26.8x) |
| reduce | 75.6ms ± 1.1ms (1.0x) | **36.4ms ± 0.5ms** (0.5x) | 2.560s ± 0.020s (33.8x) | 0.376s ± 0.017s (5.0x) |
| tak | 2.083s ± 0.024s (1.0x) | 1.958s ± 0.029s (0.9x) | — | **0.590s ± 0.010s** (0.3x) |
| transducers | **7.0ms ± 0.4ms** (1.0x) | 21.6ms ± 3.3ms (3.1x) | — | 0.375s ± 0.005s (53.6x) |

