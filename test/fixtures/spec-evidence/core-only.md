<!-- test/fixtures/spec-evidence/core-only.md -->
Arithmetic. [R-arith]

```clj-repl @R-arith
(+ 1 2)
;=> 3
(str "a" "b")
;=> "ab"
```

Output. [R-out]

<!-- evidence: @R-out -->
| form | value | out |
|---|---|---|
| `(println "hi")` | `nil` | `"hi\n"` |

A deliberate mismatch. [R-wrong]

```clj-repl @R-wrong oracle=none
(+ 1 1)
;=> 3
```

Known-bad behavior, expected to fail. [R-xfail]

```clj-repl @R-xfail expect=fail
(+ 1 1)
;=> 3
(+ 2 2)
;=> 4
;; expect: pass
```
