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

Core namespaces cover `clojure.core` (macros, lazy seqs, transducers, protocols,
records, multimethods, BigInt/BigDecimal) plus `string`, `set`, `walk`, `edn`,
`pprint`, `test`, and `core.async`, alongside let-go's own `io`, `http`, `json`,
`transit`, `os`, `System`, `syscall`, and `pods`. See
[docs/guide/clojure-compatibility.md](docs/guide/clojure-compatibility.md) for
the full per-namespace status table and the Clojure differences.

### Babashka pods

let-go can load [Babashka pods](https://github.com/babashka/pods), opening up the
whole pod ecosystem (SQLite, AWS, Docker, file watching, …) and sharing
`~/.babashka/pods/` with `bb`.

```clojure
(pods/load-pod 'org.babashka/go-sqlite3 "0.3.13")
(pod.babashka.go-sqlite3/query "app.db" ["select * from users"])
```

See [docs/guide/pods.md](docs/guide/pods.md) for a full example and the shared
pod cache.

### Portable code (`:lg` reader conditionals)

let-go ships namespaces of its own (e.g. `let-go.semver`) that JVM Clojure can't
load. To keep shared code loadable on both, guard the let-go-only parts behind
`:lg` reader conditionals in a `.cljc` file — JVM Clojure skips `:lg` branches
the same way it skips `:cljs`:

```clojure
(ns my.app
  #?(:lg (:require [let-go.semver :as semver])))   ; only let-go loads this
```

The guard is at **read** time, so a missing namespace never reaches compilation.
See [docs/guide/portability.md](docs/guide/portability.md) for the `.cljc`
resolution rule and `:lg`/`:clj` ordering gotcha.

### Version requirements (`let-go.semver`)

`let-go.semver` provides SemVer values that sort correctly, npm/cargo-style range
matching (`satisfies-range?` — comparators, x-ranges, `^`/`~`, `||`), and
`require-letgo`, which asserts at load time that the running `lg` build is new
enough and fails with one clear line instead of a "can't resolve" cascade:

```clojure
(ns my.app
  #?(:lg (:require [let-go.semver :refer [require-letgo]])))

#?(:lg (require-letgo ">=1.9.0"))   ; one clear failure line on too-old lg
```

Guard it behind [`:lg` reader conditionals](#portable-code-lg-reader-conditionals)
so shared `.cljc` stays JVM-loadable. See
[docs/guide/semver.md](docs/guide/semver.md) for the range grammar and
`require-letgo`'s detection/failure semantics.

## Known limitations

Not a drop-in JVM Clojure. The main gaps: no coordinated STM or async agents
(`ref`/`agent` are atom-backed aliases), no `clojure.spec`, unchunked lazy seqs,
no custom `*data-readers*`, no JVM host interop on `deftype`/`reify`, and no
`subseq`/`rsubseq` range queries. Behavior also differs in places — pragmatic
numeric tower, always-blocking channels, real-goroutine `go` blocks, and `re2`
(not Java) regex.

Full list with rationale:
[docs/guide/clojure-compatibility.md](docs/guide/clojure-compatibility.md).

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
brew install nooga/tap/let-go
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
lg myfile.lg a b                  # run file with arguments
lg -r myfile.lg                   # run file, then REPL
```

`*command-line-args*` holds the program's arguments — the positionals after the
script — as a seq of strings, or `nil` when there are none. It reads the same
whether you run a script or a bundled binary, so you never slice argv by hand:

```clojure
;; greet.lg — run as `lg greet.lg Alice Bob` or `./greet Alice Bob`
(doseq [name *command-line-args*]
  (println "Hello," name))
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

The WASM output is a self-contained `index.html` (~6MB, inlined and gzipped) with
a service worker for the COOP/COEP headers SharedArrayBuffer needs; `term`-using
programs get full xterm.js terminal emulation.

See [docs/guide/usage.md](docs/guide/usage.md) for the `*compiling-aot*` /
`*in-wasm*` compile-time vars, more on each output format, and project/dependency
management with [lgx](https://github.com/abogoyavlensky/lgx).

### Resources and source paths

Programs read non-source files (templates, web assets, data) via `io/resource`,
with roots set by `-resource-paths` / `LG_RESOURCE_PATHS`. Bundling with `-b`
embeds every file under those roots, so a bundled binary is self-contained.

```clojure
(when-let [r (io/resource "templates/index.html")]
  (io/slurp r))
```

`require`d namespaces resolve against `-source-paths` / `LG_SOURCE_PATHS`
(default `.`). When you set the search path it's taken as the **complete** list —
the current directory isn't added implicitly.

See
[docs/guide/resources-and-source-paths.md](docs/guide/resources-and-source-paths.md)
for path-list syntax, multi-root precedence, embedding behavior, and the
empty-value/explicit-only rules.

## nREPL

let-go ships an nREPL server that works with CIDER (Emacs), Calva (VS Code), and
Conjure (Neovim). It writes `.nrepl-port` to the working directory so editors
auto-discover it.

```bash
lg -n                             # default port 2137
lg -n -p 7888
```

See [docs/guide/nrepl.md](docs/guide/nrepl.md) for supported ops and per-editor
connect steps.

## Embedding in Go

let-go embeds cleanly as a scripting layer for Go programs: define Go values and
functions, hand them to the VM, run user-supplied Clojure against your data. Go
structs roundtrip as records, Go channels are first-class let-go channels, and Go
functions are callable from let-go.

```go
c, _ := api.NewLetGo("myapp")
c.Def("greet", func(name string) string { return "Hello, " + name })
v, _ := c.Run(`(greet "world")`)   // "Hello, world"
```

See [docs/guide/embedding-in-go.md](docs/guide/embedding-in-go.md) for struct
roundtripping, Go-channel interop, and a pointer to the full example set.

## Testing

```bash
go test ./... -count=1 -timeout 30s
```

## Contributing

After cloning, run `make install-hooks` once to register the `core_compiled.lgb`
merge driver (each clone needs this — the config lives in `.git/config`, which
isn't shared). See
[docs/regenerating-generated-artifacts.md](docs/regenerating-generated-artifacts.md)
for how generated artifacts are regenerated and kept in sync.

---

Ever wanted a 20MB pure-Go JS runtime that typechecks and runs TypeScript?
Check my other project: https://github.com/nooga/paserati

[🤓 Follow me on X](https://x.com/MGasperowicz)
[🐬 Check out monk.io](https://monk.io)
