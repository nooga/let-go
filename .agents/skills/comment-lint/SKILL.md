---
name: comment-lint
description: Act on the devlog/breadcrumb comment findings from scripts/lint.lg — decide per hit whether to rewrite or delete a comment, and rewrite it so it describes the code rather than the change that produced it. Use when running the comment linter, cleaning up flagged comments, reviewing a comment-cleanup diff, or writing a comment that describes why code has its current shape.
---

# comment-lint — acting on the devlog-comment report

Drives `scripts/lint.lg`, the report-only scan for comments that narrate a
change instead of describing the code. The script decides *what to look at*;
everything below is the judgement it deliberately leaves out. Rationale and the
open phrase-list questions: nooga/let-go#835.

## How to run

From the repository root:

```
lg scripts/lint.lg                  # default: pkg scripts
lg scripts/lint.lg pkg scripts cmd  # widen the scan
```

Report-only. It never edits source, and findings do not affect the exit status
(a run that finds 13 things still exits 0), so a hit is never a build gate and
never an instruction on its own.

## What a hit means

**A hit flags a comment block to re-read, not a phrase to delete.** The phrase
list matches tense markers — "used to be", "moved from", "this PR" — because
those are cheap to match and rarely appear in a description of working code.
The thing being flagged is the narration underneath, which can survive any
edit that only removes the matched words.

So the linter going quiet is not the goal. The goal is a comment that a reader
who has never seen a previous version of the file can use.

## The test

> Would this sentence still be true and useful to someone reading this file for
> the first time, who has no idea it was ever different?

If removing the comparison leaves the sentence incomplete, the comparison was
history and the sentence needs rewriting, not trimming.

## Counterfactual, yes; historical contrast, no

These two look alike and are not alike.

**Counterfactual — keep.** What goes wrong *without the code as it stands*.
This is the honest content of most regression-test and benchmark comments, and
it is frequently the right rewrite:

```go
// asBytes backs the binary file/stream sinks (spit, write!). The byte-array
// case is the one that needs it: without the coercion, a byte-array handed to
// spit/write! stringifies to its #byte-array[…] repr, so bytes >127 never
// reach the sink.
```

A reader can check that against the code in front of them.

**Historical contrast — rewrite.** The previous design, re-tensed. The giveaway
words are `rather than`, `instead`, `not X but Y`, `alone could do without`,
`only got away without`:

```go
// It is passed to the module scaffolder rather than baked in, because the
// AOT native path shares that code.
```

`rather than baked in` only parses if you know it was once baked in. The
second clause is subtler, and it is the reason the rule above exists: the AOT
rationale is *true* — `pkg/gomod`'s package doc states it, naming #596 and the
three behaviors — but it is a fact about `pkg/gomod`, already written on
`pkg/gomod`. Restating it on a constant in `pkg/cli` is a second copy that only
the change's author would think to put there. When a kept claim checks out,
ask where it already lives before keeping it here.

**Removing the flagged phrase without removing the history is the failure mode
this skill exists to prevent.** It passes the linter and changes nothing for
the reader.

## Rewrite or delete

Do not assume every flagged comment has a salvageable core. Three outcomes,
in rough order of frequency:

1. **Trim.** The narration is a clause on an otherwise good comment. Drop the
   clause. (`the miscompile this PR fixes` → `the miscompile`.)
2. **Delete the comment.** What remains is a fact about the code's history that
   a reader of this file does not need, or a forwarding address to code that is
   already reachable by name. A pointer saying four helpers "moved to" another
   namespace tells a reader of the file they left nothing about the file they
   are in.
3. **Rewrite from the code.** The comment records something real — an
   invariant, a hazard, why an obvious simpler form is wrong — stated as
   history. Re-derive it by reading the code, not by re-tensing the sentence.

**Every causal claim you keep must name code you can point at.** A rewrite
inherits the old comment's assertions, and those were true of a tree that has
moved. Before keeping "because X shares this" or "so that Y can Z", find X and
Y in the current tree — repository-wide, including tests and build-tagged
files, not just the package in front of you. If you cannot, the claim is
history and goes; if you can, check whether the place it is already documented
is a better home than this one.

Outcome 3 is where accuracy slips. A comment written as history was true of
code that no longer exists, so its details can be stale in ways the linter
cannot see: a benchmark comment naming the fields `root` and `curr` survived a
rename to `root` and `rootBind`, and re-tensing it would have preserved the
wrong names. **Read the code before rewriting; do not paraphrase the old
comment.**

## The one hard rule

**Never drop an invariant.** Flagged comments often carry the only written
record of why something must be done in a particular order, why a simpler form
breaks, or what a test is actually pinning. Losing that to satisfy a
report-only linter is a strictly worse outcome than leaving the comment alone.
When the narration and the invariant cannot be separated, keep the comment and
say so in the PR rather than shipping a lossy rewrite.

## Two operational gotchas

- **A block reports one phrase and stops** — the first match in
  `devlog-phrases` order, not the one appearing earliest in the comment. Fixing
  it can expose a second in the same block, so re-run the linter after a sweep
  rather than trusting one pass.
- **Editing a comment in a `.lg` file that `pkg/rt/generated.manifest` lists
  changes generated artifacts.** The manifest digests its source inputs and
  `generated.sums` digests the manifest, so a comment-only edit to one needs
  `make generate`, and both files are committed. Being git-tracked is not the
  test — most `.lg` files under `pkg/` are inputs, while `scripts/lint.lg` and
  `cmd/lginterop/lginterop.lg` are not. Grep the manifest for the path. The compiled outputs are
  byte-identical — comments never reach the reader's form tree, which is also
  why the linter has to scan raw source.

The tool reports; you exercise the judgement.
