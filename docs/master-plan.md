## let-go master plan — fastest and most useful Clojure-on-Go

This roadmap consolidates our existing plans for VM performance, Clojure-like collections, numeric and value representation, and runtime images. It defines staged milestones, success criteria, and dependencies to make let-go the fastest and most useful Clojure-on-Go implementation.

### Vision and success criteria

- **Startup**: cold start with stdlib image < 50 ms on a typical laptop; warm start < 5 ms.
- **Throughput**: 2–3x improvement on core functional pipelines (`map`/`filter`/`reduce`) over current baseline; zero-regression guardrails in CI.
- **Allocation profile**: O(1) allocations for transducer-based pipelines over vectors and ranges; no unbounded retention from calling convention.
- **Semantics**: Clojure-aligned collections and seq tower; structural equality and hashing; transients and transducers.
- **Deployability**: precompiled stdlib image by default; app/runtime images supported and versioned; deterministic serialization.

### Guiding principles

- **Correctness first**: codify Clojure semantics via interfaces and tests before aggressive optimizations.
- **Tight hot paths**: specialize VM call paths and numeric primitives; avoid reflection and unnecessary boxing.
- **Incremental rollout**: adapters and compatibility layers to avoid big-bang rewrites; keep REPL usable throughout.
- **Benchmark-driven**: every phase ships microbenchmarks and end-to-end scenarios, with perf budgets enforced in CI.

## Phases and milestones

### Phase 0 — Baseline semantics, numeric fast paths, and benchmarks

Scope:
- Introduce interfaces to align with Clojure and enable later features:
  - `Equivable`, `Hashable`, `Indexed`, `Stack`, `IMapEntry`, `Reducible`, `Reduced` in `pkg/vm/value.go`.
- Update builtins to coerce via `seq` and support `nil` like Clojure (`first`, `next`, `more`, `rest`, `seq`, `count`, `peek`, `pop`, `conj`, `contains?`, `get`).
- Implement `MapEntry` and accept `IMapEntry`/2-vec in map `conj`.
- Numeric and value representation polish:
  - Switch `Int.String()` to `strconv.Itoa`.
  - Replace `Unbox().(int)` with direct `vm.Int` assertions in numeric natives in `pkg/rt/lang.go`.
  - Ensure native registrations use `Wrap`/`WrapNoErr` instead of `Box` on hot paths.
- Establish benchmark harness and CI perf gates.

Dependencies: none.

Acceptance:
- Tests for seq coercion and `(count nil)` pass across core types.
- Benchmarks added: arithmetic loops, seq ops baselines; perf budgets recorded.
- Numeric primitives show fewer allocations and improved ns/op over baseline.

References: `docs/clojurelike-refactor-plan.md` (Phase 0), `docs/value-representation-and-numeric-performance.md` (checklist).

### Phase 1 — VM calling convention and memory retention fixes

Scope (see `docs/vm-performance-optimization.md`):
- Copy argument slices before entering bytecode callees, including tail-call growth branch.
- Extend TCO to closures (`*Closure`) in `OP_TAIL_CALL`.
- Add frame/stack pooling with `sync.Pool` and reset lifecycle.
- Introduce small-arity invoke/tail-call opcodes: `INVOKE_0/1/2/3` and `TAIL_CALL_0/1/2/3`; update compiler emission.
- Audit runtime natives to avoid reflection (`Wrap`/`WrapNoErr`).

Dependencies: Phase 0 (interface changes might be referenced by reducers later, but VM changes are largely independent).

Acceptance:
- Tail-recursive functions (including closures) do not grow stack; allocations reduced vs baseline.
- Call-heavy microbenchmarks improve (>30% fewer allocs, measurable ns/op reduction).
- No retention via arg slices in pprof heap; frame/stack pools show reuse.

### Phase 2 — PersistentVector (BPTR) + transients; switch reader for vectors

