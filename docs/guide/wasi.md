---
status: active
last-verified: 2026-07-13
human-verified:
---

# Running let-go as a WASI (`GOOS=wasip1`) module

let-go builds to a standalone WASI module with the standard Go toolchain — no
TinyGo, no build tags beyond the target. The result runs under any wasip1 host
(wasmtime, wazero, …) as a headless compute + stdio runtime, with full 64-bit
integer fidelity.

## Build

```
GOOS=wasip1 GOARCH=wasm go build -o lg.wasm .
```

The interactive REPL and terminal code are gated off wasip1 (they depend on
`chzyer/readline` and `x/sys/unix` poll/ioctl, which have no wasip1 backing), so
the wasip1 build routes to the non-interactive entry point automatically. The
`wasip1-build` CI job builds this target on every PR so it can't silently
regress.

## Run

Any wasip1 host works. Under wasmtime:

```
$ wasmtime lg.wasm -e '(println (str "wasi says " (* 6 7)))'
wasi says 42
```

Integer arithmetic is 64-bit, matching the native build:

```
$ wasmtime lg.wasm -e '(* 1000000000 1000000000)'
1000000000000000000
$ wasmtime lg.wasm -e '(bit-shift-left 1 62)'
4611686018427387904
```

## TinyGo

This guide covers the standard Go toolchain. TinyGo's `-target=wasi` is also
`GOOS=wasip1`, so the gating here applies, but building let-go under TinyGo needs
additional reflect shims (boxing and method dispatch) that are not part of this
build — without them, stock TinyGo 0.41.1 traps at initialization with
`unimplemented: (reflect.Type).Method()` before running any code. TinyGo support
is tracked separately.

When that enablement is in place, TinyGo trades the ~24 MB standard-Go module
for a much smaller one (~9 MB), but its wasm `int` is 32-bit, so 64-bit
arithmetic overflows (let-go traps on overflow rather than wrapping — loud, not
silently corrupting). Use the standard-Go build when numeric correctness matters,
which is the general case for a Clojure dialect.

## Limits

wasip1 is compute + stdio only. Specifically:

- **No `term` namespace.** The interactive terminal is excluded, so there is no
  `term` ns under wasip1; requiring it reports `unable to load namespace term`
  rather than failing hard. Fine for headless eval and embedding; a TUI-over-wasi
  would need a headless `term` stub (future work).
- **No sockets, no threads.** The wasip1 preview is single-threaded with no
  network. Compute and stdio only.

For loading the module from a Go host rather than a CLI runtime, see
[embedding in Go](embedding-in-go.md); for the I/O seams the host binds, see
[decoupling runtime I/O from the host](../design/runtime-io-host-decoupling.md).
