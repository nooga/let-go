# Native-entry frame demo (#425 Item 1)

Proves the AOT native-entry path end-to-end: `lg-compile --entry-frame`
lowers `fib` + `-main`, emits a `main.go` frame keyed off the recognized
entry, the build script places `program.lgb` beside it, and `go build`
produces a binary that boots via `rt.BootCore`, loads the program namespaces
(with override drain), then enters natively at `prog.Main`.

Frame emission is opt-in (`--entry-frame`) so package-only callers such as
Gloat keep historical behavior and own their own executable templates.

```
./examples/aot/native-entry/build.sh
./examples/aot/native-entry/out/fib.native       # => 55
./examples/aot/native-entry/out/fib.native x     # => 5702887 (fib 34)
```

On a typical machine fib(34) is ~0.07–0.08s native vs ~1.1–1.4s on the lg VM
(~16–17×), matching the #425 spike.

## What lg-compile emits

- `out/fib/fib.go` — lowered `Fib` / `Main`
- `out/main.go` — the #425 frame (`BootCore` → `LoadProgramNamespaces` →
  `prog.Main(ec, argv...)`)
- summary line: `2 fns lowered; native entry: Main ✓`

`program.lgb` is owned by the orchestrator (this `build.sh`, or gloat) — not
by lg-compile — and must sit next to `main.go` for the `//go:embed`.

## Arity / privacy coverage

The frame generator also handles `[argv]`, `[]`, and private `(defn- main …)`
(the latter gets a same-package `NativeMain` bridge). This demo uses the spike
shape: public variadic `-main`.

The entry must have a single arity. The frame calls one Go signature, so a
multi-arity `(defn -main ([] …) ([x] …))` is rejected with a diagnostic rather
than compiled; `[first & rest]` and multi-fixed `[a b]` are rejected the same
way. Supported shapes are `[]`, `[argv]`, and `[& args]`.
