---
status: active
last-verified: 2026-08-23
authoritative-for:
  - known-clojure-divergences
---

# Known Clojure divergences

This ledger records known behavioral differences between let-go and Clojure
JVM. It includes both deliberate language decisions and temporary compatibility
mismatches. Intentional entries explain why let-go does not follow Clojure;
temporary entries state the intended contract and resolution condition.

This is separate from let-go's engine-differential ledgers:
`test/parity-divergence.txt` tracks bytecode-versus-gogen output, while
`test/gogen_aot_xfail.txt` tracks known gogen fixture divergences.

## Intrinsic map traversal order

**Classification:** Intentional divergence, with related compatibility work.

### Clojure behavior

Clojure JVM's representation is observable for small maps:

- A map literal with up to eight entries is a `PersistentArrayMap` and traverses
  in insertion order.
- A map literal with nine entries is a `PersistentHashMap`.
- Adding a ninth distinct entry to a `PersistentArrayMap` with `assoc` promotes
  it to `PersistentHashMap`, which no longer preserves or guarantees insertion
  order.
- A directly constructed `(array-map ...)` remains a `PersistentArrayMap` and
  preserves insertion order even when constructed with more than eight entries.

The `PersistentArrayMap` threshold is 16 alternating key/value array slots,
which means eight map entries. Clojure's common collection chunk size of 32 is
not this promotion threshold.

### let-go decision

Traversal order for intrinsic map literals and `hash-map` is unspecified at
every size. let-go does not actively randomize traversal by default; an order
that appears stable for particular values or a particular release remains an
implementation detail.

Ordering should come from an explicitly ordered collection type. `sorted-map`
provides comparator order today; `array-map` is intended to provide insertion
order after the temporary mismatch below is resolved.

### Why not Clojure?

There are two independent reasons:

1. A generic map should not acquire an ordering contract merely because its
   current representation is small. This follows the host language Go's map
   philosophy and avoids code accidentally depending on an implementation
   detail.
2. Preserving insertion order required extra state and unfavorable access
   patterns. The order side table was deliberately removed while optimizing the
   AOT/IR pipeline in
   [PR #397](https://github.com/nooga/let-go/pull/397).

Portable code must not depend on intrinsic-map traversal. Today, sort the
entries or use `sorted-map` when order is part of the program's contract. After
the temporary mismatch is resolved, explicit `array-map` construction will also
request insertion order.

### Temporary `array-map` mismatch

let-go currently routes `array-map` through its unordered persistent-map
implementation. That is not the intended long-term contract.

The intended behavior matches Clojure:

- A directly constructed `array-map` stays array-backed and insertion ordered,
  regardless of its initial size.
- Adding a ninth distinct entry to a smaller array map promotes it to the
  intrinsic unordered map.

[Issue #763](https://github.com/nooga/let-go/issues/763) tracks the surrounding
compatibility problem, but its original acceptance criteria propose ordered
small literals and limit `array-map` ordering to eight entries. Those criteria
do not match the contract established above. The issue must be updated or
replaced by a dedicated implementation tracker before this work is considered
fully scoped.

Resolution requires preserving the intrinsic-map decision while giving explicit
`array-map` construction and growth their expected behavior.

### Shared-suite overrides

The order-independent `:lg` expectations entered the shared test suite through
[jank-lang/clojure-test-suite#927](https://github.com/jank-lang/clojure-test-suite/pull/927):

- `test/clojure/core_test/cons.cljc` has two `:lg` occurrences.
- `test/clojure/core_test/cycle.cljc` has one.
- `test/clojure/core_test/mapcat.cljc` has one.

## Character model for supplementary Unicode

**Classification:** Intentional divergence, with validation defects.

### Clojure behavior

A JVM Clojure `Character` is one UTF-16 code unit. Supplementary Unicode values
above U+FFFF are represented in strings by a high and low surrogate pair, so
one `Character` cannot hold the complete value. For example, `(char 65895)`
throws instead of producing U+10167 (`𐅧`).

### let-go decision

A let-go character represents one Unicode scalar value, using Go's rune model.
This includes supplementary values through U+10FFFF as single characters. It
does not mean one user-perceived grapheme: combining marks and multi-code-point
emoji can still contain several scalar values.

Valid character ranges are:

```text
U+0000–U+D7FF
U+E000–U+10FFFF
```

U+D800–U+DFFF contains surrogate code points reserved for UTF-16. A UTF-16
code unit whose numeric value falls in that range is a high or low surrogate
code unit, not a standalone character. A valid surrogate pair encodes another
scalar value; an isolated surrogate code unit makes an ill-formed UTF-16
sequence.

### Why not Clojure?

Clojure inherits the JVM's legacy 16-bit character type. let-go can provide a
stronger character-level Unicode model because Go and UTF-8 naturally process
supplementary scalar values without exposing UTF-16 surrogate halves.

### Temporary character-validation mismatches

`char` currently has two known validation defects:

- `Int` and `BigInt` inputs in U+D800–U+DFFF pass the outer range check even
  though the intended scalar-value contract excludes surrogate code points.
- A `BigInt` is converted with `Int64()` before validation. Values outside the
  64-bit range can therefore truncate into the accepted range; for example,
  `(int (char 18446744073709551681N))` currently produces `65`.

Resolution requires checking arbitrary-precision inputs before conversion and
rejecting the surrogate range while retaining direct support for valid
supplementary values. No dedicated issue tracks these defects yet; this entry is
the local record until an issue is created or the validation is fixed.

### Shared-suite override

The `:lg` expectation entered the shared test suite through
[jank-lang/clojure-test-suite#939](https://github.com/jank-lang/clojure-test-suite/pull/939):

- `test/clojure/core_test/char.cljc` expects `(char 65895)` to equal `𐅧`.

## Possible compatibility controls

A future feature system could make selected compatibility and diagnostic costs
opt-in. This is a possible direction, not an implemented interface or committed
design.

The likely model is mixed:

- **Runtime features** could enable strict randomized traversal for collections
  whose order is unspecified, such as intrinsic maps and sets. Ordered
  collections—including vectors, lists, `array-map`, and `sorted-map`—would
  retain their contracts.
- **Build features** could opt into heavyweight compatibility capabilities and
  their linked dependencies. JVM support through a bridge such as `gojava` is
  one possible example; it should not enlarge or complicate the default
  Go-native binary.

The default remains lean. Unspecified traversal does not require randomization,
and optional strict randomization would diagnose order dependencies rather than
change which collection types promise order.

## Inventory and maintenance

As of upstream clojure-test-suite commit `95d4a911`, the five `:lg` occurrences
in the four files named above form the two intentional divergence groups in this
ledger. No other `:lg` overrides were found.

When adding or changing a divergence:

1. Classify it as intentional or temporary.
2. For an intentional difference, explain why Clojure behavior is not adopted
   and give portable-code guidance.
3. For a temporary mismatch, state the intended contract, resolution
   condition, and any tracking reference—or say explicitly that none exists.
4. Link every affected upstream test and update this ledger in the same let-go
   change that consumes the override.
5. Remove resolved temporary entries and obsolete override references.

Audit the pinned shared-suite inventory with:

```sh
git -C test/clojure-test-suite grep -n ':lg' -- '*.cljc'
```
