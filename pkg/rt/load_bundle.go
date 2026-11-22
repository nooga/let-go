/*
 * Copyright (c) 2026 let-go contributors; see CONTRIBUTORS.
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"fmt"
	"os"

	"github.com/nooga/let-go/pkg/bytecode"
	"github.com/nooga/let-go/pkg/vm"
)

// LoadBundle loads the compiled bundle at path into the running program, as
// running it would: it runs the bundle's namespaces in order and then its
// main chunk, and returns the main chunk's value. A namespace the bundle
// marks self-contained is left undecoded until something requires it, and a
// namespace loaded before the call stays as it is. *ns* is unchanged
// afterwards.
func LoadBundle(path string) (vm.Value, error) {
	v, err := loadBundle(path)
	if err != nil {
		return vm.NIL, fmt.Errorf("load-bundle %s: %w", path, err)
	}
	return v, nil
}

func loadBundle(path string) (vm.Value, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	loaded := loadedNamespaces()
	unit, err := decodeBundleFile(data, path)
	if err != nil {
		return nil, err
	}

	// Deferred namespaces go in the table before anything runs, so an eager
	// chunk that requires one finds it.
	var deferred []string
	for _, name := range unit.NSOrder {
		if unit.DeferredNS(name) && !loaded[resolveNSAlias(name)] {
			deferred = append(deferred, name)
		}
	}
	RegisterBundleNamespaces(unit, deferred)

	saved := CurrentNS.Deref()
	defer CurrentNS.SetRoot(saved)
	mainRan := false
	for _, name := range unit.NSOrder {
		chunk := unit.NSChunks[name]
		if chunk == nil {
			continue
		}
		if err := runUnitChunk(unit, chunk); err != nil {
			return nil, fmt.Errorf("loading namespace %s: %w", name, err)
		}
		if ns := LookupNS(name); ns != nil {
			ApplyGoOverrides(ns)
		}
		if chunk == unit.MainChunk {
			mainRan = true
		}
	}
	if mainRan || unit.MainChunk == nil {
		return vm.NIL, nil
	}
	v, err := runUnitChunkValue(unit, unit.MainChunk)
	if err != nil {
		return nil, err
	}
	// A single compilation has no namespace table; its ns form left *ns* on
	// the namespace it defines.
	if ns, ok := CurrentNS.Deref().(*vm.Namespace); ok {
		ApplyGoOverrides(ns)
	}
	return v, nil
}

// decodeBundleFile decodes data asking to defer every namespace. A bundle
// whose debug information was split decodes eagerly with its companion: the
// companion attaches to chunks by index, which needs them all decoded.
func decodeBundleFile(data []byte, path string) (*bytecode.ExecUnit, error) {
	if bytecode.HasSplitDebug(data) {
		return DecodeExecUnitWithDebugFile(data, path)
	}
	return decodeProgram(func() (*bytecode.ExecUnit, error) {
		return bytecode.DecodeBundle(data, LGBVarResolver, func(string) bool { return true })
	})
}

// loadedNamespaces returns the registered namespaces that are not waiting to
// be loaded.
func loadedNamespaces() map[string]bool {
	nsMu.RLock()
	defer nsMu.RUnlock()
	out := make(map[string]bool, len(nsRegistry))
	for name := range nsRegistry {
		if !nsNeedsLoad[name] {
			out[name] = true
		}
	}
	return out
}

func runUnitChunkValue(unit *bytecode.ExecUnit, c *vm.CodeChunk) (v vm.Value, err error) {
	err = vm.WithTypeScope(unit.Types, func() error {
		f := vm.NewFrame(c, nil)
		v, err = f.RunProtected()
		vm.ReleaseFrame(f)
		return err
	})
	return v, err
}

func installLoadBundle(ns *vm.Namespace) {
	v := ns.Def("load-bundle", mustWrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("load-bundle expects 1 arg, got %d", len(vs))
		}
		path, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("load-bundle expects a path string, got %s", vs[0].Type().Name())
		}
		return LoadBundle(string(path))
	}))
	v.SetMeta(vm.NewPersistentMap([]vm.Value{
		vm.Keyword("arglists"), vm.NewList([]vm.Value{vm.NewArrayVector([]vm.Value{vm.Symbol("path")})}),
		vm.Keyword("doc"), vm.String("Loads the compiled bundle (.lgb) at path, as running it would, and returns its main chunk's value. A namespace the bundle marks self-contained stays undecoded until required; a namespace already loaded stays as it is."),
	}))
}
