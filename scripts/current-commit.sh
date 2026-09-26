#!/usr/bin/env sh
# Print the commit whose tree is checked out here, as a short git sha.
#
# `git rev-parse HEAD` cannot answer this in a jj repository. A colocated jj
# repo keeps ONE git HEAD for the whole repository, tracking whatever the
# default workspace last exported, so every `jj workspace` reports the same
# value no matter which commit it is on. A workspace sitting clean on main has
# been observed reporting a bookmarkless local commit that is not an ancestor
# of upstream main and never will be.
#
# That matters here beyond provenance: the Makefile stamps this into
# `-X main.commit`, which becomes let-go.semver/current-commit, which is what
# SHA-pinned `require-letgo` constraints evaluate against. A wrong stamp makes
# those checks fire against a commit the binary was not built from.
#
# jj stores its commits in the git object store, so a jj commit id IS a git
# sha: the output is resolvable with plain git by someone who has never used
# jj. Falls back to git, then to "none" outside a checkout, matching the
# previous behaviour exactly where jj is absent.
#
# Keep this in step with jjWorkspaceSHA in cmd/bench-ratchet/main.go — the two
# answer the same question for the same reason.
set -u

short=12
if [ "${1:-}" = "--short" ] && [ -n "${2:-}" ]; then
	short=$2
fi

# Resolves git commands with the correct git directory, handling jj workspaces
# that may not have .git in the working directory. Fall back to plain git if
# jj is not in use.
git_cmd() {
	if command -v jj >/dev/null 2>&1 && jj workspace root >/dev/null 2>&1; then
		local git_dir
		if git_dir=$(jj git root 2>/dev/null); then
			git --git-dir="$git_dir" "$@"
			return
		fi
	fi
	# Plain git fallback
	git "$@"
}

if command -v jj >/dev/null 2>&1 && jj workspace root >/dev/null 2>&1; then
	# We are in a jj workspace. Resolve the commit via jj, not git,
	# because `git rev-parse HEAD` returns the shared HEAD which may differ
	# from this workspace's commit.

	# @ when the working copy differs from its parent, the parent when it does
	# not: a clean checkout's @ is an empty commit that describes nothing and is
	# rewritten on every snapshot, whereas the parent is typically a real,
	# pushed commit. Deliberately snapshots (no --ignore-working-copy) so
	# emptiness is judged against the files actually present.
	tmpl="if(empty, parents.map(|p| p.commit_id().short($short)).join(\" \"), commit_id.short($short))"
	sha=$(jj log --no-graph -r @ -T "$tmpl" 2>/dev/null)
	# An empty @ with several parents is a merge: no single commit describes
	# that tree, so fall back to @ itself rather than picking one arbitrarily.
	case "$sha" in
		*\ *) sha=$(jj log --no-graph -r @ -T "commit_id.short($short)" 2>/dev/null) ;;
	esac
	# Only trust it if git can resolve it. jj exports lazily, so sync
	# once before giving up.
	if [ -n "$sha" ]; then
		if ! git_cmd cat-file -e "$sha^{commit}" 2>/dev/null; then
			jj git export >/dev/null 2>&1
		fi
		if git_cmd cat-file -e "$sha^{commit}" 2>/dev/null; then
			printf '%s\n' "$sha"
			exit 0
		fi
	fi
	# Resolution in jj failed. Print "none" rather than falling back to the
	# shared git HEAD, which would be wrong. Only plain git checkouts should
	# use `git rev-parse HEAD`.
	echo none
	exit 0
fi

# Plain git fallback: only for checkouts without jj
git rev-parse --short="$short" HEAD 2>/dev/null || echo none
