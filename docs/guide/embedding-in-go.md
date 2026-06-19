---
status: active
last-verified: 2026-06-19
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
