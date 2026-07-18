---
status: active
last-verified: 2026-07-17
authoritative-for:
  - ir-dynamic-vars
human-verified:
---

# IR pipeline dynamic vars

The IR compile/lowering pipeline is configured almost entirely through
`^:dynamic` vars rather than flags or options. They are scattered across
`core.lg`, `passes/pipeline.lg`, `passes/inline.lg`, `passes/fusion.lg`,
`passes/typeinfer.lg`, and `lower_go.lg`, and only one (`*strict-structured?*`)
has an environment seed — the rest are set with `binding` / `set!`. This page
is the single index.

Two kinds of var live here:

- **Knobs** — you set them (via `binding` or `set!`) to change what the
  pipeline does. Documented with defaults and semantics below.
- **Per-compile state** — `nil`/empty-initialized vars the pipeline rebinds
  as it runs. Listed at the end so they're discoverable, but they are not
  settings; binding them by hand will usually just break a compile.

For build tags, environment variables, and `lgbgen` CLI flags, see the
companion sections in the guide; this page covers the dynamic vars only.

## Compilation-mode control

The entry knobs that decide whether IR compilation runs at all and which
backend it targets.

| Var | Default | Semantics |
|---|---|---|
| `*ir-compile*` | `false` | Routes single-arity `defn`s through the IR-compile path instead of standard bytecode expansion (multi-arity and docstring-less edge cases still fall back). Enable with `(set! *ir-compile* true)` **after** `(require 'ir.passes.pipeline)` — it throws if the pipeline isn't loaded. There is no environment variable for it. `core.lg:963` |
| `*ir-compile-verbose*` | `false` | When true, the `defn` macro logs a diagnostic (name + error) each time a fn falls back to bytecode. `core.lg:968` |
| `*ir-compile-fallback-log*` | `(atom [])` | Vector of `[name error-msg]` fallback records, populated while `*ir-compile-verbose*` is true. `core.lg:973` |
| `*target*` | `:bytecode` | `:go` makes `compile-form*` route through `ir.lower-go` (native Go) instead of `ir.lower` (bytecode). Bind it before calling `compile-form`. `passes/pipeline.lg:545` |

## Pass toggles & tuning

| Var | Default | Semantics |
|---|---|---|
| `*enable-fusion*` | `true` | Transducer / deforestation fusion, placed after `cse`. On by default — measured ~20% fewer allocations across the ClojureTestSuite with backend parity and no ir-stress regression. `passes/fusion.lg:27` |
| `*enable-inline*` | `false` | Master switch for the inline pass. Opt-in: inlining supersedes the #345 direct-call path and still has rough edges (deftype-devirt codegen), so it stays off outside the AOT combinator measurement harness. `passes/inline.lg:25` |
| `*max-unroll*` | `32` | Cap on fold-over-rest unrolling (ITER-0034). A combinator call with more than this many flat rest operands is left as a runtime call with a logged skip — never silently truncated — rather than unrolled into an oversized branch chain. `passes/inline.lg:30` |
| `*typeinfer-max-drains*` | `2000000` | Backstop bound on the typeinfer fixpoint for pathological inputs (never fires on real code). The bail is sound — every assigned type is monotone and `lower-go`'s `rt.<Op>Value` path handles `:any` operands. Bind `nil` for unbounded. `passes/typeinfer.lg:494` |
| `*strict-structured?*` | `false` (seeded from `LG_STRICT_STRUCTURED`) | When true, structured-control-flow drift throws (and the caller falls the whole fn back to bytecode) instead of emitting a possibly mis-lowered `goto` body. Default off: the non-strict path stays correct via the coalesce-map interference fix; this just forbids the path. `lower_go.lg:2477` |
| `*direct-calls-disabled?*` | `false` | Forces every call through the cached-var / `InvokeValue` trampoline (which re-reads the var root each call) so runtime `alter-var-root` / `intern` overrides are observed. A baked direct call — `corefns.Count`, a lowered sibling's Go func — would otherwise ignore them. `lower_go.lg:1863` |
| `*pass-trace*` | `nil` | Bind to an atom to capture per-pass instruction traces. `passes/trace.lg:39` |

## Cross-package / exported-wrapper control (`--target=go`)

