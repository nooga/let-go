---
status: active
last-verified: 2026-08-14
authoritative-for:
  - terminal-resize-wake
human-verified: 2026-08-14
---

# Terminal resize wake contract

**Status:** active design for making Plan 9 terminal resizes wake a blocked
`term/read-key` while preserving the existing native contract.

**Decision:** terminal backends may interrupt a blocked key read with the
synthetic BEL key (`\x07`). The notification source is platform-specific, but
the observable `read-key` / `key-pending?` behavior is shared.

## 1. Problem

Interactive applications commonly block in `term/read-key` and check
`term/size` once per event-loop iteration. Native Unix already wakes that read
on `SIGWINCH`, but Plan 9's queued stdin source previously woke only for input
or EOF. Plan 9 terminals could update `/env/COLS` and `/env/LINES` immediately
while a Letgo application remained parked until the next keypress.

This is a runtime terminal concern, not an application-specific workaround:
any Letgo TUI needs a chance to observe new geometry without waiting for user
input.

## 2. Shared contract

The native contract introduced in #165 remains authoritative:

- a terminal resize may make `key-pending?` true;
- `read-key` consumes the wake as synthetic BEL (`\x07`);
- real keyboard input takes priority when input and resize are both pending;
- resize storms may coalesce into one wake; and
- EOF and reader errors take priority over a synthetic wake.

The application does not interpret BEL as geometry. It handles the event as it
normally would, then reads `term/size`; a size change triggers its ordinary
redraw path.

## 3. Platform mechanisms

### Native Unix

The existing backend receives `SIGWINCH`, writes a BEL byte to a self-pipe, and
polls that pipe alongside stdin. No behavior changes on Linux or macOS.

### Plan 9

The Plan 9 root key source wraps `queuedKeySource` and lazily watches the live
`/env/WINCH` generation. Terminals publish `COLS` and `LINES` before advancing
`WINCH`, so a changed generation is the commit signal that new dimensions are
ready. The watcher then queues an internal wake on the source.

The watcher uses a 100 ms ticker. It is event polling, not a busy loop: the
goroutine sleeps between reads and reads one small environment file per tick.
Notifications coalesce in the queued source, so resize drags cannot create an
unbounded backlog.

If `WINCH` is absent, the watcher quietly waits for it to appear. Input, EOF,
and `term/size` retain their previous behavior.

## 4. API boundary

`KeySource` remains the public two-method interface: `ReadKey` and
`KeyPending`. Adding a public `Wake` method would break existing embedders and
would mix this focused terminal fix with the separate source-lifecycle work in
#498.

Instead, `queuedKeySource` owns an unexported, coalesced wake flag and method.
The Plan 9 adapter can use that implementation detail without exposing a new
host API. WASM's SharedArrayBuffer wake protocol remains separate future work.

## 5. Lifecycle

The Plan 9 source remains the single process-lifetime root over `os.Stdin`, the
same lifecycle accepted in #461. It captures the current `WINCH` generation
when the term namespace constructs the root, before application code performs
its first `term/size` read. Its reader and resize watcher still start lazily on
the first `ReadKey` or `KeyPending` call. This split closes a startup race: a
host correction arriving between the first frame and first key read differs
from the construction-time baseline and therefore wakes that read. If an
embedder replaces `*keys*`, the root source is never consulted and neither
goroutine starts.

This does not attempt to solve repeated embedder-source teardown. That broader
ownership and `Close` question is tracked in #498.

## 6. Alternatives considered

- **Add `Wake` to `KeySource`.** Rejected for this change because it is a
  breaking public-interface expansion and is unnecessary for the Plan 9 root.
- **Inject BEL into stdin.** Rejected because concurrent byte injection can
  interleave with escape or UTF-8 sequences. Wake state belongs beside, not
  inside, the byte queue.
- **Start a new blocking read goroutine per call and select on resize.**
  Rejected because the losing read cannot be cancelled and leaks goroutines.
- **Use Plan 9 interrupt notes through `/dev/consctl winchon`.** Useful for
  terminals that expose that protocol, but not universal across Drawterm and
  redirected console namespaces. `/env/WINCH` is the portable host-published
  contract used here; note-driven wake can be a later optimization.
- **Teach the application to poll size.** Rejected because every TUI would
  need its own timer and wake loop even though the blocking primitive lives in
  the runtime.

## 7. Verification

The platform-neutral queued-source tests prove that:

- a wake interrupts a parked `ReadKey`;
- `KeyPending` observes it;
- wake bursts coalesce;
- real input is returned before a simultaneous wake; and
- the wake remains pending after that real input is consumed.

A Plan 9-only test also moves the `WINCH` generation after root-source
construction but before the first `ReadKey`, proving that the startup
correction is not lost.

CI runs the ordinary runtime suite and race detector on supported native
hosts, while a Plan 9 cross-build covers the build-tagged adapter. Interactive
acceptance testing resizes a Drawterm-backed terminal while `read-key` is
blocked and confirms that a Letgo TUI redraws without a keypress.
