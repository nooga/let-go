---
status: active
last-verified: 2026-06-01
authoritative-for:
  - xsofy-side-clojure-compat-evidence
  - intentional-jvm-divergences
shipped:
  - A1 (class-less catch — pkg/compiler/compiler.go:799-820)
  - A3 (int-array / byte-array / long-array constructors — pkg/rt/lang.go:6738-6760)
  - D2 (hex BigInt promotion — #70, in v1.8.0+)
  - D3 (unchecked-* coercion — #66, in v1.8.0+)
  - D4 (destructure-shadow-in-loop bug — #68, in v1.8.0+)
  - D5 (clojure.core shadow warnings — #69, in v1.8.0+)
remaining-open:
  - A2 aset polymorphism (recommendation in doc: keep behavior, document as intentional divergence)
---

# let-go ↔ Clojure JVM portability gaps surfaced by xsofy

Real-world portability findings from attempting to load
[xsofy](https://github.com/nooga/xsofy) (a roguelike, ~25 .lg files,
~5K LoC) onto Clojure 1.12.5 JVM after the 2026-05-22 let-go cleanups
(#66 unchecked coercion, #68 destructure-shadow, #69 core-shadow
warnings, #70 hex BigInt promotion). Every gap below was hit while
trying to `(require 'xsofy.world)` on a stock JVM Clojure with minimal
shimming.

Companion to [jvm-compat-plan.md](jvm-compat-plan.md) and
[clojure-compat-roadmap.md](clojure-compat-roadmap.md) — this doc is
the *concrete user-side evidence* of what needs work, sorted by what
let-go can fix vs. what's irreducibly let-go-specific.

## Scoreboard

22 xsofy modules, after applying fixes below:

| State | Count | Notes |
|---|---|---|
| Loads cleanly | 22/22 | After fixes below |
| Tests pass on JVM | 60/61 hash tests | 1 documented semantic divergence (see D1) |
| Headless world-gen works | Blocked | aset on byte-array (see A2) |

## A. let-go is permissive; Clojure rejects (let-go could fix)

Patterns let-go accepts that Clojure JVM doesn't. Aligning let-go to
Clojure's stricter rules would let xsofy code be portable as written.

### A1. Class-less `(catch e ...)` and `(catch _ ...)`

**Current let-go:** `(try ... (catch e body))` works as a catch-all.
**Clojure JVM:** parses `e` as a class name, fails with "Unable to
resolve classname: e".

**xsofy impact:** 7 sites (check.lg×4, check_gen.lg×3). All
straightforward catch-all error handling.

**Proposal:** parse `(catch x body)` such that if `x` resolves to a
class, use it; otherwise treat as catch-all and bind `x` to the
exception. Matches the spirit of let-go's looser typing while still
accepting Clojure-style class-typed catches.

Alternative: align strictly to Clojure (require a class name). Worse
for ergonomics but simpler to implement.

### A2. `(aset byte-array i int-literal)` polymorphism

**Current let-go:** `(aset some-byte-array 5 8)` works regardless of
the element type — `8` is implicitly narrowed to a byte.
**Clojure JVM:** `aset` reflects on the array type. On a `byte[]` it
requires the value be a `Byte`, not `Integer`. Must use
`(aset-byte arr i 8)` or `(aset arr i (byte 8))`.

**xsofy impact:** 31 sites in terrain.lg (all `(aset grid i N)` where
`grid` is a `byte-array`).

**Proposal:** let-go's behavior is arguably *better* (less type tax,
matches what users want from a Lisp). Worth keeping but documenting
as an intentional divergence so portability-conscious code can use
`aset-byte` / `(aset arr i (byte n))` explicitly.

If let-go does want strict parity: type-check `aset` against the array
element type and reject mismatched primitives. Costly.

### A3. `(ints xs)` and `(bytes xs)` as constructors

**Current let-go:** `(ints [1 -1 0 0])` constructs a primitive `int[]`
from a vector.
**Clojure JVM:** `ints` is a type *coercion* (hints `Object → int[]`)
— it doesn't construct. Use `(int-array [1 -1 0 0])` and
`(byte-array [1 2 3])`.

**xsofy impact:** 4 sites in terrain.lg.

**Proposal:** expose `(int-array xs)` / `(byte-array xs)` / `(long-array
xs)` as primitive constructors with Clojure-parity semantics. Keep
`(ints xs)` working for backward compat but document the portable
form.

This is the cheapest fix on this list — purely additive, no semantic
change required.

## B. let-go-only built-in namespaces

let-go ships built-in namespaces that have no JVM equivalent. Each
needs a shim on the JVM side.

| Namespace | xsofy surface used | JVM equivalent |
|---|---|---|
| `math` | `sqrt`, `abs`, `exp`, `pow`, `log`, `sin`, `cos`, `floor`, `ceil`, `round`, `min`, `max`, `pi`, `e` | `java.lang.Math/*` static methods |
| `io` | `open`, `write!`, `flush!`, `close!` | `java.io.BufferedWriter`/`FileWriter` |
| `test` | `deftest`, `is`, `testing`, `run-tests` | `clojure.test/*` |
| `term` | `move-cursor`, `set-fg`, `set-bg`, `write`, `clear`, `flush`, `read-key`, `size` | jline / lanterna |
| `async` | `timeout`, `<!!`, `go` | `core.async` (overlapping API) |

**Not a let-go bug**, but worth a docs section calling out:

1. Which built-in namespaces are let-go-specific.
2. For ones with direct Clojure equivalents (notably `test` →
   `clojure.test`), consider:
   - **(a)** auto-aliasing — `(require '[test])` on JVM transparently
     resolves to `clojure.test`. Letting people write code that runs
     on both with no shims.
   - **(b)** publishing a tiny `let-go-shims` jar for JVM users that
     provides `math`, `io`, `test`, etc. as thin aliases.

Option (a) is more user-friendly but only works for the subset where
the APIs match closely. `test` is the obvious candidate.

## C. let-go-injected globals

### C1. `*in-wasm*` dynamic var

**Current let-go:** auto-interned at startup; user code can reference
it for runtime-environment switching.
**Clojure JVM:** undefined symbol unless explicitly created.

**xsofy impact:** 4 sites (debug.lg, render.lg, title.lg×2).

**Proposal:** document that `*in-wasm*` is let-go-specific and that
portable code should `(def ^:dynamic *in-wasm* false)` at the top of
namespaces that use it. Alternative: provide a portable
`environment` ns (e.g., `(env/wasm?)` predicate) that's defined on
both runtimes.

## D. Known semantic divergences (probably stay as-is)

Differences where alignment would hurt let-go more than help, or
that are documentation-only.

### D1. `(bit-shift-left x 64)`

**let-go:** returns 0 (full mod-2^64 wrap on int64).
**Clojure JVM:** returns `x` (Java `<<` masks shift count by `& 63`,
so shift-by-64 is shift-by-0).

**xsofy impact:** xsofy/hash.lg's `u64-shl` test
(`u64-shl-wraps-at-2-to-64`) relies on let-go's behavior. One test
fails on JVM. Not load-bearing for xxh3 correctness — the hash
function itself never shifts by exactly 64 in practice.

**Proposal:** keep let-go's behavior. It matches u64 arithmetic
intuitions better than Java's quirk. Document under "intentional
divergences from Clojure JVM" alongside the other primitive-int
choices.

### D2 – D5. Already fixed

- **D2 hex literal BigInt promotion** — fixed by #70 (2026-05-22).
- **D3 unchecked-* coercion family** — fixed by #66 (2026-05-22).
- **D4 destructure-shadow-in-loop bug** — fixed by #68 (2026-05-22).
- **D5 clojure.core shadow warnings** — added by #69 (2026-05-22).

All four merged but not yet in a tagged let-go release.

## Recommended action sequence

In rough priority order (highest leverage first):

1. **A3** — add `int-array`, `byte-array`, `long-array` constructors.
   Purely additive, mechanical, unblocks one whole class of
   portability complaint. ~1 hour.

2. **A1** — class-less catch. Two-line reader/compiler change to
   accept `(catch sym body)` where `sym` isn't a known class. ~2
   hours.

3. **Docs** — write a "let-go vs. Clojure JVM" page in `docs/`
   covering Categories B/C/D. Lists every let-go-specific surface
   xsofy hit, what to do on the JVM side, and which divergences are
   intentional. ~2 hours.

4. **B/test alias** — when `(require '[test])` is encountered on a
   build that has `clojure.test` available, alias it. Lowest priority;
   easiest workaround already exists (a 5-line `test.clj` shim).

5. **A2** — defer or document as intentional. Cost-benefit isn't
   favorable to fix.

Total: <1 day of focused work to dramatically improve let-go's
"can I run my Clojure code on this?" story.

## Appendix: shim files used to run xsofy headlessly on JVM

For reference, the JVM-side shims I wrote to get 21/22 xsofy modules
loading on stock Clojure 1.12.5. These are what users would need to
write if let-go doesn't grow portability features.

- `src/math.clj` — 12 functions aliasing `java.lang.Math/*`.
- `src/io.clj` — 4 functions wrapping `java.io.BufferedWriter`.
- `src/test.clj` — 4 macros aliasing `clojure.test/*`.
- `src/term.clj` — 13 no-op stubs (real impl needs jline/lanterna).
- `src/async.clj` — 3 no-op stubs (real impl: core.async).
- `src/user.clj` — interns `*in-wasm*` into `clojure.core` as dynamic.

Total shim code: ~60 lines. Plus mechanical source patches (catch×7,
ints→int-array×4, aset→aset-byte×31) on the let-go side.
