---
status: active
last-verified: 2026-09-15
authoritative-for:
  - comment-lint
human-verified:
---

# comment-lint — the comment report and how to act on it

[`scripts/lint.lg`](../scripts/lint.lg) scans let-go and Go source for comments
that narrate a change rather than describe the code, plus a catalog of `.lg`
code patterns that can be simplified. Like
[`docs-status.md`](docs-status.md), it splits in two: the script decides *what
to look at*, and a skill carries the judgement it deliberately leaves out.

It never edits source.

## Running it

```sh
lg scripts/lint.lg                      # default paths: pkg scripts
lg scripts/lint.lg pkg scripts cmd      # widen the scan
lg scripts/lint.lg --edn                # machine-readable findings
lg scripts/lint.lg --gate K1[,K2...]    # exit non-zero if a named rule fires
lg scripts/lint.lg --churn HEAD~5..HEAD # added comments vs added code, for a range
```

Report-only by default: findings do not affect the exit status, so a run that
finds things still exits 0. It is not wired into CI, the Makefile, or the git
hooks — you run it by hand. `--gate` is the only way to make it fail, and it
fails only on the rule kinds you name.

## What it reports

R1–R4 and R6 are objective: computed from the source or from the corpus, with
no hand-curated wording. R5 is the original phrase-list heuristic, reported
under its own heading, kept out of the objective score, and never gateable.

| kind | reported as | scope |
|---|---|---|
| `commented-out-code` | code left in a comment | `.lg` only |
| `restatement` | comment repeats code | `.lg`, Go |
| `comment-divider` | decorative section divider | `.lg`, Go |
| `duplicated-comment` | repeated comment block | `.lg`, Go |
| `comment-density-outlier` | unusually comment-heavy definition | `.lg`, Go |
| `devlog-comment` | development-note phrase (R5, heuristic) | `.lg`, Go |
| `comment-churn` | comment additions relative to code (needs `--churn`) | `.lg`, Go |

The separate code-pattern catalog lives in
[`scripts/lint-code-rules.edn`](../scripts/lint-code-rules.edn) and applies to
`.lg` only.


## Acting on a finding

A hit flags a comment block to re-read, not a phrase to delete. The judgement —
whether to trim, delete, or rewrite, and how to write the replacement — is in
the `comment-lint` skill at
[`.agents/skills/comment-lint/SKILL.md`](../.agents/skills/comment-lint/SKILL.md),
which agents load automatically and humans can read directly.

Two things worth knowing before a sweep, both of which have bitten:

- **A block reports one phrase and stops.** Fixing it can expose a second in the
  same block, so re-run after a sweep rather than trusting one pass.
- **Editing a comment in an `.lg` file that `pkg/rt/generated.manifest` tracks
  changes generated artifacts**, so the edit needs `make generate`. Grep the
  manifest for the path; being git-tracked is not the test.

## Open questions

The rules are not settled. [#835](https://github.com/nooga/let-go/issues/835)
is where the phrase list and the surrounding conventions get decided, including
which of the skill's rules are agreed and which are still one contributor's
proposal. Read it before treating anything in the skill as policy.
