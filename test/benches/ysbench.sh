#!/usr/bin/env bash
# ITER-0035 acceptance gate (issue #352): the IR-optimized (inline-enabled)
# yamlstar parse must be measurably faster than the interpreted baseline.
# Runs ysbench.lg in both modes and compares ms-per-call.
#
#   LG_YAMLSTAR_SRC        yamlstar core/src checkout (default ~/development/yamlstar/core/src)
#   LG_YSBENCH_MIN_SPEEDUP required interp/ir ratio (default 1.10 = 10% faster)
#
# Exits 0 on pass, 1 on fail, 0 with a skip notice when yamlstar is absent.
set -euo pipefail
cd "$(dirname "$0")/../.."

YS_SRC="${LG_YAMLSTAR_SRC:-$HOME/development/yamlstar/core/src}"
if [ ! -d "$YS_SRC/yamlstar" ]; then
  echo "skip: yamlstar sources not found at $YS_SRC (set LG_YAMLSTAR_SRC)" >&2
  exit 0
fi

[ -x ./lg ] || go build -o lg .

run() {
  LG_SOURCE_PATHS="$YS_SRC" ./lg test/benches/ysbench.lg "$1" \
    | tee /dev/stderr | awk '/^ms-per-call:/{print $2}'
}

interp=$(run interp)
ir=$(run ir)
speedup=$(awk -v a="$interp" -v b="$ir" 'BEGIN{printf "%.2f", a/b}')
min="${LG_YSBENCH_MIN_SPEEDUP:-1.10}"

echo "interp: ${interp} ms/call   ir+inline: ${ir} ms/call   speedup: ${speedup}x (required: ${min}x)"
if awk -v s="$speedup" -v m="$min" 'BEGIN{exit !(s+0 >= m+0)}'; then
  echo "PASS: ITER-0035 criterion met"
else
  echo "FAIL: ITER-0035 criterion not met"
  exit 1
fi
