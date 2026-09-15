---
status: active
last-verified: 2026-09-13
---

# Git-Diff Touched-Form View Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a Git-diff-selected, form-expanded local view beside the existing whole-corpus PR delta.

**Architecture:** `quality.diff` owns Git commands and hunk/path parsing. `quality.touched` owns source-form expansion, identity pairing, and local metrics. `quality` orchestrates both views; `quality.report` prints them separately. Go AST knowledge stays in `cmd/go-callables`. The plain-Git path must not invoke jj or assume a workspace directory.

**Tech Stack:** let-go `.lg`, Go standard library `go/parser`, Git CLI, Go test harness (`TestRunner/quality_*`).

**Spec:** `docs/superpowers/specs/2026-09-13-git-diff-touched-forms-design.md`.

**Workspace:** `.workspaces/pr871-review-fix`, based on published PR #871 head `4557c532`; preserve its divergent local bookmark. A detached plain-Git verification worktree is available at `.workspaces/pr871-git-verify`. Run Go through mise 1.26.5; do not run a bare `go build` for a command that writes a binary at repo root.

---

## Chunk 1: Git input and hunk contract

### Task 1: Parse Git's exact two-sided diff without path guessing

**Files:** Create `scripts/quality/diff.lg`; create `test/quality_diff_test.lg`.

- [ ] Write table tests for `parse-name-status-z` with `M\0path\0`, `A\0path\0`, `D\0path\0`, `T\0path\0`, and `R100\0old\0new\0`; include spaces, tabs, Unicode, and a newline in a path. Assert the two path fields are preserved byte-for-byte; a type-only change remains a path entry.
- [ ] Run `rtk proxy mise exec -- go test ./test -run 'TestRunner/quality_diff_test' -count=1`; expect a missing-namespace/function failure.
- [ ] Implement `parse-name-status-z` as a NUL-token state machine. Reject truncated rename records and unknown statuses with `ex-info`; return ordered `{:status :old-path :new-path}` records.
- [ ] Write tests for `parse-patch-hunks` using `diff --git` sections, `@@ -a,b +c,d @@`, omitted counts (=1), and zero counts (=anchor). A rename-only/binary section has no hunks but retains its status. Reject a section count/order mismatch against name-status records.
- [ ] Implement `git-changes(base, head?)`: first resolve each supplied rev with `git rev-parse --verify --end-of-options <rev>^{commit}` and require a full hex object ID, so a rev beginning with `-` cannot become a Git option. Use `git diff --no-ext-diff --no-textconv --no-color --find-renames --unified=0 base-sha head-sha --` for two revisions, and `... base-sha --` for a working-tree head. Run matching `--name-status -z` with identical endpoints/options. Never interpolate revs into a shell string. Capture stderr and exit code as explicit errors. Return ordered path records with old/new line spans.
- [ ] Run the targeted test again; expect PASS. Commit `feat(quality): parse Git diff paths and hunks` in the #871 worktree.

### Task 2: Pin old/new content and edge statuses

**Files:** Modify `scripts/quality/diff.lg`; modify `test/quality_diff_test.lg`.

- [ ] Add an integration test in a temporary plain-Git fixture for a one-line edit, added file, deleted file, rename-only change, type-only change, and binary change. Assert `git-changes` emits exact paths/statuses and text content only where appropriate; binary/type-only entries are path-only. A working-tree case includes staged and unstaged tracked edits but reports untracked paths as omitted. A leading-dash revision is rejected before diff execution.
- [ ] Run the targeted test; expect failure before content loading.
- [ ] Implement `revision-text(rev,path)` via `git show rev:path`; use `slurp` only for the omitted-head working-tree side. Do not request content for a side that does not exist or is binary/type-only. Preserve empty-file content as `""`, distinct from a missing side. A failed Git call is an explicit error.
- [ ] Run the targeted test; expect PASS. Commit `feat(quality): load only changed Git file versions`.

## Chunk 2: Form expansion and local metrics

### Task 3: Expand `.lg` hunk spans and pair definitions

**Files:** Create `scripts/quality/touched.lg`; modify `scripts/quality/forms.lg`; create `test/quality_touched_test.lg`.

- [ ] Test a one-line edit inside a long definition, two hunks in one definition, one hunk crossing two definitions, a first/last-line edit with a zero-length counterpart anchor, and a comment-only edit outside forms. Add added/deleted definitions, duplicate/missing identities, and two `defmethod` forms with distinct dispatch values. Assert complete top-level spans, de-duplication, unique pairing only where proven, and explicit outside-form/unpaired entries.
- [ ] Run `rtk proxy mise exec -- go test ./test -run 'TestRunner/quality_touched_test' -count=1`; expect failure.
- [ ] Add a reader-validation helper to `quality.forms` that parses the whole `.lg` text before trusting `segments`. In `quality.touched`, build old/new indexes from `segments`; select intersecting spans, then look up a unique counterpart across the full opposite index by namespace/name/kind (and dispatch identity for `defmethod`). Strictly internal zero-length anchors may select a form; boundary anchors may not. Failed reading/segmentation returns an explicit file-level fallback reason.
- [ ] Reuse `quality.metrics` for callable `cc`/`sloc`; leave these nil for non-callable forms. Emit `:delta` only for uniquely paired metrics, never fabricate zero for added/deleted/unpaired entries. Sort deterministically by new path, old path, span.
- [ ] Run the targeted test; expect PASS. Commit `feat(quality): expand Lisp diffs to touched definitions`.

