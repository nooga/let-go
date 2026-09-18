<!-- test/fixtures/spec-evidence/basic.md -->
# Fixture

Addition holds. [R-add]

```clj-repl @R-add
(+ 1 2)
;=> 3
```

Printing holds. [R-print]

```clj-repl @R-print oracle=none
(println "x")
;; out: "x\n"
;=> nil
```

Rows hold. [R-rows]

<!-- evidence: @R-rows -->
| form | value | out |
|---|---|---|
| `(inc 1)` | `2` | |
| `(print "a")` | `nil` | `"a"` |

Expected-to-fail rows hold. [R-fail-row]

<!-- evidence: @R-fail-row -->
| form | value | out | expect |
|---|---|---|---|
| `(inc 1)` | `3` | | fail |

A test block. [R-block]

```clj-test @R-block
(ns fixture.block (:require [clojure.test :refer [deftest is]]))
(deftest t (is (= 1 1)))
```
