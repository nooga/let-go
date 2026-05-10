# Clojure Compatibility Roadmap — Execution Plan

Goal: load `hiccup` from source AND improve the
[clojure-test-suite](../test/clojure-test-suite/) score, with each step
producing a measurable delta on at least one of the two scoreboards.

This is the **execution doc** for the strategy laid out in
[jvm-compat-plan.md](jvm-compat-plan.md). That doc is the *what and
why*; this doc is the *what next*.

## Scoreboards

- **Real-world compat**: `go run . test/compat/run.lg` → histogram of
  fixable failures + raw and effective pass rates.
- **Conformance**: `go test ./test/ -count=1 -run TestClojureTestSuite -v`
  → TOTALS line + `knownFailing` whitelist in `test/zz_compat_test.go`.

After every step, re-run both, diff the deltas, and update `knownFailing`
(removing entries that now pass — the harness enforces this).

---

## Pre-flight cleanups

Five small, independent fixes. Each unblocks a handful of files; none
depends on the JVM-shim work. Knock these out first to surface cleaner
signal in the histogram before Step 1.

### P1 — reader conditionals `#?(:clj … :cljs …)`

`pkg/compiler/reader.go`. On `#?` and `#?@`, read the dispatch list
and pick the `:clj` (or `:default`) branch. Splice variant `#?@` flattens
into the surrounding form.

Test: `test/reader_conditional_test.lg` — assert reading
`#?(:cljs :a :clj :b)` returns `:b`.

Unblocks: many `.cljc` files; clears 1 syntax error in suite.

### P2 — `\uXXXX` unicode escapes in strings

`pkg/compiler/reader.go` string-literal parser. After `\u`, read 4 hex
digits and emit the rune. Surrogate pairs (`😀`) join into a
single code point.

Test: `test/string_unicode_test.lg` — assert `(count "é")` == 1.

Unblocks: 2 data.json test files.

### P3 — `(catch <ClassSymbol> e …)`

Today catch accepts only `(catch e …)`. Extend to accept any class-shaped
symbol in that position, ignoring it (or matching against `:class`-tagged
ex-info maps post Step 1). Simplest first cut: drop the symbol, treat as
catch-all.

Test: `test/catch_class_test.lg`.

Unblocks: fipp/repl.clj.

### P4 — `&env` / `&form` macro context

In macro-expansion sites these resolve to nil today. Either bind them
locally inside macro bodies (proper) or shim as no-op symbols globally
(quick). Quick wins this round.

Test: `test/macro_env_test.lg`.

Unblocks: macrovich + 2–3 others.

### P5 — Trivial stubs for `*print-level*`, `*print-length*`, `tagged-literal?`

`pkg/rt/lang.go`. Define the dyn vars with default `nil` and the
predicate as `(constantly false)` (no tagged-literal type yet).

Test: `test/print_dynvars_test.lg`.

Unblocks: fipp/edn.cljc, fipp/util.cljc.

---

## Step 1 — Class-symbol registry (highest leverage)

Dual purpose: helps hiccup AND moves the test suite. Phase A from
[jvm-compat-plan.md](jvm-compat-plan.md).

1. New `pkg/rt/jvm_shim.go` populating two namespaces on init:
   - `clojure.lang`: `Keyword`, `Symbol`, `Ratio`, `BigInt`, `IFn`,
     `LazySeq`, `Atom`, `PersistentVector`, `PersistentHashMap`,
     `PersistentHashSet` → corresponding `ValueType`.
   - `java.lang`: `String`, `Long`, `Integer`, `Double`, `Float`,
     `Boolean`, `Character`, `Throwable`, `Number`, `Object` →
     corresponding `ValueType` or sentinel.
2. Add `AnyType` sentinel to `pkg/vm/`. No real value claims it as
   `.Type()`. It is the dispatch key for `Object`-extends.
3. Auto-import `java.lang` into every namespace (Clojure's behavior —
   `Object`, `String`, `Throwable` resolve unqualified).
4. `(:import java.foo.Bar)` ns clause: refer the matching var into
   the current ns.

**Tests**:
- `test/jvm_shim_test.lg` — `(extend-protocol P clojure.lang.Keyword …)`,
  bare `Object` resolves, `(:import java.net.URI)` after Step 4 lands.
- Re-run TestClojureTestSuite — expect `ancestors` and `parents` to
  move from SKIP to RUN. Expect inline fails in `derive`/`descendants`/
  `underive` referencing `String` to pass.

**Update `knownFailing`**: remove `derive`/`descendants`/`underive`
entries that now pass, or re-comment them with what's still failing.

---

## Step 2 — Protocol fallback dispatch

`pkg/vm/protocol.go`. Extend `Protocol.Lookup`:

1. nil branch (existing).
2. Exact `target.Type()` hit (existing).
3. **NEW**: `AnyType` (`Object`) fallback.

Wire `extend-type*` (`pkg/rt/lang.go:3181`) to accept `AnyType` as a
valid type argument.

Test: `test/protocol_fallback_test.lg`.

After this: hiccup `extend-protocol ToString` clause loads in full.

---

## Step 3 — `Receiver` on built-in types

Hiccup needs `.startsWith`, `.endsWith`, `.substring` on String. Corpus
also wants `.length`, `.indexOf`, `.toLowerCase`, `.trim`, `.charAt`.

1. Add `InvokeMethod(name Symbol, args []Value) (Value, error)` to
   the relevant value types.
2. Per-type method tables delegating to existing Go fns.
3. Helper: `func registerMethods(vt ValueType, table map[Symbol]methodFn)`
   to keep registrations terse.

Method list driven by corpus grep:
```bash
grep -rEho "\.\w+" ~/.cache/let-go-compat/*/src | sort -u
```

Test: `test/string_methods_test.lg`, etc., covering each registered
method.

---

## Step 4 — `java.net` shim types (`URI`, `URLEncoder`)

1. `pkg/rt/java_net.go` — `URIType` with fields parsed from
   `net/url.URL`. `Receiver` impl exposes `.getHost`, `.getPath`,
   `.toString`, `.getScheme`, `.getAuthority`.
2. Bind `->URI` so `(URI. s)` → `(->URI s)` (compiler already rewrites
   `Class.` form).
3. `URLEncoder/encode` as a regular fn in `java.net.URLEncoder` ns.

Test: `test/java_net_shim_test.lg`.

After this: **all 9 hiccup files should pass** the compat runner.

---

## Step 5 (optional) — `deftype` with Object/protocol method bodies

Only needed if hiccup's `RawString` (`deftype` with `Object (toString …)`
block) doesn't already work via the existing record machinery.

Extend the `deftype` macro in `core.lg` to accept
`(<class-or-protocol-symbol> (method [args] body)*)` blocks and emit
`extend-type` calls after the deftype body.

Test: `test/deftype_methods_test.lg`.

---

## Done condition

- `REPOS=hiccup go run . test/compat/run.lg` reports 9/9 pass.
- TestClojureTestSuite TOTALS: ≤1 compile-skip remaining (`add_watch`,
  `remove_watch` need `agent` — out of scope here).
- `knownFailing` is shorter than today's 47 entries, with each remaining
  entry annotated with the *real* root cause (not a stale comment).

## Time estimate

| block | days |
|---|---|
| Pre-flight P1–P5 | 1 |
| Step 1 | 1 |
| Step 2 | 0.5 |
| Step 3 | 1 |
| Step 4 | 1 |
| Step 5 if needed | 0.5 |
| **Total** | **4–5 working days** |

Each step is independently shippable. Stop and re-evaluate after Step 1
once we see how much it actually moves both scores.
