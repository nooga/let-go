---
status: active
last-verified: 2026-06-05
authoritative-for:
  - docs-index
human-verified: 2026-06-20
---

# let-go docs

Design plans, execution roadmaps, and policy for the let-go implementation. Each doc carries a `status:` frontmatter line indicating whether it's the current authority or has been superseded. See "What's current" below.

## Layout

Docs are bucketed by audience, then by cluster:

- **[`guide/`](guide/)** — user-facing reference for writing and shipping let-go programs: language features, usage, Clojure compatibility.
- **[`design/`](design/)** — contributor-facing architecture: how a subsystem works or was designed (VM, value representation, IR lowering, runtime image, Go AOT, I/O host decoupling).
- **[`perf/`](perf/)** — performance baselines, the regression ratchet, and historical data.
- **`docs/` root** — cross-cutting contributor material that isn't subsystem-scoped: the master plan, roadmaps, contribution policy, and dev workflow (regeneration, frontmatter, testing).

A subdir is earned when a cluster of related docs justifies one; one-off cross-cutting docs stay at root.

## What's current

`contribution-policy.md` is the most authoritative doc on overall direction (self-host as committed direction, `gogen_ir` deployment path, CI gates, callback error contract, `:go-deps` interop schema). Where it disagrees with older docs on direction, it wins.

`master-plan.md` is still the useful phase skeleton (Phases 0–9 covering semantics, VM perf, collections, transducers, runtime images). Read it for the staging structure; cross-reference `contribution-policy.md` for any direction question.

## Topical map

| Concern | Doc(s) |
|---|---|
| Design contracts, CI gates, interop schema | `contribution-policy.md` |
| Phase skeleton, success metrics | `master-plan.md` |
| Calling convention, allocation, TCO | `design/vm-performance-optimization.md` |
| Numeric/value representation | `design/value-representation-and-numeric-performance.md` |
| Persistent collections, seq tower, transients | `clojurelike-refactor-plan.md` |
| Equality/hashing across types | `clojurelike-refactor-plan.md` (Phase 3) |
| Transducers, reduction fast paths | `clojurelike-refactor-plan.md` (Phase 4) |
| Runtime image / stdlib precompile | `design/runtime-image-and-stdlib-cache.md` |
| Go AOT / self-host deployment | `design/go-aot-backend.md` + `contribution-policy.md` §2–3 |
| JVM-shape interop (strategy) | `jvm-compat-plan.md` |
| JVM-shape interop (execution) | `clojure-compat-roadmap.md` |
| Real-world Clojure compat findings | `xsofy-portability-gaps.md` |
| Clojure-test-suite (jank) workflow | `clojure-test-suite.md` |
| Testing strategy, conformance | `testing-and-conformance.md` |
| Perf ratchet, regression checkpoints, historical baselines | `perf/ratchet.md` |
| Babashka pods (usage) | `guide/pods.md` |
| Babashka pods (host protocol / design) | `design/pods.md` |
| Portable `.cljc` / `:lg` reader conditionals | `guide/portability.md` |
| Version requirements, range matching (`let-go.semver`) | `guide/semver.md` |
| `io/resource`, `-resource-paths` / `-source-paths` resolution | `guide/resources-and-source-paths.md` |
| Embedding let-go in a Go program | `guide/embedding-in-go.md` |
| nREPL server + editor setup | `guide/nrepl.md` |
| Clojure compatibility: namespace table + differences | `guide/clojure-compatibility.md` |
| Running, compiling, WASM, project mgmt (lgx) | `guide/usage.md` |
| IR fixup / link pass | `design/els2023-ir-fixup-audit.md` |
| Parallel IR lowering + determinism | `design/parallel-lowering-and-type-cache.md` |
| Runtime I/O, host decoupling | `design/runtime-io-host-decoupling.md` |

## Reading order if starting cold

1. **`contribution-policy.md`** — current commitments, design contracts, interop schema, CI gates.
2. **`master-plan.md`** — the phased skeleton.
3. **`clojurelike-refactor-plan.md`** — semantic foundation under most of the collections/perf work.
4. Pick concerns from the topical map as needed.
5. **`clojure-compat-roadmap.md`** — only if working on Clojure-library compat specifically.

## Status conventions

Every doc carries frontmatter:

```yaml
---
status: planning | active | shipped | superseded | archived
last-verified: YYYY-MM-DD
supersedes: [...]         # optional
superseded-by: [...]      # optional
shipped: [...]            # optional, for partial-shipped docs
remaining-open: [...]     # optional, for partial-shipped docs
authoritative-for: [...]  # optional, used by this index
human-verified:           # optional, set only by explicit human action
---
```

`status` describes the *doc*, not the underlying work. A doc that proposed a feature can stay `active` after the feature ships if it's still the authority on the design rationale; the shipped items move to a `shipped:` list and the doc itself carries forward.

When a feature ships, add a dated `Shipped:` annotation in the relevant section of the doc body, pointing at the commit or PR. The doc stays in place; the status updates. When a newer doc takes over on a topic, link the supersession in both directions via the frontmatter rather than deleting the older doc.

`human-verified` carries the date a human reader last vouched for the doc as current. It's intentionally distinct from `last-verified` (which an automated maintenance pass refreshes on every touch) so a reader can see at a glance whether a human, or just tooling, last attested to the doc. Leave it blank if no human has reviewed the doc yet; populate it only by explicit human action.
