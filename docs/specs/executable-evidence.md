---
status: draft
last-verified: 2026-09-11
authoritative-for: [executable-evidence]
supersedes: []
human-verified:
---

# Executable spec evidence

A specification that carries its own evidence. Requirement sentences in a
spec are tagged with a marker; each marker is answered by a block of
runnable Clojure in the same document. `scripts/spec-evidence.lg` extracts
those blocks, turns them into portable runner programs, executes them, and
reports one line per case. The document is therefore checkable the way a
test suite is: it fails when the implementation drifts away from what the
prose claims.

This spec is itself the first document the tooling runs. Every fence below
is executed by `make spec-evidence` and by `TestSpecEvidence`, so the
document cannot describe behaviour the tooling does not have.

## Purpose

Prose rots because nothing re-reads it. The three properties this tooling
buys are:

- **Pairing.** A requirement without evidence, or evidence without a
  requirement, is a lint error — the two sides cannot drift apart silently.
- **Execution.** Evidence is source, not a transcript. It runs on `lg` on
  every `go test ./test/`.
- **Portability.** The same evidence runs on JVM Clojure, so a spec claim
  can be checked against the reference implementation rather than only
  against let-go's current behaviour.

## Grammar

A spec is ordinary markdown. Two constructs are special.

| Construct | Shape | Meaning |
|---|---|---|
| Requirement marker | `[R-<name>]` anywhere in prose | Declares a requirement named `<name>`. `<name>` is `[a-z0-9-]+`. |
| Evidence fence | a fenced block opened with <code>&#96;&#96;&#96;&lt;type&gt; @R-&lt;name&gt;</code> plus optional `key=value` attributes | The evidence answering `[R-<name>]`. |
| Evidence table | an HTML comment `<!-- evidence: @R-<name> -->` immediately followed by a pipe table | Table-shaped evidence answering `[R-<name>]`. |

Written out, with the opener indented here so this document does not
nest an evidence block inside its own prose:

    ```clj-repl @R-some-name oracle=none expect=fail
    (form)
    ;=> value
    ```

    <!-- evidence: @R-some-name -->
    | form | value | out |
    |---|---|---|
    | `(print "a")` | `nil` | `"a"` |

`<type>` selects the vocabulary (`clj-repl`, `clj-table`, `clj-test`). A
table comment may name its type as a second word; it defaults to
`clj-table`. Attributes are space-separated `key=value` pairs on the fence
opener only — the table comment carries none. A line that looks like an
evidence opener but does not match the strict grammar is a lint error
rather than silent prose.

Inside a `clj-repl` block, a case is one top-level form followed by the
expectation lines that belong to it:

| Expectation line | Asserts |
|---|---|
| `;=> <edn>` | the form's value equals `<edn>` |
| `;; out: <edn-string>` | everything the form printed, byte for byte |
| `;; throws: <Class> "<message>"` | the form threw an instance of `<Class>` with that `ex-message` |
| `;; expect: fail` | this case is expected to fail (see Expected failure) |
| `;; expect: pass` | this case is expected to pass, overriding a block-level `expect=fail` |

The value expectation is read as EDN and compared with `=`; it is not
evaluated, so map key order is irrelevant and a bare symbol compares as a
symbol. [R-value-compared-as-edn]

```clj-repl @R-value-compared-as-edn
(+ 1 2)
;=> 3
[1 2 3]
;=> [1 2 3]
{:a 1 :b 2}
;=> {:b 2 :a 1}
```

Because the expectation is data rather than code, an unquoted symbol or
list on the right of `;=>` denotes itself. [R-expectation-not-evaluated]

```clj-repl @R-expectation-not-evaluated
(symbol "a")
;=> a
(list 1 2)
;=> (1 2)
(str/upper-case "ab")
;=> "AB"
```

Captured output is compared byte-exact, so whitespace and newlines are
part of the claim. [R-out-byte-exact]

<!-- evidence: @R-out-byte-exact -->
| form | value | out |
|---|---|---|
| `(print "a b")` | `nil` | `"a b"` |
| `(println)` | `nil` | `"\n"` |
| `(do (println "x") 7)` | `7` | `"x\n"` |

A table row is positional: `form`, `value`, `out`, and an optional fourth
`expect` column. An empty cell asserts nothing, which is how a row claims
an output without pinning a value. The header and separator rows are not
validated.

A throw expectation names both the exception class and its message; both
must match, and a form that throws when no throw was expected is reported
as `threw`. [R-throws-class-and-message]

```clj-repl @R-throws-class-and-message
(throw (ex-info "boom" {:k 1}))
;; throws: ExceptionInfo "boom"
```

## Vocabulary file format

