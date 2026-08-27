---
status: active
last-verified: 2026-08-11
human-verified: 2026-08-11
---

# let-go under TinyGo

TinyGo produces much smaller binaries and roughly a third the RSS of the stock
Go build, which suits two targets: browser WASM bundles small enough to ship
inline, and constrained native hosts (edge, embedded, Raspberry Pi-class). Both
build and run — the `lg` CLI natively on linux, and real apps to WASM.

The changes that enable this are build-tag-gated (`//go:build tinygo`), so a
normal `go build` is unaffected.

## Building

Native (linux):

```
tinygo build -o lg .
```

WASM, via the `-w` bundler with the TinyGo backend:

```
LETGO_USE_TINYGO=1 lg -w -o out app.lg
```

A stack-size override is required for the WASM build — the interpreter's
recursion overflows TinyGo's default goroutine stack, and TinyGo's wasm stacks
are not growable. The bundler sets one; raise it if a deep program still faults.
The `tinygo-wasi-build` CI job builds the runtime-only binary and boots it under
wasmtime on every PR, so the target can't silently regress.

## Limitations

- **The `os` namespace is minimal.** Only `os/exit`, `os/getenv`, and `os/args`
  are registered; `os/args` is currently an empty vector. Process execution,
  filesystem and path operations, and host-information functions require a
  standard-Go build.
- **Reflect is partial.** TinyGo does not implement `reflect.Type.Method` or
  `reflect.Value.Call`, which let-go's Go-interop boxing uses. The affected paths
  are shimmed: boxed-method dispatch hand-registers the methods programs reach
  for (currently `time.Time.Sub`), and a reflect-boxed Go function installs a
  stub. Unsupported interop errors explicitly at the call site — an unshimmed
  method or an invoked stub reports TinyGo's limitation rather than returning nil
  or trapping the module at boot. Programs that stay within core plus shimmed
  interop run faithfully.
- **WASM `int` is 32-bit** (stock Go's wasm `int` is 64-bit), so 64-bit integer
  arithmetic overflows or wraps on the WASM target. Native TinyGo (arm64/amd64)
  is 64-bit and unaffected. Use the standard-Go WASI build when 64-bit fidelity
  matters.
- **`xxh3` is not bit-compatible.** The `xxh3` dependency is assembly-only on
  arm64/amd64, so TinyGo links a pure-Go substitute hash. It is deterministic but
  produces different digests, making a TinyGo build its own determinism domain.
- **No Unix domain sockets or full `os.ProcessState`**, so the nREPL Unix-socket
  transport and some process introspection are unavailable on native builds.
- **Native macOS does not link** (tinygo-org/tinygo#4794 — Darwin syscalls route
  through libSystem assembly TinyGo can't link). Linux is the native target.
- **Concurrency yes, parallelism no; `alts!` excepted.** wasm has no threads —
  `runtime.NumCPU()` and `GOMAXPROCS` are both 1 — so goroutines interleave on a
  single thread instead of running in parallel (this holds for the stock-Go wasm
  build too). The cooperative channel surface runs faithfully on native Go
  channels: `go*`, `>!`/`<!`, `chan` and its buffers, `timeout`, `pipe`, and
  `promise-chan`. The one hole is `alts!`/`alts!!`, which use `reflect.Select` —
  unimplemented in TinyGo (see the reflect limitation above), so a call panics
  rather than erroring at the call site. Parallel *throughput* also doesn't carry
  over: `pmapv` and other `NumCPU`-fanned work see a single worker and run
  sequentially — same results, no speedup. Use the standard-Go native build when
  `alts!` or real parallelism matter.
