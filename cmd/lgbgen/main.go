//go:build bootstrap

// lgbgen compiles core.lg and all embedded namespaces into a pre-compiled .lgb bundle,
// or with --target=go, generates Go source files from the Go-lowering pipeline.
//
// Usage:
//
//	go run -tags bootstrap ./cmd/lgbgen [output-path]              # .lgb bundle (default)
//	go run -tags bootstrap ./cmd/lgbgen --target=go <out-dir>       # Go source bootstrap
//
// Default .lgb output: pkg/rt/core_compiled.lgb (when run from repo root)
// Default Go output dir: pkg/rt/core_go_lowered
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nooga/let-go/pkg/bytecode"
	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/resolver"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

type embeddedNamespace struct {
	name string
	src  string
}

// goLoweringNSOrder is intentionally explicit. The bytecode bundle can use
// discovered dependency order, but Go lowering has pass-state dependencies that
// are not fully represented by namespace :require clauses. In particular,
// ir.build must lower after the type-inference passes have been loaded and
// exercised, or arithmetic like (dec (count xs)) is emitted as vm.Value - 1.
var goLoweringNSOrder = []string{
	"core",
	"walk",
	"string",
	"set",
	"pprint",
	"edn",
	"io",
	"async",
	"hash",
	"test",
	"ir.data.generated",
	"ir.zipper",
	"ir.passes",
	"ir.dominance",
	"ir.lower",
	"ir.lower-go",
	"ir.validate",
	"ir.passes.dce",
	"ir.passes.constfold",
	"ir.passes.mutability",
	"ir.passes.cse",
	"ir.passes.typeinfer",
	"ir.passes.licm",
	"ir.passes.infer-arg-types",
	"graph",
	"ir.build",
	"ir.passes.pipeline",
	"ir.dump",
	"ir.passes.trace",
}

func orderedEmbeddedNS(names []string) ([]embeddedNamespace, error) {
	out := make([]embeddedNamespace, 0, len(names))
	for _, name := range names {
		src, ok := rt.EmbeddedSource(name)
		if !ok {
			return nil, fmt.Errorf("embedded namespace %q not found", name)
		}
		out = append(out, embeddedNamespace{name: name, src: src})
	}
	return out, nil
}

// discoverEmbeddedNS walks pkg/rt's coreFS to find every namespace
// declared by a `(ns ...)` form. Returns (name, source) pairs in
// dependency order derived from each file's `:require` clause.
//
// Namespaces whose source starts with `;; lgbgen:skip` are omitted
// from the bundle — they have runtime intern blocks that need live
// function values and must be loaded from source on demand. `core`
// is always emitted first since the language depends on it.
func discoverEmbeddedNS() ([]embeddedNamespace, error) {
	names := rt.EmbeddedNSNames()
	allowed := map[string]string{}
	requires := map[string][]string{}
	for _, n := range names {
		src, ok := rt.EmbeddedSource(n)
		if !ok {
			continue
		}
		if hasLgbgenSkipDirective(src) {
			continue
		}
		allowed[n] = src
		requires[n] = parseRequires(src)
	}
	// Topological sort. Stable: visit names alphabetically so output
	// order is deterministic across builds. `core` is forced first.
	sorted, err := topoSort(allowed, requires)
	if err != nil {
		return nil, err
	}
	out := make([]embeddedNamespace, 0, len(sorted))
	for _, n := range sorted {
		out = append(out, embeddedNamespace{name: n, src: allowed[n]})
	}
	return out, nil
}

// hasLgbgenSkipDirective reports whether the source begins with a line
// containing `lgbgen:skip` (within the first ~3 lines, allowing the
// licence/blank lines a generator might prepend).
func hasLgbgenSkipDirective(src string) bool {
	for i, line := range strings.SplitN(src, "\n", 4) {
		if i > 2 {
			break
		}
		// `lgbgen:skip-go` is a distinct directive (go-lowering only) and must
		// NOT trigger the bundle skip — exclude it from this substring match.
		if strings.Contains(line, "lgbgen:skip") && !strings.Contains(line, "lgbgen:skip-go") {
			return true
		}
	}
	return false
}

// hasLgbgenSkipGoDirective reports whether the source opts out of Go-lowering
// (the --target=go path) via `;; lgbgen:skip-go`. Such a namespace is still
// compiled into the bytecode bundle; it just isn't lowered to Go. Used for
// namespaces whose forms the lowering pipeline can't emit as valid Go (e.g.
// `core`'s multi-arity `apply`, which would produce redeclared funcs) — they
// run via bytecode under -tags gogen_ir instead.
func hasLgbgenSkipGoDirective(src string) bool {
	for i, line := range strings.SplitN(src, "\n", 4) {
		if i > 2 {
			break
		}
		if strings.Contains(line, "lgbgen:skip-go") {
			return true
		}
	}
	return false
}

