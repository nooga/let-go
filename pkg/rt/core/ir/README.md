# IR framework

An SSA IR with block-arg-style branches, sitting between the source
compiler and the bytecode VM. The IR's algorithms and data live in
Lisp under this directory; Go provides only a small substrate (op
table, chunk emission, value constructors) via bridge primitives in
the `ir/` namespace.

## Architecture at a glance

```
   forms
     │
     ▼
  ir.build/build-fn  ──►  Lisp-atom IR Function  ──►  ir.passes.pipeline/optimize-fn
                                                              │
                                                              ▼
                                                       optimized IR Function
                                                              │
                                                              ▼
                                                       ir.lower/lower
                                                              │
                                                              ▼
                                                      vm.CodeChunk
```

Each layer is one file (or one directory) in this tree:

| Path                                          | Role                                                  |
|-----------------------------------------------|-------------------------------------------------------|
| `data.lg`                                     | IR data layer: Function/Block/Inst atom shape, ctors, structural mutators, uses cache. Loaded from source (intern block exposes accessors under `ir/`). |
| `data_generated.lg`                           | Mechanical field accessors, generated from `examples/go-gen/ir_data.lg`. |
| `zipper.lg`                                   | Fn-scoped cursor over Insts in block-order.           |
| `passes.lg`                                   | `defpass` DSL; binds `*current-fn*` / `*current-inst*` / `*current-zip*`. |
| `build.lg`                                    | `ir.build/build-fn` — Lisp form → IR Function.        |
| `dump.lg`                                     | Text dump for debugging.                              |
| `dominance.lg`                                | Cooper/Harvey/Kennedy immediate dominators + `dominates?`. |
| `lower.lg`                                    | IR → bytecode via chunk-emission bridge primitives.   |
| `passes/dce.lg`                               | Dead-code elimination.                                |
| `passes/constfold.lg`                         | Constant folding + algebraic identities + commutative canonicalization. |
| `passes/cse.lg`                               | Per-block common-subexpression elimination.           |
| `passes/licm.lg`                              | Loop-invariant code motion.                           |
| `passes/thread-block-args.lg`                 | Pre-lower fixup for cross-block direct refs.          |
| `passes/typeinfer.lg`                         | Abstract-interpret types.                             |
| `passes/pipeline.lg`                          | Composes the passes; `compile-form` is the e2e entry. |

The Go substrate consumed via the `ir/` namespace bridge:

- `ir/op`, `ir/refs`, `ir/aux`, `ir/block-of`, `ir/type-of` — Inst field access.
- `ir/blocks`, `ir/block-params`, `ir/block-preds`, `ir/block-insts`, `ir/block-term` — Block access.
- `ir/fn-name`, `ir/fn-arity`, `ir/fn-variadic?`, `ir/fn-entry`, `ir/fn-consts` — Function field access.
- `ir/new-fn`, `ir/new-branch-target`, `ir/new-cond-target` — constructors.
- `ir/add-block`, `ir/add-inst`, `ir/add-terminator!`, `ir/add-block-param!`, `ir/add-pred!` — structural mutators.
- `ir/replace-op!`, `ir/set-refs!`, `ir/set-aux!`, `ir/set-type!`, `ir/replace-all-uses!`, `ir/clone-inst`, `ir/remove-inst!` — pass-time mutators.
- `ir/uses`, `ir/invalidate-uses!` — def→use index (cached, invalidated by mutators).
- `ir/op-terminator?`, `ir/op-pure?`, `ir/op-cheap-load?`, `ir/op-stack-out` — op-table queries (genuine Go bridges).
- `ir/chunk-emit`, `ir/chunk-emit-with-arg`, `ir/chunk-emit-placeholder`, `ir/chunk-emit-dup-nth`, `ir/chunk-emit-recur`, `ir/chunk-update!`, `ir/chunk-length`, `ir/chunk-max-stack`, `ir/chunk-set-max-stack!`, `ir/chunk-intern-const`, `ir/chunk-add-source-info!`, `ir/new-chunk`, `ir/new-consts`, `ir/new-source-info` — VM-substrate primitives for lower.

Everything in the first three groups is implemented in `data.lg` (mostly via the generated `data_generated.lg`); only the last two groups are Go bridge primitives.

## Using the IR

End-to-end compile of a form:

```clojure
(require '[ir.passes.pipeline :as p])

(p/compile-form '(defn add [a b] (+ a b)))
;; => boxed *vm.CodeChunk
```

Steps individually:

```clojure
(require '[ir.build] '[ir.passes.pipeline :as p] '[ir.lower :as l] '[ir.dump :as d])

(let [f (ir.build/build-fn '(defn add [a b] (+ a b)))]
  (println (d/dump f))           ; pre-optimization
  (p/optimize-fn f)
  (println (d/dump f))           ; post-optimization
  (l/lower f))                   ; returns boxed *vm.CodeChunk
```

The `f` value is a Lisp atom; `@f` is a map with `:name :arity :variadic? :entry :consts :insts :blocks` etc. Inspect via `ir/op`, `ir/refs`, etc. — never reach into `@f` directly from pass code; the accessors keep the abstraction stable.

## Extending the IR

### Adding a new optimization pass

Use the `defpass` DSL. The macro binds `*current-fn*` and `*current-inst*` for the duration of each visit; helpers `replace-refs!`, `replace-aux!`, `remove!` operate on the current Inst.

```clojure
(ns ir.passes.my-pass
  (:require [ir.passes :refer [defpass replace-refs!]]))

(defpass my-pass
  "One-sentence description."
  [inst]                            ; binds inst to the current InstId
  (let [f  ir.passes/*current-fn*
        op (ir/op inst f)]
    (when (= op :add)
      ;; do something
      )))
```

