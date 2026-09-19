#!/bin/bash
# Shared helpers for the size/boot audit scripts. Source, don't run.
#
# Portable across macOS (bash 3.2, BSD stat) and Linux (GNU stat).

# fsize <path> -> bytes
fsize() { stat -f%z "$1" 2>/dev/null || stat -c%s "$1"; }

# repo_root -> absolute path of the let-go checkout we were invoked from
repo_root() { git rev-parse --show-toplevel; }

# mk_worktree <repo> <ref> -> path to a fresh detached worktree (echoes path)
# Caller must rm_worktree when done. We never checkout inside the caller's
# own tree, so a churny dev checkout is left untouched.
mk_worktree() {
  local repo="$1" ref="$2" wt
  wt="$(mktemp -d)/lg-audit"
  git -C "$repo" worktree add -f --detach "$wt" "$ref" >/dev/null 2>&1 || return 1
  echo "$wt"
}
rm_worktree() {
  local repo="$1" wt="$2"
  git -C "$repo" worktree remove --force "$wt" 2>/dev/null || true
  rmdir "$(dirname "$wt")" 2>/dev/null || true
}
