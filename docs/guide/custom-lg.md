---
status: active
last-verified: 2026-08-23
---

# Building a custom `lg`

let-go's own `lg` binary wraps whatever Go packages `pkg/rt` was compiled with.
If you need bindings for a package that isn't in the tree — a database driver, a
parser, your own library — you don't have to fork let-go. Generate the bindings
into your own module and build a binary that *is* `lg`, with your namespaces
linked in.

The whole thing is about ten lines of Go plus one generator invocation.

## The module

```
mytool/
├── go.mod
├── main.go
└── interop/           # generated, do not edit
    └── interop_sqlite3.go
```

`main.go`:

```go
package main

import (
	"os"

	"github.com/nooga/let-go/pkg/cli"

	_ "example.com/mytool/interop"  // registers the generated namespaces
	_ "github.com/mattn/go-sqlite3" // driver registers itself with database/sql
)

func main() { os.Exit(cli.Main("dev", "none")) }
```

That binary is `lg`: same flags, same REPL, same resolver, same `-c` and `-b`.
[`cli.Main`](../../pkg/cli/cli.go) returns the exit code for normal runs
instead of calling `os.Exit` itself. One carve-out: flags live on the
process-global `flag.CommandLine`, which parses with `ExitOnError` — `-h` and
a malformed flag terminate the process from inside `Main`, as they do for any
Go command — so shutdown logic that must run unconditionally has to happen
before `Main` is called, not after it returns.

The two strings you pass to `Main` are *your* binary's version and commit (the
ldflags contract), and drive `-v` output. They do not become the runtime's
identity: `(System/getProperty "let-go.version")` still reports the let-go
version your module actually links, read from Go build info.

`-w` is another exception — see [Limitation: `-w` and custom
namespaces](#limitation--w-and-custom-namespaces) below.

Both imports are blank on purpose. The generated package registers its
namespaces from `init()`, and a driver like go-sqlite3 registers itself with
`database/sql` the same way — neither exports anything you call directly.

## Generating the bindings

```
go get github.com/mattn/go-sqlite3
go run github.com/nooga/let-go/cmd/lginterop@<version> \
  -packages github.com/mattn/go-sqlite3 -out-pkg interop -out ./interop
```

`-out-pkg` is what makes the output usable outside the let-go tree; without it
the generator emits `package rt` for in-tree compilation. See
[Out-of-tree generation](go-interop.md#out-of-tree-generation) for what exactly
changes, and for the module-context caveat — the scan resolves imports against
the current directory's module, so `go get` the package first.

Re-run the generator whenever you bump the wrapped package. The generated file
records its own invocation in the header — flags, any non-default alias (as
`-packages path=alias`), and the generator version when it was run
`@<version>` — so regenerating from that line round-trips. Only `-out` is
omitted: pass the directory the file lives in.

### Why `init()` calls the installer directly

Generated in-tree files hand their installer to `rt.RegisterInstaller`.
Out-of-tree files must not: `pkg/rt` drains its installer queue during its own
package `init`, and Go runs an imported package's `init` before the importing
package's. Anything queued from your module arrives after the drain and silently
never runs — no error, just a namespace that isn't there. The generator emits a
direct call instead, which is safe because `rt` is fully initialized by then.

You only need to know this if you hand-write a binding file; `-out-pkg` handles
it.

## Using it

The generated namespace is named for the package alias, and its members keep
their Go names. Using the stdlib package the e2e wraps, to show a call whose
value you can check:

```clojure
(require '[crc32])
(println (crc32/ChecksumIEEE (.getBytes "hello")))   ; => 907060870
```

Run it the way you run `lg`:

```
./mytool script.lg
./mytool -e "(require '[crc32]) (println (crc32/ChecksumIEEE (.getBytes \"hello\")))"
./mytool -r script.lg     # REPL with your namespaces loaded
```

The raw generated namespace mirrors the Go names exactly, which keeps it easy to
cross-reference against Go docs. Most projects put a hand-written `.lg` veneer
over it for idiomatic naming — that split is covered in
[Wrapping Go packages](go-interop.md).

## Distributing a single binary

`-b` bundles a script into a standalone executable by copying a base binary and
appending bytecode. **Your custom binary must be its own base**, which is the
default when you run its own `-b`:

```
./mytool -b myapp script.lg
./myapp
```

A stock `lg` as the base would produce a binary without your natives linked, and
the `(require ...)` would fail at runtime. This is also why the custom binary
has to be full `lg` rather than a bytecode-only host: compilation itself runs
top-level forms, so a veneer's `(:require)` has to resolve the native namespace
at build time too.

For cross-OS bundling, `-bundle-base` takes a path to a target-platform build of
*your* binary — build `mytool` for each target, then bundle against each.

## Limitation: `-w` and custom namespaces

`-w` cannot carry your generated namespaces into the WASM output. Unlike `-b`,
which copies the running binary, the WASM build scaffolds a *fresh* Go module
and renders its `main.go` from a fixed template that imports let-go runtime
packages only — it has no way to know about `example.com/mytool/interop`. A
script that requires your namespace will compile into the image and then fail to
resolve it at runtime.

`-w` still works from a custom binary for scripts that stay within the stock
namespaces. If you need custom Go bindings in WASM today, build the WASM module
yourself against `pkg/wasmhost` and blank-import your interop package there.

## Build metadata

`cli.Main(version, commit)` takes your binary's own version strings; they feed
`-v`. Wire them to your release tooling the way let-go wires goreleaser:

```go
var version, commit = "dev", "none"  // -X main.version=... -X main.commit=...

func main() { os.Exit(cli.Main(version, commit)) }
```

These describe *your* module, and only your module. The runtime's own
identity is resolved separately from build info: `(System/getProperty
"let-go.version")` reports the let-go your binary actually links, and `-w`
resolves the same way (honoring a local `replace`), so passing your version
here is correct and cannot mispin either one.

## A worked, hermetic example

`test/e2e/custom_main_test.go` builds exactly this shape end to end — generate
interop for a stdlib package, build a custom main against a local let-go via a
`replace` directive, run a script through it, and bundle with `-b`. It's the
shortest complete reference if something here doesn't line up.
