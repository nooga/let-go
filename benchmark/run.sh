#!/usr/bin/env bash
set -euo pipefail

# Benchmark let-go (bytecode VM) against let-go AOT (IR-lowered to native Go),
# babashka, and Clojure JVM.
# Requires: hyperfine, bb, clj, go, python3
#
# The VM and AOT legs run the IDENTICAL driver + fixture (benchmark/aot/) —
# only the binary differs: the plain build dispatches bytecode, the lg_bench-
# tagged build (produced by benchmark/aot/build.lg) dispatches the same vars
# to IR-lowered native Go. babashka and Clojure run the plain .clj scripts,
# which contain the same workloads.

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

# Filter mode: positional args select which perf benches to run.
# In filter mode, startup/memory are skipped and results.md is NOT regenerated;
# results are printed only. Use this for iterating on a single bench.
FILTER_BENCHES=("$@")
FILTER_MODE=0
[ ${#FILTER_BENCHES[@]} -gt 0 ] && FILTER_MODE=1

if [ "$FILTER_MODE" -eq 1 ]; then
    WARMUP=10
    RUNS=50
else
    WARMUP=3
    RUNS=10
fi

# --- Build both let-go binaries ---
echo "Building let-go (VM)..."
cd "$PROJECT_DIR"
go build -ldflags="-s -w" -o "$SCRIPT_DIR/letgo" .

echo "Building let-go (AOT)..."
"$SCRIPT_DIR/letgo" "$SCRIPT_DIR/aot/build.lg" --out "benchmark/letgo-aot"

LETGO="$SCRIPT_DIR/letgo"
LETGO_AOT="$SCRIPT_DIR/letgo-aot"
AOT_SRC="$SCRIPT_DIR/aot/src"
BB="$(which bb 2>/dev/null || true)"
CLJ="$(which clj 2>/dev/null || true)"

# Collect benchmark files (apply filter if any positional args were passed)
ALL_BENCHMARKS=($(ls "$SCRIPT_DIR"/*.clj 2>/dev/null | sort))
BENCHMARKS=()
if [ "$FILTER_MODE" -eq 1 ]; then
    for want in "${FILTER_BENCHES[@]}"; do
        match="$SCRIPT_DIR/${want}.clj"
        if [ -f "$match" ]; then
            BENCHMARKS+=("$match")
        else
            echo "Unknown benchmark: $want (no $match)" >&2
            exit 1
        fi
    done
    echo "Filter mode: running ${FILTER_BENCHES[*]} (warmup=$WARMUP runs=$RUNS, no results.md update)"
else
    BENCHMARKS=("${ALL_BENCHMARKS[@]}")
fi

if [ ${#BENCHMARKS[@]} -eq 0 ]; then
    echo "No benchmark files found in $SCRIPT_DIR"
    exit 1
fi

# Per-benchmark commands for the two let-go legs (identical driver + fixture).
lg_vm_cmd()  { echo "$LETGO -source-paths $AOT_SRC $SCRIPT_DIR/aot/drivers/$1.lg"; }
lg_aot_cmd() { echo "$LETGO_AOT -source-paths $AOT_SRC $SCRIPT_DIR/aot/drivers/$1.lg"; }

# --- Gather system info ---

LETGO_SIZE=$(ls -lh "$LETGO" | awk '{print $5}')
LETGO_AOT_SIZE=$(ls -lh "$LETGO_AOT" | awk '{print $5}')
BB_VERSION=""
BB_SIZE=""
CLJ_VERSION=""
JDK_VERSION=""
JDK_SIZE=""

if [ -n "$BB" ]; then
    BB_VERSION="$(bb --version 2>&1)"
    BB_PATH="$(readlink -f "$BB" 2>/dev/null || readlink "$BB" 2>/dev/null || echo "$BB")"
    if [[ "$BB_PATH" == ../* ]]; then
        BB_PATH="$(cd "$(dirname "$BB")" && cd "$(dirname "$BB_PATH")" && pwd)/$(basename "$BB_PATH")"
    fi
    BB_SIZE="$(ls -lh "$BB_PATH" 2>/dev/null | awk '{print $5}')"
fi

if [ -n "$CLJ" ]; then
    CLJ_VERSION="$(clj --version 2>&1)"
    JDK_VERSION="$(java -version 2>&1 | head -1)"
    JAVA_HOME_DIR="$(/usr/libexec/java_home 2>/dev/null || echo "")"
    if [ -n "$JAVA_HOME_DIR" ]; then
        JDK_SIZE="$(du -sh "$JAVA_HOME_DIR" 2>/dev/null | awk '{print $1}')"
    fi
fi

CPU_INFO="$(sysctl -n machdep.cpu.brand_string 2>/dev/null || cat /proc/cpuinfo 2>/dev/null | grep 'model name' | head -1 | cut -d: -f2 | xargs || echo "unknown")"

echo ""
echo "=== Environment ==="
echo "System: $(uname -ms), $CPU_INFO"
echo "let-go VM: $LETGO_SIZE binary | let-go AOT: $LETGO_AOT_SIZE binary"
[ -n "$BB" ] && echo "babashka: $BB_VERSION ($BB_SIZE binary)"
[ -n "$CLJ" ] && echo "clojure: $CLJ_VERSION"
[ -n "$CLJ" ] && echo "JDK: $JDK_VERSION ($JDK_SIZE)"
echo ""

# --- Measure startup time ---

if [ "$FILTER_MODE" -eq 0 ]; then
echo "=== Startup Time ==="
STARTUP_JSON="/tmp/bench_startup.json"
STARTUP_ARGS=(--warmup "$WARMUP" --runs "$RUNS" --export-json "$STARTUP_JSON")
STARTUP_ARGS+=("-n" "let-go" "$LETGO -e nil")
STARTUP_ARGS+=("-n" "let-go AOT" "$LETGO_AOT -e nil")
[ -n "$BB" ] && STARTUP_ARGS+=("-n" "babashka" "bb -e nil")
[ -n "$CLJ" ] && STARTUP_ARGS+=("-n" "clojure" "clj -M -e nil")
hyperfine "${STARTUP_ARGS[@]}" 2>&1

# --- Measure peak memory ---

echo ""
echo "=== Peak Memory ==="
measure_mem() {
    local cmd="$1"
    local name="$2"
    # Run 3 times, take median
    local mems=()
    for i in 1 2 3; do
        local mem=$(/usr/bin/time -l $cmd 2>&1 >/dev/null | grep "maximum resident" | awk '{print $1}')
        mems+=($mem)
    done
    IFS=$'\n' sorted=($(sort -n <<<"${mems[*]}")); unset IFS
    local median_bytes=${sorted[1]}
    local median_mb=$(echo "scale=1; $median_bytes / 1048576" | bc)
    echo "  $name: ${median_mb}MB"
    echo "$median_mb"
}

echo "Running: nil (startup only)"
LETGO_STARTUP_MEM=$(measure_mem "$LETGO -e nil" "let-go" | tail -1)
LETGO_AOT_STARTUP_MEM=$(measure_mem "$LETGO_AOT -e nil" "let-go AOT" | tail -1)
BB_STARTUP_MEM=""
[ -n "$BB" ] && BB_STARTUP_MEM=$(measure_mem "bb -e nil" "babashka" | tail -1)
CLJ_STARTUP_MEM=""
[ -n "$CLJ" ] && CLJ_STARTUP_MEM=$(measure_mem "clj -M -e nil" "clojure" | tail -1)

echo ""
echo "Running: fib 35 (compute-heavy)"
LETGO_FIB_MEM=$(measure_mem "$(lg_vm_cmd fib)" "let-go" | tail -1)
LETGO_AOT_FIB_MEM=$(measure_mem "$(lg_aot_cmd fib)" "let-go AOT" | tail -1)
BB_FIB_MEM=""
[ -n "$BB" ] && BB_FIB_MEM=$(measure_mem "bb benchmark/fib.clj" "babashka" | tail -1)
CLJ_FIB_MEM=""
[ -n "$CLJ" ] && CLJ_FIB_MEM=$(measure_mem "clj -M -e '(load-file \"benchmark/fib.clj\")'" "clojure" | tail -1)

echo ""
echo "Running: reduce 1M (large collection)"
LETGO_REDUCE_MEM=$(measure_mem "$(lg_vm_cmd reduce)" "let-go" | tail -1)
LETGO_AOT_REDUCE_MEM=$(measure_mem "$(lg_aot_cmd reduce)" "let-go AOT" | tail -1)
BB_REDUCE_MEM=""
[ -n "$BB" ] && BB_REDUCE_MEM=$(measure_mem "bb benchmark/reduce.clj" "babashka" | tail -1)
CLJ_REDUCE_MEM=""
[ -n "$CLJ" ] && CLJ_REDUCE_MEM=$(measure_mem "clj -M -e '(load-file \"benchmark/reduce.clj\")'" "clojure" | tail -1)

fi  # end FILTER_MODE == 0 (startup + memory block)

# --- Run benchmarks ---

echo ""
echo "=== Performance Benchmarks ==="

for bench in "${BENCHMARKS[@]}"; do
    name="$(basename "$bench" .clj)"
    echo ""
    echo "--- $name ---"

    JSON="/tmp/bench_${name}.json"
    CMDS=("-n" "let-go" "$(lg_vm_cmd "$name")")
    CMDS+=("-n" "let-go AOT" "$(lg_aot_cmd "$name")")
    [ -n "$BB" ] && CMDS+=("-n" "babashka" "bb $bench")
    [ -n "$CLJ" ] && CMDS+=("-n" "clojure" "clj -M -e '(load-file \"$bench\")'")

    hyperfine --warmup "$WARMUP" --runs "$RUNS" --export-json "$JSON" "${CMDS[@]}" 2>&1
done

# --- Generate results.md ---

if [ "$FILTER_MODE" -eq 1 ]; then
    echo ""
    echo "=== Done (filter mode — results.md not updated) ==="
    rm -f "$SCRIPT_DIR/letgo" "$SCRIPT_DIR/letgo-aot"
    exit 0
fi

RESULTS_FILE="$SCRIPT_DIR/results.md"

cat > "$RESULTS_FILE" << EOF
## Benchmark Results

### Methodology

All benchmarks use [hyperfine](https://github.com/sharkdp/hyperfine) with $WARMUP warmup runs
and $RUNS timed runs per benchmark. Times shown are mean ± σ wall-clock time. Peak memory is
measured via \`/usr/bin/time -l\` (median of 3 runs).

**let-go** and **let-go AOT** run the identical driver + fixture (\`benchmark/aot/\`): the
same program, through the same binary entry point — the only difference is dispatch. The
plain build executes bytecode on the VM; the AOT build (produced by
\`benchmark/aot/build.lg\`) carries the same namespaces IR-lowered to native Go, registered
as overrides, so the same \`(bench.<name>/run)\` call dispatches to compiled Go. Wall-clock
includes binary startup and namespace load for both legs. babashka and Clojure JVM run the
plain \`.clj\` scripts, which contain the same workloads as the fixtures.

Clojure JVM times include full JVM startup (~350-500ms) which dominates short benchmarks.

**System:** $(uname -ms), $CPU_INFO

**Runtimes:**

| | let-go | let-go AOT | babashka | clojure JVM |
|---|---|---|---|---|
| **Version** | — | — | ${BB_VERSION:-—} | ${CLJ_VERSION:-—} |
| **Platform** | Go bytecode VM | IR-lowered native Go | GraalVM native | JVM (HotSpot) |
| **Binary/runtime size** | **$LETGO_SIZE** | $LETGO_AOT_SIZE | ${BB_SIZE:-—} | ${JDK_SIZE:-—} |

### Startup Time

EOF

python3 -c "
import json
with open('$STARTUP_JSON') as f:
    d = json.load(f)

def fmt(mean, stddev):
    if mean < 0.1:
        return f'{mean*1000:.1f}ms ± {stddev*1000:.1f}ms'
    return f'{mean:.3f}s ± {stddev:.3f}s'

entries = []
for r in d['results']:
    name = r['command'].strip()
    if name == 'clojure': name = 'clojure JVM'
    entries.append((name, r['mean'], r['stddev']))

lg_mean = next(e[1] for e in entries if e[0] == 'let-go')
best = min(e[1] for e in entries)
print('| Runtime | Time |')
print('|---|---|')
for name, mean, stddev in entries:
    s = fmt(mean, stddev)
    ratio = mean / lg_mean
    tag = f' ({ratio:.1f}x)' if name != 'let-go' else ' (1.0x)'
    if mean == best:
        print(f'| **{name}** | **{s}**{tag} |')
    else:
        print(f'| {name} | {s}{tag} |')
" >> "$RESULTS_FILE"

python3 -c "
def parse_mb(s):
    try: return float(s)
    except: return None

def row(label, vals):
    nums = [v[1] for v in vals if v[1] is not None]
    best = min(nums) if nums else None
    lg = vals[0][1] if vals[0][1] is not None else 1
    cells = []
    for i, (s, n) in enumerate(vals):
        if n is not None and lg:
            tag = f' ({n/lg:.1f}x)' if i > 0 else ' (1.0x)'
        else:
            tag = ''
        cells.append((f'**{s}**' if n is not None and n == best else s) + tag)
    return f'| {label} | ' + ' | '.join(cells) + ' |'

rows = [
    ('startup (nil)', [
        ('${LETGO_STARTUP_MEM}MB', parse_mb('${LETGO_STARTUP_MEM}')),
        ('${LETGO_AOT_STARTUP_MEM}MB', parse_mb('${LETGO_AOT_STARTUP_MEM}')),
        ('${BB_STARTUP_MEM:-—}MB', parse_mb('${BB_STARTUP_MEM}')),
        ('${CLJ_STARTUP_MEM:-—}MB', parse_mb('${CLJ_STARTUP_MEM}')),
    ]),
    ('fib(35)', [
        ('${LETGO_FIB_MEM}MB', parse_mb('${LETGO_FIB_MEM}')),
        ('${LETGO_AOT_FIB_MEM}MB', parse_mb('${LETGO_AOT_FIB_MEM}')),
        ('${BB_FIB_MEM:-—}MB', parse_mb('${BB_FIB_MEM}')),
        ('${CLJ_FIB_MEM:-—}MB', parse_mb('${CLJ_FIB_MEM}')),
    ]),
    ('reduce 1M', [
        ('${LETGO_REDUCE_MEM}MB', parse_mb('${LETGO_REDUCE_MEM}')),
        ('${LETGO_AOT_REDUCE_MEM}MB', parse_mb('${LETGO_AOT_REDUCE_MEM}')),
        ('${BB_REDUCE_MEM:-—}MB', parse_mb('${BB_REDUCE_MEM}')),
        ('${CLJ_REDUCE_MEM:-—}MB', parse_mb('${CLJ_REDUCE_MEM}')),
    ]),
]

print()
print('### Peak Memory Usage (RSS)')
print()
print('| Workload | let-go | let-go AOT | babashka | clojure JVM |')
print('|---|---|---|---|---|')
for label, vals in rows:
    print(row(label, vals))
" >> "$RESULTS_FILE"

cat >> "$RESULTS_FILE" << EOF

### Performance

EOF

{
echo "| Benchmark | let-go | let-go AOT | babashka | clojure JVM |"
echo "|---|---|---|---|---|"

for bench in "${BENCHMARKS[@]}"; do
    name="$(basename "$bench" .clj)"
    JSON="/tmp/bench_${name}.json"
    python3 -c "
import json
with open('$JSON') as f:
    d = json.load(f)

def fmt(mean, stddev):
    if mean < 0.1:
        return f'{mean*1000:.1f}ms ± {stddev*1000:.1f}ms'
    return f'{mean:.3f}s ± {stddev:.3f}s'

results = {}
for r in d['results']:
    cmd = r['command'].strip()
    if cmd == 'let-go': results['letgo'] = (fmt(r['mean'], r['stddev']), r['mean'])
    elif cmd == 'let-go AOT': results['aot'] = (fmt(r['mean'], r['stddev']), r['mean'])
    elif cmd == 'babashka': results['bb'] = (fmt(r['mean'], r['stddev']), r['mean'])
    elif cmd == 'clojure': results['clj'] = (fmt(r['mean'], r['stddev']), r['mean'])

best = min(v[1] for v in results.values())
lg_mean = results.get('letgo', (None, 1.0))[1]

def cell(key):
    if key not in results:
        return chr(0x2014)
    s, mean = results[key]
    ratio = mean / lg_mean if lg_mean else 0
    tag = f' ({ratio:.1f}x)' if key != 'letgo' else ' (1.0x)'
    if mean == best:
        return f'**{s}**{tag}'
    return f'{s}{tag}'

print('| $name | ' + cell('letgo') + ' | ' + cell('aot') + ' | ' + cell('bb') + ' | ' + cell('clj') + ' |')
"
done
} >> "$RESULTS_FILE"

echo "" >> "$RESULTS_FILE"

echo ""
echo "=== Done ==="
echo "Results written to $RESULTS_FILE"
echo ""
cat "$RESULTS_FILE"

# Cleanup
rm -f "$SCRIPT_DIR/letgo" "$SCRIPT_DIR/letgo-aot"
