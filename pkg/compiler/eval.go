/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package compiler

import (
	"bytes"
	"strings"

	"github.com/nooga/let-go/pkg/bytecode"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

var consts *vm.Consts

// precompiledNS holds decoded namespace chunks from the bundle.
var precompiledNS map[string]*vm.CodeChunk

// PrecompiledNSChunk returns the precompiled main chunk for a namespace, or nil.
func PrecompiledNSChunk(name string) *vm.CodeChunk {
	if precompiledNS == nil {
		return nil
	}
	return precompiledNS[name]
}

func Eval(src string) (vm.Value, error) {
	ns := rt.NS(rt.NameCoreNS)
	compiler := NewCompiler(consts, ns)

	_, out, err := compiler.CompileMultiple(strings.NewReader(src))
	if err != nil {
		return vm.NIL, err
	}

	return out, nil
}

// ReadString parses a string into a let-go Value.
func ReadString(s string) (vm.Value, error) {
	reader := NewLispReader(strings.NewReader(s), "<read-string>")
	return reader.Read()
}

func evalInit() {
	consts = vm.NewConsts()

	// Try loading pre-compiled bundle
	if len(rt.CoreCompiledLGB) > 0 {
		if err := loadPrecompiledBundle(); err == nil {
			postCoreInit()
			return
		}
		// Fall through to source compilation on error
	}

	// Original path: compile from source
	_, err := Eval(rt.CoreSrc)
	if err != nil {
		panic("core.lg compilation failed: " + err.Error())
	}
	postCoreInit()
}

func loadPrecompiledBundle() error {
	resolve := func(nsName, name string) *vm.Var {
		// Use DefNSBare to create minimal namespaces without triggering
		// the loader. This ensures vars have a home namespace but the
		// actual loading (executing precompiled chunks) happens on demand.
		n := rt.DefNSBare(nsName)
		// Use LookupLocal to check only the namespace's own registry,
		// not refers. This matches how the compiler creates vars via
		// LookupOrAdd (which also skips refers).
		v := n.LookupLocal(vm.Symbol(name))
		if v == nil {
			return n.Def(name, vm.NIL)
		}
		return v
	}
	unit, err := bytecode.DecodeToExecUnit(bytes.NewReader(rt.CoreCompiledLGB), resolve)
	if err != nil {
		return err
	}

	// Use the decoded const pool as the global pool
	consts = unit.Consts

	// Execute core's main chunk to replay all def/defn/defmacro definitions
	f := vm.NewFrame(unit.MainChunk, nil)
	_, err = f.RunProtected()
	vm.ReleaseFrame(f)
	if err != nil {
		return err
	}

	// Store remaining namespace chunks for on-demand loading by the resolver.
	// Mark non-core namespaces as needing load so LookupOrRegisterNS triggers
	// the loader even though the namespace already exists in the registry.
	if unit.NSChunks != nil {
		precompiledNS = unit.NSChunks
		for name := range precompiledNS {
			if name != "core" {
				rt.MarkNSNeedsLoad(name)
			}
		}
	}

	return nil
}

func postCoreInit() {
	// Register read-string (needs the reader which lives in the compiler package)
	readStringFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, nil
		}
		s, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, nil
		}
		return ReadString(string(s))
	})
	coreNS := rt.NS(rt.NameCoreNS)
	rsVar := coreNS.LookupOrAdd(vm.Symbol("read-string"))
	rsVar.(*vm.Var).SetRoot(readStringFn)

	// Wire up EDN reader for pod support
	rt.SetReadEDN(func(s string) (vm.Value, error) {
		return ReadString(s)
	})

	// Wire up namespace-aware eval for pod client-side code
	rt.SetEvalInNS(func(code string, ns *vm.Namespace) (vm.Value, error) {
		c := NewCompiler(consts, ns)
		_, out, err := c.CompileMultiple(strings.NewReader(code))
		return out, err
	})

	// test, walk, etc. are demand-loaded via resolver when required
}
