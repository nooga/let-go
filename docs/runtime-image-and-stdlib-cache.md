## Runtime images and precompiled stdlib

This document captures the design for dumping/loading a self-contained image of code + heap for fast startup, and precompiling the standard library to cache compilation.

### Motivation

- Instant cold-start: avoid reader/analyzer/compiler on boot; skip heavy init.
- Reproducible deployments: freeze a known-good REPL state to an image.
- Embedding: ship an image inside a Go binary (via go:embed) for turnkey apps.

### Scope of a runtime image

- Snapshot includes:

  - Namespaces: registry, aliases, current `*ns*` binding.
  - Vars: symbol → root value, metadata (if any), macro flags.
  - Code: `Func`/`Closure` with `CodeChunk` (bytecode, `maxStack`, const refs) and closed-overs.
  - Constant pool: de-duplicated `Value`s (ints, strings, booleans, nil, symbols, keywords, vectors, lists, maps, sets).
  - Intern tables for symbols/keywords (to preserve identity sharing).

- Snapshot excludes (rebind or reject on load):
  - `NativeFn` and `Boxed` host pointers (Go values): represented as externs by symbol; rebound from a host registry at load.
  - Non-serializable effects: channels, OS handles, goroutines, time; reject or replace with inert placeholders.

### Image format (outline)

- Header: magic, schema version, runtime version, endianness, optional flags.
- String table: unique strings used by symbols/keywords and values.
- Atoms: typed value section (nil, booleans, ints, chars, strings).
- Collections: vectors/lists/maps/sets encoded recursively; prefer persistent encodings when available (BPTR/HAMT) for determinism.
- Const pool: indices into value sections.
- Code: functions with `maxStack`, code bytes, const pool references, closed-over const refs.
- Namespaces: table of ns name → alias map → var table.
- Vars: name → const-id (root), flags.
- Trailer: checksum/hash and optional signature.

### Host/native binding

- Loader is provided a `HostRegistry` mapping `core-symbol` → `NativeFn`.
- During load, extern placeholders (for natives/boxed values) are resolved via the registry.
- Fail fast on missing natives or signature/arity mismatch.

### Build-time and boot-time workflows

- Precompiled stdlib cache:

  - Build step: compile stdlib once; emit `stdlib.img`.
  - App boot: `LoadImage(stdlib.img, hostRegistry)`. If invalid/mismatched, compile sources and regenerate.

- Full runtime image (REPL freeze):
  - Developer REPL: load stdlib, evaluate app code until satisfied.
  - Save: `(image/save "app.img")` dumps current namespaces/vars/code.
  - Deploy: embed or ship `app.img` and load it at startup.

### Layering

- Base image: stdlib only (fast cold start universally).
- App layer: app-specific code; can be loaded atop base image without reloading stdlib.
- Optional: multiple layers chained with conflict checks.

### Versioning and invalidation

- Image carries: schema version, runtime/VM version, and a content hash (e.g., of stdlib sources).
- Loader rejects incompatible versions; for stdlib image, rebuild when stdlib hash or runtime version changes.

### Security

- Treat images as untrusted input: validate sections and sizes, check hashes, and (optionally) verify signatures before load.
- Do not deserialize host pointers blindly; only resolve allowed externs via registry.

### Risks and mitigations

- Identity and interning: rebuild symbol/keyword interns deterministically during load.
- Map/Set determinism: Go maps have non-deterministic iteration; serialize with a stable order or adopt persistent HAMT.
- Hash caches: recompute on load; do not serialize cached hash fields.

### Minimal API sketch (loader/saver)

- `func SaveImage(w io.Writer, opts ...Option) error`
- `func LoadImage(r io.Reader, host HostRegistry, opts ...Option) (*vm.Namespace, error)`
- `type HostRegistry interface { Resolve(sym vm.Symbol) (vm.Value, bool) }`
- `func TryLoadStdlib(host HostRegistry) error` (attempt stdlib image, fallback to compile)

### Relevant files and integration points

- `pkg/rt/core/core.lg`: stdlib source to precompile.
- `pkg/rt/lang.go`: current core install; call image loader first; fallback to Eval(core).
- `pkg/compiler/eval.go`: current compile path (source → bytecode); used to regenerate images.
- `pkg/compiler/reader.go`: reader; not used on image load.
- `pkg/vm/constpool.go`: const pool structure; mirror as an encoded section.
- `pkg/vm/value.go` and scalar/collection types: serialization tags per type.
- `pkg/vm/func.go`: `Func`, `Closure`, `CodeChunk` (bytecode) encoding.
- `pkg/vm/native_func.go`: extern resolution for natives via `HostRegistry`.

### Implementation checklist

- [ ] Define image schema (header, sections, types) and versioning.
- [ ] Implement value serialization/deserialization for pure values.
- [ ] Serialize const pool and code chunks; rehydrate `Func`/`Closure` with closed-overs.
- [ ] Serialize namespace/var tables; restore `*ns*`.
- [ ] Extern/native resolution via `HostRegistry`.
- [ ] Add `image save/load` APIs and CLI hooks.
- [ ] Add stdlib build step to produce/refresh `stdlib.img`; boot tries image before compiling.
- [ ] Tests: round-trip identity of values/code; load-time validation; performance of cold boot vs compile.
