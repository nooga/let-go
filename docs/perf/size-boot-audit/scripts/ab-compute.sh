#!/bin/bash
# Interleaved A/B of a compute workload between two refs. hyperfine runs the
# two binaries interleaved with warmups and reports the ratio with error
# bars (the noise-robust method — do NOT eyeball a single `time` run).
#
# Requires: hyperfine.
# Usage:    ./ab-compute.sh <ref-a> <ref-b> <workload.clj> [min-runs]
# Example:  ./ab-compute.sh v1.7.4 main benchmark/fib.clj
# Output:   hyperfine summary + /tmp/ab-<name>.json
set -eu
command -v hyperfine >/dev/null || { echo "hyperfine not installed" >&2; exit 1; }
cd "$(dirname "$0")" && . ./lib.sh
REPO="$(repo_root)"
A="${1:?ref-a}"; B="${2:?ref-b}"; WORK="${3:?workload.clj}"; RUNS="${4:-15}"

# Freeze the workload so BOTH binaries run byte-identical input.
FROZEN="$(mktemp).clj"; cp "$REPO/$WORK" "$FROZEN" 2>/dev/null || cp "$WORK" "$FROZEN"
name="$(basename "$WORK" .clj)"

build() { # <ref> <out>
  local wt; wt="$(mk_worktree "$REPO" "$1")"
  ( cd "$wt" && go build -ldflags="-s -w" -o "$2" . )
  rm_worktree "$REPO" "$wt"
}
build "$A" "/tmp/lg-A"
build "$B" "/tmp/lg-B"

hyperfine --warmup 3 --min-runs "$RUNS" \
  -n "$name $A" "/tmp/lg-A $FROZEN" \
  -n "$name $B" "/tmp/lg-B $FROZEN" \
  --export-json "/tmp/ab-$name.json"
rm -f "$FROZEN"
