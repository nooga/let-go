---
status: active
last-verified: 2026-06-18
human-verified:
---

# Portable code (`:lg` reader conditionals)

let-go ships some namespaces of its own — e.g. `let-go.semver` — that JVM
Clojure can't load. To keep shared code loadable on both, guard the let-go-only
parts behind `:lg` reader conditionals. The reader always matches `:lg` and
`:default`, and matches `:clj` / `:bb` only when opted in. JVM Clojure has no
idea what `:lg` is, so it skips those branches entirely — the same way it skips
a `:cljs` branch:

```clojure
(ns my.app
  ;; only let-go reads the :lg branch; Clojure never tries to load let-go.semver
  #?(:lg (:require [let-go.semver :as semver])))

(defn normalize [s]
  ;; the semver alias appears only inside the :lg branch, so a non-let-go reader
  ;; never sees an unresolved symbol
  #?(:lg (semver/render (semver/version s))
     :default s))
```

This has to be guarded at **read** time: a missing namespace or an unresolved
symbol fails at compile time, before any `when`/`if` could intervene. Two
things to know:

- **Use `.cljc`.** Clojure only honors `#?` in `.cljc` files. let-go reads `#?`
  in any file and its loader resolves `.lg` → `.cljc` → `.clj`, so a shared file
  should just be `.cljc`.
- **Put `:lg` before `:clj`.** First match wins. If a let-go user opted into
  `:clj` matching to consume a Clojure library, then in `#?(:clj … :lg …)`
  let-go would take the `:clj` branch.
