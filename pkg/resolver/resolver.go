package resolver

import (
	"fmt"
	"os"
	"path"
	stdstrings "strings"

	"github.com/nooga/let-go/pkg/compiler"
	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

type NSResolver struct {
	ctx      *compiler.Context
	path     []string
	cloading map[string]bool
	// LoadedChunks captures the compiled bytecode for each user-loaded namespace.
	// Used by the compiler to serialize all namespaces into a bundle.
	LoadedChunks map[string]*vm.CodeChunk
	// LoadOrder preserves the order in which namespaces were loaded (dependency order).
	LoadOrder []string
}

// ParseSearchPaths splits a path-list string on os.PathListSeparator,
// dropping empty entries. Returns nil for empty input.
func ParseSearchPaths(raw string) []string {
	if raw == "" {
		return nil
	}
	parts := stdstrings.Split(raw, string(os.PathListSeparator))
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p != "" {
			out = append(out, p)
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// PathsFromInputs returns the namespace search path from explicit and
// fallback inputs. When explicitSet is true the explicit value wins
// even if empty. Otherwise the fallback is used.
func PathsFromInputs(explicit, fallback string, explicitSet bool) []string {
	raw := fallback
	if explicitSet {
		raw = explicit
	}
	return append([]string{"."}, ParseSearchPaths(raw)...)
}

// PathsFromDepsEdn reads deps.edn in dir and returns the :paths entries.
// Returns nil if the file doesn't exist, can't be parsed, or has no :paths.
func PathsFromDepsEdn(dir string) []string {
	depsPath := path.Join(dir, "deps.edn")
	data, err := os.ReadFile(depsPath)
	if err != nil {
		return nil
	}
	val, err := compiler.ReadString(string(data))
	if err != nil {
		return nil
	}
	m, ok := val.(*vm.PersistentMap)
	if !ok {
		return nil
	}
	if !m.Contains(vm.Keyword("paths")) {
		return nil
	}
	pathsVal := m.ValueAt(vm.Keyword("paths"))
	vec, ok := pathsVal.(vm.ArrayVector)
	if !ok {
		return nil
	}
	out := make([]string, 0, len(vec))
	for _, item := range vec {
		if s, ok := item.(vm.String); ok && s != "" {
			out = append(out, string(s))
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func NewNSResolver(ctx *compiler.Context, path []string) *NSResolver {
	return &NSResolver{
		ctx:          ctx,
		path:         path,
		cloading:     make(map[string]bool),
		LoadedChunks: make(map[string]*vm.CodeChunk),
	}
}

// DiscoverDepsEdn reads deps.edn in dir and, if present, appends its
// :paths entries to the resolver's search path. Called by runtime entry
// points after NewNSResolver when the caller wants deps.edn support.
func (r *NSResolver) DiscoverDepsEdn(dir string) {
	if depsPaths := PathsFromDepsEdn(dir); depsPaths != nil {
		r.path = append(r.path, depsPaths...)
	}
}

func (r *NSResolver) loadFile(path string) *vm.Namespace {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()
	ons := r.ctx.CurrentNS()
	freshCtx := compiler.NewCompiler(r.ctx.Consts(), ons)
	freshCtx.SetSource(path)
	chunk, _, err := freshCtx.CompileMultiple(f)
	nns := freshCtx.CurrentNS()
	r.ctx.SetCurrentNS(ons)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to load %s: %s\n", path, err)
		return nil
	}
	if chunk != nil && nns != nil {
		name := nns.Name()
		r.LoadedChunks[name] = chunk
		r.LoadOrder = append(r.LoadOrder, name)
	}
	return nns
}

func (r *NSResolver) SetPath(path []string) {
	r.path = path
}

func (r *NSResolver) Load(name string) *vm.Namespace {
	if r.cloading[name] {
		return nil
	}
	blocks := stdstrings.Split(name, ".")
	// Try embedded namespaces first
	if embedded := r.loadEmbedded(name); embedded != nil {
		return embedded
	}
	// Build candidate paths: try .lg and .cljc extensions,
	// and hyphen vs underscore variants for each path segment.
	hyphenPath := path.Join(blocks...)
	for i, b := range blocks {
		blocks[i] = stdstrings.ReplaceAll(b, "-", "_")
	}
	underscorePath := path.Join(blocks...)

	candidates := []string{
		hyphenPath + ".lg",
		underscorePath + ".lg",
		hyphenPath + ".cljc",
		underscorePath + ".cljc",
	}

	for _, dir := range r.path {
		for _, candidate := range candidates {
			cp := path.Join(dir, candidate)
			if _, err := os.Stat(cp); err == nil {
				r.cloading[name] = true
				lns := r.loadFile(cp)
				delete(r.cloading, name)
				return lns
			}
		}
	}
	return nil
}

// loadEmbedded loads bundled namespaces from embedded sources
func (r *NSResolver) loadEmbedded(name string) *vm.Namespace {
	// Try precompiled bytecode first
	if chunk := compiler.PrecompiledNSChunk(name); chunk != nil {
		return r.execPrecompiled(name, chunk)
	}

	var src string
	switch name {
	case "walk":
		src = rt.WalkSrc
	case "core":
		src = rt.CoreSrc
	case "test":
		src = rt.TestSrc
	case "string":
		src = rt.StringSrc
	case "set":
		src = rt.SetSrc
	case "pprint":
		src = rt.PprintSrc
	case "edn":
		src = rt.EdnSrc
	case "io":
		src = rt.IoSrc
	case "async":
		src = rt.AsyncSrc
	case "zip":
		src = rt.ZipSrc
	case "data":
		src = rt.DataSrc
	case "check":
		src = rt.CheckSrc
	case "term":
		// term is a pure Go namespace, already registered in init()
		return rt.NS("term")
	default:
		return nil
	}
	if src == "" {
		return nil
	}
	r.cloading[name] = true
	defer delete(r.cloading, name)
	// Save and restore CurrentNS — loading changes the global CurrentNS var
	ons := r.ctx.CurrentNS()
	freshCtx := compiler.NewCompiler(r.ctx.Consts(), ons)
	freshCtx.SetSource("<embedded:" + name + ">")
	_, _, err := freshCtx.CompileMultiple(stdstrings.NewReader(src))
	nns := freshCtx.CurrentNS()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to load embedded namespace %s: %s\n", name, err)
		r.ctx.SetCurrentNS(ons)
		return nil
	}
	// Restore the caller's namespace
	r.ctx.SetCurrentNS(ons)
	return nns
}

// execPrecompiled executes a precompiled namespace chunk.
func (r *NSResolver) execPrecompiled(name string, chunk *vm.CodeChunk) *vm.Namespace {
	r.cloading[name] = true
	defer delete(r.cloading, name)

	ons := r.ctx.CurrentNS()
	f := vm.NewFrame(chunk, nil)
	result, err := f.RunProtected()
	vm.ReleaseFrame(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to load precompiled namespace %s: %s\n", name, err)
		r.ctx.SetCurrentNS(ons)
		return nil
	}
	_ = result
	nns := r.ctx.CurrentNS()
	r.ctx.SetCurrentNS(ons)
	return nns
}

func init() {
	// Register the resolver namespace so Lisp code can call
	// resolver/deps-paths to read deps.edn :paths entries.
	ns := vm.NewNamespace("resolver")

	depsPathsFn, _ := vm.NativeFnType.Wrap(func(vs []vm.Value) (vm.Value, error) {
		if len(vs) != 1 {
			return vm.NIL, fmt.Errorf("resolver/deps-paths expects 1 arg")
		}
		dir, ok := vs[0].(vm.String)
		if !ok {
			return vm.NIL, fmt.Errorf("resolver/deps-paths expected String dir")
		}
		paths := PathsFromDepsEdn(string(dir))
		result := make([]vm.Value, len(paths))
		for i, p := range paths {
			result[i] = vm.String(p)
		}
		return vm.NewArrayVector(result), nil
	})
	ns.Def("deps-paths", depsPathsFn)

	rt.RegisterNS(ns)
}
