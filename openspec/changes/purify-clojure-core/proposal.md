# Purify `clojure.core` — split lg extras into `let-go.core`

## Why

lg's `clojure.core` is not a faithful Clojure surface: it carries **207 public
names that JVM `clojure.core` (1.12) does not have**. This is the grey-area
"nooga added a bunch of his own things to lg's clojure.core" that Ingy flagged.
Two concrete problems follow:

1. **Semantic collisions with real corpora.** Names like `gt`/`le` (and the
   string helpers) exist in lg's `clojure.core` with lg semantics, and *also*
   in third-party corpora (e.g. `ys.std`/`ys.dwim`) with different semantics. A
   program that does `(:require [ys.v0 :refer :all])` expects the referred vars,
   not lg's. The **resolution** half of this was a real bug — an explicit
   `:refer :all` did not shadow the auto-refer'd `clojure.core` — and is **fixed
   separately** in `pkg/vm/namespace.go` (see Dependencies). This proposal
   addresses the **surface hygiene** half: lg's extras should not live in
   `clojure.core` at all.

2. **`clojure.core` can't be asserted against real Clojure.** Because the
   surface is polluted, the Clojure compat suite cannot assert
   `clojure.core ⊆ JVM clojure.core`. Ingy: "the refer shadowing really needs to
   be part of the clojure compat suite, to force the issue."

An EDN key-alignment of lg's `clojure.core` (691 publics) against a JVM-Clojure
symbol→surface index (679 `clojure.core` publics + `clojure.string`/`set`/`io`/…)
classifies the 207 extras into three buckets:

- **① 17 stdlib duplicates** — names that belong in another standard namespace
  and **already exist there** in lg: `starts-with? ends-with? includes? index-of
  last-index-of lower-case upper-case trim triml trimr split` (→ `clojure.string`),
  `str-join str-replace str-replace-first` (renamed `clojure.string/join|replace|
  replace-first`), `difference intersection` (→ `clojure.set`), `delete-file`
  (→ `clojure.java.io`). Pure duplication into core.
- **② ~30 lg-specific user-facing names** — no standard Clojure home:
  `gt lt ge le`, `now spy open sleep mkdir`, `base64-encode/decode base64url-*`,
  `parse-int`, `array? big-int? bigint? range? contains-val?`,
  async `go go* go-loop chan close! <! <!! >! >!!`, `new throw let-go-new`.
- **③ ~82 internals / interop** — host-class registrations (`java.lang.*`,
  `clojure.lang.*`, `->Object`…), compiler helpers (`-dt-*`, `destructure*`,
  `make-*`, `apply*`, `lazy-*`), internal dyn vars (`*emit*`, `*ir-compile*`…)
  that should not be public `clojure.core` vars at all.

## What Changes

- **Trim `clojure.core` to the JVM `clojure.core` subset.** Delete bucket ① from
  `clojure.core` (the canonical homes already carry them).
- **Introduce a `let-go.core` namespace** holding bucket ② (lg-specific
  user-facing extras). It is **auto-refer'd alongside `clojure.core`** by
  `RegisterNS`, so existing unqualified usage keeps resolving with **zero source
  migration** for user code and the corpus. Because the refer-shadowing fix has
  landed, an explicit `(:require [lib :refer :all])` correctly shadows
  `let-go.core` too.
- **Demote bucket ③** from `clojure.core` publics: mark `^:private`, or relocate
  to a non-auto-refer'd impl namespace (`let-go.impl` / existing host-registry
  ns), so the interop/compiler machinery is no longer a user-visible core var.
- **Rewrite lg's own internal uses** of the moved comparison aliases in
  `pkg/rt/core/**/*.lg` to the standard operators (`gt`→`>`, `lt`→`<`, `ge`→`>=`,
  `le`→`<=`) so `clojure.core`'s own source does not depend on `let-go.core`.
- **Add a compat test** asserting `clojure.core`'s public set is a subset of JVM
  `clojure.core` (the "force the issue" guard Ingy asked for).
- **Polyfill the gap (converge from both sides).** Purifying removes lg-isms;
  the complementary direction is to add the **standard `clojure.core` names lg is
  missing**, so `clojure.core` converges on real Clojure from both directions.
  The gap is ~130 JVM names; excluding host-interop-only ones not feasible on lg
  (`proxy`, `gen-class`/`gen-interface`, `bean`, `aset-*`, `struct`/`defstruct`,
  `*-array` variants, `construct-proxy`…), the **polyfillable** set is pure-Clojure
  and portable — e.g. `update-vals` `update-keys` `partition-all` `partitionv`
  `partitionv-all` `splitv-at` `requiring-resolve` `with-bindings`/`with-bindings*`
  `with-local-vars`/`var-set` `halt-when` `iteration` `volatile?` `subseq`/`rsubseq`
  `sync` `newline` `printf` `replicate` `pvalues`/`pcalls` `tagged-literal`
  `reader-conditional?` `uri?` `*1`/`*2`/`*3`/`*e`. These land as ordinary
  `clojure.core` definitions (mostly `.lg`), each with a compat test.

## Impact

- **Affected surfaces**: `pkg/rt/lang.go` (`ns.Def` sites + `RegisterNS`
  auto-refer wiring), `pkg/rt/core/**/*.lg` (internal alias uses), the
  `core_compiled.lgb` bundle + `core_go_lowered/` tree + `generated.sums`
  (regenerated via `make generate`).
- **Compat**: `clojure.core` becomes assertably-Clojure. `let-go.core` is the
  documented home for lg extras.
- **Migration risk**: bounded to **zero for user source** while `let-go.core` is
  auto-refer'd. The gate is the **xsofy corpus** — baseline **281 tests / 2733
  assertions / 0 fail** on the dev `lg` (with the refer fix); every increment
  must hold that.
- **Reversible per bucket**: ① and ② and ③ land as independent, xsofy-gated
  steps; any can stop without unwinding the others.

## Dependencies

- **Refer-shadowing fix** (`pkg/vm/namespace.go` `Namespace.Lookup`): explicit
  non-core refers shadow the `clojure.core`/`let-go.core` auto-refer baseline.
  Already implemented + tested (`test/ns/refer_shadows_core_test.lg`). This
  proposal's `let-go.core` auto-refer is only safe *because* of that fix.
