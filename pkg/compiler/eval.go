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
	// Namespaces that exist BEFORE decode are native-backed (installers run
	// at package init). If the bundle also carries a chunk for one of them
	// — a HYBRID namespace like async (native fns + lg-source macros) —
	// lazy loading is broken for it: qualified symbol resolution finds the
	// already-registered ns and its decoded nil stubs without ever
	// triggering the loader. Those chunks must run eagerly (below).
	preexisting := map[string]bool{}
	for name := range rt.AllNSes() {
		preexisting[name] = true
	}

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
			if decodeTagStats {
				bytecode.NoteDecodeVarRefMiss(true)
			}
			return n.DefStub(name)
		}
		if decodeTagStats {
			bytecode.NoteDecodeVarRefHit()
		}
		return v
	}
	t0 := time.Now()
	if decodeTagStats {
		bytecode.ResetDecodeStats()
		bytecode.SetDecodeStatsEnabled(true)
		defer bytecode.SetDecodeStatsEnabled(false)
	}
	// Bytes entry: the embedded core is already a []byte, so decode keeps it
	// resident and defers per-chunk source-map materialization off the hot path.
	unit, err := bytecode.DecodeToExecUnitBytes(rt.CoreCompiledLGB, resolve)
	if err != nil {
		return err
	}
	bootMark("decode-bundle", t0)
	if decodeTagStats {
		fmt.Fprint(os.Stderr, bytecode.SnapshotDecodeStats().Summary())
	}

	// Use the decoded const pool as the global pool
	consts = unit.Consts

	// Execute core's main chunk to replay all def/defn/defmacro definitions
	t1 := time.Now()
	f := vm.NewFrame(unit.MainChunk, nil)
	_, err = f.RunProtected()
	vm.ReleaseFrame(f)
	if err != nil {
		return err
	}
	bootMark("run-core-chunk", t1)

	// Store remaining namespace chunks for on-demand loading by the resolver.
	// Mark non-core namespaces as needing load so LookupOrRegisterNS triggers
	// the loader even though the namespace already exists in the registry.
	if unit.NSChunks != nil {
		precompiledNS = unit.NSChunks
		// The lg baseline namespaces (let-go.core, let-go.types) hold lg-specific
		// extras and are auto-refer'd into every namespace (RegisterNS) but never
		// explicitly required. On-demand loading only fires through
		// LookupOrRegisterNS (qualified refs / require), which refer-resolution
		// bypasses — so their .lg chunks would never run and the .lg-defined vars
		// would stay unbound stubs. Run them eagerly right after core, whose
		// definitions they depend on. (purify-clojure-core ②)
		baselineNS := map[string]bool{}
		for _, name := range rt.LgBaselineNSNames() {
			baselineNS[name] = true
			if ch := precompiledNS[name]; ch != nil {
				lf := vm.NewFrame(ch, nil)
				_, lerr := lf.RunProtected()
				vm.ReleaseFrame(lf)
				if lerr != nil {
					return lerr
				}
			}
		}
		for name := range precompiledNS {
			if name != "core" && !baselineNS[name] {
				rt.MarkNSNeedsLoad(name)
			}
		}
		// Eagerly execute chunks of hybrid namespaces (native + bundled lg
		// source): their vars are reachable via qualified symbols without a
		// (require ...), so lazy loading would leave the bundle-defined
		// vars as nil stubs. The loader isn't configured yet at this point
		// (api.NewContext sets it later), so run the chunk directly, the
		// same way the core main chunk runs above. Note: a hybrid chunk
		// that :requires another lazy bundled ns would load it too early
		// here — fine for today's hybrids (async: macros over core only).
		for name, ch := range precompiledNS {
			if name == "core" || !preexisting[name] || ch == nil {
				continue
			}
			f := vm.NewFrame(ch, nil)
			_, err := f.RunProtected()
			vm.ReleaseFrame(f)
			if err != nil {
				return err
			}
			rt.ClearNSNeedsLoad(name)
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
}
