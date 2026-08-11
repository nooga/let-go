---
status: active
last-verified: 2026-08-11
human-verified: 2026-08-11
---

# OS and filesystem operations

The built-in `os` namespace covers process execution, environment and host
information, and common filesystem work. This page calls out the contracts that
matter for correctness; straightforward getters are grouped at the end.

## Files and paths

| Form | Result and behavior |
|---|---|
| `(os/cwd)` | Current working directory. |
| `(os/temp-dir)` | Host temporary-directory path. |
| `(os/ls path)` | Vector of entry names directly beneath `path`. |
| `(os/stat path)` | `os/FileStat` with `:name`, `:size`, `:dir?`, and `:mod-time`, or `nil` when absent. |
| `(os/unzip zip-path dest-dir)` | Safely extracts regular files under `dest-dir` and returns it. Entries cannot escape the destination; symlinks and other special entries are skipped. |
| `(os/rename old-path new-path)` | Renames and returns `new-path`. On Unix a same-filesystem rename is atomic; there is no copy/delete fallback when the host cannot rename directly. Go does not guarantee atomicity on non-Unix hosts. |
| `(os/delete-tree path)` | Recursively removes `path` and returns `nil`. An absent path succeeds; an empty path, `.`, or filesystem root is refused. Symlinks are unlinked, never followed. |

Use the two path forms according to whether the named entry already exists:

| Form | Must exist? | Resolves symlinks? | Typical use |
|---|---:|---:|---|
| `(os/absolute-path path)` | No | No | Name a file that may be created later. |
| `(os/canonical-path path)` | Yes | Yes | Compare identities or construct a cache key. |

Both return an absolute, cleaned path and reject an empty string. In particular,
`canonical-path` errors on a missing entry rather than returning a lexical
guess.

A same-directory staging file plus `rename` is the usual Unix pattern for
publishing a complete file without exposing a partial write:

```clojure
(spit "result.edn.tmp" rendered)
(os/rename "result.edn.tmp" "result.edn")
```

## Processes and environment

- `(os/sh command & args)` buffers a child process and returns
  `{:exit code :out stdout :err stderr}`.
- `(os/exec* command & args)` streams through the current `*out*` and `*err*`
  bindings, keeps stdin interactive, and returns the exit code.
- `os/exec` and `os/with-stdin` expose the lower-level Go command value when a
  caller needs to configure a process before running it.
- `os/args` is the full process argv vector. Prefer `*command-line-args*` for
  just the arguments after a script name.
- `os/getenv`, `os/setenv`, and `os/exit` read configuration, update the process
  environment, and terminate with an integer status.

`os/free-port` asks the host for an unused loopback TCP port and releases it
before returning. The result is a hint, not a reservation: another process can
bind the port before the caller does.

Host information is available through `os-name`, `arch`, `user-name`, and
`hostname`, plus the `file-separator`, `path-separator`, and `line-separator`
values.

## Platform availability

These functions expose host effects, so unsupported operations return the
underlying platform error. Two reduced environments are worth calling out:

- In browser `js/wasm`, Go's host shim does not provide a working cwd or
  filesystem lookup. `absolute-path` therefore fails for relative input, while
  `canonical-path` fails for every input; subprocess and other filesystem calls
  are similarly host-limited.
- TinyGo builds register only `os/exit`, `os/getenv`, and `os/args`; `os/args`
  is currently an empty vector. Use a standard-Go build for the rest of this
  namespace.