// parseRequires extracts the dependency ns names from the first
// `(ns name ... (:require ...))` form, using the actual Lisp reader so
// docstrings, comments, and escape sequences are handled correctly.
// Returns the empty slice when there's no (ns) form or no :require
// clause. Errors are returned as nil — the topo sort tolerates missing
// info (an unknown ns in :require just doesn't constrain order).
func parseRequires(src string) []string {
	r := compiler.NewLispReader(strings.NewReader(src), "<lgbgen:parseRequires>")
	// Read past leading comments/whitespace (LispReader.Read returns
	// vm.Void for line comments). The first real form should be (ns ...).
	var form vm.Value
	for tries := 0; tries < 64; tries++ {
		v, err := r.Read()
		if err != nil {
			return nil
		}
		if v != nil && v.Type() != vm.VoidType {
			form = v
			break
		}
	}
	if form == nil {
		return nil
	}
	list, ok := form.(*vm.List)
	if !ok {
		return nil
	}
	if list.Count() == nil {
		return nil
	}
	// First element: ns symbol
	headSym, ok := list.First().(vm.Symbol)
	if !ok || string(headSym) != "ns" {
		return nil
	}
	var out []string
	// Walk the rest looking for (:require ...) lists.
	for s := list.Next(); s != nil; s = s.Next() {
		clause, ok := s.First().(*vm.List)
		if !ok {
			continue
		}
		head, ok := clause.First().(vm.Keyword)
		if !ok || string(head) != "require" {
			continue
		}
		// Each item after :require is either a bare symbol or a
		// vector whose first element is the ns symbol.
		for cs := clause.Next(); cs != nil; cs = cs.Next() {
			switch item := cs.First().(type) {
			case vm.Symbol:
				out = append(out, string(item))
			case vm.ArrayVector:
				if len(item) > 0 {
					if sym, ok := item[0].(vm.Symbol); ok {
						out = append(out, string(sym))
					}
				}
			case vm.PersistentVector:
				if v := item.ValueAt(vm.Int(0)); v != nil {
					if sym, ok := v.(vm.Symbol); ok {
						out = append(out, string(sym))
					}
				}
			}
		}
	}
	return out
}

// topoSort returns nodes in dependency order. `core` is forced first
// (every other ns implicitly depends on it). Unknown requires are
// dropped — they don't constrain order.
func topoSort(nodes map[string]string, requires map[string][]string) ([]string, error) {
	visited := map[string]int{} // 0=unseen, 1=visiting, 2=done
	var out []string
	var visit func(n string) error
	visit = func(n string) error {
		switch visited[n] {
		case 1:
			return fmt.Errorf("cycle through namespace %q", n)
		case 2:
			return nil
		}
		visited[n] = 1
		// Iterate sorted so order is stable.
		deps := make([]string, 0, len(requires[n]))
		for _, d := range requires[n] {
			if _, ok := nodes[d]; ok {
				deps = append(deps, d)
			}
		}
		for _, d := range deps {
			if err := visit(d); err != nil {
				return err
			}
		}
		visited[n] = 2
		out = append(out, n)
		return nil
	}
	// Visit `core` first if present so it always heads the order.
	if _, ok := nodes["core"]; ok {
		if err := visit("core"); err != nil {
			return nil, err
		}
	}
	// Visit remaining in deterministic (sorted) order.
	keys := make([]string, 0, len(nodes))
	for k := range nodes {
		keys = append(keys, k)
	}
	sortStrings(keys)
	for _, k := range keys {
		if err := visit(k); err != nil {
			return nil, err
		}
	}
	return out, nil
}

func sortStrings(s []string) {
	// avoid importing sort just for this; insertion sort, n is small (≤30)
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j-1] > s[j]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}

