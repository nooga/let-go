---
status: active
last-verified: 2026-07-16
human-verified: true
reviewer: nnunley@gmail.com
---

# Wrapping Go packages for let-go (worked example: `database/sql`)

let-go programs can call into Go packages through a two-layer pattern that is
already shipping in-tree. This guide names the layers, points at the existing
exemplar, and walks through what a `database/sql` wrapper would look like —
including the parts that need a few lines of hand-written Go.

## The pattern that exists today

One Go package is already wrapped this way: **xxh3**.

1. **Generated raw namespace** — `pkg/rt/interop_xxh3.go`, produced by
   `cmd/lginterop`:

   ```
   go run ./cmd/lginterop -packages github.com/zeebo/xxh3 -opaque-structs -build-tags '!tinygo' -out pkg/rt
   ```

   The tool scans the package with `go/types` and emits an installer that
   registers each export into a namespace named after the package, using
   `vm.MustBox` to adapt Go functions to let-go values:

   ```go
   func installXxh3NS() {
       ns := vm.NewNamespace("xxh3")
       ns.Def("Hash", vm.MustBox(xxh3.Hash))
       ns.Def("New", vm.MustBox(xxh3.New))
       // ...
       RegisterNS(ns)
   }
   ```

   Boxed struct values stay opaque and dispatch methods via reflection, so
   `(.WriteString h "abc")` and `(.Sum64 h)` work on an `xxh3.Hasher` without
   any per-method wrapper.

2. **`.lg` veneer** — `pkg/rt/core/hash.lg`. The generated namespace exposes
   raw Go names (`xxh3/HashString`); the user-facing API — friendly kebab-case
   names, idiomatic shapes (a hasher as a map of callables) — is written in
   let-go on top:

   ```clojure
   (ns hash)
   (def xxh3-64-str xxh3/HashString)

   (defn xxh3-hasher []
     (hasher-map (xxh3/New)))
   ```

Keep this split in mind: **generation gets you the raw surface; the library
you'd actually want to use is a `.lg` file.**

One load-order asymmetry to know: the generated `xxh3` namespace is
registered natively at startup, so `(xxh3/HashString "abc")` works with no
`require` — but the `.lg` veneer loads on demand, so `(hash/xxh3-64-str …)`
needs `(require 'hash)` first:

```clojure
(require 'hash)
(hash/xxh3-64-str "abc")            ; => 8696274497037089104
(= (hash/xxh3-64-str "abc")
   (let [h (xxh3/New)]
     (.WriteString h "abc")
     (.Sum64 h)))                   ; => true
```

Mechanically, `cmd/lginterop` is a two-stage pipeline: the Go binary scans the
target package and extracts its exports, then drives the embedded
`cmd/lginterop/lginterop.lg` emitter to render the Go source. Both stages run
in the same process — the tool links the let-go runtime, so `gogen` is already
there — which is why it needs neither a let-go checkout nor an `lg` binary.

Related docs: [Embedding let-go in Go](embedding-in-go.md) covers the host
side (running let-go inside your own Go program), which is the other way to
expose Go functions — hand-registering them with `vm.MustBox` from your host,
no code generation involved.

## Worked example: `database/sql`

