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
| let-go | 15.3ms ± 0.4ms (1.0x) |
| babashka | 20.8ms ± 0.6ms (1.4x) |
| **joker** | **12.2ms ± 0.5ms** (0.8x) |
| go-joker | 12.7ms ± 0.5ms (0.8x) |
| gloat | 15.7ms ± 0.5ms (1.0x) |
| fennel | 49.2ms ± 2.5ms (3.2x) |
| clojure JVM | 0.381s ± 0.030s (25.0x) |

### Peak Memory Usage (RSS)

| Workload | let-go | babashka | joker | go-joker | gloat | fennel | clojure JVM |
|---|---|---|---|---|---|---|---|
| startup (nil) | 22.4MB (1.0x) | 27.0MB (1.2x) | 21.4MB (1.0x) | 23.5MB (1.0x) | 23.1MB (1.0x) | **3.2MB** (0.1x) | 97.6MB (4.4x) |
| fib(35) | 23.1MB (1.0x) | 77.4MB (3.4x) | 33.2MB (1.4x) | 24.0MB (1.0x) | 33.1MB (1.4x) | **12.8MB** (0.6x) | 118.2MB (5.1x) |
| reduce 1M | 28.2MB (1.0x) | 59.3MB (2.1x) | 33.3MB (1.2x) | **23.7MB** (0.8x) | 26.0MB (0.9x) | 886.8MB (31.4x) | 120.3MB (4.3x) |

### Performance

| Benchmark | let-go | babashka | joker | go-joker | gloat | fennel | clojure JVM |
|---|---|---|---|---|---|---|---|
| fib | 2.204s ± 0.068s (1.0x) | 1.931s ± 0.007s (0.9x) | 20.683s ± 0.201s (9.4x) | 1.444s ± 0.003s (0.7x) | 26.837s ± 0.415s (12.2x) | 1.984s ± 0.032s (0.9x) | **0.604s ± 0.010s** (0.3x) |
| loop-recur | 70.9ms ± 0.6ms (1.0x) | 65.7ms ± 0.6ms (0.9x) | 0.761s ± 0.069s (10.7x) | **40.5ms ± 33.9ms** (0.6x) | 1.031s ± 0.005s (14.5x) | 0.183s ± 0.002s (2.6x) | 0.511s ± 0.005s (7.2x) |
| map-filter | 17.8ms ± 7.1ms (1.0x) | 20.6ms ± 1.2ms (1.2x) | 13.7ms ± 0.5ms (0.8x) | **12.7ms ± 1.0ms** (0.7x) | 62.8ms ± 0.8ms (3.5x) | 1.085s ± 0.029s (61.1x) | 0.397s ± 0.003s (22.3x) |
| persistent-map | 26.5ms ± 0.5ms (1.0x) | 23.0ms ± 0.4ms (0.9x) | 50.2ms ± 0.8ms (1.9x) | **21.4ms ± 0.8ms** (0.8x) | 33.7ms ± 0.7ms (1.3x) | 3.716s ± 0.091s (140.0x) | 0.537s ± 0.006s (20.3x) |
| reduce | 74.9ms ± 1.4ms (1.0x) | 36.4ms ± 0.8ms (0.5x) | 2.593s ± 0.100s (34.6x) | **13.8ms ± 0.7ms** (0.2x) | 0.380s ± 0.011s (5.1x) | 8.275s ± 0.252s (110.4x) | 0.399s ± 0.004s (5.3x) |
| tak | 2.149s ± 0.006s (1.0x) | 1.981s ± 0.092s (0.9x) | — | 1.687s ± 0.015s (0.8x) | 22.237s ± 0.284s (10.3x) | 10.789s ± 0.139s (5.0x) | **0.619s ± 0.007s** (0.3x) |
| transducers | 52.8ms ± 0.8ms (1.0x) | 25.7ms ± 0.7ms (0.5x) | — | **16.9ms ± 0.6ms** (0.3x) | 0.198s ± 0.001s (3.8x) | 1.721s ± 0.053s (32.6x) | 0.429s ± 0.013s (8.1x) |