There is no per-type program. A vocabulary is a plain-text file at
`scripts/spec-vocabs/<type>.vocab` that maps statement shapes to Clojure
templates; `scripts/spec-evidence/vocab.lg` is the single generator that
reads one and emits runner source. A vocabulary has four directive kinds,
one per line, with `#` for comments and blank lines ignored:

| Directive | Form | Effect |
|---|---|---|
| `include` | `include <other-vocab>` | Splices that vocabulary's directives in first. |
| `boot` | `boot <form>` | Emitted once, in order, at the top of the runner. |
| `epilogue` | `epilogue <form>` | Emitted once, in order, at the bottom of the runner. |
| rule | `<pattern> => <template>` | Rewrites a matching statement line into template text. |

A pattern is a sequence of space-separated tokens. A literal token must
match verbatim; `{x}` captures one word; `{x*}` captures the raw remainder
of the line with internal spacing preserved, and is therefore only legal as
the last token. A template substitutes `{x}` and `{x*}` in one
left-to-right pass, so substituted text is never rescanned and a form
containing braces survives verbatim. The reserved pattern `form` is the
case template; `{id}` and `{expect}` are supplied by the generator rather
than captured from the line.

Three vocabularies ship:

- **`clj-repl`** — top-level forms with inline expectations. Its `boot`
  section defines `spec-case!`, which runs one form with output captured,
  compares each expectation it was given, and prints the runner-contract
  lines. Its rules are the five expectation lines tabulated above plus
  `form`. Its `epilogue` exits non-zero if any case failed.
- **`clj-table`** — `include clj-repl` plus a single `row => form` rule. A
  tangled table row is `form \t value \t out [\t expect]`; the generator
  rewrites it into the equivalent `clj-repl` statement lines and runs it
  through the `clj-repl` rules. A table case is a repl case.
- **`clj-test`** — `deftest` blocks run through `clojure.test`. The block is
  emitted verbatim after a `boot` that opens the runner's own namespace;
  the `epilogue` hops back with `in-ns`, runs the namespaces the block
  declared, and reports one line for the whole block. `expect=fail` is not
  a `clojure.test` concept and is rejected by lint on a `clj-test` block.

## Runner contract

Every runner writes one line per case to stdout and exits non-zero if any
case failed. A case id is `<name>#<n>`, where `<n>` is the 1-based position
of the form (or table row) inside its block.

| Line | Meaning |
|---|---|
| `ok <id>` | The case passed. |
| `ok <id> # XFAIL` | The case failed and was marked expected-to-fail. Counts as ok. |
| `FAIL <id> want <edn> got <edn> [<what>]` | The case failed. |

`<what>` is one of `value`, `out`, `throws`, `threw` or `xpass`. `threw`
reports an unexpected exception and renders `got` as `[<Class> "<message>"]`;
`xpass` is the fixed line `FAIL <id> want failure got pass [xpass]`. Both
`want` and `got` are `pr-str`-ed, so a FAIL line is machine-readable — this
is what `accept` consumes.

`run` relays those lines verbatim, then prints a summary line
`spec-evidence: <spec> <ok>/<total> cases ok` and exits 0 only when every
case is ok. A runner that dies before reporting is turned into a synthetic
`FAIL <name>#? engine exit <code>` line, so a crashed engine can never be
mistaken for an empty pass.

## Engines and oracle mode

The engine that executes a runner is chosen, in order, by `--engine <cmd>`,
the `CLJ_ENGINE` environment variable, or the default `./lg`. Because
runner source is portable Clojure, the same evidence can be executed by JVM
Clojure as an oracle:

```
CLJ_ENGINE="clojure -M" ./lg scripts/spec-evidence.lg run <spec> --oracle
```

`--oracle` additionally skips every block marked `oracle=none`. That
attribute means "this claim is about let-go specifically" — the evidence is
still executed by the default engine on every run, but it is never put to
the reference implementation, because the reference would legitimately
disagree. [R-oracle-none-skipped]

```clj-repl @R-oracle-none-skipped oracle=none
(str (type 1))
;=> "let-go.lang.Int"
```

On the JVM that form yields `"class java.lang.Long"`. Marking the block
rather than weakening the claim keeps the let-go-specific fact pinned
without making the oracle run red.

## Promotion

`accept` runs the spec and then rewrites the document so each failing
expectation matches what the engine actually produced. It edits the
`;=>`, `;; out:` and `;; throws:` line for the case, inserting one directly
after the form when the case had no such expectation before; for a table it
rewrites the positional cell, preserving the row's backtick convention. All
edits are computed against the pre-edit file and applied bottom-up, so line
numbers stay valid. Afterwards `last-verified:` in the masthead is bumped
to today, best-effort.

