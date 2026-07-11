# Tasks

> **Gate (every task):** `make generate` → `make check-generated` → `go test
> ./pkg/vm/ ./pkg/rt/ ./pkg/compiler/ ./test/ -run TestRunner` → **xsofy**
> `make -C ../xsofy LG=$(pwd)/lg test` must hold **281 tests / 2733 assertions /
> 0 fail**. No task may lower the xsofy count. xsofy pinned to nooga/xsofy
> `main` (#134). Regressions halt for root-cause before proceeding.

## 0. Baseline & prerequisites
- [x] 0.1 Land + test the refer-shadowing fix (`pkg/vm/namespace.go`,
      `test/ns/refer_shadows_core_test.lg`). *(done)*
- [x] 0.2 Capture xsofy baseline on dev lg: 281 tests / 2733 assert / 0 fail. *(done)*
- [x] 0.3 EDN audit oracle: JVM `symbol→surface` index + lg `clojure.core`
      alignment → the three buckets. *(done)*
- [ ] 0.4 Commit the audit reproduction script under `scripts/` (bb/clojure
      index + alignment) so the compat assertion is regenerable.
- [ ] 0.5 Decide the pre-existing `01_shadow_emits_warning.lg` regression: fix or
      confirm-intended, so the warn baseline is known-good.

## 1. Bucket ① — remove stdlib duplicates from `clojure.core`
> **Reality (validated):** these are NOT clean deletes. lg's bootstrap +
> clojure.string are built ON these core primitives. Per name: remove from core →
> qualify EVERY internal/test bare use → `make generate` → 5 gates (vm/rt/compiler
> + `.lg` + xsofy + jank). Split into 3 sub-buckets by dependency shape:
- [x] 1.a `difference` / `intersection` — independent `(defn)` dups in core.lg;
      removed. Qualified 5 bootstrap uses in `ir/build.lg` and 2 test files to
      `clojure.set`. All 5 gates green (jank 0-skip, xsofy 2733).
- [ ] 1.b 6 independent string names (`upper-case` `lower-case` `starts-with?`
      `ends-with?` `index-of` `last-index-of`): redundant Go `ns.Def` into core +
      independent `.lg` in `string.lg`. Remove the Go Defs; qualify internal/test
      uses; keep `clojure.string`.
- [ ] 1.c 5 load-bearing (`trim` `triml` `trimr` `includes?` `split`):
      `clojure.string` re-exports these FROM core (`(def trim core/trim)`) — the
      Go primitive lives in core by design. **Re-home** the primitive into the
      `clojure.string` namespace (or an impl ns) so core stops exposing it, then
      drop the `core/X` alias in `string.lg`. Bigger; needs a Go string-ns seam.
- [ ] 1.5 Add surface test: `clojure.core` publics contain none of the 13. Gate.

## 2. Bucket ② — `let-go.core` namespace for lg extras
- [x] 2.1 Create the `let-go.core` namespace (Go: `NameLetGoCoreNS`, `LetGoCoreNS`,
      `NewNamespace` at core-init finalization + `RegisterNS`).
- [x] 2.4 Auto-refer `let-go.core` alongside `CoreNS` at ALL 4 CoreNS-refer
      sites (`RegisterNS` + `LookupOrRegisterNS`/`DefNSBare`/`…NoLoad`).
- [x] 2.5b Extend the refer-shadow fix: `let-go.core` is a resolution baseline
      (`vm.SetLetGoCoreNamespace`, `isBaselineNS`), so an explicit `:refer :all`
      shadows it. Verified: `ys.std/gt` shadows `let-go.core/gt`.
- [x] 2.2a `gt lt ge le` moved to `let-go.core`; `> < >= <=` are now native
      primitives in `clojure.core` (were core.lg `(def > gt)` aliases); core.lg
      `pos?`/`neg?` rewritten to `>`/`<`. All 5 gates green (jank 0-skip, xsofy 2733).
- [x] 2.2b-i Go-defined clean isms moved → `let-go.core`: `now lines
      str-replace-first base64-encode base64-decode base64url-encode
      base64url-decode`. Fixups: `evalStr` test harness evals in `user` ns (not
      core); `time` macro expands to `let-go.core/now` (now is public — private
      would break cross-ns macro expansion). All 5 gates green (xsofy 2733, jank).
- [x] 2.2b-ii `.lg`-defined user-facing isms → `let-go.core` via new bootstrap
      source `pkg/rt/core/let-go/core.lg` (`spy str-join range? contains-val?`).
      Gogen-lowers to `core_go_lowered/let_go/core/core.go` (native preserved).
      Required `eval.go` fix: eagerly run let-go.core's bundle chunk at boot
      (auto-refer'd nss are never `require`d, so on-demand load never fired →
      vars stayed unbound stubs). pprint's one `str-join` use inlined. All 5
      gates green.
- [x] 2.2b-iii-parse `parse-int` → `let-go.core`; `parse-long` is the native
      primitive (was `(def parse-long parse-int)`). Green.
> **Remaining tail is gated on the annotation mechanism (design D7).** Recon
> showed every remaining item is entangled — `clojure.string` is a *veneer* over
> core's Go string primitives (`(defn upper-case [s] (core/upper-case (str s)))`,
> `(def trim core/trim)`), `close!`/`open` are load-bearing (core.lg `with-open`
> macro + xsofy `debug.lg`), and re-homing means moving hand-`ns.Def`'d Go
> closures. The clean fix is the `//lg:ns` / `//lg:name` / `//lg:private`
> annotation registrar (D7, → change `add-lg-primitive-annotations`), which makes
> re-home + privatize declarative. These items land on top of that.
- [ ] 2.2b-iii-str `str-replace` → `let-go.core` (lg-ism; Clojure uses
      `clojure.string/replace`). Blocked: core.lg:2149 ns-munging uses it. Via D7
      or rewrite the munge inline.