Scope (see `docs/clojurelike-refactor-plan.md`):
- Reimplement `PersistentVector` as canonical 32-ary BPTR with `Nth`, `AssocN`, `Conj`, `Pop`, `Peek`, `Seq`.
- Implement `TransientVector` with edit token; `conj!`, `assoc!`, `pop!`.
- Reader: vector literals produce `PersistentVector`. Keep `ArrayVector` for interop and builders.
- Add `Reducible` for vectors; later add `ChunkedSeq` in Phase 4.

Dependencies: Phase 0 (interfaces, semantics). Phase 1 optional but beneficial for end-to-end perf.

Acceptance:
- Boundary tests (31/32/33, multiples of 32) pass; seq traversal correct.
- `into [] xs` fast path leverages transients with clear wins over baseline.
- Benchmarks show competitive performance for random access and append vs previous vector.

### Phase 3 — PersistentHashMap/Set, structural equality and hashing

Scope:
- Implement `Equiv`/`Hash` for all scalar and collection types.
- Implement `PersistentHashMap` (HAMT) + `TransientHashMap`; `PersistentHashSet` backed by map + transient variant.
- Reader: map/set literals return persistent versions; constructors `hash-map`, `hash-set` updated.
- Map `conj` accepts `IMapEntry`/2-vec; set ops align with Clojure.

Dependencies: Phase 0 (interfaces), and vector for deterministic image encoding in Phase 5.

Acceptance:
- Equality/hash tests pass across heterogeneous types and different insertion orders.
- Map/set correctness: assoc/dissoc, collision handling, membership, `find`.
- Deterministic seq order for serialization achieved (stable encoding independent of Go map order).

### Phase 4 — Transducers and reduction fast paths

Scope:
- Add transducer APIs: `transduce`, `eduction`, `sequence`, `completing` in `pkg/rt/lang.go` and `pkg/rt/core/core.lg`.
- Implement `Reducible` for vectors and ranges; add `Reduced` early-termination.
- Add `ChunkedSeq` for vectors/ranges to reduce per-element overhead.
- Make `into`, `mapv`, `keep`, `keep-indexed`, `frequencies`, `group-by` use transducers where profitable.

Dependencies: Phase 2 (vectors) and optionally Phase 3 (for equality-dependent ops).

Acceptance:
- Pipelines like `(into [] (comp (map inc) (filter odd?)) xs)` allocate O(1) and run faster than seq-based pipelines.
- Benchmarks show 2–3x improvement on common reduce/transduce patterns vs baseline.

### Phase 5 — Runtime images and precompiled stdlib

Scope (see `docs/runtime-image-and-stdlib-cache.md`):
- Define image schema and versioning; implement value, const-pool, code, namespace, and var serialization.
- Implement host/native extern resolution via `HostRegistry`.
- Build `stdlib.img` during build; boot attempts image before compilation; fallback on mismatch.
- Add `image/save` and `image/load` APIs; support app-layer images on top of stdlib.
- Ensure deterministic encoding for persistent collections (BPTR/HAMT based) and intern tables.

Dependencies: Phase 2 and 3 (persistent collections) for determinism; Phase 1 (VM code chunk stability).

Acceptance:
- Cold start with stdlib image < 50 ms; fallback path regenerates image when invalidated.
- Round-trip tests for values/code/vars; extern resolution verified; security validations in place.

### Phase 6 — Tooling, interop, and developer experience

Scope:
- Solidify `pkg/nrepl` server ergonomics and error reporting.
- Audit all `pkg/rt/lang.go` registrations to ensure fast-native wrappers; document interop patterns and when to use `Box`.
- Improve error messages for transients misuse and `seq` coercions.
- WASM: ship a minimal runtime that loads `stdlib.img` for instant REPL/page load; document build in `wasm/`.

Dependencies: after Phase 5 for best UX.

Acceptance:
- Smooth REPL and editor integration; quick startup in WASM demo; clear interop docs.

### Phase 7 — Advanced performance polish

