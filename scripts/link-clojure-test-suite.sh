#!/usr/bin/env bash
# Link test/clojure-test-suite into a secondary worktree / jj workspace.
#
# test/clojure-test-suite is a git submodule. jj does not manage submodules
# (they go through git), and git worktrees / jj workspaces do not materialize
# them, so the path exists but is empty in every secondary workspace. Tests that
# need the pinned jank suite (make native-entry-gate ->
# TestJankSuiteDirectABIGeneratedGo) then fail loudly.
#
# This script points <workspace>/test/clojure-test-suite at the initialized
# submodule in the primary worktree. jj does not report that symlink as a
# working-copy change, so it does not pollute commits.
#
# Usage: scripts/link-clojure-test-suite.sh <workspace-path> [primary-worktree]
# Idempotent: re-running refreshes an existing symlink and refuses to clobber a
# real directory or file.
set -euo pipefail

REL="test/clojure-test-suite"
SENTINEL="$REL/test/clojure/core_test/identical_qmark.cljc"

usage() {
	echo "usage: $0 <workspace-path> [primary-worktree]" >&2
	exit 2
}

[ $# -ge 1 ] && [ $# -le 2 ] || usage

WORKSPACE=$(cd "$1" 2>/dev/null && pwd) || {
	echo "error: workspace path does not exist: $1" >&2
	exit 1
}

if [ $# -eq 2 ]; then
	PRIMARY=$(cd "$2" 2>/dev/null && pwd) || {
		echo "error: primary worktree does not exist: $2" >&2
		exit 1
	}
else
	# Default: the primary worktree is where this script lives.
	PRIMARY=$(cd "$(dirname "$0")/.." && pwd)
fi

if [ "$PRIMARY" = "$WORKSPACE" ]; then
	echo "error: workspace and primary worktree are the same path ($PRIMARY);" >&2
	echo "       run 'git submodule update --init $REL' there instead" >&2
	exit 1
fi

if [ ! -e "$PRIMARY/$SENTINEL" ]; then
	echo "error: submodule not initialized in primary worktree $PRIMARY" >&2
	echo "       run: (cd $PRIMARY && git submodule update --init $REL)" >&2
	exit 1
fi
# Resolve to the physical directory so a chain of workspace symlinks collapses
# to the one real submodule checkout.
SRC=$(cd "$PRIMARY/$REL" && pwd -P)

DST="$WORKSPACE/$REL"
if [ -L "$DST" ]; then
	CURRENT=$(readlink "$DST")
	if [ "$CURRENT" = "$SRC" ]; then
		echo "ok: $DST already links to $SRC"
		exit 0
	fi
	rm "$DST"
elif [ -e "$DST" ]; then
	if [ -d "$DST" ] && [ -z "$(ls -A "$DST")" ]; then
		# Empty submodule placeholder directory: safe to replace.
		rmdir "$DST"
	else
		echo "error: refusing to clobber existing non-empty path $DST" >&2
		exit 1
	fi
fi

mkdir -p "$(dirname "$DST")"
ln -s "$SRC" "$DST"
echo "ok: linked $DST -> $SRC"
