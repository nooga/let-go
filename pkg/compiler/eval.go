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

// precompiledCore is the decoded core bundle: the namespace chunks the
// resolver replays on (require ...), decoded on that first require when boot
// deferred them.
var precompiledCore *bytecode.ExecUnit

// PrecompiledNSChunk returns the precompiled main chunk for a namespace, or
// nil when the bundle does not define it. Decoding a deferred namespace's
// chunks happens here, so it can fail.
func PrecompiledNSChunk(name string) (*vm.CodeChunk, error) {
	if precompiledCore == nil {
		return nil, nil
	}
	return precompiledCore.NSChunk(name)
}

// evalInNSChild is the namespace-aware eval behind rt.SetEvalInNS (pod
// client-side code). Pod evals are transient, so it compiles through a
// per-call transient compiler; named rather than inline so the const-pool
// regression test can probe this exact path.
func evalInNSChild(code string, ns *vm.Namespace) (vm.Value, error) {
	// The compiler switches *ns* to ns for the duration of the compile;
	// the caller's namespace is restored so that loading a pod from inside
	// a source file leaves that file's *ns* where it was.
	prev := rt.CurrentNS.Root()
	defer rt.CurrentNS.SetRoot(prev)
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
	reader := newDataReaderWithResolvers(strings.NewReader(s), "<read-string>", nil, rootDataReaderResolver)
	return reader.ReadSkipNoValue()
}

// ReadDataString reads the first form of s with data semantics (see
// NewLispDataReader); an input with no form is an error. ReadAllDataString
// reads every top-level form the same way; EOF at a form boundary stops
// cleanly, EOF mid-form is an error. Both read only the built-in tags.
func ReadDataString(s string) (vm.Value, error) {
	return readDataString(s, nil, nil)
}

func ReadAllDataString(s string) ([]vm.Value, error) {
	return readAllDataString(s, nil)
}

// readDataString reads the first form of s, resolving tags through resolver
// (nil: only the built-in tags). When eof is non-nil an input with no form
// reads as *eof; otherwise it is an error, as EOF inside a form always is.
func readDataString(s string, resolver taggedDataReaderResolver, eof *vm.Value) (vm.Value, error) {
	reader := newDataReaderWithResolvers(strings.NewReader(s), "<read-string>", nil, resolver)
	form, found, err := readDataForm(reader)
	if err != nil {
		return vm.NIL, err
	}
	if found {
		return form, nil
	}
	if eof != nil {
		return *eof, nil
	}
	return vm.NIL, NewReaderError(reader, "EOF while reading")
}

func readAllDataString(s string, resolver taggedDataReaderResolver) ([]vm.Value, error) {
	reader := newDataReaderWithResolvers(strings.NewReader(s), "<read-all-string>", nil, resolver)
	forms := []vm.Value{}
	for {
		form, found, err := readDataForm(reader)
		if err != nil {
			return nil, err
		}
		if !found {
			return forms, nil
		}
		forms = append(forms, form)
	}
}

// readDataForm reads the next form, skipping comments and #_ discards. found
// is false when the input ends at a form boundary; EOF inside a form is an
// error. Same boundary handling as read-all-string.
func readDataForm(reader *LispReader) (form vm.Value, found bool, err error) {
	for {
		_, err := reader.eatWhitespace()
		if err != nil {
			if errors.IsCausedBy(err, io.EOF) {
				return vm.NIL, false, nil
			}
			return vm.NIL, false, err
		}
		if err := reader.unread(); err != nil {
			return vm.NIL, false, err
		}
		form, err := reader.Read()
		if err != nil {
			return vm.NIL, false, err
		}
		if form.Type() != vm.VoidType {
			return form, true, nil
		}
	}
}

// fnDataReaderResolver hands every tag the reader meets to (fn tag form), run
// in ec. The fn holds the whole policy, so every tag counts as handled and a
// nil result is an ordinary value. It never consults *data-readers*, so no
// raw entry there (such as #go) applies.
func fnDataReaderResolver(ec *vm.ExecContext, fn vm.Fn) taggedDataReaderResolver {
	return &dataReaderPolicy{
		resolve: func(tag vm.Symbol, form vm.Value) (vm.Value, bool, error) {
			out, err := ec.Invoke(fn, []vm.Value{tag, form})
			return out, true, err
		},
		rawEntry: func(string) (*RawDataReader, error) { return nil, nil },
	}
}

