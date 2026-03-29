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
| **Binary/runtime size** | **9.2M** | 68M | 26M | 304M |

### Startup Time

| Runtime | Time |
|---|---|
| let-go | 11.6ms ± 0.6ms (1.0x) |
| babashka | 20.9ms ± 3.6ms (1.8x) |
| **joker** | **11.0ms ± 1.2ms** (0.9x) |
| clojure JVM | 0.346s ± 0.016s (29.8x) |

### Peak Memory Usage (RSS)

| Workload | let-go | babashka | joker | clojure JVM |
|---|---|---|---|---|
| startup (nil) | **16.5MB** (1.0x) | 26.7MB (1.6x) | 21.2MB (1.3x) | 92.6MB (5.6x) |
| fib(35) | **17.0MB** (1.0x) | 77.0MB (4.5x) | 33.1MB (1.9x) | 112.3MB (6.6x) |
| reduce 1M | **20.9MB** (1.0x) | 58.9MB (2.8x) | 33.2MB (1.6x) | 112.2MB (5.4x) |

### Performance

| Benchmark | let-go | babashka | joker | clojure JVM |
|---|---|---|---|---|
| fib | 2.019s ± 0.020s (1.0x) | 1.944s ± 0.030s (1.0x) | 20.465s ± 0.083s (10.1x) | **0.582s ± 0.018s** (0.3x) |
| loop-recur | 69.2ms ± 10.4ms (1.0x) | **65.4ms ± 1.8ms** (0.9x) | 0.726s ± 0.014s (10.5x) | 0.488s ± 0.021s (7.0x) |
| map-filter | **12.8ms ± 0.3ms** (1.0x) | 21.5ms ± 1.5ms (1.7x) | 13.6ms ± 1.3ms (1.1x) | 0.381s ± 0.017s (29.8x) |
| persistent-map | 26.3ms ± 1.1ms (1.0x) | **23.2ms ± 1.5ms** (0.9x) | 49.4ms ± 1.5ms (1.9x) | 0.502s ± 0.005s (19.1x) |
| reduce | 93.5ms ± 15.2ms (1.0x) | **37.2ms ± 2.0ms** (0.4x) | 2.538s ± 0.025s (27.2x) | 0.381s ± 0.020s (4.1x) |
| tak | 2.075s ± 0.019s (1.0x) | 1.940s ± 0.019s (0.9x) | — | **0.600s ± 0.020s** (0.3x) |
| transducers | **12.9ms ± 0.7ms** (1.0x) | 20.4ms ± 1.1ms (1.6x) | — | 0.375s ± 0.016s (29.0x) |

