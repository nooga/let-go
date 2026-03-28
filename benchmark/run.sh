#!/usr/bin/env bash
set -euo pipefail

# Benchmark let-go against babashka and Clojure JVM
# Requires: hyperfine, bb, clj, go, python3

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

WARMUP=3
RUNS=10

# Build let-go binary
echo "Building let-go..."
cd "$PROJECT_DIR"
go build -ldflags="-s -w" -o "$SCRIPT_DIR/letgo" .

LETGO="$SCRIPT_DIR/letgo"
BB="$(which bb 2>/dev/null || true)"
CLJ="$(which clj 2>/dev/null || true)"
JOKER="$(which joker 2>/dev/null || true)"

# Benchmarks to skip for specific runtimes (too slow or unsupported)
JOKER_SKIP="tak transducers"

# Collect benchmark files
BENCHMARKS=($(ls "$SCRIPT_DIR"/*.clj 2>/dev/null | sort))

if [ ${#BENCHMARKS[@]} -eq 0 ]; then
    echo "No benchmark files found in $SCRIPT_DIR"
    exit 1
fi

# --- Gather system info ---

LETGO_SIZE=$(ls -lh "$LETGO" | awk '{print $5}')
BB_VERSION=""
BB_SIZE=""
CLJ_VERSION=""
JDK_VERSION=""
JDK_SIZE=""
JOKER_VERSION=""
JOKER_SIZE=""

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

if [ -n "$JOKER" ]; then
    JOKER_VERSION="joker $(joker --version 2>&1 | head -1)"
    JOKER_PATH="$(readlink -f "$JOKER" 2>/dev/null || readlink "$JOKER" 2>/dev/null || echo "$JOKER")"
    if [[ "$JOKER_PATH" == ../* ]]; then
        JOKER_PATH="$(cd "$(dirname "$JOKER")" && cd "$(dirname "$JOKER_PATH")" && pwd)/$(basename "$JOKER_PATH")"
    fi
    JOKER_SIZE="$(ls -lh "$JOKER_PATH" 2>/dev/null | awk '{print $5}')"
fi

CPU_INFO="$(sysctl -n machdep.cpu.brand_string 2>/dev/null || cat /proc/cpuinfo 2>/dev/null | grep 'model name' | head -1 | cut -d: -f2 | xargs || echo "unknown")"

echo ""
echo "=== Environment ==="
echo "System: $(uname -ms), $CPU_INFO"
echo "let-go: $LETGO_SIZE binary"
[ -n "$BB" ] && echo "babashka: $BB_VERSION ($BB_SIZE binary)"
[ -n "$JOKER" ] && echo "joker: $JOKER_VERSION ($JOKER_SIZE binary)"
[ -n "$CLJ" ] && echo "clojure: $CLJ_VERSION"
[ -n "$CLJ" ] && echo "JDK: $JDK_VERSION ($JDK_SIZE)"
echo ""

# --- Measure startup time ---

echo "=== Startup Time ==="
STARTUP_JSON="/tmp/bench_startup.json"
STARTUP_CMDS=("$LETGO -e nil")
STARTUP_NAMES=("-n let-go")
[ -n "$BB" ] && STARTUP_CMDS+=("bb -e nil") && STARTUP_NAMES+=("-n babashka")
[ -n "$JOKER" ] && STARTUP_CMDS+=("joker -e nil") && STARTUP_NAMES+=("-n joker")
[ -n "$CLJ" ] && STARTUP_CMDS+=("clj -M -e nil") && STARTUP_NAMES+=("-n clojure")

STARTUP_ARGS=(--warmup "$WARMUP" --runs "$RUNS" --export-json "$STARTUP_JSON")
for i in "${!STARTUP_CMDS[@]}"; do
    STARTUP_ARGS+=("${STARTUP_NAMES[$i]}" "${STARTUP_CMDS[$i]}")
done
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
    # Sort and take middle
    IFS=$'\n' sorted=($(sort -n <<<"${mems[*]}")); unset IFS
    local median_bytes=${sorted[1]}
    local median_mb=$(echo "scale=1; $median_bytes / 1048576" | bc)
    echo "  $name: ${median_mb}MB"
    echo "$median_mb"
}

echo "Running: nil (startup only)"
LETGO_STARTUP_MEM=$(measure_mem "$LETGO -e nil" "let-go" | tail -1)
BB_STARTUP_MEM=""
[ -n "$BB" ] && BB_STARTUP_MEM=$(measure_mem "bb -e nil" "babashka" | tail -1)
JOKER_STARTUP_MEM=""
[ -n "$JOKER" ] && JOKER_STARTUP_MEM=$(measure_mem "joker -e nil" "joker" | tail -1)
CLJ_STARTUP_MEM=""
[ -n "$CLJ" ] && CLJ_STARTUP_MEM=$(measure_mem "clj -M -e nil" "clojure" | tail -1)

echo ""
echo "Running: fib 35 (compute-heavy)"
LETGO_FIB_MEM=$(measure_mem "$LETGO benchmark/fib.clj" "let-go" | tail -1)
BB_FIB_MEM=""
[ -n "$BB" ] && BB_FIB_MEM=$(measure_mem "bb benchmark/fib.clj" "babashka" | tail -1)
JOKER_FIB_MEM=""
[ -n "$JOKER" ] && JOKER_FIB_MEM=$(measure_mem "joker benchmark/fib.clj" "joker" | tail -1)
CLJ_FIB_MEM=""
[ -n "$CLJ" ] && CLJ_FIB_MEM=$(measure_mem "clj -M -e '(load-file \"benchmark/fib.clj\")'" "clojure" | tail -1)

echo ""
echo "Running: reduce 1M (large collection)"
LETGO_REDUCE_MEM=$(measure_mem "$LETGO benchmark/reduce.clj" "let-go" | tail -1)
BB_REDUCE_MEM=""
[ -n "$BB" ] && BB_REDUCE_MEM=$(measure_mem "bb benchmark/reduce.clj" "babashka" | tail -1)
JOKER_REDUCE_MEM=""
[ -n "$JOKER" ] && JOKER_REDUCE_MEM=$(measure_mem "joker benchmark/reduce.clj" "joker" | tail -1)
CLJ_REDUCE_MEM=""
[ -n "$CLJ" ] && CLJ_REDUCE_MEM=$(measure_mem "clj -M -e '(load-file \"benchmark/reduce.clj\")'" "clojure" | tail -1)

# --- Run benchmarks ---

echo ""
echo "=== Performance Benchmarks ==="

# Store all benchmark JSONs
declare -A BENCH_JSONS

for bench in "${BENCHMARKS[@]}"; do
    name="$(basename "$bench" .clj)"
    echo ""
    echo "--- $name ---"

    JSON="/tmp/bench_${name}.json"
    CMDS=("-n let-go" "$LETGO $bench")
    [ -n "$BB" ] && CMDS+=("-n babashka" "bb $bench")
    if [ -n "$JOKER" ] && ! echo "$JOKER_SKIP" | grep -qw "$name"; then
        CMDS+=("-n joker" "joker $bench")
    fi
    [ -n "$CLJ" ] && CMDS+=("-n clojure" "clj -M -e '(load-file \"$bench\")'")

    hyperfine --warmup "$WARMUP" --runs "$RUNS" --export-json "$JSON" "${CMDS[@]}" 2>&1
    BENCH_JSONS[$name]="$JSON"
done

# --- Generate results.md ---

RESULTS_FILE="$SCRIPT_DIR/results.md"

# Determine which runtimes are available for column headers
RUNTIME_NAMES=("let-go" "babashka")
[ -n "$JOKER" ] && RUNTIME_NAMES+=("joker")
RUNTIME_NAMES+=("clojure JVM")
NUM_RUNTIMES=${#RUNTIME_NAMES[@]}

header_row="| |"
sep_row="|---|"
for rn in "${RUNTIME_NAMES[@]}"; do
    header_row+=" $rn |"
    sep_row+="---|"
done

cat > "$RESULTS_FILE" << EOF
## Benchmark Results

### Methodology

All benchmarks use [hyperfine](https://github.com/sharkdp/hyperfine) with $WARMUP warmup runs
and $RUNS timed runs per benchmark. Times shown are mean ± σ wall-clock time. Each benchmark
file is valid Clojure that runs unmodified on all runtimes. Peak memory is measured
via \`/usr/bin/time -l\` (median of 3 runs). Clojure JVM times include full JVM startup
(~350-500ms) which dominates short benchmarks. Joker is skipped for benchmarks that
would exceed reasonable time limits or use unsupported features (transducers).

**System:** $(uname -ms), $CPU_INFO

**Runtimes:**

$header_row
$sep_row
| **Version** | — | $BB_VERSION | ${JOKER_VERSION:-} | $CLJ_VERSION |
| **Platform** | Go bytecode VM | GraalVM native | Go tree-walk interpreter | JVM (HotSpot) |
| **Binary/runtime size** | **$LETGO_SIZE** | $BB_SIZE | ${JOKER_SIZE:-—} | $JDK_SIZE |

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

best = min(e[1] for e in entries)
print('| Runtime | Time |')
print('|---|---|')
for name, mean, stddev in entries:
    s = fmt(mean, stddev)
    if mean == best:
        print(f'| **{name}** | **{s}** |')
    else:
        print(f'| {name} | {s} |')
" >> "$RESULTS_FILE"

# Memory table
MEM_HEADER="| Workload | let-go | babashka |"
MEM_SEP="|---|---|---|"
[ -n "$JOKER" ] && MEM_HEADER+=" joker |" && MEM_SEP+="---|"
MEM_HEADER+=" clojure JVM |"
MEM_SEP+="---|"

HAS_JOKER_COL=""
[ -n "$JOKER" ] && HAS_JOKER_COL="1"

python3 -c "
has_joker = '$HAS_JOKER_COL' != ''

def bold_min_row(label, vals):
    # vals is list of (value_str, numeric_or_none)
    nums = [v[1] for v in vals if v[1] is not None]
    best = min(nums) if nums else None
    cells = []
    for s, n in vals:
        if n is not None and best is not None and n == best:
            cells.append(f'**{s}**')
        else:
            cells.append(s)
    return f'| {label} | ' + ' | '.join(cells) + ' |'

def parse_mb(s):
    try: return float(s)
    except: return None

rows_data = [
    ('startup (nil)', [
        ('${LETGO_STARTUP_MEM}MB', parse_mb('${LETGO_STARTUP_MEM}')),
        ('${BB_STARTUP_MEM:-—}MB', parse_mb('${BB_STARTUP_MEM}')),
    ] + ([('${JOKER_STARTUP_MEM:-—}MB', parse_mb('${JOKER_STARTUP_MEM}'))] if has_joker else []) + [
        ('${CLJ_STARTUP_MEM:-—}MB', parse_mb('${CLJ_STARTUP_MEM}')),
    ]),
    ('fib(35)', [
        ('${LETGO_FIB_MEM}MB', parse_mb('${LETGO_FIB_MEM}')),
        ('${BB_FIB_MEM:-—}MB', parse_mb('${BB_FIB_MEM}')),
    ] + ([('${JOKER_FIB_MEM:-—}MB', parse_mb('${JOKER_FIB_MEM}'))] if has_joker else []) + [
        ('${CLJ_FIB_MEM:-—}MB', parse_mb('${CLJ_FIB_MEM}')),
    ]),
    ('reduce 1M', [
        ('${LETGO_REDUCE_MEM}MB', parse_mb('${LETGO_REDUCE_MEM}')),
        ('${BB_REDUCE_MEM:-—}MB', parse_mb('${BB_REDUCE_MEM}')),
    ] + ([('${JOKER_REDUCE_MEM:-—}MB', parse_mb('${JOKER_REDUCE_MEM}'))] if has_joker else []) + [
        ('${CLJ_REDUCE_MEM:-—}MB', parse_mb('${CLJ_REDUCE_MEM}')),
    ]),
]

header = '| Workload | let-go | babashka |'
sep = '|---|---|---|'
if has_joker:
    header += ' joker |'
    sep += '---|'
header += ' clojure JVM |'
sep += '---|'

print()
print('### Peak Memory Usage (RSS)')
print()
print(header)
print(sep)
for label, vals in rows_data:
    print(bold_min_row(label, vals))
" >> "$RESULTS_FILE"

cat >> "$RESULTS_FILE" << EOF

### Performance

EOF

# Build performance table dynamically
{
PERF_HEADER="| Benchmark | let-go | babashka |"
PERF_SEP="|---|---|---|"
[ -n "$JOKER" ] && PERF_HEADER+=" joker |" && PERF_SEP+="---|"
PERF_HEADER+=" clojure JVM |"
PERF_SEP+="---|"
echo "$PERF_HEADER"
echo "$PERF_SEP"

for bench in "${BENCHMARKS[@]}"; do
    name="$(basename "$bench" .clj)"
    JSON="/tmp/bench_${name}.json"
    HAS_JOKER="false"
    if [ -n "$JOKER" ] && ! echo "$JOKER_SKIP" | grep -qw "$name"; then
        HAS_JOKER="true"
    fi
    python3 -c "
import json, sys
with open('$JSON') as f:
    d = json.load(f)

def fmt(mean, stddev):
    if mean < 0.1:
        return f'{mean*1000:.1f}ms ± {stddev*1000:.1f}ms'
    return f'{mean:.3f}s ± {stddev:.3f}s'

# Collect results with raw mean for comparison
results = {}  # key → (formatted, raw_mean)
for r in d['results']:
    cmd = r['command'].strip()
    if cmd == 'let-go': results['letgo'] = (fmt(r['mean'], r['stddev']), r['mean'])
    elif cmd == 'babashka': results['bb'] = (fmt(r['mean'], r['stddev']), r['mean'])
    elif cmd == 'joker': results['joker'] = (fmt(r['mean'], r['stddev']), r['mean'])
    elif cmd == 'clojure': results['clj'] = (fmt(r['mean'], r['stddev']), r['mean'])

# Find winner (lowest mean)
best = min(v[1] for v in results.values())

def cell(key):
    if key not in results:
        return '—'
    s, mean = results[key]
    if mean == best:
        return f'**{s}**'
    return s

has_joker_col = '$JOKER' != ''
row = f'| $name | {cell(\"letgo\")} | {cell(\"bb\")} |'
if has_joker_col:
    row += f' {cell(\"joker\")} |'
row += f' {cell(\"clj\")} |'
print(row)
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
rm -f "$SCRIPT_DIR/letgo"
