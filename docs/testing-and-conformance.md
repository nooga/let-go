## Testing and Conformance — framework, CI, and compatibility strategy

This document defines how we test let-go: unit/integration/perf tests, a `clojure.test`-compatible layer, conformance approach vs Clojure/Babashka/Joker, property testing, fuzzing, and CI/perf gating.

### Goals

- Provide a familiar `clojure.test`-style experience for writing tests in let-go code.
- Build a conformance suite covering the supported Clojure subset, persistent data structures, seq semantics, equality/hash, and runtime behavior.
- Integrate performance benchmarks with guardrails to avoid regressions.
- Enable portability and reproducibility: deterministic tests, golden outputs, and clear skip/focus controls.

### Test framework (user-facing)

- `clojure.test`-compatible API (subset):
  - Macros/functions: `deftest`, `testing`, `is`, `are`, `use-fixtures` (once/each), `run-tests`.
  - Output formats: human-readable default; optional TAP and JUnit XML for CI.
  - Selectors: include/exclude by ns or metadata; `:only`, `:focus`, `:skip`.
- CLI: `lg test` supports:
  - Namespace globs or files; `--watch` to re-run on changes; `--junit out.xml`; `--tap`.
  - `--fail-fast`, `--seed` for randomized/property tests.
- Implementation sketch:
  - Core under `pkg/rt/core/test.lg` (macros/helpers), runner in `pkg/rt/lang.go` or `pkg/rt/core/core.lg`.
  - JUnit/TAP encoders in Go for performance; exposed to let-go via natives.

### Conformance strategy

- Define scope: publish a spec of the supported Clojure subset and deviations.
- Build a compatibility suite:
  - Port a curated subset of `clojure.core` behavior into `clojure.test`-style tests that avoid JVM-specific features.
  - Import/adapt relevant Babashka/Joker tests where semantics overlap (with attribution and license compliance).
  - Reader/printer: golden tests with EDN fixtures; normalize map order for comparison.
  - Collections: vectors (BPTR), maps (HAMT), sets — correctness across boundary cases; equality/hash invariants.
  - Seq tower: `seq`, `first/next/more/rest`, `conj`, `count (nil)`, `peek/pop` semantics.
  - Transients/transducers: behavior and error cases.
  - Numeric: arithmetic, comparisons, edge cases (overflow per current numeric model), string conversion.
- Round-trip tests with runtime images and (optional) Go AOT embedding.

### Property testing and fuzzing

- Provide a minimal `test.check`-style library (generators + `prop/for-all`):
  - Properties: vector append/assoc/pop invariants; HAMT assoc/dissoc/idempotence; equality/hash consistency.
  - Shrinking support for common generators (ints, vectors, maps).
- Fuzzing:
  - Reader fuzzing with Go fuzz (`testing/fuzz`) for literals, numbers, and nested collections; reject invalid forms cleanly.
  - VM/bytecode: differential tests comparing VM vs Go AOT embedded functions on random inputs (where defined).

### Performance tests and guardrails

- `go test -bench` microbenchmarks covering: VM call arities, TCO recursion, vector ops, HAMT ops, numeric arithmetic, transducer pipelines, image load time.
- Benchstat thresholds enforced in CI: alert on >10% regressions in selected benchmarks.
- Baselines stored per branch; periodic refresh with manual approval.

### CI integration

- Matrix: macOS/Linux/Windows; Go versions (latest two); with/without WASM build smoke test.
- Steps:
  - Lint (`go vet`, `staticcheck`), unit/integration tests, property tests (seeded), fuzz smoke (short), benchmarks (short mode on PRs, full on scheduled runs).
  - Artifacts: JUnit XML, TAP logs, coverage, perf CSV.

### Acceptance (Phase 0–1)

- `clojure.test`-compatible layer and `lg test` CLI exist; human-readable + JUnit output.
- Conformance seed suite running in CI for core semantics (seq coercion, count nil, conj semantics, basic equality).
- Initial property tests for vectors and numeric operations.
- Bench harness wired with at least call-arity and numeric microbenchmarks.
