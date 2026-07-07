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
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"runtime/pprof"
	"sort"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/nooga/let-go/pkg/bytecode"
	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/genmanifest"
	"github.com/nooga/let-go/pkg/resolver"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"

	_ "github.com/nooga/let-go/pkg/rt/corefns"
)

type embeddedNamespace struct {
	name string
	src  string
}

// goLoweringBuildAfterPasses is the one ordering constraint that namespace
// :require clauses do not express: ir.build must be loaded after every
// optimization/inference pass. The passes register type/pass state at load
// time; if ir.build loads first that state isn't in effect and lowering emits
// untyped arithmetic (e.g. (dec (count xs)) -> vm.Value - 1). ir.build itself
// only :requires ir.data, so we inject the edge during derivation.
const goLoweringBuildNS = "ir.build"

// deriveGoLoweringOrder computes the Phase-1 compile order for the bundle by
// topo-sorting the discoverable namespaces (every (ns ...) that is not
// lgbgen:skip) by their :require edges, plus the one synthetic edge above.
//
// This replaces a hand-maintained list: a new pass namespace is now picked up
// automatically via its :require edges, in correct dependency order, instead of
// silently nil-ing at runtime when someone forgets to register it. The Go
// emission pass (runGoTarget) already topo-sorts the same way; this brings the
// bytecode/VM-state Phase 1 in line with it.
// reaches reports whether `from` can reach `target` by following :require
// edges (target included as a direct or transitive dependency). Used to skip
// synthetic edges that would close a cycle.
func reaches(from, target string, requires map[string][]string) bool {
	seen := map[string]bool{}
	var walk func(n string) bool
	walk = func(n string) bool {
		if n == target {
			return true
		}
		if seen[n] {
			return false
		}
		seen[n] = true
		for _, d := range requires[n] {
			if walk(d) {
				return true
			}
		}
		return false
	}
	for _, d := range requires[from] {
		if walk(d) {
			return true
		}
	}
	return false
}

