---
status: active
last-verified: 2026-09-29
authoritative-for:
  - custom-data-readers
  - raw-go-reader
---

# Custom data readers and raw Go fragments

let-go supports Clojure-style custom tagged literals through
`clojure.core/*data-readers*`, which is empty by default as in Clojure, with
`#inst` and `#uuid` in `default-data-readers`. A native raw `#go{...}` reader for
tools and tests that need readable Go source fragments is installed by importing
the Go package `pkg/gofragments`; the `lg` CLI imports it unless built with
`-tags lg_no_gofragments`.

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

let-go does not yet discover classpath `data_readers.clj` files.
Explicit registry entries and dynamic `*data-readers*` entries take precedence
over `default-data-readers`, which holds `#inst` and `#uuid`. Rebinding
`*data-readers*` leaves `default-data-readers` in place.

## Tags nothing registers

A tag that neither `*data-readers*` nor `default-data-readers` handles goes to
`*default-data-reader-fn*`, called with the tag symbol
and the form. When that is nil, or returns nil, reading the tag throws
`No reader function for tag <tag>`, as in Clojure.

`tagged-literal` builds the usual default: a value that keeps the tag and the
form, answers `:tag` and `:form`, compares by value, and prints back as the
literal it came from.

```clojure
(binding [*default-data-reader-fn* tagged-literal]
  (let [t (read-string "#app/unknown [1 2]")]
    [(tagged-literal? t) (:tag t) (:form t) (pr-str t)]))
;; => [true app/unknown [1 2] "#app/unknown [1 2]"]
```

The resolution order lives in `clojure.core/-read-tagged`, in `core.lg`: the
Go reader reads the tag and the form, and hands them to it.

## Raw `#go{...}` fragments

Importing `pkg/gofragments` adds a `go` entry to the root of `*data-readers*`
(`import _ "github.com/nooga/let-go/pkg/gofragments"`); the `lg` CLI imports it
unless built with `-tags lg_no_gofragments`, and an embedding program opts in. The entry's value is a native raw reader, so it
consumes the payload text itself instead of a read form. It is an ordinary
entry: binding `*data-readers*` to a map without `go` drops it, and extending
the map with `assoc` keeps it.

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
`#go` and the opening brace. A missing opening brace is a reader error. A
truncated fragment is incomplete input, like an unterminated list, so a REPL
keeps prompting for the rest.

The body is returned verbatim, including newlines: a fragment written as
`#go{`, a newline, the code, a newline, and `}` begins and ends with `\n`.
Compare such fragments byte-for-byte only after accounting for that.

Clojure-aware formatters (cljfmt, zprint, editor format-on-save) do not know
raw fragments. They treat `#go{...}` as a tag applied to a map and may insert or
remove spaces inside the payload, which changes the string it reads as. Exclude
files or regions that hold byte-exact fragments from formatting.

In an unselected reader-conditional branch, `#go` consumes its payload only when
the next non-whitespace character is `{`, at any nesting depth. Any other
payload is skipped like an ordinary tagged form, so `#?(:clj #go [1] :default 7)`
reads as `7`.

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
A raw callback receives `TaggedRawInput`, whose `NextRune` and `UnreadRune`
methods preserve the enclosing reader's source position. Raw callbacks must
consume exactly one payload and be side-effect free: an unselected reader-
conditional branch invokes them, at any nesting depth, solely to skip that
payload safely. Return an error wrapping `io.EOF` for truncated input so
callers can tell incomplete input from a malformed payload.

`NewLispReaderWithTaggedReaders` installs the same registry for callers that use
the reader directly. An explicit registry is consulted before `*data-readers*`.
