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

// declaredNSName returns the namespace a source file declares via its first
// real form — (ns NAME …) or (in-ns 'NAME) — or "" if it can't be determined
// (no ns form, read error). Read-only: it parses just the leading form and
// never evaluates the file, so it has no side effects. Used by loadFromPath to
// confirm a filename-matched candidate actually provides the requested ns.
func declaredNSName(p string) string {
	f, err := os.Open(p)
	if err != nil {
		return ""
	}
	defer f.Close()
	rdr := compiler.NewLispReader(f, p)
	var form vm.Value
	for tries := 0; tries < 64; tries++ {
		v, err := rdr.Read()
		if err != nil {
			return ""
		}
		if v != nil && v.Type() != vm.VoidType {
			form = v
			break
		}
	}
	list, ok := form.(*vm.List)
	if !ok || list.Count() == nil {
		return ""
	}
	headSym, ok := list.First().(vm.Symbol)
	if !ok || (string(headSym) != "ns" && string(headSym) != "in-ns") {
		return ""
	}
	rest := list.Next()
	if rest == nil {
		return ""
	}
	// (ns NAME …) → NAME is a bare symbol. (in-ns 'NAME) → a (quote NAME) list.
	switch nm := rest.First().(type) {
	case vm.Symbol:
		return string(nm)
	case *vm.List:
		if nx := nm.Next(); nx != nil {
			if s, ok := nx.First().(vm.Symbol); ok {
				return string(s)
			}
		}
	}
	return ""
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
	// Restore the requiring namespace on EVERY exit path. CurrentNS is a
	// process-wide var that the compile below repoints (via the scratch ns and
	// the loaded file's own `(ns …)`); if CompileMultiple panics or returns
	// early, an explicit single-site restore would be skipped and leave a stale
	// CurrentNS poisoning later compiles. Deferring makes restoration
	// unconditional.
	defer r.ctx.SetCurrentNS(ons)
	// Compile the dependency with a throwaway namespace as the initial
	// CurrentNS — NOT the requiring namespace (ons). CurrentNS is a process-wide
	// var, and the alias/refer simulations run against whatever it points at; if
	// we seed it with `ons`, any ns-op the dependency emits before its own
	// `(ns …)` switches CurrentNS lands in the REQUIRING namespace's alias/refer
	// table and clobbers it (e.g. a `[graph :as g]` in the dependency overwrites
	// the requiring ns's `[gogen :as g]`). A well-formed namespace file's first
	// form is `(ns …)`, which immediately switches CurrentNS to its own ns, so
	// the scratch ns is only transiently current and never receives real defs.
	// It still needs clojure.core referred so that very first `(ns …)` form can
	// resolve the `ns` macro itself (RegisterNS does the same for real ns's).
	scratch := vm.NewNamespace("<load:" + sourceName + ">")
	if rt.CoreNS != nil {
		scratch.Refer(rt.CoreNS, "", true)
	}
	freshCtx := compiler.NewCompiler(r.ctx.Consts(), scratch)
	freshCtx.SetSource(sourceName)
	chunk, _, err := freshCtx.CompileMultiple(reader)
	nns := freshCtx.CurrentNS()
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
	// Build candidate relative paths: try .lg, .cljc, then .clj extensions,
	// and hyphen vs underscore variants for each path segment.
	hyphenPath := path.Join(blocks...)
	ublocks := make([]string, len(blocks))
	for i, b := range blocks {
		ublocks[i] = stdstrings.ReplaceAll(b, "-", "_")
	}
	underscorePath := path.Join(ublocks...)

	candidates := []string{
		hyphenPath + ".lg",
		underscorePath + ".lg",
		hyphenPath + ".cljc",
		underscorePath + ".cljc",
		hyphenPath + ".clj",
		underscorePath + ".clj",
	}

	// loadFromPath searches the explicit source roots (r.path) — the classpath.
	loadFromPath := func() *vm.Namespace {
		for _, dir := range r.path {
			for _, candidate := range candidates {
				cp := path.Join(dir, candidate)
				if _, err := os.Stat(cp); err != nil {
					continue
				}
				// A file provides namespace `name` only if it actually declares
				// (ns name). The filename-match is necessary but not sufficient:
				// e.g. test/async.lg matches a request for `async` by filename
				// but declares (ns test.async), so loading it would shadow the
				// real (native) async with the wrong namespace. Skip such files
				// so the request falls through to the embedded/native provider.
				// A file with no determinable (ns …) form is loaded as before.
				if d := declaredNSName(cp); d != "" && d != name {
					continue
				}
				r.cloading[name] = true
				lns := r.loadFile(cp)
				delete(r.cloading, name)
				// gogen_ir: drain Go-native overrides (no-op untagged).
				rt.ApplyGoOverrides(lns)
				return lns
			}
		}
		return nil
	}

	// Explicit source roots win over the embedded/bundled copy, so a classpath
	// build (lgbgen --source-paths …) compiles what's actually on disk rather
	// than a stale baked-in snapshot — UNLESS this ns is pinned to the embedded
	// copy via LG_PREFER_EMBEDDED_NS. When r.path is empty (the common case: no
	// --source-paths), loadFromPath is a no-op and embedded serves, so default
	// behavior — and ./lg / test startup cost — is unchanged.
	if preferEmbeddedNS(name) {
		if embedded := r.loadEmbedded(name); embedded != nil {
			return embedded
		}
		return loadFromPath()
	}
	if ns := loadFromPath(); ns != nil {
		return ns
	}
	return r.loadEmbedded(name)
}

// preferEmbeddedNS reports whether namespace `name` is pinned to load from the
// embedded/bundled copy even when an explicit source root also provides it
// (the override for the explicit-sources-win default). Comma-separated list in
// LG_PREFER_EMBEDDED_NS. Mirrors forceSourceNS / LG_FORCE_SOURCE_NS.
func preferEmbeddedNS(name string) bool {
	env := os.Getenv("LG_PREFER_EMBEDDED_NS")
	if env == "" {
		return false
	}
	for _, s := range stdstrings.Split(env, ",") {
		if stdstrings.TrimSpace(s) == name {
			return true
		}
	}
	return false
}

// forceSourceNS reports whether namespace `name` is listed in the
// LG_FORCE_SOURCE_NS env var (comma-separated). Such namespaces skip the
// precompiled bundle chunk and load from the //go:embed raw .lg source, which
// `go build`/`go test` recompiles whenever the .lg file changes — giving a
// fast edit→test loop for a single namespace WITHOUT `make generate`. Dev/test
// only; empty/unset means normal (bundle-first) behavior for every namespace.
func forceSourceNS(name string) bool {
	env := os.Getenv("LG_FORCE_SOURCE_NS")
	if env == "" {
		return false
	}
	for _, s := range stdstrings.Split(env, ",") {
		if stdstrings.TrimSpace(s) == name {
			return true
		}
	}
	return false
}

// loadEmbedded loads bundled namespaces from embedded sources
func (r *NSResolver) loadEmbedded(name string) *vm.Namespace {
	// Try precompiled bytecode first — unless this ns is pinned to source
	// loading via LG_FORCE_SOURCE_NS (dev loop; see forceSourceNS).
	if !forceSourceNS(name) {
		if chunk := compiler.PrecompiledNSChunk(name); chunk != nil {
			return r.execPrecompiled(name, chunk)
		}
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