func deriveGoLoweringOrder() ([]embeddedNamespace, error) {
	names := rt.EmbeddedNSNames()
	allowed := map[string]string{}
	requires := map[string][]string{}
	for _, n := range names {
		src, ok := rt.EmbeddedSource(n)
		if !ok || hasLgbgenSkipDirective(src) {
			continue
		}
		allowed[n] = src
		requires[n] = parseRequires(src)
	}
	// Synthetic edge: ir.build after the analysis/transform passes (see
	// goLoweringBuildNS above). Only the passes that do NOT themselves depend on
	// ir.build qualify — pipeline and other orchestrators :require ir.build and
	// must stay after it; adding them here would create a cycle. Collect into a
	// SORTED slice before appending: ir.build sorts before ir.passes.*, so
	// topoSort reaches the passes THROUGH this dep list, and appending them in
	// map-iteration (randomized) order would make the emitted order — and thus
	// the bundle bytes — nondeterministic across runs.
	if _, ok := allowed[goLoweringBuildNS]; ok {
		var passEdges []string
		for n := range allowed {
			if strings.HasPrefix(n, "ir.passes.") &&
				!reaches(n, goLoweringBuildNS, requires) {
				passEdges = append(passEdges, n)
			}
		}
		sort.Strings(passEdges)
		requires[goLoweringBuildNS] = append(requires[goLoweringBuildNS], passEdges...)
	}
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

// discoverEmbeddedNS walks pkg/rt's coreFS to find every namespace
// declared by a `(ns ...)` form. Returns (name, source) pairs in
// dependency order derived from each file's `:require` clause.
//
// Namespaces whose source starts with `;; lgbgen:skip` are omitted
// from the bundle — they have runtime intern blocks that need live
// function values and must be loaded from source on demand. `core`
// is always emitted first since the language depends on it.
//
// The `ir.*` AOT/lowering pipeline (build, lower, lower-go, passes/*,
// etc.) is ALSO omitted from the bytecode bundle. It is ~10k lines /
// ~480K of source — more than the entire rest of core — yet it is a
// build-time/opt-in tool: a plain `letgo` script never touches it, and
// the runtime optimizer (`*ir-compile*`, see core.lg) loads it lazily
// via `(require 'ir.passes.pipeline)`. Decoding it eagerly on every
// process start cost ~6ms (the bundle had grown 109K→650K). Excluding
// it here keeps cold startup fast; when the pipeline IS needed it loads
// from embedded source on demand. This skip is intentionally scoped to
// the BUNDLE only — the Go-lowering path (deriveGoLoweringOrder) must
// still see ir.* so `core_go_lowered/ir/**` is generated for the
// `-tags gogen_ir` build, so it does NOT consult isIRBundleSkipped.
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

// isIRBundleSkipped reports whether namespace n belongs to the ir.* AOT
// pipeline and should be excluded from the BYTECODE bundle (but not the
// Go-lowering output). Matches the `ir` root and any `ir.` descendant.
func isIRBundleSkipped(n string) bool {
	return n == "ir" || strings.HasPrefix(n, "ir.")
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
	// lgbgen is a short-lived batch tool that allocates heavily (the interpreted
	// lowering churns ~14GB of transient persistent structures), so the default
	// GOGC=100 spends ~40-50% of the run in GC for a small live set. Raising it
	// to 800 cut a --target=go run from min 31s to 27s (~13%, min-of-3) with
	// negligible peak-memory cost (~hundreds of MB). Honor an explicit GOGC env
	// so callers can tune or disable it; GC percentage never affects output.
	if os.Getenv("GOGC") == "" {
		debug.SetGCPercent(800)
	}

	// Parse mode: --target=go <out-dir>, --target=both [bundle-path], or
	// default bytecode mode. `both` emits the .lgb bundle AND the gogen_ir
	// Go tree from a single core compile (see the dispatch below). Optional
	// -cpuprofile writes a Go CPU profile for the lgbgen process itself.
	targetGo := false
	targetBoth := false
	outPath := "pkg/rt/core_compiled.lgb"
	goOutDir := "pkg/rt/core_go_lowered"
	var target string
	var cpuProfilePath string
	var memProfilePath string
	// codeDir is the base directory for the generated //go:build gogen_ir wireup
	// files (lg_gogen_ir.go, cmd/lgbgen/main_gogen_ir.go, pkg/ir/...). It defaults
	// to "" = the repo root, the canonical generation target. Callers that lower
	// into a throwaway directory (e.g. the determinism test) set it so lgbgen
	// never rewrites the real checkout's wireup. When set, the manifest refresh
	// is skipped too, keeping the run side-effect-free outside codeDir/outDir.
	var codeDir string

	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	fs.StringVar(&target, "target", "", "generation target: go, both, or empty for bundle-only")
	fs.StringVar(&cpuProfilePath, "cpuprofile", "", "write Go CPU profile for the lgbgen process")
	fs.StringVar(&memProfilePath, "memprofile", "", "write Go allocation profile (allocs) for the lgbgen process")
	fs.StringVar(&codeDir, "code-dir", "", "base dir for generated gogen_ir wireup files (default: repo root)")
	fs.Parse(os.Args[1:])

	switch target {
	case "":
	case "go":
		targetGo = true
	case "both":
		targetBoth = true
	default:
		fmt.Fprintf(os.Stderr, "unknown --target=%q (expected go or both)\n", target)
		os.Exit(2)
	}

	args := fs.Args()
	for i := 0; i < len(args); i++ {
		if targetGo && i == 0 {
			goOutDir = args[i]
		} else {
			outPath = args[i]
		}
	}
	if cpuProfilePath != "" {
		f, err := os.Create(cpuProfilePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "create cpuprofile %s: %v\n", cpuProfilePath, err)
			os.Exit(1)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			f.Close()
			fmt.Fprintf(os.Stderr, "start cpuprofile %s: %v\n", cpuProfilePath, err)
			os.Exit(1)
		}
		var stopProfileOnce sync.Once
		stopProfile := func() {
			stopProfileOnce.Do(func() {
				pprof.StopCPUProfile()
				if err := f.Close(); err != nil {
					fmt.Fprintf(os.Stderr, "close cpuprofile %s: %v\n", cpuProfilePath, err)
				}
			})
		}
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		defer signal.Stop(sigCh)
		go func() {
			sig := <-sigCh
			fmt.Fprintf(os.Stderr, "received %s, flushing cpu profile\n", sig)
			stopProfile()
			time.Sleep(250 * time.Millisecond)
			os.Exit(130)
		}()
		defer stopProfile()
	}
	if memProfilePath != "" {
		defer func() {
			f, err := os.Create(memProfilePath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "create memprofile %s: %v\n", memProfilePath, err)
				return
			}
			defer f.Close()
			runtime.GC() // materialize up-to-date allocation stats
			if err := pprof.Lookup("allocs").WriteTo(f, 0); err != nil {
				fmt.Fprintf(os.Stderr, "write memprofile %s: %v\n", memProfilePath, err)
			}
		}()
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
	// compile time on an lgbgen:skip namespace NOT in the derived order can load
	// that source on demand instead of failing to resolve. deriveGoLoweringOrder
	// includes every discoverable (non-lgbgen:skip) namespace, so the loader
	// only covers the genuinely full-skip namespaces (e.g. ir.data, which has a
	// runtime re-intern block) that are referenced but not bundled.
	{
		coreNS := rt.NS(rt.NameCoreNS)
		loaderCtx := compiler.NewCompiler(consts, coreNS)
		rt.SetNSLoader(resolver.NewNSResolver(loaderCtx, nil))
	}

	// Phase 1: compile all namespaces from source (bytecode, sets up VM state).
	embeddedNS, err := deriveGoLoweringOrder()
	if err != nil {
		fmt.Fprintf(os.Stderr, "go-lowering ns derivation failed: %v\n", err)
		os.Exit(1)
	}
	nsChunks := make(map[string]*vm.CodeChunk)
	bundleOrder := make([]string, 0, len(embeddedNS)) // non-ir: encoded into the .lgb

	compileNS := func(ns embeddedNamespace) {
		coreNS := rt.NS(rt.NameCoreNS)
		c := compiler.NewCompiler(consts, coreNS)
		c.SetSource("<embedded:" + ns.name + ">")

		chunk, _, err := c.CompileMultiple(strings.NewReader(ns.src))
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s compilation failed: %v\n", ns.name, err)
			os.Exit(1)
		}
		nsChunks[ns.name] = chunk
		if targetGo || targetBoth {
			fmt.Printf("  compiled %-20s (%d bytecode + lowering to Go)\n", ns.name, len(chunk.Code())*4)
		} else {
			fmt.Printf("  compiled %-10s (%d bytes bytecode)\n", ns.name, len(chunk.Code())*4)
		}
	}

	// Phase 1a: compile every NON-ir namespace into the shared const pool, in
	// dependency order. The ir.* AOT/lowering pipeline is deliberately kept OUT
	// of the bytecode bundle (it loads from source on demand at runtime — see
	// isIRBundleSkipped and core.lg's *ir-compile*). It must be excluded HERE,
	// before the bundle is encoded, because EncodeBundleOrdered emits every
	// const in the pool (encoder.go) and each func const drags in its chunk —
	// so filtering only the ns-order would NOT keep ir bytecode out of the
	// .lgb; ir must never enter the pool before writeBundle.
	var irNS []embeddedNamespace
	for _, ns := range embeddedNS {
		if isIRBundleSkipped(ns.name) {
			irNS = append(irNS, ns)
			continue
		}
		compileNS(ns)
		bundleOrder = append(bundleOrder, ns.name)
	}

	// compileIRForLowering compiles the ir.* namespaces into VM state so the
	// Go-lowering pass can resolve them. Kept in the order derived above, which
	// carries the synthetic ir.build-after-passes edge (deriveGoLoweringOrder)
	// — losing it makes lowering emit untyped arithmetic. Always runs AFTER the
	// bundle is written so ir consts never reach the .lgb pool.
	compileIRForLowering := func() {
		for _, ns := range irNS {
			compileNS(ns)
		}
	}

	// Emit artifacts. `both` writes the (ir-free) bundle first, then compiles
	// ir and lowers the whole program to Go.
	if targetGo {
		compileIRForLowering()
		runGoTarget(goOutDir, codeDir)
		// Skip the manifest refresh when lowering into a throwaway codeDir —
		// the canonical generated.sums must only move under a real regen.
		if codeDir == "" {
			refreshManifest()
		}
		return
	}
	if targetBoth {
		writeBundle(outPath, consts, nsChunks, bundleOrder)
		compileIRForLowering()
		runGoTarget(goOutDir, codeDir)
		return
	}

	// Bytecode mode: write .lgb bundle (ir.* excluded).
	writeBundle(outPath, consts, nsChunks, bundleOrder)
}

