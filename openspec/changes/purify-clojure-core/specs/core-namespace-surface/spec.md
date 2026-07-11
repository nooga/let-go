# Core Namespace Surface — a Clojure-faithful `clojure.core` + `let-go.core` for lg extras

## ADDED Requirements

### Requirement: `clojure.core` carries no name outside the two-reference surface
lg's `clojure.core` public var set MUST NOT contain any name that is absent from
BOTH JVM Clojure and babashka across the standard namespaces (`clojure ∪ bb`).
The portable target — `clojure ∩ bb` — is a floor lg SHOULD cover, not a ceiling:
a JVM-`clojure.core` name lg provides that babashka lacks MAY remain. Names in
neither reference (lg-isms), and names whose home is another standard namespace,
MUST NOT be public members of `clojure.core`.

#### Scenario: no lg-specific names in clojure.core
- **WHEN** the public vars of `clojure.core` are enumerated (`ns-publics`)
- **THEN** every name is a public of JVM Clojure or babashka in some standard ns
- **AND** `gt`, `lt`, `ge`, `le`, `now`, `spy`, `parse-int` are NOT among them

#### Scenario: a stray lg-ism Def is caught
- **WHEN** a symbol present in neither JVM Clojure nor babashka is `Def`'d into
  `clojure.core`
- **THEN** the surface test fails

### Requirement: `clojure.core` leaks no JVM abstraction
`clojure.core` MUST NOT expose a var whose meaning depends on a JVM host
abstraction lg does not have (Java classes, proxies, primitive arrays). Such
interop MAY be provided by lg — the change does not prohibit it — but only from a
dedicated host-interop namespace, never `clojure.core`.

#### Scenario: host-interop names are absent from core
- **WHEN** `clojure.core` publics are enumerated
- **THEN** `proxy`, `gen-class`, `bean`, `aset-byte`, `java.lang.Long`,
  `clojure.lang.Atom` are absent from `clojure.core`

#### Scenario: interop is permitted outside core
- **WHEN** lg implements a host-interop capability
- **THEN** it is reachable from a dedicated interop namespace and is NOT a public
  of `clojure.core`

### Requirement: Standard-library names live only in their standard namespace
lg MUST expose names whose canonical Clojure home is `clojure.string`,
`clojure.set`, or `clojure.java.io` only from that namespace, not duplicated
into `clojure.core`.

#### Scenario: string helpers are not in core
- **WHEN** `clojure.core` publics are enumerated
- **THEN** `starts-with?`, `ends-with?`, `includes?`, `lower-case`,
  `upper-case`, `trim`, `split`, `index-of` are absent from `clojure.core`
- **AND** each still resolves via `clojure.string`

#### Scenario: set helpers are not in core
- **WHEN** `clojure.core` publics are enumerated
- **THEN** `difference` and `intersection` are absent from `clojure.core`
- **AND** each still resolves via `clojure.set`

### Requirement: `let-go.core` is the home for lg-specific extras
lg-specific, user-facing names that have no standard Clojure home MUST live in a
`let-go.core` namespace. `let-go.core` MUST be auto-refer'd into user namespaces
so existing unqualified usage keeps resolving without source changes.

#### Scenario: lg extras resolve unqualified via let-go.core
- **WHEN** user code in a default `(ns foo)` references `gt`, `now`, or `spy`
  unqualified
- **THEN** they resolve to the `let-go.core` vars
- **AND** `(resolve 'gt)` returns `#'let-go.core/gt`

#### Scenario: lg extras are qualified-reachable
- **WHEN** user code references `let-go.core/gt`
- **THEN** it resolves to the lg comparison helper

### Requirement: An explicit refer shadows the lg-extras baseline
An explicit refer of a lib re-exporting an lg-extra name MUST shadow the
auto-refer'd `let-go.core` (and `clojure.core`) mapping in the requiring
namespace, matching JVM Clojure refer precedence — for both `:refer :all` and
`:refer [sym]`. This holds because `let-go.core` is auto-refer'd as a baseline.

#### Scenario: :refer :all shadows an lg extra
- **WHEN** `lib` defines `gt` and a namespace does `(:require [lib :refer :all])`
- **THEN** unqualified `gt` in that namespace resolves to `lib/gt`, not
  `let-go.core/gt`

#### Scenario: unshadowed extras still resolve to let-go.core
- **WHEN** the same namespace references an lg extra `lib` does not define
- **THEN** it still resolves to the `let-go.core` var

### Requirement: Internal and interop machinery is not public in `clojure.core`
Compiler/runtime internals and host-interop registrations MUST NOT be exposed as
public members of `clojure.core`.

#### Scenario: internal dyn vars are not core publics
- **WHEN** `clojure.core` publics are enumerated
- **THEN** `*emit*`, `*ir-compile*`, `*compiling-aot*`, `*lg-trace*` are absent

#### Scenario: host-class registrations are not core publics
- **WHEN** `clojure.core` publics are enumerated
- **THEN** interop type entries (e.g. `java.lang.Long`, `clojure.lang.Atom`,
  `->Object`) are absent from the `clojure.core` public set

### Requirement: `clojure.core` polyfills the portable floor on demand
lg SHALL converge `clojure.core` toward the portable floor (`clojure.core ∩ bb`)
by adding the non-leaking names it lacks. Additions are **driven by consumer
need** (the jank suite and corpora), not added speculatively; the floor list is
the menu. A polyfilled name MUST match its JVM-Clojure behavior and MUST NOT
introduce a JVM abstraction leak.

#### Scenario: a polyfilled portable name matches Clojure
- **WHEN** a consumer needs `update-vals`, `partition-all`, `requiring-resolve`,
  or `volatile?` and it is polyfilled
- **THEN** it resolves in `clojure.core` and returns the JVM-Clojure result for
  the same inputs

#### Scenario: the jank suite defines the current requirement
- **WHEN** the jank `clojure-test-suite` runs
- **THEN** it reports 0 compile-skips from missing `clojure.core` builtins
  (currently 232 files / 5657 assertions / 0 skip)

#### Scenario: how to add a missing name is documented
- **WHEN** a contributor needs a `clojure.core` name lg lacks
- **THEN** a guide describes where it goes (`.lg` core vs interop ns), the
  portable-vs-leaking test, and the JVM-equivalence check

### Requirement: The xsofy corpus is the behavioral regression gate
Each increment of the surface reorganization MUST keep the xsofy corpus test
suite green: no reduction in passing assertions.

#### Scenario: xsofy stays green across an increment
- **WHEN** any bucket (① / ② / ③) of the reorg is applied and the bundle
  regenerated
- **THEN** `make -C ../xsofy LG=<dev lg> test` reports 0 fail / 0 error at the
  established baseline (281 tests, 2733 assertions)
