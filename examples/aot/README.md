# AOT compilation examples

Trials for compiling let-go `.lg`/`.clj(c)` namespaces **ahead-of-time to native
Go**, via `scripts/lg-compile` (the multi-namespace AOT driver with direct
cross-package calls).

Driver:

```
./lg scripts/lg-compile <out-dir> <import-prefix> <file.lg>...
```

It lowers each single-arity `defn`/`defmulti`/`defmethod` to a native Go func,
emits one Go package per namespace (canonical `gogen/ns->go-pkg` naming: leaf
segment + nested dir), and wires direct `pkg.Fn(ec, …)` calls between the
lowered packages. Variadic / non-coercible fns stay on runtime trampolines.

| Directory | What it is |
|---|---|
| `cross-package/` | Minimal two-namespace demo (`lib` + `app`) exercising a direct cross-package call. The smallest reproducible AOT example. |
| `ys-on-let-go/`  | The YS-on-let-go shim: a small YAMLScript-style stdlib plus example programs. See its own `README.md`. |
| `yamlstar/`      | Diagnostics for AOT-compiling the real YamlStar Clojure pipeline (`ys-coverage.lg` — per-file lowerability tally; `ys-lower-report.lg` — per-fn lower/fallback report). |

Generated Go output (`*/go/`, `*/go_lowered/`) is git-ignored — regenerate with
`lg-compile`.
