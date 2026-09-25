#!/bin/bash
# Interleaved A/B of a compute workload between two refs. Each measured round
# runs both binaries, with the first binary alternating by round (AB, BA, ...),
# so machine drift is balanced rather than correlated with one ref.
#
# Requires: python3.
# Usage:    ./ab-compute.sh <ref-a> <ref-b> <workload.clj> [runs]
# Example:  ./ab-compute.sh v1.7.4 ed4ecc215 benchmark/fib.clj
# Output:   paired summary + /tmp/ab-<name>.json
set -eu
command -v python3 >/dev/null || { echo "python3 not installed" >&2; exit 1; }
cd "$(dirname "$0")" && . ./lib.sh
REPO="$(repo_root)"
A="${1:?ref-a}"; B="${2:?ref-b}"; WORK="${3:?workload.clj}"; RUNS="${4:-15}"

# Freeze the workload so BOTH binaries run byte-identical input.
TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT
FROZEN="$TMP/workload.clj"; cp "$REPO/$WORK" "$FROZEN" 2>/dev/null || cp "$WORK" "$FROZEN"
name="$(basename "$WORK" .clj)"

build() { # <ref> <out>
  local wt; wt="$(mk_worktree "$REPO" "$1")"
  ( cd "$wt" && go build -ldflags="-s -w" -o "$2" . )
  rm_worktree "$REPO" "$wt"
}
build "$A" "$TMP/lg-A"
build "$B" "$TMP/lg-B"

python3 ./interleaved-ab.py \
  --name-a "$name $A" --name-b "$name $B" \
  --binary-a "$TMP/lg-A" --binary-b "$TMP/lg-B" \
  --workload "$FROZEN" --warmup 3 --runs "$RUNS" \
  --output "/tmp/ab-$name.json"
