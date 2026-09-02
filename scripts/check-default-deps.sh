#!/usr/bin/env bash
# Default-build dependency gate.
#
# Asserts that a plain, untagged build of the shipping binaries links none of
# the optional heavy subsystems. Those subsystems are meant to be reachable
# only behind a build tag; a single untagged file anywhere in their import
# closure silently puts them in every binary we ship.
#
# This exists because that is not hypothetical. On the glplat branch,
# pkg/rt/interop_glplat.go — generated, in package rt, and therefore in every
# binary — carried no build tag, so the default `lg` build linked
# golang.org/x/image and golang.org/x/text through pkg/glplat/font.go. Cost:
# +718,208 bytes, +3.60% (19,958,434 -> 20,676,642, darwin/arm64, measured
# 2026-09-02). The GL and Ebitengine backends themselves were correctly gated;
# it was the untagged font path that leaked. Nothing catches that today, which
# is the gap this closes — and it runs against the same binary-footprint work
# as #652/#658.
#
# ANCHORING IS LOAD-BEARING. `go list -deps` reports the standard library's
# vendored copies under a `vendor/` prefix, so main already lists
# vendor/golang.org/x/text/{transform,unicode/bidi,secure/bidirule,unicode/norm}
# by way of net/http's IDNA handling. Those are stdlib, not the module, and
# must not trip the gate. Every pattern below is anchored with ^ so the
# vendored copies do not match; drop the anchor and this fails on main.
set -euo pipefail

# Import-path prefixes that must not appear in an untagged build. Extend this
# list when a new optional subsystem lands, not when one trips the gate.
DENY='^golang\.org/x/image/|^golang\.org/x/text/|^github\.com/go-gl/|^github\.com/hajimehoshi/|^github\.com/ebitengine/'

TARGETS=("." "./cmd/lg-runtime")

status=0
for target in "${TARGETS[@]}"; do
    if ! deps="$(go list -deps "$target" 2>/dev/null)"; then
        echo "check-default-deps: go list -deps $target failed" >&2
        exit 2
    fi
    if found="$(printf '%s\n' "$deps" | grep -E "$DENY")"; then
        status=1
        echo "check-default-deps: $target links tag-gated dependencies in an untagged build:" >&2
        printf '  %s\n' $found >&2
    fi
done

if [ "$status" -ne 0 ]; then
    cat >&2 <<'MSG'

Every package above should be reachable only behind a build tag. The usual
cause is one untagged file in the subsystem's import closure — check for a
file with no //go:build line, including generated ones, since a generated
file in package rt puts the whole closure in every binary.
MSG
    exit 1
fi

echo "check-default-deps: OK — untagged builds link no tag-gated subsystems."
