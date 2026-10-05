---
status: active
last-verified: 2026-09-19
authoritative-for:
  - structured-concurrency
  - scope-lifecycle
---

# Structured concurrency: scopes

A **scope** is a node in a supervision tree: it owns a cancellation context
derived from its parent's, and it knows how many goroutines are live beneath it.
Cancelling a scope cancels its whole subtree and leaves its siblings running.

**Not every spawn joins the current scope.** `go` blocks and futures do — they
spawn through the calling execution's scope. The `async` namespace's own
helpers do **not**: `pipe`, `map`, `reduce`, `split`, `pub` and `timeout` spawn
into the process-wide root scope regardless of where they are called, because
they go through `vm.Goroutines` directly rather than through the caller's
context. Measured inside a `with-scope`:

```clojure
(with-scope [s] (go (<! (chan)))            (sleep 60) (scope-live s))  ;=> 1
(with-scope [s] (future (sleep 5000))       (sleep 60) (scope-live s))  ;=> 1
(with-scope [s] (async/pipe (chan) (chan))  (sleep 60) (scope-live s))  ;=> 0
(with-scope [s] (async/timeout 5000)        (sleep 60) (scope-live s))  ;=> 0
```

So closing a scope neither cancels nor drains an `async/pipe` started inside it.
That is a gap rather than a design choice — see [Known gaps](#known-gaps).

A scope supervises *lifetime* — cancelling and awaiting. It does not supervise
*failure*: a goroutine that throws is not reported to its scope. See
[What happens when spawned work throws](#what-happens-when-spawned-work-throws).

This exists because Go gives you no way to enumerate goroutines and no way to
kill one. Without a scope, work started inside a function outlives the call,
and nothing can wait for it or stop it. The nearest relatives are Erlang's
supervision trees, Trio's nurseries and Java's `StructuredTaskScope`.

There is one process-wide **root scope**, and everything runs in it unless you
open a child. Opening children is opt-in.

## `with-scope`

```clojure
(with-scope [s]
  (go (worker-one))
  (go (worker-two)))
;; both goroutines are cancelled and drained before this form returns
```

`with-scope` opens a child of the current scope, binds it, runs the body, and
then — in a `finally`, so it happens on a normal return, a `throw`, or an early
exit — cancels the child and waits for it to drain.

The body's return value is the form's value. The scope handle is an ordinary
value you can pass around, though it is only useful while the scope is open.

```clojure
(with-scope [s]
  (go (<! (chan)))
  (sleep 50)
  (scope-live s))        ;=> 1
```

## Cancellation is cooperative, and it is not an abort

This is the part worth reading twice. Go cannot force a goroutine to stop, so
cancelling a scope does not kill anything. It does two things:

1. Blocking runtime operations — `sleep`, channel take and put, `alts!` —
   **return early** instead of staying parked. `alts!` belongs here rather than
   with the `async` helpers above: it spawns nothing, selecting synchronously
   over its channels *and* the current scope's `Done()`, so an `alts!` parked
   inside a scoped `go` block is released when that scope closes and yields
   `nil`.
2. `scope-cancelled?` starts answering `true`.

A cancelled operation returns a value and execution continues with the next
form. A take that was parked returns `nil`:

```clojure
(def got (atom :none))
(with-scope [s]
  (go (reset! got [:take-returned (<! (chan))])))
(sleep 50)
@got                     ;=> [:take-returned nil]
```

The take never received anything; the scope closed underneath it, the take
returned `nil`, and the `reset!` after it still ran. Likewise a `(sleep 10000)`
inside a closing scope returns immediately and the following form executes.

**The footgun.** A loop that sleeps but never asks whether it was cancelled does
not stop on teardown — it runs to completion at full speed, because every
`sleep` now returns instantly:

```clojure
(def n (atom 0))
(with-scope [s]
  (go (loop [i 0]
        (when (< i 100)
          (swap! n inc)
          (sleep 20)
          (recur (inc i)))))
  (sleep 50))
@n                       ;=> 100   ;; all of them, not the ~2 it had time for
```

Add the check and it stops:

```clojure
(def m (atom 0))
(with-scope [s]
  (go (loop [i 0]
        (when (and (< i 100) (not (scope-cancelled?)))
          (swap! m inc)
          (sleep 20)
          (recur (inc i)))))
  (sleep 50))
@m                       ;=> 3
```

So: **any loop whose exit condition is time or iteration count should also test
`scope-cancelled?`.** A goroutine that only parks on a channel or a single sleep
needs nothing — it unparks on its own.

## Draining and the timeout

Closing waits for the subtree to drain, bounded by `*scope-drain-timeout-ms*`
(dynamic, default `5000`). If work is still live when the timeout expires,
`with-scope` prints a warning naming the number of stragglers and returns
anyway; it never blocks forever.

```clojure
(binding [*scope-drain-timeout-ms* 200]
  (with-scope [s]
    (go (uncancellable-work))))
```

Lower it when a caller must not be held up; raise it when shutdown work is
legitimately slow. A straggler is a bug in the goroutine, not in the scope —
it means work that neither parks on a cancellable operation nor polls
`scope-cancelled?`.

Put cleanup in `finally`, not in a `catch`. It runs on the cancellation path —
a parked take returns, execution continues, and the `finally` fires — so it
needs no handler. One hazard to know: a `finally` (or a `with-open` release)
that itself throws *replaces* the exception that was in flight rather than
being attached to it, so a failing cleanup can hide the failure it was cleaning
up after ([#919]). Keep release code from throwing where you can.

## Nesting

Scopes nest, and each level is independent:

```clojure
(with-scope [outer]
  (go (outer-work))
  (with-scope [inner]
    (go (inner-work)))     ;; cancelled and drained here
  (still-running? outer))  ;; outer's goroutine is untouched
```

Closing `inner` does not disturb `outer`. Cancelling `outer` cancels `inner`
too, because `inner`'s context is derived from it. Sibling scopes never affect
each other.

## API

| Form | Meaning |
|---|---|
| `(with-scope [s] body…)` | Open a child scope, run `body`, then cancel and drain it. |
| `(scope-open)` | Open a child of the current scope and install it. Prefer `with-scope`. |
| `(scope-close! s timeout-ms)` | Cancel, drain up to `timeout-ms`, restore the previous scope. |
| `(scope-cancelled?)` | Has the *current* scope been cancelled? |
| `(scope-cancelled? s)` | Has scope `s` been cancelled? |
| `(scope-live s)` | Live goroutine count for `s` and its descendants. |
| `(scope? x)` | Is `x` a scope handle? |
| `*scope-drain-timeout-ms*` | Drain bound used by `with-scope`. Default `5000`. |

`scope-open` and `scope-close!` are what `with-scope` expands to. Call them
directly only when the scope's lifetime cannot be a lexical block; you are then
responsible for closing it on every exit path.

`scope-cancelled?` is the only way let-go code can observe cancellation.
Blocking natives return early and silently, so a coordinator parked on `sleep`
cannot otherwise distinguish "woke up" from "was cancelled".

## Embedding in Go

The same tree is available from Go, and it is where the semantics are defined
(`pkg/vm/scope.go`):

```go
s := vm.Goroutines.Child()        // sub-scope of the process-wide root
s.Go(func(ctx context.Context) {  // tracked goroutine; select on ctx.Done()
    …
})
s.Shutdown(5 * time.Second)       // cancel subtree, wait for drain
```

- `Child()` — sub-scope whose context derives from this one's.
- `Go(fn)` — run `fn` in a tracked goroutine; it receives the scope's context.
- `Cancel()` — cancel this scope and its subtree; terminal, does not wait.
- `CancelAll()` — cancel the current generation and install a fresh one, so the
  scope keeps accepting work. The "drain and continue" form.
- `Await(timeout)` — block until the subtree is empty; `true` if drained. A
  non-positive timeout waits forever.
- `Shutdown(timeout)` — `Cancel` then `Await`. Use for a sub-scope.
- `Drain(timeout)` — `CancelAll` then `Await`. Use for the root.
- `Live()` / `LiveTree()` — direct, and whole-subtree, live counts.
- `Context()` — the cancellation context; blocking operations select on its
  `Done()`. One atomic load, cheap enough to call per channel operation.

A `*vm.Scope` is a `vm.Value`, so it can be handed to let-go as an opaque
handle; it prints as `#<scope live=N>`.

Spawning touches only its own scope's bookkeeping — there is no global counter
for every goroutine to contend on — and the only mutex guards the child list,
which is touched when a scope is created rather than on every spawn.

## How a goroutine finds its scope

A scope travels on the **execution context**, not on the goroutine. `ExecContext`
carries both the current scope and the dynamic-variable binding stack, and a
spawned goroutine gets a child context that inherits both. `ec.Scope()` resolves
to the root when nothing more specific is installed.

This is why `bound-fn` keeps working across a spawn, and why there is no
goroutine-ID registry to consult. The design argument is in
[`design/exec-context-threading.md`](../design/exec-context-threading.md).

## What happens when spawned work throws

A scope supervises *lifetime* — it can cancel work and wait for it. It does not
supervise *failure*. Nothing about a throw inside spawned work reaches the
scope: the scope is not notified, its live count is unaffected once the
goroutine exits, and `with-scope` returns normally.

Where the error goes depends on which form spawned the work, and the three
differ:

| Spawned by | On a throw | What the caller can observe |
|---|---|---|
| `(agent …)` + `send` | `(agent-error a)` is set, but the agent does **not** stop | The error is readable, and the agent keeps accepting actions — a later `(send a inc)` still runs and updates the value. Clojure refuses further actions until `restart-agent`. |
| `(go …)` | `ExecutionError: …` printed to **stderr** | The result channel yields `nil` — indistinguishable from a block that returned `nil`. |
| `(future …)` | **Nothing is printed** | `@f` returns `nil`, now and on every later deref. The error is lost ([#805]). |

```clojure
(def a (agent 0))
(send a (fn [_] (throw (ex-info "boom-agent" {}))))
(agent-error a)          ;=> #error {:message "boom-agent", :data {}}

(def c (go (throw (ex-info "boom-go" {}))))
(<!! c)                  ;=> nil        ;; and "boom-go" on stderr

(def f (future (throw (ex-info "boom-future" {}))))
@f                       ;=> nil        ;; silently
```

The agent path is the only one where the error is readable in code — though the
agent keeps accepting actions afterwards, so `agent-error` tells you something
failed, not that the agent stopped. If a `go` block or a future can fail and the
caller needs to know, return the failure as a value rather than throwing:

```clojure
(go (try [:ok (risky)] (catch Throwable e [:error (ex-message e)])))
```

## Known gaps

These are gaps rather than deliberate boundaries, and they are worth knowing
before you rely on a scope for anything load-bearing.

- **A failed `future` loses its error** ([#805]). Clojure rethrows on deref;
  here `@f` is `nil` and nothing is printed, so a future that throws is
  indistinguishable from one that returned `nil`. A `go` block at least prints
  to stderr, but offers no programmatic handle either.
- **A failed agent keeps running.** `agent-error` reports the error, but the
  agent does not enter a failed state: a later `send` still runs and updates the
  value, where Clojure refuses actions until `restart-agent`. `agent` and `send`
  are synchronous, atom-backed placeholders rather than spawned workers.
- **The `async` helpers ignore the current scope.** `pipe`, `map`, `reduce`,
  `split`, `pub` and `timeout` spawn into the root scope, so a `with-scope`
  around them neither bounds nor drains their workers. `go` and `future` behave
  as documented above. `alts!` is not affected — it spawns nothing, and is covered under cancellation above.
- **No exit signals, and no restarts** ([#921]). A goroutine's failure is not
  delivered anywhere a supervisor could act on it: no equivalent of Erlang's
  exit signals or linking, no per-scope handler when a child dies, and no
  restart strategies. Siblings are not told.
- **Cancellation is not observable as a failure** ([#920]). A cancelled
  operation returns `nil` rather than raising, so it cannot be distinguished
  from a legitimate `nil` except through `scope-cancelled?`. The cancel arms
  discard `ctx.Err()` even though the natives have an error slot for it.
- **Cancellation cannot be forced.** A goroutine that neither parks on a
  cancellable operation nor polls `scope-cancelled?` runs to completion. The
  drain timeout bounds how long anyone *waits*, not how long it *runs*.
- **Everything defaults to the root.** `go` blocks and futures spawn into
  whatever scope is current, which is the root unless you opened one, and the
  `async` helpers use the root even when you did. Nothing is automatically
  bounded by the function that started it.

## See also

- `pkg/vm/scope.go` — the implementation, and the semantics of each operation.
- `pkg/vm/scope_test.go` — subtree cancel with live siblings, parent cascade,
  drain reinstating a context.
- `test/e2e/scope_test.go`, `test/e2e/scope_cancelled_test.go`,
  `test/e2e/bound_fn_scope_test.go` — end-to-end teardown behaviour.
- [`design/exec-context-threading.md`](../design/exec-context-threading.md) —
  why scope and bindings share one context.

[#805]: https://github.com/nooga/let-go/issues/805
[#918]: https://github.com/nooga/let-go/issues/918
[#919]: https://github.com/nooga/let-go/issues/919
[#920]: https://github.com/nooga/let-go/issues/920
[#921]: https://github.com/nooga/let-go/issues/921
