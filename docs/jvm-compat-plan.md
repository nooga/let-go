# JVM Compatibility Shim — Plan

Goal: load real-world Clojure libraries (hiccup, fipp, medley, …) from
source without forking them. The compat runner (`test/compat/run.lg`)
is the success metric — pass rate on the corpus is what we track.

This is a plan, not a commitment. Read end-to-end before scheduling.

---

## Non-goals

- Running on the JVM. We are not embedding java; we are mapping
  *symbol shapes* libraries use into let-go runtime values.
- Supporting reflection (`Class/forName`, `(.getClass x)`).
- Full clojure.core parity — only what the corpus actually touches.
- `cljs/*` namespaces. `.cljc` is in scope (reader conditionals
  with `:clj` branch), `.cljs` is not.

---

## The concrete wall (from hiccup)

`hiccup/util.clj` blocks 7 of 9 hiccup files. It uses four JVM-shaped
constructs the runtime currently can't represent:

```clojure
(extend-protocol ToString
  clojure.lang.Keyword (to-str [k] (name k))
  clojure.lang.Ratio   (to-str [r] (str (float r)))
  java.net.URI         (to-str [u] ...)
  Object               (to-str [x] (str x))
  nil                  (to-str [_] ""))

(URI. s)            ; constructor sugar
(.getHost u)        ; method invocation
(.startsWith p "/")

(deftype RawString [^String s]
  Object
  (^String toString [this] s)
  (^boolean equals [this other] ...))
```

These four — class-symbol resolution, `Object`/default protocol
dispatch, `(.method obj args)` invocation, `(Class. args)` constructor
sugar — are the spine. The rest of the corpus uses subsets of the
same set.

---

## Architecture (layered)

The cleanest mental model is three layers, each useful on its own:

### Layer 1 — Class symbol → ValueType registry

A namespace `clojure.lang` (and friends) populated with vars whose
roots are let-go `ValueType`s:

| Clojure symbol            | let-go value          |
|---------------------------|-----------------------|
| `clojure.lang.Keyword`    | `KeywordType`         |
| `clojure.lang.Symbol`     | `SymbolType`          |
| `clojure.lang.Ratio`      | `RatioType` (or BigDecimal-ish) |
| `clojure.lang.PersistentVector` | `PersistentVectorType` |
| `clojure.lang.PersistentHashMap` | `PersistentMapType` |
| `clojure.lang.IPersistentCollection` | `AnyCollType` (sentinel) |
| `clojure.lang.IFn`        | `FnType`              |
| `java.lang.String`        | `StringType`          |
| `java.lang.Object`        | `AnyType` (new sentinel) |
| `java.lang.Throwable`     | `ErrorType`           |
| `java.lang.Number`        | `NumberType` (sentinel) |
| `java.lang.Long` / `Integer` | `IntType`          |
| `java.lang.Double`        | `FloatType`           |

Sentinel types (`AnyType`, `NumberType`, …) are not real types — they
are dispatch keys handled specially in `Protocol.Lookup`.

**Cost**: one shim namespace, ~50 vars, plus 1–2 sentinel types.
**Unblocks**: most `extend-protocol`/`extend-type` failures.

### Layer 2 — Protocol fallback dispatch

Extend `Protocol.Lookup`:

1. Direct hit on `target.Type()` — current behavior.
2. If miss, walk a list of registered "interface" sentinels
   (`NumberType`, `AnyCollType`) checking which the target satisfies.
3. If still miss, fall back to `AnyType` impl if registered.
4. nil — current `nilImpl` path.

Sentinel satisfaction is just a Go predicate per sentinel
(`func(Value) bool`). Cheap.

**Cost**: ~50 lines in `pkg/vm/protocol.go` plus sentinel registry.
**Unblocks**: `Object`-fallback `extend-protocol` clauses.

### Layer 3 — `(.method obj args)` and `(Class. args)`

**Mostly already built.** The compiler already rewrites:

- `(.member instance args)` → `(. instance 'member args)` →
  `instance.InvokeMethod(member, args)` via the `Receiver` interface
  (`pkg/compiler/compiler.go:450`, `pkg/rt/lang.go:2000`).
- `(Class. args)` → `(->Class args)` — works for any defrecord/deftype
  (`pkg/compiler/compiler.go:441`).

What's missing is just *populating* these mechanisms for the types
the corpus reaches into:

1. Built-in value types (`String`, `PersistentVector`, `PersistentMap`,
   …) currently do **not** implement `Receiver`. Today `(.startsWith
   "abc" "a")` errors with "method-invoke expected Receiver". Make
   each implement `InvokeMethod` with a small per-type method table
   delegating to existing builtins (`String.startsWith` →
   `strings.HasPrefix`, `.length` → len, etc.).
