/*
 * Copyright (c) 2026 let-go contributors; see CONTRIBUTORS.
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"github.com/nooga/let-go/pkg/bytecode"
	"github.com/nooga/let-go/pkg/vm"
)

// bundleTable maps each namespace a loaded bundle left undecoded to the unit
// that holds it, until something requires the namespace. Guarded by nsMu.
var bundleTable = map[string]*bytecode.ExecUnit{}

// RegisterBundleNamespaces records names as deferred namespaces of unit:
// each is registered and marked needs-load, so its first lookup materializes
// and runs its chunk, whatever namespace loader is installed or none. The
// caller leaves out namespaces that were loaded before the bundle was
// decoded (decoding registers every namespace it references, so the
// registry cannot tell them apart afterwards).
func RegisterBundleNamespaces(unit *bytecode.ExecUnit, names []string) {
	for _, name := range names {
		name = resolveNSAlias(name)
		DefNSBare(name)
		nsMu.Lock()
		bundleTable[name] = unit
		nsNeedsLoad[name] = true
		nsMu.Unlock()
	}
}

func bundleTableHas(name string) bool {
	nsMu.RLock()
	defer nsMu.RUnlock()
	_, ok := bundleTable[name]
	return ok
}

// loadFromBundleTable materializes and runs name's chunk from the bundle that
// deferred it. ok is false when no bundle holds name. A failure keeps the
// entry and the needs-load mark so a later require retries.
func loadFromBundleTable(name string) (ns *vm.Namespace, ok bool, err error) {
	nsMu.RLock()
	unit, ok := bundleTable[name]
	nsMu.RUnlock()
	if !ok {
		return nil, false, nil
	}
	// The chunk's ns form looks the namespace up again; clearing the mark
	// first keeps that lookup from loading it a second time.
	nsMu.Lock()
	delete(nsNeedsLoad, name)
	nsMu.Unlock()
	chunk, err := unit.NSChunk(name)
	if err == nil && chunk != nil {
		err = runUnitChunkKeepingNS(unit, chunk)
	}
	if err != nil {
		nsMu.Lock()
		nsNeedsLoad[name] = true
		nsMu.Unlock()
		return nil, true, err
	}
	ns = LookupNS(name)
	if ns != nil {
		ApplyGoOverrides(ns)
	}
	nsMu.Lock()
	delete(bundleTable, name)
	nsMu.Unlock()
	return ns, true, nil
}

// runUnitChunkKeepingNS runs a chunk under its unit's type scope and puts
// *ns* back afterwards, since the chunk's ns form moves it.
func runUnitChunkKeepingNS(unit *bytecode.ExecUnit, c *vm.CodeChunk) error {
	saved := CurrentNS.Deref()
	defer CurrentNS.SetRoot(saved)
	return runUnitChunk(unit, c)
}
