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
- explicit placeholder panes for IR / optimized bytecode / lowered Go until
  those ops are implemented on the bridge

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

Output: `examples/browser-inspector/dist/index.html`

## Serve

Any static server is fine. Example:

```bash
cd examples/browser-inspector/dist
python3 -m http.server
```

Then open `http://localhost:8000`.
