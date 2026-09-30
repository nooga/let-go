---
status: active
last-verified: 2026-09-22
authoritative-for:
  - test-runner-narrowing
  - test-parameterization
supersedes: []
human-verified:
---

# Test Runner Narrowing & Parameterized Execution

This specification defines the unified, variable-driven test execution interface for `let-go`. It standardizes test invocation across local development, CI workflows, and autonomous agents through Make and environment variables, eliminating the proliferation of bespoke single-purpose Make targets.

## Motivation & Historical Context

Between PR #249 (Feb 2026) and PRs #681, #729, and #734 (Aug 12–14, 2026), four separate PRs added bespoke Make targets (`gogen-diff`, `engine-parity-gate`, `native-entry-gate`, `test-gogen-diff-gate`, `parity-gate-phase1`, `strict-audit`, `shadow-warning-test`) rather than parameterizing `make test`.

This caused several operational issues:
1. **Developer Latency**: `make test` hardcoded the execution of the entire test corpus (~45 seconds) and ran the `test-gogen-diff-short.sh` probe on every invocation. Running an individual test file (e.g. `test/coll_test.lg`) required dropping out of Make and manually crafting subtest regexes: `go test -run 'TestRunner/coll_test.lg' ./test`.
2. **Package Scoping**: Running unit tests for a specific Go package (e.g. `pkg/rt` or `cmd/go-callables`) lacked a Make interface.
3. **Fragile Coupled Gates**: `make native-entry-gate` coupled `TestNativeEntryASTGate` with `TestJankSuiteDirectABIGeneratedGo`. In JJ workspaces and secondary worktrees lacking the `test/clojure-test-suite` git submodule, the entire gate failed loudly even when an engineer only needed to verify lowered Go ASTs.
4. **Target Bloat**: Multiple aliases and duplicate targets accumulated in `Makefile`.

## Specification

### 1. Unified `make test` Interface

`make test` is the canonical entry point for all test execution. It accepts the following Make / environment variables:

| Variable | Default | Purpose | Example |
| :--- | :--- | :--- | :--- |
| `RUN` (or `TEST`) | _(empty)_ | Test name or regex pattern passed to `go test -run` | `make test RUN=coll`<br>`make test RUN=TestNativeEntryASTGate` |
| `PKG` | `./test/...` | Target Go package path | `make test PKG=./pkg/rt`<br>`make test PKG=./test/e2e` |
| `SHORT` | `true` | Boolean toggle for `-short` (`false` or `0` enables e2e) | `make test SHORT=false PKG=./test/e2e RUN=TestNativeEntryASTGate` |
| `LG_TEST` | _(empty)_ | Fast discovery filter in `test/language_test.go` | `make test LG_TEST=quality_delta` |
| `TAGS` | _(empty)_ | Build tags passed via `-tags` | `make test TAGS=gogen_ir` |
| `COUNT` | `1` | Test repetition count passed via `-count` | `make test COUNT=5` |

### 2. Fast-Path Optimization

When `RUN`, `PKG`, or `LG_TEST` is supplied, or `COUNT` is not `1`, `make test` automatically skips the `scripts/test-gogen-diff-short.sh` probe. The probe is retained exclusively for un-narrowed default runs (`make test`), ensuring that local feedback loops complete in sub-second to low-second times.

### 3. Let-Go Harness File Filtering (`LG_TEST`)

`test/language_test.go` reads `LG_TEST` during `filepath.Walk`. When non-empty, files whose names do not contain the substring are skipped before subtest allocation. This avoids walking and constructing subtests for 120+ irrelevant `.lg` test files.

### 4. CI Modernization & Deprecation Bridges

CI workflows invoke the canonical variable-driven syntax:
```yaml
- name: Differential engine gate
  run: make test SHORT=false PKG=./test/e2e RUN=TestGogenAOTDiff
```

To preserve backward compatibility with external scripts, documentation, and developer muscle memory:
- `gogen-diff`, `engine-parity-gate`, `native-entry-gate`, and `test-gogen-diff-gate` are maintained in `Makefile`.
- When invoked, each deprecated target prints an informational notice to stderr directing the caller to the canonical `make test` syntax, and delegates execution to `go test` with `-short=false`.
- `scripts/test-gogen-diff-short.sh` continues to verify that these deprecated entry points enforce `-short=false`.