func main() {
	// Parse mode: --target=go <out-dir> or default bytecode mode
	targetGo := false
	outPath := "pkg/rt/core_compiled.lgb"
	goOutDir := "pkg/rt/core_go_lowered"

	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--target=go":
			targetGo = true
			if i+1 < len(args) {
				goOutDir = args[i+1]
				i++
			}
		default:
			outPath = args[i]
		}
	}

	// rt.init() has already run — native builtins are registered in CoreNS.
	consts := vm.NewConsts()

	// Phase 0: bootstrap ir.data (same for both modes). ir.data is
	// marked lgbgen:skip so it doesn't end up in the bundle, but its
	// consts still need to be interned in the VM so other NSes that
	// reference IR data can compile.
	{
		src, ok := rt.EmbeddedSource("ir.data")
		if !ok {
			fmt.Fprintln(os.Stderr, "ir.data bootstrap: source not found in embed FS")
			os.Exit(1)
		}
		coreNS := rt.NS(rt.NameCoreNS)
		c := compiler.NewCompiler(consts, coreNS)
		c.SetSource("<embedded:ir.data:lgbgen-bootstrap>")
		if _, _, err := c.CompileMultiple(strings.NewReader(src)); err != nil {
			fmt.Fprintf(os.Stderr, "ir.data bootstrap compilation failed: %v\n", err)
			os.Exit(1)
		}
	}

	// Register an on-demand namespace loader so a namespace that depends at
	// compile time on an lgbgen:skip namespace NOT present in
	// goLoweringNSOrder can load that source on demand instead of failing to
	// resolve. Note that some skip'd namespaces (e.g. graph, a compile-time
	// dependency of the lowered ir.* pipeline via graph/toposort-by) ARE
	// listed in goLoweringNSOrder and therefore are bundled — their skip
	// directive only governs discovery-driven bundling, not the explicit
	// go-lowering order. The loader covers the remaining skip'd namespaces
	// that are referenced but not explicitly listed.
	{
		coreNS := rt.NS(rt.NameCoreNS)
		loaderCtx := compiler.NewCompiler(consts, coreNS)
		rt.SetNSLoader(resolver.NewNSResolver(loaderCtx, nil))
	}

	// Phase 1: compile all namespaces from source (bytecode, sets up VM state).
	embeddedNS, err := orderedEmbeddedNS(goLoweringNSOrder)
	if err != nil {
		fmt.Fprintf(os.Stderr, "go-lowering ns list failed: %v\n", err)
		os.Exit(1)
	}
	nsChunks := make(map[string]*vm.CodeChunk)
	nsOrder := make([]string, 0, len(embeddedNS))

	for _, ns := range embeddedNS {
		coreNS := rt.NS(rt.NameCoreNS)
		c := compiler.NewCompiler(consts, coreNS)
		c.SetSource("<embedded:" + ns.name + ">")

		chunk, _, err := c.CompileMultiple(strings.NewReader(ns.src))
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s compilation failed: %v\n", ns.name, err)
			os.Exit(1)
		}
		nsChunks[ns.name] = chunk
		nsOrder = append(nsOrder, ns.name)
		if targetGo {
			fmt.Printf("  compiled %-20s (%d bytecode + lowering to Go)\n", ns.name, len(chunk.Code())*4)
		} else {
			fmt.Printf("  compiled %-10s (%d bytes bytecode)\n", ns.name, len(chunk.Code())*4)
		}
	}

	if targetGo {
		runGoTarget(goOutDir)
		return
	}

	// Bytecode mode: write .lgb bundle.
	f, err := os.Create(outPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "create %s: %v\n", outPath, err)
		os.Exit(1)
	}
	defer f.Close()

	if err := bytecode.EncodeBundleOrdered(f, consts, nsChunks, nsOrder); err != nil {
		fmt.Fprintf(os.Stderr, "encode failed: %v\n", err)
		os.Exit(1)
	}

	fi, _ := f.Stat()
	fmt.Printf("wrote %s (%d bytes, %d consts, %d namespaces)\n",
		outPath, fi.Size(), len(consts.Values()), len(nsChunks))
}

// nsToGoPkgName converts a let-go namespace name to a valid Go package name.
// Dots become underscores (ir.zipper → ir_zipper), hyphens become underscores.
func nsToGoPkgName(nsName string) string {
	s := strings.ReplaceAll(nsName, ".", "_")
	s = strings.ReplaceAll(s, "-", "_")
	return s
}