- [ ] 2.2b-iv IO isms `open`/`close!`/`write!` → the lg `io` ns (ions.go). Blocked:
      core.lg `with-open` macro (line 2440) + xsofy use them unqualified. Via D7.
- [x] 2.2b-v structural type predicates → new focused `let-go.types` namespace
      (visible, auto-refer'd baseline): `array? bigint? big-int?` (Go primitives,
      kept Go-only to avoid the bytecode bundle-build load-order cycle) + `range?`
      (in let-go.core, gogen-lowered). `transientable?` stays `defn-` (private).
      Baseline machinery generalized (`lgBaselineNSs`/`referBaselines`/
      `isBaselineNS`/`LgBaselineNSNames`) for N lg-baseline namespaces.
- [x] 2.6 FIX the bootstrap Go `nsMacro` (lang.go): it silently dropped :require
      (an implementation accident), so any ns loaded before core.lg's own ns macro
      (i.e. core itself) could not import deps. It now emits require/refer/alias
      forms. clojure.core imports its type predicates via a real
      `(ns core (:require [let-go.types :refer [array? bigint?]]))` — undecorated,
      no Go-side refer workaround. Portability shim repointed big-int? →
      let-go.types. All 5 gates green (jank 233, xsofy 2733).
- [ ] 2.3 (as encountered) rewrite core.lg-internal ism uses to their portable
      form so `core` source never depends on `let-go.core`.
- [ ] 2.5 Add explicit let-go.core cases to `test/ns/refer_shadows_core_test.lg`.

## 3. Bucket ③ — demote internals / interop  (via annotation mechanism, D7)
> ~82 internals in `clojure.core` publics. Wholesale privatization is exactly
> what `//lg:private` (D7) automates; do it as part of `add-lg-primitive-annotations`
> rather than hand-editing hundreds of Def sites.
- [ ] 3.1 Mark internal dyn vars + compiler helpers `^:private` / `//lg:private`
      (or move to `let-go.impl`), so they leave `ns-publics`.
- [ ] 3.2 Relocate host-class/type registrations out of the `clojure.core`
      public set (keep interop-registry reachability).
- [ ] 3.3 `make generate`. Gate + xsofy.

## 4. Force the issue — full compat guard
- [ ] 4.1 Tighten the surface test: `clojure.core` publics contain **no name
      outside `clojure ∪ bb`** (any standard ns), failing on any future stray
      lg-ism Def. Use the two-reference oracle (JVM Clojure + babashka).
- [ ] 4.2 Assert `clojure.core` leaks **no JVM abstraction** (no `proxy`/`bean`/
      `aset-*`/`java.lang.*`/`clojure.lang.*` publics).
- [ ] 4.3 (O3) Assert `let-go.core` contains exactly the intended extras.
- [ ] 4.4 Wire the surface tests into the Clojure compat suite so they run in CI.

## 5. Polyfill — converge toward the portable floor, on demand
> Scope is **consumer-driven**. jank is currently green (232 files / 5657 assert /
> 0 skip), so there is **no forced polyfill work today**; the floor is the menu.
- [x] 5.0 Author the how-to guide: `docs/guide/polyfilling-clojure-core.md`
      (two-reference test, `.lg`-first, no-leak rule, JVM-equivalence + jank gate).
- [ ] 5.1 Regenerate the floor set (`clojure.core ∩ bb` − lg), non-leaking only,
      from the two-reference oracle; keep the leaking/interop names OUT (D6).
- [ ] 5.2 When a consumer (jank skip / corpus) needs a floor name, implement it in
      the `.lg` layer per the guide — candidates: `update-vals update-keys
      partition-all partitionv partitionv-all splitv-at requiring-resolve
      with-bindings(*) with-local-vars var-set halt-when iteration volatile?
      subseq rsubseq sync newline printf replicate pvalues pcalls tagged-literal
      reader-conditional? uri?` (+ REPL `*1/*2/*3/*e`).
- [ ] 5.3 Per-name JVM-equivalence test; jank `compile=0` skips. Gate + xsofy.

## 6. Docs & archive
- [ ] 6.1 Document `let-go.core` (the lg-extras home) + the strict-`clojure.core`
      guarantee (subset ⊆ JVM, plus the polyfilled additions) in the compat notes.
- [ ] 6.2 Final full gate + xsofy; `openspec archive purify-clojure-core --yes`.