// writeBundle encodes the compiled namespace chunks into the .lgb bundle at
// outPath and closes the file before returning (so callers may proceed to the
// Go-lowering target against the same in-memory state).
func writeBundle(outPath string, consts *vm.Consts, nsChunks map[string]*vm.CodeChunk, nsOrder []string) {
	f, err := os.Create(outPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "create %s: %v\n", outPath, err)
		os.Exit(1)
	}

	if err := bytecode.EncodeBundleOrdered(f, consts, nsChunks, nsOrder); err != nil {
		f.Close()
		fmt.Fprintf(os.Stderr, "encode failed: %v\n", err)
		os.Exit(1)
	}

	// Stat is best-effort: success here is just for the byte-count in the
	// success line. If it fails, we still wrote the bundle, so report what we
	// know without dereferencing a nil FileInfo.
	if fi, err := f.Stat(); err == nil {
		fmt.Printf("wrote %s (%d bytes, %d consts, %d namespaces)\n",
			outPath, fi.Size(), len(consts.Values()), len(nsChunks))
	} else {
		fmt.Printf("wrote %s (%d consts, %d namespaces; stat failed: %v)\n",
			outPath, len(consts.Values()), len(nsChunks), err)
	}
	f.Close()
	refreshManifest()
}

// refreshManifest records the content digest of all .lg + lgbgen sources
// into pkg/rt/generated.sums, so the genmanifest staleness test and the
// check-generated CLI can tell whether committed artifacts still match the
// sources on disk.
func refreshManifest() {
	root, err := genmanifest.FindRepoRoot(".")
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: skipping manifest refresh: %v\n", err)
		return
	}

	digest, err := genmanifest.Compute(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: computing manifest digest failed: %v\n", err)
		return
	}
	if err := genmanifest.Write(root, digest); err != nil {
		fmt.Fprintf(os.Stderr, "warning: writing %s failed: %v\n", genmanifest.ManifestRelPath, err)
		return
	}
	fmt.Printf("wrote %s (%s)\n", genmanifest.ManifestRelPath, digest[:12])
}

