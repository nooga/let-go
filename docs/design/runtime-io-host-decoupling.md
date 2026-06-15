---
status: planning
last-verified: 2026-06-15
human-verified: 2026-06-15
---

# Decoupling runtime I/O from the host

**Status:** design proposal — where the runtime-I/O work is going and why
**Decision leaning:** the runtime exposes a small set of I/O *seams*; each host (native, embedder, WASM) binds the concrete stream. Output is done; input is the live frontier.

## 1. Where this came from

This direction comes from nooga/xsofy#22:

> The xterm.js bundle was sort of a temporary hack to get xsofy going. My thinking now is that we should get it out of let-go itself and provide a more decent browser interop surface instead. This should allow the client code, xsofy in this case, to build its own web shell without disturbing let-go.

This doc is what that cashes out to for I/O specifically: every place the runtime reaches for a concrete sink (`os.Stdout`, a `syscall/js` global, a JS `fs` shim) becomes a seam the host binds. The motivating consumer is xsofy, a terminal-style game running in the browser via `lg -w`, where the coupling is most acute: a native binary can hardcode `os.Stdout`, but a browser app can't, and each shortcut the runtime takes becomes a shim the client maintains downstream.

## 2. The problem

There are four trace paths for "write a string from `.lg` code," and they diverge at `os.Stdout`:

- **native `lg`** → `os.Stdout` → terminal
- **`lg -w` bundle** → `os.Stdout` → `fs.writeSync` shim → `LetGoHost._output` → the page
- **`wasm/` track** (`pkg/api`-based) → `os.Stdout` → `wasm_exec.js` default → DevTools console
- **`lg -b` native** (third-party CLIs like lgcr) → `os.Stdout` → terminal

They diverge because *the runtime doesn't know which host is using it*. The shim layer is where each host inserts itself, and that layer is duplicated, partial, or absent depending on the build. The same shape repeats for input (SAB key reads in WASM vs `os.Stdin` natively) and for emit (`js/emit` → `_lgEmit`).

## 3. The model

The runtime already had the right abstraction half-built: `*in*` / `*out*` / `*err*` are dynamic vars holding an `IOHandle`, which wraps an `io.Writer` and/or `io.Reader`. The work is finishing it: making the runtime's own ops *consult* those vars instead of going around them, then letting each host bind the concrete stream.

One contract, three bindings:

- **native** → the handle wraps `os.Stdout` / `os.Stdin`
- **embedder** (Go host via `pkg/api`) → constructor options push a per-`Run` binding
- **WASM** → a `HostWriter` / `HostReader` adapter forwards to the JS host

Two rules keep the boundary honest:

- **Host-bound capabilities live in `pkg/rt`, never `pkg/vm`.** The guest *names* the capability (`println`, `read-key`, `emit`); the host *supplies* the backend. The VM core stays free of it — and today it is (no `syscall/js` in `pkg/vm`).
- **The seam is as thin as the operation allows.** Output is a byte stream, so it's a plain `io.Writer`; no richer interface earns its keep. Input isn't (§5), so it gets a slightly larger seam. The interface widens only where the operation differs across platforms.

## 4. Status: output is decoupled

| Step | What it did | State |
|---|---|---|
| print → `*out*` | print fns consult the dynamic binding, not `os.Stdout` | #206 (merged) |
| embedder options | `WithStdout` / `WithStderr` as per-`Run` bindings | #207 (merged) |
| `term/*` → `*out*` | ANSI control ops join the same path | #223 (merged) |
| WASM host writer | `*out*` / `*err*` bound to a `HostWriter`; `fs` output shim retired | #231 (open) |

After #231, output (plain text and ANSI alike) reaches the host through one configurable sink per stream, and the browser bundle stops intercepting file descriptors. The relevant input primitives are already upstream too: a non-blocking peek (`key-pending?`, #118), the native wake-a-blocked-`read-key` mechanism (BEL self-pipe on resize, #165), and the browser-interop surface itself (`js/emit` + `LetGoHost`, #174). The direction builds on accepted ground.

## 5. Design decisions

For each: the proposal or decision, one alternative, and the tradeoff.

