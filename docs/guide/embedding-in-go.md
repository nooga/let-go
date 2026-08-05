---
status: active
last-verified: 2026-07-30
human-verified:
---

# Embedding let-go in Go

let-go embeds cleanly as a scripting layer for Go programs. Define Go values and
functions, hand them to the VM, run user-supplied Clojure against your data. Go
structs roundtrip as records, Go channels are first-class let-go channels, and Go
functions are callable from let-go.

```go
import (
    "fmt"

    "github.com/nooga/let-go/pkg/api"
    "github.com/nooga/let-go/pkg/vm" // used in the sections below
)

c, _ := api.NewLetGo("myapp")

c.Def("x", 42)
c.Def("greet", func(name string) string {
    return "Hello, " + name
})

v, _ := c.Run(`(greet "world")`)
fmt.Println(v) // "Hello, world"
```

## Structs roundtrip as records

Registered structs become records on the let-go side. Unmutated values unbox back
to the original Go type for free; mutated ones go through `vm.ToStruct[T]`.

```go
type Item struct{ Name string; Price float64; Qty int }
vm.RegisterStruct[Item]("myapp/Item")

c.Def("item", Item{Name: "Widget", Price: 9.99, Qty: 5})
c.Run(`(:name item)`)                  // "Widget"
c.Run(`(* (:price item) (:qty item))`) // 49.95
```

## Go channels are let-go channels

Go channels and `vm.Chan` plug into `go` / `<!` / `>!` directly:

```go
inch := make(chan int)
outch := make(vm.Chan)
c.Def("in", inch)
c.Def("out", outch)

c.Run(`(go (loop [i (<! in)]
             (when i
               (>! out (inc i))
               (recur (<! in)))))`)
```

[`pkg/api/interop_test.go`](../../pkg/api/interop_test.go) has the full set of
embedding examples (defs, structs, channels, function calls).

## Dropping `net/http` with `lg_no_http`

The `http` namespace is registered from an `init()`, so `net/http` — and
`crypto/tls` and `crypto/x509` behind it — is live in every binary that links
`pkg/rt`, whether or not the embedded program ever opens a socket. The linker
cannot see that it is unused.

Build with `-tags lg_no_http` to leave it out:

```sh
go build -tags lg_no_http ./cmd/lg-runtime
```

Measured on darwin/arm64 against `cmd/lg-runtime`:

| build | bytes | |
|---|---|---|
| default | 17,594,402 | |
| `lg_no_http` | 12,175,106 | −5,419,296 (−30.8%) |
| default, `-s -w` | 12,131,874 | |
| `lg_no_http`, `-s -w` | 8,382,754 | −3,749,120 (−30.9%) |

What changes in the tagged build: the `http` namespace is not installed, so
`http/get`, `http/serve` and friends do not resolve; and `io/reader` on an
`io/url` record returns an error instead of fetching. Everything else is
untouched — file and resource `slurp`, `io/reader` on a path, and the rest of
`io` behave the same.

TinyGo has never had a working `http` namespace (its `net/http` does not
compile), so that lane already builds this way and gains nothing new from the
tag. Both lanes now share one stub, `pkg/rt/http_stub.go`.