// nsToGoPkgName converts a let-go namespace name to a valid Go package name:
// the LAST dot-segment, with hyphens turned into underscores
// (ir.passes.pipeline → pipeline, ir.lower-go → lower_go, edn → edn). This is
// both the `package` identifier and the filename stem; the directory nesting
// that mirrors the full namespace is built by nsToGoRelDir. Mirrors Glojure's
// mungeID(getLastNSPart(ns)).
func nsToGoPkgName(nsName string) string {
	parts := strings.Split(nsName, ".")
	last := parts[len(parts)-1]
	return strings.ReplaceAll(last, "-", "_")
}

// nsToGoRelDir converts a let-go namespace name to the relative directory path
// (under the lowered-output root) that mirrors the namespace nesting: dots
// become path separators, hyphens become underscores
// (ir.passes.pipeline → ir/passes/pipeline, ir.lower-go → ir/lower_go).
// Mirrors Glojure's nsToPath.
func nsToGoRelDir(nsName string) string {
	s := strings.ReplaceAll(nsName, "-", "_")
	return filepath.FromSlash(strings.ReplaceAll(s, ".", "/"))
}

// goGeneratedBanner is stamped on the first line of every Go file lgbgen
// writes (both the per-package lowered files and the wireup files). It is the
// signature cleanGoOutputDir uses to recognize — and delete — only files we
// generated.
const goGeneratedBanner = "// Code generated by lgbgen --target=go. DO NOT EDIT."

