---
status: active
last-verified: 2026-06-19
human-verified:
---

# Usage: running, compiling, distributing

## Running programs

```bash
lg                                # REPL
lg -e '(+ 1 1)'                   # eval expression
lg myfile.lg                      # run file
lg myfile.lg a b                  # run file with arguments
lg -r myfile.lg                   # run file, then REPL
```

`*command-line-args*` holds the positionals after the script as a seq of strings
(or `nil`), the same whether you run a script or a bundled binary.

## Bytecode and standalone binaries

```bash
lg -c app.lgb app.lg              # compile to bytecode
lg app.lgb                        # run bytecode

lg -b myapp app.lg                # bundle into a self-contained binary
./myapp                           # runs anywhere, no lg needed
```

The standalone binary is a copy of `lg` with your bytecode appended — copy it to
another machine and it runs.

## Compiler-free deployment: lg-runtime

`lg-runtime` executes precompiled `.lgb` bytecode with no reader, compiler, or
resolver linked in. The guarantee is structural — the binary never links
`pkg/compiler` or `pkg/resolver` — so the deployed artifact has no `eval`,
`load-string`, `read-string`, or dynamic source `require`: it runs only
bytecode compiled ahead of time by a toolchain you trust.

```bash
go build -o lg-runtime ./cmd/lg-runtime      # build the runtime-only binary

lg -c app.lgb app.lg                         # compile on your machine
lg-runtime app.lgb a b                       # run it, compiler-free

lg -b myapp -bundle-base lg-runtime app.lg   # standalone AND compiler-free
./myapp
```

With `-bundle-base`, the standalone binary appends your bytecode to
`lg-runtime` instead of the full `lg`, so the shipped executable cannot
evaluate anything beyond what you compiled.

## WASM web apps

```bash
lg -w site app.lg                 # compile to a WASM web app
open site/index.html
```

The output is a self-contained `index.html` (~6MB, inlined WASM, gzipped) plus a
service worker that supplies the COOP/COEP headers GitHub Pages needs for
SharedArrayBuffer. Programs that use the `term` namespace get full terminal
emulation via xterm.js: ANSI colors, cursor positioning, raw keyboard input.

## Compile-time vars

`*compiling-aot*` is `true` during `-c`/`-b`/`-w` compilation and `false` at
runtime, useful for keeping side effects out of compile time:

```clojure
(defn -main []
  (start-server))

(when-not *compiling-aot*
  (-main))
```

`*in-wasm*` is `true` when running inside a WASM build.

## Project management with lgx

For multi-file projects with dependencies, [lgx](https://github.com/abogoyavlensky/lgx)
is a git-based package and project manager for let-go — dependency resolution,
runner, build tool, test runner, scaffolder, and task runner in one binary:

```bash
lgx new myapp                     # scaffold a project
lgx install                       # install deps from lgx.edn
lgx run                           # fetch deps, run :main
lgx build                         # bundle a standalone binary
lgx test                          # run tests under test/
```
