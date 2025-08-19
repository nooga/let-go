## Clojure-like collections, seq tower, transients and transducers — Refactor Plan

This document captures the plan to refactor VM data structures and core APIs to better align with Clojure semantics while minimizing breakage. It covers: current gaps, target interface hierarchy (seq tower), persistent data structures, transients, transducers, migration steps, tests, and risks.

### Goals

- Align collection/sequence semantics with Clojure: proper ISeq/Seqable separation, correct conj/first/next/rest behavior, and nil handling.
- Replace ad-hoc collections with persistent versions: vector (BPTR), hash map (HAMT), hash set.
- Introduce structural equality and hashing for use as map/set keys.
- Add transients for efficient bulk-building without changing public persistent APIs.
- Add transducers and reduce fast paths, preserving existing API semantics.
- Land changes incrementally with high test coverage and adapters to avoid large-scale breakage.

### Current state — key gaps

- Collections implement `Seq` directly instead of only `Seqable`. Builtins like `first`/`next` accept `Seq` instead of coercing via `seq`.
- `ArrayVector` implements `Seq` methods; should instead expose a vector-specific seq view.
- `Map` and `Set` are backed by Go `map` and implement `Seq` by converting to `List`; they are copy-on-write, not structurally persistent.
- `PersistentVector` exists but deviates from canonical 32-ary bit-partitioned trie details (variable node width, edge-case handling, pop/peek, chunked seq).
- No `MapEntry` abstraction; `conj` on map expects a 2-item vector only.
- Equality/hash are not Clojure-like: `=` uses raw equality; no value-based `equiv`; no `hash` for keys; map/set rely on Go key equality.
- Reader returns `ArrayVector` for vector literals; maps are copy-on-write.

### Target interface hierarchy (Clojure-aligned)

- Core values (existing):

  - `Value`, `ValueType`.
  - Add: `Equivable` { `Equiv(Value) bool` }, `Hashable` { `Hash() int` } for structural equality/hash across types.

- Seq tower:

  - `Seq` (ISeq): `First() Value`, `Next() Seq`, `More() Seq`, `Cons(Value) Seq`.
  - `Sequable`: `Seq() Seq`.
  - Only sequences implement `Seq`. Collections implement `Sequable`.
  - Optional later: `Reversible` { `RSeq() Seq` }, `ChunkedSeq` for performance.

- Collections:

  - `Counted` (existing): `RawCount() int`, `Count() Value`.
  - `Collection` (IPersistentCollection): `Counted`, `Empty() Collection`, `Conj(Value) Collection`.
  - `Lookup` (existing): `ValueAt(Value) Value`, `ValueAtOr(Value, Value) Value`.
  - `Associative` (existing): `Assoc(Value, Value) Associative`, `Dissoc(Value) Associative`.
  - Add: `Indexed` { `Nth(Int) Value`, `NthOr(Int, Value) Value` } for vectors.
  - Add: `Stack` { `Peek() Value`, `Pop() Collection` } (lists/vectors).
  - Add: `IMapEntry` { `Key() Value`, `Val() Value` }.

- Transients:

  - `Editable` { `AsTransient() (TransientCollection, error)` }.
  - `TransientCollection` { `ConjBang(Value) error`, `PersistentBang() (Collection, error)` }.
  - `TransientAssociative` { `AssocBang(Value, Value) error`, `DissocBang(Value) error` }.
  - `TransientVector` { `PopBang() error` } (and `Peek()` via `TransientStack` if needed).
  - Edit-token enforcement (single-owner, refuse after `persistent!`).

- Reducers/transducers:
  - `Reducible` { `Reduce(fn Fn, init Value) (Value, error)` } fast path for `reduce`.
  - `Reduced` marker wrapper for early termination.

### Behavioral semantics to match Clojure

- Seq coercion:
  - `(seq nil) => nil`. `(first x)`/`(next x)`/`(more x)` operate via `seq` on `Seqable` values; `more` returns empty list for empty seq; `next` returns `nil` at end.
- `conj` semantics:
  - List: prepend; Vector: append; Map: accepts `IMapEntry` or kv pairs; Set: adds element.