Scope:
- Finish `ChunkedSeq` across collections; add reducer fast paths for maps/sets as needed.
- Optimize iteration utilities (`range`, `repeat`, `iterate`) with tight loops and `Reducible`.
- Consider specialized opcodes for common sequence ops if warranted by profiles.
- Investigate parallel reductions/goroutine-based pipelines where semantically sound (opt-in only).

Dependencies: prior phases complete; driven by profiling.

Acceptance:
- Profiles of real programs show minimal interpreter overhead dominating, with near-zero allocations for common loops.

## Cross-cutting: benchmarks, CI, and risk management

Benchmarks to maintain:
- Micro: call overhead (arity 0–3 vs generic), TCO recursion, vector ops (conj/assoc/pop/seq), HAMT map ops, numeric arithmetic, transducer pipelines.
- End-to-end: `into []` with composed transforms, HTTP JSON handling via natives, REPL cold/warm start.

CI gates:
- Record baseline ns/op and allocs/op; alert on >10% regressions in hot benchmarks.
- Add `go test -bench` jobs and simple perf dashboards; run pprof sampling periodically.

Risks and mitigations:
- Opcode proliferation complexity → keep generic slow path; add exhaustive tests for arity-specialized ops.
- Persistent data structure correctness → comprehensive boundary tests and randomized property tests.
- Image format lock-in → versioned schema and strict validation; maintain migration path.

## Immediate workboard (next actions)

- Phase 0
  - Add interfaces in `pkg/vm/value.go` and adapters in builtins (`pkg/rt/lang.go`).
  - Implement `MapEntry`; update map `conj` semantics.
  - Numeric fast paths and `Int.String()` change.
  - Add initial benchmarks: arithmetic loops; seq coercion scenarios.

- Phase 1
  - Implement arg copy in `OP_INVOKE`/`OP_TAIL_CALL` and tail-call growth.
  - Add closure TCO; frame/stack pools; small-arity opcodes and compiler emission.
  - Benchmarks: call overhead, recursion, memory profiles.

- Phase 2–4
  - PersistentVector + transients and reader switch; then HAMT map/set + transients and reader switch; then transducers and chunked seqs.

- Phase 5
  - Image schema and loader/saver; stdlib image integration at boot.

## File index and where each phase touches

- `pkg/vm/value.go`: interfaces (`Equivable`, `Hashable`, `Indexed`, `Stack`, `IMapEntry`, `Reducible`, `Reduced`).
- `pkg/rt/lang.go`: builtins semantics, numeric primitives, transducers, transients API.
- `pkg/vm/vm.go`, `pkg/vm/func.go`, `pkg/vm/native_func.go`: VM call paths, TCO, pooling, native wrappers.
- `pkg/vm/persistent_vector.go`, `pkg/vm/vector.go`, `pkg/vm/list.go`: vectors/lists; persistent/tran sient vector.
- `pkg/vm/map.go`, `pkg/vm/set.go`: replace with persistent HAMT-based implementations and sets.
- `pkg/compiler/reader.go`: switch literals to persistent types at the right phases.
- `pkg/rt/core/core.lg`: higher-level functions relying on seq/coercion; add transducer helpers.
- `pkg/vm/constpool.go`, `pkg/vm/namespace.go`, `pkg/vm/var.go`, `pkg/vm/func.go`: image serialization.

## Success metrics snapshot (to refine with data)

- Cold start with `stdlib.img`: < 50 ms CLI; WASM demo interactive < 200 ms.
- `map`/`filter`/`reduce` pipeline on `PersistentVector 1e6`:
  - Allocations: O(1) with transducers; ns/op improved by 2–3x vs baseline seq pipeline.
- Call-heavy microbenchmark (100M nullary calls):
  - >25% ns/op reduction; >50% allocs/op reduction vs baseline.
- Vector ops (conj/assoc/pop/peek at boundaries):
  - Parity with Clojure algorithms; stable across edge cases.

## References

- `docs/vm-performance-optimization.md`
- `docs/clojurelike-refactor-plan.md`
- `docs/value-representation-and-numeric-performance.md`
- `docs/runtime-image-and-stdlib-cache.md`


