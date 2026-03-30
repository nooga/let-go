package resolver

import (
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
}

func NewNSResolver(ctx *compiler.Context, path []string) *NSResolver {
	return &NSResolver{
		ctx:      ctx,
		path:     path,
		cloading: make(map[string]bool),
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
	_, _, err = freshCtx.CompileMultiple(f)
	nns := freshCtx.CurrentNS()
	r.ctx.SetCurrentNS(ons)
	if err != nil {
		return nil
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
	p := path.Join(blocks...) + ".lg"
	// Try embedded namespaces first
	if embedded := r.loadEmbedded(name); embedded != nil {
		return embedded
	}
	for _, dir := range r.path {
		cp := path.Join(dir, p)
		if _, err := os.Stat(cp); err == nil {
			r.cloading[name] = true
			lns := r.loadFile(cp)
			delete(r.cloading, name)
			return lns
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
		r.ctx.SetCurrentNS(ons)
		return nil
	}
	_ = result
	nns := r.ctx.CurrentNS()
	r.ctx.SetCurrentNS(ons)
	return nns
}
