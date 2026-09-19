---
status: active
last-verified: 2026-09-13
---

# Quality Gitlink Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Report changed Gitlinks as path-only instead of aborting quality deltas.

**Architecture:** Extend patch metadata classification for mode `160000`; the text loader skips classified Gitlinks. Existing touched-form handling already emits path-only entries for rows without text.

**Tech Stack:** let-go, Go test runner, Git CLI.

---

## Chunk 1: Gitlink regression and fix

### Task 1: Reproduce and fix the loader failure

**Files:**
- Modify: `test/quality_diff_test.lg`
- Modify: `scripts/quality/diff.lg`

- [x] Add a temporary plain-Git repo test with a mode-`160000` entry and a local nested repo.
- [x] Assert commit and dirty-working-tree deltas retain the path without text or form metrics.
- [x] Confirm both regressions fail with `git show <sha>:<gitlink>` on the baseline.
- [x] Mark mode `160000` and mode-less dirty-submodule patch sections; skip text loading.
- [x] Add and pass a negative test for regular text resembling a submodule hunk.
- [x] Run focused checks, `go test -short -timeout 300s -count=1 ./...`, build, vet, generated-manifest, and diff checks.
- [ ] Commit only the two code files and these design/plan docs; do not push.
