---
status: active
last-verified: 2026-07-13
human-verified:
---

# Running let-go as a WASI (`GOOS=wasip1`) module

let-go builds to a standalone WASI module with the standard Go toolchain — no
TinyGo, no build tags beyond the target. The result runs under any wasip1 host
(wasmtime, wazero, …) as a headless runtime with full 64-bit integer fidelity.

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

This guide is the standard-Go build. TinyGo's `-target=wasi` is also `GOOS=wasip1`
but doesn't build let-go as-is today — see [let-go under TinyGo](tinygo.md) for
the known blocker.

## Capabilities

The module has no ambient authority; a wasip1 host grants what it needs.

- **No `term` namespace.** The interactive terminal is excluded, so there is no
  `term` ns under wasip1; requiring it reports `unable to load namespace term`
  rather than failing hard. A TUI-over-wasi would need a headless `term` stub
  (future work).
- **Filesystem is host-controlled.** Nothing is reachable by default, but the
  host can preopen directories — e.g. `wasmtime --dir=/tmp lg.wasm …` — and
  `slurp` and other file operations then see them. This is the WASI capability
  model: opt-in, host-granted, no ambient filesystem.
- **No sockets, no threads.** The wasip1 preview is single-threaded with no
  network.

For loading the module from a Go host rather than a CLI runtime, see
[embedding in Go](embedding-in-go.md); for the I/O seams the host binds, see
[decoupling runtime I/O from the host](../design/runtime-io-host-decoupling.md).
