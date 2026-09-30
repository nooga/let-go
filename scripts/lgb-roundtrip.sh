#!/usr/bin/env bash
# lgb-roundtrip.sh — program-level metamorphic check for compiled artifacts.
#
# The relation: for any program P and every way lg can serialize it,
#
#     run(P)  ==  run(deserialize(serialize(compile(P))))
#
# i.e. executing a source file directly and executing the artifact it compiles
# to must produce identical observable output. Existing coverage in
# pkg/bytecode is VALUE-level — varint, float64, string, keyword, symbol, UUID,
# bigint, the local-var table — and proves each encoder round-trips in
# isolation. It cannot catch a whole-program failure such as constant-pool or
# opcode-enum skew between writer and reader, whose symptom is "unknown
# instruction" at load time, or worse, a program that loads and quietly
# computes something else.
#
# Each fixture is checked under every mode in MODES:
#
#   c      lg -c        .lgb, run with `lg app.lgb`
#   cz     lg -z -c     DEFLATE-compressed .lgb
#   strip  lg -strip -c .lgb with its debug companion split out
#   b      lg -b        standalone executable (payload appended to lg itself)
#
# Usage:
#   scripts/lgb-roundtrip.sh [FILE...]    # FILEs, or examples/ if none given
#   scripts/lgb-roundtrip.sh --full       # examples/ + a driver per test/ suite
#   scripts/lgb-roundtrip.sh --selftest   # falsify the checker; no corpus needed
#
# Environment:
#   LG               lg binary to test (default: bin/lg)
#   MODES            space-separated subset of the modes above (default: all)
#   FIXTURE_TIMEOUT  per-invocation bound in seconds (default: 20; applied only
#                    when `timeout` or `gtimeout` is on PATH)
#
# `--full` is not green. Suites that require a namespace loaded lazily rather
# than at boot (zip, data, check, clojure.pprint, clojure.data, and the ir.*
# namespaces) diverge because their compiled artifacts carry those
# namespaces' vars with nil values: nooga/let-go#954, found by this script.
# deftest_family_test also diverges, not yet triaged. Treat `--full` as a
# research mode until those are resolved.
#
# Exit codes:
#   0  every usable fixture round-trips identically in every mode
#   1  at least one fixture DIVERGED, nothing was verified, or selftest failed
#   2  setup error (no lg binary, unknown mode, missing file)

set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
LG="${LG:-$ROOT/bin/lg}"
MODES="${MODES:-c cz strip b}"
FIXTURE_TIMEOUT="${FIXTURE_TIMEOUT:-20}"

FULL=0; SELFTEST=0; FILES=()
for arg in "$@"; do
    case "$arg" in
        --full)     FULL=1 ;;
        --selftest) SELFTEST=1 ;;
        -*)         echo "unknown flag: $arg" >&2; exit 2 ;;
        *)          FILES+=("$arg") ;;
    esac
done

for m in $MODES; do
    case "$m" in c|cz|strip|b) ;; *) echo "unknown mode: $m" >&2; exit 2 ;; esac
done
[ -x "$LG" ] || { echo "no lg binary at $LG (run: make build, or set LG=)" >&2; exit 2; }

# Stock macOS ships no `timeout`. The bound only keeps a hung fixture from
# costing a minute, so run unbounded rather than refuse to run. Every use is
# spelled ${BOUND[@]+"${BOUND[@]}"} because bash 3.2, macOS's /bin/bash,
# treats an empty array as unset under `set -u`.
if command -v timeout >/dev/null 2>&1; then BOUND=(timeout "$FIXTURE_TIMEOUT")
elif command -v gtimeout >/dev/null 2>&1; then BOUND=(gtimeout "$FIXTURE_TIMEOUT")
else BOUND=()
fi

TMP=$(mktemp -d); trap 'rm -rf "$TMP"' EXIT

# --- Classification -------------------------------------------------------
#
# A fixture is only an ORACLE if it agrees with itself. Every fixture is run
# twice directly before it is round-tripped at all; one that disagrees with its
# own previous output cannot distinguish a serialization bug from its own
# jitter, and is reported NONDET and excluded rather than counted as a pass or
# a failure. Of the programs in examples/, concurrent-primes, goroutines and
# primes are scheduler- or timing-dependent and would otherwise produce phantom
# failures.
#
# The counts are kept separate on purpose. NONDET is not a pass — a corpus that
# quietly drifted to all-NONDET would otherwise report success while verifying
# nothing. `agree` and `diverged` count (fixture, mode) pairs; the others count
# fixtures, since they are decided before any mode runs.
agree=0; diverged=0; nondet=0; errored=0; skipped=0
DIVERGED_LIST=()

reset_counts() {
    agree=0; diverged=0; nondet=0; errored=0; skipped=0; DIVERGED_LIST=()
}

