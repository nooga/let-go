package resolver

import (
	"fmt"
	"io"
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

// Embedded namespace sources are looked up via rt.EmbeddedSource, which
// derives the path from the dotted ns name (every "." is a path separator;
// hyphens in the leaf segment map to underscores). Adding a new embedded
// namespace requires only dropping a `.lg` file under `pkg/rt/core/` — no
// edits here.

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
// even if empty. Otherwise the fallback is used. The returned path is
// exactly the parsed inputs — the current directory is NOT included
// implicitly; callers that want it must list "." themselves. An empty
// input therefore yields no paths.
func PathsFromInputs(explicit, fallback string, explicitSet bool) []string {
	raw := fallback
	if explicitSet {
		raw = explicit
	}
	return ParseSearchPaths(raw)
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
	return r.loadSource(path, f, true)
}

func (r *NSResolver) loadSource(sourceName string, reader io.Reader, recordChunk bool) *vm.Namespace {
	ons := r.ctx.CurrentNS()
	freshCtx := compiler.NewCompiler(r.ctx.Consts(), ons)
	freshCtx.SetSource(sourceName)
	chunk, _, err := freshCtx.CompileMultiple(reader)
	nns := freshCtx.CurrentNS()
	r.ctx.SetCurrentNS(ons)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to load %s: %s\n", sourceName, err)
		return nil
	}
	if recordChunk && chunk != nil && nns != nil {
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
	// Build candidate paths: try .lg, .cljc, then .clj extensions,
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
		hyphenPath + ".clj",
		underscorePath + ".clj",
	}

	for _, dir := range r.path {
		for _, candidate := range candidates {
			cp := path.Join(dir, candidate)
			if _, err := os.Stat(cp); err == nil {
				r.cloading[name] = true
				lns := r.loadFile(cp)
				delete(r.cloading, name)
				// gogen_ir: drain Go-native overrides (no-op untagged).
				rt.ApplyGoOverrides(lns)
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

	if name == "term" {
		// term is a pure Go namespace, already registered in init()
		return rt.NS("term")
	}
	src, ok := rt.EmbeddedSource(name)
	if !ok || src == "" {
		return nil
	}
	r.cloading[name] = true
	defer delete(r.cloading, name)
	ns := r.loadSource("<embedded:"+name+">", stdstrings.NewReader(src), false)
	// gogen_ir: drain Go-native overrides for this namespace (no-op untagged).
	rt.ApplyGoOverrides(ns)
	return ns
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
	// gogen_ir: drain Go-native overrides registered by the lowered
	// package for this namespace (no-op on untagged builds).
	rt.ApplyGoOverrides(nns)
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