Knobs for whole-program Go lowering, where one lowered package must call into
another. Defaults keep the committed lowered tree byte-identical until the
whole-program collector binds them.

| Var | Default | Semantics |
|---|---|---|
| `*emit-exported-wrappers*` | `false` | Emit an exported thin forwarding wrapper for each direct-callable lowered fn so it is reachable from another Go package. Off keeps bootstrap codegen byte-stable; flipped on by the collector and the T3 unit test. `passes/pipeline.lg:1065` |
| `*cross-pkg-registry*` | `{}` | Whole-program `{[internal-ns name arity] -> {:go-pkg <import> :go-name "LG_<go>" …}}` of every other lowered package's direct-callable exports. Merged into the per-ns registry so a cross-package call resolves to `pkg.LG_<go>(ec, …)`. `{}` ⇒ no cross-package entries. `passes/pipeline.lg:1073` |
| `*wrapper-target-names*` | `:all` | Which fns get an exported wrapper. `:all` = every direct-callable fn (the single-ns convenience); a concrete set = exactly its members; an empty set = none. `lower-all-ns-to-go` always binds the concrete set, so a whole-program build with no cross-package references emits no dead exported API. `passes/pipeline.lg:1086` |
| `*export-name-overrides*` | `nil` | Per-namespace `{source-name -> resolved exported Go name}` bound around a namespace's collect + lower passes. PascalCase is not injective, so this remaps the loser of any collision to a distinct name. `nil` = plain PascalCase (the collision-free case). `lower_go.lg:3387` |
| `*deftype-ctor-types*` | `nil` | `{constructor-name -> deftype-name-symbol}` bound around a `:go` typeinfer pass so a call to a known constructor — `(->Square 3)` — is typed `[:dtype Square]`, carrying the concrete receiver type to field access and devirtualized dispatch. `nil` = off, zero overhead. `passes/typeinfer.lg:50` |

## Per-compile state (not knobs)

These are `nil`/empty-initialized and rebound by the pipeline as it runs.
They are listed for discoverability; setting them by hand is not a supported
configuration surface.

| Var | Location | Role |
|---|---|---|
| `*current-fn*`, `*current-inst*`, `*current-zip*` | `passes.lg:23-25` | Current traversal cursor (fn / instruction / zipper). |
| `*inline-registry*` | `passes/inline.lg:20` | Inline-candidate registry for the inline pass. |
| `*lowered-registry*` | `lower_go.lg:1342` | Registry of lowered namespaces / fns for cross-ns direct-call lowering. |
| `*native-imports-used*` | `lower_go.lg:1352` | Go imports referenced by the fn currently being emitted. |
| `*cross-ns-vars-used*` | `lower_go.lg:1371` | Cross-ns var references collected during emission (feeds the cross-package collector). |
| `*call-err-used*` | `lower_go.lg:1645` | Whether the emitted fn body needs the `callErr` plumbing. |
| `*typed-call-temps*` | `lower_go.lg:1657` | Temp bindings for typed direct calls in the current fn. |
| `*closure-arg-prefix*` | `lower_go.lg:64` | Prefix disambiguating closure-local arg names (captured-name shadowing fix). |
| `*force-needs-error*` | `lower_go.lg:2465` | Forces error plumbing on for a body regardless of inference. |
| `*deftype-ctors*` | `lower_go.lg:1906` | Deftype constructors in scope for native ctor-call emission. |
| `*protocol-methods*` | `lower_go.lg:1935` | Protocol method table for devirtualized dispatch. |
| `*protocol-method-sigs*` | `lower_go.lg:1943` | Protocol method signatures. |
| `*defmulti-dispatchers*` | `lower_go.lg:2056` | Type-dispatched `defmulti`/`defmethod` tables. |
| `*ti-counters*` | `lattice.lg:110` | Typeinfer instrumentation counters. |

## Notes

- Only `*ir-compile*` and `*target*` are ordinarily set by callers; the rest
  are pass-internal defaults that tests and the whole-program driver rebind.
- `*strict-structured?*` is the only var with an environment seed
  (`LG_STRICT_STRUCTURED`); everything else is `binding` / `set!` only.
- Line numbers are anchors, not contracts — grep the var name if a file has
  drifted since `last-verified`.
