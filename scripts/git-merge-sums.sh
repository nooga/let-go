#!/bin/bash
# Git custom merge driver for pkg/rt/generated.sums
#
# generated.sums is a one-line content digest of the .lg + lgbgen sources.
# It has no meaningful 3-way text merge — but keeping one side leaves a STALE
# signature that no longer matches the merged sources. Instead we recompute the
# digest from the merged sources on disk and write it to git's output path, so
# the merged manifest is always correct (the .lg files merge cleanly as text;
# this recomputes their digest). This is a pure recompute — no bundle or
# lowered-tree rebuild (that is the lgb driver's job for core_compiled.lgb).
#
# Git invokes this driver with:
#   $1 = %O (common ancestor)   $2 = %A (current branch — also the OUTPUT path)
#   $3 = %B (other branch)      $4 = %L (marker size)   $5 = %P (worktree path)
# We ignore all but %A (the output) and recompute from the worktree sources.
set -e
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
cd "$REPO_ROOT"
# Recompute the source digest and write the manifest to %A.
go run ./cmd/check-generated -o "$2"