// cleanGoOutputDir removes lgbgen's own previously-generated files from outDir
// before regeneration, so orphaned packages (from a renamed/removed namespace,
// or another branch's output left on disk) can't survive into this run and
// break the `go build -tags gogen_ir` step.
//
// It deletes ONLY files that carry goGeneratedBanner — files lgbgen itself
// wrote — then prunes any directory left empty. User-authored files have no
// banner and are never touched, so a mistargeted `--target=go <dir>` (e.g.
// pkg/rt) cannot destroy non-generated content. The `.`/filesystem-root/
// repo-root checks remain as a defense-in-depth backstop.
func cleanGoOutputDir(outDir string) {
	clean := filepath.Clean(outDir)
	if clean == "." || clean == string(filepath.Separator) {
		fmt.Fprintf(os.Stderr, "refusing to clean unsafe output dir %q\n", outDir)
		os.Exit(2)
	}
	if root, err := genmanifest.FindRepoRoot("."); err == nil {
		absOut, e1 := filepath.Abs(clean)
		absRoot, e2 := filepath.Abs(root)
		if e1 == nil && e2 == nil && absOut == absRoot {
			fmt.Fprintf(os.Stderr, "refusing to clean repo root %q\n", outDir)
			os.Exit(2)
		}
	}
	info, err := os.Stat(clean)
	if os.IsNotExist(err) {
		return // nothing generated yet
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "stat output dir %s: %v\n", clean, err)
		os.Exit(1)
	}
	if !info.IsDir() {
		fmt.Fprintf(os.Stderr, "output dir %s is not a directory\n", clean)
		os.Exit(1)
	}
	if err := removeGeneratedFiles(clean); err != nil {
		fmt.Fprintf(os.Stderr, "clean output dir %s: %v\n", clean, err)
		os.Exit(1)
	}
	pruneEmptyDirs(clean)
}

// removeGeneratedFiles deletes every *.go file under dir that carries
// goGeneratedBanner. Files without the banner (user-authored) are left intact.
func removeGeneratedFiles(dir string) error {
	return filepath.Walk(dir, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fi.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}
		if !fileHasGeneratedBanner(path) {
			return nil
		}
		return os.Remove(path)
	})
}

// fileHasGeneratedBanner reports whether path carries goGeneratedBanner near
// its top (the banner is the first line, before the package clause).
func fileHasGeneratedBanner(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()
	buf := make([]byte, 256)
	n, _ := f.Read(buf)
	return strings.Contains(string(buf[:n]), goGeneratedBanner)
}

// pruneEmptyDirs removes directories under root (excluding root itself) that are
// empty after generated files were deleted, deepest first.
func pruneEmptyDirs(root string) {
	var dirs []string
	filepath.Walk(root, func(path string, fi os.FileInfo, err error) error {
		if err == nil && fi.IsDir() && path != root {
			dirs = append(dirs, path)
		}
		return nil
	})
	for i := len(dirs) - 1; i >= 0; i-- {
		if entries, err := os.ReadDir(dirs[i]); err == nil && len(entries) == 0 {
			os.Remove(dirs[i])
		}
	}
}

