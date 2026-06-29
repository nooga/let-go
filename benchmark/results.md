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
| **Binary/runtime size** | **13M** | 68M | 26M | 32M | 26M | 324K | 304M |

### Startup Time

| Runtime | Time |
|---|---|
| let-go | 15.0ms ± 0.5ms (1.0x) |
| babashka | 20.4ms ± 0.5ms (1.4x) |
| **joker** | **12.2ms ± 0.7ms** (0.8x) |
| go-joker | 12.7ms ± 1.1ms (0.8x) |
| gloat | 15.8ms ± 0.6ms (1.1x) |
| fennel | 54.3ms ± 7.8ms (3.6x) |
| clojure JVM | 0.417s ± 0.030s (27.8x) |

### Peak Memory Usage (RSS)

| Workload | let-go | babashka | joker | go-joker | gloat | fennel | clojure JVM |
|---|---|---|---|---|---|---|---|
| startup (nil) | 21.5MB (1.0x) | 27.0MB (1.3x) | 21.4MB (1.0x) | 23.4MB (1.1x) | 23.0MB (1.1x) | **3.2MB** (0.1x) | 97.2MB (4.5x) |
| fib(35) | 23.1MB (1.0x) | 77.4MB (3.4x) | 33.6MB (1.5x) | 24.1MB (1.0x) | 33.1MB (1.4x) | **12.6MB** (0.5x) | 117.1MB (5.1x) |
| reduce 1M | 26.7MB (1.0x) | 59.2MB (2.2x) | 33.6MB (1.3x) | **23.7MB** (0.9x) | 26.1MB (1.0x) | 1167.7MB (43.7x) | 121.6MB (4.6x) |

### Performance

| Benchmark | let-go | babashka | joker | go-joker | gloat | fennel | clojure JVM |
|---|---|---|---|---|---|---|---|
| fib | 2.191s ± 0.030s (1.0x) | 1.944s ± 0.009s (0.9x) | 20.782s ± 0.170s (9.5x) | 1.471s ± 0.019s (0.7x) | 27.734s ± 0.677s (12.7x) | 2.209s ± 0.184s (1.0x) | **0.694s ± 0.175s** (0.3x) |
| loop-recur | 71.7ms ± 0.8ms (1.0x) | 67.5ms ± 2.2ms (0.9x) | 0.769s ± 0.049s (10.7x) | **14.8ms ± 0.9ms** (0.2x) | 1.050s ± 0.012s (14.7x) | 0.191s ± 0.006s (2.7x) | 0.548s ± 0.051s (7.6x) |
| map-filter | 14.9ms ± 0.4ms (1.0x) | 20.2ms ± 0.4ms (1.4x) | 13.5ms ± 0.8ms (0.9x) | **12.8ms ± 0.5ms** (0.9x) | 65.8ms ± 5.2ms (4.4x) | 1.104s ± 0.016s (74.1x) | 0.408s ± 0.008s (27.3x) |
| persistent-map | 26.3ms ± 0.9ms (1.0x) | 24.2ms ± 2.7ms (0.9x) | 50.5ms ± 0.9ms (1.9x) | **21.5ms ± 1.0ms** (0.8x) | 34.8ms ± 1.1ms (1.3x) | 3.745s ± 0.097s (142.4x) | 0.536s ± 0.003s (20.4x) |
| reduce | 75.1ms ± 1.7ms (1.0x) | 36.9ms ± 0.7ms (0.5x) | 2.571s ± 0.028s (34.2x) | **13.3ms ± 0.6ms** (0.2x) | 0.396s ± 0.030s (5.3x) | 8.300s ± 0.205s (110.5x) | 0.401s ± 0.003s (5.3x) |
| tak | 2.215s ± 0.048s (1.0x) | 1.947s ± 0.004s (0.9x) | — | 1.711s ± 0.024s (0.8x) | 22.337s ± 0.322s (10.1x) | 10.826s ± 0.134s (4.9x) | **0.636s ± 0.020s** (0.3x) |
| transducers | 53.1ms ± 1.1ms (1.0x) | 27.6ms ± 1.3ms (0.5x) | — | **16.9ms ± 0.6ms** (0.3x) | 0.200s ± 0.001s (3.8x) | 1.727s ± 0.039s (32.5x) | 0.416s ± 0.015s (7.8x) |

