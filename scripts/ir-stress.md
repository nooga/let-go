# ir-stress.lg — IR pipeline stress harness

`scripts/ir-stress.lg` is a corpus-agnostic harness that drives a list of
`.lg` source files through one of the IR compilation paths and reports
per-file / per-defn pass rates, error buckets, and timings. Use it to
measure IR coverage, find regressions, and drill into individual slow
defns.

## Modes

```
go run . scripts/ir-stress.lg <mode> <args>...
```

| mode | args | what it does |
|---|---|---|
| `ir-compile` | `<dir> <file>...` | `(binding [*ir-compile* true] (eval form))` per defn. Measures what users hit at load time. |
| `lower-go`   | `<dir> <file>...` | Build IR + lower to Go in `:bridge` mode, no eval. Safe against the live pipeline; this is the AOT pass-rate signal. |
| `trace`      | `<dir> <file> <defn-name>` | Drill into ONE defn. Prints per-stage ms (expand / build / optimize / lower) and the per-pass timing table from `ir.passes.trace`. |

`<dir>` is the corpus root (absolute or relative to cwd). `<file>` paths
are relative to that root.

## Environment variables

| var | values | effect |
|---|---|---|
| `LG_STRESS_TIMEOUT_MS` | integer ms (default `5000`) | per-defn timeout. Defns exceeding it are killed and bucketed `:stress/timeout`. |
| `LG_STRESS_PASSES` | `1` or `2` (default `2`) | 1 = single pass; 2 = cold + warm (lets you measure cache warm-up effects). |
| `LG_STRESS_LOG` | path | append one TSV row per defn: `file<TAB>defn<TAB>bucket<TAB>ms`. Buckets include `:ok` and the `classify-error` keys (e.g. `:missing-form/set!`, `:gogen/nth-non-int`, `:stress/timeout`). A single `grep` over the log answers "what's the gap?". |
| `LG_STRESS_AUTOTRACE` | `1` | on each `:stress/timeout`, re-run that defn under `trace-one-defn` with a 4× timeout. Embeds a one-line "expand=A build=B optimize=C lower=D ms" summary inline, and feeds per-pass cost into the top-K aggregator. Doubles cost for slow defns. Currently effective only in `lower-go` mode. |
| `LG_STRESS_TOPK` | `1` | at run end, print: (a) cumulative ms per pass across all auto-traced defns; (b) one row per traced defn naming its single worst pass. Requires `LG_STRESS_AUTOTRACE=1` to populate. |

## Common recipes

```sh
# IR-pipeline self-AOT pass rate (the canonical lower-go signal)
LG_STRESS_PASSES=1 LG_STRESS_LOG=/tmp/ir-stress.log \
  go run . scripts/ir-stress.lg lower-go pkg/rt/core \
    core.lg walk.lg string.lg set.lg pprint.lg edn.lg io.lg async.lg \
    test.lg check.lg data.lg graph.lg zip.lg \
    ir/zipper.lg ir/build.lg ir/lower.lg ir/lower_go.lg ir/passes.lg \
    ir/dominance.lg ir/dump.lg ir/validate.lg ir/data.lg \
    ir/passes/constfold.lg ir/passes/cse.lg ir/passes/dce.lg \
    ir/passes/typeinfer.lg ir/passes/licm.lg ir/passes/pipeline.lg \
    ir/passes/mutability.lg ir/passes/trace.lg

# Drill into ONE slow defn
go run . scripts/ir-stress.lg trace pkg/rt/core core.lg for-emit

# Auto-trace every timeout + top-K pass cost summary
LG_STRESS_AUTOTRACE=1 LG_STRESS_TOPK=1 LG_STRESS_PASSES=1 \
  go run . scripts/ir-stress.lg lower-go pkg/rt/core/ir \
    data.lg build.lg lower.lg ...
```

## Bucket reference

`classify-error` (in the script) maps every non-`:ok` outcome to a
keyword so the failure tally is structured rather than free-form. Most
useful buckets:

| bucket | meaning |
|---|---|
| `:ok` | defn fully lowered |
| `:stress/timeout` | exceeded `LG_STRESS_TIMEOUT_MS` |
| `:missing-form/X` | special-form X (e.g. `set!`, `var`, `def`, `try`) not handled by build |
| `:unresolved/<sym>` | namespace not loaded for that sym (often a harness `require-deps!` gap) |
| `:validate/cross-block-ref` | build emitted IR with cross-block direct ref — invariant violation |
| `:validate/no-term` | build left a block without a terminator |
| `:gogen/nth-non-int` | a gogen helper hit `nth` with a non-int (typically a nil `(ir/op nil f)`) |
| `:gogen/bad-identifier` | gogen Go-name munging rejected a Lisp identifier |
| `:lower-go/nested-closure` | by-design lower-go fallback for closure-in-closure body shapes |
| `:build/unrecognized-form` | build saw a form shape it doesn't know how to handle |
| `:destructure-rejection` | destructuring pattern unsupported |

## When to reach for which mode

- **lower-go** — your default. Cheap, idempotent, fast feedback on coverage.
- **ir-compile** — when you're testing user-visible behavior, e.g. xsofy's
  load-time IR pipeline. Slower because it actually evals each defn.
- **trace** — when one defn shows `:stress/timeout` and you need to know
  which pass is at fault before optimising blindly. Stage-level timings
  also catch surprises (e.g. expand exploding on a deeply-nested macro).

## Source

`scripts/ir-stress.lg` (the harness)
`pkg/rt/core/ir/passes/trace.lg` (the per-pass instrumentation it builds on)
