#!/usr/bin/env bash
# gogen-parity.sh — parity benchmark for the gogen_ir bootstrap.
#
# Runs end-to-end correctness + perf checks under both the untagged build
# (bytecode-replayed Lisp IR stack) and `-tags gogen_ir` (native Go IR
# stack from pkg/rt/core_go_lowered). Prints a side-by-side summary;
# returns non-zero if any suite diverges.
#
# Three suites cover different reach:
#   1. clojure-test-suite (jank)  → end-user-observable clojure.core semantics
#   2. ir-stress lower-go         → AOT lowering of IR → Go correctness
#   3. ir-stress ir-compile       → eval-mode IR compile → bytecode correctness
#
# Usage:
#   scripts/gogen-parity.sh             # default: jank + lower-go (~3 min)
#   scripts/gogen-parity.sh --quick     # jank only (~2 sec, smoke check)
#   scripts/gogen-parity.sh --full      # all three suites (~5 min)
#   scripts/gogen-parity.sh --jank-only # just the jank suite
#
# Exit codes:
#   0  all suites parity-identical (counts + buckets)
#   1  semantic divergence detected (counts or buckets differ)
#   2  setup error (missing submodule, no go binary, etc.)
#
# The wall-time delta is reported but not enforced (single-run noise is
# typically ±2%; treat as informational unless consistently >10%).

set -euo pipefail

# Bound the Go heap for the bootstrap `go run`/`go test` invocations below.
# These compile the whole .lg stdlib from source and balloon the heap;
# uncapped they OOM a 16GB machine. Soft cap — runtime GCs to stay under.
# Honors an outer GOMEMLIMIT (e.g. from the Makefile export) if already set.
export GOMEMLIMIT="${GOMEMLIMIT:-2GiB}"

MODE="${1:-default}"
case "$MODE" in
    --quick)      RUN_JANK=1; RUN_LOWERGO=0; RUN_IRCOMPILE=0; RUN_DEFTYPE=0 ;;
    --jank-only)  RUN_JANK=1; RUN_LOWERGO=0; RUN_IRCOMPILE=0; RUN_DEFTYPE=0 ;;
    # ir-compile only: the IR-optimizing bytecode path (binds *ir-compile*),
    # which runs the native passes under -tags gogen_ir and is byte-stable
    # across engines — unlike lower-go, whose AOT run trips the wall-clock
    # *typeinfer-budget-ms* and flakes. This is the CI-safe parity gate.
    --ir-compile) RUN_JANK=0; RUN_LOWERGO=0; RUN_IRCOMPILE=1; RUN_DEFTYPE=0 ;;
    --full)       RUN_JANK=1; RUN_LOWERGO=1; RUN_IRCOMPILE=1; RUN_DEFTYPE=1 ;;
    default)      RUN_JANK=1; RUN_LOWERGO=1; RUN_IRCOMPILE=0; RUN_DEFTYPE=1 ;;
    -h|--help)
        sed -n '2,/^$/p' "$0" | sed 's/^# \?//'
        exit 0
        ;;
    *)
        echo "unknown mode: $MODE (try --help)" >&2
        exit 2
        ;;
esac

# The gogen_ir lowered tree (pkg/rt/core_go_lowered/) is a gitignored build
# artifact, so regenerate it before any `-tags gogen_ir` build below. Cheap
# relative to the parity runs; non-determinism is irrelevant here (we compare
# program output across engines, not the generated Go bytes).
go run -tags bootstrap ./cmd/lgbgen --target=go >/dev/null

# IR corpus used by ir-stress. Order matters (data.lg must load first).
IR_CORPUS=(
    data.lg build.lg lower.lg lower_go.lg dump.lg
    dominance.lg validate.lg zipper.lg passes.lg
)
IR_DIR=pkg/rt/core/ir

LOG_DIR="${TMPDIR:-/tmp}/gogen-parity-$$"
mkdir -p "$LOG_DIR"
trap 'rc=$?; [ $rc -eq 0 ] && rm -rf "$LOG_DIR" || echo "logs preserved at $LOG_DIR"' EXIT

# Preflight ---------------------------------------------------------------

command -v go >/dev/null 2>&1 || { echo "go binary not on PATH" >&2; exit 2; }

