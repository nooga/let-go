---
status: active
last-verified: 2026-06-13
human-verified:
---

# Regenerating generated artifacts after `.lg` edits

Editing any `pkg/rt/core/**/*.lg` file requires regenerating **two** artifacts
that the runtime loads *instead of* the `.lg` source. Skip the regen and your
edits silently have no effect:

| Artifact | Loaded by |
|---|---|
| `pkg/rt/core_compiled.lgb` | the bytecode-VM path (default) |
| `pkg/rt/core_go_lowered/` (a tree of `.go` files) | the `-tags gogen_ir` path |

## Do this

```sh
make generate        # regenerates BOTH artifacts + refreshes the manifest
make check-generated # verify they are in sync with sources (content-based)
```

`make generate` also rewrites `pkg/rt/generated.sums`, a content digest of every
`.lg` + `cmd/lgbgen` source. That manifest is the source of truth for staleness.

## Why not `make build` / `make test`?

Those targets — and the older `check-bundle-fresh` / `check-lowered-fresh`
targets — gate on file **modification times**. mtimes are unreliable after a
`git` or `jj` checkout: VCS tools write arbitrary mtimes, so a stale bundle can
look *newer* than the sources that should have rebuilt it. This is the
long-standing "`make build` didn't actually regenerate" footgun.

`make check-generated` compares a sha256 of the sources against the digest in
`pkg/rt/generated.sums` instead, which is checkout-independent.

## Where staleness is caught

- **`go test ./...` / CI** — `TestGeneratedArtifactsAreFresh` in
  `pkg/genmanifest` fails when a source changed without `make generate`.
- **`make check-generated`** — same check as a CLI (`cmd/check-generated`).
- **git pre-commit hook** — `scripts/pre-commit` blocks a commit with stale
  artifacts once it's installed as `.git/hooks/pre-commit`. No make target wires
  it up; symlink it by hand. **Note:** `jj` does not run git hooks; jj users rely
  on the test + `make check-generated`.

## Git merge driver for `core_compiled.lgb`

`pkg/rt/core_compiled.lgb` is a binary bundle regenerated from the embedded
`.lg` sources. Git cannot meaningfully merge this binary on rebase, so we ship a
custom merge driver that regenerates it from sources *after* the `.lg` files
have been merged as text.

```sh
make install-hooks
```

`make install-hooks` registers this merge driver (a `git config merge.lgb.*`
pair). The config lives in `.git/config`, which is not shared, so each clone
needs the registration once. After it, rebases and merges that touch any embedded
`.lg` source regenerate the `.lgb` automatically — no binary merge conflicts when
stacking PRs that edit `core.lg` and friends.

## `go build` cannot regenerate

Go has no build-time codegen hook, so `go build` never regenerates. Only
`make generate` or `go generate ./cmd/lgbgen` do. The `//go:generate` directives
live in `cmd/lgbgen/generate.go`.

## Implementation

- `pkg/genmanifest` — source-hashing + staleness comparison (the single source
  of truth, shared by the test and the CLI).
- `cmd/lgbgen` — writes `pkg/rt/generated.sums` on every regen (both the bundle
  and `--target=go` paths).
- `cmd/check-generated` — the CLI used by the Makefile target and the hook.

## Bundle-only regen

If you are certain you are not touching the Go-lowered path, the bundle alone is
`go run -tags bootstrap ./cmd/lgbgen`. Prefer `make generate` regardless, so the
two artifacts never drift apart (the `-tags gogen_ir` path silently diverges
from the untagged path when only one is regenerated).
