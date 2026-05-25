# gogen determinism

This document defines what determinism `gogen` promises, what it doesn't,
where nondeterminism can sneak in, and how to detect it. It is the
contract `gogen` commits to with code that integrates its output into a
reproducible build.

## Scope

`gogen` produces Go source files from Clojure specifications. Those
files may be:

1. Committed to the repository (checked-in alongside the spec)
2. Generated at build time (via `go generate` or a `Makefile` target)
3. Loaded as Go plugins at runtime (Mode 1/2 of the AOT design)

Cases 1 and 2 both depend on **same-spec-in produces same-bytes-out**.
Case 3 additionally depends on **same-bytes-out produces same-plugin-binary**,
which is the Go toolchain's responsibility — not `gogen`'s.

This document covers cases 1 and 2.

## The determinism contract

**Given:**

- the same `gogen.lg` library code,
- the same Clojure spec data,
- the same let-go binary,
- the same Go toolchain version,

**`gogen` guarantees** that `render-go` returns byte-identical output across
runs, processes, machines, and architectures.

Three of those four conditions are out of `gogen`'s scope (you have to
pin them yourself). The first — `gogen` itself behaving deterministically
given fixed inputs — is what this document covers.

## Sources of nondeterminism, ranked

### Source 1: Go map iteration order (HIGH)

Go's spec-mandated randomized map iteration is the primary hazard. The
let-go codebase already has known sites (see the audit in the project's
research notes); `gogen` must not introduce new ones.

**In the renderer (`pkg/rt/gogen.go`):**

- `astFromMap` dispatches via a `switch` on the `:tag` field — no
  map iteration.
- `format.Node` walks `go/ast` nodes in source order — no map iteration.
- `parser.ParseExpr` parses types from a string — no map iteration.

**Status:** clean. No Go map iteration occurs anywhere in the rendering
path. Add a regression test (`gogen_determinism_test.go`) that renders
the same AST 100× and asserts byte-equal output.

**In the macro layer (`examples/go-gen/gogen.lg`):**

The walker iterates user-supplied collections positionally. If the user
passes a vector or list, order is preserved. If the user passes a
`hash-map`, iteration order is randomized and emitted code differs
per run.

**Mitigation:** see "User-spec hygiene" below.

### Source 2: User-spec hygiene (HIGH)

The user controls what data drives the codegen. If they write:

```clojure
;; HAZARD: hash-map iteration is randomized.
(def schema {:User [:id :name] :Post [:id :title]})
(doseq [[type fields] schema]
  (emit-struct type fields))
```

…the generated file's struct order changes per run. Same for sets:

```clojure
;; HAZARD: set iteration order in let-go is randomized.
(def ops #{:add :sub :mul})
(doseq [op ops] (emit-handler op))
```

**Mitigations, in priority order:**

1. **Convention:** All spec data is a vector or a `sorted-map`. Document
   this prominently in `gogen`'s README. Make the bad pattern explicit
   in examples (with a comment) so users know what to avoid.

2. **`gogen/check-spec`:** A function that walks a value and throws on
   any `hash-map` or `set` found, recursively. Users can call this on
   their spec at the top of any codegen script:

   ```clojure
   (g/check-spec schema)  ; throws if schema contains an unsorted collection
   ```