- `count`:
  - `(count nil) => 0`.
- `peek`/`pop`:
  - Use `Stack` if implemented; for lists, `peek=first`, `pop=next`; for vectors, top is the last element.
- Equality/Hash:
  - Structural `Equiv` across all collection types; numeric comparisons by value; maps/sets equality by entries; `Hash` consistent with `Equiv` and stable.

### Persistent data structures

- PersistentVector:

  - Canonical 32-ary bit-partitioned trie (BPTR): fields `count`, `shift`, `root` (node with fixed 32 slots), `tail` (<= 32), `tailOff`.
  - Ops: `Nth`, `AssocN`, `Conj` (append), `Pop`, `Peek`, `Seq` (non-chunked first, chunked later), `AsTransient` with edit token.
  - Helpers: `arrayFor(i)`, `pushTail`, `doAssoc`.

- PersistentHashMap:

  - HAMT (bitmap-indexed nodes), collision nodes for hash collisions, structural sharing.
  - Ops: `Assoc`, `Dissoc`, `ValueAt`, `Contains`, `Seq` of `IMapEntry`.
  - Backed by `Equiv`/`Hash` of keys.

- PersistentHashSet:
  - Wrapper over `PersistentHashMap` (values ignored), same semantics and seq over keys.

### Transients (performance-friendly mutability)

- API (Clojure-like): `transient`, `persistent!`, `conj!`, `assoc!`, `dissoc!`, `disj!`, `pop!`.
- Implementation: edit token stored in nodes and transient wrapper; mutate in place when the token matches; refuse operations after `persistent!`.
- Start with `PersistentVector` transient; later add map/set transients.

### Transducers and reduction fast paths

- 1-arity transformers: `(map f)`, `(filter p)`, `(take n)`, `(drop n)`, `(take-while p)`, `(remove p)`, `(cat)`, `(mapcat f)` return transducers.
- New core fns: `transduce`, `eduction`, `sequence`, `completing`.
- `reduce` prefers `Reducible` fast path; fall back to `Seq` traversal.
- `into` implemented via `transduce`, exploiting collection-specific fast paths and transients (`into [] xs`).

### Migration strategy — staged to minimize breakage

Phase 0: Interfaces, adapters, and builtin coercion (no collection rewrites)

- Add new interfaces: `Indexed`, `Stack`, `IMapEntry`, `Equivable`, `Hashable`, `Reducible`, `Reduced`.
- Update builtins to coerce via `seq` and accept `nil`:
  - `first`, `second`, `next`, `more`, `rest`, `seq`, `count`, `peek`, `pop`, `conj`, `contains?`, `get`.
- Add `MapEntry` type and make map `conj` accept `IMapEntry` or 2-vec.
- Keep existing types working (lists/vectors/maps/sets still accepted).
- Tests: seq coercion across types; `(count nil)`; `first/next/more/rest` on empty/non-empty; `conj` semantics per type.

Phase 1: PersistentVector correctness and seq

- Reimplement `PersistentVector` as canonical BPTR with `Nth`, `AssocN`, `Conj`, `Pop`, `Peek`, `Seq`.
- Add `TransientVector` with edit token and `conj!`/`assoc!`/`pop!`.
- Reader: switch vector literals to `PersistentVector`. Keep `ArrayVector` as lightweight builder/interop.
- Tests: boundary sizes (31/32/33, multiples of 32); assoc on tree/tail boundaries; pop across tail/tree; random access; seq traversal; transient bulk build vs persistent equivalence.

Phase 2: PersistentHashMap and PersistentHashSet

- Implement `Equiv`/`Hash` for core value types (numbers, strings, symbols, keywords, booleans, chars, nil, vectors, lists, maps, sets).
- Implement `PersistentHashMap` (HAMT) + `TransientHashMap`.
- Implement `PersistentHashSet` on top of map + `TransientHashSet`.
- Constructors and reader:
  - `hash-map`, `hash-set` return persistent versions.
  - Reader map/set literals build persistent versions.
