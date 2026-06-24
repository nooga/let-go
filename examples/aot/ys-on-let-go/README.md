# Running YAMLScript (YS) on let-go

[YAMLScript](https://yamlscript.org/) is a Clojure-based YAML loader and
scripting language. Its compiler emits plain Clojure forms, which means
the *output* of `ys --compile` can run on let-go with a small shim that
provides the symbols from the YS standard library (`say`, `ARGS`, `%`,
`sum`, ...).

This directory contains four worked examples. The `.ys` / `.yaml` files
are the sources; the `.lg` files are checked in so they can be run
without `ys` installed.

## Run the examples

```bash
cd examples/aot/ys-on-let-go

lg -source-paths . 01_hello.lg
lg -source-paths . 02_fizzbuzz.lg
lg -source-paths . 03_classify.lg
lg -source-paths . 04_data.lg
```

Each one prints the same output it would print under native `ys`.

## Regenerate after editing a `.ys` source

```bash
./build.sh    # requires ys in PATH
```

## What this covers and what it doesn't

The `ys-runtime.lg` shim handles enough of the YS surface to run YS
scripts that stay inside Clojure's pure-data and simple-logic
fragments: `defn`, `let`, `cond`, `doseq`, `range`, string
interpolation (compiled to `str`), `apply`, `get`, `hash-map`, plus YS
helpers like `say`/`sum`/`%`.

It does not handle:

- `babashka.fs` / `babashka.process` / `babashka.http-client`
- `java-time`
- Ordered maps (YAML round-trip preserves data but not key order; YS's
  `%` is aliased to `hash-map` here, not `flatland.ordered/ordered-map`)
- YS stdlib symbols whose let-go equivalents live under different
  names (e.g. `json/load` vs let-go's `json/read-json`). For now those
  need a hand-edit or sed pass.

The compiled-output approach is the lightest of three integration
paths discussed in [issue
#49](https://github.com/nooga/let-go/issues/49). For full embedded
yamlscript, SCI itself would need to be ported (SCI is `.cljc` and
let-go already supports `:lg` reader conditionals, so that is
tractable but not free).
