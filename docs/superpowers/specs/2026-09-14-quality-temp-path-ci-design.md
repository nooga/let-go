---
status: active
last-verified: 2026-09-14
---

# Make quality temporary paths portable and restore CI frontmatter

## Problem

#871's post-push Linux `build` job fails in `quality_cli_test.lg` and
`quality_delta_test.lg`. `os/temp-dir` returns `/tmp` there, but several
callers concatenate a filename without a separator, producing paths such as
`/tmpquality-cli-test.edn` and `/tmpquality-rev-<sha>.tar`. Those paths target
the filesystem root and cannot be written by the runner. Setting
`TMPDIR=/tmp` reproduces both failures locally. Separately, the
`docs-frontmatter` job identifies one added plan without an opening YAML
block: `docs/superpowers/plans/2026-09-13-quality-scorer-review-fixes.md`.

## Design

At each affected quality temporary-file construction, insert an explicit `/`
between the directory and filename. Cover the production Git archive tarball,
the quality CLI test EDN paths, and the quality-reap probe/pid paths found in
the same sweep. A duplicate slash on hosts whose temp directory already ends
in `/` is accepted by the filesystem and does not change the path's target.
Do not change the public `os/temp-dir` return contract or add a new path API.

Add the required `status` and `last-verified` YAML frontmatter to the named
plan, without rewriting its contents. No benchmark baseline, `jj` behavior,
CI workflow, or #870 source changes are in scope.

## Alternatives

- Change `os/temp-dir` to always include a trailing slash: central, but a
  runtime API change with unknown callers.
- Introduce a general path-join API: broader surface and migration for a
  small set of proven concatenation mistakes.

## Evidence

Before editing, `TMPDIR=/tmp go test ./test -run
'^TestRunner/(quality_cli_test.lg|quality_delta_test.lg)$' -count=1` fails
on `/tmpquality-...`; the #871 CI log shows the same paths. After editing,
that focused command, the ordinary short Go suite, build/vet, and
`scripts/docs_frontmatter_hook.py --check` on the changed plan must pass.
The actual PR CI build/frontmatter jobs are the final check after an
exact-lease push. The approved benchmark-ratchet skip applies only to that
pending push, not to tests or baseline policy.
