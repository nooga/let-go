#!/usr/bin/env bash
# Build the #425 native-entry fib demo end-to-end:
#   lg-compile → main.go frame + fib.go
#   lg -c      → program.lgb (embedded by the frame)
#   go build   → fib.native
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/../../.." && pwd)"
OUT="$(cd "$(dirname "$0")" && pwd)/out"
PREFIX=aotmod
LG="${LG:-}"
if [ -z "$LG" ]; then
  for c in "$ROOT/lg" "$ROOT/lg-bc"; do [ -x "$c" ] && LG="$c" && break; done
fi
[ -n "$LG" ] && [ -x "$LG" ] || { echo "no lg binary — build one or set LG=" >&2; exit 1; }

mkdir -p "$OUT"
# Fresh-ish out: keep go.sum if present; wipe generated sources.
rm -f "$OUT"/main.go "$OUT"/program.lgb "$OUT"/fib.native
rm -rf "$OUT/fib"

echo "==> lg-compile --entry-frame (emits fib.go + main.go native-entry frame)"
# LG_SOURCE_PATHS lets a pre-regen lg see the new entry-frame ns from disk;
# a lg built from this tree embeds it already. --entry-frame is opt-in so
# package-only callers (e.g. Gloat) keep historical behavior.
( cd "$ROOT" && LG_SOURCE_PATHS="${LG_SOURCE_PATHS:-pkg/rt/core}" "$LG" scripts/lg-compile \
    --entry-frame "$OUT" "$PREFIX" examples/aot/native-entry/fib.lg )

echo "==> lg -c program.lgb"
( cd "$ROOT" && "$LG" -c "$OUT/program.lgb" examples/aot/native-entry/fib.lg )

echo "==> go.mod + go build"
cat > "$OUT/go.mod" <<EOF
module $PREFIX

go 1.23

require github.com/nooga/let-go v0.0.0
replace github.com/nooga/let-go => $ROOT
EOF

( cd "$OUT" && go mod tidy && go build -ldflags '-s -w' -o fib.native . )
echo "built: $OUT/fib.native"
echo "try:   $OUT/fib.native       # fib 10"
echo "       $OUT/fib.native x     # fib 34 (spike)"
