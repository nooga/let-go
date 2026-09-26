---
status: experimental
last-verified: 2026-09-20
authoritative-for:
  - clojure-test-go-adapter
human-verified:
---

# Run Clojure tests with go test

The experimental `scripts/lg-test-go` adapter lowers `clojure.test` test bodies
to Go and exposes them as named subtests of `TestClojure`. Assertions, errors,
and fixture failures determine the Go test result; a test's return value does
not. An empty test is valid.

From the let-go checkout:

```sh
make build
bin/lg scripts/lg-test-go test/fixtures/go-test-adapter/framework.lg build/framework-tests
go test -v ./build/framework-tests -run '^TestClojure$/^passing$'
```

The example corpus also contains deliberate failures. Running it without
`-run` should fail. Replace its path with your own single-namespace `.lg` test
file. Use `LG_SOURCE_PATHS` when the source requires project namespaces.

The adapter writes `bodies.go`, `clojure_test.go`, and `tests.lgb`. Within the
checkout, the generated directory uses the existing Go module. Outside a
module, provide a `go.mod` requiring the matching let-go revision (or a local
`replace`) before running `go test`.

Each discovered `deftest` body must lower successfully. Generation fails on an
empty suite, duplicate test definitions, or a reserved generated-name collision.
The generated wrapper calls the lowered body directly and installs it in the
original var's `:test` metadata. Calling one test from another therefore also
uses the lowered body. Namespace initialization, helper functions, fixtures,
and the test framework may execute from embedded bytecode; this is not a claim
that the entire namespace has been lowered to Go.

`:once` fixtures wrap the selected suite; `:each` fixtures run inside each Go
subtest. Failed assertions and uncaught test errors appear in `go test -v` with
the Clojure source location, message, testing context, expected value, and actual
value. Each subtest also logs its assertion counts. An uncaught body error uses
the test declaration's location when no assertion location is available.

This first version accepts one namespace per generation and runs sequentially.
It replaces the dynamic `clojure.test/report` binding with a Go reporter; custom
reporters and tests that rebind it are outside this adapter's contract. Avoid
concurrent mutation of the runtime from other Go tests. Generation evaluates
top-level forms to load macros and definitions, as the existing AOT pipeline
does; it is intended for test-definition files rather than scripts that invoke
their own test runner at top level.

Run the adapter's integration checks with:

```sh
go test ./test/e2e -run '^TestClojureGoTestAdapter$' -count=1
```

These checks compile generated packages before checking expected failures.
They cover passing assertions, a failure before a successful return, assertion
and uncaught errors, empty and composed tests, Go filtering, and fixture
failures. A change to only the emitted Go body must also change the test result,
confirming execution of the lowered body independently of the embedded test.
