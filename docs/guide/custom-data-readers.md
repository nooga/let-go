---
status: active
last-verified: 2026-08-21
authoritative-for:
  - custom-data-readers
---

# Custom data readers

let-go supports Clojure-style custom tagged literals through
`clojure.core/*data-readers*` and explicit per-compiler registries for embedding
programs.

## Register a data reader from let-go

`*data-readers*` is a dynamic map from tag symbols to one-argument functions or
Vars holding those functions. The function receives the normally read form
after the tag and its return value becomes the value of the tagged literal.

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

Custom entries take precedence over built-in `#uuid` and `#inst` handlers.
let-go does not yet discover classpath `data_readers.clj` files and does not
implement `*default-data-reader-fn*`. An unregistered tag read without an
explicit Go registry retains let-go's legacy behavior of returning its payload.

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

`NewLispReaderWithTaggedReaders` installs the same registry for callers that use
the reader directly. With an explicit registry installed, an unknown custom tag
is an error rather than a legacy best-effort read.