The motivating use case: [HoneySQL](https://github.com/seancorfield/honeysql)
runs on let-go (see `test/compat/`), and `honey.sql/format` returns
`[sql-string & params]`. Go's `database/sql` has almost the same shape on the
other side:

```go
func (db *DB) Exec(query string, args ...any) (Result, error)
func (db *DB) Query(query string, args ...any) (*Rows, error)
```

The shapes line up, but the parameters need one small piece of Go to get
across. `Exec` is a method, and methods are not first-class values in let-go,
so `apply` has nothing to spread the parameter vector into. A shim takes the
parameters as a slice and spreads them on the Go side:

```go
func Exec(db *sql.DB, query string, args []any) (sql.Result, error) {
    return db.Exec(query, args...)
}
```

```clojure
(defn exec! [db formatted]
  (let [[q & params] formatted]
    (shim/Exec db q (vec params))))
```

See [Why `Query` needs a shim too](#why-query-needs-a-shim-too) for the full
reasoning; it applies to every variadic method.

### What generation + boxing gives you directly

- `sql/Open`, `(.Close db)`, `(.Ping db)` — connection setup/teardown.
- `(.Exec db q & args)` / `(.Query db q & args)` — statement execution;
  boxed variadic `...any` parameters accept let-go strings, ints, floats,
  nil.
- `Result` methods: `(.LastInsertId r)`, `(.RowsAffected r)`.
- `Rows` cursor movement: `(.Next rows)`, `(.Close rows)`,
  `(.Columns rows)`.

`database/sql` is standard library, so the wrapper adds **no new go.mod
dependency**. Only drivers are external, and they stay host-side (see below).

### How Go return values arrive

A boxed Go function's results map to let-go like this:

| Go signature | let-go result |
|---|---|
| `func(…)` | `nil` |
| `func(…) T` | the boxed `T` |
| `func(…) error` | the boxed error **as a value** — does not throw |
| `func(…) (T, error)` | boxed `T`, throwing when the error is non-nil |
| `func(…) (A, B)` | vector `[a b]` |
| `func(…) (A, B, …, error)` | vector `[a b …]`, throwing when the error is non-nil |

A trailing `error` becomes a throw only when some other result remains to
return. `func() error` — `Close`, `Flush`, and friends — hands the error back
as an ordinary value, so `(if (.Close f) ...)` keeps working.

The type test is on the *declared* result type, so a function returning a
concrete `*MyError` is returning a value, not signalling failure. An `error` in
any position but last is likewise an ordinary value.

Known limitation: Go **array** results (`[N]T`, as opposed to slices) are not
boxable and surface as an error. This is not specific to multi-return —
`func() [1]int` has always behaved this way.

### The one seam: `Scan`'s out-parameters

The interop layer already does the pointer work you'd expect: boxed values
hold Go pointers (`*sql.DB`, `*sql.Rows`) and dispatch every method on them
reflectively (`pkg/vm/boxed.go` builds a callable per method, variadics
included); registered structs write back through pointers
(`vm.StructMapping.RecordToStruct`); sequences convert per-element into Go
slices. So `(.Next rows)`, `(.Columns rows)`, `(.Close rows)` all work with
no wrapper at all.

The one pattern that doesn't map automatically is `Scan`'s contract:

```go
func (rs *Rows) Scan(dest ...any) error
```

`Scan` communicates results by writing **through caller-allocated pointers**,
and let-go values are immutable — `(.Scan rows x)` passes values, not
destinations. The bridge is a shim that allocates the destinations, calls
`Scan`, and returns the results as a value:

```go
// ScanRow reads the current row into a slice of values, sized from the
// column count. []any is the right return type: let-go boxes each element
// by its dynamic type, so a string column arrives as a let-go string and a
// NULL arrives as nil.
func ScanRow(rows *sql.Rows) ([]any, error) {
    cols, err := rows.Columns()
    if err != nil {
        return nil, err
    }
    dest := make([]any, len(cols))
    ptrs := make([]any, len(cols))
    for i := range dest {
        ptrs[i] = &dest[i]
    }
    if err := rows.Scan(ptrs...); err != nil {
        return nil, err
    }
    return dest, nil
}
```

This is hand-written today. The emitter has no out-parameter template, so a
wrapper needing `Scan` carries a small Go package of its own alongside the
generated one. That package registers its functions as an ordinary namespace:

```go
func init() {
    ns := vm.NewNamespace("shim")
    ns.Def("ScanRow", vm.MustBox(ScanRow))
    ns.Def("Query", vm.MustBox(Query))
    ns.Def("Exec", vm.MustBox(Exec))
    rt.RegisterNS(ns)
}
```

Call `rt.RegisterNS` directly from `init` rather than through
`rt.RegisterInstaller`, for the same reason out-of-tree generated packages do:
`pkg/rt` drains its installer queue during its own package init, which Go runs
before the importing package's, so anything queued from here would arrive
after the drain and silently never run.

Everything above that seam is `.lg`:

```clojure
(ns db)

(defn- row-seq [rows cols]
  (lazy-seq
    (when (.Next rows)
      (cons (zipmap cols (shim/ScanRow rows))
            (row-seq rows cols)))))

(defn query
  "Runs HoneySQL-formatted [q & params] against db; returns a fully
   realized seq of column-keyword → value maps."
  [db [q & params]]
  (let [rows (shim/Query db q (vec params))]
    (try
      (let [cols (map keyword (.Columns rows))]
        (doall (row-seq rows cols)))
      (finally (.Close rows)))))

(defn exec!
  "Runs HoneySQL-formatted [q & params]; returns {:rows-affected n}."
  [db [q & params]]
  (let [r (shim/Exec db q (vec params))]
    {:rows-affected (.RowsAffected r)}))
```

(Note the `doall` before `finally` closes the cursor — the seq must be
realized while `rows` is still open.)

#### Why `Query` needs a shim too

An earlier version of this guide wrote `(apply sql/Query db q params)`. That
cannot work, and the reason is worth knowing before you reach for it:
`Query` is a *method* on `*sql.DB`, methods are not first-class values in
let-go, and `apply` needs a function value to spread arguments into. Calling
`(.Query db q p1 p2)` with literal arguments is fine; a wrapper holding its
parameters in a vector has no way to spread them.

So `Query` and `Exec` take the parameters as a slice and spread them on the
Go side:

```go
func Query(db *sql.DB, query string, args []any) (*sql.Rows, error) {
    return db.Query(query, args...)
}

func Exec(db *sql.DB, query string, args []any) (sql.Result, error) {
    return db.Exec(query, args...)
}
```

`[]any` again, and again it is load-bearing: a let-go vector converts into
`[]any` element by element, and a let-go `nil` becomes a Go `nil`, so a NULL
parameter passes through like any other value.

### How collections cross the boundary

`[]any` is the natural Go type for "a row of unknown column types", and
`map[string]any` the natural one for a JSON object or an options bag, so it is
worth stating what happens to them in each direction.

**Go to let-go.** Each element is boxed by its *dynamic* type, not the
slice's static element type. A `[]any` holding `"Ada"`, `int64(36)` and
`nil` arrives as a let-go vector of a string, an int and `nil`, so ordinary
comparisons work:

```clojure
(= "Ada" (first (shim/ScanRow rows)))   ; => true
```

Elements the boxing layer has no native form for stay wrapped as opaque
values, exactly as they would outside a `[]any`. Nesting works: a `[]any`
inside a `[]any` converts recursively.

**let-go to Go.** A let-go sequence converts into a `[]any` parameter
element by element, and the conversion always succeeds, because `any` can
hold every Go value. A let-go `nil` becomes a Go `nil` rather than failing
the call, which is what makes NULL parameters work.

Element values are whatever `Value.Unbox` yields, so a let-go int arrives as
a Go `int`. Convert explicitly in the shim if a driver wants something
narrower.

#### Maps

A map crosses the same way, so a shim can declare `map[string]any` — or
`map[string]string`, or any other Go map — and let-go fills it in.

**Go to let-go.** Keys and values are boxed by their dynamic type, exactly as
slice elements are, so a `map[string]any` arrives as a let-go map of native
values. A nil Go map is `nil`; an empty one is an empty map, which is a
different value.

Go string keys become let-go **strings, not keywords**, so a map is read with
`get` and a string:

```clojure
(get (shim/Stats) "runtime")   ; => "let-go"
(:runtime (shim/Stats))        ; => nil
```

That is lossy but predictable, and it matches `clojure.data.json` with no
`:key-fn`. Keywordising automatically would be wrong: Go map keys may contain
spaces and dots, which do not make valid keywords.

**let-go to Go.** A let-go map converts into a Go map key by key. Keyword keys
become string keys — a keyword stores its name without the leading colon — so
`{:name "Ada"}` reaches Go as `map[string]any{"name": "Ada"}`. A numeric key
into a string-keyed map is rejected rather than converted: Go would turn the
key `65` into `"A"` by rune conversion, and a wrong key that looks plausible is
worse than a failure.

Both let-go map types convert, sorted maps included.

#### Nesting

Nesting works in both directions and recursively. A `[]any` inside a `[]any`,
a map inside a map, a map inside a `[]any` — each converts to its natural
shape rather than arriving as an opaque let-go value:

```clojure
{:a {:b 1}}   ; → map[string]any{"a": map[string]any{"b": 1}}
```

Two limits are worth knowing, both on the **Go to let-go** side.

Unwrapping a value out of an `any` slot follows the same allowlist as scalars,
and for maps that allowlist is exactly `map[string]any`. A nested
`map[string]string` stays an opaque value, even though the *same* map converts
fine when it is the whole return value rather than nested inside one. Return
`map[string]any` from a shim if the map has to survive nesting.

A **lazy sequence** in an `any` slot is handed over unrealized rather than
converted, because it may be infinite. Convert it in let-go first (`vec`,
`doall`) if the Go side wants a slice.

### Drivers

`sql.Open("sqlite3", ...)` only works if a driver registered itself, and
drivers register via blank import in Go:

```go
import _ "github.com/mattn/go-sqlite3"
```

There is no way to do that from `.lg` — driver selection is a **host/build
decision**. Two workable postures:

- **Embedding host**: your Go program that embeds let-go imports the driver;
  the wrapper works out of the box.
- **In-tree, opt-in**: a build tag gates the driver import (the same pattern
  `pkg/glplat` uses for its cgo/GLFW backend — `-tags glplat`), so plain
  builds of `lg` stay dependency-free and `-tags sqlite` ships a
  batteries-included binary.

### Putting it together

A contributed `database/sql` wrapper would be three small pieces:

| Piece | Where | Size |
|---|---|---|
| Generated/boxed raw surface (`sql/Open`, method dispatch on DB/Rows/Result) | `pkg/rt/interop_sql.go` | generated |
| `ScanRow` out-param shim, plus `Query`/`Exec` parameter shims | a small hand-written Go package next to the generated one | ~40 lines |
| `db` veneer: `query`, `exec!`, `with-db` | `pkg/rt/core/db.lg` (or a userland `.lg` library) | ~1 screen |

## Reference: `cmd/lginterop`

The tool links the let-go runtime and runs its emitter in-process, so it needs
neither a let-go checkout nor an `lg` binary — `go run
github.com/nooga/let-go/cmd/lginterop@<version>` works from any module (see
[Usage](usage.md) for building and running `lg` generally). Aliases must be
unique across a run — two packages
resolving to the same alias would write the same `interop_<alias>.go`, so the
tool refuses up front. Set `LGINTEROP_KEEP_SCRIPT=1` to dump the
intermediate `.lg` driver script to a temp file for inspection. The tool has
two modes:

**External-package mode** — scan a Go package and generate an interop
namespace:

```
go run ./cmd/lginterop -packages <import-path>[,<import-path>...] -out pkg/rt
```

The generated file starts with a header recording the invocation — flags, any
non-default alias (as `-packages path=alias`), and, when the generator was run
`@<version>`, a `Generated with lginterop <version>` line — so regenerating
from the header's own command round-trips byte-identically. The one thing the
header omits is `-out`: the file's own location supplies the destination, and
baking a path in would make the bytes differ per checkout. E2e golden tests
(`test/e2e/lginterop_regen_test.go`) hold `interop_xxh3.go` and a
non-default-alias output to that round trip.

Aliases that would collide with the generated file's own imports (`vm`, or
`fmt` in smart mode) are rejected up front, as are two aliases that normalize
to the same `interop_<alias>.go` (`-` and `.` become `_` in filenames). A
scanned package with no eligible exports is reported as skipped in the final
summary rather than counted as generated.

| Flag | Meaning |
|---|---|
| `-packages` | comma-separated Go import paths to wrap (overrides `deps.edn` `:gointerop`); `path=alias` pins a non-default namespace alias, like `deps.edn`'s `{"path" "alias"}` form |
| `-dir` | directory containing a `deps.edn` whose `:gointerop` key lists packages (default `.`) |
| `-out` | output directory for the generated Go files (default `.lg-interop`; use `pkg/rt` for in-tree namespaces) |
| `-smart` | generate explicit wrappers with type-specific unboxing/boxing instead of `vm.MustBox` |
| `-skeleton` | also emit a `<alias>_skeleton.lg` of `defn-` stubs to hand-customize into a veneer |
| `-opaque-structs` | skip `vm.RegisterStruct`: struct types stay `vm.Boxed` and dispatch methods reflectively, instead of flattening to field-only Records — required when the API is used through methods (xxh3's `Hasher` and its `.WriteString`/`.Sum64`) |
| `-build-tags` | emit `//go:build <constraint>` as the first line of each generated file (recorded in the header so regeneration round-trips). Used for xxh3 as `'!tinygo'`: the reflect-boxed bindings can't run under TinyGo |
| `-out-pkg` | emit a self-contained package of this name for use **outside** the let-go tree (see below). Omit it for the in-tree `package rt` output |

### Out-of-tree generation

By default the emitter writes `package rt` with unqualified calls, because the
generated file is compiled *as part of* `pkg/rt`. A third-party module can't
use that. `-out-pkg <name>` switches to a self-contained form:

```
go run github.com/nooga/let-go/cmd/lginterop@<version>   -packages github.com/mattn/go-sqlite3 -out-pkg interop -out ./interop
```

Three things change, and nothing else does:

- `package <name>` instead of `package rt`.
- `github.com/nooga/let-go/pkg/rt` joins the imports, and registration becomes
  `rt.RegisterNS(ns)` rather than the unqualified call.
- `init()` calls `install<Alias>NS()` **directly** instead of handing it to
  `RegisterInstaller`.

That last one is the load-bearing detail. `pkg/rt` drains its installer queue
during its own package init (`pkg/rt/zz_run_installers.go`), and Go runs an
imported package's `init` *before* the importing package's. An out-of-tree file
calling `RegisterInstaller` would therefore enqueue after the drain and
silently never run — no error, just a missing namespace. Calling the install
function directly is safe: `rt` is fully initialized by the time the importing
package's `init` fires, `RegisterNS` is mutex-guarded, and the ordering
relative to `LoadCore` is the same as an in-tree installer's.

The flag is recorded in the generated-by header, so regenerating from that
header's own command round-trips. `-out-pkg rt` is rejected: omitting the flag
already selects the in-tree output, and accepting `rt` would give one name two
meanings.

Blank-import the generated package from your `main` to get the namespace
registered; see [Building a custom `lg`](custom-lg.md) for the full module
layout.

#### Module-context caveat

Scanning uses the `go/types` **source** importer, which resolves imports
against the module context of the current working directory. A package is only
scannable from inside a module that already requires it:

```
go get github.com/mattn/go-sqlite3   # first
go run .../cmd/lginterop -packages github.com/mattn/go-sqlite3 -out-pkg interop
```

Scanning from an unrelated directory fails with an import error naming the
package and this hint. Standard-library packages are always resolvable.

**Primitives mode** — scan `//lg:`-annotated Go sources and generate the
internal-primitive registrar:

```
go run ./cmd/lginterop -primitives <dir> -go-pkg <import-path>
```

| Flag | Meaning |
|---|---|
| `-primitives` | directory containing `//lg:native`-annotated Go sources |
| `-primitives-out` | output file (default `pkg/rt/zz_primitives_generated.go`) |
| `-go-pkg` | Go import path of the scanned sources |

### How primitives mode works, and why it's a separate mode

Primitives mode points the other way: external-package mode wraps *someone
else's* Go package into a new namespace for let-go code to call; primitives
mode implements *let-go's own standard library* in Go — the native backing
for names like `clojure.core/name` or `clojure.string/upper-case`. A
primitive is an ordinary named Go function with a **typed signature**,
annotated with comment directives that carry everything registration needs:

```go
// Mirrors clojure.core/name — `(name x)`.
//
//lg:native
//lg:name name
func Name(v vm.Value) (string, error) { ... }
```

| Directive | Meaning |
|---|---|
| `//lg:native` | marks the function as a primitive (required) |
| `//lg:name <sym>` | the let-go symbol (default: kebab-cased Go name) |
| `//lg:ns <ns>` | target namespace (default `clojure.core`) |
| `//lg:private` | register private — dropped from `ns-publics` |

The scanner reads the signature itself: a leading `*vm.ExecContext` parameter
threads the caller's context, a trailing `error` result becomes a let-go
throw, a variadic final parameter becomes rest-args, and several functions
sharing one `//lg:name` become a multi-arity definition.

The generated registrar differs from external-mode output in kind, which is
why the modes are separate rather than one flag:

- **Typed adapters, no reflection.** External mode leans on `vm.MustBox`
  (reflective call and boxing at every invocation — fine for wrapping an
  arbitrary package's surface). Primitives are the hot path of the language,
  so the generator emits a hand-shaped adapter per function with
  compile-time unbox/box conversions.
- **Direct-call metadata.** Each (namespace, Go package) group registers a
  `NativeModule`, which the AOT lowering uses to compile call sites into
  direct Go calls — guarded by `rt.NativePrimsIntact()` so `with-redefs` /
  `alter-var-root` still win (lowered code falls back to var dispatch when a
  guarded root has been overridden).
- **Bootstrap lifecycle.** Core namespaces load `.lg` bootstrap definitions
  on demand, which would clobber the native adapters; the registrar records
  every generated binding and re-applies it after a namespace load completes
  (`pkg/rt/native_prims_lifecycle.go`), while explicit user redefinition —
  which happens after loading — is left alone.

Of these, only the **bootstrap lifecycle** is inherently a stdlib problem —
it exists because core namespaces have competing `.lg` bootstrap definitions
for the same names, which no external wrapper has. The other two are not
principled privileges of the stdlib: an external package on a hot path
deserves typed adapters and direct-call metadata just as much, and the
`-smart` flag is the (currently vestigial) seam where external mode was
meant to grow them. The two modes should converge on one emitter with the
lifecycle machinery applied only where dual `.lg`+native definitions exist;
#531 tracks that consolidation. See `pkg/rt/native_prims.go` for live
examples of the directive format.

Note that you rarely run primitives mode by hand: `make generate` invokes it
(via `scripts/generate.lg`) alongside the bundle and lowered-tree
regeneration, so annotating a function and running `make generate` is the
whole workflow.
