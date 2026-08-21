---
status: active
last-verified: 2026-08-21
authoritative-for:
  - custom-data-readers
  - raw-go-reader
---

# Custom data readers and raw Go fragments

let-go supports Clojure-style custom tagged literals through
`clojure.core/*data-readers*`. It also provides a built-in raw `#go{...}`
reader for tools and tests that need readable Go source fragments.

## Register a data reader from let-go

`*data-readers*` is a dynamic map from tag symbols to one-argument functions or
Vars holding those functions. The function receives the normally read form after the tag and its return value
becomes the value of the tagged literal.

```clojure
(binding [*data-readers*
          (assoc *data-readers* 'app/point
                 (fn [[x y]] {:x x :y y}))]
  (read-string "#app/point [10 20]"))
;; => {:x 10 :y 20}
```

Use `alter-var-root` when a registration must apply process-wide:

```clojure
(alter-var-root #'*data-readers* assoc 'app/id identity)
(read-string "#app/id {:id 42}")
;; => {:id 42}
```

Dynamic registrations are honored by `read-string`, `read-all-string`, and
`load-string`. Root registrations are also visible while ordinary source files
are compiled. A registration must exist before its tagged literal is read; a
registration nested in the same unread form cannot affect that form.

let-go does not yet discover classpath `data_readers.clj` files and does not
implement `*default-data-reader-fn*`. An unregistered tag read without an
explicit Go registry retains let-go's legacy behavior of returning its payload.
Explicit registry entries and dynamic `*data-readers*` entries take precedence
over built-in `#uuid` and `#inst` handlers and the default raw `#go` handler.

## Raw `#go{...}` fragments

`#go` consumes a balanced brace-delimited payload and returns its body as a
string, excluding the outer braces:

```clojure
#go{if ready {
  return fmt.Errorf("not ready: }")
}}
;; => "if ready {
  return fmt.Errorf("not ready: }")
}"
```

Brace balancing ignores braces inside interpreted strings, raw strings, rune
literals, `//` comments, and `/* ... */` comments. Whitespace is allowed between
`#go` and the opening brace. Truncated fragments and missing opening braces are
reader errors.

`#go` only reads source text. It does not parse, compile, or execute Go. Its
primary use is supplying readable Go fragments to Go-AST-based tooling and
production-structure tests.

## Register readers from an embedding Go program

Embedding code can install explicit per-compiler readers:

```go
registry := compiler.NewTaggedReaderRegistry()
err := registry.RegisterData("app/id", func(v vm.Value) (vm.Value, error) {
    return v, nil
})
if err != nil {
    return err
}

ctx := compiler.NewCompiler(consts, ns).SetTaggedReaders(registry)
_, result, err := ctx.CompileMultiple(source)
```

Use `RegisterRaw` when a tag has syntax that is not an ordinary let-go form.
A raw callback receives `TaggedRawInput`, whose `ReadRune` and `UnreadRune`
methods preserve the enclosing reader's source position. Raw callbacks must
consume exactly one payload and be side-effect free: an unselected reader-
conditional branch may invoke them solely to skip that payload safely.

`NewLispReaderWithTaggedReaders` installs the same registry for callers that use
the reader directly. With an explicit registry installed, an unknown custom tag
is an error rather than a legacy best-effort read.
