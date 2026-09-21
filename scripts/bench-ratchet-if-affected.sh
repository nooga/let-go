#!/usr/bin/env bash
# Run the performance ratchet only when the push can change what it measures.
#
# The ratchet used to run on every push. That made a stale bar block EVERY push
# from a machine regardless of content: a docs edit or a fix to the ratchet tool
# itself could be gated behind a baseline rebaseline it had nothing to do with.
# cmd/ratchet-scope decides scope from the Go dependency closure of the
# benchmarked packages, so "the IR or the VM, or anything they transitively
# depend on" is computed rather than approximated by a path regex.
#
# Fail-safe in every direction: if the range cannot be determined, if the
# scope tool cannot be built or run, or if anything is ambiguous, the ratchet
# RUNS. A false run costs one benchmark; a false skip defeats the gate.
set -uo pipefail

say() { printf 'bench-ratchet-scope: %s\n' "$*" >&2; }

run_ratchet() { exec make bench-ratchet; }

# The tip being pushed. prek's pre-push hook names it in PRE_COMMIT_TO_REF,
# which need not be the checkout: a push can name another ref, and a dirty
# tree differs from the commit being pushed. The decision below is computed
# from the checked-out tree (`go list` for the closure, the filesystem for
# deletions), so it is only valid when that tree IS the pushed tip. When it is
# not, or the tip cannot be resolved, the range is reported as undeterminable
# and the ratchet runs.
#
# The range starts at the fork point with upstream main, not at
# PRE_COMMIT_FROM_REF: the bar the ratchet compares against describes main,
# so everything the pushed tree adds over main is what can move the numbers.
# jj is the repository's VCS and plain `git merge-base` is refused here, so ask
# jj first and fall back to git only where jj is absent (CI checkouts).
changed_paths() {
  local tip base
  if command -v jj >/dev/null 2>&1 && jj workspace root >/dev/null 2>&1; then
    tip=$(jj log --no-graph -r "${PRE_COMMIT_TO_REF:-@}" -T 'commit_id ++ "\n"' 2>/dev/null) || return 1
    case "$tip" in ''|*$'\n'*) return 1 ;; esac
    # Deliberately NOT --ignore-working-copy: the comparison must see the tree
    # as it is now, and a stale snapshot would hide a difference.
    local drift
    drift=$(jj diff --name-only --from "$tip" --to @ 2>/dev/null) || return 1
    if [ -n "$drift" ]; then
      say "the pushed tip ${tip:0:12} differs from the checked-out tree"
      return 1
    fi
    for revset in "fork_point($tip | main@upstream)" "fork_point($tip | main)" 'main@upstream'; do
      if base=$(jj log --no-graph --ignore-working-copy -r "$revset" \
                  -T 'commit_id' 2>/dev/null) && [ -n "$base" ]; then
        jj diff --name-only --from "$base" --to "$tip" 2>/dev/null && return 0
      fi
    done
    return 1
  fi
  tip=$(git rev-parse --verify --quiet "${PRE_COMMIT_TO_REF:-HEAD}^{commit}" 2>/dev/null) || return 1
  if ! git diff --quiet "$tip" -- 2>/dev/null; then
    say "the pushed tip ${tip:0:12} differs from the checked-out tree"
    return 1
  fi
  base=$(git rev-parse --verify --quiet origin/main 2>/dev/null) || return 1
  git diff --name-only "$base"..."$tip" 2>/dev/null || return 1
}

paths=$(changed_paths) || {
  say "could not determine the pushed range; running the ratchet"
  run_ratchet
}

if [ -z "${paths//[[:space:]]/}" ]; then
  say "empty change set against the base; running the ratchet"
  run_ratchet
fi

# Build to build/, which is gitignored. Never a bare `go build`, which would
# drop a binary in the working directory where nothing ignores it. `go build
# -o` does not create the parent directory, and a fresh checkout has no
# build/, so make it first.
scope_bin="build/ratchet-scope"
mkdir -p build || { say "cannot create build/; running the ratchet"; run_ratchet; }
if ! go build -o "$scope_bin" ./cmd/ratchet-scope 2>&1; then
  say "could not build ratchet-scope; running the ratchet"
  run_ratchet
fi

decision=$(printf '%s\n' "$paths" | "./$scope_bin" -v) || {
  say "ratchet-scope failed; running the ratchet"
  run_ratchet
}

case "$decision" in
  SKIP)
    say "no pushed file is in the measured dependency closure — skipping the benchmark ratchet."
    say "  (the IR, the VM, their transitive dependencies, the generated artifacts,"
    say "   and the module/toolchain files are what require it)"
    exit 0
    ;;
  RUN)
    run_ratchet
    ;;
  *)
    say "unrecognised decision '$decision'; running the ratchet"
    run_ratchet
    ;;
esac
