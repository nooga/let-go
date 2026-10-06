---
status: active
last-verified: 2026-10-05
authoritative-for:
  - lowering-census
human-verified:
---

# lowering-census — how much of a program the Go lowerer emits

[`scripts/lowering-census.lg`](../scripts/lowering-census.lg) reports, for a
program, how many functions `ir.lower-go` lowers to Go and why the rest do not.
`lg compile` emits Go for every function the lowerer accepts and a VM bridge for
every one it does not, per function, and prints only the output path. The census
makes that split a number. It never edits source and exits 0 unless `--gate` is
given.

## Running it

```sh
lg scripts/lowering-census.lg <file.lg>...             # table per ns, totals, top reasons
lg scripts/lowering-census.lg --edn <file.lg>...       # machine-readable, one row per unit
lg scripts/lowering-census.lg --by-reason <file.lg>... # every fallback, grouped by reason
lg scripts/lowering-census.lg --gate 95 <file.lg>...   # exit 1 when the total is below 95%
lg scripts/lowering-census.lg --follow <root> main.lg  # also load the namespaces main.lg requires
```

It runs on a normally built `lg`; `-tags gogen_ir` is not needed, because the
driver (`lg.compiler`) is embedded in every binary. Like `lg compile`, it
evaluates each file's definitional forms (`def`, `defn`, `defmacro`, `require`,
and so on) so cross-namespace names resolve, and never runs other top-level
forms. A file that fails to read or evaluate stops the run with exit 2.

Lowering a whole program to a fixpoint is slow, and the time depends more on
the program than on its size. As of 2026-10-05, legmacs (672 units) took about
45 seconds and xsofy (547 units) about six minutes.

## What is counted

The unit is one arity of a `defn` or `defn-`, as the pipeline lowers it. A
`defmulti` or `defmethod` lowers through the multimethod dispatcher, not as a
function, so it is not scored; neither are `def`, `defmacro`, `deftype` or
`defprotocol`.

Whether a unit lowered is read from the package pass `lg compile` itself runs
(`ir.passes.pipeline/lower-all-ns-to-go-result`). That pass reports only what it
emitted, so the reason a unit fell back comes from lowering the form alone
(`pipeline/compile-form`, which returns `{:status :fallback :reason ...}` in
bridge mode). A unit that lowers alone but is missing from the package output is
reported as "not emitted by the package pass" and counted under
`:disagreements`. Nothing reads generated Go; `lg.compiler/count-lowered-funcs`
approximates this count with a regex over Go text, and the census is the exact
version of it.

A defn the pipeline cannot build or optimize is dropped with all its arities, so
a multi-arity defn whose one arity throws shows every arity as a fallback.

Reasons are grouped with the shared `ir.coverage/classify-error` taxonomy that
`scripts/ir-stress.lg` also buckets by, then folded further: every unresolved
symbol is one cause, and numbers in an uncategorised message become `N`. The raw
message stays on each row as `:detail` in the `--edn` output.

`scripts/ir-stress.lg lower-go` measures a corpus of definitions one at a time,
with no program context, and ratchets the failure count. The census measures a
program the way `lg compile` would build it and names the namespace, function
and arity of each fallback.

## Results, 2026-10-05

Run from each checkout root with `--follow .`, on `main.lg` and the namespaces it
requires under that root; let-go at `dbaf59cb`, xsofy at `0c38395`, legmacs at
`187fea2`.

| program | namespaces | lowered | fallback | percent |
|---|---|---|---|---|
| xsofy | 34 | 545 | 2 | 99.6% |
| legmacs | 29 | 672 | 0 | 100.0% |

xsofy's two fallbacks are both optimizer validation failures, not unsupported
shapes: `xsofy.main/read-dismiss-key!/0` (`validate after licm`: a branch passes
a different number of arguments than its target has parameters) and
`xsofy.world/bfs-path/2` (`validate/cross-block-ref`).

"Lowered" means the function was emitted as Go, not that its body is free of
VM calls or typed throughout.
