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
| **Binary/runtime size** | **9.1M** | 68M | 26M | 304M |

### Startup Time

| Runtime | Time |
|---|---|
| **let-go** | **11.5ms ± 0.6ms** |
| babashka | 21.2ms ± 2.7ms |
| joker | 12.2ms ± 1.2ms |
| clojure JVM | 0.338s ± 0.009s |

### Peak Memory Usage (RSS)

| Workload | let-go | babashka | joker | clojure JVM |
|---|---|---|---|---|
| startup (nil) | **16.6MB** | 26.8MB | 21.2MB | 93.0MB |
| fib(35) | **17.1MB** | 77.1MB | 32.9MB | 111.8MB |
| reduce 1M | 129.1MB | 58.9MB | **33.2MB** | 111.3MB |

### Performance

| Benchmark | let-go | babashka | joker | clojure JVM |
|---|---|---|---|---|
| fib | 3.051s ± 0.030s | 1.934s ± 0.017s | 20.356s ± 0.070s | **0.579s ± 0.017s** |
| loop-recur | 0.140s ± 0.002s | **71.0ms ± 8.2ms** | 0.723s ± 0.013s | 0.473s ± 0.010s |
| map-filter | 13.3ms ± 2.0ms | 22.9ms ± 2.3ms | **13.3ms ± 1.4ms** | 0.372s ± 0.008s |
| persistent-map | 30.0ms ± 7.7ms | **24.6ms ± 2.7ms** | 51.1ms ± 1.6ms | 0.564s ± 0.114s |
| reduce | 0.103s ± 0.004s | **39.4ms ± 9.3ms** | 2.523s ± 0.117s | 0.348s ± 0.004s |
| tak | 3.915s ± 0.030s | 1.920s ± 0.032s | — | **0.555s ± 0.016s** |
| transducers | **12.5ms ± 1.2ms** | 18.8ms ± 1.3ms | — | 0.335s ± 0.004s |

