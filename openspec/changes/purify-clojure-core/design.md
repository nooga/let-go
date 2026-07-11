# Design — Purify `clojure.core`

## Context

`clojure.core` in lg is the namespace named `"core"` (see `NameCoreNS = "core"`
in `pkg/rt/lang.go:600`), with an alias `"clojure.core" → "core"` in
`resolveNSAlias` (`lang.go:299`). It is populated by hundreds of Go-side
`ns.Def(...)` calls (`lang.go` ~6200–6600) plus the `pkg/rt/core/**/*.lg` source
compiled into it. `RegisterNS` (`lang.go:447`) auto-refers `CoreNS` into every
new namespace, which is how unqualified `clojure.core` names resolve.

### Oracle: the portable surface is `clojure ∩ bb` across all standard namespaces

The target surface is defined by **two reference implementations**, not one:
enumerate `ns-publics` across the standard namespaces (`clojure.core` +
`clojure.string` `clojure.set` `clojure.java.io` `clojure.walk` `clojure.edn`
`clojure.data` `clojure.zip` `clojure.pprint`) in **both JVM Clojure and
babashka**, emit each as a name-set, and take set operations against lg's
`clojure.core`. Measured (2026-07-07):

- **Portable surface** = `clojure ∩ bb` (all surfaces) = **710 names**. This is
  the proven-portable target lg should aim to cover; it auto-excludes both
  JVM-only host baggage bb lacks and bb-only isms.
- **Misplaced (bucket ①)** = lg `clojure.core` names that live in another
  standard surface in *both* refs = **13**: `difference intersection`
  (`clojure.set`) and `starts-with? ends-with? includes? index-of last-index-of
  lower-case upper-case trim triml trimr split` (`clojure.string`). Already
  present in lg's own `clojure.string`/`set` → pure duplication into core.
- **True lg-isms (bucket ②)** = lg `clojure.core` names in **neither** ref, any
  surface = **193** (e.g. `gt lt ge le`, `now spy open sleep mkdir`, `base64-*`,
  `parse-int`, async `go/chan/<!`, `str-join/str-replace*`, plus ~130
  compiler/interop internals = bucket ③). These are the genuine relocate/demote
  set. **Because they are in neither reference, moving them cannot regress
  clojure/bb compat.**
- **Polyfill floor** = `(clojure.core ∩ bb.core) − lg.core`, feasibility-ordered
  = ~52 portable names lg lacks (`update-vals update-keys partition-all
  partitionv requiring-resolve with-bindings volatile? subseq rsubseq sync
  iteration newline printf tagged-literal uri? seque replicate var-set …`).

The intersection is a **floor (must-have), not a ceiling (must-not-exceed)**:
JVM-`clojure.core` names lg has that bb happens to lack (`agent commute dosync`…)
are still real Clojure and stay — they are NOT relocate candidates.

Reproduce: `scripts/` bb+clojure index emitters (task 0.4); results in
`$SCRATCH/tmp/*` during this analysis.

## Goals / Non-goals

- **Goal**: lg `clojure.core` carries no name outside `clojure ∪ bb` (isms move
  to `let-go.core`); it covers the portable floor (`clojure ∩ bb`) as far as
  feasible; a named home (`let-go.core`) for lg extras; no user-source migration.
- **Non-goal**: changing lg semantics of any function; renaming
  `clojure.string`/`set`/`io` members; a strict "no lg extras auto-refer'd" mode
  (possible later, out of scope here).

### Guard against overcorrection (evidence, not the diff alone)