Promotion is a recording step, never a judgement. Use it only once the
observed behaviour is the behaviour you want. When an oracle is available,
the disagreement matrix decides that:

| `./lg` | `clojure -M` oracle | Reading | Action |
|---|---|---|---|
| ok | ok | Agreement; the spec is honest. | Nothing. |
| FAIL | ok | let-go diverges from the reference. | Fix let-go, or record the divergence in `known-divergences.md`. Do not `accept`. |
| ok | FAIL | The spec has pinned a let-go-only behaviour. | Mark the block `oracle=none` and say why in prose, or weaken the claim. |
| FAIL | FAIL | The written expectation is simply wrong. | `accept` — both engines agree on the observed value. |
| run | skipped | `oracle=none`: an intentional let-go-only claim. | Nothing; the claim is still executed by `./lg`. |

## Expected failure

A known gap is more useful pinned than deleted. A case can be marked
expected-to-fail three ways:

| Mark | Scope |
|---|---|
| `expect=fail` fence attribute | every case in the block |
| `;; expect: fail` expectation line | that one case |
| `fail` in a table's fourth `expect` column | that one row |

The verdicts follow: a marked case that fails reports `ok <id> # XFAIL` and
counts as ok, and a marked case that *passes* reports
`FAIL <id> want failure got pass [xpass]` — so a gap that quietly closes
fails the gate rather than going unnoticed. `accept` promotes an xpass by
deleting the case's `;; expect: fail` line (or clearing the table's
`expect` cell); a block-level `expect=fail` cannot be promoted
automatically and is reported as blocked, because dropping the attribute
would change every case in the block. [R-expected-failure-marked]

```clj-repl @R-expected-failure-marked oracle=none
(clojure.string/split "a\n" #"\n")
;=> ["a"]
;; expect: fail
```

That is a real gap, not a contrived one: on `main`, let-go's
`clojure.string/split` keeps the trailing empty field and returns a list,
so the form yields `("a" "")`. JVM Clojure returns `["a"]`, which is why
the block is `oracle=none` — under `--oracle` the case would pass and be
reported as an xpass. When the fix (#843) lands, the `./lg` run turns the
case into an xpass, the gate goes red, and `accept` drops the
`;; expect: fail` line; the `oracle=none` attribute can be removed in the
same edit.

## Tooling and wiring

| Path | Role |
|---|---|
| `scripts/spec-evidence.lg` | CLI: `lint`, `tangle`, `generate`, `run`, `accept`. |
| `scripts/spec-evidence/parse.lg` | Markdown scan: markers, blocks, pairing errors. |
| `scripts/spec-evidence/vocab.lg` | Vocabulary parsing, templating, runner generation. |
| `scripts/spec-vocabs/*.vocab` | The shipped vocabularies. |
| `test/spec_evidence_*_test.lg` | Unit tests for the tooling itself. |
| `test/spec_evidence_test.go` | `TestSpecEvidence`: the Go gate. |

The CLI needs `LG_SOURCE_PATHS=scripts` so its helper namespaces resolve:

```
LG_SOURCE_PATHS=scripts ./lg scripts/spec-evidence.lg run docs/specs/executable-evidence.md
```

`tangle` writes each block to `<outdir>/<base>.<name>.<type>` plus a
`<base>.map.edn` index; `generate` additionally writes the `.cljc` runner
for each block. `run` does both into a temporary directory and executes.

Three make targets wrap the CLI over every spec that carries evidence:
`make spec-lint`, `make spec-evidence`, and `make spec-oracle` (the last
under `CLJ_ENGINE="clojure -M" ... --oracle`). `spec-lint` also runs as a
pre-commit hook on `docs/specs/*.md`.

`TestSpecEvidence` runs one Go subtest per `docs/specs/*.md` whose text
contains `@R-`, against a freshly built `lg`, and fails on a non-zero exit
with the runner output attached. A spec opts out with `evidence: skip` in
its masthead — the escape hatch for a spec whose engine dependencies have
not landed yet.

## Coverage rule

Markers and evidence must pair both ways. `lint` reports `marker R-<name>
has no evidence` and `evidence @R-<name> has no marker`, and exits
non-zero, so neither side of a requirement can be added, renamed or deleted
alone. Lint additionally rejects a malformed evidence opener, two forms on
one line of a `clj-repl` block, a `;;` or `;=>` line that matches no rule in
the block's vocabulary, and `expect=fail` on a `clj-test` block — each of
which would otherwise assert nothing while looking like it asserted
something.

A spec is not required to tag every sentence. It is required that every
sentence it *does* tag is answered, and that every answer is claimed.
