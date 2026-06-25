/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package compiler

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/nooga/let-go/pkg/bytecode"
	"github.com/nooga/let-go/pkg/errors"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

var consts *vm.Consts

// CoreConsts returns the global const pool populated during core boot.
// Used as parent for layered child pools during user code compilation.
func CoreConsts() *vm.Consts {
	return consts
}

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

// ReadString parses a string into a let-go Value. As a single-form entry point
// it skips leading no-value forms (comments, #_ discard) so a string that opens
// with a ';;' comment yields the following form rather than the VOID sentinel.
func ReadString(s string) (vm.Value, error) {
	reader := NewLispReader(strings.NewReader(s), "<read-string>")
	return reader.ReadSkipNoValue()
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
	c := NewCompiler(consts, rt.NS(rt.NameCoreNS))
	c.SetSource("<embedded:core>")
	_, _, err := c.CompileMultiple(strings.NewReader(rt.CoreSrc))
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
			// Use DefStub to avoid spurious warn-on-shadow warnings during
			// bundle decoding (the chunk will properly Def the value later).
			return n.DefStub(name)
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

	// Source-only namespaces (precompiled bundle deliberately skips them
	// because their precompiled stubs would intern nil into dependent
	// namespaces — see cmd/lgbgen/main.go). Their vars may already exist
	// as stubs in the registry from bundle decoding's VarRef pass; mark
	// them NeedsLoad so (require 'name) re-loads from source.
	for _, name := range []string{"ir.data"} {
		rt.MarkNSNeedsLoad(name)
	}

	return nil
}

func postCoreInit() {
	// read-string: parse a single form from a string. Errors loudly on
	// wrong arity, non-string args, or parse failures.
	readStringFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("read-string: wrong number of arguments %d (expected 1)", len(vs))
		}
		s, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("read-string: expected String, got %T", vs[0])
		}
		return ReadString(string(s))
	})
	coreNS := rt.NS(rt.NameCoreNS)
	rsVar := coreNS.LookupOrAdd(vm.Symbol("read-string"))
	rsVar.(*vm.Var).SetRoot(readStringFn)

	// read-all-string: parse every top-level form from a string,
	// return as a vector. Useful for scripts that walk source
	// form-by-form (dependency analysis, codegen). EOF at a form
	// boundary stops cleanly; EOF mid-form or any other reader
	// error is propagated so callers see syntax errors instead of
	// silent truncation.
	readAllStringFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("read-all-string: wrong number of arguments %d (expected 1)", len(vs))
		}
		s, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("read-all-string: expected String, got %T", vs[0])
		}
		reader := NewLispReader(strings.NewReader(string(s)), "<read-all-string>")
		forms := []vm.Value{}
		for {
			// Peek: skip whitespace, then either give up cleanly (EOF
			// at form boundary) or put the char back so Read can see
			// the start of the next form. Distinguishes clean EOF
			// from mid-form EOF — both surface as IsCausedBy(io.EOF)
			// but only the former is acceptable.
			_, err := reader.eatWhitespace()
			if err != nil {
				if errors.IsCausedBy(err, io.EOF) {
					break
				}
				return vm.NIL, err
			}
			if err := reader.unread(); err != nil {
				return vm.NIL, err
			}
			form, err := reader.Read()
			if err != nil {
				return vm.NIL, err
			}
			forms = append(forms, form)
		}
		return vm.NewPersistentVector(forms), nil
	})
	rasVar := coreNS.LookupOrAdd(vm.Symbol("read-all-string"))
	rasVar.(*vm.Var).SetRoot(readAllStringFn)

	// load-string: compile and evaluate a string of code, returning the last value.
	loadStringFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, nil
		}
		s, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, nil
		}
		c := NewCompiler(consts, rt.NS(rt.NameCoreNS))
		_, out, err := c.CompileMultiple(strings.NewReader(string(s)))
		if err != nil {
			return vm.NIL, err
		}
		return out, nil
	})
	lsVar := coreNS.LookupOrAdd(vm.Symbol("load-string"))
	lsVar.(*vm.Var).SetRoot(loadStringFn)

	// eval: compile and evaluate a single already-read form in the current namespace.
	evalFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, nil
		}
		ns := rt.CurrentNS.Deref().(*vm.Namespace)
		c := NewCompiler(consts, ns)
		c.source = "<eval>"
		c.chunk = vm.NewCodeChunk(c.consts)
		c.resetSP()
		if err := c.compileForm(vs[0]); err != nil {
			return vm.NIL, err
		}
		c.chunk.SetMaxStack(c.spMax)
		c.emit(vm.OP_RETURN)
		f := vm.NewFrame(c.chunk, nil)
		out, err := f.RunProtected()
		vm.ReleaseFrame(f)
		if err != nil {
			return vm.NIL, err
		}
		return out, nil
	})
	evalVar := coreNS.LookupOrAdd(vm.Symbol("eval"))
	evalVar.(*vm.Var).SetRoot(evalFn)

	// set-read-clj!: opt in to matching :clj in reader conditionals.
	// Used by the real-world compat runner; off by default so the
	// conformance suite doesn't reach JVM-only :clj branches.
	setReadCljFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		v := vs[0]
		SetMatchCljConditional(v != vm.NIL && v != vm.FALSE)
		return vm.NIL, nil
	})
	coreNS.LookupOrAdd(vm.Symbol("set-read-clj!")).(*vm.Var).SetRoot(setReadCljFn)

	// set-read-bb!: opt in to matching :bb (babashka) in reader conditionals.
	// Enable alongside set-read-clj! to read babashka-compatible libraries,
	// which ship :bb fallbacks that avoid JVM-internal constructors.
	setReadBbFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		v := vs[0]
		SetMatchBbConditional(v != vm.NIL && v != vm.FALSE)
		return vm.NIL, nil
	})
	coreNS.LookupOrAdd(vm.Symbol("set-read-bb!")).(*vm.Var).SetRoot(setReadBbFn)

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

	// gogen_ir: the core bundle has now replayed every clojure.core
	// def/defn onto coreNS. Drain any Go-native overrides registered by
	// blank-imported lowered packages (lg_gogen_ir.go), clobbering the
	// bytecode-produced vars with NativeFn wrappers. No-op on untagged
	// builds — pendingGoOverrides is empty, so this is one map lookup.
	rt.ApplyGoOverrides(coreNS)
	rt.ApplyGoVarInits(coreNS, vm.RootExecContext)
}
