---
status: active
last-verified: 2026-06-05
authoritative-for:
  - parallel-lowering-findings
  - type-discovery-cache-design
human-verified:
---

# Parallel IR lowering, and a mergeable type-discovery cache

Status: design note / findings. Captures the 2026-05-30 investigation into
parallelizing the IR lowering passes, what blocked it, and the proposed
architecture (a mergeable type cache) for making parallel lowering both fast
*and* reproducible.

## Goal

`lgbgen --target=go` lowers every core namespace's defns to Go. It is the
slowest step in the bootstrap (~3 min). The per-defn work (expand → build IR →
optimize → typeinfer → lower) is embarrassingly parallel across defns, so the
question was whether we can lower a namespace's defns concurrently.

## What blocked it, and what we fixed

### 1. The future-spawn binding cost (fixed)

The obvious move — `pmap` over the defns — gave **zero** speedup (actually
slightly slower than sequential), even after fixing the gensym race below.

Cause: `pmap` is built on `future`, and every `future` calls
`SnapshotBindings` + `RunWithBindings`, which take the global `bindingsMu`
**and copy every active dynamic binding** per spawn. A future is *async* — it
may run after the caller's `binding` scope has unwound — so it must capture and
reinstate that scope. That spawn cost serializes the work.

Two prerequisites made a real fix possible:

- **Lock-free `Var.Deref`** (atomic `root`/`curr` instead of the global
  `bindingsMu`): dynamic-var *reads* no longer serialize. ~670–720× faster
  parallel deref. (Landed separately.)
- **`pmapv`** — an eager, order-preserving, *synchronous* parallel map. Because
  it blocks until every element finishes, the caller's `binding` scope is still
  live, so workers just **read** the caller's global bindings via the lock-free
  Deref. No snapshot, no `RunWithBindings`, no `bindingsMu` on spawn.

Result on the `core` lowering: future-`pmap` 213 s (0.90× — slower than
sequential 193 s), **`pmapv` 155 s (1.24×, ~1.9 effective cores)**. The win is
the binding fix; gensym was a red herring for *speed*.

Caveat: `pmapv` is not Clojure `pmap`. Clojure's `pmap` is a *lazy, chunked*
seq with bounded lookahead; `pmapv` is eager and shares the caller's dynamic
bindings — correct for "do all of these in parallel and collect," not a drop-in
for lazy chunked semantics.

### 2. Non-deterministic output (partially fixed)

Two parallel runs produced semantically-identical but textually-different Go.
Sources, in order of discovery:

- **gensym counter** — a global `int` incremented per `gensym`. Made it
  `atomic.Int64` (fixes the data race). But atomicity alone leaves the absolute
  numbers scheduling-dependent.
- **gensym numbers in emitted names** — `value-name` builds `<src>_<nid>` where
  `nid` (the IR instruction id) is already deterministic and unique per
  function, so the gensym digits in `<src>` (e.g. `doseq_seq__10369`) are pure
  non-determinism. Stripping the `__<digits>` tail makes names reproducible
  (`doseq_seq_23`). `build-args` likewise names arg temps by inst-id, not
  gensym.
- **import grouping/order** — minor; sort the import set.
- **the wall-clock typeinfer budget** — the deep one (below). *Not yet fixed.*

`order-defn-forms` was already deterministic (function order is stable).

### 3. The wall-clock typeinfer budget (root cause of the remaining non-determinism)

`typeinfer` bails when it exceeds `*typeinfer-budget-ms*` (2000 ms), returning a
**sound partial** result (a lower bound in the type lattice). Under parallel
load, N workers contend for CPU, so a given typeinfer call gets less wall-time,
hits the budget sooner, and bails with *less* type information. A call that
qualified for direct native (`corefns.X`) dispatch sequentially falls back to
generic dispatch in parallel — so even the *imports* differ run-to-run.

**A time-based budget is inherently non-deterministic under variable CPU load.**
This is the one source the per-site fixes can't reach.

## Proposed fix: a mergeable type-discovery cache

Rather than a deterministic *iteration* budget (which gives determinism but
permanently caps precision on hard functions), cache the type discoveries and
**merge** them across runs.

### Why it's sound

typeinfer is **monotone**: `set-type-if-changed!` is a lattice *join* — types
only move up (`:unknown → :int → :number → :any`), never down. The budget bails
with a lower bound. So:

- A cache entry is a point in the type lattice.
- Merging two entries is the **least upper bound** (join).
- Join is commutative, associative, idempotent ⇒ **merge order is irrelevant**.

Two runs that each got partway merge to a result at least as precise as either,
deterministically. Parallel workers can contribute discoveries to a shared
cache and the converged result is independent of scheduling. That removes the
budget-vs-load non-determinism: lowering decisions become a function of the
*converged cache*, not of who-won-the-CPU.

### Shape

- **Key:** a content hash of the function's IR (stable across runs;
  auto-invalidates when the source changes). This is the same
  "deterministic hash for elements" idea, reused as a cache key.
- **Value:** inferred signature (arg types, result type) ± per-inst types.
- **Merge:** per-field lattice join; parallel writers compare-and-join per key.
- **Convergence:** interprocedural — F's inference improves when its callees'
  cached signatures improve, so iterate to a whole-program fixpoint. Monotone +
  finite lattice height ⇒ it terminates; the cache makes re-inference
  incremental (only functions whose callees moved).
- **Committed artifact:** like the `.lgb` bundle and the lowered tree, commit
  the converged cache. A "warm to fixpoint" step (sequential, unbudgeted —
  correctness over speed) produces it; `make generate` runs it; the staleness
  manifest covers it. Lowering then *reads* the cache — fast, parallel,
  reproducible — and the wall-clock budget can be dropped from the hot path.

### Comparison

| | wall-clock budget (today) | deterministic iteration budget | mergeable type cache |
|---|---|---|---|
| deterministic | no | yes | yes (once converged) |
| reaches full precision | sometimes | no (capped) | yes (accumulates) |
| parallel-safe | — | yes | yes (join is order-free) |
| decouples discovery from use | no | no | yes |

### Honest caveats

1. **Warmup is non-deterministic until converged** — same contract as the
   committed lowered tree: deterministic *given* a committed, complete cache.
   A cache miss in CI should hard-fail (force a regen) rather than silently
   fall back to bounded inference.
2. **Invalidation is interprocedural** — keying by IR hash handles edited
   functions; a callee signature change must invalidate its callers. That's the
   fixpoint dependency tracking — real work, well-trodden.
3. **It's a feature, not a tweak** — cache structure, hashing, join, fixpoint
   driver, commit + staleness. The payoff (determinism + full precision + drops
   the budget + makes parallel lowering usable) is large.

## Sequencing

1. Land the keepers from the experiment — **atomic gensym** and **pmapv** — they
   are correct independent of any of this.
2. Land the deterministic naming (build-args inst-id, gensym-suffix strip,
   import sort) with the lowering work; regenerates the tree once.
3. Prototype the type cache: content-key + join-merge + a fixpoint driver;
   measure convergence on `core`.
4. If it converges cheaply, wire lowering to read it, drop the wall-clock
   budget, and turn on `pmapv` lowering for the speedup — now reproducible.
