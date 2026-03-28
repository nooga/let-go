<!--suppress ALL -->
<p align="center">
<img src="meta/logo.png" alt="Squishy loafer" title="Squishy loafer of let-go" />
</p>


![Tests](https://github.com/nooga/let-go/actions/workflows/go.yml/badge.svg)

# let-go

Greetings loafers! *(lambda-gophers haha, get it?)*

**let-go** is a bytecode compiler and VM for a language closely resembling Clojure, written in Go.
It ships as a single ~9MB binary with ~12ms startup time, making it suitable for scripting, CLI tools,
and embedding in Go applications.

## Goals

- Implement as much of Clojure as practical — persistent data structures, lazy sequences, protocols, transducers, core.async
- Provide two-way Go interop — pass Go structs as records, call let-go from Go and vice versa
- Stay small and fast to start — no JVM, no GraalVM, just `go build`

## Non-goals

- Being a drop-in replacement for [clojure/clojure](https://github.com/clojure/clojure)
- Matching JVM HotSpot performance on compute-heavy workloads

## Feature overview

### Language

- Macros with syntax-quote, unquote, unquote-splicing
- Destructuring (sequential, associative, `:keys`, `:as`, `:or`)
- Multi-arity and variadic functions
- `loop`/`recur` with tail-call optimization
- `try`/`catch`/`finally`, `throw`, `ex-info`
- `letfn` for mutual recursion between local functions
- Dynamic variables with `binding`
- Lazy sequences (`lazy-seq`, `iterate`, `repeat`, `cycle`)
- Transducers (`map`, `filter`, `take`, `drop`, `partition-by`, etc. all return transducers with 1-arity)
- `transduce`, `into` with xform, `completing`, `sequence`, `cat`, `dedupe`
- Protocols and `extend-type` / `extend-protocol`
- Records with `defrecord`
- Multimethods with `defmulti` / `defmethod`
- Regular expressions (Go flavor)
- Metadata on collections and vars

### Data structures

- Persistent hash maps (HAMT), vectors, sets
- Transient collections for efficient batch building
- `delay` / `force`, `promise` / `deliver`
- `atom` with watches, `volatile!` for unsynchronized mutation
- `reduced` for early termination in `reduce`/`transduce`

### Concurrency (`async` namespace)

- `go` blocks and `go-loop` — goroutine-based lightweight concurrency
- Channels with optional buffering, `<!`, `>!`, `close!`
- `alts!` — select on multiple channel operations with timeout support
- `offer!` / `poll!` — non-blocking channel ops
- `mult` / `tap` / `untap` — broadcast
- `pub` / `sub` / `unsub` — topic-based routing
- `merge`, `pipe`, `split`, `async/map`, `async/take`
- `to-chan!`, `onto-chan!`, `async/into`, `async/reduce`
- `promise-chan`, `timeout`

### IO & Networking (`io` namespace)

- Protocol-based reader/writer coercion (`IReadable`, `IWritable`)
- `io/reader`, `io/writer` — polymorphic (strings as paths, handles, buffers, URLs)
- `io/line-seq` — lazy line-by-line reading
- `io/buffer` — mutable byte buffers
- `io/copy`, `io/slurp`, `io/spit`, `io/read-lines`, `io/write-lines`
- `io/url` — parsed URL records, readable via protocol (HTTP GET)
- Encoding: `io/encode` / `io/decode` (`:base64`, `:hex`, `:url`)
- Handle-based file IO: `open`, `close!`, `read-line`, `write!`, `read-bytes`
- `with-open` macro for auto-closing resources
- `*in*`, `*out*`, `*err*` — stdin/stdout/stderr

### HTTP (`http` namespace)

- Ring-style HTTP server (`http/serve2`)
- HTTP client: `http/get`, `http/post`, `http/request`
- Streaming responses with `:as :stream`
- URL records accepted in all client functions

### JSON (`json` namespace)

- `json/read-json`, `json/write-json`
- Proper float preservation, PersistentMap/Vector support, record serialization

### OS (`os` namespace)

- `os/sh` — run shell commands, capture stdout/stderr/exit code
- `os/stat`, `os/ls`, `os/cwd`, `os/getenv`, `os/setenv`, `os/exit`

### Go interop

- `RegisterStruct[T]` — map Go structs to let-go records with cached field converters
- `ToRecord[T]` / `ToStruct[T]` — zero-cost roundtrip for unmutated records
- `BoxValue` auto-converts registered structs to records
- Boxed Go values expose methods via `.method` interop syntax
- `.field` access on records

### Core library

Comprehensive `clojure.core` coverage including:
`comp`, `partial`, `juxt`, `complement`, `constantly`, `memoize`, `trampoline`,
`map`, `filter`, `reduce`, `mapcat`, `keep`, `take`, `drop`, `take-while`, `drop-while`,
`group-by`, `frequencies`, `partition`, `partition-by`, `interpose`, `interleave`,
`flatten`, `distinct`, `dedupe`, `sort-by`, `merge-with`, `select-keys`, `update-in`,
`get-in`, `assoc-in`, `tree-seq`, `cycle`, `doall`, `dorun`, `pmap`,
`future`, `promise`, `deliver`, `add-watch`, `remove-watch`, `subvec`,
`compare`, `not-any?`, `not-every?`, `doto`, `fn?`, `replace`, `nthrest`, `nthnext`,
`bit-and`, `bit-or`, `bit-xor`, `bit-not`, `bit-shift-left`, `bit-shift-right`,
`re-find`, `re-matches`, `re-seq`, `re-groups`, and many more.

Additional namespaces: `string`, `set`, `walk`, `edn`, `pprint`, `test`.

## Benchmarks

Benchmarks compare let-go against [Babashka](https://github.com/babashka/babashka) (GraalVM native),
[Joker](https://github.com/candid82/joker) (Go tree-walk interpreter), and Clojure on the JVM.
Each benchmark file is valid Clojure that runs unmodified on all runtimes.

Run `benchmark/run.sh` to reproduce (requires `hyperfine`, `bb`, `clj`, `joker`).

| | let-go | babashka | joker | clojure JVM |
|---|---|---|---|---|
| **Platform** | Go bytecode VM | GraalVM native | Go tree-walk interpreter | JVM (HotSpot) |
| **Binary size** | **9.1M** | 68M | 26M | 304M (JDK) |
| **Startup** | **11.5ms** | 21.2ms | 12.2ms | 338ms |
| **Idle memory** | **16.6MB** | 26.8MB | 21.2MB | 93MB |

**Performance** (Apple M1 Pro, mean of 10 runs):

| Benchmark | let-go | babashka | joker | clojure JVM |
|---|---|---|---|---|
| fib(35) | 3.05s | 1.93s | 20.4s | **0.58s** |
| loop-recur 1M | 140ms | **71ms** | 723ms | 473ms |
| map+filter+take | **13.3ms** | 22.9ms | 13.3ms | 372ms |
| persistent-map 10K | 30ms | **25ms** | 51ms | 564ms |
| reduce 1M | 103ms | **39ms** | 2.52s | 348ms |
| tak(30,22,12) | 3.92s | 1.92s | — | **0.56s** |
| transducers | **12.5ms** | 18.8ms | — | 335ms |

**Takeaways:**
- **Smallest footprint** — 7x smaller binary than Babashka, 33x smaller than the JDK
- **Fastest startup** — 12ms, ideal for CLI tools and scripting
- **Wins on short tasks** — map/filter, transducers: startup cost dominates, let-go has the least
- **Competitive on data structures** — persistent maps are within ~20% of Babashka
- **~2x slower on compute** vs Babashka (GraalVM AOT), ~5-7x vs JVM (JIT) — expected for a bytecode interpreter
- **6-24x faster than Joker** on all compute benchmarks — bytecode VM vs tree-walk interpreter

Full benchmark details: [benchmark/results.md](benchmark/results.md)

## Known limitations and divergence from Clojure

### Not implemented
- **Sorted collections** (`sorted-map`, `sorted-set`)
- **Refs / STM** — atoms + channels cover practical concurrency needs
- **Agents** — use `go` blocks and channels instead
- **Chunked sequences** — lazy seqs are unchunked (simpler, slightly different perf characteristics)
- **Namespaced keywords** (`:foo/bar`)
- **Reader tagged literals** (`#inst`, `#uuid`)
- **`deftype`** — use `defrecord` instead
- **`reify`** — protocols can only be extended to named types
- **Spec** — no `clojure.spec`
- **`alter-var-root`** — vars are mutable but no `alter-var-root`

### Behavioral differences
- **`concat*` (used internally by quasiquote) is eager** — the user-facing `concat` is lazy, matching Clojure
- **All channel operations block** — `<!` and `<!!` are identical (Go channels are always blocking), same for `>!`/`>!!`
- **`go` blocks are real goroutines** — no IOC (inversion of control) state machine like Clojure's core.async; this means they're cheaper but `go` blocks can call blocking ops directly
- **No BigDecimal** — numeric tower is `int64` + `float64` + `BigInt` (no arbitrary-precision decimals)
- **Regex is Go flavor** — `re2` syntax, not Java regex
- **`finally` always runs** but there's no `finally`-only (without `catch`) guaranteed execution on uncaught exceptions that cross native function boundaries
- **`letfn` uses atoms** internally for forward references — slight overhead vs Clojure's direct binding

## Examples

See:
- [Examples](https://github.com/nooga/let-go/tree/main/examples) for small programs
- [Tests](https://github.com/nooga/let-go/tree/main/test) for comprehensive `.lg` test files covering all features

## Try online

Check out [this bare-bones online REPL](https://nooga.github.io/let-go/). It runs a WASM build of let-go in your browser!

## Prerequisites and installation

Requires Go 1.22+.

```
go install github.com/nooga/let-go@latest
```

## Running from source

```bash
go run .                           # REPL
go run . -e '(+ 1 1)'             # eval expression
go run . myfile.lg                 # run file
go run . -r myfile.lg              # run file, then REPL
go build -ldflags="-s -w" -o letgo . # ~9MB stripped binary
```

## Embedding in Go

```go
import (
    "github.com/nooga/let-go/pkg/api"
    "github.com/nooga/let-go/pkg/vm"
)

c, _ := api.NewLetGo("myapp")

// Define Go values in let-go
c.Def("x", 42)
c.Def("greet", func(name string) string {
    return "Hello, " + name
})

// Run let-go code
v, _ := c.Run(`(greet "world")`)
fmt.Println(v) // "Hello, world"

// Struct <-> Record interop
type Point struct { X, Y int }
vm.RegisterStruct[Point]("myapp/Point")
c.Def("p", Point{3, 4})
v, _ = c.Run(`(:x p)`) // 3
```

## Testing

```bash
go test ./... -count=1 -timeout 30s
```

---
[Follow me on twitter](https://twitter.com/MGasperowicz)
