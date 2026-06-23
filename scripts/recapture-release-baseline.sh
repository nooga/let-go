#!/usr/bin/env bash
# Re-derive a historical perf baseline for an OLD release tag that predates the
# bench-ratchet system, captured on THIS machine.
#
# Why this exists: BenchmarkRatchetAnchor (the calibration anchor every baseline
# normalizes against) and the bench-ratchet tool were both introduced in #119,
# AFTER v1.8.0. So an old tag has the pkg/vm micro-benchmarks but neither the
# anchor benchmark nor the harness. The committed docs/perf/historical/v1.8.0.json
# was therefore captured on a one-off (Apple M3) machine, which makes "% vs
# v1.8.0" a cross-machine/cross-arch comparison against the amd64 CI timeline.
#
# To get an honest, same-machine v1.8.0 reference, graft the (self-contained)
# anchor benchmark onto a throwaway worktree of the tag, run bench-ratchet
# against pkg/vm only, then drop the graft. Pair this with pinning the perf job
# to one runner so the recaptured baseline and the timeline share silicon.
#
# Scope: only the pkg/vm micro-benchmarks are recapturable — the modern
# aot_native / ir_bytecode (pkg/ir) suite never existed at old tags.
#
# Usage: scripts/recapture-release-baseline.sh <ref> <out.json> [benchtime=2s]
set -euo pipefail

REF="${1:?usage: recapture-release-baseline.sh <ref> <out.json> [benchtime]}"
OUT="${2:?usage: recapture-release-baseline.sh <ref> <out.json> [benchtime]}"
BENCHTIME="${3:-2s}"

# Absolute-ize OUT before we cd into the worktree.
case "$OUT" in /*) ;; *) OUT="$PWD/$OUT" ;; esac
mkdir -p "$(dirname "$OUT")"

REPO="$(git rev-parse --show-toplevel)"
ANCHOR="$REPO/pkg/vm/bench_ratchet_anchor_test.go"
[ -f "$ANCHOR" ] || { echo "anchor benchmark not found at $ANCHOR" >&2; exit 1; }

WT="$(mktemp -d)/recap"
BR="$(mktemp -d)/bench-ratchet"
cleanup() { git -C "$REPO" worktree remove --force "$WT" 2>/dev/null || true; }
trap cleanup EXIT

# --detach so it coexists with any existing checkout of the same tag.
git -C "$REPO" worktree add --detach "$WT" "$REF" >/dev/null
SHA="$(git -C "$REPO" rev-parse --short "$REF")"

# Graft the anchor (self-contained, no project deps) so ratio_to_anchor is
# computable; absent at pre-#119 tags. Build the current bench-ratchet — the
# tag doesn't have it — and run it against pkg/vm in the tag's tree.
cp "$ANCHOR" "$WT/pkg/vm/bench_ratchet_anchor_test.go"
( cd "$REPO" && go build -o "$BR" ./cmd/bench-ratchet )
( cd "$WT" && "$BR" -packages github.com/nooga/let-go/pkg/vm \
    -benchtime "$BENCHTIME" -count 1 -sha "$SHA" -baseline "$OUT" snapshot )

echo "recaptured $REF (pkg/vm) on $(uname -m) → $OUT"
