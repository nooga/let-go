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
