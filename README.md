<!--suppress ALL -->
<p align="center">
<img src="meta/logo.png" alt="Squishy loafer" title="Squishy loafer of let-go" />
</p>

![Tests](https://github.com/nooga/let-go/actions/workflows/go.yml/badge.svg)

# let-go

Greetings loafers! _(λ-gophers haha, get it?)_

This is a bytecode compiler and VM for a language closely resembling Clojure, a Clojure dialect, if you will.
The smallest and fastest-starting Clojure-family language in Go — a single ~10MB binary with ~7ms cold start.

### Why let-go?

- **Standalone executables** — compile your program into a single binary with `lg -b myapp main.lg`. No runtime needed, just distribute and run.
- **WASM web apps** — compile your program to a self-contained HTML page with `lg -w outdir main.lg`. Full terminal emulation via xterm.js, runs in any browser. Deploy to GitHub Pages or open locally.
- **Fast startup** — 6ms cold start. Pre-compiled bytecode (LGB format) makes boot near-instant even with a large standard library.
- **Small footprint** — 10MB binary, 14MB idle memory. 7x smaller than Babashka, 30x smaller than JDK.
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
- [ ] nREPL server in WASM — connect Emacs/Calva to a let-go VM running in the browser via WebSocket,
- [ ] Stretch goal: let-go bytecode -> Go translation.

Here are the non goals:

- Being a drop-in replacement for [clojure/clojure](https://github.com/clojure/clojure) at any point,
- Being a linter/formatter/tooling for Clojure in general.

## Feature overview

let-go aims to feel like day-to-day Clojure, not to be a drop-in replacement. Most idiomatic code reads,
runs, and behaves the same — but a non-trivial Clojure project will likely need some adjustments before
it runs unmodified. See [Known limitations](#known-limitations-and-divergence-from-clojure) below.

### Clojure compatibility

let-go is tested against [jank-lang/clojure-test-suite](https://github.com/jank-lang/clojure-test-suite),
a cross-dialect compliance suite:

**4696 / 4921 assertions pass (95.4%)** across 217 test files. The remaining gaps are mostly
edge cases (overflow detection on `+`/`-`/`*`/`inc`/`dec`, BigInt promotion at the Long boundary, BigDecimal
behavior) and a handful of stub namespaces — see below. Workflow guide: [docs/clojure-test-suite.md](docs/clojure-test-suite.md).

### Standard namespaces

| Namespace            | Status                                                                                                          |
| -------------------- | --------------------------------------------------------------------------------------------------------------- |
| `clojure.core`       | Macros, destructuring, lazy seqs, transducers, protocols, records, multimethods, atoms, regex, metadata, BigInt |
| `clojure.string`     | Full                                                                                                            |
| `clojure.set`        | Full                                                                                                            |
| `clojure.walk`       | `prewalk`, `postwalk`, `keywordize-keys`, `stringify-keys`, `walk`                                              |
| `clojure.edn`        | `read`, `read-string`                                                                                           |
| `clojure.pprint`     | `pprint`, `cl-format`                                                                                           |
| `clojure.test`       | `deftest`, `is`, `testing`, `are`, fixtures                                                                     |
| `clojure.core.async` | Channels, `go`/`go-loop`, `alts!`, `mult`/`pub`, `pipe`/`merge`/`split` (real goroutines, not IOC)              |
| `io`                 | Polymorphic readers/writers, `slurp`/`spit`, lazy line-seq, encoding, URLs, `with-open`                         |
| `http`               | Ring-style server + client, streaming responses                                                                 |
| `json`               | `read-json`, `write-json` — float-preserving, record-aware                                                      |
| `transit`            | transit+json codec with rolling cache                                                                           |
| `os`                 | `sh`, `stat`, `ls`, `cwd`, `getenv`/`setenv`, `exit`, `os-name`, `arch`, `user-name`, `hostname`, separators    |
| `System`             | JVM-shaped: `getProperty`, `getProperties`, `getenv`, `exit`, `lineSeparator`, `currentTimeMillis`, `nanoTime`. Exposes `let-go.version`, `let-go.commit`, `user.home`, `user.dir`, `os.name`, `os.arch`, etc. |
| `syscall`            | Direct Linux syscalls (mount, unshare, mknod, prctl, capset, seccomp, AppArmor) for systems programming         |
| `pods`               | Babashka pods over JSON / EDN / transit                                                                         |

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

## Benchmarks

Benchmarks compare let-go against [Babashka](https://github.com/babashka/babashka) (GraalVM native),
[Joker](https://github.com/candid82/joker) (Go tree-walk interpreter), and Clojure on the JVM.
Each benchmark is valid Clojure that runs unmodified on all runtimes.
Run `benchmark/run.sh` to reproduce (requires `hyperfine`, `bb`, `clj`, `joker`).

|                 | let-go         | babashka       | joker                    | clojure JVM   |
| --------------- | -------------- | -------------- | ------------------------ | ------------- |
| **Platform**    | Go bytecode VM | GraalVM native | Go tree-walk interpreter | JVM (HotSpot) |
| **Binary size** | **10M**        | 68M            | 26M                      | 304M (JDK)    |
| **Startup**     | **7ms**        | 20ms           | 12ms                     | 331ms         |
| **Idle memory** | **14MB**       | 27MB           | 21MB                     | 92MB          |

**Performance highlights** (Apple M1 Pro):

- **Smallest footprint** — 7x smaller than Babashka, 30x smaller than the JDK
- **Fastest startup** — 7ms with pre-compiled bytecode (fits in a `requestAnimationFrame`), 3x faster than Babashka, 2x faster than Joker, 48x faster than JVM
- **Wins on short-lived tasks** — map/filter and transducer pipelines: **8ms** vs bb's 19ms (2.4x faster)
- **Competitive on compute** — fib(35) within 4% of Babashka (1.98s vs 1.90s), loop-recur 1.8x faster
- **Lowest memory** — 14MB for fib(35) vs bb's 77MB (5.4x less), 20MB for reduce 1M vs bb's 59MB (3x less)
- **10x+ faster than Joker** on most compute benchmarks — bytecode VM vs tree-walk interpreter

Full results with methodology: [benchmark/results.md](benchmark/results.md)

## Known limitations and divergence from Clojure

### Not implemented

- **Refs / STM** — atoms + channels cover practical concurrency needs
- **Agents** — use `go` blocks and channels instead
- **Hierarchies** (`derive`, `underive`, `ancestors`, `descendants`, `parents`) — stub only; multimethod dispatch works, but `isa?` chains do not
- **`with-precision`** — `BigDecimal` itself works (`M` literals, `bigdec`, `decimal?`, exact arithmetic), but `with-precision` is a no-op so explicit rounding control is missing
- **Chunked sequences** — lazy seqs are unchunked (simpler, slightly different perf characteristics)
- **Reader tagged literals** (`#inst`, `#uuid`)
- **`deftype`** — use `defrecord` instead
- **`reify`** — protocols can only be extended to named types
- **Spec** — no `clojure.spec`
- **`alter-var-root`** — vars are mutable but no `alter-var-root`
- **Numeric overflow detection** — `+`/`-`/`*`/`inc`/`dec` wrap silently on int64 overflow rather than promoting to BigInt; use `+'`/`-'`/`*'` for explicit BigInt math
- **`subseq` / `rsubseq`** — sorted collections themselves work (`sorted-map`, `sorted-set`, `sorted-map-by`, `sorted-set-by`, `rseq`), but range queries on them are not yet implemented

### Known behavioral differences

- **`concat*` (used internally by quasiquote) is eager** — the user-facing `concat` is lazy, matching Clojure
- **All channel operations block** — `<!` and `<!!` are identical (Go channels are always blocking), same for `>!`/`>!!`
- **`go` blocks are real goroutines** — no IOC (inversion of control) state machine like Clojure's core.async; this means they're cheaper but `go` blocks can call blocking ops directly
- **No BigDecimal** — numeric tower is `int64` + `float64` + `BigInt` (no arbitrary-precision decimals)
- **Regex is Go flavor** — `re2` syntax, not Java regex
- **`letfn` uses atoms** internally for forward references — slight overhead vs Clojure's direct binding

## Examples

Real projects written in let-go:

- [**xsofy**](https://github.com/nooga/xsofy) — a roguelike that runs in the browser and the terminal from the same source
- [**lgcr**](https://github.com/nooga/lgcr) — a decent daemonless container runtime, built on the `syscall` namespace

In this repo:

- [examples/](https://github.com/nooga/let-go/tree/main/examples) — small programs
- [test/](https://github.com/nooga/let-go/tree/main/test) — `.lg` test files covering all features

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

## Source paths

By default `lg` looks for namespaces under the current directory. To add
extra directories use `-source-paths` flag or the `LG_SOURCE_PATHS` environment variable. 
Both accept a list separated `:` on Unix, `;` on Windows. The flag wins over
the env var. The current directory is always prepended so project code
takes precedence.

```bash
lg -source-paths /path/to/libA:/path/to/libB script.lg
LG_SOURCE_PATHS=/path/to/libA:/path/to/libB lg script.lg
```

The same paths apply to script execution, `-c` (bytecode compile), `-b`
(standalone bundle), and `-w` (WASM). Bundled binaries (`-b` output) honor
`LG_SOURCE_PATHS` at runtime; the `-source-paths` flag is a
bundle-build-time input.

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

let-go embeds cleanly as a scripting layer for Go programs — define Go values and
functions, hand them to the VM, and run user-supplied Clojure against your data.
Go structs roundtrip as records, Go channels are first-class let-go channels, and
Go functions are callable from let-go code.

```go
import (
    "github.com/nooga/let-go/pkg/api"
    "github.com/nooga/let-go/pkg/vm"
)

c, _ := api.NewLetGo("myapp")

// Expose Go values and functions to let-go
c.Def("x", 42)
c.Def("greet", func(name string) string {
    return "Hello, " + name
})

v, _ := c.Run(`(greet "world")`)
fmt.Println(v) // "Hello, world"
```

**Struct ↔ Record roundtrip.** Registered structs become records on the let-go
side. Unmutated values unbox back to the original Go type for free; mutated ones
go through `vm.ToStruct[T]`.

```go
type Item struct{ Name string; Price float64; Qty int }
vm.RegisterStruct[Item]("myapp/Item")

c.Def("item", Item{Name: "Widget", Price: 9.99, Qty: 5})
c.Run(`(:name item)`)        // "Widget"
c.Run(`(* (:price item) (:qty item))`) // 49.95

// Define a let-go function that processes Go structs
c.Run(`(defn total [it] (* (:price it) (:qty it)))`)
v, _ := c.Run(`(total item)`) // 49.95
```

**Streaming via Go channels.** A Go `chan int` and a `vm.Chan` plug straight
into `go`/`<!`/`>!` — perfect for piping events through a user-supplied script.

```go
inch := make(chan int)
outch := make(vm.Chan)
c.Def("in", inch)
c.Def("out", outch)

c.Run(`(go (loop [i (<! in)]
             (when i
               (>! out (inc i))
               (recur (<! in)))))`)
```

See [`pkg/api/interop_test.go`](pkg/api/interop_test.go) for the full set of
embedding examples (defs, structs, channels, function calls).

## Testing

```bash
go test ./... -count=1 -timeout 30s
```

---

- Ever wanted a 20MB pure-Go JS runtime that typechecks and runs TypeScript? check my other project https://github.com/nooga/paserati

[🤓 Follow me on X](https://x.com/MGasperowicz)
[🐬 Check out monk.io](https://monk.io)
