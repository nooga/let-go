---
status: active
last-verified: 2026-07-19
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

Every knob here changes what the pipeline converts to IR (or how it lowers
it), so a change to one is only as good as its coverage evidence. Measure it
with the ir-stress harness — see [Verifying a var change](#verifying-a-var-change)
at the end of this page and `scripts/ir-stress.md` for the harness itself.

## Compilation-mode control

The entry knobs that decide whether IR compilation runs at all and which
backend it targets.

| Var | Default | Semantics |
|---|---|---|
| `*ir-compile*` | `false` | Routes single-arity `defn`s through the IR-compile path instead of standard bytecode expansion (multi-arity and docstring-less edge cases still fall back). Enable with `(set! *ir-compile* true)` **after** `(require 'ir.passes.pipeline)` — it throws if the pipeline isn't loaded. There is no environment variable for it. `core.lg:963` |
| `*ir-compile-verbose*` | `false` | When true, the `defn` macro logs a diagnostic (name + error) each time a fn falls back to bytecode. `core.lg:968` |
| `*ir-compile-fallback-log*` | `(atom [])` | Vector of `[name error-msg]` fallback records, populated while `*ir-compile-verbose*` is true. `core.lg:973` |
| `*target*` | `:bytecode` | `:go` makes `compile-form*` route through `ir.lower-go` (native Go) instead of `ir.lower` (bytecode). Bind it before calling `compile-form`. `passes/pipeline.lg:545` |

**Cost of `*ir-compile*`** — it is the one knob with a real amortization
story, and turning it on is not free. Pipeline load plus per-fn compile is a
one-time cost paid up front, so it pays back only when amortized *and* only on
allocation-bound work:

- **Alloc-heavy workloads** (persistent-map, transducers, `reduce`): ~5–12×
  fewer allocations, converting to roughly −13% wall-clock per run once heap
  churn dominates — break-even around 64 runs on persistent-map.
- **Compute-bound code** (`fib`, `loop`/`recur`): a small net loss that never
  breaks even; the compile cost has nothing to amortize against.
- Most of the win is the pipeline cutting per-element boxing, not fusion:
  `reduce` over `range` drops ~78% of allocations with no fusion applicable.

Figures from @mparrett's measurement on #555, via a gctrace GC-cycle proxy on
a single machine — directional, not a contract. `BenchmarkIRCompile`
(`pkg/ir/ir_compile_bench_test.go`) is the instrument to re-measure with, and
`make ir-stress` is what tells you whether a change moved *coverage* rather
than just cost.

## Pass toggles & tuning

| Var | Default | Semantics |
|---|---|---|
| `*enable-fusion*` | `true` | Transducer / deforestation fusion, placed after `cse`. On by default — measured ~20% fewer allocations across the ClojureTestSuite with backend parity and no ir-stress regression (`make ir-stress-gate`). `passes/fusion.lg:27` |
| `*enable-inline*` | `false` | Master switch for the inline pass. Opt-in: inlining supersedes the #345 direct-call path and still has rough edges (deftype-devirt codegen), so it stays off outside the AOT combinator measurement harness. Flipping it on is a coverage change — gate it with `make ir-stress-gate` and `make parity-full`. `passes/inline.lg:25` |
| `*max-unroll*` | `32` | Cap on fold-over-rest unrolling (ITER-0034). A combinator call with more than this many flat rest operands is left as a runtime call with a logged skip — never silently truncated — rather than unrolled into an oversized branch chain. Raising it trades compile time for code size; watch `:stress/timeout` buckets via `make ir-stress`. `passes/inline.lg:30` |
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

## Verifying a var change

The vars above decide which forms convert to IR and how they lower, so
flipping a default (or adding a knob) is a coverage change until proven
otherwise. `scripts/ir-stress.lg` is the harness that measures it: it drives a
corpus of `.lg` sources through a chosen IR path and reports per-defn buckets
(`:ok`, `:missing-form/set!`, `:validate/no-term`, `:stress/timeout`, …) so a
conversion regression shows up as a bucket that moved, not as a vague slowdown.
Full harness reference — modes, env vars, bucket meanings — is in
`scripts/ir-stress.md`.

The checks, cheapest first:

| Command | What it answers |
|---|---|
| `make ir-stress` | Native-lowering pass rate over the committed corpus allow-list (`scripts/ir-stress-corpus.edn`). The everyday "did my knob drop coverage?" run. |
| `make ir-stress-gate` | Same census, ratcheted against `docs/perf/ir-stress-baseline.edn` — exits non-zero if native-lowering failures grew. This is the gate; the ratchet only tightens. |
| `make jank-stress` | Coverage over the vendored jank Clojure-compat suite, which reaches Clojure surface the internal corpus doesn't (BigDecimal literals, multimethods). |
| `make parity-full` | Both ir-stress modes (`lower-go` and `ir-compile`) plus clojure-test-suite, run once untagged and once under `-tags gogen_ir`, comparing bucket-by-symbol. Catches a var whose effect differs between the two engines. |

Which mode maps to which claim:

- **`lower-go`** — AOT conversion coverage. Use for any var read by
  `lower_go.lg` or bound for `*target* :go` (`*emit-exported-wrappers*`,
  `*cross-pkg-registry*`, `*wrapper-target-names*`, `*export-name-overrides*`,
  `*deftype-ctor-types*`, `*strict-structured?*`, `*direct-calls-disabled?*`).
- **`ir-compile`** — eval-mode conversion, i.e. what users hit at load time
  under `*ir-compile*`. Slower, since it evals each defn.
- **`trace`** — one defn, per-pass timings. Reach for it when a tuning knob
  (`*typeinfer-max-drains*`, `*max-unroll*`) produces `:stress/timeout` and you
  need to know which pass is responsible before tuning blindly.

Practical recipe for a default flip: capture a census on both settings and diff
the bucket tallies rather than eyeballing pass counts.

```sh
LG_STRESS_LOG=/tmp/before.log make ir-stress
# flip the var's default (or bind it in the harness), rebuild
LG_STRESS_LOG=/tmp/after.log  make ir-stress
# bucket-level diff: any row that moved is a conversion change
diff <(cut -f3 /tmp/before.log | sort | uniq -c) \
     <(cut -f3 /tmp/after.log  | sort | uniq -c)
```

If the change intentionally moves coverage, rebaseline the ratchet with
`make ir-stress-rebaseline` (tool-maintained — never hand-edit the EDN), review
the diff, and commit it alongside the change that caused it. This is the
evidence `*enable-fusion*` cites above ("no ir-stress regression") and the
standard any new default should meet.

## Notes

- Only `*ir-compile*` and `*target*` are ordinarily set by callers; the rest
  are pass-internal defaults that tests and the whole-program driver rebind.
- Coverage claims about any of these vars should cite an ir-stress run — see
  [Verifying a var change](#verifying-a-var-change). `make ir-stress-gate` is
  the ratcheted form of that check.
- `*strict-structured?*` is the only var with an environment seed
  (`LG_STRICT_STRUCTURED`); everything else is `binding` / `set!` only.
- Line numbers are anchors, not contracts — grep the var name if a file has
  drifted since `last-verified`.
