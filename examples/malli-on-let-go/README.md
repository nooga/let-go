# malli on let-go

Runs stock [metosin/malli](https://github.com/metosin/malli) — a real,
non-trivial Clojure library — in a let-go REPL, with **zero edits to malli's
source**. All host/JVM interop is supplied by a thin compatibility shim.

```bash
./run.sh        # builds let-go once, then drops you into the REPL (Ctrl-C to quit)
```

```clojure
;; human-readable errors
(me/humanize (m/explain [:map [:a :int] [:b :string]] {:a "x"}))
;=> {:b ["missing required key"], :a ["should be an integer"]}

;; parse a sequence into tagged values
(m/parse [:catn [:x :int] [:y :string]] [1 "hi"])   ;=> #Tags{:values {:x 1, :y "hi"}}

;; coercion, schema algebra, function schemas
(m/decode [:map [:a {:default 5} :int]] {} (mt/default-value-transformer))   ;=> {:a 5}
(m/form (mu/merge [:map [:a :int]] [:map [:b :string]]))
(m/validate [:=> [:cat :int] :int] inc)             ;=> true
```

Aliases in the REPL: `m`=malli.core, `mt`=malli.transform, `mu`=malli.util,
`me`=malli.error.

## Layout

| File | Role |
|------|------|
| `lgx.edn` | dependency manifest — malli via `:local/root`, the shim under `:paths`, `:main` = `repl.lg` |
| `repl.lg` | `:main` — requires the shim, then malli, sets aliases, prints the banner |
| `src/malli_shims/prelude.lg` | host-class / JVM-stdlib shims (loaded before malli) |
| `src/malli/sci.lg` | stub for `malli.sci` (avoids the `sci`/`borkdude.dynaload` dependency) |
| `shim-manifest.edn` | documents exactly what interop surface malli needs to load |
| `run.sh` | launcher (see below) |

## Why `run.sh` instead of `lgx run`

`lgx.edn` is the real dependency manifest. But lgx's runner currently shells out
via buffered `os/sh` with a **dead stdin**, so it can't attach an interactive
REPL (its own TODO says so). The fix needs an inherited-stdio process primitive —
added as let-go's `os/exec*` (a separate PR). Until lgx's runner adopts it,
`run.sh` is the interim: it reads `lgx.edn` to build `-source-paths`, then hands
them to `lg -r`. Once lgx uses `os/exec*`, this collapses to just `lgx run`.

## Requirements

This needs a let-go build carrying the compiler/VM fixes malli depends on
(upstreamed as PRs #144–#153) — stock `lg` won't run malli yet. `run.sh` builds
that binary from this repo automatically. malli is resolved from a local
checkout via `:local/root` in `lgx.edn`; edit that path (or switch to `:git/url`
+ `:git/sha`) to point at your malli.
