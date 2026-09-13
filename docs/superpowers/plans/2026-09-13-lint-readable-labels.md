---
status: active
last-verified: 2026-09-13
---

# Reader-Friendly Lint Labels Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make every user-facing lint finding and section understandable without decoding `R1`–`R6` or knowing a jargon-heavy rule name.

**Architecture:** Keep `:kind` and `--gate` as stable machine identifiers. A display-name map in `scripts/lint.lg` covers built-in comment kinds; code-rule catalog entries carry their own `:label`, preserving the existing data-driven extension point. An omitted label falls back to space-separated words from `:kind` for older/custom catalogs. EDN keeps existing kinds. Section summaries and skip/gate notes use the same display vocabulary. No detection threshold or gate behavior changes.

**Tech Stack:** let-go `.lg`, Go test harness `test/lint_lg_test.go`.

**Spec:** `docs/superpowers/specs/2026-09-13-readable-lint-labels-design.md` in this worktree.

**Workspace:** Use the already-prepared isolated #870 worktree (`jj workspace root` confirms the current root); do not edit the primary checkout or move `lint/comment-abuse-rules` before stack review.

---

## Chunk 1: Text display contract

### Task 1: Cover every user-visible kind with a low-context label

**Files:** Modify `scripts/lint.lg`, `scripts/lint-code-rules.edn`, and `test/lint_lg_test.go`.

- [ ] Add a table test for the complete built-in kind set: `commented-out-code`, `restatement`, `comment-divider`, `duplicated-comment`, `comment-density-outlier`, `devlog-comment`, `comment-churn`, and every `:kind` in `scripts/lint-code-rules.edn`. Assert every label is nonblank, contains no `R[1-6]` or raw hyphenated token, and no two distinct built-in kinds accidentally share one label. Assert jargon terms such as “eta expansion” are absent.
- [ ] Run `rtk proxy mise exec -- go test ./test -run 'TestLint.*Display' -count=1`; expect failure before the map exists.
- [ ] Add one `display-name` function and built-in map. Use explicit labels for opaque comment kinds: `restatement` → `comment repeats code`, `comment-divider` → `decorative section divider`, `comment-density-outlier` → `unusually comment-heavy definition`, `devlog-comment` → `development-note phrase`. Add `:label` to every current code-rule catalog entry, including `eta-expansion` → `argument-forwarding wrapper` and `composable-accessor` → `nested sequence accessor`. If a custom catalog entry omits `:label`, fall back to replacing hyphens in its `:kind` with spaces; keep `TestLintCodeVerbosityCatalogIsData` working without editing `lint.lg` for the fabricated rule.
- [ ] Run the display test; expect PASS. Commit `feat(lint): give findings reader-friendly labels`.

### Task 2: Apply labels consistently to findings, summaries, and notes

**Files:** Modify `scripts/lint.lg`; modify `test/lint_lg_test.go`.

- [ ] Add a fixture-output branch matrix: density skipped and calibrated; churn skipped, empty, finding, and zero-comment-baseline; ignored-gate notes for churn/devlog; comment findings and code-verbosity findings. Assert text includes display labels at finding and section/skip sites, has no `R1`–`R6` headings or instructions to inspect source comments, and still explains evidence. Migrate existing tests that search `[devlog-comment]`, `[restatement]`, or other raw text labels to display text or EDN kinds. Assert EDN `:kind` values and `--gate` tokens/exit behavior are unchanged.
- [ ] Run targeted `TestLint` tests; expect the old R-number/jargon output to fail the new assertions.
- [ ] Replace text-only formatting in `lint.lg` with `display-name`; remove raw kind prefixes from code-rule `:evidence` when the label already appears next to it. Keep `:kind` fields and gateable-kind set unchanged. Use plain wording for skipped, report-only, and no-calibration cases.
- [ ] Run `rtk proxy mise exec -- go test ./test -run TestLint -count=1` and `rtk proxy mise exec -- go test -tags bootstrap ./test -run TestLint -count=1`; expect PASS. Run `rtk proxy mise exec -- go build ./...`. Commit `fix(lint): use plain-language report headings`.

## Chunk 2: Verification and stack boundary

### Task 3: Audit final text and compatibility

**Files:** Modify `test/lint_lg_test.go` only if the audit finds a missing assertion.

- [ ] Search user-facing `println` and finding formatters in `scripts/lint.lg` for `R[1-6]`, bare `name (:kind ...)`, and “see section comment”. Distinguish internal comments from emitted text; require zero unexplained user-facing IDs.
- [ ] Run the targeted lint suite and `rtk proxy mise exec -- go vet ./...` with the mise Go 1.26.5 toolchain. Record the previously observed local full-e2e timeout as unresolved; do not attribute it to the base without a control run or claim full-suite success.
- [ ] Snapshot the local fix commit and create a clean detached plain-Git verification worktree for that commit. Run `rtk proxy mise exec -- go test ./test -run TestLint -count=1` there, then inspect `rtk git status --short` and `rtk git diff 16a637dded154e9e712f63731b8baf5d0e6ec4c4 HEAD --` in that Git worktree. Update `.agent/memory/working/WORKSPACE.md` with exact commands/outcomes. Do not restack or push before user review of the local stack.
