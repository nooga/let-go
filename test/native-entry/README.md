# `test/native-entry/` — native-entry AST gate corpus

Fixtures for `make native-entry-gate`
(`TestNativeEntryASTGate` in `test/e2e/native_entry_ast_gate_test.go`).

`make native-entry-gate` also runs `TestJankSuiteDirectABIGeneratedGo`
(`test/e2e/jank_direct_abi_test.go`) — see
[Running the gate: the `test/clojure-test-suite` submodule](#running-the-gate-the-testclojure-test-suite-submodule).

Each fixture is three committed files, discovered by globbing `*.lg`:

| file | role |
| --- | --- |
| `<name>.lg` | let-go source: an `-main` plus at least one lowerable `defn` |
| `<name>.expect` | exact stdout of the built binary, byte for byte |
| `<name>.goexpect.json` | the structural contract for the generated Go |

A fixture missing either sidecar is a hard failure, never a skip.

## What the gate claims

For each fixture: this source was lowered to Go by the production path
(`scripts/lg-compile --entry-frame`), the generated Go has the committed
structure, it was compiled into a real binary, and that binary produced
`<name>.expect`. All mutation sites are located through `go/ast`, never by
hardcoded temporaries or byte offsets.

Two AST-located mutants must be **killed** for every fixture:

- **oracle mutant** — a structurally perturbed copy of a committed fragment
  must be rejected by the matcher;
- **AST call-site mutant** — retargeting the entry path's call (direct callee,
  or the trampoline's var name) must break the build or the run.

## The inverted semantic-mutant contract

The third mutant wraps the returned expression of the fixture's lowered Go
function. The *same* mutation is applied to both kinds of fixture, and the
required outcomes are **opposite**. This is the point of the mutant, not a
loophole:

| `lowering` | sub-test | required outcome |
| --- | --- | --- |
| `direct` | `live_native_body_mutant_dies` | the mutant MUST change observed output |
| `vm` | `dead_native_body_mutant_survives` | the mutant MUST NOT change anything: still builds, exits 0, prints byte-exact `.expect` |

- For `direct` rows, a dying mutant is the evidence that the generated Go body
  is what the entry path actually runs. A survivor would mean the output comes
  from somewhere else (a VM trampoline, a cached var, a folded constant).
- For `vm` rows, a **surviving mutant is required evidence**, not a weakness.
  A VM-backed defn is still *emitted* as a Go function, but nothing calls it —
  the entry path goes through `ec.Invoke`. So mutating that body must be
  invisible. If a `vm` fixture's body mutant ever dies, the gate fails: the
  body would be live, and the fixture would no longer scope the direct-call
  claim.

## `goexpect.json` fields

- `namespace`, `generatedFile` — where the lowered package lands under the
  lg-compile out dir.
- `entryFn`, `lispFn`, `goFn` — the generated entry function and the fixture's
  own `defn` on both sides of lowering.
- `lowering` — `direct` (the entry path calls the lowered Go function directly
  and no `ec.Invoke` trampoline for it survives anywhere in the package) or
  `vm` (the fixture is *intentionally* still on the trampoline; the entry path
  must contain no direct call). The `vm` row is what scopes the direct-call
  claim — without it, "no trampoline" could be an accident of the corpus.
- `functionFragment`, `callFragment` — expected AST shapes, matched with
  `internal/gofragment` (exactly-one cardinality).
- `semanticWrap` — the semantic mutant. `__EXPR__` is replaced by the existing
  returned expression of `goFn`, so the mutant stays compilable (generated
  temporaries keep being used) and its outcome can only come from observed
  behaviour. Required for both `direct` and `vm` rows — see the inverted
  contract above; a `vm` row's wrap must stay compilable too, otherwise the
  "dead body" claim is untested rather than proven.

## Fixture roles

- `int_arith`, `closure_capture` — direct native lowering (unboxed ints, closure capture).
- `vm_backed` — intentionally VM-dispatched (variadic ⇒ not `:direct-callable?`); scopes
  the direct-call claim, and carries `dead_native_body_mutant_survives` as the positive
  evidence that its generated Go body is dead code.
- `toplevel_effect` — top-level side effect plus a VM-backed callee. Pins the frame's
  main-chunk replay: `lg -c` bundles carry an empty NS table, so without
  `rt.RunProgramMainChunk` the VM-backed callee resolves a Nil var and panics; with a
  double replay the banner would print twice and the byte-exact `.expect` fails.

### Regression pinned by `toplevel_effect` + `vm_backed`

The entry frame previously **never replayed the program main chunk**. Top-level
forms therefore never ran in the lowered binary: side effects were silently
dropped and any var they were supposed to define stayed Nil, so a VM-dispatched
callee resolved a Nil var and panicked at run time. `toplevel_effect` and
`vm_backed` are the regression fixtures for that defect — `toplevel_effect`
pins that the chunk is replayed exactly once (banner printed once, not zero and
not twice) and `vm_backed` pins that a VM-backed callee resolves and executes
through the trampoline afterwards.

## Adding or refreshing a fixture

```sh
go build -o lg .
LG_SOURCE_PATHS=pkg/rt/core ./lg scripts/lg-compile --entry-frame /tmp/out tmpmod test/native-entry/<name>.lg
```

Copy the emitted declaration and entry call expression out of
`/tmp/out/<generatedFile>` into the sidecar, then run `make native-entry-gate`.
The fragments are exact AST shapes, so a lowering change that renames
generated temporaries will fail this gate with an AST diff; that is intended —
update the sidecar only after reviewing the diff.

## Running the gate: the `test/clojure-test-suite` submodule

`make native-entry-gate` runs two tests:

- `TestNativeEntryASTGate` — this corpus; no external dependency;
- `TestJankSuiteDirectABIGeneratedGo` — lowers one **pinned** jank deftest
  (`clojure.core_test/identical?`) with the test-scoped generator
  `test/tools/jank-go-harness.lg` into a real Go test package and requires:
  an inline Go AST oracle match, that oracle's falsifier, a passing `go test`
  of the generated package, and the death of an AST-located semantic mutant.
  The generated package is written into `t.TempDir()` as its own Go module that
  `replace`s `github.com/nooga/let-go` with this checkout — the same shape the
  gate above uses — so generated code never lands in the repository tree.

That test reads its fixture from `test/clojure-test-suite`, which is a
**git submodule**. Two constraints follow:

1. `git submodule update --init test/clojure-test-suite` must have been run in
   the **primary git worktree**.
2. jj does not manage submodules (they go through git), and git worktrees / jj
   workspaces do not materialize them — so in every *secondary* workspace the
   directory exists but is empty.

Remedy for a jj workspace or secondary worktree — symlink it from the primary
worktree:

```sh
scripts/link-clojure-test-suite.sh <workspace-path> [primary-worktree]
```

The script is idempotent (re-running refreshes the link), verifies the
submodule really is initialized in the primary worktree, and refuses to clobber
a real non-empty directory. jj does **not** report the symlink as a
working-copy change, so it never pollutes a commit.

A missing suite is a **hard failure with both remedies in the message**, never a
skip: that test carries the direct-ABI claim, and a silent skip would turn
that claim into an unmeasured assumption.

Its semantic mutation site is located structurally through `go/ast`,
using the same `wrapLastReturnResult` helper as this gate (last `ReturnStmt` in
walk order, original expression preserved so the mutant stays compilable). A
gensym-numbered temporary such as `v78` never appears in test source.
