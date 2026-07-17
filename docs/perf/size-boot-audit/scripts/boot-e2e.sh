#!/bin/bash
# Per-release END-TO-END boot: real process launch (`lg -e 1`) via hyperfine.
# This includes the OS-spawn + Go-runtime floor (~5 ms) that the in-process
# InitFromLGB benchmark does NOT see, so it is the "what a user waits for"
# number. Complements release-compare.sh (the sensitive in-proc number).
#
# Requires: hyperfine.
# Usage:    ./boot-e2e.sh [tag ...]
# Output:   CSV on stdout (release,e2e_mean_ms,e2e_min_ms)
set -u
command -v hyperfine >/dev/null || { echo "hyperfine not installed" >&2; exit 1; }
cd "$(dirname "$0")" && . ./lib.sh
REPO="$(repo_root)"
TAGS=("$@"); [ ${#TAGS[@]} -eq 0 ] && TAGS=(v1.7.4 v1.8.0 v1.9.0 v1.10.0 v1.11.0 v1.11.1 main)

WT="$(mk_worktree "$REPO" main)" || { echo "worktree failed" >&2; exit 1; }
trap 'rm_worktree "$REPO" "$WT"' EXIT
cd "$WT"

echo "release,e2e_mean_ms,e2e_min_ms"
for tag in "${TAGS[@]}"; do
  git checkout -q --detach "$tag" 2>/dev/null && git clean -fdq 2>/dev/null || { echo "$tag,CHECKOUT_FAIL,"; continue; }
  go build -ldflags="-s -w" -o /tmp/lgh . 2>/dev/null || { echo "$tag,BUILD_FAIL,"; continue; }
  hyperfine -N --warmup 8 --min-runs 40 --export-json /tmp/h.json "/tmp/lgh -e 1" >/dev/null 2>&1
  read mean min < <(python3 -c "import json;d=json.load(open('/tmp/h.json'))['results'][0];print(round(d['mean']*1000,2),round(d['min']*1000,2))")
  echo "$tag,$mean,$min"
done
