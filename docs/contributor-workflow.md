---
status: active
last-verified: 2026-09-06
authoritative-for:
  - contributor-workflow
  - ci-gates-map
human-verified:
---

# Contributor workflow — build, test, regenerate, and what CI checks

The commands and gates for working *on* the implementation, in one place. The
README covers the language and embedding API; `contribution-policy.md` covers
direction and design contracts; this page covers the mechanics between a clone
and a green PR. Each section points at the doc that owns the detail.

## Commands

```sh
make build              # bin/lg — runs the smoke suite, then promotes build/lg
make test               # gogen diff smoke + go test -short ./test/... (the .lg test runner)
go test -short ./... -skip TestClojureTestSuite   # the Go suite, as CI runs it
go test ./test/ -run TestClojureTestSuite          # jank conformance (needs the submodule)
go test -run 'TestRunner/coll_test.lg' ./test      # one .lg test file
make lint               # golangci-lint (v2 config), same invocation as CI
make generate           # regenerate every generated artifact after editing pkg/rt/core/**/*.lg
make check-generated    # content-based freshness check; CI's generated-artifacts job
make smoke              # correctness + boot-budget smoke
make ratchets           # bench-ratchet + fanout-ratchet perf gates
git submodule update --init   # once, for TestClojureTestSuite (test/clojure-test-suite)
```

Local hooks mirroring CI: `prek install --install-hooks --hook-type pre-commit
--hook-type pre-push` (or classic `pre-commit install ...`) from
`.pre-commit-config.yaml`. `make install-hooks` registers the git merge drivers
for `core_compiled.lgb` and `generated.sums`, which otherwise conflict on every
rebase that touches `.lg` sources.

## Generated artifacts

Editing `pkg/rt/core/**/*.lg` changes nothing until `make generate` runs: the
runtime loads `pkg/rt/core_compiled.lgb`, `pkg/rt/core_go_lowered/` (under
`-tags gogen_ir`), and the `zz_primitives_generated.go` registrars instead of
the source, and `pkg/rt/generated.sums` is the content digest that decides
freshness. Never hand-edit a generated file. `make check-generated`, the
`generated-artifacts` CI job, and the `genmanifest` tests all catch a missed
regeneration. Full detail, including the merge drivers and why `make build`
cannot be trusted to regenerate: [`regenerating-generated-artifacts.md`](regenerating-generated-artifacts.md).

## Where things go

- `.lg` tests live under `test/`; Go tests next to the package they test.
  Nothing at the repository root: the `test-location` job rejects root
  `*_test.go` and `*_test.lg` (`scripts/check_test_location.py`).
- Every file under `docs/` carries YAML frontmatter (`status`, `last-verified`,
  `authoritative-for`, ...). Schema in [`README.md`](README.md); the maintenance
  hook and the `docs-frontmatter` job in [`frontmatter-hook.md`](frontmatter-hook.md).
- Natives that touch syscalls, the terminal, or JS mirror the `_linux.go` +
  `_other.go` (or `_wasm.go` + `_other.go`) split so every `GOOS` builds. CI
  builds `wasip1`, TinyGo WASI, and `-tags lg_no_http` in addition to the host.

## Facts you cannot infer from one file

- **Native registration has two paths.** Hand `ns.Def(...)` calls in
  `pkg/rt/lang.go`, and `//lg:native` markers that `cmd/lgprimgen` turns into
  `zz_primitives_generated.go` during `make generate`. A name registered both
  ways shadows one of them; `pkg/rt/shadowed_registration_test.go` fails on it.
  The typed constructors are in `pkg/vm/native_func.go` (`NewCtxNativeFn`,
  `NewArityNativeFn`).
- **Three compile paths.** The default is the direct compiler in
  `pkg/compiler`. Binding `*ir-compile*` routes forms through the IR in
  `pkg/rt/core/ir/` (written in let-go) and falls back to the direct compiler
  when lowering fails; `-tags gogen_ir` runs the Go-lowered tree instead of
  bytecode. A change to one path needs a stated position on the other two.
  Knobs: [`design/ir-dynamic-vars.md`](design/ir-dynamic-vars.md); the backend:
  [`design/go-aot-backend.md`](design/go-aot-backend.md).
- `*compiling-aot*` is true during `-c`/`-b`/`-w` and false at runtime;
  `*in-wasm*` is true in WASM builds. Guard side effects with them.
- Performance claims go through the ratchet, not a same-machine single run:
  [`perf/ratchet.md`](perf/ratchet.md). The PR A/B workflow runs only when the
  PR carries the `perf` label.

## CI gates (`.github/workflows/go.yml`)

| Job | Guards | Local equivalent |
|---|---|---|
| `lint` | golangci-lint | `make lint` |
| `test-location` | no test files at the repo root | `python3 scripts/check_test_location.py` |
| `build` | manifest fresh, regenerate, build, Go + `.lg` tests with the bundled stdlib and again with `-tags bootstrap`, jank suite both ways, lowering e2e, native-entry matrix | `make check-generated-manifest`, `go test ./...`, `go test -tags bootstrap ./...` |
| `race` | `-race` on `pkg/vm` and `pkg/rt` | `go test -race -short ./pkg/vm/... ./pkg/rt/...` |
| `default-deps` | untagged builds link no tag-gated subsystem | `make check-default-deps` |
| `no-http-build` | `cmd/lg-runtime` builds and boots with `-tags lg_no_http` and without `net/http` linked | `go build -tags lg_no_http ./cmd/lg-runtime` |
| `wasip1-build` | `GOOS=wasip1 GOARCH=wasm` builds | same |
| `tinygo-wasi-build` | runtime-only builds under TinyGo and boots in wasmtime | `tinygo build -target=wasi ./cmd/lg-runtime` |
| `gold-differential` | goldens re-derived from real Clojure match | cached on the Clojure version; runs on cache miss or dispatch |
| `gogen-diff` | Go lowering output matches the committed tree | `make gogen-diff` |
| `generated-artifacts` | every generated artifact in lockstep with source | `make check-generated` |
| `docs-frontmatter` | changed docs have well-formed frontmatter (PRs only) | `python3 scripts/docs_frontmatter_hook.py --check docs/...` |
| `docs-status` | judgement-layer docs report when a PR touches docs (report only) | `python3 scripts/docs_status.py` |

Perf workflows (`perf-pr.yml` and siblings) are opt-in via the `perf` label or
manual dispatch; see [`perf/ratchet.md`](perf/ratchet.md).

## Before opening a PR

- Read [`contribution-policy.md`](contribution-policy.md) §1 if the change
  touches a design contract (portable interpreter, embeddability, VM stays
  first-class, AOT is additive). Changes against a contract need explicit
  project-owner sign-off in the PR thread, not just green CI.
- Run `make generate` and commit the regenerated files with the `.lg` edit that
  caused them, in the same commit.
- The PR body states why before what, describes the net change, and says how
  it was verified. Titles become squash-commit subjects.
