---
status: active
last-verified: 2026-11-22
authoritative-for:
  - load-bundle
human-verified:
---

# load-bundle: running code loads a compiled bundle

Today only the command line (`runLGB` in `pkg/cli`) and core boot (`LoadCoreBundle`) read an `.lgb`. A running lg program cannot. #1145 asks for `(load-bundle path)`, which loads a bundle the way running it would and leaves its deferrable namespaces undecoded until something requires them. This note records the design.

## Why

An lg program that keeps a large generated namespace needs to load it fast. let-go-wasmgc's llvm compiler is the first case: it runs in lg and restores a cached runtime library at the start of every compile, 75 MB of EDN that `read-string` takes 845 ms to read. Compiled with `lg -c`, the same data loads from a bundle in 150 to 200 ms, nearly all of it in Go's decoder. The compiler has to load that bundle as a library, not run it as the program.

## The design

`load-bundle` takes a path string. It reads the file and decodes it with `bytecode.DecodeBundle`, asking to defer every namespace. The decoder decodes the main chunk at once even when asked to defer it, and decodes eagerly any namespace the bundle did not mark self-contained (`FlagNSRanges`). `lg -c` writes no namespace ranges, so every namespace in its bundles loads eagerly; deferral applies to bundles written with ranges, as lgbgen writes them. Under the unit's type scope (`runUnitChunk`), it then runs every eagerly decoded namespace chunk in `NSOrder`, then the main chunk if it wasn't one of them, applying `ApplyGoOverrides` after each, as `RunExecUnit` does. It returns the main chunk's value, as `load-string` returns its last form's. A read, decode or run error raises with the path in its message. A bundle whose debug information was split loads its companion the way `runLGB` does, and decodes eagerly: the companion attaches to chunks by index. Every chunk the feature runs, eager or deferred, runs with `*ns*` saved and restored, as `LoadCoreBundle` does around its replay.

Before any chunk runs, each namespace left deferred maps to its unit in a table in `pkg/rt`, guarded by `nsMu`, and is registered and marked needs-load, as core boot does, so an eager chunk that requires one finds it. A namespace loaded before the call is not marked: as under `require`, it stays as it is.

When a needs-load namespace is in that table, `LookupOrRegisterNS` materializes its chunk (`ExecUnit.NSChunk`), runs it, applies the overrides, drops the entry and clears the mark. This happens before `nsLoader` is asked, whether or not one is installed. A failure keeps the mark and records the error, so `RequireNS` reports it and a later `require` retries, as with `bytecodeNSLoader`. `RequireNS`'s failure check counts the table as a loader. Two goroutines that require the same deferred namespace at once can see it half-built or run its chunk twice, as with the namespace loaders today; a per-entry loading state that waiters block on is a follow-up, and needs a way to tell the loading goroutine's own re-entry from a second goroutine.

The native needs no compiler, so it works from the interpreter binary alone (`contribution-policy.md`, design contract 1).

## Approaches considered

- **The table in `pkg/rt`, consulted by namespace lookup (chosen).** It works under every `nsLoader` the CLI, `pkg/api`, wasmhost, the wasm render main and lgbgen install, and with none. The cost is a small change to the lookup path every `require` takes.
- **Wrap the current `nsLoader` at load time.** It leaves the lookup path alone, but a later `SetNSLoader` drops the bundles silently, and with no loader installed the lookup never asks the wrapper.
- **Fold core boot's `bytecodeNSLoader` into the same table.** It would unify the two, but it changes core boot for nothing #1145 needs.

## Tests

- A namespace with a `def` and a `defn`, compiled with `lg -c`: load its bundle, read the var, call the fn.
- A bundle with two self-contained namespaces where the main one does not require the other: after `load-bundle` the second is needs-load and undecoded; `require` runs it and its var reads correctly.
- A missing file raises an error naming the path. With no `nsLoader`, `require` of a deferred namespace still loads it.

## Not in scope

- `require` finding a namespace's `.lgb` on `-source-paths`: it changes the resolver's file lookup and gets its own issue.
- Decoding a namespace's entries on demand rather than whole, the random-access step #1145 names.
- Any change to core boot.
