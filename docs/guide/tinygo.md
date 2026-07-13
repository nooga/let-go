---
status: active
last-verified: 2026-07-13
human-verified:
---

# let-go under TinyGo

**Status: not currently supported.** This note records what's known, surfaced
while adding the standard-Go wasip1 target (which shares `GOARCH=wasm`).

TinyGo's `-target=wasi` is also `GOOS=wasip1`, so the wasip1 build gating applies
to it too — but let-go does not build and run under TinyGo as-is:

- **Traps at initialization.** With stock TinyGo 0.41.1 the module builds, but
  traps before executing any code: `unimplemented: (reflect.Type).Method()`.
  let-go's runtime boxes Go values through `reflect`, which TinyGo does not fully
  implement, so namespace installation fails at startup.

Until that's resolved, use the standard-Go build — see
[Running let-go as a WASI module](wasi.md).
