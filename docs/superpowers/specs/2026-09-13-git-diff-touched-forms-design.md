---
status: active
last-verified: 2026-09-13
---

# Git-diff touched-form view and readable lint labels

Date: 2026-09-13. Applies to the local #871 scorer and its #870 lint parent. No PR update is authorized by this design.

## Purpose and invariant

The PR report keeps its existing whole-corpus base/head delta. It adds a separate touched-form view so a small overall movement cannot hide a concentrated regression in edited code. The two views are never summed or presented as interchangeable: coverage, dead-code reachability, duplication, and other corpus terms can move outside the edited forms.

## Inputs and boundaries

For two revisions, `git diff --no-ext-diff --no-color --find-renames --unified=0 <base> <head> --` selects tracked changed lines. With no explicit head, use the same command with `<base> --` so staged and unstaged tracked working-tree edits are included. An untracked file is out of scope and the report says so. A matching `git diff --name-status -z --find-renames` call supplies NUL-delimited old/new paths and statuses; patch headers/hunk order are validated against it rather than trusting whitespace- or quote-splitting. A rename-only or binary change remains a path-level entry with no form metrics. Revision contents come from `git show <rev>:<path>` or the working tree, without a jj requirement. The same exact base/head pair drives both report views; no merge-base/triple-dot substitution occurs.

Parse hunk coordinates independently on the old and new side. Include additions, deletions, and renames. A nonzero changed range selects every enclosing top-level unit whose line span intersects it; one hunk can therefore select multiple forms. On the zero-length side of an insertion/deletion, select an old/new form only if the anchor is strictly inside that form; an anchor between adjacent forms selects neither. Selection is only the trigger: build full old/new form indexes and look up a uniquely matching counterpart by identity even when that side's hunk span selected nothing (e.g. a first-line edit or insertion at a form boundary). For `.lg`, first validate the whole file with the language reader, then expand selected lines to enclosing top-level definition/form using `quality.forms/segments`. A failed read or inconsistent segmentation forces a file-level fallback with a reason. For Go, use `go/parser`-derived declaration/function spans, not braces or regular expressions; a Go parse failure likewise forces a file-level fallback with a reason. Choose the enclosing declaration as the scoring unit; a nested expression edit belongs to its containing definition. Deduplicate overlapping hunks that touch one unit. Report edits outside any form/declaration as an explicit file-level span. Non-source files appear as changed paths but have no form score.

## Output and scoring

The touched-form view lists old/new path, old/new line span, change status, and two local structural measurements where defined: cyclomatic complexity (`cc`) and source lines of code (`sloc`). These reuse the scorer's existing callable metrics; a top-level non-callable form has a span but no invented complexity. Pair old/new declarations only on a unique stable identity (qualified `.lg` definition name/kind, including dispatch identity for `defmethod`; Go receiver/name or top-level declaration name) after applying the Git rename path mapping. If identity is absent, duplicated, or maps many-to-many, show the old and new units as unpaired entries with no signed metric delta. Added/deleted forms have a missing side, not an invented zero. It does not reuse the whole-corpus composite debt weights or assert local coverage/dead-code/duplication deltas. The existing whole-corpus result remains authoritative for those terms.

Git command failures, malformed diffs, and unavailable revisions are explicit errors. Deterministic ordering is by new path, old path, and span. No fixed `.workspaces` path or `jj` executable is required in a plain-Git checkout.

## Reader-facing lint names

Keep EDN `:kind` values and `--gate` tokens stable. Map them to plain-language text labels in CLI output and replace `R1`–`R6` in headings/skip notes with descriptive names. The priority renames are: `restatement` → “comment repeats code”; `comment-divider` → “decorative section divider”; `comment-density-outlier` → “unusually comment-heavy definition”; `devlog-comment` → “development-note phrase”; `eta-expansion` → “argument-forwarding wrapper”; `composable-accessor` → “nested sequence accessor”. Other code-rule labels receive the same reader test: can the label identify the issue without reading the source or knowing a rule number? Evidence text remains adjacent and explains why the match matters.

## Tests

- Table tests for one-line edits inside long `.lg` and Go functions, multiple hunks in one form, first/last-line boundary edits with counterpart lookup, adjacent forms and zero-length anchors, additions, deletions, renames (including rename-only), binary files, malformed `.lg`/Go, and edits outside forms.
- A plain-Git test where `jj` is absent; shallow history may leave the historical boundary empty without invalidating a head SHA.
- Golden/text assertions that every emitted lint finding and skip note has a descriptive label; EDN kinds and gate behavior remain unchanged.
- Existing quality/lint suites, `go vet ./...`, `go build ./...`; report unrelated full-suite timeouts without claiming green.
