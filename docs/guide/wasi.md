---
status: active
last-verified: 2026-07-12
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

## Toolchain: size vs. fidelity

The same 3-tag gating also lets TinyGo (`-target=wasi`, itself `GOOS=wasip1`)
build the module, and TinyGo produces a much smaller artifact — at the cost of
numeric fidelity, since its wasm `int` is 32-bit. let-go traps on overflow
rather than wrapping, so the break is loud, not silently corrupting.

| Toolchain | Module size | 64-bit `int` |
|---|---|---|
| standard Go | ~24 MB | faithful (10^18, 2^62 exact) |
| TinyGo `-target=wasi` | ~9 MB | 32-bit; `(* 1e9 1e9)` overflows |

Pick the standard-Go build when numeric correctness matters (the general case
for a Clojure dialect); reach for TinyGo only where you can accept 32-bit ints
in exchange for the smaller module.

## Limits

wasip1 is compute + stdio only. Specifically:

- **No `term` namespace.** The interactive terminal is excluded, so there is no
  `term` ns under wasip1. Fine for headless eval and embedding; a TUI-over-wasi
  would need a headless `term` stub (future work).
- **No sockets, no threads.** The wasip1 preview is single-threaded with no
  network. Compute and stdio only.

For loading the module from a Go host rather than a CLI runtime, see
[embedding in Go](embedding-in-go.md); for the I/O seams the host binds, see
[decoupling runtime I/O from the host](../design/runtime-io-host-decoupling.md).