Then add it to the pipeline (`passes/pipeline.lg`):

```clojure
(:require ... [ir.passes.my-pass :refer [my-pass]])

(defn optimize-fn [f]
  ...
  (my-pass f)
  ...
  f)
```

Wire it into the bundle by adding the embed declaration to `pkg/rt/irpasses.go` and an entry in `cmd/lgbgen/main.go`.

If your pass needs block-level state (e.g., a hash table reset per block, as CSE does) the `defpass` cursor doesn't expose block transitions cleanly — walk blocks explicitly with `(doseq [bid (ir/blocks f)] ...)`. See `passes/cse.lg` for the pattern.

### Adding a new IR op

Two-step process — the op is a Go-side enum value because it's referenced from the chunk-emission bridge.

1. **Edit `examples/go-gen/ir_ops.lg`**. Add a row to the `ops` vector:
   ```clojure
   ["MyOp" 1 1 true false OP_MY_OP false "MyOp" "what it does"]
   ;; name, stk-in, stk-out, pure?, term?, bytecode, cheap?, display, comment
   ```
   - `stack-in` / `stack-out` declare the op's stack effect.
   - `pure?` controls CSE/DCE/fold eligibility.
   - `terminator?` marks the op as ending a basic block.
   - `bytecode` is the corresponding `vm.OP_*` constant (or `nil` if not directly lowered).
   - `cheap?` means `materialize` may re-emit it inline at use sites.

2. **Run `make generate-ir-ops`**. Regenerates `pkg/ir/op_generated.go` with the new entry. The op is now usable from Lisp via its kebab-case keyword (`:my-op`).

3. **Update build.lg and lower.lg as needed** — if forms can produce this op, add a build path; if Lower needs to handle it, add a case in `lower-node!`.

Order matters: the iota values in the op enum are positional. Inserting an op in the middle shifts later constants. Always append at the end of the relevant section.

### Adding a new IR data field

The IR data is a Lisp atom holding a map. Adding a field means updating both the spec and the hand-written companion:

1. **Edit `examples/go-gen/ir_data.lg`**. Add the field to the appropriate `defirdata` block (Function/Block/Inst/BranchTarget/CondTarget):
   ```clojure
   (defirdata Inst
     :fields [[:op :kw] [:refs :vec] ... [:my-field :any]]
     :defaults {...})
   ```

2. **Run `make generate-ir-data`**. Emits new accessor + setter in `data_generated.lg`.

3. **Update `data.lg`** if structural ops need to initialize the field — most fields will fall back to `nil`, but for vector-valued fields you typically want `[]` in the constructor and in the intern block.

4. **Regen the bundle** with `go generate ./pkg/rt/`.

### Adding a new bridge primitive

Bridge primitives cross the Go ↔ Lisp boundary and are reserved for VM-substrate operations (chunk emission, op-table queries, value constructors). Adding one to a spec:

1. **Edit `examples/go-gen/ir_bridge.lg`**. Add an `:extras` entry to the relevant `defgostruct` block:
   ```clojure
   {:name "my-prim" :lisp-name "ir/my-prim"
    :args [Self Int]
    :body "/* Go code; arg0 is the receiver, arg1Int is the int */
           return vm.NIL, nil"}
   ```
   - `:args` declares each arg's type for the unbox prelude. Supported types include `Self`, `Int`, `Bool`, `String`, `OpKw`, `BoxedChunk`, `BoxedBranchTarget`, `InstIdVec`, `Value`, `BoxedAuxOrNil`.
   - `:body` is raw Go pasted into the generated function. The prelude binds `arg0`, `arg1`, etc., by the declared types.

2. **Run `make generate-ir-bridge`**. Regenerates `pkg/rt/ir_bridge_generated.go`.

If you need a new arg type not in `unbox-arg`'s switch, add a case there.

## Loading-order quirks

The Lisp IR layer has a load-time dance that's invisible most of the time but worth knowing about if you're debugging "nil is not a function" errors during compile-form:

- `data.lg` is NOT in the precompiled bundle. The precompile only replays defn STUBS (nil-valued vars); the intern block at the bottom of `data.lg` needs LIVE function values to do anything useful. So `data.lg` is loaded from source via the resolver's switch case (`pkg/resolver/resolver.go`, case `"ir.data"`).

- `ir.build` declares `(:require ir.data)` so loading build also triggers data.lg's source load. This guarantees `ir/op`, `ir/refs`, etc., resolve to Lisp impls before build's defns capture Var pointers in their compiled bytecode.

- At lgbgen (precompile) time, `cmd/lgbgen/main.go` bootstraps `data.lg` source first so subsequent bundled namespaces can resolve `ir/*` symbols at compile time.

- 3-arg `intern` was tweaked to update an existing Var's root in place rather than creating a fresh Var. Critical because compiled bytecode captures Var pointers in its const pool; recreating the Var would leave those pointers stuck on the old nil-rooted Var.

If you see `ir/something: nil is not a function`, the most likely cause is that `data.lg` didn't finish loading before the caller compiled. Either add `ir.data` to a `:require` clause in the caller, or load it explicitly via `(require '[ir.data])` at the top of your file.

## Testing

The project's IR regression tests are in `pkg/ir/`:

- `lisp_constfold_test.go` — behavioral tests for each constfold strategy.
- `pipeline_bench_test.go` — end-to-end pipeline corpus + benchmarks.

When adding a pass or op, add a corresponding entry to one or both. The `ensureLoader` helper in `lisp_constfold_test.go` shows the canonical way to wire up the namespace loader in tests.