// runGoTarget re-pipelines each namespace's defn forms through the Go
// lowering pipeline and writes .go source files to outDir.
func runGoTarget(outDir string) {
	var failed []string
	pipelineVar := rt.LookupVar("ir.passes.pipeline", "lower-ns-to-go")
	if pipelineVar == nil {
		fmt.Fprintf(os.Stderr, "fatal: ir.passes.pipeline/lower-ns-to-go not found\n")
		os.Exit(1)
	}

	if err := os.MkdirAll(outDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "create output dir %s: %v\n", outDir, err)
		os.Exit(1)
	}

	embeddedNS, err := discoverEmbeddedNS()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ns discovery failed: %v\n", err)
		os.Exit(1)
	}
	for _, ns := range embeddedNS {
		src := ns.src
		if src == "" {
			continue
		}
		// Namespaces marked lgbgen:skip-go are bundled (bytecode) but not
		// lowered to Go — their forms don't emit valid Go. They run via
		// bytecode under -tags gogen_ir.
		if hasLgbgenSkipGoDirective(src) {
			fmt.Printf("  skip-go %-20s (lgbgen:skip-go)\n", ns.name)
			continue
		}

		// Read forms from source and pick out single-arity defn/defn- forms.
		// defmacro forms are skipped for now — their bodies are macro
		// template construction code that doesn't lower cleanly. Multi-arity
		// functions are also skipped: the Go target emits one func per arity
		// under the same Go name, which Go rejects as a redeclaration. They
		// fall back to bytecode under -tags gogen_ir.
		r := compiler.NewLispReader(strings.NewReader(src), "<embedded:"+ns.name+">")
		var defnForms []vm.Value
		for {
			form, err := r.Read()
			if err != nil {
				// The reader signals end-of-input with a ReaderError wrapping
				// io.EOF, which errors.Is(io.EOF) does NOT match (ReaderError
				// has no Unwrap). Use compiler.IsErrorEOF, the same check the
				// compiler's own read loop uses, so a clean EOF isn't mistaken
				// for a syntax error.
				if compiler.IsErrorEOF(err) {
					break
				}
				fmt.Fprintf(os.Stderr, "%s: read error: %v\n", ns.name, err)
				os.Exit(1)
			}
			if isDefnOnly(form) && isSingleArityDefn(form) {
				defnForms = append(defnForms, form)
			}
		}

		if len(defnForms) == 0 {
			fmt.Printf("  skip %-20s (no defn forms)\n", ns.name)
			continue
		}

		// Each namespace maps to its own Go package in a subdirectory.
		pkgName := nsToGoPkgName(ns.name)
		pkgDir := filepath.Join(outDir, pkgName)
		if err := os.MkdirAll(pkgDir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "create pkg dir %s: %v\n", pkgDir, err)
			os.Exit(1)
		}

		// Batch all forms into one namespace file.
		formsVec := vm.NewPersistentVector(defnForms)
		args := []vm.Value{
			vm.String(pkgName),
			vm.Symbol(ns.name),
			formsVec,
		}
		result, err := pipelineVar.Invoke(args)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: Go lowering failed: %v\n", ns.name, err)
			failed = append(failed, ns.name)
			continue
		}
		goSrc, ok := result.(vm.String)
		if !ok {
			fmt.Fprintf(os.Stderr, "%s: expected string result, got %T\n", ns.name, result)
			failed = append(failed, ns.name)
			continue
		}
		filename := filepath.Join(pkgDir, pkgName+".go")
		if err := os.WriteFile(filename, []byte(goSrc), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "%s: write %s: %v\n", ns.name, filename, err)
			os.Exit(1)
		}
		fmt.Printf("  wrote %s (%d bytes, %d defns)\n", filename, len(goSrc), len(defnForms))
	}
	if len(failed) > 0 {
		fmt.Fprintf(os.Stderr, "lgbgen: %d namespace(s) failed: %v\n", len(failed), failed)
		os.Exit(1)
	}
}

// isDefnOnly returns true if form is a list beginning with defn or defn-
// (a public or private function definition), but not defmacro. Both defn
// and defn- are candidates for Go lowering: private helpers are called by
// the public functions in the same namespace, so omitting them would
// silently drop their native lowering and fall back to bytecode. Callers
// must additionally gate on isSingleArityDefn — see its note on why
// multi-arity forms cannot lower to a single Go function.
func isDefnOnly(form vm.Value) bool {
	list, ok := form.(vm.Sequable)
	if !ok {
		return false
	}
	seq := list.Seq()
	if seq == nil {
		return false
	}
	head := seq.First()
	sym, ok := head.(vm.Symbol)
	if !ok {
		return false
	}
	return string(sym) == "defn" || string(sym) == "defn-"
}

// isSingleArityDefn reports whether a defn/defn- form has exactly one arity.
// A single-arity form places its argument vector at the top level —
// (defn name docstring? meta? [args] body...) — whereas a multi-arity form
// wraps each arity in its own list — (defn name ([a] ...) ([a b] ...)) — and
// so has no top-level vector. The Go target emits one func per arity under
// the same Go name; for multi-arity that is a redeclaration Go rejects, so
// those forms are skipped and run as bytecode under -tags gogen_ir.
func isSingleArityDefn(form vm.Value) bool {
	list, ok := form.(vm.Sequable)
	if !ok {
		return false
	}
	seq := list.Seq()
	if seq == nil {
		return false
	}
	// Skip the head (defn/defn-) and the name symbol; scan the remaining
	// top-level elements for an argument vector.
	for s := seq.Next(); s != nil; s = s.Next() {
		switch s.First().(type) {
		case vm.ArrayVector, vm.PersistentVector:
			return true
		}
	}
	return false
}
