# test/benches

Standalone benchmarks and bug-repro fixtures. These are **not** run by the
Go test harness (`TestRunner` skips this dir); run them by hand with `./lg`.

## pegbench — PEG-combinator microbenchmark

Mirrors the hot pattern of a parser-combinator grammar (yamlstar's parser
shape): closures built by combinators, invoked indirectly through vars,
threading `(string, pos)` state. 200K parses of a 32-char input ≈ 16M
closure invocations. Portable across let-go / glojure / JVM Clojure (no
interop, no metadata).

| file | what | run |
|------|------|-----|
| `pegbench.lg`       | plain VM bytecode, prints result            | `./lg test/benches/pegbench.lg` |
| `pegbench-timed.lg` | VM bytecode, prints run-ms                  | `./lg test/benches/pegbench-timed.lg` |
| `pegbench-ir.lg`    | `*ir-compile*` on (IR→optimize→bytecode)    | `./lg test/benches/pegbench-ir.lg` |
| `pegbench.clj`      | reference impl for JVM Clojure / glojure    | `clojure -M …` / `glj …` |

### Recorded baselines (200K parses, startup-corrected, 2026-06-12)

| engine                                  | run time | vs VM |
|-----------------------------------------|---------:|------:|
| glojure `glj` interpreter (0.6.5-rc30)  | ~60.6s   | 15.5× slower |
| let-go VM bytecode                      | 3.95s    | 1× |
| let-go IR-on bytecode (`*ir-compile*`)  | **2.61s**| **0.66× (34% faster)** |
| JVM Clojure                             | ~0.35s   | 11× faster |

The IR-on number depends on the loop/recur lowering fix in this branch
(see the parent commit); before it, `pegbench-ir.lg` crashed. The remaining
VM→JVM gap is the lowered-Go target (tracked separately).

## ysbench — YamlStar parse benchmark (ITER-0035 gate)

The real-corpus counterpart to pegbench: times `(yamlstar.parser/parse
"foo: 42")` against yamlstar's actual grammar (212 `def`-closure rules).
Needs a [yamlstar](https://github.com/yaml/yamlstar) checkout; skips
cleanly without one.

| file | what | run |
|------|------|-----|
| `ysbench.lg` | one mode (`interp` or `ir`), prints ms-per-call | `LG_SOURCE_PATHS=$HOME/development/yamlstar/core/src ./lg test/benches/ysbench.lg ir` |
| `ysbench.sh` | EPIC-013 acceptance gate: runs both modes, requires `ir` ≥1.10× faster (`LG_YSBENCH_MIN_SPEEDUP`) | `./test/benches/ysbench.sh` |

### Recorded numbers (2026-07-13, Apple Silicon)

| mode | ms/call |
|------|--------:|
| interpreted VM bytecode | 36.4 |
| `*ir-compile*` + `ir.passes.inline` enabled | 34.1 |

Speedup 1.07× — the gate currently **fails**. The inline/specialize passes
(EPIC-013 ITER-0031..0034) don't reach this corpus yet: the runtime inline
registry seeds only from same-namespace `defn` siblings, while yamlstar's
grammar rules are top-level `def` closures calling combinators from another
namespace (`prelude`). The parse stays dispatch-bound in the combinator
tree walk. Issue #352 stays open on making this gate pass.

## repro/ — loop-lowering miscompile fixtures

Minimal reproductions of the four `ir.lower` / one `ir.build` defects fixed
in the parent commit. Each must be run as a **file** (not `lg -e '(do …)'`,
which macroexpands defns before `set! *ir-compile*` runs and masks the bug).

| file | shape | expected |
|------|-------|----------|
| `repro/min1.lg` | loop seeded from a closure param, calls captured fn | `simple-p: 3` |
| `repro/min2.lg` | same, loop var not shadowing                       | `no-shadow: 3` |
| `repro/min3.lg` | loop with no captured call                          | `no-captured-call: 3` |
| `repro/min4.lg` | direct loop / closure+loop / closure-no-loop        | `direct-loop: 3` … |

Executing regression coverage lives in `test/ir_compile_loop_test.lg`.
