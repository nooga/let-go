<!-- test/fixtures/spec-evidence/unterminated-fence.md -->
# An ordinary fence left open swallows everything after it

Real requirement. [R-swallowed]

This ```text fence is never closed, so without the unterminated-fence error
the block below it disappears from the scan -- along with its marker, so
pairing stays green and the spec silently stops running its evidence:

```text
just prose

```clj-repl @R-swallowed
(+ 1 1)
;=> 99