### Task 4: Expand Go hunks with AST-backed declaration spans

**Files:** Modify `cmd/go-callables/main.go` and its Go tests; modify `scripts/quality/touched.lg`; modify `test/quality_touched_test.lg`.

- [ ] Add Go tests for changed lines within a multiline function, adjacent declarations, method receiver identity, malformed Go alongside a valid Go file, and a non-callable top-level declaration. Require per-file `:errors` records so one bad file cannot be mistaken for an absent declaration when another parses.
- [ ] Run `rtk proxy mise exec -- go test ./cmd/go-callables -count=1`; expect failure.
- [ ] Extend `go-callables` with a declaration-span output (name/kind/start/end, plus existing callable `cc`/`sloc`) and `:errors [{:path :message}]` for every failed file, even when other files parse. Analyze only Git-diff-selected Go files: materialize their `git show`/working-tree text in a uniquely created temporary directory, run the Go tool once per side, and remove that exact temporary directory in `finally`. Remove the existing directory-name-based `.workspaces` skip in `goFiles` and rely on nested `.git`/`.jj` markers. Do not make a full revision checkout or assume a workspace directory name.
- [ ] In `quality.touched`, use these AST spans for selection and unique identity pairing; map every `:errors` record to an explicit file-level fallback. Re-run Go and `.lg` touched tests; expect PASS. Commit `feat(quality): expand Go diffs with AST spans`.

## Chunk 3: Report integration and verification

### Task 5: Put two views in one PR report without mixing their scores

**Files:** Modify `scripts/quality.lg`; modify `scripts/quality/report.lg`; modify `test/quality_delta_test.lg`; modify `test/quality_cli_test.lg`.

- [ ] Add tests asserting `q/delta` has separate `:touched` and `:debt-delta`; text says `whole corpus` and `touched forms`, prints signed `cc`/`sloc` changes only when paired, explicitly says untracked files are omitted for a working-tree head, and never adds touched metrics to composite debt. Rename-only/non-source changes remain path entries. For an explicit head, assert both views carry the same resolved base/head commit IDs; for an omitted head, assert both say `working tree` rather than inventing a head commit ID.
- [ ] Run `rtk proxy mise exec -- go test ./test -run 'TestRunner/quality_(delta|cli)_test' -count=1`; expect failure.
- [ ] At `delta` entry, resolve each supplied rev once to a full Git commit ID (only consult jj when the checkout has `.jj` and the input is a jj revset), then pass that exact ID pair to both `extract-rev`/corpus measurement and `quality.diff/git-changes`. An omitted head means working-tree content for both views. Call `quality.touched/build` on Git changes, attach `:touched` to the result, and keep `compare-results` a pure whole-corpus function. Add a distinct section in `quality.report/delta-text` after the corpus sections. Git diff errors fail the delta explicitly; parse failures are labelled file-level fallbacks.
- [ ] Run the targeted tests; expect PASS. Commit `feat(quality): report touched forms beside corpus delta`.

### Task 6: Prove plain-Git operation and scope limits

**Files:** Modify `test/quality_diff_test.lg`, `test/quality_touched_test.lg`, and docs if an observed caveat requires it.

- [ ] If testing in the local jj worktree, set `LG_QUALITY_REVISION_ROOT` to an explicit scratch parent; this is optional developer coverage, never a prerequisite for CI or the required plain-Git result.
- [ ] Snapshot the local jj commit. Verify `git cat-file -t <sha>` says `commit` and `git ls-tree -r <sha>` contains expected changed files; update only a clean detached plain-Git verification worktree to that SHA and compare an expected source file's content hash to the jj snapshot. If Git cannot inspect the object or checkout content differs, stop and report verification blocked; a Git ref cannot repair a missing object. Run `rtk proxy mise exec -- go test ./test -run 'TestRunner/quality_(diff|touched|delta|cli|inputs)_test' -count=1` there with no `.jj` and no `LG_QUALITY_REVISION_ROOT`.
- [ ] Include a one-commit shallow plain-Git fixture: `git rev-parse HEAD` must still yield a full SHA while the pre-date history boundary may be empty. Verify the touched-form Git diff accepts that SHA and reports a limited-history error only when the requested base is unavailable.
- [ ] Run `rtk proxy mise exec -- go test ./cmd/go-callables -count=1`, `rtk proxy mise exec -- go vet ./...`, and `rtk proxy mise exec -- go build ./...` in that plain-Git worktree. Confirm `scripts/check_mise_go_sync.py` passes. Record full-suite timeouts as limits, not success.
- [ ] Inspect `jj diff` and `git status` for unrelated changes; update `.agent/memory/working/WORKSPACE.md`. Do not move bookmarks, push, change CI config, or dispatch the #863 benchmark without the separate PR-stack decision.
