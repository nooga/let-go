#!/bin/bash
# Walk a commit range and record the SHIPPED (stripped) binary size + core
# bundle size at ~N sampled commits. Shows the growth curve / step changes
# that per-release snapshots smear together. Not a bisect — size grows
# gradually across many commits, so a linear walk is more informative.
#
# Usage:   ./binary-size-walk.sh [FROM_REF] [TO_REF] [SAMPLES]
# Default: v1.7.4 main 25
# Output:  CSV on stdout (idx,short,date,stripped_mb,lgb_kb,subject)
set -u
cd "$(dirname "$0")" && . ./lib.sh
REPO="$(repo_root)"
FROM="${1:-v1.7.4}"; TO="${2:-main}"; SAMPLES="${3:-25}"

WT="$(mk_worktree "$REPO" "$TO")" || { echo "worktree failed" >&2; exit 1; }
trap 'rm_worktree "$REPO" "$WT"' EXIT
cd "$WT"

LIST="$(mktemp)"; git rev-list --reverse "$FROM..$TO" > "$LIST"
N=$(wc -l < "$LIST" | tr -d ' ')
STEP=$(( (N + SAMPLES - 1) / SAMPLES )); [ "$STEP" -lt 1 ] && STEP=1

echo "idx,short,date,stripped_mb,lgb_kb,subject"
i=0
for ln in $(seq 1 "$STEP" "$N") "$N"; do
  sha=$(sed -n "${ln}p" "$LIST"); [ -z "$sha" ] && continue
  git checkout -q --detach "$sha" 2>/dev/null || continue
  [ -f pkg/rt/core_compiled.lgb ] || continue
  if ! go build -ldflags="-s -w" -o /tmp/lgw . 2>/dev/null; then
    echo "$i,${sha:0:9},BUILD_FAIL,,,"; i=$((i+1)); continue
  fi
  mb=$(echo "scale=2;$(fsize /tmp/lgw)/1048576" | bc)
  kb=$(( $(fsize pkg/rt/core_compiled.lgb)/1024 ))
  subj=$(git show -s --format=%s "$sha" | tr ',' ';' | cut -c1-45)
  echo "$i,${sha:0:9},$(git show -s --format=%cs "$sha"),$mb,$kb,$subj"
  i=$((i+1))
done
rm -f "$LIST"
