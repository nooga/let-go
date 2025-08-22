## Pods — Babashka-compatible external process integration

This document proposes implementing Babashka-compatible pods to provide an industrial-strength extension mechanism without relying on Go plugins. Pods are external processes speaking a simple EDN-based RPC over stdio (or TCP). let-go will act as a host, enabling script-friendly access to system APIs, databases, and tools while staying safe and portable.

### Goals

- Be compatible with the Babashka pod protocol so existing pods can be used from let-go.
- Provide a stable host API in let-go to load, invoke, and manage pods.
- Robust process management: timeouts, cancellation, restarts, and resource limits.
- Efficient and predictable data mapping between let-go values and EDN.

### Protocol (compatibility)

- Transport: stdio by default; optional TCP.
- Encoding: EDN messages, one per line (newline-delimited), request/response maps with an `:id` correlation key.
- Operations (as in Babashka pods):
  - `{:op :describe}` → returns namespaces, vars, arglists, and optional metadata (e.g., `:format`, `:options`).
  - `{:op :invoke :var "ns/name" :args [...] :id N}` → invokes a var with EDN-coercible args; returns `{:id N :value ...}` or `{:id N :ex ... :err <string>}`.
  - `{:op :shutdown}` → graceful shutdown.
- We will match Babashka’s shape for `:describe` and `:invoke` to the extent publicly documented, keeping room for extensions. Unknown keys are ignored.

### Host architecture in let-go

- Pod manager (`pkg/rt/pods.go`):

  - `Start(path string, opts Options) (*Pod, error)` spawns a process with pipes; sets env `BABASHKA_POD=1`.
  - Sends `:describe`, parses descriptor, and materializes a proxy namespace in let-go.
  - Maintains a goroutine reading stdout, decoding EDN messages, and dispatching to a response map keyed by `:id`.
  - Supports `Stop()` (graceful) and `Kill()` (forceful).

- Proxy namespace and vars:

  - For each described var `ns/sym`, install a `Var` in the target let-go namespace that implements `Fn`.
  - On `invoke`, serialize args to EDN, send an `:invoke` request with a unique `:id`, await response or timeout, map result back to let-go `Value`.
  - Errors are returned as let-go exceptions with pod-provided message and optional data.

- Concurrency and cancellation:
  - Each request carries a context with timeout; host can send SIGTERM on cancel and mark the call failed.
  - Optional `:interrupt` op for pods that support it; otherwise process-level cancellation.

### Data mapping (EDN ↔ let-go)

- Scalars: `nil`, booleans, ints, strings, chars, keywords, symbols map 1:1.
- Collections: vectors ↔ `PersistentVector`, lists ↔ `List`, maps ↔ `PersistentHashMap` (string/keyword/symbol keys), sets ↔ `PersistentHashSet`.
- Unsupported types: error with clear message; future extension via tagged literals.
- Performance: avoid reflection; use direct constructors and fast numeric paths.

### Security and policy

- Allowlist/denylist of pod executables and arguments.
- Resource limits: per-pod CPU timeouts, memory soft limits (OS-specific), max concurrent calls.
- Sandboxing guidance: recommend running pods in containers where needed.

### Developer ergonomics

- Host API in let-go:

  - `(pods/load path & {:as opts})` → returns an alias namespace (or installs under provided alias).
  - `(pods/unload alias)` → stops process and removes proxies.
  - `(pods/list)` → returns running pods with PIDs and descriptors.

- CLI helpers:

  - `lg pods run <path> -- <args>` → quick test runner.
  - `lg pods describe <path>` → prints descriptor.

- Library for writing pods in Go:
  - Provide a small helper in `pkg/podsdk` to implement the protocol from Go (EDN encode/decode, routing, op handlers).

### Minimal acceptance criteria (MVP)

- Load and invoke an existing Babashka pod via stdio; pass structured args; receive results and errors.
- Create proxy vars usable like normal functions: `(require '[pod.my/tool :as tool])` then `(tool/echo 1 2 3)`.
- Timeouts and graceful shutdown work.

### Stretch goals

- TCP transport and reconnect.
- Streaming operations (chunked I/O) for large payloads.
- Binary/transit encoding for throughput (configurable).

### Implementation plan

1. EDN encode/decode utilities for `vm.Value` ↔ EDN.
2. `PodManager` and `Pod` with stdio transport and request/response routing.
3. `pods/load` and proxy namespace installation; value marshaling; errors.
4. CLI tools for `describe` and `run`.
5. Concurrency: contexts, timeouts, cancellation, shutdown.
6. Tests with a simple reference pod; interop test against a known Babashka pod.
