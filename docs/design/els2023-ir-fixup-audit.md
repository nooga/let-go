---
status: active
last-verified: 2026-06-05
authoritative-for:
  - ir-link-fixup-pass-design
human-verified:
---

# ELS 2023 Architecture Audit (IR Branch) and Proposed IR Link/Fixup Pass

Date: 2026-05-23
Source paper reviewed: `els2023_final.pdf` ("Design of an Efficient Lisp Bytecode Machine and Compiler")
Audit basis: isolated `jj` workspace on `lisp-ir-pipeline` (`149fa164`), then copied into default workspace.

## What was verified

- Default workspace is active IR development and has uncommitted changes:
  - `pkg/rt/core/ir/lower.lg`
  - `pkg/rt/core_compiled.lgb`
- Isolated workspace (`.workspaces/els-audit`) was created to avoid interference.
- On `lisp-ir-pipeline`, architecture is already IR-first:
  - `build -> optimize-fn -> lower` pipeline exists.
  - Lowering already has deferred patching infrastructure for control-flow offsets:
    - lowerer state `:patches`
    - `emit-placeholder!`
    - `patch-branches!`
- Test harness still marks `pendingPhase5` cases in `pkg/ir/pipeline_bench_test.go`.

## Interpretation vs ELS paper

The ELS paper’s key practical claim is not "just patch branch offsets", but generalizing fixups to any variable-length emission decision, allowing one-pass optimistic codegen with late resolution.

Current let-go IR lower matches this direction partially:
- Already optimistic + late patching for branch displacement.
- Not yet generalized for optional instruction-group insertion/removal (e.g., closure-cell indirection decisions).

So the architecture direction is legitimate. The main gap is generalizing the existing patch mechanism into a formal IR link/fixup phase.

## What should happen first

1. Finish Phase-5 IR-lower stack/branch correctness.
2. Then expand fixups from branch offsets to generalized emission decisions.

Reason: adding generalized fixups before lower correctness stabilizes increases debugging surface and masks root-cause issues.

## Proposed pass: `ir.link` (or `lower/link!` subphase)

## Placement

Run after `optimize-fn` and after initial lowering emission, before final chunk return.

Pipeline shape:

`expand -> build -> optimize-fn -> lower(emit placeholders + records) -> link(resolve fixups) -> final chunk`

## Core design

Use a unified fixup record format in lowerer state.

Suggested lowerer additions:

- `:labels` — map of logical label id -> concrete IP
- `:fixups` — vector of fixup records
- `:decisions` — map of symbolic decision keys -> resolved policy/value

### Fixup record schema (suggested)

```clojure
{:kind        :branch-offset | :emit-optional | :replace-op | :const-index
 :site-ip      <int>          ; instruction ip or arg ip to patch
 :site-op      <kw>           ; optional sanity check
 :target       <label-or-id>  ; for branch/targeted fixups
 :width        :i32           ; future-proof if encoding evolves
 :best-size    <int>          ; optimistic assumed size
 :worst-size   <int>          ; optional bound
 :decision-key <kw-or-vector> ; for non-branch fixups
 :payload      <map>}         ; kind-specific metadata
```

For current code, most entries are `:branch-offset` and map directly from existing `:patches` usage.

## Resolution algorithm

1. **Emit phase (optimistic):**
   - Emit minimal/common-case instructions.
   - Record fixups instead of fully resolving dependent bytes.

2. **Decision phase:**
   - Compute `:decisions` (e.g., which vars/tags need indirection/cell behavior).
   - This can consult IR metadata and closure/mutation analysis artifacts.

3. **Layout phase:**
   - Recompute effective instruction positions if optional groups change size.
   - Update label -> IP mapping.

4. **Patch phase:**
   - Apply fixups in deterministic order.
   - Validate every fixup target and resulting displacement.

5. **Validation phase:**
   - Ensure no unresolved fixups.
   - Assert stack/source-map invariants still hold.

## Minimal incremental implementation plan

### Step 1: Normalize current branch patches into generic `:fixups`

- Replace ad-hoc `:patches` record shape with `:fixups` + `:kind :branch-offset`.
- Keep behavior identical.

### Step 2: Introduce explicit label table

- Materialize block/branch labels into `:labels` rather than deriving only from `:block-ip` lookups.
- Keep existing branch results unchanged.

### Step 3: Add one non-branch fixup kind

- Add `:emit-optional` as a no-op scaffolding kind first.
- Exercise it in tests without changing generated semantics.

### Step 4: Enable first real optional emission decision

- Candidate: closure-related indirection insertion/elision at lowering boundary.
- Gate behind feature flag to compare chunks and perf.

### Step 5: Expand tests

- Keep current round-trip tests.
- Add fixup-focused tests:
  - forward + backward jumps with nearby optional emissions
  - nested-if join blocks
  - loop back-edge after optional insertions
  - source-map consistency after link

## What should NOT be done yet

- Bytecode encoding overhaul (LONG-prefix / compact operand format) before IR-lower correctness and generalized fixup stability.
- Non-local-exit instruction model rewrite (ENTRY/EXIT family) before current try/catch semantics are fully stable under IR lower.

## Short recommendation

- Continue on `lisp-ir-pipeline` lineage.
- Treat existing lower patching as the seed of a formal generalized link/fixup stage.
- Make it generic in data model now, but conservative in enabled decisions until Phase-5 failures are cleared.
