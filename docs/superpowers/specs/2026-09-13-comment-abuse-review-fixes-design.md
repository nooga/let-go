---
status: active
last-verified: 2026-09-13
---

# Comment-abuse review fixes

## Scope

Repair the two confirmed blockers in PR #870 without changing unrelated lint rules or the stacked PR #871. Keep lint report-only by default.

## R4: Go definition spans

The current delimiter scan ends at the closing `)` of a Go function signature. Keep the existing Lisp form scan. Add one small Go-standard-library helper that parses all requested Go files in a single invocation and writes `path<TAB>start-line<TAB>end-line` for each `*ast.FuncDecl` with a body; `lint.lg` consumes those spans when building R4 definitions. Bodyless declarations (e.g. assembly implementations) have no span and are skipped. This keeps Go syntax ownership in `go/parser` and avoids one subprocess per file. A parse or helper failure is an explicit lint error, never a silent zero-definition result. Count declaration-to-closing-brace spans as today, excluding leading documentation. Cover multiline signatures and multiple functions in one file.

## R6: churn measurement versus gate

R6 reports a baseline-relative multiple for every range with added code, including zero added comments. No calibrated abuse threshold exists in this branch. Preserve the multiple as a measurement, but make `comment-churn` non-gateable rather than inventing a threshold. If the corpus baseline is zero, report an undefined multiple (`:measure nil`) with explicit evidence instead of dividing by zero. `--gate comment-churn` must exit zero with an explicit ignored-kind note that describes R6 as a measurement, not as the R5 heuristic. Cover zero-comment, positive-comment, zero-baseline, and missing-`--churn` cases.

## Verification

Add failing Go harness tests for both cases before production changes. Run targeted `TestLint`, then `go build ./...` and the PR's standard tests. Verify the #870 layer alone before restacking #871.