// writeGeneratedFile writes data to path, (re)creating the parent directory and
// retrying on a transient ENOENT. Generating the whole lowered tree is dozens of
// sequential file writes; under a macOS Seatbelt sandbox a freshly-created
// directory can intermittently surface a spurious "no such file or directory" on
// the immediately following write, aborting an otherwise-correct multi-minute
// regen on a different file each run. A bounded retry makes a full regen reliable
// without masking a genuine, persistent failure — a non-transient error returns
// immediately, and a still-ENOENT after the retries are exhausted still errors.
func writeGeneratedFile(path string, data []byte) error {
	var err error
	for attempt := 0; attempt < 5; attempt++ {
		if attempt > 0 {
			time.Sleep(20 * time.Millisecond)
		}
		if err = os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			continue
		}
		if err = os.WriteFile(path, data, 0644); err == nil {
			return nil
		}
		if !os.IsNotExist(err) {
			return err
		}
	}
	return err
}

// runGoTarget re-pipelines each namespace's defn forms through the Go
// lowering pipeline and writes .go source files to outDir. The gogen_ir
// wireup files go under codeDir ("" = repo root, the canonical target).
func runGoTarget(outDir, codeDir string) {
	var failed []string
	var generated []string // Go package names successfully written, for the wireup files.
	// EPIC-012/ITER-0025: lower the whole program in one call so cross-package
	// direct calls resolve (a phase-0 collector builds the global registry of
	// every package's direct-callable exports; phase-1 emits with it visible).
	pipelineVar := rt.LookupVar("ir.passes.pipeline", "lower-all-ns-to-go")
	if pipelineVar == nil {
		fmt.Fprintf(os.Stderr, "fatal: ir.passes.pipeline/lower-all-ns-to-go not found\n")
		os.Exit(1)
	}

	// Hermetic generation: clear the output dir first so orphaned packages
	// from a previous run (a renamed/removed namespace, or another branch's
	// output left on disk) can't survive into this run. Such strays are
	// gitignored — invisible to `jj`/`git status` — yet they break the
	// `go build -tags gogen_ir ./...` step (a stray package with stale or
	// invalid Go fails the build). The dir is fully recreated below.
	//
	// This lives here, not in scripts/generate.lg, because runGoTarget is the
	// single point where both `--target=go` and `--target=both` converge, so
	// every caller (make generate, make lowered, make check-generated) gets
	// the clean.
	cleanGoOutputDir(outDir)
	if err := os.MkdirAll(outDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "create output dir %s: %v\n", outDir, err)
		os.Exit(1)
	}

	embeddedNS, err := discoverEmbeddedNS()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ns discovery failed: %v\n", err)
		os.Exit(1)
	}
	// Gather specs for every lowerable namespace; the whole program is lowered
	// in ONE lower-all-ns-to-go call below so cross-package direct calls resolve.
	type goNsSpec struct {
		nsName    string
		pkgName   string
		pkgDir    string
		relDir    string
		defnCount int
	}
	var goSpecs []goNsSpec
	var specVals []vm.Value // parallel Lisp specs: [pkg-name ns-sym forms import-path]
	const loweredImportPrefix = "github.com/nooga/let-go/pkg/rt/core_go_lowered/"
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

		// Read forms from source and pick out defn/defn- forms (single- AND
		// multi-arity). defmacro forms are skipped for now — their bodies are
		// macro template construction code that doesn't lower cleanly.
		// Multi-arity defns lower per-arity: each arity emits its own
		// `<name>_<arity>` Go func (distinct name — no redeclaration; the ir-fn
		// carries :multi-arity? + :arity so lowering names it from the form),
		// registers per-arity direct-callable [ns name arity] entries, and the
		// whole fn gets one context-aware len(args)-switch var-override adapter.
		// So statically-known calls to multi-arity clojure.core fns
		// (conj/reduce/map/…) go direct instead of trampolining via rt.CachedVarFn.
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
			if isDefnOnly(form) || isMultimethodForm(form) {
				defnForms = append(defnForms, form)
			}
		}

		if len(defnForms) == 0 {
			fmt.Printf("  skip %-20s (no defn forms)\n", ns.name)
			continue
		}

		// Each namespace maps to its own Go package in a directory that
		// mirrors the namespace nesting (ir.passes.pipeline →
		// ir/passes/pipeline/), with the package and filename named after the
		// leaf segment (package pipeline, pipeline.go).
		pkgName := nsToGoPkgName(ns.name)
		relDir := nsToGoRelDir(ns.name)
		pkgDir := filepath.Join(outDir, relDir)
		if err := os.MkdirAll(pkgDir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "create pkg dir %s: %v\n", pkgDir, err)
			os.Exit(1)
		}

		// Accumulate this namespace's lowering spec; emission happens once below.
		importPath := loweredImportPrefix + filepath.ToSlash(relDir)
		specVals = append(specVals, vm.NewPersistentVector([]vm.Value{
			vm.String(pkgName),
			vm.Symbol(ns.name),
			vm.NewPersistentVector(defnForms),
			vm.String(importPath),
		}))
		goSpecs = append(goSpecs, goNsSpec{
			nsName:    ns.name,
			pkgName:   pkgName,
			pkgDir:    pkgDir,
			relDir:    relDir,
			defnCount: len(defnForms),
		})
	}

	// Whole-program lowering in one call (cross-package registry resolved inside).
	result, err := pipelineVar.Invoke([]vm.Value{vm.NewPersistentVector(specVals)})
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal: lower-all-ns-to-go failed: %v\n", err)
		os.Exit(1)
	}
	// mapv returns an indexed vector of per-ns Go source strings, in spec order;
	// a nil/non-string slot means that namespace failed to lower (per-ns
	// resilience preserved — record it as failed, keep going).
	type indexed interface {
		Count() vm.Value
		ValueAt(vm.Value) vm.Value
	}
	results, ok := result.(indexed)
	if !ok {
		fmt.Fprintf(os.Stderr, "fatal: lower-all-ns-to-go returned %T, want indexed vector\n", result)
		os.Exit(1)
	}
	n := 0
	if c, ok := results.Count().(vm.Int); ok {
		n = int(c)
	}
	for i := 0; i < n && i < len(goSpecs); i++ {
		spec := goSpecs[i]
		goSrcVal := results.ValueAt(vm.Int(i))
		goSrc, ok := goSrcVal.(vm.String)
		if !ok {
			fmt.Fprintf(os.Stderr, "%s: Go lowering failed (no source)\n", spec.nsName)
			failed = append(failed, spec.nsName)
			continue
		}
		filename := filepath.Join(spec.pkgDir, spec.pkgName+".go")
		// Stamp the generated banner so cleanGoOutputDir can later recognize
		// this as an lgbgen-owned file and never mistake a user file for one.
		banneredSrc := goGeneratedBanner + "\n\n" + string(goSrc)
		if err := writeGeneratedFile(filename, []byte(banneredSrc)); err != nil {
			fmt.Fprintf(os.Stderr, "%s: write %s: %v\n", spec.nsName, filename, err)
			os.Exit(1)
		}
		fmt.Printf("  wrote %s (%d bytes, %d defns)\n", filename, len(goSrc), spec.defnCount)
		// Wireup blank-imports by directory path (slash-separated), not the leaf
		// package name — sibling namespaces share a leaf name but live at
		// distinct paths.
		generated = append(generated, filepath.ToSlash(spec.relDir))
	}
	// Emit the //go:build gogen_ir blank-import wireup files from exactly
	// the set we just wrote — so namespaces that were skipped or failed to
	// lower are never imported (which would break the tagged build).
	writeGogenWireup(generated, codeDir)
	if len(failed) > 0 {
		fmt.Fprintf(os.Stderr, "lgbgen: %d namespace(s) failed: %v\n", len(failed), failed)
		os.Exit(1)
	}
}

