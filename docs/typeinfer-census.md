---
status: active
last-verified: 2026-10-09
authoritative-for:
  - typeinfer-census
human-verified:
---

# typeinfer-census — what type inference knows about a program

[`scripts/typeinfer-census.lg`](../scripts/typeinfer-census.lg) records the
type fact on every live instruction of every function `lg compile` lowers, and
compares two such records. It answers "did this change tighten or loosen any
types?" with a count rather than a reading of generated Go. It sits beside
[`lowering-census`](lowering-census.md), which counts which functions lower at
all.

## Running it

```sh
lg scripts/typeinfer-census.lg <file.lg>...              # table per ns, totals
lg scripts/typeinfer-census.lg --edn <file.lg>...        # snapshot, one row per function
lg scripts/typeinfer-census.lg --follow <root> main.lg   # also load the namespaces main.lg requires
lg scripts/typeinfer-census.lg --diff base.edn new.edn   # compare two snapshots
```

The facts come from the type inference of the `lg` that runs the script. To
compare two let-go trees, build `lg` in each and take a snapshot with each:

```sh
lg-base scripts/typeinfer-census.lg --edn --follow . main.lg > base.edn
lg-new  scripts/typeinfer-census.lg --edn --follow . main.lg > new.edn
lg      scripts/typeinfer-census.lg --diff base.edn new.edn
```

A run costs about as much as lowering the program with `lg compile`; `--diff`
only reads the two snapshots.

## What is recorded

The census runs the package pass `lg compile` runs
(`ir.passes.pipeline/lower-all-ns-to-go-result`) with
`ir.passes.pipeline/optimize-fn` wrapped, and keeps the IR each call leaves:
every block parameter and instruction in block order, as `[op type]`. The
package pass optimizes a defn several times: once while collecting
cross-package exports, then in pass 1 and in each pass-2 iteration while its
namespace is emitted. Only the last sweep is kept, since that is the IR the
emitted Go is built from, and of it only the arities the namespace reports as
emitted: a defn that is optimized and then falls back to the VM is not counted.

Loading the files' definitional forms and lowering them print nothing to the
census output, so `--edn` can be redirected to a file. A namespace whose
lowering throws is listed as failed in the table, in the snapshot's `:failed`
and in a diff, rather than read as a smaller program.

As in `lowering-census`, the unit is one arity of a `defn` or `defn-`,
identified by namespace, name and arity (`ns/name/1`, or `1+` when variadic),
with `#2`, `#3` added if a sweep optimizes two functions of the same name.
Deftype methods and defmethod bodies are not recorded, and neither are
lambda-lifted siblings, which skip type inference. A defn the pipeline cannot
build is not recorded either; `lowering-census` lists those.

The table counts an instruction as typed when its fact is anything other than
`:unknown`, `:any` or `:bottom`.

## Comparing

`--diff` matches functions by key. A function with the same `[op type]` pairs
in another order is counted as reordered: a shifted gensym counter can swap a
`:block-arg` and an `:invalid` in a block, and even a change to the census
script itself shifts it. For a function whose op sequence is the same in both
snapshots, each changed fact is classified with
`ir.lattice/type-join`: tighter when the new fact is below the old one, looser
when above, incomparable otherwise. `:bottom` is the lattice's least element,
but a `:bottom` left after the fixpoint means typeinfer never reached the
instruction, so a fact lost to `:bottom` counts as looser and one gained from
it as tighter. A function whose ops changed (a pass rewrote it differently)
is listed with its typed and total instruction counts,
since its instructions no longer line up one to one. Functions present on one
side only are listed by key.

## Checking that it sees type inference

A diff of all zeros should mean nothing changed, not that the census missed
the change. `test/typeinfer_census_test.lg` makes
`ir.passes.typeinfer/call-core-type` answer nil and checks that the two facts
it supplies in the fixture (a `math/sqrt` result's `:float`, a `transient`
result's `:transient-vector`) come back looser. The same check on a real
program needs code that calls those functions: legmacs calls neither, so its
census is unchanged under it.

## Results, 2026-10-09

legmacs at `187fea2`, let-go at `a47a81b`, from the legmacs checkout root with
`--follow . main.lg`: 672 functions, 45,581 instructions, 6,585 typed (14.4%),
36,485 `:unknown`, 2,511 `:any`, none `:bottom`. Two runs with the same `lg`
and script gave identical snapshots; a run after editing the script reordered
3 functions and changed no fact.
