/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package compiler

import (
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/nooga/let-go/pkg/bytecode"
	"github.com/nooga/let-go/pkg/errors"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

var bootTiming = os.Getenv("LG_BOOT_TIMING") != ""
var decodeTagStats = os.Getenv("LG_DECODE_TAG_STATS") != ""
var lookupStats = os.Getenv("LG_LOOKUP_STATS") != ""

func bootMark(label string, since time.Time) time.Time {
	if bootTiming {
		fmt.Fprintf(os.Stderr, "[boot] %-22s %8.3f ms\n", label, float64(time.Since(since).Microseconds())/1000)
	}
	return time.Now()
}

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

// evalInNSChild is the namespace-aware eval behind rt.SetEvalInNS (pod
// client-side code). Pod evals are transient, so it compiles through a
// per-call transient compiler; named rather than inline so the const-pool
// regression test can probe this exact path.
func evalInNSChild(code string, ns *vm.Namespace) (vm.Value, error) {
	c := NewTransientCompiler(consts, ns)
	_, out, err := c.CompileMultiple(strings.NewReader(code))
	return out, err
}

func Eval(src string) (vm.Value, error) {
	ns := rt.NS(rt.NameCoreNS)
	// A per-eval CHILD pool: constants this eval introduces live exactly as
	// long as the chunks/functions that reference them, instead of rooting
	// the process-global pool forever. Shared constants still dedupe against
	// the global parent; a long-lived host calling Eval in a loop no longer
	// leaks one pool entry per transient constant (e.g. regex literals).
	compiler := NewTransientCompiler(consts, ns)

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
	tStart := time.Now()

	// Bundle decode is a short, allocation-heavy burst (~5MB of transient
	// garbage). Letting the GC fire mid-decode adds latency and jitter to a
	// path that runs on every process start. Pause GC for the duration of
	// boot and restore the prior target afterward; the transient garbage is
	// reclaimed on the first collection once normal allocation resumes.
	prevGC := debug.SetGCPercent(-1)
	defer debug.SetGCPercent(prevGC)

	if lookupStats {
		vm.ResetLookupStats()
		vm.SetLookupStatsEnabled(true)
	}

	consts = vm.NewConsts()

	// Try loading pre-compiled bundle
	if len(rt.CoreCompiledLGB) > 0 {
		if err := loadPrecompiledBundle(); err == nil {
			tPost := time.Now()
			postCoreInit()
			bootMark("post-core-init", tPost)
			bootMark("evalInit-total", tStart)
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
	// Bundle-path parity (purify-clojure-core ②): the lg baseline namespaces
	// (let-go.core, …) are auto-refer'd into every namespace but never explicitly
	// required, so the on-demand loader never runs their bodies. Without this,
	// their .lg-defined lg-isms (str-join, spy, range?, …) stay unbound and
	// unqualified use fails to resolve under -tags bootstrap — even though it
	// works on the bundle path (see loadPrecompiledBundle's eager chunk-run).
	// Compile their embedded source eagerly, right after core. Go-only baselines
	// (let-go.types, whose predicates are Def'd in Go) have no source and are
	// skipped, keeping unqualified imports working without code changes.
	for _, name := range rt.LgBaselineNSNames() {
		src, ok := rt.EmbeddedSource(name)
		if !ok {
			continue
		}
		bc := NewCompiler(consts, rt.NS(name))
		bc.SetSource("<embedded:" + name + ">")
		if _, _, berr := bc.CompileMultiple(strings.NewReader(src)); berr != nil {
			panic(name + " compilation failed: " + berr.Error())
		}
	}
	postCoreInit()
}

func loadPrecompiledBundle() error {
	// The decode + replay of core, the lg baseline namespaces, and the eager
	// hybrid pass all live in the compiler-free spine rt.LoadCoreBundle (#506),
	// shared with rt.LoadCore and rt.BootCore. EagerHybrids is true here for the
	// same reason it is in BootCore: api.NewContext returns straight to user
	// code, so hybrid vars reachable via qualified symbols (which bypass the
	// on-demand loader) must be bound before then. The spine also folds in the
	// *ns* save/restore and NSOrder-deterministic hybrid order that used to be
	// BootCore-only.
	//
	// Decode diagnostics (LG_DECODE_TAG_STATS, #356) still work: the var-ref
	// hit/miss counting now lives in rt.LGBVarResolver (self-gating), so the
	// enable/reset/print wrapper here drives it across the shared decode, and
	// the OnPhase hook restores the decode-bundle / run-core-chunk bootMarks.
	if decodeTagStats {
		bytecode.ResetDecodeStats()
		bytecode.SetDecodeStatsEnabled(true)
		defer bytecode.SetDecodeStatsEnabled(false)
	}
	unit, err := rt.LoadCoreBundle(rt.CoreLoadOptions{
		EagerHybrids: true,
		OnPhase:      func(phase string, since time.Time) { bootMark(phase, since) },
	})
	if err != nil {
		return err
	}
	if decodeTagStats {
		fmt.Fprint(os.Stderr, bytecode.SnapshotDecodeStats().Summary())
	}

	// Compiler-side leftovers the runtime spine deliberately omits: the decoded
	// const pool becomes the global pool for further compilation, and the chunk
	// map is what the resolver+compiler NSLoader replays a namespace from.
	consts = unit.Consts
	precompiledNS = unit.NSChunks

	// Source-only namespaces (the precompiled bundle deliberately skips them
	// because their precompiled stubs would intern nil into dependent
	// namespaces — see cmd/lgbgen/main.go). Their vars may already exist as
	// stubs from bundle decoding's VarRef pass; mark them NeedsLoad so
	// (require 'name) re-loads from source.
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
		// Per-call transient compiler, like Eval: constants the loaded code
		// introduces (e.g. regex literals) die with its chunks instead of
		// rooting the process-global pool on every call.
		c := NewTransientCompiler(consts, rt.NS(rt.NameCoreNS))
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
		// Per-call transient compiler (see Eval): a form passed to eval interns
		// its constants into a pool owned by this one chunk, not the global pool.
		c := NewTransientCompiler(consts, ns)
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

	// Wire up namespace-aware eval for pod client-side code.
	rt.SetEvalInNS(evalInNSChild)

	// test, walk, etc. are demand-loaded via resolver when required

	// gogen_ir: the core bundle has now replayed every clojure.core
	// def/defn onto coreNS. Drain any Go-native overrides registered by
	// blank-imported lowered packages (lg_gogen_ir.go), clobbering the
	// bytecode-produced vars with NativeFn wrappers. No-op on untagged
	// builds — pendingGoOverrides is empty, so this is one map lookup.
	rt.ApplyGoOverrides(coreNS)
}
