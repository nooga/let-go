#!/usr/bin/env bash
# Promotion smoke: boot-time budget.
#
# Asserts a candidate binary still boots quickly. Startup is dominated by
# bundle decode (bootprobe reports decode ~= total), so this is effectively a
# decode budget.
#
# MEDIAN-OF-5, not a single sample: boot time varies ~3.3x run to run on a
# loaded host, which makes a single shot unusable as a gate. Measured spread
# across 9 trials on an idle M3:
#     median-of-1  3.33x     median-of-3  1.46x
#     median-of-5  1.22x     median-of-9  1.16x
# 5 is the knee — 9 buys 0.06x for double the cost. Same reasoning as the
# median-of-3 change to bench-ratchet in #561.
#
# The default budget sits ~2x above the measured ceiling (median-of-5 topped
# out at 4.10ms), so it does not flake, while still catching the regression
# class in #663 ("startup regressed 5ms -> 9ms unnoticed") — a 9ms median trips
# an 8ms budget.
set -euo pipefail

PROBE="${1:?usage: smoke-boot.sh <bootprobe-binary> [budget-ms] [samples]}"
BUDGET_MS="${2:-8}"
SAMPLES="${3:-5}"

if [ ! -x "$PROBE" ]; then
  echo "smoke-boot: $PROBE is not executable" >&2
  exit 1
fi

vals=()
for _ in $(seq "$SAMPLES"); do
  out="$("$PROBE")"
  ms="$(printf '%s\n' "$out" | sed -n 's/.*total=\([0-9.]*\).*/\1/p')"
  if [ -z "$ms" ]; then
    echo "smoke-boot: could not parse bootprobe output:" >&2
    printf '%s\n' "$out" >&2
    exit 1
  fi
  vals+=("$ms")
done

median="$(printf '%s\n' "${vals[@]}" | sort -g | awk '{a[NR]=$1} END{print a[int((NR+1)/2)]}')"

printf 'smoke-boot: samples=%s median=%sms budget=%sms\n' \
  "$(printf '%s ' "${vals[@]}")" "$median" "$BUDGET_MS"

if awk -v m="$median" -v b="$BUDGET_MS" 'BEGIN{exit !(m>b)}'; then
  echo "smoke-boot: FAIL — median boot ${median}ms exceeds ${BUDGET_MS}ms budget" >&2
  echo "  startup is dominated by bundle decode; see #663 and bootprobe -decodes N" >&2
  exit 1
fi
echo "smoke-boot: ok"
