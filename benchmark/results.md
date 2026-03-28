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
| **let-go** | **12.0ms ± 1.3ms** |
| babashka | 22.7ms ± 3.3ms |
| joker | 12.2ms ± 1.7ms |
| clojure JVM | 0.352s ± 0.017s |

### Peak Memory Usage (RSS)

| Workload | let-go | babashka | joker | clojure JVM |
|---|---|---|---|---|
| startup (nil) | **16.4MB** | 26.7MB | 21.3MB | 92.6MB |
| fib(35) | **16.8MB** | 77.1MB | 33.4MB | 111.9MB |
| reduce 1M | 114.1MB | 58.9MB | **33.1MB** | 112.5MB |

### Performance

| Benchmark | let-go | babashka | joker | clojure JVM |
|---|---|---|---|---|
| fib | 3.077s ± 0.039s | 1.947s ± 0.017s | 20.563s ± 0.089s | **0.584s ± 0.015s** |
| loop-recur | 0.149s ± 0.013s | **68.1ms ± 3.0ms** | 0.732s ± 0.013s | 0.482s ± 0.012s |
| map-filter | **13.8ms ± 1.2ms** | 22.9ms ± 2.3ms | 14.2ms ± 2.3ms | 0.381s ± 0.018s |
| persistent-map | 26.1ms ± 1.1ms | **24.8ms ± 2.6ms** | 49.7ms ± 1.5ms | 0.509s ± 0.008s |
| reduce | 0.111s ± 0.005s | **37.7ms ± 2.3ms** | 2.551s ± 0.024s | 0.377s ± 0.010s |
| tak | 3.995s ± 0.050s | 1.942s ± 0.023s | — | **0.603s ± 0.021s** |
| transducers | **13.1ms ± 2.2ms** | 23.1ms ± 3.6ms | — | 0.373s ± 0.009s |

