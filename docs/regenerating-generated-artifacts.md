---
status: active
last-verified: 2026-09-10
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
make generate        # regenerates BOTH artifacts + refreshes manifest/digest
make check-generated # verify they are in sync with sources (content-based)
```

`make generate` also refreshes the two bookkeeping files that record what went
into those artifacts:

| File | Contents |
|---|---|
| `pkg/rt/generated.manifest` | one sorted `<output> <input> <kind> <sha256>` record per dependency edge |
| `pkg/rt/generated.sums` | a single digest **of `generated.manifest`** |

The chain is ordered, each step hashing the previous one's output:

```
.lg sources -> core_compiled.lgb -> generated.manifest -> generated.sums
```

so they can only be refreshed together, in that order. `make generate`
(`scripts/generate.lg`) does exactly that; `generated.sums` is the roll-up CI
gates on.

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

## Git merge drivers for the generated artifacts

`core_compiled.lgb` is a binary bundle and `generated.sums` is a one-line
digest: neither has a meaningful 3-way merge, and keeping either side verbatim
leaves a *stale* artifact. `.gitattributes` points both at custom drivers that
recompute them from the merged `.lg` sources instead:

| Path | Driver | Action |
|---|---|---|
| `pkg/rt/core_compiled.lgb` | `merge=lgb` (`scripts/git-merge-lgb.sh`) | regenerates the bundle |
| `pkg/rt/generated.sums` | `merge=sums` (`scripts/git-merge-sums.sh`) | recomputes the digest |

```sh
make install-hooks
```

registers both (`git config merge.lgb.*` and `merge.sums.*`). That config lives
in `.git/config`, which is not shared, so **each clone needs the registration
once**. After it, rebases and merges that touch an embedded `.lg` source
regenerate these two automatically — no binary merge conflicts when stacking PRs
that edit `core.lg` and friends.

`pkg/rt/generated.manifest` deliberately has no driver: it is one sorted record
per line and normally merges as text.

**jj does not run git merge drivers**, the same gap it has with git hooks. Under
jj, reconcile all of these with `make generate` after the rebase.

### Recovering from a conflict

Without the drivers registered — or under jj — a merge touching the `.lg`
sources leaves all three files conflicted:

```
UU pkg/rt/core_compiled.lgb
UU pkg/rt/generated.manifest
UU pkg/rt/generated.sums
```

Resolve the `.lg` source conflicts first, then regenerate. Strip the conflict
markers from the two text files *before* running `make generate`: its first step
queries the manifest, and the parser rejects a marker as a bad record
(`malformed manifest line: "<<<<<<< HEAD"`). Which side you keep does not
matter — both are rewritten wholesale:

```sh
git checkout --ours pkg/rt/generated.manifest pkg/rt/generated.sums
make generate
git add pkg/rt/core_compiled.lgb pkg/rt/generated.manifest pkg/rt/generated.sums
```

`core_compiled.lgb` needs no pre-step: it is marked `binary` in
`.gitattributes`, so git leaves an intact copy in the worktree — enough to build
the `./lg` that `make generate` itself needs.

## `go build` cannot regenerate

Go has no build-time codegen hook, so `go build` never regenerates. Only
`make generate` or `go generate ./cmd/lgbgen` do. The `//go:generate` directives
live in `cmd/lgbgen/generate.go`.

## Implementation

- `pkg/genmanifest` — source-hashing + staleness comparison (the single source
  of truth, shared by the test and the CLI).
- `cmd/lgbgen` — writes `pkg/rt/generated.sums` on every regen (both the bundle
  and `--target=go` paths).
- `cmd/check-generated` — the CLI used by the Makefile target and the hook;
  `-write-manifest` rewrites `generated.manifest`, `-write` the `generated.sums`
  digest, and `-o PATH` is the entry point the `sums` merge driver calls.

## Bundle-only regen

If you are certain you are not touching the Go-lowered path, the bundle alone is
`go run -tags bootstrap ./cmd/lgbgen`. Prefer `make generate` regardless, so the
two artifacts never drift apart (the `-tags gogen_ir` path silently diverges
from the untagged path when only one is regenerated).