### 5.1 Output mechanism: root-binding primitive vs the public option

- **Decided (#231):** the WASM bundle installs the host writer via the low-level `NewWriterHandle` + `SetRoot(*out*)` primitive, not the public `WithStdout` option.
- **Alternative:** route the generated bundle main through `pkg/api` so `WithStdout` is the literal seam.
- **Tradeoff:** the bundle main runs precompiled bytecode and never touches `pkg/api`, so the public option can't reach it without teaching `pkg/api` to run a loaded unit. The option is the embedder-facing API; the primitive is the shared mechanism underneath. Reconciling them (so `WithStdout` is the one browser seam) is a reasonable later cleanup.

### 5.2 Input shape: capability op + `HostReader` vs `*in*` as `io.Reader`

- **Proposal:** keep `read-key` / `key-pending?` as `term/*` capability ops, backed by a host-supplied `HostReader` (plus a wake primitive). Treat `*in*`-as-`io.Reader` as a separate, later seam for generic stdin streaming.
- **Alternative:** make `*in*` an `io.Reader` and model `read-key` as a read off it, the symmetric dual of `*out*`.
- **Tradeoff:** there are two distinct input needs. Interactive keys are event-shaped (one keystroke = one unit), blocking, and need interruption; `io.Reader.Read` fits them awkwardly, and the siblings (`key-pending?`, `size`) aren't reads at all. Generic stdin *is* a byte stream and `io.Reader` fits it. Forcing both through one type buys symmetry at the cost of a leaky abstraction; two seams keep each focused. The interactive seam is the one the browser game needs.

### 5.3 Blocking + wake model

- **Proposal:** generalize the wake mechanism that already exists rather than invent one. Native interrupts a blocked `read-key` with a self-pipe wake-byte (returns `BEL` on resize, #165); WASM polls via `Atomics` + `key-pending?`. Make "a blocked read is interruptible by resize/EOF/stop" an explicit part of the `HostReader` seam.
- **Alternative:** non-blocking-only reads plus a callback/event model, pushing the blocking up into guest code.
- **Tradeoff:** a blocking `read-key` is what TUI/game loops want, and both platforms already solve the interrupt, so codifying it is cheaper than rewriting. The callback model inverts control and breaks the `(read-key)`-returns-a-key ergonomics `.lg` code expects. This is the gating design question for the input PR.

### 5.4 ANSI plain-vs-color: gate at emission vs strip downstream

- **Proposal:** plain-vs-color is a producer-side decision: the ops consult an `*ansi?*`-style flag and don't emit color when it's off (the `--color=auto` model). A decorating `io.Writer` that parses escapes back out is for *shaping* (truecolor → 256, strip-for-capture, log files), not the common on/off case.
- **Alternative:** always emit ANSI and strip it downstream with a decorating writer.
- **Tradeoff:** gating at emission means plain output never pays to generate-then-parse-out escapes. The downstream stripper is the right tool only when transforming a stream you already have. Cost: the ops have to honor the flag (they don't fully, yet; the near-term enhancement). Both compose at the same `*out*` seam, so the choice is in the default path, not the architecture.

### 5.5 Mechanism vs selector for the default

- **Proposal:** keep the *mechanism* (a wrapped writer, or the emission gate) separate from the *selector* (who decides plain vs color). The default can be automatic: the host checks `isatty`/tty at startup and sets the flag, so non-interactive output defaults to plain with no toggling and no always-on stripping.
- **Alternative:** a fixed default in the carrier (e.g. strip unless told otherwise).
- **Tradeoff:** separating them lets detection do the choosing once, and the useful zero-value stays pass-through. A baked-in default is simpler to state but wrong somewhere (a pipe that wants color, a TTY that wants plain) and harder to override.

### 5.6 Terminal seam: minimal vs a rich interface

- **Proposal:** keep the seam minimal — output is an `io.Writer`; input is a small capability set (`read-key`, `key-pending?`, `size`, raw-mode) behind `HostReader`. Per-platform input differences ride build tags, not a fat interface.
- **Alternative:** model the whole terminal as one rich interface (cells, styles, cursor, events as methods), the shape a full TUI library exposes.
- **Tradeoff:** ratatui is the reference for the rich shape, and its design supports keeping the seam minimal. Its `Backend` trait is output-only (`draw(cells)`, `flush`, `clear`, cursor, `size`, scroll) and handles no input; events are the application's job, pulled from the terminal library directly. And every backend is generic over a writer: `CrosstermBackend<W: Write>` wraps a `Write` and emits escape sequences to it. So a rich terminal interface isn't an alternative to the writer seam; it's a convenience layer that sits on top of one. The move here is to keep the seam at the `io.Writer` and let a ratatui-style rendering layer (buffer + diff) be built above it later, outside the runtime. Borrow the vocabulary; the runtime doesn't host the weight.

### 5.7 Where the adapters live

- **Decided:** the JS-aware adapters (`HostWriter`, `HostReader`, the emit bridge) stay in `pkg/rt` as thin, build-tagged (`js,wasm`) files; `HostWriter` already lives there as of #231. The runtime ships the adapters; the seam is the `io.Writer`/`io.Reader` boundary.
- **Alternative:** push all JS-aware code out of `pkg/rt` into the bundle/host layer, leaving `pkg/rt` with only the seams.
- **Tradeoff:** a thin build-tagged adapter in `pkg/rt` keeps the capability ops next to their backends. Fully evicting them is purer but a bigger refactor for little gain: the seam is already the `io` boundary, which is what matters.

## 6. Where this is going

1. **Input** (`HostReader`) — the output dual, and the harder half (§5.2, §5.3). The live frontier.
2. **Emit** (`HostEmitter`) — promote the existing `_lgEmit` bridge to a typed host capability, the same shape as the writer.
3. **Peer capabilities** — graphics (sixel/canvas rides `*out*` or a sibling seam), audio, controller input. Each is a guest-named capability bound by the host; none is special once the I/O seams set the pattern.

The endgame: a runtime that boots, runs bytecode, and talks to its host through a handful of bound seams, with the `lg -w` generator wiring those seams instead of installing shims.

## 7. Risks & open questions

- **The wake API.** A `wake()`/`interrupt()` method, a sentinel return (native's `BEL`-on-resize), or a context/channel? This gates the input PR.
- **`read-key` granularity.** One key per call (current), or a framed byte stream? Decides whether `io.Reader` is the right Go type for input at all.
- **Key transport.** The current WASM handoff is a single SharedArrayBuffer slot (one key in flight, 16-byte cap). A ring-buffer prototype gave noticeably more responsive input in browser builds; whether to fold that into the input PR or keep it separate (the way flush stayed per-platform) is open.
- **How much terminal-interface vocabulary to borrow** before the input seam sets, so we don't reinvent a worse version of a solved interface.
- **`vm`-level diagnostics in WASM.** A few `pkg/vm` last-resort writes (panic-recover, core-shadow, GLS-drain) can't import `pkg/rt`, so they bypass `*err*` and now land on the browser console rather than the terminal. Acceptable (a panic shouldn't scribble over the app), but noted.

## Glossary

- **Runtime:** the bytecode VM (`pkg/vm`) plus the core namespace and ops (`pkg/rt`). The decoupling goal is that it knows nothing about *where* it runs.
- **Host:** whatever embeds and drives the runtime — native CLI, a Go program using `pkg/api`, or the WASM bundle. Owns the environment (terminal, files, DOM) and supplies the concrete I/O.
- **Sink / source:** the concrete destination/origin behind a stream (`os.Stdout`, a buffer, an `xterm.js` terminal). `*out*` names the stream; the sink is what's behind it.
- **Shim:** glue the host inserts to patch around an assumption the runtime baked in (e.g. the `fs.writeSync` interceptor). The decoupling exists to delete shims in favor of bound seams.
- **Seam:** a defined boundary where the host plugs in — a dynamic var over an `io.Writer`/`io.Reader`, or a small capability interface.
- **Binding:** connecting a seam to a concrete sink/source, settled at construction; guest code can re-bind per Clojure's `(binding [*out* …] …)`.
- **Bundle:** the self-contained WASM app `lg -w` produces (generated `main` + `lg-host.js` + HTML).
