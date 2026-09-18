---
status: planning
last-verified: 2026-09-18
authoritative-for:
  - glplat-backend-contract
  - glplat-host-capability-seam
human-verified:
---

# glplat backend contract and host-capability seam

**Status:** planning, with a provisional landing agreed on 2026-09-18 (§5).
nooga/let-go#744 lands with the registry under `internal/` and the `glplat`
namespace marked experimental; the §2 seam and the behavior-changing rules
move to nooga/let-go#392's re-scope.

**Decision requested:** adopt §2 (bind the backend as a host capability at a
dynamic var, retire the package-level registry) and §3 (the backend-independent
contract). Both are written as concrete rules with defaults, so the useful
answer is yes, no, or an amendment to a numbered rule.

## 1. Where this came from

The 2026-09-07 direction on #392 set three conditions before either graphics
backend lands: re-scope onto the #572 `surface` host-capability model rather
than the `glplat` registry, put the font primitives behind their own opt-in
tag, and write a backend-independent contract covering frame orientation,
depth/occlusion, resize, input events and modifiers, and lifecycle.

The font tag is done: `glplat_fonts` is on #744, independent of both backend
tags, with `check-default-deps` green. This doc is the other two.

Two backends exist to write the contract from, which is the reason to write it
now rather than after a third: the GLFW/OpenGL backend on #392, and the
Ebitengine backend on #744. They were measured against each other on
2026-08-14 running an unmodified consumer (xsofy's GL frontend): identical
vertex-stream dumps, and a 0.3% screenshot pixel difference attributable to
glyph-edge rounding. Every place they *disagree* is a contract gap, and §3 is
the list.

## 2. The seam: bind the backend, don't register it

### 2.1 What the registry does today

`pkg/glplat/internal/registry` holds a package-level `var backend Backend`
with a `SetBackend` setter, installed from a tagged `init()`. A build has
exactly one backend for the life of the process, chosen at link time.

### 2.2 What #572 actually established

The `surface` capability is often described as "a pixel sink", but the part
worth reusing is the binding discipline, and it is three properties:

1. The capability is a Go interface in `pkg/rt`, so the implementing package
   is not the public extension point.
2. It resolves through a **dynamic var**, so `(binding [*surface* ...])` works
   and per-`Run` bindings from `pkg/api` stay isolated.
3. The zero value is a silent no-op, so guest code needs no build-tag branch.

The registry satisfies (1) only by accident of `internal/`, and fails (2) and
(3) outright: process-global, not per-`ExecContext`, no dynamic rebinding, and
`Get() == nil` surfaces as an error from every native rather than a no-op.

### 2.3 Proposal

Define `rt.Display` in `pkg/rt` with the methods `registry.Backend` has today,
bind it at `*display*`, and install it with `api.WithDisplay` — the same shape
as `rt.Surface` / `*surface*` / `api.WithSurface`. The GLFW and Ebitengine
backends become implementations of it. `pkg/glplat/internal/registry` goes
away, and the supported extension point becomes `pkg/rt` + `pkg/api`, which is
what the direction asks for.

This also adds no new native-registration mechanism, which nooga/let-go#531
counts six of already; glplat keeps its existing `lginterop` generation and
only changes where the interop resolves its backend.

It is the direction
[`runtime-io-host-decoupling.md`](runtime-io-host-decoupling.md) already
records: graphics is listed there as a peer capability, "a guest-named
capability bound by the host; none is special once the I/O seams set the
pattern", with the work scoped in nooga/let-go#255 — the issue `surface`
itself came out of.

### 2.4 On seam width

The same doc says a seam should be "as thin as the operation allows", and a
15-method `Display` invites the objection that it is not thin. The sentence
continues: "The interface widens only where the operation differs across
platforms." Presenting a finished RGBA frame is uniform across hosts, so
`Surface` is two methods. Driving a GPU is not: texture upload, matrix state,
triangle submission and window lifecycle each differ between GLFW and
Ebitengine, which is why §3 needs fifteen rules to pin them down. The width is
the platform difference, and it is the same width the registry has today — the
proposal moves the interface, it does not grow it.

### 2.5 What this is deliberately not

It is not routing the backends through `Surface.Present`. That interface
carries one RGBA frame per call, so adopting it literally means the guest
rasterizes in Clojure and the host blits — which deletes the GPU triangle
batching both backends exist to provide, and is the path the GL frontend was
written to get *off*. `Display` and `Surface` are peer capabilities over the
same binding discipline, not one wrapping the other. If the intent was the
literal reading, that changes the whole proposal and is worth saying plainly.

## 3. The contract

Rules are numbered for amendment. Each says what both backends do today, since
in several cases they already disagree.

### 3.1 Frame orientation and pixel format

**C1.** `LoadTextureRGBA` takes row-major pixel data, and texture coordinate
`v = 0` samples its **first row**. Both backends do this today: the GL backend
hands the buffer to `glTexImage2D` unflipped, and the Ebitengine backend maps
`v` to `SrcY` on a top-left-origin image. The rule is written down because it
is invisible until a backend gets it wrong, and then everything is upside down
at once.

**C2.** `SubmitTriangles` vertex positions are in the coordinate space
established by `SetMatrix`, column-major, with no implicit Y flip anywhere.

**C3.** Pixel data is **straight (non-premultiplied) alpha**. A backend whose
native format is premultiplied converts on upload.

The backends already diverge here in a way nothing states: the GL path is
straight alpha by virtue of `glTexImage2D`, and the Ebitengine backend
premultiplies in `LoadTextureRGBA` because `WritePixels` requires it. That
conversion is currently a comment in one backend rather than a contract, which
makes it the single most likely thing for a third backend to skip.

### 3.2 Depth and occlusion

**C4.** The contract is **painter's order**: triangles composite in submission
order within a frame, with no depth buffer and no implicit sorting.

Today the GL backend has a depth buffer and Ebitengine does not, which is why
the frontend's first-person mode mis-sorts on Ebitengine while top-down is
correct. C4 picks the weaker guarantee as the portable one.

**C5.** Depth-tested rendering, if it is ever wanted, arrives as a capability
query and an explicit mode, never as a backend that silently sorts differently
from its peers.

### 3.3 Resize

**C6.** `WindowSize` returns the size of the frame currently being built. It
is valid to call between `BeginFrame` and `EndFrame` and must not return a
stale value after the window has been resized.

This is one of the open safety findings on #392 as well as a contract item.

**C7.** Framebuffer and window size may differ (Retina and other scaled
displays). `WindowSize` reports **logical** size; a backend renders at device
scale on its own. The Ebitengine backend had exactly this bug — `Layout`
returned logical size, so it rendered at half the GL backend's resolution and
upscaled, fixed 2026-08-14.

### 3.4 Input events and modifiers

**C8.** `PollInputEvents` drains a queue and returns it; the same event is
never returned twice, and an empty return means no input, not "not supported."

**C9.** The event format is `"<kind>:<name>"`. Defined kinds are `key` and
`key-repeat`.

**C10.** Modifiers are part of the event, not dropped. This needs a format
extension — the proposal is `"key:shift+up"` with modifiers in a fixed order
(`ctrl`, `alt`, `shift`, `super`), lowercase, joined by `+`.

Today neither backend carries modifiers: the GL backend's `keyCallback`
receives GLFW's `mods` and discards it, so Shift-modified keys are
indistinguishable from unmodified ones. C10 is the one rule here that is a
behavior change to both backends and to any consumer parsing the strings.

**C11.** `key-repeat` timing is platform-defined and not part of the contract.
Consumers that need deterministic repeat implement it from `key` plus `Time`.
The two backends differ today (GLFW's delay/interval versus an approximated
30/3 ticks) and making them agree is not worth the coupling.

### 3.5 Lifecycle

**C12.** The legal call order is `Init` → (`BeginFrame` → … → `EndFrame`)* →
`Terminate`. Every other entry point is only valid between a successful `Init`
and a `Terminate`.

**C13.** After `Terminate`, backend state is dead: texture ids are invalid,
`ShouldClose` reports true, and calls that would touch a destroyed context
return an error rather than entering the backend. This is the other open
safety finding on #392, and §2's per-binding lifetime is what makes it
enforceable rather than conventional — a bound `Display` can be unbound.

**C14.** `Init` after `Terminate` (re-init) is **not supported** in this
revision. Both backends have process-global graphics state today and neither
has been tested for it. A backend may not silently half-work here; it returns
an error.

**C15.** Who owns the main thread is a backend property, not a guest concern.
The GL backend runs the guest on the main thread; Ebitengine owns it and runs
the guest on a goroutine, rendezvousing at `EndFrame`. Guests must not assume
either, which in practice means not assuming that `EndFrame` returns on the
same OS thread it was called from.

## 4. Cost

The seam change touches `pkg/rt` (the new capability), `pkg/api` (the option),
`pkg/glplat` (both backends lose their `init()` registration), the generated
`pkg/rt/interop_glplat.go` (regenerated, not hand-edited), and the one real
consumer, xsofy's GL frontend.

Of the contract rules, C10 and C13 are behavior changes with a visible
consumer effect, and C3 is one a third backend would have to implement; C4,
C6, C7 and C15 describe or fix what the backends should already do; the rest
are written-down status quo.

## 5. The decision

1. **Seam** (§2): adopt `rt.Display` bound at `*display*`, retiring the
   registry — yes or no. A no with the literal `Surface.Present` reading
   intended (§2.5) is a useful answer too, and a different piece of work.
2. **Contract** (§3): adopt C1–C15 as written, or amend by number.

If both carry, the implementation order is contract-visible behavior first
(C10 and C13, with their consumer changes), then the seam, then rebasing the
two backend PRs onto it. The remaining safety findings on #392 are unaffected
by this decision and can proceed in parallel.

### 5.1 Decision record, 2026-09-18

Answers collected on nooga/let-go#744.

- **nooga:** land #744 with the registry provisional and document it (yes);
  the `rt.Display` migration goes in #392's re-scope (agreed); no preference
  on C1–C15; release notes call the namespace experimental (agreed).
- **nnunley:** the larger design is a scene graph at the top, a renderer in
  the middle, and a surface at the bottom where frames land. 3D is a renderer
  choice, not a kind of surface. The §2 seam holds as long as a display does
  not imply a window and a GL context; if C1–C15 encode that assumption, that
  is the one amendment wanted.
- **Check against that condition:** no rule assumes a GL context (C4 chose
  painter's order so a backend without one qualifies; C15 makes thread
  ownership a backend property; the Ebitengine backend is the existence
  proof). No rule assumes a window either; C6 and C7 use "window" only as the
  name for the presentation target's logical size. The assumption lives in
  the interface: `Init(width, height, title)` takes a window title, and
  `PollEventsWindow`, `WindowSize`, and `ShouldClose` are named for one.
  C1–C5 assume a triangle-submitting renderer, which is the middle layer of
  the model above with presentation bundled in.
- **Outcome:** C1–C15 stand as written. The re-scope in #392 carries the §2
  seam with one amendment: the presentation target may be a bound
  `*surface*` rather than the backend's own window, and `Init`'s title and
  the `Window*` names generalise with it. C3, C4, and C10 are implemented
  there as well.