// writeGogenWireup emits the //go:build gogen_ir blank-import files that
// pull every successfully-lowered package into a binary so each package's
// init() runs and registers its Go-native overrides (rt.RegisterGoOverrides).
// The imports must live in a package above pkg/rt — the lowered packages
// import pkg/rt, so importing them from inside pkg/rt would cycle. Three
// files:
//   - lg_gogen_ir.go              → the lg binary (repo root, package main)
//   - cmd/lgbgen/main_gogen_ir.go → the lgbgen self-host (also needs bootstrap)
//   - pkg/ir/zz_gogen_ir_wire_test.go → the pkg/ir TEST binary, so
//     BenchmarkIRCompile (and any -tags gogen_ir test there) dispatches the
//     overridden native passes instead of bytecode. A _test.go file is only
//     compiled into the test build, so it never affects `go build`.
//
// All carry a "Code generated … DO NOT EDIT." banner and are kept in
// lockstep with runGoTarget's output; the import set tracks the generated
// tree exactly, including the exclusion of namespaces that failed to lower.
func writeGogenWireup(pkgNames []string, codeDir string) {
	sorted := append([]string(nil), pkgNames...)
	sort.Strings(sorted)

	const importPrefix = "github.com/nooga/let-go/pkg/rt/core_go_lowered/"
	var importBlock strings.Builder
	for _, p := range sorted {
		fmt.Fprintf(&importBlock, "\t_ \"%s%s\"\n", importPrefix, p)
	}

	render := func(pkgName, buildTag string) string {
		var b strings.Builder
		b.WriteString(goGeneratedBanner + "\n\n")
		b.WriteString("//go:build " + buildTag + "\n\n")
		b.WriteString("package " + pkgName + "\n")
		if importBlock.Len() > 0 {
			b.WriteString("\nimport (\n")
			b.WriteString(importBlock.String())
			b.WriteString(")\n")
		}
		return b.String()
	}

	files := []struct{ path, pkg, buildTag string }{
		{"lg_gogen_ir.go", "main", "gogen_ir"},
		{filepath.Join("cmd", "lgbgen", "main_gogen_ir.go"), "main", "bootstrap && gogen_ir"},
		{filepath.Join("pkg", "ir", "zz_gogen_ir_wire_test.go"), "ir_test", "gogen_ir"},
	}
	for _, f := range files {
		// codeDir defaults to "" (repo root); a caller lowering into a throwaway
		// tree sets it so the wireup never lands on the real checkout. Subdirs
		// (cmd/lgbgen, pkg/ir) are created under codeDir as needed.
		outPath := filepath.Join(codeDir, f.path)
		if dir := filepath.Dir(outPath); dir != "." {
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "mkdir %s: %v\n", dir, err)
				os.Exit(1)
			}
		}
		if err := writeGeneratedFile(outPath, []byte(render(f.pkg, f.buildTag))); err != nil {
			fmt.Fprintf(os.Stderr, "write %s: %v\n", outPath, err)
			os.Exit(1)
		}
		fmt.Printf("  wrote %s (%d pkgs)\n", outPath, len(sorted))
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

// isMultimethodForm returns true for (defmulti ...) / (defmethod ...) forms.
// They are forwarded to the Go-lowering pipeline (alongside defn forms) so
// lower-ns-to-go can capture the dispatcher, lower type-dispatched methods to
// native Go funcs, and devirtualize matching call sites to a `switch v.(type)`
// (with a runtime-vm.MultiFn default arm). The pipeline filters them by
// decl-kind; non-type-dispatched multimethods conservatively stay on the
// runtime path.
func isMultimethodForm(form vm.Value) bool {
	list, ok := form.(vm.Sequable)
	if !ok {
		return false
	}
	seq := list.Seq()
	if seq == nil {
		return false
	}
	sym, ok := seq.First().(vm.Symbol)
	if !ok {
		return false
	}
	return string(sym) == "defmulti" || string(sym) == "defmethod"
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