# Structurally unsuitable fixtures, with the reason stated rather than left for
# the next reader to rediscover by watching the harness time out. These are not
# failures and not NONDET: they cannot participate in an output-comparison
# relation at all. The test/ entries observe the compile itself (macro
# expansion, load-time bindings, source text), which a compiled artifact
# replays without by design. Measured 2026-09-26.
skip_reason() {
    case "$1" in
        examples/server.lg)     echo "runs a server; never terminates" ;;
        examples/mandelbrot.lg) echo "animation loop; still running at 90s" ;;
        examples/term_demo.lg)  echo "needs a TTY; exits 1 without one" ;;
        test/try_finally_test.lg)      echo "counts macro expansions, which happen at compile time" ;;
        test/file_var_test.lg)         echo "reads *file*, bound only while loading source" ;;
        test/clojure_repl_test.lg)     echo "source-fn reads the source file" ;;
        test/exception_classes_test.lg) echo "compares a warning printed while loading source" ;;
        *) return 1 ;;
    esac
}

# observe CMD... — the command's stdout, then its stderr, each captured whole.
# Merging them with 2>&1 would make the comparison depend on how the two
# streams happened to interleave, which varies between runs: io_binding_test
# diverged that way on a correct artifact.
observe() {
    local out rc=0
    out=$("$@" 2>"$TMP/stderr") || rc=$?
    printf '%s\n--- stderr ---\n' "$out"
    cat "$TMP/stderr"
    return "$rc"
}

# build_artifact MODE SRC OUT PATHS
#
# Every compile mode EXECUTES the program as a side effect of compiling it, so
# its stdout is discarded: comparing that against the direct run would compare
# two source-level executions and prove nothing about the artifact.
build_artifact() {
    local mode="$1" src="$2" out="$3" paths="$4" flags=()
    case "$mode" in
        c)     flags=(-c "$out") ;;
        cz)    flags=(-z -c "$out") ;;
        strip) flags=(-strip -c "$out") ;;
        b)     flags=(-b "$out") ;;
    esac
    (cd "$ROOT" && LG_SOURCE_PATHS="$paths" ${BOUND[@]+"${BOUND[@]}"} "$LG" "${flags[@]}" "$src" >/dev/null 2>&1)
}

# run_artifact MODE OUT PATHS
run_artifact() {
    local mode="$1" out="$2" paths="$3"
    if [ "$mode" = b ]; then
        (cd "$ROOT" && observe env LG_SOURCE_PATHS="$paths" ${BOUND[@]+"${BOUND[@]}"} "$out")
    else
        (cd "$ROOT" && observe env LG_SOURCE_PATHS="$paths" ${BOUND[@]+"${BOUND[@]}"} "$LG" "$out")
    fi
}

run_fixture() {
    local label="$1" src="$2" paths="$3"
    local a1 a2 b rc=0 reason mode out

    if reason=$(skip_reason "$label"); then
        printf "  %-44s SKIP (%s)\n" "$label" "$reason"
        skipped=$((skipped+1)); return 0
    fi

    # Self-control: two direct runs.
    a1=$(cd "$ROOT" && observe env LG_SOURCE_PATHS="$paths" ${BOUND[@]+"${BOUND[@]}"} "$LG" "$src") || rc=$?
    if [ "$rc" -ne 0 ]; then
        printf "  %-44s ERROR (direct run exit %s)\n" "$label" "$rc"
        errored=$((errored+1)); return 0
    fi
    a2=$(cd "$ROOT" && observe env LG_SOURCE_PATHS="$paths" ${BOUND[@]+"${BOUND[@]}"} "$LG" "$src") || true
    if [ "$a1" != "$a2" ]; then
        printf "  %-44s NONDET (excluded — not an oracle)\n" "$label"
        nondet=$((nondet+1)); return 0
    fi

    for mode in $MODES; do
        # lg picks the .lgb loader by extension, so every non-executable
        # artifact must end in .lgb or it is run as source.
        out="$TMP/$(echo "$label" | tr '/.' '__').$mode"
        [ "$mode" = b ] || out="$out.lgb"
        if ! build_artifact "$mode" "$src" "$out" "$paths"; then
            printf "  %-44s ERROR (%s: compile failed)\n" "$label" "$mode"
            errored=$((errored+1)); continue
        fi
        if ! b=$(run_artifact "$mode" "$out" "$paths"); then
            printf "  %-44s DIVERGED (%s: artifact failed to run)\n" "$label" "$mode"
            diverged=$((diverged+1)); DIVERGED_LIST+=("$label [$mode] (load/run failure)"); continue
        fi
        if [ "$a1" = "$b" ]; then
            agree=$((agree+1))
        else
            printf "  %-44s DIVERGED (%s)\n" "$label" "$mode"
            # `|| true` is load-bearing: diff exits 1 when the inputs differ,
            # and under `set -e -o pipefail` that would abort the whole run on
            # the FIRST divergence instead of reporting every one.
            diff <(printf '%s\n' "$a1") <(printf '%s\n' "$b") | sed 's/^/        /' | head -12 || true
            diverged=$((diverged+1)); DIVERGED_LIST+=("$label [$mode]")
        fi
    done
}

