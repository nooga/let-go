---
status: active
last-verified: 2026-06-18
human-verified:
---

# Version requirements (`let-go.semver`)

`let-go.semver` provides SemVer values that order correctly through `compare` /
`sort` / `sorted-set`, plus range matching and a host-version assertion.

## Range matching

`satisfies-range?` understands comparators (`>= <= < > = !=`, space-AND-composed),
bare/partial versions and x-ranges (`1.2.x`, `1.*`, `*`), npm-style caret/tilde
(`^1.2.3`, `~1.2`), and `||` OR-composition:

```clojure
(require '[let-go.semver :as semver])
(semver/satisfies-range? "1.4.0" "^1.2.3")          ; => true
(semver/satisfies-range? "2.0.0" "^1.2.3")          ; => false
(semver/satisfies-range? "1.5.0" "^1.0.0 || ^2.0.0"); => true
```

## Asserting the host build (`require-letgo`)

`require-letgo` asserts, at load time, that the running `lg` build is new enough
— failing with one clear line instead of a "can't resolve" cascade. The spec is
auto-detected: a 7–40 hex string is a commit pin (prefix-matched), anything else
is a semver range. It warns-and-passes when the build is unknown (a `dev` /
`none` build), so it never blocks REPL/dev work; known mismatches throw an
`ex-info` whose message is that one clear line and whose data is
`{:required :found :check-type}` for programmatic handling.

`require-letgo` is let-go-specific, so guard both the `:require` and the call
with [`:lg` reader conditionals](portability.md) to keep shared `.cljc` loadable
on JVM Clojure:

```clojure
(ns my.app
  #?(:lg (:require [let-go.semver :refer [require-letgo]])))

#?(:lg (require-letgo ">=1.9.0"))   ; one clear failure line on too-old lg
```
