---
status: active
last-verified: 2026-07-07
human-verified:
---

# Polyfilling `clojure.core`

How to add a `clojure.core` name that lg is missing, without leaking JVM
abstractions into the core namespace.

## The rule

`clojure.core` is a **portable, host-neutral** surface. lg does not run on the
JVM, so a core var whose meaning depends on a Java class, proxy, or primitive
array is an abstraction leak — a lie about the runtime. Two consequences:

1. **In `core`**: only names that exist in the portable surface and can be
   expressed without a host type.
2. **Interop is allowed, but not in `core`**: JVM-flavoured interop (`proxy`,
   `bean`, `aset-*`, class handles) may be implemented over time, but it lives in
   a dedicated host-interop namespace, never `clojure.core`.

## Is the name in scope? (the two-reference test)

The target surface is defined by **two** reference Clojures, not one: JVM Clojure
and babashka. A name belongs in `clojure.core` only if it is in the **portable
floor** — present in `clojure.core` of *both* references — and is non-leaking.

```bash
# Portable floor: clojure.core ∩ bb.core
clojure -M -e '(run! println (sort (map name (keys (ns-publics (quote clojure.core))))))' > /tmp/jvm.txt
bb -e '(run! println (sort (map name (keys (ns-publics (quote clojure.core))))))'          > /tmp/bb.txt
comm -12 <(LC_ALL=C sort -u /tmp/jvm.txt) <(LC_ALL=C sort -u /tmp/bb.txt)   # the floor
```

Decision table for a name `X`:

| `X` is… | Action |
|---|---|
| in the portable floor, non-leaking | **polyfill into `clojure.core`** (this guide) |
| in another standard ns (`clojure.string`/`set`/`io`) | put it there, not `core` |
| JVM-only / host-interop (`proxy`, `bean`, `aset-*`, `java.lang.*`) | a host-interop ns, **not** `core` — never auto-refer'd |
| in neither reference (an lg-ism) | `let-go.core`, not `clojure.core` |

Do not add speculatively. Additions are **driven by a consumer** — a jank
`clojure-test-suite` compile-skip, or a real corpus need. If the jank suite is
green (no missing-builtin skips), there is no forced work.

## Where it goes

Prefer the **`.lg` layer** (`pkg/rt/core/**/*.lg`) over Go — most portable core
fns are a few lines of Clojure. Reach for Go only when the fn needs a VM
primitive lg exposes natively (e.g. a new reduce fast-path).

- Most collection/seq/util fns → `pkg/rt/core/core.lg` (or the topical file).
- Anything that would need a host type → **stop**; it is not a `core` polyfill.

After editing any `pkg/rt/core/**/*.lg`, regenerate the artifacts the runtime
actually loads:

```bash
make generate        # rebuilds core_compiled.lgb + core_go_lowered/ + generated.sums
make check-generated # content-based freshness check (checkout-safe)
```

`go build` alone will NOT regenerate — the runtime loads the bundle, not the
embedded `.lg` strings, so skipping this silently no-ops your change.

## Writing the polyfill

Implement the JVM-Clojure semantics exactly. Cross-check against the reference:

```clojure
;; pkg/rt/core/core.lg
(defn update-vals
  "m f -> map with f applied to every val. Matches clojure.core/update-vals."
  [m f]
  (persistent!
    (reduce-kv (fn [acc k v] (assoc! acc k (f v)))
               (transient (empty m)) m)))
```

Match arities, laziness, and edge cases (empty input, `nil`, transient vs
persistent return type) to Clojure. When unsure, diff behavior against `clojure`
and `bb` directly.

## Testing

1. **JVM-equivalence test** — a `.lg` test under `test/` (run by `TestRunner`)
   asserting lg's result equals Clojure's for representative inputs:

   ```clojure
   ;; test/polyfill_update_vals_test.lg
   (deftest update-vals-matches-clojure
     (is (= {:a 2 :b 3} (update-vals {:a 1 :b 2} inc)))
     (is (= {} (update-vals {} inc))))
   ```

2. **jank suite** — the authoritative consumer. From the main checkout (the
   suite is a submodule; jj workspaces don't carry it):

   ```bash
   git submodule update --init   # once
   go test ./test/ -run TestClojureTestSuite -v -count=1
   # want: TOTALS ... skipped: compile=0  (0 missing-builtin skips)
   ```

3. **Regression gates** — keep them green:

   ```bash
   go test ./pkg/vm/ ./pkg/rt/ ./pkg/compiler/ ./test/ -run TestRunner
   make -C ../xsofy LG=$(pwd)/lg test    # external corpus
   ```

## Interop that is *not* a core polyfill

If the name you want is host-interop (`proxy`, `bean`, `aset-int`, a
`java.lang.*` handle): it does **not** go in `clojure.core`. lg may grow this
surface, but in a dedicated host-interop namespace that is not auto-refer'd into
user namespaces. Keeping JVM abstractions out of `core` is deliberate — see
`openspec/changes/purify-clojure-core/design.md` (D6).

## Checklist

- [ ] Name is in the portable floor (`clojure.core ∩ bb`) and non-leaking.
- [ ] There is a real consumer (jank skip / corpus need).
- [ ] Implemented in `.lg` (Go only if a VM primitive is required).
- [ ] `make generate` + `make check-generated` run.
- [ ] JVM-equivalence `.lg` test added.
- [ ] jank suite: 0 missing-builtin compile-skips.
- [ ] vm/rt/compiler/`TestRunner` + xsofy stay green.
