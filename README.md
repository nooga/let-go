<!--suppress ALL -->
<p align="center">
<img src="meta/logo.png" alt="Squishy loafer" title="Squishy loafer of let-go" />
</p>

![Tests](https://github.com/nooga/let-go/actions/workflows/go.yml/badge.svg)

# let-go

Greetings loafers! _(λ-gophers haha, get it?)_

let-go is a Clojure dialect with a bytecode compiler and stack VM, written in Go.
A single ~10.7MB binary, ~7ms cold start, no JVM. It passes the
[jank-lang test suite](https://github.com/jank-lang/clojure-test-suite).

I started this in 2021 as an elaborate joke: an excuse to write Clojure while
pretending to write Go. It turned out useful. I use it for CLIs, scripts, and
web servers, and I built [a daemonless container runtime](https://github.com/nooga/lgcr)
on top of it. You can compile let-go programs to standalone binaries or
self-contained WASM web pages. [It even runs on Plan 9](https://x.com/MGasperowicz/status/2052428420592599507?s=20),
and ReMarkable 2.

It is not a drop-in replacement for Clojure JVM. It does not load JARs and
does not aim to. Most idiomatic Clojure code runs unmodified, but a real
project with library dependencies will need adjustments. See
[Known limitations](#known-limitations) below.

## Goals (in no particular order)

- [x] Quality entertainment
- [x] Implement most of Clojure: persistent data structures, lazy seqs,
      transducers, protocols, records, multimethods, core.async, BigInts
- [x] Comfy two-way Go interop (functions, structs, channels)
- [x] AOT compilation to bytecode and standalone binaries
- [x] Boot the runtime inside a single `requestAnimationFrame` (10ms left over at 60fps)
- [x] Compile programs to self-contained WASM web pages with terminal emulation
- [ ] Make it legal to write Clojure at your Go dayjob
- [ ] nREPL in the browser (let-go VM in WASM, editor over WebSocket)
- [ ] Stretch: let-go bytecode → Go translation

Non-goals: drop-in JVM Clojure replacement; linter/formatter for Clojure-at-large.

## Benchmarks

let-go vs Babashka, Joker, [go-joker](https://github.com/rcarmo/go-joker),
[gloat](https://github.com/gloathub/gloat), and Clojure JVM. All benchmark
files are valid Clojure that runs unmodified. Apple M1 Pro.

|                 | let-go     | babashka | joker | go-joker | gloat | clojure JVM |
| --------------- | ---------- | -------- | ----- | -------- | ----- | ----------- |
| **Binary size** | **10.7MB** | 68MB     | 26MB  | 32MB     | 26MB  | 304MB (JDK) |
| **Startup**     | **6.7ms**  | 18ms     | 12ms  | 13ms     | 16ms  | 363ms       |
| **Idle memory** | **13.5MB** | 27MB     | 22MB  | 23MB     | 23MB  | 92MB        |

let-go wins decisively on the small things: smallest binary, fastest startup
(~50× under JVM, ~3× under Babashka), lowest memory. It also wins on
short-lived data work like map/filter (7.9ms vs Babashka's 21.5ms) and
persistent maps (20.8ms vs 23.7ms).

On bigger numerical workloads other implementations pull ahead. go-joker's
WASM JIT compiles inner numeric loops and beats us on fib (1.47s vs 2.08s),
tak, reduce, and transducers. The JVM dominates on long compute runs once
HotSpot warms up. We're about even with Babashka on most algorithmic
benchmarks and 10×+ faster than upstream Joker (bytecode VM vs tree-walk).

Full per-benchmark numbers and methodology:
[benchmark/results.md](benchmark/results.md).

## Compatibility

Tested against [jank-lang/clojure-test-suite](https://github.com/jank-lang/clojure-test-suite):
**5621 / 5621 assertions pass** across 232 files through the `:clj` reader
lens, with no known failures, compile skips, panic skips, or runtime skips.

### Standard namespaces

| Namespace            | Status                                                                                                                                                                                        |
| -------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `clojure.core`       | macros, destructuring, lazy seqs, transducers, protocols, records, `deftype`, `reify`, multimethods, hierarchies, atoms, regex, metadata, BigInt, BigDecimal                                   |
| `clojure.string`     | full                                                                                                                                                                                          |
| `clojure.set`        | full                                                                                                                                                                                          |
| `clojure.walk`       | `prewalk`, `postwalk`, `keywordize-keys`, `stringify-keys`, `walk`                                                                                                                            |
| `clojure.edn`        | `read`, `read-string`                                                                                                                                                                         |
| `clojure.pprint`     | `pprint`, `cl-format`                                                                                                                                                                         |
| `clojure.test`       | `deftest`, `is`, `testing`, `are`, fixtures                                                                                                                                                   |
| `clojure.core.async` | channels, `go`/`go-loop`, `alts!`, `mult`/`pub`, `pipe`/`merge`/`split` (real goroutines, not IOC)                                                                                            |
| `io`                 | polymorphic readers/writers, `slurp`/`spit`, lazy line-seq, encoding, URLs, `with-open`                                                                                                       |
| `http`               | Ring-style server + client, streaming responses                                                                                                                                               |
| `json`               | `read-json`, `write-json` (float-preserving, record-aware)                                                                                                                                    |
| `transit`            | transit+json codec with rolling cache                                                                                                                                                         |
| `os`                 | `sh`, `stat`, `ls`, `cwd`, `getenv`/`setenv`, `exit`, `os-name`, `arch`, `user-name`, `hostname`, separators                                                                                  |
| `System`             | JVM-shaped: `getProperty`, `getProperties`, `getenv`, `exit`, `currentTimeMillis`, `nanoTime`. Exposes `let-go.version`, `let-go.commit`, `user.home`, `user.dir`, `os.name`, `os.arch`, etc. |
| `syscall`            | direct Linux syscalls (mount, unshare, mknod, prctl, capset, seccomp, AppArmor)                                                                                                               |
| `pods`               | Babashka pods over JSON / EDN / transit                                                                                                                                                       |

### Babashka pods

let-go can load [Babashka pods](https://github.com/babashka/pods), which
opens up the whole pod ecosystem: SQLite, AWS, Docker, file watching, etc.

```clojure
(pods/load-pod 'org.babashka/go-sqlite3 "0.3.13")

(pod.babashka.go-sqlite3/execute! "app.db"
  ["create table users (id integer primary key, name text)"])
(pod.babashka.go-sqlite3/query "app.db"
  ["select * from users"])
;; => [{:id 1 :name "Alice"}]
```

It shares `~/.babashka/pods/` with `bb`, so install pods with babashka and use
them from `lg`. See the [pod registry](https://github.com/babashka/pod-registry)
for what's available.

## Known limitations

### Not implemented

- **STM coordination**: `ref`/`dosync`/`alter`/`commute` are atom-backed compatibility aliases, not coordinated STM
- **Asynchronous agents**: `agent`/`send`/`send-off` are synchronous atom-backed compatibility aliases
- **Chunked sequences**: lazy seqs are unchunked
- **Custom tagged literal readers**: built-in `#uuid` and `#inst` work; unknown tags read as their payload, and `*data-readers*` / `*default-data-reader-fn*` are not implemented
- **Java-style `deftype` / `reify` method bodies and host interfaces**: protocol implementations work; JVM host methods do not
- **Spec** (no `clojure.spec`)
- **`subseq` / `rsubseq`**: sorted collections work (`sorted-map`, `sorted-set`, `rseq`); range queries don't

### Behavioral differences

- `concat*` (used internally by quasiquote) is eager; user-facing `concat` is lazy
- `<!` / `<!!` are identical, same for `>!` / `>!!` (Go channels always block)
- `go` blocks are real goroutines, not IOC state machines (cheaper, and they can call blocking ops directly)
- Numeric tower is pragmatic: `int64`, `float64`, `BigInt`, ratios, and `BigDecimal`, without the JVM's full primitive/class model
- Base integer `+`/`-`/`*`/`inc`/`dec` throw on overflow; use `+'`/`-'`/`*'`/`inc'`/`dec'` for BigInt-promoting exact math
- Regex is Go flavor (`re2`), not Java regex
- `letfn` uses atoms internally for forward references

## Examples

Things written in let-go:

- [**xsofy**](https://github.com/nooga/xsofy): a roguelike that runs in the browser and the terminal from the same source
- [**lgcr**](https://github.com/nooga/lgcr): a daemonless container runtime, built on the `syscall` namespace

In this repo:

- [examples/](https://github.com/nooga/let-go/tree/main/examples): small programs
- [test/](https://github.com/nooga/let-go/tree/main/test): `.lg` test files

## Try it online

[Bare-bones browser REPL](https://nooga.github.io/let-go/), running a WASM
build of let-go.

## Install

### Homebrew (macOS / Linux)

```bash
brew tap nooga/let-go https://github.com/nooga/let-go
brew install let-go
```

### Download

Prebuilt binaries for Linux, macOS, and Plan 9 in
[Releases](https://github.com/nooga/let-go/releases).

### From source (Go 1.26+)

```bash
go install github.com/nooga/let-go@latest
```

### Usage

```bash
lg                                # REPL
lg -e '(+ 1 1)'                   # eval expression
lg myfile.lg                      # run file
lg -r myfile.lg                   # run file, then REPL
```

## Compile and distribute

let-go can compile programs to bytecode (`.lgb` files) and bundle them as
standalone executables.

```bash
lg -c app.lgb app.lg              # compile to bytecode
lg app.lgb                        # run bytecode

lg -b myapp app.lg                # bundle into a self-contained binary
./myapp                           # runs anywhere, no lg needed
```

The standalone binary is a copy of `lg` with your bytecode appended. Copy it
to another machine and it runs.

```bash
lg -w site app.lg                 # compile to a WASM web app
open site/index.html
```

The output is a self-contained `index.html` (~6MB, inlined WASM, gzipped)
plus a service worker that supplies the COOP/COEP headers GitHub Pages needs
for SharedArrayBuffer. Programs that use the `term` namespace get full
terminal emulation via xterm.js: ANSI colors, cursor positioning, raw
keyboard input.

The `*compiling-aot*` var is `true` during `-c`/`-b`/`-w` compilation and
`false` at runtime, useful for keeping side effects out of compile time:

```clojure
(defn -main []
  (start-server))

(when-not *compiling-aot*
  (-main))
```

`*in-wasm*` is `true` when running inside a WASM build.

## nREPL

let-go ships an nREPL server that works with CIDER (Emacs), Calva (VS Code),
and Conjure (Neovim).

```bash
lg -n                             # default port 2137
lg -n -p 7888
```

It writes `.nrepl-port` to the working directory so editors auto-discover it.

Supported ops: `clone`, `close`, `eval`, `load-file`, `describe`,
`completions`, `complete`, `info`, `lookup`, `ls-sessions`, `interrupt`.

- **Emacs (CIDER)**: `M-x cider-connect-clj`, `localhost`, port from `.nrepl-port`
- **VS Code (Calva)**: open a let-go project (the bundled `.vscode/settings.json` registers a connect sequence). Use "Calva: Start a Project REPL and Connect (Jack-In)" → "let-go", or "Calva: Connect to a Running REPL Server" if nREPL is already up.
- **Neovim (Conjure)**: auto-connects when `.nrepl-port` exists.

## Embedding in Go

let-go embeds cleanly as a scripting layer for Go programs. Define Go values
and functions, hand them to the VM, run user-supplied Clojure against your
data. Go structs roundtrip as records, Go channels are first-class let-go
channels, and Go functions are callable from let-go.

```go
import (
    "github.com/nooga/let-go/pkg/api"
    "github.com/nooga/let-go/pkg/vm"
)

c, _ := api.NewLetGo("myapp")

c.Def("x", 42)
c.Def("greet", func(name string) string {
    return "Hello, " + name
})

v, _ := c.Run(`(greet "world")`)
fmt.Println(v) // "Hello, world"
```

Registered structs become records on the let-go side. Unmutated values unbox
back to the original Go type for free; mutated ones go through `vm.ToStruct[T]`.

```go
type Item struct{ Name string; Price float64; Qty int }
vm.RegisterStruct[Item]("myapp/Item")

c.Def("item", Item{Name: "Widget", Price: 9.99, Qty: 5})
c.Run(`(:name item)`)                  // "Widget"
c.Run(`(* (:price item) (:qty item))`) // 49.95
```

Go channels and `vm.Chan` plug into `go` / `<!` / `>!` directly:

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

[`pkg/api/interop_test.go`](pkg/api/interop_test.go) has the full set of
embedding examples (defs, structs, channels, function calls).

## Testing

```bash
go test ./... -count=1 -timeout 30s
```

## Contributing

### Git merge driver for `core_compiled.lgb` (one-time setup)

`pkg/rt/core_compiled.lgb` is a binary bundle regenerated from the embedded
`.lg` sources. Git cannot meaningfully merge this binary on rebase, so we ship
a custom merge driver that regenerates from sources after the `.lg` files have
been merged as text.

After cloning the repo (or pulling for the first time after this driver was
added), register it locally:

```bash
make install-hooks
```

(A merge driver lives in `.git/config`, which is not shared, so each clone
needs this once. The target just runs the `git config` commands for you.)

After registration, rebases and merges that touch any embedded `.lg` source
will regenerate the `.lgb` automatically — no more binary merge conflicts when
stacking PRs that edit `core.lg` and friends.

---

Ever wanted a 20MB pure-Go JS runtime that typechecks and runs TypeScript?
Check my other project: https://github.com/nooga/paserati

[🤓 Follow me on X](https://x.com/MGasperowicz)
[🐬 Check out monk.io](https://monk.io)
