---
status: active
last-verified: 2026-06-05
authoritative-for:
  - docs-frontmatter-convention
  - frontmatter-maintenance-hook
human-verified:
---

# Docs frontmatter — convention and maintenance hook

Every doc in `docs/` carries a small YAML frontmatter block describing
its status, freshness, and the topics it's authoritative on. See
[`docs/README.md`](README.md) for the full schema.

This file covers two things: **the conventions every contributor
follows** when editing docs, and **the pre-commit hook** that keeps
the mechanical fields honest.

## Conventions

The hook handles the mechanical fields automatically (see below).
Everything else is contributor judgement:

- **`status:`** — set when creating a new doc; one of `planning`,
  `active`, `shipped`, `superseded`, `archived`. Update when the doc's
  role changes. `status` describes the *doc*, not the underlying
  work — a doc that proposed a now-shipped feature can stay `active`
  if it's still the design-rationale authority.
- **`authoritative-for:`** — the topics this doc is the current
  authority on. The README index uses these to route readers.
- **`supersedes:` / `superseded-by:`** — link displacement
  relationships in both directions when a newer doc takes over on a
  topic. Don't delete the older doc; let the supersession chain stand.
- **`shipped:` / `remaining-open:`** — when a doc tracks a multi-item
  plan, move items from `remaining-open:` (or the inline body) into
  `shipped:` with a commit or PR reference once they land. Don't flip
  `status:` unless the whole doc is now obsolete.
- **If a code change outdates another doc** — update the doc in the
  same PR rather than leaving the convention to rot.

## The human/agent asymmetry

One field behaves differently from everything else:

**`human-verified:`** — the date a human reader last vouched for the
doc as current. **Set only by explicit human action.** The
maintenance hook will not touch it. Agent contributors (LLMs) must
not autopilot-stamp it — even if a `YYYY-MM-DD`-shaped field looks
like something to fill in, leave it blank unless a human has
explicitly asked you to attest.

This is the one rule that genuinely splits humans and agents. The
goal is to keep `human-verified:` as a trustworthy signal: a doc
whose `last-verified:` is today but whose `human-verified:` is six
months old reads very differently from one a human vouched for last
week. That signal collapses if any tooling can stamp it.

## The pre-commit hook

[`scripts/docs_frontmatter_hook.py`](../scripts/docs_frontmatter_hook.py),
wired into [`.pre-commit-config.yaml`](../.pre-commit-config.yaml).
Runs on every staged `docs/**/*.md`. Stdlib-only Python — pre-commit
itself already requires Python, so no new dependency.

### Maintenance mode (default)

1. **Missing frontmatter** → prepends a minimal stub:
   ```yaml
   ---
   status: active
   last-verified: <today>
   human-verified:
   ---
   ```
   Adjust `status:` and add `authoritative-for:` etc. before
   re-staging.
2. **Existing frontmatter** → bumps `last-verified:` to today if the
   existing value is older or missing. Idempotent on docs already
   stamped today.
3. **Never touches** `status:` (when present), `authoritative-for:`,
   `supersedes:`, `superseded-by:`, `shipped:`, `remaining-open:`, or
   `human-verified:`. Only top-level `last-verified:` (no leading
   whitespace) is bumped — indented occurrences inside block scalars
   or nested mappings under human-authored keys are left alone.

The hook modifies files in place. Pre-commit then reports "files
were modified by this hook," at which point you review the changes
and re-stage.

When the hook does something, it prints a one-line note pointing back
at this file — the channel an agent contributor can't miss.

To bypass for in-progress work: `git commit --no-verify` (same
escape hatch as the other hooks).

### Check mode (CI)

```
python3 scripts/docs_frontmatter_hook.py --check <files…>
```

Validates that each path has well-formed frontmatter with a
top-level `status:` and `last-verified:`. Never mutates. Exits
non-zero on any failure.

The CI workflow runs this against the docs changed in each PR. The
check is deliberately the *floor* — frontmatter present and
parseable — not a freshness gate, so PRs that sit across dates still
pass. Staleness is a signal for readers (via `last-verified:` and
`human-verified:`), not a CI failure.

### Edge cases

- **Malformed frontmatter** (opening `---` with no closing delimiter):
  treated as a hard error in both modes. The file is not modified;
  the hook exits non-zero. Fix the delimiter and re-run.
- **UTF-8 BOM** at the start of the file: stripped on read; written
  back without one.
- **Symlinks**: skipped. A staged `docs/*.md` symlink would otherwise
  cause the hook to write through to whatever the symlink points at,
  potentially outside `docs/` or outside the repo.

## Out of scope (follow-ups)

- **Tombstone-on-delete** — when a doc is removed, appending a
  pointer to `docs/archived_files.md` with the SHA of the last
  containing commit. Tracked separately.
- **`docs-status` CLI / agent skill** — a report tool surfacing stale
  docs, contradictions, missing `authoritative-for:` claims, and
  similar judgement-level concerns. Tracked separately.
