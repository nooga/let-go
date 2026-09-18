<!-- test/fixtures/spec-evidence/lint-extra.md -->
Bad expectation. [R-badexp]

```clj-repl @R-badexp
(+ 1 2)
;; nonsense: 3
```

Two forms on one line. [R-twoforms]

```clj-repl @R-twoforms
(+ 1 2) (+ 3 4)
;=> 3
```

Test block marked expect=fail (unsupported). [R-testfail]

```clj-test @R-testfail expect=fail
(ns fixture.testfail (:require [clojure.test :refer [deftest is]]))
(deftest t (is (= 1 1)))
```
