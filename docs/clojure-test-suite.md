---
status: active
last-verified: 2026-07-05
authoritative-for:
  - clojure-test-suite-workflow
human-verified: 2026-06-07
---

# Clojure Test Suite — Workflow Guide

## Overview

The [jank-lang/clojure-test-suite](https://github.com/jank-lang/clojure-test-suite) is a cross-dialect compliance suite testing ~230 clojure.core functions with ~5000 assertions. It's a git submodule at `test/clojure-test-suite/`.

## Running

```bash
go test ./test/ -count=1 -run TestClojureTestSuite -v     # full suite, verbose
go test ./test/ -count=1 -run TestClojureTestSuite/abs     # single test file
go test ./... -count=1                                      # all tests (suite runs last)
```

The TOTALS line in verbose output shows the key metrics:
```
files=105 assertions: pass=3530 fail=454 | skipped: compile=115 panic=0 runtime=1
```

## Native lowering coverage (`make jank-stress`)

`TestClojureTestSuite` above measures *runtime conformance* — does the suite run correctly (interpreted), assertion-by-assertion. A separate axis measures *native-lowering coverage* — how many of the suite's fixtures the AOT pipeline can lower to native Go (`gogen_ir`), rather than falling back to interpretation.

```bash
make jank-stress    # lower-go AOT over the suite's deftest bodies; prints :ok/:fail per fixture + buckets
```

Fixtures are the macro-generated test fns (`deftest` bodies), enumerated by the pipeline's canonical `ir.passes.pipeline/lowerable-fn-forms` (siblings pre-declared via `defined-names`) — so the number tracks what production lowering covers, not a harness-private notion.

**Current: 273/273 (100%), 0 failures (verified 2026-07-05).** The suite lowers fully; the earlier 65% baseline (2026-06-23) and its failure buckets — float `unrecognized form` in `are`/`is` bodies, ratios, BigDecimal, `#uuid`, `gogen/func-decl` ExceptionInfo, `nth` build/typeinfer — have all been closed.

**Regression gate:** any failure (`:ok` < 273) is a regression. A newly-surfaced failure after a lowering change is a *discovered gap* — record it with its bucket, don't silently absorb.

## Architecture

- **Test runner**: `test/zz_compat_test.go` — iterates `.cljc` files, compiles through let-go, runs assertions. Named `zz_` so it runs after `language_test.go` (avoids global namespace pollution).
- **Compat shims**: `test/compat/clojure/core-test/portability.lg` — provides `when-var-exists`, `thrown?`, and portability constants that the upstream tests expect.
- **NS aliases**: `pkg/rt/lang.go` maps `clojure.test`→`test`, `clojure.string`→`string`, etc. so upstream `(:require [clojure.test ...])` works transparently.
- **Resolver**: `pkg/resolver/resolver.go` tries `.lg` then `.cljc` extensions, and hyphen→underscore path variants, so upstream files like `number_range.cljc` load from `clojure.core-test.number-range`.

## Safety

Each test runs in a goroutine with:
- **5s timeout** — kills infinite seq realization
- **512MB memory guard** — polled every 200ms, prevents OOM
- **Panic recovery** — catches reader/compiler panics, reports with stack trace

## knownFailing list

`knownFailing` in `zz_compat_test.go` tracks tests that are expected to fail. The runner:
- **Skips** tests in the list that fail (expected)
- **Errors** if a test in the list passes (remove it — it graduated!)

This keeps the list current automatically. After fixes, run the suite and look for "PASSES but is listed in knownFailing" messages.

## How to improve coverage efficiently

### Step 1: Identify the highest-value fix

Run and categorize blockers:

```bash
# Top missing symbols (compile-time)
go test ./test/ -count=1 -run TestClojureTestSuite -v 2>&1 | col -b | \
  grep "Can't resolve" | sed "s/.*Can't resolve //" | sed 's/ in this.*//' | \
  sort | uniq -c | sort -rn | head -20

# Top runtime errors (tests that compile but crash)
for t in conj cons merge sort update ...; do
  echo -n "$t: "
  go test ./test/ -count=1 -run "TestClojureTestSuite/${t}\$" -v 2>&1 | \
    col -b | grep 'caused by' | tail -1 | sed 's/.*caused by //'
done

# Reader errors (syntax the reader can't parse)
go test ./test/ -count=1 -run TestClojureTestSuite -v 2>&1 | col -b | \
  grep 'invalid number\|reading reader' | sort | uniq -c | sort -rn
```

### Step 2: Fix in order of impact

From our experience, the categories ranked by assertions-unlocked:

1. **Reader fixes** (biggest bang) — Each reader fix unlocks entire categories:
   - Number literals (M suffix, hex, octal, radix, ratio) — unlocked ~50 files
   - Reader conditional robustness (skip unreadable branches, handle comments) — unlocked ~25 files
   - Symbolic values (`##Inf`, `##NaN`) — unlocked tests using infinity

2. **Resolver/namespace fixes** — Unlock cascading dependencies:
   - `.cljc` extension + underscore path support — unlocked `number-range.cljc` which ~30 tests depend on
   - Namespace aliases (`clojure.test`→`test`) — required for any upstream test to run

3. **Missing builtins** — Each one typically unlocks 2-6 tests:
   - Trivial: `identical?`, `any?`, `ifn?`, `double?`, `instance?` (one-liners)
   - Medium: `array-map`, `some-fn`, `rseq` (small core.lg functions)
   - Heavy: `sorted-set`, `sorted-map`, `deftype` (need new VM types)

4. **Behavioral fixes** — Fix assertion failures in already-running tests:
   - `conj`/`update` on nil → treat as empty collection
   - `sort` with nil elements → nil-safe comparator
   - Cross-type seq equality → element-by-element comparison
   - `contains?` on vectors → index check
   - `count`/`subs` on strings → rune-based

### Step 3: Fix → regenerate → test → reconcile

```bash
# 1. Make the fix
vim pkg/rt/lang.go  # or core.lg, reader.go, etc.

# 2. Regenerate if core.lg changed
go generate ./pkg/rt/

# 3. Run existing .lg tests first (must not regress)
go test ./test/ -count=1 -run TestRunner

# 4. Run compat suite
go test ./test/ -count=1 -run TestClojureTestSuite -v 2>&1 | col -b | grep 'TOTALS'

# 5. Check for graduated tests
go test ./test/ -count=1 -run TestClojureTestSuite -v 2>&1 | col -b | grep 'PASSES but'

# 6. Check for new unexpected failures
go test ./test/ -count=1 -run TestClojureTestSuite -v 2>&1 | col -b | \
  grep 'FAIL:' | grep -v 'TestClojureTestSuite (' | sed 's/.*\///' | sed 's/ .*//' | sort

# 7. Update knownFailing: remove graduated, add new failures

# 8. Full green check
go test ./... -count=1
```

### Step 4: Commit with metrics

Include before/after assertion counts in the commit message so progress is trackable:
```
Result: 3530 pass / 454 fail across 105 test files (was 2230/372/107).
```

## Common pitfalls

- **Global namespace pollution**: The compat suite creates namespaces that persist. It runs AFTER `language_test.go` (file sorting: `zz_` > `language_`). The NS loader is saved/restored via defer.
- **`go generate` is required** when `core.lg` changes — the precompiled bytecode bundle must match.
- **`nth` bounds checking** can't throw without breaking internal code that relies on nil return. Many macros/builtins use `nth` internally and expect nil for out-of-bounds.
- **`drop`/`drop-while` on nil** — Clojure returns `()`, let-go returns `nil`. Changing this breaks existing `.lg` tests that assert `(nil? (drop 10 [1 2]))`.
- **BigDecimal `M` suffix** is parsed as Float — type predicates like `decimal?`, `rational?` won't match. This is a known approximation.
- **`PersistentSet` implements `Seq`** — this is wrong (sets should be Seqable, not Seq), but removing it is a large refactor. `sequential?` and `isSequentialType` explicitly exclude sets as a workaround.
- **`valueEquiv` in the HAMT** uses `ValueEquals` callback for seq comparison. Guard against infinite recursion — only delegate for Seq types, not maps/sets.
- **Reader conditional branches** may contain dialect-specific syntax (`#cpp`, `(.-field obj)`). The `skipReaderForm` function handles this by counting balanced delimiters without parsing. Comments inside `#?()` also need explicit handling.

## File locations

| File | Purpose |
|------|---------|
| `test/zz_compat_test.go` | Go test runner, knownFailing list, timeout/memory guards |
| `test/compat/clojure/core-test/portability.lg` | when-var-exists, thrown?, portability constants |
| `test/clojure-test-suite/` | Git submodule (upstream test files) |
| `pkg/rt/lang.go` | Builtins, NS aliases, valueEquals |
| `pkg/rt/core/core.lg` | Core library functions (keys, vals, find, etc.) |
| `pkg/compiler/reader.go` | Number literals, reader conditionals, symbolic values |
| `pkg/resolver/resolver.go` | .cljc extension, underscore path support |
| `pkg/vm/hash.go` | valueEquiv (HAMT key comparison) |
| `pkg/vm/persistent_set.go` | PersistentSet Lookup interface (for get) |
