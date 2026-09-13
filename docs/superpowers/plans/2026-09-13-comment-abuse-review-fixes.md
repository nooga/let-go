---
status: active
last-verified: 2026-09-13
---

# Comment-abuse Review Fixes Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make PR #870's R4 measure Go function bodies and make R6's churn measurement safe and non-gating.

**Architecture:** A single Go-standard-library helper in `scripts/lint_go_spans.go` emits Go function source spans. `scripts/lint.lg` invokes it once per run for all Go files, then reuses the existing R4 ratio logic. R6 remains a report-only measurement, with an explicit undefined result for a zero-comment corpus baseline.

**Tech Stack:** Go `go/parser`/`go/ast`, let-go `.lg`, Go integration tests.

---

## Chunk 1: Go spans

### Task 1: Go parser helper

**Files:** Create `scripts/lint_go_spans.go`, `scripts/lint_go_spans_test.go`.

- [ ] Write `TestLintGoSpans` with a temporary file whose exact contents are:

  ```go
  package fixture

  func First(
      x string,
  ) string {
      _ = "{"
      return x
  }

  func Second() {
      // }
  }

  func Asm()
  ```

  Assert exactly `[]span{{path, 3, 8}, {path, 10, 12}}` (1-based inclusive lines). From the `scripts/` test cwd, run `go run lint_go_spans.go -- <path>` and assert the CLI writes exactly `path\t3\t8\npath\t10\t12\n`.
- [ ] Run `go test ./scripts -run TestLintGoSpans -count=1 -v`; verify RED because the helper is absent.
- [ ] Implement `func spans(paths []string) ([]span, error)` with `parser.ParseFile`, `ast.File.Decls`, `*ast.FuncDecl`, `Body != nil`, and `fset.Position(Pos).Line`/`fset.Position(End()).Line`. The CLI prints one tab-separated `path/start/end` record per span and returns a nonzero exit on parse errors; reject path names containing newline or tab. Strip the leading `--` in `main`: it stops `go run` from interpreting the `.go` path arguments as additional source files.
- [ ] Re-run the targeted test and `go test ./scripts`; verify GREEN.
- [ ] Describe this focused change as `fix(lint): parse Go function body spans`.

### Task 2: Wire R4 into the existing linter

**Files:** Modify `test/lint_lg_test.go`, `scripts/lint.lg`.

- [ ] Add `TestLintR4GoFunctionBodies`: write 20 valid Go functions, each with at least five lines; use 19 with one interior comment line and one with four interior comment lines so the last is the sole outlier. Make the last signature multiline. Assert the R4 output does not say `only 0 definition(s) found` and identifies the 20th function's declaration line. Add `TestLintR4MalformedGoFailsLoudly`: write `package p\nfunc Broken( {\n` to a temp `.go` file, run `lint.lg --edn <file>`, and assert nonzero exit plus a visible `cannot parse Go file` error.
- [ ] Run `go test ./test -run 'TestLintR4(GoFunctionBodies|MalformedGoFailsLoudly)' -count=1 -v`; verify RED at the zero-definition message and the missing parse error.
- [ ] Resolve the helper path from `(os/absolute-path (nth os/args 1))`, replacing the terminal `lint.lg` filename with `lint_go_spans.go`; call `go run <absolute-helper-path> -- <all-Go-paths>` once. Parse tab-separated spans and pass them into R4 alongside the existing `.lg` `file-definitions` path. On a nonzero helper exit or malformed span, throw an explicit error.
- [ ] Re-run the targeted test and all `TestLint` tests; verify GREEN.
- [ ] Describe this focused change as `fix(lint): include Go bodies in R4 density`.

## Chunk 2: R6 measurement

### Task 3: Preserve measurement without an arbitrary gate

**Files:** Modify `test/lint_lg_test.go`, `scripts/lint.lg`.

- [ ] Add `TestLintR6ZeroCommentRangeIsDiagnostic` in a temporary Git repo: commit `package p\n// baseline\nfunc A() {}\n`, then commit a version adding `func B() {}\n` and no comments. The final corpus baseline is positive (1 comment / 3 code lines), while range ratio is zero. Pass `.` as the scan path in every invocation; this fixture has no `pkg` or `scripts` directory. Assert `--churn C1..C2 .` reports `multiple=0`; `--gate comment-churn --churn C1..C2 .` exits zero and names R6 as report-only; `--gate comment-churn .` without `--churn` exits zero, emits the ignored-kind note, and still says R6 was skipped. Add `TestLintR6ZeroBaselineIsUndefined` in a separate Git repo: commit `package p\nfunc A() {}\n`, then commit a version adding `func B() {}\n`; assert `--churn C1..C2 .` says `multiple=undefined` and `--edn --churn C1..C2 .` has `:measure nil`. Preserve the existing positive-churn test.
- [ ] Run `go test ./test -run 'TestLintR6' -count=1 -v`; verify RED on the gate assertion and/or zero-baseline case.
- [ ] Remove `:comment-churn` from `gateable-kinds`. If baseline equals zero, emit an R6 record with `:measure nil` and explicit evidence, avoiding division. Make the ignored-gate note distinguish R6 measurement from R5 heuristic/unknown names. Do not change the R6 default output without `--churn`.
- [ ] Re-run R6 tests and all `TestLint` tests; verify GREEN.
- [ ] Describe this focused change as `fix(lint): keep churn diagnostic non-gating`.

## Chunk 3: Layer verification

### Task 4: Verify #870 independently

- [ ] Run `go build ./...`, `go test ./scripts`, and `go test ./test -run TestLint -count=1`; each must exit 0.
- [ ] Run the selected `go.yml` build-job gates: `go test -short -timeout 60s ./... -skip TestClojureTestSuite`, `go test -tags bootstrap -short -timeout 60s ./... -skip TestClojureTestSuite`, and `make check-generated`; each must exit 0. Record any gate that cannot be run rather than claiming full CI equivalence.
- [ ] Run `python3 scripts/docs_frontmatter_hook.py --check docs/superpowers/specs/2026-09-13-comment-abuse-review-fixes-design.md docs/superpowers/plans/2026-09-13-comment-abuse-review-fixes.md`; it must exit 0. Inspect `jj diff --stat` and `jj status` for unrelated changes.
- [ ] Leave the verified tip in the isolated workspace. Do not move the shared bookmark, push, or update #870 until the user has reviewed the prepared fixes and the #871 CI-trigger choice.