if [ "$RUN_JANK" -eq 1 ]; then
    if [ ! -d test/clojure-test-suite/test/clojure/core_test ]; then
        echo "clojure-test-suite submodule not initialized in this worktree." >&2
        echo "In a jj worktree the submodule data isn't auto-populated. Either:" >&2
        echo "  git submodule update --init                      # in a git worktree" >&2
        echo "  ln -s <other-worktree>/test/clojure-test-suite \\" >&2
        echo "        test/clojure-test-suite                    # in a jj worktree" >&2
        exit 2
    fi
fi

# Helpers -----------------------------------------------------------------

# Run a command, capture stdout+stderr to $1, append wall time to log. The
# stdout of the function itself is the recorded wall time in seconds. The
# inner command's REAL exit code is written to "$logfile.rc" — without this
# the `wall=$(time_run ...)` substitution would only surface awk's exit (0),
# silently swallowing a compile failure / panic / timeout in the suite. The
# `if/else` keeps the failure from tripping `set -e` while still capturing rc.
time_run() {
    local logfile="$1"; shift
    local start end rc
    start=$(date +%s.%N)
    if "$@" >"$logfile" 2>&1; then rc=0; else rc=$?; fi
    end=$(date +%s.%N)
    echo "$rc" >"$logfile.rc"
    awk -v s="$start" -v e="$end" 'BEGIN{printf "%.2f\n", e-s}'
}

# Extract the TOTALS line from a jank-suite -v log.
jank_totals() { grep -m1 -E "TOTALS: " "$1" | sed 's/.*TOTALS: //'; }

# Extract the Passed/Failed line + bucket distribution from an ir-stress log.
ir_summary()  { grep -E "^Total fixtures:|^Passed:|^Failed:" "$1"; }
ir_buckets()  {
    # Scrub process-specific noise (pointer addresses in printed fn/closure
    # values) so the bucket hash captures semantic equivalence, not memory
    # layout. Without this, an error string containing `<fn foo 0xabc123>`
    # produces a different md5 every run for purely cosmetic reasons.
    awk '/^=== Failure Buckets ===/{p=1; next} p && /^=== /{exit} p' "$1" |
        sed -E 's/0x[0-9a-fA-F]{6,}/0xXXX/g'
}

# require_ran asserts a suite ACTUALLY RAN before its output is trusted.
#   $1 human label   $2 log file   $3 grep -E pattern the summary MUST contain
# A run that compiles, panics, times out, or matches 0 tests never emits its
# summary line; without this guard its empty summary would later compare
# "empty == empty" across engines and be reported as PARITY — a false pass.
# This is the structural fix for the class of bug where success is inferred
# from the absence of output (or from a pipe's exit code) rather than from a
# positive "N tests ran" signal. Exits 2 (run/setup error), distinct from the
# exit 1 used for genuine semantic divergence.
require_ran() {
    local label="$1" log="$2" pat="$3"
    local rc; rc=$(cat "$log.rc" 2>/dev/null || echo "?")
    if ! grep -qE "$pat" "$log"; then
        {
            echo "FATAL: $label produced no '$pat' summary line (inner exit=$rc)."
            echo "       The suite failed to compile, panicked, timed out, or ran 0 tests."
            echo "       This is a RUN failure, NOT a parity result — do not trust an OK."
            echo "       --- last 20 lines of $log ---"
            tail -20 "$log"
        } >&2
        trap - EXIT
        exit 2
    fi
}

# Suites ------------------------------------------------------------------

run_jank() {
    local tag_label="$1" tag_flag="$2"
    local log="$LOG_DIR/jank-${tag_label}.log"
    local wall
    wall=$(time_run "$log" go test $tag_flag ./test/ -run '^TestClojureTestSuite$' -count=1 -v -timeout 120s)
    require_ran "jank/$tag_label" "$log" 'TOTALS: '
    printf "%-10s %-50s %ss\n" "$tag_label" "$(jank_totals "$log")" "$wall"
    echo "$wall:$(jank_totals "$log")" >"$LOG_DIR/jank-${tag_label}.summary"
}

