## Benchmark Results

### Methodology

All benchmarks use [hyperfine](https://github.com/sharkdp/hyperfine) with 3 warmup runs
and 10 timed runs per benchmark. Times shown are mean ± σ wall-clock time. Peak memory is
measured via `/usr/bin/time -l` (median of 3 runs).

Benchmark files are valid Clojure that runs unmodified on let-go, babashka, joker, go-joker,
glojure, and Clojure JVM. Fennel uses equivalent implementations via
[fennel-cljlib](https://gitlab.com/andreyorst/fennel-cljlib) (lazy seqs, transducers,
persistent data structures). Gloat benchmarks are pre-compiled to native binaries via
[gloat](https://github.com/gloathub/gloat) AOT (Clojure→Go); compilation time is not
measured, only binary execution (analogous to how let-go is pre-built with `go build`).

Clojure JVM times include full JVM startup (~350-500ms) which dominates short benchmarks.
Joker is skipped for benchmarks that would exceed reasonable time limits or use unsupported
features (transducers). Binary sizes for gloat are averaged across all benchmark binaries.

**System:** Darwin arm64, Apple M1 Pro

**Runtimes:**

| | let-go | babashka | joker | go-joker | gloat | fennel | clojure JVM |
|---|---|---|---|---|---|---|---|
| **Version** | — | babashka v1.12.217 | joker v1.7.1 | go-joker v42.8.2 | gloat version 0.1.36 | Fennel 1.6.1 on PUC Lua 5.5 | Clojure CLI version 1.12.4.1618 |
| **Platform** | Go bytecode VM | GraalVM native | Go tree-walk interpreter | Go IR + WASM/wazero JIT | Go AOT (Clojure→Go) | Lua VM + cljlib | JVM (HotSpot) |
| **Binary/runtime size** | **12M** | 68M | 26M | 32M | 26M | 324K | 304M |

### Startup Time

| Runtime | Time |
|---|---|
| **let-go** | **8.2ms ± 0.4ms** (1.0x) |
| babashka | 17.7ms ± 0.6ms (2.2x) |
| joker | 11.5ms ± 0.8ms (1.4x) |
| go-joker | 12.5ms ± 0.6ms (1.5x) |
| gloat | 14.7ms ± 0.6ms (1.8x) |
| fennel | 42.9ms ± 0.9ms (5.2x) |
| clojure JVM | 0.360s ± 0.023s (43.9x) |

### Peak Memory Usage (RSS)

| Workload | let-go | babashka | joker | go-joker | gloat | fennel | clojure JVM |
|---|---|---|---|---|---|---|---|
| startup (nil) | 14.7MB (1.0x) | 27.0MB (1.8x) | 21.6MB (1.5x) | 23.7MB (1.6x) | 22.9MB (1.6x) | **3.1MB** (0.2x) | 98.0MB (6.7x) |
| fib(35) | 15.5MB (1.0x) | 77.4MB (5.0x) | 33.1MB (2.1x) | 24.1MB (1.6x) | 33.2MB (2.1x) | **12.5MB** (0.8x) | 118.2MB (7.6x) |
| reduce 1M | **21.2MB** (1.0x) | 59.2MB (2.8x) | 33.3MB (1.6x) | 23.5MB (1.1x) | 26.6MB (1.3x) | 887.3MB (41.9x) | 121.5MB (5.7x) |

### Performance

| Benchmark | let-go | babashka | joker | go-joker | gloat | fennel | clojure JVM |
|---|---|---|---|---|---|---|---|
| fib | 2.128s ± 0.030s (1.0x) | 1.906s ± 0.022s (0.9x) | 20.260s ± 0.191s (9.5x) | 1.452s ± 0.051s (0.7x) | 27.041s ± 0.636s (12.7x) | 1.940s ± 0.045s (0.9x) | **0.619s ± 0.134s** (0.3x) |
| loop-recur | 65.3ms ± 1.5ms (1.0x) | 65.8ms ± 2.9ms (1.0x) | 0.688s ± 0.005s (10.5x) | **12.6ms ± 1.1ms** (0.2x) | 1.010s ± 0.022s (15.5x) | 0.169s ± 0.003s (2.6x) | 0.454s ± 0.005s (6.9x) |
| map-filter | **7.2ms ± 0.3ms** (1.0x) | 17.0ms ± 0.6ms (2.4x) | 11.1ms ± 0.6ms (1.6x) | 13.2ms ± 1.0ms (1.8x) | 63.1ms ± 1.2ms (8.8x) | 1.011s ± 0.020s (141.3x) | 0.355s ± 0.005s (49.6x) |
| persistent-map | 20.2ms ± 1.2ms (1.0x) | **18.6ms ± 0.6ms** (0.9x) | 49.5ms ± 1.2ms (2.5x) | 19.4ms ± 0.6ms (1.0x) | 33.2ms ± 0.7ms (1.6x) | 3.635s ± 0.070s (180.3x) | 0.501s ± 0.019s (24.9x) |
| reduce | 66.9ms ± 0.7ms (1.0x) | 36.1ms ± 1.2ms (0.5x) | 2.472s ± 0.050s (37.0x) | **12.6ms ± 0.8ms** (0.2x) | 0.364s ± 0.003s (5.4x) | 8.070s ± 0.311s (120.7x) | 0.370s ± 0.014s (5.5x) |
| tak | 2.140s ± 0.043s (1.0x) | 1.929s ± 0.042s (0.9x) | — | 1.683s ± 0.035s (0.8x) | 21.978s ± 0.368s (10.3x) | 10.827s ± 0.177s (5.1x) | **0.639s ± 0.019s** (0.3x) |
| transducers | 46.9ms ± 0.7ms (1.0x) | 27.6ms ± 1.4ms (0.6x) | — | **16.7ms ± 0.6ms** (0.4x) | 0.200s ± 0.002s (4.3x) | 1.707s ± 0.052s (36.4x) | 0.388s ± 0.005s (8.3x) |

