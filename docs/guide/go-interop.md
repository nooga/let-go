---
status: active
last-verified: 2026-07-16
human-verified:
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
   go run ./cmd/lginterop -packages github.com/zeebo/xxh3 -out pkg/rt
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

Mechanically, `cmd/lginterop` is a two-stage pipeline: the Go binary scans the
target package and extracts its exports, then drives `scripts/lginterop.lg`
(via the `lg` binary) to render the Go source. You don't need to know that to
use it, but it explains why the tool wants to run from the repo root with a
built `lg`.

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

So the core of the wrapper is nearly free:

```clojure
(defn exec! [db formatted]
  (let [[q & params] formatted]
    (apply sql/Exec db q params)))     ; variadic ...any maps directly
```

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
// column count. Returned as []any so it surfaces in let-go as a sequence.
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

Nobody needs to hand-maintain this: the lginterop emitter is written in
let-go and emits Go via `gogen`, so the shim is a generation concern — an
out-param template the emitter applies to `...any`-pointer methods (with
`ScanRow` as the first instance), keeping the whole wrapper lg-first.
Everything above it is `.lg`:

```clojure
(ns db)

(defn- row-seq [rows cols]
  (lazy-seq
    (when (.Next rows)
      (cons (zipmap cols (sql/ScanRow rows))
            (row-seq rows cols)))))

(defn query
  "Runs HoneySQL-formatted [q & params] against db; returns a fully
   realized seq of column-keyword → value maps."
  [db [q & params]]
  (let [rows (apply sql/Query db q params)]
    (try
      (let [cols (map keyword (.Columns rows))]
        (doall (row-seq rows cols)))
      (finally (.Close rows)))))

(defn exec!
  "Runs HoneySQL-formatted [q & params]; returns {:rows-affected n}."
  [db [q & params]]
  (let [r (apply sql/Exec db q params)]
    {:rows-affected (.RowsAffected r)}))
```

(Note the `doall` before `finally` closes the cursor — the seq must be
realized while `rows` is still open.)

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
| `ScanRow` out-param shim | same generated file, emitted by the lg-based gogen emitter | ~20 lines, generated |
| `db` veneer: `query`, `exec!`, `with-db` | `pkg/rt/core/db.lg` (or a userland `.lg` library) | ~1 screen |

## Reference: `cmd/lginterop`

Run from the repo root with a built `lg` binary on hand (`go build -o lg .`;
see [Usage](usage.md) for building and running `lg` generally). The tool has
two modes:

**External-package mode** — scan a Go package and generate an interop
namespace:

```
go run ./cmd/lginterop -packages <import-path>[,<import-path>...] -out pkg/rt
```

| Flag | Meaning |
|---|---|
| `-packages` | comma-separated Go import paths to wrap (overrides `deps.edn` `:gointerop`) |
| `-dir` | directory containing a `deps.edn` whose `:gointerop` key lists packages (default `.`) |
| `-out` | output directory for the generated Go files (default `.lg-interop`; use `pkg/rt` for in-tree namespaces) |
| `-smart` | generate explicit wrappers with type-specific unboxing/boxing instead of `vm.MustBox` |
| `-skeleton` | also emit a `<alias>_skeleton.lg` of `defn-` stubs to hand-customize into a veneer |

**Primitives mode** — scan `//lg:`-annotated Go sources and generate the
internal-primitive registrar (see `pkg/rt/native_prims.go` for the directive
format):

```
go run ./cmd/lginterop -primitives <dir> -go-pkg <import-path>
```

| Flag | Meaning |
|---|---|
| `-primitives` | directory containing `//lg:native`-annotated Go sources |
| `-primitives-out` | output file (default `pkg/rt/zz_primitives_generated.go`) |
| `-go-pkg` | Go import path of the scanned sources |