run_ir_stress() {
    local mode="$1" tag_label="$2" tag_flag="$3"
    local log="$LOG_DIR/${mode}-${tag_label}.log"
    local wall
    wall=$(time_run "$log" go run $tag_flag . scripts/ir-stress.lg "$mode" "$IR_DIR" "${IR_CORPUS[@]}")
    require_ran "${mode}/$tag_label" "$log" '^Passed:'
    local pass fail buckets
    pass=$(grep -m1 "^Passed:" "$log" | awk '{print $2}')
    fail=$(grep -m1 "^Failed:" "$log" | awk '{print $2}')
    buckets=$(ir_buckets "$log" | md5sum | cut -c1-8)
    printf "%-10s pass=%s fail=%s buckets=%s  %ss\n" "$tag_label" "$pass" "$fail" "$buckets" "$wall"
    echo "$wall:$pass:$fail:$buckets" >"$LOG_DIR/${mode}-${tag_label}.summary"
}

compare_summaries() {
    local label="$1" untagged_file="$2" tagged_file="$3"
    local u t
    u=$(cat "$untagged_file"); t=$(cat "$tagged_file")
    # Strip leading wall time; compare the rest.
    local u_body t_body
    u_body=${u#*:}; t_body=${t#*:}
    if [ "$u_body" = "$t_body" ]; then
        echo "  $label: PARITY"
        return 0
    else
        echo "  $label: DIVERGED"
        echo "    untagged: $u"
        echo "    gogen_ir: $t"
        return 1
    fi
}

# deftype/defprotocol native lowering ------------------------------------
# Delegates to the generalized trampoline (scripts/gogen-trampoline.lg), which
# lowers every test/gogen/*.lg fixture to Go, wires it in under a build tag, and
# `require`s it so the resolver drains the Go-native override — then asserts each
# fixture's (run) matches the bytecode VM AND actually dispatched natively. This
# is the real runtime dispatch path, not a hand-written in-package call.
run_deftype_native() {
    local log="$LOG_DIR/deftype-native.log" lg="$LOG_DIR/lg-trampoline"
    # The trampoline is an lg script; build lg once, run it, and pass the same
    # binary as the bytecode-side engine it compares against.
    if go build -o "$lg" . >"$log" 2>&1 && \
       "$lg" scripts/gogen-trampoline.lg --lg "$lg" --go "$(command -v go)" >>"$log" 2>&1; then
        sed -nE 's/^  //p' "$log"   # echo the per-fixture result lines
        return 0
    fi
    echo "  deftype native trampoline FAILED — last 20 lines of $log:" >&2
    tail -20 "$log" >&2
    return 1
}

# Run --------------------------------------------------------------------

divergence=0

if [ "$RUN_JANK" -eq 1 ]; then
    echo "=== clojure-test-suite (jank) ==="
    run_jank untagged ""
    run_jank gogen_ir "-tags gogen_ir"
    echo
fi

if [ "$RUN_LOWERGO" -eq 1 ]; then
    echo "=== ir-stress lower-go (IR corpus: ${#IR_CORPUS[@]} files) ==="
    run_ir_stress lower-go untagged ""
    run_ir_stress lower-go gogen_ir "-tags gogen_ir"
    echo
fi

if [ "$RUN_IRCOMPILE" -eq 1 ]; then
    echo "=== ir-stress ir-compile (IR corpus: ${#IR_CORPUS[@]} files) ==="
    run_ir_stress ir-compile untagged ""
    run_ir_stress ir-compile gogen_ir "-tags gogen_ir"
    echo
fi

if [ "${RUN_DEFTYPE:-0}" -eq 1 ]; then
    echo "=== deftype/defprotocol native skeleton ==="
    run_deftype_native || divergence=$((divergence+1))
    echo
fi

echo "=== Parity check ==="
check() {
    local label="$1" untagged="$2" tagged="$3"
    if ! compare_summaries "$label" "$untagged" "$tagged"; then
        divergence=$((divergence+1))
    fi
}
[ "$RUN_JANK"      -eq 1 ] && check jank       "$LOG_DIR/jank-untagged.summary"       "$LOG_DIR/jank-gogen_ir.summary"
[ "$RUN_LOWERGO"   -eq 1 ] && check lower-go   "$LOG_DIR/lower-go-untagged.summary"   "$LOG_DIR/lower-go-gogen_ir.summary"
[ "$RUN_IRCOMPILE" -eq 1 ] && check ir-compile "$LOG_DIR/ir-compile-untagged.summary" "$LOG_DIR/ir-compile-gogen_ir.summary"

if [ "$divergence" -ne 0 ]; then
    echo "FAIL: $divergence suite(s) diverged. Logs at $LOG_DIR." >&2
    trap - EXIT
    exit 1
fi

echo "OK: all suites parity-identical."