The audit *classifies*; it does not *authorize* removal. lg passes the jank
`clojure-test-suite` today (**232 files / 5657 assertions / 0 skip / 0 fail**),
which is evidence the current homes serve compat. Therefore: **before relocating
or removing any name, confirm nothing depends on its current home** — the jank
suite AND xsofy stay green as hard gates, and a bare (unqualified) use in a
compat corpus blocks the move until re-homed. (During analysis, apparent bare
uses of `difference`/`now`/`gt` in jank turned out to be comments, ns names, and
`deftest` names — read the lines, don't trust grep counts.)

## Key decisions

### D1 — `let-go.core` is auto-refer'd alongside `clojure.core` (not opt-in)

Two options for how user code reaches bucket ②:

- **(A) Auto-refer `let-go.core` in `RegisterNS`** — every ns gets both
  `clojure.core` and `let-go.core` unqualified. Zero migration; existing corpus
  (xsofy) and lg's own code keep working. lg extras still shadowable by an
  explicit `:refer :all` (thanks to the landed fix).
- **(B) Opt-in** — lg extras only via `(:require [let-go.core :refer :all])`.
  Strictest hygiene, but breaks every program using `gt`/`now`/`spy`/… today.

**Decision: (A).** It delivers the surface-purity win (a program that refers only
`clojure.core` sees a clean Clojure surface) at zero migration cost. (B) can be
layered later as a `*strict-core*` toggle if desired. The pollution objection to
(A) is neutralized by the refer-shadow fix: an explicit refer now wins.

### D2 — Bucket ① is deleted, not moved

The 17 names already exist in lg's `clojure.string`/`clojure.set`/`clojure.java.io`
(verified: all 17 resolve via `ns-resolve` in those namespaces). Their presence
in `clojure.core` is pure duplication. **Delete the `clojure.core` Def/source of
each**; no new Def needed. `str-join`/`str-replace`/`str-replace-first` are lg's
own renamed aliases of `clojure.string/join|replace|replace-first`; if any lg
code depends on the short core names, repoint those uses to the qualified
`clojure.string` vars (or `let-go.core` if we choose to keep the aliases as
lg-isms — see O1).

### D3 — lg's own core source must not depend on `let-go.core`

`pkg/rt/core/**/*.lg` is compiled *into* the `core` namespace. If `gt`/`lt`/`ge`/
`le` move to `let-go.core`, any internal use of them inside core source would
create a `core → let-go.core` dependency (and a load-order/circularity hazard,
since `let-go.core` itself may want core). **Rewrite internal uses to the
standard operators** (`gt`→`>`, `lt`→`<`, `ge`→`>=`, `le`→`<=`) — they are exact
aliases, so this is mechanical and semantics-preserving. Grep for internal call
sites before deleting the aliases from the auto-refer'd surface.

### D4 — Bucket ③ demotion strategy

- Host-class/type registrations (`java.lang.*`, `clojure.lang.*`, `->Object`,
  `Boolean.`, rounding-mode enums …) are interop artifacts of `register-host-*`.
  Prefer relocating them out of the `clojure.core` *public* set (they can remain
  reachable through the interop registry) rather than deleting behavior.
- Internal dyn vars (`*emit* *ir-compile* *compiling-aot* *lg-trace* *in-wasm*
  *storage* *ansi?* *keys* *test-flag* *scope-drain-timeout-ms* *agent-errors*)
  and compiler helpers (`-dt-*`, `destructure*`, `make-*`, `apply*`, `concat*`,
  `lazy-*` internals, `*-binding!`) → `^:private` in their defining source, or an
  `let-go.impl` ns. Being private removes them from `ns-publics` and from the
  compat-subset assertion without changing runtime wiring.
- This bucket is the largest and lowest-urgency; it lands **last** and may be
  split further if it destabilizes the bundle.

### D5 — Sequencing & the xsofy gate

Land as independent, individually-revertible increments, each followed by
`make generate` + `make check-generated` + the full `.lg`/vm/rt suites + the
**xsofy corpus** (`make LG=<dev lg> test`, baseline 2733 assertions / 0 fail):

1. **① delete stdlib dups** (+ compat-subset test scoped to string/set/io names).
2. **② stand up `let-go.core`**, move the ~30 lg-isms, D3 internal rewrites, D1
   auto-refer wiring.
3. **③ demote internals/interop.**
4. **Tighten the compat test** to full `clojure.core ⊆ JVM clojure.core`.

Any step that drops an xsofy assertion halts the sequence for root-cause before
proceeding (per the 2026-05-21 "xsofy.det shadowing clojure.core/nth" lesson:
reproduce in isolation is the discriminator).

### D6 — Polyfill scope: portable names in `core`; no JVM abstraction leaks in `core`

Convergence is bidirectional: purify removes lg-isms, polyfill adds portable
standard names lg lacks. The governing constraint is **`clojure.core` must not
leak a JVM abstraction** — lg is not on the JVM, so a core var that only makes
sense in terms of Java classes/proxies/arrays is a lie about the runtime. Scope:

- **In `core` (portable, non-leaking)** — pure-Clojure semantics lg can express
  without a host type: collection/util (`update-vals` `update-keys`
  `partition-all` `partitionv` `partitionv-all` `splitv-at` `subseq` `rsubseq`
  `replicate` `halt-when` `iteration`), var/binding (`requiring-resolve`
  `with-bindings(*)` `with-local-vars` `var-set` `volatile?`), IO/print
  (`newline` `printf` `pvalues` `pcalls`), reader (`tagged-literal`
  `reader-conditional?` `uri?`), REPL (`*1` `*2` `*3` `*e`), `sync` where lg has
  the primitive. Prefer the `.lg` layer; each gets a JVM-equivalence test. (~52
  names, the `clojure∩bb` polyfill floor minus the leaking ones below.)
- **Not in `core` — but NOT prohibited** — JVM-abstraction interop (`proxy`
  `gen-class` `gen-interface` `bean` `construct-proxy` `aset-*` `*-array`
  `struct`/`defstruct`/`create-struct` `memfn` `definline`, `java.lang.*` /
  `clojure.lang.*` host-class handles, classloader/compile-path vars). Per the
  interop-can-grow direction, lg **may** implement more of these over time — but
  they live in a **dedicated host-interop namespace** (e.g. `let-go.interop` /
  the existing host-registry), **never** auto-refer'd into or defined on
  `clojure.core`. Keeping them out of `core` is the whole point of bucket ③.
- **Deferred/uncertain** — agents/refs surface (`send-via` `set-validator!`
  `get-thread-bindings` `push/pop-thread-bindings`), `stream-*` — polyfill only
  when a real consumer (jank/corpus) needs it.

**Immediate scope is bounded by the jank suite:** it is already green (232 files,
5657 assertions, 0 skip), so the polyfill has **no forced work today**. Names are
added on-demand when a corpus/consumer requires them, following the how-to guide
(task 5.x + `docs/`). The floor list is the menu, not a mandate.

The polyfill (§5) lands **after** the purify buckets so the guard (§4) is green
before we add names back toward the portable surface.

### D7 — Annotation-driven primitive export (the mechanism for the tail)

The deferred tail (re-home `clojure.string` primitives, IO isms, privatize
bucket ③) is uniformly blocked on the same thing: today a Go primitive's
**namespace, lg name, and visibility** are all decided imperatively at the
`ns.Def("name", fn)` call site (hundreds of them, `fn` usually an anonymous
`vm.NativeFnType.Wrap(func…)` closure). There is **no** comment-annotation to
control export (only `;; lgbgen:skip`/`:skip-go`, which gate bundle/lowering, and
lginterop, which routes whole *external* Go packages to a deps.edn alias). So
re-homing is manual and error-prone.

**Decision (mechanism, to build as its own change):** a `//lg:` directive set on
**named** primitive functions, plus a generator (sibling to lginterop, for
*internal* primitives) that emits the registration:

```go
//lg:ns clojure.string
//lg:name upper-case
func lgStrUpperCase(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error) { … }

//lg:ns clojure.core
//lg:private
func lgDestructureVector(…) { … }
```

- `//lg:ns <ns>` — target namespace (defaults to `clojure.core`). Re-homing a
  primitive = change one line.
- `//lg:name <lg-form-name>` — the lg-visible symbol (defaults to a munge of the
  Go name). Decouples the Go identifier from the lg name.
- `//lg:private` — registers with `isPrivate` set (drops from `ns-publics`),
  the annotation form of `defn-`. Solves bucket ③ wholesale.
- The generator produces a `zz_primitives_generated.go` registrar; `RegisterCore`
  calls it instead of the scattered `ns.Def`s. String/type primitives get
  `//lg:ns clojure.string` / `let-go.types`; internals get `//lg:private`.

This makes the tail **declarative and mechanical** — and it's the same pattern
the "bootstrap the language onto let-go" direction wants (Go seed is minimal and
self-describing; everything else is `.lg` + gogen). Prereq: convert the relevant
anonymous closures to named functions. Scope it as a separate change
(`add-lg-primitive-annotations`) since it touches the whole primitive registry.

### D8 — Coercing Fn-adapter carries `ec`; primitives are native-typed Go

Companion to D7. Today every primitive is hand-written as
`vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) { … })` — it takes the
raw `[]vm.Value`, asserts each argument's type by hand, and threads `ec`
manually. That boilerplate is why primitives are anonymous closures (blocking the
`//lg:` annotations) and why re-homing/privatizing is painful.

**Decision (mechanism, with D7):** introduce a generic **Fn-style adapter** that
makes a plain, native-typed Go function callable from the VM:

- The primitive is written over **native Go types**, e.g.
  `func(ec *vm.ExecContext, s string, n int) (string, error)` — or with no `ec`
  when it doesn't need one.
- A reflection/generics-based adapter builds the `vm.Fn` (a `NativeFn`) around it,
  performing **argument coercion at the boundary** (VM value → the Go parameter
  type) and returning a **TypeError on mismatch** — arity and type checks live in
  the adapter, once, not in every primitive.
- The **adapter carries the `ec` pointer** (and arity/coercion state), not the
  native impl. Most impls become pure Go with no VM knowledge; only those that
  genuinely need the context take `ec` as their first parameter and the adapter
  supplies it.
- The `//lg:` generator (D7) emits `adapt(nativeFn)` + the `Def` into the target
  ns with the annotated name/visibility.

Net effect: a primitive is a named, native-typed Go function with `//lg:ns` /
`//lg:name` / `//lg:private` annotations; the adapter makes it VM-callable with
coercion-or-failure; the generator wires it up. This is the minimal-Go-seed model
the project wants — Go carries only what must be native, self-described, with the
VM-boundary concerns (coercion, ec, arity) factored into one adapter. Scope with
D7 in `add-lg-primitive-annotations`.

## Open questions

- **O1**: Keep `str-join`/`str-replace`/`str-replace-first` as convenience
  aliases in `let-go.core`, or drop them entirely (force `clojure.string/join`
  etc.)? Leaning keep-in-`let-go.core` (they're widely used lg-isms).
- **O2**: `go`/`chan`/`<!`/etc. — belong in a `let-go.core` or a dedicated
  `clojure.core.async`-shaped ns? JVM Clojure puts async in a separate lib.
  Parking in `let-go.core` for now; revisit if we build a real async ns.
- **O3**: Should `let-go.core` itself be excluded from the compat-subset
  assertion surface, or asserted to contain *exactly* the intended extras (so a
  future stray `ns.Def` into core is caught)? Recommend the latter as a second
  guard.

## Risks

- **Bundle regen churn**: every increment regenerates `core_compiled.lgb` +
  `core_go_lowered/` + `generated.sums`. Use `make generate` (not mtime-based
  `make build`) and verify with `make check-generated`.
- **Hidden internal dependence on a moved name**: mitigated by grepping
  `pkg/rt/core/**` and the Go host for each moved symbol before deletion, and by
  the xsofy gate.
- **Pre-existing warn-on-shadow regression**: `01_shadow_emits_warning.lg`
  already fails in the current working copy (independent of the refer fix). It
  shares this compat area; fix it (or confirm intended) before tightening the
  compat suite, so the suite starts from a known-good warn baseline.
