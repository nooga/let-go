# Browser Inspector Example

Minimal client-owned browser shell for the `LetGoHost.request(...)` bridge.

It builds a shell-less `lg -w -w-host-eval` bundle, then injects a small
workbench UI that drives:

- `eval`
- `compile`
- `inspect-all`

Current scope:

- single embedded `"default"` session
- real REPL compile path
- bytecode disassembly pane
- IR, optimized bytecode, and lowered Go panes, each requested per cell
  through the Analysis toggles ("Inspect Active Cell" enables all four)

## Build

```bash
make -C examples/browser-inspector build
# or, with a prebuilt lg binary:
LG=./lg make -C examples/browser-inspector build
# with IR/native lowering in the wasm app:
LG_WASM_BUILD_TAGS=gogen_ir make -C examples/browser-inspector build
# from the repo root:
make browser-inspector
```

Output: `examples/browser-inspector/dist/index.html`, plus the
`coi-serviceworker.js` that `lg -w` writes beside every page it builds (the
COOP/COEP shim the page registers). `dist/` is entirely generated and
gitignored; `make clean` removes it.

## Serve

Any static server is fine. Example:

```bash
cd examples/browser-inspector/dist
python3 -m http.server
```

Then open `http://localhost:8000`.
