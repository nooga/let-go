#!/bin/bash
# Per-release comparison: shipped binary size, embedded core-bundle size,
# in-process boot cost (BenchmarkInitFromLGB), and core .lg file count.
#
# Boot is the same benchmark the ratchet gates (decode core_compiled.lgb +
# run its main chunk). We report the MIN over -count runs: scheduler / GC /
# thermal noise only ever adds time, so the floor is the truest signal.
#
# Usage:   ./release-compare.sh [tag ...]
# Default tags: v1.7.4 v1.8.0 v1.9.0 v1.10.0 v1.11.0 v1.11.1 main
# Output:  CSV on stdout (release,date,bin_mb,lgb_kb,boot_us,go_files)
#
# Absolute µs are machine-specific; trust the RATIOS across releases.
set -u
cd "$(dirname "$0")" && . ./lib.sh
REPO="$(repo_root)"
TAGS=("$@"); [ ${#TAGS[@]} -eq 0 ] && TAGS=(v1.7.4 v1.8.0 v1.9.0 v1.10.0 v1.11.0 v1.11.1 main)

WT="$(mk_worktree "$REPO" main)" || { echo "worktree failed" >&2; exit 1; }
trap 'rm_worktree "$REPO" "$WT"' EXIT
cd "$WT"

echo "release,date,bin_mb,lgb_kb,boot_us,go_files"
for tag in "${TAGS[@]}"; do
  git checkout -q --detach "$tag" 2>/dev/null && git clean -fdq 2>/dev/null || { echo "$tag,CHECKOUT_FAIL,,,,"; continue; }
  d=$(git show -s --format=%cs HEAD)
  if go build -ldflags="-s -w" -o /tmp/lgr . 2>/dev/null; then
    bin=$(echo "scale=2;$(fsize /tmp/lgr)/1048576" | bc)
  else bin=BUILD_FAIL; fi
  lgb=$(( $(fsize pkg/rt/core_compiled.lgb 2>/dev/null || echo 0) / 1024 ))
  nlg=$(find pkg/rt/core -name '*.lg' 2>/dev/null | wc -l | tr -d ' ')
  raw=$(go test -run=NONE -bench='^BenchmarkInitFromLGB$' -benchtime=200x -count=6 ./pkg/compiler 2>/dev/null \
        | awk '/InitFromLGB/{for(i=1;i<=NF;i++) if($i=="ns/op") print $(i-1)}')
  if [ -n "$raw" ]; then
    boot_us=$(echo "scale=1;$(echo "$raw" | sort -n | head -1)/1000" | bc)
  else boot_us=NA; fi
  echo "$tag,$d,$bin,$lgb,$boot_us,$nlg"
done
