package resolver

import (
	"os"
	"path"
	"strings"

	"github.com/nooga/let-go/pkg/compiler"
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
	r.ctx.SetSource(path)
	_, _, err = r.ctx.CompileMultiple(f)
	if err != nil {
		return nil
	}
	nns := r.ctx.CurrentNS()
	r.ctx.SetCurrentNS(ons)
	return nns
}

func (r *NSResolver) SetPath(path []string) {
	r.path = path
}

func (r *NSResolver) Load(name string) *vm.Namespace {
	if r.cloading[name] {
		return nil
	}
	blocks := strings.Split(name, ".")
	p := path.Join(blocks...) + ".lg"
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
