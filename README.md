<!--suppress ALL -->
<p align="center">
<img src="meta/logo.png" alt="Squishy loafer" title="Squishy loafer of let-go" />
</p>

![Tests](https://github.com/nooga/let-go/actions/workflows/go.yml/badge.svg)

# let-go

Greetings loafers! _(λ-gophers haha, get it?)_

This is a bytecode compiler and VM for a language closely resembling Clojure, a Clojure dialect, if you will.
Ships as a single ~9MB binary with ~6ms startup time.

### Why let-go?

- **Standalone executables** — compile your program into a single binary with `lg -b myapp main.lg`. No runtime needed, just distribute and run.
- **WASM web apps** — compile your program to a self-contained HTML page with `lg -w outdir main.lg`. Full terminal emulation via xterm.js, runs in any browser. Deploy to GitHub Pages or open locally.
- **Fast startup** — 6ms cold start. Pre-compiled bytecode (LGB format) makes boot near-instant even with a large standard library.
- **Small footprint** — 9MB binary, 13MB idle memory. 7x smaller than Babashka, 33x smaller than JDK.
- **Batteries included** — core.async channels, HTTP server/client, JSON, Transit, IO, Babashka pods, nREPL server.
- **Go interop** — embed let-go in Go apps, map Go structs to records, call Go functions from let-go and vice versa.
- **Broad Clojure compatibility** — macros, destructuring, protocols, records, multimethods, transducers, lazy seqs, persistent data structures, BigInts.

Here are some nebulous goals in no particular order:

- [x] Quality entertainment,
- [ ] Making it legal to write Clojure at your Go dayjob,
- [x] Implement as much of Clojure as possible — including persistent data types, true concurrency, transducers, core.async, and BigInts,
- [x] Provide comfy two-way interop for arbitrary functions and types,
- [x] AOT compilation — compile let-go programs to bytecode or standalone binaries,
- [x] Boot the entire runtime in a single `requestAnimationFrame` and still have 10ms to spare at 60fps,
- [x] Compile let-go programs to self-contained WASM web apps with terminal emulation,
- [ ] Stretch goal: let-go bytecode -> Go translation.

Here are the non goals:

- Being a drop-in replacement for [clojure/clojure](https://github.com/clojure/clojure) at any point,
- Being a linter/formatter/tooling for Clojure in general.

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

- Ring-style HTTP server (`http/serve`)
- HTTP client: `http/get`, `http/post`, `http/request`
- Streaming responses with `:as :stream`
- URL records accepted in all client functions

### JSON (`json` namespace)

- `json/read-json`, `json/write-json`
- Proper float preservation, PersistentMap/Vector support, record serialization

### Transit (`transit` namespace)

- `transit/read`, `transit/write` - transit+json codec
- Full rolling cache support for compact encoding
- Keywords, symbols, maps, vectors, sets, lists, big integers

### OS (`os` namespace)

- `os/sh` - run shell commands, capture stdout/stderr/exit code
- `os/stat`, `os/ls`, `os/cwd`, `os/getenv`, `os/setenv`, `os/exit`

### Babashka pods

let-go supports [Babashka pods](https://github.com/babashka/pods) - standalone programs that expose namespaces over a binary protocol. This gives let-go access to the entire pod ecosystem: databases, AWS, Docker, file watching, and more.

```clojure
;; Load a pod (uses babashka's shared cache)
(pods/load-pod 'org.babashka/go-sqlite3 "0.3.13")

;; Use it like any other namespace
(pod.babashka.go-sqlite3/execute! "app.db"
  ["create table users (id integer primary key, name text)"])
(pod.babashka.go-sqlite3/execute! "app.db"
  ["insert into users values (1, ?)" "Alice"])
(pod.babashka.go-sqlite3/query "app.db"
  ["select * from users"])
;; => [{:id 1 :name "Alice"}]
```

- `pods/load-pod` - load by name (PATH) or from babashka cache (symbol + version)
- Supports JSON, EDN, and transit+json payload formats
- Client-side code evaluation (pod-defined macros and wrappers)
- Async streaming via `pods/invoke` with `:handlers` for callbacks
- Shares `~/.babashka/pods/` cache - install pods with `bb`, use them from `lg`

See the [pod registry](https://github.com/babashka/pod-registry) for available pods. Install pods with babashka:

```bash
bb -e '(pods/load-pod (quote org.babashka/go-sqlite3) "0.3.13")'
```

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

Additional namespaces: `string`, `set`, `walk`, `edn`, `pprint`, `test`, `transit`, `pods`.

## Benchmarks

Benchmarks compare let-go against [Babashka](https://github.com/babashka/babashka) (GraalVM native),
[Joker](https://github.com/candid82/joker) (Go tree-walk interpreter), and Clojure on the JVM.
Each benchmark is valid Clojure that runs unmodified on all runtimes.
Run `benchmark/run.sh` to reproduce (requires `hyperfine`, `bb`, `clj`, `joker`).

|                 | let-go         | babashka       | joker                    | clojure JVM   |
| --------------- | -------------- | -------------- | ------------------------ | ------------- |
| **Platform**    | Go bytecode VM | GraalVM native | Go tree-walk interpreter | JVM (HotSpot) |
| **Binary size** | **9.4M**       | 68M            | 26M                      | 304M (JDK)    |
| **Startup**     | **6ms**        | 20ms           | 12ms                     | 353ms         |
| **Idle memory** | **13MB**       | 27MB           | 21MB                     | 98MB          |

**Performance highlights** (Apple M1 Pro):

- **Smallest footprint** — 7x smaller than Babashka, 33x smaller than the JDK
- **Fastest startup** — 6ms with pre-compiled bytecode (fits in a `requestAnimationFrame`), 3x faster than Babashka, 2x faster than Joker, 57x faster than JVM
- **Wins on short-lived tasks** — map/filter and transducer pipelines: **6-7ms** vs bb's 20ms (3x faster)
- **Competitive on compute** — fib(35) within 8% of Babashka (2.1s vs 1.9s), loop-recur 8% faster
- **Lowest memory** — 14MB for fib(35) vs bb's 77MB (5.7x less), 20MB for reduce 1M vs bb's 59MB (3x less)
- **10x faster than Joker** on all compute benchmarks — bytecode VM vs tree-walk interpreter

Full results with methodology: [benchmark/results.md](benchmark/results.md)

## Known limitations and divergence from Clojure

### Not implemented

- **Sorted collections** (`sorted-map`, `sorted-set`)
- **Refs / STM** — atoms + channels cover practical concurrency needs
- **Agents** — use `go` blocks and channels instead
- **Chunked sequences** — lazy seqs are unchunked (simpler, slightly different perf characteristics)
- **Reader tagged literals** (`#inst`, `#uuid`)
- **`deftype`** — use `defrecord` instead
- **`reify`** — protocols can only be extended to named types
- **Spec** — no `clojure.spec`
- **`alter-var-root`** — vars are mutable but no `alter-var-root`

### Known behavioral differences

- **`concat*` (used internally by quasiquote) is eager** — the user-facing `concat` is lazy, matching Clojure
- **All channel operations block** — `<!` and `<!!` are identical (Go channels are always blocking), same for `>!`/`>!!`
- **`go` blocks are real goroutines** — no IOC (inversion of control) state machine like Clojure's core.async; this means they're cheaper but `go` blocks can call blocking ops directly
- **No BigDecimal** — numeric tower is `int64` + `float64` + `BigInt` (no arbitrary-precision decimals)
- **Regex is Go flavor** — `re2` syntax, not Java regex
- **`letfn` uses atoms** internally for forward references — slight overhead vs Clojure's direct binding

## Examples

See:

- [Examples](https://github.com/nooga/let-go/tree/main/examples) for small programs
- [Tests](https://github.com/nooga/let-go/tree/main/test) for comprehensive `.lg` test files covering all features

## Try online

Check out [this bare-bones online REPL](https://nooga.github.io/let-go/). It runs a WASM build of let-go in your browser!

## Installation

### Homebrew (macOS / Linux)

```bash
brew tap nooga/let-go https://github.com/nooga/let-go
brew install let-go
```

### Download binary

Grab a prebuilt binary from [Releases](https://github.com/nooga/let-go/releases) — available for Linux, macOS, and Windows on amd64/arm64.

### From source

Requires Go 1.22+.

```bash
go install github.com/nooga/let-go@latest
```

### Usage

```bash
lg                                 # REPL
lg -e '(+ 1 1)'                   # eval expression
lg myfile.lg                       # run file
lg -r myfile.lg                    # run file, then REPL
lg -w outdir myfile.lg             # compile to WASM web app
```

### Compilation and distribution

let-go can compile programs to bytecode (`.lgb` files) and package them as standalone executables.

**Compile to bytecode** — skips the reader/parser/compiler at load time:

```bash
lg -c app.lgb app.lg               # compile to bytecode
lg app.lgb                          # run bytecode directly
```

**Create a standalone binary** — bundles the compiled bytecode into a self-contained executable:

```bash
lg -b myapp app.lg                  # compile + bundle into executable
./myapp                             # runs anywhere, no lg needed
```

The standalone binary is a copy of `lg` with your program's bytecode appended. It needs no external files or runtime — just copy it to another machine and run it.

**Build a WASM web app** — compiles your program into a single HTML page that runs in the browser:

```bash
lg -w site app.lg                   # compile to web app
open site/index.html                # open in browser
```

The output directory contains:
- `index.html` — self-contained (~6MB, inlined WASM + wasm_exec.js, gzip-compressed)
- `coi-serviceworker.js` — enables cross-origin isolation for interactive apps (needed on GitHub Pages)

Programs using the `term` namespace get full terminal emulation via xterm.js — ANSI colors, cursor positioning, raw keyboard input all work. The Go WASM runtime runs in a Web Worker with SharedArrayBuffer for blocking `term/read-key`.

For GitHub Pages deployment, just point Pages at the output directory. The service worker handles the required COOP/COEP headers automatically.

**Detecting AOT compilation** — the `*compiling-aot*` var is `true` during `-c`, `-b`, and `-w` compilation, `false` at runtime. Use it to prevent side effects (like starting a server or game loop) from running at compile time:

```clojure
(defn -main []
  (start-server))

(when-not *compiling-aot*
  (-main))
```

**Detecting WASM at runtime** — the `*in-wasm*` var is `true` when running inside a WASM web app, `false` in native mode. Use it to disable file I/O, adjust animation timing, or enable browser-specific behavior:

```clojure
(when-not *in-wasm*
  (spit "debug.log" "only in native mode"))
```

### Building from source

```bash
go run .                           # run from source
go build -ldflags="-s -w" -o lg .  # ~9MB stripped binary
```

## nREPL

let-go includes an nREPL server compatible with CIDER (Emacs), Calva (VS Code), and Conjure (Neovim).

```bash
lg -n                              # start nREPL on default port (2137)
lg -n -p 7888                      # start nREPL on port 7888
```

The server writes `.nrepl-port` in the current directory so editors auto-discover it.

**Supported ops:** `clone`, `close`, `eval`, `load-file`, `describe`, `completions`, `complete`, `info`, `lookup`, `ls-sessions`, `interrupt`

**Emacs (CIDER):** `M-x cider-connect-clj`, host `localhost`, port from `.nrepl-port`

**VS Code (Calva):** Open a let-go project — the included `.vscode/settings.json` registers a custom connect sequence. Use "Calva: Start a Project REPL and Connect (Jack-In)" and pick "let-go", or "Calva: Connect to a Running REPL Server" if the nREPL is already running.

**Neovim (Conjure):** Should auto-connect when `.nrepl-port` exists

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

[🤓 Follow me on twitter](https://twitter.com/MGasperowicz)
[🐬 Check out monk.io](https://monk.io)
