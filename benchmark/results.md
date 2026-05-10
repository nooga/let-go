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
| **Binary/runtime size** | **10M** | 68M | 26M | 32M | 26M | 324K | 304M |

### Startup Time

let-go, babashka, and fennel are measured with 10 warmup + 100 timed runs and a
5% trimmed mean (top/bottom 5% dropped) to reject scheduling outliers. Other
runtimes use the default 3 warmup + 10 runs.

| Runtime | Time |
|---|---|
| **let-go** | **6.7ms ± 0.3ms** (1.00x) |
| babashka | 18.3ms ± 0.8ms (2.74x) |
| joker | 11.8ms ± 1.1ms (1.77x) |
| go-joker | 13.1ms ± 0.9ms (1.96x) |
| gloat | 14.3ms ± 0.9ms (2.14x) |
| fennel | 40.2ms ± 4.0ms (6.02x) |
| clojure JVM | 0.307s ± 0.010s (45.96x) |

### Peak Memory Usage (RSS)

| Workload | let-go | babashka | joker | go-joker | gloat | fennel | clojure JVM |
|---|---|---|---|---|---|---|---|
| startup (nil) | 13.5MB (1.0x) | 26.7MB (2.0x) | 21.4MB (1.6x) | 23.3MB (1.7x) | 22.7MB (1.7x) | **3.1MB** (0.2x) | 92.3MB (6.8x) |
| fib(35) | 14.5MB (1.0x) | 77.1MB (5.3x) | 33.0MB (2.3x) | —MB | 32.8MB (2.3x) | **12.8MB** (0.9x) | 111.4MB (7.7x) |
| reduce 1M | **20.1MB** (1.0x) | 59.0MB (2.9x) | 33.6MB (1.7x) | 23.3MB (1.2x) | 25.7MB (1.3x) | 871.1MB (43.3x) | 118.0MB (5.9x) |

### Performance

| Benchmark | let-go | babashka | joker | go-joker | gloat | fennel | clojure JVM |
|---|---|---|---|---|---|---|---|
| fib | 1.994s ± 0.016s (1.0x) | 1.911s ± 0.018s (1.0x) | 19.445s ± 0.227s (9.8x) | — | 26.756s ± 0.595s (13.4x) | 1.922s ± 0.030s (1.0x) | **0.517s ± 0.009s** (0.3x) |
| loop-recur | 58.4ms ± 2.9ms (1.0x) | 61.6ms ± 1.3ms (1.1x) | 0.681s ± 0.005s (11.7x) | **11.9ms ± 0.5ms** (0.2x) | 1.011s ± 0.014s (17.3x) | 0.164s ± 0.003s (2.8x) | 0.426s ± 0.004s (7.3x) |
| map-filter | **6.2ms ± 0.2ms** (1.0x) | 20.0ms ± 0.6ms (3.2x) | 13.4ms ± 0.9ms (2.2x) | 10.9ms ± 0.5ms (1.8x) | 62.2ms ± 3.2ms (10.1x) | 1.010s ± 0.029s (163.6x) | 0.348s ± 0.059s (56.3x) |
| persistent-map | **18.0ms ± 0.6ms** (1.0x) | 18.3ms ± 0.4ms (1.0x) | 49.1ms ± 1.2ms (2.7x) | 19.7ms ± 1.3ms (1.1x) | 33.9ms ± 0.9ms (1.9x) | 3.506s ± 0.028s (194.8x) | 0.455s ± 0.013s (25.3x) |
| reduce | 72.3ms ± 1.8ms (1.0x) | 34.6ms ± 1.5ms (0.5x) | 2.423s ± 0.032s (33.5x) | **10.9ms ± 0.4ms** (0.2x) | 0.368s ± 0.012s (5.1x) | 7.593s ± 0.155s (105.0x) | 0.336s ± 0.009s (4.7x) |
| tak | 2.054s ± 0.031s (1.0x) | 1.932s ± 0.039s (0.9x) | — | — | 21.606s ± 0.142s (10.5x) | 10.459s ± 0.103s (5.1x) | **0.531s ± 0.011s** (0.3x) |
| transducers | **5.8ms ± 0.2ms** (1.0x) | 16.1ms ± 0.6ms (2.8x) | — | — | 14.3ms ± 0.3ms (2.4x) | 0.978s ± 0.021s (167.5x) | 0.326s ± 0.005s (55.8x) |

