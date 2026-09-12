<!-- test/fixtures/spec-evidence/fenced-marker.md -->
# Markers inside ordinary fences are documentation

Real requirement. [R-real]

```clj-repl @R-real
(+ 1 1)
;=> 2
```

The fence below SHOWS the grammar. The markers inside it are examples, not
requirements, and must not be collected:

```text
Some prose with a marker [R-phantom] in it.
[R-also-phantom]
```

A yaml fence showing the opt-out must not opt this spec out either:

```yaml
---
evidence: skip
---
```

A four-backtick fence is closed only by four backticks, so the three-backtick
fence nested inside this one does not end it early:

````markdown
```clj-repl @R-phantom-nested
[R-phantom-nested]
```
````