- Tests: assoc/dissoc, collision handling, equality across different insertion orders, membership, keys/vals, `find`, set ops.

Phase 3: Equality and hashing in core

- Replace `=` with structural `Equiv` (numbers by value; sequences element-wise; maps/sets by entries).
- Add/adjust `hash` if exposed; ensure map/set use `Hash`/`Equiv` for keys.
- Update `contains?`, `get`, set membership to use new semantics.
- Tests: cross-type equality correctness, map/set equality independent of order, hashing invariants.

Phase 4: Performance polish and transducers

- Add `ChunkedSeq` for vectors/ranges; optimize `map/reduce` pipelines.
- Implement `Reducible` for vectors, ranges, possibly maps/sets seqs.
- Add transducer functions: `transduce`, `eduction`, `sequence`, `completing`.
- Make `into`, `mapv`, `keep`, `keep-indexed`, `frequencies`, `group-by` use `transduce` where profitable.
- Tests/benchmarks: verify no intermediate allocations for common `into` cases and improved throughput.

### Compatibility considerations

- Maintain `ArrayVector` and current map/set types during transition; builtins accept both via interfaces.
- Flip reader for vectors/maps/sets only when persistent implementations and tests are in place.
- Do not change public arities/semantics of existing core fns; only expand to support `nil`/`Seqable` coercion and add new features.
- Provide clear error messages for transient misuse (after `persistent!`).

### Testing plan (high-value cases)

- Seq coercion: `(first [])`, `(next {})`, `(more nil)`, `(rest nil)`.
- Conj semantics across list/vector/map/set; map `conj` accepting `IMapEntry`.
- Vector boundaries: 31/32/33, exact multiples of 32; assoc at edges; pop through tail->tree transitions.
- Map correctness: deep updates/removals, collision nodes, equality across insert orders.
- Set operations: `difference`, `intersection`, `distinct`, membership.
- Equality/hash: structural equality across different concrete implementations; stable hash.
- Transients: vector `transient`/`conj!`/`assoc!`/`pop!`; error on use-after-persistent; persistent result equivalence.
- Transducers: `transduce (map inc)` equivalence; `into` no intermediate seqs; `eduction` correctness.
- Benchmarks: create/access/conj/assoc/seq across `ArrayVector` vs `PersistentVector`; `into` with and without transients/transducers.

### Work items checklist

- [ ] Phase 0: Add interfaces (`Indexed`, `Stack`, `IMapEntry`, `Equivable`, `Hashable`, `Reducible`, `Reduced`).
- [ ] Phase 0: Update builtins for seq coercion and nil handling.
- [ ] Phase 0: Implement `MapEntry` and adjust map `conj`.
- [ ] Phase 0: Add tests for coercion and semantics.
- [ ] Phase 1: Reimplement `PersistentVector` + tests.
- [ ] Phase 1: Implement `TransientVector` + tests; optimize `into []`.
- [ ] Phase 1: Switch reader vectors to persistent.
- [ ] Phase 2: Implement `Equiv`/`Hash` across types.
- [ ] Phase 2: Implement `PersistentHashMap`/`Set` + transients + tests.
- [ ] Phase 2: Switch reader maps/sets and constructors.
- [ ] Phase 3: Replace core equality/hash usage; update `contains?`/`get`; tests.
- [ ] Phase 4: Add chunked seqs, `Reducible` impls; introduce `transduce`/`eduction`/`sequence`/`completing`; refit `into`/`mapv`/etc.; perf tests.

### Open questions / decisions

- Final hash function details (mixing, seeded?) and cross-type hashing for numeric types.
- Equality coercions (e.g., Int vs wider numerics) — for now: strict Int only unless numeric tower expands.
- Ordering guarantees of map/set seqs (HAMT typically yields implementation-defined but stable per structure).
- How much of `IReduce` vs `IReduceInit` to model; start minimal.

### References (for implementation guidance)

- Clojure `clojure.lang` sources: `PersistentVector`, `Node`, `APersistentMap`, `PersistentHashMap`, `PersistentHashSet`, `APersistentVector$Seq`, `ISeq`, `Seqable`, `IEditableCollection`, `Transient*`, `IReduce`, `IReduceInit`, `Reduced`, `ChunkedSeq`.

