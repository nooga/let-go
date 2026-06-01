#!/bin/bash
# Git custom merge driver for pkg/rt/core_compiled.lgb
#
# The .lgb is deterministic output from the embedded .lg sources. Instead of
# attempting a 3-way binary merge (which always fails), we regenerate the .lgb
# from the merged .lg sources. The .lg files merge cleanly as text; this script
# rebuilds the binary to match.
#
# Git invokes this driver with:
#   $1 = %O (path to common ancestor's version)
#   $2 = %A (path to current branch's version — also the OUTPUT path)
#   $3 = %B (path to merging branch's version)
#   $4 = %L (merge marker size; unused)
#   $5 = %P (pathname in the worktree)
#
# We ignore all five and just regenerate.

set -e

# Find repo root by walking up from this script's location.
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
cd "$REPO_ROOT"

# Regenerate from sources. Output path comes from "$2" (%A).
# lgbgen only compiles the embedded core sources under the `bootstrap` build
# tag; without it, lgbgen refuses to run.
go run -tags bootstrap ./cmd/lgbgen "$2"