// dataReaderResolverArg reads the optional resolver argument vs[1] of the
// read-data-string natives: absent or nil leaves tag resolution to the
// reader's built-in tags.
func dataReaderResolverArg(name string, ec *vm.ExecContext, vs []vm.Value) (taggedDataReaderResolver, error) {
	if len(vs) < 2 || vs[1] == vm.NIL {
		return nil, nil
	}
	fn, ok := vs[1].(vm.Fn)
	if !ok {
		return nil, fmt.Errorf("%s: resolver must be a function, got %s", name, vs[1].Type().Name())
	}
	return fnDataReaderResolver(ec, fn), nil
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
		// A bootstrap resolve failure ("Can't resolve X") usually means a
		// generated //lg:native primitive registered into the wrong namespace.
		// Attach the non-terminating registration audit so ALL misregistered
		// primitives are reported at once (with where each actually landed),
		// instead of re-running to discover them one panic at a time.
		msg := "core.lg compilation failed: " + err.Error()
		if audit := rt.AuditGeneratedPrimitives(); len(audit) > 0 {
			msg += fmt.Sprintf("\n\ngenerated-primitive audit (%d not resolvable in canonical ns):\n  %s",
				len(audit), strings.Join(audit, "\n  "))
		}
		panic(msg)
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
	// hybrid pass all live in the compiler-free spine rt.LoadCoreBundle,
	// shared with rt.LoadCore and rt.BootCore. EagerHybrids is true here for the
	// same reason it is in BootCore: api.NewContext returns straight to user
	// code, so hybrid vars reachable via qualified symbols (which bypass the
	// on-demand loader) must be bound before then. The spine also owns the *ns*
	// save/restore and the NSOrder-deterministic hybrid order.
	//
	// Decode diagnostics (LG_DECODE_TAG_STATS) still work: the var-ref
	// hit/miss counting lives in rt.LGBVarResolver (self-gating), so the
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
	// const pool becomes the global pool for further compilation, and the unit
	// is what the resolver+compiler NSLoader replays a namespace from.
	consts = unit.Consts
	precompiledCore = unit

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
	readStringFn := vm.NewArityNativeFn("read-string", 1, false, func(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("read-string: wrong number of arguments %d (expected 1)", len(vs))
		}
		s, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("read-string: expected String, got %T", vs[0])
		}
		reader := newDataReaderWithResolvers(strings.NewReader(string(s)), "<read-string>", taggedReadersFromExecContext(ec), execContextDataReaderResolver(ec))
		return reader.ReadSkipNoValue()
	})
	coreNS := rt.NS(rt.NameCoreNS)
	rsVar := coreNS.LookupOrAdd(vm.Symbol("read-string"))
	rsVar.(*vm.Var).SetRoot(readStringFn)

	// read-data-string / read-all-data-string: Clojure data-reading semantics
	// for clojure.edn (metadata attached, real sets, discards splice nothing).
	// (read-data-string s) reads only the built-in tags and rejects an input
	// with no form; (read-data-string s resolver) hands every tag to
	// (resolver tag form) in the caller's context; (read-data-string s
	// resolver eof) reads an input with no form as eof. The option policy
	// (:readers, :default, :eof) lives in clojure.edn/read-string.
	readDataStringFn := vm.NewArityNativeFn("read-data-string", 1, true, func(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
		if len(vs) > 3 {
			return vm.NIL, fmt.Errorf("read-data-string: wrong number of arguments %d (expected 1 to 3)", len(vs))
		}
		s, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("read-data-string: expected String, got %T", vs[0])
		}
		resolver, err := dataReaderResolverArg("read-data-string", ec, vs)
		if err != nil {
			return vm.NIL, err
		}
		var eof *vm.Value
		if len(vs) == 3 {
			eof = &vs[2]
		}
		return readDataString(string(s), resolver, eof)
	})
	coreNS.LookupOrAdd(vm.Symbol("read-data-string")).(*vm.Var).SetRoot(readDataStringFn)
	readAllDataStringFn := vm.NewArityNativeFn("read-all-data-string", 1, true, func(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
		if len(vs) > 2 {
			return vm.NIL, fmt.Errorf("read-all-data-string: wrong number of arguments %d (expected 1 or 2)", len(vs))
		}
		s, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("read-all-data-string: expected String, got %T", vs[0])
		}
		resolver, err := dataReaderResolverArg("read-all-data-string", ec, vs)
		if err != nil {
			return vm.NIL, err
		}
		forms, err := readAllDataString(string(s), resolver)
		if err != nil {
			return vm.NIL, err
		}
		return vm.NewPersistentVector(forms), nil
	})
	coreNS.LookupOrAdd(vm.Symbol("read-all-data-string")).(*vm.Var).SetRoot(readAllDataStringFn)

	// read-all-string: parse every top-level form from a string,
	// return as a vector. Useful for scripts that walk source
	// form-by-form (dependency analysis, codegen). EOF at a form
	// boundary stops cleanly; EOF mid-form or any other reader
	// error is propagated so callers see syntax errors instead of
	// silent truncation.
	readAllStringFn := vm.NewArityNativeFn("read-all-string", 1, false, func(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("read-all-string: wrong number of arguments %d (expected 1)", len(vs))
		}
		s, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("read-all-string: expected String, got %T", vs[0])
		}
		reader := newDataReaderWithResolvers(strings.NewReader(string(s)), "<read-all-string>", taggedReadersFromExecContext(ec), execContextDataReaderResolver(ec))
		return readAllForms(reader)
	})
	rasVar := coreNS.LookupOrAdd(vm.Symbol("read-all-string"))
	rasVar.(*vm.Var).SetRoot(readAllStringFn)

	// read-code-string / read-all-code-string: read source the way the
	// compiler loads it, not as data. The code reader reads #{1} as the call
	// form (hash-set 1) and ^:m [1] as a (with-meta ...) form. Tools that
	// lower source (lg.compiler, the gogen fixture trampoline) read with these
	// so a macro sees the same forms the bytecode loader hands it.
	readCodeStringFn := vm.NewArityNativeFn("read-code-string", 1, false, func(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("read-code-string: wrong number of arguments %d (expected 1)", len(vs))
		}
		s, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("read-code-string: expected String, got %T", vs[0])
		}
		reader := newLispReaderWithResolvers(strings.NewReader(string(s)), "<read-code-string>", taggedReadersFromExecContext(ec), execContextDataReaderResolver(ec))
		return reader.ReadSkipNoValue()
	})
	coreNS.LookupOrAdd(vm.Symbol("read-code-string")).(*vm.Var).SetRoot(readCodeStringFn)
	readAllCodeStringFn := vm.NewArityNativeFn("read-all-code-string", 1, false, func(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("read-all-code-string: wrong number of arguments %d (expected 1)", len(vs))
		}
		s, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("read-all-code-string: expected String, got %T", vs[0])
		}
		reader := newLispReaderWithResolvers(strings.NewReader(string(s)), "<read-all-code-string>", taggedReadersFromExecContext(ec), execContextDataReaderResolver(ec))
		return readAllForms(reader)
	})
	coreNS.LookupOrAdd(vm.Symbol("read-all-code-string")).(*vm.Var).SetRoot(readAllCodeStringFn)

	// load-string: compile and evaluate a string of code, returning the last value.
	loadStringFn := vm.NewArityNativeFn("load-string", 1, false, func(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
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
		c.SetTaggedReaders(taggedReadersFromExecContext(ec))
		c.setDataReaderResolver(execContextDataReaderResolver(ec))
		c.setExecContext(ec)
		_, out, err := c.CompileMultiple(strings.NewReader(string(s)))
		if err != nil {
			return vm.NIL, err
		}
		return out, nil
	})
	lsVar := coreNS.LookupOrAdd(vm.Symbol("load-string"))
	lsVar.(*vm.Var).SetRoot(loadStringFn)

	// eval: compile and evaluate a single already-read form in the current
	// namespace. The compiled form runs in the CALLER's execution context, so
	// dynamic bindings active around eval (with-out-str, binding) and the
	// caller's structured-concurrency scope apply to the evaluated code, as
	// thread bindings do around Clojure's eval.
	evalFn := vm.NewCtxNativeFn("eval", func(ec *vm.ExecContext, vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, nil
		}
		ns := ec.Deref(rt.CurrentNS).(*vm.Namespace)
		// Per-call transient compiler (see Eval): a form passed to eval interns
		// its constants into a pool owned by this one chunk, not the global pool.
		c := NewTransientCompiler(consts, ns)
		c.SetTaggedReaders(taggedReadersFromExecContext(ec))
		c.setDataReaderResolver(execContextDataReaderResolver(ec))
		c.setExecContext(ec)
		c.source = "<eval>"
		c.chunk = vm.NewCodeChunk(c.consts)
		c.resetSP()
		if err := c.compileForm(vs[0]); err != nil {
			return vm.NIL, err
		}
		c.chunk.SetMaxStack(c.spMax)
		c.emit(vm.OP_RETURN)
		f := vm.NewFrameIn(c.chunk, nil, c.evaluationExecContext())
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

	// Wire up namespace-aware eval for pod client-side code.
	rt.SetEvalInNS(evalInNSChild)

	// clojure.repl/raw-source: the one native primitive clojure.repl.lg
	// needs (reader access) to implement source-fn/source in .lg.
	registerReplSupport()

	// test, walk, etc. are demand-loaded via resolver when required

	// gogen_ir: the core bundle has now replayed every clojure.core
	// def/defn onto coreNS. Drain any Go-native overrides registered by
	// blank-imported lowered packages (lg_gogen_ir.go), clobbering the
	// bytecode-produced vars with NativeFn wrappers. No-op on untagged
	// builds — pendingGoOverrides is empty, so this is one map lookup.
	rt.ApplyGoOverrides(coreNS)
}

// readAllForms reads every top-level form as a vector. EOF at a form boundary
// stops cleanly; EOF mid-form or any other reader error is returned, so callers
// see syntax errors instead of silent truncation.
func readAllForms(reader *LispReader) (vm.Value, error) {
	forms := []vm.Value{}
	for {
		// Skip whitespace, then either stop (EOF at a form boundary) or put
		// the char back so Read sees the start of the next form. Both EOF
		// cases surface as IsCausedBy(io.EOF); only this one is acceptable.
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
}