### Relevant file index (what they do and what to mind)

- pkg/vm/value.go: Core `Value` and interface definitions (`Seq`, `Sequable`, `Collection`, `Associative`, `Lookup`, `Fn`, etc.).

  - mind: We’ll introduce new interfaces (`Equivable`, `Hashable`, `Indexed`, `Stack`, transient interfaces). Keep existing ones stable during Phase 0; avoid breaking type assertions in `rt/lang.go`.

- pkg/vm/vector.go: `ArrayVector` (current vector) implementing `Seq`, `Collection`, `Lookup`, `Associative` in a simple slice-backed way.

  - mind: Keep for interop and tests while migrating literals/constructors to `PersistentVector`. Eventually make it only `Sequable` (not `Seq`).

- pkg/vm/persistent_vector.go: `PersistentVector` (BPTR-like) and `PersistentVectorSeq` with basic `Assoc`, `Conj`, `Seq`.

  - mind: Rework to fixed 32-slot nodes; add `arrayFor`, `Pop`, `Peek`, transients with edit token, edge-case correctness at 31/32/33 and multiples of 32. Ensure `Seq` correctness and later chunked seq.

- pkg/vm/list.go: `PersistentList` (`List`), currently correct and simple.

  - mind: Ensure it remains `Seq` and `Sequable`. Implement `Stack` behavior (`Peek`=`First`, `Pop`=`Next`).

- pkg/vm/map.go: `Map` backed by Go map, copy-on-write assoc/dissoc, ad-hoc seq via list of 2-vectors.

  - mind: Replace with `PersistentHashMap` (HAMT). Add `MapEntry` abstraction. Ensure lookup/equality/hash use new semantics. Reader and `hash-map` should target the persistent impl.

- pkg/vm/set.go: `Set` backed by Go map keys; copy-on-write.

  - mind: Replace with `PersistentHashSet` backed by persistent map; implement transients; update set ops to new semantics.

- pkg/vm/vector_test.go: Extensive tests for vector behaviors including `PersistentVectorSeq`.

  - mind: Extend with boundary and pop tests; keep `ArrayVector` vs `PersistentVector` parity where intended; add transient tests.

- pkg/vm/vector_benchmark_test.go: Benchmarks for create/access/conj/assoc/seq between Array vs Persistent vectors.

  - mind: After correctness, re-run to gauge perf; add `into`/transient/transducer benchmarks.

- pkg/rt/lang.go: Core runtime functions (`first`, `next`, `seq`, `conj`, `assoc`, `reduce`, etc.).

  - mind: Phase 0 changes here: coerce via `seq`, accept `nil`, and add transducer/transient APIs (`transient`, `persistent!`, `conj!`, `assoc!`, `transduce`, `eduction`, `sequence`, `completing`). Replace `=` with structural `Equiv` in Phase 3.

- pkg/rt/core/core.lg: Core macros and higher-level fns written in let-go.

  - mind: Depends on semantics of `first/next/rest/seq/conj/count`. After Phase 0, behavior should match Clojure closer; review places relying on current stricter `Seq` requirement.

- pkg/compiler/reader.go: Reader for literals and macros; constructs vectors, maps, sets.

  - mind: Switch vector/map/set literals to persistent types in the phases indicated. `flattenMap` assumes non-persistent `Map`; update once HAMT lands.

- pkg/vm/keyword.go, symbol.go, string.go, int.go, bool.go, char.go: Scalar types.

  - mind: Add `Equiv`/`Hash` consistently. Decide numeric tower rules (for now Int-only equivalence).

- pkg/vm/namespace.go, var.go, func.go, native_func.go: Runtime plumbing and function invocation.

  - mind: Ensure any new interfaces used by core (`Fn`, reducers) remain compatible; adjust `reduce` fast path if adding `Reducible`.

- test/language_test.go: Integration tests for the language.
  - mind: Add tests for `(count nil)`, seq coercion, conj semantics, transients/transducers once implemented.