# --- Selftest -------------------------------------------------------------
#
# Falsifies the checker rather than exercising it: asserts that an ordinary
# fixture agrees in every mode (so the checker cannot pass by being
# always-red), that a deliberately nondeterministic one is classified NONDET
# rather than silently passed, and that a wrong artifact is caught in every
# mode. Needs no corpus.
if [ "$SELFTEST" = 1 ]; then
    fails=0
    nmodes=$(set -- $MODES; echo $#)
    expect() { # expect <case> <got> <want>
        if [ "$2" = "$3" ]; then echo "  PASS  $1"
        else echo "  FAIL  $1 (got $2, want $3)" >&2; fails=$((fails+1)); fi
    }

    cat >"$TMP/ok.lg" <<'EOF'
(ns selftest-ok)
(println "answer:" (+ 1 2 3))
(println "seq:" (vec (map inc [1 2 3])))
EOF
    reset_counts
    run_fixture "selftest/ok" "$TMP/ok.lg" "$ROOT" >/dev/null
    expect "a deterministic program round-trips in every mode" "$agree" "$nmodes"

    cat >"$TMP/nd.lg" <<'EOF'
(ns selftest-nd)
(println "now:" (System/nanoTime))
EOF
    reset_counts
    run_fixture "selftest/nondet" "$TMP/nd.lg" "$ROOT" >/dev/null
    expect "a nondeterministic program is excluded as NONDET" "$nondet" 1
    expect "...and is NOT counted as agreement" "$agree" 0

    # A builder that emits a DIFFERENT program's artifact stands in for a
    # decoder that reconstructs the wrong module. Only the build is swapped:
    # the direct runs and each mode's real run path are untouched, so the two
    # sides cannot change together and agree for the wrong reason.
    cat >"$TMP/other.lg" <<'EOF'
(ns selftest-other)
(println "answer:" 999)
EOF
    eval "real_$(declare -f build_artifact)"
    build_artifact() { real_build_artifact "$1" "$TMP/other.lg" "$3" "$4"; }
    reset_counts
    run_fixture "selftest/divergent" "$TMP/ok.lg" "$ROOT" >/dev/null
    expect "a wrong artifact is caught as DIVERGED in every mode" "$diverged" "$nmodes"

    [ "$fails" -eq 0 ] || { echo "SELFTEST FAILED: $fails case(s)." >&2; exit 1; }
    echo "selftest: all cases passed (modes: $MODES)."
    exit 0
fi

# --- Corpus ---------------------------------------------------------------
if [ "${#FILES[@]}" -gt 0 ]; then
    echo "=== files ==="
    for f in "${FILES[@]}"; do
        [ -f "$f" ] || { echo "no such file: $f" >&2; exit 2; }
        abs="$(cd "$(dirname "$f")" && pwd)/$(basename "$f")"
        run_fixture "$f" "$abs" "$ROOT:$(dirname "$abs")"
    done
else
    echo "=== examples ==="
    for f in "$ROOT"/examples/*.lg; do
        [ -f "$f" ] || continue
        run_fixture "examples/$(basename "$f")" "$f" "$ROOT:$ROOT/examples"
    done
fi

if [ "$FULL" = 1 ]; then
    echo
    echo "=== test suites ==="
    # Each suite file declares a namespace but defines tests rather than
    # running them, so a driver is generated per suite: require the namespace,
    # then run its tests and print the summary map. The summary is the
    # observable output being compared, which makes an assertion that passes
    # under one path and fails under the other a visible divergence.
    for f in "$ROOT"/test/*_test.lg; do
        [ -f "$f" ] || continue
        ns=$(grep -m1 -oE '^\(ns +[A-Za-z0-9._*+!?<>=-]+' "$f" | awk '{print $2}') || true
        [ -n "$ns" ] || continue
        drv="$TMP/drv-$(basename "$f")"
        # `run-tests <ns>`, never `run-all-tests`. The latter sweeps every
        # LOADED namespace in hash-map order, so its "Testing <ns>" headers come
        # out shuffled between runs while the results stay identical, which the
        # self-control reports as nondeterminism in most suites.
        printf '(ns drv-%s (:require [test :refer :all] [%s]))\n(println "RESULT:" (run-tests (quote %s)))\n' \
            "$(basename "$f" .lg)" "$ns" "$ns" >"$drv"
        run_fixture "test/$(basename "$f")" "$drv" "$ROOT:$ROOT/test:$ROOT/pkg/rt/core"
    done
fi

echo
echo "=== round-trip (modes: $MODES) ==="
printf "  agree=%d diverged=%d nondet=%d errored=%d skipped=%d\n" \
    "$agree" "$diverged" "$nondet" "$errored" "$skipped"

# A run that verified nothing is not a pass. Without this the harness would
# report success on a corpus that had drifted to entirely NONDET or ERROR.
if [ "$agree" -eq 0 ]; then
    echo "FAIL: no fixture round-tripped — nothing was verified." >&2
    exit 1
fi
if [ "$diverged" -ne 0 ]; then
    echo "FAIL: $diverged (fixture, mode) pair(s) diverged:" >&2
    printf '  - %s\n' "${DIVERGED_LIST[@]}" >&2
    exit 1
fi
echo "OK: $agree (fixture, mode) pair(s) round-trip identically."