2. Shim "Java types" (`URI`, `URLEncoder`) need `->URI` style
   constructors. Easiest path: a tiny defrecord-shaped wrapper in a
   `java.net` shim ns that holds a parsed value and exposes
   `.getHost`, `.getPath`, etc. via `InvokeMethod`.

**Cost**: ~50 lines per built-in type for the method table; a Go-side
registration helper to keep them concise. ~30 corpus methods total.
**Unblocks**: nearly all "interop" usage in the corpus.

Reflection is not on the table — the registry is hand-curated.

### (Layer 4 — `deftype` body Java methods)

Out of scope for first pass. Workaround: `deftype` already exists in
let-go — extend its parser to accept (and ignore) `Object (^String
toString …)` style blocks, treating them as protocol-style impls
against the Object sentinel.

---

## Phasing

Each phase is shippable on its own and gives a measurable corpus
delta.

**Phase A** — Layer 1 (class registry).
- Add `pkg/rt/jvm_shim.go` defining `clojure.lang` + `java.lang` +
  `java.net`/`java.io` shim namespaces with class→ValueType bindings.
- Mark non-existing types with a stub sentinel that errors helpfully
  when used.
- Run corpus → expect platform→fixable shifts and a few new passes.

**Phase B** — Layer 2 (protocol fallback).
- Add `AnyType` + `NumberType` sentinels.
- Extend `Protocol.Lookup` walk.
- Wire sentinel resolution in `extend-type*`.
- Run corpus → expect hiccup/util.clj to load if Layer 3 isn't needed
  for it (it is — the URI body needs `.getHost`).

**Phase C** — Layer 3 (populate Receiver for built-in types).
- Reader/compiler dispatch already exists — only need to add
  `InvokeMethod` to `String`, the persistent collection types, and
  whatever else the corpus invokes methods on.
- Hand-register the ~30 methods the corpus actually uses.
- Add `java.net`/`java.io` shim types (defrecord-shaped) for the
  classes that show up in `(Class. args)` form.
- Run corpus → hiccup loads end-to-end.

Stop after C and re-evaluate. Layer 4 (deftype Java methods) only if
the remaining corpus failures still pin on it.

---

## Test strategy

- `test/compat/run.lg` is the regression suite. After each phase,
  diff the histogram before/after.
- Per-phase: a focused `.lg` integration test under `test/` exercising
  the new mechanism without third-party code (e.g., `test/jvm_shim_test.lg`
  with `(extend-protocol ToString clojure.lang.Keyword …)` and
  `(.startsWith "abc" "a")`).
- No mocks. Real types, real protocols.
- Bench: load-string of a representative file (medley/core.cljc)
  before/after — must not regress meaningfully.

---

## Open decisions (need answers before building)

1. **Sentinel granularity**. `AnyType` is enough for `Object`. Do we
   also need `NumberType`, `AnyCollType`, `SeqableType`? Probably yes
   for `extend-protocol` to types like `clojure.lang.Sequential`.
   Decision: start minimal, add as corpus demands.

2. **Constructor namespace**. `(URI. s)` — where does `URI` resolve?
   Option (a): a `java.net` shim ns with `URI` bound to a let-go
   record-like type with a factory. Option (b): a global registry
   keyed on the symbol. Decision: (a) — it's discoverable and
   composes with refer/import.

3. **`(:import …)` ns clause**. Currently silently ignored. Should
   it instead refer the imported class symbols into the current ns
   so `URI` (without prefix) resolves? Yes, post Phase A. Cheap.

4. **Reflection ducking**. Some methods like `.getClass`, `.hashCode`
   are everywhere. Mapping them is mostly trivial (`type` and `hash`).
   Decision: include in Phase C's hand-registered set.

5. **What about `gen-class`, `proxy`, `reify`**? `reify` is the only
   one widely used in the corpus. Out of scope for the first pass.
   Document as a known gap.

---

## What we explicitly will not do

- Auto-generate the class registry from JVM reflection. Hand-curated
  is faster and more honest about what we actually support.
- Try to make let-go a Clojure replacement. The bar is "load and run
  pure-Clojure libraries", not "be Clojure".
- Support multi-arity Java method overloading. Pick the most common
  signature; reject the rest.

---

## Estimate

Phase A: 1 day. Phase B: 1 day. Phase C: 1–2 days (smaller than the
original estimate — the compiler-side dispatch is already built;
only the receiver tables need writing). Layer 4 if needed: 2 days.
Total: ~3–5 working days to land the spine, plus ongoing work each
time the corpus surfaces a new method/class.