3. **`gogen/sort-spec`:** A defensive coercion that converts any
   `hash-map`s in a spec to `sorted-map`s, recursively. Slower than
   `check-spec` and silently hides bugs (because the *user's* code
   probably still treats the spec as if order doesn't matter), so
   `check-spec` is preferred.

### Source 3: Go toolchain version drift (MEDIUM)

`go/format.Node`'s output is deterministic *within a Go release* but
can change between releases. Examples:

- Go 1.17 changed struct tag rendering in some edge cases.
- Go 1.18 added generics syntax; types with type parameters format
  differently.
- Go 1.19 tweaked some `embed`-adjacent formatting.

These changes are tiny but real. If your build uses Go 1.20 and CI
uses Go 1.21, you can get diff noise even with identical inputs.

**Mitigations:**

1. **Pin Go version** in `go.mod` (`go 1.x`) and let-go's CI config.
2. **Pin toolchain** with `toolchain go1.x.y` in `go.mod` (Go 1.21+).
   This makes Go itself enforce that the matching toolchain is used.
3. Use `go/format.Source` instead of `format.Node` where possible —
   it's the same code path as `gofmt`, which has stronger stability
   guarantees as a user-facing tool.

`gogen` doesn't promise stability across Go versions. That's the
caller's responsibility, just as it is for any user of `go/format`.

### Source 4: `parser.ParseExpr` for types (LOW)

`gogen` parses type specs from strings (`"[]float64"`, `"map[K]V"`).
`parser.ParseExpr` is deterministic. The risk would be if the *spec
strings themselves* are constructed nondeterministically — e.g.,
joining items from an unordered set. This reduces to Source 2.

### Source 5: Timestamps, hostnames, PIDs (LOW)

`gogen` MUST NOT embed any of:

- `time.Now()` results
- `os.Hostname()`
- `os.Getpid()`
- `runtime.GOMAXPROCS()` results
- Random UUIDs unless explicitly given by the user
- Any other ambient process state

A useful convention many codegen tools follow is including a comment
like `// Code generated by X; DO NOT EDIT.` (the standard Go marker)
but **never** including a timestamp in that comment. The fact that the
code was generated is stable; *when* it was generated is not, and
embedding it produces diff churn.

### Source 6: Floating-point literal formatting (LOW)

`strconv.FormatFloat(v, 'g', -1, 64)` is deterministic, but produces
shortest-roundtrip output. Two semantically equivalent floats from
different upstream sources (e.g., `0.1 + 0.2` vs literal `0.3`) format
differently. This is a *user* concern: don't let floating-point
arithmetic into your spec values. Use literals or sorted parsed
constants.

## Guarantees `gogen` provides

1. **Renderer is pure:** `render-go(ast)` is a referentially transparent
   function of its input AST map. No global state, no I/O, no ambient
   reads.

2. **AST construction is pure:** `walk-decl-form`, `walk-stmt-form`,
   etc. depend only on their input form. The macro expansion is a
   compile-time pure function of the user's source.

3. **No accidental map iteration:** The renderer dispatch is `switch`-based.
   The walker iterates user collections positionally. We don't iterate
   `map[K]V` anywhere internally.

4. **No timestamps, no PIDs, no hostnames** in generated output.

5. **Floating-point formatting is `strconv.FormatFloat(v, 'g', -1, 64)`**
   for `:float-lit`s — shortest-roundtrip, deterministic per value.

6. **Integer formatting is `strconv.FormatInt(v, 10)`** for `:int-lit`s
   — base 10, deterministic.

## Guarantees `gogen` does NOT provide

1. **Cross-Go-version stability.** Same input on Go 1.20 and Go 1.21 may
   produce different output if `format.Node` changed.

2. **Cross-let-go-version stability of the macro expansion.** If we
   change `walk-expr-form` to coerce literals differently, output
   changes. We will note such changes in the let-go CHANGELOG.

3. **Protection against unhygienic user specs.** If you iterate a
   `hash-map` in your spec, you lose. We provide `check-spec` to detect
   this, but you have to call it.

## Verification

Three layered checks, in increasing strength:

### Test 1: rendering is deterministic per AST

In `pkg/rt/gogen_test.go`:

```go
func TestRenderIsDeterministic(t *testing.T) {
    ast := someComplexAST(t)
    first := render(t, ast)
    for i := 0; i < 100; i++ {
        if render(t, ast) != first {
            t.Fatalf("render %d differs from first", i)
        }
    }
}
```

This catches Source 1 in the renderer.

### Test 2: macro expansion is deterministic per spec

In `test/gogen_macro_test.lg`:

```clojure
(deftest macro-expansion-deterministic
  (let [spec '[[OpAdd "+" foo bar] [OpSub "-" baz qux]]
        first (g/render-go (handler-for (first spec)))]
    (dotimes [_ 100]
      (is (= first (g/render-go (handler-for (first spec))))))))
```

This catches Source 1 in the macro layer.

### Test 3: end-to-end output hash

For any spec that emits files into the source tree, the build pipeline
compares the output's hash against a committed `.goldenhash` file:

```bash
# Makefile
.PHONY: check-gen
check-gen:
	@./scripts/regenerate-numeric-ops.lg | sha256sum -c numeric_ops.goldenhash
```

This catches every nondeterminism source — including Source 3 (Go
version drift) — and forces an explicit committed update when intended
output changes.

## `gogen/check-spec` — the recommended discipline

For any codegen pipeline that affects committed code or the build:

```clojure
(require '[gogen :as g])

(def spec
  ;; A vector of vectors. Stable. Auditable in git diff.
  '[[OpAdd "+" checkedAddInt NumAdd]
    [OpSub "-" checkedSubInt NumSub]])

(g/check-spec spec)   ; throws if spec contains hash-map/set

(doseq [op-row spec]
  (let [src (g/render-go (handler-for op-row))]
    (spit (str "gen/" (first op-row) ".go") src)))
```

`check-spec` is opt-in but cheap. The cost of not calling it is a flaky
build that mysteriously produces different bytes on every regeneration
— exactly the failure mode that ate JVM reproducible-build efforts for
years.

## A note on `(ns ...)` and let-go's existing hazards

This document is about `gogen`'s determinism. It is layered on top of
let-go, which has its own determinism hazards (see the audit notes —
default bytecode encoder ordering, `Map` type seq order, namespace
refer collision resolution). Those are out of `gogen`'s scope, but
they do mean that:

- The bytecode produced from a `gogen`-generated file may not be
  bit-stable across compiles even if the source is bit-stable.
- The `core_compiled.lgb` regeneration path uses `EncodeBundleOrdered`
  to work around this.
- If `gogen` is ever used to generate code that flows back through
  the same path, the same workaround applies.

Until those let-go-level hazards are fixed, the strongest reproducibility
guarantee you can make for `gogen`-generated code is "the source file
is byte-stable" — not "the resulting `.lgb` bytecode is byte-stable."

## Implementation plan

The pieces of this document that don't exist yet:

1. **`gogen/check-spec`** — recursive walker that throws on `hash-map`
   or `set`. ~40 lines of Clojure.

2. **`gogen_test.go` with `TestRenderIsDeterministic`** — render-100x
   determinism test. ~30 lines of Go.

3. **Macro-expansion determinism test in `test/`** — same shape at the
   Clojure layer. ~20 lines.

4. **`gogen/sort-spec`** (optional, if `check-spec` proves too strict)
   — defensive coercion. ~30 lines.

5. **README section** documenting the contract and pointing at this
   document.

None of these block the v1 of `gogen`. They block any use of `gogen`
that affects build reproducibility, which is the only use that matters
in practice.
