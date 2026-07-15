/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"fmt"
	"time"

	"github.com/nooga/let-go/pkg/bytecode"
)

// CoreLoadOptions tunes the shared core-boot spine (LoadCoreBundle).
//
// It intentionally exposes a SINGLE behavioral axis (EagerHybrids). If a future
// caller wants a third replay behavior, that's the signal this shared spine is
// turning into the wrong abstraction — split it rather than growing more flags.
type CoreLoadOptions struct {
	// EagerHybrids replays HYBRID namespace chunks (native fns + bundled lg
	// source, e.g. async) immediately after core+baseline instead of leaving
	// them to the on-demand loader. Their vars are reachable via qualified
	// symbols WITHOUT a (require ...), which bypasses the loader, so a caller
	// that returns straight to user code (the compiler boot path,
	// BootCore) must run them eagerly or those vars stay nil stubs. A caller
	// that installs a namespace loader and only reaches hybrids through it
	// (LoadCore + bytecodeNSLoader) can leave this false.
	EagerHybrids bool

	// OnPhase, if set, is called after each boot phase with the phase name and
	// the time it started, so a caller can record fine-grained boot timing
	// (the compiler path's bootMark). Phases: "decode-bundle", "run-core-chunk".
	// nil disables timing.
	OnPhase func(phase string, since time.Time)
}

// LoadCoreBundle is the single compiler-free core-boot path: decode the
// embedded clojure.core bundle (CoreCompiledLGB) and replay it onto the live
// registry — core's main chunk (all of core's def/defn/defmacro), then the lg
// baseline namespaces (let-go.core, let-go.types) which are auto-refer'd
// everywhere but never explicitly required, then optionally the hybrid
// namespace chunks. Non-core/non-baseline namespaces are marked NeedsLoad so a
// configured loader picks them up on (require ...). Returns the decoded unit so
// a caller can read its const pool or chunk map.
//
// It folds together what were three parallel implementations (#506):
// compiler.loadPrecompiledBundle, rt.LoadCore, and rt.BootCore. Two invariants
// that used to live in only some of them are now unconditional here:
//
//   - *ns* is saved before the replay and restored after. Each chunk runs its
//     (ns …) form under a frame with no ExecContext, so in-ns falls through to
//     CurrentNS.SetRoot and mutates the global root; unrestored, a caller that
//     returns straight to user code (BootCore) would deref a *ns* left pointing
//     at whichever chunk ran last. loadPrecompiledBundle only got away without
//     this because api.NewContext re-establishes *ns* before user code.
//   - the eager hybrid loop iterates unit.NSOrder (dependency order), not the
//     NSChunks map, so replay order is deterministic, and calls
//     ReapplyGeneratedPrimitives after each chunk — the chunk's bootstrap defs
//     overwrite the native adapters #438 Def'd at init, so without the reapply
//     a bundle that eager-re-Defs one would strand its callers on the
//     trampoline and lose the direct-call fast path.
func LoadCoreBundle(opts CoreLoadOptions) (*bytecode.ExecUnit, error) {
	if len(CoreCompiledLGB) == 0 {
		return nil, fmt.Errorf("LoadCoreBundle: embedded core is empty (built -tags bootstrap?)")
	}

	// Restore the pre-boot *ns* root: the eager (ns …) forms below mutate it.
	savedNS := CurrentNS.Deref()
	defer CurrentNS.SetRoot(savedNS)

	// Namespaces present before decode are native-backed (installers ran at
	// package init). A bundle chunk for one of them is a HYBRID namespace whose
	// chunk must run eagerly (qualified refs bypass the on-demand loader).
	preexisting := map[string]bool{}
	for name := range AllNSes() {
		preexisting[name] = true
	}

	tDecode := time.Now()
	unit, err := bytecode.DecodeToExecUnitBytes(CoreCompiledLGB, LGBVarResolver)
	if err != nil {
		return nil, fmt.Errorf("decode core bundle: %w", err)
	}
	if opts.OnPhase != nil {
		opts.OnPhase("decode-bundle", tDecode)
	}

	// Replay core's def/defn/defmacro definitions.
	tCore := time.Now()
	if err := runChunk(unit.MainChunk); err != nil {
		return nil, fmt.Errorf("run core chunk: %w", err)
	}
	if opts.OnPhase != nil {
		opts.OnPhase("run-core-chunk", tCore)
	}

	if unit.NSChunks == nil {
		return unit, nil
	}

	// lg baseline namespaces depend on core's defs — run eagerly, right after.
	// A Go-only baseline (let-go.types, predicates Def'd in Go) has no chunk;
	// rt.runChunk does not nil-guard, so skip those here.
	baseline := map[string]bool{}
	for _, name := range LgBaselineNSNames() {
		baseline[name] = true
		if ch := unit.NSChunks[name]; ch != nil {
			if err := runChunk(ch); err != nil {
				return nil, fmt.Errorf("run baseline %s: %w", name, err)
			}
		}
	}

	// Everything else loads on demand.
	for name := range unit.NSChunks {
		if name != NameCoreNS && !baseline[name] {
			MarkNSNeedsLoad(name)
		}
	}

	if opts.EagerHybrids {
		for _, name := range unit.NSOrder {
			if name == NameCoreNS || baseline[name] || !preexisting[name] {
				continue
			}
			ch := unit.NSChunks[name]
			if ch == nil {
				continue
			}
			if err := runChunk(ch); err != nil {
				return nil, fmt.Errorf("run hybrid %s: %w", name, err)
			}
			ClearNSNeedsLoad(name)
			ReapplyGeneratedPrimitives(name)
		}
	}

	return unit, nil
}
