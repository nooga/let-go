---
status: active
last-verified: 2026-06-19
human-verified:
---

# Clojure compatibility and differences

let-go is a Clojure dialect, not a drop-in JVM Clojure: most idiomatic Clojure
runs unmodified, but it doesn't load JARs and the host-interop and concurrency
models differ. Tested against
[jank-lang/clojure-test-suite](https://github.com/jank-lang/clojure-test-suite):
**5621 / 5621 assertions pass** across 232 files through the `:clj` reader lens,
with no known failures, compile skips, panic skips, or runtime skips.

## Standard namespaces

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
| `io`                 | polymorphic readers/writers, `slurp`/`spit`, lazy line-seq, encoding, URLs, `with-open`, `resource` (filesystem in dev, embedded in `-b` binaries)                                                |
| `http`               | Ring-style server + client, streaming responses                                                                                                                                               |
| `json`               | `read-json`, `write-json` (float-preserving, record-aware)                                                                                                                                    |
| `transit`            | transit+json codec with rolling cache                                                                                                                                                         |
| `os`                 | `sh`, `stat`, `ls`, `cwd`, `getenv`/`setenv`, `exit`, `os-name`, `arch`, `user-name`, `hostname`, separators                                                                                  |
| `System`             | JVM-shaped: `getProperty`, `getProperties`, `getenv`, `exit`, `currentTimeMillis`, `nanoTime`. Exposes `let-go.version`, `let-go.commit`, `user.home`, `user.dir`, `os.name`, `os.arch`, etc. |
| `syscall`            | direct Linux syscalls (mount, unshare, mknod, prctl, capset, seccomp, AppArmor)                                                                                                               |
| `pods`               | Babashka pods over JSON / EDN / transit                                                                                                                                                       |

## Not implemented

- **STM coordination**: `ref`/`dosync`/`alter`/`commute` are atom-backed compatibility aliases, not coordinated STM
- **Asynchronous agents**: `agent`/`send`/`send-off` are synchronous atom-backed compatibility aliases
- **Chunked sequences**: lazy seqs are unchunked
- **Custom tagged literal readers**: built-in `#uuid` and `#inst` work; unknown tags read as their payload, and `*data-readers*` / `*default-data-reader-fn*` are not implemented
- **Java-style `deftype` / `reify` method bodies and host interfaces**: protocol implementations work; JVM host methods do not
- **Spec** (no `clojure.spec`)
- **`subseq` / `rsubseq`**: sorted collections work (`sorted-map`, `sorted-set`, `rseq`); range queries don't

## Behavioral differences

- `concat*` (used internally by quasiquote) is eager; user-facing `concat` is lazy
- `<!` / `<!!` are identical, same for `>!` / `>!!` (Go channels always block)
- `go` blocks are real goroutines, not IOC state machines (cheaper, and they can call blocking ops directly)
- Numeric tower is pragmatic: `int64`, `float64`, `BigInt`, ratios, and `BigDecimal`, without the JVM's full primitive/class model
- Base integer `+`/`-`/`*`/`inc`/`dec` throw on overflow; use `+'`/`-'`/`*'`/`inc'`/`dec'` for BigInt-promoting exact math
- Regex is Go flavor (`re2`), not Java regex
- `letfn` uses atoms internally for forward references
