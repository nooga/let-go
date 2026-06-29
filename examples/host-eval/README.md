# host-eval browser demo

A client-owned shell driving a live let-go image in the browser through
`window.LetGoHost.eval` — the `-w-host-eval` bundle feature.

The hosted REPL under `wasm/` is a bespoke Go-WASM build with its own
`window.Eval`. This is the other path: a plain `lg -w` bundle where the page
owns the shell and binds the runtime only through the public `window.LetGoHost`
surface. It demonstrates two seams at once:

- **eval in** — `LetGoHost.eval(code)` runs a form in the loaded image and
  resolves to its value. The REPL pane shows the returned values.
- **stdout out** — `*out*` routes through a `HostWriter` to `onOutput`. The
  stdout pane shows what the image prints. The `time` and `fizzbuzz` chips
  exercise both at once (a return value *and* printed lines).

## Build

```
LG=lg ./build.sh
```

This produces a single `dist/index.html` with the wasm inlined. It still pulls
xterm and the JetBrains Mono web font from a CDN at runtime, so it needs network
access (vendor those assets if you want a fully offline page). Serve it **with
cross-origin isolation** — the wasm input ring is `SharedArrayBuffer`-backed, so
the page must be `crossOriginIsolated`, which needs:

```
Cross-Origin-Opener-Policy: same-origin
Cross-Origin-Embedder-Policy: require-corp
```

Without those headers the runtime can't read input and boot fails.

| File | What it is |
|---|---|
| `main.lg`    | Minimal program — prints a banner, then the image stays live for `eval`. |
| `shell.html` | The client shell: an inline REPL plus a read-only stdout pane, bound only via `window.LetGoHost`. |
| `build.sh`   | Builds the bundle (`-w-shell none -w-host-eval`) and grafts `shell.html` into the page. |
